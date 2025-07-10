package errors

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestSlurmError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *SlurmError
		expected string
	}{
		{
			name: "error with details",
			err: &SlurmError{
				Code:    ErrorCodeNetworkTimeout,
				Message: "Network operation timed out",
				Details: "Connection to server failed after 30s",
			},
			expected: "[NETWORK_TIMEOUT] Network operation timed out: Connection to server failed after 30s",
		},
		{
			name: "error without details",
			err: &SlurmError{
				Code:    ErrorCodeInvalidCredentials,
				Message: "Authentication failed",
			},
			expected: "[INVALID_CREDENTIALS] Authentication failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("SlurmError.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSlurmError_Unwrap(t *testing.T) {
	originalErr := errors.New("original error")
	slurmErr := NewSlurmErrorWithCause(ErrorCodeNetworkTimeout, "timeout", originalErr)

	if unwrapped := slurmErr.Unwrap(); unwrapped != originalErr {
		t.Errorf("SlurmError.Unwrap() = %v, want %v", unwrapped, originalErr)
	}
}

func TestSlurmError_Is(t *testing.T) {
	err1 := NewSlurmError(ErrorCodeNetworkTimeout, "timeout 1")
	err2 := NewSlurmError(ErrorCodeNetworkTimeout, "timeout 2")
	err3 := NewSlurmError(ErrorCodeInvalidCredentials, "auth error")

	// Same error code should match
	if !err1.Is(err2) {
		t.Error("Expected err1.Is(err2) to be true for same error codes")
	}

	// Different error codes should not match
	if err1.Is(err3) {
		t.Error("Expected err1.Is(err3) to be false for different error codes")
	}

	// Non-SlurmError should not match
	if err1.Is(errors.New("regular error")) {
		t.Error("Expected err1.Is(regular error) to be false")
	}
}

func TestSlurmError_IsRetryable(t *testing.T) {
	tests := []struct {
		name      string
		code      ErrorCode
		retryable bool
	}{
		{"network timeout", ErrorCodeNetworkTimeout, true},
		{"connection refused", ErrorCodeConnectionRefused, true},
		{"server internal", ErrorCodeServerInternal, true},
		{"rate limited", ErrorCodeRateLimited, true},
		{"invalid credentials", ErrorCodeInvalidCredentials, false},
		{"validation failed", ErrorCodeValidationFailed, false},
		{"resource not found", ErrorCodeResourceNotFound, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewSlurmError(tt.code, "test message")
			if got := err.IsRetryable(); got != tt.retryable {
				t.Errorf("SlurmError.IsRetryable() = %v, want %v", got, tt.retryable)
			}
		})
	}
}

func TestSlurmError_IsTemporary(t *testing.T) {
	tests := []struct {
		name      string
		code      ErrorCode
		temporary bool
	}{
		{"network timeout", ErrorCodeNetworkTimeout, true},
		{"server internal", ErrorCodeServerInternal, true},
		{"resource exhausted", ErrorCodeResourceExhausted, true},
		{"rate limited", ErrorCodeRateLimited, true},
		{"invalid credentials", ErrorCodeInvalidCredentials, false},
		{"validation failed", ErrorCodeValidationFailed, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewSlurmError(tt.code, "test message")
			if got := err.IsTemporary(); got != tt.temporary {
				t.Errorf("SlurmError.IsTemporary() = %v, want %v", got, tt.temporary)
			}
		})
	}
}

func TestNewSlurmError(t *testing.T) {
	before := time.Now()
	err := NewSlurmError(ErrorCodeNetworkTimeout, "timeout error")
	after := time.Now()

	if err.Code != ErrorCodeNetworkTimeout {
		t.Errorf("Expected code %v, got %v", ErrorCodeNetworkTimeout, err.Code)
	}

	if err.Message != "timeout error" {
		t.Errorf("Expected message 'timeout error', got %v", err.Message)
	}

	if err.Category != CategoryNetwork {
		t.Errorf("Expected category %v, got %v", CategoryNetwork, err.Category)
	}

	if !err.Retryable {
		t.Error("Expected retryable to be true for network timeout")
	}

	if err.Timestamp.Before(before) || err.Timestamp.After(after) {
		t.Error("Timestamp should be set to current time")
	}
}

