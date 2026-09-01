package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"ads-bff/internal/business/auth/service"
	"ads-bff/internal/business/stats/publisher"
	cacheclient "ads-bff/internal/core/client/cache"
	"ads-bff/internal/core/config"
	"ads-bff/internal/core/exception"
	"ads-bff/internal/core/middleware"
	"ads-bff/internal/core/ratelimit"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const maxEventBody = 8 << 10

type ingestRequest struct {
	AdID       int64  `json:"ad_id"`
	Event      string `json:"event"`
	ViewerID   string `json:"viewer_id"`
	OccurredAt string `json:"occurred_at"`
}

type outboundEvent struct {
	AdID          int64     `json:"ad_id"`
	Event         string    `json:"event"`
	ViewerID      string    `json:"viewer_id"`
	OccurredAt    time.Time `json:"occurred_at"`
	SessionUserID *int64    `json:"session_user_id,omitempty"`
}

type StatsHandler struct {
	cfg     *config.Config
	auth    service.AuthService
	pub     publisher.Publisher
	limiter *ratelimit.Limiter
}

func NewStatsHandler(cfg *config.Config, auth service.AuthService, pub publisher.Publisher, limiter *ratelimit.Limiter) *StatsHandler {
	if pub == nil {
		pub = publisher.Noop()
	}
	if limiter == nil {
		limiter = ratelimit.New(time.Minute, 60)
	}
	return &StatsHandler{cfg: cfg, auth: auth, pub: pub, limiter: limiter}
}

func (h *StatsHandler) Ingest(c *gin.Context) {
	if !h.limiter.Allow(c.ClientIP()) {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "RATE_LIMITED", "statusCode": http.StatusTooManyRequests})
		return
	}

	raw, err := io.ReadAll(io.LimitReader(c.Request.Body, maxEventBody+1))
	if err != nil || len(raw) == 0 || len(raw) > maxEventBody {
		middleware.HandleError(c, exception.NewAppError("INVALID_REQUEST", http.StatusBadRequest), 0)
		return
	}

	var req ingestRequest
	if err := json.Unmarshal(raw, &req); err != nil {
		middleware.HandleError(c, exception.NewAppError("INVALID_REQUEST", http.StatusBadRequest), 0)
		return
	}

	if req.AdID <= 0 || !validEvent(req.Event) {
		middleware.HandleError(c, exception.NewAppError("INVALID_REQUEST", http.StatusBadRequest), 0)
		return
	}
	if _, err := uuid.Parse(req.ViewerID); err != nil {
		middleware.HandleError(c, exception.NewAppError("INVALID_REQUEST", http.StatusBadRequest), 0)
		return
	}
	occurredAt, err := time.Parse(time.RFC3339Nano, req.OccurredAt)
	if err != nil {
		occurredAt, err = time.Parse(time.RFC3339, req.OccurredAt)
	}
	if err != nil || occurredAt.IsZero() {
		middleware.HandleError(c, exception.NewAppError("INVALID_REQUEST", http.StatusBadRequest), 0)
		return
	}

	payload, err := json.Marshal(outboundEvent{
		AdID:          req.AdID,
		Event:         req.Event,
		ViewerID:      req.ViewerID,
		OccurredAt:    occurredAt.UTC(),
		SessionUserID: h.optionalSessionUserID(c),
	})
	if err != nil {
		middleware.HandleError(c, exception.NewAppError("INVALID_REQUEST", http.StatusBadRequest), 0)
		return
	}

	_ = h.pub.Publish(c.Request.Context(), payload)
	c.Status(http.StatusNoContent)
}

func (h *StatsHandler) optionalSessionUserID(c *gin.Context) *int64 {
	if h.auth == nil || h.cfg == nil {
		return nil
	}
	sessionID, err := c.Cookie(h.cfg.SessionCookieName)
	if err != nil || sessionID == "" {
		return nil
	}
	user, err := h.auth.GetCurrentUser(c.Request.Context(), sessionID)
	if err != nil {
		if cacheclient.IsSessionNotFound(err) {
			return nil
		}
		return nil
	}
	if user == nil || user.ID <= 0 {
		return nil
	}
	id := user.ID
	return &id
}

func validEvent(name string) bool {
	switch name {
	case "view", "contact_reveal", "call":
		return true
	default:
		return false
	}
}
