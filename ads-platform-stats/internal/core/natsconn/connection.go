package natsconn

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

type Connection struct {
	conn *nats.Conn
}

func Connect(url string) (*Connection, error) {
	conn, err := nats.Connect(url,
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2*time.Second),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			if err != nil {
				log.Printf("NATS disconnected: %v", err)
			}
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("NATS reconnected to %s", nc.ConnectedUrl())
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

func (c *Connection) QueueSubscribe(subject, queue string, handler nats.MsgHandler) (*nats.Subscription, error) {
	return c.conn.QueueSubscribe(subject, queue, handler)
}
