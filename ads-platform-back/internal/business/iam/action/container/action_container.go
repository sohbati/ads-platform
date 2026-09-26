package container

import (
	"ads-platform/internal/business/iam/action/handler"
	repoimpl "ads-platform/internal/business/iam/action/repository/impl"
	serviceimpl "ads-platform/internal/business/iam/action/service/impl"

	"gorm.io/gorm"
)

type ActionContainer struct {
	Handler *handler.ActionHandler
}

func NewActionContainer(db *gorm.DB) *ActionContainer {
	repo := repoimpl.NewActionRepository(db)
	svc := serviceimpl.NewActionService(repo)
	return &ActionContainer{Handler: handler.NewActionHandler(svc)}
}
