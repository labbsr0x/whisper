package api

import (
	"github.com/labbsr0x/goh/gohserver"
	"github.com/labbsr0x/goh/gohtypes"
	whisper "github.com/labbsr0x/whisper-client/client"
	"github.com/labbsr0x/whisper/misc"
	"github.com/labbsr0x/whisper/web/api/types"
	"github.com/labbsr0x/whisper/web/config"
	"github.com/labbsr0x/whisper/web/ui"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
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
func (dapi *DefaultConsentAPI) InitFromWebBuilder(webBuilder *config.WebBuilder) *DefaultConsentAPI {
	dapi.WebBuilder = webBuilder
	return dapi
}

// ConsentPOSTHandler post form handler for app authorization
func (dapi *DefaultConsentAPI) ConsentPOSTHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		consentRequest := new(types.ConsentRequestPayload).InitFromRequest(r)
		logrus.Debugf("Consent request payload '%v'", consentRequest)
		if consentRequest.Accept {
			info := dapi.Self.GetConsentRequestInfo(consentRequest.Challenge)
			logrus.Debugf("Consent request info: '%v'", info)
			if info != nil {
				acceptInfo := dapi.Self.AcceptConsentRequest(
					consentRequest.Challenge,
					whisper.AcceptConsentRequestPayload{
						GrantAccessTokenAudience: misc.ConvertInterfaceArrayToStringArray(info["requested_access_token_audience"].([]interface{})),
						GrantScope:               consentRequest.GrantScope,
						Remember:                 consentRequest.Remember,
						RememberFor:              3600,
					})

				logrus.Debugf("Consent Accept Info: '%v'", acceptInfo)
				if acceptInfo != nil {
					gohserver.WriteJSONResponse(map[string]interface{}{
						"redirect_to": acceptInfo["redirect_to"],
					}, http.StatusOK, w)
					return
				}
			}
		} else {
			rejectInfo := dapi.Self.RejectConsentRequest(consentRequest.Challenge, whisper.RejectConsentRequestPayload{Error: "access_denied", ErrorDescription: "The resource owner denied the request"})
			if rejectInfo != nil {
				http.Redirect(w, r, rejectInfo["redirect_to"].(string), http.StatusFound)
				return
			}
		}
		panic(gohtypes.Error{Code: http.StatusInternalServerError, Message: "Unable to process consent request"})
	})
}

// ConsentGETHandler prompts the browser to the consent UI or redirects it to hydra
func (dapi *DefaultConsentAPI) ConsentGETHandler(route string) http.Handler {
	return http.StripPrefix(route, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		challenge, err := url.QueryUnescape(r.URL.Query().Get("consent_challenge"))
		gohtypes.PanicIfError("Unable to parse the consent_challenge parameter", http.StatusBadRequest, err)
		info := dapi.Self.GetConsentRequestInfo(challenge)
		logrus.Debugf("Consent Request Info: '%v'", info)
		if info["skip"].(bool) {
			info = dapi.Self.AcceptConsentRequest(
				challenge,
				whisper.AcceptConsentRequestPayload{
					GrantScope:               misc.ConvertInterfaceArrayToStringArray(info["requested_scope"].([]interface{})),
					GrantAccessTokenAudience: misc.ConvertInterfaceArrayToStringArray(info["requested_access_token_audience"].([]interface{}))},
			)

			if info != nil {
				logrus.Debugf("Consent request skipped for '%v'", info)
				http.Redirect(w, r, info["redirect_to"].(string), http.StatusFound)
			}
		} else {
			page := getConsentPage(info, dapi.GrantScopes)
			view := ui.BuildPage(dapi.BaseUIPath, ui.Consent, &page)
			Render(w, view)
		}
	}))
}

// getConsentPageInfo builds the data structure for a consent page
func getConsentPage(consentRequestInfo map[string]interface{}, scopes misc.GrantScopes) types.ConsentPage {
	consentPageInfo := types.ConsentPage{ClientName: "Unknown", ClientURI: "#", RequestedScopes: make([]misc.GrantScope, 0)}

	if clientName, ok := consentRequestInfo["client_name"].(string); ok {
		consentPageInfo.ClientName = clientName
	}

	if clientURI, ok := consentRequestInfo["client_uri"].(string); ok {
		consentPageInfo.ClientURI = clientURI
	}

	if i, ok := consentRequestInfo["requested_scope"].([]interface{}); ok {
		requestedScopes := misc.ConvertInterfaceArrayToStringArray(i)

		for _, scope := range requestedScopes {
			consentPageInfo.RequestedScopes = append(consentPageInfo.RequestedScopes, scopes[scope])
		}
	}

	return consentPageInfo
}
