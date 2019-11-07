package misc

import "html/template"

// IPage defines what a page should expose
type IPage interface {
	SetHTML(html template.HTML)
}

// BasePage holds the basic information from a page
type BasePage struct {
	HTML template.HTML // a page should have an HTML
}

// SetHTML exposes the attribute HTML
func (p *BasePage) SetHTML(html template.HTML) {
	p.HTML = html
}
