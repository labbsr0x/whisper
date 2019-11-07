package types

import (
	"fmt"
	"github.com/labbsr0x/whisper/misc"
	"html/template"
)

// LoginPage defines the data needed to build a consent page
type LoginPage struct {
	misc.BasePage
	ClientURI       string
	ClientName      string
	RequestedScopes []misc.GrantScope
	Challenge       string
}

// SetHTML exposes the HTML from base page
func (p *LoginPage) SetHTML(html template.HTML) {
	p.HTML = html
}

// RequestLoginPayload holds the data that defines a login request to Whisper
type RequestLoginPayload struct {
	Username  string
	Password  string
	Challenge string
	Remember  bool
}

// Check validates payload
func (payload *RequestLoginPayload) Check() error {
	if len(payload.Challenge) == 0 || len(payload.Password) == 0 || len(payload.Username) == 0 {
		return fmt.Errorf("incomplete fields")
	}

	return nil
}
