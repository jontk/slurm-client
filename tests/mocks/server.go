// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package mocks

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"

	// Import generated builders
	builderv0040 "github.com/jontk/slurm-client/tests/mocks/generated/v0_0_40"
	builderv0042 "github.com/jontk/slurm-client/tests/mocks/generated/v0_0_42"
	builderv0043 "github.com/jontk/slurm-client/tests/mocks/generated/v0_0_43"
	builderv0044 "github.com/jontk/slurm-client/tests/mocks/generated/v0_0_44"
)

// MockSlurmServer represents a mock SLURM REST API server for testing
type MockSlurmServer struct {
	server     *httptest.Server
	router     *mux.Router
	storage    *MockStorage
	config     *ServerConfig
	jobHandler JobHandler // Version-specific job handler
	mu         sync.RWMutex
}

// ServerConfig holds configuration for the mock server
type ServerConfig struct {
	APIVersion          string
	SlurmVersion        string
	EnableAuth          bool
	AuthToken           string
	ErrorResponses      map[string]ErrorResponse
	ResponseDelay       time.Duration
	JobIDCounter        int64
	SupportedOperations map[string]bool
}

// ErrorResponse represents a mock error response
type ErrorResponse struct {
	StatusCode int
	Body       interface{}
}

// MockStorage holds the mock data
// Jobs now stores version-specific OpenAPI types (interface{} to support all versions)
type MockStorage struct {
	Jobs        map[string]interface{} // Can be *v0_0_40.V0040JobInfo, *v0_0_42.V0042JobInfo, etc.
	Nodes       map[string]*MockNode
	Partitions  map[string]*MockPartition
	ClusterInfo *MockClusterInfo
	mu          sync.RWMutex
	version     string // API version for this storage
}

// MockJob represents a mock job
// Deprecated: MockJob is being phased out in favor of using real OpenAPI types
// (e.g., V0040JobInfo) which are created via generated builders. This eliminates
// field name mismatches and ensures type safety. The mock server storage now holds
// interface{} which can be real OpenAPI types. Consider using generated builders
// from tests/mocks/generated/v0_0_40 instead.
type MockJob struct {
	JobID       int32             `json:"job_id"`
	Name        string            `json:"name"`
	UserID      int32             `json:"user_id"`
	State       string            `json:"state"`
	Partition   string            `json:"partition"`
	CPUs        int               `json:"cpus"`
	Memory      int64             `json:"memory"`
	TimeLimit   int               `json:"time_limit"`
	SubmitTime  int64             `json:"submit_time"`
	StartTime   *int64            `json:"start_time,omitempty"`
	EndTime     *int64            `json:"end_time,omitempty"`
	ExitCode    *int              `json:"exit_code,omitempty"`
	WorkingDir  string            `json:"working_directory"`
	Script      string            `json:"script"`
	Environment map[string]string `json:"environment"`
}

// MockNode represents a mock node
type MockNode struct {
	Name         string   `json:"name"`
	State        []string `json:"state"` // Array of state flags (e.g., ["IDLE"], ["ALLOCATED", "COMPLETING"])
	CPUs         int      `json:"cpus"`
	Memory       int64    `json:"memory"`
	Features     []string `json:"features"`
	Partition    string   `json:"partition"`
	Architecture string   `json:"architecture"`
	OS           string   `json:"operating_system"`
}

// MockPartition represents a mock partition
type MockPartition struct {
	Name        string `json:"name"`
	State       string `json:"state"`
	MaxNodes    int    `json:"max_nodes"`
	MinNodes    int    `json:"min_nodes"`
	DefaultTime int    `json:"default_time"`
	MaxTime     int    `json:"max_time"`
	Nodes       *struct {
		AllowedAllocation *string `json:"allowed_allocation,omitempty"`
		Configured        *string `json:"configured,omitempty"`
		Total             *int32  `json:"total,omitempty"`
	} `json:"nodes,omitempty"`
	Features []string `json:"features"`
}

// MockClusterInfo represents mock cluster information
type MockClusterInfo struct {
	ClusterName   string `json:"cluster_name"`
	SlurmVersion  string `json:"slurm_version"`
	APIVersion    string `json:"api_version"`
	TotalNodes    int    `json:"total_nodes"`
	TotalCPUs     int    `json:"total_cpus"`
	RunningJobs   int    `json:"running_jobs"`
	PendingJobs   int    `json:"pending_jobs"`
	LastHeartbeat int64  `json:"last_heartbeat"`
}

