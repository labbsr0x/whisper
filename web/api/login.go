package api

import (
	"encoding/json"
	"net/http"

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
	Username  string `json:"username"`
	Password  string `json:"password"`
	Challenge string `json:"challenge"`
	Remember  bool   `json:"remember"`
}

// DefaultLoginAPI holds the default implementation of the User API interface
type DefaultLoginAPI struct {
	HydraClient *misc.HydraClient
	BaseUIPath  string
}

// Init initializes a default login api instance
func (api *DefaultLoginAPI) Init(hydraClient *misc.HydraClient, baseUIPath string) *DefaultLoginAPI {
	api.HydraClient = hydraClient
	api.BaseUIPath = baseUIPath
	return api
}

// LoginPOSTHandler REST POST api handler for logging in users
func (api *DefaultLoginAPI) LoginPOSTHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var loginRequest LoginRequestPayload
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&loginRequest)
		if err == nil {
			logrus.Infof("Login request. Payload '%v'", loginRequest)
			if loginRequest.Password == "foobar" && loginRequest.Username == "foo@bar.com" { // TODO validation BL
				info := api.HydraClient.AcceptLoginRequest(
					loginRequest.Challenge,
					misc.AcceptLoginRequestPayload{ACR: "0", Remember: loginRequest.Remember, RememberFor: 3600, Subject: loginRequest.Username},
				)
				if info != nil {
					http.Redirect(w, r, info["redirect_to"].(string), 302)
				}
			}
			panic(gohtypes.Error{Code: 403, Message: "Unable to authenticate user"})
		}
		panic(gohtypes.Error{Err: err, Code: 400, Message: "Unable to read request login payload."})
	})
}

// LoginGETHandler redirects the browser appropriately given
func (api *DefaultLoginAPI) LoginGETHandler(route string) http.Handler {
	return http.StripPrefix(route, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		challenge := r.URL.Query().Get("login_challenge")
		logrus.Infof("challenge: %v", challenge)
		info := api.HydraClient.GetLoginRequestInfo(challenge)
		if info["skip"].(bool) {
			subject := info["subject"].(string)
			info = api.HydraClient.AcceptLoginRequest(
				challenge,
				misc.AcceptLoginRequestPayload{Subject: subject},
			)
			if info != nil {
				logrus.Infof("Login request skipped for subject '%v'", subject)
				http.Redirect(w, r, info["redirect_to"].(string), 302)
			}
		} else {
			http.ServeFile(w, r, api.BaseUIPath)
		}
	}))
}
