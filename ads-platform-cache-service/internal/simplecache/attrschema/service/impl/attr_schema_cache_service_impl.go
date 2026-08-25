package impl

import (
	"context"
	"encoding/json"
	"strings"
	"sync"

	"cache-service/internal/cachestore"
	"cache-service/internal/core/exception"
	"cache-service/internal/simplecache/attrschema/client"
	"cache-service/internal/simplecache/attrschema/errorcode"
	"cache-service/internal/simplecache/attrschema/model"
	"cache-service/internal/simplecache/attrschema/service"
)

type attrSchemaCacheService struct {
	nameToSchema *cachestore.Cache[string, string] // name → schema JSON
	cdn          client.CDNAttrSchemaClient
	loadMu       sync.Mutex
}

func NewAttrSchemaCacheService(
	nameToSchema *cachestore.Cache[string, string],
	cdn client.CDNAttrSchemaClient,
) service.AttrSchemaCacheService {
	return &attrSchemaCacheService{
		nameToSchema: nameToSchema,
		cdn:          cdn,
	}
}

func (s *attrSchemaCacheService) GetByNames(ctx context.Context, names []string) ([]model.AttrSchema, error) {
	clean := make([]string, 0, len(names))
	for _, name := range names {
		name = strings.TrimSpace(name)
		if name != "" {
			clean = append(clean, name)
		}
	}
	if len(clean) == 0 {
		return nil, exception.NewAppError(errorcode.ErrNamesEmpty.Code, errorcode.ErrNamesEmpty.HttpStatus)
	}
	if err := s.ensureLoaded(ctx); err != nil {
		return nil, err
	}

	out := make([]model.AttrSchema, 0, len(clean))
	for _, name := range clean {
		raw, err := s.nameToSchema.Get(ctx, name)
		if err != nil {
			continue
		}
		var schema model.AttrSchema
		if err := json.Unmarshal([]byte(raw), &schema); err != nil {
			continue
		}
		out = append(out, schema)
	}
	return out, nil
}

func (s *attrSchemaCacheService) ensureLoaded(ctx context.Context) error {
	if s.nameToSchema.Count() > 0 {
		return nil
	}

	s.loadMu.Lock()
	defer s.loadMu.Unlock()

	if s.nameToSchema.Count() > 0 {
		return nil
	}

	items, err := s.cdn.ListAttrSchemas(ctx)
	if err != nil {
		return exception.NewAppError(
			errorcode.ErrCDNUnavailable.Code,
			errorcode.ErrCDNUnavailable.HttpStatus,
		).WithCause(err)
	}
	if len(items) == 0 {
		return exception.NewAppError(errorcode.ErrCacheEmpty.Code, errorcode.ErrCacheEmpty.HttpStatus)
	}

	for _, item := range items {
		if item.Name == "" {
			continue
		}
		raw, err := json.Marshal(item)
		if err != nil {
			continue
		}
		_ = s.nameToSchema.Set(ctx, item.Name, string(raw))
	}
	return nil
}
