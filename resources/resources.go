package resources

import (
	"github.com/labbsr0x/whisper/mail"
)

var Outbox chan<- mail.Mail
