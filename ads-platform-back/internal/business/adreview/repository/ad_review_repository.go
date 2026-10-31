package repository

import (
	"context"

	"ads-platform/internal/business/adreview/model"
)

type AdReviewRepository interface {
	List(ctx context.Context, filter model.ReviewFilter) ([]model.ReviewListItem, int64, error)
	Get(ctx context.Context, id int64) (*model.ReviewAd, error)
	// Approve and Reject change only ads that are still pending; they report whether a row changed,
	// so two reviewers acting at once cannot both decide the same ad.
	Approve(ctx context.Context, id, staffID int64) (bool, error)
	Reject(ctx context.Context, id, staffID int64, customerReason string, internalReason *string) (bool, error)
	Exists(ctx context.Context, id int64) (bool, error)
}
