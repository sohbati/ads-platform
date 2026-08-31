package handler

import (
	"net/http"
	"strconv"
	"strings"

	"ads-platform-ui/internal/business/queryads/service"
	"ads-platform-ui/internal/core/bff"
	"ads-platform-ui/internal/core/config"
	"ads-platform-ui/internal/core/i18n"

	"github.com/gin-gonic/gin"
)

type PageHandler struct {
	service service.QueryAdsService
	config  *config.Config
	bff     *bff.Client
}

func NewPageHandler(svc service.QueryAdsService, cfg *config.Config, bffClient *bff.Client) *PageHandler {
	return &PageHandler{service: svc, config: cfg, bff: bffClient}
}

func (h *PageHandler) Index(c *gin.Context) {
	pageNum, _ := strconv.Atoi(c.Query("page"))
	page := h.service.BuildSearchPage(
		c.Request.Context(),
		i18n.FromContext(c),
		h.config.AppName,
		i18n.CityFromContext(c),
		c.Request.URL.Path,
		i18n.LocationsFromContext(c),
		searchParams(c, pageNum),
	)
	page.Page.DefaultCountryCode = h.config.DefaultCountryCode
	c.HTML(http.StatusOK, "query_ads", page)
}

func (h *PageHandler) Show(c *gin.Context) {
	adID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	page := h.service.BuildDetailPage(
		c.Request.Context(),
		i18n.FromContext(c),
		h.config.AppName,
		i18n.CityFromContext(c),
		c.Request.URL.Path,
		i18n.LocationsFromContext(c),
		adID,
	)
	page.Page.DefaultCountryCode = h.config.DefaultCountryCode
	status := http.StatusOK
	if page.NotFound {
		status = http.StatusNotFound
	}
	c.HTML(status, "ad_detail", page)
}

func searchParams(c *gin.Context, pageNum int) service.SearchParams {
	return service.SearchParams{
		Query:    strings.TrimSpace(c.Query("q")),
		Category: strings.TrimSpace(c.Query("category")),
		Page:     pageNum,
	}
}
