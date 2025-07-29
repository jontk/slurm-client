package v0_0_41

import (
	"fmt"
	"strings"
	"time"

	"github.com/jontk/slurm-client/internal/common/types"
	api "github.com/jontk/slurm-client/internal/api/v0_0_41"
)

// convertAPIPartitionToCommon converts a v0.0.41 API Partition to common Partition type
func (a *PartitionAdapter) convertAPIPartitionToCommon(apiPartition interface{}) (*types.Partition, error) {
	// Type assertion to handle the anonymous struct
	partData, ok := apiPartition.(struct {
		Accounts               *struct {
			Allowed *string `json:"allowed,omitempty"`
			Deny    *string `json:"deny,omitempty"`
		} `json:"accounts,omitempty"`
		AlternativeName        *string `json:"alternative_name,omitempty"`
		BillingWeights         *string `json:"billing_weights,omitempty"`
		DefaultMemoryPerCpu    *api.V0041OpenapiPartitionRespPartitionsDefaultMemoryPerCpu `json:"default_memory_per_cpu,omitempty"`
		DefaultTimeLimit       *api.V0041OpenapiPartitionRespPartitionsDefaultTimeLimit `json:"default_time_limit,omitempty"`
		DenyAccounts           *string `json:"deny_accounts,omitempty"`
		DenyQos                *string `json:"deny_qos,omitempty"`
		GraceTime              *int32  `json:"grace_time,omitempty"`
		Maximums               *struct {
			CpusPerNode           *api.V0041OpenapiPartitionRespPartitionsMaximumsCpusPerNode           `json:"cpus_per_node,omitempty"`
			CpusPerSocket         *api.V0041OpenapiPartitionRespPartitionsMaximumsCpusPerSocket         `json:"cpus_per_socket,omitempty"`
			MemoryPerCpu          *api.V0041OpenapiPartitionRespPartitionsMaximumsMemoryPerCpu          `json:"memory_per_cpu,omitempty"`
			MemoryPerNode         *api.V0041OpenapiPartitionRespPartitionsMaximumsMemoryPerNode         `json:"memory_per_node,omitempty"`
			Nodes                 *api.V0041OpenapiPartitionRespPartitionsMaximumsNodes                 `json:"nodes,omitempty"`
			OversubscribeFlags    *[]api.V0041OpenapiPartitionRespPartitionsMaximumsOversubscribeFlags  `json:"oversubscribe_flags,omitempty"`
			OversubscribeJobs     *api.V0041OpenapiPartitionRespPartitionsMaximumsOversubscribeJobs     `json:"oversubscribe_jobs,omitempty"`
			PartitionWall         *api.V0041OpenapiPartitionRespPartitionsMaximumsPartitionWall         `json:"partition_wall,omitempty"`
			SharedSpace           *api.V0041OpenapiPartitionRespPartitionsMaximumsSharedSpace           `json:"shared_space,omitempty"`
			Time                  *api.V0041OpenapiPartitionRespPartitionsMaximumsTime                  `json:"time,omitempty"`
		} `json:"maximums,omitempty"`
		MaximumMemoryPerNode   *api.V0041OpenapiPartitionRespPartitionsMaximumMemoryPerNode `json:"maximum_memory_per_node,omitempty"`
		MaximumMemoryPerCpu    *api.V0041OpenapiPartitionRespPartitionsMaximumMemoryPerCpu  `json:"maximum_memory_per_cpu,omitempty"`
		Minimums               *struct {
			Nodes *api.V0041OpenapiPartitionRespPartitionsMinimumsNodes `json:"nodes,omitempty"`
		} `json:"minimums,omitempty"`
		Name                   *string   `json:"name,omitempty"`
		Nodes                  *string   `json:"nodes,omitempty"`
		AllowedAllocationNodes *string   `json:"allowed_allocation_nodes,omitempty"`
		PartitionState         *[]api.V0041OpenapiPartitionRespPartitionsPartitionState `json:"partition_state,omitempty"`
		PreemptMode            *[]string `json:"preempt_mode,omitempty"`
		Priority               *struct {
			JobFactor *api.V0041OpenapiPartitionRespPartitionsPriorityJobFactor `json:"job_factor,omitempty"`
			Tier      *api.V0041OpenapiPartitionRespPartitionsPriorityTier      `json:"tier,omitempty"`
		} `json:"priority,omitempty"`
		Qos                    *struct {
			Allowed *string `json:"allowed,omitempty"`
			Deny    *string `json:"deny,omitempty"`
			Forced  *string `json:"forced,omitempty"`
		} `json:"qos,omitempty"`
		DefaultQos             *string `json:"default_qos,omitempty"`
		SelectType             *[]api.V0041OpenapiPartitionRespPartitionsSelectType `json:"select_type,omitempty"`
		TresEnabled            *string `json:"tres_enabled,omitempty"`
		Timeouts               *struct {
			ResumeTimeout  *api.V0041OpenapiPartitionRespPartitionsTimeoutsResumeTimeout  `json:"resume_timeout,omitempty"`
			SuspendTimeout *api.V0041OpenapiPartitionRespPartitionsTimeoutsSuspendTimeout `json:"suspend_timeout,omitempty"`
			SuspendTime    *api.V0041OpenapiPartitionRespPartitionsTimeoutsSuspendTime    `json:"suspend_time,omitempty"`
		} `json:"timeouts,omitempty"`
		SuspendTime            *api.V0041OpenapiPartitionRespPartitionsSuspendTime `json:"suspend_time,omitempty"`
	})
	if !ok {
		return nil, fmt.Errorf("unexpected partition data type")
	}

	partition := &types.Partition{}

	// Basic fields
	if partData.Name != nil {
		partition.Name = *partData.Name
	}
	if partData.Nodes != nil {
		partition.Nodes = *partData.Nodes
	}
	if partData.AlternativeName != nil {
		partition.AlternateName = *partData.AlternativeName
	}

	// State
	if partData.PartitionState != nil && len(*partData.PartitionState) > 0 {
		// Convert state to string
		partition.State = string((*partData.PartitionState)[0])
	}

	// Time limits
	if partData.DefaultTimeLimit != nil && partData.DefaultTimeLimit.Number != nil {
		partition.DefaultTime = time.Duration(*partData.DefaultTimeLimit.Number) * time.Minute
	}
	if partData.Maximums != nil && partData.Maximums.Time != nil && partData.Maximums.Time.Number != nil {
		partition.MaxTime = time.Duration(*partData.Maximums.Time.Number) * time.Minute
	}

	// Memory settings
	if partData.DefaultMemoryPerCpu != nil && partData.DefaultMemoryPerCpu.Number != nil {
		partition.DefMemPerCPU = uint64(*partData.DefaultMemoryPerCpu.Number)
	}
	if partData.MaximumMemoryPerCpu != nil && partData.MaximumMemoryPerCpu.Number != nil {
		partition.MaxMemPerCPU = uint64(*partData.MaximumMemoryPerCpu.Number)
	}
	if partData.MaximumMemoryPerNode != nil && partData.MaximumMemoryPerNode.Number != nil {
		partition.MaxMemPerNode = uint64(*partData.MaximumMemoryPerNode.Number)
	}

	// Node counts
	if partData.Maximums != nil && partData.Maximums.Nodes != nil && partData.Maximums.Nodes.Number != nil {
		partition.MaxNodes = uint32(*partData.Maximums.Nodes.Number)
	}
	if partData.Minimums != nil && partData.Minimums.Nodes != nil && partData.Minimums.Nodes.Number != nil {
		partition.MinNodes = uint32(*partData.Minimums.Nodes.Number)
	}

	// Priority
	if partData.Priority != nil {
		if partData.Priority.Tier != nil && partData.Priority.Tier.Number != nil {
			partition.PriorityTier = *partData.Priority.Tier.Number
		}
		if partData.Priority.JobFactor != nil && partData.Priority.JobFactor.Number != nil {
			partition.PriorityJobFactor = *partData.Priority.JobFactor.Number
		}
	}

	// QoS
	if partData.DefaultQos != nil {
		partition.QoS = *partData.DefaultQos
	}
	if partData.Qos != nil {
		if partData.Qos.Allowed != nil {
			partition.AllowQos = strings.Split(*partData.Qos.Allowed, ",")
		}
		if partData.Qos.Deny != nil {
			partition.DenyQos = strings.Split(*partData.Qos.Deny, ",")
		}
	}

	// Accounts
	if partData.Accounts != nil {
		if partData.Accounts.Allowed != nil {
			partition.AllowAccounts = strings.Split(*partData.Accounts.Allowed, ",")
		}
		if partData.Accounts.Deny != nil {
			partition.DenyAccounts = strings.Split(*partData.Accounts.Deny, ",")
		}
	}

	// Grace time
	if partData.GraceTime != nil {
		partition.GraceTime = uint32(*partData.GraceTime)
	}

	// Preempt mode
	if partData.PreemptMode != nil {
		partition.PreemptMode = strings.Join(*partData.PreemptMode, ",")
	}

	// Oversubscribe
	if partData.Maximums != nil && partData.Maximums.OversubscribeFlags != nil && len(*partData.Maximums.OversubscribeFlags) > 0 {
		// Convert oversubscribe flags to string
		var flags []string
		for _, flag := range *partData.Maximums.OversubscribeFlags {
			flags = append(flags, string(flag))
		}
		partition.OverSubscribe = strings.Join(flags, ",")
	}

	// TREs
	if partData.TresEnabled != nil {
		partition.TRESBillingWeights = *partData.TresEnabled
	}

	// Select type
	if partData.SelectType != nil && len(*partData.SelectType) > 0 {
		// Convert select type to string
		var selectTypes []string
		for _, st := range *partData.SelectType {
			selectTypes = append(selectTypes, string(st))
		}
		partition.SelectTypeParameters = strings.Join(selectTypes, ",")
	}

	// Timeouts
	if partData.Timeouts != nil {
		if partData.Timeouts.ResumeTimeout != nil && partData.Timeouts.ResumeTimeout.Number != nil {
			partition.ResumeTimeout = uint32(*partData.Timeouts.ResumeTimeout.Number)
		}
		if partData.Timeouts.SuspendTimeout != nil && partData.Timeouts.SuspendTimeout.Number != nil {
			partition.SuspendTimeout = uint32(*partData.Timeouts.SuspendTimeout.Number)
		}
		if partData.Timeouts.SuspendTime != nil && partData.Timeouts.SuspendTime.Number != nil {
			partition.SuspendTime = uint32(*partData.Timeouts.SuspendTime.Number)
		}
	}

	// Billing weights
	if partData.BillingWeights != nil {
		partition.BillingWeights = *partData.BillingWeights
	}

	// Calculate total nodes if possible
	if partition.Nodes != "" {
		nodeList := parseNodeList(partition.Nodes)
		partition.TotalNodes = uint32(len(nodeList))
	}

	return partition, nil
}

// convertCommonToAPIPartition converts common Partition to v0.0.41 API request (not supported)
func (a *PartitionAdapter) convertCommonToAPIPartition(partition *types.Partition) interface{} {
	// v0.0.41 doesn't support partition creation/update through the API
	// This function is here for interface compliance but returns nil
	return nil
}

// convertCommonToAPIPartitionUpdate converts common PartitionUpdate to v0.0.41 API request (not supported)
func (a *PartitionAdapter) convertCommonToAPIPartitionUpdate(update *types.PartitionUpdate) interface{} {
	// v0.0.41 doesn't support partition update through the API
	// This function is here for interface compliance but returns nil
	return nil
}