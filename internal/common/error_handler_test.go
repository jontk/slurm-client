package common

import (
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

func (m mockResponse) StatusCode() int                    { return m.statusCode }
func (m mockResponse) HasErrors() bool                    { return m.hasErrors }
func (m mockResponse) GetErrorResponse() ErrorResponse   { return m.errorResponse }

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
				assert.True(t, errors.IsClientError(err))
				assert.Contains(t, err.Error(), tt.operation)
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
				assert.Contains(t, result.Error(), tt.version)
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
			assert.True(t, errors.IsClientError(err))
			assert.Contains(t, err.Error(), tt.resourceType)
			
			if clientErr, ok := err.(*errors.SlurmError); ok && tt.expectDetail {
				assert.NotEmpty(t, clientErr.Details)
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

