package types

import (
	"encoding/json"
	"github.com/labbsr0x/whisper/misc"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/labbsr0x/goh/gohtypes"
	"github.com/sirupsen/logrus"
)

// RegistrationPage defines the information needed to load a registration page
type RegistrationPage struct {
	Page
	LoginChallenge string
}

func (p *RegistrationPage) SetHTML(html template.HTML) {
	p.HTML = html
}

type EmailConfirmationPage struct {
	Page
	Successful bool
	Message    string
	RedirectTo string
}

func (p *EmailConfirmationPage) SetHTML(html template.HTML) {
	p.HTML = html
}

// UpdatePage defines the information needed to load a update user credentials page
type UpdatePage struct {
	Page
	Username   string
	Email      string
	Token      string
	RedirectTo string
}

func (p *UpdatePage) SetHTML(html template.HTML) {
	p.HTML = html
}

// AddUserCredentialRequestPayload defines the payload for adding a user
type AddUserCredentialRequestPayload struct {
	Email                string `json:"email"`
	Username             string `json:"username"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
	Challenge            string `json:"challenge"`
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
	data, err := ioutil.ReadAll(r.Body)
	gohtypes.PanicIfError("Not possible to parse registration payload", http.StatusBadRequest, err)

	err = json.Unmarshal(data, &payload)
	gohtypes.PanicIfError("Not possible to unmarshal update payload", http.StatusBadRequest, err)
	
	logrus.Debugf("Payload: '%v'", payload)

	payload.check()

	return payload
}

// check verifies if the login request payload is ok
func (payload *AddUserCredentialRequestPayload) check() {
	if len(payload.Username) == 0 || len(payload.Password) == 0 || len(payload.PasswordConfirmation) == 0 || len(payload.Email) == 0 {
		gohtypes.Panic("All fields are required", http.StatusBadRequest)
	}

	if payload.Password != payload.PasswordConfirmation {
		gohtypes.Panic("Wrong password confirmation", http.StatusBadRequest)
	}

	err := misc.ValidatePassword(payload.Password, payload.Username, payload.Email)
	gohtypes.PanicIfError("Password not valid", http.StatusBadRequest, err)

	verifyEmail(payload.Email)
}

// InitFromRequest initializes the update request payload from an http request form
func (payload *UpdateUserCredentialRequestPayload) InitFromRequest(r *http.Request) *UpdateUserCredentialRequestPayload {
	data, err := ioutil.ReadAll(r.Body)
	gohtypes.PanicIfError("Not possible to parse update payload", http.StatusBadRequest, err)

	err = json.Unmarshal(data, &payload)
	gohtypes.PanicIfError("Not possible to unmarshal update payload", http.StatusBadRequest, err)

	logrus.Debugf("Payload: '%v'", payload)

	payload.check()

	return payload
}

// check verifies if the login request payload is ok
func (payload *UpdateUserCredentialRequestPayload) check() {
	if len(payload.OldPassword) == 0 || len(payload.NewPassword) == 0 || len(payload.NewPasswordConfirmation) == 0 || len(payload.Email) == 0 {
		gohtypes.Panic("All fields must not be empty", http.StatusBadRequest)
	}

	if payload.NewPassword != payload.NewPasswordConfirmation {
		gohtypes.Panic("Wrong password confirmation", http.StatusBadRequest)
	}

	verifyEmail(payload.Email)
}

func verifyEmail(email string) {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if !re.MatchString(email) {
		gohtypes.Panic("Invalid email", http.StatusBadRequest)
	}
}
