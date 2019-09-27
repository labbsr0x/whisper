package api

import (
	"bytes"
	"github.com/labbsr0x/goh/gohserver"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/labbsr0x/goh/gohtypes"
	whisper "github.com/labbsr0x/whisper-client/client"
	"github.com/sirupsen/logrus"

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
		req := new(types.RequestLoginPayload).InitFromRequest(r)
		logrus.Debugf("Login request payload '%v'", req)

		dapi.UserCredentialsDAO.CheckCredentials(req.Username, req.Password, req.Challenge)

		info := dapi.Self.AcceptLoginRequest(
			req.Challenge,
			whisper.AcceptLoginRequestPayload{ACR: "0", Remember: req.Remember, RememberFor: 3600, Subject: req.Username},
		)
		logrus.Debugf("Accept login request info: %v", info)
		if info != nil {
			gohserver.WriteJSONResponse(map[string]interface{}{
				"redirect_to": info["redirect_to"],
			}, 200, w)
			return
		}
	})
}

// LoginGETHandler prompts the browser to the login UI or redirects it to hydra
func (dapi *DefaultLoginAPI) LoginGETHandler(route string) http.Handler {
	return http.StripPrefix(route, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		challenge, err := url.QueryUnescape(r.URL.Query().Get("login_challenge"))
		if err == nil {
			info := dapi.Self.GetLoginRequestInfo(challenge)
			logrus.Debugf("Login Request Info: %v", info)
			if info["skip"].(bool) {
				subject := info["subject"].(string)
				info = dapi.Self.AcceptLoginRequest(
					challenge,
					whisper.AcceptLoginRequestPayload{Subject: subject},
				)
				if info != nil {
					logrus.Debugf("Login request skipped for subject '%v'", subject)
					http.Redirect(w, r, info["redirect_to"].(string), 302)
				}
			} else {
				templ, info := dapi.getLoginPageTemplateAndInfo(challenge)
				templ.Execute(w, info)
			}
			return
		}
		panic(gohtypes.Error{Code: 400, Err: err, Message: "Unable to parse the login_challenge"})
	}))
}

// getLoginPageTemplateAndInfo gets the login page html and its defining payload
func (dapi *DefaultLoginAPI) getLoginPageTemplateAndInfo(challenge string) (*template.Template, types.LoginPage) {
	loginPage := types.LoginPage{Challenge: challenge}

	buf := new(bytes.Buffer)
	template.Must(template.ParseFiles(path.Join(dapi.BaseUIPath, "login.html"))).Execute(buf, loginPage)
	html, _ := ioutil.ReadAll(buf)

	loginPage.HTML = template.HTML(html)

	return template.Must(template.ParseFiles(path.Join(dapi.BaseUIPath, "index.html"))), loginPage
}
