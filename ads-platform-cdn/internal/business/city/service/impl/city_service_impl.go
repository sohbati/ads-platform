package impl

import (
	"time"

	"ads-platform-cdn/internal/business/city/model"
	"ads-platform-cdn/internal/core/jsonstore"
)

type CityServiceImpl struct {
	store *jsonstore.Store[model.City]
}

func NewCityService(jsonPath string) *CityServiceImpl {
	return &CityServiceImpl{
		store: jsonstore.New[model.City](jsonPath, 30*time.Second),
	}
}

func (s *CityServiceImpl) List() ([]model.City, error) {
	return s.store.Get()
}
