package impl

import (
	"context"
	"testing"
	"time"

	"ads-platform/internal/business/adreview/model"
	"ads-platform/internal/business/adreview/service"
	adsmodel "ads-platform/internal/business/ads/model"
	searchclient "ads-platform/internal/business/search/client"
	"ads-platform/internal/core/exception"
)

type fakeAd struct {
	status         string
	customerReason *string
	internalReason *string
	reviewedBy     *int64
	claimedBy      *int64
	claimedAt      *time.Time
	expiresAt      *time.Time
}

func (ad *fakeAd) heldBy(staffID int64, now time.Time) bool {
	return ad.status == adsmodel.AdStatusPending && ad.claimedBy != nil && *ad.claimedBy == staffID &&
		ad.expiresAt.After(now)
}

func (ad *fakeAd) free(now time.Time) bool {
	return ad.claimedBy == nil || !ad.expiresAt.After(now)
}

func (ad *fakeAd) clearClaim() {
	ad.claimedBy, ad.claimedAt, ad.expiresAt = nil, nil, nil
}

type fakeRepo struct {
	ads        map[int64]*fakeAd
	lastQuery  model.ReviewQuery
	listItems  []model.ReviewListItem
	listTotal  int64
	countLimit int
}

func (f *fakeRepo) List(_ context.Context, q model.ReviewQuery, countLimit int) ([]model.ReviewListItem, int64, error) {
	f.lastQuery, f.countLimit = q, countLimit
	items := f.listItems
	if len(items) > q.Limit {
		items = items[:q.Limit]
	}
	return items, f.listTotal, nil
}

func (f *fakeRepo) Get(_ context.Context, id int64) (*model.ReviewAd, error) {
	ad, ok := f.ads[id]
	if !ok {
		return nil, nil
	}
	return &model.ReviewAd{
		ReviewListItem: model.ReviewListItem{
			ID: id, Status: ad.status,
			ClaimedByStaffID: ad.claimedBy, ClaimedAt: ad.claimedAt, ClaimExpiresAt: ad.expiresAt,
		},
		RejectionReasonCustomer: ad.customerReason,
		RejectionReasonInternal: ad.internalReason,
		ReviewedByStaffID:       ad.reviewedBy,
	}, nil
}

func (f *fakeRepo) Approve(_ context.Context, id, staffID int64, now time.Time) (bool, error) {
	ad, ok := f.ads[id]
	if !ok || !ad.heldBy(staffID, now) {
		return false, nil
	}
	ad.status, ad.customerReason, ad.internalReason, ad.reviewedBy = adsmodel.AdStatusActive, nil, nil, &staffID
	ad.clearClaim()
	return true, nil
}

func (f *fakeRepo) Reject(_ context.Context, id, staffID int64, customer string, internal *string, now time.Time) (bool, error) {
	ad, ok := f.ads[id]
	if !ok || !ad.heldBy(staffID, now) {
		return false, nil
	}
	ad.status, ad.customerReason, ad.internalReason, ad.reviewedBy = adsmodel.AdStatusDraft, &customer, internal, &staffID
	ad.clearClaim()
	return true, nil
}

// ClaimNext takes the lowest pending id, standing in for the oldest pending_since.
func (f *fakeRepo) ClaimNext(_ context.Context, staffID int64, now, expiresAt time.Time) (int64, error) {
	var next int64
	for id, ad := range f.ads {
		if ad.heldBy(staffID, now) {
			return id, nil
		}
		if ad.status == adsmodel.AdStatusPending && ad.free(now) && (next == 0 || id < next) {
			next = id
		}
	}
	if next != 0 {
		f.ads[next].claimedBy, f.ads[next].claimedAt, f.ads[next].expiresAt = &staffID, &now, &expiresAt
	}
	return next, nil
}

func (f *fakeRepo) Claim(_ context.Context, id, staffID int64, now, expiresAt time.Time) (bool, error) {
	ad, ok := f.ads[id]
	if !ok || ad.status != adsmodel.AdStatusPending || !(ad.free(now) || ad.heldBy(staffID, now)) {
		return false, nil
	}
	_ = f.ReleaseAll(context.Background(), staffID)
	ad.claimedBy, ad.claimedAt, ad.expiresAt = &staffID, &now, &expiresAt
	return true, nil
}

func (f *fakeRepo) Extend(_ context.Context, id, staffID int64, now, expiresAt time.Time) (bool, error) {
	ad, ok := f.ads[id]
	if !ok || !ad.heldBy(staffID, now) {
		return false, nil
	}
	ad.expiresAt = &expiresAt
	return true, nil
}

