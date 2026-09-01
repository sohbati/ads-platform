package container

import (
	"message-broker/internal/business/event/handler"
	serviceimpl "message-broker/internal/business/event/service/impl"
	"message-broker/internal/core/broker"
)

type EventContainer struct {
	EventHandler *handler.EventHandler
}

func NewEventContainer(conn *broker.Connection) *EventContainer {
	eventService := serviceimpl.NewEventService(conn)
	eventHandler := handler.NewEventHandler(eventService)

	return &EventContainer{
		EventHandler: eventHandler,
	}
}
