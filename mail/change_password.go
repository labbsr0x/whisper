package mail

import (
	"fmt"
	"github.com/labbsr0x/whisper/misc"
)

// Enum
const (
	changePasswordMail = "change_password_mail.html"
)

type changePasswordMailContent struct {
	Link     string
	Username string
}

// GetChangePasswordMail render the mail for changing password
func GetChangePasswordMail(baseUIPath, secret, publicAddress, username, email, challenge string) Mail {
	to := []string{email}
	token := misc.GetEmailConfirmationToken(secret, username, challenge)
	link := fmt.Sprintf("%v/change-password?token=%v", publicAddress, token)
	page := changePasswordMailContent{Link: link, Username: username}
	content := render(baseUIPath, changePasswordMail, &page)

	return Mail{To: to, Content: content}
}
