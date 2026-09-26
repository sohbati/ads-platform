package handler

import (
	"net/http"
	"strconv"

	"ads-platform/internal/business/iam/action/errorcode"
	"ads-platform/internal/business/iam/action/service"
	"ads-platform/internal/core/exception"
	"ads-platform/internal/core/middleware"

	"github.com/gin-gonic/gin"
)

type ActionHandler struct {
	svc service.ActionService
}

func NewActionHandler(svc service.ActionService) *ActionHandler {
	return &ActionHandler{svc: svc}
}

type namedRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type renameRequest struct {
	Name string `json:"name"`
}

func (h *ActionHandler) Create(c *gin.Context) {
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

func (h *ActionHandler) List(c *gin.Context) {
	items, err := h.svc.List(c.Request.Context())
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *ActionHandler) Get(c *gin.Context) {
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

func (h *ActionHandler) Update(c *gin.Context) {
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

func (h *ActionHandler) Delete(c *gin.Context) {
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

func parseID(raw string) (int64, error) {
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, exception.NewAppError(errorcode.ErrInvalidID.Code, errorcode.ErrInvalidID.HttpStatus, raw)
	}
	return id, nil
}
