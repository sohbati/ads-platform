package container

import (
	"ads-platform-ui/internal/business/category/handler"
	serviceimpl "ads-platform-ui/internal/business/category/service/impl"
	"ads-platform-ui/internal/core/cdn"
)

type CategoryContainer struct {
	APIHandler *handler.APIHandler
}

func NewCategoryContainer(cdnClient *cdn.Client) *CategoryContainer {
	svc := serviceimpl.NewCategoryService(cdnClient)
	return &CategoryContainer{
		APIHandler: handler.NewAPIHandler(svc),
	}
}
