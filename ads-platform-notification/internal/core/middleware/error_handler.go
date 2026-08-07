package middleware

import (
	"log"
	"net/http"
	"ads-platform-notification/internal/core/exception"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error      string   `json:"error"`
	Message    string   `json:"message,omitempty"`
	StatusCode int      `json:"statusCode"`
	Params     []string `json:"params,omitempty"`
}

func GlobalErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			statusCode, errorKey, params := handleError(err)

			c.JSON(statusCode, ErrorResponse{
				Error:      errorKey,
				StatusCode: statusCode,
				Params:     params,
			})
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

func handleError(ginErr *gin.Error) (int, string, []string) {
	err := ginErr.Err

	if appErr, ok := exception.AsAppError(err); ok {
		return appErr.StatusCode, appErr.ErrorCode, appErr.Params
	}

	return http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", nil
}
