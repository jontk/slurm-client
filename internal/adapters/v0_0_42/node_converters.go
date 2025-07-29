package v0_0_42

import (
	"strings"
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
		node.Hostname = *apiNode.Hostname
	}
	if apiNode.Address != nil {
		node.Address = *apiNode.Address
	}
	if apiNode.Architecture != nil {
		node.Architecture = *apiNode.Architecture
	}
	if apiNode.OperatingSystem != nil {
		node.OperatingSystem = *apiNode.OperatingSystem
	}

	// State
	if apiNode.State != nil {
		node.State = a.convertNodeStatesToString(apiNode.State)
	}
	if apiNode.StateFlags != nil {
		node.StateFlags = a.convertNodeStateFlagsToStrings(apiNode.StateFlags)
	}
	if apiNode.Reason != nil {
		node.Reason = *apiNode.Reason
	}
	if apiNode.ReasonTime != nil {
		node.ReasonTime = time.Unix(int64(apiNode.ReasonTime.Number), 0)
	}
	if apiNode.ReasonSetByUser != nil {
		node.ReasonUID = uint32(apiNode.ReasonSetByUser.Number)
	}

	// Resources
	if apiNode.Cpus != nil {
		node.CPUs = uint32(*apiNode.Cpus)
	}
	if apiNode.EffectiveCpus != nil {
		node.CPUsEffective = uint32(*apiNode.EffectiveCpus)
	}
	if apiNode.AllocCpus != nil {
		node.CPUsAllocated = uint32(*apiNode.AllocCpus)
	}
	if apiNode.AllocIdleCpus != nil {
		node.CPUsIdle = uint32(*apiNode.AllocIdleCpus)
	}
	if apiNode.Cores != nil {
		node.Cores = uint32(*apiNode.Cores)
	}
	if apiNode.Boards != nil {
		node.Boards = uint32(*apiNode.Boards)
	}
	if apiNode.Sockets != nil {
		node.Sockets = uint32(*apiNode.Sockets)
	}
	if apiNode.ThreadsPerCore != nil {
		node.ThreadsPerCore = uint32(*apiNode.ThreadsPerCore)
	}

	// Memory
	if apiNode.RealMemory != nil {
		node.RealMemory = uint64(*apiNode.RealMemory)
	}
	if apiNode.AllocMemory != nil {
		node.AllocMemory = uint64(*apiNode.AllocMemory)
	}
	if apiNode.FreeMem != nil {
		node.FreeMemory = apiNode.FreeMem.Number
	}

	// Features
	if apiNode.Features != nil && apiNode.Features.String != nil {
		node.Features = strings.Split(*apiNode.Features.String, ",")
	}
	if apiNode.ActiveFeatures != nil && apiNode.ActiveFeatures.String != nil {
		node.ActiveFeatures = strings.Split(*apiNode.ActiveFeatures.String, ",")
	}

	// Partitions
	if apiNode.Partitions != nil && apiNode.Partitions.String != nil {
		node.Partitions = strings.Split(*apiNode.Partitions.String, ",")
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
	if apiNode.BootTime != nil && apiNode.BootTime.Number > 0 {
		node.BootTime = time.Unix(int64(apiNode.BootTime.Number), 0)
	}

	// Last busy time
	if apiNode.LastBusy != nil && apiNode.LastBusy.Number > 0 {
		node.LastBusy = time.Unix(int64(apiNode.LastBusy.Number), 0)
	}

	// Slurm versions
	if apiNode.SlurmdVersion != nil {
		node.SlurmdVersion = *apiNode.SlurmdVersion
	}
	if apiNode.SlurmVersion != nil {
		node.Version = *apiNode.SlurmVersion
	}

	// Network
	if apiNode.Port != nil {
		node.Port = uint32(*apiNode.Port)
	}
	if apiNode.BurstbufferNetworkAddress != nil {
		node.BurstBufferNetworkAddress = *apiNode.BurstbufferNetworkAddress
	}

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
	if apiNode.Extra != nil {
		node.Extra = *apiNode.Extra
	}

	// Weight
	if apiNode.Weight != nil {
		node.Weight = uint32(apiNode.Weight.Number)
	}

	// CPU load
	if apiNode.CpuLoad != nil {
		node.CPULoad = float64(*apiNode.CpuLoad) / 100.0 // Convert from centipercent
	}

	// Temporary disk
	if apiNode.TmpDisk != nil {
		node.TmpDisk = uint64(*apiNode.TmpDisk)
	}

	// Energy
	if apiNode.Energy != nil {
		// Store energy data as extra info
		if node.Extra == "" {
			node.Extra = "energy_available"
		}
	}

	// Cloud instance info
	if apiNode.InstanceId != nil {
		node.CloudInstanceID = *apiNode.InstanceId
	}
	if apiNode.InstanceType != nil {
		node.CloudInstanceType = *apiNode.InstanceType
	}

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
	apiReq := &api.SlurmV0042PostNodeJSONRequestBody{
		Nodes: &[]api.V0042UpdateNodeMsg{
			{
				NodeNames: &[]string{nodeName},
			},
		},
	}

	update := &(*apiReq.Nodes)[0]

	// State update
	if req.State != nil {
		state := api.V0042NodeStates{*req.State}
		update.State = &state
	}

	// Reason for state change
	if req.Reason != nil {
		update.Reason = req.Reason
		
		// Set reason user if provided
		if req.ReasonUID != nil {
			reasonUID := api.V0042Uint32NoValStruct{
				Set:    true,
				Number: uint64(*req.ReasonUID),
			}
			update.ReasonUid = &reasonUID
		}
	}

	// Features
	if req.Features != nil {
		featuresStr := strings.Join(req.Features, ",")
		update.Features = &featuresStr
	}

	// Weight
	if req.Weight != nil {
		weight := api.V0042Uint32NoValStruct{
			Set:    true,
			Number: uint64(*req.Weight),
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

	// Extra
	if req.Extra != nil {
		update.Extra = req.Extra
	}

	return apiReq, nil
}