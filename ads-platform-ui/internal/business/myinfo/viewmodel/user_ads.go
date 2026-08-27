package viewmodel

import (
	queryadsvm "ads-platform-ui/internal/business/queryads/viewmodel"
	"ads-platform-ui/internal/core/i18n"
)

// UserAdsPage is the logged-in user's ad list.
type UserAdsPage struct {
	i18n.Page
	Ads         []queryadsvm.SearchAd
	Unavailable bool
}
