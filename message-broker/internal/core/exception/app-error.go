package exception

import (
	"errors"
	"fmt"
)

type AppError struct {
	StatusCode int
	ErrorCode  string
	Params     []string
	Cause      error
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.ErrorCode, e.Cause)
	}
	return e.ErrorCode
}

func (e *AppError) Unwrap() error {
	return e.Cause
}

func NewAppError(code string, statusCode int, params ...string) *AppError {
	return &AppError{
		StatusCode: statusCode,
		ErrorCode:  code,
		Params:     params,
	}
}

func (e *AppError) WithCause(err error) *AppError {
	e.Cause = err
	return e
}

func AsAppError(err error) (*AppError, bool) {
	var be *AppError
	if errors.As(err, &be) {
		return be, true
	}
	return nil, false
}
