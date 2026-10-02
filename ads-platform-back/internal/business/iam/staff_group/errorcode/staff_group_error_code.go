package errorcode

import "net/http"

type ErrorCode struct {
	Code       string
	HttpStatus int
}

func New(code string, httpStatus int) *ErrorCode {
	return &ErrorCode{Code: code, HttpStatus: httpStatus}
}

func (e *ErrorCode) Error() string { return e.Code }

var (
	ErrInvalidCode        = New("STAFF_GROUP_INVALID_CODE", http.StatusBadRequest)
	ErrInvalidName        = New("STAFF_GROUP_INVALID_NAME", http.StatusBadRequest)
	ErrInvalidBody        = New("STAFF_GROUP_INVALID_BODY", http.StatusBadRequest)
	ErrInvalidID          = New("STAFF_GROUP_INVALID_ID", http.StatusBadRequest)
	ErrDuplicateCode      = New("STAFF_GROUP_DUPLICATE_CODE", http.StatusConflict)
	ErrNotFound           = New("STAFF_GROUP_NOT_FOUND", http.StatusNotFound)
	ErrMembershipNotFound = New("STAFF_GROUP_MEMBERSHIP_NOT_FOUND", http.StatusNotFound)
	ErrStaffNotFound       = New("STAFF_NOT_FOUND", http.StatusNotFound)
)
