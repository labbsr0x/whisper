package mail

import (
	"fmt"
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
func GetEmailConfirmationMail(baseUIPath, secret, publicAddress, username, email, challenge string) Mail {
	to := []string{email}
	token := misc.GetEmailConfirmationToken(secret, username, challenge)
	link := fmt.Sprintf("%v/email-confirmation?token=%v", publicAddress, token)
	page := emailConfirmationMailContent{Link: link, Username: username}
	content := render(baseUIPath, emailConfirmationMail, &page)

	return Mail{To: to, Content: content}
}