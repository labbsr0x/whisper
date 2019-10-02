package api

import (
	"fmt"
	"github.com/labbsr0x/goh/gohtypes"
	whisper "github.com/labbsr0x/whisper-client/client"
	"github.com/labbsr0x/whisper/mail"
	"github.com/labbsr0x/whisper/misc"
	"github.com/labbsr0x/whisper/resources"
	"github.com/labbsr0x/whisper/web/api/types"
	"github.com/labbsr0x/whisper/web/config"
	"github.com/labbsr0x/whisper/web/ui"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
)

// UserCredentialsAPI defines the available user apis
type UserCredentialsAPI interface {
	POSTHandler() http.Handler
	PUTHandler() http.Handler
	GETEmailConfirmationPageHandler(route string) http.Handler
	GETRegistrationPageHandler(route string) http.Handler
	GETUpdatePageHandler(route string) http.Handler
}

// DefaultUserCredentialsAPI holds the default implementation of the User API interface
type DefaultUserCredentialsAPI struct {
	*config.WebBuilder
}

// InitFromWebBuilder initializes the default user credentials API from a WebBuilder
func (dapi *DefaultUserCredentialsAPI) InitFromWebBuilder(builder *config.WebBuilder) *DefaultUserCredentialsAPI {
	dapi.WebBuilder = builder
	return dapi
}

// POSTHandler handles post requests to create user credentials
func (dapi *DefaultUserCredentialsAPI) POSTHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := new(types.AddUserCredentialRequestPayload).InitFromRequest(r)
		userID, err := dapi.UserCredentialsDAO.CreateUserCredential(payload.Username, payload.Password, payload.Email)
		gohtypes.PanicIfError("Not possible to create user", http.StatusInternalServerError, err)
		logrus.Infof("User created: %v", userID)

		to, content := misc.GetEmailConfirmationMail(payload.Username, payload.Email, payload.Challenge)

		resources.Outbox <- mail.Mail{To: to, Content: content}

		w.WriteHeader(http.StatusOK)
	})
}

// PUTHandler handles put requests to update user credentials
func (dapi *DefaultUserCredentialsAPI) PUTHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := new(types.UpdateUserCredentialRequestPayload).InitFromRequest(r)

		if token, ok := r.Context().Value(whisper.TokenKey).(whisper.Token); ok {
			dapi.UserCredentialsDAO.CheckCredentials(token.Subject, payload.OldPassword, "")

			err := dapi.UserCredentialsDAO.UpdateUserCredential(token.Subject, payload.Email, payload.NewPassword)
			gohtypes.PanicIfError("Error updating user credential info", http.StatusInternalServerError, err)

			w.WriteHeader(http.StatusOK)
		}
	})
}

// GETRegistrationPageHandler builds the page where new credentials will be inserted
func (dapi *DefaultUserCredentialsAPI) GETRegistrationPageHandler(route string) http.Handler {
	return http.StripPrefix(route, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		challenge, err := url.QueryUnescape(r.URL.Query().Get("login_challenge"))
		gohtypes.PanicIfError("Unable to parse the login_challenge parameter", http.StatusBadRequest, err)

		page := types.RegistrationPage{LoginChallenge: challenge}
		view := ui.BuildPage(ui.Registration, &page)
		Render(w, view)
	}))
}

func getRedirectionLink(challenge, username string, api *DefaultUserCredentialsAPI) string {
	if len(challenge) > 0 {
		payload := whisper.AcceptLoginRequestPayload{ACR: "0", Remember: false, Subject: username}
		info := api.Self.AcceptLoginRequest(challenge, payload)
		if info == nil {
			gohtypes.Panic("Unable to accept token login request", http.StatusInternalServerError)
		}

		return info["redirect_to"].(string)
	}

	return "/login"
}

// GETEmailConfirmationPageHandler builds the page where new credentials will be inserted
func (dapi *DefaultUserCredentialsAPI) GETEmailConfirmationPageHandler(route string) http.Handler {
	return http.StripPrefix(route, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		LoadErrorPage := func() {
			if rec := recover(); rec != nil {
				page := types.EmailConfirmationPage{Successful: false, Message: rec.(gohtypes.Error).Message}
				view := ui.BuildPage(ui.EmailConfirmation, &page)
				Render(w, view)
			}
		}

		defer LoadErrorPage()

		claims := misc.ExtractClaimsTokenFromRequest(r)
		username, challenge := misc.UnmarshalEmailConfirmationToken(claims)

		dapi.UserCredentialsDAO.AuthenticateUserCredential(username)

		link := getRedirectionLink(challenge, username, dapi)
		page := types.EmailConfirmationPage{Successful: true, Message: "Your email has been confirmed", RedirectTo: link}
		view := ui.BuildPage(ui.EmailConfirmation, &page)
		Render(w, view)
	}))
}

// GETUpdatePageHandler builds the page where credentials will be updated
func (dapi *DefaultUserCredentialsAPI) GETUpdatePageHandler(route string) http.Handler {
	return http.StripPrefix(route, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		redirectTo, err := url.QueryUnescape(r.URL.Query().Get("redirect_to"))
		gohtypes.PanicIfError("Unable to parse the redirect_to parameter", http.StatusBadRequest, err)

		if token, ok := r.Context().Value(whisper.TokenKey).(whisper.Token); ok {
			userCredentials, err := dapi.UserCredentialsDAO.GetUserCredential(token.Subject)
			gohtypes.PanicIfError(fmt.Sprintf("Could not find credentials with username '%v'", token.Subject), http.StatusInternalServerError, err)

			page := types.UpdatePage{RedirectTo: redirectTo, Username: userCredentials.Username, Email: userCredentials.Email}
			view := ui.BuildPage(ui.Update, &page)
			Render(w, view)

			return
		}
		gohtypes.Panic("Unauthorized: token not found", http.StatusUnauthorized)
	}))
}
