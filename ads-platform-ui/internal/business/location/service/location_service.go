package service

import (
	"context"

	"ads-platform-ui/internal/core/cdn"
)

type LocationService interface {
	ListCities(ctx context.Context) ([]cdn.City, error)
}
