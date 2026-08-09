package container

import (
	"ads-platform-ui/internal/business/otp/handler"
	"ads-platform-ui/internal/core/bff"
)

type OtpContainer struct {
	APIHandler *handler.APIHandler
}

func NewOtpContainer(client *bff.Client) *OtpContainer {
	return &OtpContainer{
		APIHandler: handler.NewAPIHandler(client),
	}
}
