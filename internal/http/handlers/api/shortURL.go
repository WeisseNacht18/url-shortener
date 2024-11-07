package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/WeisseNacht18/url-shortener/internal/logger"
	"github.com/WeisseNacht18/url-shortener/internal/storage"
)

var (
	BaseURL string
)

func New(baseURL string) {
	BaseURL = baseURL
}

type (
	ShortenRequest struct {
		URL string `json:"url"`
	}

	ShortenResponse struct {
		Result string `json:"result"`
	}
)

func CreateShortURLWithAPIHandler(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		var content ShortenRequest

		err := json.NewDecoder(r.Body).Decode(&content)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		userID := r.Header.Get("x-user-id")

		logger.Logger.Infoln(userID, "|", content.URL)

		shortLink, hasURL := storage.AddURLToStorage(userID, content.URL)

		response := ShortenResponse{
			Result: BaseURL + "/" + shortLink,
		}

		responseContent, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if hasURL {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusCreated)
		}

		w.Write(responseContent)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
