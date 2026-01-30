// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_43

import (
	"strings"
	"time"

	api "github.com/jontk/slurm-client/internal/api/v0_0_43"
	"github.com/jontk/slurm-client/internal/common/types"
)

// convertAPINodeToCommon converts a v0.0.43 API Node to common Node type
func (a *NodeAdapter) convertAPINodeToCommon(apiNode api.V0043Node) *types.Node {
	node := &types.Node{}

	// Basic fields
	if apiNode.Name != nil {
		node.Name = *apiNode.Name
	}
	if apiNode.Architecture != nil {
		node.Arch = *apiNode.Architecture
	}
	if apiNode.BurstbufferNetworkAddress != nil {
		node.BcastAddress = *apiNode.BurstbufferNetworkAddress
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
	if apiNode.CpuBinding != nil {
		node.CPUBinding = *apiNode.CpuBinding
	}
	if apiNode.CpuLoad != nil {
		node.CPULoad = float64(*apiNode.CpuLoad)
	}
	if apiNode.Cpus != nil {
		node.CPUs = *apiNode.Cpus
	}
	if apiNode.EffectiveCpus != nil {
		node.CPUsEffective = *apiNode.EffectiveCpus
	}
	if apiNode.Sockets != nil {
		node.Sockets = *apiNode.Sockets
	}

	// Memory information
	if apiNode.RealMemory != nil {
		node.Memory = *apiNode.RealMemory
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
	if apiNode.Address != nil {
		node.NodeAddress = *apiNode.Address
	}
	if apiNode.Hostname != nil {
		node.NodeHostname = *apiNode.Hostname
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
	// SLURM API returns state as an array (e.g. ["IDLE", "DRAIN"])
	// Concatenate all states with "+" to preserve all flags (e.g. "IDLE+DRAIN")
	if apiNode.State != nil && len(*apiNode.State) > 0 {
		states := *apiNode.State
		if len(states) == 1 {
			node.State = types.NodeState(states[0])
		} else {
			// Join multiple states with "+" (e.g. "IDLE+DRAIN")
			stateStrings := make([]string, len(states))
			for i, s := range states {
				stateStrings[i] = string(s)
			}
			node.State = types.NodeState(strings.Join(stateStrings, "+"))
		}
	}
	// StateFlags not available in v0_0_43
	if apiNode.NextStateAfterReboot != nil && len(*apiNode.NextStateAfterReboot) > 0 {
		node.NextStateAfterReboot = types.NodeState((*apiNode.NextStateAfterReboot)[0])
	}
	// NextStateAfterRebootFlags not available in v0_0_43

	// Reason information
	if apiNode.Reason != nil {
		node.Reason = *apiNode.Reason
	}
	if apiNode.ReasonSetByUser != nil {
		node.ReasonSetByUser = *apiNode.ReasonSetByUser
	}

	// Resource information
	if apiNode.TresUsed != nil {
		node.TresUsed = *apiNode.TresUsed
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
		if apiNode.Energy.AverageWatts != nil {
			energy.AveWatts = int64(*apiNode.Energy.AverageWatts)
		}
		if apiNode.Energy.BaseConsumedEnergy != nil {
			energy.BaseConsumedEnergy = *apiNode.Energy.BaseConsumedEnergy
		}
		if apiNode.Energy.ConsumedEnergy != nil {
			energy.ConsumedEnergy = *apiNode.Energy.ConsumedEnergy
		}
		if apiNode.Energy.CurrentWatts != nil && apiNode.Energy.CurrentWatts.Number != nil {
			energy.CurrentWatts = int64(*apiNode.Energy.CurrentWatts.Number)
		}
		if apiNode.Energy.LastCollected != nil {
			energy.LastCollected = time.Unix(*apiNode.Energy.LastCollected, 0)
		}
		node.Energy = energy
	}

	return node
}

// convertCommonNodeUpdateToAPI converts common NodeUpdate to v0.0.43 API format
func (a *NodeAdapter) convertCommonNodeUpdateToAPI(existing *types.Node, update *types.NodeUpdate) *api.V0043Node {
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
		nextStates := []api.V0043NodeNextStateAfterReboot{api.V0043NodeNextStateAfterReboot(*update.NextStateAfterReboot)}
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

	return apiNode
}
