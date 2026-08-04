package impl

import (
	"ads-platform-ui/internal/business/queryads/viewmodel"
	"ads-platform-ui/internal/core/cities"
	"ads-platform-ui/internal/core/i18n"
	"ads-platform-ui/internal/domain"
)

type QueryAdsServiceImpl struct {
	i18n   *i18n.Registry
	cities *cities.Catalog
}

func NewQueryAdsService(reg *i18n.Registry, catalog *cities.Catalog) *QueryAdsServiceImpl {
	return &QueryAdsServiceImpl{i18n: reg, cities: catalog}
}

func (s *QueryAdsServiceImpl) BuildPage(loc i18n.Locale, appName, citySlug, currentPath string) viewmodel.QueryAdsPage {
	page := i18n.BuildPage(s.i18n, s.cities, loc, appName, citySlug, currentPath)
	page.Title = appName + " — " + s.i18n.MessagesFor(loc).Nav.QueryAds
	return viewmodel.QueryAdsPage{
		Page:          page,
		Categories:    s.categories(loc),
		PopularCities: s.popularCities(loc),
	}
}

func (s *QueryAdsServiceImpl) categories(loc i18n.Locale) []domain.Category {
	t := s.i18n.MessagesFor(loc)
	base := []struct {
		id, icon, slug string
	}{
		{"real-estate", "home", "real-estate"},
		{"vehicles", "car", "vehicles"},
		{"digital", "device", "digital"},
		{"home", "sofa", "home"},
		{"services", "wrench", "services"},
		{"jobs", "briefcase", "jobs"},
		{"personal", "shirt", "personal"},
		{"leisure", "ball", "leisure"},
	}
	out := make([]domain.Category, 0, len(base))
	for _, c := range base {
		item := t.CategoryItems[c.id]
		out = append(out, domain.Category{
			ID:          c.id,
			Name:        item.Name,
			Description: item.Description,
			Icon:        c.icon,
			Slug:        c.slug,
			Href:        "/query-ads?category=" + c.slug,
		})
	}
	return out
}

func (s *QueryAdsServiceImpl) popularCities(loc i18n.Locale) []domain.City {
	records := s.cities.PopularRecords(8)
	out := make([]domain.City, 0, len(records))
	for _, r := range records {
		slug := r.Slug
		out = append(out, domain.City{
			ID:   slug,
			Name: i18n.CityDisplayName(s.i18n, s.cities, loc, slug),
			Slug: slug,
			Href: "/query-ads?city=" + slug,
		})
	}
	return out
}
