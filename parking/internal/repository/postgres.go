package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/h4x4d/parking_net/parking/internal/utils"
	"github.com/h4x4d/parking_net/pkg/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresParkingRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresParkingRepository(pool *pgxpool.Pool) ParkingRepository {
	return &PostgresParkingRepository{pool: pool}
}

func (r *PostgresParkingRepository) Create(ctx context.Context, parking *domain.ParkingPlace) (*domain.ParkingPlace, error) {
	if err := utils.ValidateOwnerID(parking.OwnerID); err != nil {
		return nil, fmt.Errorf("invalid owner ID")
	}

	if err := utils.ValidateString(parking.Name, "name"); err != nil {
		return nil, fmt.Errorf("invalid name")
	}

	if err := utils.ValidateString(parking.City, "city"); err != nil {
		return nil, fmt.Errorf("invalid city")
	}

	if err := utils.ValidateString(parking.Address, "address"); err != nil {
		return nil, fmt.Errorf("invalid address")
	}

	if err := utils.ValidateParkingType(string(parking.Type)); err != nil {
		return nil, fmt.Errorf("invalid parking type")
	}

	if err := utils.ValidateHourlyRate(int64(parking.HourlyRate)); err != nil {
		return nil, fmt.Errorf("invalid hourly rate")
	}

	if err := utils.ValidateCapacity(int64(parking.Capacity)); err != nil {
		return nil, fmt.Errorf("invalid capacity")
	}

	query := `INSERT INTO parking_places (name, city, address, parking_type, hourly_rate, capacity, owner_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	err := r.pool.QueryRow(ctx, query,
		parking.Name,
		parking.City,
		parking.Address,
		string(parking.Type),
		parking.HourlyRate,
		parking.Capacity,
		parking.OwnerID,
	).Scan(&parking.ID)

	if err != nil {
		return nil, fmt.Errorf("failed to create parking place")
	}

	return parking, nil
}

func (r *PostgresParkingRepository) GetByID(ctx context.Context, id int64) (*domain.ParkingPlace, error) {
	query := `SELECT id, name, city, address, parking_type, hourly_rate, capacity, owner_id
		FROM parking_places WHERE id = $1`

	var parking domain.ParkingPlace
	var parkingType string

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&parking.ID,
		&parking.Name,
		&parking.City,
		&parking.Address,
		&parkingType,
		&parking.HourlyRate,
		&parking.Capacity,
		&parking.OwnerID,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get parking place by id")
	}

	parking.Type = domain.ParkingType(parkingType)
	return &parking, nil
}

func (r *PostgresParkingRepository) GetAll(ctx context.Context, filters ParkingFilters) ([]*domain.ParkingPlace, error) {
	query := `SELECT id, name, city, address, parking_type, hourly_rate, capacity, owner_id
		FROM parking_places`

	var clauses []string
	var args []interface{}
	argIndex := 1

	if filters.City != nil {
		clauses = append(clauses, fmt.Sprintf("city = $%d", argIndex))
		args = append(args, *filters.City)
		argIndex++
	}
	if filters.Name != nil {
		clauses = append(clauses, fmt.Sprintf("name ILIKE $%d", argIndex))
		args = append(args, "%"+*filters.Name+"%")
		argIndex++
	}
	if filters.ParkingType != nil {
		clauses = append(clauses, fmt.Sprintf("parking_type = $%d", argIndex))
		args = append(args, string(*filters.ParkingType))
		argIndex++
	}
	if filters.OwnerID != nil {
		clauses = append(clauses, fmt.Sprintf("owner_id = $%d", argIndex))
		args = append(args, *filters.OwnerID)
		argIndex++
	}

	if len(clauses) > 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query parking places")
	}
	defer rows.Close()

	var parkings []*domain.ParkingPlace
	for rows.Next() {
		var parking domain.ParkingPlace
		var parkingType string

		err := rows.Scan(
			&parking.ID,
			&parking.Name,
			&parking.City,
			&parking.Address,
			&parkingType,
			&parking.HourlyRate,
			&parking.Capacity,
			&parking.OwnerID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan parking place")
		}

		parking.Type = domain.ParkingType(parkingType)
		parkings = append(parkings, &parking)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating parking places")
	}

	return parkings, nil
}

func (r *PostgresParkingRepository) Update(ctx context.Context, parking *domain.ParkingPlace) error {
	if err := utils.ValidateParkingID(parking.ID); err != nil {
		return fmt.Errorf("invalid parking ID")
	}

	if err := utils.ValidateOwnerID(parking.OwnerID); err != nil {
		return fmt.Errorf("invalid owner ID")
	}

	if parking.Name != "" {
		if err := utils.ValidateString(parking.Name, "name"); err != nil {
			return fmt.Errorf("invalid name")
		}
	}

	if parking.City != "" {
		if err := utils.ValidateString(parking.City, "city"); err != nil {
			return fmt.Errorf("invalid city")
		}
	}

	if parking.Address != "" {
		if err := utils.ValidateString(parking.Address, "address"); err != nil {
			return fmt.Errorf("invalid address")
		}
	}

	if parking.Type != "" {
		if err := utils.ValidateParkingType(string(parking.Type)); err != nil {
			return fmt.Errorf("invalid parking type")
		}
	}

	if parking.HourlyRate > 0 {
		if err := utils.ValidateHourlyRate(int64(parking.HourlyRate)); err != nil {
			return fmt.Errorf("invalid hourly rate")
		}
	}

	if parking.Capacity > 0 {
		if err := utils.ValidateCapacity(int64(parking.Capacity)); err != nil {
			return fmt.Errorf("invalid capacity")
		}
	}

	query := `UPDATE parking_places 
		SET name = $1, city = $2, address = $3, parking_type = $4, hourly_rate = $5, capacity = $6
		WHERE id = $7 AND owner_id = $8`

	result, err := r.pool.Exec(ctx, query,
		parking.Name,
		parking.City,
		parking.Address,
		string(parking.Type),
		parking.HourlyRate,
		parking.Capacity,
		parking.ID,
		parking.OwnerID,
	)

	if err != nil {
		return fmt.Errorf("failed to update parking place")
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("parking place not found or access denied")
	}

	return nil
}

func (r *PostgresParkingRepository) Exists(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM parking_places WHERE id = $1)`

	var exists bool
	err := r.pool.QueryRow(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check parking place existence")
	}

	return exists, nil
}

func (r *PostgresParkingRepository) GetByOwnerID(ctx context.Context, ownerID string) ([]*domain.ParkingPlace, error) {
	return r.GetAll(ctx, ParkingFilters{OwnerID: &ownerID})
}

func (r *PostgresParkingRepository) Delete(ctx context.Context, id int64, ownerID string) error {
	query := `DELETE FROM parking_places WHERE id = $1 AND owner_id = $2`

	result, err := r.pool.Exec(ctx, query, id, ownerID)
	if err != nil {
		return fmt.Errorf("failed to delete parking place")
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("parking place not found or access denied")
	}

	return nil
}
