package service

import (
	"context"

	"ads-platform/internal/business/adreview/model"
)

// MaxReasonLength caps each rejection reason, in characters.
const MaxReasonLength = 1000

type RejectInput struct {
	// CustomerReason is shown to the ad owner and is required.
	CustomerReason string
	// InternalReason is visible only to staff and managers; empty means none.
	InternalReason string
}

type AdReviewService interface {
	List(ctx context.Context, filter model.ReviewFilter) (*model.ReviewPage, error)
	Get(ctx context.Context, id int64) (*model.ReviewAd, error)
	Approve(ctx context.Context, id, staffID int64) (*model.ReviewAd, error)
	Reject(ctx context.Context, id, staffID int64, input RejectInput) (*model.ReviewAd, error)
}
