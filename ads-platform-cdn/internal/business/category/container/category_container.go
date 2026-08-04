package container

import (
	"ads-platform-cdn/internal/business/category/handler"
	serviceimpl "ads-platform-cdn/internal/business/category/service/impl"
)

type CategoryContainer struct {
	APIHandler *handler.APIHandler
}

func NewCategoryContainer(categoryJSONPath string) *CategoryContainer {
	svc := serviceimpl.NewCategoryService(categoryJSONPath)
	return &CategoryContainer{
		APIHandler: handler.NewAPIHandler(svc),
	}
}