func (f *fakeRepo) Release(_ context.Context, id, staffID int64) error {
	if ad, ok := f.ads[id]; ok && ad.claimedBy != nil && *ad.claimedBy == staffID {
		ad.clearClaim()
	}
	return nil
}

func (f *fakeRepo) ReleaseAll(_ context.Context, staffID int64) error {
	for _, ad := range f.ads {
		if ad.claimedBy != nil && *ad.claimedBy == staffID {
			ad.clearClaim()
		}
	}
	return nil
}

func (f *fakeRepo) ForceRelease(_ context.Context, id int64) error {
	if ad, ok := f.ads[id]; ok {
		ad.clearClaim()
	}
	return nil
}

func (f *fakeRepo) ReviewerStats(_ context.Context, _ time.Time) ([]model.ReviewerStats, error) {
	return []model.ReviewerStats{}, nil
}

func (f *fakeRepo) Exists(_ context.Context, id int64) (bool, error) {
	_, ok := f.ads[id]
	return ok, nil
}

func newRepo() *fakeRepo {
	return &fakeRepo{ads: map[int64]*fakeAd{
		1: {status: adsmodel.AdStatusPending},
		2: {status: adsmodel.AdStatusActive},
		3: {status: adsmodel.AdStatusPending},
	}}
}

// clock is a settable time source for the service.
type clock struct{ at time.Time }

func (c *clock) now() time.Time { return c.at }

func (c *clock) advance(d time.Duration) { c.at = c.at.Add(d) }

func newService(repo *fakeRepo) (*adReviewService, *clock) {
	c := &clock{at: time.Date(2026, 11, 4, 9, 0, 0, 0, time.UTC)}
	svc := NewAdReviewService(repo, nil).(*adReviewService)
	svc.now = c.now
	return svc, c
}

// claimed returns a service where staff 5 holds ad 1.
func claimed(t *testing.T) (*adReviewService, *fakeRepo, *clock) {
	t.Helper()
	repo := newRepo()
	svc, c := newService(repo)
	if _, err := svc.Claim(context.Background(), 1, 5); err != nil {
		t.Fatal(err)
	}
	return svc, repo, c
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
	svc, _, _ := claimed(t)
	ad, err := svc.Approve(context.Background(), 1, 5)
	if err != nil {
		t.Fatal(err)
	}
	if ad.Status != adsmodel.AdStatusActive || ad.ReviewedByStaffID == nil || *ad.ReviewedByStaffID != 5 {
		t.Fatalf("got %+v", ad)
	}
	if ad.ClaimedByStaffID != nil {
		t.Fatalf("claim kept after decision: %+v", ad)
	}
}

