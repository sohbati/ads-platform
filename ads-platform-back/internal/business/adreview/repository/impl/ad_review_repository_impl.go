package impl

import (
	"context"
	"errors"
	"strings"
	"time"

	"ads-platform/internal/business/adreview/model"
	"ads-platform/internal/business/adreview/repository"
	adsmodel "ads-platform/internal/business/ads/model"

	"gorm.io/gorm"
)

type adReviewRepository struct {
	db *gorm.DB
}

func NewAdReviewRepository(db *gorm.DB) repository.AdReviewRepository {
	return &adReviewRepository{db: db}
}

const listColumns = `a.id, a.title, a.status, a.category_id, a.city_id,
	a.price_amount, a.price_type, a.currency,
	a.customer_id, c.name AS customer_name, c.mobile AS customer_mobile,
	a.pending_since, a.published_at, a.reviewed_at,
	(a.rejection_reason_customer IS NOT NULL) AS was_rejected,
	a.created_at, a.updated_at`

func (r *adReviewRepository) base(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).
		Table("ads AS a").
		Joins("JOIN customer AS c ON c.id = a.customer_id").
		Where("a.status <> ?", adsmodel.AdStatusDeleted)
}

func (r *adReviewRepository) List(ctx context.Context, filter model.ReviewFilter) ([]model.ReviewListItem, int64, error) {
	query := r.base(ctx)
	if filter.Status != "" {
		query = query.Where("a.status = ?", filter.Status)
	}
	if term := strings.TrimSpace(filter.Search); term != "" {
		like := "%" + term + "%"
		query = query.Where(
			"a.title ILIKE ? OR CAST(a.id AS TEXT) = ? OR c.mobile LIKE ? OR c.name ILIKE ?",
			like, term, like, like,
		)
	}

	var total int64
	if err := query.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// The review queue is first-come, first-served; other lists show recent activity first.
	order := "a.updated_at DESC, a.id DESC"
	if filter.Status == adsmodel.AdStatusPending {
		order = "a.pending_since ASC NULLS LAST, a.id ASC"
	}

	items := []model.ReviewListItem{}
	err := query.Select(listColumns).
		Order(order).
		Limit(filter.PageSize).
		Offset((filter.Page - 1) * filter.PageSize).
		Scan(&items).Error
	return items, total, err
}

func (r *adReviewRepository) Get(ctx context.Context, id int64) (*model.ReviewAd, error) {
	var ad model.ReviewAd
	err := r.base(ctx).
		Joins("LEFT JOIN staff AS s ON s.id = a.reviewed_by_staff_id").
		Where("a.id = ?", id).
		Select(listColumns + `, a.description, a.attrs, a.media, a.contact, a.location,
			a.rejection_reason_customer, a.rejection_reason_internal,
			a.reviewed_by_staff_id, s.name AS reviewed_by_name`).
		Take(&ad).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &ad, nil
}

func (r *adReviewRepository) Approve(ctx context.Context, id, staffID int64) (bool, error) {
	now := time.Now().UTC()
	res := r.db.WithContext(ctx).Model(&adsmodel.Ad{}).
		Where("id = ? AND status = ?", id, adsmodel.AdStatusPending).
		Updates(map[string]any{
			"status": adsmodel.AdStatusActive,
			// Keep the first publish date so an edited ad does not jump to the top of search.
			"published_at":              gorm.Expr("COALESCE(published_at, ?)", now),
			"pending_since":             nil,
			"rejection_reason_customer": nil,
			"rejection_reason_internal": nil,
			"reviewed_by_staff_id":      staffID,
			"reviewed_at":               now,
		})
	return res.RowsAffected > 0, res.Error
}

func (r *adReviewRepository) Reject(ctx context.Context, id, staffID int64, customerReason string, internalReason *string) (bool, error) {
	now := time.Now().UTC()
	res := r.db.WithContext(ctx).Model(&adsmodel.Ad{}).
		Where("id = ? AND status = ?", id, adsmodel.AdStatusPending).
		Updates(map[string]any{
			"status":                    adsmodel.AdStatusDraft,
			"pending_since":             nil,
			"rejection_reason_customer": customerReason,
			"rejection_reason_internal": internalReason,
			"reviewed_by_staff_id":      staffID,
			"reviewed_at":               now,
		})
	return res.RowsAffected > 0, res.Error
}

func (r *adReviewRepository) Exists(ctx context.Context, id int64) (bool, error) {
	var count int64
	err := r.base(ctx).Where("a.id = ?", id).Count(&count).Error
	return count > 0, err
}
