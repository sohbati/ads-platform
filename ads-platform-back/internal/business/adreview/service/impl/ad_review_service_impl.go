package impl

import (
	"context"
	"strings"
	"time"
	"unicode/utf8"

	"ads-platform/internal/business/adreview/errorcode"
	"ads-platform/internal/business/adreview/model"
	"ads-platform/internal/business/adreview/repository"
	"ads-platform/internal/business/adreview/service"
	adsmodel "ads-platform/internal/business/ads/model"
	searchclient "ads-platform/internal/business/search/client"
	"ads-platform/internal/core/exception"
)

const (
	defaultPageSize = 20
	maxPageSize     = 100
)

var listableStatuses = map[string]bool{
	"":                        true,
	adsmodel.AdStatusDraft:    true,
	adsmodel.AdStatusPending:  true,
	adsmodel.AdStatusActive:   true,
	adsmodel.AdStatusRejected: true,
	adsmodel.AdStatusExpired:  true,
}

// catalogLookup resolves display names; a nil lookup or a failed call leaves names empty.
type catalogLookup interface {
	CategoriesByIDs(ctx context.Context, ids []int, includeDescendants bool) ([]searchclient.Category, error)
	CitiesByIDs(ctx context.Context, ids []int) ([]searchclient.City, error)
}

type adReviewService struct {
	repo    repository.AdReviewRepository
	catalog catalogLookup
	now     func() time.Time
}

func NewAdReviewService(repo repository.AdReviewRepository, catalog catalogLookup) service.AdReviewService {
	return &adReviewService{repo: repo, catalog: catalog, now: func() time.Time { return time.Now().UTC() }}
}

func (s *adReviewService) List(ctx context.Context, filter model.ReviewFilter) (*model.ReviewPage, error) {
	q, err := s.buildQuery(ctx, filter)
	if err != nil {
		return nil, err
	}

	pageSize := q.Limit
	q.Limit++ // one extra row tells whether another page follows
	items, total, err := s.repo.List(ctx, q, service.MaxCountedTotal+1)
	if err != nil {
		return nil, err
	}

	page := &model.ReviewPage{Total: total, PageSize: pageSize}
	if total > service.MaxCountedTotal {
		page.Total, page.TotalCapped = service.MaxCountedTotal, true
	}
	if len(items) > pageSize {
		items = items[:pageSize]
		page.NextCursor = encodeCursor(cursorAfter(items[pageSize-1], q.Status == adsmodel.AdStatusPending))
	}
	now := s.now()
	for i := range items {
		hideStaleClaim(&items[i], now)
	}
	s.fillNames(ctx, items)
	page.Items = items
	return page, nil
}

// buildQuery validates a list request. Outside the review queue it needs at least one filter, so
// the panel never pages through every ad on the platform.
func (s *adReviewService) buildQuery(ctx context.Context, filter model.ReviewFilter) (model.ReviewQuery, error) {
	fail := func(code *errorcode.ErrorCode, params ...string) (model.ReviewQuery, error) {
		return model.ReviewQuery{}, exception.NewAppError(code.Code, code.HttpStatus, params...)
	}

	q := model.ReviewQuery{
		Status: strings.ToLower(strings.TrimSpace(filter.Status)),
		CityID: filter.CityID,
		Limit:  filter.PageSize,
	}
	if !listableStatuses[q.Status] {
		return fail(errorcode.ErrInvalidStatus, q.Status)
	}
	if q.Limit < 1 {
		q.Limit = defaultPageSize
	}
	if q.Limit > maxPageSize {
		q.Limit = maxPageSize
	}

	terms, ok := parseSearch(filter.Search)
	if !ok {
		return fail(errorcode.ErrSearchTooShort)
	}
	q.AdID, q.Mobiles, q.Title = terms.adID, terms.mobiles, terms.title

	if filter.CategoryID > 0 {
		ids, err := s.categoryWithDescendants(ctx, filter.CategoryID)
		if err != nil {
			return model.ReviewQuery{}, err
		}
		q.CategoryIDs = ids
	}

	after, ok := decodeCursor(filter.Cursor)
	if !ok {
		return fail(errorcode.ErrInvalidCursor)
	}
	q.After = after

	narrowed := q.AdID > 0 || len(q.Mobiles) > 0 || q.Title != "" || len(q.CategoryIDs) > 0 || q.CityID > 0
	if q.Status != adsmodel.AdStatusPending && !narrowed {
		return fail(errorcode.ErrFilterRequired)
	}
	return q, nil
}

