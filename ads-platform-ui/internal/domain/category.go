package domain

// Category represents a top-level listing group on the marketplace home page.
type Category struct {
	ID          string
	Name        string
	Description string
	Icon        string // icon key resolved in templates
	Slug        string
	Href        string
}
