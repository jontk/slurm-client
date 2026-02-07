// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_41

import (
	"encoding/json"
	"fmt"
	"time"

	types "github.com/jontk/slurm-client/api"
)

// convertAPIJobToCommon converts a v0.0.41 API Job to common Job type
func (a *JobAdapter) convertAPIJobToCommon(apiJob interface{}) (*types.Job, error) {
	// v0.0.41 uses anonymous structs, so we need to handle reflection
	// Try to convert to map first (in case it's already been marshalled/unmarshalled)
	jobData, ok := apiJob.(map[string]interface{})
	if !ok {
		// If not a map, try to convert via JSON marshalling
		jsonBytes, err := json.Marshal(apiJob)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal job data: %w", err)
		}
		if err := json.Unmarshal(jsonBytes, &jobData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal job data to map: %w", err)
		}
	}
	job := &types.Job{}
	// Basic fields - using safe type assertions
	if v, ok := jobData["job_id"]; ok {
		if jobID, ok := v.(float64); ok {
			id := int32(jobID)
			job.JobID = &id
		}
	}
	if v, ok := jobData["name"]; ok {
		if name, ok := v.(string); ok {
			n := name
			job.Name = &n
		}
	}
	if v, ok := jobData["user_name"]; ok {
		if userName, ok := v.(string); ok {
			u := userName
			job.UserName = &u
		}
	}
	if v, ok := jobData["account"]; ok {
		if account, ok := v.(string); ok {
			a := account
			job.Account = &a
		}
	}
	if v, ok := jobData["partition"]; ok {
		if partition, ok := v.(string); ok {
			p := partition
			job.Partition = &p
		}
	}
	if v, ok := jobData["qos"]; ok {
		if qos, ok := v.(string); ok {
			q := qos
			job.QoS = &q
		}
	}
	// Job state
	if v, ok := jobData["job_state"]; ok {
		if states, ok := v.([]interface{}); ok {
			jobStates := make([]types.JobState, 0, len(states))
			for _, s := range states {
				if state, ok := s.(string); ok {
					jobStates = append(jobStates, types.JobState(state))
				}
			}
			job.JobState = jobStates
		}
	}
	if v, ok := jobData["state_reason"]; ok {
		if reason, ok := v.(string); ok {
			r := reason
			job.StateReason = &r
		}
	}
	// Time fields - handle both direct numbers and structured time objects
	if v, ok := jobData["submit_time"]; ok {
		if timeStruct, ok := v.(map[string]interface{}); ok {
			if number, ok := timeStruct["number"].(float64); ok {
				job.SubmitTime = time.Unix(int64(number), 0)
			}
		} else if timestamp, ok := v.(float64); ok {
			job.SubmitTime = time.Unix(int64(timestamp), 0)
		}
	}
	if v, ok := jobData["start_time"]; ok {
		if timeStruct, ok := v.(map[string]interface{}); ok {
			if number, ok := timeStruct["number"].(float64); ok {
				job.StartTime = time.Unix(int64(number), 0)
			}
		}
	}
	if v, ok := jobData["end_time"]; ok {
		if timeStruct, ok := v.(map[string]interface{}); ok {
			if number, ok := timeStruct["number"].(float64); ok {
				job.EndTime = time.Unix(int64(number), 0)
			}
		}
	}
	// Resource requirements
	if v, ok := jobData["node_count"]; ok {
		if nodeStruct, ok := v.(map[string]interface{}); ok {
			if number, ok := nodeStruct["number"].(float64); ok {
				nc := uint32(number)
				job.NodeCount = &nc
			}
		}
	}
	if v, ok := jobData["cpus"]; ok {
		if cpuStruct, ok := v.(map[string]interface{}); ok {
			if number, ok := cpuStruct["number"].(float64); ok {
				cpus := uint32(number)
				job.CPUs = &cpus
			}
		}
	}
	// Time limit
	if v, ok := jobData["time_limit"]; ok {
		if timeStruct, ok := v.(map[string]interface{}); ok {
			if number, ok := timeStruct["number"].(float64); ok {
				tl := uint32(number)
				job.TimeLimit = &tl
			}
		}
	}
	// Priority
	if v, ok := jobData["priority"]; ok {
		if priorityStruct, ok := v.(map[string]interface{}); ok {
			if number, ok := priorityStruct["number"].(float64); ok {
				pri := uint32(number)
				job.Priority = &pri
			}
		}
	}
	// Node information
	if v, ok := jobData["nodes"]; ok {
		if nodes, ok := v.(string); ok {
			n := nodes
			job.Nodes = &n
		}
	}
	// Standard I/O
	if v, ok := jobData["standard_input"]; ok {
		if stdIn, ok := v.(string); ok {
			si := stdIn
			job.StandardInput = &si
		}
	}
	if v, ok := jobData["standard_output"]; ok {
		if stdOut, ok := v.(string); ok {
			so := stdOut
			job.StandardOutput = &so
		}
	}
	if v, ok := jobData["standard_error"]; ok {
		if stdErr, ok := v.(string); ok {
			se := stdErr
			job.StandardError = &se
		}
	}
	// Working directory
	if v, ok := jobData["current_working_directory"]; ok {
		if workDir, ok := v.(string); ok {
			cwd := workDir
			job.CurrentWorkingDirectory = &cwd
		}
	}
	// Comment
	if v, ok := jobData["comment"]; ok {
		if comment, ok := v.(string); ok {
			c := comment
			job.Comment = &c
		}
	}
	// Flags - convert using helper function
	if v, ok := jobData["flags"]; ok {
		if flags, ok := v.([]interface{}); ok {
			jobFlags := make([]types.FlagsValue, 0, len(flags))
			for _, f := range flags {
				if flag, ok := f.(string); ok {
					jobFlags = append(jobFlags, types.FlagsValue(flag))
				}
			}
			job.Flags = jobFlags
		}
	}
	// MailType - convert using helper function
	if v, ok := jobData["mail_type"]; ok {
		if mailTypes, ok := v.([]interface{}); ok {
			jobMailTypes := make([]types.MailTypeValue, 0, len(mailTypes))
			for _, mt := range mailTypes {
				if mailType, ok := mt.(string); ok {
					jobMailTypes = append(jobMailTypes, types.MailTypeValue(mailType))
				}
			}
			job.MailType = jobMailTypes
		}
	}
	// Profile - convert using helper function
	if v, ok := jobData["profile"]; ok {
		if profiles, ok := v.([]interface{}); ok {
			jobProfiles := make([]types.ProfileValue, 0, len(profiles))
			for _, p := range profiles {
				if profile, ok := p.(string); ok {
					jobProfiles = append(jobProfiles, types.ProfileValue(profile))
				}
			}
			job.Profile = jobProfiles
		}
	}
	// Shared - convert using helper function
	if v, ok := jobData["shared"]; ok {
		if shareds, ok := v.([]interface{}); ok {
			jobShareds := make([]types.SharedValue, 0, len(shareds))
			for _, s := range shareds {
				if shared, ok := s.(string); ok {
					jobShareds = append(jobShareds, types.SharedValue(shared))
				}
			}
			job.Shared = jobShareds
		}
	}
	// ExitCode - convert using helper function
	if v, ok := jobData["exit_code"]; ok {
		if exitCodeData, ok := v.(map[string]interface{}); ok {
			job.ExitCode = convertExitCodeFromMap(exitCodeData)
		}
	}
	// DerivedExitCode - convert using helper function
	if v, ok := jobData["derived_exit_code"]; ok {
		if derivedExitCodeData, ok := v.(map[string]interface{}); ok {
			job.DerivedExitCode = convertExitCodeFromMap(derivedExitCodeData)
		}
	}
	// JobResources - convert complex nested structure
	if v, ok := jobData["job_resources"]; ok {
		if jobResourcesData, ok := v.(map[string]interface{}); ok {
			job.JobResources = convertJobResourcesFromMap(jobResourcesData)
		}
	}
	return job, nil
}

