package types

import (
	"fmt"
	"net/url"
)

// App defines the structure of an app
type App struct {
	ClientID          string `json:"id"`
	ClientName        string `json:"name"`
	ClientSecret      string `json:"secret"`
	Address           string `json:"address"`
	LoginRedirectURL  string `json:"loginRedirectURL"`
	LogoutRedirectURL string `json:"logoutRedirectURL"`
	Type              string `json:"type"`
}

var appTypes = map[string]string{
	"WEB": "WEB",
	"API": "API",
}

// AddAppInitialRequestPayload defines the payload for adding a new app
type AddAppInitialRequestPayload struct {
	Name              string `json:"name"`
	Type              string `json:"type"`
	LoginRedirectURL  string `json:"loginRedirectURL"`
	LogoutRedirectURL string `json:"logoutRedirectURL"`
}

// UpdateAppRequestPayload defines the payload for adding new scopes to an app
type UpdateAppRequestPayload struct {
	ClientID          string `json:"id"`
	Name              string `json:"name"`
	Address           string `json:"address"`
	Secret            string `json:"secret"`
	LoginRedirectURL  string `json:"loginRedirectURL"`
	LogoutRedirectURL string `json:"logoutRedirectURL"`
	Type              string `json:"type"`
}

// Check validates add app initial request payload
func (payload *AddAppInitialRequestPayload) Check() error {
	if len(payload.Name) == 0 || len(payload.Type) == 0 {
		return fmt.Errorf("Name and Type cannot be empty")
	}

	if appTypes[payload.Type] == "" {
		return fmt.Errorf("Application type '%v' not supported", appTypes[payload.Type])
	}

	return checkLoginAndLogoutURLs(payload.LoginRedirectURL, payload.LogoutRedirectURL, payload.Type)
}

// Check validates update app request payload
func (payload *UpdateAppRequestPayload) Check() error {
	return checkLoginAndLogoutURLs(payload.LoginRedirectURL, payload.LogoutRedirectURL, payload.Type)
}

func checkLoginAndLogoutURLs(loginURL, logoutURL, appType string) error {
	if appType == "WEB" {
		if len(loginURL) > 0 && len(logoutURL) > 0 {
			li, err := url.ParseRequestURI(loginURL)
			if err != nil {
				return fmt.Errorf("LoginRedirectURL '%v' not valid", loginURL)
			}
			lo, err := url.ParseRequestURI(logoutURL)
			if err != nil {
				return fmt.Errorf("LogoutRedirectURL '%v' not valid", logoutURL)
			}

			if li.Hostname() != lo.Hostname() {
				return fmt.Errorf("Login and Logout URL must point to the same Host")
			}
		} else {
			return fmt.Errorf("Login and Logout redirect URLs mandatory for a Web Application")
		}
	}
	return nil
}
