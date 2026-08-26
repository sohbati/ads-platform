package cities

import (
	"strconv"
	"strings"
)

const (
	ContextKey          = "city"
	ContextLocationsKey = "locations"
	CookieName          = "city"
	LocationsCookieName = "locations"
)

// PopularSlugs are the default cities shown on the home page (order preserved).
var PopularSlugs = []string{
	"tehran", "mashhad", "karaj", "shiraz", "isfahan",
	"ahvaz", "tabriz", "kermanshah", "qom", "rasht",
}

// Catalog is the in-memory index of cities loaded from the CDN.
type Catalog struct {
	bySlug   map[string]Record
	byID     map[int]Record
	children map[int][]Record
	all      []Record
}

// Get returns a city by slug and whether it exists.
func (c *Catalog) Get(slug string) (Record, bool) {
	if c == nil {
		return Record{}, false
	}
	r, ok := c.bySlug[normalizeSlug(slug)]
	return r, ok
}

// GetByID returns a record by numeric id.
func (c *Catalog) GetByID(id int) (Record, bool) {
	if c == nil {
		return Record{}, false
	}
	r, ok := c.byID[id]
	return r, ok
}

// Children returns direct children of a parent id, in catalog order.
func (c *Catalog) Children(parentID int) []Record {
	if c == nil {
		return nil
	}
	return c.children[parentID]
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

// ParseLocationSlugs splits a cookie or query value into known slugs.
func (c *Catalog) ParseLocationSlugs(raw string) []string {
	if c == nil || strings.TrimSpace(raw) == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		slug := c.NormalizeSlug(part)
		if slug == "" {
			continue
		}
		if _, ok := c.Get(slug); !ok {
			continue
		}
		if _, dup := seen[slug]; dup {
			continue
		}
		seen[slug] = struct{}{}
		out = append(out, slug)
	}
	return out
}

// PrimaryCitySlug is the first type-2 city implied by the selection (province → first city).
func (c *Catalog) PrimaryCitySlug(slugs []string, fallback string) string {
	ids := c.ExpandToCityIDs(slugs)
	if len(ids) > 0 {
		if r, ok := c.GetByID(ids[0]); ok && r.Slug != "" {
			return r.Slug
		}
	}
	if fallback != "" {
		return c.NormalizeSlug(fallback)
	}
	return "tehran"
}

// ExpandToCityIDs turns selected provinces/cities/areas into type-2 city ids.
func (c *Catalog) ExpandToCityIDs(slugs []string) []int {
	if c == nil || len(slugs) == 0 {
		return nil
	}
	seen := make(map[int]struct{})
	out := make([]int, 0)
	add := func(id int) {
		if id == 0 {
			return
		}
		if _, ok := seen[id]; ok {
			return
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	for _, slug := range slugs {
		r, ok := c.Get(slug)
		if !ok {
			continue
		}
		switch r.Type {
		case TypeCity:
			add(r.ID)
		case TypeProvince:
			for _, child := range c.Children(r.ID) {
				if child.Type == TypeCity {
					add(child.ID)
				}
			}
		case TypeArea:
			if r.Parent != nil {
				if parent, ok := c.GetByID(*r.Parent); ok && parent.Type == TypeCity {
					add(parent.ID)
				}
			}
		}
	}
	return out
}

// SearchPlace maps a multi-location selection onto the search API place + cities filter.
// One city → that slug; several → iran with a CSV of city ids.
func (c *Catalog) SearchPlace(slugs []string, primary string) (place string, citiesCSV string) {
	ids := c.ExpandToCityIDs(slugs)
	if len(ids) == 0 {
		if primary != "" {
			return c.NormalizeSlug(primary), ""
		}
		return "iran", ""
	}
	if len(ids) == 1 {
		if r, ok := c.GetByID(ids[0]); ok && r.Slug != "" {
			return r.Slug, ""
		}
	}
	parts := make([]string, len(ids))
	for i, id := range ids {
		parts[i] = strconv.Itoa(id)
	}
	return "iran", strings.Join(parts, ",")
}

func normalizeSlug(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}
