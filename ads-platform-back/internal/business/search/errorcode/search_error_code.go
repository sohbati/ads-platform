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
	ErrInvalidPlace       = New("SEARCH_INVALID_PLACE", http.StatusBadRequest)
	ErrInvalidCategory    = New("SEARCH_INVALID_CATEGORY", http.StatusBadRequest)
	ErrInvalidCities      = New("SEARCH_INVALID_CITIES", http.StatusBadRequest)
	ErrCitiesRequired     = New("SEARCH_CITIES_REQUIRED", http.StatusBadRequest)
	ErrCatalogUnavailable = New("SEARCH_CATALOG_UNAVAILABLE", http.StatusServiceUnavailable)
)
