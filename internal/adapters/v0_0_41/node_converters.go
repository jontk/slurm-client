// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

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
	// Use map interface for handling anonymous structs in v0.0.41
	nodeData, ok := apiNode.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected node data type: %T", apiNode)
	}

	node := &types.Node{}

	// Basic fields - using safe type assertions
	if v, ok := nodeData["name"]; ok {
		if name, ok := v.(string); ok {
			node.Name = name
		}
	}
	if v, ok := nodeData["architecture"]; ok {
		if arch, ok := v.(string); ok {
			node.Arch = arch
		}
	}
	if v, ok := nodeData["operating_system"]; ok {
		if os, ok := v.(string); ok {
			node.OS = os
		}
	}
	if v, ok := nodeData["address"]; ok {
		if addr, ok := v.(string); ok {
			node.NodeAddress = addr
		}
	}
	if v, ok := nodeData["hostname"]; ok {
		if hostname, ok := v.(string); ok {
			node.NodeHostname = hostname
		}
	}

	// State
	if v, ok := nodeData["state"]; ok {
		if states, ok := v.([]interface{}); ok && len(states) > 0 {
			if state, ok := states[0].(string); ok {
				node.State = types.NodeState(state)
			}
		}
	}

	// Reason
	if v, ok := nodeData["reason"]; ok {
		if reason, ok := v.(string); ok {
			node.Reason = reason
		}
	}

	// Resources
	if v, ok := nodeData["cpus"]; ok {
		if cpus, ok := v.(float64); ok {
			node.CPUs = int32(cpus)
		}
	}
	if v, ok := nodeData["boards"]; ok {
		if boards, ok := v.(float64); ok {
			node.Boards = int32(boards)
		}
	}
	if v, ok := nodeData["sockets"]; ok {
		if sockets, ok := v.(float64); ok {
			node.Sockets = int32(sockets)
		}
	}
	if v, ok := nodeData["cores"]; ok {
		if cores, ok := v.(float64); ok {
			node.Cores = int32(cores)
		}
	}
	if v, ok := nodeData["threads_per_core"]; ok {
		if threads, ok := v.(float64); ok {
			node.ThreadsPerCore = int32(threads)
		}
	}

	// Memory
	if v, ok := nodeData["real_memory"]; ok {
		if mem, ok := v.(float64); ok {
			node.RealMemory = int64(mem)
		}
	}
	if v, ok := nodeData["alloc_memory"]; ok {
		if mem, ok := v.(float64); ok {
			node.AllocMemory = int64(mem)
		}
	}
	if v, ok := nodeData["free_memory"]; ok {
		if mem, ok := v.(float64); ok {
			node.FreeMemory = int64(mem)
		}
	}

	// CPU allocation
	if v, ok := nodeData["alloc_cpus"]; ok {
		if cpus, ok := v.(float64); ok {
			node.AllocCPUs = int32(cpus)
		}
	}
	if v, ok := nodeData["alloc_idle_cpus"]; ok {
		if cpus, ok := v.(float64); ok {
			node.AllocIdleCPUs = int32(cpus)
		}
	}

	// Time fields - handle both direct numbers and structured time objects
	if v, ok := nodeData["boot_time"]; ok {
		if timeStruct, ok := v.(map[string]interface{}); ok {
			if number, ok := timeStruct["number"].(float64); ok && number > 0 {
				bootTime := time.Unix(int64(number), 0)
				node.BootTime = &bootTime
			}
		}
	}
	if v, ok := nodeData["last_busy"]; ok {
		if timeStruct, ok := v.(map[string]interface{}); ok {
			if number, ok := timeStruct["number"].(float64); ok && number > 0 {
				lastBusy := time.Unix(int64(number), 0)
				node.LastBusy = &lastBusy
			}
		}
	}

	// Features
	if v, ok := nodeData["features"]; ok {
		if features, ok := v.(string); ok {
			node.Features = strings.Split(features, ",")
		}
	}
	if v, ok := nodeData["active_features"]; ok {
		if features, ok := v.(string); ok {
			node.ActiveFeatures = strings.Split(features, ",")
		}
	}

	// Partitions
	if v, ok := nodeData["partitions"]; ok {
		if partitions, ok := v.(string); ok {
			node.Partitions = strings.Split(partitions, ",")
		}
	}

	// GRES
	if v, ok := nodeData["gres"]; ok {
		if gres, ok := v.(string); ok {
			node.Gres = gres
		}
	}
	if v, ok := nodeData["gres_drained"]; ok {
		if gres, ok := v.(string); ok {
			node.GresDrained = gres
		}
	}
	if v, ok := nodeData["gres_used"]; ok {
		if gres, ok := v.(string); ok {
			node.GresUsed = gres
		}
	}

	// Other fields
	if v, ok := nodeData["comment"]; ok {
		if comment, ok := v.(string); ok {
			node.Comment = comment
		}
	}
	if v, ok := nodeData["owner"]; ok {
		if owner, ok := v.(string); ok {
			node.Owner = owner
		}
	}
	if v, ok := nodeData["mcs_label"]; ok {
		if label, ok := v.(string); ok {
			node.MCSLabel = label
		}
	}
	if v, ok := nodeData["weight"]; ok {
		if weight, ok := v.(float64); ok {
			node.Weight = int32(weight)
		}
	}
	if v, ok := nodeData["port"]; ok {
		if port, ok := v.(float64); ok {
			node.Port = int32(port)
		}
	}

	// CPU load
	if v, ok := nodeData["cpu_load"]; ok {
		if cpuStruct, ok := v.(map[string]interface{}); ok {
			if number, ok := cpuStruct["number"].(float64); ok {
				node.CPULoad = number / 100.0 // Convert from centipercent
			}
		}
	}

	// Version
	if v, ok := nodeData["version"]; ok {
		if version, ok := v.(string); ok {
			node.Version = version
		}
	}

	return node, nil
}

// convertCommonToAPINodeUpdate converts common node update request to v0.0.41 API format
func (a *NodeAdapter) convertCommonToAPINodeUpdate(update *types.NodeUpdate) *api.SlurmV0041PostNodeJSONRequestBody {
	// Create a basic update request structure
	updateReq := &api.SlurmV0041PostNodeJSONRequestBody{}

	// Note: The exact structure for node updates in v0.0.41 may be different
	// This is a placeholder implementation that would need to be adjusted
	// based on the actual API structure

	// For now, return an empty request
	_ = update
	return updateReq
}
