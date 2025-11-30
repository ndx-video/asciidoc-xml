package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// LogLevel represents the severity level of a log entry
type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// ParseLevel converts a string to a LogLevel
func ParseLevel(level string) LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "warn", "warning":
		return LevelWarn
	case "error":
		return LevelError
	case "fatal":
		return LevelFatal
	default:
		return LevelInfo
	}
}

// LogConfig holds configuration for the logger
type LogConfig struct {
	Level  string      `json:"level"`
	Format string      `json:"format"` // "text" or "json"
	File   FileConfig  `json:"file"`
	Console ConsoleConfig `json:"console"`
}

// FileConfig holds file logging configuration
type FileConfig struct {
	Enabled  bool   `json:"enabled"`
	Path     string `json:"path"`
	Filename string `json:"filename"`
	MaxSize  int64  `json:"maxSize"`
	MaxFiles int    `json:"maxFiles"`
	Rotation string `json:"rotation"` // "size", "daily", "hourly"
}

// ConsoleConfig holds console logging configuration
type ConsoleConfig struct {
	Enabled bool   `json:"enabled"`
	Level   string `json:"level"`
}

// Logger provides structured logging with multiple outputs and rotation
type Logger struct {
	level      LogLevel
	consoleLevel LogLevel
	writers    []io.Writer
	fileWriter *os.File
	mu         sync.Mutex
	format     string // "text" or "json"
	maxSize    int64
	maxFiles   int
	logDir     string
	baseName   string
	rotation   string
	lastRotate time.Time
	rotateTicker *time.Ticker
	stopRotate   chan struct{}
}

// NewLogger creates a new logger with the given configuration
func NewLogger(config LogConfig) (*Logger, error) {
	logger := &Logger{
		level:        ParseLevel(config.Level),
		consoleLevel: ParseLevel(config.Console.Level),
		format:       config.Format,
		maxSize:      config.File.MaxSize,
		maxFiles:     config.File.MaxFiles,
		logDir:       config.File.Path,
		baseName:     config.File.Filename,
		rotation:     config.File.Rotation,
		lastRotate:   time.Now(),
		stopRotate:   make(chan struct{}),
	}

	if logger.format == "" {
		logger.format = "text"
	}
	if logger.consoleLevel == 0 && config.Console.Level == "" {
		logger.consoleLevel = logger.level
	}

	var writers []io.Writer

	if config.Console.Enabled {
		writers = append(writers, os.Stdout)
	}

	if config.File.Enabled {
		if err := logger.setupLogFile(); err != nil {
			return nil, fmt.Errorf("failed to setup log file: %w", err)
		}
		if logger.fileWriter != nil {
			writers = append(writers, logger.fileWriter)
		}
	}

	logger.writers = writers

	// Start time-based rotation if needed
	if config.File.Enabled && (logger.rotation == "daily" || logger.rotation == "hourly") {
		go logger.startTimeBasedRotation()
	}

	return logger, nil
}

// setupLogFile creates and opens the log file
func (l *Logger) setupLogFile() error {
	if l.logDir == "" {
		l.logDir = "./logs"
	}
	if l.baseName == "" {
		l.baseName = "adc.log"
	}

	if err := os.MkdirAll(l.logDir, 0755); err != nil {
		return err
	}

	logPath := filepath.Join(l.logDir, l.baseName)
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	l.fileWriter = file
	return nil
}

// startTimeBasedRotation starts a background goroutine for time-based rotation
func (l *Logger) startTimeBasedRotation() {
	var tickerDuration time.Duration
	switch l.rotation {
	case "daily":
		tickerDuration = 1 * time.Hour // Check every hour, rotate at midnight
	case "hourly":
		tickerDuration = 1 * time.Minute // Check every minute, rotate at top of hour
	default:
		return
	}

	l.rotateTicker = time.NewTicker(tickerDuration)
	defer l.rotateTicker.Stop()

	for {
		select {
		case <-l.rotateTicker.C:
			l.checkTimeBasedRotation()
		case <-l.stopRotate:
			return
		}
	}
}

// checkTimeBasedRotation checks if time-based rotation is needed
func (l *Logger) checkTimeBasedRotation() {
	now := time.Now()
	shouldRotate := false

	switch l.rotation {
	case "daily":
		// Rotate if we've crossed midnight
		if now.Day() != l.lastRotate.Day() || now.Month() != l.lastRotate.Month() || now.Year() != l.lastRotate.Year() {
			shouldRotate = true
		}
	case "hourly":
		// Rotate if we've crossed the hour
		if now.Hour() != l.lastRotate.Hour() || now.Day() != l.lastRotate.Day() {
			shouldRotate = true
		}
	}

	if shouldRotate {
		if err := l.rotate(); err != nil {
			// Log error to stderr since file rotation failed
			fmt.Fprintf(os.Stderr, "Log rotation failed: %v\n", err)
		} else {
			l.lastRotate = now
		}
	}
}

