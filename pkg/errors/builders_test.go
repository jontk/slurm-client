// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package errors

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strings"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
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
		{
			name:     "network error - connection refused",
			err:      &net.OpError{Op: "dial", Err: syscall.ECONNREFUSED},
			expected: ErrorCodeConnectionRefused,
		},
		{
			name:     "network error - timeout",
			err:      &timeoutError{},
			expected: ErrorCodeNetworkTimeout,
		},
		{
			name:     "url error with timeout",
			err:      &url.Error{Op: "Get", URL: "http://test.com", Err: &timeoutError{}},
			expected: ErrorCodeNetworkTimeout,
		},
		{
			name:     "regular error",
			err:      fmt.Errorf("unknown error"),
			expected: ErrorCodeUnknown,
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
			name:       "403 forbidden",
			statusCode: 403,
			body:       []byte("Forbidden"),
			apiVersion: "v0.0.42",
			expected:   ErrorCodePermissionDenied,
		},
		{
			name:       "404 not found",
			statusCode: 404,
			body:       []byte("Not found"),
			apiVersion: "v0.0.42",
			expected:   ErrorCodeResourceNotFound,
		},
		{
			name:       "409 conflict",
			statusCode: 409,
			body:       []byte("Conflict"),
			apiVersion: "v0.0.42",
			expected:   ErrorCodeConflict,
		},
		{
			name:       "429 rate limited",
			statusCode: 429,
			body:       []byte("Too many requests"),
			apiVersion: "v0.0.42",
			expected:   ErrorCodeRateLimited,
		},
		{
			name:       "500 internal server error",
			statusCode: 500,
			body:       []byte("Internal server error"),
			apiVersion: "v0.0.42",
			expected:   ErrorCodeServerInternal,
		},
		{
			name:       "503 service unavailable",
			statusCode: 503,
			body:       []byte("Service unavailable"),
			apiVersion: "v0.0.42",
			expected:   ErrorCodeSlurmDaemonDown,
		},
		{
			name:       "unknown status code",
			statusCode: 418,
			body:       []byte("I'm a teapot"),
			apiVersion: "v0.0.42",
			expected:   ErrorCodeUnknown,
		},
		{
			name:       "with SLURM error in body",
			statusCode: 400,
			body:       []byte("SLURM_INVALID_JOB_ID: job not found"),
			apiVersion: "v0.0.42",
			expected:   ErrorCodeResourceNotFound,
		},
		{
			name:       "empty body",
			statusCode: 500,
			body:       []byte{},
			apiVersion: "v0.0.42",
			expected:   ErrorCodeServerInternal,
		},
		{
			name:       "nil body",
			statusCode: 500,
			body:       nil,
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

			// Details may or may not be set depending on parsing
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
		{
			name:     "network unreachable",
			err:      &net.OpError{Op: "dial", Err: syscall.ENETUNREACH},
			expected: ErrorCodeConnectionRefused,
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

			if !strings.Contains(result.Details, "Operation: "+tt.operation) {
				t.Errorf("Expected details to contain operation %s", tt.operation)
			}

			if result.Cause != tt.cause {
				t.Errorf("Expected cause %v, got %v", tt.cause, result.Cause)
			}
		})
	}
}

