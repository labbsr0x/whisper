package mail

import (
	"bytes"
	"encoding/base64"
	"github.com/google/uuid"
	"github.com/labbsr0x/goh/gohtypes"
	"github.com/sirupsen/logrus"
	"html/template"
	"io/ioutil"
	"mime/quotedprintable"
	"net/http"
	"net/smtp"
	"os"
	"path"
)

// API defines what should the mail expose
type API interface {
	Init(user, password, host, port string, inbox <-chan Mail) API
	Run()
}

// Mail holds the content necessary for a email
type Mail struct {
	To      []string
	Content []byte
}

var _sendMail = smtp.SendMail

// DefaultHandler holds the default implementation of mail
type DefaultHandler struct {
	user    string
	address string
	auth    smtp.Auth
	Inbox   <-chan Mail
}

// Init initializes a default email api instance
func (mh *DefaultHandler) Init(user, password, host, port string, inbox <-chan Mail) API {
	mh.user = user
	mh.address = host + ":" + port
	mh.auth = smtp.PlainAuth("", user, password, host)
	mh.Inbox = inbox

	return mh
}

// Run start the goroutine that listen for emails to send
func (mh *DefaultHandler) Run() {
	go func() {
		for mail := range mh.Inbox {
			err := _sendMail(mh.address, mh.auth, mh.user, mail.To, mail.Content)

			if err != nil {
				logrus.Error(err)
			}
		}
	}()
}

func render(baseUIPath, htmlFile string, mailContent interface{}) []byte {
	var body bytes.Buffer

	boundary := uuid.New().String()

	// Add headers
	body.WriteString("Subject: Whisper\n")
	body.WriteString("MIME-version: 1.0;\n")
	body.WriteString("Content-Type: multipart/mixed; boundary=" + boundary + "\n") // Tells the content has multiple parts, each initiated with the '-- boundary'

	// Add HTML to body
	body.WriteString("\n\n--" + boundary + "\n")
	body.WriteString("Content-Type: text/html; charset=\"UTF-8\";\n")
	body.WriteString("Content-Transfer-Encoding: quoted-printable\n\n")
	body.WriteString(getHTMLBytes(baseUIPath, htmlFile, mailContent))

	// Embed logo
	logoFile := "spy-black.png"
	logoName := "logo"
	body.WriteString("\n\n--" + boundary + "\n")
	body.WriteString("Content-Type: image/png; name=\"" + logoFile + "\"\n")
	body.WriteString("Content-Transfer-Encoding: base64\n")
	body.WriteString("Content-ID: <" + logoName + ">\n")
	body.WriteString("Content-Disposition: inline; filename=\"" + logoFile + "\"\n")
	body.WriteString("X-Attachment-Id: " + logoName + "\n\n")
	body.Write(getLogoBytes(baseUIPath, logoFile))

	// End multiple parts
	body.WriteString("\n\n--" + boundary + "--\n")

	return body.Bytes()
}

func getHTMLBytes(baseUIPath, htmlFile string, mailContent interface{}) string {
	tmpl, err := template.ParseFiles(path.Join(baseUIPath, htmlFile))
	gohtypes.PanicIfError("Unable to open mail content", http.StatusInternalServerError, err)

	buff := new(bytes.Buffer)
	err = tmpl.Execute(buff, mailContent)
	gohtypes.PanicIfError("Unable to load mail content", http.StatusInternalServerError, err)

	res := new(bytes.Buffer)
	w := quotedprintable.NewWriter(res)

	_, err = w.Write(buff.Bytes())
	gohtypes.PanicIfError("Unable to quote print mail content", http.StatusInternalServerError, err)

	return string(res.Bytes())
}

func getLogoBytes(baseUIPath, logoName string) []byte {
	logo := baseUIPath + "/static/images/" + logoName
	file, err := os.Open(logo)
	gohtypes.PanicIfError("Unable to open email images", http.StatusInternalServerError, err)

	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	gohtypes.PanicIfError("Unable to load email images", http.StatusInternalServerError, err)

	buff := make([]byte, base64.StdEncoding.EncodedLen(len(fileBytes)))
	base64.StdEncoding.Encode(buff, fileBytes)

	return buff
}
