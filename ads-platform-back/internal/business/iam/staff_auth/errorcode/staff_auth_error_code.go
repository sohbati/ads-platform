package errorcode

import "net/http"

type ErrorCode struct {
	Code       string
	HttpStatus int
}

func New(code string, httpStatus int) *ErrorCode {
	return &ErrorCode{Code: code, HttpStatus: httpStatus}
}

func (e *ErrorCode) Error() string { return e.Code }

var (
	ErrInvalidBody        = New("STAFF_AUTH_INVALID_BODY", http.StatusBadRequest)
	ErrInvalidCredentials = New("STAFF_AUTH_INVALID_CREDENTIALS", http.StatusUnauthorized)
	ErrSessionInvalid     = New("STAFF_AUTH_SESSION_INVALID", http.StatusUnauthorized)
)
