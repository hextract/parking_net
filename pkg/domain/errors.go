package domain

import "errors"

var (
	ErrInvalidParkingName    = errors.New("parking name is required")
	ErrInvalidParkingCity    = errors.New("parking city is required")
	ErrInvalidParkingAddress = errors.New("parking address is required")
	ErrInvalidHourlyRate      = errors.New("hourly rate must be greater than 0")
	ErrInvalidCapacity       = errors.New("capacity must be greater than 0")
	ErrInvalidParkingType     = errors.New("parking type is required")
	ErrInvalidOwnerID         = errors.New("owner ID is required")
	ErrInvalidDateFrom        = errors.New("date from is required")
	ErrInvalidDateTo          = errors.New("date to is required")
	ErrInvalidDateRange       = errors.New("date from must be before date to")
	ErrInvalidParkingPlaceID  = errors.New("parking place ID is required")
	ErrInvalidUserID          = errors.New("user ID is required")
	ErrParkingNotFound       = errors.New("parking place not found")
	ErrBookingNotFound       = errors.New("booking not found")
	ErrUnauthorizedAccess     = errors.New("unauthorized access")
	ErrForbiddenAction        = errors.New("forbidden action")
)

