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

// UpdateUserHandler REST PUT api handler for updating a user's info
func (u *MockUserAPI) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}
