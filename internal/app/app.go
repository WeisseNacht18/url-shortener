package app

import (
	"net/http"

	"github.com/WeisseNacht18/url-shortener/internal/config"
	"github.com/WeisseNacht18/url-shortener/internal/http/handlers"
	"github.com/WeisseNacht18/url-shortener/internal/http/middlewares"
	"github.com/WeisseNacht18/url-shortener/internal/logger"
	"github.com/WeisseNacht18/url-shortener/internal/storage"
	"github.com/go-chi/chi/v5"
)

func Run(config config.Config) {
	logger.Init()

	storage.NewURLStorage(config.FileStoragePath, config.DatabaseDSN)
	defer storage.Close()
	handlers.New(config.BaseURL)

	router := chi.NewRouter()

	middlewares.AddMiddlewaresToRouter(router)

	handlers.AddHandlersToRouter(router)

	err := http.ListenAndServe(config.ServerHost, router)

	if err != nil {
		logger.Logger.Fatalf("Server error: %v", err)
	}
}
