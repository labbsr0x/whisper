package types

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/labbsr0x/goh/gohtypes"
	"github.com/labbsr0x/whisper/misc"
	"github.com/sirupsen/logrus"
)

// ConsentPage defines the data needed to build a consent page
type ConsentPage struct {
	ClientURI       string
	ClientName      string
	RequestedScopes []misc.GrantScope
}

// ConsentRequestPayload holds the data that defines a consent request to Whisper
type ConsentRequestPayload struct {
	Accept     bool
	Challenge  string
	GrantScope []string
	Remember   bool
}

// InitFromRequest initializes the consent payload from an http request
func (payload *ConsentRequestPayload) InitFromRequest(r *http.Request) *ConsentRequestPayload {
	err := r.ParseForm()
	if err == nil {
		logrus.Debugf("Form sent: '%v'", r.Form)
		if err := payload.check(r.Form); err == nil {
			payload.Accept = r.Form["accept"][0] == "true"
			payload.Challenge = r.Form["challenge"][0]
			payload.GrantScope = r.Form["grant-scope"]
			payload.Remember = true

			return payload
		}
		panic(gohtypes.Error{Code: 400, Message: "Bad Request", Err: err})
	}
	panic(gohtypes.Error{Code: 400, Err: err, Message: "Not possible to parse http form"})
}

// check verifies if the consent payload is ok
func (payload *ConsentRequestPayload) check(form url.Values) error {
	if len(form["challenge"]) == 0 && len(form["accept"]) > 0 {
		return fmt.Errorf("Incomplete form data")
	}
	return nil
}
