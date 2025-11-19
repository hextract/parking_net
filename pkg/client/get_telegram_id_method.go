package client

import (
	"context"
	"errors"
	"fmt"
	"strconv"
)

func (c Client) GetTelegramId(ctx context.Context, userId string) (tgId int, err error) {
	if userId == "" {
		return 0, errors.New("user ID is required")
	}

	adminToken, err := c.GetAdminToken(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get admin token: %w", err)
	}
	user, err := c.Client.GetUserByID(ctx, adminToken.AccessToken, c.Config.Realm, userId)
	if err != nil {
		return 0, fmt.Errorf("failed to get user by ID: %w", err)
	}
	if user == nil {
		return 0, errors.New("user not found")
	}
	if user.Attributes == nil {
		return 0, errors.New("user attributes are missing")
	}
	telegramIDAttr, exists := (*user.Attributes)["telegram_id"]
	if !exists || len(telegramIDAttr) == 0 || telegramIDAttr[0] == "" {
		return 0, errors.New("telegram ID is missing")
	}
	tgId, err = strconv.Atoi(telegramIDAttr[0])
	if err != nil {
		return 0, fmt.Errorf("invalid telegram ID format: %w", err)
	}
	return tgId, nil
}
