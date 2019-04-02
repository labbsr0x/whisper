package api

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"path"

	"github.com/sirupsen/logrus"

	"github.com/abilioesteves/goh/gohtypes"

	"github.com/abilioesteves/whisper/misc"
)

// LoginAPI defines the available user apis
type LoginAPI interface {
	LoginGETHandler(route string) http.Handler
	LoginPOSTHandler() http.Handler
}

// LoginRequestPayload holds the data that defines a login request to Whisper
type LoginRequestPayload struct {
	Username  string
	Password  string
	Challenge string
	Remember  bool
}

// DefaultLoginAPI holds the default implementation of the User API interface
type DefaultLoginAPI struct {
	HydraClient *misc.HydraClient
	BaseUIPath  string
}

// InitFromRequest initializes the login request payload from an http request form
func (payload *LoginRequestPayload) InitFromRequest(r *http.Request) *LoginRequestPayload {
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
func (payload *LoginRequestPayload) check(form url.Values) error {
	if len(form["challenge"]) == 0 || len(form["password"]) == 0 || len(form["username"]) == 0 {
		return fmt.Errorf("Incomplete form data")
	}

	return nil
}

// Init initializes a default login api instance
func (api *DefaultLoginAPI) Init(hydraClient *misc.HydraClient, baseUIPath string) *DefaultLoginAPI {
	api.HydraClient = hydraClient
	api.BaseUIPath = baseUIPath
	return api
}

// LoginPOSTHandler post form handler for logging in users
func (api *DefaultLoginAPI) LoginPOSTHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loginRequest := new(LoginRequestPayload).InitFromRequest(r)
		logrus.Debugf("Login request payload '%v'", loginRequest)
		if loginRequest.Password == "foobar" && loginRequest.Username == "foo@bar.com" { // TODO validation BL
			info := api.HydraClient.AcceptLoginRequest(
				loginRequest.Challenge,
				misc.AcceptLoginRequestPayload{ACR: "0", Remember: loginRequest.Remember, RememberFor: 3600, Subject: loginRequest.Username},
			)
			logrus.Debugf("Accept login request info: %v", info)
			if info != nil {
				http.Redirect(w, r, info["redirect_to"].(string), 302)
				return
			}

		}
		panic(gohtypes.Error{Code: 403, Message: "Unable to authenticate user"})
	})
}

// LoginGETHandler prompts the browser to the login UI or redirects it to hydra
func (api *DefaultLoginAPI) LoginGETHandler(route string) http.Handler {
	return http.StripPrefix(route, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		challenge, err := url.QueryUnescape(r.URL.Query().Get("login_challenge"))
		if err == nil {
			info := api.HydraClient.GetLoginRequestInfo(challenge)
			logrus.Debugf("Login Request Info: %v", info)
			if info["skip"].(bool) {
				subject := info["subject"].(string)
				info = api.HydraClient.AcceptLoginRequest(
					challenge,
					misc.AcceptLoginRequestPayload{Subject: subject},
				)
				if info != nil {
					logrus.Debugf("Login request skipped for subject '%v'", subject)
					http.Redirect(w, r, info["redirect_to"].(string), 302)
				}
			} else {
				templ := template.Must(template.ParseFiles(path.Join(api.BaseUIPath, "index.html")))
				templ.Execute(w, nil)
			}
			return
		}
		panic(gohtypes.Error{Code: 500, Err: err, Message: "Unable to parse the login_challenge"})
	}))
}
