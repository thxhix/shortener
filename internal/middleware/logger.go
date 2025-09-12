package middleware

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

type (
	// Структура для хранения сведений об ответе
	responseData struct {
		status int
		size   int
	}
	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

// Write writes the response body and records its size in bytes.
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

// WriteHeader writes the HTTP status code to the underlying ResponseWriter
// and records it for logging.
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// WithLogging returns a middleware that logs HTTP requests and responses
// using the provided zap.SugaredLogger.
//
// It logs the following information for each request:
//   - URI
//   - HTTP method
//   - Response status code
//   - Duration of request processing
//   - Response size in bytes
func WithLogging(logger *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

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

			logger.Infoln(
				"uri", r.URL.String(),
				"method", r.Method,
				"status", responseData.status, // http_code
				"duration", duration,
				"size", responseData.size, // Размер
			)
		})
	}
}
