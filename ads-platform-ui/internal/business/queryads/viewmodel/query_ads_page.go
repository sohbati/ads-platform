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
	PublishedAt string `json:"published_at"`
	Views       int    `json:"views,omitempty"`
	Calls       int    `json:"calls,omitempty"`
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
	Description string
	Images      []string
	HasPhone    bool
	PhoneMasked string
}
