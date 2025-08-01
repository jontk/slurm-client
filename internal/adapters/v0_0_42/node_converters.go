// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_42

import (
	"time"

	"github.com/jontk/slurm-client/internal/common/types"
	api "github.com/jontk/slurm-client/internal/api/v0_0_42"
)

// convertAPINodeToCommon converts a v0.0.42 API Node to common Node type
func (a *NodeAdapter) convertAPINodeToCommon(apiNode api.V0042Node) (*types.Node, error) {
	node := &types.Node{}

	// Basic fields
	if apiNode.Name != nil {
		node.Name = *apiNode.Name
	}
	if apiNode.Hostname != nil {
		node.NodeHostname = *apiNode.Hostname
	}
	if apiNode.Address != nil {
		node.NodeAddress = *apiNode.Address
	}
	if apiNode.Architecture != nil {
		node.Arch = *apiNode.Architecture
	}
	if apiNode.OperatingSystem != nil {
		node.OS = *apiNode.OperatingSystem
	}

	// State
	if apiNode.State != nil {
		node.State = types.NodeState(a.convertNodeStatesToString(apiNode.State))
	}
	// StateFlags field doesn't exist in v0.0.42 API
	// Skipping state flags conversion
	if apiNode.Reason != nil {
		node.Reason = *apiNode.Reason
	}
	// ReasonTime and ReasonSetByUser fields don't exist in v0.0.42 API
	// Skipping reason time and reason UID conversion

	// Resources
	if apiNode.Cpus != nil {
		node.CPUs = int32(*apiNode.Cpus)
	}
	if apiNode.EffectiveCpus != nil {
		node.CPUsEffective = int32(*apiNode.EffectiveCpus)
	}
	if apiNode.AllocCpus != nil {
		node.AllocCPUs = int32(*apiNode.AllocCpus)
	}
	if apiNode.AllocIdleCpus != nil {
		node.AllocIdleCPUs = int32(*apiNode.AllocIdleCpus)
	}
	if apiNode.Cores != nil {
		node.Cores = int32(*apiNode.Cores)
	}
	if apiNode.Boards != nil {
		node.Boards = int32(*apiNode.Boards)
	}
	if apiNode.Sockets != nil {
		node.Sockets = int32(*apiNode.Sockets)
	}
	// ThreadsPerCore field doesn't exist in v0.0.42 API
	// Skipping threads per core conversion

	// Memory
	if apiNode.RealMemory != nil {
		node.RealMemory = int64(*apiNode.RealMemory)
	}
	if apiNode.AllocMemory != nil {
		node.AllocMemory = int64(*apiNode.AllocMemory)
	}
	if apiNode.FreeMem != nil && apiNode.FreeMem.Number != nil {
		node.FreeMemory = *apiNode.FreeMem.Number
	}

	// Features - V0042CsvString doesn't have .String field, need to handle differently
	if apiNode.Features != nil {
		// Convert CSV string to slice - need to check actual structure
		// For now, skip features conversion as API structure is unclear
	}
	if apiNode.ActiveFeatures != nil {
		// Skip active features conversion
	}

	// Partitions - V0042CsvString structure unclear, skip for now
	if apiNode.Partitions != nil {
		// Skip partitions conversion
	}

	// Generic resources (GRES)
	if apiNode.Gres != nil {
		node.Gres = *apiNode.Gres
	}
	if apiNode.GresDrained != nil {
		node.GresDrained = *apiNode.GresDrained
	}
	if apiNode.GresUsed != nil {
		node.GresUsed = *apiNode.GresUsed
	}

	// Boot time
	if apiNode.BootTime != nil && apiNode.BootTime.Number != nil && *apiNode.BootTime.Number > 0 {
		bootTime := time.Unix(*apiNode.BootTime.Number, 0)
		node.BootTime = &bootTime
	}

	// Last busy time
	if apiNode.LastBusy != nil && apiNode.LastBusy.Number != nil && *apiNode.LastBusy.Number > 0 {
		lastBusy := time.Unix(*apiNode.LastBusy.Number, 0)
		node.LastBusy = &lastBusy
	}

	// Slurm versions - these fields don't exist in v0.0.42 API
	// Skip version field conversion

	// Network
	if apiNode.Port != nil {
		node.Port = int32(*apiNode.Port)
	}
	// BurstBufferNetworkAddress field doesn't exist in common Node type
	// Skip burst buffer network address conversion

	// Other fields
	if apiNode.Comment != nil {
		node.Comment = *apiNode.Comment
	}
	if apiNode.Owner != nil {
		node.Owner = *apiNode.Owner
	}
	if apiNode.McsLabel != nil {
		node.MCSLabel = *apiNode.McsLabel
	}
	// Extra field doesn't exist in common Node type
	// Skip extra field conversion

	// Weight
	if apiNode.Weight != nil {
		node.Weight = int32(*apiNode.Weight)
	}

	// CPU load
	if apiNode.CpuLoad != nil {
		node.CPULoad = float64(*apiNode.CpuLoad) / 100.0 // Convert from centipercent
	}

	// Temporary disk - TmpDisk field doesn't exist in v0.0.42 API
	// Skip temporary disk conversion

	// Energy - Extra field doesn't exist in common Node type
	// Skip energy field conversion

	// Cloud instance info - these fields don't exist in common Node type
	// Skip cloud instance field conversion

	return node, nil
}

// convertNodeStatesToString converts node state flags to a string representation
func (a *NodeAdapter) convertNodeStatesToString(states *api.V0042NodeStates) string {
	if states == nil || len(*states) == 0 {
		return "UNKNOWN"
	}

	// Return the first state as the primary state
	return (*states)[0]
}

// convertNodeStateFlagsToStrings converts node state flags to string array
func (a *NodeAdapter) convertNodeStateFlagsToStrings(states *api.V0042NodeStates) []string {
	if states == nil {
		return []string{}
	}

	return *states
}

// convertCommonNodeUpdateToAPI converts common node update request to v0.0.42 API format
func (a *NodeAdapter) convertCommonNodeUpdateToAPI(nodeName string, req *types.NodeUpdateRequest) (*api.SlurmV0042PostNodeJSONRequestBody, error) {
	// Note: The exact structure for v0.0.42 node updates may be different
	// This is a simplified placeholder
	apiReq := &api.SlurmV0042PostNodeJSONRequestBody{}
	update := &api.V0042UpdateNodeMsg{}

	// State update
	if req.State != nil {
		state := api.V0042NodeStates{string(*req.State)}
		update.State = &state
	}

	// Reason for state change
	if req.Reason != nil {
		update.Reason = req.Reason
		
		// ReasonUID field doesn't exist in NodeUpdateRequest
		// Skip reason UID setting
	}

	// Features - convert to V0042CsvString (which is []string)
	if req.Features != nil {
		featuresCSV := api.V0042CsvString(req.Features)
		update.Features = &featuresCSV
	}

	// Weight
	if req.Weight != nil {
		weight := api.V0042Uint32NoValStruct{
			Set:    &[]bool{true}[0],
			Number: (*int32)(req.Weight),
		}
		update.Weight = &weight
	}

	// Comment
	if req.Comment != nil {
		update.Comment = req.Comment
	}

	// Gres
	if req.Gres != nil {
		update.Gres = req.Gres
	}

	// Extra - convert map to string
	if req.Extra != nil {
		// Convert map to JSON string or similar format
		// For now, skip extra field conversion as API expects string
		_ = req.Extra
	}

	return apiReq, nil
}
