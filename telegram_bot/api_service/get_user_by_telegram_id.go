package api_service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Nerzal/gocloak/v13"
)

type UserInfo struct {
	UserID     string
	Login      string
	Email      string
	Role       string
	TelegramID int64
	Token      string
}

func (s *Service) GetUserByTelegramID(ctx context.Context, telegramID int64) (*UserInfo, error) {
	adminToken, err := s.Client.GetAdminToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get admin token: %w", err)
	}

	telegramIDStr := strconv.FormatInt(telegramID, 10)

	maxUsers := 1000
	first := 0

	for {
		params := gocloak.GetUsersParams{
			First: &first,
			Max:   &maxUsers,
		}

		users, err := s.Client.Client.GetUsers(ctx, adminToken.AccessToken, s.Client.Config.Realm, params)
		if err != nil {
			return nil, fmt.Errorf("failed to get users: %w", err)
		}

		if len(users) == 0 {
			break
		}

		for _, user := range users {
			if user.ID == nil {
				continue
			}

			fullUser, err := s.Client.Client.GetUserByID(ctx, adminToken.AccessToken, s.Client.Config.Realm, *user.ID)
			if err != nil {
				continue
			}

			if fullUser.Attributes == nil {
				continue
			}

			tgIDAttr, exists := (*fullUser.Attributes)["telegram_id"]
			if !exists || len(tgIDAttr) == 0 {
				continue
			}

			if tgIDAttr[0] == telegramIDStr {
				groups, err := s.Client.Client.GetUserGroups(ctx, adminToken.AccessToken, s.Client.Config.Realm, *user.ID, gocloak.GetGroupsParams{})
				if err != nil || len(groups) == 0 {
					continue
				}

				role := "driver"
				if groups[0].Name != nil {
					role = *groups[0].Name
				}

				login := ""
				if fullUser.Username != nil {
					login = *fullUser.Username
				}

				email := ""
				if fullUser.Email != nil {
					email = *fullUser.Email
				}

				token, err := s.DatabaseService.GetToken(telegramID)
				if err != nil {
					return nil, fmt.Errorf("failed to get token: %w", err)
				}

				tokenValue := ""
				if token != nil {
					tokenValue = token.Value
				}

				return &UserInfo{
					UserID:     *user.ID,
					Login:      login,
					Email:      email,
					Role:       role,
					TelegramID: telegramID,
					Token:      tokenValue,
				}, nil
			}
		}

		if len(users) < maxUsers {
			break
		}

		first += maxUsers
	}

	return nil, fmt.Errorf("user not found")
}
