package app

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_CreateShortUrl(t *testing.T) {
	type want struct {
		code        int
		method      string
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
			name: "positive test #1",
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
			name: "positive test #2",
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
			name: "negative test #1",
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
	shortUrls = map[string]string{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.data.url))
			request.Header.Set("Content-Type", test.data.contentType)

			w := httptest.NewRecorder()
			rootHandler(w, request)

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
			name: "positive test #1",
			url:  "/abcdefgh",
			want: want{
				code: 307,
			},
		},
		{
			name: "positive test #2",
			url:  "/dcbahgfe",
			want: want{
				code: 307,
			},
		},
		{
			name: "negative test #1",
			url:  "/ffffffff",
			want: want{
				code: 400,
			},
		},
	}
	shortUrls = map[string]string{
		"abcdefgh": "https://ya.ru",
		"dcbahgfe": "https://mail.ru",
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, test.url, nil)

			w := httptest.NewRecorder()
			rootHandler(w, request)

			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, test.want.code, res.StatusCode)
		})
	}
}
