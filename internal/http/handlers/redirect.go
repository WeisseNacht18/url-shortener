package handlers

import (
	"net/http"

	"github.com/WeisseNacht18/url-shortener/internal/storage"
	"github.com/go-chi/chi"
)

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("x-user-id")
	shortLink := chi.URLParam(r, "id")
	value, ok := storage.GetURLFromStorage(userID, shortLink)
	if ok {
		http.Redirect(w, r, value, http.StatusTemporaryRedirect)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
