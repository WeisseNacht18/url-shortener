package middlewares

import (
	"net/http"
	"time"

	"github.com/WeisseNacht18/url-shortener/internal/logger"
)

type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func WithLogging(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		uri := r.RequestURI
		method := r.Method

		responseData := &responseData{
			status: 0,
			size:   0,
		}

		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		h.ServeHTTP(&lw, r)

		duration := time.Since(start)

		logger.Logger.Infoln(
			"request:",
			"uri", uri,
			"method", method,
			"duration", duration,
		)

		logger.Logger.Infoln(
			"response:",
			"status", responseData.status,
			"content-length", responseData.size,
		)
	}

	return http.HandlerFunc(logFn)
}
