package v0_0_43

import (
	"time"

	"github.com/jontk/slurm-client/internal/common/types"
	api "github.com/jontk/slurm-client/internal/api/v0_0_43"
)

// convertAPINodeToCommon converts a v0.0.43 API Node to common Node type
func (a *NodeAdapter) convertAPINodeToCommon(apiNode api.V0043Node) (*types.Node, error) {
	node := &types.Node{}

	// Basic fields
	if apiNode.Name != nil {
		node.Name = *apiNode.Name
	}
	if apiNode.Architecture != nil {
		node.Arch = *apiNode.Architecture
	}
	if apiNode.BurstBufferNetworkAddress != nil {
		node.BcastAddress = *apiNode.BurstBufferNetworkAddress
	}
	if apiNode.Boards != nil {
		node.Boards = *apiNode.Boards
	}
	if apiNode.ClusterName != nil {
		node.ClusterName = *apiNode.ClusterName
	}
	if apiNode.Comment != nil {
		node.Comment = *apiNode.Comment
	}

	// CPU and core information
	if apiNode.Cores != nil {
		node.Cores = *apiNode.Cores
	}
	if apiNode.CoreSpecCount != nil {
		node.CoreSpecCount = *apiNode.CoreSpecCount
	}
	if apiNode.CpuBinding != nil {
		node.CPUBinding = *apiNode.CpuBinding
	}
	if apiNode.CpuLoad != nil && apiNode.CpuLoad.Number != nil {
		node.CPULoad = *apiNode.CpuLoad.Number
	}
	if apiNode.Cpus != nil {
		node.CPUs = *apiNode.Cpus
	}
	if apiNode.CpusEffective != nil {
		node.CPUsEffective = *apiNode.CpusEffective
	}
	if apiNode.ThreadsPerCore != nil {
		node.ThreadsPerCore = *apiNode.ThreadsPerCore
	}
	if apiNode.Sockets != nil {
		node.Sockets = *apiNode.Sockets
	}

	// Memory information
	if apiNode.FreeMemory != nil && apiNode.FreeMemory.Number != nil {
		node.FreeMemory = *apiNode.FreeMemory.Number
	}
	if apiNode.Memory != nil {
		node.Memory = *apiNode.Memory
	}
	if apiNode.MemorySpecLimit != nil {
		node.MemorySpecLimit = *apiNode.MemorySpecLimit
	}
	if apiNode.RealMemory != nil {
		node.RealMemory = *apiNode.RealMemory
	}
	if apiNode.AllocMemory != nil {
		node.AllocMemory = *apiNode.AllocMemory
	}

	// Features and GRES
	if apiNode.Features != nil {
		node.Features = *apiNode.Features
	}
	if apiNode.ActiveFeatures != nil {
		node.ActiveFeatures = *apiNode.ActiveFeatures
	}
	if apiNode.Gres != nil {
		node.Gres = *apiNode.Gres
	}
	if apiNode.GresDrained != nil {
		node.GresDrained = *apiNode.GresDrained
	}
	if apiNode.GresUsed != nil {
		node.GresUsed = *apiNode.GresUsed
	}

	// Time fields
	if apiNode.BootTime != nil && apiNode.BootTime.Number != nil && *apiNode.BootTime.Number > 0 {
		bootTime := time.Unix(*apiNode.BootTime.Number, 0)
		node.BootTime = &bootTime
	}
	if apiNode.LastBusy != nil && apiNode.LastBusy.Number != nil && *apiNode.LastBusy.Number > 0 {
		lastBusy := time.Unix(*apiNode.LastBusy.Number, 0)
		node.LastBusy = &lastBusy
	}
	if apiNode.ReasonChangedAt != nil && apiNode.ReasonChangedAt.Number != nil && *apiNode.ReasonChangedAt.Number > 0 {
		reasonChanged := time.Unix(*apiNode.ReasonChangedAt.Number, 0)
		node.ReasonChangedAt = &reasonChanged
	}
	if apiNode.ResumeAfter != nil && apiNode.ResumeAfter.Number != nil && *apiNode.ResumeAfter.Number > 0 {
		resumeAfter := time.Unix(*apiNode.ResumeAfter.Number, 0)
		node.ResumeAfter = &resumeAfter
	}
	if apiNode.SlurmdStartTime != nil && apiNode.SlurmdStartTime.Number != nil && *apiNode.SlurmdStartTime.Number > 0 {
		slurmdStart := time.Unix(*apiNode.SlurmdStartTime.Number, 0)
		node.SlurmdStartTime = &slurmdStart
	}

	// Network information
	if apiNode.NodeAddress != nil {
		node.NodeAddress = *apiNode.NodeAddress
	}
	if apiNode.NodeHostname != nil {
		node.NodeHostname = *apiNode.NodeHostname
	}
	if apiNode.Port != nil {
		node.Port = *apiNode.Port
	}

	// System information
	if apiNode.OperatingSystem != nil {
		node.OS = *apiNode.OperatingSystem
	}
	if apiNode.Owner != nil {
		node.Owner = *apiNode.Owner
	}
	if apiNode.Version != nil {
		node.Version = *apiNode.Version
	}
	if apiNode.Weight != nil {
		node.Weight = *apiNode.Weight
	}

	// Partitions
	if apiNode.Partitions != nil {
		node.Partitions = *apiNode.Partitions
	}

	// State information
	if apiNode.State != nil && len(*apiNode.State) > 0 {
		node.State = types.NodeState((*apiNode.State)[0])
	}
	if apiNode.StateFlags != nil && len(*apiNode.StateFlags) > 0 {
		stateFlags := make([]string, len(*apiNode.StateFlags))
		for i, flag := range *apiNode.StateFlags {
			stateFlags[i] = string(flag)
		}
		node.StateFlags = stateFlags
	}
	if apiNode.NextStateAfterReboot != nil && len(*apiNode.NextStateAfterReboot) > 0 {
		node.NextStateAfterReboot = types.NodeState((*apiNode.NextStateAfterReboot)[0])
	}
	if apiNode.NextStateAfterRebootFlags != nil && len(*apiNode.NextStateAfterRebootFlags) > 0 {
		nextStateFlags := make([]string, len(*apiNode.NextStateAfterRebootFlags))
		for i, flag := range *apiNode.NextStateAfterRebootFlags {
			nextStateFlags[i] = string(flag)
		}
		node.NextStateAfterRebootFlags = nextStateFlags
	}

	// Reason information
	if apiNode.Reason != nil {
		node.Reason = *apiNode.Reason
	}
	if apiNode.ReasonSetByUser != nil {
		node.ReasonSetByUser = *apiNode.ReasonSetByUser
	}

	// Resource information
	if apiNode.TmpDisk != nil {
		node.TmpDisk = *apiNode.TmpDisk
	}
	if apiNode.TresUsed != nil {
		node.TresUsed = *apiNode.TresUsed
	}
	if apiNode.TresFmtStr != nil {
		node.TresFmtStr = *apiNode.TresFmtStr
	}

	// Allocation information
	if apiNode.AllocCpus != nil {
		node.AllocCPUs = *apiNode.AllocCpus
	}
	if apiNode.AllocIdleCpus != nil {
		node.AllocIdleCPUs = *apiNode.AllocIdleCpus
	}

	// MCS Label
	if apiNode.McsLabel != nil {
		node.MCSLabel = *apiNode.McsLabel
	}

	// Energy information (if available)
	if apiNode.Energy != nil {
		energy := &types.NodeEnergy{}
		if apiNode.Energy.AverageWatts != nil && apiNode.Energy.AverageWatts.Number != nil {
			energy.AveWatts = *apiNode.Energy.AverageWatts.Number
		}
		if apiNode.Energy.BaseConsumedEnergy != nil && apiNode.Energy.BaseConsumedEnergy.Number != nil {
			energy.BaseConsumedEnergy = *apiNode.Energy.BaseConsumedEnergy.Number
		}
		if apiNode.Energy.ConsumedEnergy != nil && apiNode.Energy.ConsumedEnergy.Number != nil {
			energy.ConsumedEnergy = *apiNode.Energy.ConsumedEnergy.Number
		}
		if apiNode.Energy.CurrentWatts != nil && apiNode.Energy.CurrentWatts.Number != nil {
			energy.CurrentWatts = *apiNode.Energy.CurrentWatts.Number
		}
		if apiNode.Energy.LastCollected != nil && apiNode.Energy.LastCollected.Number != nil {
			energy.LastCollected = time.Unix(*apiNode.Energy.LastCollected.Number, 0)
		}
		node.Energy = energy
	}

	return node, nil
}

