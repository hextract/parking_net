package impl

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/Nerzal/gocloak/v13"
	"github.com/h4x4d/parking_net/auth/internal/restapi/operations"
	"github.com/h4x4d/parking_net/auth/internal/utils"
	"github.com/h4x4d/parking_net/pkg/client"
)

func CreateUser(ctx context.Context, clt *client.Client, fields operations.PostAuthRegisterBody) (*string, error) {
	if fields.Login == nil || fields.Email == nil || fields.Password == nil || fields.Role == nil {
		return nil, fmt.Errorf("required fields are missing")
	}

	// Set default telegram_id to 0 if not provided (optional field)
	telegramID := int64(0)
	if fields.TelegramID != nil {
		telegramID = *fields.TelegramID
	}

	if err := utils.ValidateLogin(*fields.Login); err != nil {
		return nil, fmt.Errorf("invalid login: %v", err)
	}

	if err := utils.ValidateEmail(*fields.Email); err != nil {
		return nil, fmt.Errorf("invalid email: %v", err)
	}

	if err := utils.ValidatePassword(*fields.Password); err != nil {
		return nil, fmt.Errorf("invalid password: %v", err)
	}

	if err := utils.ValidateRole(*fields.Role); err != nil {
		return nil, fmt.Errorf("invalid role: %v", err)
	}

	if err := utils.ValidateTelegramID(telegramID); err != nil {
		return nil, fmt.Errorf("invalid telegram ID: %v", err)
	}

	user := gocloak.User{
		Email:    fields.Email,
		Enabled:  gocloak.BoolP(true),
		Username: fields.Login,
		Attributes: &map[string][]string{
			"telegram_id": {strconv.FormatInt(telegramID, 10)},
		},
		Groups: &[]string{*fields.Role},
	}

	token, err := clt.GetAdminToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get admin token")
	}

	userId, err := clt.Client.CreateUser(ctx, token.AccessToken, clt.Config.Realm, user)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") || strings.Contains(err.Error(), "duplicate") {
			return nil, fmt.Errorf("user already exists")
		}
		return nil, fmt.Errorf("failed to create user")
	}

	groups, err := clt.Client.GetGroups(ctx, token.AccessToken, clt.Config.Realm, gocloak.GetGroupsParams{
		Search: fields.Role,
	})
	if err == nil && len(groups) > 0 {
		groupID := *groups[0].ID
		err = clt.Client.AddUserToGroup(ctx, token.AccessToken, clt.Config.Realm, userId, groupID)
		if err != nil {
			slog.Warn("failed to assign user to group",
				slog.String("user_id", userId),
				slog.String("group", *fields.Role),
				slog.String("error", err.Error()),
			)
		}
	}

	err = clt.Client.SetPassword(ctx, token.AccessToken, userId, clt.Config.Realm, *fields.Password, false)
	if err != nil {
		return nil, fmt.Errorf("failed to set password")
	}

	userToken, err := clt.Client.Login(ctx, clt.Config.Client, clt.Config.ClientSecret, clt.Config.Realm, *fields.Login, *fields.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to login after registration")
	}

	return &userToken.AccessToken, nil
}
