package api

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"path"

	"github.com/abilioesteves/goh/gohtypes"
	"github.com/abilioesteves/whisper/misc"
	"github.com/sirupsen/logrus"
)

// ConsentAPI defines the available user apis
type ConsentAPI interface {
	ConsentGETHandler(route string) http.Handler
	ConsentPOSTHandler() http.Handler
}

// ConsentPage defines the data needed to build a consent page
type ConsentPage struct {
	ClientURI       string
	ClientName      string
	RequestedScopes []GrantScope
}

// GrantScope defines the structure of a grant scope
type GrantScope struct {
	Description string
	Details     string
	Scope       string
}

// ConsentRequestPayload holds the data that defines a consent request to Whisper
type ConsentRequestPayload struct {
	Accept     bool
	Challenge  string
	GrantScope []string
	Remember   bool
}

// InitFromRequest initializes the consent payload from an http request
func (payload *ConsentRequestPayload) InitFromRequest(r *http.Request) *ConsentRequestPayload {
	err := r.ParseForm()
	if err == nil {
		logrus.Debugf("Form sent: '%v'", r.Form)
		if err := payload.check(r.Form); err == nil {
			payload.Accept = r.Form["accept"][0] == "true"
			payload.Challenge = r.Form["challenge"][0]
			payload.GrantScope = r.Form["grant-scope"]
			payload.Remember = len(r.Form["remember"]) > 0 && r.Form["remember"][0] == "on"

			return payload
		}
		panic(gohtypes.Error{Code: 400, Message: "Bad Request", Err: err})
	}
	panic(gohtypes.Error{Code: 400, Err: err, Message: "Not possible to parse http form"})
}

// check verifies if the consent payload is ok
func (payload *ConsentRequestPayload) check(form url.Values) error {
	if len(form["challenge"]) == 0 && len(form["accept"]) > 0 {
		return fmt.Errorf("Incomplete form data")
	}
	return nil
}

// DefaultConsentAPI holds the default implementation of the User API interface
type DefaultConsentAPI struct {
	HydraClient *misc.HydraClient
	BaseUIPath  string
	GrantScopes map[string]GrantScope
}

// Init initializes a default consent api instance
func (api *DefaultConsentAPI) Init(hydraClient *misc.HydraClient, baseUIPath string, grantScopes map[string]GrantScope) *DefaultConsentAPI {
	api.HydraClient = hydraClient
	api.BaseUIPath = baseUIPath
	api.GrantScopes = grantScopes
	return api
}

// ConsentPOSTHandler post form handler for app authorization
func (api *DefaultConsentAPI) ConsentPOSTHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		consentRequest := new(ConsentRequestPayload).InitFromRequest(r)
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
		panic(gohtypes.Error{Code: 500, Message: "Unable process consent request"})
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
				misc.AcceptConsentRequestPayload{GrantScope: info["requested_scope"].([]string), GrantAccessTokenAudience: misc.ConvertInterfaceArrayToStringArray(info["requested_access_token_audience"].([]interface{}))},
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

func (api *DefaultConsentAPI) getConsentPageInfo(consentRequestInfo map[string]interface{}) ConsentPage {
	toReturn := ConsentPage{ClientName: "Unknown", ClientURI: "#", RequestedScopes: make([]GrantScope, 0)}
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