// convertCommonNodeUpdateToAPI converts common NodeUpdate to v0.0.43 API format
func (a *NodeAdapter) convertCommonNodeUpdateToAPI(existing *types.Node, update *types.NodeUpdate) (*api.V0043Node, error) {
	apiNode := &api.V0043Node{}

	// Always include the node name for updates
	apiNode.Name = &existing.Name

	// Apply updates to fields
	comment := existing.Comment
	if update.Comment != nil {
		comment = *update.Comment
	}
	if comment != "" {
		apiNode.Comment = &comment
	}

	// CPU binding
	cpuBinding := existing.CPUBinding
	if update.CPUBinding != nil {
		cpuBinding = *update.CPUBinding
	}
	if cpuBinding != 0 {
		apiNode.CpuBinding = &cpuBinding
	}

	// Features
	features := existing.Features
	if len(update.Features) > 0 {
		features = update.Features
	}
	if len(features) > 0 {
		apiNode.Features = &features
	}

	// Active features
	activeFeatures := existing.ActiveFeatures
	if len(update.ActiveFeatures) > 0 {
		activeFeatures = update.ActiveFeatures
	}
	if len(activeFeatures) > 0 {
		apiNode.ActiveFeatures = &activeFeatures
	}

	// GRES
	gres := existing.Gres
	if update.Gres != nil {
		gres = *update.Gres
	}
	if gres != "" {
		apiNode.Gres = &gres
	}

	// Next state after reboot
	if update.NextStateAfterReboot != nil {
		nextStates := []api.V0043NodeState{api.V0043NodeState(*update.NextStateAfterReboot)}
		apiNode.NextStateAfterReboot = &nextStates
	}

	// Reason
	reason := existing.Reason
	if update.Reason != nil {
		reason = *update.Reason
	}
	if reason != "" {
		apiNode.Reason = &reason
	}

	// Resume after
	if update.ResumeAfter != nil {
		setTrue := true
		resumeAfter := update.ResumeAfter.Unix()
		apiNode.ResumeAfter = &api.V0043Uint64NoValStruct{
			Set:    &setTrue,
			Number: &resumeAfter,
		}
	}

	// State
	if update.State != nil {
		states := []api.V0043NodeState{api.V0043NodeState(*update.State)}
		apiNode.State = &states
	}

	// Weight
	weight := existing.Weight
	if update.Weight != nil {
		weight = *update.Weight
	}
	if weight != 0 {
		apiNode.Weight = &weight
	}

	return apiNode, nil
}