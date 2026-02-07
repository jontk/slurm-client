// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_41

import (
	"fmt"

	types "github.com/jontk/slurm-client/api"
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
			n := name
			partition.Name = &n
		}
	}
	// Nodes - nested structure
	if v, ok := partitionData["nodes"]; ok {
		if nodesStr, ok := v.(string); ok {
			partition.Nodes = &types.PartitionNodes{
				Configured: &nodesStr,
			}
		} else if nodesMap, ok := v.(map[string]interface{}); ok {
			partitionNodes := &types.PartitionNodes{}
			if configured, ok := nodesMap["configured"].(string); ok {
				partitionNodes.Configured = &configured
			}
			if total, ok := nodesMap["total"].(float64); ok {
				t := int32(total)
				partitionNodes.Total = &t
			}
			partition.Nodes = partitionNodes
		}
	}
	// State - nested in partition field
	if v, ok := partitionData["state"]; ok {
		if states, ok := v.([]interface{}); ok {
			stateValues := make([]types.StateValue, 0, len(states))
			for _, s := range states {
				if state, ok := s.(string); ok {
					stateValues = append(stateValues, types.StateValue(state))
				}
			}
			partition.Partition = &types.PartitionPartition{
				State: stateValues,
			}
		}
	}
	// Maximums - nested structure
	maximums := &types.PartitionMaximums{}
	hasMaximums := false
	if v, ok := partitionData["max_cpus_per_node"]; ok {
		if maxCpus, ok := v.(float64); ok {
			mc := uint32(maxCpus)
			maximums.CPUsPerNode = &mc
			hasMaximums = true
		}
	}
	if v, ok := partitionData["max_memory_per_node"]; ok {
		if maxMem, ok := v.(float64); ok {
			mm := int64(maxMem)
			maximums.MemoryPerCPU = &mm
			hasMaximums = true
		}
	}
	if v, ok := partitionData["max_nodes"]; ok {
		if maxNodes, ok := v.(float64); ok {
			mn := uint32(maxNodes)
			maximums.Nodes = &mn
			hasMaximums = true
		}
	}
	if v, ok := partitionData["max_time"]; ok {
		if timeStruct, ok := v.(map[string]interface{}); ok {
			if number, ok := timeStruct["number"].(float64); ok {
				mt := uint32(number)
				maximums.Time = &mt
				hasMaximums = true
			}
		} else if maxTime, ok := v.(float64); ok {
			mt := uint32(maxTime)
			maximums.Time = &mt
			hasMaximums = true
		}
	}
	if v, ok := partitionData["over_time_limit"]; ok {
		if overTime, ok := v.(float64); ok {
			ot := uint16(overTime)
			maximums.OverTimeLimit = &ot
			hasMaximums = true
		}
	}
	if hasMaximums {
		partition.Maximums = maximums
	}
	// Minimums - nested structure
	if v, ok := partitionData["min_nodes"]; ok {
		if minNodes, ok := v.(float64); ok {
			mn := int32(minNodes)
			partition.Minimums = &types.PartitionMinimums{
				Nodes: &mn,
			}
		}
	}
	// Defaults - nested structure
	if v, ok := partitionData["default_time"]; ok {
		var defaultTime uint32
		if timeStruct, ok := v.(map[string]interface{}); ok {
			if number, ok := timeStruct["number"].(float64); ok {
				defaultTime = uint32(number)
			}
		} else if dt, ok := v.(float64); ok {
			defaultTime = uint32(dt)
		}
		if defaultTime > 0 {
			partition.Defaults = &types.PartitionDefaults{
				Time: &defaultTime,
			}
		}
	}
	// Priority - nested structure
	priority := &types.PartitionPriority{}
	hasPriority := false
	if v, ok := partitionData["priority_job_factor"]; ok {
		if pf, ok := v.(float64); ok {
			factor := int32(pf)
			priority.JobFactor = &factor
			hasPriority = true
		}
	}
	if v, ok := partitionData["priority_tier"]; ok {
		if tier, ok := v.(float64); ok {
			t := int32(tier)
			priority.Tier = &t
			hasPriority = true
		}
	}
	if hasPriority {
		partition.Priority = priority
	}
	// Groups - nested structure
	if v, ok := partitionData["allowed_groups"]; ok {
		if groups, ok := v.(string); ok {
			partition.Groups = &types.PartitionGroups{
				Allowed: &groups,
			}
		}
	}
	// Accounts - nested structure
	if v, ok := partitionData["allowed_accounts"]; ok {
		if accounts, ok := v.(string); ok {
			partition.Accounts = &types.PartitionAccounts{
				Allowed: &accounts,
			}
		}
	}
	// QoS - nested structure
	if v, ok := partitionData["allowed_qos"]; ok {
		if qos, ok := v.(string); ok {
			partition.QoS = &types.PartitionQoS{
				Allowed: &qos,
			}
		}
	}
	// Grace time
	if v, ok := partitionData["grace_time"]; ok {
		if graceTime, ok := v.(float64); ok {
			gt := int32(graceTime)
			partition.GraceTime = &gt
		}
	}
	return partition, nil
}
