package impl

import (
	"context"
	"fmt"
	"github.com/Nerzal/gocloak/v13"
	"github.com/h4x4d/parking_net/auth/internal/restapi/operations"
	"github.com/h4x4d/parking_net/pkg/client"
)

func ChangePasswordUser(ctx context.Context, clt *client.Client, fields operations.PostChangePasswordBody) (*string, error) {
	if fields.Login == nil || fields.OldPassword == nil || fields.NewPassword == nil {
		return nil, fmt.Errorf("login, old password, and new password are required")
	}

	_, err := clt.Client.Login(ctx, clt.Config.Client, clt.Config.ClientSecret,
		clt.Config.Realm, *fields.Login, *fields.OldPassword)
	if err != nil {
		return nil, fmt.Errorf("invalid old password: %w", err)
	}

	token, err := clt.GetAdminToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get admin token: %w", err)
	}

	params := gocloak.GetUsersParams{
		Username: fields.Login,
	}
	users, err := clt.Client.GetUsers(ctx, token.AccessToken, clt.Config.Realm, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	if users[0].ID == nil {
		return nil, fmt.Errorf("invalid user data")
	}

	err = clt.Client.SetPassword(ctx, token.AccessToken, *users[0].ID,
		clt.Config.Realm, *fields.NewPassword, false)
	if err != nil {
		return nil, fmt.Errorf("failed to set new password: %w", err)
	}

	userToken, err := clt.Client.Login(ctx, clt.Config.Client, clt.Config.ClientSecret,
		clt.Config.Realm, *fields.Login, *fields.NewPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to login with new password: %w", err)
	}

	return &userToken.AccessToken, nil
}
