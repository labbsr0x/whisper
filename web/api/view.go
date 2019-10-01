package api

import (
	"github.com/labbsr0x/goh/gohtypes"
	"net/http"
)

// Render render a page in the response
func Render(w http.ResponseWriter, view []byte) {
	_, err := w.Write(view)
	gohtypes.PanicIfError("Unable to render", http.StatusInternalServerError, err)
}
