package service

import (
	"context"

	"ads-platform-ui/internal/core/cdn"
)

type LocationService interface {
	ListCities(ctx context.Context) ([]cdn.City, error)
}

type GeoPoint struct {
	Lat    float64 `json:"lat"`
	Lng    float64 `json:"lng"`
	Source string  `json:"source,omitempty"`
}

type ReverseResult struct {
	Lat          float64 `json:"lat"`
	Lng          float64 `json:"lng"`
	Neighborhood string  `json:"neighborhood,omitempty"`
	DisplayName  string  `json:"display_name,omitempty"`
}

type GeoService interface {
	CityCenter(ctx context.Context, slug, name, lang string) (GeoPoint, error)
	Reverse(ctx context.Context, lat, lng float64, lang string) (*ReverseResult, error)
}
