package handler

import (
	"net/http"

	"ads-platform/internal/business/otp/model"
	"ads-platform/internal/business/otp/service"
	"ads-platform/internal/core/middleware"

	"github.com/gin-gonic/gin"
)

type OtpHandler struct {
	otpService service.OtpService
}

func NewOtpHandler(otpService service.OtpService) *OtpHandler {
	return &OtpHandler{otpService: otpService}
}

func (h *OtpHandler) SendOTP(c *gin.Context) {
	mobile := c.Param("mobile")

	resp, err := h.otpService.SendOTP(c.Request.Context(), mobile)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *OtpHandler) VerifyOTP(c *gin.Context) {
	mobile := c.Param("mobile")

	var req model.VerifyOtpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, err, http.StatusBadRequest)
		return
	}

	resp, err := h.otpService.VerifyOTP(c.Request.Context(), mobile, req.Otp)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}

	c.JSON(http.StatusOK, resp)
}
