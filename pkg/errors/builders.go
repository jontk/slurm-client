// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package errors

import (
	"context"
	stderrors "errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"syscall"
)

// WrapError converts a generic error into a structured SlurmError
func WrapError(err error) *SlurmError {
	if err == nil {
		return nil
	}

	// If already a SlurmError, return as-is
	var slurmErr *SlurmError
	if stderrors.As(err, &slurmErr) {
		return slurmErr
	}

	// Check for context errors first
	if stderrors.Is(err, context.Canceled) {
		return NewSlurmErrorWithCause(ErrorCodeContextCanceled, "Operation was canceled", err)
	}
	if stderrors.Is(err, context.DeadlineExceeded) {
		return NewSlurmErrorWithCause(ErrorCodeDeadlineExceeded, "Operation timed out", err)
	}

	// Check for network errors
	if netErr := classifyNetworkError(err); netErr != nil {
		return netErr
	}

	// Check for URL errors
	var urlErr *url.Error
	if stderrors.As(err, &urlErr) {
		return classifyURLError(urlErr)
	}

	// Default to unknown error
	return NewSlurmErrorWithCause(ErrorCodeUnknown, err.Error(), err)
}

// WrapHTTPError converts an HTTP response error into a structured SlurmError
func WrapHTTPError(statusCode int, body []byte, apiVersion string) *SlurmError {
	// Try to parse Slurm API error response
	if apiErr := parseSlurmAPIError(statusCode, body, apiVersion); apiErr != nil {
		return apiErr.SlurmError
	}

	// Fall back to HTTP status code mapping
	code := mapHTTPStatusToErrorCode(statusCode)
	message := fmt.Sprintf("HTTP %d: %s", statusCode, http.StatusText(statusCode))

	slurmErr := NewSlurmError(code, message)
	slurmErr.StatusCode = statusCode
	slurmErr.APIVersion = apiVersion

	if len(body) > 0 && len(body) < 1000 { // Include response body if reasonable size
		slurmErr.Details = string(body)
	}

	return slurmErr
}

// classifyNetworkError identifies and wraps network-related errors
func classifyNetworkError(err error) *SlurmError {
	if err == nil {
		return nil
	}

	// Check for context errors first (before net.Error check)
	// because context.DeadlineExceeded also implements net.Error with Timeout() == true
	// Use errors.Is() to handle wrapped errors
	if stderrors.Is(err, context.Canceled) {
		return NewSlurmErrorWithCause(ErrorCodeContextCanceled, "Operation was canceled", err)
	}
	if stderrors.Is(err, context.DeadlineExceeded) {
		return NewSlurmErrorWithCause(ErrorCodeDeadlineExceeded, "Operation deadline exceeded", err)
	}

	errStr := err.Error()

	// Check for specific network error types
	var netErr net.Error
	if stderrors.As(err, &netErr) {
		if netErr.Timeout() {
			return NewSlurmErrorWithCause(ErrorCodeNetworkTimeout, "Network operation timed out", err)
		}
		// Note: netErr.Temporary() is deprecated since Go 1.18
		// We classify common temporary network errors by error string patterns
		errorStr := err.Error()
		if strings.Contains(errorStr, "connection reset") ||
			strings.Contains(errorStr, "broken pipe") ||
			strings.Contains(errorStr, "network is unreachable") ||
			strings.Contains(errorStr, "temporary") {
			return NewSlurmErrorWithCause(ErrorCodeConnectionRefused, "Temporary network failure", err)
		}
	}

	// Check for specific error patterns
	switch {
	case strings.Contains(errStr, "connection refused"):
		return NewSlurmErrorWithCause(ErrorCodeConnectionRefused, "Connection refused by server", err)
	case strings.Contains(errStr, "no such host"):
		return NewSlurmErrorWithCause(ErrorCodeDNSResolution, "DNS resolution failed", err)
	case strings.Contains(errStr, "timeout"):
		return NewSlurmErrorWithCause(ErrorCodeNetworkTimeout, "Network timeout", err)
	case strings.Contains(errStr, "tls"):
		return NewSlurmErrorWithCause(ErrorCodeTLSHandshake, "TLS handshake failed", err)
	case strings.Contains(errStr, "certificate"):
		return NewSlurmErrorWithCause(ErrorCodeTLSHandshake, "TLS certificate error", err)
	}

	// Check for syscall errors
	var opErr *net.OpError
	if stderrors.As(err, &opErr) {
		var dnsErr *net.DNSError
		if stderrors.As(opErr.Err, &dnsErr) {
			return NewSlurmErrorWithCause(ErrorCodeDNSResolution, "DNS lookup failed", dnsErr)
		}
		var syscallErr syscall.Errno
		if stderrors.As(opErr.Err, &syscallErr) {
			switch syscallErr {
			case syscall.ECONNREFUSED:
				return NewSlurmErrorWithCause(ErrorCodeConnectionRefused, "Connection refused", err)
			case syscall.ETIMEDOUT:
				return NewSlurmErrorWithCause(ErrorCodeNetworkTimeout, "Connection timeout", err)
			case syscall.ENETUNREACH:
				return NewSlurmErrorWithCause(ErrorCodeDNSResolution, "Network unreachable", err)
			}
		}
	}

	return nil
}

