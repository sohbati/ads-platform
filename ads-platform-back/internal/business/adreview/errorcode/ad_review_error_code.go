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
	ErrInvalidID      = New("AD_REVIEW_INVALID_ID", http.StatusBadRequest)
	ErrInvalidBody    = New("AD_REVIEW_INVALID_BODY", http.StatusBadRequest)
	ErrInvalidStatus  = New("AD_REVIEW_INVALID_STATUS", http.StatusBadRequest)
	ErrReasonRequired = New("AD_REVIEW_REASON_REQUIRED", http.StatusBadRequest)
	ErrReasonTooLong  = New("AD_REVIEW_REASON_TOO_LONG", http.StatusBadRequest)
	ErrStaffRequired  = New("AD_REVIEW_STAFF_REQUIRED", http.StatusUnauthorized)
	ErrNotFound       = New("AD_REVIEW_NOT_FOUND", http.StatusNotFound)
	ErrNotPending     = New("AD_REVIEW_NOT_PENDING", http.StatusConflict)
)
