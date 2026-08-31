package handler

import (
	"net/http"
	"strconv"
	"strings"

	"ads-platform-ui/internal/core/bff"
	"ads-platform-ui/internal/core/i18n"

	"github.com/gin-gonic/gin"
)

func (h *PageHandler) SearchJSON(c *gin.Context) {
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
	if page.Search == nil {
		c.JSON(http.StatusOK, gin.H{"ads": []any{}, "page": 1, "total_pages": 0, "has_more": false})
		return
	}
	c.JSON(http.StatusOK, page.Search)
}

func (h *PageHandler) ContactJSON(c *gin.Context) {
	if h.bff == nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "BACKEND_UNAVAILABLE"})
		return
	}
	id := strings.TrimSpace(c.Param("id"))
	result, err := h.bff.Get(c.Request.Context(), "/api/v1/ads/"+id+"/contact", bff.RequestCookies(c.Request))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "BACKEND_UNAVAILABLE"})
		return
	}
	bff.ForwardResponse(c.Writer, result)
}
