package model

// SearchFilter is the resolved query used by repository/service.
type SearchFilter struct {
	PlaceSlug      string
	CategorySlug   string
	CategoryIDs    []int
	CityIDs        []int
	Query          string
	PriceMin       *int64
	PriceMax       *int64
	HasPhoto       *bool
	Neighborhood   string
	Sort           string // newest | cheapest | expensive
	Page           int
	Limit          int
}

// SearchResponse is a listing payload.
type SearchResponse struct {
	Place      string         `json:"place"`
	Category   string         `json:"category"`
	CategoryTitle string      `json:"category_title"`
	Filters    SearchFilters  `json:"filters"`
	Pagination Pagination     `json:"pagination"`
	Ads        []SearchAdItem `json:"ads"`
}

type SearchFilters struct {
	Cities       []int  `json:"cities,omitempty"`
	Query        string `json:"q,omitempty"`
	PriceMin     *int64 `json:"price_min,omitempty"`
	PriceMax     *int64 `json:"price_max,omitempty"`
	HasPhoto     *bool  `json:"has_photo,omitempty"`
	Neighborhood string `json:"neighborhood,omitempty"`
	Sort         string `json:"sort"`
}

type Pagination struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

// SearchAdItem is a list-card projection.
type SearchAdItem struct {
	ID           int64   `json:"id"`
	Title        string  `json:"title"`
	PriceAmount  *int64  `json:"price_amount"`
	PriceType    string  `json:"price_type"`
	Currency     string  `json:"currency"`
	CityID       int     `json:"city_id"`
	CityName     string  `json:"city_name,omitempty"`
	Neighborhood string  `json:"neighborhood,omitempty"`
	Thumbnail    string  `json:"thumbnail,omitempty"`
	HasPhoto     bool    `json:"has_photo"`
	CategoryID   int     `json:"category_id"`
	Slug         string  `json:"slug,omitempty"`
	PublishedAt  *string `json:"published_at,omitempty"`
}
