package broker

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

type Handler func(data []byte)

type Subscription interface {
	Unsubscribe() error
}

type Connection struct {
	conn *nats.Conn
}

func Connect(url string) (*Connection, error) {
	conn, err := nats.Connect(url,
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2*time.Second),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			if err != nil {
				log.Printf("message broker disconnected: %v", err)
			}
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("message broker reconnected to %s", nc.ConnectedUrl())
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("connect to message broker: %w", err)
	}

	log.Printf("Connected to message broker at %s", conn.ConnectedUrl())
	return &Connection{conn: conn}, nil
}

func (c *Connection) IsConnected() bool {
	return c.conn != nil && c.conn.IsConnected()
}

func (c *Connection) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Connection) Subscribe(subject string, handler Handler) (Subscription, error) {
	return c.conn.Subscribe(subject, func(msg *nats.Msg) {
		handler(msg.Data)
	})
}
