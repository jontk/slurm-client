// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_43

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewErrorAdapter(t *testing.T) {
	adapter := NewErrorAdapter()
	assert.NotNil(t, adapter)
}

func TestErrorAdapter_HandleAPIResponse(t *testing.T) {
	adapter := NewErrorAdapter()
	assert.NotNil(t, adapter)

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
			operation:     "TestOperation",
			expectedError: false,
		},
		{
			name:          "client error 400",
			statusCode:    400,
			body:          []byte(`{"error": "Bad Request"}`),
			operation:     "TestOperation",
			expectedError: true,
			expectedInMsg: "400",
		},
		{
			name:          "not found 404",
			statusCode:    404,
			body:          []byte(`{"error": "Not Found"}`),
			operation:     "TestOperation",
			expectedError: true,
			expectedInMsg: "404",
		},
		{
			name:          "server error 500",
			statusCode:    500,
			body:          []byte(`{"error": "Internal Server Error"}`),
			operation:     "TestOperation",
			expectedError: true,
			expectedInMsg: "500",
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

func TestErrorAdapter_HandleAPIError(t *testing.T) {
	adapter := NewErrorAdapter()

	// ErrorAdapter doesn't have HandleAPIError method
	// Instead test the HandleAPIResponse method with standard error
	result := adapter.HandleAPIResponse(500, []byte(`{"error": "Internal Server Error"}`), "TestOperation")
	assert.Error(t, result)
	assert.Contains(t, result.Error(), "500")
}
