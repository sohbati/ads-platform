package impl

import (
	"bytes"
	"encoding/json"
	"math"
	"strconv"
	"strings"

	"ads-platform-ui/internal/business/queryads/viewmodel"
	"ads-platform-ui/internal/core/cdn"
	"ads-platform-ui/internal/core/i18n"
)

type attrSchemaDoc struct {
	Properties json.RawMessage `json:"properties"`
}

type attrProp struct {
	Type       string `json:"type"`
	Title      string `json:"title"`
	EnumVocab  string `json:"x-enumVocab"`
}

func labeledAttrs(loc i18n.Locale, t i18n.Messages, categoryID int, attrs json.RawMessage, categories []cdn.Category, schemas []cdn.AttrSchema, enums json.RawMessage) []viewmodel.AdAttr {
	if categoryID <= 0 || len(attrs) == 0 {
		return nil
	}
	var values map[string]any
	if err := json.Unmarshal(attrs, &values); err != nil || len(values) == 0 {
		return nil
	}

	templateName := ""
	for _, cat := range categories {
		if cat.ID != categoryID || cat.AdsAttrsJSONSchemaTemplateName == nil {
			continue
		}
		templateName = strings.TrimSpace(*cat.AdsAttrsJSONSchemaTemplateName)
		break
	}
	if templateName == "" {
		return nil
	}

	var schema cdn.AttrSchema
	found := false
	for _, item := range schemas {
		if item.Name == templateName {
			schema = item
			found = true
			break
		}
	}
	if !found {
		return nil
	}

	var doc attrSchemaDoc
	if err := json.Unmarshal(schema.JSONSchema, &doc); err != nil || len(doc.Properties) == 0 {
		return nil
	}
	var props map[string]attrProp
	if err := json.Unmarshal(doc.Properties, &props); err != nil {
		return nil
	}

	var enumMap map[string]map[string]map[string]string
	_ = json.Unmarshal(enums, &enumMap)

	locale := string(loc)
	out := make([]viewmodel.AdAttr, 0, len(props))
	for _, name := range jsonObjectKeys(doc.Properties) {
		prop, ok := props[name]
		if !ok {
			continue
		}
		val, ok := values[name]
		if !ok || val == nil {
			continue
		}
		formatted := formatAttrValue(prop, val, locale, enumMap, t)
		if formatted == "" {
			continue
		}
		label := strings.TrimSpace(prop.Title)
		if label == "" {
			label = name
		}
		out = append(out, viewmodel.AdAttr{Label: label, Value: formatted})
	}
	return out
}

func formatAttrValue(prop attrProp, val any, locale string, enums map[string]map[string]map[string]string, t i18n.Messages) string {
	switch prop.Type {
	case "boolean":
		b, ok := val.(bool)
		if !ok {
			return ""
		}
		if b {
			return t.AdDetail.Yes
		}
		return t.AdDetail.No
	case "integer", "number":
		n, ok := attrAsFloat(val)
		if !ok {
			return ""
		}
		if n == math.Trunc(n) {
			return strconv.FormatInt(int64(n), 10)
		}
		return strconv.FormatFloat(n, 'f', -1, 64)
	default:
		s, ok := val.(string)
		if !ok {
			return ""
		}
		s = strings.TrimSpace(s)
		if s == "" {
			return ""
		}
		if prop.EnumVocab != "" && enums != nil {
			if byToken, ok := enums[prop.EnumVocab]; ok {
				if byLoc, ok := byToken[s]; ok {
					if label := strings.TrimSpace(byLoc[locale]); label != "" {
						return label
					}
					if label := strings.TrimSpace(byLoc["fa"]); label != "" {
						return label
					}
					if label := strings.TrimSpace(byLoc["en"]); label != "" {
						return label
					}
				}
			}
		}
		return s
	}
}

func attrAsFloat(val any) (float64, bool) {
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

func jsonObjectKeys(raw json.RawMessage) []string {
	dec := json.NewDecoder(bytes.NewReader(raw))
	tok, err := dec.Token()
	if err != nil {
		return nil
	}
	delim, ok := tok.(json.Delim)
	if !ok || delim != '{' {
		return nil
	}
	var names []string
	for dec.More() {
		key, err := dec.Token()
		if err != nil {
			break
		}
		name, _ := key.(string)
		var skip json.RawMessage
		if err := dec.Decode(&skip); err != nil {
			break
		}
		if name != "" {
			names = append(names, name)
		}
	}
	return names
}