func TestNewSlurmErrorWithCause(t *testing.T) {
	originalErr := errors.New("original error")
	err := NewSlurmErrorWithCause(ErrorCodeNetworkTimeout, "timeout error", originalErr)

	if err.Cause != originalErr {
		t.Errorf("Expected cause %v, got %v", originalErr, err.Cause)
	}

	if err.Unwrap() != originalErr {
		t.Errorf("Expected Unwrap() to return %v, got %v", originalErr, err.Unwrap())
	}
}

func TestNetworkError(t *testing.T) {
	originalErr := errors.New("connection failed")
	netErr := NewNetworkError(ErrorCodeConnectionRefused, "connection refused", "localhost", 6820, originalErr)

	if netErr.Host != "localhost" {
		t.Errorf("Expected host 'localhost', got %v", netErr.Host)
	}

	if netErr.Port != 6820 {
		t.Errorf("Expected port 6820, got %v", netErr.Port)
	}

	if netErr.Code != ErrorCodeConnectionRefused {
		t.Errorf("Expected code %v, got %v", ErrorCodeConnectionRefused, netErr.Code)
	}

	if netErr.Category != CategoryNetwork {
		t.Errorf("Expected category %v, got %v", CategoryNetwork, netErr.Category)
	}
}

func TestAuthenticationError(t *testing.T) {
	originalErr := errors.New("auth failed")
	authErr := NewAuthenticationError(ErrorCodeInvalidCredentials, "invalid token", "bearer", "jwt", originalErr)

	if authErr.AuthMethod != "bearer" {
		t.Errorf("Expected auth method 'bearer', got %v", authErr.AuthMethod)
	}

	if authErr.TokenType != "jwt" {
		t.Errorf("Expected token type 'jwt', got %v", authErr.TokenType)
	}

	if authErr.Code != ErrorCodeInvalidCredentials {
		t.Errorf("Expected code %v, got %v", ErrorCodeInvalidCredentials, authErr.Code)
	}

	if authErr.Category != CategoryAuthentication {
		t.Errorf("Expected category %v, got %v", CategoryAuthentication, authErr.Category)
	}
}

func TestValidationError(t *testing.T) {
	originalErr := errors.New("validation failed")
	valErr := NewValidationError(ErrorCodeValidationFailed, "invalid field", "name", "", originalErr)

	if valErr.Field != "name" {
		t.Errorf("Expected field 'name', got %v", valErr.Field)
	}

	if valErr.Value != "" {
		t.Errorf("Expected value '', got %v", valErr.Value)
	}

	if valErr.Code != ErrorCodeValidationFailed {
		t.Errorf("Expected code %v, got %v", ErrorCodeValidationFailed, valErr.Code)
	}

	if valErr.Category != CategoryValidation {
		t.Errorf("Expected category %v, got %v", CategoryValidation, valErr.Category)
	}
}

func TestNewSlurmAPIError(t *testing.T) {
	details := []SlurmAPIErrorDetail{
		{
			ErrorNumber: 401,
			ErrorCode:   "SLURM_AUTHENTICATION_ERROR",
			Source:      "auth",
			Description: "Authentication failed",
		},
	}

	apiErr := NewSlurmAPIError(http.StatusUnauthorized, "v0.0.42", details)

	if apiErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, apiErr.StatusCode)
	}

	if apiErr.APIVersion != "v0.0.42" {
		t.Errorf("Expected API version 'v0.0.42', got %v", apiErr.APIVersion)
	}

	if len(apiErr.Errors) != 1 {
		t.Errorf("Expected 1 error detail, got %d", len(apiErr.Errors))
	}

	if apiErr.Code != ErrorCodeInvalidCredentials {
		t.Errorf("Expected code %v, got %v", ErrorCodeInvalidCredentials, apiErr.Code)
	}

	if apiErr.Message != "Authentication failed" {
		t.Errorf("Expected message 'Authentication failed', got %v", apiErr.Message)
	}
}

