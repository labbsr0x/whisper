package types

import (
	"fmt"
	"html/template"

	"github.com/labbsr0x/whisper/misc"
)

// ConsentPage defines the data needed to build a consent page
type ConsentPage struct {
	misc.BasePage
	ClientURI       string
	ClientName      string
	RequestedScopes []misc.GrantScope
}

// SetHTML exposes the HTML from base page
func (p *ConsentPage) SetHTML(html template.HTML) {
	p.HTML = html
}

// ConsentRequestPayload holds the data that defines a consent request to Whisper
type ConsentRequestPayload struct {
	Accept     bool
	Challenge  string
	GrantScope []string
	Remember   bool
}

// Check validates payload
func (payload *ConsentRequestPayload) Check() error {
	payload.Remember = true

	if len(payload.Challenge) == 0 {
		return fmt.Errorf("there must be a challenge")
	}
	return nil
}
