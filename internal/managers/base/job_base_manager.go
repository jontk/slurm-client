// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package base

import (
	"fmt"
	"strings"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/pkg/errors"
)

// JobBaseManager provides common job management functionality
type JobBaseManager struct {
	*CRUDManager
}

// NewJobBaseManager creates a new job base manager
func NewJobBaseManager(version string) *JobBaseManager {
	return &JobBaseManager{
		CRUDManager: NewCRUDManager(version, "Job"),
	}
}

// ValidateJobCreate validates job creation data
func (m *JobBaseManager) ValidateJobCreate(job *types.JobCreate) error {
	if job == nil {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Job data is required",
			"job", job, nil,
		)
	}

	// Validate required fields
	if job.Command == "" && job.Script == "" {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Either command or script is required",
			"job.Command/Script", nil, nil,
		)
	}

	// Validate numeric fields
	if err := m.ValidateNonNegative(int(job.CPUs), "job.CPUs"); err != nil {
		return err
	}
	if err := m.ValidateNonNegative(int(job.Nodes), "job.Nodes"); err != nil {
		return err
	}
	if err := m.ValidateNonNegative(int(job.Tasks), "job.Tasks"); err != nil {
		return err
	}
	if err := m.ValidateNonNegative(int(job.TimeLimit), "job.TimeLimit"); err != nil {
		return err
	}

	// Validate resource requests
	if err := m.ValidateResourceRequests(&job.ResourceRequests); err != nil {
		return err
	}

	// Validate array job string format
	if job.ArrayString != "" {
		if err := m.ValidateArrayString(job.ArrayString); err != nil {
			return err
		}
	}

	// Validate dependencies
	if len(job.Dependencies) > 0 {
		if err := m.ValidateJobDependencies(job.Dependencies); err != nil {
			return err
		}
	}

	// Validate mail type
	if len(job.MailType) > 0 {
		if err := m.ValidateMailTypes(job.MailType); err != nil {
			return err
		}
	}

	// Validate CPU frequency
	if job.CPUFrequencyMin > 0 && job.CPUFrequencyMax > 0 {
		if job.CPUFrequencyMin > job.CPUFrequencyMax {
			return errors.NewValidationError(
				errors.ErrorCodeValidationFailed,
				"CPU frequency min cannot be greater than max",
				"job.CPUFrequency", nil, nil,
			)
		}
	}

	return nil
}

// ValidateJobUpdate validates job update data
func (m *JobBaseManager) ValidateJobUpdate(update *types.JobUpdate) error {
	if update == nil {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Update data is required",
			"update", update, nil,
		)
	}

	// Validate time limit if provided
	if update.TimeLimit != nil {
		if err := m.ValidateNonNegative(int(*update.TimeLimit), "update.TimeLimit"); err != nil {
			return err
		}
	}

	// Validate priority if provided
	if update.Priority != nil {
		if err := m.ValidateNonNegative(int(*update.Priority), "update.Priority"); err != nil {
			return err
		}
	}

	// Validate node counts if provided
	if update.MinNodes != nil && update.MaxNodes != nil {
		if *update.MinNodes > *update.MaxNodes {
			return errors.NewValidationError(
				errors.ErrorCodeValidationFailed,
				"Min nodes cannot be greater than max nodes",
				"update.MinNodes/MaxNodes", nil, nil,
			)
		}
	}

	return nil
}

// ValidateResourceRequests validates resource request specifications
func (m *JobBaseManager) ValidateResourceRequests(req *types.ResourceRequests) error {
	if req.Memory < 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Memory must be non-negative",
			"ResourceRequests.Memory", req.Memory, nil,
		)
	}
	if req.MemoryPerCPU < 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Memory per CPU must be non-negative",
			"ResourceRequests.MemoryPerCPU", req.MemoryPerCPU, nil,
		)
	}
	if req.MemoryPerGPU < 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Memory per GPU must be non-negative",
			"ResourceRequests.MemoryPerGPU", req.MemoryPerGPU, nil,
		)
	}
	if req.TmpDisk < 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Temporary disk space must be non-negative",
			"ResourceRequests.TmpDisk", req.TmpDisk, nil,
		)
	}
	return nil
}

// ValidateArrayString validates job array specification string
func (m *JobBaseManager) ValidateArrayString(arrayStr string) error {
	// Basic validation of array string format
	// Format examples: "1-10", "1-10:2", "1,3,5,7", "1-5,10-15"
	if arrayStr == "" {
		return nil
	}

	// Check for valid characters
	for _, ch := range arrayStr {
		if !((ch >= '0' && ch <= '9') || ch == '-' || ch == ',' || ch == ':' || ch == '%') {
			return errors.NewValidationError(
				errors.ErrorCodeValidationFailed,
				fmt.Sprintf("Invalid character '%c' in array string", ch),
				"arrayString", arrayStr, nil,
			)
		}
	}

	return nil
}

// ValidateJobDependencies validates job dependency specifications
func (m *JobBaseManager) ValidateJobDependencies(deps []types.JobDependency) error {
	validTypes := map[string]bool{
		"after":        true,
		"afterany":     true,
		"aftercorr":    true,
		"afternotok":   true,
		"afterok":      true,
		"expand":       true,
		"singleton":    true,
		"burstbuffer":  true,
	}

	for _, dep := range deps {
		if _, ok := validTypes[strings.ToLower(dep.Type)]; !ok {
			return errors.NewValidationError(
				errors.ErrorCodeValidationFailed,
				fmt.Sprintf("Invalid dependency type: %s", dep.Type),
				"dependency.Type", dep.Type, nil,
			)
		}

		// Validate job IDs are positive
		for _, jobID := range dep.JobIDs {
			if jobID <= 0 {
				return errors.NewValidationError(
					errors.ErrorCodeValidationFailed,
					"Job ID in dependency must be positive",
					"dependency.JobID", jobID, nil,
				)
			}
		}
	}

	return nil
}

