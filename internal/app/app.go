package app

import (
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

var (
	shortUrls map[string]string
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

func createShortUrlHandler(w http.ResponseWriter, r *http.Request) {
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
		w.Write([]byte("http://" + r.Host + "/" + shortLink))
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func Run() {
	shortUrls = map[string]string{}

	router := chi.NewRouter()

	router.Post("/", createShortUrlHandler)
	router.Get("/{id}", redirectHandler)

	err := http.ListenAndServe(":8080", router)

	if err != nil {
		panic(err)
	}
}
