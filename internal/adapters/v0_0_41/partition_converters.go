// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_41

import (
	"fmt"
	"strings"

	"github.com/jontk/slurm-client/internal/common/types"
)

// convertAPIPartitionToCommon converts a v0.0.41 API Partition to common Partition type
func (a *PartitionAdapter) convertAPIPartitionToCommon(apiPartition interface{}) (*types.Partition, error) {
	// Use map interface for handling anonymous structs in v0.0.41
	partitionData, ok := apiPartition.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected partition data type: %T", apiPartition)
	}

	partition := &types.Partition{}

	// Basic fields - using safe type assertions
	if v, ok := partitionData["name"]; ok {
		if name, ok := v.(string); ok {
			partition.Name = name
		}
	}
	if v, ok := partitionData["nodes"]; ok {
		if nodes, ok := v.(string); ok {
			partition.Nodes = nodes
		}
	}
	if v, ok := partitionData["state"]; ok {
		if states, ok := v.([]interface{}); ok && len(states) > 0 {
			if state, ok := states[0].(string); ok {
				partition.State = types.PartitionState(state)
			}
		}
	}

	// Resource limits
	if v, ok := partitionData["max_cpus_per_node"]; ok {
		if maxCpus, ok := v.(float64); ok {
			partition.MaxCPUsPerNode = int32(maxCpus)
		}
	}
	if v, ok := partitionData["max_memory_per_node"]; ok {
		if maxMem, ok := v.(float64); ok {
			partition.MaxMemPerNode = int64(maxMem)
		}
	}
	if v, ok := partitionData["max_nodes"]; ok {
		if maxNodes, ok := v.(float64); ok {
			partition.MaxNodes = int32(maxNodes)
		}
	}
	if v, ok := partitionData["min_nodes"]; ok {
		if minNodes, ok := v.(float64); ok {
			partition.MinNodes = int32(minNodes)
		}
	}

	// Time limits - handle structured objects
	if v, ok := partitionData["max_time"]; ok {
		if timeStruct, ok := v.(map[string]interface{}); ok {
			if number, ok := timeStruct["number"].(float64); ok {
				partition.MaxTime = int32(number)
			}
		} else if maxTime, ok := v.(float64); ok {
			partition.MaxTime = int32(maxTime)
		}
	}
	if v, ok := partitionData["default_time"]; ok {
		if timeStruct, ok := v.(map[string]interface{}); ok {
			if number, ok := timeStruct["number"].(float64); ok {
				partition.DefaultTime = int32(number)
			}
		} else if defaultTime, ok := v.(float64); ok {
			partition.DefaultTime = int32(defaultTime)
		}
	}

	// Priority and preemption
	if v, ok := partitionData["priority_job_factor"]; ok {
		if priority, ok := v.(float64); ok {
			partition.PriorityJobFactor = int32(priority)
		}
	}
	if v, ok := partitionData["priority_tier"]; ok {
		if tier, ok := v.(float64); ok {
			partition.PriorityTier = int32(tier)
		}
	}
	if v, ok := partitionData["preempt_mode"]; ok {
		if modes, ok := v.([]interface{}); ok && len(modes) > 0 {
			modeStrings := make([]string, len(modes))
			for i, mode := range modes {
				if modeStr, ok := mode.(string); ok {
					modeStrings[i] = modeStr
				}
			}
			partition.PreemptMode = modeStrings
		}
	}

	// Allowed groups and accounts
	if v, ok := partitionData["allowed_groups"]; ok {
		if groups, ok := v.(string); ok {
			partition.AllowGroups = strings.Split(groups, ",")
		}
	}
	if v, ok := partitionData["allowed_accounts"]; ok {
		if accounts, ok := v.(string); ok {
			partition.AllowAccounts = strings.Split(accounts, ",")
		}
	}
	if v, ok := partitionData["allowed_qos"]; ok {
		if qos, ok := v.(string); ok {
			partition.AllowQoS = strings.Split(qos, ",")
		}
	}

	// Boolean flags
	// Default field doesn't exist in common Partition type
	// Skip default flag conversion
	if v, ok := partitionData["root_only"]; ok {
		if rootOnly, ok := v.(bool); ok {
			partition.RootOnly = rootOnly
		}
	}
	// Shared field doesn't exist in common Partition type
	// Skip shared flag conversion

	// Other fields
	if v, ok := partitionData["grace_time"]; ok {
		if graceTime, ok := v.(float64); ok {
			partition.GraceTime = int32(graceTime)
		}
	}
	if v, ok := partitionData["over_time_limit"]; ok {
		if overTime, ok := v.(float64); ok {
			partition.OverTimeLimit = int32(overTime)
		}
	}

	// CPU allocation ratio - CPUBind field doesn't exist in common Partition type
	// Skip cpu_bind conversion

	// Features - RequiredFeatures field doesn't exist in common Partition type
	// Skip required_features conversion

	return partition, nil
}
