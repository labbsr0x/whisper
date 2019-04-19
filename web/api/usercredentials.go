package api

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"net/http"
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
	dapi.UserCredentialsDAO = new(db.DefaultUserCredentialsDAO)
	return nil
}

// POSTHandler handles post requests to create user credentials
func (dapi *DefaultUserCredentialsAPI) POSTHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := new(types.AddUserCredentialRequestPayload).InitFromRequest(r)
		logrus.Debugf("%v", payload)
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
		page := types.RegistrationPage{}
		buf, err := ioutil.ReadFile(path.Join(dapi.BaseUIPath, "registration.html"))
		if err != nil {
			panic(err)
		}
		page.HTML = template.HTML(buf)

		tmpl := template.Must(template.ParseFiles(path.Join(dapi.BaseUIPath, "index.html")))
		tmpl.Execute(w, page)
	}))
}

// GETUpdatePageHandler builds the page where credentials will be updated
func (dapi *DefaultUserCredentialsAPI) GETUpdatePageHandler(route string) http.Handler {
	return http.StripPrefix(route, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := types.UpdatePage{}

		if token, ok := r.Context().Value(middleware.TokenKey).(hydra.Token); ok {
			page.Username = token.Subject
			buf := new(bytes.Buffer)
			err := template.Must(template.ParseFiles(path.Join(dapi.BaseUIPath, "update.html"))).Execute(buf, page)
			gohtypes.PanicIfError("Error building update page", http.StatusInternalServerError, err)

			err = template.Must(template.ParseFiles(path.Join(dapi.BaseUIPath, "index.html"))).Execute(w, page)
			gohtypes.PanicIfError("Error building update page", http.StatusInternalServerError, err)
		}
	}))
}
