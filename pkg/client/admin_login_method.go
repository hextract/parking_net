package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/Nerzal/gocloak/v13"
)

func (c Client) GetAdminToken(ctx context.Context) (*gocloak.JWT, error) {
	if c.Config.Admin == "" {
		return nil, errors.New("admin username is required")
	}
	if c.Config.AdminPassword == "" {
		return nil, errors.New("admin password is required")
	}

	token, err := c.Client.LoginAdmin(ctx, c.Config.Admin, c.Config.AdminPassword, c.Config.MasterRealm)
	if err != nil {
		return nil, fmt.Errorf("failed to login as admin: %w", err)
	}
	if token == nil {
		return nil, errors.New("received nil token")
	}
	return token, nil
}
