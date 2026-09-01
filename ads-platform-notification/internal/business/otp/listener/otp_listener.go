package listener

import (
	"context"
	"encoding/json"
	"log"

	"ads-platform-notification/internal/business/otp/model"
	"ads-platform-notification/internal/business/otp/service"
	"ads-platform-notification/internal/core/broker"
)

type OtpListener struct {
	subscription broker.Subscription
}

func NewOtpListener(brokerConn *broker.Connection, subject string, otpService service.OtpNotificationService) (*OtpListener, error) {
	sub, err := brokerConn.Subscribe(subject, func(data []byte) {
		var event model.OtpEvent
		if err := json.Unmarshal(data, &event); err != nil {
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
