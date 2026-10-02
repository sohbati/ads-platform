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
	ErrInvalidCustomer      = New("PROFILE_INVALID_CUSTOMER", http.StatusBadRequest)
	ErrInvalidLocation      = New("PROFILE_INVALID_LOCATION", http.StatusBadRequest)
	ErrTooManyLocations     = New("PROFILE_TOO_MANY_LOCATIONS", http.StatusBadRequest)
	ErrCustomerNotFound         = New("CUSTOMER_NOT_FOUND", http.StatusNotFound)
)
