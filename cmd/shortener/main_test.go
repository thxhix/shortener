// пакеты исполняемых приложений должны называться main
package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func Test_shortLink(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
	}

	tests := []struct {
		name   string
		want   want
		action string
		method string
		answer string
		body   string
	}{
		{
			name: "Empty link request",
			want: want{
				contentType: "text/plain",
				statusCode:  400,
			},
			action: "/",
			method: http.MethodPost,
			body:   "",
		},
		{
			name: "Success request",
			want: want{
				contentType: "text/plain",
				statusCode:  201,
			},
			action: "/",
			method: http.MethodPost,
			body:   "https://ya.ru",
		},
	}
	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			rBody := bytes.NewBufferString(tt.body)

			r := httptest.NewRequest(tt.method, tt.action, rBody)
			w := httptest.NewRecorder()

			hostParams.Set("localhost:8080")
			destinationURL.Set("http://localhost:8080")

			shortLink(w, r)

			// Код ответа
			require.Equal(t, tt.want.statusCode, w.Code, "Код ответа не совпадает с ожидаемым")

			// Проверяем ответ, если вернулся 201
			body, err := io.ReadAll(w.Body)
			require.NoError(t, err, "Не удалось получить ответ")

			parsedURL, err := url.Parse(string(body))
			require.NoError(t, err, "Не удалось прочитать ответ")

			fmt.Println("Ответ сервера:", parsedURL)
		})
	}
}

func Test_getFullLink(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		header      string
	}
	tests := []struct {
		name     string
		want     want
		action   string
		method   string
		database map[string]string
	}{
		{
			name: "No redirect request",
			want: want{
				contentType: "text/plain",
				statusCode:  400,
			},
			action:   "/no-redirect",
			method:   http.MethodGet,
			database: map[string]string{},
		},
		{
			name: "Has redirect request",
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusTemporaryRedirect,
				header:      "https://ya.ru",
			},
			action:   "/link1",
			method:   http.MethodGet,
			database: map[string]string{"link1": "https://ya.ru"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Post("/", shortLink)
			r.Get("/{id}", getFullLink)

			Database = tt.database

			req, err := http.NewRequest(tt.method, tt.action, nil)
			w := httptest.NewRecorder()
			if err != nil {
				panic(err)
			}

			r.ServeHTTP(w, req)

			require.Equal(t, tt.want.statusCode, w.Code, "Код ответа не совпадает с ожидаемым")
			require.Equal(t, tt.want.header, w.Header().Get("Location"), "Header Location не совпадает с ожидаемым")
		})
	}
}
