package impl

import (
	"context"

	"ads-platform-ui/internal/core/cdn"
)

type CategoryServiceImpl struct {
	cdn *cdn.Client
}

func NewCategoryService(client *cdn.Client) *CategoryServiceImpl {
	return &CategoryServiceImpl{cdn: client}
}

func (s *CategoryServiceImpl) List(ctx context.Context) ([]cdn.Category, error) {
	return s.cdn.GetCategories(ctx)
}
