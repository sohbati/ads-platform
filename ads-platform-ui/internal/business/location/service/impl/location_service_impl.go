package impl

import (
	"context"

	"ads-platform-ui/internal/core/cdn"
)

type LocationServiceImpl struct {
	cdn *cdn.Client
}

func NewLocationService(client *cdn.Client) *LocationServiceImpl {
	return &LocationServiceImpl{cdn: client}
}

func (s *LocationServiceImpl) ListCities(ctx context.Context) ([]cdn.City, error) {
	return s.cdn.GetCities(ctx)
}
