package api

import (
	"github.com/labbsr0x/goh/gohtypes"
	"github.com/labbsr0x/whisper/web/config"
	"net/http"
	"net/url"
)

// LogoutAPI defines the available user apis
type LogoutAPI interface {
	LogoutGETHandler() http.Handler
}

// DefaultLogoutAPI holds the default implementation of the User API interface
type DefaultLogoutAPI struct {
	*config.WebBuilder
}

// InitFromWebBuilder initializes a default logout api instance
func (dapi *DefaultLogoutAPI) InitFromWebBuilder(w *config.WebBuilder) *DefaultLogoutAPI {
	dapi.WebBuilder = w
	return dapi
}

// LogoutGETHandler prompts the browser to the logout UI or redirects it to hydra
func (dapi *DefaultLogoutAPI) LogoutGETHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		challenge, err := url.QueryUnescape(r.URL.Query().Get("challenge"))
		gohtypes.PanicIfError("Unable to retrieve logout challenge", http.StatusBadRequest, err)

		info := dapi.HydraHelper.AcceptLogoutRequest(challenge)
		if info == nil {
			gohtypes.Panic("Unable to accept logout", http.StatusInternalServerError)
		}

		http.Redirect(w, r, info["redirect_to"].(string), http.StatusOK)
	})
}

