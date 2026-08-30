package viewmodel

import "ads-platform-ui/internal/core/i18n"

// SettingPage is the account settings / appearance screen.
type SettingPage struct {
	i18n.Page
	Themes []ThemeOption
}

// ThemeOption is one choice in the appearance picker.
type ThemeOption struct {
	ID   string
	Name string
}

// ThemesFor returns the product appearance options.
func ThemesFor(t i18n.Messages) []ThemeOption {
	return []ThemeOption{
		{ID: "light", Name: t.Appearance.Light},
		{ID: "tide", Name: t.Appearance.Tide},
		{ID: "dark", Name: t.Appearance.Dark},
	}
}
