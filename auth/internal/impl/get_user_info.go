package impl

import (
	"context"
	"fmt"
	"strconv"

	"github.com/h4x4d/parking_net/auth/internal/restapi/operations"
	"github.com/h4x4d/parking_net/pkg/client"
)

func GetUserInfo(ctx context.Context, clt *client.Client, params operations.GetAuthMeParams) (*operations.GetAuthMeOKBody, error) {
	// Get token from params (swagger binds it from header)
	token := params.APIKey

	if token == "" {
		return nil, fmt.Errorf("token is required")
	}

	// Use CheckToken to get user info
	user, err := clt.CheckToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info")
	}

	// Get user details from Keycloak
	adminToken, err := clt.GetAdminToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get admin token")
	}

	keycloakUser, err := clt.Client.GetUserByID(ctx, adminToken.AccessToken, clt.Config.Realm, user.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user details")
	}

	login := ""
	if keycloakUser.Username != nil {
		login = *keycloakUser.Username
	}

	email := ""
	if keycloakUser.Email != nil {
		email = *keycloakUser.Email
	}

	telegramID := int64(0)
	if keycloakUser.Attributes != nil {
		if tgIDAttr, exists := (*keycloakUser.Attributes)["telegram_id"]; exists && len(tgIDAttr) > 0 {
			if tgID, err := strconv.ParseInt(tgIDAttr[0], 10, 64); err == nil {
				telegramID = tgID
			}
		}
	}

	// Determine role
	role := user.Role
	if role != "driver" && role != "owner" && role != "admin" {
		role = "driver" // default
	}

	return &operations.GetAuthMeOKBody{
		UserID:     user.UserID,
		Login:      login,
		Email:      email,
		Role:       role,
		TelegramID: telegramID,
	}, nil
}
