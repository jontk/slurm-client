package v0_0_41

import (
	"fmt"
	"strings"
	"time"

	"github.com/jontk/slurm-client/internal/common/types"
	api "github.com/jontk/slurm-client/internal/api/v0_0_41"
)

// convertAPINodeToCommon converts a v0.0.41 API Node to common Node type
func (a *NodeAdapter) convertAPINodeToCommon(apiNode interface{}) (*types.Node, error) {
	// Type assertion to handle the anonymous struct
	nodeData, ok := apiNode.(struct {
		Architecture           *string `json:"architecture,omitempty"`
		BurstbufferNetworkAddress *string `json:"burstbuffer_network_address,omitempty"`
		Boards                 *int32  `json:"boards,omitempty"`
		BootTime               *api.V0041OpenapiNodesRespNodesBootTime `json:"boot_time,omitempty"`
		ClusterName            *string `json:"cluster_name,omitempty"`
		Cores                  *int32  `json:"cores,omitempty"`
		SpecializedCores       *int32  `json:"specialized_cores,omitempty"`
		CpuBinding             *int32  `json:"cpu_binding,omitempty"`
		CpuLoad                *api.V0041OpenapiNodesRespNodesCpuLoad `json:"cpu_load,omitempty"`
		SpecializedCpus        *string `json:"specialized_cpus,omitempty"`
		EffectiveCpus          *int32  `json:"effective_cpus,omitempty"`
		SpecializedMemory      *int64  `json:"specialized_memory,omitempty"`
		Energy                 *struct {
			AccumulatedEnergy *api.V0041OpenapiNodesRespNodesEnergyAccumulatedEnergy `json:"accumulated_energy,omitempty"`
			AveragePower      *int32 `json:"average_power,omitempty"`
			BasePower         *int32 `json:"base_power,omitempty"`
			CurrentPower      *api.V0041OpenapiNodesRespNodesEnergyCurrentPower `json:"current_power,omitempty"`
		} `json:"energy,omitempty"`
		ExternalSensors        *struct {
			AccumulatedEnergy *api.V0041OpenapiNodesRespNodesExternalSensorsAccumulatedEnergy `json:"accumulated_energy,omitempty"`
			CurrentPower      *api.V0041OpenapiNodesRespNodesExternalSensorsCurrentPower `json:"current_power,omitempty"`
			Temperature       *api.V0041OpenapiNodesRespNodesExternalSensorsTemperature `json:"temperature,omitempty"`
		} `json:"external_sensors,omitempty"`
		Extra                  *string `json:"extra,omitempty"`
		PowerCapabilities      *int32  `json:"power_capabilities,omitempty"`
		PowerConfigured        *int32  `json:"power_configured,omitempty"`
		Features               *string `json:"features,omitempty"`
		ActiveFeatures         *string `json:"active_features,omitempty"`
		Gres                   *string `json:"gres,omitempty"`
		GresDrained            *string `json:"gres_drained,omitempty"`
		GresUsed               *string `json:"gres_used,omitempty"`
		LastBusy               *api.V0041OpenapiNodesRespNodesLastBusy `json:"last_busy,omitempty"`
		McsLabel               *string `json:"mcs_label,omitempty"`
		SpecializedMemoryLimit *int64  `json:"specialized_memory_limit,omitempty"`
		Name                   *string `json:"name,omitempty"`
		NextStateAfterReboot   *[]api.V0041OpenapiNodesRespNodesNextStateAfterReboot `json:"next_state_after_reboot,omitempty"`
		Address                *string `json:"address,omitempty"`
		Hostname               *string `json:"hostname,omitempty"`
		State                  *[]api.V0041OpenapiNodesRespNodesState `json:"state,omitempty"`
		OperatingSystem        *string `json:"operating_system,omitempty"`
		Owner                  *string `json:"owner,omitempty"`
		Port                   *int32  `json:"port,omitempty"`
		Memory                 *struct {
			Allocated *int64 `json:"allocated,omitempty"`
			Real      *int64 `json:"real,omitempty"`
		} `json:"memory,omitempty"`
		Partitions             *string `json:"partitions,omitempty"`
		RebootRequested        *api.V0041OpenapiNodesRespNodesRebootRequested `json:"reboot_requested,omitempty"`
		Reason                 *string `json:"reason,omitempty"`
		ReasonChangedAt        *api.V0041OpenapiNodesRespNodesReasonChangedAt `json:"reason_changed_at,omitempty"`
		ReasonSetByUser        *string `json:"reason_set_by_user,omitempty"`
		ResumeAfter            *api.V0041OpenapiNodesRespNodesResumeAfter `json:"resume_after,omitempty"`
		ReservationName        *string `json:"reservation_name,omitempty"`
		Alloc                  *struct {
			IdleCpus   *int32 `json:"idle_cpus,omitempty"`
			AllocCpus  *int32 `json:"alloc_cpus,omitempty"`
			AllocMemory *int64 `json:"alloc_memory,omitempty"`
		} `json:"alloc,omitempty"`
		Size                   *int32  `json:"size,omitempty"`
		SlurmdStartTime        *api.V0041OpenapiNodesRespNodesSlurmdStartTime `json:"slurmd_start_time,omitempty"`
		Sockets                *int32  `json:"sockets,omitempty"`
		Threads                *int32  `json:"threads,omitempty"`
		TemporaryDisk          *int32  `json:"temporary_disk,omitempty"`
		Weight                 *int32  `json:"weight,omitempty"`
		Tres                   *string `json:"tres,omitempty"`
		TresUsed               *string `json:"tres_used,omitempty"`
		TresWeighted           *float64 `json:"tres_weighted,omitempty"`
		Version                *string `json:"version,omitempty"`
		AllocIdleCpus          *int32  `json:"alloc_idle_cpus,omitempty"`
		AllocCpus              *int32  `json:"alloc_cpus,omitempty"`
		IdleCpus               *int32  `json:"idle_cpus,omitempty"`
		RealMemory             *int64  `json:"real_memory,omitempty"`
		Comment                *string `json:"comment,omitempty"`
	})
	if !ok {
		return nil, fmt.Errorf("unexpected node data type")
	}

	node := &types.Node{}

	// Basic fields
	if nodeData.Name != nil {
		node.Name = *nodeData.Name
	}
	if nodeData.Hostname != nil {
		node.NodeHostname = *nodeData.Hostname
	}
	if nodeData.Address != nil {
		node.NodeAddr = *nodeData.Address
	}
	if nodeData.Port != nil {
		node.Port = uint16(*nodeData.Port)
	}

	// State
	if nodeData.State != nil && len(*nodeData.State) > 0 {
		// Convert state to string
		node.State = string((*nodeData.State)[0])
		
		// Parse state flags
		var stateFlags []string
		for _, s := range *nodeData.State {
			stateFlags = append(stateFlags, string(s))
		}
		node.StateFlags = strings.Join(stateFlags, ",")
	}

	// Resources
	if nodeData.Cpus != nil && nodeData.Cpus != nil {
		totalCpus := int32(0)
		if nodeData.EffectiveCpus != nil {
			totalCpus = *nodeData.EffectiveCpus
		} else if nodeData.Cores != nil && nodeData.Threads != nil && nodeData.Sockets != nil {
			totalCpus = *nodeData.Cores * *nodeData.Threads * *nodeData.Sockets
		}
		node.CPUs = uint32(totalCpus)
	}

	if nodeData.Cores != nil {
		node.CoresPerSocket = uint32(*nodeData.Cores)
	}
	if nodeData.Sockets != nil {
		node.Sockets = uint32(*nodeData.Sockets)
	}
	if nodeData.Threads != nil {
		node.ThreadsPerCore = uint32(*nodeData.Threads)
	}

	// Memory
	if nodeData.Memory != nil && nodeData.Memory.Real != nil {
		node.RealMemory = uint64(*nodeData.Memory.Real)
	} else if nodeData.RealMemory != nil {
		node.RealMemory = uint64(*nodeData.RealMemory)
	}

	if nodeData.Memory != nil && nodeData.Memory.Allocated != nil {
		node.AllocMemory = uint64(*nodeData.Memory.Allocated)
	}

	// CPU allocation info
	if nodeData.Alloc != nil {
		if nodeData.Alloc.AllocCpus != nil {
			node.AllocCPUs = uint32(*nodeData.Alloc.AllocCpus)
		}
		if nodeData.Alloc.IdleCpus != nil {
			node.IdleCPUs = uint32(*nodeData.Alloc.IdleCpus)
		}
	} else {
		if nodeData.AllocCpus != nil {
			node.AllocCPUs = uint32(*nodeData.AllocCpus)
		}
		if nodeData.IdleCpus != nil {
			node.IdleCPUs = uint32(*nodeData.IdleCpus)
		}
	}

	// Features
	if nodeData.Features != nil {
		node.Features = *nodeData.Features
	}
	if nodeData.ActiveFeatures != nil {
		node.ActiveFeatures = *nodeData.ActiveFeatures
	}

	// GRES
	if nodeData.Gres != nil {
		node.Gres = *nodeData.Gres
	}
	if nodeData.GresDrained != nil {
		node.GresDrained = *nodeData.GresDrained
	}
	if nodeData.GresUsed != nil {
		node.GresUsed = *nodeData.GresUsed
	}

	// Partitions
	if nodeData.Partitions != nil {
		node.Partitions = strings.Split(*nodeData.Partitions, ",")
	}

	// Architecture and OS
	if nodeData.Architecture != nil {
		node.Arch = *nodeData.Architecture
	}
	if nodeData.OperatingSystem != nil {
		node.OS = *nodeData.OperatingSystem
	}

	// Version
	if nodeData.Version != nil {
		node.Version = *nodeData.Version
	}

	// Boot time
	if nodeData.BootTime != nil && nodeData.BootTime.Number != nil {
		node.BootTime = time.Unix(*nodeData.BootTime.Number, 0)
	}

	// Slurmd start time
	if nodeData.SlurmdStartTime != nil && nodeData.SlurmdStartTime.Number != nil {
		node.SlurmdStartTime = time.Unix(*nodeData.SlurmdStartTime.Number, 0)
	}

	// Last busy time
	if nodeData.LastBusy != nil && nodeData.LastBusy.Number != nil {
		node.LastBusy = time.Unix(*nodeData.LastBusy.Number, 0)
	}

	// Reason and reason changed
	if nodeData.Reason != nil {
		node.Reason = *nodeData.Reason
	}
	if nodeData.ReasonChangedAt != nil && nodeData.ReasonChangedAt.Number != nil {
		node.ReasonTime = time.Unix(*nodeData.ReasonChangedAt.Number, 0)
	}
	if nodeData.ReasonSetByUser != nil {
		node.ReasonUID = *nodeData.ReasonSetByUser
	}

	// Temporary disk
	if nodeData.TemporaryDisk != nil {
		node.TmpDisk = uint32(*nodeData.TemporaryDisk)
	}

	// Weight
	if nodeData.Weight != nil {
		node.Weight = uint32(*nodeData.Weight)
	}

	// CPU load
	if nodeData.CpuLoad != nil && nodeData.CpuLoad.Number != nil {
		node.CPULoad = float64(*nodeData.CpuLoad.Number) / 100.0 // Convert from basis points
	}

	// TREs
	if nodeData.Tres != nil {
		node.TRESFmt = *nodeData.Tres
	}
	if nodeData.TresUsed != nil {
		node.TRESUsed = *nodeData.TresUsed
	}

	// Owner
	if nodeData.Owner != nil {
		node.Owner = *nodeData.Owner
	}

	// MCS label
	if nodeData.McsLabel != nil {
		node.MCSLabel = *nodeData.McsLabel
	}

	// Comment
	if nodeData.Comment != nil {
		node.Comment = *nodeData.Comment
	}

	// Extra
	if nodeData.Extra != nil {
		node.Extra = *nodeData.Extra
	}

	// Power information
	if nodeData.PowerCapabilities != nil {
		node.PowerCapWatts = uint32(*nodeData.PowerCapabilities)
	}
	if nodeData.PowerConfigured != nil {
		node.PowerConfiguredWatts = uint32(*nodeData.PowerConfigured)
	}

	// Energy information
	if nodeData.Energy != nil {
		if nodeData.Energy.CurrentPower != nil && nodeData.Energy.CurrentPower.Number != nil {
			node.CurrentWatts = uint32(*nodeData.Energy.CurrentPower.Number)
		}
		if nodeData.Energy.AveragePower != nil {
			node.AverageWatts = uint32(*nodeData.Energy.AveragePower)
		}
		if nodeData.Energy.AccumulatedEnergy != nil && nodeData.Energy.AccumulatedEnergy.Number != nil {
			node.ConsumedEnergy = uint64(*nodeData.Energy.AccumulatedEnergy.Number)
		}
	}

	// External sensors
	if nodeData.ExternalSensors != nil {
		if nodeData.ExternalSensors.CurrentPower != nil && nodeData.ExternalSensors.CurrentPower.Number != nil {
			node.ExtSensorsWatts = uint32(*nodeData.ExternalSensors.CurrentPower.Number)
		}
		if nodeData.ExternalSensors.Temperature != nil && nodeData.ExternalSensors.Temperature.Number != nil {
			node.ExtSensorsTemp = uint32(*nodeData.ExternalSensors.Temperature.Number)
		}
	}

	return node, nil
}

