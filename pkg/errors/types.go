// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package errors

import (
	"fmt"
	"net/http"
	"time"
)

// ErrorCode represents structured error codes for the Slurm client
type ErrorCode string

const (
	// Network and connectivity errors
	ErrorCodeNetworkTimeout    ErrorCode = "NETWORK_TIMEOUT"
	ErrorCodeConnectionRefused ErrorCode = "CONNECTION_REFUSED"
	ErrorCodeDNSResolution     ErrorCode = "DNS_RESOLUTION"
	ErrorCodeTLSHandshake      ErrorCode = "TLS_HANDSHAKE"

	// Authentication and authorization errors
	ErrorCodeInvalidCredentials ErrorCode = "INVALID_CREDENTIALS"
	ErrorCodeTokenExpired       ErrorCode = "TOKEN_EXPIRED"
	ErrorCodePermissionDenied   ErrorCode = "PERMISSION_DENIED"
	ErrorCodeUnauthorized       ErrorCode = "UNAUTHORIZED"

	// API and request errors
	ErrorCodeInvalidRequest   ErrorCode = "INVALID_REQUEST"
	ErrorCodeValidationFailed ErrorCode = "VALIDATION_FAILED"
	ErrorCodeResourceNotFound ErrorCode = "RESOURCE_NOT_FOUND"
	ErrorCodeConflict         ErrorCode = "CONFLICT"
	ErrorCodeRateLimited      ErrorCode = "RATE_LIMITED"

	// Server and Slurm errors
	ErrorCodeServerInternal       ErrorCode = "SERVER_INTERNAL"
	ErrorCodeSlurmDaemonDown      ErrorCode = "SLURM_DAEMON_DOWN"
	ErrorCodeServiceUnavailable   ErrorCode = "SERVICE_UNAVAILABLE"
	ErrorCodeResourceExhausted    ErrorCode = "RESOURCE_EXHAUSTED"
	ErrorCodeJobQueueFull         ErrorCode = "JOB_QUEUE_FULL"
	ErrorCodePartitionUnavailable ErrorCode = "PARTITION_UNAVAILABLE"

	// Client and configuration errors
	ErrorCodeClientNotInitialized ErrorCode = "CLIENT_NOT_INITIALIZED"
	ErrorCodeInvalidConfiguration ErrorCode = "INVALID_CONFIGURATION"
	ErrorCodeVersionMismatch      ErrorCode = "VERSION_MISMATCH"
	ErrorCodeUnsupportedOperation ErrorCode = "UNSUPPORTED_OPERATION"

	// Context and cancellation errors
	ErrorCodeContextCanceled  ErrorCode = "CONTEXT_CANCELED"
	ErrorCodeDeadlineExceeded ErrorCode = "DEADLINE_EXCEEDED"

	// Unknown or unclassified errors
	ErrorCodeUnknown ErrorCode = "UNKNOWN"
)

// ErrorCategory groups related error codes for easier handling
type ErrorCategory string

const (
	CategoryNetwork        ErrorCategory = "NETWORK"
	CategoryAuthentication ErrorCategory = "AUTHENTICATION"
	CategoryValidation     ErrorCategory = "VALIDATION"
	CategoryResource       ErrorCategory = "RESOURCE"
	CategoryServer         ErrorCategory = "SERVER"
	CategoryClient         ErrorCategory = "CLIENT"
	CategoryContext        ErrorCategory = "CONTEXT"
	CategoryUnknown        ErrorCategory = "UNKNOWN"
)

// SlurmError represents a structured error from the Slurm client
type SlurmError struct {
	Code       ErrorCode     `json:"code"`
	Category   ErrorCategory `json:"category"`
	Message    string        `json:"message"`
	Details    string        `json:"details,omitempty"`
	Timestamp  time.Time     `json:"timestamp"`
	StatusCode int           `json:"status_code,omitempty"`
	APIVersion string        `json:"api_version,omitempty"`
	RequestID  string        `json:"request_id,omitempty"`
	Retryable  bool          `json:"retryable"`
	Cause      error         `json:"-"` // Original error, not serialized
}

