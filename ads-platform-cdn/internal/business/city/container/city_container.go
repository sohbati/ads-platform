package container

import (
	"ads-platform-cdn/internal/business/city/handler"
	serviceimpl "ads-platform-cdn/internal/business/city/service/impl"
)

type CityContainer struct {
	APIHandler *handler.APIHandler
}

func NewCityContainer(citiesJSONPath string) *CityContainer {
	svc := serviceimpl.NewCityService(citiesJSONPath)
	return &CityContainer{
		APIHandler: handler.NewAPIHandler(svc),
	}
}
