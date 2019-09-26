package ui

import (
	"bytes"
	"github.com/labbsr0x/goh/gohtypes"
	"github.com/labbsr0x/whisper/web/api/types"
	"html/template"
	"io/ioutil"
	"net/http"
	"path"
)

// Handler defines the handler for ui requests
func Handler(baseUIPath string) http.Handler {
	uiPath := path.Join(baseUIPath, "")
	return http.FileServer(http.Dir(uiPath))
}

// LoadPage is used to load a page with the standard layout
func LoadPage(baseUIPath string, htmlFile string, page types.IPage, w http.ResponseWriter) {
	buf := new(bytes.Buffer)
	content := template.Must(template.ParseFiles(path.Join(baseUIPath, htmlFile)))

	err := content.Execute(buf, page)
	gohtypes.PanicIfError("Unable to load page", http.StatusInternalServerError, err)

	html, err := ioutil.ReadAll(buf)
	gohtypes.PanicIfError("Unable to read page", http.StatusInternalServerError, err)

	page.SetHTML(template.HTML(html))

	layout := template.Must(template.ParseFiles(path.Join(baseUIPath, "index.html")))
	err = layout.Execute(w, page)
	gohtypes.PanicIfError("Unable to load layout", http.StatusInternalServerError, err)
}