// NewMockSlurmServer creates a new mock SLURM server
func NewMockSlurmServer(config *ServerConfig) *MockSlurmServer {
	if config == nil {
		config = DefaultServerConfig()
	}

	storage := &MockStorage{
		Jobs:       make(map[string]interface{}),
		Nodes:      make(map[string]*MockNode),
		Partitions: make(map[string]*MockPartition),
		version:    config.APIVersion,
		ClusterInfo: &MockClusterInfo{
			ClusterName:   "test-cluster",
			SlurmVersion:  config.SlurmVersion,
			APIVersion:    config.APIVersion,
			TotalNodes:    4,
			TotalCPUs:     64,
			RunningJobs:   0,
			PendingJobs:   0,
			LastHeartbeat: time.Now().Unix(),
		},
	}

	// Initialize default data
	storage.initializeDefaultData()

	mock := &MockSlurmServer{
		storage:    storage,
		config:     config,
		jobHandler: NewJobHandler(config.APIVersion),
	}

	mock.setupRouter()
	mock.server = httptest.NewServer(mock.router)

	return mock
}

// DefaultServerConfig returns a default server configuration
func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		APIVersion:     "v0.0.42",
		SlurmVersion:   "25.05",
		EnableAuth:     false,
		AuthToken:      "test-token-v42",
		ErrorResponses: make(map[string]ErrorResponse),
		ResponseDelay:  0,
		JobIDCounter:   1000,
		SupportedOperations: map[string]bool{
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
			// Analytics endpoints
			"jobs.utilization":          true,
			"jobs.efficiency":           true,
			"jobs.performance":          true,
			"jobs.live_metrics":         true,
			"jobs.resource_trends":      true,
			"jobs.step_utilization":     true,
			"jobs.performance_history":  true,
			"jobs.performance_trends":   true,
			"jobs.efficiency_trends":    true,
			"jobs.compare_performance":  true,
			"jobs.similar_performance":  true,
			"jobs.analyze_batch":        true,
			"jobs.workflow_performance": true,
			"jobs.efficiency_report":    true,
		},
	}
}

// URL returns the mock server URL
func (m *MockSlurmServer) URL() string {
	return m.server.URL
}

// Close shuts down the mock server
func (m *MockSlurmServer) Close() {
	m.server.Close()
}

// SetError configures an error response for a specific endpoint
func (m *MockSlurmServer) SetError(endpoint string, statusCode int, body interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.config.ErrorResponses[endpoint] = ErrorResponse{
		StatusCode: statusCode,
		Body:       body,
	}
}

// ClearError removes an error response for a specific endpoint
func (m *MockSlurmServer) ClearError(endpoint string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.config.ErrorResponses, endpoint)
}

// GetConfig returns the server configuration
func (m *MockSlurmServer) GetConfig() *ServerConfig {
	return m.config
}

// GetJob returns a job from storage
// Deprecated: This function works with MockJob which is being phased out.
// Access m.storage.Jobs directly with proper type assertions instead.
func (m *MockSlurmServer) GetJob(jobID string) *MockJob {
	m.storage.mu.RLock()
	defer m.storage.mu.RUnlock()

	// Try to get job from storage
	if job, exists := m.storage.Jobs[jobID]; exists {
		// If it's already a MockJob, return it
		if mockJob, ok := job.(*MockJob); ok {
			return mockJob
		}
		// TODO: Convert from real OpenAPI type to MockJob for backward compatibility
		// For now, return nil as these are incompatible types
	}
	return nil
}

// AddJob adds a job to storage
// Deprecated: This function works with MockJob which is being phased out.
// Use generated builders to create jobs and add them directly to storage.
func (m *MockSlurmServer) AddJob(job *MockJob) {
	m.storage.mu.Lock()
	defer m.storage.mu.Unlock()
	m.storage.Jobs[strconv.Itoa(int(job.JobID))] = job
}

