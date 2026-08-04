package handler

import (
	"encoding/json"
	"net/http"
	"nats-message-broker/internal/business/event/errorcode"
	"nats-message-broker/internal/business/event/model"
	"nats-message-broker/internal/business/event/service"
	"nats-message-broker/internal/core/exception"
	"nats-message-broker/internal/core/middleware"

	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	eventService service.EventService
}

func NewEventHandler(eventService service.EventService) *EventHandler {
	return &EventHandler{eventService: eventService}
}

func (h *EventHandler) Publish(c *gin.Context) {
	subject := c.Param("subject")

	var req model.PublishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		exc := exception.NewAppError(
			errorcode.ErrFailToReadBody.Code, errorcode.ErrFailToReadBody.HttpStatus, subject).WithCause(err)
		middleware.HandleError(c, exc, http.StatusBadRequest)
		return
	}

	payload, err := json.Marshal(req.Data)
	if err != nil {
		middleware.HandleError(c, err, http.StatusBadRequest)
		return
	}

	resp, err := h.eventService.Publish(c.Request.Context(), subject, payload)
	if err != nil {
		middleware.HandleError(c, err, http.StatusMethodNotAllowed)
		return
	}

	c.JSON(http.StatusOK, resp)
}
