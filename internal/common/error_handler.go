// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"fmt"
	"reflect"
	
	"github.com/jontk/slurm-client/pkg/errors"
)

// ErrorResponse represents a generic interface for API error responses
type ErrorResponse interface {
	GetErrors() []ErrorDetail
}

// ErrorDetail represents a generic error detail interface
type ErrorDetail interface {
	GetErrorNumber() *int
	GetError() *string
	GetSource() *string
	GetDescription() *string
}

// ResponseWithErrors represents a response that may contain errors
type ResponseWithErrors interface {
	StatusCode() int
	HasErrors() bool
	GetErrorResponse() ErrorResponse
}

// HandleAPIResponse processes common API response error patterns
func HandleAPIResponse(resp ResponseWithErrors, version string) error {
	// Check HTTP status
	if resp.StatusCode() == 200 {
		return nil
	}

	// Try to extract detailed error information
	if resp.HasErrors() {
		errResp := resp.GetErrorResponse()
		if errResp != nil {
			errs := errResp.GetErrors()
			if len(errs) > 0 {
				apiErrors := make([]errors.SlurmAPIErrorDetail, len(errs))
				for i, apiErr := range errs {
					apiErrors[i] = extractErrorDetail(apiErr)
				}
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), version, apiErrors)
				return apiError.SlurmError
			}
		}
	}

	// Fall back to HTTP error
	return errors.WrapHTTPError(resp.StatusCode(), nil, version)
}

// extractErrorDetail converts a generic error detail to SlurmAPIErrorDetail
func extractErrorDetail(err ErrorDetail) errors.SlurmAPIErrorDetail {
	detail := errors.SlurmAPIErrorDetail{}
	
	if num := err.GetErrorNumber(); num != nil {
		detail.ErrorNumber = *num
		
		// Enhance error description with SLURM-specific information
		if desc := err.GetDescription(); desc != nil {
			enhancedDesc := EnhanceErrorMessage(int32(*num), *desc)
			detail.Description = enhancedDesc
		} else {
			// If no description provided, use SLURM error description
			detail.Description = GetErrorDescription(int32(*num))
		}
		
		// Add error category
		if info := GetErrorInfo(int32(*num)); info != nil {
			detail.ErrorCode = info.Name
		}
	} else {
		// No error number, use original description
		if desc := err.GetDescription(); desc != nil {
			detail.Description = *desc
		}
		
		if code := err.GetError(); code != nil {
			detail.ErrorCode = *code
		}
	}
	
	if src := err.GetSource(); src != nil {
		detail.Source = *src
	}
	
	return detail
}

// CheckNilResponse checks if a response is nil and returns appropriate error
func CheckNilResponse(response interface{}, operation string) error {
	if response == nil {
		return errors.NewClientError(
			errors.ErrorCodeServerInternal, 
			"Unexpected response format", 
			"Expected JSON response but got nil for "+operation,
		)
	}
	
	// Use reflection to check for typed nil pointers
	if isNilPointer(response) {
		return errors.NewClientError(
			errors.ErrorCodeServerInternal, 
			"Unexpected response format", 
			"Expected JSON response but got nil for "+operation,
		)
	}
	
	return nil
}

// isNilPointer checks if an interface contains a typed nil pointer
func isNilPointer(i interface{}) bool {
	if i == nil {
		return true
	}
	
	// Use reflection to check if it's a nil pointer
	v := reflect.ValueOf(i)
	return v.Kind() == reflect.Ptr && v.IsNil()
}

// WrapAndEnhanceError wraps an error and enhances it with version information
func WrapAndEnhanceError(err error, version string) error {
	if err == nil {
		return nil
	}
	wrappedErr := errors.WrapError(err)
	return errors.EnhanceErrorWithVersion(wrappedErr, version)
}

// HandleConversionError creates a standardized conversion error
func HandleConversionError(err error, resourceType string, resourceID interface{}) error {
	conversionErr := errors.NewClientError(
		errors.ErrorCodeServerInternal, 
		"Failed to convert "+resourceType+" data",
	)
	conversionErr.Cause = err
	if resourceID != nil {
		conversionErr.Details = "Error converting "+resourceType+" ID "+formatResourceID(resourceID)
	}
	return conversionErr
}

// formatResourceID formats a resource ID for error messages
func formatResourceID(id interface{}) string {
	if id == nil {
		return "<nil>"
	}
	
	switch v := id.(type) {
	case *int32:
		if v != nil {
			return fmt.Sprintf("%d", *v)
		}
	case *string:
		if v != nil {
			return *v
		}
	case string:
		return v
	case int32:
		return fmt.Sprintf("%d", v)
	case int:
		return fmt.Sprintf("%d", v)
	}
	
	return fmt.Sprintf("%v", id)
}

// CheckClientInitialized verifies the API client is initialized
func CheckClientInitialized(client interface{}) error {
	if client == nil || isNilPointer(client) {
		return errors.NewClientError(
			errors.ErrorCodeClientNotInitialized, 
			"API client not initialized",
		)
	}
	return nil
}

// NewResourceNotFoundError creates a standardized resource not found error
func NewResourceNotFoundError(resourceType string, identifier interface{}) error {
	return errors.NewClientError(
		errors.ErrorCodeResourceNotFound,
		fmt.Sprintf("%s not found", resourceType),
		fmt.Sprintf("%s '%v' was not found", resourceType, identifier),
	)
}

// NewValidationError creates a standardized validation error
func NewValidationError(message, field string, value interface{}) error {
	return errors.NewValidationError(
		errors.ErrorCodeValidationFailed,
		message,
		field,
		value,
		nil,
	)
}
