package publisher

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

type Publisher interface {
	Publish(ctx context.Context, payload []byte) error
}

type brokerPublisher struct {
	conn    *nats.Conn
	subject string
}

func New(brokerURL, subject string) (Publisher, error) {
	if brokerURL == "" {
		return Noop(), nil
	}
	if subject == "" {
		subject = "ads.stats.event"
	}

	conn, err := nats.Connect(brokerURL,
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("connect to message broker: %w", err)
	}
	log.Printf("stats publisher connected to message broker at %s", conn.ConnectedUrl())
	return &brokerPublisher{conn: conn, subject: subject}, nil
}

func (p *brokerPublisher) Publish(_ context.Context, payload []byte) error {
	if err := p.conn.Publish(p.subject, payload); err != nil {
		return fmt.Errorf("publish stats event: %w", err)
	}
	return nil
}

type noopPublisher struct{}

func Noop() Publisher {
	return noopPublisher{}
}

func (noopPublisher) Publish(context.Context, []byte) error {
	return nil
}
