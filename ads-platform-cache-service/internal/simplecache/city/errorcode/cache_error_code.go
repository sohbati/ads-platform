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
	ErrFailToReadBody   = New("FAIL_TO_READ_BODY", http.StatusBadRequest)
	ErrIDsEmpty         = New("CITY_IDS_EMPTY", http.StatusBadRequest)
	ErrSlugsEmpty       = New("CITY_SLUGS_EMPTY", http.StatusBadRequest)
	ErrCDNUnavailable   = New("CITY_CDN_UNAVAILABLE", http.StatusBadGateway)
	ErrCityCacheEmpty   = New("CITY_CACHE_EMPTY", http.StatusServiceUnavailable)
)
