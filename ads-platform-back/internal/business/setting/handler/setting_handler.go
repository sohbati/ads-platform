package handler

import (
	"net/http"

	"ads-platform/internal/business/setting/errorcode"
	"ads-platform/internal/business/setting/service"
	"ads-platform/internal/core/exception"
	"ads-platform/internal/core/middleware"
	"ads-platform/internal/core/panelctx"

	"github.com/gin-gonic/gin"
)

type SettingHandler struct {
	svc service.SettingService
}

func NewSettingHandler(svc service.SettingService) *SettingHandler {
	return &SettingHandler{svc: svc}
}

type updateRequest struct {
	Value string `json:"value"`
}

func (h *SettingHandler) List(c *gin.Context) {
	items, err := h.svc.List(c.Request.Context())
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *SettingHandler) Update(c *gin.Context) {
	var req updateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, exception.NewAppError(errorcode.ErrInvalidBody.Code, errorcode.ErrInvalidBody.HttpStatus), 0)
		return
	}
	staffID, _ := panelctx.StaffID(c)
	item, err := h.svc.Update(c.Request.Context(), c.Param("key"), req.Value, staffID)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.JSON(http.StatusOK, item)
}
