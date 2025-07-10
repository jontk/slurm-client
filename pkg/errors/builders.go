package errors

import (
	"context"
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
	if slurmErr, ok := err.(*SlurmError); ok {
		return slurmErr
	}

	// Check for context errors first
	if err == context.Canceled {
		return NewSlurmErrorWithCause(ErrorCodeContextCanceled, "Operation was canceled", err)
	}
	if err == context.DeadlineExceeded {
		return NewSlurmErrorWithCause(ErrorCodeDeadlineExceeded, "Operation timed out", err)
	}

	// Check for network errors
	if netErr := classifyNetworkError(err); netErr != nil {
		return netErr
	}

	// Check for URL errors
	if urlErr, ok := err.(*url.Error); ok {
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

	errStr := err.Error()

	// Check for specific network error types
	if netErr, ok := err.(net.Error); ok {
		if netErr.Timeout() {
			return NewSlurmErrorWithCause(ErrorCodeNetworkTimeout, "Network operation timed out", err)
		}
		if netErr.Temporary() {
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
	if opErr, ok := err.(*net.OpError); ok {
		if syscallErr, ok := opErr.Err.(*net.DNSError); ok {
			return NewSlurmErrorWithCause(ErrorCodeDNSResolution, "DNS lookup failed", syscallErr)
		}
		if syscallErr, ok := opErr.Err.(syscall.Errno); ok {
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
			fmt.Sscanf(u.Port(), "%d", &port)
		}
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
	return NewSlurmErrorWithCause(ErrorCodeNetworkTimeout, fmt.Sprintf("URL error: %s", urlErr.Op), urlErr)
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
	if slurmErr, ok := err.(*SlurmError); ok {
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
	if slurmErr, ok := err.(*SlurmError); ok {
		return slurmErr.IsTemporary()
	}
	
	// Check for net.Error interface
	if netErr, ok := err.(net.Error); ok {
		return netErr.Temporary()
	}
	
	return false
}

// GetErrorCode extracts the error code from any error
func GetErrorCode(err error) ErrorCode {
	if slurmErr, ok := err.(*SlurmError); ok {
		return slurmErr.Code
	}
	return ErrorCodeUnknown
}

// GetErrorCategory extracts the error category from any error
func GetErrorCategory(err error) ErrorCategory {
	if slurmErr, ok := err.(*SlurmError); ok {
		return slurmErr.Category
	}
	return CategoryUnknown
}