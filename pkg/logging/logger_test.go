// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLogger(t *testing.T) {
	t.Run("with config", func(t *testing.T) {
		config := &Config{
			Level:   slog.LevelDebug,
			Format:  FormatJSON,
			Output:  os.Stdout,
			Version: "1.0.0",
		}

		logger := NewLogger(config)
		require.NotNil(t, logger)
		
		// Verify it's the expected type
		slogLogger, ok := logger.(*slogLogger)
		assert.True(t, ok)
		assert.NotNil(t, slogLogger.logger)
	})

	t.Run("with nil config", func(t *testing.T) {
		logger := NewLogger(nil)
		require.NotNil(t, logger)
		
		// Should use default config
		slogLogger, ok := logger.(*slogLogger)
		assert.True(t, ok)
		assert.NotNil(t, slogLogger.logger)
	})
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	
	require.NotNil(t, config)
	assert.Equal(t, slog.LevelInfo, config.Level)
	assert.Equal(t, FormatText, config.Format)
	assert.Equal(t, os.Stdout, config.Output)
	assert.Equal(t, "unknown", config.Version)
}

func TestSlogLogger_LogMethods(t *testing.T) {
	config := &Config{
		Level:   slog.LevelDebug,
		Format:  FormatJSON,
		Output:  os.Stdout,
		Version: "test",
	}

	logger := NewLogger(config)
	
	// Test that methods don't panic
	logger.Debug("debug message", "key", "value")
	logger.Info("info message", "key", "value")
	logger.Warn("warn message", "key", "value")
	logger.Error("error message", "key", "value")
}

func TestSlogLogger_With(t *testing.T) {
	config := &Config{
		Level:   slog.LevelDebug,
		Format:  FormatText,
		Output:  os.Stdout,
		Version: "test",
	}

	logger := NewLogger(config)
	
	// Create a new logger with additional fields
	newLogger := logger.With("component", "test", "user_id", 123)
	
	// Should return a new logger instance
	assert.NotEqual(t, logger, newLogger)
	assert.IsType(t, &slogLogger{}, newLogger)
}

func TestSlogLogger_WithContext(t *testing.T) {
	config := &Config{
		Level:   slog.LevelDebug,
		Format:  FormatText,
		Output:  os.Stdout,
		Version: "test",
	}

	logger := NewLogger(config)
	
	t.Run("context with values", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, "trace_id", "trace-123")
		ctx = context.WithValue(ctx, "request_id", "req-456")
		ctx = context.WithValue(ctx, "user", "john@example.com")
		
		contextLogger := logger.WithContext(ctx)
		
		// Should return a new logger instance with context values
		assert.NotEqual(t, logger, contextLogger)
		assert.IsType(t, &slogLogger{}, contextLogger)
	})

	t.Run("context without values", func(t *testing.T) {
		ctx := context.Background()
		
		contextLogger := logger.WithContext(ctx)
		
		// Should return the same logger since no context values to extract
		assert.Equal(t, logger, contextLogger)
	})

	t.Run("context with some values", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, "trace_id", "trace-123")
		ctx = context.WithValue(ctx, "other_key", "other_value") // This won't be extracted
		
		contextLogger := logger.WithContext(ctx)
		
		// Should return a new logger with only the trace_id
		assert.NotEqual(t, logger, contextLogger)
		assert.IsType(t, &slogLogger{}, contextLogger)
	})
}

func TestLogOperation(t *testing.T) {
	config := &Config{
		Level:   slog.LevelDebug,
		Format:  FormatText,
		Output:  os.Stdout,
		Version: "test",
	}

	logger := NewLogger(config)
	
	operationLogger := LogOperation(logger, "test-operation", "extra", "field")
	
	// Should return a new logger with operation fields
	assert.NotEqual(t, logger, operationLogger)
	assert.IsType(t, &slogLogger{}, operationLogger)
}

func TestLogAPICall(t *testing.T) {
	config := &Config{
		Level:   slog.LevelDebug,
		Format:  FormatText,
		Output:  os.Stdout,
		Version: "test",
	}

	logger := NewLogger(config)
	
	apiLogger := LogAPICall(logger, "GET", "/api/v1/jobs", "extra", "field")
	
	// Should return a new logger with API call fields
	assert.NotEqual(t, logger, apiLogger)
	assert.IsType(t, &slogLogger{}, apiLogger)
}

func TestLogDuration(t *testing.T) {
	config := &Config{
		Level:   slog.LevelDebug,
		Format:  FormatText,
		Output:  os.Stdout,
		Version: "test",
	}

	logger := NewLogger(config)
	
	start := time.Now().Add(-100 * time.Millisecond)
	
	// Should not panic
	LogDuration(logger, start, "test-operation")
}

func TestLogError(t *testing.T) {
	config := &Config{
		Level:   slog.LevelDebug,
		Format:  FormatText,
		Output:  os.Stdout,
		Version: "test",
	}

	logger := NewLogger(config)
	
	t.Run("with error", func(t *testing.T) {
		err := errors.New("test error")
		
		// Should not panic
		LogError(logger, err, "test-operation", "extra", "field")
	})

	t.Run("with nil error", func(t *testing.T) {
		// Should not panic and should not log anything
		LogError(logger, nil, "test-operation", "extra", "field")
	})
}

