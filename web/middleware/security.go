package middleware

import (
	"context"
	"net/http"

	"github.com/labbsr0x/whisper-client/hydra"

	"github.com/gorilla/mux"

	"github.com/labbsr0x/whisper/misc"
)

type key string

const (
	// TokenKey defines the key that shall be used to store a token in a requests' context
	TokenKey key = "token"
)

// GetSecurityMiddleware verifies if the client is authorized to make this request
func GetSecurityMiddleware(hydraClient *hydra.Client) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var tokenString string
			var token hydra.Token
			var err error

			if tokenString, err = misc.GetAccessTokenFromRequest(r); err == nil {
				if token, err = hydraClient.IntrospectToken(tokenString); err == nil {
					if token.Active {
						newR := r.WithContext(context.WithValue(r.Context(), TokenKey, token))
						next.ServeHTTP(w, newR)
						return
					}
				}
			}
			w.WriteHeader(401)
		})
	}

}
