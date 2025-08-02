// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package errors

import (
	"reflect"
	"strings"
	"testing"
)


func TestGetVersionMapping(t *testing.T) {
	tests := []struct {
		version string
		isValid bool
	}{
		{"v0.0.40", true},
		{"v0.0.41", true},
		{"v0.0.42", true},
		{"v0.0.43", true},
		{"unknown", true}, // Should fall back to latest
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			mapping := GetVersionMapping(tt.version)

			if mapping == nil {
				t.Fatal("Expected non-nil mapping")
			}

			if mapping.SlurmErrorMappings == nil {
				t.Error("Expected non-nil SlurmErrorMappings")
			}

			if mapping.HTTPStatusMappings == nil {
				t.Error("Expected non-nil HTTPStatusMappings")
			}

			if mapping.FeatureSupport == nil {
				t.Error("Expected non-nil FeatureSupport")
			}

			// Check that basic error mappings exist
			if _, exists := mapping.SlurmErrorMappings["SLURM_AUTHENTICATION_ERROR"]; !exists {
				t.Error("Expected SLURM_AUTHENTICATION_ERROR mapping")
			}

			if _, exists := mapping.HTTPStatusMappings[401]; !exists {
				t.Error("Expected 401 HTTP status mapping")
			}
		})
	}
}

func TestVersionSpecificMappings(t *testing.T) {
	// Test v0.0.40 specific features
	v040 := GetVersionMapping("v0.0.40")
	if !v040.FeatureSupport["frontend_mode"] {
		t.Error("Expected frontend_mode to be supported in v0.0.40")
	}
	if v040.FeatureSupport["reservation_support"] {
		t.Error("Expected reservation_support to not be supported in v0.0.40")
	}

	// Test v0.0.41 specific features
	v041 := GetVersionMapping("v0.0.41")
	if !v041.FeatureSupport["reservation_support"] {
		t.Error("Expected reservation_support to be supported in v0.0.41")
	}
	if !v041.FeatureSupport["enhanced_job_deps"] {
		t.Error("Expected enhanced_job_deps to be supported in v0.0.41")
	}

	// Test v0.0.42 specific features (field removals)
	v042 := GetVersionMapping("v0.0.42")
	if v042.FeatureSupport["exclusive_field"] {
		t.Error("Expected exclusive_field to not be supported in v0.0.42")
	}
	if v042.FeatureSupport["oversubscribe_field"] {
		t.Error("Expected oversubscribe_field to not be supported in v0.0.42")
	}
	if !v042.FeatureSupport["strict_validation"] {
		t.Error("Expected strict_validation to be supported in v0.0.42")
	}

	// Test v0.0.43 specific features
	v043 := GetVersionMapping("v0.0.43")
	if v043.FeatureSupport["frontend_mode"] {
		t.Error("Expected frontend_mode to not be supported in v0.0.43")
	}
	if !v043.FeatureSupport["advanced_reservations"] {
		t.Error("Expected advanced_reservations to be supported in v0.0.43")
	}
	if !v043.FeatureSupport["multi_cluster"] {
		t.Error("Expected multi_cluster to be supported in v0.0.43")
	}
}

func TestMapSlurmErrorForVersion(t *testing.T) {
	tests := []struct {
		name           string
		slurmErrorCode string
		apiVersion     string
		httpStatusCode int
		expectedCode   ErrorCode
	}{
		{
			name:           "known slurm error",
			slurmErrorCode: "SLURM_AUTHENTICATION_ERROR",
			apiVersion:     "v0.0.42",
			httpStatusCode: 401,
			expectedCode:   ErrorCodeInvalidCredentials,
		},
		{
			name:           "unknown slurm error falls back to HTTP",
			slurmErrorCode: "UNKNOWN_ERROR",
			apiVersion:     "v0.0.42",
			httpStatusCode: 500,
			expectedCode:   ErrorCodeServerInternal,
		},
		{
			name:           "version specific error v0.0.41",
			slurmErrorCode: "SLURM_RESERVATION_INVALID",
			apiVersion:     "v0.0.41",
			httpStatusCode: 404,
			expectedCode:   ErrorCodeResourceNotFound,
		},
		{
			name:           "version specific error v0.0.43",
			slurmErrorCode: "SLURM_ADVANCED_RESERVATION_ERROR",
			apiVersion:     "v0.0.43",
			httpStatusCode: 404,
			expectedCode:   ErrorCodeResourceNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MapSlurmErrorForVersion(tt.slurmErrorCode, tt.apiVersion, tt.httpStatusCode)
			if result != tt.expectedCode {
				t.Errorf("Expected error code %v, got %v", tt.expectedCode, result)
			}
		})
	}
}

