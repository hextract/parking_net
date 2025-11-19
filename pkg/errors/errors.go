package errors

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func Wrap(err error, code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

var (
	ErrNotFound            = New(http.StatusNotFound, "resource not found")
	ErrBadRequest          = New(http.StatusBadRequest, "bad request")
	ErrUnauthorized        = New(http.StatusUnauthorized, "unauthorized")
	ErrForbidden           = New(http.StatusForbidden, "forbidden")
	ErrInternalServer      = New(http.StatusInternalServerError, "internal server error")
	ErrConflict            = New(http.StatusConflict, "resource conflict")
	ErrValidation          = New(http.StatusBadRequest, "validation error")
)

func NotFound(resource string) *AppError {
	return New(http.StatusNotFound, fmt.Sprintf("%s not found", resource))
}

func BadRequest(message string) *AppError {
	return New(http.StatusBadRequest, message)
}

func Validation(message string) *AppError {
	return New(http.StatusBadRequest, fmt.Sprintf("validation error: %s", message))
}

func Internal(err error) *AppError {
	return Wrap(err, http.StatusInternalServerError, "internal server error")
}

