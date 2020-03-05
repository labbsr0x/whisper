package api

import (
	"net/http"

	"github.com/labbsr0x/whisper/web/api/types"
	"github.com/labbsr0x/whisper/web/config"
	"github.com/labbsr0x/whisper/web/ui"
)

// HomeAPI defines the available Home APIs
type HomeAPI interface {
	GETHandler(path string) http.Handler
	InitFromWebBuilder(webBuilder *config.WebBuilder) HomeAPI
}

// DefaultHomeAPI holds the default implementation of the Home API
type DefaultHomeAPI struct {
	BaseUIPath string
}

// InitFromWebBuilder defines the DefaultHomeAPI initialization logic from a constructed webBuilder struct
func (api *DefaultHomeAPI) InitFromWebBuilder(w *config.WebBuilder) HomeAPI {
	api.BaseUIPath = w.BaseUIPath
	return api
}

// GETHandler defines the HTTP GET Handler for Home
func (api *DefaultHomeAPI) GETHandler(path string) http.Handler {
	return http.StripPrefix(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := types.HomePage{}
		ui.WritePage(w, api.BaseUIPath, ui.Home, &page)
	}))
}
