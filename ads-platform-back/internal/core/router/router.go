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
		customers := api.Group("/customers")
		{
			customers.GET("/mobile/:mobile", r.container.Customer.Handler.GetCustomerByMobile)
			customers.POST("/register-by-mobile/:mobile", r.container.Customer.Handler.RegisterByMobile)
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
			self.GET("/:customerId/profile", r.container.CustomerProfile.Handler.Get)
			self.PUT("/:customerId/profile", r.container.CustomerProfile.Handler.Put)
			self.GET("/:customerId/ads", r.container.Ads.AdHandler.ListByCustomer)
			self.GET("/:customerId/ad-stats", r.container.Ads.AdHandler.ListStats)
			self.GET("/:customerId/ads/:adId", r.container.Ads.AdHandler.GetForOwner)
			self.PUT("/:customerId/ads/:adId", r.container.Ads.AdHandler.Update)
		}

		panel := api.Group("/panel")
		{
			panelStaff := panel.Group("/staff")
			{
				panelStaff.GET("/:staffId/roles", r.container.StaffRole.Handler.ListRoles)
				panelStaff.POST("/:staffId/roles/:roleId", r.container.StaffRole.Handler.Add)
				panelStaff.DELETE("/:staffId/roles/:roleId", r.container.StaffRole.Handler.Remove)
				panelStaff.GET("/:staffId/groups", r.container.StaffGroup.Handler.ListByStaff)
				panelStaff.GET("/:staffId/actions", r.container.StaffRole.Handler.ListActions)
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

			groups := panel.Group("/staff-groups")
			{
				groups.GET("", r.container.StaffGroup.Handler.List)
				groups.POST("", r.container.StaffGroup.Handler.Create)
				groups.GET("/:id", r.container.StaffGroup.Handler.Get)
				groups.PUT("/:id", r.container.StaffGroup.Handler.Update)
				groups.DELETE("/:id", r.container.StaffGroup.Handler.Delete)
				groups.GET("/:id/staff", r.container.StaffGroup.Handler.ListStaff)
				groups.POST("/:id/staff/:staffId", r.container.StaffGroup.Handler.AddStaff)
				groups.DELETE("/:id/staff/:staffId", r.container.StaffGroup.Handler.RemoveStaff)
			}
		}
	}

	return router
}
