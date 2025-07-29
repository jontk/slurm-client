package v0_0_40

import (
	"strings"
	"time"

	"github.com/jontk/slurm-client/internal/common/types"
	api "github.com/jontk/slurm-client/internal/api/v0_0_40"
)

// convertAPIPartitionToCommon converts a v0.0.40 API Partition to common Partition type
func (a *PartitionAdapter) convertAPIPartitionToCommon(apiPartition api.V0040PartitionInfo) (*types.Partition, error) {
	partition := &types.Partition{}

	// Basic fields
	if apiPartition.Name != nil {
		partition.Name = *apiPartition.Name
	}

	// State
	if apiPartition.State != nil && len(*apiPartition.State) > 0 {
		// Convert state array to single state
		stateStr := (*apiPartition.State)[0]
		partition.State = types.PartitionState(strings.ToUpper(stateStr))
	}

	// Time limits
	if apiPartition.MaxTime != nil && apiPartition.MaxTime.Number != nil {
		partition.MaxTime = *apiPartition.MaxTime.Number
	}
	if apiPartition.DefaultTime != nil && apiPartition.DefaultTime.Number != nil {
		partition.DefaultTime = *apiPartition.DefaultTime.Number
	}

	// Node limits
	if apiPartition.MaxNodesPerJob != nil && apiPartition.MaxNodesPerJob.Number != nil {
		partition.MaxNodes = *apiPartition.MaxNodesPerJob.Number
	}
	if apiPartition.MinNodesPerJob != nil && apiPartition.MinNodesPerJob.Number != nil {
		partition.MinNodes = *apiPartition.MinNodesPerJob.Number
	}

	// Default flag
	if apiPartition.Defaults != nil && apiPartition.Defaults.Default != nil {
		partition.Default = *apiPartition.Defaults.Default
	}

	// Node lists
	if apiPartition.NodeNames != nil {
		partition.Nodes = *apiPartition.NodeNames
	}
	if apiPartition.AllocNodes != nil {
		partition.AllocNodes = *apiPartition.AllocNodes
	}
	if apiPartition.DeniedNodes != nil {
		partition.DeniedNodes = *apiPartition.DeniedNodes
	}

	// Groups and accounts
	if apiPartition.AllowedGroups != nil && len(*apiPartition.AllowedGroups) > 0 {
		partition.AllowGroups = *apiPartition.AllowedGroups
	}
	if apiPartition.AllowedAccounts != nil && len(*apiPartition.AllowedAccounts) > 0 {
		partition.AllowAccounts = *apiPartition.AllowedAccounts
	}
	if apiPartition.DeniedAccounts != nil && len(*apiPartition.DeniedAccounts) > 0 {
		partition.DenyAccounts = *apiPartition.DeniedAccounts
	}

	// QoS
	if apiPartition.AllowedQos != nil && len(*apiPartition.AllowedQos) > 0 {
		// Join QoS list into comma-separated string
		partition.QoS = strings.Join(*apiPartition.AllowedQos, ",")
	}
	if apiPartition.DeniedQos != nil && len(*apiPartition.DeniedQos) > 0 {
		partition.DenyQoS = strings.Join(*apiPartition.DeniedQos, ",")
	}

	// Flags
	if apiPartition.Flags != nil && len(*apiPartition.Flags) > 0 {
		partition.Flags = make([]string, len(*apiPartition.Flags))
		for i, flag := range *apiPartition.Flags {
			partition.Flags[i] = string(flag)
		}
	}

	// Preemption mode
	if apiPartition.PreemptionMode != nil && len(*apiPartition.PreemptionMode) > 0 {
		partition.PreemptMode = strings.Join(*apiPartition.PreemptionMode, ",")
	}

	// Job defaults
	if apiPartition.Defaults != nil {
		if apiPartition.Defaults.MemoryPerCpu != nil && apiPartition.Defaults.MemoryPerCpu.Number != nil {
			memPerCPU := int64(*apiPartition.Defaults.MemoryPerCpu.Number)
			partition.DefMemPerCPU = &memPerCPU
		}
		if apiPartition.Defaults.Time != nil && apiPartition.Defaults.Time.Number != nil {
			partition.DefaultTime = *apiPartition.Defaults.Time.Number
		}
	}

	// Miscellaneous fields
	if apiPartition.MaximumMemoryPerNode != nil && apiPartition.MaximumMemoryPerNode.Number != nil {
		maxMem := int64(*apiPartition.MaximumMemoryPerNode.Number)
		partition.MaxMemPerNode = &maxMem
	}
	if apiPartition.MaximumCpusPerNode != nil && apiPartition.MaximumCpusPerNode.Number != nil {
		partition.MaxCPUsPerNode = apiPartition.MaximumCpusPerNode.Number
	}
	if apiPartition.Priority != nil && apiPartition.Priority.Number != nil {
		partition.Priority = apiPartition.Priority.Number
	}
	if apiPartition.Suspend != nil && apiPartition.Suspend.Number != nil {
		suspendTime := *apiPartition.Suspend.Number
		partition.SuspendTime = &suspendTime
	}

	// Grace time
	if apiPartition.GraceTime != nil {
		partition.GraceTime = *apiPartition.GraceTime
	}

	// Set modification time to current time (API doesn't provide this)
	partition.LastUpdate = time.Now()

	return partition, nil
}

