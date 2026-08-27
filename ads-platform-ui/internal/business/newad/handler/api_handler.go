package handler

import (
	"io"
	"net/http"

	"ads-platform-ui/internal/core/bff"

	"github.com/gin-gonic/gin"
)

type APIHandler struct {
	bff *bff.Client
}

func NewAPIHandler(client *bff.Client) *APIHandler {
	return &APIHandler{bff: client}
}

func (h *APIHandler) Create(c *gin.Context) {
	payload, err := io.ReadAll(io.LimitReader(c.Request.Body, 32<<20))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_REQUEST", "statusCode": http.StatusBadRequest})
		return
	}

	result, err := h.bff.Do(c.Request.Context(), http.MethodPost, "/api/v1/ads", payload, c.Request.Header.Get("Content-Type"), bff.RequestCookies(c.Request))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "BACKEND_UNAVAILABLE", "statusCode": http.StatusBadGateway})
		return
	}
	bff.ForwardResponse(c.Writer, result)
}

func (h *APIHandler) Update(c *gin.Context) {
	id := c.Param("id")
	payload, err := io.ReadAll(io.LimitReader(c.Request.Body, 32<<20))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_REQUEST", "statusCode": http.StatusBadRequest})
		return
	}

	result, err := h.bff.Do(c.Request.Context(), http.MethodPut, "/api/v1/me/ads/"+id, payload, c.Request.Header.Get("Content-Type"), bff.RequestCookies(c.Request))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "BACKEND_UNAVAILABLE", "statusCode": http.StatusBadGateway})
		return
	}
	bff.ForwardResponse(c.Writer, result)
}
