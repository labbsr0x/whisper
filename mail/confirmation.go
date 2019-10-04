package mail

import (
	"github.com/labbsr0x/whisper/misc"
)

// Enum
const (
	emailConfirmationMail = "email_confirmation_mail.html"
)

type emailConfirmationMailContent struct {
	Link     string
	Username string
}

// GetEmailConfirmationMail render the mail for email confirmation
func GetEmailConfirmationMail(baseUIPath, secret, username, email, challenge string) Mail {
	to := []string{email}
	token := misc.GetEmailConfirmationToken(secret, username, challenge)
	link := "http://localhost:7070/email-confirmation?token=" + token
	page := emailConfirmationMailContent{Link: link, Username: username}
	content := render(baseUIPath, emailConfirmationMail, &page)

	return Mail{To: to, Content: content}
}