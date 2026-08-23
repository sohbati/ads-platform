package model

// Category matches cdn/json/category.json entries.
type Category struct {
	ID     int    `json:"id"`
	Parent *int   `json:"parent"`
	Order  int    `json:"order"`
	Title  string `json:"title"`
	Slug   string `json:"slug"`
	Path   string `json:"path"` // root → … → self ids, e.g. "3,15"
	IsLeaf bool   `json:"isLeaf"`
	// AdsAttrsJSONSchemaTemplateName names the JSON Schema for this category's
	// ad attributes; null until a template is assigned.
	AdsAttrsJSONSchemaTemplateName *string `json:"adsAttrsJsonSchemaTemplateName"`
}
