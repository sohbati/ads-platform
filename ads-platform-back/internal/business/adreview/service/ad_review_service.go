package service

import (
	"context"
	"time"

	"ads-platform/internal/business/adreview/model"
)

// MaxReasonLength caps each rejection reason, in characters.
const MaxReasonLength = 1000

// MaxCountedTotal caps the match count in list responses; beyond it the panel shows "1000+".
const MaxCountedTotal = 1000

const (
	// ClaimLease is how long a claim lasts without being extended; an abandoned ad returns to the
	// queue after this. The panel extends it every minute while the ad stays open.
	ClaimLease = 10 * time.Minute
	// MaxClaimHold caps one claim, so an ad left open in an idle tab still goes back to the queue.
	// Claiming the ad again starts a new hold.
	MaxClaimHold = 30 * time.Minute
)

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

	// Next returns the ad staffID is reviewing, or claims the oldest waiting one.
	// It returns nil when no ad is waiting.
	Next(ctx context.Context, staffID int64) (*model.ReviewAd, error)
	// Claim takes a specific pending ad for review, giving up the reviewer's current one.
	Claim(ctx context.Context, id, staffID int64) (*model.ReviewAd, error)
	// Extend keeps the reviewer's claim alive while they work, up to MaxClaimHold.
	Extend(ctx context.Context, id, staffID int64) (*model.ReviewAd, error)
	// Release hands the ad back to the queue; releasing an ad you do not hold does nothing.
	Release(ctx context.Context, id, staffID int64) error
	// ReleaseAll hands back every ad staffID holds, e.g. when they sign out.
	ReleaseAll(ctx context.Context, staffID int64) error
	// ForceRelease frees an ad whoever holds it; for managers clearing a stuck claim.
	ForceRelease(ctx context.Context, id int64) error
	// TodayStats counts each reviewer's decisions since midnight, Tehran time.
	TodayStats(ctx context.Context) ([]model.ReviewerStats, error)
}
