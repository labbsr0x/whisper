package types

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/abilioesteves/goh/gohtypes"
	"github.com/sirupsen/logrus"
)

// RequestLoginPayload holds the data that defines a login request to Whisper
type RequestLoginPayload struct {
	Username  string
	Password  string
	Challenge string
	Remember  bool
}

// InitFromRequest initializes the login request payload from an http request form
func (payload *RequestLoginPayload) InitFromRequest(r *http.Request) *RequestLoginPayload {
	err := r.ParseForm()
	if err == nil {
		logrus.Debugf("Form sent: '%v'", r.Form)
		if err := payload.check(r.Form); err == nil {
			payload.Challenge = r.Form["challenge"][0]
			payload.Password = r.Form["password"][0]
			payload.Username = r.Form["username"][0]
			payload.Remember = len(r.Form["remember"]) > 0 && r.Form["remember"][0] == "on"

			return payload
		}
		panic(gohtypes.Error{Code: 400, Message: "Bad Request", Err: err})
	}
	panic(gohtypes.Error{Code: 400, Message: "Not possible to parse http form", Err: err})
}

// check verifies if the login request payload is ok
func (payload *RequestLoginPayload) check(form url.Values) error {
	if len(form["challenge"]) == 0 || len(form["password"]) == 0 || len(form["username"]) == 0 {
		return fmt.Errorf("Incomplete form data")
	}

	return nil
}
