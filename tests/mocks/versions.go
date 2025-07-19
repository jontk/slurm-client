package mocks

import "time"

// VersionConfigs holds configurations for different API versions
var VersionConfigs = map[string]*ServerConfig{
	"v0.0.40": {
		APIVersion:          "v0.0.40",
		SlurmVersion:        "24.05",
		EnableAuth:          false,
		AuthToken:           "test-token-v40",
		ErrorResponses:      make(map[string]ErrorResponse),
		ResponseDelay:       0,
		JobIDCounter:        1000,
		SupportedOperations: v040SupportedOperations(),
	},
	"v0.0.41": {
		APIVersion:          "v0.0.41",
		SlurmVersion:        "24.11",
		EnableAuth:          false,
		AuthToken:           "test-token-v41",
		ErrorResponses:      make(map[string]ErrorResponse),
		ResponseDelay:       0,
		JobIDCounter:        1000,
		SupportedOperations: v041SupportedOperations(),
	},
	"v0.0.42": {
		APIVersion:          "v0.0.42",
		SlurmVersion:        "25.05",
		EnableAuth:          false,
		AuthToken:           "test-token-v42",
		ErrorResponses:      make(map[string]ErrorResponse),
		ResponseDelay:       0,
		JobIDCounter:        1000,
		SupportedOperations: v042SupportedOperations(),
	},
	"v0.0.43": {
		APIVersion:          "v0.0.43",
		SlurmVersion:        "25.05",
		EnableAuth:          false,
		AuthToken:           "test-token-v43",
		ErrorResponses:      make(map[string]ErrorResponse),
		ResponseDelay:       0,
		JobIDCounter:        1000,
		SupportedOperations: v043SupportedOperations(),
	},
}

// v040SupportedOperations returns supported operations for API v0.0.40
func v040SupportedOperations() map[string]bool {
	return map[string]bool{
		"jobs.list":         true,
		"jobs.get":          true,
		"jobs.submit":       true,
		"jobs.cancel":       true,
		"jobs.update":       false, // Limited update support in v0.0.40
		"jobs.steps":        true,
		"nodes.list":        true,
		"nodes.get":         true,
		"nodes.update":      false, // No node update in v0.0.40
		"partitions.list":   true,
		"partitions.get":    true,
		"partitions.update": false, // No partition update in v0.0.40
		"info.get":          true,
		"info.ping":         true,
		"info.stats":        true,
		"info.version":      true,
		// Basic analytics support in v0.0.40
		"jobs.utilization":      true,
		"jobs.efficiency":       true,
		"jobs.performance":      true,
	}
}

// v041SupportedOperations returns supported operations for API v0.0.41
func v041SupportedOperations() map[string]bool {
	return map[string]bool{
		"jobs.list":         true,
		"jobs.get":          true,
		"jobs.submit":       true,
		"jobs.cancel":       true,
		"jobs.update":       true, // Enhanced update support in v0.0.41
		"jobs.steps":        true,
		"nodes.list":        true,
		"nodes.get":         true,
		"nodes.update":      false, // Limited node update in v0.0.41
		"partitions.list":   true,
		"partitions.get":    true,
		"partitions.update": false, // No partition update in v0.0.41
		"info.get":          true,
		"info.ping":         true,
		"info.stats":        true,
		"info.version":      true,
		// Enhanced analytics support in v0.0.41
		"jobs.utilization":      true,
		"jobs.efficiency":       true,
		"jobs.performance":      true,
		"jobs.live_metrics":     true,
		"jobs.resource_trends":  true,
		"jobs.step_utilization": true,
	}
}

// v042SupportedOperations returns supported operations for API v0.0.42 (stable)
func v042SupportedOperations() map[string]bool {
	return map[string]bool{
		"jobs.list":         true,
		"jobs.get":          true,
		"jobs.submit":       true,
		"jobs.cancel":       true,
		"jobs.update":       true,
		"jobs.steps":        true,
		"nodes.list":        true,
		"nodes.get":         true,
		"nodes.update":      true, // Full node update support in v0.0.42
		"partitions.list":   true,
		"partitions.get":    true,
		"partitions.update": true, // Full partition update support in v0.0.42
		"info.get":          true,
		"info.ping":         true,
		"info.stats":        true,
		"info.version":      true,
		// Full analytics support in v0.0.42
		"jobs.utilization":         true,
		"jobs.efficiency":          true,
		"jobs.performance":         true,
		"jobs.live_metrics":        true,
		"jobs.resource_trends":     true,
		"jobs.step_utilization":    true,
		"jobs.performance_history": true,
		"jobs.performance_trends":  true,
		"jobs.efficiency_trends":   true,
		"jobs.compare_performance": true,
		"jobs.similar_performance": true,
		"jobs.analyze_batch":       true,
		"jobs.workflow_performance": true,
		"jobs.efficiency_report":   true,
	}
}

// v043SupportedOperations returns supported operations for API v0.0.43 (latest)
func v043SupportedOperations() map[string]bool {
	return map[string]bool{
		"jobs.list":         true,
		"jobs.get":          true,
		"jobs.submit":       true,
		"jobs.cancel":       true,
		"jobs.update":       true,
		"jobs.steps":        true,
		"nodes.list":        true,
		"nodes.get":         true,
		"nodes.update":      true,
		"partitions.list":   true,
		"partitions.get":    true,
		"partitions.update": true,
		"info.get":          true,
		"info.ping":         true,
		"info.stats":        true,
		"info.version":      true,
		// Complete analytics support in v0.0.43
		"jobs.utilization":         true,
		"jobs.efficiency":          true,
		"jobs.performance":         true,
		"jobs.live_metrics":        true,
		"jobs.resource_trends":     true,
		"jobs.step_utilization":    true,
		"jobs.performance_history": true,
		"jobs.performance_trends":  true,
		"jobs.efficiency_trends":   true,
		"jobs.compare_performance": true,
		"jobs.similar_performance": true,
		"jobs.analyze_batch":       true,
		"jobs.workflow_performance": true,
		"jobs.efficiency_report":   true,
		// Additional features in v0.0.43
		"reservations.list": false, // Would be true if we implemented reservations
		"qos.list":          false, // Would be true if we implemented QoS
	}
}

