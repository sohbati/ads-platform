package cachestore

import (
	"net/http"
)

type ErrorCode struct {
	Code       string
	HttpStatus int
}

var (
	ErrCacheMiss  = NewErrorCode("key_not_found", http.StatusNotFound)
	ErrKeyExpired = NewErrorCode("key-has-expired", http.StatusNotFound)
	ErrKeyEmpty   = NewErrorCode("key_cannot_be_empty", http.StatusNotFound)
)

func NewErrorCode(code string, httpStatus int) *ErrorCode {
	return &ErrorCode{
		Code:       code,
		HttpStatus: httpStatus,
	}
}

func (e *ErrorCode) Error() string {
	return e.Code
}
