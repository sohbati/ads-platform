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
	ErrUnknownKey    = New("SETTING_UNKNOWN_KEY", http.StatusNotFound)
	ErrInvalidValue  = New("SETTING_INVALID_VALUE", http.StatusBadRequest)
	ErrInvalidBody   = New("SETTING_INVALID_BODY", http.StatusBadRequest)
	ErrStaffRequired = New("SETTING_STAFF_REQUIRED", http.StatusUnauthorized)
)