// categoryWithDescendants expands a category to the leaf categories ads are filed under.
func (s *adReviewService) categoryWithDescendants(ctx context.Context, id int) ([]int, error) {
	if s.catalog == nil {
		return []int{id}, nil
	}
	found, err := s.catalog.CategoriesByIDs(ctx, []int{id}, true)
	if err != nil {
		return nil, exception.NewAppError(errorcode.ErrCatalogUnavailable.Code, errorcode.ErrCatalogUnavailable.HttpStatus)
	}
	// The cache service lists the category and everything under it in DescendantIDs.
	for _, category := range found {
		if category.ID == id && len(category.DescendantIDs) > 0 {
			return category.DescendantIDs, nil
		}
	}
	return []int{id}, nil
}

func (s *adReviewService) Get(ctx context.Context, id int64) (*model.ReviewAd, error) {
	if id <= 0 {
		return nil, exception.NewAppError(errorcode.ErrInvalidID.Code, errorcode.ErrInvalidID.HttpStatus)
	}
	ad, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if ad == nil {
		return nil, exception.NewAppError(errorcode.ErrNotFound.Code, errorcode.ErrNotFound.HttpStatus)
	}
	hideStaleClaim(&ad.ReviewListItem, s.now())
	items := []model.ReviewListItem{ad.ReviewListItem}
	s.fillNames(ctx, items)
	ad.ReviewListItem = items[0]
	return ad, nil
}

func (s *adReviewService) Approve(ctx context.Context, id, staffID int64) (*model.ReviewAd, error) {
	if err := requireStaff(staffID); err != nil {
		return nil, err
	}
	if id <= 0 {
		return nil, exception.NewAppError(errorcode.ErrInvalidID.Code, errorcode.ErrInvalidID.HttpStatus)
	}
	changed, err := s.repo.Approve(ctx, id, staffID, s.now())
	if err != nil {
		return nil, err
	}
	if !changed {
		return nil, s.claimError(ctx, id, staffID)
	}
	return s.Get(ctx, id)
}

func (s *adReviewService) Reject(ctx context.Context, id, staffID int64, input service.RejectInput) (*model.ReviewAd, error) {
	if err := requireStaff(staffID); err != nil {
		return nil, err
	}
	if id <= 0 {
		return nil, exception.NewAppError(errorcode.ErrInvalidID.Code, errorcode.ErrInvalidID.HttpStatus)
	}
	customerReason := strings.TrimSpace(input.CustomerReason)
	if customerReason == "" {
		return nil, exception.NewAppError(errorcode.ErrReasonRequired.Code, errorcode.ErrReasonRequired.HttpStatus)
	}
	internalReason := strings.TrimSpace(input.InternalReason)
	if utf8.RuneCountInString(customerReason) > service.MaxReasonLength ||
		utf8.RuneCountInString(internalReason) > service.MaxReasonLength {
		return nil, exception.NewAppError(errorcode.ErrReasonTooLong.Code, errorcode.ErrReasonTooLong.HttpStatus)
	}
	var internal *string
	if internalReason != "" {
		internal = &internalReason
	}

	changed, err := s.repo.Reject(ctx, id, staffID, customerReason, internal, s.now())
	if err != nil {
		return nil, err
	}
	if !changed {
		return nil, s.claimError(ctx, id, staffID)
	}
	return s.Get(ctx, id)
}

func (s *adReviewService) Next(ctx context.Context, staffID int64) (*model.ReviewAd, error) {
	if err := requireStaff(staffID); err != nil {
		return nil, err
	}
	now := s.now()
	id, err := s.repo.ClaimNext(ctx, staffID, now, now.Add(service.ClaimLease))
	if err != nil || id == 0 {
		return nil, err
	}
	return s.Get(ctx, id)
}

func (s *adReviewService) Claim(ctx context.Context, id, staffID int64) (*model.ReviewAd, error) {
	if err := validateClaim(id, staffID); err != nil {
		return nil, err
	}
	now := s.now()
	changed, err := s.repo.Claim(ctx, id, staffID, now, now.Add(service.ClaimLease))
	if err != nil {
		return nil, err
	}
	if !changed {
		return nil, s.claimError(ctx, id, staffID)
	}
	return s.Get(ctx, id)
}

func (s *adReviewService) Extend(ctx context.Context, id, staffID int64) (*model.ReviewAd, error) {
	if err := validateClaim(id, staffID); err != nil {
		return nil, err
	}
	ad, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	now := s.now()
	if ad == nil || !holdsLiveClaim(&ad.ReviewListItem, staffID, now) {
		return nil, s.claimError(ctx, id, staffID)
	}
	expiresAt := now.Add(service.ClaimLease)
	if limit := ad.ClaimedAt.Add(service.MaxClaimHold); limit.Before(expiresAt) {
		expiresAt = limit
	}
	changed, err := s.repo.Extend(ctx, id, staffID, now, expiresAt)
	if err != nil {
		return nil, err
	}
	if !changed {
		return nil, s.claimError(ctx, id, staffID)
	}
	return s.Get(ctx, id)
}