// UpdateJobState updates a job's state
// Deprecated: This function works with MockJob which is being phased out.
// Access m.storage.Jobs directly with proper type assertions to update state.
func (m *MockSlurmServer) UpdateJobState(jobID, state string) {
	m.storage.mu.Lock()
	defer m.storage.mu.Unlock()

	if job, exists := m.storage.Jobs[jobID]; exists {
		// If it's a MockJob, update it
		if mockJob, ok := job.(*MockJob); ok {
			mockJob.State = state
			now := time.Now().Unix()
			switch state {
			case "RUNNING":
				mockJob.StartTime = &now
			case "COMPLETED", "FAILED", "CANCELLED":
				mockJob.EndTime = &now
			}
		}
		// TODO: Handle real OpenAPI types
	}
}

// setupRouter configures the HTTP router with all endpoints
func (m *MockSlurmServer) setupRouter() {
	m.router = mux.NewRouter().StrictSlash(false)
	// Add middleware
	m.router.Use(m.trailingSlashMiddleware)
	m.router.Use(m.loggingMiddleware)
	m.router.Use(m.authMiddleware)
	m.router.Use(m.delayMiddleware)
	m.router.Use(m.errorMiddleware)

	// API version prefix - make trailing slash optional on the prefix
	apiRouter := m.router.PathPrefix("/slurm/" + m.config.APIVersion).Subrouter().StrictSlash(false)
	// Apply trailing slash middleware to subrouter as well
	apiRouter.Use(m.trailingSlashMiddleware)

	// Job endpoints - register both with and without trailing slash for compatibility
	apiRouter.HandleFunc("/jobs", m.handleJobsList).Methods("GET")
	apiRouter.HandleFunc("/jobs/", m.handleJobsList).Methods("GET")
	apiRouter.HandleFunc("/job/{job_id}", m.handleJobsGet).Methods("GET")
	apiRouter.HandleFunc("/job/submit", m.handleJobsSubmit).Methods("POST")
	apiRouter.HandleFunc("/job/{job_id}", m.handleJobsCancel).Methods("DELETE")
	apiRouter.HandleFunc("/job/{job_id}", m.handleJobsUpdate).Methods("PATCH", "POST") // Support both PATCH and POST for updates
	apiRouter.HandleFunc("/job/{job_id}/steps/", m.handleJobsSteps).Methods("GET")

	// Job Analytics endpoints - parameterized routes without trailing slash, collections with trailing slash
	apiRouter.HandleFunc("/job/{job_id}/utilization", m.handleJobUtilization).Methods("GET")
	apiRouter.HandleFunc("/job/{job_id}/efficiency", m.handleJobEfficiency).Methods("GET")
	apiRouter.HandleFunc("/job/{job_id}/performance", m.handleJobPerformance).Methods("GET")
	apiRouter.HandleFunc("/job/{job_id}/live_metrics", m.handleJobLiveMetrics).Methods("GET")
	apiRouter.HandleFunc("/job/{job_id}/resource_trends", m.handleJobResourceTrends).Methods("GET")
	apiRouter.HandleFunc("/job/{job_id}/step/{step_id}/utilization", m.handleJobStepUtilization).Methods("GET")
	apiRouter.HandleFunc("/jobs/performance/history", m.handleJobPerformanceHistory).Methods("GET")
	apiRouter.HandleFunc("/jobs/performance/trends", m.handlePerformanceTrends).Methods("GET")
	apiRouter.HandleFunc("/jobs/efficiency/trends", m.handleUserEfficiencyTrends).Methods("GET")
	apiRouter.HandleFunc("/jobs/performance/compare", m.handleCompareJobPerformance).Methods("POST")
	apiRouter.HandleFunc("/jobs/performance/similar", m.handleSimilarJobsPerformance).Methods("GET")
	apiRouter.HandleFunc("/jobs/performance/analyze_batch", m.handleAnalyzeBatchJobs).Methods("POST")
	apiRouter.HandleFunc("/jobs/workflow/performance", m.handleWorkflowPerformance).Methods("GET")
	apiRouter.HandleFunc("/jobs/efficiency/report", m.handleGenerateEfficiencyReport).Methods("GET")

	// Node endpoints - collection has trailing slash, parameterized don't
	apiRouter.HandleFunc("/nodes/", m.handleNodesList).Methods("GET")
	apiRouter.HandleFunc("/node/{node_name}", m.handleNodesGet).Methods("GET")
	apiRouter.HandleFunc("/node/{node_name}", m.handleNodesUpdate).Methods("PATCH", "POST") // Support both PATCH and POST for updates

	// Partition endpoints - collection has trailing slash, parameterized don't
	apiRouter.HandleFunc("/partitions/", m.handlePartitionsList).Methods("GET")
	apiRouter.HandleFunc("/partition/{partition_name}", m.handlePartitionsGet).Methods("GET")
	apiRouter.HandleFunc("/partition/{partition_name}", m.handlePartitionsUpdate).Methods("PATCH", "POST") // Support both PATCH and POST for updates

	// Info endpoints - register with trailing slashes to match OpenAPI client
	apiRouter.HandleFunc("/diag/", m.handleInfoGet).Methods("GET")
	apiRouter.HandleFunc("/ping/", m.handleInfoPing).Methods("GET")
	apiRouter.HandleFunc("/stats/", m.handleInfoStats).Methods("GET")
	apiRouter.HandleFunc("/", m.handleInfoVersion).Methods("GET")
}

