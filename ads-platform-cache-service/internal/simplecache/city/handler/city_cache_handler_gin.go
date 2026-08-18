package handler

import (
	"net/http"
	"strconv"
	"strings"

	"cache-service/internal/core/exception"
	"cache-service/internal/simplecache/city/errorcode"
	"cache-service/internal/simplecache/city/service"

	"github.com/gin-gonic/gin"
)

type CityCacheHandler struct {
	service service.CityCacheService
}

func NewCityCacheHandler(svc service.CityCacheService) *CityCacheHandler {
	return &CityCacheHandler{service: svc}
}

// GetByIDs returns cities for the given ids.
// GET /api/v1/caches/cities/by-ids?ids=1,869
func (h *CityCacheHandler) GetByIDs(c *gin.Context) {
	ids, err := parseIDs(c.Query("ids"))
	if err != nil {
		_ = c.Error(exception.NewAppError(
			errorcode.ErrFailToReadBody.Code,
			errorcode.ErrFailToReadBody.HttpStatus,
		).WithCause(err))
		c.Abort()
		return
	}

	cities, err := h.service.GetByIDs(c.Request.Context(), ids)
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, cities)
}

// GetBySlugs resolves slug→id then id→city.
// GET /api/v1/caches/cities/by-slugs?slugs=tehran,abyek
func (h *CityCacheHandler) GetBySlugs(c *gin.Context) {
	slugs := parseCSV(c.Query("slugs"))
	cities, err := h.service.GetBySlugs(c.Request.Context(), slugs)
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, cities)
}

func parseCSV(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func parseIDs(raw string) ([]int, error) {
	parts := parseCSV(raw)
	ids := make([]int, 0, len(parts))
	for _, p := range parts {
		id, err := strconv.Atoi(p)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