func TestIsFeatureSupportedInVersion(t *testing.T) {
	tests := []struct {
		feature   string
		version   string
		supported bool
	}{
		{"frontend_mode", "v0.0.40", true},
		{"frontend_mode", "v0.0.43", false},
		{"reservation_support", "v0.0.40", false},
		{"reservation_support", "v0.0.41", true},
		{"exclusive_field", "v0.0.41", true},
		{"exclusive_field", "v0.0.42", false},
		{"advanced_reservations", "v0.0.42", false},
		{"advanced_reservations", "v0.0.43", true},
		{"unknown_feature", "v0.0.42", false},
	}

	for _, tt := range tests {
		t.Run(tt.feature+"_"+tt.version, func(t *testing.T) {
			result := IsFeatureSupportedInVersion(tt.feature, tt.version)
			if result != tt.supported {
				t.Errorf("Expected feature %s in version %s to be %v, got %v",
					tt.feature, tt.version, tt.supported, result)
			}
		})
	}
}

func TestGetBreakingChanges(t *testing.T) {
	tests := []struct {
		name        string
		fromVersion string
		toVersion   string
		expectEmpty bool
	}{
		{
			name:        "v0.0.40 to v0.0.41",
			fromVersion: "v0.0.40",
			toVersion:   "v0.0.41",
			expectEmpty: false,
		},
		{
			name:        "v0.0.41 to v0.0.42",
			fromVersion: "v0.0.41",
			toVersion:   "v0.0.42",
			expectEmpty: false,
		},
		{
			name:        "v0.0.42 to v0.0.43",
			fromVersion: "v0.0.42",
			toVersion:   "v0.0.43",
			expectEmpty: false,
		},
		{
			name:        "v0.0.40 to v0.0.43 (multiple versions)",
			fromVersion: "v0.0.40",
			toVersion:   "v0.0.43",
			expectEmpty: false,
		},
		{
			name:        "same version",
			fromVersion: "v0.0.42",
			toVersion:   "v0.0.42",
			expectEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			changes := GetBreakingChanges(tt.fromVersion, tt.toVersion)

			if tt.expectEmpty && len(changes) > 0 {
				t.Errorf("Expected no breaking changes, got %v", changes)
			}

			if !tt.expectEmpty && len(changes) == 0 {
				t.Error("Expected breaking changes, got none")
			}
		})
	}

	// Test specific breaking changes
	changes := GetBreakingChanges("v0.0.40", "v0.0.41")
	found := false
	for _, change := range changes {
		if contains(change, "minimum_switches") && contains(change, "required_switches") {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected field rename breaking change from v0.0.40 to v0.0.41")
	}

	changes = GetBreakingChanges("v0.0.41", "v0.0.42")
	found = false
	for _, change := range changes {
		if contains(change, "exclusive") && contains(change, "oversubscribe") {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected field removal breaking change from v0.0.41 to v0.0.42")
	}

	changes = GetBreakingChanges("v0.0.42", "v0.0.43")
	found = false
	for _, change := range changes {
		if contains(change, "FrontEnd") {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected FrontEnd mode removal breaking change from v0.0.42 to v0.0.43")
	}
}

func TestValidateVersionCompatibility(t *testing.T) {
	tests := []struct {
		name          string
		clientVersion string
		serverVersion string
		expectError   bool
	}{
		{
			name:          "same version",
			clientVersion: "v0.0.42",
			serverVersion: "v0.0.42",
			expectError:   false,
		},
		{
			name:          "client older than server with breaking changes",
			clientVersion: "v0.0.40",
			serverVersion: "v0.0.42",
			expectError:   true,
		},
		{
			name:          "unsupported client version",
			clientVersion: "v0.0.39",
			serverVersion: "v0.0.42",
			expectError:   true,
		},
		{
			name:          "unsupported server version",
			clientVersion: "v0.0.42",
			serverVersion: "v0.0.39",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateVersionCompatibility(tt.clientVersion, tt.serverVersion)

			if tt.expectError && err == nil {
				t.Error("Expected compatibility error, got nil")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Expected no compatibility error, got %v", err)
			}

			if err != nil {
				// Verify it's a SlurmError with appropriate code
				if slurmErr, ok := err.(*SlurmError); ok {
					if slurmErr.Code != ErrorCodeVersionMismatch {
						t.Errorf("Expected ErrorCodeVersionMismatch, got %v", slurmErr.Code)
					}
				} else {
					t.Errorf("Expected SlurmError, got %T", err)
				}
			}
		})
	}
}

func TestEnhanceErrorWithVersion(t *testing.T) {
	tests := []struct {
		name       string
		err        *SlurmError
		apiVersion string
		expectMore bool
	}{
		{
			name:       "nil error",
			err:        nil,
			apiVersion: "v0.0.42",
			expectMore: false,
		},
		{
			name:       "unsupported operation in v0.0.43",
			err:        NewSlurmError(ErrorCodeUnsupportedOperation, "unsupported"),
			apiVersion: "v0.0.43",
			expectMore: true,
		},
		{
			name:       "validation failed in v0.0.42",
			err:        NewSlurmError(ErrorCodeValidationFailed, "validation error"),
			apiVersion: "v0.0.42",
			expectMore: true,
		},
		{
			name:       "version mismatch",
			err:        NewSlurmError(ErrorCodeVersionMismatch, "version issue"),
			apiVersion: "v0.0.42",
			expectMore: true,
		},
		{
			name:       "network error (no enhancement)",
			err:        NewSlurmError(ErrorCodeNetworkTimeout, "timeout"),
			apiVersion: "v0.0.42",
			expectMore: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalDetails := ""
			if tt.err != nil {
				originalDetails = tt.err.Details
			}

			result := EnhanceErrorWithVersion(tt.err, tt.apiVersion)

			if tt.err == nil {
				if result != nil {
					t.Error("Expected nil result for nil error")
				}
				return
			}

			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			if result.APIVersion != tt.apiVersion {
				t.Errorf("Expected API version %s, got %s", tt.apiVersion, result.APIVersion)
			}

			if tt.expectMore {
				if result.Details == originalDetails {
					t.Error("Expected enhanced details, but details unchanged")
				}
			}
		})
	}
}

