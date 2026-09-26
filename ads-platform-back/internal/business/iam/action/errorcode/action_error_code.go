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
	ErrInvalidCode   = New("ACTION_INVALID_CODE", http.StatusBadRequest)
	ErrInvalidName   = New("ACTION_INVALID_NAME", http.StatusBadRequest)
	ErrInvalidBody   = New("ACTION_INVALID_BODY", http.StatusBadRequest)
	ErrInvalidID     = New("ACTION_INVALID_ID", http.StatusBadRequest)
	ErrDuplicateCode = New("ACTION_DUPLICATE_CODE", http.StatusConflict)
	ErrNotFound      = New("ACTION_NOT_FOUND", http.StatusNotFound)
)
