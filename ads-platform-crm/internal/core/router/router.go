package router

import (
	"net/http"

	appContainer "ads-platform-crm/internal/core/container"
	"ads-platform-crm/internal/core/middleware"

	"github.com/gin-gonic/gin"
)

type Router struct {
	container *appContainer.AppContainer
}

func NewRouter(c *appContainer.AppContainer) *Router {
	return &Router{container: c}
}

func (r *Router) SetupRoutes() *gin.Engine {
	engine := gin.New()
	engine.Use(middleware.CustomRecovery())
	engine.Use(middleware.GlobalErrorHandler())

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "ads-platform-crm",
			"app":     r.container.Config.AppName,
		})
	})

	engine.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":      "ROUTE_NOT_DEFINED",
			"statusCode": http.StatusNotFound,
		})
	})

	return engine
}
