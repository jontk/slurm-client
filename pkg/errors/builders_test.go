package errors

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strings"
	"syscall"
	"testing"
)

func TestWrapError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected ErrorCode
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: "",
		},
		{
			name:     "context canceled",
			err:      context.Canceled,
			expected: ErrorCodeContextCanceled,
		},
		{
			name:     "context deadline exceeded",
			err:      context.DeadlineExceeded,
			expected: ErrorCodeDeadlineExceeded,
		},
		{
			name:     "existing SlurmError",
			err:      NewSlurmError(ErrorCodeNetworkTimeout, "timeout"),
			expected: ErrorCodeNetworkTimeout,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WrapError(tt.err)

			if tt.err == nil {
				if result != nil {
					t.Errorf("Expected nil for nil error, got %v", result)
				}
				return
			}

			if result == nil {
				t.Fatal("Expected non-nil error result")
			}

			if result.Code != tt.expected {
				t.Errorf("Expected error code %v, got %v", tt.expected, result.Code)
			}
		})
	}
}

func TestWrapHTTPError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       []byte
		apiVersion string
		expected   ErrorCode
	}{
		{
			name:       "400 bad request",
			statusCode: 400,
			body:       []byte("Bad request"),
			apiVersion: "v0.0.42",
			expected:   ErrorCodeInvalidRequest,
		},
		{
			name:       "401 unauthorized",
			statusCode: 401,
			body:       []byte("Unauthorized"),
			apiVersion: "v0.0.42",
			expected:   ErrorCodeUnauthorized,
		},
		{
			name:       "404 not found",
			statusCode: 404,
			body:       []byte("Not found"),
			apiVersion: "v0.0.42",
			expected:   ErrorCodeResourceNotFound,
		},
		{
			name:       "500 internal server error",
			statusCode: 500,
			body:       []byte("Internal server error"),
			apiVersion: "v0.0.42",
			expected:   ErrorCodeServerInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WrapHTTPError(tt.statusCode, tt.body, tt.apiVersion)

			if result.Code != tt.expected {
				t.Errorf("Expected error code %v, got %v", tt.expected, result.Code)
			}

			if result.StatusCode != tt.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.statusCode, result.StatusCode)
			}

			if result.APIVersion != tt.apiVersion {
				t.Errorf("Expected API version %s, got %s", tt.apiVersion, result.APIVersion)
			}

			// Note: Details might be set by parseSlurmAPIError or HTTP fallback
			// For plain text bodies in tests, we expect them to be included
			if len(tt.body) > 0 && len(tt.body) < 1000 && result.Details == "" {
				t.Errorf("Expected details to be set for body %s, but got empty details", string(tt.body))
			}
		})
	}
}

func TestClassifyNetworkError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected ErrorCode
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: "",
		},
		{
			name:     "connection refused",
			err:      &net.OpError{Op: "dial", Err: syscall.ECONNREFUSED},
			expected: ErrorCodeConnectionRefused,
		},
		{
			name:     "timeout error",
			err:      &timeoutError{},
			expected: ErrorCodeNetworkTimeout,
		},
		{
			name:     "temporary error",
			err:      &temporaryError{},
			expected: ErrorCodeConnectionRefused,
		},
		{
			name:     "DNS error",
			err:      &net.OpError{Op: "dial", Err: &net.DNSError{Name: "example.com"}},
			expected: ErrorCodeDNSResolution,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := classifyNetworkError(tt.err)

			if tt.expected == "" {
				if result != nil {
					t.Errorf("Expected nil for %s, got %v", tt.name, result)
				}
				return
			}

			if result == nil {
				t.Fatalf("Expected non-nil error for %s", tt.name)
			}

			if result.Code != tt.expected {
				t.Errorf("Expected error code %v for %s, got %v", tt.expected, tt.name, result.Code)
			}
		})
	}
}

func TestClassifyURLError(t *testing.T) {
	tests := []struct {
		name     string
		urlErr   *url.Error
		expected ErrorCode
	}{
		{
			name: "URL with connection refused",
			urlErr: &url.Error{
				Op:  "Get",
				URL: "https://localhost:6820/api",
				Err: syscall.ECONNREFUSED,
			},
			expected: ErrorCodeConnectionRefused,
		},
		{
			name: "URL with timeout",
			urlErr: &url.Error{
				Op:  "Get",
				URL: "https://localhost:6820/api",
				Err: &timeoutError{},
			},
			expected: ErrorCodeNetworkTimeout,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := classifyURLError(tt.urlErr)

			if result.Code != tt.expected {
				t.Errorf("Expected error code %v, got %v", tt.expected, result.Code)
			}
		})
	}
}