func TestRejectMovesAdToDraftWithBothReasons(t *testing.T) {
	svc, _, _ := claimed(t)
	ad, err := svc.Reject(context.Background(), 1, 5, service.RejectInput{
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
	svc, _, _ := claimed(t)
	ad, err := svc.Reject(context.Background(), 1, 5, service.RejectInput{CustomerReason: "Price missing"})
	if err != nil {
		t.Fatal(err)
	}
	if ad.RejectionReasonInternal != nil {
		t.Fatalf("internal reason %q", *ad.RejectionReasonInternal)
	}
}

func TestRejectValidation(t *testing.T) {
	svc, _, _ := claimed(t)
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
	svc, _ := newService(newRepo())
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
	// One row beyond the page shows whether another page follows.
	if repo.lastQuery.Status != "pending" || repo.lastQuery.Limit != maxPageSize+1 {
		t.Fatalf("query %+v", repo.lastQuery)
	}
}

func TestListOutsideQueueNeedsAFilter(t *testing.T) {
	repo := newRepo()
	svc := NewAdReviewService(repo, nil)
	ctx := context.Background()
	for _, status := range []string{"", "active", "draft"} {
		if _, err := svc.List(ctx, model.ReviewFilter{Status: status}); code(t, err) != "AD_REVIEW_FILTER_REQUIRED" {
			t.Fatalf("status %q without filter: %v", status, err)
		}
	}
	for _, filter := range []model.ReviewFilter{
		{Status: "active", Search: "59"},
		{Status: "active", CategoryID: 11},
		{CityID: 1},
		{Search: "پژو ۲۰۶"},
	} {
		if _, err := svc.List(ctx, filter); err != nil {
			t.Fatalf("%+v: %v", filter, err)
		}
	}
}

func TestParseSearch(t *testing.T) {
	mobiles := []string{"+989121234567", "09121234567", "989121234567"}
	cases := []struct {
		in      string
		adID    int64
		mobiles []string
		title   string
		ok      bool
	}{
		{in: "", ok: true},
		{in: "59", adID: 59, ok: true},
		{in: "#59", adID: 59, ok: true},
		{in: "۵۹", adID: 59, ok: true},
		{in: "09121234567", mobiles: mobiles, ok: true},
		{in: "۰۹۱۲ ۱۲۳ ۴۵۶۷", mobiles: mobiles, ok: true},
		{in: "+98 912-123-4567", mobiles: mobiles, ok: true},
		{in: "9121234567", mobiles: mobiles, ok: true},
		{in: "  آپارتمان  ", title: "آپارتمان", ok: true},
		{in: "پژو ۲۰۶", title: "پژو 206", ok: true},
		{in: "ab", ok: false},
		{in: "99999999999999999999", ok: false},
	}
	for _, tc := range cases {
		got, ok := parseSearch(tc.in)
		if ok != tc.ok || got.adID != tc.adID || got.title != tc.title || len(got.mobiles) != len(tc.mobiles) {
			t.Errorf("%q: got %+v ok=%v", tc.in, got, ok)
			continue
		}
		for i := range tc.mobiles {
			if got.mobiles[i] != tc.mobiles[i] {
				t.Errorf("%q: mobiles %v", tc.in, got.mobiles)
			}
		}
	}
}

func TestListPagesWithCursorAndCapsTotal(t *testing.T) {
	repo := newRepo()
	base := time.Date(2026, 11, 4, 9, 0, 0, 0, time.UTC)
	for i := int64(1); i <= 3; i++ {
		at := base.Add(time.Duration(i) * time.Minute)
		repo.listItems = append(repo.listItems, model.ReviewListItem{ID: i, Status: "pending", PendingSince: &at, UpdatedAt: at})
	}
	repo.listTotal = service.MaxCountedTotal + 1
	svc := NewAdReviewService(repo, nil)
	ctx := context.Background()

	page, err := svc.List(ctx, model.ReviewFilter{Status: "pending", PageSize: 2})
	if err != nil {
		t.Fatal(err)
	}
	if len(page.Items) != 2 || page.NextCursor == "" || !page.TotalCapped || page.Total != service.MaxCountedTotal {
		t.Fatalf("page %+v", page)
	}
	if repo.countLimit != service.MaxCountedTotal+1 {
		t.Fatalf("count limit %d", repo.countLimit)
	}

	if _, err := svc.List(ctx, model.ReviewFilter{Status: "pending", PageSize: 2, Cursor: page.NextCursor}); err != nil {
		t.Fatal(err)
	}
	after := repo.lastQuery.After
	if after == nil || after.ID != 2 || !after.At.Equal(*repo.listItems[1].PendingSince) {
		t.Fatalf("cursor %+v", after)
	}

	if _, err := svc.List(ctx, model.ReviewFilter{Status: "pending", Cursor: "not-a-cursor"}); code(t, err) != "AD_REVIEW_INVALID_CURSOR" {
		t.Fatalf("bad cursor: %v", err)
	}
}

func TestShortTitleSearchIsRejected(t *testing.T) {
	_, err := NewAdReviewService(newRepo(), nil).List(context.Background(), model.ReviewFilter{Search: "ab"})
	if code(t, err) != "AD_REVIEW_SEARCH_TOO_SHORT" {
		t.Fatalf("got %v", err)
	}
}

func TestDecisionsRequireTheReviewersClaim(t *testing.T) {
	svc, _, _ := claimed(t)
	ctx := context.Background()
	if _, err := svc.Approve(ctx, 1, 6); code(t, err) != "AD_REVIEW_CLAIMED_BY_OTHER" {
		t.Fatalf("other reviewer approving: %v", err)
	}
	if _, err := svc.Approve(ctx, 3, 5); code(t, err) != "AD_REVIEW_NOT_CLAIMED" {
		t.Fatalf("approving an unclaimed ad: %v", err)
	}
}

func TestNextGivesEachReviewerADifferentAd(t *testing.T) {
	svc, _ := newService(newRepo())
	ctx := context.Background()
	first, err := svc.Next(ctx, 5)
	if err != nil {
		t.Fatal(err)
	}
	second, err := svc.Next(ctx, 6)
	if err != nil {
		t.Fatal(err)
	}
	if first.ID != 1 || second.ID != 3 {
		t.Fatalf("got %d and %d", first.ID, second.ID)
	}
	again, _ := svc.Next(ctx, 5)
	if again.ID != 1 {
		t.Fatalf("a reviewer asking again should keep their ad, got %d", again.ID)
	}
	if none, err := svc.Next(ctx, 7); err != nil || none != nil {
		t.Fatalf("empty queue: %+v %v", none, err)
	}
}

func TestClaimExpiresAndReturnsToQueue(t *testing.T) {
	svc, _, c := claimed(t)
	ctx := context.Background()
	c.advance(service.ClaimLease + time.Second)

	if _, err := svc.Approve(ctx, 1, 5); code(t, err) != "AD_REVIEW_CLAIM_EXPIRED" {
		t.Fatalf("approve after expiry: %v", err)
	}
	ad, err := svc.Get(ctx, 1)
	if err != nil || ad.ClaimedByStaffID != nil {
		t.Fatalf("expired claim still shown: %+v %v", ad, err)
	}
	if _, err := svc.Claim(ctx, 1, 6); err != nil {
		t.Fatalf("another reviewer should get the expired ad: %v", err)
	}
}

func TestExtendStopsAtMaxHold(t *testing.T) {
	svc, _, c := claimed(t)
	ctx := context.Background()
	start := c.at
	for c.at.Sub(start) < service.MaxClaimHold-time.Minute {
		c.advance(time.Minute)
		if _, err := svc.Extend(ctx, 1, 5); err != nil {
			t.Fatalf("extend at %v: %v", c.at.Sub(start), err)
		}
	}
	ad, _ := svc.Get(ctx, 1)
	if !ad.ClaimExpiresAt.Equal(start.Add(service.MaxClaimHold)) {
		t.Fatalf("expiry %v, want hold limit", ad.ClaimExpiresAt.Sub(start))
	}
	c.advance(2 * time.Minute)
	if _, err := svc.Extend(ctx, 1, 5); code(t, err) != "AD_REVIEW_CLAIM_EXPIRED" {
		t.Fatalf("extend past the hold limit: %v", err)
	}
}

func TestClaimingAnotherAdReleasesTheFirst(t *testing.T) {
	svc, repo, _ := claimed(t)
	if _, err := svc.Claim(context.Background(), 3, 5); err != nil {
		t.Fatal(err)
	}
	if repo.ads[1].claimedBy != nil {
		t.Fatal("reviewer still holds their first ad")
	}
}

func TestReleaseAndForceRelease(t *testing.T) {
	svc, repo, _ := claimed(t)
	ctx := context.Background()
	if err := svc.Release(ctx, 1, 6); err != nil || repo.ads[1].claimedBy == nil {
		t.Fatalf("another reviewer must not release the claim: %v", err)
	}
	if err := svc.ForceRelease(ctx, 1); err != nil || repo.ads[1].claimedBy != nil {
		t.Fatalf("force release: %v", err)
	}
	if err := svc.ForceRelease(ctx, 99); code(t, err) != "AD_REVIEW_NOT_FOUND" {
		t.Fatalf("force release missing: %v", err)
	}
}

func TestStartOfDayUsesTehranMidnight(t *testing.T) {
	// 22:00 UTC on Nov 3 is already 01:30 on Nov 4 in Tehran.
	got := startOfDay(time.Date(2026, 11, 3, 22, 0, 0, 0, time.UTC))
	want := time.Date(2026, 11, 3, 20, 30, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

type fakeCatalog struct{}

func (fakeCatalog) CategoriesByIDs(_ context.Context, ids []int, includeDescendants bool) ([]searchclient.Category, error) {
	category := searchclient.Category{ID: ids[0]}
	if includeDescendants && ids[0] == 1 {
		category.DescendantIDs = []int{1, 11, 12}
	}
	return []searchclient.Category{category}, nil
}

func (fakeCatalog) CitiesByIDs(context.Context, []int) ([]searchclient.City, error) {
	return nil, nil
}

func TestCategoryFilterIncludesSubcategories(t *testing.T) {
	repo := newRepo()
	svc := NewAdReviewService(repo, fakeCatalog{})
	if _, err := svc.List(context.Background(), model.ReviewFilter{CategoryID: 1}); err != nil {
		t.Fatal(err)
	}
	if got := repo.lastQuery.CategoryIDs; len(got) != 3 || got[1] != 11 {
		t.Fatalf("category ids %v", got)
	}
	if _, err := svc.List(context.Background(), model.ReviewFilter{CategoryID: 11}); err != nil {
		t.Fatal(err)
	}
	if got := repo.lastQuery.CategoryIDs; len(got) != 1 || got[0] != 11 {
		t.Fatalf("leaf category ids %v", got)
	}
}
