package router

import (
	"message-broker/internal/core/brokerserver"
	appContainer "message-broker/internal/core/container"
	"message-broker/internal/core/middleware"
	"net/http"

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
		brokerServerStatus := "running"
		brokerClientStatus := "connected"

		if !r.container.BrokerServer.IsRunning() {
			status = "degraded"
			brokerServerStatus = "stopped"
		}
		if !r.container.Broker.IsConnected() {
			status = "degraded"
			brokerClientStatus = "disconnected"
		}

		c.JSON(http.StatusOK, gin.H{
			"status":         status,
			"service":        "message-broker",
			"brokerServer":   brokerServerStatus,
			"brokerClient":   brokerClientStatus,
			"brokerPort":     r.container.BrokerServer.Port(),
			"brokerURL":      r.container.BrokerServer.ClientURL(),
			"brokerPanelURL": brokerserver.MonitorPath + "/",
		})
	})

	r.mountBrokerMonitor(router)

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

func (r *Router) mountBrokerMonitor(router *gin.Engine) {
	handler := r.container.BrokerServer.HTTPHandler()
	if handler == nil {
		return
	}

	monitor := http.StripPrefix(brokerserver.MonitorPath, handler)

	router.GET(brokerserver.MonitorPath, func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, brokerserver.MonitorPath+"/")
	})
	router.Any(brokerserver.MonitorPath+"/*path", gin.WrapH(monitor))
}