// NewMockSlurmServerForVersion creates a mock server for a specific API version
func NewMockSlurmServerForVersion(version string) *MockSlurmServer {
	config, exists := VersionConfigs[version]
	if !exists {
		// Default to v0.0.42 if version not found
		config = VersionConfigs["v0.0.42"]
	}

	// Clone config to avoid modifying the original
	configCopy := *config
	configCopy.ErrorResponses = make(map[string]ErrorResponse)

	return NewMockSlurmServer(&configCopy)
}

// GetVersionDifferences returns the differences between API versions
func GetVersionDifferences() map[string]map[string]interface{} {
	return map[string]map[string]interface{}{
		"v0.0.40": {
			"breaking_changes": []string{
				"Limited job update capabilities",
				"No node update support",
				"No partition update support",
			},
			"new_features": []string{
				"Basic job management",
				"Node information retrieval",
				"Partition information retrieval",
			},
			"field_changes": map[string]string{
				"jobs.minimum_switches": "Renamed to required_switches in v0.0.41",
			},
		},
		"v0.0.41": {
			"breaking_changes": []string{
				"minimum_switches field renamed to required_switches",
				"Some job output fields restructured",
			},
			"new_features": []string{
				"Enhanced job update capabilities",
				"Improved error handling",
				"Extended job submission options",
			},
			"field_changes": map[string]string{
				"jobs.minimum_switches": "Renamed from minimum_switches",
				"jobs.required_switches": "New field replacing minimum_switches",
			},
		},
		"v0.0.42": {
			"breaking_changes": []string{
				"exclusive and oversubscribe fields removed from job outputs",
				"Some endpoint URL patterns changed",
			},
			"new_features": []string{
				"Full node update support",
				"Full partition update support",
				"Enhanced filtering capabilities",
				"Improved performance optimizations",
			},
			"field_changes": map[string]string{
				"jobs.exclusive":     "Removed field",
				"jobs.oversubscribe": "Removed field",
			},
		},
		"v0.0.43": {
			"breaking_changes": []string{
				"FrontEnd mode support removed",
				"Some legacy endpoint deprecations",
			},
			"new_features": []string{
				"Reservation management support",
				"Enhanced QoS integration",
				"New cluster statistics endpoints",
				"WebSocket support for real-time updates",
			},
			"field_changes": map[string]string{
				"cluster.frontend_mode": "Removed field",
			},
		},
	}
}

// CreateVersionSpecificErrorScenarios creates error scenarios for testing version differences
func CreateVersionSpecificErrorScenarios() map[string]map[string]ErrorResponse {
	return map[string]map[string]ErrorResponse{
		"v0.0.40": {
			"PATCH /slurm/v0.0.40/job/1001": {
				StatusCode: 501,
				Body: map[string]string{
					"error": "Job update not fully supported in v0.0.40",
				},
			},
			"PATCH /slurm/v0.0.40/node/node001": {
				StatusCode: 501,
				Body: map[string]string{
					"error": "Node update not supported in v0.0.40",
				},
			},
		},
		"v0.0.41": {
			"GET /slurm/v0.0.41/job/1001": {
				StatusCode: 400,
				Body: map[string]interface{}{
					"error": "Field minimum_switches deprecated, use required_switches",
					"code":  "DEPRECATED_FIELD",
				},
			},
		},
		"v0.0.42": {
			"GET /slurm/v0.0.42/job/1001": {
				StatusCode: 400,
				Body: map[string]interface{}{
					"error": "Fields exclusive and oversubscribe no longer supported",
					"code":  "REMOVED_FIELD",
				},
			},
		},
		"v0.0.43": {
			"GET /slurm/v0.0.43/cluster/frontend": {
				StatusCode: 404,
				Body: map[string]interface{}{
					"error": "FrontEnd mode no longer supported",
					"code":  "FEATURE_REMOVED",
				},
			},
		},
	}
}

// MockServerPool manages multiple mock servers for different versions
type MockServerPool struct {
	servers map[string]*MockSlurmServer
}

// NewMockServerPool creates a pool of mock servers for all versions
func NewMockServerPool() *MockServerPool {
	pool := &MockServerPool{
		servers: make(map[string]*MockSlurmServer),
	}

	for version := range VersionConfigs {
		pool.servers[version] = NewMockSlurmServerForVersion(version)
	}

	return pool
}

// GetServer returns a mock server for the specified version
func (p *MockServerPool) GetServer(version string) *MockSlurmServer {
	return p.servers[version]
}

// GetURL returns the URL for a specific version's mock server
func (p *MockServerPool) GetURL(version string) string {
	if server := p.servers[version]; server != nil {
		return server.URL()
	}
	return ""
}

// Close closes all mock servers in the pool
func (p *MockServerPool) Close() {
	for _, server := range p.servers {
		server.Close()
	}
}

// SetResponseDelay sets a response delay for all servers (useful for timeout testing)
func (p *MockServerPool) SetResponseDelay(delay time.Duration) {
	for _, server := range p.servers {
		server.config.ResponseDelay = delay
	}
}

// EnableAuth enables authentication for all servers
func (p *MockServerPool) EnableAuth(enable bool) {
	for _, server := range p.servers {
		server.config.EnableAuth = enable
	}
}