package handler

import (
	"net/http"
	"strings"

	"ads-platform/internal/business/otp/errorcode"
	"ads-platform/internal/business/otp/model"
	"ads-platform/internal/business/otp/service"
	"ads-platform/internal/core/exception"
	"ads-platform/internal/core/mobile"
	"ads-platform/internal/core/middleware"

	"github.com/gin-gonic/gin"
)

type OtpHandler struct {
	otpService service.OtpService
	mobileNorm *mobile.Normalizer
}

func NewOtpHandler(otpService service.OtpService, mobileNorm *mobile.Normalizer) *OtpHandler {
	return &OtpHandler{otpService: otpService, mobileNorm: mobileNorm}
}

func (h *OtpHandler) SendOTP(c *gin.Context) {
	mobileNumber, err := h.normalizeMobile(c.Param("mobile"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}

	resp, err := h.otpService.SendOTP(c.Request.Context(), mobileNumber)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *OtpHandler) VerifyOTP(c *gin.Context) {
	mobileNumber, err := h.normalizeMobile(c.Param("mobile"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}

	var req model.VerifyOtpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, err, http.StatusBadRequest)
		return
	}

	resp, err := h.otpService.VerifyOTP(c.Request.Context(), mobileNumber, req.Otp)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *OtpHandler) normalizeMobile(raw string) (string, error) {
	if strings.TrimSpace(raw) == "" {
		return "", exception.NewAppError(
			errorcode.ErrMobileEmpty.Code, errorcode.ErrMobileEmpty.HttpStatus, raw)
	}

	normalized, err := h.mobileNorm.Normalize(raw)
	if err != nil {
		return "", exception.NewAppError(
			errorcode.ErrInvalidMobile.Code, errorcode.ErrInvalidMobile.HttpStatus, raw)
	}
	return normalized, nil
}
