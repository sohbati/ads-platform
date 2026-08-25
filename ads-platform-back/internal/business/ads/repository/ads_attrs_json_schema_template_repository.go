package repository

import (
	"context"

	"ads-platform/internal/business/ads/model"
)

type AdsAttrsJSONSchemaTemplateRepository interface {
	Create(ctx context.Context, template *model.AdsAttrsJSONSchemaTemplate) error
	GetByID(ctx context.Context, id int64) (*model.AdsAttrsJSONSchemaTemplate, error)
	GetByName(ctx context.Context, name string) (*model.AdsAttrsJSONSchemaTemplate, error)
	Update(ctx context.Context, template *model.AdsAttrsJSONSchemaTemplate) error
}
