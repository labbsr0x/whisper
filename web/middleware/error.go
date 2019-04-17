package middleware

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/labbsr0x/goh/gohserver"
	"github.com/sirupsen/logrus"
)

// GetErrorMiddleware deals with erros in a graceful way
func GetErrorMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sw := new(gohserver.StatusWriter).Init(w)
			func() {
				defer gohserver.HandleError(sw)
				next.ServeHTTP(sw, r)
			}()
			logrus.Debugf("Done processing request to '%v' with status %v", r.RequestURI, sw.StatusCode) // logs
		})
	}
}
