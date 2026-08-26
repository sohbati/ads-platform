package i18n

import (
	"strings"

	"ads-platform-ui/internal/core/cities"

	"github.com/gin-gonic/gin"
)

const ContextLocaleKey = "locale"

// Middleware resolves city from the request, sets locale from the city catalog, and stores both on context.
// Priority: ?city= → locations cookie → city cookie → /s/:slug path → default city from config.
func Middleware(reg *Registry, catalog *cities.Catalog, defaultCity string) gin.HandlerFunc {
	defaultCity = catalog.NormalizeSlug(defaultCity)
	if defaultCity == "" {
		defaultCity = "tehran"
	}

	return func(c *gin.Context) {
		slugs, primary := resolveLocations(c, catalog, defaultCity)
		if q := c.Query("city"); q != "" {
			if parsed := catalog.ParseLocationSlugs(q); len(parsed) > 0 {
				setLocationCookies(c, parsed, catalog.PrimaryCitySlug(parsed, defaultCity))
			}
		}

		loc := LocaleForCity(catalog, primary)
		c.Set(cities.ContextKey, primary)
		c.Set(cities.ContextLocationsKey, slugs)
		c.Set(ContextLocaleKey, loc)
		c.Next()
	}
}

func resolveLocations(c *gin.Context, catalog *cities.Catalog, fallback string) (slugs []string, primary string) {
	if q := c.Query("city"); q != "" {
		if parsed := catalog.ParseLocationSlugs(q); len(parsed) > 0 {
			return parsed, catalog.PrimaryCitySlug(parsed, fallback)
		}
	}

	if cookie, err := c.Cookie(cities.LocationsCookieName); err == nil {
		if parsed := catalog.ParseLocationSlugs(cookie); len(parsed) > 0 {
			return parsed, catalog.PrimaryCitySlug(parsed, fallback)
		}
	}

	if cookie, err := c.Cookie(cities.CookieName); err == nil {
		if parsed := catalog.ParseLocationSlugs(cookie); len(parsed) > 0 {
			return parsed, catalog.PrimaryCitySlug(parsed, fallback)
		}
	}

	if slug := cityFromPath(c.Request.URL.Path); slug != "" {
		if parsed := catalog.ParseLocationSlugs(slug); len(parsed) > 0 {
			return parsed, catalog.PrimaryCitySlug(parsed, fallback)
		}
	}

	parsed := catalog.ParseLocationSlugs(fallback)
	if len(parsed) == 0 {
		parsed = []string{fallback}
	}
	return parsed, catalog.PrimaryCitySlug(parsed, fallback)
}

func cityFromPath(path string) string {
	path = strings.Trim(path, "/")
	if strings.HasPrefix(path, "s/") {
		parts := strings.SplitN(path, "/", 2)
		if len(parts) == 2 {
			return parts[1]
		}
	}
	return ""
}

func setLocationCookies(c *gin.Context, slugs []string, primary string) {
	c.SetCookie(cities.CookieName, primary, 365*24*3600, "/", "", false, false)
	c.SetCookie(cities.LocationsCookieName, strings.Join(slugs, ","), 365*24*3600, "/", "", false, false)
}

// FromContext reads the locale derived from the active city.
func FromContext(c *gin.Context) Locale {
	if v, ok := c.Get(ContextLocaleKey); ok {
		if loc, ok := v.(Locale); ok {
			return loc
		}
	}
	return DefaultLocale
}

// CityFromContext reads the primary city slug.
func CityFromContext(c *gin.Context) string {
	if v, ok := c.Get(cities.ContextKey); ok {
		if slug, ok := v.(string); ok {
			return slug
		}
	}
	return "tehran"
}

// LocationsFromContext reads the selected location slugs (provinces and/or cities).
func LocationsFromContext(c *gin.Context) []string {
	if v, ok := c.Get(cities.ContextLocationsKey); ok {
		if slugs, ok := v.([]string); ok && len(slugs) > 0 {
			return slugs
		}
	}
	return []string{CityFromContext(c)}
}
