package impl

import (
	"slices"
	"sort"
	"time"

	"ads-platform-cdn/internal/business/category/model"
	"ads-platform-cdn/internal/core/jsonstore"
)

type CategoryServiceImpl struct {
	store *jsonstore.Store[model.Category]
}

func NewCategoryService(jsonPath string) *CategoryServiceImpl {
	return &CategoryServiceImpl{
		store: jsonstore.New[model.Category](jsonPath, 30*time.Second),
	}
}

func (s *CategoryServiceImpl) List() ([]model.Category, error) {
	items, err := s.store.Get()
	if err != nil {
		return nil, err
	}

	items = slices.Clone(items)
	sort.SliceStable(items, func(i, j int) bool {
		return items[i].Order < items[j].Order
	})
	return items, nil
}
