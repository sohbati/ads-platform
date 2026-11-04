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
	ErrInvalidID     = New("AD_REVIEW_INVALID_ID", http.StatusBadRequest)
	ErrInvalidBody   = New("AD_REVIEW_INVALID_BODY", http.StatusBadRequest)
	ErrInvalidStatus = New("AD_REVIEW_INVALID_STATUS", http.StatusBadRequest)
	ErrInvalidCursor = New("AD_REVIEW_INVALID_CURSOR", http.StatusBadRequest)
	// ErrFilterRequired: lists other than the pending queue need a search or a filter.
	ErrFilterRequired = New("AD_REVIEW_FILTER_REQUIRED", http.StatusBadRequest)
	// ErrSearchTooShort: a title search needs at least three letters.
	ErrSearchTooShort = New("AD_REVIEW_SEARCH_TOO_SHORT", http.StatusBadRequest)
	// ErrCatalogUnavailable: the category filter could not be expanded to its subcategories.
	ErrCatalogUnavailable = New("AD_REVIEW_CATALOG_UNAVAILABLE", http.StatusBadGateway)
	ErrReasonRequired     = New("AD_REVIEW_REASON_REQUIRED", http.StatusBadRequest)
	ErrReasonTooLong      = New("AD_REVIEW_REASON_TOO_LONG", http.StatusBadRequest)
	ErrStaffRequired      = New("AD_REVIEW_STAFF_REQUIRED", http.StatusUnauthorized)
	ErrNotFound           = New("AD_REVIEW_NOT_FOUND", http.StatusNotFound)
	ErrNotPending         = New("AD_REVIEW_NOT_PENDING", http.StatusConflict)
	// ErrNotClaimed: the reviewer must take the ad for review before deciding it.
	ErrNotClaimed = New("AD_REVIEW_NOT_CLAIMED", http.StatusConflict)
	// ErrClaimExpired: the reviewer held the ad, but too long passed and it went back to the queue.
	ErrClaimExpired = New("AD_REVIEW_CLAIM_EXPIRED", http.StatusConflict)
	// ErrClaimedByOther: another reviewer is working on the ad.
	ErrClaimedByOther = New("AD_REVIEW_CLAIMED_BY_OTHER", http.StatusConflict)
)