// rotateIfNeeded checks if rotation is needed and performs it
func (l *Logger) rotateIfNeeded() error {
	if l.fileWriter == nil || l.rotation != "size" {
		return nil
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	info, err := l.fileWriter.Stat()
	if err != nil {
		return err
	}

	if info.Size() < l.maxSize {
		return nil
	}

	return l.rotate()
}

// rotate performs log file rotation
func (l *Logger) rotate() error {
	if l.fileWriter == nil {
		return nil
	}

	// Close current file
	if err := l.fileWriter.Close(); err != nil {
		return err
	}

	currentPath := filepath.Join(l.logDir, l.baseName)

	// Rotate existing files (move .1 to .2, .2 to .3, etc.)
	for i := l.maxFiles - 1; i > 0; i-- {
		oldPath := filepath.Join(l.logDir, fmt.Sprintf("%s.%d", l.baseName, i))
		newPath := filepath.Join(l.logDir, fmt.Sprintf("%s.%d", l.baseName, i+1))
		if _, err := os.Stat(oldPath); err == nil {
			os.Rename(oldPath, newPath)
		}
	}

	// Move current to .1
	rotatedPath := filepath.Join(l.logDir, fmt.Sprintf("%s.1", l.baseName))
	if _, err := os.Stat(currentPath); err == nil {
		if err := os.Rename(currentPath, rotatedPath); err != nil {
			return err
		}
	}

	// Delete files beyond maxFiles
	for i := l.maxFiles + 1; i <= l.maxFiles+10; i++ {
		oldPath := filepath.Join(l.logDir, fmt.Sprintf("%s.%d", l.baseName, i))
		os.Remove(oldPath)
	}

	// Create new file
	file, err := os.Create(currentPath)
	if err != nil {
		return err
	}
	l.fileWriter = file

	// Update writers slice
	var writers []io.Writer
	for _, w := range l.writers {
		if w != os.Stdout && w != os.Stderr {
			continue
		}
		writers = append(writers, w)
	}
	writers = append(writers, file)
	l.writers = writers

	return nil
}

// getRequestID extracts request ID from context
func getRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if id, ok := ctx.Value("request_id").(string); ok {
		return id
	}
	return ""
}

// formatFields converts key-value pairs to a map
func formatFields(fields ...interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for i := 0; i < len(fields)-1; i += 2 {
		if key, ok := fields[i].(string); ok {
			result[key] = fields[i+1]
		}
	}
	return result
}

