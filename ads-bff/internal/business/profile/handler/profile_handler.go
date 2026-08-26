package handler

import (
	"io"
	"net/http"

	"ads-bff/internal/business/auth/service"
	backendclient "ads-bff/internal/core/client/backend"
	cacheclient "ads-bff/internal/core/client/cache"
	"ads-bff/internal/core/config"
	"ads-bff/internal/core/exception"
	"ads-bff/internal/core/middleware"

	"github.com/gin-gonic/gin"
)

const maxProfileBody = 1 << 16

type ProfileHandler struct {
	cfg     *config.Config
	auth    service.AuthService
	backend *backendclient.Client
}

func NewProfileHandler(cfg *config.Config, auth service.AuthService, backend *backendclient.Client) *ProfileHandler {
	return &ProfileHandler{cfg: cfg, auth: auth, backend: backend}
}

func (h *ProfileHandler) Get(c *gin.Context) {
	userID, err := h.sessionUserID(c)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}

	status, body, err := h.backend.GetUserProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "BACKEND_UNAVAILABLE", "statusCode": http.StatusBadGateway})
		return
	}
	c.Data(status, "application/json", body)
}

func (h *ProfileHandler) Put(c *gin.Context) {
	userID, err := h.sessionUserID(c)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}

	raw, err := io.ReadAll(io.LimitReader(c.Request.Body, maxProfileBody))
	if err != nil {
		middleware.HandleError(c, exception.NewAppError("INVALID_REQUEST", http.StatusBadRequest).WithCause(err), 0)
		return
	}

	status, body, err := h.backend.PutUserProfile(c.Request.Context(), userID, raw)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "BACKEND_UNAVAILABLE", "statusCode": http.StatusBadGateway})
		return
	}
	c.Data(status, "application/json", body)
}

func (h *ProfileHandler) sessionUserID(c *gin.Context) (int64, error) {
	sessionID, err := c.Cookie(h.cfg.SessionCookieName)
	if err != nil || sessionID == "" {
		return 0, exception.NewAppError("AUTH_REQUIRED", http.StatusUnauthorized)
	}

	user, err := h.auth.GetCurrentUser(c.Request.Context(), sessionID)
	if err != nil {
		if cacheclient.IsSessionNotFound(err) {
			return 0, exception.NewAppError("AUTH_REQUIRED", http.StatusUnauthorized)
		}
		return 0, exception.NewAppError("SESSION_LOOKUP_FAILED", http.StatusBadGateway).WithCause(err)
	}
	if user == nil || user.ID <= 0 {
		return 0, exception.NewAppError("AUTH_REQUIRED", http.StatusUnauthorized)
	}
	return user.ID, nil
}
