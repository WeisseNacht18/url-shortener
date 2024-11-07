package handlers

import (
	"net/http"

	"github.com/WeisseNacht18/url-shortener/internal/logger"
	"github.com/WeisseNacht18/url-shortener/internal/storage"
)

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("x-user-id")
	shortLink := r.RequestURI[1:]
	value, ok := storage.GetURLFromStorage(userID, shortLink)

	logger.Logger.Infoln(userID, "|", shortLink, "|", value, "|", ok)

	if ok {
		http.Redirect(w, r, value, http.StatusTemporaryRedirect)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
