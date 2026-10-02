package handler

import (
	"ads-platform/internal/business/customer/errorcode"
	"ads-platform/internal/business/customer/service"
	"ads-platform/internal/core/exception"
	"ads-platform/internal/core/middleware"
	"ads-platform/internal/core/mobile"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	customers  service.CustomerService
	mobileNorm *mobile.Normalizer
}

func NewCustomerHandler(customers service.CustomerService, mobileNorm *mobile.Normalizer) *CustomerHandler {
	return &CustomerHandler{
		customers:  customers,
		mobileNorm: mobileNorm,
	}
}

func (h *CustomerHandler) GetCustomerByMobile(c *gin.Context) {
	mobileNumber, err := h.normalizeMobile(c.Param("mobile"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}

	customer, err := h.customers.GetCustomerByMobile(c.Request.Context(), mobileNumber)
	if err != nil {
		middleware.HandleError(c, err, http.StatusMethodNotAllowed)
		return
	}

	c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) RegisterByMobile(c *gin.Context) {
	mobileNumber, err := h.normalizeMobile(c.Param("mobile"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}

	customer, err := h.customers.RegisterByMobile(c.Request.Context(), mobileNumber)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}

	c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) normalizeMobile(raw string) (string, error) {
	normalized, err := h.mobileNorm.Normalize(raw)
	if err != nil {
		return "", exception.NewAppError(
			errorcode.ErrInvalidMobile.Code, errorcode.ErrInvalidMobile.HttpStatus, raw)
	}
	return normalized, nil
}
