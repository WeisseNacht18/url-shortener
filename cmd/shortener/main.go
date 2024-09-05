package main

import (
	"io"
	"net/http"

	"github.com/google/uuid"
)

var (
	server_url string = "http://localhost:8080/"
	short_urls map[string]string
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost && r.Header.Get("Content-Type") == "text/plain" {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		link := string(body)
		id := uuid.New()
		short_link := id.String()[:8]
		short_urls[short_link] = link
		w.Write([]byte(server_url + short_link))
	} else if r.Method == http.MethodGet {
		id := r.URL.Path[1:]
		value, ok := short_urls[id]
		if ok {
			http.Redirect(w, r, value, http.StatusTemporaryRedirect)
			return
		}
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
