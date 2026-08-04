package router

import (
	appContainer "cache-service/internal/core/container"
	"net/http"

	"cache-service/internal/core/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// CORS middleware
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

	// Apply global error handler and recovery (similar to @ControllerAdvice in Spring Boot)
	router.Use(middleware.CustomRecovery())
	router.Use(middleware.GlobalErrorHandler())

	// Apply CORS middleware to all routes
	router.Use(corsMiddleware())

	// Handle 404: http.StatusNotFound routes
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":      "ROUTE_NOT_DEFINED",
			"statusCode": http.StatusNotFound,
		})
	})

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "ads-platform"})
	})

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Business API routes
	api := router.Group("/api/v1")
	{
		// Cache routes
		caches := api.Group("/caches")
		{
			caches.GET("/otp/:key", r.container.OtpCacheContainer.OtpCacheHandler.GetCacheByKey)
			caches.POST("/otp/:key", r.container.OtpCacheContainer.OtpCacheHandler.AddItem)

		}

	}

	return router
}
