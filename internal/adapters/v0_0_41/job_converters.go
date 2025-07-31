package v0_0_41

import (
	"fmt"
	"strings"
	"time"

	"github.com/jontk/slurm-client/internal/common/types"
)

// convertAPIJobToCommon converts a v0.0.41 API Job to common Job type
func (a *JobAdapter) convertAPIJobToCommon(apiJob interface{}) (*types.Job, error) {
	// Use map interface for handling anonymous structs in v0.0.41
	jobData, ok := apiJob.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected job data type: %T", apiJob)
	}

	job := &types.Job{}

	// Basic fields - using safe type assertions
	if v, ok := jobData["job_id"]; ok {
		if jobID, ok := v.(float64); ok {
			job.JobID = int32(jobID)
		}
	}
	if v, ok := jobData["name"]; ok {
		if name, ok := v.(string); ok {
			job.Name = name
		}
	}
	if v, ok := jobData["user_name"]; ok {
		if userName, ok := v.(string); ok {
			job.UserName = userName
		}
	}
	if v, ok := jobData["account"]; ok {
		if account, ok := v.(string); ok {
			job.Account = account
		}
	}
	if v, ok := jobData["partition"]; ok {
		if partition, ok := v.(string); ok {
			job.Partition = partition
		}
	}
	if v, ok := jobData["qos"]; ok {
		if qos, ok := v.(string); ok {
			job.QoS = qos
		}
	}

	// Job state
	if v, ok := jobData["job_state"]; ok {
		if states, ok := v.([]interface{}); ok && len(states) > 0 {
			if state, ok := states[0].(string); ok {
				job.State = types.JobState(state)
			}
		}
	}
	if v, ok := jobData["state_reason"]; ok {
		if reason, ok := v.(string); ok {
			job.StateReason = reason
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
				startTime := time.Unix(int64(number), 0)
				job.StartTime = &startTime
			}
		}
	}
	if v, ok := jobData["end_time"]; ok {
		if timeStruct, ok := v.(map[string]interface{}); ok {
			if number, ok := timeStruct["number"].(float64); ok {
				endTime := time.Unix(int64(number), 0)
				job.EndTime = &endTime
			}
		}
	}

	// Resource requirements
	if v, ok := jobData["node_count"]; ok {
		if nodeStruct, ok := v.(map[string]interface{}); ok {
			if number, ok := nodeStruct["number"].(float64); ok {
				job.Nodes = int32(number)
			}
		}
	}
	if v, ok := jobData["cpus"]; ok {
		if cpuStruct, ok := v.(map[string]interface{}); ok {
			if number, ok := cpuStruct["number"].(float64); ok {
				job.CPUs = int32(number)
			}
		}
	}

	// Time limit
	if v, ok := jobData["time_limit"]; ok {
		if timeStruct, ok := v.(map[string]interface{}); ok {
			if number, ok := timeStruct["number"].(float64); ok {
				job.TimeLimit = int32(number)
			}
		}
	}

	// Priority
	if v, ok := jobData["priority"]; ok {
		if priorityStruct, ok := v.(map[string]interface{}); ok {
			if number, ok := priorityStruct["number"].(float64); ok {
				job.Priority = int32(number)
			}
		}
	}

	// Node information
	if v, ok := jobData["nodes"]; ok {
		if nodes, ok := v.(string); ok {
			job.NodeList = nodes
		}
	}

	// Standard I/O
	if v, ok := jobData["standard_input"]; ok {
		if stdIn, ok := v.(string); ok {
			job.StandardInput = stdIn
		}
	}
	if v, ok := jobData["standard_output"]; ok {
		if stdOut, ok := v.(string); ok {
			job.StandardOutput = stdOut
		}
	}
	if v, ok := jobData["standard_error"]; ok {
		if stdErr, ok := v.(string); ok {
			job.StandardError = stdErr
		}
	}

	// Working directory
	if v, ok := jobData["current_working_directory"]; ok {
		if workDir, ok := v.(string); ok {
			job.WorkingDirectory = workDir
		}
	}

	// Environment - convert from []string to map[string]string
	if v, ok := jobData["environment"]; ok {
		if env, ok := v.([]interface{}); ok {
			envMap := make(map[string]string)
			for _, e := range env {
				if envStr, ok := e.(string); ok {
					parts := strings.SplitN(envStr, "=", 2)
					if len(parts) == 2 {
						envMap[parts[0]] = parts[1]
					}
				}
			}
			job.Environment = envMap
		}
	}

	// Comment
	if v, ok := jobData["comment"]; ok {
		if comment, ok := v.(string); ok {
			job.Comment = comment
		}
	}

	return job, nil
}

// convertCommonToAPIJobUpdate converts common JobUpdate to v0.0.41 API request
func (a *JobAdapter) convertCommonToAPIJobUpdate(update *types.JobUpdate) map[string]interface{} {
	req := make(map[string]interface{})

	// Set fields that can be updated
	if update.Comment != nil {
		req["comment"] = *update.Comment
	}
	if update.Priority != nil {
		req["priority"] = *update.Priority
	}
	if update.QoS != nil {
		req["qos"] = *update.QoS
	}
	if update.TimeLimit != nil {
		req["time_limit"] = *update.TimeLimit
	}

	return req
}