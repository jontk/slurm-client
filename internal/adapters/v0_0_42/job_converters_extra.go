// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_42

import (
	types "github.com/jontk/slurm-client/api"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_42"
)

// convertCommonJobCreateToAPI converts JobCreate to the API request body type
func (a *JobAdapter) convertCommonJobCreateToAPI(input *types.JobCreate) api.SlurmV0042PostJobSubmitJSONRequestBody {
	if input == nil {
		return api.SlurmV0042PostJobSubmitJSONRequestBody{}
	}
	jobDesc := &api.V0042JobDescMsg{}
	// Basic job metadata
	if input.Name != nil {
		jobDesc.Name = input.Name
	}
	if input.Account != nil {
		jobDesc.Account = input.Account
	}
	if input.Partition != nil {
		jobDesc.Partition = input.Partition
	}
	if input.Script != nil {
		jobDesc.Script = input.Script
	}
	if input.Comment != nil {
		jobDesc.Comment = input.Comment
	}
	if input.QoS != nil {
		jobDesc.Qos = input.QoS
	}
	// Time limits
	if input.TimeLimit != nil {
		jobDesc.TimeLimit = &api.V0042Uint32NoValStruct{
			Set:    ptrBool(true),
			Number: ptrInt32(int32(*input.TimeLimit)),
		}
	}
	// Node requirements
	if input.MinimumNodes != nil {
		jobDesc.MinimumNodes = input.MinimumNodes
	}
	if input.MaximumNodes != nil {
		jobDesc.MaximumNodes = input.MaximumNodes
	}
	if input.Nodes != nil {
		jobDesc.Nodes = input.Nodes
	}
	// CPU requirements
	if input.MinimumCPUs != nil {
		jobDesc.MinimumCpus = input.MinimumCPUs
	}
	if input.CPUsPerTask != nil {
		jobDesc.CpusPerTask = input.CPUsPerTask
	}
	// Task counts
	if input.Tasks != nil {
		jobDesc.Tasks = input.Tasks
	}
	if input.TasksPerNode != nil {
		jobDesc.TasksPerNode = input.TasksPerNode
	}
	// Memory requirements
	if input.MemoryPerNode != nil {
		jobDesc.MemoryPerNode = &api.V0042Uint64NoValStruct{
			Set:    ptrBool(true),
			Number: ptrInt64(int64(*input.MemoryPerNode)),
		}
	}
	if input.MemoryPerCPU != nil {
		jobDesc.MemoryPerCpu = &api.V0042Uint64NoValStruct{
			Set:    ptrBool(true),
			Number: ptrInt64(int64(*input.MemoryPerCPU)),
		}
	}
	// Environment variables
	if len(input.Environment) > 0 {
		envArray := api.V0042StringArray(input.Environment)
		jobDesc.Environment = &envArray
	}
	// Output/error files
	if input.StandardOutput != nil {
		jobDesc.StandardOutput = input.StandardOutput
	}
	if input.StandardError != nil {
		jobDesc.StandardError = input.StandardError
	}
	if input.StandardInput != nil {
		jobDesc.StandardInput = input.StandardInput
	}
	// Working directory
	if input.CurrentWorkingDirectory != nil {
		jobDesc.CurrentWorkingDirectory = input.CurrentWorkingDirectory
	}
	// Constraints and dependencies
	if input.Constraints != nil {
		jobDesc.Constraints = input.Constraints
	}
	if input.Dependency != nil {
		jobDesc.Dependency = input.Dependency
	}
	// Job array
	if input.Array != nil {
		jobDesc.Array = input.Array
	}
	// Hold status
	if input.Hold != nil {
		jobDesc.Hold = input.Hold
	}
	// Reservation
	if input.Reservation != nil {
		jobDesc.Reservation = input.Reservation
	}
	return api.SlurmV0042PostJobSubmitJSONRequestBody{
		Job:    jobDesc,
		Script: jobDesc.Script,
	}
}
