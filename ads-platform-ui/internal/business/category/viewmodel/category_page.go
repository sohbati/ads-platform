package viewmodel

import "ads-platform-ui/internal/core/i18n"

// CategoryItem is one category link in the browse UI.
type CategoryItem struct {
	ID    int
	Title string
	Slug  string
	Href  string
	Icon  string
}

// CategoryColumn is one vertical column in the mega-menu panel.
type CategoryColumn struct {
	Title string
	Href  string
	Links []CategoryItem
}

// CategoryGroup is a top-level category with optional children.
type CategoryGroup struct {
	CategoryItem
	Children []CategoryItem
	Columns  []CategoryColumn
}

// CategoryPage is the view model for /category.
type CategoryPage struct {
	i18n.Page
	Groups    []CategoryGroup
	ActiveID  int
	LoadError string
}
