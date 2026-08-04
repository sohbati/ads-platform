package errorcode

import (
	"net/http"
)

type ErrorCode struct {
	Code       string
	HttpStatus int
}

var (
	ErrCacheNotFound  = New("CACHE_ITEM_NOT_FOUND", http.StatusNotFound)
	ErrCacheMiss      = New("key_not_found", http.StatusNotFound)
	ErrKeyExpired     = New("key-has-expired", http.StatusNotFound)
	ErrKeyEmpty       = New("key_cannot_be_empty", http.StatusNotFound)
	ErrFailToReadBody = New("FAIL_TO_READ_BODY", http.StatusBadRequest)
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
