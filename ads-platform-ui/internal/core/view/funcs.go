package view

import (
	"html/template"
	"strings"

	"ads-platform-ui/internal/core/i18n"
)

// FuncMap returns custom template functions shared across all pages.
func FuncMap() template.FuncMap {
	return template.FuncMap{
		"format":    i18n.Formatf,
		"hasPrefix": strings.HasPrefix,
	}
}