func TestCommonHTTPMappings(t *testing.T) {
	mappings := getCommonHTTPMappings()

	expectedMappings := map[int]ErrorCode{
		200: ErrorCodeUnknown,
		400: ErrorCodeInvalidRequest,
		401: ErrorCodeUnauthorized,
		403: ErrorCodePermissionDenied,
		404: ErrorCodeResourceNotFound,
		409: ErrorCodeConflict,
		422: ErrorCodeValidationFailed,
		429: ErrorCodeRateLimited,
		500: ErrorCodeServerInternal,
		502: ErrorCodeSlurmDaemonDown,
		503: ErrorCodeSlurmDaemonDown,
		504: ErrorCodeNetworkTimeout,
	}

	if !reflect.DeepEqual(mappings, expectedMappings) {
		t.Errorf("HTTP mappings mismatch.\nExpected: %v\nGot: %v", expectedMappings, mappings)
	}
}

func TestVersionSpecificErrorMappings(t *testing.T) {
	// Test that newer versions include mappings from older versions
	v040 := GetVersionMapping("v0.0.40")
	v043 := GetVersionMapping("v0.0.43")

	// Check that common errors exist in both
	commonErrors := []string{
		"SLURM_AUTHENTICATION_ERROR",
		"SLURM_ACCESS_DENIED",
		"SLURM_INVALID_JOB_ID",
		"SLURM_INVALID_PARTITION_NAME",
	}

	for _, errorCode := range commonErrors {
		if _, exists := v040.SlurmErrorMappings[errorCode]; !exists {
			t.Errorf("Expected error code %s in v0.0.40", errorCode)
		}
		if _, exists := v043.SlurmErrorMappings[errorCode]; !exists {
			t.Errorf("Expected error code %s in v0.0.43", errorCode)
		}
	}

	// Check that v0.0.43 has additional errors not in v0.0.40
	v043SpecificErrors := []string{
		"SLURM_ADVANCED_RESERVATION_ERROR",
		"SLURM_MULTI_CLUSTER_ERROR",
		"SLURM_FEDERATION_ERROR",
	}

	for _, errorCode := range v043SpecificErrors {
		if _, exists := v040.SlurmErrorMappings[errorCode]; exists {
			t.Errorf("Did not expect error code %s in v0.0.40", errorCode)
		}
		if _, exists := v043.SlurmErrorMappings[errorCode]; !exists {
			t.Errorf("Expected error code %s in v0.0.43", errorCode)
		}
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
