package v0_0_43

import (
	"strings"

	"github.com/jontk/slurm-client/internal/common/types"
	api "github.com/jontk/slurm-client/internal/api/v0_0_43"
)

// convertAPIPartitionToCommon converts a v0.0.43 API Partition to common Partition type
func (a *PartitionAdapter) convertAPIPartitionToCommon(apiPartition api.V0043PartitionInfo) (*types.Partition, error) {
	partition := &types.Partition{}

	// Basic fields
	if apiPartition.Name != nil {
		partition.Name = *apiPartition.Name
	}
	if apiPartition.AllocNodes != nil {
		partition.AllocNodes = *apiPartition.AllocNodes
	}
	if apiPartition.AllowAccounts != nil {
		partition.AllowAccounts = *apiPartition.AllowAccounts
	}
	if apiPartition.AllowAllocNodes != nil {
		partition.AllowAllocNodes = *apiPartition.AllowAllocNodes
	}
	if apiPartition.AllowGroups != nil {
		partition.AllowGroups = *apiPartition.AllowGroups
	}
	if apiPartition.AllowQos != nil {
		partition.AllowQoS = *apiPartition.AllowQos
	}
	if apiPartition.DenyAccounts != nil {
		partition.DenyAccounts = *apiPartition.DenyAccounts
	}
	if apiPartition.DenyQos != nil {
		partition.DenyQoS = *apiPartition.DenyQos
	}

	// Memory settings
	if apiPartition.DefaultMemoryPerCpu != nil && apiPartition.DefaultMemoryPerCpu.Number != nil {
		partition.DefaultMemPerCPU = *apiPartition.DefaultMemoryPerCpu.Number
	}
	if apiPartition.DefaultMemoryPerNode != nil && apiPartition.DefaultMemoryPerNode.Number != nil {
		partition.DefaultMemPerNode = *apiPartition.DefaultMemoryPerNode.Number
	}
	if apiPartition.DefMemPerNode != nil && apiPartition.DefMemPerNode.Number != nil {
		partition.DefMemPerNode = *apiPartition.DefMemPerNode.Number
	}
	if apiPartition.MaxMemoryPerCpu != nil && apiPartition.MaxMemoryPerCpu.Number != nil {
		partition.MaxMemPerCPU = *apiPartition.MaxMemoryPerCpu.Number
	}
	if apiPartition.MaxMemoryPerNode != nil && apiPartition.MaxMemoryPerNode.Number != nil {
		partition.MaxMemPerNode = *apiPartition.MaxMemoryPerNode.Number
	}

	// Time settings
	if apiPartition.DefaultTime != nil && apiPartition.DefaultTime.Number != nil {
		partition.DefaultTime = *apiPartition.DefaultTime.Number
	}
	if apiPartition.DefaultTimeStr != nil {
		partition.DefaultTimeStr = *apiPartition.DefaultTimeStr
	}
	if apiPartition.MaxTime != nil && apiPartition.MaxTime.Number != nil {
		partition.MaxTime = *apiPartition.MaxTime.Number
	}
	if apiPartition.MaxTimeStr != nil {
		partition.MaxTimeStr = *apiPartition.MaxTimeStr
	}
	if apiPartition.GraceTime != nil && apiPartition.GraceTime.Number != nil {
		partition.GraceTime = *apiPartition.GraceTime.Number
	}
	if apiPartition.OverTimeLimit != nil && apiPartition.OverTimeLimit.Number != nil {
		partition.OverTimeLimit = *apiPartition.OverTimeLimit.Number
	}

	// Node settings
	if apiPartition.MaxCpusPerNode != nil && apiPartition.MaxCpusPerNode.Number != nil {
		partition.MaxCPUsPerNode = *apiPartition.MaxCpusPerNode.Number
	}
	if apiPartition.MaxNodes != nil && apiPartition.MaxNodes.Number != nil {
		partition.MaxNodes = *apiPartition.MaxNodes.Number
	}
	if apiPartition.MinNodes != nil && apiPartition.MinNodes.Number != nil {
		partition.MinNodes = *apiPartition.MinNodes.Number
	}
	if apiPartition.Nodes != nil {
		partition.Nodes = *apiPartition.Nodes
	}
	if apiPartition.TotalNodes != nil {
		partition.TotalNodes = *apiPartition.TotalNodes
	}
	if apiPartition.TotalCpus != nil {
		partition.TotalCPUs = *apiPartition.TotalCpus
	}

	// Priority settings
	if apiPartition.Priority != nil && apiPartition.Priority.Number != nil {
		partition.Priority = *apiPartition.Priority.Number
	}
	if apiPartition.PriorityJobFactor != nil && apiPartition.PriorityJobFactor.Number != nil {
		partition.PriorityJobFactor = *apiPartition.PriorityJobFactor.Number
	}
	if apiPartition.PriorityTier != nil && apiPartition.PriorityTier.Number != nil {
		partition.PriorityTier = *apiPartition.PriorityTier.Number
	}

	// QoS settings
	if apiPartition.Qos != nil {
		partition.QoS = *apiPartition.Qos
	}

	// State and reason
	if apiPartition.State != nil && len(*apiPartition.State) > 0 {
		partition.State = types.PartitionState((*apiPartition.State)[0])
	}
	if apiPartition.StateReason != nil {
		partition.StateReason = *apiPartition.StateReason
	}

	// TRES and billing
	if apiPartition.TresStr != nil {
		partition.TresStr = *apiPartition.TresStr
	}
	if apiPartition.BillingWeights != nil {
		partition.BillingWeightStr = *apiPartition.BillingWeights
	}

	// Preempt mode
	if apiPartition.PreemptMode != nil && len(*apiPartition.PreemptMode) > 0 {
		preemptModes := make([]string, len(*apiPartition.PreemptMode))
		for i, mode := range *apiPartition.PreemptMode {
			preemptModes[i] = string(mode)
		}
		partition.PreemptMode = preemptModes
	}

	// Job defaults
	if apiPartition.JobDefaults != nil && len(*apiPartition.JobDefaults) > 0 {
		jobDefaults := make(map[string]string)
		// Parse job defaults format - this might be implementation specific
		for _, defStr := range *apiPartition.JobDefaults {
			parts := strings.SplitN(defStr, "=", 2)
			if len(parts) == 2 {
				jobDefaults[parts[0]] = parts[1]
			}
		}
		partition.JobDefaults = jobDefaults
	}

	// Timeout settings
	if apiPartition.ResumeTimeout != nil && apiPartition.ResumeTimeout.Number != nil {
		partition.ResumeTimeout = *apiPartition.ResumeTimeout.Number
	}
	if apiPartition.SuspendTime != nil && apiPartition.SuspendTime.Number != nil {
		partition.SuspendTime = *apiPartition.SuspendTime.Number
	}
	if apiPartition.SuspendTimeout != nil && apiPartition.SuspendTimeout.Number != nil {
		partition.SuspendTimeout = *apiPartition.SuspendTimeout.Number
	}

	// Boolean flags
	if apiPartition.Flags != nil {
		for _, flag := range *apiPartition.Flags {
			switch flag {
			case api.V0043PartitionInfoFlagsHIDDEN:
				partition.Hidden = true
			case api.V0043PartitionInfoFlagsROOTONLY:
				partition.RootOnly = true
			case api.V0043PartitionInfoFlagsREQRESV:
				partition.ReqResv = true
			case api.V0043PartitionInfoFlagsLLN:
				partition.LLN = true
			case api.V0043PartitionInfoFlagsEXCLUSIVEUSER:
				partition.ExclusiveUser = true
			}
		}
	}

	return partition, nil
}

