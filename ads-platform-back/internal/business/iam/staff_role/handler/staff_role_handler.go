package handler

import (
	"net/http"
	"strconv"

	"ads-platform/internal/business/iam/staff_role/errorcode"
	"ads-platform/internal/business/iam/staff_role/service"
	"ads-platform/internal/core/exception"
	"ads-platform/internal/core/middleware"

	"github.com/gin-gonic/gin"
)

type StaffRoleHandler struct {
	svc service.StaffRoleService
}

func NewStaffRoleHandler(svc service.StaffRoleService) *StaffRoleHandler {
	return &StaffRoleHandler{svc: svc}
}

func (h *StaffRoleHandler) ListRoles(c *gin.Context) {
	staffID, err := parseID(c.Param("staffId"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	items, err := h.svc.ListRoles(c.Request.Context(), staffID)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *StaffRoleHandler) Add(c *gin.Context) {
	staffID, err := parseID(c.Param("staffId"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	roleID, err := parseID(c.Param("roleId"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	items, err := h.svc.Add(c.Request.Context(), staffID, roleID)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *StaffRoleHandler) Remove(c *gin.Context) {
	staffID, err := parseID(c.Param("staffId"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	roleID, err := parseID(c.Param("roleId"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	if err := h.svc.Remove(c.Request.Context(), staffID, roleID); err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *StaffRoleHandler) ListActions(c *gin.Context) {
	staffID, err := parseID(c.Param("staffId"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	items, err := h.svc.ListActions(c.Request.Context(), staffID)
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
