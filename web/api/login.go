package api

import (
	"net/http"

	"github.com/abilioesteves/whisper/web/ui"

	"github.com/abilioesteves/whisper/misc"
)

// LoginAPI defines the available user apis
type LoginAPI interface {
	LoginGETHandler() http.Handler
	LoginPOSTHandler(w http.ResponseWriter, r *http.Request)
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
func (api *DefaultLoginAPI) LoginPOSTHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

// LoginGETHandler redirects the browser appropriately given
func (api *DefaultLoginAPI) LoginGETHandler() http.Handler {
	// TODO verify skip and accept accordingly
	return ui.Handler(api.BaseUIPath)
}