// convertJobResourcesFromMap converts job resources data from a map to common JobResources type.
// This handles the anonymous struct pattern used in v0_0_41 API responses.
func convertJobResourcesFromMap(data map[string]interface{}) *types.JobResources {
	if data == nil {
		return nil
	}
	result := &types.JobResources{}

	// Convert CPUs
	if cpus, ok := data["cpus"].(float64); ok {
		result.CPUs = int32(cpus)
	}

	// Convert SelectType
	if selectTypeData, ok := data["select_type"].([]interface{}); ok {
		selectType := make([]types.SelectTypeValue, 0, len(selectTypeData))
		for _, st := range selectTypeData {
			if stStr, ok := st.(string); ok {
				selectType = append(selectType, types.SelectTypeValue(stStr))
			}
		}
		result.SelectType = selectType
	}

	// Convert ThreadsPerCore (NoValStruct)
	if tpcData, ok := data["threads_per_core"].(map[string]interface{}); ok {
		if set, ok := tpcData["set"].(bool); ok && set {
			if number, ok := tpcData["number"].(float64); ok {
				result.ThreadsPerCore = uint16(number)
			}
		}
	}

	// Convert Nodes
	if nodesData, ok := data["nodes"].(map[string]interface{}); ok {
		result.Nodes = convertJobResourcesNodesFromMap(nodesData)
	}

	return result
}

