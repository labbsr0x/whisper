package api

import (
	"net/http"

	"github.com/labbsr0x/whisper/db"

	"github.com/labbsr0x/whisper/web/config"
)

// UserCredentialsAPI defines the available user apis
type UserCredentialsAPI interface {
	POSTHandler() http.Handler
	PUTHandler() http.Handler
	GETPageHandler(route string) http.Handler
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
	return nil
}

// PUTHandler handles put requests to update user credentials
func (dapi *DefaultUserCredentialsAPI) PUTHandler() http.Handler {
	return nil
}

// GETPageHandler builds the page where new credentials will be inserted
func (dapi *DefaultUserCredentialsAPI) GETPageHandler(route string) http.Handler {
	return http.StripPrefix(route, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	}))
}
