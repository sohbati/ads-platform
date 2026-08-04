package handler

import (
	"net/http"

	"ads-platform-ui/internal/business/category/service"

	"github.com/gin-gonic/gin"
)

type APIHandler struct {
	service service.CategoryService
}

func NewAPIHandler(svc service.CategoryService) *APIHandler {
	return &APIHandler{service: svc}
}

// List proxies GET /api/categories from ads-platform-cdn.
func (h *APIHandler) List(c *gin.Context) {
	items, err := h.service.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to load categories from CDN"})
		return
	}
	c.JSON(http.StatusOK, items)
}
