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
	ErrInvalidCode        = New("USER_GROUP_INVALID_CODE", http.StatusBadRequest)
	ErrInvalidName        = New("USER_GROUP_INVALID_NAME", http.StatusBadRequest)
	ErrInvalidBody        = New("USER_GROUP_INVALID_BODY", http.StatusBadRequest)
	ErrInvalidID          = New("USER_GROUP_INVALID_ID", http.StatusBadRequest)
	ErrDuplicateCode      = New("USER_GROUP_DUPLICATE_CODE", http.StatusConflict)
	ErrNotFound           = New("USER_GROUP_NOT_FOUND", http.StatusNotFound)
	ErrMembershipNotFound = New("USER_GROUP_MEMBERSHIP_NOT_FOUND", http.StatusNotFound)
	ErrUserNotFound       = New("USER_NOT_FOUND", http.StatusNotFound)
)
