package i18n

import (
	"fmt"
	"strings"

	"ads-platform-ui/internal/core/cities"
)

// Page carries locale, direction, and translated strings for templates.
type Page struct {
	Title              string
	Heading            string
	Locale             Locale
	Lang               string
	Dir                string
	AppName            string
	CitySlug           string
	CityDisplayName    string
	LocationSlugs      []string
	LocationSlugsCSV   string
	CurrentPath        string
	NextPath           string
	SessionUserName    string
	SessionUserMobile  string
	IsAuthenticated    bool
	DefaultCountryCode string
	// SearchQuery pre-fills the header search input on result pages.
	SearchQuery string
	T           Messages
}

// BuildPage creates shared template data for HTML pages.
func BuildPage(reg *Registry, catalog *cities.Catalog, loc Locale, appName, citySlug, currentPath string, locationSlugs []string) Page {
	b := reg.Get(loc)
	slug := catalog.NormalizeSlug(citySlug)
	slugs := locationSlugs
	if len(slugs) == 0 && slug != "" {
		slugs = []string{slug}
	}
	return Page{
		Locale:           loc,
		Lang:             loc.Lang(),
		Dir:              b.Dir,
		AppName:          appName,
		CitySlug:         slug,
		CityDisplayName:  LocationHeaderLabel(reg, catalog, loc, slugs, slug),
		LocationSlugs:    slugs,
		LocationSlugsCSV: strings.Join(slugs, ","),
		CurrentPath:      currentPath,
		T:                b.Messages,
	}
}

// Formatf formats a translated pattern: {{ format .T.Hero.Title .CityDisplayName }}
func Formatf(pattern string, args ...interface{}) string {
	return fmt.Sprintf(pattern, args...)
}
