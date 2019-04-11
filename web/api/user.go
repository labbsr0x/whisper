package api

import (
	"net/http"

	"github.com/abilioesteves/goh/gohserver"
)

// UserAPI defines the available user apis
type UserAPI interface {
	AddUserHandler(w http.ResponseWriter, r *http.Request)
	RemoveUserHandler(w http.ResponseWriter, r *http.Request)
	UpdateUserHandler(w http.ResponseWriter, r *http.Request)
}

// DefaultUserAPI holds the default implementation of the User API interface
type DefaultUserAPI struct {
}

// AddUserHandler REST POST api handler for adding new users
func (u *DefaultUserAPI) AddUserHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
	gohserver.WriteJSONResponse("AddUser: This is just a test", 200, w)
}

// RemoveUserHandler REST POST api handler for removing users
func (u *DefaultUserAPI) RemoveUserHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
	gohserver.WriteJSONResponse("RemoveUser: This is just a test", 200, w)
}

// UpdateUserHandler REST PUT api handler for updating users
func (u *DefaultUserAPI) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
	gohserver.WriteJSONResponse("UpdateUser: This is just a test", 200, w)
}
