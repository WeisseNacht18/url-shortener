package app

import (
	"net/http"

	"github.com/WeisseNacht18/url-shortener/internal/config"
	"github.com/WeisseNacht18/url-shortener/internal/database"
	"github.com/WeisseNacht18/url-shortener/internal/http/handlers"
	"github.com/WeisseNacht18/url-shortener/internal/http/middlewares"
	"github.com/WeisseNacht18/url-shortener/internal/logger"
	"github.com/WeisseNacht18/url-shortener/internal/storage"
	"github.com/go-chi/chi/v5"
)

func Run(config config.Config) {
	logger.Init()

	if config.DatabaseDSN != "" {
		err := database.NewConnection(config.DatabaseDSN)
		if err != nil {
			logger.Logger.Fatalf("Database connection error: %v", err)
		}
	}
	defer database.CloseConnection()

	storage.NewURLStorage(config.FileStoragePath, config.DatabaseDSN)
	handlers.New(config.BaseURL)

	router := chi.NewRouter()

	middlewares.AddMiddlewaresToRouter(router)

	handlers.AddHandlersToRouter(router)

	err := http.ListenAndServe(config.ServerHost, router)

	if err != nil {
		logger.Logger.Fatalf("Server error: %v", err)
	}
}
