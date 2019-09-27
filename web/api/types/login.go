package types

import (
	"encoding/json"
	"fmt"
	"github.com/labbsr0x/goh/gohtypes"
	"github.com/labbsr0x/whisper/misc"
	"github.com/sirupsen/logrus"
	"html/template"
	"io/ioutil"
	"net/http"
)

// LoginPage defines the data needed to build a consent page
type LoginPage struct {
	Page
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

// InitFromRequest initializes the login request payload from an http request form
func (payload *RequestLoginPayload) InitFromRequest(r *http.Request) *RequestLoginPayload {
	data, err := ioutil.ReadAll(r.Body)
	gohtypes.PanicIfError("Not possible to parse registration payload", http.StatusBadRequest, err)

	json.Unmarshal(data, &payload)
	logrus.Debugf("Payload: '%v'", payload)

	payload.check()

	return payload
}

// check verifies if the login request payload is ok
func (payload *RequestLoginPayload) check() error {
	if len(payload.Challenge) == 0 || len(payload.Password) == 0 || len(payload.Username) == 0 {
		return fmt.Errorf("Incomplete fields")
	}

	return nil
}
