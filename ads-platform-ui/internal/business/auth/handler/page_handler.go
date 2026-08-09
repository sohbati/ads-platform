package handler

import (
	"net/http"

	"ads-platform-ui/internal/core/cities"
	"ads-platform-ui/internal/core/config"
	"ads-platform-ui/internal/core/i18n"

	"github.com/gin-gonic/gin"
)

type PageHandler struct {
	config *config.Config
	i18n   *i18n.Registry
	cities *cities.Catalog
}

func NewPageHandler(cfg *config.Config, reg *i18n.Registry, catalog *cities.Catalog) *PageHandler {
	return &PageHandler{config: cfg, i18n: reg, cities: catalog}
}

func (h *PageHandler) Login(c *gin.Context) {
	next := c.Query("next")
	if next == "" {
		next = "/my-info"
	}

	pageData := i18n.BuildPage(h.i18n, h.cities, i18n.FromContext(c), h.config.AppName, i18n.CityFromContext(c), c.Request.URL.Path)
	pageData.Title = h.config.AppName + " — " + pageData.T.Auth.LoginTitle
	pageData.Heading = pageData.T.Auth.LoginTitle
	pageData.NextPath = next
	c.HTML(http.StatusOK, "login", pageData)
}
