package slurm

import (
	"context"
	"fmt"

	"github.com/jontk/slurm-client/internal/factory"
	"github.com/jontk/slurm-client/internal/versioning"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/pkg/config"
	"github.com/jontk/slurm-client/pkg/retry"
)

// ClientOption represents a configuration option for creating a Slurm client
type ClientOption func(*factory.ClientFactory) error

// NewClient creates a new Slurm REST API client with automatic version detection
func NewClient(ctx context.Context, options ...ClientOption) (SlurmClient, error) {
	factoryOptions := make([]factory.FactoryOption, 0, len(options))
	
	for _, option := range options {
		factoryOptions = append(factoryOptions, factory.FactoryOption(option))
	}
	
	clientFactory, err := factory.NewClientFactory(factoryOptions...)
	if err != nil {
		return nil, fmt.Errorf("failed to create client factory: %w", err)
	}
	
	factoryClient, err := clientFactory.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return newClientBridge(factoryClient), nil
}

// NewClientWithVersion creates a new Slurm REST API client for a specific version
func NewClientWithVersion(ctx context.Context, version string, options ...ClientOption) (SlurmClient, error) {
	factoryOptions := make([]factory.FactoryOption, 0, len(options))
	
	for _, option := range options {
		factoryOptions = append(factoryOptions, factory.FactoryOption(option))
	}
	
	clientFactory, err := factory.NewClientFactory(factoryOptions...)
	if err != nil {
		return nil, fmt.Errorf("failed to create client factory: %w", err)
	}
	
	factoryClient, err := clientFactory.NewClientWithVersion(ctx, version)
	if err != nil {
		return nil, err
	}
	return newClientBridge(factoryClient), nil
}

// NewClientForSlurmVersion creates a client compatible with a specific Slurm version
func NewClientForSlurmVersion(ctx context.Context, slurmVersion string, options ...ClientOption) (SlurmClient, error) {
	factoryOptions := make([]factory.FactoryOption, 0, len(options))
	
	for _, option := range options {
		factoryOptions = append(factoryOptions, factory.FactoryOption(option))
	}
	
	clientFactory, err := factory.NewClientFactory(factoryOptions...)
	if err != nil {
		return nil, fmt.Errorf("failed to create client factory: %w", err)
	}
	
	factoryClient, err := clientFactory.NewClientForSlurmVersion(ctx, slurmVersion)
	if err != nil {
		return nil, err
	}
	return newClientBridge(factoryClient), nil
}

// Convenience option functions

// WithConfig sets the client configuration
func WithConfig(cfg *config.Config) ClientOption {
	return func(f *factory.ClientFactory) error {
		return factory.WithConfig(cfg)(f)
	}
}

// WithAuth sets the authentication provider
func WithAuth(auth auth.Provider) ClientOption {
	return func(f *factory.ClientFactory) error {
		return factory.WithAuth(auth)(f)
	}
}

// WithRetryPolicy sets the retry policy
func WithRetryPolicy(policy retry.Policy) ClientOption {
	return func(f *factory.ClientFactory) error {
		return factory.WithRetryPolicy(policy)(f)
	}
}

// WithBaseURL sets the base URL for the Slurm REST API
func WithBaseURL(baseURL string) ClientOption {
	return func(f *factory.ClientFactory) error {
		return factory.WithBaseURL(baseURL)(f)
	}
}

// Version information functions

// SupportedVersions returns all supported API versions
func SupportedVersions() []string {
	versions := make([]string, len(versioning.SupportedVersions))
	for i, v := range versioning.SupportedVersions {
		versions[i] = v.String()
	}
	return versions
}

// LatestVersion returns the latest supported API version
func LatestVersion() string {
	return versioning.LatestVersion().String()
}

// StableVersion returns the stable API version
func StableVersion() string {
	return versioning.StableVersion().String()
}

// IsVersionSupported checks if a version is supported
func IsVersionSupported(version string) bool {
	_, err := versioning.FindBestVersion(version)
	return err == nil
}

// GetVersionCompatibility returns version compatibility information
func GetVersionCompatibility() *versioning.VersionCompatibilityMatrix {
	return versioning.DefaultCompatibilityMatrix()
}

// Error types

// SlurmError represents an error from the Slurm REST API
type SlurmError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Source  string `json:"source"`
	Version string `json:"version,omitempty"`
}

func (e *SlurmError) Error() string {
	if e.Version != "" {
		return fmt.Sprintf("Slurm API %s error %d: %s", e.Version, e.Code, e.Message)
	}
	return fmt.Sprintf("Slurm API error %d: %s", e.Code, e.Message)
}

// VersionError represents a version-related error
type VersionError struct {
	RequestedVersion string
	SupportedVersions []string
	Message          string
}

func (e *VersionError) Error() string {
	return fmt.Sprintf("version error: %s (requested: %s, supported: %v)", 
		e.Message, e.RequestedVersion, e.SupportedVersions)
}

// CompatibilityError represents a compatibility error between versions
type CompatibilityError struct {
	ClientVersion string
	ServerVersion string
	Message       string
}

func (e *CompatibilityError) Error() string {
	return fmt.Sprintf("compatibility error: %s (client: %s, server: %s)", 
		e.Message, e.ClientVersion, e.ServerVersion)
}