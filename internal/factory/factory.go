// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package factory

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	types "github.com/jontk/slurm-client/api"
	"github.com/jontk/slurm-client/internal/versioning"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/pkg/config"
	"github.com/jontk/slurm-client/pkg/retry"
)

// ClientFactory creates version-specific Slurm clients
type ClientFactory struct {
	config      *config.Config
	httpClient  *http.Client
	auth        auth.Provider
	retryPolicy retry.Policy
	baseURL     string

	// Version detection cache
	detectedVersion *versioning.APIVersion
	compatibility   *versioning.VersionCompatibilityMatrix

	// Enhanced options for new features
	enhanced *EnhancedOptions
}

// NewClientFactory creates a new client factory
func NewClientFactory(options ...Option) (*ClientFactory, error) {
	cfg := config.NewDefault()
	factory := &ClientFactory{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		// Create retry policy using config values
		retryPolicy: retry.NewHTTPExponentialBackoff().
			WithMaxRetries(cfg.MaxRetries).
			WithMinWaitTime(cfg.RetryWaitMin).
			WithMaxWaitTime(cfg.RetryWaitMax),
		compatibility: versioning.DefaultCompatibilityMatrix(),
	}

	for _, option := range options {
		if err := option(factory); err != nil {
			return nil, err
		}
	}

	if factory.baseURL == "" {
		factory.baseURL = factory.config.BaseURL
	}

	// Re-apply retry config if config was changed by options
	if factory.config != cfg {
		factory.retryPolicy = retry.NewHTTPExponentialBackoff().
			WithMaxRetries(factory.config.MaxRetries).
			WithMinWaitTime(factory.config.RetryWaitMin).
			WithMaxWaitTime(factory.config.RetryWaitMax)
	}

	// Wire config fields to enhanced options
	factory.applyConfigToEnhancedOptions()

	return factory, nil
}

// applyConfigToEnhancedOptions applies pkg/config.Config fields to EnhancedOptions.
// This ensures config fields like Timeout, UserAgent, MaxRetries, RetryPolicy are actually used.
func (f *ClientFactory) applyConfigToEnhancedOptions() {
	if f.config == nil {
		return
	}

	// Initialize enhanced options if needed
	if f.enhanced == nil {
		f.enhanced = &EnhancedOptions{}
	}

	// Wire config fields to enhanced options (only if not already set)
	if f.config.Timeout > 0 && f.enhanced.DefaultTimeout == 0 {
		f.enhanced.DefaultTimeout = f.config.Timeout
	}
	if f.config.UserAgent != "" && f.enhanced.UserAgent == "" {
		f.enhanced.UserAgent = f.config.UserAgent
	}
	if f.config.MaxRetries > 0 && f.enhanced.MaxRetries == 0 {
		f.enhanced.MaxRetries = f.config.MaxRetries
	}
	if f.config.Debug {
		f.enhanced.Debug = true
	}

	// Wire retry policy to enhanced options (this ensures retryPolicy is actually used)
	if f.retryPolicy != nil && f.enhanced.RetryBackoff == nil {
		f.enhanced.RetryBackoff = f.retryPolicy
	}
}

// Option represents a configuration option for the ClientFactory
type Option func(*ClientFactory) error

// FactoryOption is a deprecated alias for Option, kept for backward compatibility
type FactoryOption = Option

