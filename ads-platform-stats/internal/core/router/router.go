package router

import (
	"net/http"

	appContainer "ads-platform-stats/internal/core/container"
	"ads-platform-stats/internal/core/middleware"

	"github.com/gin-gonic/gin"
)

type Router struct {
	container *appContainer.StatsContainer
	natsOK    func() bool
}

func NewRouter(c *appContainer.StatsContainer, natsOK func() bool) *Router {
	return &Router{container: c, natsOK: natsOK}
}

func (r *Router) SetupRoutes() *gin.Engine {
	engine := gin.New()
	engine.Use(middleware.CustomRecovery())
	engine.Use(middleware.GlobalErrorHandler())

	engine.GET("/health", func(c *gin.Context) {
		status := "ok"
		natsStatus := "connected"
		if r.natsOK != nil && !r.natsOK() {
			status = "degraded"
			natsStatus = "disconnected"
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  status,
			"service": "ads-platform-stats",
			"nats":    natsStatus,
		})
	})
	return engine
}
