package router

import (
	"net/http"

	appContainer "ads-platform-stats/internal/core/container"
	"ads-platform-stats/internal/core/middleware"

	"github.com/gin-gonic/gin"
)

type Router struct {
	container *appContainer.StatsContainer
	brokerOK  func() bool
}

func NewRouter(c *appContainer.StatsContainer, brokerOK func() bool) *Router {
	return &Router{container: c, brokerOK: brokerOK}
}

func (r *Router) SetupRoutes() *gin.Engine {
	engine := gin.New()
	engine.Use(middleware.CustomRecovery())
	engine.Use(middleware.GlobalErrorHandler())

	engine.GET("/health", func(c *gin.Context) {
		status := "ok"
		brokerStatus := "connected"
		if r.brokerOK != nil && !r.brokerOK() {
			status = "degraded"
			brokerStatus = "disconnected"
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  status,
			"service": "ads-platform-stats",
			"broker":  brokerStatus,
		})
	})
	return engine
}
