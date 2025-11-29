package database_service

import (
	"context"
	"errors"
	"github.com/h4x4d/parking_net/payment/internal/models"
	"github.com/jackc/pgx/v5"
)

func (ds *DatabaseService) GetBalance(userID string) (*models.Balance, error) {
	var balance models.Balance
	balance.UserID = &userID
	currency := "USD"
	balance.Currency = &currency

	var balanceValue int64
	err := ds.pool.QueryRow(context.Background(),
		"SELECT balance FROM balances WHERE user_id = $1", userID).Scan(&balanceValue)

	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, err
		}
		if errors.Is(err, pgx.ErrNoRows) {
			_, insertErr := ds.pool.Exec(context.Background(),
				"INSERT INTO balances (user_id, balance, currency) VALUES ($1, 0, 'USD') ON CONFLICT DO NOTHING",
				userID)
			if insertErr != nil {
				return nil, insertErr
			}
			balanceValue = 0
		} else {
			return nil, err
		}
	}

	balance.Balance = &balanceValue
	return &balance, nil
}

