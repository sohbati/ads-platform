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
	ErrInvalidName         = New("STAFF_INVALID_NAME", http.StatusBadRequest)
	ErrInvalidMobile       = New("STAFF_INVALID_MOBILE", http.StatusBadRequest)
	ErrInvalidNationalId   = New("STAFF_INVALID_NATIONAL_ID", http.StatusBadRequest)
	ErrInvalidBody         = New("STAFF_INVALID_BODY", http.StatusBadRequest)
	ErrInvalidID           = New("STAFF_INVALID_ID", http.StatusBadRequest)
	ErrDuplicateMobile     = New("STAFF_DUPLICATE_MOBILE", http.StatusConflict)
	ErrDuplicateNationalId = New("STAFF_DUPLICATE_NATIONAL_ID", http.StatusConflict)
	ErrNotFound            = New("STAFF_NOT_FOUND", http.StatusNotFound)
	ErrInvalidPassword     = New("STAFF_INVALID_PASSWORD", http.StatusBadRequest)
	ErrSuperAdminProtected = New("STAFF_SUPER_ADMIN_PROTECTED", http.StatusConflict)
)
