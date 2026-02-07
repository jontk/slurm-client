// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_44

import (
	"encoding/json"
	stderrors "errors"
	"fmt"
	"net/http"

	api "github.com/jontk/slurm-client/internal/openapi/v0_0_44"
	"github.com/jontk/slurm-client/pkg/errors"
)

// ErrorAdapter handles API-specific error conversion for v0.0.44
type ErrorAdapter struct{}

// NewErrorAdapter creates a new error adapter for v0.0.44
func NewErrorAdapter() *ErrorAdapter {
	return &ErrorAdapter{}
}

// HandleAPIResponse processes the API response and returns appropriate errors
func (e *ErrorAdapter) HandleAPIResponse(statusCode int, body []byte, operation string) error {
	// Check for successful responses
	if statusCode >= 200 && statusCode < 300 {
		return nil
	}
	// Try to parse the error response
	var apiResp struct {
		Meta     *api.V0044OpenapiMeta     `json:"meta,omitempty"`
		Errors   *api.V0044OpenapiErrors   `json:"errors,omitempty"`
		Warnings *api.V0044OpenapiWarnings `json:"warnings,omitempty"`
	}
	var details []errors.SlurmAPIErrorDetail
	if err := json.Unmarshal(body, &apiResp); err == nil {
		// Check for errors in the meta field (v0.0.44 doesn't have Errors in meta)
		// Skip this check for v0.0.44
		// Also check for errors in the errors field
		if apiResp.Errors != nil {
			for _, apiErr := range *apiResp.Errors {
				detail := errors.SlurmAPIErrorDetail{}
				if apiErr.ErrorNumber != nil {
					detail.ErrorNumber = int(*apiErr.ErrorNumber)
				}
				if apiErr.Error != nil {
					detail.ErrorCode = *apiErr.Error
				}
				if apiErr.Description != nil {
					detail.Description = *apiErr.Description
				} else if apiErr.Error != nil {
					detail.Description = *apiErr.Error
				}
				if apiErr.Source != nil {
					detail.Source = *apiErr.Source
				}
				details = append(details, detail)
			}
		}
	}
	// Create a structured error using the error package
	if len(details) > 0 {
		apiErr := errors.NewSlurmAPIError(statusCode, "v0.0.44", details)
		apiErr.SlurmError.Details = "Operation: " + operation
		return apiErr
	}
	// Handle specific HTTP status codes
	switch statusCode {
	case http.StatusUnauthorized:
		return errors.NewAuthenticationError(
			errors.ErrorCodeUnauthorized,
			"Authentication failed",
			"",
			"",
			fmt.Errorf("HTTP 401: %s", string(body)),
		)
	case http.StatusForbidden:
		return errors.NewAuthenticationError(
			errors.ErrorCodePermissionDenied,
			"Permission denied",
			"",
			"",
			fmt.Errorf("HTTP 403: %s", string(body)),
		)
	case http.StatusNotFound:
		err := errors.NewSlurmError(errors.ErrorCodeResourceNotFound, operation+": resource not found (404)")
		err.StatusCode = statusCode
		err.Details = string(body)
		return err
	case http.StatusConflict:
		err := errors.NewSlurmError(errors.ErrorCodeConflict, operation+": resource conflict (409)")
		err.StatusCode = statusCode
		err.Details = string(body)
		return err
	case http.StatusUnprocessableEntity:
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			operation+": validation failed",
			"",
			nil,
			fmt.Errorf("HTTP 422: %s", string(body)),
		)
	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable:
		err := errors.NewSlurmError(errors.ErrorCodeServerInternal, fmt.Sprintf("%s: server error (%d)", operation, statusCode))
		err.StatusCode = statusCode
		err.Details = string(body)
		return err
	default:
		err := errors.NewSlurmError(errors.ErrorCodeUnknown, fmt.Sprintf("%s failed with status %d", operation, statusCode))
		err.StatusCode = statusCode
		err.Details = string(body)
		return err
	}
}

// ParseSlurmError attempts to extract SLURM-specific error information
func (e *ErrorAdapter) ParseSlurmError(err error) (string, string, int) {
	// Default values
	code := "UNKNOWN"
	message := err.Error()
	errno := -1
	// Try to extract more specific error information
	var apiErr *errors.SlurmAPIError
	if stderrors.As(err, &apiErr) {
		if len(apiErr.Errors) > 0 {
			code = apiErr.Errors[0].ErrorCode
			message = apiErr.Errors[0].Description
			errno = apiErr.Errors[0].ErrorNumber
		}
	}
	return code, message, errno
}
