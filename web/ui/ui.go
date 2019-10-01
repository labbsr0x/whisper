package ui

import (
	"bytes"
	"encoding/base64"
	"github.com/labbsr0x/goh/gohtypes"
	"github.com/labbsr0x/whisper/resources"
	"html/template"
	"io/ioutil"
	"mime/quotedprintable"
	"net/http"
	"os"
	"path"
)

// Enum
const (
	Consent           = "consent.html"
	EmailConfirmation = "email_confirmation_mail.html"
	Layout            = "index.html"
	Login             = "login.html"
	Registration      = "registration.html"
	Update            = "update.html"
)

type IPage interface {
	SetHTML(html template.HTML)
}

// Handler defines the handler for ui requests
func Handler(baseUIPath string) http.Handler {
	uiPath := path.Join(baseUIPath, "")
	return http.FileServer(http.Dir(uiPath))
}

// BuildPage is used to load a page with the standard layout
func BuildPage(htmlFile string, page IPage) []byte {
	buf := new(bytes.Buffer)
	content := template.Must(template.ParseFiles(path.Join(resources.BaseUIPath, htmlFile)))

	err := content.Execute(buf, page)
	gohtypes.PanicIfError("Unable to load page", http.StatusInternalServerError, err)

	html, err := ioutil.ReadAll(buf)
	gohtypes.PanicIfError("Unable to read page", http.StatusInternalServerError, err)

	page.SetHTML(template.HTML(html))

	layout := template.Must(template.ParseFiles(path.Join(resources.BaseUIPath, Layout)))
	err = layout.Execute(buf, page)
	gohtypes.PanicIfError("Unable to load layout", http.StatusInternalServerError, err)

	return buf.Bytes()
}

func BuildMail(htmlFile string, mailContent interface{}) []byte {
	var body bytes.Buffer

	boundary := "f46d043c813270fc6b04c2d223da"

	// Add headers
	body.WriteString("Subject: Whisper\n")
	body.WriteString("MIME-version: 1.0;\n")
	body.WriteString("Content-Type: multipart/mixed; boundary=" + boundary + "\n") // Tells the content has multiple parts, each initiated with the '-- boundary'

	// Add HTML to body
	body.WriteString("\n\n--" + boundary + "\n")
	body.WriteString("Content-Type: text/html; charset=\"UTF-8\";\n")
	body.WriteString("Content-Transfer-Encoding: quoted-printable\n\n")
	body.WriteString(getHTMLBytes(htmlFile, mailContent))

	// Embed logo
	logoName := "spyblack"
	body.WriteString("\n\n--" + boundary + "\n")
	body.WriteString("Content-Type: image/png; name=\"" + logoName + ".png\"\n")
	body.WriteString("Content-Transfer-Encoding: base64\n")
	body.WriteString("Content-ID: <" + logoName + ">\n")
	body.WriteString("Content-Disposition: inline; filename=\"" + logoName + ".png\"\n")
	body.WriteString("X-Attachment-Id: " + logoName + "\n\n")
	body.Write(getLogoBytes(logoName + ".png"))

	// End multiple parts
	body.WriteString("\n\n--" + boundary + "--\n")

	return body.Bytes()
}

func getHTMLBytes(htmlFile string, mailContent interface{}) string {
	tmpl, err := template.ParseFiles(path.Join(resources.BaseUIPath, htmlFile))
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

func getLogoBytes(logoName string) []byte {
	logo := resources.BaseUIPath + "/static/images/" + logoName
	file, err := os.Open(logo)
	gohtypes.PanicIfError("Unable to open email images", http.StatusInternalServerError, err)

	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	gohtypes.PanicIfError("Unable to load email images", http.StatusInternalServerError, err)

	buff := make([]byte, base64.StdEncoding.EncodedLen(len(fileBytes)))
	base64.StdEncoding.Encode(buff, fileBytes)

	return buff
}
