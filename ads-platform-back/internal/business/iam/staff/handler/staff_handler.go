package handler

import (
	"net/http"
	"strconv"

	"ads-platform/internal/business/iam/staff/errorcode"
	"ads-platform/internal/business/iam/staff/service"
	"ads-platform/internal/core/exception"
	"ads-platform/internal/core/middleware"

	"github.com/gin-gonic/gin"
)

type StaffHandler struct {
	staffService service.StaffService
}

func NewStaffHandler(staffService service.StaffService) *StaffHandler {
	return &StaffHandler{staffService: staffService}
}

type staffRequest struct {
	Name       string `json:"name"`
	Mobile     string `json:"mobile"`
	NationalId string `json:"national_id"`
	Password   string `json:"password"`
	IsActive   *bool  `json:"is_active"`
}

func (r staffRequest) input() service.StaffInput {
	return service.StaffInput{
		Name:       r.Name,
		Mobile:     r.Mobile,
		NationalId: r.NationalId,
		Password:   r.Password,
		IsActive:   r.IsActive,
	}
}

type passwordRequest struct {
	Password string `json:"password"`
}

func (h *StaffHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	result, err := h.staffService.List(c.Request.Context(), c.Query("search"), page, pageSize)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *StaffHandler) Create(c *gin.Context) {
	var req staffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, exception.NewAppError(errorcode.ErrInvalidBody.Code, errorcode.ErrInvalidBody.HttpStatus), 0)
		return
	}
	item, err := h.staffService.Create(c.Request.Context(), req.input())
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *StaffHandler) Get(c *gin.Context) {
	id, err := parseID(c.Param("staffId"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	item, err := h.staffService.GetStaffByID(c.Request.Context(), id)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *StaffHandler) Update(c *gin.Context) {
	id, err := parseID(c.Param("staffId"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	var req staffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, exception.NewAppError(errorcode.ErrInvalidBody.Code, errorcode.ErrInvalidBody.HttpStatus), 0)
		return
	}
	item, err := h.staffService.Update(c.Request.Context(), id, req.input())
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *StaffHandler) Delete(c *gin.Context) {
	id, err := parseID(c.Param("staffId"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	if err := h.staffService.Delete(c.Request.Context(), id); err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *StaffHandler) SetPassword(c *gin.Context) {
	id, err := parseID(c.Param("staffId"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	var req passwordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, exception.NewAppError(errorcode.ErrInvalidBody.Code, errorcode.ErrInvalidBody.HttpStatus), 0)
		return
	}
	if err := h.staffService.SetPassword(c.Request.Context(), id, req.Password); err != nil {
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
