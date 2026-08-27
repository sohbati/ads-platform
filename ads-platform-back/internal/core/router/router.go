package router

import (
	appContainer "ads-platform/internal/core/container"
	"net/http"

	"ads-platform/internal/core/middleware"

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
		// User routes
		users := api.Group("/users")
		{
			users.GET("/mobile/:mobile", r.container.User.UserHandler.GetUserByMobile)
			users.POST("/register-by-mobile/:mobile", r.container.User.UserHandler.RegisterByMobile)
			users.GET("/:userId/profile", r.container.UserProfile.Handler.Get)
			users.PUT("/:userId/profile", r.container.UserProfile.Handler.Put)
			users.GET("/:userId/ads", r.container.Ads.AdHandler.ListByUser)
		}

		// OTP routes
		otp := api.Group("/otp")
		{
			otp.POST("/:mobile/send", r.container.Otp.OtpHandler.SendOTP)
			otp.POST("/:mobile/verify", r.container.Otp.OtpHandler.VerifyOTP)
		}

		// Search routes
		api.GET("/q/:place/:category", r.container.Search.SearchHandler.SearchAds)

		api.POST("/ads", r.container.Ads.AdHandler.Create)
	}

	return router
}
