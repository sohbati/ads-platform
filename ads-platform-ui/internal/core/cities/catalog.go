package cities

import "strings"

const (
	ContextKey = "city"
	CookieName = "city"
)

// PopularSlugs are the default cities shown on the home page (order preserved).
var PopularSlugs = []string{
	"tehran", "mashhad", "karaj", "shiraz", "isfahan",
	"ahvaz", "tabriz", "kermanshah", "qom", "rasht",
}

// Catalog is the in-memory index of cities loaded from the CDN.
type Catalog struct {
	bySlug map[string]Record
	all    []Record
}

// Get returns a city by slug and whether it exists.
func (c *Catalog) Get(slug string) (Record, bool) {
	r, ok := c.bySlug[normalizeSlug(slug)]
	return r, ok
}

// NormalizeSlug returns the canonical slug if known, otherwise a normalized input.
func (c *Catalog) NormalizeSlug(raw string) string {
	slug := normalizeSlug(raw)
	if _, ok := c.bySlug[slug]; ok {
		return slug
	}
	return slug
}

// Name returns the display name from the CDN (e.g. Persian).
func (c *Catalog) Name(slug string) string {
	if r, ok := c.Get(slug); ok {
		return r.Name
	}
	return slug
}

// LocaleCode returns the UI locale for a city slug (Iran catalog → fa).
func (c *Catalog) LocaleCode(slug string) string {
	if _, ok := c.Get(slug); ok {
		return "fa"
	}
	return "fa"
}

// PopularRecords returns type-2 city records for the home page chips.
func (c *Catalog) PopularRecords(limit int) []Record {
	if limit <= 0 {
		limit = 8
	}
	out := make([]Record, 0, limit)
	for _, slug := range PopularSlugs {
		if len(out) >= limit {
			break
		}
		r, ok := c.Get(slug)
		if !ok || r.Type != TypeCity {
			continue
		}
		out = append(out, r)
	}
	return out
}

// CitiesByType returns all records of a given type (e.g. TypeCity).
func (c *Catalog) CitiesByType(cityType string) []Record {
	out := make([]Record, 0)
	for _, r := range c.all {
		if r.Type == cityType {
			out = append(out, r)
		}
	}
	return out
}

func normalizeSlug(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}
