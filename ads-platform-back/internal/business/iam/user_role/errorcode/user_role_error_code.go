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
	ErrInvalidID    = New("USER_ROLE_INVALID_ID", http.StatusBadRequest)
	ErrInvalidCode  = New("USER_ROLE_INVALID_CODE", http.StatusBadRequest)
	ErrNotFound     = New("USER_ROLE_NOT_FOUND", http.StatusNotFound)
	ErrUserNotFound = New("USER_NOT_FOUND", http.StatusNotFound)
	ErrRoleNotFound = New("ROLE_NOT_FOUND", http.StatusNotFound)
	ErrActionDenied = New("ACTION_DENIED", http.StatusForbidden)
)
