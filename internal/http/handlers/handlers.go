package handlers

import (
	"github.com/WeisseNacht18/url-shortener/internal/http/handlers/api"
	"github.com/WeisseNacht18/url-shortener/internal/http/handlers/api/user"
	"github.com/go-chi/chi/v5"
)

func AddHandlersToRouter(router *chi.Mux) {
	api.New(BaseURL)
	router.Post("/", CreateShortURLHandler)
	router.Post("/api/shorten", api.CreateShortURLWithAPIHandler)
	router.Post("/api/shorten/batch", api.CreateShortURLBatchHandler)
	router.Get("/ping", PingHandler)
	router.Get("/{id}", RedirectHandler)
	router.Get("/api/user/urls", user.URLsHandler)
	router.Delete("/api/user/urls", user.DeleteURLs)
}

var (
	BaseURL string
)

func New(baseURL string) {
	BaseURL = baseURL
}
