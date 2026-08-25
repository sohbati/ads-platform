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
	// AdsAttrsJSONSchemaTemplateName is the key in attr-schemas.json for this
	// category's ad-attrs form; null on non-leaf categories.
	AdsAttrsJSONSchemaTemplateName *string `json:"adsAttrsJsonSchemaTemplateName"`
}