// classifyURLError handles URL-specific errors
func classifyURLError(urlErr *url.Error) *SlurmError {
	// Extract host and port for network errors
	var host string
	var port int
	if u, err := url.Parse(urlErr.URL); err == nil {
		host = u.Hostname()
		if u.Port() != "" {
			_, _ = fmt.Sscanf(u.Port(), "%d", &port) // Ignore error, port parsing is best-effort
		}
	}

	// Check for context errors first (before network classification)
	if stderrors.Is(urlErr.Err, context.Canceled) {
		return NewSlurmErrorWithCause(ErrorCodeContextCanceled, "Operation was canceled", urlErr)
	}
	if stderrors.Is(urlErr.Err, context.DeadlineExceeded) {
		return NewSlurmErrorWithCause(ErrorCodeDeadlineExceeded, "Operation deadline exceeded", urlErr)
	}

	// Check underlying error
	if netErr := classifyNetworkError(urlErr.Err); netErr != nil {
		if host != "" {
			networkErr := &NetworkError{
				SlurmError: netErr,
				Host:       host,
				Port:       port,
			}
			return networkErr.SlurmError
		}
		return netErr
	}

	// Default URL error handling
	return NewSlurmErrorWithCause(ErrorCodeNetworkTimeout, "URL error: "+urlErr.Op, urlErr)
}

// NewClientError creates errors for client-side issues
func NewClientError(code ErrorCode, message string, details ...string) *SlurmError {
	err := NewSlurmError(code, message)
	if len(details) > 0 {
		err.Details = strings.Join(details, "; ")
	}
	return err
}

// NewAuthError creates authentication-related errors
func NewAuthError(authMethod, tokenType string, cause error) *AuthenticationError {
	var code ErrorCode
	var message string

	switch {
	case strings.Contains(cause.Error(), "401"):
		code = ErrorCodeInvalidCredentials
		message = "Invalid credentials provided"
	case strings.Contains(cause.Error(), "403"):
		code = ErrorCodePermissionDenied
		message = "Access denied"
	case strings.Contains(cause.Error(), "expired"):
		code = ErrorCodeTokenExpired
		message = "Authentication token has expired"
	default:
		code = ErrorCodeUnauthorized
		message = "Authentication failed"
	}

	return NewAuthenticationError(code, message, authMethod, tokenType, cause)
}

// NewValidationErrorf creates a validation error with formatted message
func NewValidationErrorf(field string, value interface{}, format string, args ...interface{}) *ValidationError {
	message := fmt.Sprintf(format, args...)
	return NewValidationError(ErrorCodeValidationFailed, message, field, value, nil)
}

