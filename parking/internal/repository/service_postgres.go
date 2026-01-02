package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/h4x4d/parking_net/pkg/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresServiceRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresServiceRepository(pool *pgxpool.Pool) ServiceRepository {
	return &PostgresServiceRepository{pool: pool}
}

func (r *PostgresServiceRepository) Create(ctx context.Context, service *domain.Service) (*domain.Service, error) {
	if service.Name == "" {
		return nil, fmt.Errorf("service name is required")
	}
	if service.Price < 0 {
		return nil, fmt.Errorf("service price cannot be negative")
	}

	query := `INSERT INTO services (parking_id, name, description, price)
		VALUES ($1, $2, $3, $4) RETURNING id`

	err := r.pool.QueryRow(ctx, query,
		service.ParkingID,
		service.Name,
		service.Description,
		service.Price,
	).Scan(&service.ID)

	if err != nil {
		return nil, fmt.Errorf("failed to create service: %w", err)
	}

	return service, nil
}

func (r *PostgresServiceRepository) GetByID(ctx context.Context, id int64) (*domain.Service, error) {
	query := `SELECT id, parking_id, name, description, price
		FROM services WHERE id = $1`

	var service domain.Service
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&service.ID,
		&service.ParkingID,
		&service.Name,
		&service.Description,
		&service.Price,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get service by id: %w", err)
	}

	return &service, nil
}

func (r *PostgresServiceRepository) GetByParkingID(ctx context.Context, parkingID int64) ([]*domain.Service, error) {
	query := `SELECT id, parking_id, name, description, price
		FROM services WHERE parking_id = $1 ORDER BY id`

	rows, err := r.pool.Query(ctx, query, parkingID)
	if err != nil {
		return nil, fmt.Errorf("failed to query services: %w", err)
	}
	defer rows.Close()

	var services []*domain.Service
	for rows.Next() {
		var service domain.Service
		err := rows.Scan(
			&service.ID,
			&service.ParkingID,
			&service.Name,
			&service.Description,
			&service.Price,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan service: %w", err)
		}
		services = append(services, &service)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating services: %w", err)
	}

	return services, nil
}

func (r *PostgresServiceRepository) Update(ctx context.Context, service *domain.Service) error {
	if service.Name == "" {
		return fmt.Errorf("service name is required")
	}
	if service.Price < 0 {
		return fmt.Errorf("service price cannot be negative")
	}

	query := `UPDATE services 
		SET name = $1, description = $2, price = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4 AND parking_id = $5`

	result, err := r.pool.Exec(ctx, query,
		service.Name,
		service.Description,
		service.Price,
		service.ID,
		service.ParkingID,
	)

	if err != nil {
		return fmt.Errorf("failed to update service: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("service not found or access denied")
	}

	return nil
}

func (r *PostgresServiceRepository) Delete(ctx context.Context, id int64, parkingID int64) error {
	query := `DELETE FROM services WHERE id = $1 AND parking_id = $2`

	result, err := r.pool.Exec(ctx, query, id, parkingID)
	if err != nil {
		return fmt.Errorf("failed to delete service: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("service not found or access denied")
	}

	return nil
}

func (r *PostgresServiceRepository) Exists(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM services WHERE id = $1)`

	var exists bool
	err := r.pool.QueryRow(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check service existence: %w", err)
	}

	return exists, nil
}
