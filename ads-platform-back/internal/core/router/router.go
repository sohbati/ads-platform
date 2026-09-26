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

	// Business API routes.
	// Public: catalog + identity. Self: owner resources, no IAM. Panel: CRM, IAM only.
	api := router.Group("/api/v1")
	{
		users := api.Group("/users")
		{
			users.GET("/mobile/:mobile", r.container.User.UserHandler.GetUserByMobile)
			users.POST("/register-by-mobile/:mobile", r.container.User.UserHandler.RegisterByMobile)
		}

		otp := api.Group("/otp")
		{
			otp.POST("/:mobile/send", r.container.Otp.OtpHandler.SendOTP)
			otp.POST("/:mobile/verify", r.container.Otp.OtpHandler.VerifyOTP)
		}

		api.GET("/q/:place/:category", r.container.Search.SearchHandler.SearchAds)
		api.GET("/media/status", r.container.Ads.AdHandler.MediaStatus)
		api.GET("/ads/:id", r.container.Ads.AdHandler.GetPublic)
		api.GET("/ads/:id/contact", r.container.Ads.AdHandler.GetPublicContact)

		self := api.Group("/self")
		{
			self.POST("/ads", r.container.Ads.AdHandler.Create)
			self.GET("/:userId/profile", r.container.UserProfile.Handler.Get)
			self.PUT("/:userId/profile", r.container.UserProfile.Handler.Put)
			self.GET("/:userId/ads", r.container.Ads.AdHandler.ListByUser)
			self.GET("/:userId/ad-stats", r.container.Ads.AdHandler.ListStats)
			self.GET("/:userId/ads/:adId", r.container.Ads.AdHandler.GetForOwner)
			self.PUT("/:userId/ads/:adId", r.container.Ads.AdHandler.Update)
		}

		panel := api.Group("/panel")
		{
			panelUsers := panel.Group("/users")
			{
				panelUsers.GET("/:userId/roles", r.container.UserRole.Handler.ListRoles)
				panelUsers.POST("/:userId/roles/:roleId", r.container.UserRole.Handler.Add)
				panelUsers.DELETE("/:userId/roles/:roleId", r.container.UserRole.Handler.Remove)
				panelUsers.GET("/:userId/groups", r.container.UserGroup.Handler.ListByUser)
				panelUsers.GET("/:userId/actions", r.container.UserRole.Handler.ListActions)
			}

			actions := panel.Group("/actions")
			{
				actions.GET("", r.container.Action.Handler.List)
				actions.POST("", r.container.Action.Handler.Create)
				actions.GET("/:id", r.container.Action.Handler.Get)
				actions.PUT("/:id", r.container.Action.Handler.Update)
				actions.DELETE("/:id", r.container.Action.Handler.Delete)
			}

			roles := panel.Group("/roles")
			{
				roles.GET("", r.container.Role.Handler.List)
				roles.POST("", r.container.Role.Handler.Create)
				roles.GET("/:id", r.container.Role.Handler.Get)
				roles.PUT("/:id", r.container.Role.Handler.Update)
				roles.DELETE("/:id", r.container.Role.Handler.Delete)
				roles.GET("/:id/actions", r.container.RoleAction.Handler.List)
				roles.PUT("/:id/actions", r.container.RoleAction.Handler.Set)
				roles.POST("/:id/actions/:actionId", r.container.RoleAction.Handler.Add)
				roles.DELETE("/:id/actions/:actionId", r.container.RoleAction.Handler.Remove)
			}

			groups := panel.Group("/user-groups")
			{
				groups.GET("", r.container.UserGroup.Handler.List)
				groups.POST("", r.container.UserGroup.Handler.Create)
				groups.GET("/:id", r.container.UserGroup.Handler.Get)
				groups.PUT("/:id", r.container.UserGroup.Handler.Update)
				groups.DELETE("/:id", r.container.UserGroup.Handler.Delete)
				groups.GET("/:id/users", r.container.UserGroup.Handler.ListUsers)
				groups.POST("/:id/users/:userId", r.container.UserGroup.Handler.AddUser)
				groups.DELETE("/:id/users/:userId", r.container.UserGroup.Handler.RemoveUser)
			}
		}
	}

	return router
}
