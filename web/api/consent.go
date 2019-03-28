package api

import (
	"net/http"

	"github.com/abilioesteves/whisper/web/ui"

	"github.com/abilioesteves/whisper/misc"
)

// ConsentAPI defines the available user apis
type ConsentAPI interface {
	ConsentGETHandler() http.Handler
	ConsentPOSTHandler(w http.ResponseWriter, r *http.Request)
}

// DefaultConsentAPI holds the default implementation of the User API interface
type DefaultConsentAPI struct {
	HydraClient *misc.HydraClient
	BaseUIPath  string
}

// Init initializes a default consent api instance
func (api *DefaultConsentAPI) Init(hydraClient *misc.HydraClient, baseUIPath string) *DefaultConsentAPI {
	api.HydraClient = hydraClient
	api.BaseUIPath = baseUIPath
	return api
}

// ConsentPOSTHandler REST POST api handler for app authorization
func (api *DefaultConsentAPI) ConsentPOSTHandler(w http.ResponseWriter, r *http.Request) {
	// TODO accept/deny consent request with hydra and redirect to its redirect_to response
	w.WriteHeader(200)
}

// ConsentGETHandler redirects the browser appropriately given
func (api *DefaultConsentAPI) ConsentGETHandler() http.Handler {
	// Verify skip and accept accordingly
	return ui.Handler(api.BaseUIPath)
}
