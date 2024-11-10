package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/WeisseNacht18/url-shortener/internal/storage"
)

type (
	ShortenBatchRequest struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}

	ShortenBatchResponse struct {
		CorrelationID string `json:"correlation_id"`
		ShortURL      string `json:"short_url"`
	}
)

func CreateShortURLBatchHandler(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		var content []ShortenBatchRequest

		err := json.NewDecoder(r.Body).Decode(&content)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		data := map[string]string{}

		for _, v := range content {
			data[v.CorrelationID] = v.OriginalURL
		}

		userID := r.Header.Get("x-user-id")

		result, err := storage.AddArrayOfURLToStorage(userID, data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var response []ShortenBatchResponse

		for correlationID, shortURL := range result {
			response = append(response, ShortenBatchResponse{
				CorrelationID: correlationID,
				ShortURL:      BaseURL + "/" + shortURL,
			})
		}

		responseContent, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		w.Write(responseContent)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
