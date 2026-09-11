package handler

import (
	"net/http"
	"strconv"

	"ads-platform-ui/internal/business/location/service"

	"github.com/gin-gonic/gin"
)

type GeoHandler struct {
	geo service.GeoService
}

func NewGeoHandler(geo service.GeoService) *GeoHandler {
	return &GeoHandler{geo: geo}
}

func (h *GeoHandler) CityCenter(c *gin.Context) {
	pt, err := h.geo.CityCenter(c.Request.Context(), c.Query("slug"), c.Query("name"), langFrom(c))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "GEO_UNAVAILABLE"})
		return
	}
	c.JSON(http.StatusOK, pt)
}

func (h *GeoHandler) Reverse(c *gin.Context) {
	lat, err1 := strconv.ParseFloat(c.Query("lat"), 64)
	lng, err2 := strconv.ParseFloat(c.Query("lng"), 64)
	if err1 != nil || err2 != nil || lat < -90 || lat > 90 || lng < -180 || lng > 180 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_COORDINATES"})
		return
	}
	result, err := h.geo.Reverse(c.Request.Context(), lat, lng, langFrom(c))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "GEO_UNAVAILABLE"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func langFrom(c *gin.Context) string {
	lang := c.Query("lang")
	if lang == "" {
		lang = "fa"
	}
	return lang
}
