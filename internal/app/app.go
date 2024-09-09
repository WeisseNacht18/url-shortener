package app

import (
	"io"
	"net/http"
	"strings"

	"github.com/WeisseNacht18/url-shortener/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

var (
	shortUrls map[string]string
	baseURL   string
)

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[1:]
	value, ok := shortUrls[id]
	if ok {
		w.Header().Set("Location", value)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func createShortURLHandler(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.Header.Get("Content-Type"), "text/plain") {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		link := string(body)
		id := uuid.New()
		shortLink := id.String()[:8]
		shortUrls[shortLink] = link
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(baseURL + shortLink))
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func Run(config config.Config) {
	shortUrls = map[string]string{}

	baseURL = config.BaseURL

	router := chi.NewRouter()

	router.Post("/", createShortURLHandler)
	router.Get("/{id}", redirectHandler)

	err := http.ListenAndServe(config.ServerHost, router)

	if err != nil {
		panic(err)
	}
}
