package repository

import (
	"context"
	"github.com/h4x4d/parking_net/pkg/domain"
)

type ParkingRepository interface {
	Create(ctx context.Context, parking *domain.ParkingPlace) (*domain.ParkingPlace, error)
	GetByID(ctx context.Context, id int64) (*domain.ParkingPlace, error)
	GetAll(ctx context.Context, filters ParkingFilters) ([]*domain.ParkingPlace, error)
	Update(ctx context.Context, parking *domain.ParkingPlace) error
	Delete(ctx context.Context, id int64, ownerID string) error
	Exists(ctx context.Context, id int64) (bool, error)
	GetByOwnerID(ctx context.Context, ownerID string) ([]*domain.ParkingPlace, error)
}

type ParkingFilters struct {
	City       *string
	Name       *string
	ParkingType *domain.ParkingType
	OwnerID    *string
}

