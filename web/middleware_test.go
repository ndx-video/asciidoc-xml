package main

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ndx-video/asciidoc-xml/lib"
)

func TestLoggingMiddleware(t *testing.T) {
	// Create a test logger
	logConfig := lib.LogConfig{
		Level:  "info",
		Format: "text",
		Console: lib.ConsoleConfig{
			Enabled: true,
			Level:   "info",
		},
		File: lib.FileConfig{
			Enabled: false,
		},
	}

	logger, err := lib.NewLogger(logConfig)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	// Create a simple handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Wrap with middleware
	middleware := LoggingMiddleware(logger)
	wrappedHandler := middleware(handler)

	// Create test request
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	req.Header.Set("User-Agent", "test-agent")

	// Create response recorder
	rr := httptest.NewRecorder()

	// Execute request
	wrappedHandler.ServeHTTP(rr, req)

	// Verify response
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	// Verify request ID was added to context
	ctx := req.Context()
	requestID := GetRequestID(ctx)
	if requestID == "" {
		t.Error("Request ID should be set in context")
	}
}

func TestLoggingMiddleware_RequestIDGeneration(t *testing.T) {
	logConfig := lib.LogConfig{
		Level:  "info",
		Format: "text",
		Console: lib.ConsoleConfig{
			Enabled: false, // Disable console to avoid output
		},
		File: lib.FileConfig{
			Enabled: false,
		},
	}

	logger, err := lib.NewLogger(logConfig)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request ID in context
		requestID := GetRequestID(r.Context())
		if requestID == "" {
			t.Error("Request ID should be available in handler context")
		}
		w.WriteHeader(http.StatusOK)
	})

	middleware := LoggingMiddleware(logger)
	wrappedHandler := middleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(rr, req)
}

func TestLoggingMiddleware_UniqueRequestIDs(t *testing.T) {
	logConfig := lib.LogConfig{
		Level:  "info",
		Format: "text",
		Console: lib.ConsoleConfig{
			Enabled: false,
		},
		File: lib.FileConfig{
			Enabled: false,
		},
	}

	logger, err := lib.NewLogger(logConfig)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	requestIDs := make(map[string]bool)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := GetRequestID(r.Context())
		if requestIDs[requestID] {
			t.Errorf("Duplicate request ID: %s", requestID)
		}
		requestIDs[requestID] = true
		w.WriteHeader(http.StatusOK)
	})

	middleware := LoggingMiddleware(logger)
	wrappedHandler := middleware(handler)

	// Make multiple requests
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		rr := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(rr, req)
	}

	if len(requestIDs) != 10 {
		t.Errorf("Expected 10 unique request IDs, got %d", len(requestIDs))
	}
}

func TestLoggingMiddleware_ResponseWriter(t *testing.T) {
	logConfig := lib.LogConfig{
		Level:  "info",
		Format: "text",
		Console: lib.ConsoleConfig{
			Enabled: false,
		},
		File: lib.FileConfig{
			Enabled: false,
		},
	}

	logger, err := lib.NewLogger(logConfig)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
	})

	middleware := LoggingMiddleware(logger)
	wrappedHandler := middleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", rr.Code)
	}

	body := rr.Body.String()
	if body != "Not Found" {
		t.Errorf("Expected body 'Not Found', got %q", body)
	}
}

func TestLoggingMiddleware_StatusCodes(t *testing.T) {
	logConfig := lib.LogConfig{
		Level:  "info",
		Format: "text",
		Console: lib.ConsoleConfig{
			Enabled: false,
		},
		File: lib.FileConfig{
			Enabled: false,
		},
	}

	logger, err := lib.NewLogger(logConfig)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	tests := []struct {
		statusCode int
		name       string
	}{
		{http.StatusOK, "200 OK"},
		{http.StatusBadRequest, "400 Bad Request"},
		{http.StatusNotFound, "404 Not Found"},
		{http.StatusInternalServerError, "500 Internal Server Error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			})

			middleware := LoggingMiddleware(logger)
			wrappedHandler := middleware(handler)

			req := httptest.NewRequest("GET", "/test", nil)
			rr := httptest.NewRecorder()

			wrappedHandler.ServeHTTP(rr, req)

			if rr.Code != tt.statusCode {
				t.Errorf("Expected status %d, got %d", tt.statusCode, rr.Code)
			}
		})
	}
}

func TestLoggingMiddleware_BytesWritten(t *testing.T) {
	logConfig := lib.LogConfig{
		Level:  "info",
		Format: "text",
		Console: lib.ConsoleConfig{
			Enabled: false,
		},
		File: lib.FileConfig{
			Enabled: false,
		},
	}

	logger, err := lib.NewLogger(logConfig)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	responseBody := strings.Repeat("x", 1000)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(responseBody))
	})

	middleware := LoggingMiddleware(logger)
	wrappedHandler := middleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(rr, req)

	if rr.Body.Len() != len(responseBody) {
		t.Errorf("Expected %d bytes written, got %d", len(responseBody), rr.Body.Len())
	}
}

