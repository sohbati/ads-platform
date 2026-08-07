package router

import (
	"net/http"
	appContainer "ads-platform-notification/internal/core/container"

	"github.com/gin-gonic/gin"
	"ads-platform-notification/internal/core/middleware"
)

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

type Router struct {
	container *appContainer.AppContainer
}

func NewRouter(c *appContainer.AppContainer) *Router {
	return &Router{container: c}
}

func (r *Router) SetupRoutes() *gin.Engine {
	router := gin.New()

	router.Use(middleware.CustomRecovery())
	router.Use(middleware.GlobalErrorHandler())
	router.Use(corsMiddleware())

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":      "ROUTE_NOT_DEFINED",
			"statusCode": http.StatusNotFound,
		})
	})

	router.GET("/health", func(c *gin.Context) {
		status := "ok"
		natsStatus := "connected"
		if !r.container.Nats.IsConnected() {
			status = "degraded"
			natsStatus = "disconnected"
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  status,
			"service": "ads-platform-notification",
			"nats":    natsStatus,
		})
	})

	return router
}
