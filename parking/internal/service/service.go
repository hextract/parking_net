package service

import (
	"context"

	"github.com/h4x4d/parking_net/parking/internal/repository"
	"github.com/h4x4d/parking_net/parking/internal/utils"
	"github.com/h4x4d/parking_net/pkg/domain"
	"github.com/h4x4d/parking_net/pkg/errors"
)

type ServiceService struct {
	serviceRepo repository.ServiceRepository
	parkingRepo repository.ParkingRepository
}

func NewServiceService(serviceRepo repository.ServiceRepository, parkingRepo repository.ParkingRepository) *ServiceService {
	return &ServiceService{
		serviceRepo: serviceRepo,
		parkingRepo: parkingRepo,
	}
}

func (s *ServiceService) CreateService(ctx context.Context, service *domain.Service, user *domain.User) (*domain.Service, *errors.AppError) {
	if !user.IsOwner() {
		return nil, errors.ErrForbidden
	}

	parking, err := s.parkingRepo.GetByID(ctx, service.ParkingID)
	if err != nil {
		return nil, errors.Internal(utils.SanitizeError(err))
	}

	if parking == nil {
		return nil, errors.NotFound("parking place")
	}

	if parking.OwnerID != user.ID {
		return nil, errors.ErrForbidden
	}

	if service.Name == "" {
		return nil, errors.Validation("service name is required")
	}

	if service.Price < 0 {
		return nil, errors.Validation("service price cannot be negative")
	}

	created, err := s.serviceRepo.Create(ctx, service)
	if err != nil {
		return nil, errors.Internal(utils.SanitizeError(err))
	}

	return created, nil
}

func (s *ServiceService) GetServiceByID(ctx context.Context, id int64) (*domain.Service, *errors.AppError) {
	service, err := s.serviceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.Internal(utils.SanitizeError(err))
	}

	if service == nil {
		return nil, errors.NotFound("service")
	}

	return service, nil
}

func (s *ServiceService) GetServicesByParkingID(ctx context.Context, parkingID int64) ([]*domain.Service, *errors.AppError) {
	services, err := s.serviceRepo.GetByParkingID(ctx, parkingID)
	if err != nil {
		return nil, errors.Internal(utils.SanitizeError(err))
	}

	return services, nil
}

func (s *ServiceService) UpdateService(ctx context.Context, id int64, service *domain.Service, user *domain.User) *errors.AppError {
	if !user.IsOwner() {
		return errors.ErrForbidden
	}

	existing, err := s.serviceRepo.GetByID(ctx, id)
	if err != nil {
		return errors.Internal(utils.SanitizeError(err))
	}

	if existing == nil {
		return errors.NotFound("service")
	}

	parking, err := s.parkingRepo.GetByID(ctx, existing.ParkingID)
	if err != nil {
		return errors.Internal(utils.SanitizeError(err))
	}

	if parking == nil || parking.OwnerID != user.ID {
		return errors.ErrForbidden
	}

	service.ID = id
	service.ParkingID = existing.ParkingID

	if service.Name == "" {
		return errors.Validation("service name is required")
	}

	if service.Price < 0 {
		return errors.Validation("service price cannot be negative")
	}

	if err := s.serviceRepo.Update(ctx, service); err != nil {
		return errors.Internal(utils.SanitizeError(err))
	}

	return nil
}

func (s *ServiceService) DeleteService(ctx context.Context, id int64, user *domain.User) *errors.AppError {
	if !user.IsOwner() {
		return errors.ErrForbidden
	}

	existing, err := s.serviceRepo.GetByID(ctx, id)
	if err != nil {
		return errors.Internal(utils.SanitizeError(err))
	}

	if existing == nil {
		return errors.NotFound("service")
	}

	parking, err := s.parkingRepo.GetByID(ctx, existing.ParkingID)
	if err != nil {
		return errors.Internal(utils.SanitizeError(err))
	}

	if parking == nil || parking.OwnerID != user.ID {
		return errors.ErrForbidden
	}

	if err := s.serviceRepo.Delete(ctx, id, existing.ParkingID); err != nil {
		return errors.Internal(utils.SanitizeError(err))
	}

	return nil
}
