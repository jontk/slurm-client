// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_44

import (
	"github.com/jontk/slurm-client/internal/common"
)

// ResponseAdapter implements common.ResponseWithErrors for v0.0.44 responses
type ResponseAdapter struct {
	statusCode int
	errors     *V0044OpenapiErrors
}

// NewResponseAdapter creates a new response adapter
func NewResponseAdapter(statusCode int, errors *V0044OpenapiErrors) *ResponseAdapter {
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

// GetErrorResponse returns the error response implementation
func (r *ResponseAdapter) GetErrorResponse() common.ErrorResponse {
	if !r.HasErrors() {
		return nil
	}
	return &ErrorResponseImpl{errors: r.errors}
}

// ErrorResponseImpl implements common.ErrorResponse
type ErrorResponseImpl struct {
	errors *V0044OpenapiErrors
}

// GetErrors extracts error details from the response
func (e *ErrorResponseImpl) GetErrors() []common.ErrorDetail {
	if e.errors == nil || len(*e.errors) == 0 {
		return nil
	}

	apiErrors := make([]common.ErrorDetail, len(*e.errors))
	for i, apiErr := range *e.errors {
		apiErrors[i] = &errorDetailImpl{
			errorNumber: func() int {
				if apiErr.ErrorNumber != nil {
					return int(*apiErr.ErrorNumber)
				}
				return 0
			}(),
			errorCode: func() string {
				if apiErr.Error != nil {
					return *apiErr.Error
				}
				return ""
			}(),
			source: func() string {
				if apiErr.Source != nil {
					return *apiErr.Source
				}
				return ""
			}(),
			description: func() string {
				if apiErr.Description != nil {
					return *apiErr.Description
				}
				return ""
			}(),
		}
	}

	return apiErrors
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

// Implement common.ResponseWithErrors interface
var _ common.ResponseWithErrors = (*ResponseAdapter)(nil)
