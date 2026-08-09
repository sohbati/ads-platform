package handler

import (
	"net/http"

	"cache-service/internal/core/exception"
	"cache-service/internal/simplecache/session/errorcode"
	"cache-service/internal/simplecache/session/model"
	"cache-service/internal/simplecache/session/service"

	"github.com/gin-gonic/gin"
)

type SessionCacheHandler struct {
	cacheService service.SessionCacheService
}

func NewSessionCacheHandler(cacheService service.SessionCacheService) *SessionCacheHandler {
	return &SessionCacheHandler{cacheService: cacheService}
}

func (h *SessionCacheHandler) GetCacheByKey(c *gin.Context) {
	key := c.Param("key")

	data, err := h.cacheService.GetCacheByKey(c.Request.Context(), key)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *SessionCacheHandler) AddItem(c *gin.Context) {
	key := c.Param("key")

	var req model.SessionHttpRequestModel
	if err := c.ShouldBindJSON(&req); err != nil {
		exc := exception.NewAppError(
			errorcode.ErrFailToReadBody.Code, errorcode.ErrFailToReadBody.HttpStatus, key).WithCause(err)
		c.Error(exc)
		c.Abort()
		return
	}

	if key == "" {
		exc := exception.NewAppError(errorcode.ErrKeyEmpty.Code, errorcode.ErrKeyEmpty.HttpStatus, key)
		c.Error(exc)
		c.Abort()
		return
	}

	if err := h.cacheService.AddItem(c.Request.Context(), key, req.Data); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, "")
}

func (h *SessionCacheHandler) DeleteItem(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		exc := exception.NewAppError(errorcode.ErrKeyEmpty.Code, errorcode.ErrKeyEmpty.HttpStatus, key)
		c.Error(exc)
		c.Abort()
		return
	}

	if err := h.cacheService.DeleteItem(c.Request.Context(), key); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Status(http.StatusNoContent)
}
