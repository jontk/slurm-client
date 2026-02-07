// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_41

import (
	"fmt"
	"strings"
	"time"

	types "github.com/jontk/slurm-client/api"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_41"
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
			n := name
			node.Name = &n
		}
	}
	if v, ok := nodeData["architecture"]; ok {
		if arch, ok := v.(string); ok {
			a := arch
			node.Architecture = &a
		}
	}
	if v, ok := nodeData["operating_system"]; ok {
		if osVal, ok := v.(string); ok {
			o := osVal
			node.OperatingSystem = &o
		}
	}
	if v, ok := nodeData["address"]; ok {
		if addr, ok := v.(string); ok {
			a := addr
			node.Address = &a
		}
	}
	if v, ok := nodeData["hostname"]; ok {
		if hostname, ok := v.(string); ok {
			h := hostname
			node.Hostname = &h
		}
	}
	// State - SLURM API returns state as an array (e.g. ["IDLE", "DRAIN"])
	if v, ok := nodeData["state"]; ok {
		if states, ok := v.([]interface{}); ok {
			nodeStates := make([]types.NodeState, 0, len(states))
			for _, s := range states {
				if state, ok := s.(string); ok {
					nodeStates = append(nodeStates, types.NodeState(state))
				}
			}
			node.State = nodeStates
		}
	}
	// Reason
	if v, ok := nodeData["reason"]; ok {
		if reason, ok := v.(string); ok {
			r := reason
			node.Reason = &r
		}
	}
	// Resources
	if v, ok := nodeData["cpus"]; ok {
		if cpus, ok := v.(float64); ok {
			c := int32(cpus)
			node.CPUs = &c
		}
	}
	if v, ok := nodeData["boards"]; ok {
		if boards, ok := v.(float64); ok {
			b := int32(boards)
			node.Boards = &b
		}
	}
	if v, ok := nodeData["sockets"]; ok {
		if sockets, ok := v.(float64); ok {
			s := int32(sockets)
			node.Sockets = &s
		}
	}
	if v, ok := nodeData["cores"]; ok {
		if cores, ok := v.(float64); ok {
			c := int32(cores)
			node.Cores = &c
		}
	}
	if v, ok := nodeData["threads_per_core"]; ok {
		if threads, ok := v.(float64); ok {
			t := int32(threads)
			node.Threads = &t
		}
	}
	// Memory
	if v, ok := nodeData["real_memory"]; ok {
		if mem, ok := v.(float64); ok {
			m := int64(mem)
			node.RealMemory = &m
		}
	}
	if v, ok := nodeData["alloc_memory"]; ok {
		if mem, ok := v.(float64); ok {
			m := int64(mem)
			node.AllocMemory = &m
		}
	}
	if v, ok := nodeData["free_memory"]; ok {
		if mem, ok := v.(float64); ok {
			m := uint64(mem)
			node.FreeMem = &m
		}
	}
	// CPU allocation
	if v, ok := nodeData["alloc_cpus"]; ok {
		if cpus, ok := v.(float64); ok {
			c := int32(cpus)
			node.AllocCPUs = &c
		}
	}
	if v, ok := nodeData["alloc_idle_cpus"]; ok {
		if cpus, ok := v.(float64); ok {
			c := int32(cpus)
			node.AllocIdleCPUs = &c
		}
	}
	// Time fields - handle both direct numbers and structured time objects
	if v, ok := nodeData["boot_time"]; ok {
		if timeStruct, ok := v.(map[string]interface{}); ok {
			if number, ok := timeStruct["number"].(float64); ok && number > 0 {
				node.BootTime = time.Unix(int64(number), 0)
			}
		}
	}
	if v, ok := nodeData["last_busy"]; ok {
		if timeStruct, ok := v.(map[string]interface{}); ok {
			if number, ok := timeStruct["number"].(float64); ok && number > 0 {
				node.LastBusy = time.Unix(int64(number), 0)
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
			g := gres
			node.GRES = &g
		}
	}
	if v, ok := nodeData["gres_drained"]; ok {
		if gres, ok := v.(string); ok {
			g := gres
			node.GRESDrained = &g
		}
	}
	if v, ok := nodeData["gres_used"]; ok {
		if gres, ok := v.(string); ok {
			g := gres
			node.GRESUsed = &g
		}
	}
	// Other fields
	if v, ok := nodeData["comment"]; ok {
		if comment, ok := v.(string); ok {
			c := comment
			node.Comment = &c
		}
	}
	if v, ok := nodeData["owner"]; ok {
		if owner, ok := v.(string); ok {
			o := owner
			node.Owner = &o
		}
	}
	if v, ok := nodeData["mcs_label"]; ok {
		if label, ok := v.(string); ok {
			l := label
			node.MCSLabel = &l
		}
	}
	if v, ok := nodeData["weight"]; ok {
		if weight, ok := v.(float64); ok {
			w := int32(weight)
			node.Weight = &w
		}
	}
	if v, ok := nodeData["port"]; ok {
		if port, ok := v.(float64); ok {
			p := int32(port)
			node.Port = &p
		}
	}
	// CPU load
	if v, ok := nodeData["cpu_load"]; ok {
		if cpuStruct, ok := v.(map[string]interface{}); ok {
			if number, ok := cpuStruct["number"].(float64); ok {
				cl := int32(number / 100.0) // Convert from centipercent
				node.CPULoad = &cl
			}
		}
	}
	// Version
	if v, ok := nodeData["version"]; ok {
		if version, ok := v.(string); ok {
			ver := version
			node.Version = &ver
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
