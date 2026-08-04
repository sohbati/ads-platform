package router

import (
	"net/http"

	appContainer "ads-platform-cdn/internal/core/container"

	"github.com/gin-gonic/gin"
)

type Router struct {
	container *appContainer.AppContainer
}

func NewRouter(c *appContainer.AppContainer) *Router {
	return &Router{container: c}
}

func (r *Router) SetupRoutes() *gin.Engine {
	router := gin.Default()

	cfg := r.container.Config
	router.Static("/json", cfg.StaticDir+"/json")
	router.Static("/static", cfg.StaticDir)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "ads-platform-cdn"})
	})

	api := router.Group("/api")
	{
		api.GET("/categories", r.container.Category.APIHandler.List)
		api.GET("/cities", r.container.City.APIHandler.List)
	}

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":      "ROUTE_NOT_DEFINED",
			"statusCode": http.StatusNotFound,
		})
	})

	return router
}
