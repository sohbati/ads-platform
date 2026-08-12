package handler

import (
	"net/http"

	"ads-platform-ui/internal/business/queryads/service"
	"ads-platform-ui/internal/core/config"
	"ads-platform-ui/internal/core/i18n"

	"github.com/gin-gonic/gin"
)

type PageHandler struct {
	service service.QueryAdsService
	config  *config.Config
}

func NewPageHandler(svc service.QueryAdsService, cfg *config.Config) *PageHandler {
	return &PageHandler{service: svc, config: cfg}
}

func (h *PageHandler) Index(c *gin.Context) {
	loc := i18n.FromContext(c)
	city := i18n.CityFromContext(c)
	page := h.service.BuildPage(loc, h.config.AppName, city, c.Request.URL.Path)
	page.Page.DefaultCountryCode = h.config.DefaultCountryCode
	c.HTML(http.StatusOK, "query_ads", page)
}
