package handler

import (
	"ads-platform/internal/business/user/service"
	"ads-platform/internal/core/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserHandler handles HTTP requests for user-related operations
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// CheckUserExists handles GET /users/check/:user_id
// Checks if a user exists for the given user ID
func (h *UserHandler) GetUserByMobile(c *gin.Context) {
	mobile := c.Param("mobile")

	user, err := h.userService.GetUserByMobile(c.Request.Context(), mobile)
	if err != nil {
		middleware.HandleError(c, err, http.StatusMethodNotAllowed)
		return
	}

	c.JSON(http.StatusOK, user)
}
