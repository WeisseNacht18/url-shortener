package handlers

import (
	"net/http"

	"github.com/WeisseNacht18/url-shortener/internal/storage"
)

func PingHandler(w http.ResponseWriter, r *http.Request) {
	err := storage.CheckConnection()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
