package database_service

import (
	"context"
	"errors"
	"fmt"
	"github.com/h4x4d/parking_net/payment/internal/models"
	"github.com/h4x4d/parking_net/payment/internal/utils"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
)

func (ds *DatabaseService) ProcessRefund(ctx context.Context, bookingID int64, driverID string, ownerID string, amount int64) (*models.TransactionResponse, error) {
	if err := utils.ValidateAmount(amount); err != nil {
		return &models.TransactionResponse{
			Status:  "failed",
			Message: "invalid amount",
		}, nil
	}

	if err := utils.ValidateUserID(driverID); err != nil {
		return &models.TransactionResponse{
			Status:  "failed",
			Message: "invalid driver ID",
		}, nil
	}

	if err := utils.ValidateUserID(ownerID); err != nil {
		return &models.TransactionResponse{
			Status:  "failed",
			Message: "invalid owner ID",
		}, nil
	}

	if err := utils.ValidateBookingID(bookingID); err != nil {
		return &models.TransactionResponse{
			Status:  "failed",
			Message: "invalid booking ID",
		}, nil
	}

	tracer := otel.Tracer("Payment")
	ctx, span := tracer.Start(ctx, "process_refund")
	defer span.End()

	tx, err := ds.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var ownerBalance int64
	err = tx.QueryRow(ctx, "SELECT balance FROM balances WHERE user_id = $1 FOR UPDATE", ownerID).Scan(&ownerBalance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("owner balance not found")
		}
		return nil, fmt.Errorf("failed to get owner balance: %w", err)
	}

	if ownerBalance < amount {
		return &models.TransactionResponse{
			Status:  "failed",
			Message: "owner has insufficient funds for refund",
		}, nil
	}

	var driverBalance int64
	err = tx.QueryRow(ctx, "SELECT balance FROM balances WHERE user_id = $1 FOR UPDATE", driverID).Scan(&driverBalance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			_, err = tx.Exec(ctx, "INSERT INTO balances (user_id, balance, currency) VALUES ($1, 0, 'USD')", driverID)
			if err != nil {
				return nil, fmt.Errorf("failed to create driver balance: %w", err)
			}
			driverBalance = 0
		} else {
			return nil, fmt.Errorf("failed to get driver balance: %w", err)
		}
	}

	newOwnerBalance, err := utils.SafeSubtractBalance(ownerBalance, amount)
	if err != nil {
		return &models.TransactionResponse{
			Status:  "failed",
			Message: err.Error(),
		}, nil
	}

	newDriverBalance, err := utils.SafeAddBalance(driverBalance, amount)
	if err != nil {
		return &models.TransactionResponse{
			Status:  "failed",
			Message: "refund failed",
		}, nil
	}

	_, err = tx.Exec(ctx, "UPDATE balances SET balance = $1 WHERE user_id = $2", newOwnerBalance, ownerID)
	if err != nil {
		return nil, fmt.Errorf("failed to update owner balance: %w", err)
	}

	_, err = tx.Exec(ctx, "UPDATE balances SET balance = $1 WHERE user_id = $2", newDriverBalance, driverID)
	if err != nil {
		return nil, fmt.Errorf("failed to update driver balance: %w", err)
	}

	var refundTransactionID int64
	err = tx.QueryRow(ctx,
		"INSERT INTO transactions (booking_id, user_id, amount, transaction_type, status, description) VALUES ($1, $2, $3, 'refund', 'completed', $4) RETURNING id",
		bookingID, driverID, amount, fmt.Sprintf("Refund for booking %d", bookingID)).Scan(&refundTransactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to create refund transaction: %w", err)
	}

	var chargebackTransactionID int64
	err = tx.QueryRow(ctx,
		"INSERT INTO transactions (booking_id, user_id, amount, transaction_type, status, description) VALUES ($1, $2, $3, 'charge', 'completed', $4) RETURNING id",
		bookingID, ownerID, -amount, fmt.Sprintf("Chargeback for booking %d refund", bookingID)).Scan(&chargebackTransactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to create chargeback transaction: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &models.TransactionResponse{
		TransactionID: refundTransactionID,
		Status:        "completed",
		Message:       "refund completed successfully",
	}, nil
}

