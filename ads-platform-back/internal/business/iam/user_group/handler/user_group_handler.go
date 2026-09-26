package handler

import (
	"net/http"
	"strconv"

	"ads-platform/internal/business/iam/user_group/errorcode"
	"ads-platform/internal/business/iam/user_group/service"
	"ads-platform/internal/core/exception"
	"ads-platform/internal/core/middleware"

	"github.com/gin-gonic/gin"
)

type UserGroupHandler struct {
	svc service.UserGroupService
}

func NewUserGroupHandler(svc service.UserGroupService) *UserGroupHandler {
	return &UserGroupHandler{svc: svc}
}

type namedRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type renameRequest struct {
	Name string `json:"name"`
}

func (h *UserGroupHandler) Create(c *gin.Context) {
	var req namedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, exception.NewAppError(errorcode.ErrInvalidBody.Code, errorcode.ErrInvalidBody.HttpStatus), 0)
		return
	}
	item, err := h.svc.Create(c.Request.Context(), req.Code, req.Name)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *UserGroupHandler) List(c *gin.Context) {
	items, err := h.svc.List(c.Request.Context())
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *UserGroupHandler) Get(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	item, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *UserGroupHandler) Update(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	var req renameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, exception.NewAppError(errorcode.ErrInvalidBody.Code, errorcode.ErrInvalidBody.HttpStatus), 0)
		return
	}
	item, err := h.svc.Update(c.Request.Context(), id, req.Name)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *UserGroupHandler) Delete(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *UserGroupHandler) ListUsers(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	ids, err := h.svc.ListUsers(c.Request.Context(), id)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.JSON(http.StatusOK, gin.H{"user_ids": ids})
}

func (h *UserGroupHandler) AddUser(c *gin.Context) {
	groupID, err := parseID(c.Param("id"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	userID, err := parseID(c.Param("userId"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	if err := h.svc.AddUser(c.Request.Context(), groupID, userID); err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *UserGroupHandler) RemoveUser(c *gin.Context) {
	groupID, err := parseID(c.Param("id"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	userID, err := parseID(c.Param("userId"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	if err := h.svc.RemoveUser(c.Request.Context(), groupID, userID); err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *UserGroupHandler) ListByUser(c *gin.Context) {
	userID, err := parseID(c.Param("userId"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	items, err := h.svc.ListByUser(c.Request.Context(), userID)
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
