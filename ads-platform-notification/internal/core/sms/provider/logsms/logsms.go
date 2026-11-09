// Package logsms is a development provider that only logs messages.
package logsms

import (
	"context"
	"log"

	"ads-platform-notification/internal/core/sms"
)

type sender struct{}

func New() sms.Sender { return &sender{} }

func (s *sender) Name() string { return "log" }

func (s *sender) SendTemplate(_ context.Context, msg sms.TemplateMessage) (*sms.Result, error) {
	log.Printf("[SMS:log] to=%s template=%s params=%v", msg.Mobile, msg.Template, msg.Params)
	return &sms.Result{Provider: s.Name()}, nil
}

func (s *sender) SendText(_ context.Context, msg sms.TextMessage) (*sms.Result, error) {
	log.Printf("[SMS:log] to=%s text=%q", msg.Mobile, msg.Text)
	return &sms.Result{Provider: s.Name()}, nil
}
