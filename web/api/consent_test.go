package api

import (
	"net/http"
)

// MockConsentAPI holds the default implementation of the User API interface
type MockConsentAPI struct {
}

func (mock *MockConsentAPI) ConsentGETHandler(route string) http.Handler {
	return nil
}

func (mock *MockConsentAPI) ConsentPOSTHandler() http.Handler {
	return nil
}
