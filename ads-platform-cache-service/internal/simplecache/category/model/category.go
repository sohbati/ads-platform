package model

// Category matches CDN /api/categories JSON entries.
type Category struct {
	ID                             int     `json:"id"`
	Parent                         *int    `json:"parent"`
	Order                          int     `json:"order"`
	Title                          string  `json:"title"`
	Slug                           string  `json:"slug"`
	Path                           string  `json:"path,omitempty"`
	IsLeaf                         bool    `json:"isLeaf"`
	AdsAttrsJSONSchemaTemplateName *string `json:"adsAttrsJsonSchemaTemplateName"`
	// DescendantIDs holds self + all descendant category ids.
	// Populated only when include_descendants=true is requested.
	DescendantIDs []int `json:"descendant_ids,omitempty"`
}