// NewJobError creates job-specific errors
func NewJobError(jobID uint32, operation string, cause error) *SlurmError {
	var code ErrorCode
	var message string

	causeStr := cause.Error()
	switch {
	case strings.Contains(causeStr, "not found") || strings.Contains(causeStr, "invalid job id"):
		code = ErrorCodeResourceNotFound
		message = fmt.Sprintf("Job %d not found", jobID)
	case strings.Contains(causeStr, "already complete"):
		code = ErrorCodeConflict
		message = fmt.Sprintf("Job %d is already completed", jobID)
	case strings.Contains(causeStr, "permission denied"):
		code = ErrorCodePermissionDenied
		message = fmt.Sprintf("Permission denied for job %d", jobID)
	case strings.Contains(causeStr, "queue full"):
		code = ErrorCodeJobQueueFull
		message = "Job queue is full"
	default:
		code = ErrorCodeServerInternal
		message = fmt.Sprintf("Job %s failed for job %d", operation, jobID)
	}

	err := NewSlurmErrorWithCause(code, message, cause)
	err.Details = fmt.Sprintf("Job ID: %d, Operation: %s", jobID, operation)
	return err
}

// NewNodeError creates node-specific errors
func NewNodeError(nodeNames []string, operation string, cause error) *SlurmError {
	var code ErrorCode
	var message string

	causeStr := cause.Error()
	switch {
	case strings.Contains(causeStr, "not found"):
		code = ErrorCodeResourceNotFound
		message = fmt.Sprintf("Nodes not found: %v", nodeNames)
	case strings.Contains(causeStr, "not available"):
		code = ErrorCodeResourceExhausted
		message = fmt.Sprintf("Nodes not available: %v", nodeNames)
	case strings.Contains(causeStr, "permission denied"):
		code = ErrorCodePermissionDenied
		message = fmt.Sprintf("Permission denied for nodes: %v", nodeNames)
	default:
		code = ErrorCodeServerInternal
		message = fmt.Sprintf("Node %s failed", operation)
	}

	err := NewSlurmErrorWithCause(code, message, cause)
	err.Details = fmt.Sprintf("Nodes: %v, Operation: %s", nodeNames, operation)
	return err
}

// NewPartitionError creates partition-specific errors
func NewPartitionError(partitionName, operation string, cause error) *SlurmError {
	var code ErrorCode
	var message string

	causeStr := cause.Error()
	switch {
	case strings.Contains(causeStr, "not found") || strings.Contains(causeStr, "invalid partition"):
		code = ErrorCodeResourceNotFound
		message = fmt.Sprintf("Partition '%s' not found", partitionName)
	case strings.Contains(causeStr, "unavailable"):
		code = ErrorCodePartitionUnavailable
		message = fmt.Sprintf("Partition '%s' is unavailable", partitionName)
	case strings.Contains(causeStr, "permission denied"):
		code = ErrorCodePermissionDenied
		message = fmt.Sprintf("Permission denied for partition '%s'", partitionName)
	default:
		code = ErrorCodeServerInternal
		message = fmt.Sprintf("Partition %s failed", operation)
	}

	err := NewSlurmErrorWithCause(code, message, cause)
	err.Details = fmt.Sprintf("Partition: %s, Operation: %s", partitionName, operation)
	return err
}

// IsRetryableError checks if an error is retryable
func IsRetryableError(err error) bool {
	var slurmErr *SlurmError
	if stderrors.As(err, &slurmErr) {
		return slurmErr.IsRetryable()
	}

	// Check for known retryable error patterns
	if err != nil {
		errStr := err.Error()
		return strings.Contains(errStr, "timeout") ||
			strings.Contains(errStr, "connection refused") ||
			strings.Contains(errStr, "temporary failure") ||
			strings.Contains(errStr, "service unavailable")
	}

	return false
}

