package impl

import (
	"context"
	"testing"

	"ads-platform/internal/business/adreview/model"
	"ads-platform/internal/business/adreview/service"
	adsmodel "ads-platform/internal/business/ads/model"
	"ads-platform/internal/core/exception"
)

type fakeAd struct {
	status         string
	customerReason *string
	internalReason *string
	reviewedBy     *int64
}

type fakeRepo struct {
	ads        map[int64]*fakeAd
	lastFilter model.ReviewFilter
}

func (f *fakeRepo) List(_ context.Context, filter model.ReviewFilter) ([]model.ReviewListItem, int64, error) {
	f.lastFilter = filter
	return []model.ReviewListItem{}, 0, nil
}

func (f *fakeRepo) Get(_ context.Context, id int64) (*model.ReviewAd, error) {
	ad, ok := f.ads[id]
	if !ok {
		return nil, nil
	}
	return &model.ReviewAd{
		ReviewListItem:          model.ReviewListItem{ID: id, Status: ad.status},
		RejectionReasonCustomer: ad.customerReason,
		RejectionReasonInternal: ad.internalReason,
		ReviewedByStaffID:       ad.reviewedBy,
	}, nil
}

func (f *fakeRepo) Approve(_ context.Context, id, staffID int64) (bool, error) {
	ad, ok := f.ads[id]
	if !ok || ad.status != adsmodel.AdStatusPending {
		return false, nil
	}
	ad.status, ad.customerReason, ad.internalReason, ad.reviewedBy = adsmodel.AdStatusActive, nil, nil, &staffID
	return true, nil
}

func (f *fakeRepo) Reject(_ context.Context, id, staffID int64, customer string, internal *string) (bool, error) {
	ad, ok := f.ads[id]
	if !ok || ad.status != adsmodel.AdStatusPending {
		return false, nil
	}
	ad.status, ad.customerReason, ad.internalReason, ad.reviewedBy = adsmodel.AdStatusDraft, &customer, internal, &staffID
	return true, nil
}

func (f *fakeRepo) Exists(_ context.Context, id int64) (bool, error) {
	_, ok := f.ads[id]
	return ok, nil
}

func newRepo() *fakeRepo {
	return &fakeRepo{ads: map[int64]*fakeAd{
		1: {status: adsmodel.AdStatusPending},
		2: {status: adsmodel.AdStatusActive},
	}}
}

func code(t *testing.T, err error) string {
	t.Helper()
	app, ok := exception.AsAppError(err)
	if !ok {
		t.Fatalf("want AppError, got %v", err)
	}
	return app.ErrorCode
}

func TestApprovePendingAd(t *testing.T) {
	ad, err := NewAdReviewService(newRepo(), nil).Approve(context.Background(), 1, 5)
	if err != nil {
		t.Fatal(err)
	}
	if ad.Status != adsmodel.AdStatusActive || ad.ReviewedByStaffID == nil || *ad.ReviewedByStaffID != 5 {
		t.Fatalf("got %+v", ad)
	}
}

func TestRejectMovesAdToDraftWithBothReasons(t *testing.T) {
	ad, err := NewAdReviewService(newRepo(), nil).Reject(context.Background(), 1, 5, service.RejectInput{
		CustomerReason: "  عکس‌ها مربوط به کالا نیستند  ",
		InternalReason: "possible duplicate of #12",
	})
	if err != nil {
		t.Fatal(err)
	}
	if ad.Status != adsmodel.AdStatusDraft {
		t.Fatalf("status %s", ad.Status)
	}
	if *ad.RejectionReasonCustomer != "عکس‌ها مربوط به کالا نیستند" || *ad.RejectionReasonInternal != "possible duplicate of #12" {
		t.Fatalf("reasons %q / %q", *ad.RejectionReasonCustomer, *ad.RejectionReasonInternal)
	}
}

func TestRejectWithoutInternalReasonStoresNull(t *testing.T) {
	ad, err := NewAdReviewService(newRepo(), nil).Reject(context.Background(), 1, 5, service.RejectInput{CustomerReason: "Price missing"})
	if err != nil {
		t.Fatal(err)
	}
	if ad.RejectionReasonInternal != nil {
		t.Fatalf("internal reason %q", *ad.RejectionReasonInternal)
	}
}

func TestRejectValidation(t *testing.T) {
	svc := NewAdReviewService(newRepo(), nil)
	long := make([]rune, service.MaxReasonLength+1)
	for i := range long {
		long[i] = 'x'
	}
	cases := map[string]service.RejectInput{
		"AD_REVIEW_REASON_REQUIRED": {CustomerReason: "   "},
		"AD_REVIEW_REASON_TOO_LONG": {CustomerReason: string(long)},
	}
	for want, input := range cases {
		_, err := svc.Reject(context.Background(), 1, 5, input)
		if got := code(t, err); got != want {
			t.Errorf("want %s, got %s", want, got)
		}
	}
}

func TestDecisionsOnlyApplyToPendingAds(t *testing.T) {
	svc := NewAdReviewService(newRepo(), nil)
	if _, err := svc.Approve(context.Background(), 2, 5); code(t, err) != "AD_REVIEW_NOT_PENDING" {
		t.Fatalf("approve active: %v", err)
	}
	if _, err := svc.Reject(context.Background(), 2, 5, service.RejectInput{CustomerReason: "x"}); code(t, err) != "AD_REVIEW_NOT_PENDING" {
		t.Fatalf("reject active: %v", err)
	}
	if _, err := svc.Approve(context.Background(), 99, 5); code(t, err) != "AD_REVIEW_NOT_FOUND" {
		t.Fatalf("approve missing: %v", err)
	}
}

func TestDecisionsRequireStaff(t *testing.T) {
	if _, err := NewAdReviewService(newRepo(), nil).Approve(context.Background(), 1, 0); code(t, err) != "AD_REVIEW_STAFF_REQUIRED" {
		t.Fatalf("got %v", err)
	}
}

func TestListValidatesAndClampsFilter(t *testing.T) {
	repo := newRepo()
	svc := NewAdReviewService(repo, nil)
	if _, err := svc.List(context.Background(), model.ReviewFilter{Status: "bogus"}); code(t, err) != "AD_REVIEW_INVALID_STATUS" {
		t.Fatalf("got %v", err)
	}
	if _, err := svc.List(context.Background(), model.ReviewFilter{Status: " Pending ", PageSize: 5000}); err != nil {
		t.Fatal(err)
	}
	if repo.lastFilter.Status != "pending" || repo.lastFilter.Page != 1 || repo.lastFilter.PageSize != maxPageSize {
		t.Fatalf("filter %+v", repo.lastFilter)
	}
}
