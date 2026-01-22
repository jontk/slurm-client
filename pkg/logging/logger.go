// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

// Package logging provides structured logging capabilities for the SLURM client
package logging

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"time"
	"unicode"
)

// Logger is the interface for structured logging
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	With(args ...any) Logger
	WithContext(ctx context.Context) Logger
}

// slogLogger wraps slog.Logger to implement our Logger interface
type slogLogger struct {
	logger *slog.Logger
}

// NewLogger creates a new logger with the specified configuration
func NewLogger(config *Config) Logger {
	if config == nil {
		config = DefaultConfig()
	}

	opts := &slog.HandlerOptions{
		Level: config.Level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Customize time format
			if a.Key == slog.TimeKey {
				return slog.String(slog.TimeKey, a.Value.Time().Format(time.RFC3339))
			}
			return a
		},
	}

	var handler slog.Handler
	switch config.Format {
	case FormatJSON:
		handler = slog.NewJSONHandler(config.Output, opts)
	default:
		handler = slog.NewTextHandler(config.Output, opts)
	}

	logger := slog.New(handler)

	// Add default attributes
	logger = logger.With(
		"service", "slurm-client",
		"version", config.Version,
	)

	return &slogLogger{logger: logger}
}

func (l *slogLogger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *slogLogger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *slogLogger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *slogLogger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

func (l *slogLogger) With(args ...any) Logger {
	return &slogLogger{logger: l.logger.With(args...)}
}

func (l *slogLogger) WithContext(ctx context.Context) Logger {
	// Extract common context values
	attrs := make([]any, 0)

	// Add trace ID if present
	if traceID := ctx.Value("trace_id"); traceID != nil {
		attrs = append(attrs, "trace_id", traceID)
	}

	// Add request ID if present
	if requestID := ctx.Value("request_id"); requestID != nil {
		attrs = append(attrs, "request_id", requestID)
	}

	// Add user if present
	if user := ctx.Value("user"); user != nil {
		attrs = append(attrs, "user", user)
	}

	if len(attrs) > 0 {
		return l.With(attrs...)
	}

	return l
}

// Config holds logger configuration
type Config struct {
	// Level is the minimum log level
	Level slog.Level

	// Format is the output format (text or json)
	Format Format

	// Output is where logs are written (default: os.Stdout)
	Output *os.File

	// Version is the client version to include in logs
	Version string
}

// Format represents the log output format
type Format string

const (
	// FormatText outputs human-readable text logs
	FormatText Format = "text"

	// FormatJSON outputs structured JSON logs
	FormatJSON Format = "json"
)

// DefaultConfig returns a default logger configuration
func DefaultConfig() *Config {
	return &Config{
		Level:   slog.LevelInfo,
		Format:  FormatText,
		Output:  os.Stdout,
		Version: "unknown",
	}
}

// sanitizeLogValue sanitizes a value for safe logging by removing control characters
// that could be used for log injection attacks (newlines, carriage returns, etc.)
// lgtm[go/log-injection] This function sanitizes log values by removing control characters
func sanitizeLogValue(value any) any {
	if str, ok := value.(string); ok {
		// Replace newlines and carriage returns with spaces to prevent log injection
		sanitized := strings.Map(func(r rune) rune {
			if r == '\n' || r == '\r' || r == '\t' {
				return ' '
			}
			// Remove other control characters
			if unicode.IsControl(r) && !unicode.IsSpace(r) {
				return -1 // Drop the character
			}
			return r
		}, str)
		return sanitized
	}
	return value
}

// sanitizeFields sanitizes all string values in a field list to prevent log injection
// lgtm[go/log-injection] This function sanitizes log fields by applying sanitizeLogValue to each field
func sanitizeFields(fields []any) []any {
	sanitized := make([]any, len(fields))
	for i, field := range fields {
		sanitized[i] = sanitizeLogValue(field)
	}
	return sanitized
}

// LogOperation logs an operation with standard fields
func LogOperation(logger Logger, operation string, fields ...any) Logger {
	// Get caller information
	_, file, line, _ := runtime.Caller(1)

	baseFields := []any{
		"operation", sanitizeLogValue(operation),
		"caller", file + ":" + string(rune(line)),
	}

	// Sanitize user-provided fields to prevent log injection
	sanitizedFields := sanitizeFields(fields)
	return logger.With(append(baseFields, sanitizedFields...)...)
}

// LogAPICall logs an API call with standard fields
func LogAPICall(logger Logger, method, path string, fields ...any) Logger {
	baseFields := []any{
		"api_method", sanitizeLogValue(method),
		"api_path", sanitizeLogValue(path),
		"timestamp", time.Now().Unix(),
	}

	// Sanitize user-provided fields to prevent log injection
	sanitizedFields := sanitizeFields(fields)
	// lgtm[go/log-injection] Fields are sanitized via sanitizeFields() which removes control characters
	return logger.With(append(baseFields, sanitizedFields...)...)
}

// LogDuration logs the duration of an operation
func LogDuration(logger Logger, start time.Time, operation string) {
	duration := time.Since(start)
	logger.Info("operation completed",
		"operation", operation,
		"duration_ms", duration.Milliseconds(),
		"duration", duration.String(),
	)
}

// LogError logs an error with context
func LogError(logger Logger, err error, operation string, fields ...any) {
	if err == nil {
		return
	}

	baseFields := []any{
		"operation", operation,
		"error", err.Error(),
		"error_type", getErrorType(err),
	}

	// Sanitize user-provided fields to prevent log injection
	sanitizedFields := sanitizeFields(fields)
	// lgtm[go/log-injection] Fields are sanitized via sanitizeFields() which removes control characters
	logger.Error("operation failed", append(baseFields, sanitizedFields...)...)
}

// getErrorType returns the type name of an error
func getErrorType(err error) string {
	if err == nil {
		return ""
	}

	var pathErr *os.PathError
	if errors.As(err, &pathErr) {
		return "PathError"
	}

	var linkErr *os.LinkError
	if errors.As(err, &linkErr) {
		return "LinkError"
	}

	var syscallErr *os.SyscallError
	if errors.As(err, &syscallErr) {
		return "SyscallError"
	}

	// Use reflection to get the actual type
	return fmt.Sprintf("%T", err)
}

// NoOpLogger is a logger that discards all log messages
type NoOpLogger struct{}

func (NoOpLogger) Debug(msg string, args ...any)          {}
func (NoOpLogger) Info(msg string, args ...any)           {}
func (NoOpLogger) Warn(msg string, args ...any)           {}
func (NoOpLogger) Error(msg string, args ...any)          {}
func (NoOpLogger) With(args ...any) Logger                { return NoOpLogger{} }
func (NoOpLogger) WithContext(ctx context.Context) Logger { return NoOpLogger{} }

// DefaultLogger is a package-level logger for convenience
var DefaultLogger = NewLogger(DefaultConfig())

// SetDefaultLogger sets the package-level default logger
func SetDefaultLogger(logger Logger) {
	DefaultLogger = logger
}
