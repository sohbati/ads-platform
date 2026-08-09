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
	path := "/api/v1/auth/otp/" + mobile + "/send"

	result, err := h.bff.PostJSON(c.Request.Context(), path, nil, bff.RequestCookies(c.Request))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to reach BFF"})
		return
	}
	bff.ForwardResponse(c.Writer, result)
}

func (h *APIHandler) VerifyOTP(c *gin.Context) {
	mobile := c.Param("mobile")
	path := "/api/v1/auth/otp/" + mobile + "/verify"

	payload, err := io.ReadAll(io.LimitReader(c.Request.Body, 1<<20))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	result, err := h.bff.PostJSON(c.Request.Context(), path, payload, bff.RequestCookies(c.Request))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to reach BFF"})
		return
	}
	bff.ForwardResponse(c.Writer, result)
}

func (h *APIHandler) Me(c *gin.Context) {
	result, err := h.bff.Get(c.Request.Context(), "/api/v1/auth/me", bff.RequestCookies(c.Request))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to reach BFF"})
		return
	}
	bff.ForwardResponse(c.Writer, result)
}

func (h *APIHandler) Logout(c *gin.Context) {
	result, err := h.bff.PostJSON(c.Request.Context(), "/api/v1/auth/logout", nil, bff.RequestCookies(c.Request))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to reach BFF"})
		return
	}
	bff.ForwardResponse(c.Writer, result)
}
