package utils

import (
	"errors"
	"fmt"
	"regexp"
	"unicode/utf8"
)

const (
	MaxAmount          = 1000000000000
	MinAmount          = 1
	MaxPromoCodeLength = 100
	MinPromoCodeLength = 4
)

var (
	ErrInvalidAmount    = errors.New("invalid amount")
	ErrAmountTooLarge   = errors.New("amount exceeds maximum allowed")
	ErrAmountTooSmall   = errors.New("amount must be positive")
	ErrInvalidUserID    = errors.New("invalid user ID")
	ErrInvalidPromoCode = errors.New("invalid promocode format")
	ErrEmptyUserID      = errors.New("user ID cannot be empty")
	ErrEmptyPromoCode   = errors.New("promocode cannot be empty")
)

func ValidateAmount(amount int64) error {
	if amount <= 0 {
		return ErrAmountTooSmall
	}
	if amount > MaxAmount {
		return ErrAmountTooLarge
	}
	return nil
}

func ValidateUserID(userID string) error {
	if userID == "" {
		return ErrEmptyUserID
	}
	if len(userID) > 128 {
		return fmt.Errorf("%w: user ID too long", ErrInvalidUserID)
	}
	uuidRegex := regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`)
	if !uuidRegex.MatchString(userID) {
		return fmt.Errorf("%w: invalid UUID format", ErrInvalidUserID)
	}
	return nil
}

func ValidatePromoCode(code string) error {
	if code == "" {
		return ErrEmptyPromoCode
	}
	codeLen := utf8.RuneCountInString(code)
	if codeLen < MinPromoCodeLength || codeLen > MaxPromoCodeLength {
		return fmt.Errorf("%w: length must be between %d and %d characters", ErrInvalidPromoCode, MinPromoCodeLength, MaxPromoCodeLength)
	}
	alphanumericRegex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	if !alphanumericRegex.MatchString(code) {
		return fmt.Errorf("%w: must contain only alphanumeric characters", ErrInvalidPromoCode)
	}
	return nil
}

func ValidateBookingID(bookingID int64) error {
	if bookingID <= 0 {
		return errors.New("booking ID must be positive")
	}
	return nil
}

func SafeAddBalance(current int64, amount int64) (int64, error) {
	if amount < 0 {
		return 0, ErrAmountTooSmall
	}
	if current > MaxAmount-amount {
		return 0, fmt.Errorf("balance would exceed maximum allowed value")
	}
	return current + amount, nil
}

func SafeSubtractBalance(current int64, amount int64) (int64, error) {
	if amount < 0 {
		return 0, ErrAmountTooSmall
	}
	if current < amount {
		return 0, errors.New("insufficient funds")
	}
	return current - amount, nil
}
