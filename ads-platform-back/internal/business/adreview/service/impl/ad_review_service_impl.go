package impl

import (
	"context"
	"strings"
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
}

func NewAdReviewService(repo repository.AdReviewRepository, catalog catalogLookup) service.AdReviewService {
	return &adReviewService{repo: repo, catalog: catalog}
}

func (s *adReviewService) List(ctx context.Context, filter model.ReviewFilter) (*model.ReviewPage, error) {
	filter.Status = strings.ToLower(strings.TrimSpace(filter.Status))
	if !listableStatuses[filter.Status] {
		return nil, exception.NewAppError(errorcode.ErrInvalidStatus.Code, errorcode.ErrInvalidStatus.HttpStatus, filter.Status)
	}
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 {
		filter.PageSize = defaultPageSize
	}
	if filter.PageSize > maxPageSize {
		filter.PageSize = maxPageSize
	}

	items, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	s.fillNames(ctx, items)
	return &model.ReviewPage{Items: items, Total: total, Page: filter.Page, PageSize: filter.PageSize}, nil
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
	changed, err := s.repo.Approve(ctx, id, staffID)
	if err != nil {
		return nil, err
	}
	if !changed {
		return nil, s.notPendingError(ctx, id)
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

	changed, err := s.repo.Reject(ctx, id, staffID, customerReason, internal)
	if err != nil {
		return nil, err
	}
	if !changed {
		return nil, s.notPendingError(ctx, id)
	}
	return s.Get(ctx, id)
}

// notPendingError explains why a decision changed nothing: the ad is gone, or already decided.
func (s *adReviewService) notPendingError(ctx context.Context, id int64) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return exception.NewAppError(errorcode.ErrNotFound.Code, errorcode.ErrNotFound.HttpStatus)
	}
	return exception.NewAppError(errorcode.ErrNotPending.Code, errorcode.ErrNotPending.HttpStatus)
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
