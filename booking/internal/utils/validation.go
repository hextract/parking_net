package utils

import (
	"errors"
	"fmt"
	"regexp"
	"time"
	"unicode/utf8"
)

const (
	MaxBookingID     = 9223372036854775807
	MinBookingID     = 1
	MaxFullCost      = 1000000000000
	MinFullCost      = 0
	MaxHoursDuration = 8760
	MinHoursDuration = 0
	MaxStringLength  = 500
	MinStringLength  = 1
)

var (
	ErrInvalidBookingID     = errors.New("invalid booking ID")
	ErrInvalidParkingPlaceID = errors.New("invalid parking place ID")
	ErrInvalidUserID        = errors.New("invalid user ID")
	ErrInvalidFullCost      = errors.New("invalid full cost")
	ErrInvalidDateRange     = errors.New("invalid date range")
	ErrDateTooFarInFuture   = errors.New("date too far in future")
	ErrDateInPast           = errors.New("date cannot be in the past")
	ErrInvalidStringLength  = errors.New("invalid string length")
)

func ValidateBookingID(bookingID int64) error {
	if bookingID < MinBookingID || bookingID > MaxBookingID {
		return fmt.Errorf("%w: must be between %d and %d", ErrInvalidBookingID, MinBookingID, MaxBookingID)
	}
	return nil
}

func ValidateParkingPlaceID(parkingPlaceID *int64) error {
	if parkingPlaceID == nil {
		return fmt.Errorf("%w: cannot be nil", ErrInvalidParkingPlaceID)
	}
	if *parkingPlaceID < MinBookingID || *parkingPlaceID > MaxBookingID {
		return fmt.Errorf("%w: must be between %d and %d", ErrInvalidParkingPlaceID, MinBookingID, MaxBookingID)
	}
	return nil
}

func ValidateUserID(userID string) error {
	if userID == "" {
		return fmt.Errorf("%w: cannot be empty", ErrInvalidUserID)
	}
	if len(userID) > 128 {
		return fmt.Errorf("%w: too long", ErrInvalidUserID)
	}
	uuidRegex := regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`)
	if !uuidRegex.MatchString(userID) {
		return fmt.Errorf("%w: invalid UUID format", ErrInvalidUserID)
	}
	return nil
}

func ValidateFullCost(fullCost int64) error {
	if fullCost < MinFullCost {
		return fmt.Errorf("%w: cannot be negative", ErrInvalidFullCost)
	}
	if fullCost > MaxFullCost {
		return fmt.Errorf("%w: exceeds maximum allowed value", ErrInvalidFullCost)
	}
	return nil
}

func ValidateDateRange(dateFrom, dateTo *time.Time) error {
	if dateFrom == nil || dateTo == nil {
		return fmt.Errorf("%w: dates cannot be nil", ErrInvalidDateRange)
	}

	if dateFrom.After(*dateTo) || dateFrom.Equal(*dateTo) {
		return fmt.Errorf("%w: date_from must be before date_to", ErrInvalidDateRange)
	}

	now := time.Now()
	maxFutureDate := now.AddDate(1, 0, 0)

	if dateFrom.After(maxFutureDate) || dateTo.After(maxFutureDate) {
		return fmt.Errorf("%w: cannot be more than 1 year in the future", ErrDateTooFarInFuture)
	}

	duration := dateTo.Sub(*dateFrom)
	hours := duration.Hours()
	if hours < MinHoursDuration || hours > MaxHoursDuration {
		return fmt.Errorf("%w: duration must be between %d and %d hours", ErrInvalidDateRange, MinHoursDuration, MaxHoursDuration)
	}

	return nil
}

func ValidateString(str string, fieldName string) error {
	length := utf8.RuneCountInString(str)
	if length < MinStringLength {
		return fmt.Errorf("%w: %s cannot be empty", ErrInvalidStringLength, fieldName)
	}
	if length > MaxStringLength {
		return fmt.Errorf("%w: %s exceeds maximum length of %d", ErrInvalidStringLength, fieldName, MaxStringLength)
	}
	return nil
}

func SanitizeError(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("operation failed")
}

