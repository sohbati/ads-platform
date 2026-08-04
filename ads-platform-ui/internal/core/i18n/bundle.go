package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

//go:embed bundles/*.json
var bundleFS embed.FS

// Bundle is one loaded language file (e.g. bundles/fa.json).
type Bundle struct {
	Locale   Locale   `json:"locale"`
	Dir      string   `json:"dir"`
	Label    string   `json:"label"`
	Messages Messages `json:"messages"`
}

type bundleFile struct {
	Locale   string   `json:"locale"`
	Dir      string   `json:"dir"`
	Label    string   `json:"label"`
	Messages Messages `json:"messages"`
}

// Registry holds all bundles keyed by locale.
type Registry struct {
	byLocale map[Locale]*Bundle
	fallback Locale
}

// LoadRegistry reads every bundles/<locale>.json file.
// Add a new language by creating bundles/<code>.json — no Go code changes required.
func LoadRegistry(fallback Locale) (*Registry, error) {
	reg := &Registry{
		byLocale: make(map[Locale]*Bundle),
		fallback: fallback,
	}

	entries, err := fs.ReadDir(bundleFS, "bundles")
	if err != nil {
		return nil, fmt.Errorf("i18n: read bundles dir: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		path := filepath.Join("bundles", entry.Name())
		data, err := bundleFS.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("i18n: read %s: %w", path, err)
		}
		var raw bundleFile
		if err := json.Unmarshal(data, &raw); err != nil {
			return nil, fmt.Errorf("i18n: parse %s: %w", path, err)
		}
		loc, ok := Parse(raw.Locale)
		if !ok {
			return nil, fmt.Errorf("i18n: unsupported locale in %s: %q", path, raw.Locale)
		}
		if raw.Dir != "rtl" && raw.Dir != "ltr" {
			return nil, fmt.Errorf("i18n: invalid dir in %s: %q (use rtl or ltr)", path, raw.Dir)
		}
		reg.byLocale[loc] = &Bundle{
			Locale:   loc,
			Dir:      raw.Dir,
			Label:    raw.Label,
			Messages: raw.Messages,
		}
	}

	if len(reg.byLocale) == 0 {
		return nil, fmt.Errorf("i18n: no bundle files found in bundles/")
	}
	if _, ok := reg.byLocale[fallback]; !ok {
		for loc := range reg.byLocale {
			reg.fallback = loc
			break
		}
	}
	return reg, nil
}

// Get returns the bundle for a locale, falling back when missing.
func (r *Registry) Get(loc Locale) *Bundle {
	if b, ok := r.byLocale[loc]; ok {
		return b
	}
	if b, ok := r.byLocale[r.fallback]; ok {
		return b
	}
	for _, b := range r.byLocale {
		return b
	}
	return nil
}

// Supported lists locales that have a bundle file.
func (r *Registry) Supported() []Locale {
	out := make([]Locale, 0, len(r.byLocale))
	for loc := range r.byLocale {
		out = append(out, loc)
	}
	return out
}

// MessagesFor returns messages for the locale.
func (r *Registry) MessagesFor(loc Locale) Messages {
	return r.Get(loc).Messages
}
