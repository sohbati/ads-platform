package model

// Category matches CDN /api/categories JSON entries.
type Category struct {
	ID     int    `json:"id"`
	Parent *int   `json:"parent"`
	Order  int    `json:"order"`
	Title  string `json:"title"`
	Slug   string `json:"slug"`
	Path   string `json:"path,omitempty"`
}
