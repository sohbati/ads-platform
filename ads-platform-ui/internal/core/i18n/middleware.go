package i18n

import (
	"strings"

	"ads-platform-ui/internal/core/cities"

	"github.com/gin-gonic/gin"
)

const ContextLocaleKey = "locale"

// Middleware resolves city from the request, sets locale from the city catalog, and stores both on context.
// Priority: ?city= → city cookie → /s/:slug path → default city from config.
func Middleware(reg *Registry, catalog *cities.Catalog, defaultCity string) gin.HandlerFunc {
	defaultCity = catalog.NormalizeSlug(defaultCity)
	if defaultCity == "" {
		defaultCity = "tehran"
	}

	return func(c *gin.Context) {
		city := resolveCity(c, catalog, defaultCity)
		if q := c.Query("city"); q != "" {
			if slug := catalog.NormalizeSlug(q); slug != "" {
				setCityCookie(c, slug)
			}
		}

		loc := LocaleForCity(catalog, city)
		c.Set(cities.ContextKey, city)
		c.Set(ContextLocaleKey, loc)
		c.Next()
	}
}

func resolveCity(c *gin.Context, catalog *cities.Catalog, fallback string) string {
	if q := c.Query("city"); q != "" {
		if slug := catalog.NormalizeSlug(q); slug != "" {
			return slug
		}
	}

	if cookie, err := c.Cookie(cities.CookieName); err == nil {
		if slug := catalog.NormalizeSlug(cookie); slug != "" {
			return slug
		}
	}

	if slug := cityFromPath(c.Request.URL.Path); slug != "" {
		return catalog.NormalizeSlug(slug)
	}

	return fallback
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

func setCityCookie(c *gin.Context, slug string) {
	c.SetCookie(cities.CookieName, slug, 365*24*3600, "/", "", false, false)
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

// CityFromContext reads the active city slug.
func CityFromContext(c *gin.Context) string {
	if v, ok := c.Get(cities.ContextKey); ok {
		if slug, ok := v.(string); ok {
			return slug
		}
	}
	return "tehran"
}
