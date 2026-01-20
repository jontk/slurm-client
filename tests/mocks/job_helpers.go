// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package mocks

import (
	v0_0_40 "github.com/jontk/slurm-client/internal/api/v0_0_40"
	v0_0_42 "github.com/jontk/slurm-client/internal/api/v0_0_42"
	v0_0_43 "github.com/jontk/slurm-client/internal/api/v0_0_43"
	v0_0_44 "github.com/jontk/slurm-client/internal/api/v0_0_44"
)

// extractCPUs extracts the CPU count from any job version
func extractCPUs(jobInterface interface{}) int64 {
	defaultCPUs := int64(4)

	switch job := jobInterface.(type) {
	case *v0_0_40.V0040JobInfo:
		if job.Cpus != nil && job.Cpus.Number != nil {
			return int64(*job.Cpus.Number)
		}
	case *v0_0_42.V0042JobInfo:
		if job.Cpus != nil && job.Cpus.Number != nil {
			return int64(*job.Cpus.Number)
		}
	case *v0_0_43.V0043JobInfo:
		if job.Cpus != nil && job.Cpus.Number != nil {
			return int64(*job.Cpus.Number)
		}
	case *v0_0_44.V0044JobInfo:
		if job.Cpus != nil && job.Cpus.Number != nil {
			return int64(*job.Cpus.Number)
		}
	}

	return defaultCPUs
}

// extractMemory extracts the memory per node from any job version
func extractMemory(jobInterface interface{}) int64 {
	defaultMemory := int64(4 * 1024 * 1024 * 1024) // 4GB

	switch job := jobInterface.(type) {
	case *v0_0_40.V0040JobInfo:
		if job.MemoryPerNode != nil && job.MemoryPerNode.Number != nil {
			return int64(*job.MemoryPerNode.Number)
		}
	case *v0_0_42.V0042JobInfo:
		if job.MemoryPerNode != nil && job.MemoryPerNode.Number != nil {
			return int64(*job.MemoryPerNode.Number)
		}
	case *v0_0_43.V0043JobInfo:
		if job.MemoryPerNode != nil && job.MemoryPerNode.Number != nil {
			return int64(*job.MemoryPerNode.Number)
		}
	case *v0_0_44.V0044JobInfo:
		if job.MemoryPerNode != nil && job.MemoryPerNode.Number != nil {
			return int64(*job.MemoryPerNode.Number)
		}
	}

	return defaultMemory
}

// extractJobState extracts the job state from any job version
func extractJobState(jobInterface interface{}) string {
	defaultState := "UNKNOWN"

	switch job := jobInterface.(type) {
	case *v0_0_40.V0040JobInfo:
		if job.JobState != nil && len(*job.JobState) > 0 {
			return string((*job.JobState)[0])
		}
	case *v0_0_42.V0042JobInfo:
		if job.JobState != nil && len(*job.JobState) > 0 {
			return string((*job.JobState)[0])
		}
	case *v0_0_43.V0043JobInfo:
		if job.JobState != nil && len(*job.JobState) > 0 {
			return string((*job.JobState)[0])
		}
	case *v0_0_44.V0044JobInfo:
		if job.JobState != nil && len(*job.JobState) > 0 {
			return string((*job.JobState)[0])
		}
	}

	return defaultState
}
