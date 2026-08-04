package model

// City matches cdn/json/cities.json entries.
type City struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Parent      *int   `json:"parent"`
	Type        string `json:"type"`
	CitiesCount *int   `json:"cities_count,omitempty"`
}
