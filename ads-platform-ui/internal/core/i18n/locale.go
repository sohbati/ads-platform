package i18n

import "strings"

// Locale is a BCP 47 primary language tag matching bundles/<locale>.json.
type Locale string

const (
	FA Locale = "fa"
	EN Locale = "en"
)

// DefaultLocale is used when a city has no explicit mapping.
const DefaultLocale = FA

// Parse validates a locale string (e.g. "fa-IR" → fa).
func Parse(raw string) (Locale, bool) {
	raw = strings.ToLower(strings.TrimSpace(raw))
	if raw == "" {
		return "", false
	}
	if i := strings.IndexByte(raw, '-'); i > 0 {
		raw = raw[:i]
	}
	if len(raw) < 2 || len(raw) > 8 {
		return "", false
	}
	for _, r := range raw {
		if (r < 'a' || r > 'z') && r != '_' {
			return "", false
		}
	}
	return Locale(raw), true
}

// Lang returns the HTML lang attribute.
func (l Locale) Lang() string {
	return string(l)
}
