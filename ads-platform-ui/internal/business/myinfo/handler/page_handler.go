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

func (h *PageHandler) render(c *gin.Context, templateName, heading string) {
	page.HTML(c, h.config, h.i18n, h.cities, templateName, h.config.AppName+" — "+heading, heading)
}

func (h *PageHandler) Index(c *gin.Context) {
	h.render(c, "myinfo", h.i18n.MessagesFor(i18n.FromContext(c)).Nav.MyInfo)
}

func (h *PageHandler) UserDetails(c *gin.Context) {
	h.render(c, "myinfo_user_details", h.i18n.MessagesFor(i18n.FromContext(c)).Nav.UserDetails)
}

func (h *PageHandler) UserAds(c *gin.Context) {
	h.render(c, "myinfo_user_ads", h.i18n.MessagesFor(i18n.FromContext(c)).Nav.UserAds)
}

func (h *PageHandler) MarkedAds(c *gin.Context) {
	h.render(c, "myinfo_marked_ads", h.i18n.MessagesFor(i18n.FromContext(c)).Nav.MarkedAds)
}

func (h *PageHandler) Setting(c *gin.Context) {
	h.render(c, "myinfo_setting", h.i18n.MessagesFor(i18n.FromContext(c)).Nav.Setting)
}
