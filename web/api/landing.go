package api

import (
	"net/http"

	"github.com/labbsr0x/whisper-client/client"
	"github.com/labbsr0x/whisper/web/config"
)

// LandingAPI defines the Landing Page API
type LandingAPI interface {
	GETHandler() http.Handler
	InitFromWebBuilder(w *config.WebBuilder) LandingAPI
}

// DefaultLandingAPI defines the default landing API
type DefaultLandingAPI struct {
	wClient *client.WhisperClient
}

// InitFromWebBuilder defines the DefaultLandingAPI initialization logic from a WebBuilder struct
func (api *DefaultLandingAPI) InitFromWebBuilder(w *config.WebBuilder) LandingAPI {
	api.wClient = w.Self
	return api
}

// GETHandler defines the HTTP GET handler for the landing page
func (api *DefaultLandingAPI) GETHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loginURL, codeVerifier, state := api.wClient.GetOAuth2LoginParams()
		http.SetCookie(w, &http.Cookie{
			Name:  "CODE_VERIFIER",
			Value: codeVerifier,
		})
		http.SetCookie(w, &http.Cookie{
			Name:  "STATE",
			Value: state,
		})

		http.Redirect(w, r, loginURL, http.StatusFound)
	})
}
