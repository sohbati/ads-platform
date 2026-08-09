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

func (h *APIHandler) SendOTP(c *gin.Context) {
	mobile := c.Param("mobile")
	path := "/api/v1/otp/" + mobile + "/send"

	status, body, err := h.bff.PostJSON(c.Request.Context(), path, nil)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to reach backend via BFF"})
		return
	}

	c.Data(status, "application/json", body)
}

func (h *APIHandler) VerifyOTP(c *gin.Context) {
	mobile := c.Param("mobile")
	path := "/api/v1/otp/" + mobile + "/verify"

	payload, err := io.ReadAll(io.LimitReader(c.Request.Body, 1<<20))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	status, body, err := h.bff.PostJSON(c.Request.Context(), path, payload)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to reach backend via BFF"})
		return
	}

	c.Data(status, "application/json", body)
}
