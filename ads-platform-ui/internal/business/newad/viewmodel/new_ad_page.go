package viewmodel

import (
	"encoding/json"

	"ads-platform-ui/internal/core/cdn"
	"ads-platform-ui/internal/core/i18n"
)

type NewAdPage struct {
	i18n.Page
	CityID    int
	LoadError string
	Bootstrap Bootstrap
}

type Bootstrap struct {
	Locale      string           `json:"locale"`
	Mode        string           `json:"mode,omitempty"`
	AdID        int64            `json:"adId,omitempty"`
	CityID      int              `json:"cityId"`
	CitySlug    string           `json:"citySlug"`
	CityName    string           `json:"cityName"`
	MaxPictures int              `json:"maxPictures"`
	SuccessHref string           `json:"successHref"`
	Prefill     *Prefill         `json:"prefill,omitempty"`
	Categories  []cdn.Category   `json:"categories"`
	Schemas     []cdn.AttrSchema `json:"schemas"`
	Enums       json.RawMessage  `json:"enums"`
}

type Prefill struct {
	CategoryID   int             `json:"category_id"`
	Title        string          `json:"title"`
	Description  string          `json:"description"`
	PriceAmount  *int64          `json:"price_amount,omitempty"`
	PriceType    string          `json:"price_type,omitempty"`
	Neighborhood string          `json:"neighborhood,omitempty"`
	Attrs        json.RawMessage `json:"attrs,omitempty"`
	Media        []PrefillMedia  `json:"media,omitempty"`
}

type PrefillMedia struct {
	URL         string `json:"url,omitempty"`
	Thumb       string `json:"thumb,omitempty"`
	StoredURL   string `json:"stored_url,omitempty"`
	StoredThumb string `json:"stored_thumb,omitempty"`
}
