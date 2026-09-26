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
	ErrInvalidCode   = New("ROLE_INVALID_CODE", http.StatusBadRequest)
	ErrInvalidName   = New("ROLE_INVALID_NAME", http.StatusBadRequest)
	ErrInvalidBody   = New("ROLE_INVALID_BODY", http.StatusBadRequest)
	ErrInvalidID     = New("ROLE_INVALID_ID", http.StatusBadRequest)
	ErrDuplicateCode = New("ROLE_DUPLICATE_CODE", http.StatusConflict)
	ErrNotFound      = New("ROLE_NOT_FOUND", http.StatusNotFound)
)
