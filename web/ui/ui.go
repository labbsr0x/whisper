package ui

import (
	"bytes"
	"github.com/labbsr0x/goh/gohtypes"
	"github.com/labbsr0x/whisper/misc"
	"html/template"
	"io/ioutil"
	"net/http"
	"path"
)

// Enum
const (
	ChangePasswordStep1 = "change_password_step_1.html"
	ChangePasswordStep2 = "change_password_step_2.html"
	Consent             = "consent.html"
	EmailConfirmation   = "email_confirmation.html"
	Layout              = "index.html"
	Login               = "login.html"
	Registration        = "registration.html"
	Update              = "update.html"
)

// Handler defines the handler for ui requests
func Handler(baseUIPath string) http.Handler {
	uiPath := path.Join(baseUIPath, "")
	return http.FileServer(http.Dir(uiPath))
}

// Render render a page in the response
func WritePage(w http.ResponseWriter, baseUIPath, htmlFile string, page misc.IPage) {
	if page == nil {
		page = &misc.BasePage{}
	}

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