func TestGetErrorCategory(t *testing.T) {
	tests := []struct {
		code     ErrorCode
		category ErrorCategory
	}{
		{ErrorCodeNetworkTimeout, CategoryNetwork},
		{ErrorCodeConnectionRefused, CategoryNetwork},
		{ErrorCodeInvalidCredentials, CategoryAuthentication},
		{ErrorCodeTokenExpired, CategoryAuthentication},
		{ErrorCodeInvalidRequest, CategoryValidation},
		{ErrorCodeValidationFailed, CategoryValidation},
		{ErrorCodeResourceNotFound, CategoryResource},
		{ErrorCodeConflict, CategoryResource},
		{ErrorCodeServerInternal, CategoryServer},
		{ErrorCodeSlurmDaemonDown, CategoryServer},
		{ErrorCodeClientNotInitialized, CategoryClient},
		{ErrorCodeVersionMismatch, CategoryClient},
		{ErrorCodeContextCanceled, CategoryContext},
		{ErrorCodeDeadlineExceeded, CategoryContext},
		{ErrorCodeUnknown, CategoryUnknown},
	}

	for _, tt := range tests {
		t.Run(string(tt.code), func(t *testing.T) {
			if got := getErrorCategory(tt.code); got != tt.category {
				t.Errorf("getErrorCategory(%v) = %v, want %v", tt.code, got, tt.category)
			}
		})
	}
}

func TestMapHTTPStatusToErrorCode(t *testing.T) {
	tests := []struct {
		status int
		code   ErrorCode
	}{
		{http.StatusBadRequest, ErrorCodeInvalidRequest},
		{http.StatusUnauthorized, ErrorCodeUnauthorized},
		{http.StatusForbidden, ErrorCodePermissionDenied},
		{http.StatusNotFound, ErrorCodeResourceNotFound},
		{http.StatusConflict, ErrorCodeConflict},
		{http.StatusUnprocessableEntity, ErrorCodeValidationFailed},
		{http.StatusTooManyRequests, ErrorCodeRateLimited},
		{http.StatusInternalServerError, ErrorCodeServerInternal},
		{http.StatusBadGateway, ErrorCodeSlurmDaemonDown},
		{http.StatusServiceUnavailable, ErrorCodeSlurmDaemonDown},
		{http.StatusGatewayTimeout, ErrorCodeSlurmDaemonDown},
		{999, ErrorCodeUnknown}, // Unknown status code
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("status_%d", tt.status), func(t *testing.T) {
			if got := mapHTTPStatusToErrorCode(tt.status); got != tt.code {
				t.Errorf("mapHTTPStatusToErrorCode(%d) = %v, want %v", tt.status, got, tt.code)
			}
		})
	}
}

func TestMapSlurmErrorCodeToClientCode(t *testing.T) {
	tests := []struct {
		slurmCode  string
		httpStatus int
		expected   ErrorCode
	}{
		{"SLURM_NO_CHANGE_IN_DATA", 404, ErrorCodeResourceNotFound},
		{"SLURM_PROTOCOL_VERSION_ERROR", 400, ErrorCodeVersionMismatch},
		{"SLURM_AUTHENTICATION_ERROR", 401, ErrorCodeInvalidCredentials},
		{"SLURM_ACCESS_DENIED", 403, ErrorCodePermissionDenied},
		{"SLURM_INVALID_JOB_ID", 404, ErrorCodeResourceNotFound},
		{"SLURM_JOB_ALREADY_COMPLETE", 409, ErrorCodeConflict},
		{"UNKNOWN_SLURM_ERROR", 500, ErrorCodeServerInternal}, // Falls back to HTTP status
	}

	for _, tt := range tests {
		t.Run(tt.slurmCode, func(t *testing.T) {
			if got := mapSlurmErrorCodeToClientCode(tt.slurmCode, tt.httpStatus); got != tt.expected {
				t.Errorf("mapSlurmErrorCodeToClientCode(%s, %d) = %v, want %v", tt.slurmCode, tt.httpStatus, got, tt.expected)
			}
		})
	}
}