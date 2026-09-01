package model

import (
	"encoding/json"
	"time"
)

// Ad maps to ads_platform_schema.ads
type Ad struct {
	ID          int64           `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID      int64           `json:"user_id" gorm:"not null"`
	CategoryID  int             `json:"category_id" gorm:"not null"`
	CityID      int             `json:"city_id" gorm:"not null"`
	Title       string          `json:"title" gorm:"type:varchar(120);not null"`
	Description string          `json:"description" gorm:"type:text;not null"`
	Status      string          `json:"status" gorm:"type:varchar(20);not null;default:draft"`
	PriceAmount *int64          `json:"price_amount"`
	PriceType   string          `json:"price_type" gorm:"type:varchar(20);not null;default:fixed"`
	Currency    string          `json:"currency" gorm:"type:char(3);not null;default:IRR"`
	Attrs       json.RawMessage `json:"attrs" gorm:"type:jsonb;not null"`
	Media       json.RawMessage `json:"media" gorm:"type:jsonb;not null"`
	Contact     json.RawMessage `json:"contact" gorm:"type:jsonb;not null"`
	Location    json.RawMessage `json:"location" gorm:"type:jsonb;not null"`
	Slug        *string         `json:"slug" gorm:"type:varchar(160)"`
	PublishedAt *time.Time      `json:"published_at"`
	ExpiresAt   *time.Time      `json:"expires_at"`
	CreatedAt   time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
}

const (
	AdStatusDraft    = "draft"
	AdStatusPending  = "pending"
	AdStatusActive   = "active"
	AdStatusRejected = "rejected"
	AdStatusExpired  = "expired"
	AdStatusDeleted  = "deleted"
)

func (Ad) TableName() string {
	return "ads"
}

// UserAdItem is a list-card projection of an ad owned by the current user.
type UserAdItem struct {
	ID           int64   `json:"id"`
	Title        string  `json:"title"`
	PriceAmount  *int64  `json:"price_amount"`
	PriceType    string  `json:"price_type"`
	Currency     string  `json:"currency"`
	CityID       int     `json:"city_id"`
	CityName     string  `json:"city_name,omitempty"`
	Neighborhood string  `json:"neighborhood,omitempty"`
	Thumbnail    string  `json:"thumbnail,omitempty"`
	HasPhoto     bool    `json:"has_photo"`
	CategoryID   int     `json:"category_id"`
	Slug         string  `json:"slug,omitempty"`
	Status       string  `json:"status"`
	PublishedAt  *string `json:"published_at,omitempty"`
}

// PublicAd is the public details view of an active listing.
type PublicAd struct {
	ID           int64         `json:"id"`
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	PriceAmount  *int64        `json:"price_amount"`
	PriceType    string        `json:"price_type"`
	Currency     string        `json:"currency"`
	CityID       int           `json:"city_id"`
	CityName     string        `json:"city_name,omitempty"`
	Neighborhood string        `json:"neighborhood,omitempty"`
	Media        []PublicMedia `json:"media"`
	PublishedAt  *string       `json:"published_at,omitempty"`
	HasPhone     bool          `json:"has_phone"`
	PhoneMasked  string        `json:"phone_masked,omitempty"`
}

// PublicContact is the seller phone, returned only after the viewer is logged in.
type PublicContact struct {
	Phone string `json:"phone"`
}

// AdStatsItem is owner-scoped totals for one ad over a day range.
type AdStatsItem struct {
	AdID           int64 `json:"ad_id"`
	Views          int   `json:"views"`
	UniqueViewers  int   `json:"unique_viewers"`
	ContactReveals int   `json:"contact_reveals"`
	Calls          int   `json:"calls"`
}

// AdStatsResponse is GET /api/v1/users/:userId/ad-stats.
type AdStatsResponse struct {
	From string        `json:"from"`
	To   string        `json:"to"`
	Ads  []AdStatsItem `json:"ads"`
}

// PublicMedia is one photo on the public details page. URL is the full WebP.
type PublicMedia struct {
	URL     string `json:"url"`
	Thumb   string `json:"thumb,omitempty"`
	IsCover bool   `json:"is_cover"`
}
