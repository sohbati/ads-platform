package page

import (
	"net/http"

	"ads-platform-ui/internal/core/cities"
	"ads-platform-ui/internal/core/config"
	"ads-platform-ui/internal/core/i18n"

	"github.com/gin-gonic/gin"
)

// Base builds shared i18n page data for any HTML route.
func Base(c *gin.Context, cfg *config.Config, reg *i18n.Registry, catalog *cities.Catalog, title, heading string) i18n.Page {
	loc := i18n.FromContext(c)
	city := i18n.CityFromContext(c)
	p := i18n.BuildPage(reg, catalog, loc, cfg.AppName, city, c.Request.URL.Path)
	p.Title = title
	p.Heading = heading
	return p
}

// HTML renders a named template with a base page model.
func HTML(c *gin.Context, cfg *config.Config, reg *i18n.Registry, catalog *cities.Catalog, templateName, title, heading string) {
	c.HTML(http.StatusOK, templateName, Base(c, cfg, reg, catalog, title, heading))
}
