package middlewares

import (
	"github.com/go-chi/chi/v5"
)

func AddMiddlewaresToRouter(router *chi.Mux) {
	router.Use(WithLogging)
	router.Use(GzipHandle)
}
