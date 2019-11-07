package misc

import "html/template"

type IPage interface {
	SetHTML(html template.HTML)
}

// BasePage defines the information expected from a page
type BasePage struct {
	HTML template.HTML // a page should have an HTML
}

func (p *BasePage) SetHTML(html template.HTML) {
	p.HTML = html
}
