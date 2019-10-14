package api

import (
	"github.com/labbsr0x/goh/gohserver"
	"github.com/labbsr0x/whisper/web/config"
	"net/http"
)

type HydraAPI interface {
	HydraGETHandler(route string) http.Handler
}

// DefaultHydraAPI holds the default implementation of the Hydra API interface
type DefaultHydraAPI struct {
	*config.WebBuilder
}

// InitFromWebBuilder initializes a default hydra api instance from a web builder instance
func (dapi *DefaultHydraAPI) InitFromWebBuilder(webBuilder *config.WebBuilder) *DefaultHydraAPI {
	dapi.WebBuilder = webBuilder
	return dapi
}

// HydraGETHandler provides the Hydra URLs
func (dapi *DefaultHydraAPI) HydraGETHandler(route string) http.Handler {
	return http.StripPrefix(route, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"hydraAdminUrl":  dapi.HydraAdminURL,
			"hydraPublicUrl": dapi.HydraPublicURL,
		}

		gohserver.WriteJSONResponse(data, http.StatusOK, w)
	}))
}
