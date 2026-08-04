package handler

import (
	"net/http"

	"ads-platform-ui/internal/business/location/service"

	"github.com/gin-gonic/gin"
)

type APIHandler struct {
	service service.LocationService
}

func NewAPIHandler(svc service.LocationService) *APIHandler {
	return &APIHandler{service: svc}
}

// ListCities proxies GET /api/cities from ads-platform-cdn.
func (h *APIHandler) ListCities(c *gin.Context) {
	items, err := h.service.ListCities(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to load cities from CDN"})
		return
	}
	c.JSON(http.StatusOK, items)
}
