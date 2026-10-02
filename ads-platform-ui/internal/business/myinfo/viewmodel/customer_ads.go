package viewmodel

import (
	queryadsvm "ads-platform-ui/internal/business/queryads/viewmodel"
	"ads-platform-ui/internal/core/i18n"
)

// CustomerAdsPage is the logged-in customer's ad list.
type CustomerAdsPage struct {
	i18n.Page
	Ads         []queryadsvm.SearchAd
	Unavailable bool
}
