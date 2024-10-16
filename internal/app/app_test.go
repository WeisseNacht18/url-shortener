package app

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/WeisseNacht18/url-shortener/internal/http/handlers"
	"github.com/WeisseNacht18/url-shortener/internal/http/handlers/api"
	"github.com/WeisseNacht18/url-shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_CreateShortUrl(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	type data struct {
		url         string
		contentType string
	}
	tests := []struct {
		name string
		data data
		want want
	}{
		{
			name: "create short URL With valid input",
			data: data{
				url:         "https://ya.ru/",
				contentType: "text/plain",
			},
			want: want{
				code:        201,
				contentType: "text/plain",
			},
		},
		{
			name: "create short URL with empty input",
			data: data{
				url:         "",
				contentType: "text/plain",
			},
			want: want{
				code:        201,
				contentType: "text/plain",
			},
		},
		{
			name: "create short URL with incorrect Content-Type",
			data: data{
				url:         "",
				contentType: "application/json",
			},
			want: want{
				code:        400,
				contentType: "text/plain",
			},
		},
	}
	storage.NewEmpty()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.data.url))
			request.Header.Set("Content-Type", test.data.contentType)

			w := httptest.NewRecorder()
			handlers.CreateShortURLHandler(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)

			if res.StatusCode != 400 {
				defer res.Body.Close()
				resBody, err := io.ReadAll(res.Body)

				require.NoError(t, err)
				assert.Contains(t, res.Header.Get("Content-Type"), test.want.contentType)
				assert.NotEqual(t, "", string(resBody))
			}
		})
	}
}

func TestHandler_CreateShortUrlWithAPI(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	type data struct {
		url         string
		contentType string
	}
	tests := []struct {
		name string
		data data
		want want
	}{
		{
			name: "create short URL (API method) With valid input",
			data: data{
				url:         "https://ya.ru/",
				contentType: "application/json",
			},
			want: want{
				code:        201,
				contentType: "application/json",
			},
		},
		{
			name: "create short URL (API method) with empty input",
			data: data{
				url:         "",
				contentType: "application/json",
			},
			want: want{
				code:        201,
				contentType: "application/json",
			},
		},
		{
			name: "create short URL (API method) with incorrect Content-Type",
			data: data{
				url:         "",
				contentType: "text/plain",
			},
			want: want{
				code:        400,
				contentType: "application/json",
			},
		},
	}
	storage.NewEmpty()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			requestBody := api.ShortenRequest{
				URL: test.data.url,
			}
			requestBytes, err := json.Marshal(requestBody)
			assert.NoError(t, err)
			request := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(requestBytes))
			request.Header.Set("Content-Type", test.data.contentType)
			assert.NotNil(t, request)
			w := httptest.NewRecorder()
			api.CreateShortURLWithAPIHandler(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)

			if res.StatusCode != 400 {
				defer res.Body.Close()
				resBody, err := io.ReadAll(res.Body)

				require.NoError(t, err)
				assert.Contains(t, res.Header.Get("Content-Type"), test.want.contentType)
				assert.NotEqual(t, "", string(resBody))
			}
		})
	}
}

func TestHandler_RedirectShortUrl(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "redirect short URL with valid id #1",
			url:  "/abcdef",
			want: want{
				code: 307,
			},
		},
		{
			name: "redirect short URL with valid id #2",
			url:  "/dcbahg",
			want: want{
				code: 307,
			},
		},
		{
			name: "redirect short URL with non-existent id",
			url:  "/ffffffff",
			want: want{
				code: 400,
			},
		},
	}
	shortUrls := map[string]string{
		"abcdef": "https://ya.ru",
		"dcbahg": "https://mail.ru",
	}
	storage.NewWithMap(shortUrls)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, test.url, nil)

			w := httptest.NewRecorder()
			handlers.RedirectHandler(w, request)

			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, test.want.code, res.StatusCode)
		})
	}
}
