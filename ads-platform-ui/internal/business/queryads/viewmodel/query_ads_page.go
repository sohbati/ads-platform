package viewmodel

import (
	"ads-platform-ui/internal/core/i18n"
	"ads-platform-ui/internal/domain"
)

// QueryAdsPage is the browse/search ads page. Search is nil on the landing
// page and set when the request carries a query or category.
type QueryAdsPage struct {
	i18n.Page
	Categories    []domain.Category
	PopularCities []domain.City
	Search        *SearchResults
}

// SearchResults holds rendered search API results.
type SearchResults struct {
	Query         string
	CategorySlug  string
	CategoryTitle string
	Total         int64
	Page          int
	TotalPages    int
	PrevURL       string
	NextURL       string
	Ads           []SearchAd
	Unavailable   bool
}

// SearchAd is one result card.
type SearchAd struct {
	Title       string
	Price       string
	Location    string
	Thumbnail   string
	HasPhoto    bool
	PublishedAt string
}
