package repository

import (
	"context"
	"github.com/h4x4d/parking_net/pkg/domain"
)

// BookingRepository defines the interface for booking data access
type BookingRepository interface {
	Create(ctx context.Context, booking *domain.Booking) (*domain.Booking, error)
	GetByID(ctx context.Context, id int64) (*domain.Booking, error)
	GetAll(ctx context.Context, filters BookingFilters) ([]*domain.Booking, error)
	Update(ctx context.Context, booking *domain.Booking) error
	Delete(ctx context.Context, id int64) error
	GetByParkingPlaceID(ctx context.Context, parkingPlaceID int64) ([]*domain.Booking, error)
	GetByUserID(ctx context.Context, userID string) ([]*domain.Booking, error)
}

// BookingFilters represents filters for querying bookings
type BookingFilters struct {
	ParkingPlaceID *int64
	UserID         *string
	Status         *domain.BookingStatus
}

