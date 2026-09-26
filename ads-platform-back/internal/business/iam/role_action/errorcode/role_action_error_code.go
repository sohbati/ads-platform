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
	ErrInvalidID      = New("ROLE_ACTION_INVALID_ID", http.StatusBadRequest)
	ErrInvalidBody    = New("ROLE_ACTION_INVALID_BODY", http.StatusBadRequest)
	ErrNotFound       = New("ROLE_ACTION_NOT_FOUND", http.StatusNotFound)
	ErrRoleNotFound   = New("ROLE_NOT_FOUND", http.StatusNotFound)
	ErrActionNotFound = New("ACTION_NOT_FOUND", http.StatusNotFound)
)
