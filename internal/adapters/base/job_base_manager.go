// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package base

import (
	"fmt"
	"strings"
	"time"

	types "github.com/jontk/slurm-client/api"
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
	// Validate required fields - Script is required
	if job.Script == nil || *job.Script == "" {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Script is required",
			"job.Script", nil, nil,
		)
	}
	// Validate numeric fields (only if set)
	if job.MinimumCPUs != nil {
		if err := m.ValidateNonNegative(int(*job.MinimumCPUs), "job.MinimumCPUs"); err != nil {
			return err
		}
	}
	if job.MinimumNodes != nil {
		if err := m.ValidateNonNegative(int(*job.MinimumNodes), "job.MinimumNodes"); err != nil {
			return err
		}
	}
	if job.Tasks != nil {
		if err := m.ValidateNonNegative(int(*job.Tasks), "job.Tasks"); err != nil {
			return err
		}
	}
	if job.TimeLimit != nil {
		if err := m.ValidateNonNegative(int(*job.TimeLimit), "job.TimeLimit"); err != nil {
			return err
		}
	}
	// Validate memory fields
	if job.MemoryPerCPU != nil {
		if err := m.ValidateNonNegative(int(*job.MemoryPerCPU), "job.MemoryPerCPU"); err != nil {
			return err
		}
	}
	if job.MemoryPerNode != nil {
		if err := m.ValidateNonNegative(int(*job.MemoryPerNode), "job.MemoryPerNode"); err != nil {
			return err
		}
	}
	// Validate array job string format
	if job.Array != nil && *job.Array != "" {
		if err := m.ValidateArrayString(*job.Array); err != nil {
			return err
		}
	}
	// Validate dependency string format (SLURM format: "afterok:123:456")
	if job.Dependency != nil && *job.Dependency != "" {
		if err := m.ValidateDependencyString(*job.Dependency); err != nil {
			return err
		}
	}
	// Validate mail type
	if len(job.MailType) > 0 {
		if err := m.ValidateMailTypeValues(job.MailType); err != nil {
			return err
		}
	}
	return nil
}

