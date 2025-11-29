package database_service

import (
	"context"
	"errors"
	"fmt"
	"github.com/h4x4d/parking_net/payment/internal/models"
	"github.com/h4x4d/parking_net/payment/internal/utils"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"time"
)

func (ds *DatabaseService) ActivatePromocode(ctx context.Context, userID string, code string) (*models.Balance, error) {
	if err := utils.ValidateUserID(userID); err != nil {
		return nil, fmt.Errorf("invalid user ID")
	}

	if err := utils.ValidatePromoCode(code); err != nil {
		return nil, fmt.Errorf("invalid promocode format")
	}

	tracer := otel.Tracer("Payment")
	ctx, span := tracer.Start(ctx, "activate_promocode")
	defer span.End()

	tx, err := ds.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var promocodeAmount int64
	var maxUses int
	var usedCount int
	var expiresAt *time.Time

	err = tx.QueryRow(ctx,
		"SELECT amount, max_uses, used_count, expires_at FROM promocodes WHERE code = $1 FOR UPDATE",
		code).Scan(&promocodeAmount, &maxUses, &usedCount, &expiresAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("promocode not found")
		}
		return nil, fmt.Errorf("failed to get promocode: %w", err)
	}

	if usedCount >= maxUses {
		return nil, fmt.Errorf("promocode has reached maximum uses")
	}

	if expiresAt != nil && expiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("promocode has expired")
	}

	var balance int64
	err = tx.QueryRow(ctx, "SELECT balance FROM balances WHERE user_id = $1 FOR UPDATE", userID).Scan(&balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			_, err = tx.Exec(ctx, "INSERT INTO balances (user_id, balance, currency) VALUES ($1, 0, 'USD')", userID)
			if err != nil {
				return nil, fmt.Errorf("failed to create balance: %w", err)
			}
			balance = 0
		} else {
			return nil, fmt.Errorf("failed to get balance: %w", err)
		}
	}

	newBalance, err := utils.SafeAddBalance(balance, promocodeAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to update balance")
	}

	_, err = tx.Exec(ctx, "UPDATE balances SET balance = $1 WHERE user_id = $2", newBalance, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}

	newUsedCount := usedCount + 1
	_, err = tx.Exec(ctx, "UPDATE promocodes SET used_count = $1 WHERE code = $2", newUsedCount, code)
	if err != nil {
		return nil, fmt.Errorf("failed to update promocode: %w", err)
	}

	_, err = tx.Exec(ctx,
		"INSERT INTO transactions (user_id, amount, transaction_type, status, description) VALUES ($1, $2, 'promocode_activate', 'completed', $3)",
		userID, promocodeAmount, fmt.Sprintf("Activated promocode %s", code))
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	currency := "USD"
	return &models.Balance{
		UserID:   &userID,
		Balance:  &newBalance,
		Currency: &currency,
	}, nil
}

