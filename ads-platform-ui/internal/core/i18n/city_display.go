package i18n

import "ads-platform-ui/internal/core/cities"

// CityDisplayName returns the localized city label: bundle override, else cities.json name.
func CityDisplayName(reg *Registry, catalog *cities.Catalog, loc Locale, slug string) string {
	slug = catalog.NormalizeSlug(slug)
	msg := reg.MessagesFor(loc)
	if name, ok := msg.CityNames[slug]; ok && name != "" {
		return name
	}
	return catalog.Name(slug)
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