// initializeDefaultData sets up default test data using generated builders
func (s *MockStorage) initializeDefaultData() {
	// Initialize jobs based on API version
	switch s.version {
	case "v0.0.40":
		s.initializeV0040Jobs()
	case "v0.0.41":
		s.initializeV0041Jobs()
	case "v0.0.42":
		s.initializeV0042Jobs()
	case "v0.0.43":
		s.initializeV0043Jobs()
	case "v0.0.44":
		s.initializeV0044Jobs()
	default:
		// Fallback to v0.0.40 for unknown versions
		s.initializeV0040Jobs()
	}

	// Initialize nodes and partitions (version-independent for now)
	s.initializeNodesAndPartitions()
}

// initializeV0040Jobs initializes jobs using v0.0.40 generated builders
func (s *MockStorage) initializeV0040Jobs() {
	// Job 1001: Running job
	s.Jobs["1001"] = builderv0040.NewJobInfo().
		WithJobId(1001).
		WithName("test-job-1").
		WithUserId(1000).
		WithJobState("RUNNING").
		WithPartition("compute").
		WithCpus(4).
		WithMemoryPerNode(4 * 1024 * 1024 * 1024).
		WithTimeLimit(60).
		WithSubmitTime(time.Now().Add(-10 * time.Minute).Unix()).
		WithCurrentWorkingDirectory("/tmp").
		WithCommand("#!/bin/bash\necho 'Hello World'").
		Build()

	// Job 1002: Pending job
	s.Jobs["1002"] = builderv0040.NewJobInfo().
		WithJobId(1002).
		WithName("test-job-2").
		WithUserId(1000).
		WithJobState("PENDING").
		WithPartition("gpu").
		WithCpus(8).
		WithMemoryPerNode(8 * 1024 * 1024 * 1024).
		WithTimeLimit(120).
		WithSubmitTime(time.Now().Add(-5 * time.Minute).Unix()).
		WithCurrentWorkingDirectory("/home/testuser").
		WithCommand("#!/bin/bash\npython train_model.py").
		Build()
}

// initializeV0041Jobs initializes jobs using v0.0.41 (fallback to v0.0.40 for now)
func (s *MockStorage) initializeV0041Jobs() {
	// v0.0.41 builders not yet generated, use v0.0.40 as fallback
	s.initializeV0040Jobs()
}

// initializeV0042Jobs initializes jobs using v0.0.42 generated builders
func (s *MockStorage) initializeV0042Jobs() {
	// Job 1001: Running job
	s.Jobs["1001"] = builderv0042.NewJobInfo().
		WithJobId(1001).
		WithName("test-job-1").
		WithUserId(1000).
		WithJobState("RUNNING").
		WithPartition("compute").
		WithCpus(4).
		WithMemoryPerNode(4 * 1024 * 1024 * 1024).
		WithTimeLimit(60).
		WithSubmitTime(time.Now().Add(-10 * time.Minute).Unix()).
		WithCurrentWorkingDirectory("/tmp").
		WithCommand("#!/bin/bash\necho 'Hello World'").
		Build()

	// Job 1002: Pending job
	s.Jobs["1002"] = builderv0042.NewJobInfo().
		WithJobId(1002).
		WithName("test-job-2").
		WithUserId(1000).
		WithJobState("PENDING").
		WithPartition("gpu").
		WithCpus(8).
		WithMemoryPerNode(8 * 1024 * 1024 * 1024).
		WithTimeLimit(120).
		WithSubmitTime(time.Now().Add(-5 * time.Minute).Unix()).
		WithCurrentWorkingDirectory("/home/testuser").
		WithCommand("#!/bin/bash\npython train_model.py").
		Build()
}

