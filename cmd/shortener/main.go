package main

import (
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

var (
	short_urls map[string]string
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost && strings.Contains(r.Header.Get("Content-Type"), "text/plain") {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		link := string(body)
		id := uuid.New()
		short_link := id.String()[:8]
		short_urls[short_link] = link
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("http://" + r.Host + "/" + short_link))
		return
	} else if r.Method == http.MethodGet {
		id := r.URL.Path[1:]
		value, ok := short_urls[id]
		if ok {
			w.Header().Set("Location", value)
		}
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}

func main() {
	short_urls = map[string]string{}

	router := http.NewServeMux()
	router.HandleFunc(`/`, RootHandler)

	err := http.ListenAndServe(":8080", router)

	if err != nil {
		panic(err)
	}
}
