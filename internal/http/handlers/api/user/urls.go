package user

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/WeisseNacht18/url-shortener/internal/http/handlers/api"
	"github.com/WeisseNacht18/url-shortener/internal/logger"
	"github.com/WeisseNacht18/url-shortener/internal/storage"
)

type Response struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func URLsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("x-user-id")

	logger.Logger.Infoln(userID)

	urls := storage.GetAllURLsFromStorage(userID)

	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	data := []Response{}

	for shortURL, originalURL := range urls {
		data = append(data, Response{
			ShortURL:    api.BaseURL + "/" + shortURL,
			OriginalURL: originalURL,
		})
	}

	content, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Logger.Infoln(string(content))

	w.Header().Set("Content-Type", "application/json")

	w.Write(content)

	w.WriteHeader(http.StatusOK)
}

func DeleteURLs(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		var content []string

		err := json.NewDecoder(r.Body).Decode(&content)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		userID := r.Header.Get("x-user-id")

		//попробовать сделать fanIn
		for _, url := range content {
			storage.DeleteURL(userID, url)
		}

		w.WriteHeader(http.StatusAccepted)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
