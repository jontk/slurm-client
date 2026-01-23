// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_41

import (
	"encoding/json"
	stderrors "errors"
	"fmt"
	"net/http"

	"github.com/jontk/slurm-client/pkg/errors"
)

// ErrorAdapter handles API-specific error conversion for v0.0.41
type ErrorAdapter struct{}

// NewErrorAdapter creates a new error adapter for v0.0.41
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
		Meta *struct {
			Plugin *struct {
				AccountingStorage *string `json:"accounting_storage,omitempty"`
				DataParser        *string `json:"data_parser,omitempty"`
				Name              *string `json:"name,omitempty"`
				Type              *string `json:"type,omitempty"`
			} `json:"plugin,omitempty"`
			SlurmVersion *struct {
				Major *struct {
					Infinite *bool  `json:"infinite,omitempty"`
					Number   *int64 `json:"number,omitempty"`
					Set      *bool  `json:"set,omitempty"`
				} `json:"major,omitempty"`
				Micro *struct {
					Infinite *bool  `json:"infinite,omitempty"`
					Number   *int64 `json:"number,omitempty"`
					Set      *bool  `json:"set,omitempty"`
				} `json:"micro,omitempty"`
				Minor *struct {
					Infinite *bool  `json:"infinite,omitempty"`
					Number   *int64 `json:"number,omitempty"`
					Set      *bool  `json:"set,omitempty"`
				} `json:"minor,omitempty"`
			} `json:"slurm_version,omitempty"`
		} `json:"meta,omitempty"`
		Errors []struct {
			Description *string `json:"description,omitempty"`
			Error       *string `json:"error,omitempty"`
			ErrorNumber *int32  `json:"error_number,omitempty"`
			Source      *string `json:"source,omitempty"`
		} `json:"errors,omitempty"`
		Warnings []struct {
			Description *string `json:"description,omitempty"`
			Source      *string `json:"source,omitempty"`
		} `json:"warnings,omitempty"`
	}

	var details []errors.SlurmAPIErrorDetail
	if err := json.Unmarshal(body, &apiResp); err == nil {
		// Check for errors in the errors field
		if len(apiResp.Errors) > 0 {
			for _, apiErr := range apiResp.Errors {
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
		apiErr := errors.NewSlurmAPIError(statusCode, "v0.0.41", details)
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
		err := errors.NewSlurmError(errors.ErrorCodeResourceNotFound, operation+": resource not found")
		err.StatusCode = statusCode
		err.Details = string(body)
		return err
	case http.StatusConflict:
		err := errors.NewSlurmError(errors.ErrorCodeConflict, operation+": resource conflict")
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
		err := errors.NewSlurmError(errors.ErrorCodeServerInternal, operation+": server error")
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
