// Package allowlist wraps an sms.Sender so only listed numbers receive messages.
// Outside production it keeps test traffic away from real customers, even when
// the wrapped provider is live.
package allowlist

import (
	"context"
	"fmt"
	"log"

	"ads-platform-notification/internal/core/sms"
)

type sender struct {
	next    sms.Sender
	allowed map[string]struct{}
}

// Wrap returns a Sender that forwards only to the given numbers. An empty list
// blocks every recipient (fail closed). Entries that are not valid mobile numbers
// return an error, so typos are caught at startup.
func Wrap(next sms.Sender, numbers []string) (sms.Sender, error) {
	allowed := make(map[string]struct{}, len(numbers))
	for _, n := range numbers {
		norm, ok := sms.NormalizeMobile(n)
		if !ok {
			return nil, fmt.Errorf("SMS_ALLOWLIST: %q is not a valid mobile number", n)
		}
		allowed[norm] = struct{}{}
	}
	return &sender{next: next, allowed: allowed}, nil
}

func (s *sender) Name() string { return s.next.Name() + "+allowlist" }

func (s *sender) check(mobile string) error {
	if norm, ok := sms.NormalizeMobile(mobile); ok {
		if _, listed := s.allowed[norm]; listed {
			return nil
		}
	}
	log.Printf("[SMS:allowlist] BLOCKED to=%s (not in SMS_ALLOWLIST)", mobile)
	return fmt.Errorf("%w: %s", sms.ErrRecipientNotAllowed, mobile)
}

func (s *sender) SendTemplate(ctx context.Context, msg sms.TemplateMessage) (*sms.Result, error) {
	if err := s.check(msg.Mobile); err != nil {
		return nil, err
	}
	return s.next.SendTemplate(ctx, msg)
}

func (s *sender) SendText(ctx context.Context, msg sms.TextMessage) (*sms.Result, error) {
	if err := s.check(msg.Mobile); err != nil {
		return nil, err
	}
	return s.next.SendText(ctx, msg)
}
