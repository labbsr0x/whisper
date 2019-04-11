package middleware

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/abilioesteves/whisper/misc"
)

// SecurityMiddleware verifies if the client is authorized to make this request
func SecurityMiddleware(next http.Handler, hydraClient *misc.HydraClient) mux.MiddlewareFunc {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var tokenString string
			var token misc.HydraToken
			var err error

			if tokenString, err = misc.GetAccessTokenFromRequest(r); err == nil {
				if token, err = hydraClient.IntrospectToken(tokenString); err == nil {
					if token.Active {
						newR := r.WithContext(context.WithValue(r.Context(), "token", token))
						next.ServeHTTP(w, newR)
						return
					}
				}
			}
			w.WriteHeader(401)
		})
	}

}
