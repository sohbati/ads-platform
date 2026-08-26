package i18n

import (
	"fmt"
	"strconv"
	"strings"

	"ads-platform-ui/internal/core/cities"
)

// CityDisplayName returns the localized city label: bundle override, else cities.json name.
func CityDisplayName(reg *Registry, catalog *cities.Catalog, loc Locale, slug string) string {
	slug = catalog.NormalizeSlug(slug)
	msg := reg.MessagesFor(loc)
	if name, ok := msg.CityNames[slug]; ok && name != "" {
		return name
	}
	return catalog.Name(slug)
}

// LocationHeaderLabel is the topbar city control: one name, or "N cities" when several are selected.
func LocationHeaderLabel(reg *Registry, catalog *cities.Catalog, loc Locale, slugs []string, fallbackSlug string) string {
	if len(slugs) <= 1 {
		slug := fallbackSlug
		if len(slugs) == 1 {
			slug = slugs[0]
		}
		return CityDisplayName(reg, catalog, loc, slug)
	}
	pattern := reg.MessagesFor(loc).Location.HeaderCount
	if pattern == "" {
		return withLocaleDigits(loc, strconv.Itoa(len(slugs)))
	}
	return withLocaleDigits(loc, fmt.Sprintf(pattern, len(slugs)))
}

func withLocaleDigits(loc Locale, s string) string {
	if loc != FA {
		return s
	}
	return strings.NewReplacer(
		"0", "۰", "1", "۱", "2", "۲", "3", "۳", "4", "۴",
		"5", "۵", "6", "۶", "7", "۷", "8", "۸", "9", "۹",
	).Replace(s)
}

// LocaleForCity resolves a locale from the city catalog.
func LocaleForCity(catalog *cities.Catalog, slug string) Locale {
	code := catalog.LocaleCode(slug)
	loc, ok := Parse(code)
	if ok {
		return loc
	}
	return DefaultLocale
}