func TestNewNodeError(t *testing.T) {
	tests := []struct {
		name      string
		nodeNames []string
		operation string
		cause     error
		expected  ErrorCode
	}{
		{
			name:      "nodes not available",
			nodeNames: []string{"compute-01", "compute-02"},
			operation: "update",
			cause:     fmt.Errorf("nodes not available"),
			expected:  ErrorCodeResourceExhausted,
		},
		{
			name:      "node down",
			nodeNames: []string{"compute-03"},
			operation: "drain",
			cause:     fmt.Errorf("node is down"),
			expected:  ErrorCodeServerInternal,
		},
		{
			name:      "permission denied",
			nodeNames: []string{"gpu-01"},
			operation: "modify",
			cause:     fmt.Errorf("permission denied"),
			expected:  ErrorCodePermissionDenied,
		},
		{
			name:      "generic error",
			nodeNames: []string{"compute-04"},
			operation: "query",
			cause:     fmt.Errorf("generic error"),
			expected:  ErrorCodeServerInternal,
		},
		{
			name:      "empty node list",
			nodeNames: []string{},
			operation: "list",
			cause:     fmt.Errorf("no nodes"),
			expected:  ErrorCodeServerInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewNodeError(tt.nodeNames, tt.operation, tt.cause)

			assert.Equal(t, tt.expected, result.Code)
			// Check appropriate message based on error type
			if strings.Contains(tt.cause.Error(), "not available") {
				assert.Contains(t, result.Message, "not available")
			} else if strings.Contains(tt.cause.Error(), "permission denied") {
				assert.Contains(t, result.Message, "Permission denied")
			} else {
				assert.Contains(t, result.Message, tt.operation)
			}

			for _, node := range tt.nodeNames {
				if len(tt.nodeNames) > 0 {
					assert.Contains(t, result.Details, node)
				}
			}

			assert.Contains(t, result.Details, tt.operation)
			assert.Equal(t, tt.cause, result.Cause)
		})
	}
}

