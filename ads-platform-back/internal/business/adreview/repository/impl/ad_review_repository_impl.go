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
	a.claimed_by_staff_id, cs.name AS claimed_by_name, a.claimed_at, a.claim_expires_at,
	a.created_at, a.updated_at`

// claimLockKey namespaces the per-reviewer advisory lock that keeps one reviewer to one claim.
const claimLockKey = 7301

// noClaim clears a claim; decisions apply it too, since a decided ad leaves the queue.
func noClaim() map[string]any {
	return map[string]any{
		"claimed_by_staff_id": nil,
		"claimed_at":          nil,
		"claim_expires_at":    nil,
	}
}

func withoutClaim(fields map[string]any) map[string]any {
	for key, value := range noClaim() {
		fields[key] = value
	}
	return fields
}

func (r *adReviewRepository) base(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).
		Table("ads AS a").
		Joins("JOIN customer AS c ON c.id = a.customer_id").
		Joins("LEFT JOIN staff AS cs ON cs.id = a.claimed_by_staff_id").
		Where("a.status <> ?", adsmodel.AdStatusDeleted)
}

// filtered applies a query's filters. Each one can use an index: the primary key, the customer
// mobile unique key, the title trigram index, or the category and city indexes.
func (r *adReviewRepository) filtered(ctx context.Context, q model.ReviewQuery) *gorm.DB {
	query := r.base(ctx)
	if q.Status != "" {
		query = query.Where("a.status = ?", q.Status)
	}
	if q.AdID > 0 {
		query = query.Where("a.id = ?", q.AdID)
	}
	if len(q.Mobiles) > 0 {
		query = query.Where("c.mobile IN ?", q.Mobiles)
	}
	if q.Title != "" {
		query = query.Where("a.title ILIKE ?", "%"+escapeLike(q.Title)+"%")
	}
	if len(q.CategoryIDs) > 0 {
		query = query.Where("a.category_id IN ?", q.CategoryIDs)
	}
	if q.CityID > 0 {
		query = query.Where("a.city_id = ?", q.CityID)
	}
	return query
}

func (r *adReviewRepository) List(ctx context.Context, q model.ReviewQuery, countLimit int) ([]model.ReviewListItem, int64, error) {
	// Count at most countLimit rows: an exact total over millions of ads costs as much as the search.
	var total int64
	capped := r.filtered(ctx, q).Select("1").Limit(countLimit)
	if err := r.db.WithContext(ctx).Table("(?) AS matches", capped).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pages continue from the last row's sort position rather than an OFFSET, so a deep page costs
	// the same as the first. The review queue is first-come, first-served; other lists show recent
	// activity first.
	query := r.filtered(ctx, q)
	order := "a.updated_at DESC, a.id DESC"
	if q.Status == adsmodel.AdStatusPending {
		order = "a.pending_since ASC, a.id ASC"
		if q.After != nil {
			query = query.Where("(a.pending_since, a.id) > (?, ?)", q.After.At, q.After.ID)
		}
	} else if q.After != nil {
		query = query.Where("(a.updated_at, a.id) < (?, ?)", q.After.At, q.After.ID)
	}

	items := []model.ReviewListItem{}
	err := query.Select(listColumns).Order(order).Limit(q.Limit).Scan(&items).Error
	return items, total, err
}

var likeEscaper = strings.NewReplacer(`\`, `\\`, "%", `\%`, "_", `\_`)

// escapeLike makes user text match literally inside an ILIKE pattern.
func escapeLike(term string) string {
	return likeEscaper.Replace(term)
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

func (r *adReviewRepository) Approve(ctx context.Context, id, staffID int64, now time.Time) (bool, error) {
	res := r.claimedBy(ctx, id, staffID, now).
		Updates(withoutClaim(map[string]any{
			"status": adsmodel.AdStatusActive,
			// Keep the first publish date so an edited ad does not jump to the top of search.
			"published_at":              gorm.Expr("COALESCE(published_at, ?)", now),
			"pending_since":             nil,
			"rejection_reason_customer": nil,
			"rejection_reason_internal": nil,
			"reviewed_by_staff_id":      staffID,
			"reviewed_at":               now,
		}))
	return res.RowsAffected > 0, res.Error
}

func (r *adReviewRepository) Reject(ctx context.Context, id, staffID int64, customerReason string, internalReason *string, now time.Time) (bool, error) {
	res := r.claimedBy(ctx, id, staffID, now).
		Updates(withoutClaim(map[string]any{
			"status":                    adsmodel.AdStatusDraft,
			"pending_since":             nil,
			"rejection_reason_customer": customerReason,
			"rejection_reason_internal": internalReason,
			"reviewed_by_staff_id":      staffID,
			"reviewed_at":               now,
		}))
	return res.RowsAffected > 0, res.Error
}

// claimedBy scopes an update to a pending ad that staffID holds a live claim on.
func (r *adReviewRepository) claimedBy(ctx context.Context, id, staffID int64, now time.Time) *gorm.DB {
	return r.db.WithContext(ctx).Model(&adsmodel.Ad{}).
		Where("id = ? AND status = ? AND claimed_by_staff_id = ? AND claim_expires_at > ?",
			id, adsmodel.AdStatusPending, staffID, now)
}

func (r *adReviewRepository) Exists(ctx context.Context, id int64) (bool, error) {
	var count int64
	err := r.base(ctx).Where("a.id = ?", id).Count(&count).Error
	return count > 0, err
}

// Claim changes below use UpdateColumns: taking an ad is not an edit, so updated_at stays.

// lockReviewer serializes one reviewer's claims, so a double click cannot claim two ads.
func lockReviewer(tx *gorm.DB, staffID int64) error {
	return tx.Exec("SELECT pg_advisory_xact_lock(?, ?)", claimLockKey, int32(staffID)).Error
}

func (r *adReviewRepository) ClaimNext(ctx context.Context, staffID int64, now, expiresAt time.Time) (int64, error) {
	var id int64
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := lockReviewer(tx, staffID); err != nil {
			return err
		}
		var held []int64
		if err := tx.Model(&adsmodel.Ad{}).
			Where("status = ? AND claimed_by_staff_id = ? AND claim_expires_at > ?",
				adsmodel.AdStatusPending, staffID, now).
			Order("claimed_at DESC").Limit(1).Pluck("id", &held).Error; err != nil {
			return err
		}
		if len(held) > 0 {
			id = held[0]
			return nil
		}
		// SKIP LOCKED: reviewers asking at the same moment each get a different ad, without waiting.
		var claimed []int64
		if err := tx.Raw(`
			WITH next AS (
				SELECT id FROM ads
				WHERE status = ? AND (claim_expires_at IS NULL OR claim_expires_at <= ?)
				ORDER BY pending_since ASC NULLS LAST, id ASC
				LIMIT 1
				FOR UPDATE SKIP LOCKED
			)
			UPDATE ads AS a
			SET claimed_by_staff_id = ?, claimed_at = ?, claim_expires_at = ?
			FROM next WHERE a.id = next.id
			RETURNING a.id`,
			adsmodel.AdStatusPending, now, staffID, now, expiresAt,
		).Scan(&claimed).Error; err != nil {
			return err
		}
		if len(claimed) > 0 {
			id = claimed[0]
		}
		return nil
	})
	return id, err
}

func (r *adReviewRepository) Claim(ctx context.Context, id, staffID int64, now, expiresAt time.Time) (bool, error) {
	changed := false
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := lockReviewer(tx, staffID); err != nil {
			return err
		}
		if err := tx.Model(&adsmodel.Ad{}).
			Where("claimed_by_staff_id = ? AND id <> ?", staffID, id).
			UpdateColumns(noClaim()).Error; err != nil {
			return err
		}
		res := tx.Model(&adsmodel.Ad{}).
			Where("id = ? AND status = ?", id, adsmodel.AdStatusPending).
			Where("claimed_by_staff_id IS NULL OR claimed_by_staff_id = ? OR claim_expires_at <= ?", staffID, now).
			UpdateColumns(map[string]any{
				"claimed_by_staff_id": staffID,
				"claimed_at":          now,
				"claim_expires_at":    expiresAt,
			})
		changed = res.RowsAffected > 0
		return res.Error
	})
	return changed, err
}

func (r *adReviewRepository) Extend(ctx context.Context, id, staffID int64, now, expiresAt time.Time) (bool, error) {
	res := r.claimedBy(ctx, id, staffID, now).UpdateColumn("claim_expires_at", expiresAt)
	return res.RowsAffected > 0, res.Error
}

func (r *adReviewRepository) Release(ctx context.Context, id, staffID int64) error {
	return r.db.WithContext(ctx).Model(&adsmodel.Ad{}).
		Where("id = ? AND claimed_by_staff_id = ?", id, staffID).
		UpdateColumns(noClaim()).Error
}

func (r *adReviewRepository) ReleaseAll(ctx context.Context, staffID int64) error {
	return r.db.WithContext(ctx).Model(&adsmodel.Ad{}).
		Where("claimed_by_staff_id = ?", staffID).
		UpdateColumns(noClaim()).Error
}

func (r *adReviewRepository) ForceRelease(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Model(&adsmodel.Ad{}).
		Where("id = ? AND claimed_by_staff_id IS NOT NULL", id).
		UpdateColumns(noClaim()).Error
}

func (r *adReviewRepository) ReviewerStats(ctx context.Context, since time.Time) ([]model.ReviewerStats, error) {
	// A rejected ad keeps its customer reason until it is approved, which tells the two apart.
	stats := []model.ReviewerStats{}
	err := r.db.WithContext(ctx).
		Table("ads AS a").
		Joins("JOIN staff AS s ON s.id = a.reviewed_by_staff_id").
		Where("a.reviewed_at >= ?", since).
		Select(`s.id AS staff_id, s.name,
			COUNT(*) FILTER (WHERE a.rejection_reason_customer IS NULL) AS approved,
			COUNT(*) FILTER (WHERE a.rejection_reason_customer IS NOT NULL) AS rejected`).
		Group("s.id, s.name").
		Order("COUNT(*) DESC, s.name").
		Scan(&stats).Error
	return stats, err
}
