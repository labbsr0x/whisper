package types

import (
	"fmt"
	"github.com/labbsr0x/whisper/misc"
	"html/template"
)

// ConsentPage defines the data needed to build a consent page
type ConsentPage struct {
	misc.BasePage
	ClientURI       string
	ClientName      string
	RequestedScopes []misc.GrantScope
}

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

func (payload *ConsentRequestPayload) Check() error {
	payload.Remember = true

	if len(payload.Challenge) == 0 {
		return fmt.Errorf("Incomplete form data")
	}
	return nil
}
