package api

import (
	"html/template"
	"net/http"
	"net/url"
	"path"

	"github.com/sirupsen/logrus"

	"github.com/abilioesteves/goh/gohtypes"

	"github.com/abilioesteves/whisper/misc"
	"github.com/abilioesteves/whisper/web/api/types"
	"github.com/abilioesteves/whisper/web/config"
)

// LoginAPI defines the available user apis
type LoginAPI interface {
	LoginGETHandler(route string) http.Handler
	LoginPOSTHandler() http.Handler
}

// DefaultLoginAPI holds the default implementation of the User API interface
type DefaultLoginAPI struct {
	*config.WebBuilder
}

// InitFromWebBuilder initializes a default login api instance
func (api *DefaultLoginAPI) InitFromWebBuilder(webBuilder *config.WebBuilder) *DefaultLoginAPI {
	api.WebBuilder = webBuilder
	return api
}

// LoginPOSTHandler post form handler for logging in users
func (api *DefaultLoginAPI) LoginPOSTHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loginRequest := new(types.RequestLoginPayload).InitFromRequest(r)
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