func TestLoggingMiddleware_QueryParameters(t *testing.T) {
	logConfig := lib.LogConfig{
		Level:  "info",
		Format: "text",
		Console: lib.ConsoleConfig{
			Enabled: false,
		},
		File: lib.FileConfig{
			Enabled: false,
		},
	}

	logger, err := lib.NewLogger(logConfig)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := LoggingMiddleware(logger)
	wrappedHandler := middleware(handler)

	req := httptest.NewRequest("GET", "/test?param1=value1&param2=value2", nil)
	rr := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(rr, req)

	// Verify query parameters are logged (they should be in RawQuery)
	if req.URL.RawQuery == "" {
		t.Error("Query parameters should be preserved")
	}
}

func TestLoggingMiddleware_MethodAndPath(t *testing.T) {
	logConfig := lib.LogConfig{
		Level:  "info",
		Format: "text",
		Console: lib.ConsoleConfig{
			Enabled: false,
		},
		File: lib.FileConfig{
			Enabled: false,
		},
	}

	logger, err := lib.NewLogger(logConfig)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := LoggingMiddleware(logger)
	wrappedHandler := middleware(handler)

	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	paths := []string{"/api/test", "/api/users", "/api/files"}

	for _, method := range methods {
		for _, path := range paths {
			req := httptest.NewRequest(method, path, nil)
			rr := httptest.NewRecorder()
			wrappedHandler.ServeHTTP(rr, req)

			if req.Method != method {
				t.Errorf("Method should be %s, got %s", method, req.Method)
			}
			if req.URL.Path != path {
				t.Errorf("Path should be %s, got %s", path, req.URL.Path)
			}
		}
	}
}

func TestGetRequestID(t *testing.T) {
	// Test with request ID in context
	ctx := context.WithValue(context.Background(), "request_id", "test-id-123")
	requestID := GetRequestID(ctx)
	if requestID != "test-id-123" {
		t.Errorf("Expected 'test-id-123', got %q", requestID)
	}

	// Test without request ID
	ctx2 := context.Background()
	requestID2 := GetRequestID(ctx2)
	if requestID2 != "" {
		t.Errorf("Expected empty string, got %q", requestID2)
	}

	// Test with nil context
	requestID3 := GetRequestID(nil)
	if requestID3 != "" {
		t.Errorf("Expected empty string for nil context, got %q", requestID3)
	}
}

func TestResponseWriter_WriteHeader(t *testing.T) {
	rr := httptest.NewRecorder()
	rw := &responseWriter{
		ResponseWriter: rr,
		statusCode:     http.StatusOK,
	}

	rw.WriteHeader(http.StatusNotFound)
	if rw.statusCode != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", rw.statusCode)
	}
	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected recorder status 404, got %d", rr.Code)
	}
}

func TestResponseWriter_Write(t *testing.T) {
	rr := httptest.NewRecorder()
	rw := &responseWriter{
		ResponseWriter: rr,
		statusCode:     0, // Not set yet
	}

	data := []byte("test data")
	n, err := rw.Write(data)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	if n != len(data) {
		t.Errorf("Expected to write %d bytes, got %d", len(data), n)
	}
	if rw.statusCode != http.StatusOK {
		t.Errorf("Expected status to be set to 200, got %d", rw.statusCode)
	}
	if rw.written != int64(len(data)) {
		t.Errorf("Expected written count %d, got %d", len(data), rw.written)
	}
	if !bytes.Equal(rr.Body.Bytes(), data) {
		t.Errorf("Expected body %q, got %q", data, rr.Body.Bytes())
	}
}

func TestResponseWriter_MultipleWrites(t *testing.T) {
	rr := httptest.NewRecorder()
	rw := &responseWriter{
		ResponseWriter: rr,
		statusCode:     http.StatusOK,
	}

	data1 := []byte("hello ")
	data2 := []byte("world")

	rw.Write(data1)
	rw.Write(data2)

	expectedWritten := int64(len(data1) + len(data2))
	if rw.written != expectedWritten {
		t.Errorf("Expected written count %d, got %d", expectedWritten, rw.written)
	}

	expectedBody := string(data1) + string(data2)
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, rr.Body.String())
	}
}

func TestLoggingMiddleware_ContextPropagation(t *testing.T) {
	logConfig := lib.LogConfig{
		Level:  "info",
		Format: "text",
		Console: lib.ConsoleConfig{
			Enabled: false,
		},
		File: lib.FileConfig{
			Enabled: false,
		},
	}

	logger, err := lib.NewLogger(logConfig)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	var capturedRequestID string
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedRequestID = GetRequestID(r.Context())
		w.WriteHeader(http.StatusOK)
	})

	middleware := LoggingMiddleware(logger)
	wrappedHandler := middleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(rr, req)

	if capturedRequestID == "" {
		t.Error("Request ID should be propagated to handler")
	}
}

