package viewmodel

import "ads-platform-ui/internal/core/i18n"

// QueryAdsPage is the browse/search ads page. Search is always set.
type QueryAdsPage struct {
	i18n.Page
	Search *SearchResults
}

// SearchResults holds rendered search API results.
type SearchResults struct {
	Query         string     `json:"query"`
	CategorySlug  string     `json:"category"`
	CategoryTitle string     `json:"category_title"`
	Total         int64      `json:"total"`
	Page          int        `json:"page"`
	TotalPages    int        `json:"total_pages"`
	HasMore       bool       `json:"has_more"`
	PrevURL       string     `json:"-"`
	NextURL       string     `json:"-"`
	Ads           []SearchAd `json:"ads"`
	Unavailable   bool       `json:"unavailable"`
}

// SearchAd is one result card.
type SearchAd struct {
	ID          int64  `json:"id,omitempty"`
	Href        string `json:"href,omitempty"`
	Title       string `json:"title"`
	Price       string `json:"price"`
	Location    string `json:"location"`
	Thumbnail   string `json:"thumbnail"`
	HasPhoto    bool   `json:"has_photo"`
	// PublishedAt is a UTC date for pages without JavaScript; local_time.js replaces it with
	// PublishedAtISO shown in the visitor's own time zone and language.
	PublishedAt    string `json:"published_at"`
	PublishedAtISO string `json:"published_at_iso,omitempty"`
	Views          int    `json:"views,omitempty"`
	Calls       int    `json:"calls,omitempty"`
	// Owner-only fields for the my-ads list.
	StatusLabel     string `json:"status_label,omitempty"`
	StatusTone      string `json:"status_tone,omitempty"`
	RejectionReason string `json:"rejection_reason,omitempty"`
}

// AdDetailPage is the public listing details page.
type AdDetailPage struct {
	i18n.Page
	Ad          *AdDetail
	NotFound    bool
	Unavailable bool
}

// AdDetail is one public listing.
type AdDetail struct {
	ID          int64
	Title       string
	Price       string
	Location    string
	PublishedAt string
	// PublishedAtISO is the exact moment, shown in the visitor's zone (see SearchAd).
	PublishedAtISO string
	Description    string
	Images      []string
	HasPhone    bool
	PhoneMasked string
	HasMap      bool
	MapLat      float64
	MapLng      float64
	Attrs       []AdAttr
}

// AdAttr is one labeled listing feature from the category schema.
type AdAttr struct {
	Label string
	Value string
}
