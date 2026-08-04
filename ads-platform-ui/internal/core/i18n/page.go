package i18n

import (
	"fmt"

	"ads-platform-ui/internal/core/cities"
)

// Page carries locale, direction, and translated strings for templates.
type Page struct {
	Title           string
	Heading         string
	Locale          Locale
	Lang            string
	Dir             string
	AppName         string
	CitySlug        string
	CityDisplayName string
	CurrentPath     string
	T               Messages
}

// BuildPage creates shared template data for HTML pages.
func BuildPage(reg *Registry, catalog *cities.Catalog, loc Locale, appName, citySlug, currentPath string) Page {
	b := reg.Get(loc)
	slug := catalog.NormalizeSlug(citySlug)
	return Page{
		Locale:          loc,
		Lang:            loc.Lang(),
		Dir:             b.Dir,
		AppName:         appName,
		CitySlug:        slug,
		CityDisplayName: CityDisplayName(reg, catalog, loc, slug),
		CurrentPath:     currentPath,
		T:               b.Messages,
	}
}

// Formatf formats a translated pattern: {{ format .T.Hero.Title .CityDisplayName }}
func Formatf(pattern string, args ...interface{}) string {
	return fmt.Sprintf(pattern, args...)
}
