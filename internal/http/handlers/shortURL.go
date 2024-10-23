package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/WeisseNacht18/url-shortener/internal/storage"
)

func CreateShortURLHandler(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.Header.Get("Content-Type"), "text/plain") ||
		strings.Contains(r.Header.Get("Content-Type"), "application/x-gzip") ||
		strings.Contains(r.Header.Get("Content-Type"), "application/gzip") {

		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		link := string(body)
		shortLink, hasURL := storage.AddURLToStorage(link)
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))

		if hasURL {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusCreated)
		}

		content := []byte(BaseURL + "/" + shortLink)
		w.Write(content)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
