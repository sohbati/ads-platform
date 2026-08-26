package cities

// Record matches one entry in ads-platform-cdn/cdn/json/cities.json.
type Record struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Parent      *int   `json:"parent"`
	Type        string `json:"type"`
	CitiesCount *int   `json:"cities_count"`
}

// Type values used in cities.json.
const (
	TypeCountry  = "0"
	TypeProvince = "1"
	TypeCity     = "2"
	TypeArea     = "4"
)
