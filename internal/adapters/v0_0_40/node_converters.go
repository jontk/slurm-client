package v0_0_40

import (
	"strings"
	"time"

	"github.com/jontk/slurm-client/internal/common/types"
	api "github.com/jontk/slurm-client/internal/api/v0_0_40"
)

// convertAPINodeToCommon converts a v0.0.40 API Node to common Node type
func (a *NodeAdapter) convertAPINodeToCommon(apiNode api.V0040Node) (*types.Node, error) {
	node := &types.Node{}

	// Basic fields
	if apiNode.Name != nil {
		node.Name = *apiNode.Name
	}
	if apiNode.Hostname != nil {
		node.Hostname = *apiNode.Hostname
	}
	if apiNode.Address != nil {
		node.Address = *apiNode.Address
	}
	if apiNode.Architecture != nil {
		node.Architecture = *apiNode.Architecture
	}
	if apiNode.OperatingSystem != nil {
		node.OS = *apiNode.OperatingSystem
	}

	// State
	if apiNode.State != nil && len(*apiNode.State) > 0 {
		// Convert state array to single state
		stateStr := (*apiNode.State)[0]
		node.State = types.NodeState(strings.ToUpper(stateStr))
	}
	if apiNode.StateFlags != nil && len(*apiNode.StateFlags) > 0 {
		node.StateFlags = make([]string, len(*apiNode.StateFlags))
		for i, flag := range *apiNode.StateFlags {
			node.StateFlags[i] = string(flag)
		}
	}
	if apiNode.NextStateAfterReboot != nil && len(*apiNode.NextStateAfterReboot) > 0 {
		nextState := (*apiNode.NextStateAfterReboot)[0]
		node.NextState = types.NodeState(strings.ToUpper(nextState))
	}

	// Resources
	if apiNode.Cpus != nil {
		node.CPUs = *apiNode.Cpus
	}
	if apiNode.EffectiveCpus != nil {
		node.AllocCPUs = *apiNode.EffectiveCpus
	}
	if apiNode.RealMemory != nil {
		realMem := int64(*apiNode.RealMemory)
		node.RealMemory = &realMem
	}
	if apiNode.AllocMemory != nil {
		allocMem := int64(*apiNode.AllocMemory)
		node.AllocMemory = &allocMem
	}
	if apiNode.FreeMemory != nil && apiNode.FreeMemory.Number != nil {
		freeMem := int64(*apiNode.FreeMemory.Number)
		node.FreeMemory = &freeMem
	}
	if apiNode.TemporaryDisk != nil {
		tmpDisk := int64(*apiNode.TemporaryDisk)
		node.TmpDisk = &tmpDisk
	}

	// Sockets, cores, threads
	if apiNode.Sockets != nil {
		node.Sockets = *apiNode.Sockets
	}
	if apiNode.Cores != nil {
		node.CoresPerSocket = *apiNode.Cores
	}
	if apiNode.Threads != nil {
		node.ThreadsPerCore = *apiNode.Threads
	}

	// Features and GRES
	if apiNode.Features != nil && len(*apiNode.Features) > 0 {
		node.Features = strings.Join(*apiNode.Features, ",")
	}
	if apiNode.ActiveFeatures != nil && len(*apiNode.ActiveFeatures) > 0 {
		node.ActiveFeatures = strings.Join(*apiNode.ActiveFeatures, ",")
	}
	if apiNode.Gres != nil {
		node.GRES = *apiNode.Gres
	}

	// Partitions
	if apiNode.Partitions != nil {
		node.Partitions = *apiNode.Partitions
	}

	// Time fields
	if apiNode.BootTime != nil && apiNode.BootTime.Number != nil && *apiNode.BootTime.Number > 0 {
		node.BootTime = time.Unix(*apiNode.BootTime.Number, 0)
	}
	if apiNode.LastBusyTime != nil && apiNode.LastBusyTime.Number != nil && *apiNode.LastBusyTime.Number > 0 {
		lastBusy := time.Unix(*apiNode.LastBusyTime.Number, 0)
		node.LastBusy = &lastBusy
	}
	if apiNode.SlurmdStartTime != nil && apiNode.SlurmdStartTime.Number != nil && *apiNode.SlurmdStartTime.Number > 0 {
		slurmdStart := time.Unix(*apiNode.SlurmdStartTime.Number, 0)
		node.SlurmdStartTime = &slurmdStart
	}

	// Additional fields
	if apiNode.Reason != nil {
		node.Reason = *apiNode.Reason
	}
	if apiNode.ReasonSetTime != nil && apiNode.ReasonSetTime.Number != nil && *apiNode.ReasonSetTime.Number > 0 {
		reasonTime := time.Unix(*apiNode.ReasonSetTime.Number, 0)
		node.ReasonTime = &reasonTime
	}
	if apiNode.ReasonUid != nil && apiNode.ReasonUid.Number != nil {
		reasonUID := *apiNode.ReasonUid.Number
		node.ReasonUID = &reasonUID
	}

	// Network and power
	if apiNode.McsLabel != nil {
		node.MCSLabel = *apiNode.McsLabel
	}
	if apiNode.Power != nil {
		// Handle power management if needed
	}

	// Version
	if apiNode.Version != nil {
		node.Version = *apiNode.Version
	}

	// Weight
	if apiNode.Weight != nil && apiNode.Weight.Number != nil {
		weight := *apiNode.Weight.Number
		node.Weight = &weight
	}

	// Energy
	if apiNode.Energy != nil {
		if apiNode.Energy.CurrentWatts != nil && apiNode.Energy.CurrentWatts.Number != nil {
			currentPower := *apiNode.Energy.CurrentWatts.Number
			node.CurrentWatts = &currentPower
		}
		if apiNode.Energy.AverageWatts != nil {
			avgPower := *apiNode.Energy.AverageWatts
			node.AverageWatts = &avgPower
		}
		if apiNode.Energy.ConsumedEnergy != nil {
			consumed := *apiNode.Energy.ConsumedEnergy
			node.ConsumedEnergy = &consumed
		}
	}

	// Owner
	if apiNode.Owner != nil {
		node.Owner = *apiNode.Owner
	}

	// Comment
	if apiNode.Comment != nil {
		node.Comment = *apiNode.Comment
	}

	// Extra
	if apiNode.Extra != nil {
		node.Extra = *apiNode.Extra
	}

	return node, nil
}

// convertCommonNodeUpdateToAPI converts common NodeUpdate to v0.0.40 API format
func (a *NodeAdapter) convertCommonNodeUpdateToAPI(nodeName string, update *types.NodeUpdate) (*api.V0040Node, error) {
	apiNode := &api.V0040Node{}

	// Name (required)
	apiNode.Name = &nodeName

	// Apply updates
	if update.State != nil {
		state := []string{string(*update.State)}
		apiNode.State = &state
	}

	if update.Reason != nil {
		apiNode.Reason = update.Reason
		
		// Set reason time to current time
		now := time.Now().Unix()
		reasonTime := api.V0040Uint64NoVal{
			Set:    boolPtr(true),
			Number: &now,
		}
		apiNode.ReasonSetTime = &reasonTime
	}

	if update.Comment != nil {
		apiNode.Comment = update.Comment
	}

	if update.Features != nil {
		features := strings.Split(*update.Features, ",")
		apiNode.Features = &features
	}

	if update.GRES != nil {
		apiNode.Gres = update.GRES
	}

	if update.Weight != nil {
		weight := api.V0040Uint32NoVal{
			Set:    boolPtr(true),
			Number: update.Weight,
		}
		apiNode.Weight = &weight
	}

	return apiNode, nil
}