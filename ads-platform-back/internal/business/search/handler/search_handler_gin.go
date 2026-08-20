package handler

import (
	"net/http"
	"strconv"
	"strings"

	"ads-platform/internal/business/search/service"
	"ads-platform/internal/core/middleware"

	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	searchService service.SearchService
}

func NewSearchHandler(searchService service.SearchService) *SearchHandler {
	return &SearchHandler{searchService: searchService}
}

// SearchAds handles:
//
//	GET /api/v1/q/:place/:category
//	GET /api/v1/q/tehran/electronic-devices
//	GET /api/v1/q/iran/electronic-devices?cities=1,869
//
// Query params: q, cities, price_min, price_max, has_photo, neighborhood, sort, page, limit
func (h *SearchHandler) SearchAds(c *gin.Context) {
	place := c.Param("place")
	category := c.Param("category")

	priceMin, err := parseOptionalInt64(c.Query("price_min"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "SEARCH_INVALID_PRICE_MIN", "statusCode": http.StatusBadRequest})
		return
	}
	priceMax, err := parseOptionalInt64(c.Query("price_max"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "SEARCH_INVALID_PRICE_MAX", "statusCode": http.StatusBadRequest})
		return
	}

	hasPhoto, err := parseOptionalBool(c.Query("has_photo"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "SEARCH_INVALID_HAS_PHOTO", "statusCode": http.StatusBadRequest})
		return
	}

	page := parseIntDefault(c.Query("page"), 1)
	limit := parseIntDefault(c.Query("limit"), 24)

	resp, err := h.searchService.Search(
		c.Request.Context(),
		place,
		category,
		c.Query("cities"),
		c.Query("q"),
		priceMin,
		priceMax,
		hasPhoto,
		c.Query("neighborhood"),
		c.DefaultQuery("sort", "newest"),
		page,
		limit,
	)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func parseOptionalInt64(raw string) (*int64, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	v, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func parseOptionalBool(raw string) (*bool, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	v, err := strconv.ParseBool(raw)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func parseIntDefault(raw string, def int) int {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return def
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return def
	}
	return v
}
