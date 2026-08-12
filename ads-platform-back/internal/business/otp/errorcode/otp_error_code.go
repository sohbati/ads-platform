package errorcode

import (
	"net/http"
)

type ErrorCode struct {
	Code       string
	HttpStatus int
}

var (
	ErrMobileEmpty    = New("MOBILE_EMPTY", http.StatusBadRequest)
	ErrInvalidMobile  = New("INVALID_MOBILE", http.StatusBadRequest)
	ErrInvalidOTP     = New("INVALID_OTP", http.StatusBadRequest)
	ErrOTPNotFound    = New("OTP_NOT_FOUND", http.StatusNotFound)
	ErrOTPExpired     = New("OTP_EXPIRED", http.StatusNotFound)
	ErrOTPVerifyFailed = New("OTP_VERIFY_FAILED", http.StatusUnauthorized)
	ErrCacheUnavailable = New("CACHE_SERVICE_UNAVAILABLE", http.StatusServiceUnavailable)
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
