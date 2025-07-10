package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// Job endpoint handlers

func (m *MockSlurmServer) handleJobsList(w http.ResponseWriter, r *http.Request) {
	if !m.config.SupportedOperations["jobs.list"] {
		m.writeErrorResponse(w, http.StatusNotImplemented, "Operation not supported")
		return
	}

	// Parse query parameters
	userID := parseQueryParam(r, "user_id")
	partition := parseQueryParam(r, "partition")
	states := parseQueryParam(r, "state")
	limit := parseQueryParamInt(r, "limit", 100)
	offset := parseQueryParamInt(r, "offset", 0)

	stateList := []string{}
	if states != "" {
		stateList = strings.Split(states, ",")
	}

	m.storage.mu.RLock()
	defer m.storage.mu.RUnlock()

	jobs := []*MockJob{}
	for _, job := range m.storage.Jobs {
		// Apply filters
		if userID != "" && job.UserID != userID {
			continue
		}
		if partition != "" && job.Partition != partition {
			continue
		}
		if len(stateList) > 0 {
			found := false
			for _, state := range stateList {
				if strings.TrimSpace(state) == job.State {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		jobs = append(jobs, job)
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
		jobs = []*MockJob{}
	}

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

	m.storage.mu.RLock()
	job, exists := m.storage.Jobs[jobID]
	m.storage.mu.RUnlock()

	if !exists {
		m.writeErrorResponse(w, http.StatusNotFound, fmt.Sprintf("Job %s not found", jobID))
		return
	}

	response := map[string]interface{}{
		"job": job,
	}

	m.writeJSONResponse(w, response)
}

func (m *MockSlurmServer) handleJobsSubmit(w http.ResponseWriter, r *http.Request) {
	if !m.config.SupportedOperations["jobs.submit"] {
		m.writeErrorResponse(w, http.StatusNotImplemented, "Operation not supported")
		return
	}

	var submission struct {
		Name        string            `json:"name"`
		Script      string            `json:"script"`
		Partition   string            `json:"partition"`
		CPUs        int               `json:"cpus"`
		Memory      int64             `json:"memory"`
		TimeLimit   int               `json:"time_limit"`
		WorkingDir  string            `json:"working_directory"`
		Environment map[string]string `json:"environment"`
	}

	if err := json.NewDecoder(r.Body).Decode(&submission); err != nil {
		m.writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Validate submission
	if submission.Name == "" {
		m.writeErrorResponse(w, http.StatusBadRequest, "Job name is required")
		return
	}
	if submission.Script == "" {
		m.writeErrorResponse(w, http.StatusBadRequest, "Job script is required")
		return
	}

	// Create new job
	jobID := m.generateJobID()
	job := &MockJob{
		JobID:       jobID,
		Name:        submission.Name,
		UserID:      "testuser", // In real implementation, this would come from auth
		State:       "PENDING",
		Partition:   submission.Partition,
		CPUs:        submission.CPUs,
		Memory:      submission.Memory,
		TimeLimit:   submission.TimeLimit,
		SubmitTime:  time.Now().Unix(),
		WorkingDir:  submission.WorkingDir,
		Script:      submission.Script,
		Environment: submission.Environment,
	}

	if job.Partition == "" {
		job.Partition = "compute" // Default partition
	}
	if job.CPUs == 0 {
		job.CPUs = 1 // Default CPU count
	}
	if job.TimeLimit == 0 {
		job.TimeLimit = 60 // Default time limit
	}
	if job.WorkingDir == "" {
		job.WorkingDir = "/tmp" // Default working directory
	}

	m.storage.mu.Lock()
	m.storage.Jobs[jobID] = job
	m.storage.mu.Unlock()

	response := map[string]interface{}{
		"job_id": jobID,
		"result": "SUCCESS",
	}

	w.WriteHeader(http.StatusCreated)
	m.writeJSONResponse(w, response)
}

func (m *MockSlurmServer) handleJobsCancel(w http.ResponseWriter, r *http.Request) {
	if !m.config.SupportedOperations["jobs.cancel"] {
		m.writeErrorResponse(w, http.StatusNotImplemented, "Operation not supported")
		return
	}

	vars := mux.Vars(r)
	jobID := vars["job_id"]

	m.storage.mu.Lock()
	job, exists := m.storage.Jobs[jobID]
	if exists {
		if job.State == "RUNNING" || job.State == "PENDING" {
			job.State = "CANCELLED"
			now := time.Now().Unix()
			job.EndTime = &now
		} else {
			m.storage.mu.Unlock()
			m.writeErrorResponse(w, http.StatusConflict, fmt.Sprintf("Job %s cannot be cancelled (state: %s)", jobID, job.State))
			return
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
	job, exists := m.storage.Jobs[jobID]
	if exists {
		if update.Name != nil {
			job.Name = *update.Name
		}
		if update.TimeLimit != nil {
			job.TimeLimit = *update.TimeLimit
		}
		// Priority field would be added to MockJob if needed
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
				if strings.TrimSpace(state) == node.State {
					found = true
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

	response := map[string]interface{}{
		"node": node,
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
		State *string `json:"state,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		m.writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	m.storage.mu.Lock()
	node, exists := m.storage.Nodes[nodeName]
	if exists && update.State != nil {
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

	response := map[string]interface{}{
		"partition": partition,
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

	response := map[string]interface{}{
		"status": "OK",
		"time":   time.Now().Unix(),
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
	for _, job := range m.storage.Jobs {
		switch job.State {
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