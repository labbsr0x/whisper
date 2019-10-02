package misc

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labbsr0x/goh/gohtypes"
	"github.com/spf13/viper"
	"net/http"
	"net/url"
	"time"
)

// GenerateToken generates a jwt token
func GenerateToken(data jwt.MapClaims) (string, error) {
	secret := viper.GetString("secret-key")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)

	return token.SignedString([]byte(secret))
}

// ExtractClaimsTokenFromRequest extract a jwt token from a given request
func ExtractClaimsTokenFromRequest(r *http.Request) jwt.MapClaims {
	secret := viper.GetString("secret-key")
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	}

	tokenString, err := url.QueryUnescape(r.URL.Query().Get("email_confirmation_token"))
	gohtypes.PanicIfError("Unable to retrieve the email confirmation token", http.StatusBadRequest, err)

	token, err := jwt.Parse(tokenString, keyFunc)
	gohtypes.PanicIfError("Unable to parse the email confirmation token", http.StatusBadRequest, err)

	if !token.Valid {
		gohtypes.Panic("Token is not valid", http.StatusBadRequest)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		gohtypes.Panic("Unable to parse claims from the email confirmation token", http.StatusBadRequest)
	}

	exp, ok := claims["exp"].(int64)
	if ok && exp > time.Now().Unix() {
		gohtypes.Panic("Expired token", http.StatusBadRequest)
	}

	return claims
}

// UnmarshalEmailConfirmationToken verify it is an email confirmation token and extract the main confirmation
func UnmarshalEmailConfirmationToken(claims jwt.MapClaims) (username, challenge string) {
	emt, ok := claims["emt"].(bool)
	if !ok || !emt {
		gohtypes.Panic("Email confirmation token not valid", http.StatusNotAcceptable)
	}

	username, ok = claims["sub"].(string)
	if !ok {
		gohtypes.Panic("Unable to find the user", http.StatusNotFound)
	}

	challenge, ok = claims["challenge"].(string)
	if !ok {
		gohtypes.Panic("Unable to find the login challenge", http.StatusNotFound)
	}

	return
}

func GetEmailConfirmationToken(username, challenge string) string {
	claims := jwt.MapClaims{
		"sub":       username,                                // Subject
		"exp":       time.Now().Add(10 * time.Minute).Unix(), // Expiration
		"challenge": challenge,                               // Login Challenge
		"emt":       true,                                    // Email Confirmation Token
		"iat":       time.Now().Unix(),                       // Issued At
	}

	token, err := GenerateToken(claims)
	gohtypes.PanicIfError("Not possible to create token", http.StatusInternalServerError, err)

	return token
}
