package repository

import (
	"context"

	"github.com/h4x4d/parking_net/pkg/domain"
)

type ServiceRepository interface {
	Create(ctx context.Context, service *domain.Service) (*domain.Service, error)
	GetByID(ctx context.Context, id int64) (*domain.Service, error)
	GetByParkingID(ctx context.Context, parkingID int64) ([]*domain.Service, error)
	Update(ctx context.Context, service *domain.Service) error
	Delete(ctx context.Context, id int64, parkingID int64) error
	Exists(ctx context.Context, id int64) (bool, error)
}
