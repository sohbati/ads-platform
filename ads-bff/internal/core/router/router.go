package router

import (
	"context"
	"net/http"
	"strings"
	"time"

	appContainer "ads-bff/internal/core/container"
	"ads-bff/internal/core/middleware"

	"github.com/gin-gonic/gin"
)

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
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
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api/v1/") && !strings.HasPrefix(path, "/api/v1/auth") {
			r.container.Backend.ServeHTTP(c.Writer, c.Request)
			return
		}

		c.JSON(http.StatusNotFound, gin.H{
			"error":      "ROUTE_NOT_DEFINED",
			"statusCode": http.StatusNotFound,
		})
	})

	router.GET("/health", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		backendOK := r.container.Backend.BackendHealthy(ctx)
		status := "ok"
		if !backendOK {
			status = "degraded"
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  status,
			"service": "ads-bff",
			"backend": map[string]any{
				"reachable": backendOK,
				"url":       r.container.Config.BackendAPIBaseURL,
			},
		})
	})

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/otp/:mobile/send", r.container.Auth.AuthHandler.SendOTP)
			auth.POST("/otp/:mobile/verify", r.container.Auth.AuthHandler.VerifyOTP)
			auth.GET("/me", r.container.Auth.AuthHandler.Me)
			auth.POST("/logout", r.container.Auth.AuthHandler.Logout)
		}

		api.POST("/ads", r.container.Ads.AdHandler.Create)
		api.GET("/me/profile", r.container.Profile.Handler.Get)
		api.PUT("/me/profile", r.container.Profile.Handler.Put)
	}

	return router
}
