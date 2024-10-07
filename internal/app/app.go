package app

import (
	"log"
	"net/http"

	"github.com/WeisseNacht18/url-shortener/internal/config"
	"github.com/WeisseNacht18/url-shortener/internal/http/handlers"
	"github.com/WeisseNacht18/url-shortener/internal/logger"
	"github.com/WeisseNacht18/url-shortener/internal/storage"
	"github.com/go-chi/chi/v5"
)

func Run(config config.Config) {
	logger.Init()
	storage.Init()
	handlers.Init(config.BaseURL)

	router := chi.NewRouter()

	router.Post("/", handlers.WithLogging(handlers.GzipHandle(http.HandlerFunc(handlers.CreateShortURLHandler))))
	router.Post("/api/shorten", handlers.WithLogging(handlers.GzipHandle(http.HandlerFunc(handlers.CreateShortURLWithAPIHandler))))
	router.Get("/{id}", handlers.WithLogging(http.HandlerFunc(handlers.RedirectHandler)))

	err := http.ListenAndServe(config.ServerHost, router)

	log.Printf("Starting server on %s", config.ServerHost)
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
