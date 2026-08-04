package container

import (
	"nats-message-broker/internal/business/event/handler"
	serviceimpl "nats-message-broker/internal/business/event/service/impl"
	"nats-message-broker/internal/core/natsconn"
)

type EventContainer struct {
	EventHandler *handler.EventHandler
}

func NewEventContainer(nats *natsconn.Connection) *EventContainer {
	eventService := serviceimpl.NewEventService(nats)
	eventHandler := handler.NewEventHandler(eventService)

	return &EventContainer{
		EventHandler: eventHandler,
	}
}
