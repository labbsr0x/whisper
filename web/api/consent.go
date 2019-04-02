package api

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/abilioesteves/whisper/web/ui"

	"github.com/abilioesteves/goh/gohtypes"
	"github.com/abilioesteves/whisper/misc"
	"github.com/sirupsen/logrus"
)

// ConsentAPI defines the available user apis
type ConsentAPI interface {
	ConsentGETHandler(route string) http.Handler
	ConsentPOSTHandler() http.Handler
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
}

// Init initializes a default consent api instance
func (api *DefaultConsentAPI) Init(hydraClient *misc.HydraClient, baseUIPath string) *DefaultConsentAPI {
	api.HydraClient = hydraClient
	api.BaseUIPath = baseUIPath
	return api
}

// ConsentPOSTHandler REST POST api handler for app authorization
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

// ConsentGETHandler redirects the browser appropriately
func (api *DefaultConsentAPI) ConsentGETHandler(route string) http.Handler {
	return http.StripPrefix(route, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		challenge := r.URL.Query().Get("consent_challenge")
		info := api.HydraClient.GetConsentRequestInfo(challenge)
		if info["skip"].(bool) {
			info = api.HydraClient.AcceptConsentRequest(
				challenge,
				misc.AcceptConsentRequestPayload{GrantScope: info["requested_scope"].([]string), GrantAccessTokenAudience: info["requested_access_token_audience"].([]string)},
			)

			if info != nil {
				logrus.Debugf("Consent request skipped for '%v'", info)
				http.Redirect(w, r, info["redirect_to"].(string), 302)
			}
		} else {
			ui.Handler(api.BaseUIPath).ServeHTTP(w, r)
		}
	}))
}
