package middleware

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"

	"{{index .App "git"}}/pkg/logger"
)

// responseWriter wraps http.ResponseWriter to capture status code and response body.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
		body:           &bytes.Buffer{},
	}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

// Logging returns middleware that logs HTTP requests and responses.
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Read request body
		var requestBody []byte
		if r.Body != nil {
			requestBody, _ = io.ReadAll(r.Body)
			// Restore the body for the next handler
			r.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Wrap response writer to capture response
		rw := newResponseWriter(w)

		// Process request
		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		// Log request and response
		logger.Info(r.Context(), "HTTP request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("query", r.URL.RawQuery),
			zap.String("remote_addr", r.RemoteAddr),
			zap.String("user_agent", r.UserAgent()),
			zap.String("request_body", string(requestBody)),
			zap.Int("status", rw.statusCode),
			zap.String("response_body", rw.body.String()),
			zap.Int("response_size", len(rw.body.String())),
			zap.Duration("duration", duration),
		)
	})
}
