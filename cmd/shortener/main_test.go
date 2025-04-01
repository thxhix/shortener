// пакеты исполняемых приложений должны называться main
package main

import (
	"bytes"
	"github.com/thxhix/shortener/internal/app/database"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"github.com/thxhix/shortener/internal/app/config"
	"github.com/thxhix/shortener/internal/app/router"
)

var cfg config.Config
var r *chi.Mux

func TestMain(m *testing.M) {
	cfg = *config.NewConfig()
	db, err := database.NewDatabase(&cfg)
	if err != nil {
		panic(err)
	}
	r = router.NewRouter(&cfg, db)

	os.Exit(m.Run())
}

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
				statusCode:  http.StatusCreated,
			},
			action: "/",
			method: http.MethodPost,
			body:   "https://ya.ru",
		},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			body := bytes.NewBufferString(tt.body)

			req, err := http.NewRequest(tt.method, tt.action, body)
			if err != nil {
				panic(err)
			}
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			require.Equal(t, tt.want.statusCode, w.Code, "Код ответа не совпадает с ожидаемым")
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
		name   string
		want   want
		action string
		method string
	}{
		{
			name: "No redirect request",
			want: want{
				contentType: "text/plain",
				statusCode:  400,
			},
			action: "/no-redirect",
			method: http.MethodGet,
		},
		{
			name: "Has redirect request",
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusTemporaryRedirect,
				header:      "https://ya.ru",
			},
			action: "/testHash",
			method: http.MethodGet,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.action, nil)
			if err != nil {
				panic(err)
			}
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			require.Equal(t, tt.want.statusCode, w.Code, "Код ответа не совпадает с ожидаемым")
			require.Equal(t, tt.want.header, w.Header().Get("Location"), "Header Location не совпадает с ожидаемым")
		})
	}
}

func Test_APIStoreLink(t *testing.T) {
	type want struct {
		contentType  string
		statusCode   int
		jsonResponse string
	}

	tests := []struct {
		name        string
		action      string
		method      string
		body        string
		contentType string
		want        want
	}{
		{
			name:        "API store link request",
			action:      "/api/shorten",
			method:      http.MethodPost,
			body:        "{\"url\": \"https://test.ru\"}",
			contentType: "application/json",

			want: want{
				contentType:  "application/json",
				statusCode:   http.StatusCreated,
				jsonResponse: "https://ya.ru",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.action, strings.NewReader(tt.body))
			if err != nil {
				panic(err)
			}
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			require.Equal(t, tt.want.statusCode, w.Code, "Код ответа не совпадает с ожидаемым")
			require.Equal(t, tt.want.contentType, w.Header().Get("Content-Type"), "Content-Type ответа не совпадает с ожидаемым")
		})
	}
}
