package workers

import (
	"github.com/labbsr0x/whisper/web/config"
	"github.com/sirupsen/logrus"
	"net/smtp"
)

// MailWorkerAPI defines what a mail worker should expose
type MailWorkerAPI interface {
	Start()
	Send (to []string, content []byte)
	Stop()
}

// MailWork defines the email
type MailWork struct {
	to      []string
	content []byte
}

// DefaultMailWorkerAPI holds the default implementation of the Mail Worker interface
type DefaultMailWorkerAPI struct {
	*config.WebBuilder
	JobChannel chan MailWork
	EndChannel chan bool
}

var MailWorker DefaultMailWorkerAPI

// InitFromWebBuilder initializes a default login api instance
func (w *DefaultMailWorkerAPI) InitFromWebBuilder(webBuilder *config.WebBuilder) *DefaultMailWorkerAPI {
	w.WebBuilder = webBuilder
	return w
}

// Start inits a goroutine that awaits for MailWorkerMessages
func (w *DefaultMailWorkerAPI) Start() {
	logrus.Info("Worker is starting")

	auth := smtp.PlainAuth(w.MailIdentity, w.MailUser, w.MailPassword, w.MailHost)
	address := w.MailHost + ":" + w.MailPort

	go serve(w.MailUser, address, auth, w.JobChannel, w.EndChannel)
}

// Send pass some work to the worker
func (w *DefaultMailWorkerAPI) Send (to []string, content []byte) {
	logrus.Info("Worker is being required to perform some work")
	w.JobChannel <- MailWork{to: to, content: content}
}

// Stop interrupts the worker
func (w *DefaultMailWorkerAPI) Stop() {
	logrus.Info("Worker is stopping")
	w.EndChannel <- true
}

func serve(user string, address string, auth smtp.Auth, ch chan MailWork, cancel chan bool) {
	for {
		select {
		case job := <-ch:
			err := smtp.SendMail(address, auth, user, job.to, job.content)

			if err != nil {
				logrus.Error(err)
			}
		case <-cancel:
			return
		}
	}
}

// valores default
// mail-user "alfredcoinworth@gmail.com"
// mail-password "tudosemprepioraantesdemelhorar"
// mail-host "smtp.gmail.com"
// mail-port "587"
// mail-identiy ""
