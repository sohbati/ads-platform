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
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		ginErr := c.Errors.Last()
		statusCode, errorKey, params := handleError(ginErr)

		log.Printf("[ERROR] %s %s status=%d error=%q params=%v cause=%v",
			c.Request.Method,
			c.Request.URL.Path,
			statusCode,
			errorKey,
			params,
			ginErr.Err,
		)

		c.JSON(statusCode, ErrorResponse{
			Error:      errorKey,
			StatusCode: statusCode,
			Params:     params,
		})
		c.Abort()
	}
}

// CustomRecovery handles panics (similar to Spring Boot's exception handling)
func CustomRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[PANIC] %s %s recovered=%v", c.Request.Method, c.Request.URL.Path, err)

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

	if statusCode, ok := ginErr.Meta.(int); ok && statusCode != 0 {
		return statusCode, mapErrorKey(statusCode, err), nil
	}

	return http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", nil
}

func mapErrorKey(statusCode int, err error) string {
	switch statusCode {
	case http.StatusBadRequest:
		return "BAD_REQUEST"
	case http.StatusUnauthorized:
		return "UNAUTHORIZED"
	case http.StatusForbidden:
		return "FORBIDDEN"
	case http.StatusNotFound:
		return "NOT_FOUND"
	case http.StatusMethodNotAllowed:
		return "METHOD_NOT_ALLOWED"
	default:
		if err != nil {
			return err.Error()
		}
		return "INTERNAL_SERVER_ERROR"
	}
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
