package factory

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	v040 "github.com/jontk/slurm-client/internal/api/v0_0_40"
	v041 "github.com/jontk/slurm-client/internal/api/v0_0_41"
	v042 "github.com/jontk/slurm-client/internal/api/v0_0_42"
	v043 "github.com/jontk/slurm-client/internal/api/v0_0_43"
	"github.com/jontk/slurm-client/internal/interfaces"
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
	
	// Use adapters instead of wrapper clients
	useAdapters bool
}

// NewClientFactory creates a new client factory
func NewClientFactory(options ...FactoryOption) (*ClientFactory, error) {
	factory := &ClientFactory{
		config: config.NewDefault(),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		retryPolicy:   retry.NewHTTPExponentialBackoff(),
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

	return factory, nil
}

// FactoryOption represents a configuration option for the ClientFactory
type FactoryOption func(*ClientFactory) error

// WithConfig sets the factory configuration
func WithConfig(cfg *config.Config) FactoryOption {
	return func(f *ClientFactory) error {
		f.config = cfg
		return nil
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) FactoryOption {
	return func(f *ClientFactory) error {
		f.httpClient = httpClient
		return nil
	}
}

// WithAuth sets the authentication provider
func WithAuth(auth auth.Provider) FactoryOption {
	return func(f *ClientFactory) error {
		f.auth = auth
		return nil
	}
}

// WithRetryPolicy sets the retry policy
func WithRetryPolicy(policy retry.Policy) FactoryOption {
	return func(f *ClientFactory) error {
		f.retryPolicy = policy
		return nil
	}
}

// WithBaseURL sets the base URL for the Slurm REST API
func WithBaseURL(baseURL string) FactoryOption {
	return func(f *ClientFactory) error {
		f.baseURL = baseURL
		return nil
	}
}

// WithUseAdapters enables the use of adapter implementations instead of wrapper clients
func WithUseAdapters(useAdapters bool) FactoryOption {
	return func(f *ClientFactory) error {
		f.useAdapters = useAdapters
		return nil
	}
}

// NewClient creates a new Slurm client with automatic version detection
func (f *ClientFactory) NewClient(ctx context.Context) (SlurmClient, error) {
	return f.NewClientWithVersion(ctx, "")
}

// NewClientWithVersion creates a new Slurm client for a specific version
func (f *ClientFactory) NewClientWithVersion(ctx context.Context, version string) (SlurmClient, error) {
	var targetVersion *versioning.APIVersion
	var err error

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

	return f.createClient(targetVersion)
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

	return f.createClient(compatibleVersion)
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
	req, err := http.NewRequestWithContext(ctx, "GET", f.baseURL+"/openapi/v3", nil)
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
	defer resp.Body.Close()

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

	version, err := versioning.ParseVersion(detectedVersionStr)
	if err != nil {
		return nil, fmt.Errorf("invalid detected version %s: %w", detectedVersionStr, err)
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

// createClient creates a version-specific client implementation
func (f *ClientFactory) createClient(version *versioning.APIVersion) (SlurmClient, error) {
	switch version.String() {
	case "v0.0.40":
		return f.createV0_0_40Client()
	case "v0.0.41":
		return f.createV0_0_41Client()
	case "v0.0.42":
		return f.createV0_0_42Client()
	case "v0.0.43":
		return f.createV0_0_43Client()
	default:
		return nil, fmt.Errorf("unsupported API version: %s", version.String())
	}
}

// Version-specific client creation methods (to be implemented with generated code)

func (f *ClientFactory) createV0_0_40Client() (SlurmClient, error) {
	// Create enhanced HTTP client with all features
	httpClient := f.buildEnhancedHTTPClient()
	
	// Apply authentication if needed
	if f.auth != nil {
		httpClient = createAuthenticatedHTTPClient(httpClient, f.auth)
	}

	// Check if adapters should be used
	if f.useAdapters {
		// Create adapter client config
		/*config := &ClientConfig{
			BaseURL:    f.baseURL,
			HTTPClient: httpClient,
			APIKey:     "",    // Not used when we have auth provider
			Debug:      f.config.Debug,
		}
		return NewAdapterClient("v0.0.40", config)*/
		return nil, fmt.Errorf("adapter implementation is incomplete and disabled")
	}

	// Create the wrapper client config
	config := &interfaces.ClientConfig{
		BaseURL:    f.baseURL,
		HTTPClient: httpClient,
		APIKey:     "",    // Not used when we have auth provider
		Debug:      f.config.Debug,
	}

	// Return the wrapper client
	return v040.NewWrapperClient(config)
}

func (f *ClientFactory) createV0_0_41Client() (SlurmClient, error) {
	// Create enhanced HTTP client with all features
	httpClient := f.buildEnhancedHTTPClient()
	
	// Apply authentication if needed
	if f.auth != nil {
		httpClient = createAuthenticatedHTTPClient(httpClient, f.auth)
	}

	// Check if adapters should be used
	if f.useAdapters {
		// Create adapter client config
		/*config := &ClientConfig{
			BaseURL:    f.baseURL,
			HTTPClient: httpClient,
			APIKey:     "",    // Not used when we have auth provider
			Debug:      f.config.Debug,
		}
		return NewAdapterClient("v0.0.41", config)*/
		return nil, fmt.Errorf("adapter implementation is incomplete and disabled")
	}

	// Create the wrapper client config
	config := &interfaces.ClientConfig{
		BaseURL:    f.baseURL,
		HTTPClient: httpClient,
		APIKey:     "",    // Not used when we have auth provider
		Debug:      f.config.Debug,
	}

	// Return the wrapper client
	return v041.NewWrapperClient(config)
}

func (f *ClientFactory) createV0_0_42Client() (SlurmClient, error) {
	// Create enhanced HTTP client with all features
	httpClient := f.buildEnhancedHTTPClient()
	
	// Apply authentication if needed
	if f.auth != nil {
		httpClient = createAuthenticatedHTTPClient(httpClient, f.auth)
	}

	// Check if adapters should be used
	if f.useAdapters {
		// Create adapter client config
		/*config := &ClientConfig{
			BaseURL:    f.baseURL,
			HTTPClient: httpClient,
			APIKey:     "",    // Not used when we have auth provider
			Debug:      f.config.Debug,
		}
		return NewAdapterClient("v0.0.42", config)*/
		return nil, fmt.Errorf("adapter implementation is incomplete and disabled")
	}

	// Create the wrapper client config
	config := &interfaces.ClientConfig{
		BaseURL:    f.baseURL,
		HTTPClient: httpClient,
		APIKey:     "",    // Not used when we have auth provider
		Debug:      f.config.Debug,
	}

	// Return the wrapper client
	return v042.NewWrapperClient(config)
}

func (f *ClientFactory) createV0_0_43Client() (SlurmClient, error) {
	// Create enhanced HTTP client with all features
	httpClient := f.buildEnhancedHTTPClient()
	
	// Apply authentication if needed
	if f.auth != nil {
		httpClient = createAuthenticatedHTTPClient(httpClient, f.auth)
	}

	// Check if adapters should be used
	if f.useAdapters {
		// Create adapter client config
		/*config := &ClientConfig{
			BaseURL:    f.baseURL,
			HTTPClient: httpClient,
			APIKey:     "",    // Not used when we have auth provider
			Debug:      f.config.Debug,
		}
		return NewAdapterClient("v0.0.43", config)*/
		return nil, fmt.Errorf("adapter implementation is incomplete and disabled")
	}

	// Create the wrapper client config
	config := &interfaces.ClientConfig{
		BaseURL:    f.baseURL,
		HTTPClient: httpClient,
		APIKey:     "",    // Not used when we have auth provider
		Debug:      f.config.Debug,
	}

	// Return the wrapper client
	return v043.NewWrapperClient(config)
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
