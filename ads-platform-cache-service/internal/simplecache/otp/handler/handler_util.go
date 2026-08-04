package handler

import (
	"cache-service/internal/core/exception"
	"cache-service/internal/simplecache/otp/errorcode"

	"github.com/gin-gonic/gin"
)

func GetBodyAsString(c *gin.Context) (string, error) {
	bodyBytes, err := c.GetRawData()
	key := c.Param("key")
	if err != nil {
		return "", exception.NewAppError(
			errorcode.ErrFailToReadBody.Code, errorcode.ErrFailToReadBody.HttpStatus, key).WithCause(err)

	}

	return string(bodyBytes), nil
}
