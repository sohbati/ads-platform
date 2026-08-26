package handler

import (
	"net/http"
	"strconv"
	"strings"

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
	locations := i18n.LocationsFromContext(c)

	query := strings.TrimSpace(c.Query("q"))
	category := strings.TrimSpace(c.Query("category"))
	if query != "" || category != "" {
		pageNum, _ := strconv.Atoi(c.Query("page"))
		page := h.service.BuildSearchPage(c.Request.Context(), loc, h.config.AppName, city, c.Request.URL.Path, locations, service.SearchParams{
			Query:    query,
			Category: category,
			Page:     pageNum,
		})
		page.Page.DefaultCountryCode = h.config.DefaultCountryCode
		c.HTML(http.StatusOK, "query_ads", page)
		return
	}

	page := h.service.BuildPage(loc, h.config.AppName, city, c.Request.URL.Path, locations)
	page.Page.DefaultCountryCode = h.config.DefaultCountryCode
	c.HTML(http.StatusOK, "query_ads", page)
}
