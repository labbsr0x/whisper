package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/labbsr0x/goh/gohtypes"

	"github.com/labbsr0x/whisper-client/hydra"

	"github.com/gorilla/mux"
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

			if tokenString, err = getAccessTokenFromRequest(r); err == nil {
				if token, err = hydraClient.IntrospectToken(tokenString); err == nil {
					if token.Active {
						newR := r.WithContext(context.WithValue(r.Context(), TokenKey, token))
						next.ServeHTTP(w, newR)
						return
					}
				}
			}
			gohtypes.PanicIfError("Unauthorized user", 401, err)
		})
	}
}

// getAccessTokenFromRequest is a helper method to recover an Access Token from a http request
func getAccessTokenFromRequest(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	authURLParam := r.URL.Query().Get("token")
	var t string

	if len(authHeader) == 0 && len(authURLParam) == 0 {
		return "", fmt.Errorf("No Authorization Header or URL Param found")
	}

	if len(authHeader) > 0 {
		data := strings.Split(authHeader, " ")

		if len(data) != 2 {
			return "", fmt.Errorf("Bad Authorization Header")
		}

		t = data[0]

		if len(t) == 0 || t != "Bearer" {
			return "", fmt.Errorf("No Bearer Token found")
		}

		t = data[1]

	} else {
		t = authURLParam
	}

	if len(t) == 0 {
		return "", fmt.Errorf("Bad Authorization Token")
	}

	return t, nil
}
