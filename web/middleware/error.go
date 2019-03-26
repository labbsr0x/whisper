package middleware

import (
	"net/http"

	"github.com/abilioesteves/goh/gohserver"
	"github.com/sirupsen/logrus"
)

// ErrorMiddleware deals with erros in a graceful way
func ErrorMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := new(gohserver.StatusWriter).Init(w)
		next.ServeHTTP(sw, r)
		gohserver.HandleError(sw)
		logrus.Infof("Done processing request to '%v' with status %v", r.RequestURI, sw.StatusCode) // logs
	})
}
