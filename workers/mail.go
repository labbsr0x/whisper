package mail_worker

import (
	"github.com/labbsr0x/whisper/web/config"
	"net/smtp"
)

type MailWorkerAPI interface {
	Send(to []string, message []byte) error
}

type MailWorkerMessage struct {
	to []string
	content []byte
}

// MailWorker holds the default implementation of the Worker interface
type MailWorker struct {
	*config.WebBuilder
}

// InitFromWebBuilder initializes a default login api instance
func (w *MailWorker) InitFromWebBuilder(webBuilder *config.WebBuilder) *MailWorker {
	w.WebBuilder = webBuilder
	return w
}



//type smtpServer struct {
//	host string
//	port string
//}
//
//var (
//	server = smtpServer{host: "smtp.gmail.com", port: "587"}
//	from   = "alfredcoinworth@gmail.com"
//	pass   = "tudosemprepioraantesdemelhorar"
//)

func (w *MailWorker) InitWorker (ch chan MailWorkerMessage) {
	auth := smtp.PlainAuth("", w.MailUser, w.MailPassword, w.MailHost)
	address := w.MailHost + ":" + w.MailPort

	
}

// Send use a smtp server to send an email
func (w *MailWorker) Send(to []string, message []byte) error {


	return smtp.SendMail(address, auth, w.MailUser, to, message)
}
