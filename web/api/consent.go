package api

import (
	"encoding/json"
	"net/http"

	"github.com/abilioesteves/goh/gohtypes"
	"github.com/abilioesteves/whisper/misc"
	"github.com/gorilla/mux"
)

// ConsentAPI defines the available user apis
type ConsentAPI interface {
	ConsentGETHandler() http.Handler
	ConsentPOSTHandler(w http.ResponseWriter, r *http.Request)
}

// ConsentRequestPayload holds the data that defines a consent request to Whisper
type ConsentRequestPayload struct {
	Accept     bool     `json:"accept"`
	Challenge  string   `json:"challenge"`
	GrantScope []string `json:"grant_scope"`
	Remember   bool
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
	var consentRequest ConsentRequestPayload
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&consentRequest)
	if err == nil {
		if consentRequest.Accept {
			info := api.HydraClient.GetConsentRequestInfo(consentRequest.Challenge)
			if info != nil {
				acceptInfo := api.HydraClient.AcceptConsentRequest(
					consentRequest.Challenge,
					misc.AcceptConsentRequestPayload{
						GrantAccessTokenAudience: info["requested_access_token_audience"].(string),
						GrantScope:               consentRequest.GrantScope,
						Remember:                 consentRequest.Remember,
						RememberFor:              3600,
					})

				if acceptInfo != nil {
					http.Redirect(w, r, acceptInfo["redirec_to"].(string), 302)
				}
			}

		} else {
			rejectInfo := api.HydraClient.RejectConsentRequest(consentRequest.Challenge, misc.RejectConsentRequestPayload{Error: "access_denied", ErrorDescription: "The resource owner denied the request"})
			if rejectInfo != nil {
				http.Redirect(w, r, rejectInfo["redirect_to"].(string), 302)
			}
		}
	}
	panic(gohtypes.Error{Err: err, Code: 400, Message: "Unable to read request login payload."})
}

// ConsentGETHandler redirects the browser appropriately
func (api *DefaultConsentAPI) ConsentGETHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		challenge := vars["consent_challenge"]
		info := api.HydraClient.GetConsentRequestInfo(challenge)
		if info["skip"].(bool) {
			info = api.HydraClient.AcceptConsentRequest(
				challenge,
				misc.AcceptConsentRequestPayload{GrantScope: info["requested_scope"].([]string), GrantAccessTokenAudience: info["requested_access_token_audience"]},
			)

			if info != nil {
				http.Redirect(w, r, info["redirect_to"].(string), 302)
			}
		} else {
			http.ServeFile(w, r, api.BaseUIPath)
		}
	})
}
