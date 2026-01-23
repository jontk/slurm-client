// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package mocks

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"

	// Still used by submit/update/cancel handlers (not yet converted)
	v0_0_40 "github.com/jontk/slurm-client/internal/api/v0_0_40"

	// Import generated builders for job creation
	builderv0040 "github.com/jontk/slurm-client/tests/mocks/generated/v0_0_40"
	builderv0042 "github.com/jontk/slurm-client/tests/mocks/generated/v0_0_42"
	builderv0043 "github.com/jontk/slurm-client/tests/mocks/generated/v0_0_43"
)

// Job endpoint handlers

func (m *MockSlurmServer) handleJobsList(w http.ResponseWriter, r *http.Request) {
	if !m.config.SupportedOperations["jobs.list"] {
		m.writeErrorResponse(w, http.StatusNotImplemented, "Operation not supported")
		return
	}

	// Parse pagination parameters
	limit := parseQueryParamInt(r, "limit", 100)
	offset := parseQueryParamInt(r, "offset", 0)

	// Use version-specific handler to get filtered jobs
	jobs, err := m.jobHandler.ListJobs(r, m.storage)
	if err != nil {
		m.writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Apply pagination
	total := len(jobs)
	start := offset
	end := offset + limit
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}
	if start < end {
		jobs = jobs[start:end]
	} else {
		jobs = []interface{}{}
	}

	// Return version-specific types directly - JSON encoder handles serialization
	response := map[string]interface{}{
		"jobs": jobs,
		"meta": map[string]interface{}{
			"total":  total,
			"offset": offset,
			"limit":  limit,
		},
	}

	m.writeJSONResponse(w, response)
}

