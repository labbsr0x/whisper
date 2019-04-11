package api

import "net/http"

type MockUserCredentialsAPI struct {
}

// AddUserCredentialHandler REST POST api handler for adding new users
func (u *MockUserCredentialsAPI) AddUserCredentialHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

// RemoveUserCredentialHandler REST POST api handler for removing users
func (u *MockUserCredentialsAPI) RemoveUserCredentialHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

// UpdateUserCredentialHandler REST PUT api handler for updating a user's info
func (u *MockUserCredentialsAPI) UpdateUserCredentialHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}
