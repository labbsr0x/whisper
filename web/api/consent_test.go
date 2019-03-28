package api

import (
	"net/http"
)

// MockConsentAPI holds the default implementation of the User API interface
type MockConsentAPI struct {
}

// ConsentPOSTHandler REST POST api handler for app authorization
func (u *MockConsentAPI) ConsentPOSTHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
	w.WriteHeader(200)
}
