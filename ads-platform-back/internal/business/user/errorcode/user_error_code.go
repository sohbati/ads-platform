package errorcode

import (
	"net/http"
)

type ErrorCode struct {
	Code       string
	HttpStatus int
}

var (
	ErrUserNotFound    = New("USER_NOT_FOUND", http.StatusNotFound)
	ErrMobileEmpty     = New("MOBILE_EMPTY", http.StatusBadRequest)
	ErrDuplicateMobile = New("DUPLICATE_MOBILE", http.StatusConflict)
)

func New(code string, httpStatus int) *ErrorCode {
	return &ErrorCode{
		Code:       code,
		HttpStatus: httpStatus,
	}
}

func (e *ErrorCode) Error() string {
	return e.Code
}
