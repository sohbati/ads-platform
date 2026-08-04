package service

import "ads-platform-cdn/internal/business/city/model"

type CityService interface {
	List() ([]model.City, error)
}
