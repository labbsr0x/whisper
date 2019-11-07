package types

import (
	"fmt"
	"github.com/labbsr0x/whisper/misc"
	"html/template"
)

// RegistrationPage defines the information needed to load a registration page
type RegistrationPage struct {
	misc.BasePage
	LoginChallenge              string
	PasswordTooltip             string
	PasswordMinCharacters       int
	PasswordMaxCharacters       int
	PasswordMinUniqueCharacters int
}

func (p *RegistrationPage) SetHTML(html template.HTML) {
	p.HTML = html
}

type EmailConfirmationPage struct {
	misc.BasePage
	Successful bool
	Message    string
	RedirectTo string
}

func (p *EmailConfirmationPage) SetHTML(html template.HTML) {
	p.HTML = html
}

type ChangePasswordPage struct {
	misc.BasePage
	Username                    string
	Email                       string
	PasswordTooltip             string
	PasswordMinCharacters       int
	PasswordMaxCharacters       int
	PasswordMinUniqueCharacters int
}

func (p *ChangePasswordPage) SetHTML(html template.HTML) {
	p.HTML = html
}

// UpdatePage defines the information needed to load a update user credentials page
type UpdatePage struct {
	misc.BasePage
	Username                    string
	Email                       string
	Token                       string
	RedirectTo                  string
	PasswordTooltip             string
	PasswordMinCharacters       int
	PasswordMaxCharacters       int
	PasswordMinUniqueCharacters int
}

func (p *UpdatePage) SetHTML(html template.HTML) {
	p.HTML = html
}

// AddUserCredentialResponsePayload defines the response payload after adding a user
type AddUserCredentialResponsePayload struct {
	UserCredentialID string
}

// AddUserCredentialRequestPayload defines the payload for adding a user
type AddUserCredentialRequestPayload struct {
	Email                string `json:"email"`
	Username             string `json:"username"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
	Challenge            string `json:"challenge"`
}

func (payload *AddUserCredentialRequestPayload) Check() error {
	if len(payload.Username) == 0 || len(payload.Password) == 0 || len(payload.PasswordConfirmation) == 0 || len(payload.Email) == 0 {
		return fmt.Errorf("All fields are required")
	}

	if payload.Password != payload.PasswordConfirmation {
		return fmt.Errorf("Wrong password confirmation")
	}

	err := misc.ValidatePassword(payload.Password, payload.Username, payload.Email)
	if err != nil {
		return err
	}

	return misc.VerifyEmail(payload.Email)
}

// UpdateUserCredentialRequestPayload defines the payload for updating a user
type UpdateUserCredentialRequestPayload struct {
	Email                   string `json:"email"`
	NewPassword             string `json:"newPassword"`
	NewPasswordConfirmation string `json:"newPasswordConfirmation"`
	OldPassword             string `json:"oldPassword"`
}

func (payload *UpdateUserCredentialRequestPayload) Check() error {
	if len(payload.OldPassword) == 0 || len(payload.NewPassword) == 0 || len(payload.NewPasswordConfirmation) == 0 || len(payload.Email) == 0 {
		return fmt.Errorf("All fields must not be empty")
	}

	if payload.NewPassword != payload.NewPasswordConfirmation {
		return fmt.Errorf("Wrong password confirmation")
	}

	return misc.VerifyEmail(payload.Email)
}

type ChangePasswordInitUserCredentialRequestPayload struct {
	RedirectTo string `json:"redirect_to"`
	Email      string `json:"email"`
}

func (payload *ChangePasswordInitUserCredentialRequestPayload) Check() error {
	if len(payload.Email) == 0 {
		return fmt.Errorf("All fields must not be empty")
	}

	return misc.VerifyEmail(payload.Email)
}

type ChangePasswordUserCredentialRequestPayload struct {
	Token                   string `json:"token"`
	NewPassword             string `json:"newPassword"`
	NewPasswordConfirmation string `json:"newPasswordConfirmation"`
}

func (payload *ChangePasswordUserCredentialRequestPayload) Check() error {
	if len(payload.Token) == 0 || len(payload.NewPassword) == 0 || len(payload.NewPasswordConfirmation) == 0 {
		return fmt.Errorf("All fields must not be empty")
	}

	if payload.NewPassword != payload.NewPasswordConfirmation {
		return fmt.Errorf("Wrong password confirmation")
	}

	return nil
}
