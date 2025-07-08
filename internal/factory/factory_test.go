package factory

import (
	"testing"

	"github.com/jontk/slurm-client/pkg/config"
	"github.com/jontk/slurm-client/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestNewClientFactory(t *testing.T) {
	tests := []struct {
		name        string
		options     []FactoryOption
		expectError bool
	}{
		{
			name: "create factory with default options",
			options: []FactoryOption{
				WithBaseURL("https://example.com"),
			},
			expectError: false,
		},
		{
			name: "create factory with custom config",
			options: []FactoryOption{
				WithConfig(config.NewDefault()),
				WithBaseURL("https://example.com"),
			},
			expectError: false,
		},
		{
			name:        "create factory with no options",
			options:     []FactoryOption{},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory, err := NewClientFactory(tt.options...)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, factory)
			} else {
				helpers.AssertNoError(t, err)
				helpers.AssertNotNil(t, factory)
			}
		})
	}
}

func TestClientFactory_NewClientWithVersion(t *testing.T) {
	ctx := helpers.TestContext(t)
	
	factory, err := NewClientFactory(
		WithBaseURL("https://example.com"),
	)
	helpers.RequireNoError(t, err)
	
	tests := []struct {
		name        string
		version     string
		expectError bool
	}{
		{
			name:        "create client with v0.0.40",
			version:     "v0.0.40",
			expectError: false,
		},
		{
			name:        "create client with v0.0.41",
			version:     "v0.0.41",
			expectError: true, // Bridge not implemented yet
		},
		{
			name:        "create client with v0.0.42",
			version:     "v0.0.42",
			expectError: false,
		},
		{
			name:        "create client with v0.0.43",
			version:     "v0.0.43",
			expectError: true, // Bridge not implemented yet
		},
		{
			name:        "create client with latest",
			version:     "latest",
			expectError: true, // Latest is v0.0.43, not implemented
		},
		{
			name:        "create client with stable",
			version:     "stable",
			expectError: false, // Stable is v0.0.42
		},
		{
			name:        "create client with unsupported version",
			version:     "v0.0.39",
			expectError: true, // Compatible version (v0.0.43) not implemented
		},
		{
			name:        "create client with invalid version",
			version:     "invalid",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := factory.NewClientWithVersion(ctx, tt.version)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, client)
			} else {
				helpers.AssertNoError(t, err)
				helpers.AssertNotNil(t, client)
				
				// Verify client has expected interface methods
				assert.NotNil(t, client.Jobs())
				assert.NotNil(t, client.Nodes())
				assert.NotNil(t, client.Partitions())
				assert.NotNil(t, client.Info())
			}
		})
	}
}

func TestClientFactory_NewClient(t *testing.T) {
	ctx := helpers.TestContext(t)
	
	factory, err := NewClientFactory(
		WithBaseURL("https://example.com"),
	)
	helpers.RequireNoError(t, err)
	
	// Test with empty version (should auto-detect, fallback to stable)
	client, err := factory.NewClient(ctx)
	
	// Should succeed with stable version since we can't actually detect
	helpers.AssertNoError(t, err)
	helpers.AssertNotNil(t, client)
}

func TestClientFactory_NewClientForSlurmVersion(t *testing.T) {
	ctx := helpers.TestContext(t)
	
	factory, err := NewClientFactory(
		WithBaseURL("https://example.com"),
	)
	helpers.RequireNoError(t, err)
	
	tests := []struct {
		name         string
		slurmVersion string
		expectError  bool
	}{
		{
			name:         "Slurm 24.05",
			slurmVersion: "24.05",
			expectError:  false, // Should use v0.0.40
		},
		{
			name:         "Slurm 25.05",
			slurmVersion: "25.05",
			expectError:  true, // Would use v0.0.43 (latest), but not implemented
		},
		{
			name:         "unsupported Slurm version",
			slurmVersion: "20.02",
			expectError:  true,
		},
		{
			name:         "invalid Slurm version",
			slurmVersion: "invalid",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := factory.NewClientForSlurmVersion(ctx, tt.slurmVersion)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, client)
			} else {
				helpers.AssertNoError(t, err)
				helpers.AssertNotNil(t, client)
			}
		})
	}
}

func TestClientFactory_ListSupportedVersions(t *testing.T) {
	factory, err := NewClientFactory(
		WithBaseURL("https://example.com"),
	)
	helpers.RequireNoError(t, err)
	
	versions := factory.ListSupportedVersions()
	
	// Verify we have the expected versions
	expectedVersions := []string{"v0.0.40", "v0.0.41", "v0.0.42", "v0.0.43"}
	
	helpers.AssertEqual(t, len(expectedVersions), len(versions))
	
	for i, version := range versions {
		helpers.AssertEqual(t, expectedVersions[i], version.String())
	}
}

func TestClientFactory_GetVersionCompatibility(t *testing.T) {
	factory, err := NewClientFactory(
		WithBaseURL("https://example.com"),
	)
	helpers.RequireNoError(t, err)
	
	matrix := factory.GetVersionCompatibility()
	
	helpers.AssertNotNil(t, matrix)
	helpers.AssertNotNil(t, matrix.SlurmVersions)
	helpers.AssertNotNil(t, matrix.BreakingChanges)
	
	// Verify all supported versions have Slurm version mappings
	for _, version := range []string{"v0.0.40", "v0.0.41", "v0.0.42", "v0.0.43"} {
		slurmVersions, exists := matrix.SlurmVersions[version]
		assert.True(t, exists, "Version %s should have Slurm version mapping", version)
		assert.NotEmpty(t, slurmVersions, "Version %s should have at least one Slurm version", version)
	}
}

func TestClientInterfaces(t *testing.T) {
	ctx := helpers.TestContext(t)
	
	factory, err := NewClientFactory(
		WithBaseURL("https://example.com"),
	)
	helpers.RequireNoError(t, err)
	
	client, err := factory.NewClientWithVersion(ctx, "v0.0.42")
	helpers.AssertNoError(t, err)
	helpers.AssertNotNil(t, client)
	
	// Test that all manager interfaces are available
	jobManager := client.Jobs()
	nodeManager := client.Nodes()
	partitionManager := client.Partitions()
	infoManager := client.Info()
	
	helpers.AssertNotNil(t, jobManager)
	helpers.AssertNotNil(t, nodeManager)
	helpers.AssertNotNil(t, partitionManager)
	helpers.AssertNotNil(t, infoManager)
	
	// Test client version
	version := client.Version()
	helpers.AssertEqual(t, "v0.0.42", version)
}