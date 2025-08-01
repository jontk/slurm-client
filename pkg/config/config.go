// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds configuration for the Slurm client
type Config struct {
	// BaseURL is the base URL for the Slurm REST API
	BaseURL string

	// Timeout is the request timeout
	Timeout time.Duration

	// UserAgent is the user agent string
	UserAgent string

	// MaxRetries is the maximum number of retries
	MaxRetries int

	// RetryWaitMin is the minimum wait time between retries
	RetryWaitMin time.Duration

	// RetryWaitMax is the maximum wait time between retries
	RetryWaitMax time.Duration

	// APIVersion is the Slurm API version to use
	APIVersion string

	// Debug enables debug logging
	Debug bool

	// InsecureSkipVerify skips TLS certificate verification
	InsecureSkipVerify bool
}

// NewDefault creates a new configuration with default values
func NewDefault() *Config {
	return &Config{
		BaseURL:            getEnvOrDefault("SLURM_REST_URL", "http://localhost:6820"),
		Timeout:            30 * time.Second,
		UserAgent:          "slurm-client/1.0",
		MaxRetries:         3,
		RetryWaitMin:       1 * time.Second,
		RetryWaitMax:       30 * time.Second,
		APIVersion:         "v0.0.39",
		Debug:              getEnvBoolOrDefault("SLURM_DEBUG", false),
		InsecureSkipVerify: getEnvBoolOrDefault("SLURM_INSECURE_SKIP_VERIFY", false),
	}
}

// Load loads configuration from environment variables
func (c *Config) Load() {
	if url := os.Getenv("SLURM_REST_URL"); url != "" {
		c.BaseURL = url
	}

	if timeout := os.Getenv("SLURM_TIMEOUT"); timeout != "" {
		if d, err := time.ParseDuration(timeout); err == nil {
			c.Timeout = d
		}
	}

	if userAgent := os.Getenv("SLURM_USER_AGENT"); userAgent != "" {
		c.UserAgent = userAgent
	}

	if maxRetries := os.Getenv("SLURM_MAX_RETRIES"); maxRetries != "" {
		if i, err := strconv.Atoi(maxRetries); err == nil {
			c.MaxRetries = i
		}
	}

	if apiVersion := os.Getenv("SLURM_API_VERSION"); apiVersion != "" {
		c.APIVersion = apiVersion
	}

	c.Debug = getEnvBoolOrDefault("SLURM_DEBUG", c.Debug)
	c.InsecureSkipVerify = getEnvBoolOrDefault("SLURM_INSECURE_SKIP_VERIFY", c.InsecureSkipVerify)
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.BaseURL == "" {
		return ErrMissingBaseURL
	}

	if c.Timeout <= 0 {
		return ErrInvalidTimeout
	}

	if c.MaxRetries < 0 {
		return ErrInvalidMaxRetries
	}

	return nil
}

// getEnvOrDefault returns the environment variable value or a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvBoolOrDefault returns the environment variable value as a boolean or a default value
func getEnvBoolOrDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if b, err := strconv.ParseBool(value); err == nil {
			return b
		}
	}
	return defaultValue
}
