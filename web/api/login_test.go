package api

import (
	"net/http"
)

// MockLoginAPI holds the default implementation of the User API interface
type MockLoginAPI struct {
}

func (mock *MockLoginAPI) LoginGETHandler(route string) http.Handler {
	return nil
}

func (mock *MockLoginAPI) LoginPOSTHandler() http.Handler {
	return nil
}
