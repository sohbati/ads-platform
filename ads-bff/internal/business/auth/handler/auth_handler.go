package handler

import (
	"net/http"
	"time"

	"ads-bff/internal/business/auth/model"
	"ads-bff/internal/business/auth/service"
	cacheclient "ads-bff/internal/core/client/cache"
	"ads-bff/internal/core/config"
	"ads-bff/internal/core/middleware"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	cfg     *config.Config
	service service.AuthService
}

func NewAuthHandler(cfg *config.Config, svc service.AuthService) *AuthHandler {
	return &AuthHandler{cfg: cfg, service: svc}
}

func (h *AuthHandler) SendOTP(c *gin.Context) {
	mobile := c.Param("mobile")
	status, body, err := h.service.SendOTP(c.Request.Context(), mobile)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "BACKEND_UNAVAILABLE", "statusCode": http.StatusBadGateway})
		return
	}
	c.Data(status, "application/json", body)
}

func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	mobile := c.Param("mobile")

	var req model.VerifyOtpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, err, http.StatusBadRequest)
		return
	}

	loginResp, sessionID, status, body, err := h.service.LoginWithOTP(c.Request.Context(), mobile, req.Otp)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "BACKEND_UNAVAILABLE", "statusCode": http.StatusBadGateway})
		return
	}
	if status != http.StatusOK {
		c.Data(status, "application/json", body)
		return
	}

	h.setSessionCookie(c, sessionID)
	c.JSON(http.StatusOK, loginResp)
}

func (h *AuthHandler) Me(c *gin.Context) {
	sessionID, err := c.Cookie(h.cfg.SessionCookieName)
	if err != nil || sessionID == "" {
		c.JSON(http.StatusOK, model.MeResponse{Authenticated: false})
		return
	}

	user, err := h.service.GetCurrentUser(c.Request.Context(), sessionID)
	if err != nil {
		if cacheclient.IsSessionNotFound(err) {
			h.clearSessionCookie(c)
			c.JSON(http.StatusOK, model.MeResponse{Authenticated: false})
			return
		}
		c.JSON(http.StatusBadGateway, gin.H{"error": "SESSION_LOOKUP_FAILED", "statusCode": http.StatusBadGateway})
		return
	}

	c.JSON(http.StatusOK, model.MeResponse{Authenticated: true, User: user})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	sessionID, _ := c.Cookie(h.cfg.SessionCookieName)
	_ = h.service.Logout(c.Request.Context(), sessionID)
	h.clearSessionCookie(c)
	c.JSON(http.StatusOK, gin.H{"message": "logged_out"})
}

func (h *AuthHandler) setSessionCookie(c *gin.Context, sessionID string) {
	maxAge := int(h.cfg.SessionTTL.Seconds())
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(h.cfg.SessionCookieName, sessionID, maxAge, "/", "", h.cfg.CookieSecure, true)
}

func (h *AuthHandler) clearSessionCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(h.cfg.SessionCookieName, "", -1, "/", "", h.cfg.CookieSecure, true)
}

// SessionCookieName exposes cookie name for middleware/tests.
func (h *AuthHandler) SessionCookieName() string {
	return h.cfg.SessionCookieName
}

// SessionMaxAge exposes TTL for proxy clients.
func (h *AuthHandler) SessionMaxAge() time.Duration {
	return h.cfg.SessionTTL
}
