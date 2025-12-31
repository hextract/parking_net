package client

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/Nerzal/gocloak/v13"
	"github.com/h4x4d/parking_net/pkg/models"
)

func (c Client) CheckToken(ctx context.Context, token string) (user *models.User, err error) {
	if token == "" {
		return nil, errors.New("token is required")
	}

	slog.Info("CheckToken: trying GetUserInfo", slog.String("realm", c.Config.Realm))
	usrInfo, err := c.Client.GetUserInfo(ctx, token, c.Config.Realm)
	var userId string
	var email string

	if err != nil {
		slog.Warn("CheckToken: GetUserInfo failed, trying to decode JWT", slog.String("error", err.Error()))
		// If GetUserInfo fails, decode JWT token to get user ID
		parts := strings.Split(token, ".")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid token format")
		}

		// Decode the payload (second part)
		payload, err := base64.RawURLEncoding.DecodeString(parts[1])
		if err != nil {
			return nil, fmt.Errorf("failed to decode token payload: %w", err)
		}

		var claims map[string]interface{}
		if err := json.Unmarshal(payload, &claims); err != nil {
			return nil, fmt.Errorf("failed to parse token claims: %w", err)
		}

		// Extract user ID from 'sub' claim
		sub, ok := claims["sub"].(string)
		if !ok || sub == "" {
			return nil, errors.New("token missing 'sub' claim")
		}
		userId = sub

		// Extract email if available
		if emailVal, ok := claims["email"].(string); ok {
			email = emailVal
		}

		slog.Info("CheckToken: decoded JWT token", slog.String("user_id", userId), slog.String("email", email))
	} else {
		slog.Info("CheckToken: GetUserInfo succeeded")
		// GetUserInfo worked, use email to find user
		if usrInfo.Email != nil {
			email = *usrInfo.Email
		}
	}

	adminToken, err := c.GetAdminToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get admin token: %w", err)
	}

	var users []*gocloak.User
	if userId != "" {
		// We have user ID from JWT, get user directly
		keycloakUser, err := c.Client.GetUserByID(ctx, adminToken.AccessToken, c.Config.Realm, userId)
		if err != nil {
			return nil, fmt.Errorf("failed to get user by ID: %w", err)
		}
		users = []*gocloak.User{keycloakUser}
	} else if email != "" {
		// Use email to find user
		exact := true
		params := gocloak.GetUsersParams{
			Email: &email,
			Exact: &exact,
		}
		users, err = c.Client.GetUsers(ctx, adminToken.AccessToken, c.Config.Realm, params)
		if err != nil {
			return nil, fmt.Errorf("failed to get users: %w", err)
		}
		if len(users) == 0 {
			return nil, errors.New("user not found")
		}
	} else {
		return nil, errors.New("cannot determine user: no user ID or email")
	}

	if users[0].ID == nil {
		return nil, errors.New("user ID is missing")
	}
	userId = *users[0].ID

	tgId := 0
	if users[0].Attributes != nil {
		telegramIDAttr, exists := (*users[0].Attributes)["telegram_id"]
		if exists && len(telegramIDAttr) > 0 && telegramIDAttr[0] != "" {
			parsedId, err := strconv.Atoi(telegramIDAttr[0])
			if err == nil {
				tgId = parsedId
			}
		}
	}

	groups, err := c.Client.GetUserGroups(ctx, adminToken.AccessToken, c.Config.Realm, userId, gocloak.GetGroupsParams{})
	if err != nil {
		return nil, fmt.Errorf("failed to get user groups: %w", err)
	}

	var role string
	if len(groups) == 0 {
		// User has no groups, assign default "driver" group
		driverGroups, err := c.Client.GetGroups(ctx, adminToken.AccessToken, c.Config.Realm, gocloak.GetGroupsParams{
			Search: gocloak.StringP("driver"),
		})
		if err != nil || len(driverGroups) == 0 {
			return nil, errors.New("failed to find default driver group")
		}
		driverGroupID := *driverGroups[0].ID

		// Check if user is already in the driver group before adding
		userGroups, err := c.Client.GetUserGroups(ctx, adminToken.AccessToken, c.Config.Realm, userId, gocloak.GetGroupsParams{})
		if err == nil {
			for _, group := range userGroups {
				if group.ID != nil && *group.ID == driverGroupID {
					slog.Info("CheckToken: user already in driver group", slog.String("user_id", userId))
					role = "driver"
					return &models.User{
						UserID:     userId,
						TelegramID: tgId,
						Role:       role,
					}, nil
				}
			}
		}

		err = c.Client.AddUserToGroup(ctx, adminToken.AccessToken, c.Config.Realm, userId, driverGroupID)
		if err != nil {
			// Check if user is already in the group (409 Conflict or similar)
			if strings.Contains(err.Error(), "409") || strings.Contains(err.Error(), "Conflict") || strings.Contains(err.Error(), "Duplicate") {
				slog.Info("CheckToken: user already in driver group (detected via error), ignoring", slog.String("user_id", userId))
				role = "driver"
			} else {
				return nil, fmt.Errorf("failed to assign user to driver group: %w", err)
			}
		} else {
			role = "driver"
		}
	} else {
		if groups[0].Name == nil {
			return nil, errors.New("group name is missing")
		}
		role = *groups[0].Name
	}
	return &models.User{
		UserID:     userId,
		TelegramID: tgId,
		Role:       role,
	}, nil
}
