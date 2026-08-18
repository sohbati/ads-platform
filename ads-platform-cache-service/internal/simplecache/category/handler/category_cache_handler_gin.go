package handler

import (
	"net/http"
	"strconv"
	"strings"

	"cache-service/internal/core/exception"
	"cache-service/internal/simplecache/category/errorcode"
	"cache-service/internal/simplecache/category/service"

	"github.com/gin-gonic/gin"
)

type CategoryCacheHandler struct {
	service service.CategoryCacheService
}

func NewCategoryCacheHandler(svc service.CategoryCacheService) *CategoryCacheHandler {
	return &CategoryCacheHandler{service: svc}
}

// GetBySlugs returns categories for the given slugs.
// GET /api/v1/caches/categories/by-slugs?slugs=digital,cars
func (h *CategoryCacheHandler) GetBySlugs(c *gin.Context) {
	slugs := parseCSV(c.Query("slugs"))
	categories, err := h.service.GetBySlugs(c.Request.Context(), slugs)
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, categories)
}

// GetByIDs resolves id→slug then slug→category.
// GET /api/v1/caches/categories/by-ids?ids=3,13
func (h *CategoryCacheHandler) GetByIDs(c *gin.Context) {
	ids, err := parseIDs(c.Query("ids"))
	if err != nil {
		_ = c.Error(exception.NewAppError(
			errorcode.ErrIDsEmpty.Code,
			errorcode.ErrIDsEmpty.HttpStatus,
		).WithCause(err))
		c.Abort()
		return
	}

	categories, err := h.service.GetByIDs(c.Request.Context(), ids)
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, categories)
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
