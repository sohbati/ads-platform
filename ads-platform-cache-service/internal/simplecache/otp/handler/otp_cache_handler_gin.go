package handler

import (
	"cache-service/internal/core/exception"
	"cache-service/internal/simplecache/otp/errorcode"
	"cache-service/internal/simplecache/otp/model"
	"cache-service/internal/simplecache/otp/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CacheHandler handles HTTP requests for cache-related operations
type OtpCacheHandler struct {
	cacheService service.OtpCacheService
}

// NewCacheHandler creates a new cache handler
func NewOtpCacheHandler(cacheService service.OtpCacheService) *OtpCacheHandler {
	return &OtpCacheHandler{
		cacheService: cacheService,
	}
}

// CheckCacheExists handles GET /caches/check/:cache_id
// Checks if a cache exists for the given cache ID
func (h *OtpCacheHandler) GetCacheByKey(c *gin.Context) {
	key := c.Param("key")

	cache, err := h.cacheService.GetCacheByKey(c.Request.Context(), key)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, cache)
}

// CheckCacheExists handles GET /caches/check/:cache_id
// Checks if a cache exists for the given cache ID
func (h *OtpCacheHandler) AddItem(c *gin.Context) {
	key := c.Param("key")

	var req model.OtpHttpRequestModel
	if err := c.ShouldBindJSON(&req); err != nil {
		exc := exception.NewAppError(
			errorcode.ErrFailToReadBody.Code, errorcode.ErrFailToReadBody.HttpStatus, key, ).WithCause(err)
		c.Error(exc)
		c.Abort()
		return
	}

	if key == "" {
		exc := exception.NewAppError(
			errorcode.ErrKeyEmpty.Code, errorcode.ErrKeyEmpty.HttpStatus, key)
		c.Error(exc)
		c.Abort()
		return
	}
	err := h.cacheService.AddItem(c.Request.Context(), key, req.Otp)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, "")
}
