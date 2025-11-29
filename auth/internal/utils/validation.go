package utils

import (
	"errors"
	"fmt"
	"regexp"
	"unicode"
	"unicode/utf8"
)

const (
	MaxLoginLength    = 50
	MinLoginLength    = 3
	MaxPasswordLength = 128
	MinPasswordLength = 8
	MaxEmailLength    = 254
	MinEmailLength    = 5
	MaxTelegramID     = 9223372036854775807
	MinTelegramID     = 1
)

var (
	ErrInvalidLogin      = errors.New("invalid login")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrInvalidEmail      = errors.New("invalid email")
	ErrInvalidTelegramID = errors.New("invalid telegram ID")
	ErrInvalidRole       = errors.New("invalid role")
)

func ValidateLogin(login string) error {
	if login == "" {
		return fmt.Errorf("%w: cannot be empty", ErrInvalidLogin)
	}
	length := utf8.RuneCountInString(login)
	if length < MinLoginLength || length > MaxLoginLength {
		return fmt.Errorf("%w: must be between %d and %d characters", ErrInvalidLogin, MinLoginLength, MaxLoginLength)
	}
	alphanumericRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !alphanumericRegex.MatchString(login) {
		return fmt.Errorf("%w: must contain only letters, numbers, and underscores", ErrInvalidLogin)
	}
	return nil
}

func ValidatePassword(password string) error {
	if password == "" {
		return fmt.Errorf("%w: cannot be empty", ErrInvalidPassword)
	}
	length := utf8.RuneCountInString(password)
	if length < MinPasswordLength || length > MaxPasswordLength {
		return fmt.Errorf("%w: must be between %d and %d characters", ErrInvalidPassword, MinPasswordLength, MaxPasswordLength)
	}

	var hasUpper, hasLower, hasNumber bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber {
		return fmt.Errorf("%w: must contain at least one uppercase letter, one lowercase letter, and one number", ErrInvalidPassword)
	}

	return nil
}

func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("%w: cannot be empty", ErrInvalidEmail)
	}
	length := utf8.RuneCountInString(email)
	if length < MinEmailLength || length > MaxEmailLength {
		return fmt.Errorf("%w: must be between %d and %d characters", ErrInvalidEmail, MinEmailLength, MaxEmailLength)
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("%w: invalid email format", ErrInvalidEmail)
	}
	return nil
}

func ValidateTelegramID(telegramID int64) error {
	if telegramID < MinTelegramID {
		return fmt.Errorf("%w: must be positive", ErrInvalidTelegramID)
	}
	if telegramID > MaxTelegramID {
		return fmt.Errorf("%w: exceeds maximum value", ErrInvalidTelegramID)
	}
	return nil
}

func ValidateRole(role string) error {
	validRoles := map[string]bool{
		"driver": true,
		"owner":  true,
	}
	if !validRoles[role] {
		return fmt.Errorf("%w: must be 'driver' or 'owner'", ErrInvalidRole)
	}
	return nil
}

func SanitizeError(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("operation failed")
}
