package api

import (
	"fmt"
	"net/http"

	"github.com/labbsr0x/goh/gohserver"
	"github.com/labbsr0x/goh/gohtypes"
	whisper "github.com/labbsr0x/whisper-client/client"
	"github.com/labbsr0x/whisper/db"
	"github.com/labbsr0x/whisper/misc"
	"github.com/labbsr0x/whisper/web/api/types"
	"github.com/labbsr0x/whisper/web/config"
)

// AppAPI defines the available user apis
type AppAPI interface {
	Init(w *config.WebBuilder) AppAPI
	POSTHandler() http.Handler
	GETHandler() http.Handler
}

// DefaultAppAPI holds the default implementation for the app api
type DefaultAppAPI struct {
	*config.WebBuilder
	appDAO db.AppDAO
}

var _dao = new(db.DefaultAppDAO)

// Init initializes the API
func (dapi *DefaultAppAPI) Init(w *config.WebBuilder) AppAPI {
	dapi.WebBuilder = w
	dapi.appDAO = _dao.Init(w.DB)

	return dapi
}

// POSTHandler defines the function to handle HTTP Post requests to create a new App
func (dapi *DefaultAppAPI) POSTHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload types.AddAppInitialRequestPayload
		err := misc.UnmarshalPayloadFromRequest(&payload, r)
		gohtypes.PanicIfError(fmt.Sprintf("Unable to unmarshal the request: %v", err), http.StatusBadRequest, err)

		if token, ok := r.Context().Value(whisper.TokenKey).(whisper.Token); ok {
			app, err := dapi.appDAO.InsertAppData(payload, token.Subject)
			gohtypes.PanicIfError("Unable to create App", http.StatusInternalServerError, err)
			gohserver.WriteJSONResponse(app, http.StatusCreated, w)
		}
	})
}

// PUTHandler defines the function to handle HTTP Put requests to update an existing App
func (dapi *DefaultAppAPI) PUTHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload types.UpdateAppRequestPayload
		err := misc.UnmarshalPayloadFromRequest(&payload, r)
		gohtypes.PanicIfError(fmt.Sprintf("Unable to unmarshal the request: %v", err), http.StatusBadRequest, err)

		if token, ok := r.Context().Value(whisper.TokenKey).(whisper.Token); ok {
			app, err := dapi.appDAO.UpdateAppData(payload, token.Subject)
			gohtypes.PanicIfError("Unable to update app", http.StatusInternalServerError, err)
			gohserver.WriteJSONResponse(app, http.StatusOK, w)
		}
	})
}

// GETHandler defines the function to handle HTTP Get requests
func (dapi *DefaultAppAPI) GETHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if token, ok := r.Context().Value(whisper.TokenKey).(whisper.Token); ok {
			list, err := dapi.appDAO.ListApps(token.Subject)
			gohtypes.PanicIfError(fmt.Sprintf("Unable to list the Apps from user '%s'", token.Subject), 500, err)

			gohserver.WriteJSONResponse(list, http.StatusCreated, w)
		}
	})
}
