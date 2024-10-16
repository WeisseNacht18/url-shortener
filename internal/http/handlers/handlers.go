package handlers

import (
	"github.com/WeisseNacht18/url-shortener/internal/http/handlers/api"
	"github.com/go-chi/chi/v5"
)

func AddHandlersToRouter(router *chi.Mux) {
	api.New(BaseURL)
	router.Post("/", CreateShortURLHandler)
	router.Post("/api/shorten", api.CreateShortURLWithAPIHandler)
	router.Get("/{id}", RedirectHandler)
}

var (
	BaseURL string
)

func New(baseURL string) {
	BaseURL = baseURL
}