// convertCommonPartitionCreateToAPI converts common PartitionCreate to v0.0.40 API format
func (a *PartitionAdapter) convertCommonPartitionCreateToAPI(partition *types.PartitionCreate) (*api.V0040PartitionInfo, error) {
	apiPartition := &api.V0040PartitionInfo{}

	// Basic fields
	apiPartition.Name = &partition.Name

	// State
	if partition.State != "" {
		state := []string{string(partition.State)}
		apiPartition.State = &state
	}

	// Time limits
	if partition.MaxTime > 0 {
		maxTime := api.V0040Uint32NoVal{
			Set:    boolPtr(true),
			Number: int32Ptr(partition.MaxTime),
		}
		apiPartition.MaxTime = &maxTime
	}
	if partition.DefaultTime > 0 {
		defaultTime := api.V0040Uint32NoVal{
			Set:    boolPtr(true),
			Number: int32Ptr(partition.DefaultTime),
		}
		apiPartition.DefaultTime = &defaultTime
	}

	// Node limits
	if partition.MaxNodes > 0 {
		maxNodes := api.V0040Uint32NoVal{
			Set:    boolPtr(true),
			Number: int32Ptr(partition.MaxNodes),
		}
		apiPartition.MaxNodesPerJob = &maxNodes
	}
	if partition.MinNodes > 0 {
		minNodes := api.V0040Uint32NoVal{
			Set:    boolPtr(true),
			Number: int32Ptr(partition.MinNodes),
		}
		apiPartition.MinNodesPerJob = &minNodes
	}

	// Default flag
	defaults := &api.V0040PartitionInfoDefaults{
		Default: &partition.Default,
	}
	apiPartition.Defaults = defaults

	// Node lists
	if len(partition.Nodes) > 0 {
		nodeNames := strings.Join(partition.Nodes, ",")
		apiPartition.NodeNames = &nodeNames
	}

	// Groups and accounts
	if len(partition.AllowGroups) > 0 {
		apiPartition.AllowedGroups = &partition.AllowGroups
	}
	if len(partition.AllowAccounts) > 0 {
		apiPartition.AllowedAccounts = &partition.AllowAccounts
	}
	if len(partition.DenyGroups) > 0 {
		// v0.0.40 doesn't have DeniedGroups, so we might need to handle this differently
	}
	if len(partition.DenyAccounts) > 0 {
		apiPartition.DeniedAccounts = &partition.DenyAccounts
	}

	// QoS
	if partition.QoS != "" {
		qosList := strings.Split(partition.QoS, ",")
		apiPartition.AllowedQos = &qosList
	}
	if partition.DenyQoS != "" {
		qosList := strings.Split(partition.DenyQoS, ",")
		apiPartition.DeniedQos = &qosList
	}

	// Flags
	if len(partition.Flags) > 0 {
		flags := make([]api.V0040PartitionInfoFlags, len(partition.Flags))
		for i, flag := range partition.Flags {
			flags[i] = api.V0040PartitionInfoFlags(flag)
		}
		apiPartition.Flags = &flags
	}

	// Preemption mode
	if partition.PreemptMode != "" {
		preemptModes := strings.Split(partition.PreemptMode, ",")
		apiPartition.PreemptionMode = &preemptModes
	}

	// Memory per CPU
	if partition.DefMemPerCPU != nil && *partition.DefMemPerCPU > 0 {
		memPerCPU := api.V0040Uint64NoVal{
			Set:    boolPtr(true),
			Number: int64Ptr(*partition.DefMemPerCPU),
		}
		if apiPartition.Defaults == nil {
			apiPartition.Defaults = &api.V0040PartitionInfoDefaults{}
		}
		apiPartition.Defaults.MemoryPerCpu = &memPerCPU
	}

	// Priority
	if partition.Priority != nil && *partition.Priority > 0 {
		priority := api.V0040Uint32NoVal{
			Set:    boolPtr(true),
			Number: partition.Priority,
		}
		apiPartition.Priority = &priority
	}

	// Grace time
	if partition.GraceTime > 0 {
		apiPartition.GraceTime = &partition.GraceTime
	}

	return apiPartition, nil
}

