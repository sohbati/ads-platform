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
	CityID      int              `json:"cityId"`
	CitySlug    string           `json:"citySlug"`
	CityName    string           `json:"cityName"`
	MaxPictures int              `json:"maxPictures"`
	SuccessHref string           `json:"successHref"`
	Categories  []cdn.Category   `json:"categories"`
	Schemas     []cdn.AttrSchema `json:"schemas"`
	Enums       json.RawMessage  `json:"enums"`
}