// initializeV0043Jobs initializes jobs using v0.0.43 generated builders
func (s *MockStorage) initializeV0043Jobs() {
	// Job 1001: Running job
	s.Jobs["1001"] = builderv0043.NewJobInfo().
		WithJobId(1001).
		WithName("test-job-1").
		WithUserId(1000).
		WithJobState("RUNNING").
		WithPartition("compute").
		WithCpus(4).
		WithMemoryPerNode(4 * 1024 * 1024 * 1024).
		WithTimeLimit(60).
		WithSubmitTime(time.Now().Add(-10 * time.Minute).Unix()).
		WithCurrentWorkingDirectory("/tmp").
		WithCommand("#!/bin/bash\necho 'Hello World'").
		Build()

	// Job 1002: Pending job
	s.Jobs["1002"] = builderv0043.NewJobInfo().
		WithJobId(1002).
		WithName("test-job-2").
		WithUserId(1000).
		WithJobState("PENDING").
		WithPartition("gpu").
		WithCpus(8).
		WithMemoryPerNode(8 * 1024 * 1024 * 1024).
		WithTimeLimit(120).
		WithSubmitTime(time.Now().Add(-5 * time.Minute).Unix()).
		WithCurrentWorkingDirectory("/home/testuser").
		WithCommand("#!/bin/bash\npython train_model.py").
		Build()
}

// initializeV0044Jobs initializes jobs using v0.0.44 generated builders
func (s *MockStorage) initializeV0044Jobs() {
	// Job 1001: Running job
	s.Jobs["1001"] = builderv0044.NewJobInfo().
		WithJobId(1001).
		WithName("test-job-1").
		WithUserId(1000).
		WithJobState("RUNNING").
		WithPartition("compute").
		WithCpus(4).
		WithMemoryPerNode(4 * 1024 * 1024 * 1024).
		WithTimeLimit(60).
		WithSubmitTime(time.Now().Add(-10 * time.Minute).Unix()).
		WithCurrentWorkingDirectory("/tmp").
		WithCommand("#!/bin/bash\necho 'Hello World'").
		Build()

	// Job 1002: Pending job
	s.Jobs["1002"] = builderv0044.NewJobInfo().
		WithJobId(1002).
		WithName("test-job-2").
		WithUserId(1000).
		WithJobState("PENDING").
		WithPartition("gpu").
		WithCpus(8).
		WithMemoryPerNode(8 * 1024 * 1024 * 1024).
		WithTimeLimit(120).
		WithSubmitTime(time.Now().Add(-5 * time.Minute).Unix()).
		WithCurrentWorkingDirectory("/home/testuser").
		WithCommand("#!/bin/bash\npython train_model.py").
		Build()
}

// initializeNodesAndPartitions initializes nodes and partitions (version-independent)
func (s *MockStorage) initializeNodesAndPartitions() {
	// Default nodes
	s.Nodes["node001"] = &MockNode{
		Name:         "node001",
		State:        []string{"IDLE"},
		CPUs:         16,
		Memory:       32 * 1024 * 1024 * 1024, // 32GB
		Features:     []string{"intel", "avx2"},
		Partition:    "compute",
		Architecture: "x86_64",
		OS:           "Linux",
	}

	s.Nodes["node002"] = &MockNode{
		Name:         "node002",
		State:        []string{"ALLOCATED"},
		CPUs:         16,
		Memory:       32 * 1024 * 1024 * 1024, // 32GB
		Features:     []string{"intel", "avx2"},
		Partition:    "compute",
		Architecture: "x86_64",
		OS:           "Linux",
	}

	s.Nodes["gpu001"] = &MockNode{
		Name:         "gpu001",
		State:        []string{"IDLE"},
		CPUs:         16,
		Memory:       64 * 1024 * 1024 * 1024, // 64GB
		Features:     []string{"gpu", "nvidia", "v100"},
		Partition:    "gpu",
		Architecture: "x86_64",
		OS:           "Linux",
	}

	// Default partitions
	computeNodesConfig := "node[001-002]"
	computeNodesTotal := int32(2)
	s.Partitions["compute"] = &MockPartition{
		Name:        "compute",
		State:       "UP",
		MaxNodes:    2,
		MinNodes:    1,
		DefaultTime: 60,
		MaxTime:     1440,
		Nodes: &struct {
			AllowedAllocation *string `json:"allowed_allocation,omitempty"`
			Configured        *string `json:"configured,omitempty"`
			Total             *int32  `json:"total,omitempty"`
		}{
			Configured: &computeNodesConfig,
			Total:      &computeNodesTotal,
		},
		Features: []string{"intel", "avx2"},
	}

	gpuNodesConfig := "gpu001"
	gpuNodesTotal := int32(1)
	s.Partitions["gpu"] = &MockPartition{
		Name:        "gpu",
		State:       "UP",
		MaxNodes:    1,
		MinNodes:    1,
		DefaultTime: 120,
		MaxTime:     2880,
		Nodes: &struct {
			AllowedAllocation *string `json:"allowed_allocation,omitempty"`
			Configured        *string `json:"configured,omitempty"`
			Total             *int32  `json:"total,omitempty"`
		}{
			Configured: &gpuNodesConfig,
			Total:      &gpuNodesTotal,
		},
		Features: []string{"gpu", "nvidia"},
	}
}

