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
	ErrAdInvalidUser        = New("AD_INVALID_USER", http.StatusBadRequest)
	ErrAdInvalidTitle       = New("AD_INVALID_TITLE", http.StatusBadRequest)
	ErrAdInvalidDescription = New("AD_INVALID_DESCRIPTION", http.StatusBadRequest)
	ErrAdInvalidCategory    = New("AD_INVALID_CATEGORY", http.StatusBadRequest)
	ErrAdCategoryNotLeaf    = New("AD_CATEGORY_NOT_LEAF", http.StatusBadRequest)
	ErrAdInvalidCity        = New("AD_INVALID_CITY", http.StatusBadRequest)
	ErrAdInvalidLocation    = New("AD_INVALID_LOCATION", http.StatusBadRequest)
	ErrAdInvalidPrice       = New("AD_INVALID_PRICE", http.StatusBadRequest)
	ErrAdInvalidAttrs       = New("AD_INVALID_ATTRS", http.StatusBadRequest)
	ErrAdTooManyPictures    = New("AD_TOO_MANY_PICTURES", http.StatusBadRequest)
	ErrAdInvalidPicture     = New("AD_INVALID_PICTURE", http.StatusBadRequest)
	ErrAdPictureTooLarge    = New("AD_PICTURE_TOO_LARGE", http.StatusRequestEntityTooLarge)
	ErrAdCatalogUnavailable = New("AD_CATALOG_UNAVAILABLE", http.StatusServiceUnavailable)
	ErrAdStorageUnavailable = New("AD_STORAGE_UNAVAILABLE", http.StatusServiceUnavailable)
	ErrAdNotFound           = New("AD_NOT_FOUND", http.StatusNotFound)
)
