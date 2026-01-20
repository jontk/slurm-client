// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"testing"
)

func TestGetErrorInfo(t *testing.T) {
	tests := []struct {
		name     string
		code     int32
		expected SlurmErrorInfo
	}{
		{
			name: "Success code",
			code: 0,
			expected: SlurmErrorInfo{
				Code:        SlurmSuccess,
				Name:        "SUCCESS",
				Description: "Operation completed successfully",
				Category:    "Success",
			},
		},
		{
			name: "Batch job submit failed",
			code: 2063,
			expected: SlurmErrorInfo{
				Code:        SlurmErrorBatchJobSubmitFailed,
				Name:        "BATCH_JOB_SUBMIT_FAILED",
				Description: "Batch job submission failed - often due to missing required fields like working directory",
				Category:    "Job Submission",
			},
		},
		{
			name: "Invalid partition",
			code: 2001,
			expected: SlurmErrorInfo{
				Code:        SlurmErrorInvalidPartition,
				Name:        "INVALID_PARTITION",
				Description: "The specified partition does not exist or is not available",
				Category:    "Job Submission",
			},
		},
		{
			name: "Account not found",
			code: 4001,
			expected: SlurmErrorInfo{
				Code:        SlurmErrorAccountNotFound,
				Name:        "ACCOUNT_NOT_FOUND",
				Description: "The requested account does not exist",
				Category:    "Account Management",
			},
		},
		{
			name: "Unknown error code",
			code: 99999,
			expected: SlurmErrorInfo{
				Code:        SlurmErrorCode(99999),
				Name:        "UNKNOWN_ERROR",
				Description: "Unknown SLURM error code",
				Category:    "Unknown",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetErrorInfo(tt.code)
			if result == nil {
				t.Fatal("GetErrorInfo returned nil")
			}

			if result.Code != tt.expected.Code {
				t.Errorf("Expected code %v, got %v", tt.expected.Code, result.Code)
			}
			if result.Name != tt.expected.Name {
				t.Errorf("Expected name %s, got %s", tt.expected.Name, result.Name)
			}
			if result.Description != tt.expected.Description {
				t.Errorf("Expected description %s, got %s", tt.expected.Description, result.Description)
			}
			if result.Category != tt.expected.Category {
				t.Errorf("Expected category %s, got %s", tt.expected.Category, result.Category)
			}
		})
	}
}

func TestIsKnownError(t *testing.T) {
	tests := []struct {
		name     string
		code     int32
		expected bool
	}{
		{
			name:     "Known error - success",
			code:     0,
			expected: true,
		},
		{
			name:     "Known error - batch job submit failed",
			code:     2063,
			expected: true,
		},
		{
			name:     "Known error - invalid account",
			code:     2002,
			expected: true, // Now in our map
		},
		{
			name:     "Unknown error",
			code:     99999,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsKnownError(tt.code)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGetErrorCategory(t *testing.T) {
	tests := []struct {
		name     string
		code     int32
		expected string
	}{
		{
			name:     "Success category",
			code:     0,
			expected: "Success",
		},
		{
			name:     "Job submission category",
			code:     2063,
			expected: "Job Submission",
		},
		{
			name:     "Account management category",
			code:     4001,
			expected: "Account Management",
		},
		{
			name:     "Unknown category",
			code:     99999,
			expected: "Unknown",
		},
		{
			name:     "QoS category",
			code:     5001,
			expected: "QoS Management",
		},
		{
			name:     "Reservation category",
			code:     6001,
			expected: "Reservation Management",
		},
		{
			name:     "Authentication category",
			code:     7001,
			expected: "Authentication",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetErrorCategory(tt.code)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestGetErrorDescription(t *testing.T) {
	tests := []struct {
		name     string
		code     int32
		expected string
	}{
		{
			name:     "Success description",
			code:     0,
			expected: "Operation completed successfully",
		},
		{
			name:     "Batch job submit failed description",
			code:     2063,
			expected: "Batch job submission failed - often due to missing required fields like working directory",
		},
		{
			name:     "Unknown error description",
			code:     99999,
			expected: "Unknown SLURM error code",
		},
		{
			name:     "QoS not found description",
			code:     5001,
			expected: "The requested QoS does not exist",
		},
		{
			name:     "Reservation not found description",
			code:     6001,
			expected: "The requested reservation does not exist",
		},
		{
			name:     "Authentication failed description",
			code:     7001,
			expected: "Authentication failed - check credentials or token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetErrorDescription(tt.code)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestEnhanceErrorMessage(t *testing.T) {
	tests := []struct {
		name            string
		errorCode       int32
		originalMessage string
		expected        string
	}{
		{
			name:            "Known error - batch job submit failed",
			errorCode:       2063,
			originalMessage: "Job submission failed",
			expected:        "Batch job submission failed - often due to missing required fields like working directory (SLURM Error BATCH_JOB_SUBMIT_FAILED)",
		},
		{
			name:            "Known error - invalid partition",
			errorCode:       2001,
			originalMessage: "Partition error",
			expected:        "The specified partition does not exist or is not available (SLURM Error INVALID_PARTITION)",
		},
		{
			name:            "Unknown error",
			errorCode:       99999,
			originalMessage: "Something went wrong",
			expected:        "Something went wrong",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EnhanceErrorMessage(tt.errorCode, tt.originalMessage)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestSlurmErrorCodeConstants(t *testing.T) {
	// Test that constants have expected values
	tests := []struct {
		name     string
		code     SlurmErrorCode
		expected int32
	}{
		{"SlurmSuccess", SlurmSuccess, 0},
		{"SlurmErrorBatchJobSubmitFailed", SlurmErrorBatchJobSubmitFailed, 2063},
		{"SlurmErrorInvalidPartition", SlurmErrorInvalidPartition, 2001},
		{"SlurmErrorAccountNotFound", SlurmErrorAccountNotFound, 4001},
		{"SlurmErrorQoSNotFound", SlurmErrorQoSNotFound, 5001},
		{"SlurmErrorReservationNotFound", SlurmErrorReservationNotFound, 6001},
		{"SlurmErrorAuthenticationFailed", SlurmErrorAuthenticationFailed, 7001},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int32(tt.code) != tt.expected {
				t.Errorf("Expected %s to equal %d, got %d", tt.name, tt.expected, int32(tt.code))
			}
		})
	}
}