// convertCommonPartitionCreateToAPI converts common PartitionCreate type to v0.0.43 API format
func (a *PartitionAdapter) convertCommonPartitionCreateToAPI(create *types.PartitionCreate) (*api.V0043PartitionInfo, error) {
	apiPartition := &api.V0043PartitionInfo{}

	// Required fields
	apiPartition.Name = &create.Name

	// Basic fields
	if create.AllocNodes != "" {
		apiPartition.AllocNodes = &create.AllocNodes
	}
	if len(create.AllowAccounts) > 0 {
		apiPartition.AllowAccounts = &create.AllowAccounts
	}
	if create.AllowAllocNodes != "" {
		apiPartition.AllowAllocNodes = &create.AllowAllocNodes
	}
	if len(create.AllowGroups) > 0 {
		apiPartition.AllowGroups = &create.AllowGroups
	}
	if len(create.AllowQoS) > 0 {
		apiPartition.AllowQos = &create.AllowQoS
	}
	if len(create.DenyAccounts) > 0 {
		apiPartition.DenyAccounts = &create.DenyAccounts
	}
	if len(create.DenyQoS) > 0 {
		apiPartition.DenyQos = &create.DenyQoS
	}

	// Memory settings
	if create.DefaultMemPerCPU > 0 {
		setTrue := true
		memory := create.DefaultMemPerCPU
		apiPartition.DefaultMemoryPerCpu = &api.V0043Uint64NoValStruct{
			Set:    &setTrue,
			Number: &memory,
		}
	}
	if create.DefaultMemPerNode > 0 {
		setTrue := true
		memory := create.DefaultMemPerNode
		apiPartition.DefaultMemoryPerNode = &api.V0043Uint64NoValStruct{
			Set:    &setTrue,
			Number: &memory,
		}
	}
	if create.DefMemPerNode > 0 {
		setTrue := true
		memory := create.DefMemPerNode
		apiPartition.DefMemPerNode = &api.V0043Uint64NoValStruct{
			Set:    &setTrue,
			Number: &memory,
		}
	}
	if create.MaxMemPerCPU > 0 {
		setTrue := true
		memory := create.MaxMemPerCPU
		apiPartition.MaxMemoryPerCpu = &api.V0043Uint64NoValStruct{
			Set:    &setTrue,
			Number: &memory,
		}
	}
	if create.MaxMemPerNode > 0 {
		setTrue := true
		memory := create.MaxMemPerNode
		apiPartition.MaxMemoryPerNode = &api.V0043Uint64NoValStruct{
			Set:    &setTrue,
			Number: &memory,
		}
	}

	// Time settings
	if create.DefaultTime > 0 {
		setTrue := true
		time := int32(create.DefaultTime)
		apiPartition.DefaultTime = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &time,
		}
	}
	if create.MaxTime > 0 {
		setTrue := true
		time := int32(create.MaxTime)
		apiPartition.MaxTime = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &time,
		}
	}
	if create.GraceTime > 0 {
		setTrue := true
		graceTime := int32(create.GraceTime)
		apiPartition.GraceTime = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &graceTime,
		}
	}
	if create.OverTimeLimit > 0 {
		setTrue := true
		overTime := int32(create.OverTimeLimit)
		apiPartition.OverTimeLimit = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &overTime,
		}
	}

	// Node settings
	if create.MaxCPUsPerNode > 0 {
		setTrue := true
		cpus := int32(create.MaxCPUsPerNode)
		apiPartition.MaxCpusPerNode = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &cpus,
		}
	}
	if create.MaxNodes > 0 {
		setTrue := true
		nodes := int32(create.MaxNodes)
		apiPartition.MaxNodes = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &nodes,
		}
	}
	if create.MinNodes > 0 {
		setTrue := true
		nodes := int32(create.MinNodes)
		apiPartition.MinNodes = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &nodes,
		}
	}
	if create.Nodes != "" {
		apiPartition.Nodes = &create.Nodes
	}

	// Priority settings
	if create.Priority > 0 {
		setTrue := true
		priority := int32(create.Priority)
		apiPartition.Priority = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &priority,
		}
	}
	if create.PriorityJobFactor > 0 {
		setTrue := true
		factor := int32(create.PriorityJobFactor)
		apiPartition.PriorityJobFactor = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &factor,
		}
	}
	if create.PriorityTier > 0 {
		setTrue := true
		tier := int32(create.PriorityTier)
		apiPartition.PriorityTier = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &tier,
		}
	}

	// QoS settings
	if create.QoS != "" {
		apiPartition.Qos = &create.QoS
	}

	// State
	if create.State != "" {
		states := []api.V0043PartitionState{api.V0043PartitionState(create.State)}
		apiPartition.State = &states
	}

	// TRES and billing
	if create.TresStr != "" {
		apiPartition.TresStr = &create.TresStr
	}
	if create.BillingWeightStr != "" {
		apiPartition.BillingWeights = &create.BillingWeightStr
	}

	// Preempt mode
	if len(create.PreemptMode) > 0 {
		preemptModes := make([]api.V0043PartitionPreemptMode, len(create.PreemptMode))
		for i, mode := range create.PreemptMode {
			preemptModes[i] = api.V0043PartitionPreemptMode(mode)
		}
		apiPartition.PreemptMode = &preemptModes
	}

	// Job defaults
	if len(create.JobDefaults) > 0 {
		jobDefaults := make([]string, 0, len(create.JobDefaults))
		for key, value := range create.JobDefaults {
			jobDefaults = append(jobDefaults, key+"="+value)
		}
		apiPartition.JobDefaults = &jobDefaults
	}

	// Timeout settings
	if create.ResumeTimeout > 0 {
		setTrue := true
		timeout := int32(create.ResumeTimeout)
		apiPartition.ResumeTimeout = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &timeout,
		}
	}
	if create.SuspendTime > 0 {
		setTrue := true
		suspendTime := int32(create.SuspendTime)
		apiPartition.SuspendTime = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &suspendTime,
		}
	}
	if create.SuspendTimeout > 0 {
		setTrue := true
		timeout := int32(create.SuspendTimeout)
		apiPartition.SuspendTimeout = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &timeout,
		}
	}

	// Boolean flags
	var flags []api.V0043PartitionInfoFlags
	if create.Hidden {
		flags = append(flags, api.V0043PartitionInfoFlagsHIDDEN)
	}
	if create.RootOnly {
		flags = append(flags, api.V0043PartitionInfoFlagsROOTONLY)
	}
	if create.ReqResv {
		flags = append(flags, api.V0043PartitionInfoFlagsREQRESV)
	}
	if create.LLN {
		flags = append(flags, api.V0043PartitionInfoFlagsLLN)
	}
	if create.ExclusiveUser {
		flags = append(flags, api.V0043PartitionInfoFlagsEXCLUSIVEUSER)
	}
	if len(flags) > 0 {
		apiPartition.Flags = &flags
	}

	return apiPartition, nil
}

