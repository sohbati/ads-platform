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
	ErrImageNotFound         = New("IMAGE_NOT_FOUND", http.StatusNotFound)
	ErrImageInvalidUser      = New("IMAGE_INVALID_USER", http.StatusBadRequest)
	ErrImageInvalidFilename  = New("IMAGE_INVALID_FILENAME", http.StatusBadRequest)
	ErrImageInvalidType      = New("IMAGE_INVALID_CONTENT_TYPE", http.StatusBadRequest)
	ErrImageInvalidSize      = New("IMAGE_INVALID_SIZE", http.StatusBadRequest)
	ErrImageTooLarge         = New("IMAGE_TOO_LARGE", http.StatusRequestEntityTooLarge)
	ErrImageInvalidStatus    = New("IMAGE_INVALID_STATUS", http.StatusConflict)
	ErrImageChecksumRequired = New("IMAGE_CHECKSUM_REQUIRED", http.StatusBadRequest)
)
