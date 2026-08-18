package model

// City matches CDN /api/cities JSON entries.
type City struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Parent      *int   `json:"parent"`
	Type        string `json:"type"`
	CitiesCount *int   `json:"cities_count,omitempty"`
	Path        string `json:"path,omitempty"`
}