// convertCommonPartitionUpdateToAPI converts common PartitionUpdate to v0.0.40 API format
func (a *PartitionAdapter) convertCommonPartitionUpdateToAPI(existingPartition *types.Partition, update *types.PartitionUpdate) (*api.V0040PartitionInfo, error) {
	apiPartition := &api.V0040PartitionInfo{}

	// Name (required)
	apiPartition.Name = &existingPartition.Name

	// Apply updates
	if update.State != nil {
		state := []string{string(*update.State)}
		apiPartition.State = &state
	}

	if update.MaxTime != nil {
		maxTime := api.V0040Uint32NoVal{
			Set:    boolPtr(true),
			Number: update.MaxTime,
		}
		apiPartition.MaxTime = &maxTime
	}

	if update.DefaultTime != nil {
		defaultTime := api.V0040Uint32NoVal{
			Set:    boolPtr(true),
			Number: update.DefaultTime,
		}
		apiPartition.DefaultTime = &defaultTime
	}

	if update.MaxNodes != nil {
		maxNodes := api.V0040Uint32NoVal{
			Set:    boolPtr(true),
			Number: update.MaxNodes,
		}
		apiPartition.MaxNodesPerJob = &maxNodes
	}

	if update.MinNodes != nil {
		minNodes := api.V0040Uint32NoVal{
			Set:    boolPtr(true),
			Number: update.MinNodes,
		}
		apiPartition.MinNodesPerJob = &minNodes
	}

	if update.Default != nil {
		defaults := &api.V0040PartitionInfoDefaults{
			Default: update.Default,
		}
		apiPartition.Defaults = defaults
	}

	if update.Nodes != nil {
		nodeNames := strings.Join(*update.Nodes, ",")
		apiPartition.NodeNames = &nodeNames
	}

	if update.AllowGroups != nil {
		apiPartition.AllowedGroups = update.AllowGroups
	}

	if update.AllowAccounts != nil {
		apiPartition.AllowedAccounts = update.AllowAccounts
	}

	if update.DenyAccounts != nil {
		apiPartition.DeniedAccounts = update.DenyAccounts
	}

	if update.QoS != nil {
		qosList := strings.Split(*update.QoS, ",")
		apiPartition.AllowedQos = &qosList
	}

	if update.Flags != nil {
		flags := make([]api.V0040PartitionInfoFlags, len(*update.Flags))
		for i, flag := range *update.Flags {
			flags[i] = api.V0040PartitionInfoFlags(flag)
		}
		apiPartition.Flags = &flags
	}

	return apiPartition, nil
}