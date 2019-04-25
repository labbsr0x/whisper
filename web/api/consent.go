package api

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/labbsr0x/whisper-client/hydra"

	"github.com/labbsr0x/goh/gohtypes"
	"github.com/labbsr0x/whisper/misc"
	"github.com/labbsr0x/whisper/web/api/types"
	"github.com/labbsr0x/whisper/web/config"
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
			info := dapi.HydraClient.GetConsentRequestInfo(consentRequest.Challenge)
			logrus.Debugf("Consent request info: '%v'", info)
			if info != nil {
				acceptInfo := dapi.HydraClient.AcceptConsentRequest(
					consentRequest.Challenge,
					hydra.AcceptConsentRequestPayload{
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
			rejectInfo := dapi.HydraClient.RejectConsentRequest(consentRequest.Challenge, hydra.RejectConsentRequestPayload{Error: "access_denied", ErrorDescription: "The resource owner denied the request"})
			if rejectInfo != nil {
				http.Redirect(w, r, rejectInfo["redirect_to"].(string), 302)
				return
			}
		}
		panic(gohtypes.Error{Code: 500, Message: "Unable to process consent request"})
	})
}

// ConsentGETHandler prompts the browser to the consent UI or redirects it to hydra
func (dapi *DefaultConsentAPI) ConsentGETHandler(route string) http.Handler {
	return http.StripPrefix(route, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		challenge, err := url.QueryUnescape(r.URL.Query().Get("consent_challenge"))
		gohtypes.PanicIfError("Unable to parse the consent_challenge parameter", 400, err)
		info := dapi.HydraClient.GetConsentRequestInfo(challenge)
		logrus.Debugf("Consent Request Info: '%v'", info)
		if info["skip"].(bool) {
			info = dapi.HydraClient.AcceptConsentRequest(
				challenge,
				hydra.AcceptConsentRequestPayload{
					GrantScope:               misc.ConvertInterfaceArrayToStringArray(info["requested_scope"].([]interface{})),
					GrantAccessTokenAudience: misc.ConvertInterfaceArrayToStringArray(info["requested_access_token_audience"].([]interface{}))},
			)

			if info != nil {
				logrus.Debugf("Consent request skipped for '%v'", info)
				http.Redirect(w, r, info["redirect_to"].(string), 302)
			}
		} else {
			templ, consentPageInfo := dapi.getConsentPageTemplateAndInfo(info, challenge)
			templ.Execute(w, consentPageInfo)
		}
	}))
}

// getConsentPageInfo builds the data structure for a consent page
func (dapi *DefaultConsentAPI) getConsentPageTemplateAndInfo(consentRequestInfo map[string]interface{}, challenge string) (*template.Template, types.ConsentPage) {
	consentPageInfo := types.ConsentPage{ClientName: "Unknown", ClientURI: "#", RequestedScopes: make([]misc.GrantScope, 0), Challenge: challenge}

	if clientName, ok := consentRequestInfo["client_name"].(string); ok {
		consentPageInfo.ClientName = clientName
	}

	if clientURI, ok := consentRequestInfo["client_uri"].(string); ok {
		consentPageInfo.ClientURI = clientURI
	}

	if i, ok := consentRequestInfo["requested_scope"].([]interface{}); ok {
		requestedScopes := misc.ConvertInterfaceArrayToStringArray(i)

		for _, scope := range requestedScopes {
			consentPageInfo.RequestedScopes = append(consentPageInfo.RequestedScopes, dapi.GrantScopes[scope])
		}
	}

	buf := new(bytes.Buffer)
	template.Must(template.ParseFiles(path.Join(dapi.BaseUIPath, "consent.html"))).Execute(buf, consentPageInfo)
	html, _ := ioutil.ReadAll(buf)

	consentPageInfo.HTML = template.HTML(html)

	return template.Must(template.ParseFiles(path.Join(dapi.BaseUIPath, "index.html"))), consentPageInfo
}
