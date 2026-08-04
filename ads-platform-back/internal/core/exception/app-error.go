package exception

import (
	"errors"
	"fmt"
)

type AppError struct {
	StatusCode int
	ErrorCode string
	Params    []string
	Cause     error
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.ErrorCode, e.Cause)
	}
	return e.ErrorCode
}

// Allow errors.Unwrap / errors.Is to work
func (e *AppError) Unwrap() error {
	return e.Cause
}

// Constructor
func NewAppError(code string, statusCode int, params ...string) *AppError {
	return &AppError{
		StatusCode: statusCode,
		ErrorCode: code,
		Params:    params,
	}
}

// Attach cause
func (e *AppError) WithCause(err error) *AppError {
	e.Cause = err
	return e
}

// Helper to check if an error is AppError
func AsAppError(err error) (*AppError, bool) {
	var be *AppError
	if errors.As(err, &be) {
		return be, true
	}
	return nil, false
}