// WithConfig sets the factory configuration
func WithConfig(cfg *config.Config) Option {
	return func(f *ClientFactory) error {
		f.config = cfg
		return nil
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) Option {
	return func(f *ClientFactory) error {
		f.httpClient = httpClient
		return nil
	}
}

// WithAuth sets the authentication provider
func WithAuth(auth auth.Provider) Option {
	return func(f *ClientFactory) error {
		f.auth = auth
		return nil
	}
}

// WithRetryPolicy sets the retry policy
func WithRetryPolicy(policy retry.Policy) Option {
	return func(f *ClientFactory) error {
		f.retryPolicy = policy
		return nil
	}
}

// WithBaseURL sets the base URL for the Slurm REST API
func WithBaseURL(baseURL string) Option {
	return func(f *ClientFactory) error {
		f.baseURL = baseURL
		return nil
	}
}

// SetTimeout modifies the timeout of the existing HTTP client without replacing it
// This preserves TLS configuration and custom transport settings
func (f *ClientFactory) SetTimeout(timeout time.Duration) error {
	if f.httpClient == nil {
		f.httpClient = &http.Client{}
	}
	f.httpClient.Timeout = timeout
	return nil
}

// NewClient creates a new Slurm client with automatic version detection
func (f *ClientFactory) NewClient(ctx context.Context) (SlurmClient, error) {
	return f.NewClientWithVersion(ctx, "")
}

// NewClientWithVersion creates a new Slurm client for a specific version
func (f *ClientFactory) NewClientWithVersion(ctx context.Context, version string) (SlurmClient, error) {
	var targetVersion *versioning.APIVersion
	var err error

	// Use config.APIVersion if no version specified and config has one set
	if version == "" && f.config.APIVersion != "" {
		version = f.config.APIVersion
	}

	if version == "" {
		// Auto-detect version
		targetVersion, err = f.detectVersion(ctx)
		if err != nil {
			// Fallback to stable version
			if f.config.Debug {
				fmt.Printf("Version detection failed, using stable version: %v\n", err)
			}
			targetVersion = versioning.StableVersion()
		}
	} else {
		// Use specified version
		targetVersion, err = versioning.FindBestVersion(version)
		if err != nil {
			return nil, fmt.Errorf("invalid version %s: %w", version, err)
		}
	}

	return f.createClient(ctx, targetVersion)
}

// NewClientForSlurmVersion creates a client compatible with a specific Slurm version
func (f *ClientFactory) NewClientForSlurmVersion(ctx context.Context, slurmVersion string) (SlurmClient, error) {
	// Find compatible API version for the Slurm version
	var compatibleVersion *versioning.APIVersion

	for _, apiVersion := range versioning.SupportedVersions {
		if f.compatibility.IsSlurmVersionSupported(apiVersion.String(), slurmVersion) {
			if compatibleVersion == nil || apiVersion.Compare(compatibleVersion) > 0 {
				compatibleVersion = apiVersion
			}
		}
	}

	if compatibleVersion == nil {
		return nil, fmt.Errorf("no compatible API version found for Slurm %s", slurmVersion)
	}

	return f.createClient(ctx, compatibleVersion)
}

// ListSupportedVersions returns all supported API versions
func (f *ClientFactory) ListSupportedVersions() []*versioning.APIVersion {
	return versioning.SupportedVersions
}

// GetVersionCompatibility returns version compatibility information
func (f *ClientFactory) GetVersionCompatibility() *versioning.VersionCompatibilityMatrix {
	return f.compatibility
}

// detectVersion detects the API version by querying the OpenAPI endpoint
func (f *ClientFactory) detectVersion(ctx context.Context) (*versioning.APIVersion, error) {
	if f.detectedVersion != nil {
		return f.detectedVersion, nil
	}

	// Try to get OpenAPI spec to detect version
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, f.baseURL+"/openapi/v3", http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create version detection request: %w", err)
	}

	// Add authentication if available
	if f.auth != nil {
		if err := f.auth.Authenticate(ctx, req); err != nil {
			if f.config.Debug {
				fmt.Printf("Authentication failed during version detection: %v\n", err)
			}
		}
	}

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to detect version: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("version detection failed with status %d", resp.StatusCode)
	}

	var openAPISpec struct {
		Info struct {
			Version string `json:"version"`
		} `json:"info"`
		Servers []struct {
			URL string `json:"url"`
		} `json:"servers"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&openAPISpec); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI spec: %w", err)
	}

	// Extract version from server URLs or info
	var detectedVersionStr string
	if openAPISpec.Info.Version != "" {
		detectedVersionStr = openAPISpec.Info.Version
	} else if len(openAPISpec.Servers) > 0 {
		// Try to extract version from server URL
		// Example: /slurm/v0.0.42/ -> v0.0.42
		for _, server := range openAPISpec.Servers {
			if version := extractVersionFromURL(server.URL); version != "" {
				detectedVersionStr = version
				break
			}
		}
	}

	if detectedVersionStr == "" {
		return nil, fmt.Errorf("could not determine API version from OpenAPI spec")
	}

	// Check if this is a SLURM version string (Slurm-x.y.z format)
	// and map it to a compatible API version
	var version *versioning.APIVersion
	if strings.HasPrefix(detectedVersionStr, "Slurm-") {
		slurmVersion := strings.TrimPrefix(detectedVersionStr, "Slurm-")
		version, err = f.findCompatibleAPIVersion(slurmVersion)
		if err != nil {
			return nil, fmt.Errorf("invalid detected SLURM version %s: %w", detectedVersionStr, err)
		}
	} else {
		version, err = versioning.ParseVersion(detectedVersionStr)
		if err != nil {
			return nil, fmt.Errorf("invalid detected version %s: %w", detectedVersionStr, err)
		}
	}

	// Verify this version is supported
	supported := false
	for _, supportedVersion := range versioning.SupportedVersions {
		if version.Compare(supportedVersion) == 0 {
			supported = true
			break
		}
	}

	if !supported {
		return nil, fmt.Errorf("detected version %s is not supported", version.String())
	}

	f.detectedVersion = version
	return version, nil
}

// findCompatibleAPIVersion finds a compatible API version for the given SLURM version
func (f *ClientFactory) findCompatibleAPIVersion(slurmVersion string) (*versioning.APIVersion, error) {
	var compatibleVersion *versioning.APIVersion

	// Find the best compatible API version for this SLURM version
	for _, apiVersion := range versioning.SupportedVersions {
		if f.compatibility.IsSlurmVersionSupported(apiVersion.String(), slurmVersion) {
			if compatibleVersion == nil || apiVersion.Compare(compatibleVersion) > 0 {
				compatibleVersion = apiVersion
			}
		}
	}

	if compatibleVersion == nil {
		return nil, fmt.Errorf("no compatible API version found for SLURM %s", slurmVersion)
	}

	return compatibleVersion, nil
}

// createClient creates a version-specific client implementation
func (f *ClientFactory) createClient(ctx context.Context, version *versioning.APIVersion) (SlurmClient, error) {
	switch version.String() {
	case "v0.0.40":
		return f.createV0_0_40Client(ctx)
	case "v0.0.41":
		return f.createV0_0_41Client(ctx)
	case "v0.0.42":
		return f.createV0_0_42Client(ctx)
	case "v0.0.43":
		return f.createV0_0_43Client(ctx)
	case "v0.0.44":
		return f.createV0_0_44Client(ctx)
	default:
		return nil, fmt.Errorf("unsupported API version: %s", version.String())
	}
}

// Version-specific client creation methods (to be implemented with generated code)

func (f *ClientFactory) createV0_0_40Client(ctx context.Context) (SlurmClient, error) {
	// Create enhanced HTTP client with all features
	httpClient := f.buildEnhancedHTTPClient(ctx)

	// Apply authentication if needed
	if f.auth != nil {
		httpClient = createAuthenticatedHTTPClient(httpClient, f.auth)
	}

	// Create adapter client config
	config := &types.ClientConfig{
		BaseURL:    f.baseURL,
		HTTPClient: httpClient,
		Debug:      f.config.Debug,
	}
	client, err := NewAdapterClient("v0.0.40", config)
	if err != nil {
		return nil, err
	}
	// Set the pool for proper cleanup on Close()
	if ac, ok := client.(*AdapterClient); ok && f.enhanced != nil && f.enhanced.ConnectionPool != nil {
		ac.SetPool(f.enhanced.ConnectionPool)
	}
	return client, nil
}

func (f *ClientFactory) createV0_0_41Client(ctx context.Context) (SlurmClient, error) {
	// Create enhanced HTTP client with all features
	httpClient := f.buildEnhancedHTTPClient(ctx)

	// Apply authentication if needed
	if f.auth != nil {
		httpClient = createAuthenticatedHTTPClient(httpClient, f.auth)
	}

	// Create adapter client config
	config := &types.ClientConfig{
		BaseURL:    f.baseURL,
		HTTPClient: httpClient,
		Debug:      f.config.Debug,
	}
	client, err := NewAdapterClient("v0.0.41", config)
	if err != nil {
		return nil, err
	}
	// Set the pool for proper cleanup on Close()
	if ac, ok := client.(*AdapterClient); ok && f.enhanced != nil && f.enhanced.ConnectionPool != nil {
		ac.SetPool(f.enhanced.ConnectionPool)
	}
	return client, nil
}

func (f *ClientFactory) createV0_0_42Client(ctx context.Context) (SlurmClient, error) {
	// Create enhanced HTTP client with all features
	httpClient := f.buildEnhancedHTTPClient(ctx)

	// Apply authentication if needed
	if f.auth != nil {
		httpClient = createAuthenticatedHTTPClient(httpClient, f.auth)
	}

	// Create adapter client config
	config := &types.ClientConfig{
		BaseURL:    f.baseURL,
		HTTPClient: httpClient,
		Debug:      f.config.Debug,
	}
	client, err := NewAdapterClient("v0.0.42", config)
	if err != nil {
		return nil, err
	}
	// Set the pool for proper cleanup on Close()
	if ac, ok := client.(*AdapterClient); ok && f.enhanced != nil && f.enhanced.ConnectionPool != nil {
		ac.SetPool(f.enhanced.ConnectionPool)
	}
	return client, nil
}

func (f *ClientFactory) createV0_0_43Client(ctx context.Context) (SlurmClient, error) {
	// Create enhanced HTTP client with all features
	httpClient := f.buildEnhancedHTTPClient(ctx)

	// Apply authentication if needed
	if f.auth != nil {
		httpClient = createAuthenticatedHTTPClient(httpClient, f.auth)
	}

	// Create adapter client config
	config := &types.ClientConfig{
		BaseURL:    f.baseURL,
		HTTPClient: httpClient,
		Debug:      f.config.Debug,
	}
	client, err := NewAdapterClient("v0.0.43", config)
	if err != nil {
		return nil, err
	}
	// Set the pool for proper cleanup on Close()
	if ac, ok := client.(*AdapterClient); ok && f.enhanced != nil && f.enhanced.ConnectionPool != nil {
		ac.SetPool(f.enhanced.ConnectionPool)
	}
	return client, nil
}

func (f *ClientFactory) createV0_0_44Client(ctx context.Context) (SlurmClient, error) {
	// Create enhanced HTTP client with all features
	httpClient := f.buildEnhancedHTTPClient(ctx)

	// Apply authentication if needed
	if f.auth != nil {
		httpClient = createAuthenticatedHTTPClient(httpClient, f.auth)
	}

	// Use adapters for v0.0.44 as they are now implemented
	config := &types.ClientConfig{
		BaseURL:    f.baseURL,
		HTTPClient: httpClient,
		Debug:      f.config.Debug,
	}
	client, err := NewAdapterClient("v0.0.44", config)
	if err != nil {
		return nil, err
	}
	// Set the pool for proper cleanup on Close()
	if ac, ok := client.(*AdapterClient); ok && f.enhanced != nil && f.enhanced.ConnectionPool != nil {
		ac.SetPool(f.enhanced.ConnectionPool)
	}
	return client, nil
}

// extractVersionFromURL extracts version from a URL like "/slurm/v0.0.42/"
func extractVersionFromURL(url string) string {
	parts := strings.Split(strings.Trim(url, "/"), "/")
	for _, part := range parts {
		if strings.HasPrefix(part, "v") && strings.Count(part, ".") == 2 {
			return part
		}
	}
	return ""
}

// Helper method to create common client configuration
