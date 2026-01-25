// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_42

import (
	"encoding/json"
	stderrors "errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	api "github.com/jontk/slurm-client/internal/api/v0_0_42"
	"github.com/jontk/slurm-client/pkg/errors"
)

func TestNewErrorAdapter(t *testing.T) {
	adapter := NewErrorAdapter()
	assert.NotNil(t, adapter)
}

func TestErrorAdapter_HandleAPIResponse(t *testing.T) {
	adapter := NewErrorAdapter()

	tests := []struct {
		name          string
		statusCode    int
		body          []byte
		operation     string
		expectedError bool
		expectedInMsg string
	}{
		{
			name:          "success response 200",
			statusCode:    200,
			body:          nil,
			operation:     "TestOperation",
			expectedError: false,
		},
		{
			name:          "success response 201",
			statusCode:    201,
			body:          nil,
			operation:     "Create",
			expectedError: false,
		},
		{
			name:          "success response 204",
			statusCode:    204,
			body:          nil,
			operation:     "Delete",
			expectedError: false,
		},
		{
			name:       "client error with API error response",
			statusCode: 400,
			body: func() []byte {
				resp := struct {
					Errors []api.V0042OpenapiError `json:"errors"`
				}{
					Errors: []api.V0042OpenapiError{
						{
							ErrorNumber: ptrInt32(123),
							Error:       ptrString("INVALID_REQUEST"),
							Description: ptrString("Invalid request parameters"),
							Source:      ptrString("slurmrestd"),
						},
					},
				}
				b, _ := json.Marshal(resp)
				return b
			}(),
			operation:     "TestOperation",
			expectedError: true,
			expectedInMsg: "INVALID_REQUEST",
		},
		{
			name:       "multiple errors in response",
			statusCode: 400,
			body: func() []byte {
				resp := struct {
					Errors []api.V0042OpenapiError `json:"errors"`
				}{
					Errors: []api.V0042OpenapiError{
						{
							Error:       ptrString("VALIDATION_ERROR"),
							Description: ptrString("Field validation failed"),
						},
						{
							Error:       ptrString("PERMISSION_ERROR"),
							Description: ptrString("Insufficient permissions"),
						},
					},
				}
				b, _ := json.Marshal(resp)
				return b
			}(),
			operation:     "TestOperation",
			expectedError: true,
			expectedInMsg: "VALIDATION_ERROR",
		},
		{
			name:          "unauthorized 401",
			statusCode:    401,
			body:          []byte("Unauthorized"),
			operation:     "TestOperation",
			expectedError: true,
			expectedInMsg: "Authentication failed",
		},
		{
			name:          "forbidden 403",
			statusCode:    403,
			body:          []byte("Forbidden"),
			operation:     "TestOperation",
			expectedError: true,
			expectedInMsg: "Permission denied",
		},
		{
			name:          "not found 404",
			statusCode:    404,
			body:          []byte("Not found"),
			operation:     "GetResource",
			expectedError: true,
			expectedInMsg: "GetResource: resource not found",
		},
		{
			name:          "conflict 409",
			statusCode:    409,
			body:          []byte("Resource exists"),
			operation:     "CreateResource",
			expectedError: true,
			expectedInMsg: "CreateResource: resource conflict",
		},
		{
			name:          "validation error 422",
			statusCode:    422,
			body:          []byte("Invalid data"),
			operation:     "UpdateResource",
			expectedError: true,
			expectedInMsg: "UpdateResource: validation failed",
		},
		{
			name:          "internal server error 500",
			statusCode:    500,
			body:          []byte("Internal error"),
			operation:     "TestOperation",
			expectedError: true,
			expectedInMsg: "TestOperation: server error",
		},
		{
			name:          "bad gateway 502",
			statusCode:    502,
			body:          []byte("Bad gateway"),
			operation:     "TestOperation",
			expectedError: true,
			expectedInMsg: "TestOperation: server error",
		},
		{
			name:          "service unavailable 503",
			statusCode:    503,
			body:          []byte("Service unavailable"),
			operation:     "TestOperation",
			expectedError: true,
			expectedInMsg: "TestOperation: server error",
		},
		{
			name:          "unknown error code",
			statusCode:    418,
			body:          []byte("I'm a teapot"),
			operation:     "TestOperation",
			expectedError: true,
			expectedInMsg: "TestOperation failed with status 418",
		},
		{
			name:          "empty body error",
			statusCode:    400,
			body:          []byte{},
			operation:     "TestOperation",
			expectedError: true,
			expectedInMsg: "TestOperation failed with status 400",
		},
		{
			name:          "invalid json body",
			statusCode:    400,
			body:          []byte(`{invalid json`),
			operation:     "TestOperation",
			expectedError: true,
			expectedInMsg: "TestOperation failed with status 400",
		},
		{
			name:       "error without description uses error field",
			statusCode: 400,
			body: func() []byte {
				resp := struct {
					Errors []api.V0042OpenapiError `json:"errors"`
				}{
					Errors: []api.V0042OpenapiError{
						{
							Error: ptrString("SIMPLE_ERROR"),
							// No Description field
						},
					},
				}
				b, _ := json.Marshal(resp)
				return b
			}(),
			operation:     "TestOperation",
			expectedError: true,
			expectedInMsg: "SIMPLE_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.HandleAPIResponse(tt.statusCode, tt.body, tt.operation)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.expectedInMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedInMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestErrorAdapter_ParseSlurmError(t *testing.T) {
	adapter := NewErrorAdapter()

	tests := []struct {
		name          string
		inputError    error
		expectedCode  string
		expectedMsg   string
		expectedErrno int
	}{
		{
			name:          "generic error",
			inputError:    assert.AnError,
			expectedCode:  "UNKNOWN",
			expectedMsg:   assert.AnError.Error(),
			expectedErrno: -1,
		},
		{
			name: "SLURM API error with details",
			inputError: &errors.SlurmAPIError{
				SlurmError: &errors.SlurmError{
					StatusCode: 400,
				},
				Source: "slurmrestd",
				Errors: []errors.SlurmAPIErrorDetail{
					{
						ErrorNumber: 123,
						ErrorCode:   "INVALID_JOB_ID",
						Description: "Job ID is invalid",
					},
				},
			},
			expectedCode:  "INVALID_JOB_ID",
			expectedMsg:   "Job ID is invalid",
			expectedErrno: 123,
		},
		{
			name: "SLURM API error with multiple errors (returns first)",
			inputError: &errors.SlurmAPIError{
				SlurmError: &errors.SlurmError{
					StatusCode: 400,
				},
				Source: "slurmrestd",
				Errors: []errors.SlurmAPIErrorDetail{
					{
						ErrorNumber: 100,
						ErrorCode:   "FIRST_ERROR",
						Description: "First error message",
					},
					{
						ErrorNumber: 200,
						ErrorCode:   "SECOND_ERROR",
						Description: "Second error message",
					},
				},
			},
			expectedCode:  "FIRST_ERROR",
			expectedMsg:   "First error message",
			expectedErrno: 100,
		},
		{
			name: "SLURM API error with no details",
			inputError: &errors.SlurmAPIError{
				SlurmError: &errors.SlurmError{
					Message:    "Generic API error",
					StatusCode: 500,
				},
				Source: "slurmrestd",
				Errors: []errors.SlurmAPIErrorDetail{},
			},
			expectedCode:  "UNKNOWN",
			expectedMsg:   "Generic API error",
			expectedErrno: -1,
		},
		{
			name: "SLURM API error with empty error details",
			inputError: &errors.SlurmAPIError{
				SlurmError: &errors.SlurmError{
					Message:    "Empty details error",
					StatusCode: 400,
				},
				Source: "slurmrestd",
				Errors: []errors.SlurmAPIErrorDetail{
					{
						// Empty details
					},
				},
			},
			expectedCode:  "",
			expectedMsg:   "",
			expectedErrno: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, message, errno := adapter.ParseSlurmError(tt.inputError)

			assert.Equal(t, tt.expectedCode, code)
			assert.Equal(t, tt.expectedMsg, message)
			assert.Equal(t, tt.expectedErrno, errno)
		})
	}
}

func TestErrorAdapter_HandleAPIResponse_ErrorTypesValidation(t *testing.T) {
	adapter := NewErrorAdapter()

	tests := []struct {
		name         string
		statusCode   int
		body         []byte
		operation    string
		expectedType string
	}{
		{
			name:         "authentication error type for 401",
			statusCode:   401,
			body:         []byte("auth error"),
			operation:    "test",
			expectedType: "*errors.AuthenticationError",
		},
		{
			name:         "authentication error type for 403",
			statusCode:   403,
			body:         []byte("forbidden"),
			operation:    "test",
			expectedType: "*errors.AuthenticationError",
		},
		{
			name:         "slurm error type for 404",
			statusCode:   404,
			body:         []byte("not found"),
			operation:    "test",
			expectedType: "*errors.SlurmError",
		},
		{
			name:         "slurm error type for 409",
			statusCode:   409,
			body:         []byte("conflict"),
			operation:    "test",
			expectedType: "*errors.SlurmError",
		},
		{
			name:         "validation error type for 422",
			statusCode:   422,
			body:         []byte("validation failed"),
			operation:    "test",
			expectedType: "*errors.ValidationError",
		},
		{
			name:         "slurm error type for 500",
			statusCode:   500,
			body:         []byte("server error"),
			operation:    "test",
			expectedType: "*errors.SlurmError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.HandleAPIResponse(tt.statusCode, tt.body, tt.operation)
			assert.Error(t, err)

			// Check the specific error type
			switch tt.expectedType {
			case "*errors.AuthenticationError":
				var authErr *errors.AuthenticationError
				assert.True(t, stderrors.As(err, &authErr), "Expected AuthenticationError but got %T", err)
			case "*errors.ValidationError":
				var valErr *errors.ValidationError
				assert.True(t, stderrors.As(err, &valErr), "Expected ValidationError but got %T", err)
			case "*errors.SlurmError":
				var slErr *errors.SlurmError
				assert.True(t, stderrors.As(err, &slErr), "Expected SlurmError but got %T", err)
			}
		})
	}
}

