package api_service

import (
	"context"
	"fmt"
)

func (s *Service) GetOrCreateToken(ctx context.Context, userID string, login string, email string, telegramID int64) (string, error) {
	token, err := s.DatabaseService.GetToken(telegramID)
	if err == nil && token != nil && token.Value != "" {
		return token.Value, nil
	}

	return "", fmt.Errorf("token not found, please use /login command to authorize")
}
