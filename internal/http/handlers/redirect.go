package handlers

import (
	"net/http"
	"strings"

	"github.com/WeisseNacht18/url-shortener/internal/storage"
)

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/")
	value, ok := storage.GetURLFromStorage(id)
	if ok {
		http.Redirect(w, r, value, http.StatusTemporaryRedirect)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