// convertJobResourcesNodesFromMap converts nodes data from a map to common JobResourcesNodes.
func convertJobResourcesNodesFromMap(data map[string]interface{}) *types.JobResourcesNodes {
	if data == nil {
		return nil
	}
	result := &types.JobResourcesNodes{}

	// Count
	if count, ok := data["count"].(float64); ok {
		c := int32(count)
		result.Count = &c
	}

	// List
	if list, ok := data["list"].(string); ok {
		result.List = &list
	}

	// Whole
	if whole, ok := data["whole"].(bool); ok {
		result.Whole = &whole
	}

	// SelectType
	if selectTypeData, ok := data["select_type"].([]interface{}); ok {
		selectType := make([]types.JobResourcesNodesSelectTypeValue, 0, len(selectTypeData))
		for _, st := range selectTypeData {
			if stStr, ok := st.(string); ok {
				selectType = append(selectType, types.JobResourcesNodesSelectTypeValue(stStr))
			}
		}
		result.SelectType = selectType
	}

	// Allocation (slice of nodes)
	if allocationData, ok := data["allocation"].([]interface{}); ok {
		allocation := make([]types.JobResNode, 0, len(allocationData))
		for _, nodeInterface := range allocationData {
			if nodeMap, ok := nodeInterface.(map[string]interface{}); ok {
				node := convertJobResNodeFromMap(nodeMap)
				allocation = append(allocation, node)
			}
		}
		result.Allocation = allocation
	}

	return result
}

// convertJobResNodeFromMap converts a node from map to JobResNode.
func convertJobResNodeFromMap(data map[string]interface{}) types.JobResNode {
	result := types.JobResNode{}

	// Name
	if name, ok := data["name"].(string); ok {
		result.Name = name
	}

	// CPUs
	if cpusData, ok := data["cpus"].(map[string]interface{}); ok {
		result.CPUs = &types.JobResNodeCPUs{}
		if count, ok := cpusData["count"].(float64); ok {
			c := int32(count)
			result.CPUs.Count = &c
		}
		if used, ok := cpusData["used"].(float64); ok {
			u := int32(used)
			result.CPUs.Used = &u
		}
	}

	// Memory
	if memoryData, ok := data["memory"].(map[string]interface{}); ok {
		result.Memory = &types.JobResNodeMemory{}
		if allocated, ok := memoryData["allocated"].(float64); ok {
			a := int64(allocated)
			result.Memory.Allocated = &a
		}
	}

	// Sockets
	if socketsData, ok := data["sockets"].([]interface{}); ok {
		sockets := make([]types.JobResSocket, 0, len(socketsData))
		for _, socketInterface := range socketsData {
			if socketMap, ok := socketInterface.(map[string]interface{}); ok {
				socket := convertJobResSocketFromMap(socketMap)
				sockets = append(sockets, socket)
			}
		}
		result.Sockets = sockets
	}

	return result
}

// convertJobResSocketFromMap converts a socket from map to JobResSocket.
func convertJobResSocketFromMap(data map[string]interface{}) types.JobResSocket {
	result := types.JobResSocket{}

	// Index
	if index, ok := data["index"].(float64); ok {
		result.Index = int32(index)
	}

	// Cores
	if coresData, ok := data["cores"].([]interface{}); ok {
		cores := make([]types.JobResCore, 0, len(coresData))
		for _, coreInterface := range coresData {
			if coreMap, ok := coreInterface.(map[string]interface{}); ok {
				core := convertJobResCoreFromMap(coreMap)
				cores = append(cores, core)
			}
		}
		result.Cores = cores
	}

	return result
}

// convertJobResCoreFromMap converts a core from map to JobResCore.
func convertJobResCoreFromMap(data map[string]interface{}) types.JobResCore {
	result := types.JobResCore{}

	// Index
	if index, ok := data["index"].(float64); ok {
		result.Index = int32(index)
	}

	// Status
	if statusData, ok := data["status"].([]interface{}); ok {
		status := make([]types.JobResCoreStatusValue, 0, len(statusData))
		for _, s := range statusData {
			if statusStr, ok := s.(string); ok {
				status = append(status, types.JobResCoreStatusValue(statusStr))
			}
		}
		result.Status = status
	}

	return result
}

// convertExitCodeFromMap converts exit code data from a map to common ExitCode type.
// This handles the anonymous struct pattern used in v0_0_41 API responses.
func convertExitCodeFromMap(data map[string]interface{}) *types.ExitCode {
	if data == nil {
		return nil
	}
	result := &types.ExitCode{}
	// Convert return code
	if returnCode, ok := data["return_code"].(map[string]interface{}); ok {
		if number, ok := returnCode["number"].(float64); ok {
			rc := uint32(number)
			result.ReturnCode = &rc
		}
	}
	// Convert signal
	if signal, ok := data["signal"].(map[string]interface{}); ok {
		result.Signal = &types.ExitCodeSignal{}
		if id, ok := signal["id"].(map[string]interface{}); ok {
			if number, ok := id["number"].(float64); ok {
				sigID := uint16(number)
				result.Signal.ID = &sigID
			}
		}
		if name, ok := signal["name"].(string); ok {
			result.Signal.Name = &name
		}
	}
	// Convert status
	if status, ok := data["status"].([]interface{}); ok {
		for _, s := range status {
			if statusStr, ok := s.(string); ok {
				result.Status = append(result.Status, types.StatusValue(statusStr))
			}
		}
	}
	return result
}
