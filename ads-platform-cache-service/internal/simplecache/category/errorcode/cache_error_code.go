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
	ErrSlugsEmpty     = New("CATEGORY_SLUGS_EMPTY", http.StatusBadRequest)
	ErrIDsEmpty       = New("CATEGORY_IDS_EMPTY", http.StatusBadRequest)
	ErrCDNUnavailable = New("CATEGORY_CDN_UNAVAILABLE", http.StatusBadGateway)
	ErrCacheEmpty     = New("CATEGORY_CACHE_EMPTY", http.StatusServiceUnavailable)
)
