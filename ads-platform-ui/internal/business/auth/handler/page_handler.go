package handler

import (
	"net/http"
	"net/url"

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

	target := "/query-ads?login=1&next=" + url.QueryEscape(next)
	c.Redirect(http.StatusFound, target)
}
