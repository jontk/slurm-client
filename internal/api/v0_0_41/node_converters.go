// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_41

import (
	"strings"
	"time"

	"github.com/jontk/slurm-client/interfaces"
)

// convertNodeFromAPI converts v0.0.41 API node struct to interfaces.Node
func convertNodeFromAPI(apiNode interface{}) interfaces.Node {
	// Type assertion to handle the actual API response structure
	nodeData, ok := apiNode.(struct {
		ActiveFeatures *[]string `json:"active_features,omitempty"`
		Address        *string   `json:"address,omitempty"`
		AllocCpus      *int32    `json:"alloc_cpus,omitempty"`
		AllocIdleCpus  *int32    `json:"alloc_idle_cpus,omitempty"`
		AllocMemory    *int64    `json:"alloc_memory,omitempty"`
		Architecture   *string   `json:"architecture,omitempty"`
		Boards         *int32    `json:"boards,omitempty"`
		BootTime       *struct {
			Infinite *bool  `json:"infinite,omitempty"`
			Number   *int64 `json:"number,omitempty"`
			Set      *bool  `json:"set,omitempty"`
		} `json:"boot_time,omitempty"`
		BurstbufferNetworkAddress *string   `json:"burstbuffer_network_address,omitempty"`
		ClusterName               *string   `json:"cluster_name,omitempty"`
		Comment                   *string   `json:"comment,omitempty"`
		Cores                     *int32    `json:"cores,omitempty"`
		CpuBinding                *int32    `json:"cpu_binding,omitempty"`
		CpuLoad                   *int32    `json:"cpu_load,omitempty"`
		Cpus                      *int32    `json:"cpus,omitempty"`
		EffectiveCpus             *int32    `json:"effective_cpus,omitempty"`
		Features                  *[]string `json:"features,omitempty"`
		FreeMem                   *struct {
			Infinite *bool  `json:"infinite,omitempty"`
			Number   *int64 `json:"number,omitempty"`
			Set      *bool  `json:"set,omitempty"`
		} `json:"free_mem,omitempty"`
		Hostname             *string                                           `json:"hostname,omitempty"`
		Name                 *string                                           `json:"name,omitempty"`
		NextStateAfterReboot *[]V0041OpenapiNodesRespNodesNextStateAfterReboot `json:"next_state_after_reboot,omitempty"`
		OperatingSystem      *string                                           `json:"operating_system,omitempty"`
		Owner                *string                                           `json:"owner,omitempty"`
		Partitions           *[]string                                         `json:"partitions,omitempty"`
		Port                 *int32                                            `json:"port,omitempty"`
		RealMemory           *int64                                            `json:"real_memory,omitempty"`
		Reason               *string                                           `json:"reason,omitempty"`
		ReasonChangedAt      *struct {
			Infinite *bool  `json:"infinite,omitempty"`
			Number   *int64 `json:"number,omitempty"`
			Set      *bool  `json:"set,omitempty"`
		} `json:"reason_changed_at,omitempty"`
		ReasonSetByUser *string                            `json:"reason_set_by_user,omitempty"`
		Sockets         *int32                             `json:"sockets,omitempty"`
		State           *[]V0041OpenapiNodesRespNodesState `json:"state,omitempty"`
		TemporaryDisk   *int32                             `json:"temporary_disk,omitempty"`
		Threads         *int32                             `json:"threads,omitempty"`
		Weight          *int32                             `json:"weight,omitempty"`
	})
	if !ok {
		// If type assertion fails, return empty node
		return interfaces.Node{}
	}

	node := interfaces.Node{
		// Initialize metadata
		Metadata: make(map[string]interface{}),
	}

	// Basic info
	if nodeData.Name != nil {
		node.Name = *nodeData.Name
	}

	// State
	if nodeData.State != nil && len(*nodeData.State) > 0 {
		node.State = string((*nodeData.State)[0])
	} else {
		node.State = "UNKNOWN"
	}

	// Resources
	if nodeData.Cpus != nil {
		node.CPUs = int(*nodeData.Cpus)
	}
	if nodeData.RealMemory != nil {
		node.Memory = int(*nodeData.RealMemory) * 1024 * 1024 // Convert MB to bytes
	}

	// Features
	if nodeData.Features != nil {
		node.Features = *nodeData.Features
	}

	// Partitions
	if nodeData.Partitions != nil {
		node.Partitions = *nodeData.Partitions
	}

	// Reason
	if nodeData.Reason != nil {
		node.Reason = *nodeData.Reason
	}

	// Architecture info - store in metadata
	if nodeData.Architecture != nil {
		node.Architecture = *nodeData.Architecture
	}

	// CPU Load - convert int32 to float64
	if nodeData.CpuLoad != nil {
		node.CPULoad = float64(*nodeData.CpuLoad)
	}

	// Allocated CPUs
	if nodeData.AllocCpus != nil {
		node.AllocCPUs = *nodeData.AllocCpus
	}

	// Allocated Memory (already in MB from API)
	if nodeData.AllocMemory != nil {
		node.AllocMemory = *nodeData.AllocMemory
	}

	// Free Memory (in MB from API)
	if nodeData.FreeMem != nil && nodeData.FreeMem.Set != nil && *nodeData.FreeMem.Set &&
		nodeData.FreeMem.Number != nil {
		node.FreeMemory = *nodeData.FreeMem.Number
	}

	// Store additional fields in metadata
	if nodeData.Hostname != nil {
		node.Metadata["hostname"] = *nodeData.Hostname
	}
	if nodeData.Address != nil {
		node.Metadata["address"] = *nodeData.Address
	}
	if nodeData.TemporaryDisk != nil {
		node.Metadata["tmp_disk"] = *nodeData.TemporaryDisk
	}
	if nodeData.OperatingSystem != nil {
		node.Metadata["os"] = *nodeData.OperatingSystem
	}
	if nodeData.ActiveFeatures != nil {
		node.Metadata["active_features"] = strings.Join(*nodeData.ActiveFeatures, ",")
	}
	if nodeData.ReasonChangedAt != nil && nodeData.ReasonChangedAt.Number != nil {
		node.Metadata["reason_set_time"] = time.Unix(*nodeData.ReasonChangedAt.Number, 0)
	}
	if nodeData.ReasonSetByUser != nil {
		node.Metadata["reason_user"] = *nodeData.ReasonSetByUser
	}
	if nodeData.Sockets != nil {
		sockets := int(*nodeData.Sockets)
		if nodeData.Boards != nil {
			sockets *= int(*nodeData.Boards)
		}
		node.Metadata["sockets"] = sockets
	}
	if nodeData.Cores != nil {
		node.Metadata["cores_per_socket"] = *nodeData.Cores
	}
	if nodeData.Threads != nil {
		node.Metadata["threads_per_core"] = *nodeData.Threads
	}
	if nodeData.BootTime != nil && nodeData.BootTime.Number != nil {
		bootTime := time.Unix(*nodeData.BootTime.Number, 0)
		node.LastBusy = &bootTime // Use LastBusy as proxy for boot time
	}
	if nodeData.Weight != nil {
		node.Metadata["weight"] = *nodeData.Weight
	}
	if nodeData.Owner != nil {
		node.Metadata["owner"] = *nodeData.Owner
	}
	if nodeData.Port != nil {
		node.Metadata["port"] = *nodeData.Port
	}

	return node
}