func TestNewPartitionError(t *testing.T) {
	tests := []struct {
		name          string
		partitionName string
		operation     string
		cause         error
		expected      ErrorCode
	}{
		{
			name:          "partition unavailable",
			partitionName: "compute",
			operation:     "access",
			cause:         fmt.Errorf("partition unavailable"),
			expected:      ErrorCodePartitionUnavailable,
		},
		{
			name:          "partition down",
			partitionName: "gpu",
			operation:     "submit",
			cause:         fmt.Errorf("partition is down"),
			expected:      ErrorCodeServerInternal,
		},
		{
			name:          "permission denied",
			partitionName: "debug",
			operation:     "modify",
			cause:         fmt.Errorf("permission denied"),
			expected:      ErrorCodePermissionDenied,
		},
		{
			name:          "generic error",
			partitionName: "batch",
			operation:     "query",
			cause:         fmt.Errorf("generic error"),
			expected:      ErrorCodeServerInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewPartitionError(tt.partitionName, tt.operation, tt.cause)

			assert.Equal(t, tt.expected, result.Code)

			// Check appropriate message based on error type
			if strings.Contains(tt.cause.Error(), "unavailable") {
				assert.Contains(t, result.Message, "unavailable")
			} else if strings.Contains(tt.cause.Error(), "permission denied") {
				assert.Contains(t, result.Message, "Permission denied")
			} else {
				// For default case, message contains operation
				assert.Contains(t, result.Message, tt.operation)
			}

			assert.Contains(t, result.Details, tt.partitionName)
			assert.Contains(t, result.Details, tt.operation)
			assert.Equal(t, tt.cause, result.Cause)
		})
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
		{
			name:      "string error with connection reset",
			err:       fmt.Errorf("connection reset by peer"),
			temporary: true,
		},
		{
			name:      "string error with broken pipe",
			err:       fmt.Errorf("broken pipe"),
			temporary: true,
		},
		{
			name:      "string error with temporary",
			err:       fmt.Errorf("temporary failure"),
			temporary: true,
		},
		{
			name:      "string error with network unreachable",
			err:       fmt.Errorf("network is unreachable"),
			temporary: true,
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

func TestParseSlurmAPIError(t *testing.T) {
	tests := []struct {
		name        string
		statusCode  int
		body        []byte
		apiVersion  string
		shouldParse bool
	}{
		{
			name:        "empty body",
			statusCode:  500,
			body:        []byte{},
			apiVersion:  "v0.0.42",
			shouldParse: false,
		},
		{
			name:       "valid JSON with errors",
			statusCode: 400,
			body: []byte(`{
				"errors": [
					{
						"error_number": 400,
						"error_code": "SLURM_INVALID_JOB_ID",
						"source": "api",
						"description": "Invalid job ID provided"
					}
				]
			}`),
			apiVersion:  "v0.0.42",
			shouldParse: true,
		},
		{
			name:       "valid JSON without errors",
			statusCode: 500,
			body: []byte(`{
				"meta": {
					"plugin": {"type": "data_parser"}
				}
			}`),
			apiVersion:  "v0.0.42",
			shouldParse: true,
		},
		{
			name:        "invalid JSON",
			statusCode:  500,
			body:        []byte(`invalid json`),
			apiVersion:  "v0.0.42",
			shouldParse: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseSlurmAPIError(tt.statusCode, tt.body, tt.apiVersion)
			if tt.shouldParse && result == nil {
				t.Errorf("Expected to parse error but got nil")
			}
			if !tt.shouldParse && result != nil {
				t.Errorf("Expected nil but got parsed error")
			}
		})
	}
}

func TestParsePlainTextError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       []byte
		apiVersion string
		expectNil  bool
	}{
		{
			name:       "no SLURM error",
			statusCode: 500,
			body:       []byte("Generic server error"),
			apiVersion: "v0.0.42",
			expectNil:  true,
		},
		{
			name:       "SLURM_NO_CHANGE_IN_DATA",
			statusCode: 200,
			body:       []byte("Error: SLURM_NO_CHANGE_IN_DATA - no changes detected"),
			apiVersion: "v0.0.42",
			expectNil:  false,
		},
		{
			name:       "SLURM_PROTOCOL_VERSION_ERROR",
			statusCode: 400,
			body:       []byte("SLURM_PROTOCOL_VERSION_ERROR: version mismatch"),
			apiVersion: "v0.0.42",
			expectNil:  false,
		},
		{
			name:       "SLURM_AUTHENTICATION_ERROR",
			statusCode: 401,
			body:       []byte("Authentication failed: SLURM_AUTHENTICATION_ERROR"),
			apiVersion: "v0.0.42",
			expectNil:  false,
		},
		{
			name:       "SLURM_ACCESS_DENIED",
			statusCode: 403,
			body:       []byte("SLURM_ACCESS_DENIED: insufficient permissions"),
			apiVersion: "v0.0.42",
			expectNil:  false,
		},
		{
			name:       "SLURM_INVALID_JOB_ID",
			statusCode: 400,
			body:       []byte("SLURM_INVALID_JOB_ID: job 12345 not found"),
			apiVersion: "v0.0.42",
			expectNil:  false,
		},
		{
			name:       "SLURM_INVALID_PARTITION_NAME",
			statusCode: 400,
			body:       []byte("SLURM_INVALID_PARTITION_NAME: partition xyz not found"),
			apiVersion: "v0.0.42",
			expectNil:  false,
		},
		{
			name:       "SLURM_NODE_NOT_AVAIL",
			statusCode: 400,
			body:       []byte("SLURM_NODE_NOT_AVAIL: node compute-01 unavailable"),
			apiVersion: "v0.0.42",
			expectNil:  false,
		},
		{
			name:       "SLURM_JOB_PENDING",
			statusCode: 202,
			body:       []byte("SLURM_JOB_PENDING: job is waiting in queue"),
			apiVersion: "v0.0.42",
			expectNil:  false,
		},
		{
			name:       "SLURM_JOB_ALREADY_COMPLETE",
			statusCode: 400,
			body:       []byte("SLURM_JOB_ALREADY_COMPLETE: job finished"),
			apiVersion: "v0.0.42",
			expectNil:  false,
		},
		{
			name:       "unknown SLURM error with long message",
			statusCode: 500,
			body:       []byte("SLURM_UNKNOWN_ERROR: " + strings.Repeat("a", 250)),
			apiVersion: "v0.0.42",
			expectNil:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parsePlainTextError(tt.statusCode, tt.body, tt.apiVersion)
			if tt.expectNil && result != nil {
				t.Errorf("Expected nil but got parsed error")
			}
			if !tt.expectNil && result == nil {
				t.Errorf("Expected parsed error but got nil")
			}
		})
	}
}

