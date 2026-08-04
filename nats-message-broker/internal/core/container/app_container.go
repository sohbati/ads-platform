package container

import (
	eventContainer "nats-message-broker/internal/business/event/container"
	"nats-message-broker/internal/core/natsconn"
	"nats-message-broker/internal/core/natsserver"
)

type AppContainer struct {
	NatsServer *natsserver.EmbeddedServer
	Nats       *natsconn.Connection
	Event      *eventContainer.EventContainer
}

func NewAppContainer(natsServer *natsserver.EmbeddedServer, nats *natsconn.Connection) *AppContainer {
	return &AppContainer{
		NatsServer: natsServer,
		Nats:       nats,
		Event:      eventContainer.NewEventContainer(nats),
	}
}
