package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/WeisseNacht18/url-shortener/internal/logger"
	"github.com/WeisseNacht18/url-shortener/internal/storage"
)

var (
	BaseURL string
)

func Init(baseURL string) {
	BaseURL = baseURL
}

type (
	// берём структуру для хранения сведений об ответе
	responseData struct {
		status int
		size   int
	}

	// добавляем реализацию http.ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
	}

	ShortenRequest struct {
		Url string `json:"url"`
	}

	ShortenResponse struct {
		Result string `json:"result"`
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}

func WithLogging(fn func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	h := http.HandlerFunc(fn)
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		uri := r.RequestURI
		method := r.Method

		responseData := &responseData{
			status: 0,
			size:   0,
		}

		lw := loggingResponseWriter{
			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
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

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/")
	value, ok := storage.GetURLFromStorage(id)
	if ok {
		w.Header().Set("Location", value)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func CreateShortURLHandler(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.Header.Get("Content-Type"), "text/plain") {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		link := string(body)
		shortLink := storage.AddURLToStorage(link)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		content := []byte(BaseURL + "/" + shortLink)
		w.Write(content)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func CreateShortURLWithAPIHandler(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var content ShortenRequest
		err = json.Unmarshal(body, &content)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		shortLink := storage.AddURLToStorage(content.Url)

		response := ShortenResponse{
			Result: BaseURL + "/" + shortLink,
		}

		responseContent, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", strconv.Itoa(len(string(responseContent))))
		w.WriteHeader(http.StatusCreated)

		w.Write(responseContent)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
