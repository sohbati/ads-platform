package handler

import (
	"net/http"
	"strconv"

	"ads-platform/internal/business/userprofile/errorcode"
	"ads-platform/internal/business/userprofile/service"
	"ads-platform/internal/core/exception"
	"ads-platform/internal/core/middleware"

	"github.com/gin-gonic/gin"
)

type UserProfileHandler struct {
	service service.UserProfileService
}

func NewUserProfileHandler(svc service.UserProfileService) *UserProfileHandler {
	return &UserProfileHandler{service: svc}
}

type putProfileRequest struct {
	LocationSlugs []string `json:"location_slugs"`
}

func (h *UserProfileHandler) Get(c *gin.Context) {
	userID, err := parseUserID(c.Param("userId"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}

	profile, err := h.service.Get(c.Request.Context(), userID)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.JSON(http.StatusOK, profile)
}

func (h *UserProfileHandler) Put(c *gin.Context) {
	userID, err := parseUserID(c.Param("userId"))
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

	profile, err := h.service.Put(c.Request.Context(), userID, req.LocationSlugs)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.JSON(http.StatusOK, profile)
}

func parseUserID(raw string) (int64, error) {
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, exception.NewAppError(errorcode.ErrInvalidUser.Code, errorcode.ErrInvalidUser.HttpStatus, raw)
	}
	return id, nil
}
