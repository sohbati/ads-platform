package model

import (
	"encoding/json"
	"time"
)

// ReviewListItem is one row of the panel ads list.
type ReviewListItem struct {
	ID             int64      `json:"id"`
	Title          string     `json:"title"`
	Status         string     `json:"status"`
	CategoryID     int        `json:"category_id"`
	CategoryTitle  string     `json:"category_title,omitempty" gorm:"-"`
	CityID         int        `json:"city_id"`
	CityName       string     `json:"city_name,omitempty" gorm:"-"`
	PriceAmount    *int64     `json:"price_amount"`
	PriceType      string     `json:"price_type"`
	Currency       string     `json:"currency"`
	CustomerID     int64      `json:"customer_id"`
	CustomerName   string     `json:"customer_name"`
	CustomerMobile string     `json:"customer_mobile"`
	PendingSince   *time.Time `json:"pending_since"`
	PublishedAt    *time.Time `json:"published_at"`
	ReviewedAt     *time.Time `json:"reviewed_at"`
	// WasRejected marks ads whose last review was a rejection (back in draft, or resubmitted).
	WasRejected bool `json:"was_rejected"`
	// The reviewer working on the ad now; all nil when nobody holds a live claim.
	ClaimedByStaffID *int64     `json:"claimed_by_staff_id"`
	ClaimedByName    *string    `json:"claimed_by_name"`
	ClaimedAt        *time.Time `json:"claimed_at"`
	ClaimExpiresAt   *time.Time `json:"claim_expires_at"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// ReviewAd is the full panel view of one ad, including internal review notes.
type ReviewAd struct {
	ReviewListItem
	Description             string          `json:"description"`
	Attrs                   json.RawMessage `json:"attrs" gorm:"type:jsonb"`
	Media                   json.RawMessage `json:"media" gorm:"type:jsonb"`
	Contact                 json.RawMessage `json:"contact" gorm:"type:jsonb"`
	Location                json.RawMessage `json:"location" gorm:"type:jsonb"`
	RejectionReasonCustomer *string         `json:"rejection_reason_customer"`
	RejectionReasonInternal *string         `json:"rejection_reason_internal"`
	ReviewedByStaffID       *int64          `json:"reviewed_by_staff_id"`
	ReviewedByName          *string         `json:"reviewed_by_name"`
}

// ReviewFilter is the panel list request as it arrives.
type ReviewFilter struct {
	Status string
	// Search is an ad ID, a customer mobile number, or part of a title.
	Search     string
	CategoryID int
	CityID     int
	// Cursor continues a previous page; empty starts from the top.
	Cursor   string
	PageSize int
}

// ReviewQuery is a validated ReviewFilter, ready for the database.
type ReviewQuery struct {
	Status string
	AdID   int64
	// Mobiles holds every stored form of one customer mobile (+989…, 09…, 989…).
	Mobiles []string
	Title   string
	// CategoryIDs is the chosen category and everything under it.
	CategoryIDs []int
	CityID      int
	After       *ReviewCursor
	Limit       int
}

// ReviewCursor is the sort position of the last row on the previous page.
type ReviewCursor struct {
	At time.Time `json:"t"`
	ID int64     `json:"id"`
}

type ReviewPage struct {
	Items []ReviewListItem `json:"items"`
	// Total counts matches up to MaxCountedTotal; TotalCapped means there are more.
	Total       int64  `json:"total"`
	TotalCapped bool   `json:"total_capped"`
	NextCursor  string `json:"next_cursor,omitempty"`
	PageSize    int    `json:"pageSize"`
}

// ReviewerStats counts one reviewer's decisions since the start of today (Tehran time).
type ReviewerStats struct {
	StaffID  int64  `json:"staff_id"`
	Name     string `json:"name"`
	Approved int64  `json:"approved"`
	Rejected int64  `json:"rejected"`
}
