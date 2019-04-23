package api

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/labbsr0x/goh/gohtypes"

	"github.com/labbsr0x/whisper-client/hydra"
	"github.com/labbsr0x/whisper/web/middleware"

	"github.com/labbsr0x/whisper/db"
	"github.com/sirupsen/logrus"

	"github.com/labbsr0x/whisper/web/api/types"
	"github.com/labbsr0x/whisper/web/config"
)

// UserCredentialsAPI defines the available user apis
type UserCredentialsAPI interface {
	POSTHandler() http.Handler
	PUTHandler() http.Handler
	GETRegistrationPageHandler(route string) http.Handler
	GETUpdatePageHandler(route string) http.Handler
}

// DefaultUserCredentialsAPI holds the default implementation of the User API interface
type DefaultUserCredentialsAPI struct {
	*config.WebBuilder
	UserCredentialsDAO db.UserCredentialsDAO
}

// InitFromWebBuilder initializes the default user credentials API from a WebBuilder
func (dapi *DefaultUserCredentialsAPI) InitFromWebBuilder(builder *config.WebBuilder) *DefaultUserCredentialsAPI {
	dapi.WebBuilder = builder
	dapi.UserCredentialsDAO = new(db.DefaultUserCredentialsDAO).InitFromWebBuilder(builder)
	return dapi
}

// POSTHandler handles post requests to create user credentials
func (dapi *DefaultUserCredentialsAPI) POSTHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := new(types.AddUserCredentialRequestPayload).InitFromRequest(r)

		userID, err := dapi.UserCredentialsDAO.CreateUserCredential(payload.Username, payload.Password, payload.Email)
		gohtypes.PanicIfError("Not possible to create user", 500, err)
		logrus.Infof("User created: %v", userID)

		http.Redirect(w, r, fmt.Sprintf("/login?first_login=true&username=%v&login_challenge=%v", payload.Username, payload.LoginChallenge), 302)
	})
}

// PUTHandler handles put requests to update user credentials
func (dapi *DefaultUserCredentialsAPI) PUTHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := new(types.UpdateUserCredentialRequestPayload).InitFromRequest(r)
		logrus.Debugf("%v", payload)
	})
}

// GETRegistrationPageHandler builds the page where new credentials will be inserted
func (dapi *DefaultUserCredentialsAPI) GETRegistrationPageHandler(route string) http.Handler {
	return http.StripPrefix(route, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		challenge, err := url.QueryUnescape(r.URL.Query().Get("login_challenge"))
		gohtypes.PanicIfError("Unable to parse the login_challenge parameter", 400, err)
		page := types.RegistrationPage{LoginChallenge: challenge}

		buf := new(bytes.Buffer)
		template.Must(template.ParseFiles(path.Join(dapi.BaseUIPath, "registration.html"))).Execute(buf, page)
		html, _ := ioutil.ReadAll(buf)

		page.HTML = template.HTML(html)

		tmpl := template.Must(template.ParseFiles(path.Join(dapi.BaseUIPath, "index.html")))
		tmpl.Execute(w, page)
	}))
}

// GETUpdatePageHandler builds the page where credentials will be updated
func (dapi *DefaultUserCredentialsAPI) GETUpdatePageHandler(route string) http.Handler {
	return http.StripPrefix(route, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		redirectTo, err := url.QueryUnescape(r.URL.Query().Get("redirect_to"))
		gohtypes.PanicIfError("Unable to parse the redirect_to parameter", 400, err)
		page := types.UpdatePage{RedirectTo: redirectTo}

		if token, ok := r.Context().Value(middleware.TokenKey).(hydra.Token); ok {
			page.Username = token.Subject

			// TODO get usercredential information from DB; token.Subject will store the user id

			buf := new(bytes.Buffer)
			err := template.Must(template.ParseFiles(path.Join(dapi.BaseUIPath, "update.html"))).Execute(buf, page)
			gohtypes.PanicIfError("Error building update page", http.StatusInternalServerError, err)

			html, _ := ioutil.ReadAll(buf)
			page.HTML = template.HTML(html)

			template.Must(template.ParseFiles(path.Join(dapi.BaseUIPath, "index.html"))).Execute(w, page)
			return
		}
		gohtypes.Panic("Unauthorized: token not found", 401)
	}))
}
