package impl

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/Nerzal/gocloak/v13"
	"github.com/h4x4d/parking_net/auth/internal/restapi/operations"
	"github.com/h4x4d/parking_net/pkg/client"
)

func CreateUser(ctx context.Context, clt *client.Client, fields operations.PostRegisterBody) (*string, error) {
	if fields.Login == nil || fields.Email == nil || fields.Password == nil ||
		fields.Role == nil || fields.TelegramID == nil {
		return nil, fmt.Errorf("all fields are required")
	}

	user := gocloak.User{
		Email:    fields.Email,
		Enabled:  gocloak.BoolP(true),
		Username: fields.Login,
		Attributes: &map[string][]string{
			"telegram_id": {strconv.FormatInt(*fields.TelegramID, 10)},
		},
		Groups: &[]string{*fields.Role},
	}

	token, err := clt.GetAdminToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get admin token: %w", err)
	}

	userId, err := clt.Client.CreateUser(ctx, token.AccessToken, clt.Config.Realm, user)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") || strings.Contains(err.Error(), "duplicate") {
			return nil, fmt.Errorf("user already exists")
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
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
		return nil, fmt.Errorf("failed to set password: %w", err)
	}

	userToken, err := clt.Client.Login(ctx, clt.Config.Client, clt.Config.ClientSecret, clt.Config.Realm, *fields.Login, *fields.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to login after registration: %w", err)
	}

	return &userToken.AccessToken, nil
}
