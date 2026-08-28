package impl

import (
	"context"
	"errors"
	"net/url"
	"strconv"
	"strings"
	"time"

	searchclient "ads-platform-ui/internal/business/queryads/client"
	"ads-platform-ui/internal/business/queryads/service"
	"ads-platform-ui/internal/business/queryads/viewmodel"
	"ads-platform-ui/internal/core/cities"
	"ads-platform-ui/internal/core/i18n"
	"ads-platform-ui/internal/core/media"
)

const searchPageSize = 24

type QueryAdsServiceImpl struct {
	i18n     *i18n.Registry
	cities   *cities.Catalog
	search   *searchclient.SearchClient
	mediaCDN string
}

func NewQueryAdsService(reg *i18n.Registry, catalog *cities.Catalog, search *searchclient.SearchClient, mediaCDN string) *QueryAdsServiceImpl {
	return &QueryAdsServiceImpl{i18n: reg, cities: catalog, search: search, mediaCDN: mediaCDN}
}

// BuildSearchPage calls the search API (via BFF) and renders results on the
// query-ads page. Multiple selected cities are sent as iran + cities ids.
func (s *QueryAdsServiceImpl) BuildSearchPage(ctx context.Context, loc i18n.Locale, appName, citySlug, currentPath string, locationSlugs []string, params service.SearchParams) viewmodel.QueryAdsPage {
	page := i18n.BuildPage(s.i18n, s.cities, loc, appName, citySlug, currentPath, locationSlugs)
	t := s.i18n.MessagesFor(loc)
	if params.Query != "" {
		page.Title = appName + " — " + t.Search.ResultsTitle
	} else {
		page.Title = appName + " — " + t.Nav.QueryAds
	}
	page.SearchQuery = params.Query

	category := strings.ToLower(strings.TrimSpace(params.Category))
	if category == "" {
		category = "all"
	}
	place, citiesCSV := s.cities.SearchPlace(locationSlugs, page.CitySlug)
	if place == "" {
		place = "iran"
	}
	if params.Page < 1 {
		params.Page = 1
	}

	results := &viewmodel.SearchResults{
		Query:        params.Query,
		CategorySlug: category,
		Page:         params.Page,
	}
	vm := viewmodel.QueryAdsPage{Page: page, Search: results}

	resp, err := s.search.Search(ctx, place, category, params.Query, params.Page, citiesCSV)
	if err != nil {
		// Unknown place/category reads as "no results"; everything else is an outage.
		results.Unavailable = !errors.Is(err, searchclient.ErrNotFound)
		return vm
	}

	results.CategoryTitle = resp.CategoryTitle
	results.Total = resp.Pagination.Total
	results.Page = resp.Pagination.Page
	limit := resp.Pagination.Limit
	if limit < 1 {
		limit = searchPageSize
	}
	results.TotalPages = int((resp.Pagination.Total + int64(limit) - 1) / int64(limit))
	results.HasMore = results.Page < results.TotalPages
	if results.Page > 1 {
		results.PrevURL = searchURL(params.Query, params.Category, results.Page-1)
	}
	if results.HasMore {
		results.NextURL = searchURL(params.Query, params.Category, results.Page+1)
	}

	results.Ads = make([]viewmodel.SearchAd, 0, len(resp.Ads))
	for _, ad := range resp.Ads {
		results.Ads = append(results.Ads, toSearchAd(ad, t, s.mediaCDN))
	}
	return vm
}

func (s *QueryAdsServiceImpl) BuildDetailPage(ctx context.Context, loc i18n.Locale, appName, citySlug, currentPath string, locationSlugs []string, adID int64) viewmodel.AdDetailPage {
	page := i18n.BuildPage(s.i18n, s.cities, loc, appName, citySlug, currentPath, locationSlugs)
	t := s.i18n.MessagesFor(loc)
	vm := viewmodel.AdDetailPage{Page: page}

	if adID <= 0 {
		vm.NotFound = true
		page.Title = appName + " — " + t.AdDetail.NotFound
		vm.Page = page
		return vm
	}

	ad, err := s.search.GetAd(ctx, adID)
	if err != nil {
		if errors.Is(err, searchclient.ErrNotFound) {
			vm.NotFound = true
			page.Title = appName + " — " + t.AdDetail.NotFound
		} else {
			vm.Unavailable = true
			page.Title = appName + " — " + t.AdDetail.Unavailable
		}
		vm.Page = page
		return vm
	}

	page.Title = appName + " — " + ad.Title
	vm.Page = page
	vm.Ad = toAdDetail(ad, t, s.mediaCDN)
	return vm
}

func toAdDetail(ad *searchclient.PublicAd, t i18n.Messages, mediaCDN string) *viewmodel.AdDetail {
	out := &viewmodel.AdDetail{
		ID:          ad.ID,
		Title:       ad.Title,
		Description: ad.Description,
		Location:    ad.CityName,
		Images:      make([]string, 0, len(ad.Media)),
	}
	if ad.Neighborhood != "" {
		if out.Location != "" {
			out.Location += "، " + ad.Neighborhood
		} else {
			out.Location = ad.Neighborhood
		}
	}
	if ad.PriceAmount != nil {
		out.Price = formatAmount(*ad.PriceAmount) + " " + t.Search.Currency
	} else {
		out.Price = t.Search.Negotiable
	}
	if ad.PublishedAt != nil {
		if ts, err := time.Parse(time.RFC3339, *ad.PublishedAt); err == nil {
			out.PublishedAt = ts.Format("2006-01-02")
		}
	}
	for _, m := range ad.Media {
		if u := media.PublicURL(mediaCDN, m.URL); u != "" {
			out.Images = append(out.Images, u)
		}
	}
	return out
}

func toSearchAd(ad searchclient.Ad, t i18n.Messages, mediaCDN string) viewmodel.SearchAd {
	out := viewmodel.SearchAd{
		ID:        ad.ID,
		Title:     ad.Title,
		Location:  ad.CityName,
		Thumbnail: media.PublicURL(mediaCDN, ad.Thumbnail),
		HasPhoto:  ad.HasPhoto,
	}
	if ad.ID > 0 {
		out.Href = "/ad/" + strconv.FormatInt(ad.ID, 10)
	}
	if ad.Neighborhood != "" {
		if out.Location != "" {
			out.Location += "، " + ad.Neighborhood
		} else {
			out.Location = ad.Neighborhood
		}
	}
	if ad.PriceAmount != nil {
		out.Price = formatAmount(*ad.PriceAmount) + " " + t.Search.Currency
	} else {
		out.Price = t.Search.Negotiable
	}
	if ad.PublishedAt != nil {
		if ts, err := time.Parse(time.RFC3339, *ad.PublishedAt); err == nil {
			out.PublishedAt = ts.Format("2006-01-02")
		}
	}
	return out
}

func searchURL(query, category string, page int) string {
	q := url.Values{}
	if query != "" {
		q.Set("q", query)
	}
	if category != "" {
		q.Set("category", category)
	}
	if page > 1 {
		q.Set("page", strconv.Itoa(page))
	}
	if encoded := q.Encode(); encoded != "" {
		return "/query-ads?" + encoded
	}
	return "/query-ads"
}

// formatAmount groups digits: 2500000 -> "2,500,000".
func formatAmount(v int64) string {
	s := strconv.FormatInt(v, 10)
	neg := strings.HasPrefix(s, "-")
	if neg {
		s = s[1:]
	}
	var b strings.Builder
	for i, r := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(r)
	}
	if neg {
		return "-" + b.String()
	}
	return b.String()
}
