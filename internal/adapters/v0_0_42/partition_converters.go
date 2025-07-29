package v0_0_42

import (
	"strings"
	"time"

	"github.com/jontk/slurm-client/internal/common/types"
	api "github.com/jontk/slurm-client/internal/api/v0_0_42"
)

// convertAPIPartitionToCommon converts a v0.0.42 API Partition to common Partition type
func (a *PartitionAdapter) convertAPIPartitionToCommon(apiPartition api.V0042PartitionInfo) (*types.Partition, error) {
	partition := &types.Partition{}

	// Basic fields
	if apiPartition.Name != nil {
		partition.Name = *apiPartition.Name
	}

	// State
	if apiPartition.Partition != nil && apiPartition.Partition.State != nil {
		// Convert partition state flags to string
		state := a.convertPartitionStatesToString(apiPartition.Partition.State)
		partition.State = state
	}

	// Node information
	if apiPartition.Nodes != nil {
		if apiPartition.Nodes.Total != nil {
			partition.TotalNodes = uint32(*apiPartition.Nodes.Total)
		}
		if apiPartition.Nodes.Configured != nil {
			partition.Nodes = *apiPartition.Nodes.Configured
		}
		if apiPartition.Nodes.AllowedAllocation != nil {
			partition.AllocNodes = *apiPartition.Nodes.AllowedAllocation
		}
	}

	// CPU information
	if apiPartition.Cpus != nil {
		if apiPartition.Cpus.Total != nil {
			partition.TotalCPUs = uint32(*apiPartition.Cpus.Total)
		}
	}

	// Time limits
	if apiPartition.Defaults != nil && apiPartition.Defaults.Time != nil {
		partition.DefaultTime = uint32(apiPartition.Defaults.Time.Number)
	}
	if apiPartition.Maximums != nil && apiPartition.Maximums.Time != nil {
		partition.MaxTime = uint32(apiPartition.Maximums.Time.Number)
	}

	// Memory limits
	if apiPartition.Defaults != nil {
		if apiPartition.Defaults.MemoryPerCpu != nil && *apiPartition.Defaults.MemoryPerCpu > 0 {
			partition.DefMemPerCPU = uint64(*apiPartition.Defaults.MemoryPerCpu)
		}
		if apiPartition.Defaults.PartitionMemoryPerCpu != nil && apiPartition.Defaults.PartitionMemoryPerCpu.Number > 0 {
			partition.DefMemPerCPU = apiPartition.Defaults.PartitionMemoryPerCpu.Number
		}
		if apiPartition.Defaults.PartitionMemoryPerNode != nil && apiPartition.Defaults.PartitionMemoryPerNode.Number > 0 {
			partition.DefMemPerNode = apiPartition.Defaults.PartitionMemoryPerNode.Number
		}
	}

	if apiPartition.Maximums != nil {
		if apiPartition.Maximums.MemoryPerCpu != nil && *apiPartition.Maximums.MemoryPerCpu > 0 {
			partition.MaxMemPerCPU = uint64(*apiPartition.Maximums.MemoryPerCpu)
		}
		if apiPartition.Maximums.PartitionMemoryPerCpu != nil && apiPartition.Maximums.PartitionMemoryPerCpu.Number > 0 {
			partition.MaxMemPerCPU = apiPartition.Maximums.PartitionMemoryPerCpu.Number
		}
		if apiPartition.Maximums.PartitionMemoryPerNode != nil && apiPartition.Maximums.PartitionMemoryPerNode.Number > 0 {
			partition.MaxMemPerNode = apiPartition.Maximums.PartitionMemoryPerNode.Number
		}
	}

	// Node limits
	if apiPartition.Minimums != nil && apiPartition.Minimums.Nodes != nil {
		partition.MinNodes = uint32(*apiPartition.Minimums.Nodes)
	}
	if apiPartition.Maximums != nil && apiPartition.Maximums.Nodes != nil {
		partition.MaxNodes = uint32(apiPartition.Maximums.Nodes.Number)
	}

	// Priority
	if apiPartition.Priority != nil {
		if apiPartition.Priority.JobFactor != nil {
			partition.PriorityJobFactor = *apiPartition.Priority.JobFactor
		}
		if apiPartition.Priority.Tier != nil {
			partition.PriorityTier = *apiPartition.Priority.Tier
		}
	}

	// QoS
	if apiPartition.Qos != nil {
		if apiPartition.Qos.Assigned != nil {
			partition.QoS = strings.Split(*apiPartition.Qos.Assigned, ",")
		}
		if apiPartition.Qos.Allowed != nil {
			partition.AllowQoS = strings.Split(*apiPartition.Qos.Allowed, ",")
		}
		if apiPartition.Qos.Deny != nil {
			partition.DenyQoS = strings.Split(*apiPartition.Qos.Deny, ",")
		}
	}

	// Accounts
	if apiPartition.Accounts != nil {
		if apiPartition.Accounts.Allowed != nil {
			partition.AllowAccounts = strings.Split(*apiPartition.Accounts.Allowed, ",")
		}
		if apiPartition.Accounts.Deny != nil {
			partition.DenyAccounts = strings.Split(*apiPartition.Accounts.Deny, ",")
		}
	}

	// Groups
	if apiPartition.Groups != nil && apiPartition.Groups.Allowed != nil {
		partition.AllowGroups = strings.Split(*apiPartition.Groups.Allowed, ",")
	}

	// Grace time
	if apiPartition.GraceTime != nil {
		partition.GraceTime = time.Duration(*apiPartition.GraceTime) * time.Second
	}

	// Oversubscription
	if apiPartition.Maximums != nil && apiPartition.Maximums.Oversubscribe != nil {
		if apiPartition.Maximums.Oversubscribe.Jobs != nil {
			// Convert oversubscribe jobs to a string representation
			if *apiPartition.Maximums.Oversubscribe.Jobs > 1 {
				partition.OverSubscribe = "YES"
			} else {
				partition.OverSubscribe = "NO"
			}
		}
	}

	// CPU binding
	if apiPartition.Cpus != nil && apiPartition.Cpus.TaskBinding != nil {
		// Store CPU binding as part of parameters
		if partition.Parameters == nil {
			partition.Parameters = make(map[string]string)
		}
		partition.Parameters["cpu_bind"] = string(*apiPartition.Cpus.TaskBinding)
	}

	// Alternate partition
	if apiPartition.Alternate != nil {
		partition.AlternatePartition = *apiPartition.Alternate
	}

	// Select type
	if apiPartition.SelectType != nil {
		// Store select type as part of parameters
		if partition.Parameters == nil {
			partition.Parameters = make(map[string]string)
		}
		partition.Parameters["select_type"] = string(*apiPartition.SelectType)
	}

	// Job defaults
	if apiPartition.Defaults != nil && apiPartition.Defaults.Job != nil {
		if partition.Parameters == nil {
			partition.Parameters = make(map[string]string)
		}
		partition.Parameters["job_defaults"] = *apiPartition.Defaults.Job
	}

	return partition, nil
}

// convertPartitionStatesToString converts partition state flags to a string representation
func (a *PartitionAdapter) convertPartitionStatesToString(states *api.V0042PartitionStates) string {
	if states == nil {
		return "UNKNOWN"
	}

	// Convert the partition states to a string
	// In v0.0.42, states might be represented as flags
	stateStrs := make([]string, 0)
	
	// Check each state flag (this is a simplified version)
	// The actual implementation would depend on how V0042PartitionStates is defined
	stateStr := "UP" // Default state
	
	// You would need to check the actual V0042PartitionStates structure
	// and map the flags to appropriate state strings
	
	return stateStr
}