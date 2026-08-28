package media

import "strings"

// PublicURL turns a stored media path into a browser URL.
// Stored values are host-free ("/ads-media/ads/1/1_1.webp"). Absolute http(s)
// URLs from older rows are returned unchanged.
func PublicURL(base, path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return ""
	}
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}
	base = strings.TrimRight(strings.TrimSpace(base), "/")
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if base == "" {
		return path
	}
	return base + path
}
