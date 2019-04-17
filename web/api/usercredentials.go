package api

import (
	"encoding/json"
	"net/http"

	"github.com/labbsr0x/whisper-client/hydra"

	"github.com/labbsr0x/whisper/web/middleware"

	"github.com/labbsr0x/whisper/db"

	"github.com/labbsr0x/whisper/web/config"

	"github.com/gorilla/mux"

	"github.com/labbsr0x/goh/gohtypes"

	"github.com/labbsr0x/goh/gohserver"
	"github.com/labbsr0x/whisper/web/api/types"
)

// UserCredentialsAPI defines the available user apis
type UserCredentialsAPI interface {
	AddUserCredentialHandler(w http.ResponseWriter, r *http.Request)
	RemoveUserCredentialHandler(w http.ResponseWriter, r *http.Request)
}

// DefaultUserCredentialsAPI holds the default implementation of the User API interface
type DefaultUserCredentialsAPI struct {
	*config.WebBuilder
	UserCredentialsDAO db.UserCredentialsDAO
}

// InitFromWebBuilder initializes the default user credentials API from a WebBuilder
func (api *DefaultUserCredentialsAPI) InitFromWebBuilder(builder *config.WebBuilder) *DefaultUserCredentialsAPI {
	api.UserCredentialsDAO = new(db.DefaultUserCredentialsDAO)
	return nil
}

// AddUserCredentialHandler REST POST api handler for adding new users
func (api *DefaultUserCredentialsAPI) AddUserCredentialHandler(w http.ResponseWriter, r *http.Request) {
	var payload types.AddUserCredentialRequestPayload
	var ucID string
	var err error

	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&payload); err != nil {
		panic(gohtypes.Error{Code: 400, Err: err, Message: "Unable to decode payload"})
	}

	if token, ok := r.Context().Value(middleware.TokenKey).(hydra.Token); ok {
		if ucID, err = api.UserCredentialsDAO.CreateUserCredential(payload.Username, payload.Password, token.ClientID); err == nil {
			gohserver.WriteJSONResponse(types.AddUserCredentialResponsePayload{UserCredentialID: ucID}, http.StatusOK, w)
		}
	}

	panic(err)
}

// RemoveUserCredentialHandler REST DELETE api handler for removing users
func (api *DefaultUserCredentialsAPI) RemoveUserCredentialHandler(w http.ResponseWriter, r *http.Request) {
	if userCredentialID := mux.Vars(r)["userCredentialID"]; len(userCredentialID) == 0 {
		panic(gohtypes.Error{Code: 400, Message: "Missing userCredentialID"})
	}

	gohserver.WriteJSONResponse("RemoveUserCredentialHandler: This is just a test", 200, w)
}
