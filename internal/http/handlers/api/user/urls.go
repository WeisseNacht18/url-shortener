package user

import (
	"encoding/json"
	"net/http"

	"github.com/WeisseNacht18/url-shortener/internal/storage"
)

type Response struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func URLsHandler(w http.ResponseWriter, r *http.Request) {
	userID := 1

	urls := storage.GetAllURLsFromStorage(userID)

	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	data := []Response{}

	for shortURL, originalURL := range urls {
		data = append(data, Response{
			ShortURL:    shortURL,
			OriginalURL: originalURL,
		})
	}

	content, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(content)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
