package api

import (
	"net/http"

	"github.com/labbsr0x/whisper/web/config"
)

// LogoutAPI defines the Logout APIs
type LogoutAPI interface {
	GETHandler() http.Handler
	InitFromWebBuilder(w *config.WebBuilder) LogoutAPI
}

// DefaultLogoutAPI defines the default implementation for the Logout APIs
type DefaultLogoutAPI struct {
}

// InitFromWebBuilder defines the initialization logic for the Default Logout API package
func (api *DefaultLogoutAPI) InitFromWebBuilder(w *config.WebBuilder) LogoutAPI {
	return api
}

// GETHandler defines the HTTP GET Handler for the logout API
func (api *DefaultLogoutAPI) GETHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO
	})
}
