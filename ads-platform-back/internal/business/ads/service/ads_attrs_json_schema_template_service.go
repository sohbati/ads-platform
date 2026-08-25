package service

import (
	"context"

	"ads-platform/internal/business/ads/model"
)

// AdsAttrsJSONSchemaTemplateService looks up JSON Schema templates that drive
// the category-specific attrs form. HTTP APIs will sit on top of this later.
type AdsAttrsJSONSchemaTemplateService interface {
	GetByName(ctx context.Context, name string) (*model.AdsAttrsJSONSchemaTemplate, error)
	GetByID(ctx context.Context, id int64) (*model.AdsAttrsJSONSchemaTemplate, error)
	Create(ctx context.Context, template *model.AdsAttrsJSONSchemaTemplate) error
	Update(ctx context.Context, template *model.AdsAttrsJSONSchemaTemplate) error
}
