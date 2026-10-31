package handler

import (
	"net/http"
	"strconv"

	"ads-platform/internal/business/iam/staff_auth/errorcode"
	"ads-platform/internal/business/iam/staff_auth/service"
	"ads-platform/internal/core/exception"
	"ads-platform/internal/core/middleware"

	"github.com/gin-gonic/gin"
)

type StaffAuthHandler struct {
	svc service.StaffAuthService
}

func NewStaffAuthHandler(svc service.StaffAuthService) *StaffAuthHandler {
	return &StaffAuthHandler{svc: svc}
}

type loginRequest struct {
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
}

func (h *StaffAuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, exception.NewAppError(errorcode.ErrInvalidBody.Code, errorcode.ErrInvalidBody.HttpStatus), 0)
		return
	}
	session, err := h.svc.Login(c.Request.Context(), req.Mobile, req.Password)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.JSON(http.StatusOK, session)
}

func (h *StaffAuthHandler) Session(c *gin.Context) {
	staffID, err := strconv.ParseInt(c.Param("staffId"), 10, 64)
	if err != nil {
		staffID = 0
	}
	session, err := h.svc.Session(c.Request.Context(), staffID)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.JSON(http.StatusOK, session)
}