func TestNewClientError(t *testing.T) {
	err := NewClientError(ErrorCodeClientNotInitialized, "Client not ready", "detail1", "detail2")

	if err.Code != ErrorCodeClientNotInitialized {
		t.Errorf("Expected code %v, got %v", ErrorCodeClientNotInitialized, err.Code)
	}

	if err.Message != "Client not ready" {
		t.Errorf("Expected message 'Client not ready', got %v", err.Message)
	}

	expectedDetails := "detail1; detail2"
	if err.Details != expectedDetails {
		t.Errorf("Expected details %s, got %s", expectedDetails, err.Details)
	}

	if err.Category != CategoryClient {
		t.Errorf("Expected category %v, got %v", CategoryClient, err.Category)
	}
}

func TestNewAuthError(t *testing.T) {
	tests := []struct {
		name       string
		cause      error
		authMethod string
		tokenType  string
		expected   ErrorCode
	}{
		{
			name:       "401 error",
			cause:      fmt.Errorf("401 Unauthorized"),
			authMethod: "bearer",
			tokenType:  "jwt",
			expected:   ErrorCodeInvalidCredentials,
		},
		{
			name:       "403 error",
			cause:      fmt.Errorf("403 Forbidden"),
			authMethod: "basic",
			tokenType:  "",
			expected:   ErrorCodePermissionDenied,
		},
		{
			name:       "expired token",
			cause:      fmt.Errorf("token expired"),
			authMethod: "bearer",
			tokenType:  "jwt",
			expected:   ErrorCodeTokenExpired,
		},
		{
			name:       "generic auth error",
			cause:      fmt.Errorf("auth failed"),
			authMethod: "api_key",
			tokenType:  "",
			expected:   ErrorCodeUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewAuthError(tt.authMethod, tt.tokenType, tt.cause)

			if result.Code != tt.expected {
				t.Errorf("Expected error code %v, got %v", tt.expected, result.Code)
			}

			if result.AuthMethod != tt.authMethod {
				t.Errorf("Expected auth method %s, got %s", tt.authMethod, result.AuthMethod)
			}

			if result.TokenType != tt.tokenType {
				t.Errorf("Expected token type %s, got %s", tt.tokenType, result.TokenType)
			}

			if result.Cause != tt.cause {
				t.Errorf("Expected cause %v, got %v", tt.cause, result.Cause)
			}
		})
	}
}

func TestNewValidationErrorf(t *testing.T) {
	result := NewValidationErrorf("name", "", "field %s cannot be empty", "name")

	if result.Code != ErrorCodeValidationFailed {
		t.Errorf("Expected code %v, got %v", ErrorCodeValidationFailed, result.Code)
	}

	expectedMessage := "field name cannot be empty"
	if result.Message != expectedMessage {
		t.Errorf("Expected message %s, got %s", expectedMessage, result.Message)
	}

	if result.Field != "name" {
		t.Errorf("Expected field 'name', got %s", result.Field)
	}

	if result.Value != "" {
		t.Errorf("Expected value '', got %v", result.Value)
	}
}

func TestNewJobError(t *testing.T) {
	tests := []struct {
		name      string
		jobID     uint32
		operation string
		cause     error
		expected  ErrorCode
	}{
		{
			name:      "job not found",
			jobID:     12345,
			operation: "get",
			cause:     fmt.Errorf("job not found"),
			expected:  ErrorCodeResourceNotFound,
		},
		{
			name:      "job already complete",
			jobID:     12346,
			operation: "cancel",
			cause:     fmt.Errorf("job already complete"),
			expected:  ErrorCodeConflict,
		},
		{
			name:      "permission denied",
			jobID:     12347,
			operation: "cancel",
			cause:     fmt.Errorf("permission denied"),
			expected:  ErrorCodePermissionDenied,
		},
		{
			name:      "queue full",
			jobID:     12348,
			operation: "submit",
			cause:     fmt.Errorf("queue full"),
			expected:  ErrorCodeJobQueueFull,
		},
		{
			name:      "generic error",
			jobID:     12349,
			operation: "submit",
			cause:     fmt.Errorf("generic error"),
			expected:  ErrorCodeServerInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewJobError(tt.jobID, tt.operation, tt.cause)

			if result.Code != tt.expected {
				t.Errorf("Expected error code %v, got %v", tt.expected, result.Code)
			}

			if !strings.Contains(result.Details, fmt.Sprintf("Job ID: %d", tt.jobID)) {
				t.Errorf("Expected details to contain job ID %d", tt.jobID)
			}

			if !strings.Contains(result.Details, fmt.Sprintf("Operation: %s", tt.operation)) {
				t.Errorf("Expected details to contain operation %s", tt.operation)
			}

			if result.Cause != tt.cause {
				t.Errorf("Expected cause %v, got %v", tt.cause, result.Cause)
			}
		})
	}
}

