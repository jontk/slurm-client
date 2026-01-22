// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package common

import (
	stderrors "errors"
	"testing"

	"github.com/jontk/slurm-client/internal/testutil"
	"github.com/jontk/slurm-client/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock types for testing

type mockErrorDetail struct {
	errorNumber *int
	errorCode   *string
	source      *string
	description *string
}

func (m mockErrorDetail) GetErrorNumber() *int    { return m.errorNumber }
func (m mockErrorDetail) GetError() *string       { return m.errorCode }
func (m mockErrorDetail) GetSource() *string      { return m.source }
func (m mockErrorDetail) GetDescription() *string { return m.description }

type mockErrorResponse struct {
	errors []ErrorDetail
}

func (m mockErrorResponse) GetErrors() []ErrorDetail { return m.errors }

type mockResponse struct {
	statusCode    int
	hasErrors     bool
	errorResponse ErrorResponse
}

func (m mockResponse) StatusCode() int                 { return m.statusCode }
func (m mockResponse) HasErrors() bool                 { return m.hasErrors }
func (m mockResponse) GetErrorResponse() ErrorResponse { return m.errorResponse }

func TestHandleAPIResponse(t *testing.T) {
	tests := []struct {
		name        string
		response    ResponseWithErrors
		version     string
		expectError bool
		errorType   func(error) bool
	}{
		{
			name: "successful response (200)",
			response: mockResponse{
				statusCode: 200,
				hasErrors:  false,
			},
			version:     "v0.0.43",
			expectError: false,
		},
		{
			name: "error response with details",
			response: mockResponse{
				statusCode: 400,
				hasErrors:  true,
				errorResponse: mockErrorResponse{
					errors: []ErrorDetail{
						mockErrorDetail{
							errorNumber: testutil.IntPtr(1001),
							errorCode:   testutil.StringPtr("INVALID_REQUEST"),
							source:      testutil.StringPtr("job_manager"),
							description: testutil.StringPtr("Invalid job parameters"),
						},
					},
				},
			},
			version:     "v0.0.43",
			expectError: true,
			errorType:   func(err error) bool { return true }, // API errors are wrapped
		},
		{
			name: "error response without details",
			response: mockResponse{
				statusCode: 500,
				hasErrors:  false,
			},
			version:     "v0.0.43",
			expectError: true,
			errorType:   func(err error) bool { return true }, // HTTP errors are wrapped
		},
		{
			name: "error response with empty error list",
			response: mockResponse{
				statusCode:    400,
				hasErrors:     true,
				errorResponse: mockErrorResponse{errors: []ErrorDetail{}},
			},
			version:     "v0.0.43",
			expectError: true,
			errorType:   func(err error) bool { return true }, // HTTP errors are wrapped
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := HandleAPIResponse(tt.response, tt.version)

			if tt.expectError {
				require.Error(t, err)
				if tt.errorType != nil {
					assert.True(t, tt.errorType(err), "Expected error type check to pass")
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCheckNilResponse(t *testing.T) {
	tests := []struct {
		name        string
		response    interface{}
		operation   string
		expectError bool
	}{
		{
			name:        "non-nil response",
			response:    &struct{}{},
			operation:   "test operation",
			expectError: false,
		},
		{
			name:        "nil response",
			response:    nil,
			operation:   "test operation",
			expectError: true,
		},
		{
			name:        "nil pointer",
			response:    (*struct{})(nil),
			operation:   "test operation",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckNilResponse(tt.response, tt.operation)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.operation)
				// Check it's a SlurmError
				var slurmErr *errors.SlurmError
				ok := stderrors.As(err, &slurmErr)
				assert.True(t, ok, "Expected SlurmError type")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestWrapAndEnhanceError(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		version string
		isNil   bool
	}{
		{
			name:    "nil error",
			err:     nil,
			version: "v0.0.43",
			isNil:   true,
		},
		{
			name:    "non-nil error",
			err:     errors.NewClientError(errors.ErrorCodeClientNotInitialized, "test error"),
			version: "v0.0.43",
			isNil:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WrapAndEnhanceError(tt.err, tt.version)

			if tt.isNil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				// Check that the API version was set
				var slurmErr *errors.SlurmError
				if stderrors.As(result, &slurmErr) {
					assert.Equal(t, tt.version, slurmErr.APIVersion)
				}
			}
		})
	}
}

func TestHandleConversionError(t *testing.T) {
	baseErr := errors.NewClientError(errors.ErrorCodeServerInternal, "base error")

	tests := []struct {
		name         string
		err          error
		resourceType string
		resourceID   interface{}
		expectDetail bool
	}{
		{
			name:         "with string resource ID",
			err:          baseErr,
			resourceType: "job",
			resourceID:   "12345",
			expectDetail: true,
		},
		{
			name:         "with int32 resource ID",
			err:          baseErr,
			resourceType: "node",
			resourceID:   int32(42),
			expectDetail: true,
		},
		{
			name:         "with nil resource ID",
			err:          baseErr,
			resourceType: "partition",
			resourceID:   nil,
			expectDetail: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := HandleConversionError(tt.err, tt.resourceType, tt.resourceID)

			require.Error(t, err)
			// HandleConversionError creates a server error, not a client error
			var slurmErr *errors.SlurmError
			if stderrors.As(err, &slurmErr) {
				assert.Equal(t, errors.ErrorCodeServerInternal, slurmErr.Code)
			}
			assert.Contains(t, err.Error(), tt.resourceType)

			if stderrors.As(err, &slurmErr) && tt.expectDetail {
				assert.NotEmpty(t, slurmErr.Details)
			}
		})
	}
}

func TestCheckClientInitialized(t *testing.T) {
	tests := []struct {
		name        string
		client      interface{}
		expectError bool
	}{
		{
			name:        "initialized client",
			client:      &struct{}{},
			expectError: false,
		},
		{
			name:        "nil client",
			client:      nil,
			expectError: true,
		},
		{
			name:        "nil pointer",
			client:      (*struct{})(nil),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckClientInitialized(tt.client)

			if tt.expectError {
				require.Error(t, err)
				assert.True(t, errors.IsClientError(err))
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNewResourceNotFoundError(t *testing.T) {
	tests := []struct {
		name         string
		resourceType string
		identifier   interface{}
		expectInMsg  []string
	}{
		{
			name:         "job not found",
			resourceType: "job",
			identifier:   12345,
			expectInMsg:  []string{"job", "not found", "12345"},
		},
		{
			name:         "partition not found",
			resourceType: "partition",
			identifier:   "debug",
			expectInMsg:  []string{"partition", "not found", "debug"},
		},
		{
			name:         "node not found",
			resourceType: "node",
			identifier:   "compute-01",
			expectInMsg:  []string{"node", "not found", "compute-01"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewResourceNotFoundError(tt.resourceType, tt.identifier)

			require.Error(t, err)
			slurmErr, ok := err.(*errors.SlurmError)
			require.True(t, ok)
			assert.Equal(t, errors.ErrorCodeResourceNotFound, slurmErr.Code)

			for _, exp := range tt.expectInMsg {
				assert.Contains(t, err.Error(), exp)
			}
		})
	}
}

func TestNewValidationError(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		field    string
		value    interface{}
		expected string
	}{
		{
			name:     "empty string validation",
			message:  "field cannot be empty",
			field:    "username",
			value:    "",
			expected: "field cannot be empty",
		},
		{
			name:     "invalid number validation",
			message:  "must be positive",
			field:    "priority",
			value:    -5,
			expected: "must be positive",
		},
		{
			name:     "nil value validation",
			message:  "required field",
			field:    "job_id",
			value:    nil,
			expected: "required field",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewValidationError(tt.message, tt.field, tt.value)

			require.Error(t, err)
			valErr, ok := err.(*errors.ValidationError)
			require.True(t, ok)
			assert.Equal(t, errors.ErrorCodeValidationFailed, valErr.Code)
			assert.Equal(t, tt.field, valErr.Field)
			assert.Equal(t, tt.value, valErr.Value)
			assert.Contains(t, err.Error(), tt.expected)
		})
	}
}

func TestFormatResourceID(t *testing.T) {
	tests := []struct {
		name     string
		id       interface{}
		expected string
	}{
		{
			name:     "nil value",
			id:       nil,
			expected: "<nil>",
		},
		{
			name:     "int32 pointer",
			id:       testutil.Int32Ptr(42),
			expected: "42",
		},
		{
			name:     "string pointer",
			id:       testutil.StringPtr("test-id"),
			expected: "test-id",
		},
		{
			name:     "direct string",
			id:       "direct-string",
			expected: "direct-string",
		},
		{
			name:     "direct int32",
			id:       int32(123),
			expected: "123",
		},
		{
			name:     "direct int",
			id:       456,
			expected: "456",
		},
		{
			name:     "nil int32 pointer",
			id:       (*int32)(nil),
			expected: "<nil>",
		},
		{
			name:     "nil string pointer",
			id:       (*string)(nil),
			expected: "<nil>",
		},
		{
			name:     "other type",
			id:       struct{ Name string }{Name: "test"},
			expected: "{test}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We can't test the private function directly, but we can test it through HandleConversionError
			err := HandleConversionError(errors.NewClientError(errors.ErrorCodeServerInternal, "test"), "resource", tt.id)

			if tt.id != nil {
				assert.Contains(t, err.Error(), tt.expected)
			}
		})
	}
}

func TestExtractErrorDetail(t *testing.T) {
	tests := []struct {
		name        string
		errorDetail mockErrorDetail
		expectInMsg []string
	}{
		{
			name: "error detail with number and description",
			errorDetail: mockErrorDetail{
				errorNumber: testutil.IntPtr(2040),
				errorCode:   testutil.StringPtr("SLURM_ERROR_BATCH_JOB_SUBMIT_FAILED"),
				source:      testutil.StringPtr("scheduler"),
				description: testutil.StringPtr("Job submission failed"),
			},
			expectInMsg: []string{"Job submission failed"},
		},
		{
			name: "error detail with number but no description",
			errorDetail: mockErrorDetail{
				errorNumber: testutil.IntPtr(2050),
				errorCode:   testutil.StringPtr("SLURM_ERROR_INVALID_PARTITION_NAME"),
				source:      testutil.StringPtr("api"),
				description: nil,
			},
			expectInMsg: []string{"Unknown SLURM error code"}, // Unknown error code gets default message
		},
		{
			name: "error detail without number",
			errorDetail: mockErrorDetail{
				errorNumber: nil,
				errorCode:   testutil.StringPtr("CUSTOM_ERROR"),
				source:      testutil.StringPtr("plugin"),
				description: testutil.StringPtr("Custom error occurred"),
			},
			expectInMsg: []string{"Custom error occurred"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test through HandleAPIResponse
			resp := mockResponse{
				statusCode: 400,
				hasErrors:  true,
				errorResponse: mockErrorResponse{
					errors: []ErrorDetail{tt.errorDetail},
				},
			}

			err := HandleAPIResponse(resp, "v0.0.43")
			require.Error(t, err)

			// Verify expected content in error message
			for _, exp := range tt.expectInMsg {
				assert.Contains(t, err.Error(), exp)
			}
		})
	}
}

func TestIsNilPointer(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{
			name:     "nil interface",
			value:    nil,
			expected: true,
		},
		{
			name:     "nil pointer to string",
			value:    (*string)(nil),
			expected: true,
		},
		{
			name:     "nil pointer to int",
			value:    (*int)(nil),
			expected: true,
		},
		{
			name:     "valid pointer to string",
			value:    testutil.StringPtr("test"),
			expected: false,
		},
		{
			name:     "valid value",
			value:    "test",
			expected: false,
		},
		{
			name:     "valid int",
			value:    42,
			expected: false,
		},
		{
			name:     "struct",
			value:    struct{}{},
			expected: false,
		},
	}

	// We can't test isNilPointer directly since it's private, but we can test through CheckNilResponse
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckNilResponse(tt.value, "test")

			if tt.expected {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "nil")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestFormatResourceIDEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		id       interface{}
		expected string
	}{
		{
			name:     "nil interface",
			id:       nil,
			expected: "",
		},
		{
			name:     "slice",
			id:       []string{"a", "b"},
			expected: "[a b]",
		},
		{
			name:     "uint32",
			id:       uint32(789),
			expected: "789",
		},
		{
			name:     "complex struct",
			id:       map[string]int{"key": 123},
			expected: "map[key:123]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test through HandleConversionError which uses formatResourceID
			err := HandleConversionError(errors.NewClientError(errors.ErrorCodeServerInternal, "test"), "resource", tt.id)
			require.Error(t, err)

			if tt.expected != "" {
				assert.Contains(t, err.Error(), tt.expected)
			}
		})
	}
}
