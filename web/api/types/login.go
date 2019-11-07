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

// check verifies if the login request payload is ok
func (payload *RequestLoginPayload) Check() error {
	if len(payload.Challenge) == 0 || len(payload.Password) == 0 || len(payload.Username) == 0 {
		return fmt.Errorf("Incomplete fields")
	}

	return nil
}
