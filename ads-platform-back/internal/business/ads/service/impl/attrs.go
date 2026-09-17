package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"ads-platform/internal/business/ads/errorcode"
	searchclient "ads-platform/internal/business/search/client"
	"ads-platform/internal/core/exception"
)

type attrSchemaDoc struct {
	Properties map[string]attrProp `json:"properties"`
	Required   []string            `json:"required"`
}

type attrProp struct {
	Type    string    `json:"type"`
	Enum    []any     `json:"enum"`
	Minimum *float64  `json:"minimum"`
	Maximum *float64  `json:"maximum"`
}

func (s *adService) prepareAttrs(ctx context.Context, categoryID int, raw json.RawMessage) (json.RawMessage, error) {
	cleaned := normalizeJSONObject(raw)
	if s.catalog == nil {
		return cleaned, nil
	}
	cats, err := s.catalog.CategoriesByIDs(ctx, []int{categoryID}, false)
	if err != nil || len(cats) == 0 {
		return cleaned, nil
	}
	name := ""
	if cats[0].AdsAttrsJSONSchemaTemplateName != nil {
		name = strings.TrimSpace(*cats[0].AdsAttrsJSONSchemaTemplateName)
	}
	if name == "" {
		return cleaned, nil
	}
	schemas, err := s.catalog.AttrSchemasByNames(ctx, []string{name})
	if err != nil || len(schemas) == 0 {
		return cleaned, nil
	}
	return filterAttrsBySchema(cleaned, schemas[0])
}

func filterAttrsBySchema(raw json.RawMessage, schema searchclient.AttrSchema) (json.RawMessage, error) {
	var doc attrSchemaDoc
	if err := json.Unmarshal(schema.JSONSchema, &doc); err != nil || len(doc.Properties) == 0 {
		return raw, nil
	}
	var incoming map[string]any
	if err := json.Unmarshal(raw, &incoming); err != nil {
		return nil, exception.NewAppError(errorcode.ErrAdInvalidAttrs.Code, errorcode.ErrAdInvalidAttrs.HttpStatus)
	}
	if incoming == nil {
		incoming = map[string]any{}
	}

	out := make(map[string]any, len(doc.Properties))
	for name, prop := range doc.Properties {
		val, ok := incoming[name]
		if !ok || val == nil {
			continue
		}
		coerced, err := coerceAttr(prop, val)
		if err != nil {
			return nil, exception.NewAppError(errorcode.ErrAdInvalidAttrs.Code, errorcode.ErrAdInvalidAttrs.HttpStatus)
		}
		if coerced != nil {
			out[name] = coerced
		}
	}
	for _, name := range doc.Required {
		if _, ok := out[name]; !ok {
			return nil, exception.NewAppError(errorcode.ErrAdInvalidAttrs.Code, errorcode.ErrAdInvalidAttrs.HttpStatus)
		}
	}
	body, err := json.Marshal(out)
	if err != nil {
		return json.RawMessage("{}"), nil
	}
	return body, nil
}

func coerceAttr(prop attrProp, val any) (any, error) {
	switch prop.Type {
	case "boolean":
		b, ok := val.(bool)
		if !ok {
			return nil, fmt.Errorf("not boolean")
		}
		return b, nil
	case "integer":
		n, ok := asFloat(val)
		if !ok || math.Trunc(n) != n {
			return nil, fmt.Errorf("not integer")
		}
		if err := checkBounds(prop, n); err != nil {
			return nil, err
		}
		if len(prop.Enum) > 0 && !enumContains(prop.Enum, fmt.Sprint(int64(n))) {
			return nil, fmt.Errorf("invalid enum")
		}
		return int64(n), nil
	case "number":
		n, ok := asFloat(val)
		if !ok {
			return nil, fmt.Errorf("not number")
		}
		if err := checkBounds(prop, n); err != nil {
			return nil, err
		}
		if len(prop.Enum) > 0 && !enumContains(prop.Enum, strings.TrimSpace(fmt.Sprint(n))) {
			return nil, fmt.Errorf("invalid enum")
		}
		return n, nil
	default:
		s, ok := val.(string)
		if !ok {
			return nil, fmt.Errorf("not string")
		}
		s = strings.TrimSpace(s)
		if s == "" {
			return nil, nil
		}
		if len(prop.Enum) > 0 && !enumContains(prop.Enum, s) {
			return nil, fmt.Errorf("invalid enum")
		}
		return s, nil
	}
}

func asFloat(val any) (float64, bool) {
	switch n := val.(type) {
	case float64:
		return n, true
	case json.Number:
		f, err := n.Float64()
		return f, err == nil
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	default:
		return 0, false
	}
}

func checkBounds(prop attrProp, n float64) error {
	if prop.Minimum != nil && n < *prop.Minimum {
		return fmt.Errorf("below min")
	}
	if prop.Maximum != nil && n > *prop.Maximum {
		return fmt.Errorf("above max")
	}
	return nil
}

func enumContains(enum []any, s string) bool {
	for _, item := range enum {
		if strings.TrimSpace(fmt.Sprint(item)) == s {
			return true
		}
	}
	return false
}