func TestNewNodeError(t *testing.T) {
	nodeNames := []string{"compute-01", "compute-02"}
	operation := "update"
	cause := fmt.Errorf("nodes not available")

	result := NewNodeError(nodeNames, operation, cause)

	if result.Code != ErrorCodeResourceExhausted {
		t.Errorf("Expected error code %v, got %v", ErrorCodeResourceExhausted, result.Code)
	}

	if !strings.Contains(result.Details, "compute-01") {
		t.Error("Expected details to contain node names")
	}

	if !strings.Contains(result.Details, operation) {
		t.Errorf("Expected details to contain operation %s", operation)
	}

	if result.Cause != cause {
		t.Errorf("Expected cause %v, got %v", cause, result.Cause)
	}
}

func TestNewPartitionError(t *testing.T) {
	partitionName := "compute"
	operation := "access"
	cause := fmt.Errorf("partition unavailable")

	result := NewPartitionError(partitionName, operation, cause)

	if result.Code != ErrorCodePartitionUnavailable {
		t.Errorf("Expected error code %v, got %v", ErrorCodePartitionUnavailable, result.Code)
	}

	if !strings.Contains(result.Message, partitionName) {
		t.Errorf("Expected message to contain partition name %s", partitionName)
	}

	if !strings.Contains(result.Details, partitionName) {
		t.Errorf("Expected details to contain partition name %s", partitionName)
	}

	if result.Cause != cause {
		t.Errorf("Expected cause %v, got %v", cause, result.Cause)
	}
}

func TestIsRetryableError(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		retryable bool
	}{
		{
			name:      "retryable SlurmError",
			err:       NewSlurmError(ErrorCodeNetworkTimeout, "timeout"),
			retryable: true,
		},
		{
			name:      "non-retryable SlurmError",
			err:       NewSlurmError(ErrorCodeInvalidCredentials, "bad auth"),
			retryable: false,
		},
		{
			name:      "timeout string error",
			err:       fmt.Errorf("connection timeout"),
			retryable: true,
		},
		{
			name:      "connection refused string error",
			err:       fmt.Errorf("connection refused"),
			retryable: true,
		},
		{
			name:      "non-retryable string error",
			err:       fmt.Errorf("invalid input"),
			retryable: false,
		},
		{
			name:      "nil error",
			err:       nil,
			retryable: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsRetryableError(tt.err); got != tt.retryable {
				t.Errorf("IsRetryableError() = %v, want %v", got, tt.retryable)
			}
		})
	}
}

func TestIsTemporaryError(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		temporary bool
	}{
		{
			name:      "temporary SlurmError",
			err:       NewSlurmError(ErrorCodeNetworkTimeout, "timeout"),
			temporary: true,
		},
		{
			name:      "non-temporary SlurmError",
			err:       NewSlurmError(ErrorCodeInvalidCredentials, "bad auth"),
			temporary: false,
		},
		{
			name:      "temporary network error",
			err:       &temporaryError{},
			temporary: true,
		},
		{
			name:      "non-temporary error",
			err:       fmt.Errorf("permanent error"),
			temporary: false,
		},
		{
			name:      "nil error",
			err:       nil,
			temporary: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsTemporaryError(tt.err); got != tt.temporary {
				t.Errorf("IsTemporaryError() = %v, want %v", got, tt.temporary)
			}
		})
	}
}

func TestGetErrorCode(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected ErrorCode
	}{
		{
			name:     "SlurmError",
			err:      NewSlurmError(ErrorCodeNetworkTimeout, "timeout"),
			expected: ErrorCodeNetworkTimeout,
		},
		{
			name:     "regular error",
			err:      fmt.Errorf("regular error"),
			expected: ErrorCodeUnknown,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: ErrorCodeUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetErrorCode(tt.err); got != tt.expected {
				t.Errorf("GetErrorCode() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetErrorCategoryFromError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected ErrorCategory
	}{
		{
			name:     "SlurmError",
			err:      NewSlurmError(ErrorCodeNetworkTimeout, "timeout"),
			expected: CategoryNetwork,
		},
		{
			name:     "regular error",
			err:      fmt.Errorf("regular error"),
			expected: CategoryUnknown,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: CategoryUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetErrorCategory(tt.err); got != tt.expected {
				t.Errorf("GetErrorCategory() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// Test helper types
type timeoutError struct{}

func (e *timeoutError) Error() string   { return "timeout" }
func (e *timeoutError) Timeout() bool   { return true }
func (e *timeoutError) Temporary() bool { return false }

type temporaryError struct{}

func (e *temporaryError) Error() string   { return "temporary" }
func (e *temporaryError) Timeout() bool   { return false }
func (e *temporaryError) Temporary() bool { return true }
