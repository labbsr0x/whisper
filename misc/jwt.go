package misc

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labbsr0x/goh/gohtypes"
	"net/http"
	"net/url"
	"time"
)

// GenerateToken generates a jwt token
func GenerateToken(secret string, data jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)

	return token.SignedString([]byte(secret))
}

// ExtractClaimsTokenFromRequest extract a jwt token from a given request
func ExtractClaimsTokenFromRequest(secret string, r *http.Request) jwt.MapClaims {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	}

	tokenString, err := url.QueryUnescape(r.URL.Query().Get("token"))
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

// UnmarshalChangePasswordToken verify it is an change password token and extract the main confirmation
func UnmarshalChangePasswordToken(claims jwt.MapClaims) (username, challenge string) {
	cp, ok := claims["cp"].(bool)
	if !ok || !cp {
		gohtypes.Panic("Change password token not valid", http.StatusNotAcceptable)
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

func GetEmailConfirmationToken(secret, username, challenge string) string {
	claims := jwt.MapClaims{
		"sub":       username,                                // Subject
		"exp":       time.Now().Add(10 * time.Minute).Unix(), // Expiration
		"challenge": challenge,                               // Login Challenge
		"emt":       true,                                    // Email Confirmation Token
		"iat":       time.Now().Unix(),                       // Issued At
	}

	token, err := GenerateToken(secret, claims)
	gohtypes.PanicIfError("Not possible to create token", http.StatusInternalServerError, err)

	return token
}

func GetChangePasswordToken(secret, username, challenge string) string {
	claims := jwt.MapClaims{
		"sub":       username,                                // Subject
		"exp":       time.Now().Add(10 * time.Minute).Unix(), // Expiration
		"challenge": challenge,                               // Login Challenge
		"cp":       true,                                    // Change Password Token
		"iat":       time.Now().Unix(),                       // Issued At
	}

	token, err := GenerateToken(secret, claims)
	gohtypes.PanicIfError("Not possible to create token", http.StatusInternalServerError, err)

	return token
}