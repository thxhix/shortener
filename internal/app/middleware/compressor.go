package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"
)

type compressedResponseWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

func (g compressedResponseWriter) Write(b []byte) (int, error) {
	return g.Writer.Write(b)
}

func CompressorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			compress, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "не удалось разжать gzip", http.StatusBadRequest)
				return
			}
			defer compress.Close()

			r.Body = compress
		}

		// Не сжимаем если неверный Accept-Encoding
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// Не сжимаем если неверный Content-Type
		if !(strings.Contains(w.Header().Get("Content-Type"), "text/html") || strings.Contains(w.Header().Get("Content-Type"), "application/json")) {
			next.ServeHTTP(w, r)
			return
		}

		gz := gzip.NewWriter(w)
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		wrw := compressedResponseWriter{Writer: gz, ResponseWriter: w}
		next.ServeHTTP(wrw, r)
	})
}
