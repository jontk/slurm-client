// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_42

import (
	"github.com/jontk/slurm-client/internal/common"
)

// ResponseAdapter implements common.ErrorResponse for v0.0.42 responses
type ResponseAdapter struct {
	statusCode int
	errors     *V0042OpenapiErrors
}

// NewResponseAdapter creates a new response adapter
func NewResponseAdapter(statusCode int, errors *V0042OpenapiErrors) *ResponseAdapter {
	return &ResponseAdapter{
		statusCode: statusCode,
		errors:     errors,
	}
}

// StatusCode returns the HTTP status code
func (r *ResponseAdapter) StatusCode() int {
	return r.statusCode
}

// HasErrors checks if the response contains errors
func (r *ResponseAdapter) HasErrors() bool {
	return r.errors != nil && len(*r.errors) > 0
}

// ErrorResponseImpl implements common.ErrorResponse
type ErrorResponseImpl struct {
	errors *V0042OpenapiErrors
}

// errorDetailImpl implements common.ErrorDetail
type errorDetailImpl struct {
	errorNumber int
	errorCode   string
	source      string
	description string
}

func (e *errorDetailImpl) GetErrorNumber() *int {
	if e.errorNumber == 0 {
		return nil
	}
	return &e.errorNumber
}

func (e *errorDetailImpl) GetError() *string {
	if e.errorCode == "" {
		return nil
	}
	return &e.errorCode
}

func (e *errorDetailImpl) GetSource() *string {
	if e.source == "" {
		return nil
	}
	return &e.source
}

func (e *errorDetailImpl) GetDescription() *string {
	if e.description == "" {
		return nil
	}
	return &e.description
}

// GetErrors extracts error details from the response
func (e *ErrorResponseImpl) GetErrors() []common.ErrorDetail {
	if e.errors == nil || len(*e.errors) == 0 {
		return nil
	}

	apiErrors := make([]common.ErrorDetail, len(*e.errors))
	for i, apiErr := range *e.errors {
		var errorNumber int
		if apiErr.ErrorNumber != nil {
			errorNumber = int(*apiErr.ErrorNumber)
		}
		var errorCode string
		if apiErr.Error != nil {
			errorCode = *apiErr.Error
		}
		var source string
		if apiErr.Source != nil {
			source = *apiErr.Source
		}
		var description string
		if apiErr.Description != nil {
			description = *apiErr.Description
		}

		apiErrors[i] = &errorDetailImpl{
			errorNumber: errorNumber,
			errorCode:   errorCode,
			source:      source,
			description: description,
		}
	}

	return apiErrors
}

// GetErrorResponse returns the error response implementation
func (r *ResponseAdapter) GetErrorResponse() common.ErrorResponse {
	if !r.HasErrors() {
		return nil
	}
	return &ErrorResponseImpl{errors: r.errors}
}
