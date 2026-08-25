package errorcode

import "net/http"

type ErrorCode struct {
	Code       string
	HttpStatus int
}

func New(code string, httpStatus int) *ErrorCode {
	return &ErrorCode{Code: code, HttpStatus: httpStatus}
}

func (e *ErrorCode) Error() string {
	return e.Code
}

var (
	ErrNamesEmpty     = New("ATTR_SCHEMA_NAMES_EMPTY", http.StatusBadRequest)
	ErrCDNUnavailable = New("ATTR_SCHEMA_CDN_UNAVAILABLE", http.StatusBadGateway)
	ErrCacheEmpty     = New("ATTR_SCHEMA_CACHE_EMPTY", http.StatusServiceUnavailable)
)
