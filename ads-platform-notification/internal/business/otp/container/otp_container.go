package container

import (
	"ads-platform-notification/internal/business/otp/listener"
	serviceimpl "ads-platform-notification/internal/business/otp/service/impl"
	smsimpl "ads-platform-notification/internal/business/otp/sms/impl"
	"ads-platform-notification/internal/core/broker"
)

type OtpContainer struct {
	OtpListener *listener.OtpListener
}

func NewOtpContainer(brokerConn *broker.Connection, subject string) (*OtpContainer, error) {
	smsProvider := smsimpl.NewLogProvider()
	otpService := serviceimpl.NewOtpNotificationService(smsProvider)

	otpListener, err := listener.NewOtpListener(brokerConn, subject, otpService)
	if err != nil {
		return nil, err
	}

	return &OtpContainer{
		OtpListener: otpListener,
	}, nil
}
