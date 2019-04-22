package types

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/labbsr0x/goh/gohtypes"
	"github.com/sirupsen/logrus"
)

// RegistrationPage defines the information needed to load a registration page
type RegistrationPage struct {
	Page
	LoginChallenge string
}

// UpdatePage defins the information needed to load a update user credentials page
type UpdatePage struct {
	Page
	Username   string
	RedirectTo string
}

// AddUserCredentialRequestPayload defines the payload for adding a user
type AddUserCredentialRequestPayload struct {
	Email                string
	Username             string
	Password             string
	PasswordConfirmation string
	LoginChallenge       string
}

// AddUserCredentialResponsePayload defines the response payload after adding a user
type AddUserCredentialResponsePayload struct {
	UserCredentialID string
}

// UpdateUserCredentialRequestPayload defines the payload for updating a user
type UpdateUserCredentialRequestPayload struct {
	Email                   string
	NewPassword             string
	NewPasswordConfirmation string
	OldPassword             string
}

// InitFromRequest initializes the login request payload from an http request form
func (payload *AddUserCredentialRequestPayload) InitFromRequest(r *http.Request) *AddUserCredentialRequestPayload {
	err := r.ParseForm()
	if err == nil {
		logrus.Debugf("Form sent: '%v'", r.Form)
		if err := payload.check(r.Form); err == nil {
			payload.Email = r.Form["email"][0]
			payload.PasswordConfirmation = r.Form["password-confirmation"][0]
			payload.Password = r.Form["password"][0]
			payload.Username = r.Form["username"][0]
			payload.LoginChallenge = r.Form["login-challenge"][0]
			return payload
		}
		panic(gohtypes.Error{Code: 400, Message: "Bad Request", Err: err})
	}
	panic(gohtypes.Error{Code: 400, Message: "Not possible to parse http form", Err: err})
}

// check verifies if the login request payload is ok
func (payload *AddUserCredentialRequestPayload) check(form url.Values) error {
	if len(form["username"]) == 0 || len(form["password"]) == 0 || len(form["password-confirmation"]) == 0 || len(form["email"]) == 0 || len(form["login-challenge"]) == 0 {
		return errors.New("all form fields are required")
	}

	if form["password"][0] != form["password-confirmation"][0] {
		return errors.New("wrong password confirmation")
	}

	return nil
}

// InitFromRequest initializes the update request payload from an http request form
func (payload *UpdateUserCredentialRequestPayload) InitFromRequest(r *http.Request) *UpdateUserCredentialRequestPayload {
	err := r.ParseForm()
	if err == nil {
		logrus.Debugf("Form sent: '%v'", r.Form)
		if err := payload.check(r.Form); err == nil {
			payload.Email = r.Form["email"][0]
			payload.NewPasswordConfirmation = r.Form["new-password-confirmation"][0]
			payload.NewPassword = r.Form["new-password"][0]
			payload.OldPassword = r.Form["old-password"][0]
			return payload
		}
		panic(gohtypes.Error{Code: 400, Message: "Bad Request", Err: err})
	}
	panic(gohtypes.Error{Code: 400, Message: "Not possible to parse http form", Err: err})
}

// check verifies if the login request payload is ok
func (payload *UpdateUserCredentialRequestPayload) check(form url.Values) error {
	if len(form["old-password"]) == 0 || len(form["new-password"]) == 0 || len(form["new-password-confirmation"]) == 0 || len(form["email"]) == 0 {
		return errors.New("all fields must not be empty")
	}

	if form["new-password"][0] != form["new-password-confirmation"][0] {
		return errors.New("wrong password confirmation")
	}

	return nil
}
