package service

import (
	"ads-platform-ui/internal/business/queryads/viewmodel"
	"ads-platform-ui/internal/core/i18n"
)

type QueryAdsService interface {
	BuildPage(loc i18n.Locale, appName, citySlug, currentPath string) viewmodel.QueryAdsPage
}