func (m *MockSlurmServer) handleJobsGet(w http.ResponseWriter, r *http.Request) {
	if !m.config.SupportedOperations["jobs.get"] {
		m.writeErrorResponse(w, http.StatusNotImplemented, "Operation not supported")
		return
	}

	vars := mux.Vars(r)
	jobID := vars["job_id"]

	// Use version-specific handler to get job
	job, err := m.jobHandler.GetJob(jobID, m.storage)
	if err != nil {
		m.writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if job == nil {
		m.writeErrorResponse(w, http.StatusNotFound, fmt.Sprintf("Job %s not found", jobID))
		return
	}

	// Return version-specific type directly - API expects "jobs" array
	response := map[string]interface{}{
		"jobs": []interface{}{job},
	}

	m.writeJSONResponse(w, response)
}

func (m *MockSlurmServer) handleJobsSubmit(w http.ResponseWriter, r *http.Request) {
	if !m.config.SupportedOperations["jobs.submit"] {
		m.writeErrorResponse(w, http.StatusNotImplemented, "Operation not supported")
		return
	}

	// SLURM API sends job data nested under "job" key
	var requestBody struct {
		Job struct {
			Name      string `json:"name"`
			Script    string `json:"script"`
			Partition string `json:"partition"`
			CPUs      int    `json:"cpus"`
			Memory    int64  `json:"memory"`
			TimeLimit *struct {
				Number *int32 `json:"number,omitempty"`
			} `json:"time_limit,omitempty"`
			WorkingDir  string   `json:"working_directory"`
			Environment []string `json:"environment"` // Array of "KEY=value" strings
		} `json:"job"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		log.Printf("Job submit decode error: %v", err)
		m.writeErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Invalid JSON: %v", err))
		return
	}

	submission := requestBody.Job

	// Validate submission
	if submission.Name == "" {
		log.Printf("Job submit validation failed: name is empty")
		m.writeErrorResponse(w, http.StatusBadRequest, "Job name is required")
		return
	}
	if submission.Script == "" {
		log.Printf("Job submit validation failed: script is empty")
		m.writeErrorResponse(w, http.StatusBadRequest, "Job script is required")
		return
	}

	// Convert environment array to map
	envMap := make(map[string]string)
	for _, env := range submission.Environment {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			envMap[parts[0]] = parts[1]
		}
	}

	// Create new job
	jobID := m.generateJobID()

	// Extract time limit from V0040Uint32NoVal structure
	timeLimit := 60 // Default
	if submission.TimeLimit != nil && submission.TimeLimit.Number != nil {
		timeLimit = int(*submission.TimeLimit.Number)
	}

	// Set defaults
	partition := submission.Partition
	if partition == "" {
		partition = "compute"
	}
	cpus := submission.CPUs
	if cpus == 0 {
		cpus = 1
	}
	workingDir := submission.WorkingDir
	if workingDir == "" {
		workingDir = "/tmp"
	}

	// Create version-specific job using builders
	var job interface{}
	switch m.config.APIVersion {
	case "v0.0.40", "v0.0.41":
		// v0.0.41 falls back to v0.0.40 builders
		job = builderv0040.NewJobInfo().
			WithJobId(jobID).
			WithName(submission.Name).
			WithUserId(1000).
			WithJobState("PENDING").
			WithPartition(partition).
			WithCpus(int64(cpus)).
			WithMemoryPerNode(submission.Memory).
			WithTimeLimit(int64(timeLimit)).
			WithSubmitTime(time.Now().Unix()).
			WithCurrentWorkingDirectory(workingDir).
			WithCommand(submission.Script).
			Build()
	case "v0.0.42":
		job = builderv0042.NewJobInfo().
			WithJobId(jobID).
			WithName(submission.Name).
			WithUserId(1000).
			WithJobState("PENDING").
			WithPartition(partition).
			WithCpus(int32(cpus)).
			WithMemoryPerNode(submission.Memory).
			WithTimeLimit(int32(timeLimit)).
			WithSubmitTime(time.Now().Unix()).
			WithCurrentWorkingDirectory(workingDir).
			WithCommand(submission.Script).
			Build()
	case "v0.0.43", "v0.0.44":
		job = builderv0043.NewJobInfo().
			WithJobId(jobID).
			WithName(submission.Name).
			WithUserId(1000).
			WithJobState("PENDING").
			WithPartition(partition).
			WithCpus(int32(cpus)).
			WithMemoryPerNode(submission.Memory).
			WithTimeLimit(int32(timeLimit)).
			WithSubmitTime(time.Now().Unix()).
			WithCurrentWorkingDirectory(workingDir).
			WithCommand(submission.Script).
			Build()
	default:
		// Fallback to v0.0.40 for unknown versions
		job = builderv0040.NewJobInfo().
			WithJobId(jobID).
			WithName(submission.Name).
			WithUserId(1000).
			WithJobState("PENDING").
			WithPartition(partition).
			WithCpus(int64(cpus)).
			WithMemoryPerNode(submission.Memory).
			WithTimeLimit(int64(timeLimit)).
			WithSubmitTime(time.Now().Unix()).
			WithCurrentWorkingDirectory(workingDir).
			WithCommand(submission.Script).
			Build()
	}

	m.storage.mu.Lock()
	m.storage.Jobs[strconv.Itoa(int(jobID))] = job
	m.storage.mu.Unlock()

	// SLURM API response format matches OpenAPI spec
	response := map[string]interface{}{
		"job_id":              jobID,
		"step_id":             "batch",
		"job_submit_user_msg": "job submitted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Use 200 OK for consistency
	//nolint:errchkjson // Ignore JSON encoding error in test mock - response already committed
	_ = json.NewEncoder(w).Encode(response)
}

func (m *MockSlurmServer) handleJobsCancel(w http.ResponseWriter, r *http.Request) {
	if !m.config.SupportedOperations["jobs.cancel"] {
		m.writeErrorResponse(w, http.StatusNotImplemented, "Operation not supported")
		return
	}

	vars := mux.Vars(r)
	jobID := vars["job_id"]

	m.storage.mu.Lock()
	jobInterface, exists := m.storage.Jobs[jobID]
	if exists {
		// Type assert to v0.0.40 job type
		if job, ok := jobInterface.(*v0_0_40.V0040JobInfo); ok {
			// Check current state (JobState is an array in the API)
			currentState := ""
			if job.JobState != nil && len(*job.JobState) > 0 {
				currentState = (*job.JobState)[0]
			}

			if currentState == "RUNNING" || currentState == "PENDING" {
				// Update to CANCELLED state
				cancelledState := []string{"CANCELLED"}
				job.JobState = &cancelledState
				// Note: EndTime is a NoVal field in the real API
				// For now, skipping this update as it requires proper NoVal handling
			} else {
				m.storage.mu.Unlock()
				m.writeErrorResponse(w, http.StatusConflict, fmt.Sprintf("Job %s cannot be cancelled (state: %s)", jobID, currentState))
				return
			}
		}
	}
	m.storage.mu.Unlock()

	if !exists {
		m.writeErrorResponse(w, http.StatusNotFound, fmt.Sprintf("Job %s not found", jobID))
		return
	}

	response := map[string]interface{}{
		"result": "SUCCESS",
	}

	m.writeJSONResponse(w, response)
}

func (m *MockSlurmServer) handleJobsUpdate(w http.ResponseWriter, r *http.Request) {
	if !m.config.SupportedOperations["jobs.update"] {
		m.writeErrorResponse(w, http.StatusNotImplemented, "Operation not supported")
		return
	}

	vars := mux.Vars(r)
	jobID := vars["job_id"]

	var update struct {
		Name      *string `json:"name,omitempty"`
		Priority  *int    `json:"priority,omitempty"`
		TimeLimit *int    `json:"time_limit,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		m.writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	m.storage.mu.Lock()
	jobInterface, exists := m.storage.Jobs[jobID]
	if exists {
		// Type assert to v0.0.40 job type
		if job, ok := jobInterface.(*v0_0_40.V0040JobInfo); ok {
			if update.Name != nil {
				job.Name = update.Name
			}
			if update.TimeLimit != nil {
				// TimeLimit is a NoVal field, needs proper wrapping
				setTrue := true
				timeLimit := int64(*update.TimeLimit)
				job.TimeLimit = &v0_0_40.V0040Uint32NoVal{
					Set:    &setTrue,
					Number: &timeLimit,
				}
			}
			// Priority field is available in V0040JobInfo as Priority (NoVal type)
		}
	}
	m.storage.mu.Unlock()

	if !exists {
		m.writeErrorResponse(w, http.StatusNotFound, fmt.Sprintf("Job %s not found", jobID))
		return
	}

	response := map[string]interface{}{
		"result": "SUCCESS",
	}

	m.writeJSONResponse(w, response)
}

func (m *MockSlurmServer) handleJobsSteps(w http.ResponseWriter, r *http.Request) {
	if !m.config.SupportedOperations["jobs.steps"] {
		m.writeErrorResponse(w, http.StatusNotImplemented, "Operation not supported")
		return
	}

	vars := mux.Vars(r)
	jobID := vars["job_id"]

	m.storage.mu.RLock()
	_, exists := m.storage.Jobs[jobID]
	m.storage.mu.RUnlock()

	if !exists {
		m.writeErrorResponse(w, http.StatusNotFound, fmt.Sprintf("Job %s not found", jobID))
		return
	}

	// Mock job steps response
	steps := []map[string]interface{}{
		{
			"step_id":    "0",
			"name":       "batch",
			"state":      "RUNNING",
			"start_time": time.Now().Add(-5 * time.Minute).Unix(),
		},
	}

	response := map[string]interface{}{
		"job_steps": steps,
	}

	m.writeJSONResponse(w, response)
}

// Node endpoint handlers

func (m *MockSlurmServer) handleNodesList(w http.ResponseWriter, r *http.Request) {
	if !m.config.SupportedOperations["nodes.list"] {
		m.writeErrorResponse(w, http.StatusNotImplemented, "Operation not supported")
		return
	}

	// Parse query parameters
	partition := parseQueryParam(r, "partition")
	states := parseQueryParam(r, "state")
	features := parseQueryParam(r, "features")
	limit := parseQueryParamInt(r, "limit", 100)
	offset := parseQueryParamInt(r, "offset", 0)

	stateList := []string{}
	if states != "" {
		stateList = strings.Split(states, ",")
	}

	featureList := []string{}
	if features != "" {
		featureList = strings.Split(features, ",")
	}

	m.storage.mu.RLock()
	defer m.storage.mu.RUnlock()

	nodes := []*MockNode{}
	for _, node := range m.storage.Nodes {
		// Apply filters
		if partition != "" && node.Partition != partition {
			continue
		}
		if len(stateList) > 0 {
			found := false
			for _, state := range stateList {
				requestedState := strings.TrimSpace(state)
				// Check if any of the node's states match the requested state
				for _, nodeState := range node.State {
					if nodeState == requestedState {
						found = true
						break
					}
				}
				if found {
					break
				}
			}
			if !found {
				continue
			}
		}
		if len(featureList) > 0 {
			hasAllFeatures := true
			for _, requiredFeature := range featureList {
				requiredFeature = strings.TrimSpace(requiredFeature)
				found := false
				for _, nodeFeature := range node.Features {
					if nodeFeature == requiredFeature {
						found = true
						break
					}
				}
				if !found {
					hasAllFeatures = false
					break
				}
			}
			if !hasAllFeatures {
				continue
			}
		}
		nodes = append(nodes, node)
	}

	// Apply pagination
	total := len(nodes)
	start := offset
	end := offset + limit
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}
	if start < end {
		nodes = nodes[start:end]
	} else {
		nodes = []*MockNode{}
	}

	response := map[string]interface{}{
		"nodes": nodes,
		"meta": map[string]interface{}{
			"total":  total,
			"offset": offset,
			"limit":  limit,
		},
	}

	m.writeJSONResponse(w, response)
}

func (m *MockSlurmServer) handleNodesGet(w http.ResponseWriter, r *http.Request) {
	if !m.config.SupportedOperations["nodes.get"] {
		m.writeErrorResponse(w, http.StatusNotImplemented, "Operation not supported")
		return
	}

	vars := mux.Vars(r)
	nodeName := vars["node_name"]

	m.storage.mu.RLock()
	node, exists := m.storage.Nodes[nodeName]
	m.storage.mu.RUnlock()

	if !exists {
		m.writeErrorResponse(w, http.StatusNotFound, fmt.Sprintf("Node %s not found", nodeName))
		return
	}

	// OpenAPI spec expects nodes array even for single node get
	lastUpdate := time.Now().Unix()
	response := map[string]interface{}{
		"nodes": []*MockNode{node},
		"last_update": map[string]interface{}{
			"number": lastUpdate,
		},
	}

	m.writeJSONResponse(w, response)
}

func (m *MockSlurmServer) handleNodesUpdate(w http.ResponseWriter, r *http.Request) {
	if !m.config.SupportedOperations["nodes.update"] {
		m.writeErrorResponse(w, http.StatusNotImplemented, "Operation not supported")
		return
	}

	vars := mux.Vars(r)
	nodeName := vars["node_name"]

	var update struct {
		State *[]string `json:"state,omitempty"` // State is an array of strings in the OpenAPI spec
	}

	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		m.writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	m.storage.mu.Lock()
	node, exists := m.storage.Nodes[nodeName]
	if exists && update.State != nil {
		// Update node state with the provided array
		node.State = *update.State
	}
	m.storage.mu.Unlock()

	if !exists {
		m.writeErrorResponse(w, http.StatusNotFound, fmt.Sprintf("Node %s not found", nodeName))
		return
	}

	response := map[string]interface{}{
		"result": "SUCCESS",
	}

	m.writeJSONResponse(w, response)
}

// Partition endpoint handlers

func (m *MockSlurmServer) handlePartitionsList(w http.ResponseWriter, r *http.Request) {
	if !m.config.SupportedOperations["partitions.list"] {
		m.writeErrorResponse(w, http.StatusNotImplemented, "Operation not supported")
		return
	}

	limit := parseQueryParamInt(r, "limit", 100)
	offset := parseQueryParamInt(r, "offset", 0)

	m.storage.mu.RLock()
	defer m.storage.mu.RUnlock()

	partitions := []*MockPartition{}
	for _, partition := range m.storage.Partitions {
		partitions = append(partitions, partition)
	}

	// Apply pagination
	total := len(partitions)
	start := offset
	end := offset + limit
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}
	if start < end {
		partitions = partitions[start:end]
	} else {
		partitions = []*MockPartition{}
	}

	response := map[string]interface{}{
		"partitions": partitions,
		"meta": map[string]interface{}{
			"total":  total,
			"offset": offset,
			"limit":  limit,
		},
	}

	m.writeJSONResponse(w, response)
}

