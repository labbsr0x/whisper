package ui

import (
	"net/http"
)

// Handler defines the handler for ui requests
func Handler() http.Handler {
	return http.FileServer(http.Dir("/"))
}
