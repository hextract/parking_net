package database_service

import (
	"context"
	"fmt"
)

func (ds *DatabaseService) AddBookingServices(ctx context.Context, bookingID int64, services []BookingService) error {
	if len(services) == 0 {
		return nil
	}

	tx, err := ds.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	for _, service := range services {
		query := `INSERT INTO booking_services (booking_id, service_id, quantity, price)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (booking_id, service_id) DO UPDATE SET quantity = $3, price = $4`

		_, err := tx.Exec(ctx, query, bookingID, service.ServiceID, service.Quantity, service.Price)
		if err != nil {
			return fmt.Errorf("failed to insert booking service: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (ds *DatabaseService) GetBookingServices(ctx context.Context, bookingID int64) ([]BookingService, error) {
	query := `SELECT booking_id, service_id, quantity, price
		FROM booking_services WHERE booking_id = $1`

	rows, err := ds.pool.Query(ctx, query, bookingID)
	if err != nil {
		return nil, fmt.Errorf("failed to query booking services: %w", err)
	}
	defer rows.Close()

	var services []BookingService
	for rows.Next() {
		var service BookingService
		err := rows.Scan(
			&service.BookingID,
			&service.ServiceID,
			&service.Quantity,
			&service.Price,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan booking service: %w", err)
		}
		services = append(services, service)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating booking services: %w", err)
	}

	return services, nil
}

type BookingService struct {
	BookingID int64
	ServiceID int64
	Quantity  int64
	Price     int64
}

type Service struct {
	ID          int64
	ParkingID   int64
	Name        string
	Description string
	Price       int64
}
