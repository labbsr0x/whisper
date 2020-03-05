package api

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/labbsr0x/goh/gohserver"
	"github.com/labbsr0x/goh/gohtypes"
	"github.com/labbsr0x/whisper/db"
	"github.com/labbsr0x/whisper/hydra"
	"github.com/labbsr0x/whisper/mail"
	"github.com/labbsr0x/whisper/misc"
	"github.com/labbsr0x/whisper/web/ui"
	"github.com/sirupsen/logrus"

	"github.com/labbsr0x/whisper/web/api/types"
	"github.com/labbsr0x/whisper/web/config"
)

// LoginAPI defines the available user apis
type LoginAPI interface {
	LoginGETHandler(route string) http.Handler
	LoginPOSTHandler() http.Handler
	PostLoginCallbackGETHandler() http.Handler
	InitFromWebBuilder(w *config.WebBuilder) LoginAPI
}

// DefaultLoginAPI holds the default implementation of the User API interface
type DefaultLoginAPI struct {
	*config.WebBuilder
	UserCredentialsDAO db.UserCredentialsDAO
}

// InitFromWebBuilder initializes a default login api instance
func (dapi *DefaultLoginAPI) InitFromWebBuilder(w *config.WebBuilder) LoginAPI {
	dapi.WebBuilder = w
	dapi.UserCredentialsDAO = new(db.DefaultUserCredentialsDAO).Init(w.SecretKey, w.BaseUIPath, w.PublicURL, w.Outbox, w.DB)
	return dapi
}

// LoginPOSTHandler post form handler for logging in users
func (dapi *DefaultLoginAPI) LoginPOSTHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload types.RequestLoginPayload

		err := misc.UnmarshalPayloadFromRequest(&payload, r)
		gohtypes.PanicIfError("Unable to unmarshal the request", http.StatusBadRequest, err)

		userCredential := dapi.UserCredentialsDAO.CheckCredentials(payload.Username, payload.Password)

		if !userCredential.EmailValidated {
			dapi.Outbox <- mail.GetEmailConfirmationMail(dapi.BaseUIPath, dapi.SecretKey, dapi.PublicURL, userCredential.Username, userCredential.Email, payload.Challenge)
			gohtypes.Panic("This account email is not authenticated, an email was sent to you confirm your email", http.StatusUnauthorized)
		}

		info := dapi.HydraHelper.AcceptLoginRequest(
			payload.Challenge,
			hydra.AcceptLoginRequestPayload{ACR: "0", Remember: payload.Remember, RememberFor: 3600, Subject: payload.Username},
		)
		logrus.Debugf("Accept login request info: %v", info)
		if info != nil {
			gohserver.WriteJSONResponse(map[string]interface{}{
				"redirect_to": info["redirect_to"],
			}, http.StatusOK, w)
			return
		}
	})
}

// LoginGETHandler prompts the browser to the login UI or redirects it to hydra
func (dapi *DefaultLoginAPI) LoginGETHandler(route string) http.Handler {
	return http.StripPrefix(route, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		challenge, err := url.QueryUnescape(r.URL.Query().Get("login_challenge"))
		if err == nil {
			info := dapi.HydraHelper.GetLoginRequestInfo(challenge)
			logrus.Debugf("Login Request Info: %v", info)
			if info["skip"].(bool) {
				subject := info["subject"].(string)
				info = dapi.HydraHelper.AcceptLoginRequest(
					challenge,
					hydra.AcceptLoginRequestPayload{Subject: subject},
				)
				if info != nil {
					logrus.Debugf("Login request skipped for subject '%v'", subject)
					http.Redirect(w, r, info["redirect_to"].(string), http.StatusFound)
				}
			} else {
				page := types.LoginPage{Challenge: challenge}
				ui.WritePage(w, dapi.BaseUIPath, ui.Login, &page)
			}
			return
		}
		panic(gohtypes.Error{Code: http.StatusBadRequest, Err: err, Message: "Unable to parse the login_challenge"})
	}))
}

// PostLoginCallbackGETHandler defines the logic for handling behavior in whisper after successful login
func (dapi *DefaultLoginAPI) PostLoginCallbackGETHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code, err := url.QueryUnescape(r.URL.Query().Get("code"))
		if err == nil && code != "" {
			// exchange the code retrieved by tokens that identify the user and more
			codeVerifierCookie, err := r.Cookie("CODE_VERIFIER")
			gohtypes.PanicIfError("Unable to exchange code for tokens", 500, err)
			stateCookie, err := r.Cookie("STATE")
			gohtypes.PanicIfError("Unable to exchange code for tokens", 500, err)

			tokens, err := dapi.Self.ExchangeCodeForToken(code, codeVerifierCookie.Value, stateCookie.Value)
			gohtypes.PanicIfError(fmt.Sprintf("Unable to exchange code '%v' for tokens: %s", code, tokens), 500, err)

			logrus.Infof("Exchanged code '%v' for the following tokens: %v", tokens)
			http.SetCookie(w, &http.Cookie{
				Name:  "ACCESS_TOKEN",
				Value: tokens.AccessToken,
			})
			http.SetCookie(w, &http.Cookie{
				Name:    "CODE_VERIFIER",
				Value:   "",
				Expires: time.Unix(0, 0),
			})
			http.SetCookie(w, &http.Cookie{
				Name:    "STATE",
				Value:   "",
				Expires: time.Unix(0, 0),
			})

			// Redirect to the specified url
			http.Redirect(w, r, "/home", http.StatusFound)
		}
	})
}