// IsTemporaryError checks if an error is temporary
func IsTemporaryError(err error) bool {
	if err == nil {
		return false
	}

	var slurmErr *SlurmError
	if stderrors.As(err, &slurmErr) {
		return slurmErr.IsTemporary()
	}

	// Check for net.Error interface
	// Note: netErr.Temporary() is deprecated since Go 1.18
	// We classify common temporary errors by timeout or error string patterns
	var netErr net.Error
	if stderrors.As(err, &netErr) {
		if netErr.Timeout() {
			return true
		}
	}

	// Check for common temporary error patterns
	errorStr := err.Error()
	if strings.Contains(errorStr, "connection reset") ||
		strings.Contains(errorStr, "broken pipe") ||
		strings.Contains(errorStr, "network is unreachable") ||
		strings.Contains(errorStr, "temporary") {
		return true
	}

	return false
}

// GetErrorCode extracts the error code from any error
func GetErrorCode(err error) ErrorCode {
	var slurmErr *SlurmError
	if stderrors.As(err, &slurmErr) {
		return slurmErr.Code
	}
	return ErrorCodeUnknown
}

// GetErrorCategory extracts the error category from any error
func GetErrorCategory(err error) ErrorCategory {
	var slurmErr *SlurmError
	if stderrors.As(err, &slurmErr) {
		return slurmErr.Category
	}
	return CategoryUnknown
}

// IsNetworkError checks if an error is a network-related error
func IsNetworkError(err error) bool {
	if err == nil {
		return false
	}

	// Check if it's a SlurmError with network category
	var slurmErr *SlurmError
	if stderrors.As(err, &slurmErr) {
		return slurmErr.Category == CategoryNetwork
	}

	// Check if it's a direct network error
	var netErr net.Error
	if stderrors.As(err, &netErr) {
		return true
	}

	// Check for URL errors
	var urlErr *url.Error
	if stderrors.As(err, &urlErr) {
		return true
	}

	// Check for specific network error patterns
	errMsg := strings.ToLower(err.Error())
	networkPatterns := []string{
		"connection refused",
		"connection reset",
		"no such host",
		"network unreachable",
		"timeout",
		"tls handshake",
		"dns",
	}

	for _, pattern := range networkPatterns {
		if strings.Contains(errMsg, pattern) {
			return true
		}
	}

	return false
}

// IsAuthenticationError checks if an error is an authentication-related error
func IsAuthenticationError(err error) bool {
	if err == nil {
		return false
	}

	// Check if it's a SlurmError with authentication category
	var slurmErr *SlurmError
	if stderrors.As(err, &slurmErr) {
		return slurmErr.Category == CategoryAuthentication
	}

	// Check for specific authentication error patterns
	errMsg := strings.ToLower(err.Error())
	authPatterns := []string{
		"unauthorized",
		"authentication failed",
		"invalid token",
		"expired token",
		"permission denied",
		"access denied",
		"forbidden",
	}

	for _, pattern := range authPatterns {
		if strings.Contains(errMsg, pattern) {
			return true
		}
	}

	return false
}

// NewNotImplementedError creates errors for operations not yet implemented
func NewNotImplementedError(operation, version string) *SlurmError {
	message := fmt.Sprintf("Operation '%s' not implemented in version %s", operation, version)
	err := NewSlurmError(ErrorCodeUnsupportedOperation, message)
	err.Details = "Version: " + version
	return err
}

// IsNotImplementedError checks if an error is a not implemented error
func IsNotImplementedError(err error) bool {
	var slurmErr *SlurmError
	if stderrors.As(err, &slurmErr) {
		return slurmErr.Code == ErrorCodeUnsupportedOperation
	}
	return false
}

// IsClientError checks if an error is a client-side error
func IsClientError(err error) bool {
	// Check if it's a SlurmError with client category
	var slurmErr *SlurmError
	if stderrors.As(err, &slurmErr) {
		return slurmErr.Category == CategoryClient
	}
	return false
}

// IsValidationError checks if an error is a validation error
func IsValidationError(err error) bool {
	// Check if it's directly a ValidationError
	var valErr *ValidationError
	if stderrors.As(err, &valErr) {
		return true
	}
	// Check if it's a SlurmError with validation category
	var slurmErr *SlurmError
	if stderrors.As(err, &slurmErr) {
		return slurmErr.Category == CategoryValidation
	}
	return false
}
