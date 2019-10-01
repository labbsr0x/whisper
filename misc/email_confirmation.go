package misc

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labbsr0x/goh/gohtypes"
	"github.com/labbsr0x/whisper/web/ui"
	"net/http"
	"time"
)

// Enum
const (
	emailConfirmationMail = "email_confirmation_mail.html"
)

type emailConfirmationMailContent struct {
	Link     string
	Username string
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

// GetEmailConfirmationMail build the mail for email confirmation
func GetEmailConfirmationMail(username, email, challenge string) (to []string, content []byte) {
	to = []string{email}
	content = getEmailConfirmationMailContent(username, challenge)

	return
}

func getEmailConfirmationMailContent(username, challenge string) []byte {
	token := getEmailConfirmationToken(username, challenge)
	link := "http://localhost:7070/email-confirmation?email_confirmation_token=" + token

	return ui.BuildMail(emailConfirmationMail, emailConfirmationMailContent{Link: link, Username: username})
}

func getEmailConfirmationToken(username, challenge string) string {
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
