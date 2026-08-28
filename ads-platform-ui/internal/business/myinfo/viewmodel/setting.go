package viewmodel

import "ads-platform-ui/internal/core/i18n"

// SettingPage is the account settings / look-and-feel screen.
type SettingPage struct {
	i18n.Page
	Seas []SeaOption
}

// SeaOption is one swatch in the sea palette picker.
type SeaOption struct {
	ID   string
	Hex  string
	Name string
}

var seaOrder = []struct {
	ID  string
	Hex string
}{
	{"ruab", "#0d9488"},
	{"pale-blue", "#AFEEEE"},
	{"sky-blue", "#87CEEB"},
	{"aquamarine", "#7FFFD4"},
	{"aqua-blue", "#00FFFF"},
	{"cyan", "#00BFFF"},
	{"azure", "#007FFF"},
	{"deep-blue", "#0000FF"},
	{"turquoise", "#40E0D0"},
	{"deep-turquoise", "#00CED1"},
	{"teal", "#008080"},
	{"dark-cyan", "#008B8B"},
	{"cerulean", "#007BA7"},
	{"sea-green", "#2E8B57"},
	{"deep-sea-blue", "#005B7D"},
	{"steel-blue", "#4682B4"},
	{"royal-blue", "#4169E1"},
	{"cobalt-blue", "#0047AB"},
	{"ocean-blue", "#4F42B5"},
	{"indigo", "#4B0082"},
	{"navy-blue", "#000080"},
	{"midnight-blue", "#191970"},
}

// SeasFor returns palette swatches with localized names.
func SeasFor(t i18n.Messages) []SeaOption {
	out := make([]SeaOption, 0, len(seaOrder))
	for _, s := range seaOrder {
		name := t.Appearance.Reset
		if s.ID != "ruab" {
			name = t.Appearance.Seas[s.ID]
			if name == "" {
				name = s.ID
			}
		}
		out = append(out, SeaOption{ID: s.ID, Hex: s.Hex, Name: name})
	}
	return out
}
