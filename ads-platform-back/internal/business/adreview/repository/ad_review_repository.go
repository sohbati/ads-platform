package repository

import (
	"context"
	"time"

	"ads-platform/internal/business/adreview/model"
)

type AdReviewRepository interface {
	// List returns up to q.Limit ads and the number of matches, counted no further than countLimit.
	List(ctx context.Context, q model.ReviewQuery, countLimit int) ([]model.ReviewListItem, int64, error)
	Get(ctx context.Context, id int64) (*model.ReviewAd, error)
	// Approve and Reject change only ads that are still pending and claimed by staffID until after
	// now; they report whether a row changed, so two reviewers cannot both decide the same ad.
	Approve(ctx context.Context, id, staffID int64, now time.Time) (bool, error)
	Reject(ctx context.Context, id, staffID int64, customerReason string, internalReason *string, now time.Time) (bool, error)
	Exists(ctx context.Context, id int64) (bool, error)

	// ClaimNext returns the ad staffID already holds, or claims the oldest unclaimed pending ad
	// until expiresAt. It returns 0 when the queue is empty.
	ClaimNext(ctx context.Context, staffID int64, now, expiresAt time.Time) (int64, error)
	// Claim takes one pending ad that nobody else holds, releasing staffID's other claims.
	Claim(ctx context.Context, id, staffID int64, now, expiresAt time.Time) (bool, error)
	// Extend moves the expiry of a live claim staffID holds; it reports whether one was extended.
	Extend(ctx context.Context, id, staffID int64, now, expiresAt time.Time) (bool, error)
	// Release drops staffID's claim on the ad, if any.
	Release(ctx context.Context, id, staffID int64) error
	// ReleaseAll drops every claim staffID holds.
	ReleaseAll(ctx context.Context, staffID int64) error
	// ForceRelease drops the claim on the ad, whoever holds it.
	ForceRelease(ctx context.Context, id int64) error
	// ReviewerStats counts decisions made since the given time, per reviewer.
	ReviewerStats(ctx context.Context, since time.Time) ([]model.ReviewerStats, error)
}