func (m *MockSlurmServer) handlePartitionsGet(w http.ResponseWriter, r *http.Request) {
	if !m.config.SupportedOperations["partitions.get"] {
		m.writeErrorResponse(w, http.StatusNotImplemented, "Operation not supported")
		return
	}

	vars := mux.Vars(r)
	partitionName := vars["partition_name"]

	m.storage.mu.RLock()
	partition, exists := m.storage.Partitions[partitionName]
	m.storage.mu.RUnlock()

	if !exists {
		m.writeErrorResponse(w, http.StatusNotFound, fmt.Sprintf("Partition %s not found", partitionName))
		return
	}

	// OpenAPI spec expects partitions array even for single partition get
	lastUpdate := time.Now().Unix()
	response := map[string]interface{}{
		"partitions": []*MockPartition{partition},
		"last_update": map[string]interface{}{
			"number": lastUpdate,
		},
	}

	m.writeJSONResponse(w, response)
}

func (m *MockSlurmServer) handlePartitionsUpdate(w http.ResponseWriter, r *http.Request) {
	if !m.config.SupportedOperations["partitions.update"] {
		m.writeErrorResponse(w, http.StatusNotImplemented, "Operation not supported")
		return
	}

	vars := mux.Vars(r)
	partitionName := vars["partition_name"]

	var update struct {
		State       *string `json:"state,omitempty"`
		DefaultTime *int    `json:"default_time,omitempty"`
		MaxTime     *int    `json:"max_time,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		m.writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	m.storage.mu.Lock()
	partition, exists := m.storage.Partitions[partitionName]
	if exists {
		if update.State != nil {
			partition.State = *update.State
		}
		if update.DefaultTime != nil {
			partition.DefaultTime = *update.DefaultTime
		}
		if update.MaxTime != nil {
			partition.MaxTime = *update.MaxTime
		}
	}
	m.storage.mu.Unlock()

	if !exists {
		m.writeErrorResponse(w, http.StatusNotFound, fmt.Sprintf("Partition %s not found", partitionName))
		return
	}

	response := map[string]interface{}{
		"result": "SUCCESS",
	}

	m.writeJSONResponse(w, response)
}

// Info endpoint handlers

func (m *MockSlurmServer) handleInfoGet(w http.ResponseWriter, r *http.Request) {
	if !m.config.SupportedOperations["info.get"] {
		m.writeErrorResponse(w, http.StatusNotImplemented, "Operation not supported")
		return
	}

	m.storage.mu.RLock()
	info := m.storage.ClusterInfo
	m.storage.mu.RUnlock()

	response := map[string]interface{}{
		"cluster": info,
	}

	m.writeJSONResponse(w, response)
}

func (m *MockSlurmServer) handleInfoPing(w http.ResponseWriter, r *http.Request) {
	if !m.config.SupportedOperations["info.ping"] {
		m.writeErrorResponse(w, http.StatusNotImplemented, "Operation not supported")
		return
	}

	// Parse Slurm version from config (e.g., "24.05" -> major: "24", minor: "05")
	versionParts := strings.Split(m.config.SlurmVersion, ".")
	major := "0"
	minor := "0"
	micro := "0"
	if len(versionParts) >= 1 {
		major = versionParts[0]
	}
	if len(versionParts) >= 2 {
		minor = versionParts[1]
	}
	if len(versionParts) >= 3 {
		micro = versionParts[2]
	}

	// Return proper OpenAPI ping response with meta information
	response := map[string]interface{}{
		"meta": map[string]interface{}{
			"slurm": map[string]interface{}{
				"cluster": "test-cluster",
				"release": m.config.SlurmVersion,
				"version": map[string]interface{}{
					"major": major,
					"minor": minor,
					"micro": micro,
				},
			},
		},
		"pings": []interface{}{
			map[string]interface{}{
				"hostname": "controller",
				"ping":     "UP",
				"latency":  int64(1), // 1ms
			},
		},
	}

	m.writeJSONResponse(w, response)
}

func (m *MockSlurmServer) handleInfoStats(w http.ResponseWriter, r *http.Request) {
	if !m.config.SupportedOperations["info.stats"] {
		m.writeErrorResponse(w, http.StatusNotImplemented, "Operation not supported")
		return
	}

	m.storage.mu.RLock()
	runningJobs := 0
	pendingJobs := 0
	for _, jobInterface := range m.storage.Jobs {
		// Use version-agnostic helper to get job state
		state := extractJobState(jobInterface)
		switch state {
		case "RUNNING":
			runningJobs++
		case "PENDING":
			pendingJobs++
		}
	}

	totalNodes := len(m.storage.Nodes)
	totalCPUs := 0
	for _, node := range m.storage.Nodes {
		totalCPUs += node.CPUs
	}
	m.storage.mu.RUnlock()

	stats := map[string]interface{}{
		"running_jobs": runningJobs,
		"pending_jobs": pendingJobs,
		"total_nodes":  totalNodes,
		"total_cpus":   totalCPUs,
		"timestamp":    time.Now().Unix(),
	}

	response := map[string]interface{}{
		"statistics": stats,
	}

	m.writeJSONResponse(w, response)
}

func (m *MockSlurmServer) handleInfoVersion(w http.ResponseWriter, r *http.Request) {
	if !m.config.SupportedOperations["info.version"] {
		m.writeErrorResponse(w, http.StatusNotImplemented, "Operation not supported")
		return
	}

	version := map[string]interface{}{
		"api_version":   m.config.APIVersion,
		"slurm_version": m.config.SlurmVersion,
		"server_time":   time.Now().Unix(),
	}

	response := map[string]interface{}{
		"version": version,
	}

	m.writeJSONResponse(w, response)
}

// Job Analytics endpoint handlers

func (m *MockSlurmServer) handleJobUtilization(w http.ResponseWriter, r *http.Request) {
	if !m.config.SupportedOperations["jobs.utilization"] {
		m.writeErrorResponse(w, http.StatusNotImplemented, "Operation not supported")
		return
	}

	vars := mux.Vars(r)
	jobID := vars["job_id"]

	m.storage.mu.RLock()
	jobInterface, exists := m.storage.Jobs[jobID]
	m.storage.mu.RUnlock()

	if !exists {
		m.writeErrorResponse(w, http.StatusNotFound, fmt.Sprintf("Job %s not found", jobID))
		return
	}

	// Extract values using version-agnostic helpers
	cpus := extractCPUs(jobInterface)
	memory := extractMemory(jobInterface)

	// Generate mock utilization data based on job properties
	utilization := map[string]interface{}{
		"job_id": jobID,
		"cpu_utilization": map[string]interface{}{
			"allocated_cores":     cpus,
			"used_cores":          float64(cpus) * 0.85, // 85% utilization
			"utilization_percent": 85.0,
			"efficiency_percent":  82.0,
		},
		"memory_utilization": map[string]interface{}{
			"allocated_bytes":     memory,
			"used_bytes":          int64(float64(memory) * 0.78), // 78% utilization
			"utilization_percent": 78.0,
			"efficiency_percent":  76.0,
		},
		"gpu_utilization": map[string]interface{}{
			"device_count":        0, // No GPUs for basic jobs
			"utilization_percent": 0.0,
		},
		"io_utilization": map[string]interface{}{
			"read_bytes":          int64(1024 * 1024 * 100), // 100MB read
			"write_bytes":         int64(1024 * 1024 * 50),  // 50MB write
			"utilization_percent": 65.0,
		},
		"network_utilization": map[string]interface{}{
			"total_bandwidth": map[string]interface{}{
				"used_max":   float64(1000), // 1Gbps
				"efficiency": 45.0,
			},
		},
	}

	response := map[string]interface{}{
		"utilization": utilization,
	}

	m.writeJSONResponse(w, response)
}

func (m *MockSlurmServer) handleJobEfficiency(w http.ResponseWriter, r *http.Request) {
	if !m.config.SupportedOperations["jobs.efficiency"] {
		m.writeErrorResponse(w, http.StatusNotImplemented, "Operation not supported")
		return
	}

	vars := mux.Vars(r)
	jobID := vars["job_id"]

	m.storage.mu.RLock()
	_, exists := m.storage.Jobs[jobID]
	m.storage.mu.RUnlock()

	if !exists {
		m.writeErrorResponse(w, http.StatusNotFound, fmt.Sprintf("Job %s not found", jobID))
		return
	}

	efficiency := map[string]interface{}{
		"job_id":                   jobID,
		"overall_efficiency_score": 79.5,
		"cpu_efficiency":           82.0,
		"memory_efficiency":        76.0,
		"gpu_efficiency":           0.0,
		"io_efficiency":            65.0,
		"network_efficiency":       45.0,
		"energy_efficiency":        88.0,
		"resource_waste": map[string]interface{}{
			"cpu_core_hours":  0.6,  // 0.6 core-hours wasted
			"memory_gb_hours": 1.2,  // 1.2 GB-hours wasted
			"cpu_percent":     15.0, // 15% CPU waste
			"memory_percent":  22.0, // 22% memory waste
		},
		"optimization_recommendations": []map[string]interface{}{
			{
				"resource":    "Memory",
				"type":        "reduction",
				"current":     4,
				"recommended": 3,
				"reason":      "Low memory utilization detected",
				"confidence":  0.85,
			},
		},
	}

	response := map[string]interface{}{
		"efficiency": efficiency,
	}

	m.writeJSONResponse(w, response)
}

func (m *MockSlurmServer) handleJobPerformance(w http.ResponseWriter, r *http.Request) {
	if !m.config.SupportedOperations["jobs.performance"] {
		m.writeErrorResponse(w, http.StatusNotImplemented, "Operation not supported")
		return
	}

	vars := mux.Vars(r)
	jobID := vars["job_id"]

	m.storage.mu.RLock()
	_, exists := m.storage.Jobs[jobID]
	m.storage.mu.RUnlock()

	if !exists {
		m.writeErrorResponse(w, http.StatusNotFound, fmt.Sprintf("Job %s not found", jobID))
		return
	}

	performance := map[string]interface{}{
		"job_id": jobID,
		"cpu_analytics": map[string]interface{}{
			"allocated_cores":     4,
			"used_cores":          3.4,
			"utilization_percent": 85.0,
			"efficiency_percent":  82.0,
			"average_frequency":   2800.0,
			"max_frequency":       3200.0,
		},
		"memory_analytics": map[string]interface{}{
			"allocated_bytes":     int64(4 * 1024 * 1024 * 1024), // 4GB
			"used_bytes":          int64(3 * 1024 * 1024 * 1024), // 3GB
			"utilization_percent": 75.0,
			"efficiency_percent":  73.0,
		},
		"io_analytics": map[string]interface{}{
			"read_bytes":              int64(100 * 1024 * 1024), // 100MB
			"write_bytes":             int64(50 * 1024 * 1024),  // 50MB
			"read_operations":         int64(1000),
			"write_operations":        int64(500),
			"average_read_bandwidth":  120.5,
			"average_write_bandwidth": 85.2,
		},
		"overall_efficiency": 79.5,
	}

	response := map[string]interface{}{
		"performance": performance,
	}

	m.writeJSONResponse(w, response)
}

func (m *MockSlurmServer) handleJobLiveMetrics(w http.ResponseWriter, r *http.Request) {
	if !m.config.SupportedOperations["jobs.live_metrics"] {
		m.writeErrorResponse(w, http.StatusNotImplemented, "Operation not supported")
		return
	}

	vars := mux.Vars(r)
	jobID := vars["job_id"]

	m.storage.mu.RLock()
	jobInterface, exists := m.storage.Jobs[jobID]
	m.storage.mu.RUnlock()

	if !exists {
		m.writeErrorResponse(w, http.StatusNotFound, fmt.Sprintf("Job %s not found", jobID))
		return
	}

	// Check if job is running using version-agnostic helper
	currentState := extractJobState(jobInterface)
	if currentState != "RUNNING" {
		m.writeErrorResponse(w, http.StatusConflict, fmt.Sprintf("Job %s is not running (state: %s)", jobID, currentState))
		return
	}

	liveMetrics := map[string]interface{}{
		"job_id":    jobID,
		"timestamp": time.Now().Unix(),
		"cpu_usage": map[string]interface{}{
			"current":     85.2,
			"average":     82.5,
			"peak":        95.1,
			"utilization": 85.0,
		},
		"memory_usage": map[string]interface{}{
			"current":     int64(3221225472), // ~3GB
			"average":     int64(3000000000), // ~2.8GB
			"peak":        int64(3400000000), // ~3.2GB
			"utilization": 75.0,
		},
		"network_usage": map[string]interface{}{
			"in_rate_mbps":  125.5,
			"out_rate_mbps": 87.3,
		},
		"disk_usage": map[string]interface{}{
			"read_rate_mbps":  45.2,
			"write_rate_mbps": 28.7,
		},
	}

	response := map[string]interface{}{
		"live_metrics": liveMetrics,
	}

	m.writeJSONResponse(w, response)
}

func (m *MockSlurmServer) handleJobResourceTrends(w http.ResponseWriter, r *http.Request) {
	if !m.config.SupportedOperations["jobs.resource_trends"] {
		m.writeErrorResponse(w, http.StatusNotImplemented, "Operation not supported")
		return
	}

	vars := mux.Vars(r)
	jobID := vars["job_id"]

	m.storage.mu.RLock()
	_, exists := m.storage.Jobs[jobID]
	m.storage.mu.RUnlock()

	if !exists {
		m.writeErrorResponse(w, http.StatusNotFound, fmt.Sprintf("Job %s not found", jobID))
		return
	}

	// Generate mock trend data for the last hour
	trends := []map[string]interface{}{}
	baseTime := time.Now().Add(-1 * time.Hour)
	for i := 0; i < 12; i++ { // 12 data points, 5 minutes apart
		timestamp := baseTime.Add(time.Duration(i*5) * time.Minute)
		trends = append(trends, map[string]interface{}{
			"timestamp":          timestamp.Unix(),
			"cpu_utilization":    75.0 + float64(i)*2.0,  // Gradual increase
			"memory_utilization": 70.0 + float64(i)*1.5,  // Gradual increase
			"io_bandwidth":       120.0 - float64(i)*0.5, // Gradual decrease
		})
	}

	response := map[string]interface{}{
		"job_id": jobID,
		"trends": trends,
		"analysis": map[string]interface{}{
			"cpu_trend":    "increasing",
			"memory_trend": "increasing",
			"io_trend":     "decreasing",
		},
	}

	m.writeJSONResponse(w, response)
}

func (m *MockSlurmServer) handleJobStepUtilization(w http.ResponseWriter, r *http.Request) {
	if !m.config.SupportedOperations["jobs.step_utilization"] {
		m.writeErrorResponse(w, http.StatusNotImplemented, "Operation not supported")
		return
	}

	vars := mux.Vars(r)
	jobID := vars["job_id"]
	stepID := vars["step_id"]

	m.storage.mu.RLock()
	_, exists := m.storage.Jobs[jobID]
	m.storage.mu.RUnlock()

	if !exists {
		m.writeErrorResponse(w, http.StatusNotFound, fmt.Sprintf("Job %s not found", jobID))
		return
	}

	stepUtilization := map[string]interface{}{
		"job_id":  jobID,
		"step_id": stepID,
		"cpu_utilization": map[string]interface{}{
			"utilization_percent": 88.5,
			"efficiency_percent":  85.2,
		},
		"memory_utilization": map[string]interface{}{
			"utilization_percent": 82.1,
			"efficiency_percent":  79.8,
		},
	}

	response := map[string]interface{}{
		"step_utilization": stepUtilization,
	}

	m.writeJSONResponse(w, response)
}

func (m *MockSlurmServer) handleJobPerformanceHistory(w http.ResponseWriter, r *http.Request) {
	if !m.config.SupportedOperations["jobs.performance_history"] {
		m.writeErrorResponse(w, http.StatusNotImplemented, "Operation not supported")
		return
	}

	jobID := parseQueryParam(r, "job_id")
	if jobID == "" {
		m.writeErrorResponse(w, http.StatusBadRequest, "job_id parameter is required")
		return
	}

	// Generate mock historical data
	baseTime := time.Now().Add(-24 * time.Hour)
	timeSeriesData := []map[string]interface{}{}
	for i := 0; i < 24; i++ { // Hourly data for 24 hours
		timestamp := baseTime.Add(time.Duration(i) * time.Hour)
		timeSeriesData = append(timeSeriesData, map[string]interface{}{
			"timestamp":          timestamp.Unix(),
			"cpu_utilization":    75.0 + float64(i%6)*5.0,   // Varies between 75-100%
			"memory_utilization": 65.0 + float64(i%4)*7.5,   // Varies between 65-87.5%
			"io_bandwidth":       100.0 + float64(i%8)*12.5, // Varies between 100-187.5MB/s
			"efficiency":         70.0 + float64(i%5)*6.0,   // Varies between 70-94%
		})
	}

	history := map[string]interface{}{
		"job_id":           jobID,
		"start_time":       baseTime.Unix(),
		"end_time":         time.Now().Unix(),
		"time_series_data": timeSeriesData,
		"statistics": map[string]interface{}{
			"average_cpu":        85.5,
			"average_memory":     76.25,
			"average_io":         143.75,
			"average_efficiency": 82.0,
			"peak_cpu":           100.0,
			"peak_memory":        87.5,
			"peak_io":            187.5,
		},
		"trends": map[string]interface{}{
			"cpu_trend": map[string]interface{}{
				"direction":  "stable",
				"slope":      0.15,
				"confidence": 0.75,
			},
			"memory_trend": map[string]interface{}{
				"direction":  "increasing",
				"slope":      0.85,
				"confidence": 0.82,
			},
		},
		"anomalies": []map[string]interface{}{},
	}

	response := map[string]interface{}{
		"performance_history": history,
	}

	m.writeJSONResponse(w, response)
}

func (m *MockSlurmServer) handlePerformanceTrends(w http.ResponseWriter, r *http.Request) {
	if !m.config.SupportedOperations["jobs.performance_trends"] {
		m.writeErrorResponse(w, http.StatusNotImplemented, "Operation not supported")
		return
	}

	trends := map[string]interface{}{
		"cluster_performance": map[string]interface{}{
			"average_efficiency":  78.5,
			"total_jobs_analyzed": 1247,
			"efficiency_trend":    "improving",
			"trend_period_days":   30,
		},
		"resource_trends": map[string]interface{}{
			"cpu_efficiency_avg":     82.1,
			"memory_efficiency_avg":  75.6,
			"io_efficiency_avg":      68.9,
			"network_efficiency_avg": 45.2,
		},
		"partition_trends": []map[string]interface{}{
			{
				"partition":          "compute",
				"average_efficiency": 79.8,
				"jobs_count":         856,
			},
			{
				"partition":          "gpu",
				"average_efficiency": 75.2,
				"jobs_count":         391,
			},
		},
	}

	response := map[string]interface{}{
		"performance_trends": trends,
	}

	m.writeJSONResponse(w, response)
}

func (m *MockSlurmServer) handleUserEfficiencyTrends(w http.ResponseWriter, r *http.Request) {
	if !m.config.SupportedOperations["jobs.efficiency_trends"] {
		m.writeErrorResponse(w, http.StatusNotImplemented, "Operation not supported")
		return
	}

	userID := parseQueryParam(r, "user_id")
	if userID == "" {
		userID = "testuser" // Default test user
	}

	trends := map[string]interface{}{
		"user_id": userID,
		"efficiency_trends": map[string]interface{}{
			"current_avg_efficiency":  81.2,
			"previous_avg_efficiency": 78.9,
			"improvement_percent":     2.9,
			"trend_direction":         "improving",
		},
		"monthly_data": []map[string]interface{}{
			{"month": "2024-01", "avg_efficiency": 76.5, "jobs_count": 45},
			{"month": "2024-02", "avg_efficiency": 78.9, "jobs_count": 52},
			{"month": "2024-03", "avg_efficiency": 81.2, "jobs_count": 48},
		},
		"recommendations": []map[string]interface{}{
			{
				"category": "Memory",
				"message":  "Consider reducing memory allocation by 15-20% based on usage patterns",
				"impact":   "High",
			},
		},
	}

	response := map[string]interface{}{
		"user_efficiency_trends": trends,
	}

	m.writeJSONResponse(w, response)
}

func (m *MockSlurmServer) handleCompareJobPerformance(w http.ResponseWriter, r *http.Request) {
	if !m.config.SupportedOperations["jobs.compare_performance"] {
		m.writeErrorResponse(w, http.StatusNotImplemented, "Operation not supported")
		return
	}

	var request struct {
		JobAID string `json:"job_a_id"`
		JobBID string `json:"job_b_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		m.writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	comparison := map[string]interface{}{
		"job_a_id": request.JobAID,
		"job_b_id": request.JobBID,
		"metrics": map[string]interface{}{
			"overall_efficiency_delta": 5.2, // Job B is 5.2% more efficient
			"cpu_efficiency_delta":     3.1,
			"memory_efficiency_delta":  -2.8, // Job A is better at memory
			"runtime_ratio":            0.85, // Job B took 85% of Job A's time
		},
		"resource_differences": map[string]interface{}{
			"cpu_delta":       2,    // Job B used 2 more CPUs
			"memory_delta_gb": -1.5, // Job B used 1.5GB less memory
			"gpu_delta":       0,
		},
		"winner": "job_b",
		"winner_reasons": []string{
			"Better overall efficiency",
			"Faster execution time",
			"More efficient CPU usage",
		},
	}

	response := map[string]interface{}{
		"performance_comparison": comparison,
	}

	m.writeJSONResponse(w, response)
}

