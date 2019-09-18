package types

import (
	"encoding/json"
	"fmt"
	"github.com/labbsr0x/goh/gohtypes"
	"github.com/labbsr0x/whisper/misc"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

// ConsentPage defines the data needed to build a consent page
type ConsentPage struct {
	Page
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
	data, err := ioutil.ReadAll(r.Body)
	gohtypes.PanicIfError("Not possible to parse registration payload", 400, err)

	err = json.Unmarshal(data, &payload)
	logrus.Debugf("Payload: '%v' Error: %v", payload, err)

	payload.Remember = true
	payload.check()

	return payload
}

// check verifies if the consent payload is ok
func (payload *ConsentRequestPayload) check() error {
	if len(payload.Challenge) == 0 {
		return fmt.Errorf("Incomplete form data")
	}
	return nil
}