// log writes a log entry
func (l *Logger) log(level LogLevel, ctx context.Context, msg string, fields ...interface{}) {
	// Check if this level should be logged
	if level < l.level {
		return
	}

	requestID := getRequestID(ctx)
	fieldMap := formatFields(fields...)
	timestamp := time.Now().Format(time.RFC3339)

	var logLine string
	if l.format == "json" {
		logLine = l.formatJSON(timestamp, level, requestID, msg, fieldMap)
	} else {
		logLine = l.formatText(timestamp, level, requestID, msg, fieldMap)
	}

	// Rotate if needed (size-based)
	if l.rotation == "size" {
		if err := l.rotateIfNeeded(); err != nil {
			fmt.Fprintf(os.Stderr, "Log rotation failed: %v\n", err)
		}
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// Write to all writers
	for _, writer := range l.writers {
		// For console, check console level
		if writer == os.Stdout || writer == os.Stderr {
			if level < l.consoleLevel {
				continue
			}
		}
		fmt.Fprintln(writer, logLine)
	}

	// For fatal, exit after logging
	if level == LevelFatal {
		os.Exit(1)
	}
}

// formatText formats a log entry as text
func (l *Logger) formatText(timestamp string, level LogLevel, requestID, msg string, fields map[string]interface{}) string {
	var parts []string
	parts = append(parts, fmt.Sprintf("[%s]", timestamp))
	parts = append(parts, fmt.Sprintf("[%s]", level.String()))
	if requestID != "" {
		parts = append(parts, fmt.Sprintf("[%s]", requestID))
	}
	parts = append(parts, msg)

	for k, v := range fields {
		parts = append(parts, fmt.Sprintf("%s=%v", k, v))
	}

	return strings.Join(parts, " ")
}

// formatJSON formats a log entry as JSON
func (l *Logger) formatJSON(timestamp string, level LogLevel, requestID, msg string, fields map[string]interface{}) string {
	entry := map[string]interface{}{
		"timestamp": timestamp,
		"level":      level.String(),
		"message":    msg,
	}
	if requestID != "" {
		entry["request_id"] = requestID
	}
	for k, v := range fields {
		entry[k] = v
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Sprintf(`{"error":"failed to marshal log entry: %v"}`, err)
	}
	return string(data)
}

// Debug logs a debug message
func (l *Logger) Debug(ctx context.Context, msg string, fields ...interface{}) {
	l.log(LevelDebug, ctx, msg, fields...)
}

// Info logs an info message
func (l *Logger) Info(ctx context.Context, msg string, fields ...interface{}) {
	l.log(LevelInfo, ctx, msg, fields...)
}

// Warn logs a warning message
func (l *Logger) Warn(ctx context.Context, msg string, fields ...interface{}) {
	l.log(LevelWarn, ctx, msg, fields...)
}

// Error logs an error message
func (l *Logger) Error(ctx context.Context, msg string, fields ...interface{}) {
	l.log(LevelError, ctx, msg, fields...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(ctx context.Context, msg string, fields ...interface{}) {
	l.log(LevelFatal, ctx, msg, fields...)
}

// Debugf logs a formatted debug message without context
func (l *Logger) Debugf(msg string, args ...interface{}) {
	l.Debug(nil, fmt.Sprintf(msg, args...))
}

// Infof logs a formatted info message without context
func (l *Logger) Infof(msg string, args ...interface{}) {
	l.Info(nil, fmt.Sprintf(msg, args...))
}

// Warnf logs a formatted warning message without context
func (l *Logger) Warnf(msg string, args ...interface{}) {
	l.Warn(nil, fmt.Sprintf(msg, args...))
}

// Errorf logs a formatted error message without context
func (l *Logger) Errorf(msg string, args ...interface{}) {
	l.Error(nil, fmt.Sprintf(msg, args...))
}

// Fatalf logs a formatted fatal message and exits
func (l *Logger) Fatalf(msg string, args ...interface{}) {
	l.Fatal(nil, fmt.Sprintf(msg, args...))
}

// WithFields returns a logger with pre-set fields (for chaining)
func (l *Logger) WithFields(fields map[string]interface{}) *LoggerWithFields {
	return &LoggerWithFields{
		logger: l,
		fields: fields,
	}
}

// WithField returns a logger with a single pre-set field
func (l *Logger) WithField(key string, value interface{}) *LoggerWithFields {
	return &LoggerWithFields{
		logger: l,
		fields: map[string]interface{}{key: value},
	}
}

// LoggerWithFields is a logger with pre-set fields
type LoggerWithFields struct {
	logger *Logger
	fields map[string]interface{}
}

// mergeFields merges pre-set fields with new fields
func (l *LoggerWithFields) mergeFields(newFields ...interface{}) []interface{} {
	result := make([]interface{}, 0, len(l.fields)*2+len(newFields))
	for k, v := range l.fields {
		result = append(result, k, v)
	}
	result = append(result, newFields...)
	return result
}

// Debug logs with pre-set fields
func (l *LoggerWithFields) Debug(ctx context.Context, msg string, fields ...interface{}) {
	l.logger.Debug(ctx, msg, l.mergeFields(fields...)...)
}

// Info logs with pre-set fields
func (l *LoggerWithFields) Info(ctx context.Context, msg string, fields ...interface{}) {
	l.logger.Info(ctx, msg, l.mergeFields(fields...)...)
}

// Warn logs with pre-set fields
func (l *LoggerWithFields) Warn(ctx context.Context, msg string, fields ...interface{}) {
	l.logger.Warn(ctx, msg, l.mergeFields(fields...)...)
}

// Error logs with pre-set fields
func (l *LoggerWithFields) Error(ctx context.Context, msg string, fields ...interface{}) {
	l.logger.Error(ctx, msg, l.mergeFields(fields...)...)
}

// Fatal logs with pre-set fields
func (l *LoggerWithFields) Fatal(ctx context.Context, msg string, fields ...interface{}) {
	l.logger.Fatal(ctx, msg, l.mergeFields(fields...)...)
}

// SetLevel changes the minimum log level
func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// SetFormat changes the log format
func (l *Logger) SetFormat(format string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.format = format
}

// Close closes the logger and cleans up resources
func (l *Logger) Close() error {
	if l.rotateTicker != nil {
		l.rotateTicker.Stop()
		close(l.stopRotate)
	}
	if l.fileWriter != nil {
		return l.fileWriter.Close()
	}
	return nil
}

// GetStack returns the current stack trace
func GetStack() string {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

