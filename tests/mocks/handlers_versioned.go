// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package mocks

import (
	"net/http"
	"strconv"
	"strings"

	v0_0_40 "github.com/jontk/slurm-client/internal/api/v0_0_40"
	v0_0_42 "github.com/jontk/slurm-client/internal/api/v0_0_42"
	v0_0_43 "github.com/jontk/slurm-client/internal/api/v0_0_43"
	v0_0_44 "github.com/jontk/slurm-client/internal/api/v0_0_44"
)

// JobHandler provides version-specific job handling
type JobHandler interface {
	// ListJobs returns filtered jobs for the version
	ListJobs(r *http.Request, storage *MockStorage) ([]interface{}, error)

	// GetJob returns a specific job
	GetJob(jobID string, storage *MockStorage) (interface{}, error)
}

// V0040JobHandler handles v0.0.40 jobs
type V0040JobHandler struct{}

func (h *V0040JobHandler) ListJobs(r *http.Request, storage *MockStorage) ([]interface{}, error) {
	// Parse filters
	userID := parseQueryParam(r, "user_id")
	partition := parseQueryParam(r, "partition")
	states := parseQueryParam(r, "state")

	stateList := []string{}
	if states != "" {
		stateList = strings.Split(states, ",")
	}

	storage.mu.RLock()
	defer storage.mu.RUnlock()

	var jobs []interface{}
	for _, jobInterface := range storage.Jobs {
		job, ok := jobInterface.(*v0_0_40.V0040JobInfo)
		if !ok {
			continue
		}

		// Apply filters
		if !h.matchesFilters(job, userID, partition, stateList) {
			continue
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}

func (h *V0040JobHandler) GetJob(jobID string, storage *MockStorage) (interface{}, error) {
	storage.mu.RLock()
	defer storage.mu.RUnlock()

	jobInterface, exists := storage.Jobs[jobID]
	if !exists {
		return nil, nil //nolint:nilnil // Intentional: nil job with nil error means "not found" (caller checks job == nil)
	}

	job, ok := jobInterface.(*v0_0_40.V0040JobInfo)
	if !ok {
		return nil, nil //nolint:nilnil // Intentional: wrong type treated as not found
	}

	return job, nil
}

func (h *V0040JobHandler) matchesFilters(job *v0_0_40.V0040JobInfo, userID, partition string, states []string) bool {
	// User filter
	if userID != "" {
		uid, err := strconv.ParseInt(userID, 10, 32)
		if err == nil && job.UserId != nil && *job.UserId != int32(uid) {
			return false
		}
	}

	// Partition filter
	if partition != "" && job.Partition != nil && *job.Partition != partition {
		return false
	}

	// State filter
	if len(states) > 0 {
		found := false
		if job.JobState != nil {
			for _, state := range states {
				for _, jobState := range *job.JobState {
					if strings.TrimSpace(state) == jobState {
						found = true
						break
					}
				}
				if found {
					break
				}
			}
		}
		if !found {
			return false
		}
	}

	return true
}

// V0042JobHandler handles v0.0.42 jobs
type V0042JobHandler struct{}

func (h *V0042JobHandler) ListJobs(r *http.Request, storage *MockStorage) ([]interface{}, error) {
	userID := parseQueryParam(r, "user_id")
	partition := parseQueryParam(r, "partition")
	states := parseQueryParam(r, "state")

	stateList := []string{}
	if states != "" {
		stateList = strings.Split(states, ",")
	}

	storage.mu.RLock()
	defer storage.mu.RUnlock()

	var jobs []interface{}
	for _, jobInterface := range storage.Jobs {
		job, ok := jobInterface.(*v0_0_42.V0042JobInfo)
		if !ok {
			continue
		}

		if !h.matchesFilters(job, userID, partition, stateList) {
			continue
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}

func (h *V0042JobHandler) GetJob(jobID string, storage *MockStorage) (interface{}, error) {
	storage.mu.RLock()
	defer storage.mu.RUnlock()

	jobInterface, exists := storage.Jobs[jobID]
	if !exists {
		return nil, nil //nolint:nilnil // Intentional: nil job with nil error means "not found" (caller checks job == nil)
	}

	job, ok := jobInterface.(*v0_0_42.V0042JobInfo)
	if !ok {
		return nil, nil //nolint:nilnil // Intentional: wrong type treated as not found
	}

	return job, nil
}

func (h *V0042JobHandler) matchesFilters(job *v0_0_42.V0042JobInfo, userID, partition string, states []string) bool {
	if userID != "" {
		uid, err := strconv.ParseInt(userID, 10, 32)
		if err == nil && job.UserId != nil && *job.UserId != int32(uid) {
			return false
		}
	}

	if partition != "" && job.Partition != nil && *job.Partition != partition {
		return false
	}

	if len(states) > 0 {
		found := false
		if job.JobState != nil {
			for _, state := range states {
				for _, jobState := range *job.JobState {
					if strings.TrimSpace(state) == jobState {
						found = true
						break
					}
				}
				if found {
					break
				}
			}
		}
		if !found {
			return false
		}
	}

	return true
}

// V0043JobHandler handles v0.0.43 jobs
type V0043JobHandler struct{}

func (h *V0043JobHandler) ListJobs(r *http.Request, storage *MockStorage) ([]interface{}, error) {
	userID := parseQueryParam(r, "user_id")
	partition := parseQueryParam(r, "partition")
	states := parseQueryParam(r, "state")

	stateList := []string{}
	if states != "" {
		stateList = strings.Split(states, ",")
	}

	storage.mu.RLock()
	defer storage.mu.RUnlock()

	var jobs []interface{}
	for _, jobInterface := range storage.Jobs {
		job, ok := jobInterface.(*v0_0_43.V0043JobInfo)
		if !ok {
			continue
		}

		if !h.matchesFilters(job, userID, partition, stateList) {
			continue
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}

func (h *V0043JobHandler) GetJob(jobID string, storage *MockStorage) (interface{}, error) {
	storage.mu.RLock()
	defer storage.mu.RUnlock()

	jobInterface, exists := storage.Jobs[jobID]
	if !exists {
		return nil, nil //nolint:nilnil // Intentional: nil job with nil error means "not found" (caller checks job == nil)
	}

	job, ok := jobInterface.(*v0_0_43.V0043JobInfo)
	if !ok {
		return nil, nil //nolint:nilnil // Intentional: wrong type treated as not found
	}

	return job, nil
}

func (h *V0043JobHandler) matchesFilters(job *v0_0_43.V0043JobInfo, userID, partition string, states []string) bool {
	if userID != "" {
		uid, err := strconv.ParseInt(userID, 10, 32)
		if err == nil && job.UserId != nil && *job.UserId != int32(uid) {
			return false
		}
	}

	if partition != "" && job.Partition != nil && *job.Partition != partition {
		return false
	}

	if len(states) > 0 {
		found := false
		if job.JobState != nil {
			for _, state := range states {
				for _, jobState := range *job.JobState {
					// v0.0.43 JobState is an enum type, convert to string
					if strings.TrimSpace(state) == string(jobState) {
						found = true
						break
					}
				}
				if found {
					break
				}
			}
		}
		if !found {
			return false
		}
	}

	return true
}

// V0044JobHandler handles v0.0.44 jobs
type V0044JobHandler struct{}

func (h *V0044JobHandler) ListJobs(r *http.Request, storage *MockStorage) ([]interface{}, error) {
	userID := parseQueryParam(r, "user_id")
	partition := parseQueryParam(r, "partition")
	states := parseQueryParam(r, "state")

	stateList := []string{}
	if states != "" {
		stateList = strings.Split(states, ",")
	}

	storage.mu.RLock()
	defer storage.mu.RUnlock()

	var jobs []interface{}
	for _, jobInterface := range storage.Jobs {
		job, ok := jobInterface.(*v0_0_44.V0044JobInfo)
		if !ok {
			continue
		}

		if !h.matchesFilters(job, userID, partition, stateList) {
			continue
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}

func (h *V0044JobHandler) GetJob(jobID string, storage *MockStorage) (interface{}, error) {
	storage.mu.RLock()
	defer storage.mu.RUnlock()

	jobInterface, exists := storage.Jobs[jobID]
	if !exists {
		return nil, nil //nolint:nilnil // Intentional: nil job with nil error means "not found" (caller checks job == nil)
	}

	job, ok := jobInterface.(*v0_0_44.V0044JobInfo)
	if !ok {
		return nil, nil //nolint:nilnil // Intentional: wrong type treated as not found
	}

	return job, nil
}

func (h *V0044JobHandler) matchesFilters(job *v0_0_44.V0044JobInfo, userID, partition string, states []string) bool {
	if userID != "" {
		uid, err := strconv.ParseInt(userID, 10, 32)
		if err == nil && job.UserId != nil && *job.UserId != int32(uid) {
			return false
		}
	}

	if partition != "" && job.Partition != nil && *job.Partition != partition {
		return false
	}

	if len(states) > 0 {
		found := false
		if job.JobState != nil {
			for _, state := range states {
				for _, jobState := range *job.JobState {
					// v0.0.44 JobState is an enum type, convert to string
					if strings.TrimSpace(state) == string(jobState) {
						found = true
						break
					}
				}
				if found {
					break
				}
			}
		}
		if !found {
			return false
		}
	}

	return true
}

// NewJobHandler returns the appropriate handler for the version
func NewJobHandler(version string) JobHandler {
	switch version {
	case "v0.0.40":
		return &V0040JobHandler{}
	case "v0.0.42":
		return &V0042JobHandler{}
	case "v0.0.43":
		return &V0043JobHandler{}
	case "v0.0.44":
		return &V0044JobHandler{}
	default:
		// Fallback to v0.0.40
		return &V0040JobHandler{}
	}
}
