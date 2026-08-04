package service

import (
	"context"

	"ads-platform-ui/internal/core/cdn"
)

type CategoryService interface {
	List(ctx context.Context) ([]cdn.Category, error)
}
