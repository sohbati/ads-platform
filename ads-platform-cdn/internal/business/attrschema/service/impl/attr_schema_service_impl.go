package impl

import (
	"sort"
	"time"

	"ads-platform-cdn/internal/business/attrschema/model"
	"ads-platform-cdn/internal/core/jsonstore"
)

type AttrSchemaServiceImpl struct {
	store *jsonstore.MapStore[model.FileEntry]
}

func NewAttrSchemaService(jsonPath string) *AttrSchemaServiceImpl {
	return &AttrSchemaServiceImpl{
		store: jsonstore.NewMap[model.FileEntry](jsonPath, 30*time.Second),
	}
}

func (s *AttrSchemaServiceImpl) List() ([]model.AttrSchema, error) {
	items, err := s.store.Get()
	if err != nil {
		return nil, err
	}

	out := make([]model.AttrSchema, 0, len(items))
	for name, entry := range items {
		out = append(out, model.AttrSchema{
			Name:       name,
			Title:      entry.Title,
			JSONSchema: entry.JSONSchema,
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}
