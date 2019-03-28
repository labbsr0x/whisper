package ui

import (
	"net/http"
	"path"
)

// Handler defines the handler for ui requests
func Handler(baseUIPath string) http.Handler {
	uiPath := path.Join(baseUIPath, "")
	return http.FileServer(http.Dir(uiPath))
}
