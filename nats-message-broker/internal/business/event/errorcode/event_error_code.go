package errorcode

import (
	"net/http"
)

type ErrorCode struct {
	Code       string
	HttpStatus int
}

var (
	ErrSubjectEmpty     = New("SUBJECT_EMPTY", http.StatusBadRequest)
	ErrPayloadEmpty     = New("PAYLOAD_EMPTY", http.StatusBadRequest)
	ErrPublishFailed    = New("PUBLISH_FAILED", http.StatusInternalServerError)
	ErrNatsUnavailable  = New("NATS_UNAVAILABLE", http.StatusServiceUnavailable)
	ErrFailToReadBody   = New("FAIL_TO_READ_BODY", http.StatusBadRequest)
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