func (m *MockSlurmServer) handleSimilarJobsPerformance(w http.ResponseWriter, r *http.Request) {
	if !m.config.SupportedOperations["jobs.similar_performance"] {
		m.writeErrorResponse(w, http.StatusNotImplemented, "Operation not supported")
		return
	}

	referenceJobID := parseQueryParam(r, "reference_job_id")
	if referenceJobID == "" {
		m.writeErrorResponse(w, http.StatusBadRequest, "reference_job_id parameter is required")
		return
	}

	similarJobs := []map[string]interface{}{
		{
			"job_id":           "1234",
			"similarity_score": 0.92,
			"efficiency_delta": 3.5, // 3.5% more efficient
			"runtime_ratio":    0.88,
			"performance_rank": 1,
		},
		{
			"job_id":           "1235",
			"similarity_score": 0.87,
			"efficiency_delta": -1.2, // 1.2% less efficient
			"runtime_ratio":    1.15,
			"performance_rank": 2,
		},
		{
			"job_id":           "1236",
			"similarity_score": 0.83,
			"efficiency_delta": 7.8, // 7.8% more efficient
			"runtime_ratio":    0.72,
			"performance_rank": 3,
		},
	}

	analysis := map[string]interface{}{
		"reference_job_id": referenceJobID,
		"similar_jobs":     similarJobs,
		"analysis_summary": map[string]interface{}{
			"best_performer":           "1236",
			"worst_performer":          "1235",
			"average_efficiency_delta": 3.37,
			"recommendations": []string{
				"Consider adopting configuration from job 1236",
				"Review memory allocation patterns",
			},
		},
	}

	response := map[string]interface{}{
		"similar_jobs_performance": analysis,
	}

	m.writeJSONResponse(w, response)
}

