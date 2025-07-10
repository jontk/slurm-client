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
)

// MockSlurmServer represents a mock SLURM REST API server for testing
type MockSlurmServer struct {
	server  *httptest.Server
	router  *mux.Router
	storage *MockStorage
	config  *ServerConfig
	mu      sync.RWMutex
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
type MockStorage struct {
	Jobs       map[string]*MockJob
	Nodes      map[string]*MockNode
	Partitions map[string]*MockPartition
	ClusterInfo *MockClusterInfo
	mu         sync.RWMutex
}

// MockJob represents a mock job
type MockJob struct {
	JobID       string            `json:"job_id"`
	Name        string            `json:"name"`
	UserID      string            `json:"user_id"`
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
	State        string   `json:"state"`
	CPUs         int      `json:"cpus"`
	Memory       int64    `json:"memory"`
	Features     []string `json:"features"`
	Partition    string   `json:"partition"`
	Architecture string   `json:"architecture"`
	OS           string   `json:"operating_system"`
}

// MockPartition represents a mock partition
type MockPartition struct {
	Name        string   `json:"name"`
	State       string   `json:"state"`
	MaxNodes    int      `json:"max_nodes"`
	MinNodes    int      `json:"min_nodes"`
	DefaultTime int      `json:"default_time"`
	MaxTime     int      `json:"max_time"`
	Nodes       []string `json:"nodes"`
	Features    []string `json:"features"`
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
		Jobs:       make(map[string]*MockJob),
		Nodes:      make(map[string]*MockNode),
		Partitions: make(map[string]*MockPartition),
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
		storage: storage,
		config:  config,
	}

	mock.setupRouter()
	mock.server = httptest.NewServer(mock.router)

	return mock
}

