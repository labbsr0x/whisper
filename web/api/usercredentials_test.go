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

// GETRegistrationPageHandler builds the page where new credentials will be inserted
func (mock *MockUserCredentialsAPI) GETRegistrationPageHandler(route string) http.Handler {
	return nil
}

// GETUpdatePageHandler builder the page where credentials will be updated
func (mock *MockUserCredentialsAPI) GETUpdatePageHandler(route string) http.Handler {
	return nil
}