// convertCommonPartitionUpdateToAPI converts common PartitionUpdate to v0.0.43 API format
func (a *PartitionAdapter) convertCommonPartitionUpdateToAPI(existing *types.Partition, update *types.PartitionUpdate) (*api.V0043PartitionInfo, error) {
	apiPartition := &api.V0043PartitionInfo{}

	// Always include the partition name for updates
	apiPartition.Name = &existing.Name

	// Apply updates to fields
	allocNodes := existing.AllocNodes
	if update.AllocNodes != nil {
		allocNodes = *update.AllocNodes
	}
	if allocNodes != "" {
		apiPartition.AllocNodes = &allocNodes
	}

	// Allow/Deny lists
	allowAccounts := existing.AllowAccounts
	if len(update.AllowAccounts) > 0 {
		allowAccounts = update.AllowAccounts
	}
	if len(allowAccounts) > 0 {
		apiPartition.AllowAccounts = &allowAccounts
	}

	allowAllocNodes := existing.AllowAllocNodes
	if update.AllowAllocNodes != nil {
		allowAllocNodes = *update.AllowAllocNodes
	}
	if allowAllocNodes != "" {
		apiPartition.AllowAllocNodes = &allowAllocNodes
	}

	allowGroups := existing.AllowGroups
	if len(update.AllowGroups) > 0 {
		allowGroups = update.AllowGroups
	}
	if len(allowGroups) > 0 {
		apiPartition.AllowGroups = &allowGroups
	}

	allowQoS := existing.AllowQoS
	if len(update.AllowQoS) > 0 {
		allowQoS = update.AllowQoS
	}
	if len(allowQoS) > 0 {
		apiPartition.AllowQos = &allowQoS
	}

	denyAccounts := existing.DenyAccounts
	if len(update.DenyAccounts) > 0 {
		denyAccounts = update.DenyAccounts
	}
	if len(denyAccounts) > 0 {
		apiPartition.DenyAccounts = &denyAccounts
	}

	denyQoS := existing.DenyQoS
	if len(update.DenyQoS) > 0 {
		denyQoS = update.DenyQoS
	}
	if len(denyQoS) > 0 {
		apiPartition.DenyQos = &denyQoS
	}

	// Memory settings
	defaultMemPerCPU := existing.DefaultMemPerCPU
	if update.DefaultMemPerCPU != nil {
		defaultMemPerCPU = *update.DefaultMemPerCPU
	}
	if defaultMemPerCPU > 0 {
		setTrue := true
		apiPartition.DefaultMemoryPerCpu = &api.V0043Uint64NoValStruct{
			Set:    &setTrue,
			Number: &defaultMemPerCPU,
		}
	}

	defaultMemPerNode := existing.DefaultMemPerNode
	if update.DefaultMemPerNode != nil {
		defaultMemPerNode = *update.DefaultMemPerNode
	}
	if defaultMemPerNode > 0 {
		setTrue := true
		apiPartition.DefaultMemoryPerNode = &api.V0043Uint64NoValStruct{
			Set:    &setTrue,
			Number: &defaultMemPerNode,
		}
	}

	// Time settings
	defaultTime := existing.DefaultTime
	if update.DefaultTime != nil {
		defaultTime = *update.DefaultTime
	}
	if defaultTime > 0 {
		setTrue := true
		time := int32(defaultTime)
		apiPartition.DefaultTime = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &time,
		}
	}

	maxTime := existing.MaxTime
	if update.MaxTime != nil {
		maxTime = *update.MaxTime
	}
	if maxTime > 0 {
		setTrue := true
		time := int32(maxTime)
		apiPartition.MaxTime = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &time,
		}
	}

	// Node settings
	maxNodes := existing.MaxNodes
	if update.MaxNodes != nil {
		maxNodes = *update.MaxNodes
	}
	if maxNodes > 0 {
		setTrue := true
		nodes := int32(maxNodes)
		apiPartition.MaxNodes = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &nodes,
		}
	}

	minNodes := existing.MinNodes
	if update.MinNodes != nil {
		minNodes = *update.MinNodes
	}
	if minNodes > 0 {
		setTrue := true
		nodes := int32(minNodes)
		apiPartition.MinNodes = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &nodes,
		}
	}

	// Priority
	priority := existing.Priority
	if update.Priority != nil {
		priority = *update.Priority
	}
	if priority > 0 {
		setTrue := true
		priorityInt32 := int32(priority)
		apiPartition.Priority = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &priorityInt32,
		}
	}

	// State
	state := existing.State
	if update.State != nil {
		state = *update.State
	}
	if state != "" {
		states := []api.V0043PartitionState{api.V0043PartitionState(state)}
		apiPartition.State = &states
	}

	// Boolean flags - merge existing with updates
	var flags []api.V0043PartitionInfoFlags
	
	hidden := existing.Hidden
	if update.Hidden != nil {
		hidden = *update.Hidden
	}
	if hidden {
		flags = append(flags, api.V0043PartitionInfoFlagsHIDDEN)
	}

	rootOnly := existing.RootOnly
	if update.RootOnly != nil {
		rootOnly = *update.RootOnly
	}
	if rootOnly {
		flags = append(flags, api.V0043PartitionInfoFlagsROOTONLY)
	}

	reqResv := existing.ReqResv
	if update.ReqResv != nil {
		reqResv = *update.ReqResv
	}
	if reqResv {
		flags = append(flags, api.V0043PartitionInfoFlagsREQRESV)
	}

	lln := existing.LLN
	if update.LLN != nil {
		lln = *update.LLN
	}
	if lln {
		flags = append(flags, api.V0043PartitionInfoFlagsLLN)
	}

	exclusiveUser := existing.ExclusiveUser
	if update.ExclusiveUser != nil {
		exclusiveUser = *update.ExclusiveUser
	}
	if exclusiveUser {
		flags = append(flags, api.V0043PartitionInfoFlagsEXCLUSIVEUSER)
	}

	if len(flags) > 0 {
		apiPartition.Flags = &flags
	}

	return apiPartition, nil
}