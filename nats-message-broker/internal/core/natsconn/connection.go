package natsconn

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

type Connection struct {
	conn *nats.Conn
}

func Connect(url string) (*Connection, error) {
	conn, err := nats.Connect(url,
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			if err != nil {
				log.Printf("NATS disconnected: %v", err)
			}
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("NATS reconnected to %s", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(_ *nats.Conn) {
			log.Println("NATS connection closed")
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("connect to nats: %w", err)
	}

	log.Printf("Connected to NATS at %s", conn.ConnectedUrl())
	return &Connection{conn: conn}, nil
}

func (c *Connection) Conn() *nats.Conn {
	return c.conn
}

func (c *Connection) IsConnected() bool {
	return c.conn != nil && c.conn.IsConnected()
}

func (c *Connection) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Connection) Publish(subject string, data []byte) error {
	return c.conn.Publish(subject, data)
}

func (c *Connection) Subscribe(subject string, handler nats.MsgHandler) (*nats.Subscription, error) {
	return c.conn.Subscribe(subject, handler)
}
