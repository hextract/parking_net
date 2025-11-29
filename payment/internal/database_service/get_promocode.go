package database_service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-openapi/strfmt"
	"github.com/h4x4d/parking_net/payment/internal/models"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"time"
)

func (ds *DatabaseService) GetPromocode(ctx context.Context, code string) (*models.PromocodeInfo, error) {
	tracer := otel.Tracer("Payment")
	ctx, span := tracer.Start(ctx, "get_promocode")
	defer span.End()

	var promocode models.PromocodeInfo
	var expiresAt sql.NullTime

	err := ds.pool.QueryRow(ctx,
		"SELECT code, amount, max_uses, used_count, expires_at FROM promocodes WHERE code = $1",
		code).Scan(&promocode.Code, &promocode.Amount, &promocode.MaxUses, &promocode.UsedCount, &expiresAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if expiresAt.Valid {
		promocode.ExpiresAt = strfmt.DateTime(expiresAt.Time)
	}

	promocode.RemainingUses = promocode.MaxUses - promocode.UsedCount
	promocode.IsActive = promocode.RemainingUses > 0
	if !promocode.ExpiresAt.IsZero() {
		if time.Time(promocode.ExpiresAt).Before(time.Now()) {
			promocode.IsActive = false
		}
	}

	return &promocode, nil
}

