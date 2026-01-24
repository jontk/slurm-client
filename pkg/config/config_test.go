// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"testing"
	"time"

	"github.com/jontk/slurm-client/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestNewDefault(t *testing.T) {
	config := NewDefault()

	helpers.AssertNotNil(t, config)

	// Check default values
	helpers.AssertEqual(t, false, config.Debug)
	helpers.AssertEqual(t, false, config.InsecureSkipVerify)
	helpers.AssertEqual(t, "slurm-client/1.0", config.UserAgent)
	helpers.AssertEqual(t, "v0.0.39", config.APIVersion)

	// Verify defaults are reasonable
	assert.Greater(t, config.Timeout, time.Duration(0))
	assert.Positive(t, config.MaxRetries)
	assert.Greater(t, config.RetryWaitMin, time.Duration(0))
	assert.Greater(t, config.RetryWaitMax, time.Duration(0))
}

func TestConfigLoad(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected func(*Config)
	}{
		{
			name: "base URL from environment",
			envVars: map[string]string{
				"SLURM_REST_URL": "https://slurm.example.com:6820",
			},
			expected: func(config *Config) {
				helpers.AssertEqual(t, "https://slurm.example.com:6820", config.BaseURL)
			},
		},
		{
			name: "timeout from environment",
			envVars: map[string]string{
				"SLURM_TIMEOUT": "60s",
			},
			expected: func(config *Config) {
				assert.Greater(t, config.Timeout, time.Duration(0))
			},
		},
		{
			name: "user agent from environment",
			envVars: map[string]string{
				"SLURM_USER_AGENT": "custom-client/2.0",
			},
			expected: func(config *Config) {
				helpers.AssertEqual(t, "custom-client/2.0", config.UserAgent)
			},
		},
		{
			name: "max retries from environment",
			envVars: map[string]string{
				"SLURM_MAX_RETRIES": "5",
			},
			expected: func(config *Config) {
				helpers.AssertEqual(t, 5, config.MaxRetries)
			},
		},
		{
			name: "API version from environment",
			envVars: map[string]string{
				"SLURM_API_VERSION": "v0.0.42",
			},
			expected: func(config *Config) {
				helpers.AssertEqual(t, "v0.0.42", config.APIVersion)
			},
		},
		{
			name: "debug from environment",
			envVars: map[string]string{
				"SLURM_DEBUG": "true",
			},
			expected: func(config *Config) {
				helpers.AssertEqual(t, true, config.Debug)
			},
		},
		{
			name: "insecure skip verify from environment",
			envVars: map[string]string{
				"SLURM_INSECURE_SKIP_VERIFY": "true",
			},
			expected: func(config *Config) {
				helpers.AssertEqual(t, true, config.InsecureSkipVerify)
			},
		},
		{
			name: "all environment variables",
			envVars: map[string]string{
				"SLURM_REST_URL":             "https://slurm.example.com:6820",
				"SLURM_TIMEOUT":              "120s",
				"SLURM_USER_AGENT":           "test-client/1.0",
				"SLURM_MAX_RETRIES":          "10",
				"SLURM_API_VERSION":          "v0.0.42",
				"SLURM_DEBUG":                "true",
				"SLURM_INSECURE_SKIP_VERIFY": "true",
			},
			expected: func(config *Config) {
				helpers.AssertEqual(t, "https://slurm.example.com:6820", config.BaseURL)
				helpers.AssertEqual(t, "test-client/1.0", config.UserAgent)
				helpers.AssertEqual(t, 10, config.MaxRetries)
				helpers.AssertEqual(t, "v0.0.42", config.APIVersion)
				helpers.AssertEqual(t, true, config.Debug)
				helpers.AssertEqual(t, true, config.InsecureSkipVerify)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for key, value := range tt.envVars {
				t.Setenv(key, value)
			}

			config := NewDefault()
			config.Load()

			helpers.AssertNotNil(t, config)
			tt.expected(config)
		})
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
		expectedErr error
	}{
		{
			name: "valid config",
			config: &Config{
				BaseURL:    "https://example.com",
				Timeout:    30 * time.Second,
				MaxRetries: 3,
			},
			expectError: false,
		},
		{
			name: "missing base URL",
			config: &Config{
				Timeout:    30 * time.Second,
				MaxRetries: 3,
			},
			expectError: true,
			expectedErr: ErrMissingBaseURL,
		},
		{
			name: "empty base URL",
			config: &Config{
				BaseURL:    "",
				Timeout:    30 * time.Second,
				MaxRetries: 3,
			},
			expectError: true,
			expectedErr: ErrMissingBaseURL,
		},
		{
			name: "invalid timeout",
			config: &Config{
				BaseURL:    "https://example.com",
				Timeout:    -1 * time.Second,
				MaxRetries: 3,
			},
			expectError: true,
			expectedErr: ErrInvalidTimeout,
		},
		{
			name: "invalid max retries",
			config: &Config{
				BaseURL:    "https://example.com",
				Timeout:    30 * time.Second,
				MaxRetries: -1,
			},
			expectError: true,
			expectedErr: ErrInvalidMaxRetries,
		},
		{
			name: "zero timeout",
			config: &Config{
				BaseURL:    "https://example.com",
				Timeout:    0,
				MaxRetries: 3,
			},
			expectError: true,
			expectedErr: ErrInvalidTimeout,
		},
		{
			name: "zero max retries (should be valid)",
			config: &Config{
				BaseURL:    "https://example.com",
				Timeout:    30 * time.Second,
				MaxRetries: 0,
			},
			expectError: false, // 0 retries means no retries, which is valid
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedErr != nil {
					helpers.AssertEqual(t, tt.expectedErr, err)
				}
			} else {
				helpers.AssertNoError(t, err)
			}
		})
	}
}

func TestConfigMutation(t *testing.T) {
	config := NewDefault()

	// Test that we can modify config fields directly
	config.BaseURL = "https://example.com"
	helpers.AssertEqual(t, "https://example.com", config.BaseURL)

	config.Timeout = 60 * time.Second
	helpers.AssertEqual(t, 60*time.Second, config.Timeout)

	config.MaxRetries = 5
	helpers.AssertEqual(t, 5, config.MaxRetries)

	config.Debug = true
	helpers.AssertEqual(t, true, config.Debug)

	config.InsecureSkipVerify = true
	helpers.AssertEqual(t, true, config.InsecureSkipVerify)

	config.UserAgent = "test-client/1.0"
	helpers.AssertEqual(t, "test-client/1.0", config.UserAgent)

	config.APIVersion = "v0.0.42"
	helpers.AssertEqual(t, "v0.0.42", config.APIVersion)
}

func TestConfigDefaults(t *testing.T) {
	// Test that NewDefault returns expected defaults
	config := NewDefault()

	// Should have default localhost URL
	helpers.AssertEqual(t, "http://localhost:6820", config.BaseURL)

	// Should have reasonable timeout
	helpers.AssertEqual(t, 30*time.Second, config.Timeout)

	// Should have default user agent
	helpers.AssertEqual(t, "slurm-client/1.0", config.UserAgent)

	// Should have default max retries
	helpers.AssertEqual(t, 3, config.MaxRetries)

	// Should have default API version
	helpers.AssertEqual(t, "v0.0.39", config.APIVersion)

	// Should have default boolean values
	helpers.AssertEqual(t, false, config.Debug)
	helpers.AssertEqual(t, false, config.InsecureSkipVerify)
}
