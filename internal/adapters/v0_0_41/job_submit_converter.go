// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_41

import (
	"encoding/json"

	types "github.com/jontk/slurm-client/api"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_41"
)

// convertCommonJobCreateToAPI converts JobCreate to the v0.0.41 API request body type.
// v0.0.41 uses inline anonymous structs which requires JSON marshaling/unmarshaling
// as a workaround for Go's strict type system.
func (a *JobAdapter) convertCommonJobCreateToAPI(input *types.JobCreate) (api.SlurmV0041PostJobSubmitJSONRequestBody, error) {
	if input == nil {
		return api.SlurmV0041PostJobSubmitJSONRequestBody{}, nil
	}

	// Build an intermediate representation that matches the API structure
	// Using a generic map allows us to bypass Go's type system limitations
	// with the anonymous structs in v0.0.41
	jobMap := make(map[string]interface{})

	// Set basic string fields
	if input.Account != nil {
		jobMap["account"] = *input.Account
	}
	if input.Name != nil {
		jobMap["name"] = *input.Name
	}
	if input.Partition != nil {
		jobMap["partition"] = *input.Partition
	}
	if input.Script != nil {
		jobMap["script"] = *input.Script
	}
	if input.Comment != nil {
		jobMap["comment"] = *input.Comment
	}
	if input.QoS != nil {
		jobMap["qos"] = *input.QoS
	}
	if input.Constraints != nil {
		jobMap["constraints"] = *input.Constraints
	}
	if input.Dependency != nil {
		jobMap["dependency"] = *input.Dependency
	}
	if input.Array != nil {
		jobMap["array"] = *input.Array
	}
	if input.Reservation != nil {
		jobMap["reservation"] = *input.Reservation
	}
	if input.StandardOutput != nil {
		jobMap["standard_output"] = *input.StandardOutput
	}
	if input.StandardError != nil {
		jobMap["standard_error"] = *input.StandardError
	}
	if input.StandardInput != nil {
		jobMap["standard_input"] = *input.StandardInput
	}
	if input.CurrentWorkingDirectory != nil {
		jobMap["current_working_directory"] = *input.CurrentWorkingDirectory
	}
	if input.Nodes != nil {
		jobMap["nodes"] = *input.Nodes
	}

	// Set numeric fields
	if input.MinimumNodes != nil {
		jobMap["minimum_nodes"] = *input.MinimumNodes
	}
	if input.MaximumNodes != nil {
		jobMap["maximum_nodes"] = *input.MaximumNodes
	}
	if input.MinimumCPUs != nil {
		jobMap["minimum_cpus"] = *input.MinimumCPUs
	}
	if input.CPUsPerTask != nil {
		jobMap["cpus_per_task"] = *input.CPUsPerTask
	}
	if input.Tasks != nil {
		jobMap["tasks"] = *input.Tasks
	}
	if input.TasksPerNode != nil {
		jobMap["tasks_per_node"] = *input.TasksPerNode
	}

	// Set boolean fields
	if input.Hold != nil {
		jobMap["hold"] = *input.Hold
	}

	// Set complex fields with number wrappers (v0.0.41 uses set/number/infinite structs)
	if input.TimeLimit != nil {
		jobMap["time_limit"] = map[string]interface{}{
			"set":    true,
			"number": int32(*input.TimeLimit),
		}
	}

	if input.MemoryPerNode != nil {
		jobMap["memory_per_node"] = map[string]interface{}{
			"set":    true,
			"number": int64(*input.MemoryPerNode),
		}
	}

	if input.MemoryPerCPU != nil {
		jobMap["memory_per_cpu"] = map[string]interface{}{
			"set":    true,
			"number": int64(*input.MemoryPerCPU),
		}
	}

	// Set environment if present
	if len(input.Environment) > 0 {
		jobMap["environment"] = input.Environment
	}

	// Build the request body structure
	bodyMap := map[string]interface{}{
		"job": jobMap,
	}
	if input.Script != nil {
		bodyMap["script"] = *input.Script
	}

	// Marshal to JSON and unmarshal to the API type
	jsonBytes, err := json.Marshal(bodyMap)
	if err != nil {
		return api.SlurmV0041PostJobSubmitJSONRequestBody{}, err
	}

	var body api.SlurmV0041PostJobSubmitJSONRequestBody
	if err := json.Unmarshal(jsonBytes, &body); err != nil {
		return api.SlurmV0041PostJobSubmitJSONRequestBody{}, err
	}

	return body, nil
}
