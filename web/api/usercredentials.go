package api

import (
	"net/http"

	"github.com/abilioesteves/goh/gohserver"
)

// UserCredentialsAPI defines the available user apis
type UserCredentialsAPI interface {
	AddUserCredentialHandler(w http.ResponseWriter, r *http.Request)
	RemoveUserCredentialHandler(w http.ResponseWriter, r *http.Request)
	UpdateUserCredentialHandler(w http.ResponseWriter, r *http.Request)
}

// DefaultUserCredentialAPI holds the default implementation of the User API interface
type DefaultUserCredentialsAPI struct {
}

// AddUserCredentialHandler REST POST api handler for adding new users
func (u *DefaultUserCredentialsAPI) AddUserCredentialHandler(w http.ResponseWriter, r *http.Request) {
	// TODO

	gohserver.WriteJSONResponse("AddUserCredentialHandler: This is just a test", 200, w)
}

// RemoveUserCredentialHandler REST DELETE api handler for removing users
func (u *DefaultUserCredentialsAPI) RemoveUserCredentialHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
	gohserver.WriteJSONResponse("RemoveUserCredentialHandler: This is just a test", 200, w)
}

// UpdateUserCredentialHandler REST PUT api handler for updating users
func (u *DefaultUserCredentialsAPI) UpdateUserCredentialHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
	gohserver.WriteJSONResponse("UpdateUserCredentialHandler: This is just a test", 200, w)
}
