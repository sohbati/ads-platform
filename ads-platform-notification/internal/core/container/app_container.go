package container

import (
	"log"

	otpContainer "ads-platform-notification/internal/business/otp/container"
	"ads-platform-notification/internal/core/broker"
	"ads-platform-notification/internal/core/config"
	"ads-platform-notification/internal/core/sms"
	smsfactory "ads-platform-notification/internal/core/sms/factory"
)

type AppContainer struct {
	Broker *broker.Connection
	SMS    sms.Sender
	Otp    *otpContainer.OtpContainer
}

func NewAppContainer(cfg *config.Config, brokerConn *broker.Connection) (*AppContainer, error) {
	smsSender, err := smsfactory.New(cfg)
	if err != nil {
		return nil, err
	}
	log.Printf("env=%s sms=%s allowlist=%d numbers", cfg.AppEnv, smsSender.Name(), len(cfg.SMS.Allowlist))

	otp, err := otpContainer.NewOtpContainer(brokerConn, cfg.OtpSubject, smsSender)
	if err != nil {
		return nil, err
	}

	return &AppContainer{
		Broker: brokerConn,
		SMS:    smsSender,
		Otp:    otp,
	}, nil
}
