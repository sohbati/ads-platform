package service

import (
	"context"

	"ads-platform-ui/internal/business/queryads/viewmodel"
	"ads-platform-ui/internal/core/i18n"
)

// SearchParams carries the search request coming from the search bar or
// category/city links.
type SearchParams struct {
	Query    string
	Category string
	Page     int
}

type QueryAdsService interface {
	BuildPage(loc i18n.Locale, appName, citySlug, currentPath string, locationSlugs []string) viewmodel.QueryAdsPage
	BuildSearchPage(ctx context.Context, loc i18n.Locale, appName, citySlug, currentPath string, locationSlugs []string, params SearchParams) viewmodel.QueryAdsPage
}