func (m *MockSlurmServer) handleAnalyzeBatchJobs(w http.ResponseWriter, r *http.Request) {
	if !m.config.SupportedOperations["jobs.analyze_batch"] {
		m.writeErrorResponse(w, http.StatusNotImplemented, "Operation not supported")
		return
	}

	var request struct {
		JobIDs []string `json:"job_ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		m.writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Build individual job analyses
	jobAnalyses := make([]map[string]interface{}, 0, len(request.JobIDs))
	analyzedCount := 0
	failedCount := 0

	for _, jobID := range request.JobIDs {
		// Determine status based on job ID
		status := "completed"
		if jobID == "invalid_job" || jobID == "99999" {
			status = "failed"
			failedCount++
		} else {
			analyzedCount++
		}

		jobAnalysis := map[string]interface{}{
			"job_id": jobID,
			"status": status,
		}

		if status == "completed" {
			jobAnalysis["efficiency"] = 85.5
			jobAnalysis["cpu_utilization"] = 90.0
			jobAnalysis["memory_utilization"] = 75.0
			jobAnalysis["runtime_efficiency"] = 88.0
		} else {
			jobAnalysis["issues"] = []string{"Job not found or failed to analyze"}
		}

		jobAnalyses = append(jobAnalyses, jobAnalysis)
	}

	analysis := map[string]interface{}{
		"job_count":      len(request.JobIDs),
		"analyzed_count": analyzedCount,
		"failed_count":   failedCount,
		"job_analyses":   jobAnalyses,
		"analysis_summary": map[string]interface{}{
			"average_efficiency":   76.8,
			"efficiency_std_dev":   12.3,
			"best_performing_job":  request.JobIDs[0],
			"worst_performing_job": request.JobIDs[len(request.JobIDs)-1],
			"total_resource_waste": map[string]interface{}{
				"cpu_core_hours":  245.6,
				"memory_gb_hours": 1024.3,
			},
		},
		"efficiency_distribution": map[string]interface{}{
			"excellent": 15, // Jobs with >90% efficiency
			"good":      45, // Jobs with 75-90% efficiency
			"fair":      25, // Jobs with 60-75% efficiency
			"poor":      15, // Jobs with <60% efficiency
		},
		"recommendations": []map[string]interface{}{
			{
				"category": "Resource Optimization",
				"message":  "25% of jobs show significant memory over-allocation",
				"priority": "High",
			},
			{
				"category": "Performance Tuning",
				"message":  "CPU efficiency could be improved by 12% on average",
				"priority": "Medium",
			},
		},
	}

	response := map[string]interface{}{
		"batch_analysis": analysis,
	}

	m.writeJSONResponse(w, response)
}

func (m *MockSlurmServer) handleWorkflowPerformance(w http.ResponseWriter, r *http.Request) {
	if !m.config.SupportedOperations["jobs.workflow_performance"] {
		m.writeErrorResponse(w, http.StatusNotImplemented, "Operation not supported")
		return
	}

	workflowID := parseQueryParam(r, "workflow_id")
	if workflowID == "" {
		m.writeErrorResponse(w, http.StatusBadRequest, "workflow_id parameter is required")
		return
	}

	workflow := map[string]interface{}{
		"workflow_id": workflowID,
		"total_jobs":  5,
		"performance_summary": map[string]interface{}{
			"total_runtime_hours": 12.5,
			"total_cpu_hours":     62.5,
			"average_efficiency":  81.2,
			"critical_path_jobs":  []string{"job_1", "job_3", "job_5"},
			"bottleneck_job":      "job_3",
		},
		"job_performance": []map[string]interface{}{
			{
				"job_id":        "job_1",
				"efficiency":    85.2,
				"runtime_hours": 2.1,
				"cpu_hours":     8.4,
				"critical_path": true,
			},
			{
				"job_id":        "job_2",
				"efficiency":    78.9,
				"runtime_hours": 1.8,
				"cpu_hours":     7.2,
				"critical_path": false,
			},
			{
				"job_id":        "job_3",
				"efficiency":    72.1,
				"runtime_hours": 4.5,
				"cpu_hours":     18.0,
				"critical_path": true,
			},
			{
				"job_id":        "job_4",
				"efficiency":    83.6,
				"runtime_hours": 2.2,
				"cpu_hours":     8.8,
				"critical_path": false,
			},
			{
				"job_id":        "job_5",
				"efficiency":    86.1,
				"runtime_hours": 1.9,
				"cpu_hours":     7.6,
				"critical_path": true,
			},
		},
		"optimization_opportunities": []map[string]interface{}{
			{
				"job_id":                "job_3",
				"potential_improvement": "Optimize job_3 to reduce workflow runtime by 15%",
				"resource":              "CPU",
				"impact":                "High",
			},
		},
	}

	response := map[string]interface{}{
		"workflow_performance": workflow,
	}

	m.writeJSONResponse(w, response)
}

func (m *MockSlurmServer) handleGenerateEfficiencyReport(w http.ResponseWriter, r *http.Request) {
	if !m.config.SupportedOperations["jobs.efficiency_report"] {
		m.writeErrorResponse(w, http.StatusNotImplemented, "Operation not supported")
		return
	}

	// Parse query parameters for report scope
	userID := parseQueryParam(r, "user_id")
	partition := parseQueryParam(r, "partition")
	days := parseQueryParamInt(r, "days", 30)

	report := map[string]interface{}{
		"report_metadata": map[string]interface{}{
			"generated_at": time.Now().Unix(),
			"period_days":  days,
			"user_id":      userID,
			"partition":    partition,
		},
		"executive_summary": map[string]interface{}{
			"total_jobs_analyzed": 1247,
			"average_efficiency":  78.5,
			"total_resource_waste": map[string]interface{}{
				"cpu_core_hours":  2450.6,
				"memory_gb_hours": 8192.3,
				"estimated_cost":  "$1,245.50",
			},
			"efficiency_trend": "improving",
			"top_opportunities": []string{
				"Memory over-allocation in compute partition",
				"CPU underutilization in ML workloads",
				"I/O inefficiencies in data processing jobs",
			},
		},
		"detailed_metrics": map[string]interface{}{
			"resource_efficiency": map[string]interface{}{
				"cpu_efficiency":     82.1,
				"memory_efficiency":  75.6,
				"gpu_efficiency":     68.9,
				"io_efficiency":      65.2,
				"network_efficiency": 45.8,
			},
			"partition_breakdown": []map[string]interface{}{
				{
					"partition":       "compute",
					"jobs_count":      856,
					"avg_efficiency":  79.8,
					"waste_cpu_hours": 1245.2,
					"waste_memory_gb": 4096.1,
				},
				{
					"partition":       "gpu",
					"jobs_count":      391,
					"avg_efficiency":  75.2,
					"waste_cpu_hours": 856.4,
					"waste_memory_gb": 2048.7,
				},
			},
		},
		"recommendations": []map[string]interface{}{
			{
				"priority":    "High",
				"category":    "Memory Optimization",
				"description": "Reduce default memory allocation by 20% for compute partition",
				"impact":      "Save $500/month in resource costs",
			},
			{
				"priority":    "Medium",
				"category":    "CPU Optimization",
				"description": "Implement CPU affinity for parallel jobs",
				"impact":      "Improve efficiency by 8-12%",
			},
			{
				"priority":    "Low",
				"category":    "I/O Optimization",
				"description": "Consider local SSD storage for high I/O workloads",
				"impact":      "Reduce job runtime by 15-25%",
			},
		},
	}

	response := map[string]interface{}{
		"efficiency_report": report,
	}

	m.writeJSONResponse(w, response)
}