// ValidateJobUpdate validates job update data
// Note: JobUpdate is an alias for JobCreate since SLURM uses the same job_desc_msg
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
	if update.MinimumNodes != nil && update.MaximumNodes != nil {
		if *update.MinimumNodes > *update.MaximumNodes {
			return errors.NewValidationError(
				errors.ErrorCodeValidationFailed,
				"Min nodes cannot be greater than max nodes",
				"update.MinimumNodes/MaximumNodes", nil, nil,
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
	if req.TemporaryDisk < 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Temporary disk space must be non-negative",
			"ResourceRequests.TemporaryDisk", req.TemporaryDisk, nil,
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
		"after":       true,
		"afterany":    true,
		"aftercorr":   true,
		"afternotok":  true,
		"afterok":     true,
		"expand":      true,
		"singleton":   true,
		"burstbuffer": true,
	}
	for _, dep := range deps {
		if _, ok := validTypes[strings.ToLower(dep.Type)]; !ok {
			return errors.NewValidationError(
				errors.ErrorCodeValidationFailed,
				"Invalid dependency type: "+dep.Type,
				"dependency.Type", dep.Type, nil,
			)
		}
		// Validate job IDs are positive
		for _, jobID := range dep.JobIDs {
			if jobID <= 0 {
				return errors.NewValidationError(
					errors.ErrorCodeValidationFailed,
					"Job ID in dependency must be positive",
					"dependency.JobId", jobID, nil,
				)
			}
		}
	}
	return nil
}

// ValidateMailTypes validates mail notification types (legacy string-based)
func (m *JobBaseManager) ValidateMailTypes(mailTypes []string) error {
	validTypes := map[string]bool{
		"NONE":          true,
		"BEGIN":         true,
		"END":           true,
		"FAIL":          true,
		"REQUEUE":       true,
		"ALL":           true,
		"STAGE_OUT":     true,
		"TIME_LIMIT":    true,
		"TIME_LIMIT_90": true,
		"TIME_LIMIT_80": true,
		"TIME_LIMIT_50": true,
		"ARRAY_TASKS":   true,
	}
	for _, mailType := range mailTypes {
		if _, ok := validTypes[strings.ToUpper(mailType)]; !ok {
			return errors.NewValidationError(
				errors.ErrorCodeValidationFailed,
				"Invalid mail type: "+mailType,
				"mailType", mailType, nil,
			)
		}
	}
	return nil
}

// ValidateMailTypeValues validates mail notification types using the enum type
func (m *JobBaseManager) ValidateMailTypeValues(mailTypes []types.MailTypeValue) error {
	validTypes := map[types.MailTypeValue]bool{
		types.MailTypeBegin:             true,
		types.MailTypeEnd:               true,
		types.MailTypeFail:              true,
		types.MailTypeRequeue:           true,
		types.MailTypeTime100:           true,
		types.MailTypeTime90:            true,
		types.MailTypeTime80:            true,
		types.MailTypeTime50:            true,
		types.MailTypeStageOut:          true,
		types.MailTypeArrayTasks:        true,
		types.MailTypeInvalidDependency: true,
	}
	for _, mailType := range mailTypes {
		if _, ok := validTypes[mailType]; !ok {
			return errors.NewValidationError(
				errors.ErrorCodeValidationFailed,
				"Invalid mail type: "+string(mailType),
				"mailType", mailType, nil,
			)
		}
	}
	return nil
}

// ValidateDependencyString validates a SLURM dependency specification string
// Format examples: "afterok:123:456", "afterany:789", "singleton"
func (m *JobBaseManager) ValidateDependencyString(depStr string) error {
	if depStr == "" {
		return nil
	}
	validTypes := map[string]bool{
		"after":       true,
		"afterany":    true,
		"aftercorr":   true,
		"afternotok":  true,
		"afterok":     true,
		"expand":      true,
		"singleton":   true,
		"burstbuffer": true,
	}
	// Split by comma for multiple dependencies
	parts := strings.Split(depStr, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		// Split by colon to get type and job IDs
		colonParts := strings.SplitN(part, ":", 2)
		depType := strings.ToLower(colonParts[0])
		if _, ok := validTypes[depType]; !ok {
			return errors.NewValidationError(
				errors.ErrorCodeValidationFailed,
				"Invalid dependency type: "+colonParts[0],
				"dependency", depStr, nil,
			)
		}
	}
	return nil
}

// ApplyJobDefaults applies default values to job create request
func (m *JobBaseManager) ApplyJobDefaults(job *types.JobCreate) *types.JobCreate {
	// Apply default values (using pointers)
	if job.MinimumCPUs == nil {
		cpus := int32(1)
		job.MinimumCPUs = &cpus
	}
	if job.MinimumNodes == nil {
		nodes := int32(1)
		job.MinimumNodes = &nodes
	}
	if job.Tasks == nil {
		tasks := int32(1)
		job.Tasks = &tasks
	}
	// Initialize arrays if nil
	if job.Environment == nil {
		job.Environment = []string{}
	}
	if job.MailType == nil {
		job.MailType = []types.MailTypeValue{}
	}
	// Set default working directory if not specified
	if job.CurrentWorkingDirectory == nil {
		cwd := "/tmp"
		job.CurrentWorkingDirectory = &cwd
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
	return m.checkJobIDFilter(opts.JobIDs, getJobID(&job)) &&
		m.checkStringFilter(opts.JobNames, derefString(job.Name)) &&
		m.checkStringFilter(opts.Accounts, derefString(job.Account)) &&
		m.checkStringFilter(opts.Users, derefString(job.UserName)) &&
		m.checkStateFilter(opts.States, getJobState(&job)) &&
		m.checkStringFilter(opts.Partitions, derefString(job.Partition)) &&
		m.checkStringFilter(opts.QoS, derefString(job.QoS)) &&
		m.checkTimeRange(&job.SubmitTime, opts.StartTime, opts.EndTime)
}
func (m *JobBaseManager) checkJobIDFilter(filterIDs []int32, jobID int32) bool {
	if len(filterIDs) == 0 {
		return true
	}
	for _, id := range filterIDs {
		if jobID == id {
			return true
		}
	}
	return false
}
func (m *JobBaseManager) checkStringFilter(filters []string, value string) bool {
	if len(filters) == 0 {
		return true
	}
	for _, filter := range filters {
		if strings.EqualFold(value, filter) {
			return true
		}
	}
	return false
}
func (m *JobBaseManager) checkStateFilter(states []types.JobState, jobState types.JobState) bool {
	if len(states) == 0 {
		return true
	}
	for _, state := range states {
		if jobState == state {
			return true
		}
	}
	return false
}
func (m *JobBaseManager) checkTimeRange(submitTime, startTime, endTime *time.Time) bool {
	if startTime != nil && submitTime.Before(*startTime) {
		return false
	}
	if endTime != nil && submitTime.After(*endTime) {
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
	if req.JobId <= 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Job ID must be positive",
			"request.JobId", req.JobId, nil,
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
	if req.JobId <= 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Job ID must be positive",
			"request.JobId", req.JobId, nil,
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
	if req.JobId <= 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Job ID must be positive",
			"request.JobId", req.JobId, nil,
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
