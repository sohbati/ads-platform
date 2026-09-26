package handler

import (
	"net/http"
	"strconv"

	"ads-platform/internal/business/iam/role_action/errorcode"
	"ads-platform/internal/business/iam/role_action/service"
	"ads-platform/internal/core/exception"
	"ads-platform/internal/core/middleware"

	"github.com/gin-gonic/gin"
)

type RoleActionHandler struct {
	svc service.RoleActionService
}

func NewRoleActionHandler(svc service.RoleActionService) *RoleActionHandler {
	return &RoleActionHandler{svc: svc}
}

type setActionsRequest struct {
	ActionIDs []int64 `json:"action_ids"`
}

func (h *RoleActionHandler) List(c *gin.Context) {
	roleID, err := parseID(c.Param("id"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	items, err := h.svc.List(c.Request.Context(), roleID)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *RoleActionHandler) Set(c *gin.Context) {
	roleID, err := parseID(c.Param("id"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	var req setActionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, exception.NewAppError(errorcode.ErrInvalidBody.Code, errorcode.ErrInvalidBody.HttpStatus), 0)
		return
	}
	items, err := h.svc.Set(c.Request.Context(), roleID, req.ActionIDs)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *RoleActionHandler) Add(c *gin.Context) {
	roleID, err := parseID(c.Param("id"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	actionID, err := parseID(c.Param("actionId"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	items, err := h.svc.Add(c.Request.Context(), roleID, actionID)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *RoleActionHandler) Remove(c *gin.Context) {
	roleID, err := parseID(c.Param("id"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	actionID, err := parseID(c.Param("actionId"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	if err := h.svc.Remove(c.Request.Context(), roleID, actionID); err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.Status(http.StatusNoContent)
}

func parseID(raw string) (int64, error) {
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, exception.NewAppError(errorcode.ErrInvalidID.Code, errorcode.ErrInvalidID.HttpStatus, raw)
	}
	return id, nil
}
