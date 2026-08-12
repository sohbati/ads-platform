package i18n

import "fmt"

// ResolveAPIError maps a backend error code to localized UI text.
// Patterns may include fmt-style verbs (e.g. %s) filled from params in order.
func ResolveAPIError(catalog map[string]string, code string, params []string, fallback string) string {
	pattern := ""
	if catalog != nil {
		pattern = catalog[code]
		if pattern == "" {
			pattern = catalog["_default"]
		}
	}
	if pattern == "" {
		pattern = fallback
	}
	if pattern == "" {
		return code
	}
	if len(params) == 0 {
		return pattern
	}

	args := make([]any, len(params))
	for i, param := range params {
		args[i] = param
	}
	return fmt.Sprintf(pattern, args...)
}
