package client

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/Nerzal/gocloak/v13"
	"github.com/h4x4d/parking_net/pkg/models"
)

func (c Client) CheckToken(ctx context.Context, token string) (user *models.User, err error) {
	if token == "" {
		return nil, errors.New("token is required")
	}

	usrInfo, err := c.Client.GetUserInfo(ctx, token, c.Config.Realm)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	exact := true
	params := gocloak.GetUsersParams{
		Email: usrInfo.Email,
		Exact: &exact,
	}
	adminToken, err := c.GetAdminToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get admin token: %w", err)
	}
	users, err := c.Client.GetUsers(ctx, adminToken.AccessToken, c.Config.Realm, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	if len(users) == 0 {
		return nil, errors.New("user not found")
	}
	if users[0].ID == nil {
		return nil, errors.New("user ID is missing")
	}
	userId := *users[0].ID

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
	if len(groups) == 0 {
		return nil, errors.New("user has no groups assigned")
	}
	if groups[0].Name == nil {
		return nil, errors.New("group name is missing")
	}
	role := *groups[0].Name
	return &models.User{
		UserID:     userId,
		TelegramID: tgId,
		Role:       role,
	}, nil
}
