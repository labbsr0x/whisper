package api

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/labbsr0x/whisper-client/hydra"

	"github.com/sirupsen/logrus"

	"github.com/labbsr0x/goh/gohtypes"

	"github.com/labbsr0x/whisper/web/api/types"
	"github.com/labbsr0x/whisper/web/config"
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
func (dapi *DefaultLoginAPI) InitFromWebBuilder(webBuilder *config.WebBuilder) *DefaultLoginAPI {
	dapi.WebBuilder = webBuilder
	return dapi
}

// LoginPOSTHandler post form handler for logging in users
func (dapi *DefaultLoginAPI) LoginPOSTHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loginRequest := new(types.RequestLoginPayload).InitFromRequest(r)
		logrus.Debugf("Login request payload '%v'", loginRequest)
		if loginRequest.Password == "foobar" && loginRequest.Username == "foo@bar.com" { // TODO validation BL
			info := dapi.HydraClient.AcceptLoginRequest(
				loginRequest.Challenge,
				hydra.AcceptLoginRequestPayload{ACR: "0", Remember: loginRequest.Remember, RememberFor: 3600, Subject: loginRequest.Username},
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
func (dapi *DefaultLoginAPI) LoginGETHandler(route string) http.Handler {
	return http.StripPrefix(route, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		challenge, err := url.QueryUnescape(r.URL.Query().Get("login_challenge"))
		if err == nil {
			info := dapi.HydraClient.GetLoginRequestInfo(challenge)
			logrus.Debugf("Login Request Info: %v", info)
			if info["skip"].(bool) {
				subject := info["subject"].(string)
				info = dapi.HydraClient.AcceptLoginRequest(
					challenge,
					hydra.AcceptLoginRequestPayload{Subject: subject},
				)
				if info != nil {
					logrus.Debugf("Login request skipped for subject '%v'", subject)
					http.Redirect(w, r, info["redirect_to"].(string), 302)
				}
			} else {
				templ, info := dapi.getLoginPageTemplateAndInfo()
				templ.Execute(w, info)

			}
			return
		}
		panic(gohtypes.Error{Code: 500, Err: err, Message: "Unable to parse the login_challenge"})
	}))
}

// getLoginPageTemplateAndInfo gets the login page html and its defining payload
func (dapi *DefaultLoginAPI) getLoginPageTemplateAndInfo() (*template.Template, types.LoginPage) {
	loginPage := types.LoginPage{}
	pathToLoginHTML := path.Join(dapi.BaseUIPath, "login.html")
	buf, err := ioutil.ReadFile(pathToLoginHTML)
	if err != nil {
		panic(err)
	}

	loginPage.HTML = template.HTML(buf)

	return template.Must(template.ParseFiles(path.Join(dapi.BaseUIPath, "index.html"))), loginPage
}
