package impl

import (
	"context"
	"fmt"

	"github.com/h4x4d/parking_net/auth/internal/restapi/operations"
	"github.com/h4x4d/parking_net/auth/internal/utils"
	"github.com/h4x4d/parking_net/pkg/client"
)

func LoginUser(ctx context.Context, clt *client.Client, fields operations.PostAuthLoginBody) (*string, error) {
	if fields.Login == nil || fields.Password == nil {
		return nil, fmt.Errorf("login and password are required")
	}

	if err := utils.ValidateLogin(*fields.Login); err != nil {
		return nil, fmt.Errorf("invalid login")
	}

	if *fields.Password == "" {
		return nil, fmt.Errorf("password cannot be empty")
	}

	userToken, err := clt.Client.Login(ctx, clt.Config.Client, clt.Config.ClientSecret,
		clt.Config.Realm, *fields.Login, *fields.Password)
	if err != nil {
		return nil, fmt.Errorf("authentication failed")
	}
	return &userToken.AccessToken, nil
}
