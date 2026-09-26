package handler

import (
	"net/http"
	"strconv"

	"ads-platform/internal/business/iam/user_role/errorcode"
	"ads-platform/internal/business/iam/user_role/service"
	"ads-platform/internal/core/exception"
	"ads-platform/internal/core/middleware"

	"github.com/gin-gonic/gin"
)

type UserRoleHandler struct {
	svc service.UserRoleService
}

func NewUserRoleHandler(svc service.UserRoleService) *UserRoleHandler {
	return &UserRoleHandler{svc: svc}
}

func (h *UserRoleHandler) ListRoles(c *gin.Context) {
	userID, err := parseID(c.Param("userId"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	items, err := h.svc.ListRoles(c.Request.Context(), userID)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *UserRoleHandler) Add(c *gin.Context) {
	userID, err := parseID(c.Param("userId"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	roleID, err := parseID(c.Param("roleId"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	items, err := h.svc.Add(c.Request.Context(), userID, roleID)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *UserRoleHandler) Remove(c *gin.Context) {
	userID, err := parseID(c.Param("userId"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	roleID, err := parseID(c.Param("roleId"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	if err := h.svc.Remove(c.Request.Context(), userID, roleID); err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *UserRoleHandler) ListActions(c *gin.Context) {
	userID, err := parseID(c.Param("userId"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	items, err := h.svc.ListActions(c.Request.Context(), userID)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.JSON(http.StatusOK, items)
}

func parseID(raw string) (int64, error) {
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, exception.NewAppError(errorcode.ErrInvalidID.Code, errorcode.ErrInvalidID.HttpStatus, raw)
	}
	return id, nil
}
