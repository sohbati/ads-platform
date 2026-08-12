package handler

import (
	"ads-platform/internal/business/user/errorcode"
	"ads-platform/internal/business/user/service"
	"ads-platform/internal/core/exception"
	"ads-platform/internal/core/mobile"
	"ads-platform/internal/core/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserHandler handles HTTP requests for user-related operations
type UserHandler struct {
	userService service.UserService
	mobileNorm  *mobile.Normalizer
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService service.UserService, mobileNorm *mobile.Normalizer) *UserHandler {
	return &UserHandler{
		userService: userService,
		mobileNorm:  mobileNorm,
	}
}

func (h *UserHandler) GetUserByMobile(c *gin.Context) {
	mobileNumber, err := h.normalizeMobile(c.Param("mobile"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}

	user, err := h.userService.GetUserByMobile(c.Request.Context(), mobileNumber)
	if err != nil {
		middleware.HandleError(c, err, http.StatusMethodNotAllowed)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) RegisterByMobile(c *gin.Context) {
	mobileNumber, err := h.normalizeMobile(c.Param("mobile"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}

	user, err := h.userService.RegisterByMobile(c.Request.Context(), mobileNumber)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) normalizeMobile(raw string) (string, error) {
	normalized, err := h.mobileNorm.Normalize(raw)
	if err != nil {
		return "", exception.NewAppError(
			errorcode.ErrInvalidMobile.Code, errorcode.ErrInvalidMobile.HttpStatus, raw)
	}
	return normalized, nil
}
