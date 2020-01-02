package api

import (
	"fmt"
	"net/http"
	"strings"

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
		var payload types.AddAppRequestPayload
		err := misc.UnmarshalPayloadFromRequest(&payload, r)
		gohtypes.PanicIfError(fmt.Sprintf("Unable to unmarshal the request: %v", err), http.StatusBadRequest, err)

		if token, ok := r.Context().Value(whisper.TokenKey).(whisper.Token); ok {
			err = dapi.appDAO.CreateApp(payload, token.Subject)
			gohtypes.PanicIfError("Unable to create App", 500, err)

			w.WriteHeader(http.StatusOK)
		}
	})
}

// GETHandler defines the function to handle HTTP Get requests
func (dapi *DefaultAppAPI) GETHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if token, ok := r.Context().Value(whisper.TokenKey).(whisper.Token); ok {
			list, err := dapi.appDAO.ListApps(token.Subject)
			gohtypes.PanicIfError(fmt.Sprintf("Unable to list the Apps from user '%s'", token.Subject), 500, err)

			toReturn := []types.App{}
			for _, app := range list {
				toReturn = append(toReturn, types.App{
					ID:                app.ClientID,
					Name:              app.ClientName,
					URL:               app.ClientURL,
					LoginRedirectURL:  app.LoginRedirectURL,
					LogoutRedirectURL: app.LogoutRedirectURL,
					GrantTypes:        strings.Split(app.GrantTypes, ","),
					Scopes:            misc.GrantScopes{}, // TODO recover scopes
				})
			}
			gohserver.WriteJSONResponse(toReturn, 200, w)
			return
		}
	})
}