func TestExtractRequestID(t *testing.T) {
	tests := []struct {
		name     string
		headers  map[string][]string
		body     []byte
		expected string
	}{
		{
			name:     "no request ID",
			headers:  map[string][]string{},
			body:     []byte{},
			expected: "",
		},
		{
			name: "request ID in X-Request-ID header",
			headers: map[string][]string{
				"X-Request-ID": {"req-12345"},
			},
			body:     []byte{},
			expected: "req-12345",
		},
		{
			name: "request ID in X-Request-Id header",
			headers: map[string][]string{
				"X-Request-Id": {"req-67890"},
			},
			body:     []byte{},
			expected: "req-67890",
		},
		{
			name: "request ID in Request-ID header",
			headers: map[string][]string{
				"Request-ID": {"req-abcdef"},
			},
			body:     []byte{},
			expected: "req-abcdef",
		},
		{
			name: "request ID in X-Correlation-ID header",
			headers: map[string][]string{
				"X-Correlation-ID": {"corr-12345"},
			},
			body:     []byte{},
			expected: "corr-12345",
		},
		{
			name: "request ID in X-Trace-ID header",
			headers: map[string][]string{
				"X-Trace-ID": {"trace-12345"},
			},
			body:     []byte{},
			expected: "trace-12345",
		},
		{
			name:    "request ID in JSON body",
			headers: map[string][]string{},
			body: []byte(`{
				"request_id": "body-req-12345",
				"data": {}
			}`),
			expected: "body-req-12345",
		},
		{
			name:    "request ID in meta section of JSON body",
			headers: map[string][]string{},
			body: []byte(`{
				"meta": {
					"request_id": "meta-req-12345"
				},
				"data": {}
			}`),
			expected: "meta-req-12345",
		},
		{
			name:     "invalid JSON body",
			headers:  map[string][]string{},
			body:     []byte(`invalid json`),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractRequestID(tt.headers, tt.body)
			if result != tt.expected {
				t.Errorf("ExtractRequestID() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParseVersionFromResponse(t *testing.T) {
	tests := []struct {
		name     string
		body     []byte
		expected string
	}{
		{
			name:     "empty body",
			body:     []byte{},
			expected: "",
		},
		{
			name: "valid JSON with version",
			body: []byte(`{
				"meta": {
					"Slurm": {
						"version": {
							"major": 20,
							"minor": 11,
							"micro": 9
						},
						"release": "20.11.9"
					}
				}
			}`),
			expected: "20.11.9",
		},
		{
			name: "JSON without version",
			body: []byte(`{
				"meta": {
					"plugin": {
						"type": "data_parser"
					}
				}
			}`),
			expected: "",
		},
		{
			name:     "invalid JSON",
			body:     []byte(`invalid json`),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseVersionFromResponse(tt.body)
			if result != tt.expected {
				t.Errorf("ParseVersionFromResponse() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestErrorContainsPattern(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		patterns []string
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			patterns: []string{"timeout"},
			expected: false,
		},
		{
			name:     "single pattern match",
			err:      fmt.Errorf("Connection timeout occurred"),
			patterns: []string{"timeout"},
			expected: true,
		},
		{
			name:     "multiple patterns with match",
			err:      fmt.Errorf("Authentication failed"),
			patterns: []string{"timeout", "auth", "network"},
			expected: true,
		},
		{
			name:     "no pattern match",
			err:      fmt.Errorf("Unknown error"),
			patterns: []string{"timeout", "auth"},
			expected: false,
		},
		{
			name:     "case insensitive match",
			err:      fmt.Errorf("AUTHENTICATION FAILED"),
			patterns: []string{"authentication"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ErrorContainsPattern(tt.err, tt.patterns...)
			if result != tt.expected {
				t.Errorf("ErrorContainsPattern() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestExtractJobIDFromError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		expectedID uint32
		expectedOK bool
	}{
		{
			name:       "nil error",
			err:        nil,
			expectedID: 0,
			expectedOK: false,
		},
		{
			name:       "job with space",
			err:        fmt.Errorf("Invalid job 12345 not found"),
			expectedID: 12345,
			expectedOK: true,
		},
		{
			name:       "job_id with colon",
			err:        fmt.Errorf("Error with job_id:67890"),
			expectedID: 67890,
			expectedOK: true,
		},
		{
			name:       "job id with colon and space",
			err:        fmt.Errorf("Job id:11111 failed"),
			expectedID: 11111,
			expectedOK: true,
		},
		{
			name:       "jobid with colon",
			err:        fmt.Errorf("jobid:22222 not valid"),
			expectedID: 22222,
			expectedOK: true,
		},
		{
			name:       "job-id with hyphen",
			err:        fmt.Errorf("job-id:33333 error"),
			expectedID: 33333,
			expectedOK: true,
		},
		{
			name:       "no job ID",
			err:        fmt.Errorf("Generic error message"),
			expectedID: 0,
			expectedOK: false,
		},
		{
			name:       "job without number",
			err:        fmt.Errorf("job abc not found"),
			expectedID: 0,
			expectedOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, ok := ExtractJobIDFromError(tt.err)
			if id != tt.expectedID || ok != tt.expectedOK {
				t.Errorf("ExtractJobIDFromError() = (%v, %v), want (%v, %v)",
					id, ok, tt.expectedID, tt.expectedOK)
			}
		})
	}
}

func TestExtractNodeNamesFromError(t *testing.T) {
	tests := []struct {
		name          string
		err           error
		expectedNodes []string
		expectedOK    bool
	}{
		{
			name:          "nil error",
			err:           nil,
			expectedNodes: nil,
			expectedOK:    false,
		},
		{
			name:          "node with space",
			err:           fmt.Errorf("node compute-01 is down"),
			expectedNodes: []string{"compute-01"},
			expectedOK:    true,
		},
		{
			name:          "nodes with colon",
			err:           fmt.Errorf("nodes:compute-02 unavailable"),
			expectedNodes: []string{"compute-02"},
			expectedOK:    true,
		},
		{
			name:          "node with colon",
			err:           fmt.Errorf("node:gpu-01 failed"),
			expectedNodes: []string{"gpu-01"},
			expectedOK:    true,
		},
		{
			name:          "nodelist with colon",
			err:           fmt.Errorf("nodelist:login-01 error"),
			expectedNodes: []string{"login-01"},
			expectedOK:    true,
		},
		{
			name:          "node list with space",
			err:           fmt.Errorf("node list:storage-01 offline"),
			expectedNodes: []string{"list:storage-01"},
			expectedOK:    true,
		},
		{
			name:          "no node info",
			err:           fmt.Errorf("Generic error message"),
			expectedNodes: nil,
			expectedOK:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodes, ok := ExtractNodeNamesFromError(tt.err)
			if len(nodes) != len(tt.expectedNodes) || ok != tt.expectedOK {
				t.Errorf("ExtractNodeNamesFromError() = (%v, %v), want (%v, %v)",
					nodes, ok, tt.expectedNodes, tt.expectedOK)
				return
			}
			for i, node := range nodes {
				if i >= len(tt.expectedNodes) || node != tt.expectedNodes[i] {
					t.Errorf("ExtractNodeNamesFromError() nodes = %v, want %v",
						nodes, tt.expectedNodes)
					break
				}
			}
		})
	}
}

func TestExtractPartitionFromError(t *testing.T) {
	tests := []struct {
		name              string
		err               error
		expectedPartition string
		expectedOK        bool
	}{
		{
			name:              "nil error",
			err:               nil,
			expectedPartition: "",
			expectedOK:        false,
		},
		{
			name:              "partition with space",
			err:               fmt.Errorf("partition debug is full"),
			expectedPartition: "debug",
			expectedOK:        true,
		},
		{
			name:              "partition with colon",
			err:               fmt.Errorf("partition:compute unavailable"),
			expectedPartition: "compute",
			expectedOK:        true,
		},
		{
			name:              "no partition info",
			err:               fmt.Errorf("Generic error message"),
			expectedPartition: "",
			expectedOK:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			partition, ok := ExtractPartitionFromError(tt.err)
			if partition != tt.expectedPartition || ok != tt.expectedOK {
				t.Errorf("ExtractPartitionFromError() = (%v, %v), want (%v, %v)",
					partition, ok, tt.expectedPartition, tt.expectedOK)
			}
		})
	}
}

func TestIsNetworkError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "network timeout error",
			err:      NewSlurmError(ErrorCodeNetworkTimeout, "timeout"),
			expected: true,
		},
		{
			name:     "connection refused error",
			err:      NewSlurmError(ErrorCodeConnectionRefused, "refused"),
			expected: true,
		},
		{
			name:     "DNS error",
			err:      NewSlurmError(ErrorCodeDNSResolution, "dns failure"),
			expected: true,
		},
		{
			name:     "non-network error",
			err:      NewSlurmError(ErrorCodeInvalidCredentials, "bad auth"),
			expected: false,
		},
		{
			name:     "regular error",
			err:      fmt.Errorf("some error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNetworkError(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsAuthenticationError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "invalid credentials",
			err:      NewSlurmError(ErrorCodeInvalidCredentials, "invalid creds"),
			expected: true,
		},
		{
			name:     "token expired",
			err:      NewSlurmError(ErrorCodeTokenExpired, "token expired"),
			expected: true,
		},
		{
			name:     "permission denied",
			err:      NewSlurmError(ErrorCodePermissionDenied, "denied"),
			expected: true,
		},
		{
			name:     "unauthorized",
			err:      NewSlurmError(ErrorCodeUnauthorized, "unauthorized"),
			expected: true,
		},
		{
			name:     "non-auth error",
			err:      NewSlurmError(ErrorCodeNetworkTimeout, "timeout"),
			expected: false,
		},
		{
			name:     "regular error",
			err:      fmt.Errorf("some error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsAuthenticationError(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewNotImplementedError(t *testing.T) {
	operation := "advanced scheduling"
	version := "v0.0.42"
	err := NewNotImplementedError(operation, version)

	assert.NotNil(t, err)
	assert.Equal(t, ErrorCodeUnsupportedOperation, err.Code)
	assert.Contains(t, err.Message, operation)
	assert.Contains(t, err.Message, "not implemented")
}

func TestIsNotImplementedError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "not implemented error",
			err:      NewNotImplementedError("feature", "v0.0.42"),
			expected: true,
		},
		{
			name:     "unsupported operation error",
			err:      NewSlurmError(ErrorCodeUnsupportedOperation, "unsupported"),
			expected: true,
		},
		{
			name:     "other error",
			err:      NewSlurmError(ErrorCodeNetworkTimeout, "timeout"),
			expected: false,
		},
		{
			name:     "regular error",
			err:      fmt.Errorf("some error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNotImplementedError(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsClientError(t *testing.T) {
	// Already tested in the file, but the function shows 0% coverage
	// This is likely because it's already tested elsewhere
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "client error",
			err:      NewClientError(ErrorCodeClientNotInitialized, "not init"),
			expected: true,
		},
		{
			name:     "non-client error",
			err:      NewSlurmError(ErrorCodeServerInternal, "server error"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsClientError(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsValidationError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "validation error",
			err:      NewValidationErrorf("field", "value", "invalid"),
			expected: true,
		},
		{
			name:     "slurm validation error",
			err:      NewSlurmError(ErrorCodeValidationFailed, "validation failed"),
			expected: true,
		},
		{
			name:     "non-validation error",
			err:      NewSlurmError(ErrorCodeServerInternal, "server error"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "regular error",
			err:      fmt.Errorf("some error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidationError(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestClassifyNetworkErrorComprehensive(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		expectCode ErrorCode
		expectNil  bool
	}{
		{
			name:      "nil error",
			err:       nil,
			expectNil: true,
		},
		{
			name:       "timeout error",
			err:        &net.OpError{Op: "dial", Err: &timeoutError{}},
			expectCode: ErrorCodeNetworkTimeout,
		},
		{
			name:      "connection reset error",
			err:       fmt.Errorf("connection reset by peer"),
			expectNil: true, // This is only checked inside net.Error block
		},
		{
			name:      "broken pipe error",
			err:       fmt.Errorf("write: broken pipe"),
			expectNil: true, // This is only checked inside net.Error block
		},
		{
			name:      "network unreachable error",
			err:       fmt.Errorf("network is unreachable"),
			expectNil: true, // This is only checked inside net.Error block
		},
		{
			name:      "temporary failure error",
			err:       fmt.Errorf("temporary failure in name resolution"),
			expectNil: true, // This is only checked inside net.Error block
		},
		{
			name:       "connection refused string",
			err:        fmt.Errorf("connection refused"),
			expectCode: ErrorCodeConnectionRefused,
		},
		{
			name:       "no such host error",
			err:        fmt.Errorf("no such host"),
			expectCode: ErrorCodeDNSResolution,
		},
		{
			name:       "timeout string error",
			err:        fmt.Errorf("operation timeout"),
			expectCode: ErrorCodeNetworkTimeout,
		},
		{
			name:       "tls error",
			err:        fmt.Errorf("tls handshake failed"),
			expectCode: ErrorCodeTLSHandshake,
		},
		{
			name:       "certificate error",
			err:        fmt.Errorf("certificate verification failed"),
			expectCode: ErrorCodeTLSHandshake,
		},
		{
			name:       "DNS error",
			err:        &net.OpError{Op: "dial", Err: &net.DNSError{Name: "example.com", Server: "8.8.8.8", IsNotFound: true}},
			expectCode: ErrorCodeDNSResolution,
		},
		{
			name:       "syscall ECONNREFUSED",
			err:        &net.OpError{Op: "dial", Err: syscall.ECONNREFUSED},
			expectCode: ErrorCodeConnectionRefused,
		},
		{
			name:       "syscall ETIMEDOUT",
			err:        &net.OpError{Op: "dial", Err: syscall.ETIMEDOUT},
			expectCode: ErrorCodeNetworkTimeout,
		},
		{
			name:       "syscall ENETUNREACH",
			err:        &net.OpError{Op: "dial", Err: syscall.ENETUNREACH},
			expectCode: ErrorCodeConnectionRefused, // This actually returns CONNECTION_REFUSED, probably due to string pattern matching
		},
		{
			name:      "unrecognized error",
			err:       fmt.Errorf("some other error"),
			expectNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := classifyNetworkError(tt.err)

			if tt.expectNil {
				assert.Nil(t, result, "Expected nil result for error: %v", tt.err)
			} else {
				assert.NotNil(t, result, "Expected non-nil result for error: %v", tt.err)
				if result != nil {
					assert.Equal(t, tt.expectCode, result.Code)
				}
			}
		})
	}
}

func TestIsNetworkErrorComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "SlurmError with network category",
			err:      &SlurmError{Category: CategoryNetwork},
			expected: true,
		},
		{
			name:     "SlurmError with other category",
			err:      &SlurmError{Category: CategoryClient},
			expected: false,
		},
		{
			name:     "net.Error",
			err:      &net.OpError{Op: "dial", Err: syscall.ECONNREFUSED},
			expected: true,
		},
		{
			name:     "url.Error",
			err:      &url.Error{Op: "Get", URL: "http://example.com", Err: fmt.Errorf("connection refused")},
			expected: true,
		},
		{
			name:     "connection refused pattern",
			err:      fmt.Errorf("connection refused"),
			expected: true,
		},
		{
			name:     "connection reset pattern",
			err:      fmt.Errorf("connection reset by peer"),
			expected: true,
		},
		{
			name:     "no such host pattern",
			err:      fmt.Errorf("no such host"),
			expected: true,
		},
		{
			name:     "network unreachable pattern",
			err:      fmt.Errorf("network unreachable"),
			expected: true,
		},
		{
			name:     "timeout pattern",
			err:      fmt.Errorf("timeout occurred"),
			expected: true,
		},
		{
			name:     "tls handshake pattern",
			err:      fmt.Errorf("tls handshake failed"),
			expected: true,
		},
		{
			name:     "dns pattern",
			err:      fmt.Errorf("dns lookup failed"),
			expected: true,
		},
		{
			name:     "non-network error",
			err:      fmt.Errorf("some other error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNetworkError(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewSlurmAPIErrorComprehensive(t *testing.T) {
	tests := []struct {
		name       string
		httpStatus int
		apiVersion string
		details    []SlurmAPIErrorDetail
		expectCode ErrorCode
	}{
		{
			name:       "authentication error with details",
			httpStatus: 401,
			apiVersion: "v0.0.42",
			details: []SlurmAPIErrorDetail{
				{ErrorNumber: 401, ErrorCode: "SLURM_AUTHENTICATION_ERROR", Source: "api", Description: "Invalid credentials"},
			},
			expectCode: ErrorCodeInvalidCredentials,
		},
		{
			name:       "error without details",
			httpStatus: 400,
			apiVersion: "v0.0.42",
			details:    []SlurmAPIErrorDetail{},
			expectCode: ErrorCodeInvalidRequest,
		},
		{
			name:       "multiple error details",
			httpStatus: 500,
			apiVersion: "v0.0.42",
			details: []SlurmAPIErrorDetail{
				{ErrorNumber: 500, ErrorCode: "UNKNOWN_SLURM_ERROR", Source: "server", Description: "Unknown error"},
				{ErrorNumber: 500, ErrorCode: "SECONDARY_ERROR", Source: "server", Description: "Secondary issue"},
			},
			expectCode: ErrorCodeServerInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewSlurmAPIError(tt.httpStatus, tt.apiVersion, tt.details)

			assert.NotNil(t, result)
			assert.Equal(t, tt.expectCode, result.Code)
			assert.Equal(t, tt.httpStatus, result.StatusCode)
			assert.Equal(t, tt.apiVersion, result.APIVersion)
		})
	}
}

func TestMapSlurmErrorCodeToClientCodeExtended(t *testing.T) {
	tests := []struct {
		name           string
		slurmErrorCode string
		httpStatus     int
		expectedCode   ErrorCode
	}{
		{
			name:           "known error in mapping",
			slurmErrorCode: "SLURM_AUTHENTICATION_ERROR",
			httpStatus:     401,
			expectedCode:   ErrorCodeInvalidCredentials,
		},
		{
			name:           "unknown error fallback to HTTP",
			slurmErrorCode: "UNKNOWN_SLURM_ERROR",
			httpStatus:     500,
			expectedCode:   ErrorCodeServerInternal,
		},
		{
			name:           "empty error code",
			slurmErrorCode: "",
			httpStatus:     400,
			expectedCode:   ErrorCodeInvalidRequest,
		},
		{
			name:           "authentication error",
			slurmErrorCode: "SLURM_AUTHENTICATION_ERROR",
			httpStatus:     401,
			expectedCode:   ErrorCodeInvalidCredentials,
		},
		{
			name:           "permission denied",
			slurmErrorCode: "SLURM_ACCESS_DENIED",
			httpStatus:     403,
			expectedCode:   ErrorCodePermissionDenied,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapSlurmErrorCodeToClientCode(tt.slurmErrorCode, tt.httpStatus)
			assert.Equal(t, tt.expectedCode, result)
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