// Error implements the error interface
func (e *SlurmError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns the underlying cause error
func (e *SlurmError) Unwrap() error {
	return e.Cause
}

// Is checks if the error matches a specific error code
func (e *SlurmError) Is(target error) bool {
	if targetErr, ok := target.(*SlurmError); ok {
		return e.Code == targetErr.Code
	}
	return false
}

// IsRetryable returns true if the error indicates the operation can be retried
func (e *SlurmError) IsRetryable() bool {
	return e.Retryable
}

// IsTemporary returns true if the error is likely temporary
func (e *SlurmError) IsTemporary() bool {
	return e.Category == CategoryNetwork ||
		e.Code == ErrorCodeServerInternal ||
		e.Code == ErrorCodeResourceExhausted ||
		e.Code == ErrorCodeRateLimited
}

// NetworkError represents network-related errors
type NetworkError struct {
	*SlurmError
	Host string `json:"host,omitempty"`
	Port int    `json:"port,omitempty"`
}

// AuthenticationError represents authentication/authorization errors
type AuthenticationError struct {
	*SlurmError
	AuthMethod string `json:"auth_method,omitempty"`
	TokenType  string `json:"token_type,omitempty"`
}

// ValidationError represents request validation errors
type ValidationError struct {
	*SlurmError
	Field string      `json:"field,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

// SlurmAPIError represents errors returned by the Slurm REST API
type SlurmAPIError struct {
	*SlurmError
	ErrorNumber int                   `json:"error_number,omitempty"`
	ErrorCode   string                `json:"error_code,omitempty"`
	Source      string                `json:"source,omitempty"`
	Errors      []SlurmAPIErrorDetail `json:"errors,omitempty"`
}

// SlurmAPIErrorDetail represents detailed error information from Slurm API responses
type SlurmAPIErrorDetail struct {
	ErrorNumber int    `json:"error_number"`
	ErrorCode   string `json:"error_code"`
	Source      string `json:"source"`
	Description string `json:"description"`
}

// NewSlurmError creates a new structured Slurm error
func NewSlurmError(code ErrorCode, message string) *SlurmError {
	return &SlurmError{
		Code:      code,
		Category:  getErrorCategory(code),
		Message:   message,
		Timestamp: time.Now(),
		Retryable: isRetryable(code),
	}
}

// NewSlurmErrorWithCause creates a new structured Slurm error with an underlying cause
func NewSlurmErrorWithCause(code ErrorCode, message string, cause error) *SlurmError {
	return &SlurmError{
		Code:      code,
		Category:  getErrorCategory(code),
		Message:   message,
		Timestamp: time.Now(),
		Retryable: isRetryable(code),
		Cause:     cause,
	}
}

// NewNetworkError creates a new network error
func NewNetworkError(code ErrorCode, message, host string, port int, cause error) *NetworkError {
	return &NetworkError{
		SlurmError: NewSlurmErrorWithCause(code, message, cause),
		Host:       host,
		Port:       port,
	}
}

// NewAuthenticationError creates a new authentication error
func NewAuthenticationError(code ErrorCode, message, authMethod, tokenType string, cause error) *AuthenticationError {
	return &AuthenticationError{
		SlurmError: NewSlurmErrorWithCause(code, message, cause),
		AuthMethod: authMethod,
		TokenType:  tokenType,
	}
}

// NewValidationError creates a new validation error
func NewValidationError(code ErrorCode, message, field string, value interface{}, cause error) *ValidationError {
	return &ValidationError{
		SlurmError: NewSlurmErrorWithCause(code, message, cause),
		Field:      field,
		Value:      value,
	}
}

// NewSlurmAPIError creates a new Slurm API error from API response
func NewSlurmAPIError(statusCode int, apiVersion string, details []SlurmAPIErrorDetail) *SlurmAPIError {
	var code ErrorCode
	var message string

	if len(details) > 0 {
		// Use first error for primary classification
		primary := details[0]
		code = mapSlurmErrorCodeToClientCode(primary.ErrorCode, statusCode)
		message = primary.Description
	} else {
		code = mapHTTPStatusToErrorCode(statusCode)
		message = http.StatusText(statusCode)
	}

	return &SlurmAPIError{
		SlurmError: &SlurmError{
			Code:       code,
			Category:   getErrorCategory(code),
			Message:    message,
			Timestamp:  time.Now(),
			StatusCode: statusCode,
			APIVersion: apiVersion,
			Retryable:  isRetryable(code),
		},
		Errors: details,
	}
}

// getErrorCategory maps error codes to categories
func getErrorCategory(code ErrorCode) ErrorCategory {
	switch code {
	case ErrorCodeNetworkTimeout, ErrorCodeConnectionRefused, ErrorCodeDNSResolution, ErrorCodeTLSHandshake:
		return CategoryNetwork
	case ErrorCodeInvalidCredentials, ErrorCodeTokenExpired, ErrorCodePermissionDenied, ErrorCodeUnauthorized:
		return CategoryAuthentication
	case ErrorCodeInvalidRequest, ErrorCodeValidationFailed:
		return CategoryValidation
	case ErrorCodeResourceNotFound, ErrorCodeConflict, ErrorCodeResourceExhausted, ErrorCodeJobQueueFull, ErrorCodePartitionUnavailable:
		return CategoryResource
	case ErrorCodeServerInternal, ErrorCodeSlurmDaemonDown, ErrorCodeRateLimited:
		return CategoryServer
	case ErrorCodeClientNotInitialized, ErrorCodeInvalidConfiguration, ErrorCodeVersionMismatch, ErrorCodeUnsupportedOperation:
		return CategoryClient
	case ErrorCodeContextCanceled, ErrorCodeDeadlineExceeded:
		return CategoryContext
	default:
		return CategoryUnknown
	}
}

// isRetryable determines if an error code indicates a retryable operation
func isRetryable(code ErrorCode) bool {
	switch code {
	case ErrorCodeNetworkTimeout, ErrorCodeConnectionRefused, ErrorCodeDNSResolution,
		ErrorCodeServerInternal, ErrorCodeSlurmDaemonDown, ErrorCodeResourceExhausted,
		ErrorCodeRateLimited:
		return true
	default:
		return false
	}
}

// mapHTTPStatusToErrorCode maps HTTP status codes to Slurm error codes
func mapHTTPStatusToErrorCode(statusCode int) ErrorCode {
	switch statusCode {
	case http.StatusBadRequest:
		return ErrorCodeInvalidRequest
	case http.StatusUnauthorized:
		return ErrorCodeUnauthorized
	case http.StatusForbidden:
		return ErrorCodePermissionDenied
	case http.StatusNotFound:
		return ErrorCodeResourceNotFound
	case http.StatusConflict:
		return ErrorCodeConflict
	case http.StatusUnprocessableEntity:
		return ErrorCodeValidationFailed
	case http.StatusTooManyRequests:
		return ErrorCodeRateLimited
	case http.StatusInternalServerError:
		return ErrorCodeServerInternal
	case http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return ErrorCodeSlurmDaemonDown
	default:
		return ErrorCodeUnknown
	}
}

// mapSlurmErrorCodeToClientCode maps Slurm-specific error codes to client error codes
func mapSlurmErrorCodeToClientCode(slurmErrorCode string, statusCode int) ErrorCode {
	// Map known Slurm error codes
	switch slurmErrorCode {
	case "SLURM_NO_CHANGE_IN_DATA":
		return ErrorCodeResourceNotFound
	case "SLURM_PROTOCOL_VERSION_ERROR":
		return ErrorCodeVersionMismatch
	case "SLURM_UNEXPECTED_MSG_ERROR":
		return ErrorCodeInvalidRequest
	case "SLURM_COMMUNICATIONS_CONNECTION_ERROR":
		return ErrorCodeConnectionRefused
	case "SLURM_COMMUNICATIONS_SEND_ERROR", "SLURM_COMMUNICATIONS_RECEIVE_ERROR":
		return ErrorCodeNetworkTimeout
	case "SLURM_AUTHENTICATION_ERROR":
		return ErrorCodeInvalidCredentials
	case "SLURM_ACCESS_DENIED":
		return ErrorCodePermissionDenied
	case "SLURM_JOB_PENDING":
		return ErrorCodeResourceExhausted
	case "SLURM_INVALID_PARTITION_NAME":
		return ErrorCodePartitionUnavailable
	case "SLURM_INVALID_JOB_ID":
		return ErrorCodeResourceNotFound
	case "SLURM_JOB_ALREADY_COMPLETE":
		return ErrorCodeConflict
	case "SLURM_NODE_NOT_AVAIL":
		return ErrorCodeResourceExhausted
	default:
		// Fall back to HTTP status code mapping
		return mapHTTPStatusToErrorCode(statusCode)
	}
}
