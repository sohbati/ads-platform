package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

const OtpSubject = "notifications.otp.send"

type OtpEvent struct {
	Mobile string `json:"mobile"`
	Otp    string `json:"otp"`
}

type OtpEventPublisher interface {
	PublishOtpEvent(ctx context.Context, mobile string, otp string) error
}

type otpEventPublisher struct {
	conn    *nats.Conn
	subject string
}

func NewOtpEventPublisher(natsURL string, subject string) (OtpEventPublisher, error) {
	if natsURL == "" {
		return &noopOtpEventPublisher{}, nil
	}

	if subject == "" {
		subject = OtpSubject
	}

	conn, err := nats.Connect(natsURL,
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("connect to nats: %w", err)
	}

	log.Printf("OTP event publisher connected to NATS at %s", conn.ConnectedUrl())

	return &otpEventPublisher{
		conn:    conn,
		subject: subject,
	}, nil
}

func (p *otpEventPublisher) PublishOtpEvent(ctx context.Context, mobile string, otp string) error {
	event := OtpEvent{Mobile: mobile, Otp: otp}
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal otp event: %w", err)
	}

	if err := p.conn.Publish(p.subject, payload); err != nil {
		return fmt.Errorf("publish otp event: %w", err)
	}

	log.Printf("OTP event published subject=%s mobile=%s", p.subject, mobile)
	return nil
}

type noopOtpEventPublisher struct{}

func (p *noopOtpEventPublisher) PublishOtpEvent(_ context.Context, _ string, _ string) error {
	return nil
}
