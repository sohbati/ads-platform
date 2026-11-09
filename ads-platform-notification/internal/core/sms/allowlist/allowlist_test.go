package allowlist

import (
	"context"
	"errors"
	"testing"

	"ads-platform-notification/internal/core/sms"
)

type recorder struct{ sent []string }

func (r *recorder) Name() string { return "rec" }
func (r *recorder) SendTemplate(_ context.Context, m sms.TemplateMessage) (*sms.Result, error) {
	r.sent = append(r.sent, m.Mobile)
	return &sms.Result{}, nil
}
func (r *recorder) SendText(_ context.Context, m sms.TextMessage) (*sms.Result, error) {
	r.sent = append(r.sent, m.Mobile)
	return &sms.Result{}, nil
}

func TestAllowlist(t *testing.T) {
	rec := &recorder{}
	s, err := Wrap(rec, []string{"+98 912 111 2233"})
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	for _, m := range []string{"09121112233", "9121112233", "00989121112233", "۰۹۱۲۱۱۱۲۲۳۳"} {
		if _, err := s.SendTemplate(ctx, sms.TemplateMessage{Mobile: m}); err != nil {
			t.Errorf("%s should be allowed: %v", m, err)
		}
	}
	for _, m := range []string{"09129999999", "", "abc"} {
		if _, err := s.SendText(ctx, sms.TextMessage{Mobile: m}); !errors.Is(err, sms.ErrRecipientNotAllowed) {
			t.Errorf("%q should be blocked, got %v", m, err)
		}
	}
	if len(rec.sent) != 4 {
		t.Errorf("provider got %d messages, want 4", len(rec.sent))
	}
}

func TestEmptyListBlocksEveryone(t *testing.T) {
	rec := &recorder{}
	s, _ := Wrap(rec, nil)
	if _, err := s.SendTemplate(context.Background(), sms.TemplateMessage{Mobile: "09121112233"}); !errors.Is(err, sms.ErrRecipientNotAllowed) {
		t.Fatalf("got %v", err)
	}
	if len(rec.sent) != 0 {
		t.Fatal("nothing should reach the provider")
	}
}

func TestInvalidEntryFailsAtStartup(t *testing.T) {
	if _, err := Wrap(&recorder{}, []string{"0912"}); err == nil {
		t.Fatal("expected error for invalid entry")
	}
}
