package viewmodel

import (
	"ads-platform-ui/internal/core/i18n"
	"ads-platform-ui/internal/domain"
)

// QueryAdsPage is the browse/search ads page.
type QueryAdsPage struct {
	i18n.Page
	Categories    []domain.Category
	PopularCities []domain.City
}