// convertCommonToAPINodeUpdate converts common NodeUpdate to v0.0.41 API request
func (a *NodeAdapter) convertCommonToAPINodeUpdate(update *types.NodeUpdate) *api.V0041UpdateNodeMsg {
	msg := &api.V0041UpdateNodeMsg{}

	// Set state
	if update.State != nil {
		state := convertNodeStateToAPI(*update.State)
		msg.State = &state
	}

	// Set reason
	if update.Reason != nil {
		msg.Reason = update.Reason
	}

	// Set weight
	if update.Weight != nil {
		weight := int32(*update.Weight)
		msg.Weight = &weight
	}

	// Set features
	if update.Features != nil {
		msg.Features = update.Features
	}

	// Set gres
	if update.Gres != nil {
		msg.Gres = update.Gres
	}

	// Set comment
	if update.Comment != nil {
		msg.Comment = update.Comment
	}

	// Set extra
	if update.Extra != nil {
		msg.Extra = update.Extra
	}

	return msg
}

// convertNodeStateToAPI converts common NodeState to API state
func convertNodeStateToAPI(state types.NodeState) api.V0041UpdateNodeMsgState {
	switch state {
	case types.NodeStateDown:
		return api.V0041UpdateNodeMsgStateDOWN
	case types.NodeStateDrain:
		return api.V0041UpdateNodeMsgStateDRAIN
	case types.NodeStateFail:
		return api.V0041UpdateNodeMsgStateFAIL
	case types.NodeStateFailResp:
		return api.V0041UpdateNodeMsgStateFAILINGRESPONDING
	case types.NodeStateFuture:
		return api.V0041UpdateNodeMsgStateFUTURE
	case types.NodeStateNoResp:
		return api.V0041UpdateNodeMsgStateNOTRESPONDING
	case types.NodeStatePowerDown:
		return api.V0041UpdateNodeMsgStatePOWERDOWN
	case types.NodeStatePowerUp:
		return api.V0041UpdateNodeMsgStatePOWERUP
	case types.NodeStateResume:
		return api.V0041UpdateNodeMsgStateRESUME
	case types.NodeStateUndrain:
		return api.V0041UpdateNodeMsgStateUNDRAIN
	default:
		// Default to resume for unknown states
		return api.V0041UpdateNodeMsgStateRESUME
	}
}