func TestGetErrorType(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: "",
		},
		{
			name:     "generic error",
			err:      errors.New("test error"),
			expected: "*errors.errorString",
		},
		{
			name:     "path error",
			err:      &os.PathError{Op: "open", Path: "/test", Err: errors.New("not found")},
			expected: "PathError",
		},
		{
			name:     "link error",
			err:      &os.LinkError{Op: "link", Old: "/old", New: "/new", Err: errors.New("failed")},
			expected: "LinkError",
		},
		{
			name:     "syscall error",
			err:      &os.SyscallError{Syscall: "test", Err: errors.New("failed")},
			expected: "SyscallError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getErrorType(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNoOpLogger(t *testing.T) {
	logger := NoOpLogger{}
	
	// All methods should not panic
	logger.Debug("test")
	logger.Info("test")
	logger.Warn("test")
	logger.Error("test")
	
	// With should return NoOpLogger
	withLogger := logger.With("key", "value")
	assert.Equal(t, NoOpLogger{}, withLogger)
	
	// WithContext should return NoOpLogger
	ctx := context.Background()
	contextLogger := logger.WithContext(ctx)
	assert.Equal(t, NoOpLogger{}, contextLogger)
}

func TestDefaultLogger(t *testing.T) {
	// DefaultLogger should be initialized
	assert.NotNil(t, DefaultLogger)
	
	// Should be able to use it
	DefaultLogger.Info("test message")
}

func TestSetDefaultLogger(t *testing.T) {
	originalLogger := DefaultLogger
	
	// Create a new logger
	newLogger := NoOpLogger{}
	
	// Set it as default
	SetDefaultLogger(newLogger)
	
	// Verify it was set
	assert.Equal(t, newLogger, DefaultLogger)
	
	// Restore original logger
	SetDefaultLogger(originalLogger)
}

func TestFormat(t *testing.T) {
	// Test that format constants have expected values
	assert.Equal(t, Format("text"), FormatText)
	assert.Equal(t, Format("json"), FormatJSON)
}

func TestLoggerInterface(t *testing.T) {
	// Verify that slogLogger implements Logger interface
	var _ Logger = (*slogLogger)(nil)
	
	// Verify that NoOpLogger implements Logger interface
	var _ Logger = NoOpLogger{}
}

// TestLoggerOutput tests that the logger actually produces output
func TestLoggerOutput(t *testing.T) {
	t.Run("text format", func(t *testing.T) {
		var buf bytes.Buffer
		
		// Create a custom handler that writes to our buffer
		handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
		logger := &slogLogger{
			logger: slog.New(handler).With("service", "slurm-client", "version", "test"),
		}
		
		logger.Info("test message", "key", "value")
		
		output := buf.String()
		assert.Contains(t, output, "test message")
		assert.Contains(t, output, "key=value")
		assert.Contains(t, output, "service=slurm-client")
	})

	t.Run("json format", func(t *testing.T) {
		var buf bytes.Buffer
		
		// Create a custom handler that writes to our buffer
		handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
		logger := &slogLogger{
			logger: slog.New(handler).With("service", "slurm-client", "version", "test"),
		}
		
		logger.Info("test message", "key", "value")
		
		output := buf.String()
		assert.True(t, json.Valid([]byte(output)), "Output should be valid JSON")
		assert.Contains(t, output, "test message")
		assert.Contains(t, output, "\"key\":\"value\"")
		assert.Contains(t, output, "\"service\":\"slurm-client\"")
	})
}

// TestLogLevels tests that different log levels work correctly
func TestLogLevels(t *testing.T) {
	tests := []struct {
		name        string
		level       slog.Level
		shouldLog   []string
		shouldntLog []string
	}{
		{
			name:        "debug level",
			level:       slog.LevelDebug,
			shouldLog:   []string{"debug", "info", "warn", "error"},
			shouldntLog: []string{},
		},
		{
			name:        "info level",
			level:       slog.LevelInfo,
			shouldLog:   []string{"info", "warn", "error"},
			shouldntLog: []string{"debug"},
		},
		{
			name:        "warn level",
			level:       slog.LevelWarn,
			shouldLog:   []string{"warn", "error"},
			shouldntLog: []string{"debug", "info"},
		},
		{
			name:        "error level",
			level:       slog.LevelError,
			shouldLog:   []string{"error"},
			shouldntLog: []string{"debug", "info", "warn"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			
			handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
				Level: tt.level,
			})
			logger := &slogLogger{
				logger: slog.New(handler),
			}
			
			// Log at different levels
			logger.Debug("debug message")
			logger.Info("info message")
			logger.Warn("warn message")
			logger.Error("error message")
			
			output := buf.String()
			
			for _, should := range tt.shouldLog {
				assert.Contains(t, output, should+" message", "should log %s at level %v", should, tt.level)
			}
			
			for _, shouldnt := range tt.shouldntLog {
				assert.NotContains(t, output, shouldnt+" message", "should not log %s at level %v", shouldnt, tt.level)
			}
		})
	}
}