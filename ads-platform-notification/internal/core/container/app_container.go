package container

import (
	otpContainer "ads-platform-notification/internal/business/otp/container"
	"ads-platform-notification/internal/core/broker"
	"ads-platform-notification/internal/core/config"
)

type AppContainer struct {
	Broker *broker.Connection
	Otp    *otpContainer.OtpContainer
}

func NewAppContainer(cfg *config.Config, brokerConn *broker.Connection) (*AppContainer, error) {
	otp, err := otpContainer.NewOtpContainer(brokerConn, cfg.OtpSubject)
	if err != nil {
		return nil, err
	}

	return &AppContainer{
		Broker: brokerConn,
		Otp:    otp,
	}, nil
}
