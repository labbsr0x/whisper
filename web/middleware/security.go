package middleware

import (
	"net/http"
)

// SecurityMiddleware verifies if the client is authorized to make this request
func SecurityMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO validate token
		w.WriteHeader(403)
	})

}
