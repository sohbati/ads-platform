package router

import (
	"net/http"
	appContainer "nats-message-broker/internal/core/container"
	"nats-message-broker/internal/core/middleware"
	"nats-message-broker/internal/core/natsserver"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
		natsServerStatus := "running"
		natsClientStatus := "connected"

		if !r.container.NatsServer.IsRunning() {
			status = "degraded"
			natsServerStatus = "stopped"
		}
		if !r.container.Nats.IsConnected() {
			status = "degraded"
			natsClientStatus = "disconnected"
		}

		c.JSON(http.StatusOK, gin.H{
			"status":       status,
			"service":      "nats-message-broker",
			"natsServer":   natsServerStatus,
			"natsClient":   natsClientStatus,
			"natsPort":     r.container.NatsServer.Port(),
			"natsPanelURL": natsserver.MonitorPath + "/",
		})
	})

	r.mountNatsMonitor(router)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api/v1")
	{
		events := api.Group("/events")
		{
			events.POST("/:subject", r.container.Event.EventHandler.Publish)
		}
	}

	return router
}

func (r *Router) mountNatsMonitor(router *gin.Engine) {
	handler := r.container.NatsServer.HTTPHandler()
	if handler == nil {
		return
	}

	monitor := http.StripPrefix(natsserver.MonitorPath, handler)

	router.GET(natsserver.MonitorPath, func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, natsserver.MonitorPath+"/")
	})
	router.GET("/nat", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, natsserver.MonitorPath+"/")
	})
	router.Any(natsserver.MonitorPath+"/*path", gin.WrapH(monitor))
}
