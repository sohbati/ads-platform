package listener

import (
	"context"
	"encoding/json"
	"log"

	"ads-platform-notification/internal/business/otp/model"
	"ads-platform-notification/internal/business/otp/service"
	"ads-platform-notification/internal/core/natsconn"

	"github.com/nats-io/nats.go"
)

type OtpListener struct {
	subscription *nats.Subscription
}

func NewOtpListener(natsConn *natsconn.Connection, subject string, otpService service.OtpNotificationService) (*OtpListener, error) {
	sub, err := natsConn.Subscribe(subject, func(msg *nats.Msg) {
		var event model.OtpEvent
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Printf("[OTP] failed to parse message on %s: %v", subject, err)
			return
		}

		if err := otpService.HandleOtpEvent(context.Background(), &event); err != nil {
			log.Printf("[OTP] failed to handle event mobile=%s: %v", event.Mobile, err)
		}
	})
	if err != nil {
		return nil, err
	}

	log.Printf("[OTP] listening on subject %s", subject)
	return &OtpListener{subscription: sub}, nil
}

func (l *OtpListener) Stop() {
	if l.subscription != nil {
		if err := l.subscription.Unsubscribe(); err != nil {
			log.Printf("[OTP] unsubscribe error: %v", err)
		}
	}
}
