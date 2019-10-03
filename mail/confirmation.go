package mail

import (
	"github.com/labbsr0x/whisper/misc"
	"github.com/labbsr0x/whisper/web/ui"
)

// Enum
const (
	emailConfirmationMail = "email_confirmation_mail.html"
)

type emailConfirmationMailContent struct {
	Link     string
	Username string
}

// GetEmailConfirmationMail build the mail for email confirmation
func GetEmailConfirmationMail(baseUIPath, username, email, challenge string) Mail {
	to := []string{email}
	token := misc.GetEmailConfirmationToken(username, challenge)
	link := "http://localhost:7070/email-confirmation?email_confirmation_token=" + token
	content := ui.BuildMail(baseUIPath, emailConfirmationMail, emailConfirmationMailContent{Link: link, Username: username})

	return Mail{To: to, Content: content}
}