// DefaultServerConfig returns a default server configuration
func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		APIVersion:          "v0.0.42",
		SlurmVersion:        "25.05",
		EnableAuth:          false,
		AuthToken:           "test-token",
		ErrorResponses:      make(map[string]ErrorResponse),
		ResponseDelay:       0,
		JobIDCounter:        1000,
		SupportedOperations: map[string]bool{
			"jobs.list":        true,
			"jobs.get":         true,
			"jobs.submit":      true,
			"jobs.cancel":      true,
			"jobs.update":      true,
			"jobs.steps":       true,
			"nodes.list":       true,
			"nodes.get":        true,
			"nodes.update":     true,
			"partitions.list":  true,
			"partitions.get":   true,
			"partitions.update": true,
			"info.get":         true,
			"info.ping":        true,
			"info.stats":       true,
			"info.version":     true,
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
func (m *MockSlurmServer) GetJob(jobID string) *MockJob {
	m.storage.mu.RLock()
	defer m.storage.mu.RUnlock()
	return m.storage.Jobs[jobID]
}

// AddJob adds a job to storage
func (m *MockSlurmServer) AddJob(job *MockJob) {
	m.storage.mu.Lock()
	defer m.storage.mu.Unlock()
	m.storage.Jobs[job.JobID] = job
}

// UpdateJobState updates a job's state
func (m *MockSlurmServer) UpdateJobState(jobID, state string) {
	m.storage.mu.Lock()
	defer m.storage.mu.Unlock()
	if job, exists := m.storage.Jobs[jobID]; exists {
		job.State = state
		now := time.Now().Unix()
		switch state {
		case "RUNNING":
			job.StartTime = &now
		case "COMPLETED", "FAILED", "CANCELLED":
			job.EndTime = &now
		}
	}
}

// setupRouter configures the HTTP router with all endpoints
func (m *MockSlurmServer) setupRouter() {
	m.router = mux.NewRouter()

	// Add middleware
	m.router.Use(m.loggingMiddleware)
	m.router.Use(m.authMiddleware)
	m.router.Use(m.delayMiddleware)
	m.router.Use(m.errorMiddleware)

	// API version prefix
	apiRouter := m.router.PathPrefix("/slurm/" + m.config.APIVersion).Subrouter()

	// Job endpoints
	apiRouter.HandleFunc("/jobs", m.handleJobsList).Methods("GET")
	apiRouter.HandleFunc("/job/{job_id}", m.handleJobsGet).Methods("GET")
	apiRouter.HandleFunc("/job/submit", m.handleJobsSubmit).Methods("POST")
	apiRouter.HandleFunc("/job/{job_id}", m.handleJobsCancel).Methods("DELETE")
	apiRouter.HandleFunc("/job/{job_id}", m.handleJobsUpdate).Methods("PATCH")
	apiRouter.HandleFunc("/job/{job_id}/steps", m.handleJobsSteps).Methods("GET")

	// Node endpoints
	apiRouter.HandleFunc("/nodes", m.handleNodesList).Methods("GET")
	apiRouter.HandleFunc("/node/{node_name}", m.handleNodesGet).Methods("GET")
	apiRouter.HandleFunc("/node/{node_name}", m.handleNodesUpdate).Methods("PATCH")

	// Partition endpoints
	apiRouter.HandleFunc("/partitions", m.handlePartitionsList).Methods("GET")
	apiRouter.HandleFunc("/partition/{partition_name}", m.handlePartitionsGet).Methods("GET")
	apiRouter.HandleFunc("/partition/{partition_name}", m.handlePartitionsUpdate).Methods("PATCH")

	// Info endpoints
	apiRouter.HandleFunc("/diag", m.handleInfoGet).Methods("GET")
	apiRouter.HandleFunc("/ping", m.handleInfoPing).Methods("GET")
	apiRouter.HandleFunc("/stats", m.handleInfoStats).Methods("GET")
	apiRouter.HandleFunc("/", m.handleInfoVersion).Methods("GET")
}

// initializeDefaultData sets up default test data
func (s *MockStorage) initializeDefaultData() {
	// Default jobs
	s.Jobs["1001"] = &MockJob{
		JobID:      "1001",
		Name:       "test-job-1",
		UserID:     "testuser",
		State:      "RUNNING",
		Partition:  "compute",
		CPUs:       4,
		Memory:     4 * 1024 * 1024 * 1024, // 4GB
		TimeLimit:  60,
		SubmitTime: time.Now().Add(-10 * time.Minute).Unix(),
		WorkingDir: "/tmp",
		Script:     "#!/bin/bash\necho 'Hello World'",
		Environment: map[string]string{
			"PATH": "/usr/local/bin:/usr/bin:/bin",
		},
	}

	s.Jobs["1002"] = &MockJob{
		JobID:      "1002",
		Name:       "test-job-2",
		UserID:     "testuser",
		State:      "PENDING",
		Partition:  "gpu",
		CPUs:       8,
		Memory:     8 * 1024 * 1024 * 1024, // 8GB
		TimeLimit:  120,
		SubmitTime: time.Now().Add(-5 * time.Minute).Unix(),
		WorkingDir: "/home/testuser",
		Script:     "#!/bin/bash\npython train_model.py",
		Environment: map[string]string{
			"CUDA_VISIBLE_DEVICES": "0",
		},
	}

	// Default nodes
	s.Nodes["node001"] = &MockNode{
		Name:         "node001",
		State:        "IDLE",
		CPUs:         16,
		Memory:       32 * 1024 * 1024 * 1024, // 32GB
		Features:     []string{"intel", "avx2"},
		Partition:    "compute",
		Architecture: "x86_64",
		OS:           "Linux",
	}

	s.Nodes["node002"] = &MockNode{
		Name:         "node002",
		State:        "ALLOCATED",
		CPUs:         16,
		Memory:       32 * 1024 * 1024 * 1024, // 32GB
		Features:     []string{"intel", "avx2"},
		Partition:    "compute",
		Architecture: "x86_64",
		OS:           "Linux",
	}

	s.Nodes["gpu001"] = &MockNode{
		Name:         "gpu001",
		State:        "IDLE",
		CPUs:         16,
		Memory:       64 * 1024 * 1024 * 1024, // 64GB
		Features:     []string{"gpu", "nvidia", "v100"},
		Partition:    "gpu",
		Architecture: "x86_64",
		OS:           "Linux",
	}

	// Default partitions
	s.Partitions["compute"] = &MockPartition{
		Name:        "compute",
		State:       "UP",
		MaxNodes:    2,
		MinNodes:    1,
		DefaultTime: 60,
		MaxTime:     1440,
		Nodes:       []string{"node001", "node002"},
		Features:    []string{"intel", "avx2"},
	}

	s.Partitions["gpu"] = &MockPartition{
		Name:        "gpu",
		State:       "UP",
		MaxNodes:    1,
		MinNodes:    1,
		DefaultTime: 120,
		MaxTime:     2880,
		Nodes:       []string{"gpu001"},
		Features:    []string{"gpu", "nvidia"},
	}
}

// Middleware functions

func (m *MockSlurmServer) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Mock SLURM API: %s %s", r.Method, r.URL.Path)
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
		endpoint := r.Method + " " + r.URL.Path
		m.mu.RLock()
		errorResponse, hasError := m.config.ErrorResponses[endpoint]
		m.mu.RUnlock()

		if hasError {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(errorResponse.StatusCode)
			json.NewEncoder(w).Encode(errorResponse.Body)
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
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func (m *MockSlurmServer) generateJobID() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.config.JobIDCounter++
	return strconv.FormatInt(m.config.JobIDCounter, 10)
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