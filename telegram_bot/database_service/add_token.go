package database_service

import (
	"context"
	"fmt"
	"telegram_bot/models"
)

func (ds *DatabaseService) AddToken(telegramID int64, token *models.ApiToken) error {
	if token == nil {
		return fmt.Errorf("token cannot be nil")
	}

	query := `INSERT INTO telegram (token, telegram_id) VALUES ($1, $2)
		ON CONFLICT (telegram_id) DO UPDATE SET token = $1`
	_, errInsert := ds.pool.Exec(context.Background(), query, token.Value, telegramID)
	if errInsert != nil {
		return errInsert
	}

	return nil
}
