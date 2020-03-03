package api

import (
	"net/http"

	"github.com/labbsr0x/goh/gohserver"
	"github.com/labbsr0x/whisper/web/config"
)

// HydraAPI hydra info interface
type HydraAPI interface {
	HydraGETHandler() http.Handler
	InitFromWebBuilder(webBuilder *config.WebBuilder) HydraAPI
}

// DefaultHydraAPI holds the default implementation of the Hydra API interface
type DefaultHydraAPI struct {
	*config.WebBuilder
}

// InitFromWebBuilder initializes a default hydra api instance from a web builder instance
func (dapi *DefaultHydraAPI) InitFromWebBuilder(webBuilder *config.WebBuilder) HydraAPI {
	dapi.WebBuilder = webBuilder
	return dapi
}

// HydraGETHandler provides the Hydra URLs
func (dapi *DefaultHydraAPI) HydraGETHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"hydraAdminUrl":  dapi.HydraAdminURL,
			"hydraPublicUrl": dapi.HydraPublicURL,
		}

		gohserver.WriteJSONResponse(data, http.StatusOK, w)
	})
}
