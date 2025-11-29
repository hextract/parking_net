package utils

import (
	"errors"
	"fmt"
	"regexp"
	"unicode/utf8"
)

const (
	MaxParkingID    = 9223372036854775807
	MinParkingID    = 1
	MaxHourlyRate   = 1000000
	MinHourlyRate   = 0
	MaxCapacity     = 100000
	MinCapacity     = 1
	MaxStringLength = 500
	MinStringLength = 1
)

var (
	ErrInvalidParkingID    = errors.New("invalid parking ID")
	ErrInvalidOwnerID      = errors.New("invalid owner ID")
	ErrInvalidHourlyRate   = errors.New("invalid hourly rate")
	ErrInvalidCapacity     = errors.New("invalid capacity")
	ErrInvalidStringLength = errors.New("invalid string length")
	ErrInvalidParkingType  = errors.New("invalid parking type")
)

func ValidateParkingID(parkingID int64) error {
	if parkingID < MinParkingID || parkingID > MaxParkingID {
		return fmt.Errorf("%w: must be between %d and %d", ErrInvalidParkingID, MinParkingID, MaxParkingID)
	}
	return nil
}

func ValidateOwnerID(ownerID string) error {
	if ownerID == "" {
		return fmt.Errorf("%w: cannot be empty", ErrInvalidOwnerID)
	}
	if len(ownerID) > 128 {
		return fmt.Errorf("%w: too long", ErrInvalidOwnerID)
	}
	uuidRegex := regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`)
	if !uuidRegex.MatchString(ownerID) {
		return fmt.Errorf("%w: invalid UUID format", ErrInvalidOwnerID)
	}
	return nil
}

func ValidateHourlyRate(hourlyRate int64) error {
	if hourlyRate < MinHourlyRate {
		return fmt.Errorf("%w: cannot be negative", ErrInvalidHourlyRate)
	}
	if hourlyRate > MaxHourlyRate {
		return fmt.Errorf("%w: exceeds maximum allowed value of %d", ErrInvalidHourlyRate, MaxHourlyRate)
	}
	return nil
}

func ValidateCapacity(capacity int64) error {
	if capacity < MinCapacity {
		return fmt.Errorf("%w: must be at least %d", ErrInvalidCapacity, MinCapacity)
	}
	if capacity > MaxCapacity {
		return fmt.Errorf("%w: exceeds maximum allowed value of %d", ErrInvalidCapacity, MaxCapacity)
	}
	return nil
}

func ValidateString(str string, fieldName string) error {
	if str == "" {
		return fmt.Errorf("%w: %s cannot be empty", ErrInvalidStringLength, fieldName)
	}
	length := utf8.RuneCountInString(str)
	if length > MaxStringLength {
		return fmt.Errorf("%w: %s exceeds maximum length of %d", ErrInvalidStringLength, fieldName, MaxStringLength)
	}
	xssPattern := regexp.MustCompile(`(?i)(<script|javascript:|onerror=|onclick=)`)
	if xssPattern.MatchString(str) {
		return fmt.Errorf("%w: %s contains potentially dangerous content", ErrInvalidStringLength, fieldName)
	}
	return nil
}

func ValidateParkingType(parkingType string) error {
	validTypes := map[string]bool{
		"outdoor":     true,
		"covered":     true,
		"underground": true,
		"multi-level": true,
	}
	if !validTypes[parkingType] {
		return fmt.Errorf("%w: must be one of: outdoor, covered, underground, multi-level", ErrInvalidParkingType)
	}
	return nil
}

func SanitizeError(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("operation failed")
}
