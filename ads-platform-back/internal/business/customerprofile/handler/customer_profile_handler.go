package handler

import (
	"net/http"
	"strconv"

	"ads-platform/internal/business/customerprofile/errorcode"
	"ads-platform/internal/business/customerprofile/service"
	"ads-platform/internal/core/exception"
	"ads-platform/internal/core/middleware"

	"github.com/gin-gonic/gin"
)

type CustomerProfileHandler struct {
	service service.CustomerProfileService
}

func NewCustomerProfileHandler(svc service.CustomerProfileService) *CustomerProfileHandler {
	return &CustomerProfileHandler{service: svc}
}

type putProfileRequest struct {
	LocationSlugs []string `json:"location_slugs"`
}

func (h *CustomerProfileHandler) Get(c *gin.Context) {
	customerID, err := parseCustomerID(c.Param("customerId"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}

	profile, err := h.service.Get(c.Request.Context(), customerID)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.JSON(http.StatusOK, profile)
}

func (h *CustomerProfileHandler) Put(c *gin.Context) {
	customerID, err := parseCustomerID(c.Param("customerId"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}

	var req putProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, exception.NewAppError(
			errorcode.ErrInvalidLocation.Code, errorcode.ErrInvalidLocation.HttpStatus), 0)
		return
	}

	profile, err := h.service.Put(c.Request.Context(), customerID, req.LocationSlugs)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.JSON(http.StatusOK, profile)
}

func parseCustomerID(raw string) (int64, error) {
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, exception.NewAppError(errorcode.ErrInvalidCustomer.Code, errorcode.ErrInvalidCustomer.HttpStatus, raw)
	}
	return id, nil
}
