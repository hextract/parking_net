package database_service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/h4x4d/parking_net/payment/internal/models"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
)

func (ds *DatabaseService) GeneratePromocode(ctx context.Context, userID string, amount int64) (*models.PromocodeResponse, error) {
	tracer := otel.Tracer("Payment")
	ctx, span := tracer.Start(ctx, "generate_promocode")
	defer span.End()

	tx, err := ds.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var balance int64
	err = tx.QueryRow(ctx, "SELECT balance FROM balances WHERE user_id = $1 FOR UPDATE", userID).Scan(&balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			_, insertErr := tx.Exec(ctx,
				"INSERT INTO balances (user_id, balance, currency) VALUES ($1, 0, 'USD') ON CONFLICT DO NOTHING",
				userID)
			if insertErr != nil {
				return nil, fmt.Errorf("failed to create balance: %w", insertErr)
			}
			balance = 0
		} else {
			return nil, fmt.Errorf("failed to get balance: %w", err)
		}
	}

	if balance < amount {
		return nil, fmt.Errorf("insufficient funds")
	}

	newBalance := balance - amount
	_, err = tx.Exec(ctx, "UPDATE balances SET balance = $1 WHERE user_id = $2", newBalance, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}

	code, err := generateUniqueCodeTx(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate code: %w", err)
	}

	_, err = tx.Exec(ctx,
		"INSERT INTO promocodes (code, amount, max_uses, created_by) VALUES ($1, $2, 1, $3)",
		code, amount, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to create promocode: %w", err)
	}

	_, err = tx.Exec(ctx,
		"INSERT INTO transactions (user_id, amount, transaction_type, status, description) VALUES ($1, $2, 'promocode_generate', 'completed', $3)",
		userID, -amount, fmt.Sprintf("Generated promocode %s", code))
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &models.PromocodeResponse{
		Code:          code,
		Amount:        amount,
		MaxUses:       1,
		RemainingUses: 1,
	}, nil
}

func generateUniqueCodeTx(ctx context.Context, tx pgx.Tx) (string, error) {
	for i := 0; i < 10; i++ {
		bytes := make([]byte, 8)
		if _, err := rand.Read(bytes); err != nil {
			return "", err
		}
		code := hex.EncodeToString(bytes)
		code = code[:16]

		var exists bool
		err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM promocodes WHERE code = $1)", code).Scan(&exists)
		if err != nil {
			continue
		}
		if !exists {
			return code, nil
		}
	}
	return "", fmt.Errorf("failed to generate unique code")
}

