package api

import (
	"net/http"
)

// UserAPI defines the available user apis
type UserAPI interface {
	AddUserHandler(w http.ResponseWriter, r *http.Request)
	RemoveUserHandler(w http.ResponseWriter, r *http.Request)
	ListUsersHandler(w http.ResponseWriter, r *http.Request)
	GetUserHandler(w http.ResponseWriter, r *http.Request)
}

// DefaultUserAPI holds the default implementation of the User API interface
type DefaultUserAPI struct {
}

// AddUserHandler REST POST api handler for adding new users
func (u *DefaultUserAPI) AddUserHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

// RemoveUserHandler REST POST api handler for removing users
func (u *DefaultUserAPI) RemoveUserHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

// ListUsersHandler REST GET api handler for listing available users.
func (u *DefaultUserAPI) ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

// GetUserHandler REST GET api handler for getting a user's info
func (u *DefaultUserAPI) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}
