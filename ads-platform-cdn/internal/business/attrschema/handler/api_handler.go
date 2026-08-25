package handler

import (
	"net/http"

	"ads-platform-cdn/internal/business/attrschema/service"

	"github.com/gin-gonic/gin"
)

type APIHandler struct {
	service service.AttrSchemaService
}

func NewAPIHandler(svc service.AttrSchemaService) *APIHandler {
	return &APIHandler{service: svc}
}

func (h *APIHandler) List(c *gin.Context) {
	items, err := h.service.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load attr schemas"})
		return
	}
	c.Header("Cache-Control", "public, max-age=60")
	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, items)
}
