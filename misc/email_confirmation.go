package misc

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labbsr0x/goh/gohtypes"
	"github.com/labbsr0x/whisper/mail"
	"net/http"
	"time"
)

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

func GetEmailConfirmationMail(username, email, challenge string) mail.Mail {
	to := []string{email}
	content := GetEmailConfirmationMailContent(username, challenge)

	return mail.Mail{To: to, Content: content}
}

func GetEmailConfirmationMailContent(username, challenge string) []byte {
	token := GetEmailConfirmationToken(username, challenge)
	link := GetEmailConfirmationLink(token)
	content := GetEmailConfirmationMessage(username, link)

	return []byte(content)
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

func GetEmailConfirmationLink(token string) string {
	return "localhost:7070/email-confirmation?email_confirmation_token=" + token
}

func GetEmailConfirmationMessage(username, redirect_to string) string {
	return fmt.Sprintf("Hi %v,\nClick on the link below to authenticate your email.\n\n %v\n\nThanks,\nWhisper Developers\n", username, redirect_to)
}
