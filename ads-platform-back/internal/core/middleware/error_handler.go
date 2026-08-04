package middleware

import (
	"ads-platform/internal/core/exception"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Error      string   `json:"error"`
	Message    string   `json:"message,omitempty"`
	StatusCode int      `json:"statusCode"`
	Params     []string `json:"params,omitempty"`
}

// GlobalErrorHandler is similar to Spring Boot's @ControllerAdvice
// It catches all errors from handlers and formats them consistently
func GlobalErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process the request
		c.Next()
		log.Printf("********")
		// Check for errors set by handlers
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			// Determine HTTP status code and format error response
			statusCode, errorKey, params := handleError(err)

			// Return standardized error response
			c.JSON(statusCode, ErrorResponse{
				Error:      errorKey,
				StatusCode: statusCode,
				Params:     params,
			})
			c.Abort()
			return
		}
	}
}

// CustomRecovery handles panics (similar to Spring Boot's exception handling)
func CustomRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)

				errorMsg := "INTERNAL_SERVER_ERRROR"

				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Error:      errorMsg,
					StatusCode: http.StatusInternalServerError,
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

// handleError processes errors and determines status code and error key
// Similar to Spring Boot's exception handler methods
func handleError(ginErr *gin.Error) (int, string, []string) {
	err := ginErr.Err

	// Handle application errors
	if appErr, ok := exception.AsAppError(err); ok {
		return appErr.StatusCode, appErr.ErrorCode, appErr.Params
	}
	// Default: internal server error
	return http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", nil
}

// HandleError is a helper for handlers to add errors to context
// Handlers should use this instead of directly calling c.JSON for errors
func HandleError(c *gin.Context, err error, statusCode int) {
	if statusCode == 0 {
		statusCode = http.StatusInternalServerError
	}
	c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(statusCode)
	c.Abort()
}
