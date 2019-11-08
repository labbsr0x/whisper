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
	PasswordMinChar       int
	PasswordMaxChar       int
	PasswordMinUniqueChar int
}

// SetHTML exposes the HTML from base page
func (p *RegistrationPage) SetHTML(html template.HTML) {
	p.HTML = html
}

// EmailConfirmationPage defines the information needed to load the email confirmation page
type EmailConfirmationPage struct {
	misc.BasePage
	Successful bool
	Message    string
	RedirectTo string
}

// SetHTML exposes the HTML from base page
func (p *EmailConfirmationPage) SetHTML(html template.HTML) {
	p.HTML = html
}

// ChangePasswordStep2Page defines the information needed to load the second step of change password page
type ChangePasswordStep2Page struct {
	misc.BasePage
	Username                    string
	Email                       string
	PasswordTooltip             string
	PasswordMinCharacters       int
	PasswordMaxCharacters       int
	PasswordMinUniqueCharacters int
}

// SetHTML exposes the HTML from base page
func (p *ChangePasswordStep2Page) SetHTML(html template.HTML) {
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
	PasswordMinChar       int
	PasswordMaxChar       int
	PasswordMinUniqueChar int
}

// SetHTML exposes the HTML from base page
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

// Check validates payload
func (payload *AddUserCredentialRequestPayload) Check() error {
	if len(payload.Username) == 0 || len(payload.Password) == 0 || len(payload.PasswordConfirmation) == 0 || len(payload.Email) == 0 {
		return fmt.Errorf("only challenge field can be empty")
	}

	if payload.Password != payload.PasswordConfirmation {
		return fmt.Errorf("wrong password confirmation")
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

// Check validates payload
func (payload *UpdateUserCredentialRequestPayload) Check() error {
	if len(payload.OldPassword) == 0 || len(payload.NewPassword) == 0 || len(payload.NewPasswordConfirmation) == 0 || len(payload.Email) == 0 {
		return fmt.Errorf("no field should be empty")
	}

	if payload.NewPassword != payload.NewPasswordConfirmation {
		return fmt.Errorf("wrong password confirmation")
	}

	return misc.VerifyEmail(payload.Email)
}

// ChangePasswordStep1UserCredentialRequestPayload defines the payload for start changing password
type ChangePasswordStep1UserCredentialRequestPayload struct {
	RedirectTo string `json:"redirect_to"`
	Email      string `json:"email"`
}

// Check validates payload
func (payload *ChangePasswordStep1UserCredentialRequestPayload) Check() error {
	if len(payload.Email) == 0 {
		return fmt.Errorf("email field should not be empty")
	}

	return misc.VerifyEmail(payload.Email)
}

// ChangePasswordStep2UserCredentialRequestPayload defines the payload for finish changing password
type ChangePasswordStep2UserCredentialRequestPayload struct {
	Token                   string `json:"token"`
	NewPassword             string `json:"newPassword"`
	NewPasswordConfirmation string `json:"newPasswordConfirmation"`
}

// Check validates payload
func (payload *ChangePasswordStep2UserCredentialRequestPayload) Check() error {
	if len(payload.Token) == 0 || len(payload.NewPassword) == 0 || len(payload.NewPasswordConfirmation) == 0 {
		return fmt.Errorf("no field should be empty")
	}

	if payload.NewPassword != payload.NewPasswordConfirmation {
		return fmt.Errorf("wrong password confirmation")
	}

	return nil
}
