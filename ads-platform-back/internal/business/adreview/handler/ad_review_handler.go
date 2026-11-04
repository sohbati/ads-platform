package handler

import (
	"context"
	"net/http"
	"strconv"

	"ads-platform/internal/business/adreview/errorcode"
	"ads-platform/internal/business/adreview/model"
	"ads-platform/internal/business/adreview/service"
	"ads-platform/internal/core/exception"
	"ads-platform/internal/core/middleware"
	"ads-platform/internal/core/panelctx"

	"github.com/gin-gonic/gin"
)

type AdReviewHandler struct {
	svc service.AdReviewService
}

func NewAdReviewHandler(svc service.AdReviewService) *AdReviewHandler {
	return &AdReviewHandler{svc: svc}
}

type rejectRequest struct {
	CustomerReason string `json:"customer_reason"`
	InternalReason string `json:"internal_reason"`
}

func (h *AdReviewHandler) List(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	categoryID, _ := strconv.Atoi(c.Query("category_id"))
	cityID, _ := strconv.Atoi(c.Query("city_id"))
	result, err := h.svc.List(c.Request.Context(), model.ReviewFilter{
		Status:     c.Query("status"),
		Search:     c.Query("search"),
		CategoryID: categoryID,
		CityID:     cityID,
		Cursor:     c.Query("cursor"),
		PageSize:   pageSize,
	})
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *AdReviewHandler) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	ad, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.JSON(http.StatusOK, ad)
}

func (h *AdReviewHandler) Approve(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	staffID, _ := panelctx.StaffID(c)
	ad, err := h.svc.Approve(c.Request.Context(), id, staffID)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.JSON(http.StatusOK, ad)
}

func (h *AdReviewHandler) Reject(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req rejectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, exception.NewAppError(errorcode.ErrInvalidBody.Code, errorcode.ErrInvalidBody.HttpStatus), 0)
		return
	}
	staffID, _ := panelctx.StaffID(c)
	ad, err := h.svc.Reject(c.Request.Context(), id, staffID, service.RejectInput{
		CustomerReason: req.CustomerReason,
		InternalReason: req.InternalReason,
	})
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.JSON(http.StatusOK, ad)
}

// Next answers 204 when no ad is waiting.
func (h *AdReviewHandler) Next(c *gin.Context) {
	staffID, _ := panelctx.StaffID(c)
	ad, err := h.svc.Next(c.Request.Context(), staffID)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	if ad == nil {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, ad)
}

func (h *AdReviewHandler) Claim(c *gin.Context) {
	h.claimAction(c, h.svc.Claim)
}

func (h *AdReviewHandler) Extend(c *gin.Context) {
	h.claimAction(c, h.svc.Extend)
}

func (h *AdReviewHandler) claimAction(c *gin.Context, action func(ctx context.Context, id, staffID int64) (*model.ReviewAd, error)) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	staffID, _ := panelctx.StaffID(c)
	ad, err := action(c.Request.Context(), id, staffID)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.JSON(http.StatusOK, ad)
}

func (h *AdReviewHandler) Release(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	staffID, _ := panelctx.StaffID(c)
	respondEmpty(c, h.svc.Release(c.Request.Context(), id, staffID))
}

func (h *AdReviewHandler) ReleaseAll(c *gin.Context) {
	staffID, _ := panelctx.StaffID(c)
	respondEmpty(c, h.svc.ReleaseAll(c.Request.Context(), staffID))
}

func (h *AdReviewHandler) ForceRelease(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	respondEmpty(c, h.svc.ForceRelease(c.Request.Context(), id))
}

func (h *AdReviewHandler) TodayStats(c *gin.Context) {
	stats, err := h.svc.TodayStats(c.Request.Context())
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": stats})
}

func respondEmpty(c *gin.Context, err error) {
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.Status(http.StatusNoContent)
}

func parseID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("adId"), 10, 64)
	if err != nil || id <= 0 {
		middleware.HandleError(c, exception.NewAppError(errorcode.ErrInvalidID.Code, errorcode.ErrInvalidID.HttpStatus, c.Param("adId")), 0)
		return 0, false
	}
	return id, true
}
