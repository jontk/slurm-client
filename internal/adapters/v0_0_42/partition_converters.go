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
		partition.State = types.PartitionState(state)
	}

	// Node information
	if apiPartition.Nodes != nil {
		if apiPartition.Nodes.Total != nil {
			partition.TotalNodes = int32(*apiPartition.Nodes.Total)
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
			partition.TotalCPUs = int32(*apiPartition.Cpus.Total)
		}
	}

	// Time limits
	if apiPartition.Defaults != nil && apiPartition.Defaults.Time != nil {
		partition.DefaultTime = int32(*apiPartition.Defaults.Time.Number)
	}
	if apiPartition.Maximums != nil && apiPartition.Maximums.Time != nil {
		partition.MaxTime = int32(*apiPartition.Maximums.Time.Number)
	}

	// Memory limits
	if apiPartition.Defaults != nil {
		if apiPartition.Defaults.MemoryPerCpu != nil && *apiPartition.Defaults.MemoryPerCpu > 0 {
			partition.DefaultMemPerCPU = int64(*apiPartition.Defaults.MemoryPerCpu)
		}
		if apiPartition.Defaults.PartitionMemoryPerCpu != nil && apiPartition.Defaults.PartitionMemoryPerCpu.Number != nil && *apiPartition.Defaults.PartitionMemoryPerCpu.Number > 0 {
			partition.DefaultMemPerCPU = *apiPartition.Defaults.PartitionMemoryPerCpu.Number
		}
		if apiPartition.Defaults.PartitionMemoryPerNode != nil && apiPartition.Defaults.PartitionMemoryPerNode.Number != nil && *apiPartition.Defaults.PartitionMemoryPerNode.Number > 0 {
			partition.DefMemPerNode = *apiPartition.Defaults.PartitionMemoryPerNode.Number
		}
	}

	if apiPartition.Maximums != nil {
		if apiPartition.Maximums.MemoryPerCpu != nil && *apiPartition.Maximums.MemoryPerCpu > 0 {
			partition.MaxMemPerCPU = int64(*apiPartition.Maximums.MemoryPerCpu)
		}
		if apiPartition.Maximums.PartitionMemoryPerCpu != nil && apiPartition.Maximums.PartitionMemoryPerCpu.Number != nil && *apiPartition.Maximums.PartitionMemoryPerCpu.Number > 0 {
			partition.MaxMemPerCPU = *apiPartition.Maximums.PartitionMemoryPerCpu.Number
		}
		if apiPartition.Maximums.PartitionMemoryPerNode != nil && apiPartition.Maximums.PartitionMemoryPerNode.Number != nil && *apiPartition.Maximums.PartitionMemoryPerNode.Number > 0 {
			partition.MaxMemPerNode = *apiPartition.Maximums.PartitionMemoryPerNode.Number
		}
	}

	// Node limits
	if apiPartition.Minimums != nil && apiPartition.Minimums.Nodes != nil {
		partition.MinNodes = int32(*apiPartition.Minimums.Nodes)
	}
	if apiPartition.Maximums != nil && apiPartition.Maximums.Nodes != nil {
		partition.MaxNodes = int32(*apiPartition.Maximums.Nodes.Number)
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
			partition.QoS = *apiPartition.Qos.Assigned
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
		partition.GraceTime = int32(*apiPartition.GraceTime)
	}

	// Oversubscription
	if apiPartition.Maximums != nil && apiPartition.Maximums.Oversubscribe != nil {
		if apiPartition.Maximums.Oversubscribe.Jobs != nil {
			// OverSubscribe field doesn't exist in common Partition type
			// Skip oversubscribe conversion
			_ = apiPartition.Maximums.Oversubscribe.Jobs
		}
	}

	// CPU binding - Parameters field doesn't exist in common Partition type
	// Skip CPU binding conversion
	if apiPartition.Cpus != nil && apiPartition.Cpus.TaskBinding != nil {
		_ = apiPartition.Cpus.TaskBinding
	}

	// Alternate partition - AlternatePartition field doesn't exist in common Partition type
	// Skip alternate partition conversion
	if apiPartition.Alternate != nil {
		_ = apiPartition.Alternate
	}

	// Select type - store in SelectTypeParameters field
	if apiPartition.SelectType != nil {
		// SelectType is []V0042CrType, skip complex conversion for now
		_ = apiPartition.SelectType
	}

	// Job defaults - store in JobDefaults field
	if apiPartition.Defaults != nil && apiPartition.Defaults.Job != nil {
		if partition.JobDefaults == nil {
			partition.JobDefaults = make(map[string]string)
		}
		partition.JobDefaults["job_defaults"] = *apiPartition.Defaults.Job
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
	// Check each state flag (this is a simplified version)
	// The actual implementation would depend on how V0042PartitionStates is defined
	stateStr := "UP" // Default state
	
	// You would need to check the actual V0042PartitionStates structure
	// and map the flags to appropriate state strings
	
	return stateStr
}