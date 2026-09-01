package container

import (
	eventContainer "message-broker/internal/business/event/container"
	"message-broker/internal/core/broker"
	"message-broker/internal/core/brokerserver"
)

type AppContainer struct {
	BrokerServer *brokerserver.EmbeddedServer
	Broker       *broker.Connection
	Event        *eventContainer.EventContainer
}

func NewAppContainer(brokerServer *brokerserver.EmbeddedServer, conn *broker.Connection) *AppContainer {
	return &AppContainer{
		BrokerServer: brokerServer,
		Broker:       conn,
		Event:        eventContainer.NewEventContainer(conn),
	}
}
