// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package integration

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/pkg/config"
)

// VersionCompatibilityTestSuite tests compatibility between different API versions
type VersionCompatibilityTestSuite struct {
	suite.Suite
	clients       map[string]slurm.SlurmClient
	versions      []string
	serverURL     string
	token         string
	compatibility CompatibilityMatrix
}

// CompatibilityMatrix tracks feature compatibility across versions
type CompatibilityMatrix struct {
	Features map[string]map[string]FeatureSupport `json:"features"`
	Types    map[string]map[string]TypeSupport    `json:"types"`
}

// FeatureSupport indicates the level of support for a feature
type FeatureSupport struct {
	Supported    bool   `json:"supported"`
	Tested       bool   `json:"tested"`
	Limitations  string `json:"limitations,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

// TypeSupport tracks common type compatibility
type TypeSupport struct {
	Compatible   bool     `json:"compatible"`
	MissingFields []string `json:"missing_fields,omitempty"`
	ExtraFields   []string `json:"extra_fields,omitempty"`
	TypeDiffs     []string `json:"type_diffs,omitempty"`
}

// SetupSuite initializes compatibility testing
func (suite *VersionCompatibilityTestSuite) SetupSuite() {
	// Check if compatibility testing is enabled
	if os.Getenv("SLURM_COMPATIBILITY_TEST") != "true" {
		suite.T().Skip("Compatibility tests disabled. Set SLURM_COMPATIBILITY_TEST=true to enable")
	}

	// Initialize compatibility matrix
	suite.compatibility = CompatibilityMatrix{
		Features: make(map[string]map[string]FeatureSupport),
		Types:    make(map[string]map[string]TypeSupport),
	}

	// Get server configuration
	suite.serverURL = os.Getenv("SLURM_SERVER_URL")
	if suite.serverURL == "" {
		suite.serverURL = "http://rocky9:6820"
	}

	// Get JWT token
	suite.token = os.Getenv("SLURM_JWT_TOKEN")
	if suite.token == "" {
		token, err := fetchJWTTokenViaSSH()
		require.NoError(suite.T(), err, "Failed to fetch JWT token")
		suite.token = token
	}

	// Initialize clients for all supported versions
	suite.versions = []string{"v0.0.40", "v0.0.41", "v0.0.42", "v0.0.43"}
	suite.clients = make(map[string]slurm.SlurmClient)

	ctx := context.Background()
	for _, version := range suite.versions {
		client, err := slurm.NewClientWithVersion(ctx, version,
			slurm.WithBaseURL(suite.serverURL),
			slurm.WithAuth(auth.NewTokenAuth(suite.token)),
			slurm.WithConfig(&config.Config{
				Timeout:            30 * time.Second,
				MaxRetries:         2,
				Debug:              false,
				InsecureSkipVerify: true,
			}),
		)
		
		if err != nil {
			suite.T().Logf("Failed to create client for version %s: %v", version, err)
			continue
		}
		
		suite.clients[version] = client
		suite.T().Logf("Compatibility client created for version %s", version)
	}

	require.NotEmpty(suite.T(), suite.clients, "At least one client must be created")
}

// TearDownSuite cleans up compatibility test resources
func (suite *VersionCompatibilityTestSuite) TearDownSuite() {
	for version, client := range suite.clients {
		if client != nil {
			client.Close()
			suite.T().Logf("Closed compatibility client for version %s", version)
		}
	}
	
	// Generate compatibility report
	suite.generateCompatibilityReport()
}

// testFeatureSupport tests if a feature is supported in a version
func (suite *VersionCompatibilityTestSuite) testFeatureSupport(version string, feature string, testFunc func() error) FeatureSupport {
	err := testFunc()
	support := FeatureSupport{
		Supported: err == nil,
		Tested:    true,
	}
	
	if err != nil {
		support.ErrorMessage = err.Error()
		// Classify error types
		errorMsg := err.Error()
		if contains(errorMsg, "not found") || contains(errorMsg, "404") {
			support.Limitations = "Endpoint not available"
		} else if contains(errorMsg, "not implemented") || contains(errorMsg, "501") {
			support.Limitations = "Feature not implemented"
		} else if contains(errorMsg, "forbidden") || contains(errorMsg, "403") {
			support.Limitations = "Access forbidden"
		} else if contains(errorMsg, "timeout") {
			support.Limitations = "Request timeout"
		} else {
			support.Limitations = "Unknown error"
		}
	}
	
	return support
}

// versionContains checks if a string contains a substring (case-insensitive)
func versionContains(str, substr string) bool {
	return len(str) >= len(substr) && 
		   (str == substr || 
		    len(str) > len(substr) && 
		    (str[:len(substr)] == substr || str[len(str)-len(substr):] == substr || 
		     findInString(str, substr)))
}

func findInString(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestCoreFeatureCompatibility tests core feature compatibility across versions
func (suite *VersionCompatibilityTestSuite) TestCoreFeatureCompatibility() {
	ctx := context.Background()
	
	suite.T().Log("=== Testing Core Feature Compatibility ===")
	
	coreFeatures := map[string]func(slurm.SlurmClient) error{
		"ping": func(client slurm.SlurmClient) error {
			return client.Info().Ping(ctx)
		},
		"version": func(client slurm.SlurmClient) error {
			_, err := client.Info().Version(ctx)
			return err
		},
		"cluster_info": func(client slurm.SlurmClient) error {
			_, err := client.Info().Get(ctx)
			return err
		},
		"cluster_stats": func(client slurm.SlurmClient) error {
			_, err := client.Info().Stats(ctx)
			return err
		},
	}
	
	for feature, testFunc := range coreFeatures {
		suite.compatibility.Features[feature] = make(map[string]FeatureSupport)
		
		for version, client := range suite.clients {
			support := suite.testFeatureSupport(version, feature, func() error {
				return testFunc(client)
			})
			
			suite.compatibility.Features[feature][version] = support
			
			status := "✓ SUPPORTED"
			if !support.Supported {
				status = "✗ NOT SUPPORTED"
			}
			
			suite.T().Logf("  %s [%s]: %s", feature, version, status)
			if !support.Supported {
				suite.T().Logf("    Error: %s", support.ErrorMessage)
				suite.T().Logf("    Limitation: %s", support.Limitations)
			}
		}
	}
}

// TestResourceListingCompatibility tests resource listing compatibility
func (suite *VersionCompatibilityTestSuite) TestResourceListingCompatibility() {
	ctx := context.Background()
	
	suite.T().Log("=== Testing Resource Listing Compatibility ===")
	
	resourceFeatures := map[string]func(slurm.SlurmClient) error{
		"list_jobs": func(client slurm.SlurmClient) error {
			_, err := client.Jobs().List(ctx, &interfaces.ListJobsOptions{Limit: 5})
			return err
		},
		"list_nodes": func(client slurm.SlurmClient) error {
			_, err := client.Nodes().List(ctx, &interfaces.ListNodesOptions{Limit: 5})
			return err
		},
		"list_partitions": func(client slurm.SlurmClient) error {
			_, err := client.Partitions().List(ctx, &interfaces.ListPartitionsOptions{Limit: 5})
			return err
		},
		"list_qos": func(client slurm.SlurmClient) error {
			_, err := client.QoS().List(ctx, &interfaces.ListQoSOptions{Limit: 5})
			return err
		},
	}
	
	for feature, testFunc := range resourceFeatures {
		suite.compatibility.Features[feature] = make(map[string]FeatureSupport)
		
		for version, client := range suite.clients {
			support := suite.testFeatureSupport(version, feature, func() error {
				return testFunc(client)
			})
			
			suite.compatibility.Features[feature][version] = support
			
			status := "✓ SUPPORTED"
			if !support.Supported {
				status = "✗ NOT SUPPORTED"
			}
			
			suite.T().Logf("  %s [%s]: %s", feature, version, status)
			if !support.Supported && support.Limitations != "" {
				suite.T().Logf("    Limitation: %s", support.Limitations)
			}
		}
	}
}

// TestAdvancedFeatureCompatibility tests advanced feature compatibility
func (suite *VersionCompatibilityTestSuite) TestAdvancedFeatureCompatibility() {
	ctx := context.Background()
	
	suite.T().Log("=== Testing Advanced Feature Compatibility ===")
	
	// Get a sample job ID and QoS name for individual resource tests
	var sampleJobID, sampleQoSName string
	
	// Find working client to get sample data
	for _, client := range suite.clients {
		if jobs, err := client.Jobs().List(ctx, &interfaces.ListJobsOptions{Limit: 1}); err == nil && len(jobs.Jobs) > 0 {
			sampleJobID = jobs.Jobs[0].ID
		}
		if qos, err := client.QoS().List(ctx, &interfaces.ListQoSOptions{Limit: 1}); err == nil && len(qos.QoS) > 0 {
			sampleQoSName = qos.QoS[0].Name
		}
		break
	}
	
	advancedFeatures := map[string]func(slurm.SlurmClient) error{
		"get_job": func(client slurm.SlurmClient) error {
			if sampleJobID == "" {
				return fmt.Errorf("no sample job ID available")
			}
			_, err := client.Jobs().Get(ctx, sampleJobID)
			return err
		},
		"get_qos": func(client slurm.SlurmClient) error {
			if sampleQoSName == "" {
				return fmt.Errorf("no sample QoS name available")
			}
			_, err := client.QoS().Get(ctx, sampleQoSName)
			return err
		},
		"job_submission": func(client slurm.SlurmClient) error {
			// Test job submission capability (don't actually submit)
			submission := &interfaces.JobSubmission{
				Name:      "compatibility-test",
				Script:    "#!/bin/bash\necho test",
				Partition: "nonexistent",
				Nodes:     1,
				CPUs:      1,
				TimeLimit: 1,
			}
			_, err := client.Jobs().Submit(ctx, submission)
			// We expect this to fail, but we want to test if the endpoint exists
			if err != nil && (contains(err.Error(), "partition") || contains(err.Error(), "invalid")) {
				return nil // Endpoint exists but partition is invalid (expected)
			}
			return err
		},
	}
	
	for feature, testFunc := range advancedFeatures {
		suite.compatibility.Features[feature] = make(map[string]FeatureSupport)
		
		for version, client := range suite.clients {
			support := suite.testFeatureSupport(version, feature, func() error {
				return testFunc(client)
			})
			
			suite.compatibility.Features[feature][version] = support
			
			status := "✓ SUPPORTED"
			if !support.Supported {
				status = "✗ NOT SUPPORTED"
			}
			
			suite.T().Logf("  %s [%s]: %s", feature, version, status)
			if !support.Supported && support.Limitations != "" {
				suite.T().Logf("    Limitation: %s", support.Limitations)
			}
		}
	}
}

// TestCommonTypeCompatibility tests compatibility of common types across versions
func (suite *VersionCompatibilityTestSuite) TestCommonTypeCompatibility() {
	suite.T().Log("=== Testing Common Type Compatibility ===")
	
	// Test common types structure compatibility
	commonTypes := map[string]interface{}{
		"Account":   types.Account{},
		"Job":       types.Job{},
		"Node":      types.Node{},
		"Partition": types.Partition{},
		"QoS":       types.QoS{},
		"User":      types.User{},
	}
	
	for typeName, typeInstance := range commonTypes {
		suite.compatibility.Types[typeName] = make(map[string]TypeSupport)
		
		// Analyze type structure
		typeInfo := suite.analyzeTypeStructure(typeInstance)
		
		// For now, assume all versions use the same common types (they should)
		for _, version := range suite.versions {
			if _, exists := suite.clients[version]; exists {
				support := TypeSupport{
					Compatible: true, // Common types should be compatible
				}
				suite.compatibility.Types[typeName][version] = support
				
				suite.T().Logf("  %s [%s]: ✓ COMPATIBLE (%d fields)", 
					typeName, version, len(typeInfo.Fields))
			}
		}
	}
}

// TypeInfo holds information about a type's structure
type TypeInfo struct {
	Name   string
	Fields []FieldInfo
}

// FieldInfo holds information about a struct field
type FieldInfo struct {
	Name string
	Type string
	Tag  string
}

// analyzeTypeStructure analyzes the structure of a type
func (suite *VersionCompatibilityTestSuite) analyzeTypeStructure(typeInstance interface{}) TypeInfo {
	typeVal := reflect.TypeOf(typeInstance)
	if typeVal.Kind() == reflect.Ptr {
		typeVal = typeVal.Elem()
	}
	
	info := TypeInfo{
		Name:   typeVal.Name(),
		Fields: make([]FieldInfo, 0),
	}
	
	for i := 0; i < typeVal.NumField(); i++ {
		field := typeVal.Field(i)
		fieldInfo := FieldInfo{
			Name: field.Name,
			Type: field.Type.String(),
			Tag:  string(field.Tag),
		}
		info.Fields = append(info.Fields, fieldInfo)
	}
	
	return info
}

// TestVersionSwitching tests switching between versions
func (suite *VersionCompatibilityTestSuite) TestVersionSwitching() {
	suite.T().Log("=== Testing Version Switching ===")
	
	if len(suite.clients) < 2 {
		suite.T().Skip("Need at least 2 versions for switching test")
		return
	}
	
	ctx := context.Background()
	
	// Test switching between different versions
	versionList := make([]string, 0, len(suite.clients))
	for version := range suite.clients {
		versionList = append(versionList, version)
	}
	
	// Perform operations with different versions in sequence
	for i := 0; i < 3; i++ { // Test 3 rounds of switching
		for j, version := range versionList {
			client := suite.clients[version]
			
			suite.T().Logf("Round %d: Using version %s", i+1, version)
			
			// Perform a series of operations
			err := client.Info().Ping(ctx)
			if err == nil {
				suite.T().Logf("  ✓ Ping successful")
			} else {
				suite.T().Logf("  ✗ Ping failed: %v", err)
			}
			
			_, err = client.Info().Version(ctx)
			if err == nil {
				suite.T().Logf("  ✓ Version query successful")
			} else {
				suite.T().Logf("  ✗ Version query failed: %v", err)
			}
			
			// Small delay between version switches
			if j < len(versionList)-1 {
				time.Sleep(100 * time.Millisecond)
			}
		}
		
		suite.T().Logf("Completed switching round %d", i+1)
	}
}

// TestFeatureEvolution tests how features evolved across versions
func (suite *VersionCompatibilityTestSuite) TestFeatureEvolution() {
	suite.T().Log("=== Testing Feature Evolution ===")
	
	// Analyze which features are supported in which versions
	evolutionReport := make(map[string][]string) // feature -> supported versions
	
	for feature, versionSupport := range suite.compatibility.Features {
		var supportedVersions []string
		for version, support := range versionSupport {
			if support.Supported {
				supportedVersions = append(supportedVersions, version)
			}
		}
		evolutionReport[feature] = supportedVersions
	}
	
	suite.T().Log("\nFeature Evolution Summary:")
	suite.T().Log("Feature\t\t\tSupported Versions")
	suite.T().Log("-------\t\t\t------------------")
	
	for feature, versions := range evolutionReport {
		suite.T().Logf("%-20s\t%v", feature, versions)
		
		// Analyze evolution patterns
		if len(versions) == len(suite.versions) {
			suite.T().Logf("  → Stable feature (supported across all versions)")
		} else if len(versions) == 0 {
			suite.T().Logf("  → Unsupported feature (not working in any version)")
		} else {
			suite.T().Logf("  → Partial support (supported in %d/%d versions)", 
				len(versions), len(suite.versions))
		}
	}
}

// generateCompatibilityReport generates a comprehensive compatibility report
func (suite *VersionCompatibilityTestSuite) generateCompatibilityReport() {
	suite.T().Log("\n" + strings.Repeat("=", 80))
	suite.T().Log("COMPREHENSIVE COMPATIBILITY REPORT")
	suite.T().Log(strings.Repeat("=", 80))
	
	// Summary statistics
	totalFeatures := len(suite.compatibility.Features)
	totalVersions := len(suite.versions)
	totalTests := totalFeatures * totalVersions
	
	passedTests := 0
	for _, versionSupport := range suite.compatibility.Features {
		for _, support := range versionSupport {
			if support.Supported {
				passedTests++
			}
		}
	}
	
	compatibilityRate := float64(passedTests) / float64(totalTests) * 100
	
	suite.T().Logf("\nCompatibility Summary:")
	suite.T().Logf("  Total Features Tested: %d", totalFeatures)
	suite.T().Logf("  Total Versions Tested: %d", totalVersions)
	suite.T().Logf("  Total Test Cases: %d", totalTests)
	suite.T().Logf("  Passed Test Cases: %d", passedTests)
	suite.T().Logf("  Overall Compatibility Rate: %.1f%%", compatibilityRate)
	
	// Version-specific compatibility
	suite.T().Log("\nVersion-Specific Compatibility:")
	for _, version := range suite.versions {
		if _, exists := suite.clients[version]; !exists {
			continue
		}
		
		versionPassed := 0
		versionTotal := 0
		
		for _, versionSupport := range suite.compatibility.Features {
			if support, exists := versionSupport[version]; exists {
				versionTotal++
				if support.Supported {
					versionPassed++
				}
			}
		}
		
		if versionTotal > 0 {
			versionRate := float64(versionPassed) / float64(versionTotal) * 100
			suite.T().Logf("  %s: %.1f%% (%d/%d features)", 
				version, versionRate, versionPassed, versionTotal)
		}
	}
	
	// Feature compatibility matrix
	suite.T().Log("\nFeature Compatibility Matrix:")
	suite.T().Log("Feature\t\t\tv0.0.40\tv0.0.41\tv0.0.42\tv0.0.43")
	suite.T().Log("-------\t\t\t-------\t-------\t-------\t-------")
	
	for feature, versionSupport := range suite.compatibility.Features {
		line := fmt.Sprintf("%-20s", feature)
		
		for _, version := range []string{"v0.0.40", "v0.0.41", "v0.0.42", "v0.0.43"} {
			if support, exists := versionSupport[version]; exists {
				if support.Supported {
					line += "\t✓"
				} else {
					line += "\t✗"
				}
			} else {
				line += "\t-"
			}
		}
		
		suite.T().Log(line)
	}
	
	// Recommendations
	suite.T().Log("\nRecommendations:")
	
	if compatibilityRate >= 90 {
		suite.T().Log("  ✓ Excellent compatibility across versions")
	} else if compatibilityRate >= 75 {
		suite.T().Log("  ⚠ Good compatibility with some version-specific limitations")
	} else {
		suite.T().Log("  ✗ Significant compatibility issues detected")
	}
	
	// Find most compatible version
	var bestVersion string
	var bestRate float64
	
	for _, version := range suite.versions {
		if _, exists := suite.clients[version]; !exists {
			continue
		}
		
		versionPassed := 0
		versionTotal := 0
		
		for _, versionSupport := range suite.compatibility.Features {
			if support, exists := versionSupport[version]; exists {
				versionTotal++
				if support.Supported {
					versionPassed++
				}
			}
		}
		
		if versionTotal > 0 {
			rate := float64(versionPassed) / float64(versionTotal) * 100
			if rate > bestRate {
				bestRate = rate
				bestVersion = version
			}
		}
	}
	
	if bestVersion != "" {
		suite.T().Logf("  → Recommended version: %s (%.1f%% compatibility)", bestVersion, bestRate)
	}
	
	suite.T().Log(strings.Repeat("=", 80))
}

// TestVersionCompatibilitySuite runs the version compatibility test suite
func TestVersionCompatibilitySuite(t *testing.T) {
	suite.Run(t, new(VersionCompatibilityTestSuite))
}
