package database_service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/go-openapi/strfmt"
	"github.com/h4x4d/parking_net/payment/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	"time"
)

func (ds *DatabaseService) CreatePromocode(ctx context.Context, adminID string, amount int64, maxUses int64, customCode *string, expiresAt *strfmt.DateTime) (*models.PromocodeResponse, error) {
	tracer := otel.Tracer("Payment")
	ctx, span := tracer.Start(ctx, "create_promocode")
	defer span.End()

	var code string
	var err error

	if customCode != nil && *customCode != "" {
		var exists bool
		err = ds.pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM promocodes WHERE code = $1)", *customCode).Scan(&exists)
		if err != nil {
			return nil, fmt.Errorf("failed to check code existence: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("promocode already exists")
		}
		code = *customCode
	} else {
		code, err = generateUniqueCodeStandalone(ctx, ds.pool)
		if err != nil {
			return nil, fmt.Errorf("failed to generate code: %w", err)
		}
	}

	var expiresAtTime *time.Time
	if expiresAt != nil {
		t := time.Time(*expiresAt)
		expiresAtTime = &t
	}

	_, err = ds.pool.Exec(ctx,
		"INSERT INTO promocodes (code, amount, max_uses, created_by, expires_at) VALUES ($1, $2, $3, $4, $5)",
		code, amount, maxUses, adminID, expiresAtTime)
	if err != nil {
		return nil, fmt.Errorf("failed to create promocode: %w", err)
	}

	var expiresAtDT *strfmt.DateTime
	if expiresAtTime != nil {
		dt := strfmt.DateTime(*expiresAtTime)
		expiresAtDT = &dt
	}

	response := &models.PromocodeResponse{
		Code:          code,
		Amount:        amount,
		MaxUses:       maxUses,
		RemainingUses: maxUses,
	}
	if expiresAtDT != nil {
		response.ExpiresAt = *expiresAtDT
	}
	return response, nil
}

func generateUniqueCodeStandalone(ctx context.Context, pool *pgxpool.Pool) (string, error) {
	for i := 0; i < 10; i++ {
		bytes := make([]byte, 8)
		if _, err := rand.Read(bytes); err != nil {
			return "", err
		}
		code := hex.EncodeToString(bytes)
		code = code[:16]

		var exists bool
		err := pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM promocodes WHERE code = $1)", code).Scan(&exists)
		if err != nil {
			continue
		}
		if !exists {
			return code, nil
		}
	}
	return "", fmt.Errorf("failed to generate unique code")
}