func TestErrorAdapter_HandleAPIResponse_StatusCodeRanges(t *testing.T) {
	adapter := NewErrorAdapter()

	successCodes := []int{200, 201, 202, 204, 299}
	for _, code := range successCodes {
		t.Run(fmt.Sprintf("success_code_%d", code), func(t *testing.T) {
			err := adapter.HandleAPIResponse(code, nil, "test")
			assert.NoError(t, err)
		})
	}

	errorCodes := []int{300, 400, 401, 403, 404, 409, 422, 500, 502, 503, 600}
	for _, code := range errorCodes {
		t.Run(fmt.Sprintf("error_code_%d", code), func(t *testing.T) {
			err := adapter.HandleAPIResponse(code, []byte("error"), "test")
			assert.Error(t, err)
		})
	}
}

func TestErrorAdapter_HandleAPIResponse_ConcurrentAccess(t *testing.T) {
	adapter := NewErrorAdapter()

	// Test concurrent access to error handling (should be safe)
	done := make(chan bool)
	for i := range 10 {
		go func(index int) {
			_ = adapter.HandleAPIResponse(400, []byte("test"), fmt.Sprintf("operation_%d", index))
			done <- true
		}(i)
	}

	for range 10 {
		<-done
	}

	// If we get here without panic, concurrent access is safe
	assert.True(t, true)
}
