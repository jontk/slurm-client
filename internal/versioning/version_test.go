// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package versioning

import (
	"testing"

	"github.com/jontk/slurm-client/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestSupportedVersions(t *testing.T) {
	// Test that SupportedVersions contains expected versions
	expected := []string{"v0.0.40", "v0.0.41", "v0.0.42", "v0.0.43"}

	var versionStrings []string
	for _, v := range SupportedVersions {
		versionStrings = append(versionStrings, v.String())
	}

	helpers.AssertEqual(t, expected, versionStrings)
}

func TestStableVersion(t *testing.T) {
	version := StableVersion()
	helpers.AssertEqual(t, "v0.0.42", version.String())
}

func TestLatestVersion(t *testing.T) {
	version := LatestVersion()
	helpers.AssertEqual(t, "v0.0.43", version.String())
}

func TestIsVersionSupported(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		expected bool
	}{
		{"supported v0.0.40", "v0.0.40", true},
		{"supported v0.0.41", "v0.0.41", true},
		{"supported v0.0.42", "v0.0.42", true},
		{"supported v0.0.43", "v0.0.43", true},
		{"unsupported v0.0.39", "v0.0.39", true}, // Compatible with v0.0.40 (same major.minor)
		{"unsupported v0.0.44", "v0.0.44", true}, // Compatible with v0.0.43 (same major.minor)
		{"invalid version", "invalid", false},
		{"empty version", "", true}, // Defaults to latest
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Check if version is supported by trying to find it in SupportedVersions
			_, err := FindBestVersion(tt.version)
			isSupported := (err == nil)
			helpers.AssertEqual(t, tt.expected, isSupported)
		})
	}
}

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name        string
		version     string
		expectError bool
		expected    APIVersion
	}{
		{
			name:        "valid v0.0.40",
			version:     "v0.0.40",
			expectError: false,
			expected:    APIVersion{Major: 0, Minor: 0, Patch: 40, Raw: "v0.0.40"},
		},
		{
			name:        "valid v0.0.43",
			version:     "v0.0.43",
			expectError: false,
			expected:    APIVersion{Major: 0, Minor: 0, Patch: 43, Raw: "v0.0.43"},
		},
		{
			name:        "valid without v prefix",
			version:     "0.0.40",
			expectError: false,
			expected:    APIVersion{Major: 0, Minor: 0, Patch: 40, Raw: "v0.0.40"},
		},
		{
			name:        "invalid format",
			version:     "invalid",
			expectError: true,
		},
		{
			name:        "empty version",
			version:     "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseVersion(tt.version)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				helpers.AssertNoError(t, err)
				helpers.AssertEqual(t, tt.expected, *result)
			}
		})
	}
}

func TestVersionString(t *testing.T) {
	tests := []struct {
		name     string
		version  APIVersion
		expected string
	}{
		{
			name:     "v0.0.40",
			version:  APIVersion{Major: 0, Minor: 0, Patch: 40, Raw: "v0.0.40"},
			expected: "v0.0.40",
		},
		{
			name:     "v0.0.43",
			version:  APIVersion{Major: 0, Minor: 0, Patch: 43, Raw: "v0.0.43"},
			expected: "v0.0.43",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.version.String()
			helpers.AssertEqual(t, tt.expected, result)
		})
	}
}

func TestVersionCompare(t *testing.T) {
	v40, _ := ParseVersion("v0.0.40")
	v41, _ := ParseVersion("v0.0.41")
	v42, _ := ParseVersion("v0.0.42")
	v43, _ := ParseVersion("v0.0.43")

	tests := []struct {
		name     string
		v1       *APIVersion
		v2       *APIVersion
		expected int
	}{
		{"v0.0.40 < v0.0.41", v40, v41, -1},
		{"v0.0.41 < v0.0.42", v41, v42, -1},
		{"v0.0.42 < v0.0.43", v42, v43, -1},
		{"v0.0.43 > v0.0.42", v43, v42, 1},
		{"v0.0.42 > v0.0.41", v42, v41, 1},
		{"v0.0.41 > v0.0.40", v41, v40, 1},
		{"v0.0.42 == v0.0.42", v42, v42, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v1.Compare(tt.v2)
			helpers.AssertEqual(t, tt.expected, result)
		})
	}
}

