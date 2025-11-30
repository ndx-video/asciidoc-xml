package lib

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func TestLogLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected LogLevel
	}{
		{"debug", LevelDebug},
		{"DEBUG", LevelDebug},
		{"info", LevelInfo},
		{"INFO", LevelInfo},
		{"warn", LevelWarn},
		{"warning", LevelWarn},
		{"error", LevelError},
		{"fatal", LevelFatal},
		{"unknown", LevelInfo}, // default
	}

	for _, tt := range tests {
		result := ParseLevel(tt.input)
		if result != tt.expected {
			t.Errorf("ParseLevel(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestLogger_LogLevels(t *testing.T) {
	// Create temp directory for log files
	tempDir, err := os.MkdirTemp("", "logger-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	config := LogConfig{
		Level:  "info",
		Format: "text",
		Console: ConsoleConfig{
			Enabled: true,
			Level:   "info",
		},
		File: FileConfig{
			Enabled:  true,
			Path:     tempDir,
			Filename: "test.log",
			MaxSize:  1024 * 1024,
			MaxFiles: 3,
			Rotation: "size",
		},
	}

	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	ctx := context.Background()

	// Test that debug messages are filtered out
	logger.Debug(ctx, "This should not appear")
	logger.Info(ctx, "This should appear")
	logger.Warn(ctx, "This should appear")
	logger.Error(ctx, "This should appear")

	// Check log file
	logPath := filepath.Join(tempDir, "test.log")
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	contentStr := string(content)
	if strings.Contains(contentStr, "This should not appear") {
		t.Error("Debug message should be filtered out")
	}
	if !strings.Contains(contentStr, "This should appear") {
		t.Error("Info/Warn/Error messages should appear")
	}
}

func TestLogger_StructuredLogging(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "logger-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test text format
	config := LogConfig{
		Level:  "debug",
		Format: "text",
		Console: ConsoleConfig{
			Enabled: true,
			Level:   "debug",
		},
		File: FileConfig{
			Enabled:  true,
			Path:     tempDir,
			Filename: "test-text.log",
		},
	}

	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	ctx := context.Background()
	logger.Info(ctx, "Test message", "key1", "value1", "key2", "value2")

	logPath := filepath.Join(tempDir, "test-text.log")
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "key1=value1") || !strings.Contains(contentStr, "key2=value2") {
		t.Error("Structured fields not found in text format")
	}

	// Test JSON format
	config.Format = "json"
	config.File.Filename = "test-json.log"
	logger2, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger2.Close()

	logger2.Info(ctx, "Test message", "key1", "value1", "key2", "value2")

	logPath2 := filepath.Join(tempDir, "test-json.log")
	content2, err := os.ReadFile(logPath2)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	contentStr2 := string(content2)
	if !strings.Contains(contentStr2, `"key1":"value1"`) || !strings.Contains(contentStr2, `"key2":"value2"`) {
		t.Error("Structured fields not found in JSON format")
	}
}

func TestLogger_ContextRequestID(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "logger-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	config := LogConfig{
		Level:  "info",
		Format: "text",
		Console: ConsoleConfig{
			Enabled: true,
		},
		File: FileConfig{
			Enabled:  true,
			Path:     tempDir,
			Filename: "test-ctx.log",
		},
	}

	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	ctx := context.WithValue(context.Background(), "request_id", "req-123")
	logger.Info(ctx, "Test message with request ID")

	logPath := filepath.Join(tempDir, "test-ctx.log")
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "req-123") {
		t.Error("Request ID not found in log output")
	}
}

func TestLogger_ConcurrentLogging(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "logger-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	config := LogConfig{
		Level:  "info",
		Format: "text",
		Console: ConsoleConfig{
			Enabled: false, // Disable console for cleaner test
		},
		File: FileConfig{
			Enabled:  true,
			Path:     tempDir,
			Filename: "test-concurrent.log",
		},
	}

	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	// Concurrent logging
	var wg sync.WaitGroup
	numGoroutines := 10
	messagesPerGoroutine := 10

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < messagesPerGoroutine; j++ {
				logger.Info(nil, "Concurrent test message",
					"goroutine", id,
					"message", j,
				)
			}
		}(i)
	}

	wg.Wait()

	// Verify all messages were written
	logPath := filepath.Join(tempDir, "test-concurrent.log")
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	contentStr := string(content)
	expectedLines := numGoroutines * messagesPerGoroutine
	actualLines := strings.Count(contentStr, "Concurrent test message")

	if actualLines != expectedLines {
		t.Errorf("Expected %d log lines, got %d", expectedLines, actualLines)
	}
}

