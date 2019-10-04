package ui

import (
	"bytes"
	"github.com/labbsr0x/goh/gohtypes"
	"html/template"
	"io/ioutil"
	"net/http"
	"path"
)

// Enum
const (
	Consent           = "consent.html"
	EmailConfirmation = "email_confirmation.html"
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

// Render render a page in the response
func WritePage(w http.ResponseWriter, baseUIPath, htmlFile string, page IPage) {
	buf := new(bytes.Buffer)
	content := template.Must(template.ParseFiles(path.Join(baseUIPath, htmlFile)))

	err := content.Execute(buf, page)
	gohtypes.PanicIfError("Unable to load page", http.StatusInternalServerError, err)

	html, err := ioutil.ReadAll(buf)
	gohtypes.PanicIfError("Unable to read page", http.StatusInternalServerError, err)

	page.SetHTML(template.HTML(html))

	layout := template.Must(template.ParseFiles(path.Join(baseUIPath, Layout)))
	err = layout.Execute(buf, page)
	gohtypes.PanicIfError("Unable to load layout", http.StatusInternalServerError, err)

	_, err = w.Write(buf.Bytes())
	gohtypes.PanicIfError("Unable to render", http.StatusInternalServerError, err)
}