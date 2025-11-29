package database_service

import (
	"context"
	"database/sql"
	"github.com/go-openapi/strfmt"
	"github.com/h4x4d/parking_net/payment/internal/models"
	"go.opentelemetry.io/otel"
	"time"
)

func (ds *DatabaseService) GetTransactions(ctx context.Context, userID string, limit int64, offset int64) ([]*models.Transaction, error) {
	tracer := otel.Tracer("Payment")
	ctx, span := tracer.Start(ctx, "get_transactions")
	defer span.End()

	rows, err := ds.pool.Query(ctx,
		"SELECT id, booking_id, amount, transaction_type, status, description, created_at FROM transactions WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3",
		userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*models.Transaction
	for rows.Next() {
		var t models.Transaction
		var bookingID sql.NullInt64
		var createdAt time.Time

		err := rows.Scan(&t.ID, &bookingID, &t.Amount, &t.TransactionType, &t.Status, &t.Description, &createdAt)
		if err != nil {
			return nil, err
		}

		if bookingID.Valid {
			t.BookingID = bookingID.Int64
		}

		t.UserID = userID
		t.CreatedAt = strfmt.DateTime(createdAt)

		transactions = append(transactions, &t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

