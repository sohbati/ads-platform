package view

import (
	"encoding/json"
	"html/template"
	"strings"

	"ads-platform-ui/internal/core/i18n"
)

// FuncMap returns custom template functions shared across all pages.
func FuncMap() template.FuncMap {
	return template.FuncMap{
		"format":    i18n.Formatf,
		"hasPrefix": strings.HasPrefix,
		"json": func(v any) template.JS {
			b, err := json.Marshal(v)
			if err != nil {
				return template.JS("{}")
			}
			return template.JS(b)
		},
	}
}
