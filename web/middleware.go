package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/ndx-video/asciidoc-xml/lib"
)

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    int64
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.statusCode == 0 {
		rw.statusCode = http.StatusOK
	}
	n, err := rw.ResponseWriter.Write(b)
	rw.written += int64(n)
	return n, err
}

// generateRequestID generates a unique request ID
func generateRequestID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// LoggingMiddleware creates HTTP logging middleware
func LoggingMiddleware(logger *lib.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			requestID := generateRequestID()
			ctx := context.WithValue(r.Context(), "request_id", requestID)

			// Log request
			logger.Info(ctx, "HTTP request",
				"method", r.Method,
				"path", r.URL.Path,
				"query", r.URL.RawQuery,
				"remote_addr", r.RemoteAddr,
				"user_agent", r.UserAgent(),
			)

			// Wrap response writer to capture status code
			rw := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			// Process request
			next.ServeHTTP(rw, r.WithContext(ctx))

			// Log response
			duration := time.Since(start)
			logLevel := lib.LevelInfo
			if rw.statusCode >= 500 {
				logLevel = lib.LevelError
			} else if rw.statusCode >= 400 {
				logLevel = lib.LevelWarn
			}

			fields := []interface{}{
				"method", r.Method,
				"path", r.URL.Path,
				"status", rw.statusCode,
				"duration_ms", duration.Milliseconds(),
				"bytes_written", rw.written,
			}

			switch logLevel {
			case lib.LevelError:
				logger.Error(ctx, "HTTP response", fields...)
			case lib.LevelWarn:
				logger.Warn(ctx, "HTTP response", fields...)
			default:
				logger.Info(ctx, "HTTP response", fields...)
			}
		})
	}
}

// GetRequestID extracts request ID from context
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value("request_id").(string); ok {
		return id
	}
	return ""
}

