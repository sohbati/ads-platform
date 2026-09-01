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

type natsPublisher struct {
	conn    *nats.Conn
	subject string
}

func New(natsURL, subject string) (Publisher, error) {
	if natsURL == "" {
		return Noop(), nil
	}
	if subject == "" {
		subject = "ads.stats.event"
	}

	conn, err := nats.Connect(natsURL,
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("connect to nats: %w", err)
	}
	log.Printf("stats publisher connected to NATS at %s", conn.ConnectedUrl())
	return &natsPublisher{conn: conn, subject: subject}, nil
}

func (p *natsPublisher) Publish(_ context.Context, payload []byte) error {
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
