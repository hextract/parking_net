package service

import (
	"context"

	"github.com/h4x4d/parking_net/parking/internal/repository"
	"github.com/h4x4d/parking_net/pkg/domain"
	"github.com/h4x4d/parking_net/pkg/errors"
)

type ParkingService struct {
	repo repository.ParkingRepository
}

func NewParkingService(repo repository.ParkingRepository) *ParkingService {
	return &ParkingService{repo: repo}
}

func (s *ParkingService) CreateParking(ctx context.Context, parking *domain.ParkingPlace, user *domain.User) (*domain.ParkingPlace, *errors.AppError) {
	if !user.IsOwner() {
		return nil, errors.ErrForbidden
	}

	if err := parking.IsValid(); err != nil {
		return nil, errors.Validation(err.Error())
	}

	parking.OwnerID = user.ID

	created, err := s.repo.Create(ctx, parking)
	if err != nil {
		return nil, errors.Internal(err)
	}

	return created, nil
}

func (s *ParkingService) GetParkingByID(ctx context.Context, id int64) (*domain.ParkingPlace, *errors.AppError) {
	parking, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.Internal(err)
	}

	if parking == nil {
		return nil, errors.NotFound("parking place")
	}

	return parking, nil
}

func (s *ParkingService) GetParkings(ctx context.Context, filters repository.ParkingFilters) ([]*domain.ParkingPlace, *errors.AppError) {
	parkings, err := s.repo.GetAll(ctx, filters)
	if err != nil {
		return nil, errors.Internal(err)
	}

	return parkings, nil
}

type ParkingFilters = repository.ParkingFilters

func (s *ParkingService) UpdateParking(ctx context.Context, id int64, parking *domain.ParkingPlace, user *domain.User) *errors.AppError {
	if !user.IsOwner() && !user.IsAdmin() {
		return errors.ErrForbidden
	}

	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.Internal(err)
	}

	if existing == nil {
		return errors.NotFound("parking place")
	}

	if !user.IsAdmin() && existing.OwnerID != user.ID {
		return errors.ErrForbidden
	}

	parking.ID = id
	if !user.IsAdmin() {
		parking.OwnerID = user.ID
	} else {
		parking.OwnerID = existing.OwnerID
	}

	if err := parking.IsValid(); err != nil {
		return errors.Validation(err.Error())
	}

	if err := s.repo.Update(ctx, parking); err != nil {
		return errors.Internal(err)
	}

	return nil
}

func (s *ParkingService) DeleteParking(ctx context.Context, id int64, user *domain.User) *errors.AppError {
	if !user.IsOwner() && !user.IsAdmin() {
		return errors.ErrForbidden
	}

	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.Internal(err)
	}

	if existing == nil {
		return errors.NotFound("parking place")
	}

	if !user.IsAdmin() && existing.OwnerID != user.ID {
		return errors.ErrForbidden
	}

	ownerID := user.ID
	if user.IsAdmin() {
		ownerID = existing.OwnerID
	}

	if err := s.repo.Delete(ctx, id, ownerID); err != nil {
		return errors.Internal(err)
	}

	return nil
}
