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

func ParseToken(tokenString, secret string) (jwt.MapClaims, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	}

	token, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse the email confirmation token")
	}

	if !token.Valid {
		return nil, fmt.Errorf("Token is not valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("Unable to parse claims from the email confirmation token")
	}

	exp, ok := claims["exp"].(int64)
	if ok && exp > time.Now().Unix() {
		return nil, fmt.Errorf("Expired token")
	}

	return claims, nil
}

// ExtractClaimsTokenFromRequest extract a jwt token from a given request
func ExtractClaimsTokenFromRequest(secret string, r *http.Request) (jwt.MapClaims, error) {
	tokenString, err := url.QueryUnescape(r.URL.Query().Get("token"))
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve the email confirmation token")
	}

	return ParseToken(tokenString, secret)
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
func UnmarshalChangePasswordToken(claims jwt.MapClaims) (string, string, error) {
	cp, ok := claims["cp"].(bool)
	if !ok || !cp {
		return "", "", fmt.Errorf("Change password token not valid")
	}

	username, ok := claims["sub"].(string)
	if !ok {
		return "", "", fmt.Errorf("Unable to find the user")
	}

	redirectTo, ok := claims["redirect_to"].(string)
	if !ok {
		return "", "", fmt.Errorf("Unable to find the user")
	}

	return username, redirectTo, nil
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

func GetChangePasswordToken(secret, username, redirectTo string) string {
	if len(redirectTo) == 0 {
		redirectTo = "/login"
	}

	claims := jwt.MapClaims{
		"sub":         username,                                // Subject
		"redirect_to": redirectTo,                              // Redirect Back To
		"exp":         time.Now().Add(10 * time.Minute).Unix(), // Expiration
		"cp":          true,                                    // Change Password Token
		"iat":         time.Now().Unix(),                       // Issued At
	}

	token, err := GenerateToken(secret, claims)
	gohtypes.PanicIfError("Not possible to create token", http.StatusInternalServerError, err)

	return token
}
