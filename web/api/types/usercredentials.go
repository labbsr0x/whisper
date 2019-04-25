package types

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"

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
	Email      string
	Token      string
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
	Email                   string `json:"email"`
	NewPassword             string `json:"newPassword"`
	NewPasswordConfirmation string `json:"newPasswordConfirmation"`
	OldPassword             string `json:"oldPassword"`
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
		gohtypes.Panic(err.Error(), 400)
	}
	panic(gohtypes.Error{Code: 400, Message: "Not possible to parse http form", Err: err})
}

// check verifies if the login request payload is ok
func (payload *AddUserCredentialRequestPayload) check(form url.Values) error {
	if len(form["username"]) == 0 || len(form["password"]) == 0 || len(form["password-confirmation"]) == 0 || len(form["email"]) == 0 || len(form["login-challenge"]) == 0 {
		return errors.New("All form fields are required")
	}

	if form["password"][0] != form["password-confirmation"][0] {
		return errors.New("Wrong password confirmation")
	}

	return verifyEmail(form["email"][0])
}

// InitFromRequest initializes the update request payload from an http request form
func (payload *UpdateUserCredentialRequestPayload) InitFromRequest(r *http.Request) *UpdateUserCredentialRequestPayload {
	data, err := ioutil.ReadAll(r.Body)
	gohtypes.PanicIfError("Not possible to parse update POST payload", 400, err)

	json.Unmarshal(data, &payload)
	logrus.Debugf("Payload: '%v'", payload)

	err = payload.check()
	gohtypes.PanicIfError(err.Error(), 400, err)

	return payload
}

// check verifies if the login request payload is ok
func (payload *UpdateUserCredentialRequestPayload) check() error {
	if len(payload.OldPassword) == 0 || len(payload.NewPassword) == 0 || len(payload.NewPasswordConfirmation) == 0 || len(payload.Email) == 0 {
		return errors.New("All fields must not be empty")
	}

	if payload.NewPassword != payload.NewPasswordConfirmation {
		return errors.New("Wrong password confirmation")
	}

	return verifyEmail(payload.Email)
}

func verifyEmail(email string) error {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if !re.MatchString(email) {
		return errors.New("Invalid email")
	}

	return nil
}
