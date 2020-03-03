package api

import (
	"net/http"

	"github.com/labbsr0x/whisper/web/config"
)

// HomeAPI defines the available Home APIs
type HomeAPI interface {
	GETHandler() http.Handler
	InitFromWebBuilder(webBuilder *config.WebBuilder) HomeAPI
}

// DefaultHomeAPI holds the default implementation of the Home API
type DefaultHomeAPI struct {
	*config.WebBuilder
}

// InitFromWebBuilder defines the DefaultHomeAPI initialization logic from a constructed webBuilder struct
func (api *DefaultHomeAPI) InitFromWebBuilder(w *config.WebBuilder) HomeAPI {
	return api
}

// GETHandler defines the HTTP GET Handler for Home
func (api *DefaultHomeAPI) GETHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO
	})
}
