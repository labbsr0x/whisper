package api

import "net/http"

type MockUserAPI struct {
}

// AddUserHandler REST POST api handler for adding new users
func (u *MockUserAPI) AddUserHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

// RemoveUserHandler REST POST api handler for removing users
func (u *MockUserAPI) RemoveUserHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

// ListUsersHandler REST GET api handler for listing available users.
func (u *MockUserAPI) ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

// GetUserHandler REST GET api handler for getting a user's info
func (u *MockUserAPI) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}
