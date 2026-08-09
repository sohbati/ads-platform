package handler

import (
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

func (h *APIHandler) GetByMobile(c *gin.Context) {
	mobile := c.Param("mobile")
	path := "/api/v1/users/mobile/" + mobile

	status, body, err := h.bff.Get(c.Request.Context(), path)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to reach backend via BFF"})
		return
	}

	c.Data(status, "application/json", body)
}
