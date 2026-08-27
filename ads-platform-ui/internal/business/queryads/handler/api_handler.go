package handler

import (
	"net/http"
	"strconv"

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
