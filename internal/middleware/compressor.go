package middleware

import (
	"compress/gzip"
	"log"
	"net/http"
	"strings"
)

type compressedResponseWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

// Write writes the given bytes to the underlying gzip.Writer,
// compressing the response body before sending it to the client.
// It satisfies the http.ResponseWriter interface.
func (g compressedResponseWriter) Write(b []byte) (int, error) {
	return g.Writer.Write(b)
}

// CompressorMiddleware is an HTTP middleware that provides gzip
// compression and decompression for requests and responses.
//
//   - If the request has Content-Encoding: gzip, the body is decompressed
//     before passing it to the next handler.
//   - If the client supports Accept-Encoding: gzip and the response Content-Type
//     is "text/html" or "application/json", the response is compressed.
//
// If gzip initialization or closing fails, the middleware writes an error response
// or logs a fatal error.
func CompressorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			compress, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "не удалось разжать gzip", http.StatusBadRequest)
				return
			}
			defer func() {
				err := compress.Close()
				if err != nil {
					log.Fatalf("Error defer close: %v", err)
				}
			}()

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
		defer func() {
			err := gz.Close()
			if err != nil {
				log.Fatalf("Error defer close: %v", err)
			}
		}()

		w.Header().Set("Content-Encoding", "gzip")
		wrw := compressedResponseWriter{Writer: gz, ResponseWriter: w}
		next.ServeHTTP(wrw, r)
	})
}
