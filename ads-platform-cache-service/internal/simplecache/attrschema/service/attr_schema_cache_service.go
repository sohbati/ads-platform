package service

import (
	"context"

	"cache-service/internal/simplecache/attrschema/model"
)

type AttrSchemaCacheService interface {
	GetByNames(ctx context.Context, names []string) ([]model.AttrSchema, error)
}
