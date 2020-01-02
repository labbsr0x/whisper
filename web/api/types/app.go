package types

import (
	"fmt"
	"github.com/labbsr0x/whisper/misc"
)

var grantTypes = map[string]bool{"authorization_code": true, "refresh_token": true, "client_credentials": true}

// AddAppRequestPayload defines the payload for adding a new app
type AddAppRequestPayload struct {
	ID                string   `json:"id"`
	Secret            string   `json:"secret"`
	Name              string   `json:"name"`
	URL               string   `json:"url"`
	LoginRedirectURL  string   `json:"loginRedirectURL"`
	LogoutRedirectURL string   `json:"logoutRedirectURL"`
	GrantTypes        []string `json:"grantTypes"`
}

// AddScopesRequestPayload defines the payload for adding new scopes to an app
type AddScopesRequestPayload struct {
	*misc.GrantScope
}

// Check validates payload
func (payload *AddAppRequestPayload) Check() error {
	if len(payload.ID) == 0 || len(payload.Secret) == 0 || len(payload.Name) == 0 || len(payload.URL) == 0 {
		return fmt.Errorf("only challenge field can be empty")
	}

	if len(payload.GrantTypes) > 0 {
		for _, gt := range payload.GrantTypes {
			if gt == "authorization_code" && len(payload.LoginRedirectURL) == 0 || len(payload.LogoutRedirectURL) == 0 {
				return fmt.Errorf("login and logout redirect URLs must be non empty when an authorization_code grant is requested")
			}
			if !grantTypes[gt] {
				return fmt.Errorf("the grant type '%v' is not supported", gt)
			}
		}
	} else {
		return fmt.Errorf("at least one grant type must be selected")
	}

	return nil
}

// App defines the structure of an app
type App struct {
	ID                string           `json:"id"`
	Name              string           `json:"name"`
	URL               string           `json:"url"`
	LoginRedirectURL  string           `json:"loginRedirectURL"`
	LogoutRedirectURL string           `json:"logoutRedirectURL"`
	GrantTypes        []string         `json:"grantTypes"`
	Scopes            misc.GrantScopes `json:"scopes"`
}
