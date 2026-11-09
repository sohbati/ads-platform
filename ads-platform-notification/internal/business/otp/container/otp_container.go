package container

import (
	"ads-platform-notification/internal/business/otp/listener"
	serviceimpl "ads-platform-notification/internal/business/otp/service/impl"
	"ads-platform-notification/internal/core/broker"
	"ads-platform-notification/internal/core/sms"
)

type OtpContainer struct {
	OtpListener *listener.OtpListener
}

func NewOtpContainer(brokerConn *broker.Connection, subject string, smsSender sms.Sender) (*OtpContainer, error) {
	otpService := serviceimpl.NewOtpNotificationService(smsSender)

	otpListener, err := listener.NewOtpListener(brokerConn, subject, otpService)
	if err != nil {
		return nil, err
	}

	return &OtpContainer{
		OtpListener: otpListener,
	}, nil
}
