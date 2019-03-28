package api

import (
	"net/http"
)

// MockLoginAPI holds the default implementation of the User API interface
type MockLoginAPI struct {
}

// LoginPOSTHandler REST POST api handler for logging in users
func (u *MockLoginAPI) LoginPOSTHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