// Middleware functions

// sanitizeForLog sanitizes a string value for safe logging by removing control characters
// that could be used for log injection attacks
// lgtm[go/log-injection] This function sanitizes log values by removing control characters
func sanitizeForLog(value string) string {
	// Replace newlines, carriage returns, and tabs with spaces
	sanitized := strings.Map(func(r rune) rune {
		if r == '\n' || r == '\r' || r == '\t' {
			return ' '
		}
		return r
	}, value)
	return sanitized
}

func (m *MockSlurmServer) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Sanitize method and path to prevent log injection
		method := sanitizeForLog(r.Method)
		path := sanitizeForLog(r.URL.Path)
		// lgtm[go/log-injection] Values are sanitized via sanitizeForLog() which removes control characters
		log.Printf("Mock SLURM API: %s %s", method, path)
		next.ServeHTTP(w, r)
	})
}

func (m *MockSlurmServer) trailingSlashMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Strip trailing slash from path (but keep root path "/")
		if len(r.URL.Path) > 1 && r.URL.Path[len(r.URL.Path)-1] == '/' {
			r.URL.Path = r.URL.Path[:len(r.URL.Path)-1]
		}
		next.ServeHTTP(w, r)
	})
}

func (m *MockSlurmServer) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.config.EnableAuth {
			// Check for SLURM token header
			token := r.Header.Get("X-SLURM-USER-TOKEN")
			if token != m.config.AuthToken {
				http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (m *MockSlurmServer) delayMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.config.ResponseDelay > 0 {
			time.Sleep(m.config.ResponseDelay)
		}
		next.ServeHTTP(w, r)
	})
}

func (m *MockSlurmServer) errorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Strip trailing slash for consistent error matching
		path := r.URL.Path
		if len(path) > 1 && path[len(path)-1] == '/' {
			path = path[:len(path)-1]
		}

		endpoint := r.Method + " " + path
		m.mu.RLock()
		errorResponse, hasError := m.config.ErrorResponses[endpoint]
		m.mu.RUnlock()

		if hasError {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(errorResponse.StatusCode)
			_ = json.NewEncoder(w).Encode(errorResponse.Body) // Ignore error during HTTP response
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Helper functions

func (m *MockSlurmServer) writeJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (m *MockSlurmServer) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	//nolint:errchkjson // Ignore JSON encoding error in test mock - already writing error response
	// Use SLURM REST API error format with errors array
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"errors": []map[string]interface{}{
			{
				"error":        "operation_not_supported",
				"description":  message,
				"error_number": statusCode,
			},
		},
	})
}

func (m *MockSlurmServer) generateJobID() int32 {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.config.JobIDCounter++
	return int32(m.config.JobIDCounter)
}

func parseQueryParam(r *http.Request, param string) string {
	return strings.TrimSpace(r.URL.Query().Get(param))
}

func parseQueryParamInt(r *http.Request, param string, defaultValue int) int {
	value := parseQueryParam(r, param)
	if value == "" {
		return defaultValue
	}
	if intValue, err := strconv.Atoi(value); err == nil {
		return intValue
	}
	return defaultValue
}
