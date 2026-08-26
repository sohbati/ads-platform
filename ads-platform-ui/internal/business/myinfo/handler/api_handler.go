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

func (h *APIHandler) GetProfile(c *gin.Context) {
	result, err := h.bff.Get(c.Request.Context(), "/api/v1/me/profile", bff.RequestCookies(c.Request))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "BACKEND_UNAVAILABLE", "statusCode": http.StatusBadGateway})
		return
	}
	bff.ForwardResponse(c.Writer, result)
}

func (h *APIHandler) PutProfile(c *gin.Context) {
	payload, err := io.ReadAll(io.LimitReader(c.Request.Body, 1<<16))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_REQUEST", "statusCode": http.StatusBadRequest})
		return
	}
	result, err := h.bff.Do(c.Request.Context(), http.MethodPut, "/api/v1/me/profile", payload, "application/json", bff.RequestCookies(c.Request))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "BACKEND_UNAVAILABLE", "statusCode": http.StatusBadGateway})
		return
	}
	bff.ForwardResponse(c.Writer, result)
}
