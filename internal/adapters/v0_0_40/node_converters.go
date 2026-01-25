// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_40

import (
	"strings"

	api "github.com/jontk/slurm-client/internal/api/v0_0_40"
	"github.com/jontk/slurm-client/internal/common/types"
)

// convertAPINodeToCommon converts a v0.0.40 API Node to common Node type
func (a *NodeAdapter) convertAPINodeToCommon(apiNode api.V0040Node) *types.Node {
	node := &types.Node{}

	// Essential fields only for v0.0.40 - minimal conversion to get building
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

	// State - handle safely
	if apiNode.State != nil && len(*apiNode.State) > 0 {
		stateStr := (*apiNode.State)[0]
		node.State = types.NodeState(strings.ToUpper(stateStr))
	}

	// Resources - handle safely
	if apiNode.Cpus != nil {
		node.CPUs = *apiNode.Cpus
	}
	if apiNode.RealMemory != nil {
		node.RealMemory = *apiNode.RealMemory
	}
	if apiNode.AllocMemory != nil {
		node.AllocMemory = *apiNode.AllocMemory
	}

	// GRES - handle safely
	if apiNode.Gres != nil {
		node.Gres = *apiNode.Gres
	}

	// Features - V0040CsvString is []string already
	if apiNode.Features != nil {
		node.Features = *apiNode.Features
	}
	if apiNode.ActiveFeatures != nil {
		node.ActiveFeatures = *apiNode.ActiveFeatures
	}

	// Basic string fields
	if apiNode.Reason != nil {
		node.Reason = *apiNode.Reason
	}
	if apiNode.Comment != nil {
		node.Comment = *apiNode.Comment
	}

	return node
}

// convertCommonNodeUpdateToAPI converts common NodeUpdate to v0.0.40 API format
func (a *NodeAdapter) convertCommonNodeUpdateToAPI(nodeName string, update *types.NodeUpdate) *api.V0040UpdateNodeMsg {
	apiNode := &api.V0040UpdateNodeMsg{}

	// Name (required) - V0040HostlistString is []string
	nameList := api.V0040HostlistString{nodeName}
	apiNode.Name = &nameList

	// Basic fields - only what's essential for v0.0.40
	if update.State != nil {
		state := api.V0040NodeStates{string(*update.State)}
		apiNode.State = &state
	}

	if update.Reason != nil {
		apiNode.Reason = update.Reason
	}

	if update.Comment != nil {
		apiNode.Comment = update.Comment
	}

	if update.Gres != nil {
		apiNode.Gres = update.Gres
	}

	// Features - V0040CsvString is []string
	if len(update.Features) > 0 {
		features := update.Features
		apiNode.Features = &features
	}

	if len(update.ActiveFeatures) > 0 {
		activeFeatures := update.ActiveFeatures
		apiNode.FeaturesAct = &activeFeatures
	}

	return apiNode
}
