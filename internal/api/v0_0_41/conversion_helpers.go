// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_41

import (
	"fmt"

	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
)

// convertJobSubmissionToAPI converts interfaces.JobSubmission to V0041JobDescMsg
func convertJobSubmissionToAPI(job *interfaces.JobSubmission) (*V0041JobDescMsg, error) {
	if job == nil {
		return nil, errors.NewClientError(errors.ErrorCodeInvalidRequest, "Job submission cannot be nil")
	}

	jobDesc := &V0041JobDescMsg{}

	// Basic fields
	if job.Name != "" {
		jobDesc.Name = &job.Name
	}

	if job.Script != "" {
		jobDesc.Script = &job.Script
	}

	if job.Partition != "" {
		jobDesc.Partition = &job.Partition
	}

	if job.WorkingDir != "" {
		jobDesc.CurrentWorkingDirectory = &job.WorkingDir
	} else {
		// SLURM requires a working directory - default to /tmp if not specified
		defaultWorkDir := "/tmp"
		jobDesc.CurrentWorkingDirectory = &defaultWorkDir
	}

	// Resource requirements
	if job.CPUs > 0 {
		cpus := int32(job.CPUs)
		jobDesc.CpusPerTask = &cpus
	}

	if job.Memory > 0 {
		// Convert bytes to MB (Slurm expects MB)
		memoryMB := int64(job.Memory / (1024 * 1024))
		set := true
		jobDesc.MemoryPerNode = &struct {
			Infinite *bool  `json:"infinite,omitempty"`
			Number   *int64 `json:"number,omitempty"`
			Set      *bool  `json:"set,omitempty"`
		}{
			Number: &memoryMB,
			Set:    &set,
		}
	}

	if job.TimeLimit > 0 {
		timeLimit := int32(job.TimeLimit)
		set := true
		jobDesc.TimeLimit = &struct {
			Infinite *bool  `json:"infinite,omitempty"`
			Number   *int32 `json:"number,omitempty"`
			Set      *bool  `json:"set,omitempty"`
		}{
			Number: &timeLimit,
			Set:    &set,
		}
	}

	if job.Nodes > 0 {
		nodes := int32(job.Nodes)
		jobDesc.MinimumNodes = &nodes
	}

	if job.Priority > 0 {
		priority := int32(job.Priority)
		set := true
		jobDesc.Priority = &struct {
			Infinite *bool  `json:"infinite,omitempty"`
			Number   *int32 `json:"number,omitempty"`
			Set      *bool  `json:"set,omitempty"`
		}{
			Number: &priority,
			Set:    &set,
		}
	}

	// Environment variables
	// Always provide at least minimal environment to avoid SLURM write errors
	envVars := make([]string, 0)
	
	// Add default PATH if not provided
	hasPath := false
	for key := range job.Environment {
		if key == "PATH" {
			hasPath = true
			break
		}
	}
	
	if !hasPath {
		envVars = append(envVars, "PATH=/usr/bin:/bin")
	}
	
	// Add user-provided environment
	for key, value := range job.Environment {
		envVars = append(envVars, fmt.Sprintf("%s=%s", key, value))
	}
	
	jobDesc.Environment = &envVars

	// Args
	if len(job.Args) > 0 {
		jobDesc.Argv = &job.Args
	}

	return jobDesc, nil
}
