package container

import (
	"ads-platform/internal/business/setting/handler"
	repoimpl "ads-platform/internal/business/setting/repository/impl"
	"ads-platform/internal/business/setting/service"
	serviceimpl "ads-platform/internal/business/setting/service/impl"

	"gorm.io/gorm"
)

type SettingContainer struct {
	Handler *handler.SettingHandler
	// Service is exposed for in-process readers such as the pending-ads job.
	Service service.SettingService
}

func NewSettingContainer(db *gorm.DB) *SettingContainer {
	svc := serviceimpl.NewSettingService(repoimpl.NewSettingRepository(db))
	return &SettingContainer{Handler: handler.NewSettingHandler(svc), Service: svc}
}
