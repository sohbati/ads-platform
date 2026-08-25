package model

import "encoding/json"

// FileEntry is one value in cdn/json/attr-schemas.json (keyed by template name).
type FileEntry struct {
	Title      string          `json:"title"`
	JSONSchema json.RawMessage `json:"jsonSchema"`
}

// AttrSchema is the CDN API shape: file entry plus the map key as Name.
type AttrSchema struct {
	Name       string          `json:"name"`
	Title      string          `json:"title"`
	JSONSchema json.RawMessage `json:"jsonSchema"`
}