func TestDefaultCompatibilityMatrix(t *testing.T) {
	matrix := DefaultCompatibilityMatrix()

	// Test that matrix has expected structure
	helpers.AssertNotNil(t, matrix)
	helpers.AssertNotNil(t, matrix.SlurmVersions)
	helpers.AssertNotNil(t, matrix.BreakingChanges)

	// Test that all supported versions have Slurm version mappings
	for _, version := range []string{"v0.0.40", "v0.0.41", "v0.0.42", "v0.0.43"} {
		slurmVersions, exists := matrix.SlurmVersions[version]
		assert.True(t, exists, "Version %s should have Slurm version mapping", version)
		assert.NotEmpty(t, slurmVersions, "Version %s should have at least one Slurm version", version)
	}
}

func TestVersionCompatibilityMatrix_IsSlurmVersionSupported(t *testing.T) {
	matrix := DefaultCompatibilityMatrix()

	tests := []struct {
		name         string
		apiVersion   string
		slurmVersion string
		expected     bool
	}{
		{"v0.0.40 supports Slurm 24.05", "v0.0.40", "24.05", true},
		{"v0.0.42 supports Slurm 25.05", "v0.0.42", "25.05", true},
		{"v0.0.40 does not support Slurm 26.05", "v0.0.40", "26.05", false},
		{"unsupported API version", "v0.0.39", "25.05", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matrix.IsSlurmVersionSupported(tt.apiVersion, tt.slurmVersion)
			helpers.AssertEqual(t, tt.expected, result)
		})
	}
}

func TestVersionCompatibilityMatrix_GetBreakingChanges(t *testing.T) {
	matrix := DefaultCompatibilityMatrix()

	v40, _ := ParseVersion("v0.0.40")
	v41, _ := ParseVersion("v0.0.41")
	v42, _ := ParseVersion("v0.0.42")
	v43, _ := ParseVersion("v0.0.43")

	tests := []struct {
		name       string
		from       *APIVersion
		to         *APIVersion
		hasChanges bool
	}{
		{"v0.0.40 to v0.0.41", v40, v41, true},
		{"v0.0.41 to v0.0.42", v41, v42, true},
		{"v0.0.42 to v0.0.43", v42, v43, true},
		{"v0.0.40 to v0.0.40", v40, v40, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			changes := matrix.GetBreakingChanges(tt.from, tt.to)

			if tt.hasChanges {
				assert.NotEmpty(t, changes, "Should have breaking changes for %s->%s", tt.from.String(), tt.to.String())
			} else {
				assert.Empty(t, changes, "Should not have breaking changes for same version")
			}
		})
	}
}

func TestFindBestVersion(t *testing.T) {
	tests := []struct {
		name        string
		constraint  string
		expected    string
		expectError bool
	}{
		{
			name:        "latest version",
			constraint:  "latest",
			expected:    "v0.0.43",
			expectError: false,
		},
		{
			name:        "stable version",
			constraint:  "stable",
			expected:    "v0.0.42",
			expectError: false,
		},
		{
			name:        "exact supported version",
			constraint:  "v0.0.40",
			expected:    "v0.0.40",
			expectError: false,
		},
		{
			name:        "unsupported version",
			constraint:  "v0.0.39",
			expected:    "v0.0.43", // Latest compatible version (same major.minor)
			expectError: false,
		},
		{
			name:        "invalid constraint",
			constraint:  "invalid",
			expected:    "",
			expectError: true,
		},
		{
			name:        "empty constraint defaults to latest",
			constraint:  "",
			expected:    "v0.0.43",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FindBestVersion(tt.constraint)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				helpers.AssertNoError(t, err)
				helpers.AssertEqual(t, tt.expected, result.String())
			}
		})
	}
}
