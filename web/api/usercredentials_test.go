package api

import "net/http"

type MockUserCredentialsAPI struct {
}

// POSTHandler handles post requests to create user credentials
func (mock *MockUserCredentialsAPI) POSTHandler() http.Handler {
	return nil
}

// PUTHandler handles put requests to update user credentials
func (mock *MockUserCredentialsAPI) PUTHandler() http.Handler {
	return nil
}

// GETPageHandler builds the page where new credentials will be inserted
func (mock *MockUserCredentialsAPI) GETPageHandler(route string) http.Handler {
	return nil
}