// ValidateMailTypes validates mail notification types
func (m *JobBaseManager) ValidateMailTypes(mailTypes []string) error {
	validTypes := map[string]bool{
		"NONE":       true,
		"BEGIN":      true,
		"END":        true,
		"FAIL":       true,
		"REQUEUE":    true,
		"ALL":        true,
		"STAGE_OUT":  true,
		"TIME_LIMIT": true,
		"TIME_LIMIT_90": true,
		"TIME_LIMIT_80": true,
		"TIME_LIMIT_50": true,
		"ARRAY_TASKS": true,
	}

	for _, mailType := range mailTypes {
		if _, ok := validTypes[strings.ToUpper(mailType)]; !ok {
			return errors.NewValidationError(
				errors.ErrorCodeValidationFailed,
				fmt.Sprintf("Invalid mail type: %s", mailType),
				"mailType", mailType, nil,
			)
		}
	}

	return nil
}

// ApplyJobDefaults applies default values to job create request
func (m *JobBaseManager) ApplyJobDefaults(job *types.JobCreate) *types.JobCreate {
	// Apply default values
	if job.CPUs == 0 {
		job.CPUs = 1
	}
	if job.Nodes == 0 {
		job.Nodes = 1
	}
	if job.Tasks == 0 {
		job.Tasks = 1
	}

	// Initialize arrays if nil
	if job.Environment == nil {
		job.Environment = make(map[string]string)
	}
	if job.Dependencies == nil {
		job.Dependencies = []types.JobDependency{}
	}
	if job.MailType == nil {
		job.MailType = []string{}
	}
	if job.Features == nil {
		job.Features = []string{}
	}
	if job.ClusterFeatures == nil {
		job.ClusterFeatures = []string{}
	}
	if job.Profile == nil {
		job.Profile = []string{}
	}

	// Set default working directory if not specified
	if job.WorkingDirectory == "" {
		job.WorkingDirectory = "/tmp"
	}

	return job
}

// FilterJobList applies filtering to a job list
func (m *JobBaseManager) FilterJobList(jobs []types.Job, opts *types.JobListOptions) []types.Job {
	if opts == nil {
		return jobs
	}

	filtered := make([]types.Job, 0, len(jobs))
	for _, job := range jobs {
		if m.matchesJobFilters(job, opts) {
			filtered = append(filtered, job)
		}
	}

	return filtered
}

// matchesJobFilters checks if a job matches the given filters
func (m *JobBaseManager) matchesJobFilters(job types.Job, opts *types.JobListOptions) bool {
	// Filter by job IDs
	if len(opts.JobIDs) > 0 {
		found := false
		for _, id := range opts.JobIDs {
			if job.JobID == id {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Filter by job names
	if len(opts.JobNames) > 0 {
		found := false
		for _, name := range opts.JobNames {
			if strings.EqualFold(job.Name, name) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Filter by accounts
	if len(opts.Accounts) > 0 {
		found := false
		for _, account := range opts.Accounts {
			if strings.EqualFold(job.Account, account) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Filter by users
	if len(opts.Users) > 0 {
		found := false
		for _, user := range opts.Users {
			if strings.EqualFold(job.UserName, user) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Filter by states
	if len(opts.States) > 0 {
		found := false
		for _, state := range opts.States {
			if job.State == state {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Filter by partitions
	if len(opts.Partitions) > 0 {
		found := false
		for _, partition := range opts.Partitions {
			if strings.EqualFold(job.Partition, partition) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Filter by QoS
	if len(opts.QoS) > 0 {
		found := false
		for _, qos := range opts.QoS {
			if strings.EqualFold(job.QoS, qos) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Filter by time range
	if opts.StartTime != nil && job.SubmitTime.Before(*opts.StartTime) {
		return false
	}
	if opts.EndTime != nil && job.SubmitTime.After(*opts.EndTime) {
		return false
	}

	return true
}

// ValidateJobSignal validates job signal request
func (m *JobBaseManager) ValidateJobSignal(req *types.JobSignalRequest) error {
	if req == nil {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Signal request is required",
			"request", req, nil,
		)
	}

	if req.JobID <= 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Job ID must be positive",
			"request.JobID", req.JobID, nil,
		)
	}

	if req.Signal == "" {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Signal is required",
			"request.Signal", req.Signal, nil,
		)
	}

	return nil
}

// ValidateJobHold validates job hold request
func (m *JobBaseManager) ValidateJobHold(req *types.JobHoldRequest) error {
	if req == nil {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Hold request is required",
			"request", req, nil,
		)
	}

	if req.JobID <= 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Job ID must be positive",
			"request.JobID", req.JobID, nil,
		)
	}

	return nil
}

// ValidateJobNotify validates job notify request
func (m *JobBaseManager) ValidateJobNotify(req *types.JobNotifyRequest) error {
	if req == nil {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Notify request is required",
			"request", req, nil,
		)
	}

	if req.JobID <= 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Job ID must be positive",
			"request.JobID", req.JobID, nil,
		)
	}

	if req.Message == "" {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Message is required",
			"request.Message", req.Message, nil,
		)
	}

	return nil
}
