package middleware

import (
	"log"
	"net/http"

	"ads-platform-stats/internal/core/exception"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error      string   `json:"error"`
	StatusCode int      `json:"statusCode"`
	Params     []string `json:"params,omitempty"`
}

func GlobalErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			statusCode, errorKey, params := http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", []string(nil)
			if appErr, ok := exception.AsAppError(err.Err); ok {
				statusCode, errorKey, params = appErr.StatusCode, appErr.ErrorCode, appErr.Params
			}
			c.JSON(statusCode, ErrorResponse{Error: errorKey, StatusCode: statusCode, Params: params})
			c.Abort()
		}
	}
}

func CustomRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Error:      "INTERNAL_SERVER_ERROR",
					StatusCode: http.StatusInternalServerError,
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
