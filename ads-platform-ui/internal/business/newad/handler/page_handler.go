package handler

import (
	"ads-platform-ui/internal/core/cities"
	"ads-platform-ui/internal/core/config"
	"ads-platform-ui/internal/core/i18n"
	"ads-platform-ui/internal/web/page"

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

func (h *PageHandler) Index(c *gin.Context) {
	heading := h.i18n.MessagesFor(i18n.FromContext(c)).Nav.NewAd
	page.HTML(c, h.config, h.i18n, h.cities, "new_ad", h.config.AppName+" — "+heading, heading)
}
