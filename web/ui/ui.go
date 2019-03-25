package ui

import (
	"net/http"
	"path"

	"github.com/sirupsen/logrus"
)

// Handler defines the handler for ui requests
func Handler(baseUIPath string) http.Handler {
	uiPath := path.Join(baseUIPath, "/static")
	logrus.Infof("UI Path: %v", uiPath)
	return http.FileServer(http.Dir(uiPath))
}
