package handler

import (
	"net/http"
	"strings"

	"cache-service/internal/simplecache/attrschema/service"

	"github.com/gin-gonic/gin"
)

type AttrSchemaCacheHandler struct {
	service service.AttrSchemaCacheService
}

func NewAttrSchemaCacheHandler(svc service.AttrSchemaCacheService) *AttrSchemaCacheHandler {
	return &AttrSchemaCacheHandler{service: svc}
}

// GetByNames returns attr JSON Schema templates for the given names.
// GET /api/v1/caches/attr-schemas/by-names?names=cars,laptop
func (h *AttrSchemaCacheHandler) GetByNames(c *gin.Context) {
	names := parseCSV(c.Query("names"))
	items, err := h.service.GetByNames(c.Request.Context(), names)
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, items)
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
