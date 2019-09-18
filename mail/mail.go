package mail

import (
	"net/smtp"
)

type smtpServer struct {
	host string
	port string
}

var (
	server = smtpServer{host: "smtp.gmail.com", port: "587"}
	from   = "alfredcoinworth@gmail.com"
	pass   = "tudosemprepioraantesdemelhorar"
)

func (s *smtpServer) address() string {
	return s.host + ":" + s.port
}

// Send use a smtp server to send an email
func Send(to []string, message []byte) error {
	auth := smtp.PlainAuth("", from, pass, server.host)

	return smtp.SendMail(server.address(), auth, from, to, message)
}
