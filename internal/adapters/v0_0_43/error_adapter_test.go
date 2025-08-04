// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_43

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	api "github.com/jontk/slurm-client/internal/api/v0_0_43"
)

func TestNewErrorAdapter(t *testing.T) {
	adapter := NewErrorAdapter()
	assert.NotNil(t, adapter)
	assert.NotNil(t, adapter.errorMappings)
	assert.Greater(t, len(adapter.errorMappings), 0)
}

func TestErrorAdapter_HandleAPIResponse(t *testing.T) {
	adapter := NewErrorAdapter()

	tests := []struct {
		name           string
		statusCode     int
		body           []byte
		operation      string
		expectedError  bool
		expectedInMsg  string
	}{
		{
			name:          "success response",
			statusCode:    200,
			body:          nil,
			operation:     "TestOperation",
			expectedError: false,
		},
		{
			name:          "client error",
			statusCode:    400,
			body:          []byte(`{"errors":[{"error":"Bad request"}]}`),
			operation:     "TestOperation",
			expectedError: true,
			expectedInMsg: "Bad request",
		},
		{
			name:          "server error",
			statusCode:    500,
			body:          []byte(`{"errors":[{"error":"Internal server error"}]}`),
			operation:     "TestOperation",
			expectedError: true,
			expectedInMsg: "Internal server error",
		},
		{
			name:          "unauthorized",
			statusCode:    401,
			body:          []byte(`{"errors":[{"error":"Unauthorized"}]}`),
			operation:     "TestOperation",
			expectedError: true,
			expectedInMsg: "Unauthorized",
		},
		{
			name:          "not found",
			statusCode:    404,
			body:          []byte(`{"errors":[{"error":"Not found"}]}`),
			operation:     "TestOperation",
			expectedError: true,
			expectedInMsg: "Not found",
		},
		{
			name:          "rate limit",
			statusCode:    429,
			body:          []byte(`{"errors":[{"error":"Too many requests"}]}`),
			operation:     "TestOperation",
			expectedError: true,
			expectedInMsg: "Too many requests",
		},
		{
			name:          "empty body error",
			statusCode:    400,
			body:          []byte{},
			operation:     "TestOperation",
			expectedError: true,
			expectedInMsg: "status 400",
		},
		{
			name:          "invalid json",
			statusCode:    400,
			body:          []byte(`{invalid json`),
			operation:     "TestOperation",
			expectedError: true,
			expectedInMsg: "status 400",
		},
		{
			name:          "created response",
			statusCode:    201,
			body:          nil,
			operation:     "Create",
			expectedError: false,
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

func TestErrorAdapter_ParseErrorResponse(t *testing.T) {
	adapter := NewErrorAdapter()

	tests := []struct {
		name           string
		body           []byte
		expectedErrors []string
		expectedEmpty  bool
	}{
		{
			name: "single error",
			body: []byte(`{"errors":[{"error":"Test error"}]}`),
			expectedErrors: []string{"Test error"},
		},
		{
			name: "multiple errors",
			body: []byte(`{"errors":[{"error":"Error 1"},{"error":"Error 2"}]}`),
			expectedErrors: []string{"Error 1", "Error 2"},
		},
		{
			name: "error with description",
			body: []byte(`{"errors":[{"error":"Test error","description":"Detailed description"}]}`),
			expectedErrors: []string{"Test error: Detailed description"},
		},
		{
			name: "error without message",
			body: []byte(`{"errors":[{}]}`),
			expectedErrors: []string{"Unknown error"},
		},
		{
			name:          "empty errors array",
			body:          []byte(`{"errors":[]}`),
			expectedEmpty: true,
		},
		{
			name:          "invalid json",
			body:          []byte(`{invalid}`),
			expectedEmpty: true,
		},
		{
			name:          "nil body",
			body:          nil,
			expectedEmpty: true,
		},
		{
			name:          "empty body",
			body:          []byte{},
			expectedEmpty: true,
		},
		{
			name: "nested error structure",
			body: []byte(`{"errors":[{"error":"Main error","nested":{"detail":"More info"}}]}`),
			expectedErrors: []string{"Main error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := adapter.ParseErrorResponse(tt.body)

			if tt.expectedEmpty {
				assert.Empty(t, errors)
			} else {
				assert.Equal(t, len(tt.expectedErrors), len(errors))
				for i, expected := range tt.expectedErrors {
					assert.Equal(t, expected, errors[i])
				}
			}
		})
	}
}

func TestErrorAdapter_HandleHTTPResponse(t *testing.T) {
	adapter := NewErrorAdapter()

	tests := []struct {
		name          string
		resp          *http.Response
		body          []byte
		expectedError bool
		expectedInMsg string
	}{
		{
			name: "nil response",
			resp: nil,
			body: nil,
			expectedError: true,
			expectedInMsg: "nil HTTP response",
		},
		{
			name: "success response",
			resp: &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
			},
			body:          nil,
			expectedError: false,
		},
		{
			name: "error with body",
			resp: &http.Response{
				StatusCode: 400,
				Status:     "400 Bad Request",
			},
			body:          []byte(`{"errors":[{"error":"Invalid parameter"}]}`),
			expectedError: true,
			expectedInMsg: "Invalid parameter",
		},
		{
			name: "error without body",
			resp: &http.Response{
				StatusCode: 500,
				Status:     "500 Internal Server Error",
			},
			body:          nil,
			expectedError: true,
			expectedInMsg: "500 Internal Server Error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.HandleHTTPResponse(tt.resp, tt.body)

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

func TestErrorAdapter_MapErrorCode(t *testing.T) {
	adapter := NewErrorAdapter()

	tests := []struct {
		name         string
		errorCode    string
		expectedType string
		expectedMsg  string
		shouldExist  bool
	}{
		{
			name:         "invalid job id",
			errorCode:    "INVALID_JOB_ID",
			expectedType: "ValidationError",
			shouldExist:  true,
		},
		{
			name:         "access denied",
			errorCode:    "ACCESS_DENIED",
			expectedType: "AuthenticationError",
			shouldExist:  true,
		},
		{
			name:         "unknown error code",
			errorCode:    "UNKNOWN_ERROR_CODE_XYZ",
			shouldExist:  false,
		},
		{
			name:         "empty error code",
			errorCode:    "",
			shouldExist:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errType, msg := adapter.MapErrorCode(tt.errorCode)

			if tt.shouldExist {
				assert.NotEmpty(t, errType)
				assert.NotEmpty(t, msg)
				if tt.expectedType != "" {
					assert.Equal(t, tt.expectedType, errType)
				}
			} else {
				assert.Equal(t, "UnknownError", errType)
				assert.Equal(t, fmt.Sprintf("Unknown error code: %s", tt.errorCode), msg)
			}
		})
	}
}

func TestErrorAdapter_EnhanceError(t *testing.T) {
	adapter := NewErrorAdapter()

	tests := []struct {
		name        string
		err         error
		operation   string
		context     map[string]interface{}
		expectedMsg string
	}{
		{
			name:        "nil error",
			err:         nil,
			operation:   "TestOp",
			context:     nil,
			expectedMsg: "",
		},
		{
			name:        "simple error",
			err:         fmt.Errorf("test error"),
			operation:   "CreateJob",
			context:     nil,
			expectedMsg: "CreateJob failed: test error",
		},
		{
			name:      "error with context",
			err:       fmt.Errorf("validation failed"),
			operation: "UpdateUser",
			context: map[string]interface{}{
				"user": "testuser",
				"field": "email",
			},
			expectedMsg: "UpdateUser failed: validation failed",
		},
		{
			name:        "empty operation",
			err:         fmt.Errorf("some error"),
			operation:   "",
			context:     nil,
			expectedMsg: "some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := adapter.EnhanceError(tt.err, tt.operation, tt.context)

			if tt.err == nil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Contains(t, result.Error(), tt.expectedMsg)
			}
		})
	}
}

func TestErrorAdapter_ErrorResponseHandling(t *testing.T) {
	adapter := NewErrorAdapter()

	// Test with actual V0043OpenapiError structure
	errorResp := api.V0043OpenapiResp{
		Errors: &[]api.V0043OpenapiError{
			{
				Error:       ptrString("Test error"),
				Description: ptrString("Test description"),
			},
			{
				Error: ptrString("Second error"),
			},
		},
	}

	// Marshal to JSON to test parsing
	body, err := json.Marshal(errorResp)
	assert.NoError(t, err)

	errors := adapter.ParseErrorResponse(body)
	assert.Len(t, errors, 2)
	assert.Equal(t, "Test error: Test description", errors[0])
	assert.Equal(t, "Second error", errors[1])
}

func TestErrorAdapter_ConcurrentAccess(t *testing.T) {
	adapter := NewErrorAdapter()

	// Test concurrent access to error mappings
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			_, _ = adapter.MapErrorCode("INVALID_JOB_ID")
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	// If we get here without panic, concurrent access is safe
	assert.True(t, true)
}