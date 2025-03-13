// пакеты исполняемых приложений должны называться main
package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"github.com/thxhix/shortener/internal/app/config"
	"github.com/thxhix/shortener/internal/app/router"
)

var cfg config.Config
var r *chi.Mux

func TestMain(m *testing.M) {
	r = router.NewRouter(&cfg)
	cfg = *config.NewConfig()

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