func (s *adReviewService) Release(ctx context.Context, id, staffID int64) error {
	if err := validateClaim(id, staffID); err != nil {
		return err
	}
	return s.repo.Release(ctx, id, staffID)
}

func (s *adReviewService) ReleaseAll(ctx context.Context, staffID int64) error {
	if err := requireStaff(staffID); err != nil {
		return err
	}
	return s.repo.ReleaseAll(ctx, staffID)
}

func (s *adReviewService) ForceRelease(ctx context.Context, id int64) error {
	if id <= 0 {
		return exception.NewAppError(errorcode.ErrInvalidID.Code, errorcode.ErrInvalidID.HttpStatus)
	}
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return exception.NewAppError(errorcode.ErrNotFound.Code, errorcode.ErrNotFound.HttpStatus)
	}
	return s.repo.ForceRelease(ctx, id)
}

func (s *adReviewService) TodayStats(ctx context.Context) ([]model.ReviewerStats, error) {
	return s.repo.ReviewerStats(ctx, startOfDay(s.now()))
}

// claimError explains why a claim or decision changed nothing.
func (s *adReviewService) claimError(ctx context.Context, id, staffID int64) error {
	ad, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	fail := func(code *errorcode.ErrorCode) error {
		return exception.NewAppError(code.Code, code.HttpStatus)
	}
	now := s.now()
	switch {
	case ad == nil:
		return fail(errorcode.ErrNotFound)
	case ad.Status != adsmodel.AdStatusPending:
		return fail(errorcode.ErrNotPending)
	case ad.ClaimedByStaffID != nil && *ad.ClaimedByStaffID != staffID && claimLive(&ad.ReviewListItem, now):
		return fail(errorcode.ErrClaimedByOther)
	case ad.ClaimedByStaffID != nil && *ad.ClaimedByStaffID == staffID:
		return fail(errorcode.ErrClaimExpired)
	default:
		return fail(errorcode.ErrNotClaimed)
	}
}

func claimLive(item *model.ReviewListItem, now time.Time) bool {
	return item.ClaimedByStaffID != nil && item.ClaimExpiresAt != nil && item.ClaimExpiresAt.After(now)
}

func holdsLiveClaim(item *model.ReviewListItem, staffID int64, now time.Time) bool {
	return item.Status == adsmodel.AdStatusPending && claimLive(item, now) &&
		*item.ClaimedByStaffID == staffID && item.ClaimedAt != nil
}

// hideStaleClaim drops a claim that no longer counts, so the panel shows only who is reviewing now.
func hideStaleClaim(item *model.ReviewListItem, now time.Time) {
	if item.Status != adsmodel.AdStatusPending || !claimLive(item, now) {
		item.ClaimedByStaffID, item.ClaimedByName, item.ClaimedAt, item.ClaimExpiresAt = nil, nil, nil, nil
	}
}

// startOfDay is midnight today in Tehran, as the UTC wall-clock time the database stores.
func startOfDay(now time.Time) time.Time {
	loc, err := time.LoadLocation("Asia/Tehran")
	if err != nil {
		loc = time.FixedZone("Tehran", 3*3600+30*60)
	}
	local := now.In(loc)
	return time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, loc).UTC()
}

func validateClaim(id, staffID int64) error {
	if err := requireStaff(staffID); err != nil {
		return err
	}
	if id <= 0 {
		return exception.NewAppError(errorcode.ErrInvalidID.Code, errorcode.ErrInvalidID.HttpStatus)
	}
	return nil
}

func requireStaff(staffID int64) error {
	if staffID <= 0 {
		return exception.NewAppError(errorcode.ErrStaffRequired.Code, errorcode.ErrStaffRequired.HttpStatus)
	}
	return nil
}

func (s *adReviewService) fillNames(ctx context.Context, items []model.ReviewListItem) {
	if s.catalog == nil || len(items) == 0 {
		return
	}
	cityIDs, categoryIDs := map[int]bool{}, map[int]bool{}
	for _, item := range items {
		cityIDs[item.CityID] = true
		categoryIDs[item.CategoryID] = true
	}

	cities := map[int]string{}
	if found, err := s.catalog.CitiesByIDs(ctx, keys(cityIDs)); err == nil {
		for _, city := range found {
			cities[city.ID] = city.Name
		}
	}
	categories := map[int]string{}
	if found, err := s.catalog.CategoriesByIDs(ctx, keys(categoryIDs), false); err == nil {
		for _, category := range found {
			categories[category.ID] = category.Title
		}
	}
	for i := range items {
		items[i].CityName = cities[items[i].CityID]
		items[i].CategoryTitle = categories[items[i].CategoryID]
	}
}

func keys(set map[int]bool) []int {
	out := make([]int, 0, len(set))
	for id := range set {
		if id > 0 {
			out = append(out, id)
		}
	}
	return out
}
