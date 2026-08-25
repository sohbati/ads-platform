package model

import "encoding/json"

// AttrSchema matches CDN GET /api/attr-schemas entries.
type AttrSchema struct {
	Name       string          `json:"name"`
	Title      string          `json:"title"`
	JSONSchema json.RawMessage `json:"jsonSchema"`
}
