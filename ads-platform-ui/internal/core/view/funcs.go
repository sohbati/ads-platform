package view

import (
	"encoding/json"
	"html/template"
	"strings"

	"ads-platform-ui/internal/core/i18n"
)

// FuncMap returns custom template functions shared across all pages.
// asset resolves a static file path (relative to the static root) to its
// versioned public URL; pass nil to fall back to unversioned /static URLs.
func FuncMap(asset func(string) string) template.FuncMap {
	if asset == nil {
		asset = func(rel string) string { return "/static/" + rel }
	}
	return template.FuncMap{
		"asset":     asset,
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