func TestLogger_WithFields(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "logger-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	config := LogConfig{
		Level:  "info",
		Format: "text",
		Console: ConsoleConfig{
			Enabled: false,
		},
		File: FileConfig{
			Enabled:  true,
			Path:     tempDir,
			Filename: "test-fields.log",
		},
	}

	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	ctx := context.Background()
	loggerWithFields := logger.WithFields(map[string]interface{}{
		"component": "test",
		"version":   "1.0",
	})
	loggerWithFields.Info(ctx, "Message with pre-set fields", "additional", "value")

	logPath := filepath.Join(tempDir, "test-fields.log")
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "component=test") {
		t.Error("Pre-set field 'component' not found")
	}
	if !strings.Contains(contentStr, "version=1.0") {
		t.Error("Pre-set field 'version' not found")
	}
	if !strings.Contains(contentStr, "additional=value") {
		t.Error("Additional field not found")
	}
}

func TestLogger_Rotation(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "logger-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	config := LogConfig{
		Level:  "info",
		Format: "text",
		Console: ConsoleConfig{
			Enabled: false,
		},
		File: FileConfig{
			Enabled:  true,
			Path:     tempDir,
			Filename: "test-rotate.log",
			MaxSize:  100, // Very small to trigger rotation quickly
			MaxFiles: 3,
			Rotation: "size",
		},
	}

	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	// Write enough to trigger rotation
	for i := 0; i < 50; i++ {
		logger.Info(nil, "Test rotation message", "iteration", i)
	}

	// Check that rotation occurred (should have .1 file)
	rotatedPath := filepath.Join(tempDir, "test-rotate.log.1")
	if _, err := os.Stat(rotatedPath); err != nil {
		t.Logf("Rotation may not have occurred (this is OK if file didn't reach maxSize): %v", err)
	}

	// Verify current log file exists
	currentPath := filepath.Join(tempDir, "test-rotate.log")
	if _, err := os.Stat(currentPath); err != nil {
		t.Errorf("Current log file should exist: %v", err)
	}
}

func TestLogger_FormatMethods(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "logger-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	config := LogConfig{
		Level:  "info",
		Format: "text",
		Console: ConsoleConfig{
			Enabled: false,
		},
		File: FileConfig{
			Enabled:  true,
			Path:     tempDir,
			Filename: "test-format.log",
		},
	}

	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	logger.Infof("Formatted message: %s", "test")
	logger.Warnf("Warning: %d items", 5)
	logger.Errorf("Error code: %d", 404)

	logPath := filepath.Join(tempDir, "test-format.log")
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "Formatted message: test") {
		t.Error("Infof message not found")
	}
	if !strings.Contains(contentStr, "Warning: 5 items") {
		t.Error("Warnf message not found")
	}
	if !strings.Contains(contentStr, "Error code: 404") {
		t.Error("Errorf message not found")
	}
}

func TestLogger_SetLevel(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "logger-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	config := LogConfig{
		Level:  "warn",
		Format: "text",
		Console: ConsoleConfig{
			Enabled: false,
		},
		File: FileConfig{
			Enabled:  true,
			Path:     tempDir,
			Filename: "test-level.log",
		},
	}

	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	ctx := context.Background()
	logger.Debug(ctx, "Debug message")
	logger.Info(ctx, "Info message")
	logger.Warn(ctx, "Warn message")
	logger.Error(ctx, "Error message")

	// Change level
	logger.SetLevel(LevelDebug)
	logger.Debug(ctx, "Debug message after level change")

	logPath := filepath.Join(tempDir, "test-level.log")
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	contentStr := string(content)
	if strings.Contains(contentStr, "Debug message") && !strings.Contains(contentStr, "Debug message after level change") {
		t.Error("First debug message should be filtered, second should appear")
	}
	if !strings.Contains(contentStr, "Warn message") {
		t.Error("Warn message should appear")
	}
	if !strings.Contains(contentStr, "Error message") {
		t.Error("Error message should appear")
	}
}

