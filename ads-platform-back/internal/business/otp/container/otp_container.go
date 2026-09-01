package container

import (
	"log"
	"time"

	"ads-platform/internal/business/otp/client"
	"ads-platform/internal/business/otp/handler"
	serviceimpl "ads-platform/internal/business/otp/service/impl"
	"ads-platform/internal/core/mobile"
)

type OtpContainer struct {
	OtpHandler *handler.OtpHandler
}

func NewOtpContainer(cacheServiceURL, brokerURL, otpSubject, defaultCountryCode string, resendAfter time.Duration) *OtpContainer {
	cacheClient := client.NewOtpCacheClient(cacheServiceURL, nil)

	eventPublisher, err := client.NewOtpEventPublisher(brokerURL, otpSubject)
	if err != nil {
		log.Printf("OTP message broker publisher disabled: %v", err)
		eventPublisher, _ = client.NewOtpEventPublisher("", otpSubject)
	}

	otpService := serviceimpl.NewOtpService(cacheClient, eventPublisher, resendAfter)
	mobileNorm := mobile.NewNormalizer(defaultCountryCode)
	otpHandler := handler.NewOtpHandler(otpService, mobileNorm)

	return &OtpContainer{
		OtpHandler: otpHandler,
	}
}
