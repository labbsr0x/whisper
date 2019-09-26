package types

import "html/template"

type IPage interface {
	SetHTML(html template.HTML)
}

// Page defines the information expected from a page
type Page struct {
	HTML template.HTML // a page should have an HTML
}
