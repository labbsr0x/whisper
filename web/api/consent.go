package api

import (
	"html/template"
	"net/http"
	"path"

	"github.com/abilioesteves/goh/gohtypes"
	"github.com/abilioesteves/whisper/misc"
	"github.com/abilioesteves/whisper/web/api/types"
	"github.com/abilioesteves/whisper/web/config"
	"github.com/sirupsen/logrus"
)

// ConsentAPI defines the available user apis
type ConsentAPI interface {
	ConsentGETHandler(route string) http.Handler
	ConsentPOSTHandler() http.Handler
}

// DefaultConsentAPI holds the default implementation of the User API interface
type DefaultConsentAPI struct {
	*config.WebBuilder
}

// InitFromWebBuilder initializes a default consent api instance from a web builder instance
func (api *DefaultConsentAPI) InitFromWebBuilder(webBuilder *config.WebBuilder) *DefaultConsentAPI {
	api.WebBuilder = webBuilder
	return api
}

// ConsentPOSTHandler post form handler for app authorization
func (api *DefaultConsentAPI) ConsentPOSTHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		consentRequest := new(types.ConsentRequestPayload).InitFromRequest(r)
		logrus.Debugf("Consent request payload '%v'", consentRequest)
		if consentRequest.Accept {
			info := api.HydraClient.GetConsentRequestInfo(consentRequest.Challenge)
			logrus.Debugf("Consent request info: '%v'", info)
			if info != nil {
				acceptInfo := api.HydraClient.AcceptConsentRequest(
					consentRequest.Challenge,
					misc.AcceptConsentRequestPayload{
						GrantAccessTokenAudience: misc.ConvertInterfaceArrayToStringArray(info["requested_access_token_audience"].([]interface{})),
						GrantScope:               consentRequest.GrantScope,
						Remember:                 consentRequest.Remember,
						RememberFor:              3600,
					})

				logrus.Debugf("Consent Accept Info: '%v'", acceptInfo)
				if acceptInfo != nil {
					http.Redirect(w, r, acceptInfo["redirect_to"].(string), 302)
					return
				}
			}
		} else {
			rejectInfo := api.HydraClient.RejectConsentRequest(consentRequest.Challenge, misc.RejectConsentRequestPayload{Error: "access_denied", ErrorDescription: "The resource owner denied the request"})
			if rejectInfo != nil {
				http.Redirect(w, r, rejectInfo["redirect_to"].(string), 302)
				return
			}
		}
		panic(gohtypes.Error{Code: 500, Message: "Unable to process consent request"})
	})
}

// ConsentGETHandler prompts the browser to the consent UI or redirects it to hydra
func (api *DefaultConsentAPI) ConsentGETHandler(route string) http.Handler {
	return http.StripPrefix(route, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		challenge := r.URL.Query().Get("consent_challenge")
		info := api.HydraClient.GetConsentRequestInfo(challenge)
		logrus.Debugf("Consent Request Info: '%v'", info)
		if info["skip"].(bool) {
			info = api.HydraClient.AcceptConsentRequest(
				challenge,
				misc.AcceptConsentRequestPayload{
					GrantScope:               misc.ConvertInterfaceArrayToStringArray(info["requested_scope"].([]interface{})),
					GrantAccessTokenAudience: misc.ConvertInterfaceArrayToStringArray(info["requested_access_token_audience"].([]interface{}))},
			)

			if info != nil {
				logrus.Debugf("Consent request skipped for '%v'", info)
				http.Redirect(w, r, info["redirect_to"].(string), 302)
			}
		} else {
			templ := template.Must(template.ParseFiles(path.Join(api.BaseUIPath, "index.html")))
			templ.Execute(w, api.getConsentPageInfo(info))
		}
	}))
}

// getConsentPageInfo builds the data structure for a consent page
func (api *DefaultConsentAPI) getConsentPageInfo(consentRequestInfo map[string]interface{}) types.ConsentPage {
	toReturn := types.ConsentPage{ClientName: "Unknown", ClientURI: "#", RequestedScopes: make([]config.GrantScope, 0)}
	if clientName, ok := consentRequestInfo["client_name"].(string); ok {
		toReturn.ClientName = clientName
	}

	if clientURI, ok := consentRequestInfo["client_uri"].(string); ok {
		toReturn.ClientURI = clientURI
	}

	if i, ok := consentRequestInfo["requested_scope"].([]interface{}); ok {
		requestedScopes := misc.ConvertInterfaceArrayToStringArray(i)

		for _, scope := range requestedScopes {
			toReturn.RequestedScopes = append(toReturn.RequestedScopes, api.GrantScopes[scope])
		}
	}

	logrus.Debugf("Consent page info: '%v'", toReturn)
	return toReturn
}
