package model

// Category matches cdn/json/category.json entries.
type Category struct {
	ID     int    `json:"id"`
	Parent *int   `json:"parent"`
	Order  int    `json:"order"`
	Title  string `json:"title"`
	Slug   string `json:"slug"`
}
