// Code generated from OpenAPI spec. DO NOT EDIT.
// SPDX-FileCopyrightText: 2026 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package api

// NodeUpdate represents a SLURM NodeUpdate.
type NodeUpdate struct {
	Address []string `json:"address,omitempty"` // NodeAddr, used to establish a communication path
	Comment *string `json:"comment,omitempty"` // Arbitrary comment
	CPUBind *int32 `json:"cpu_bind,omitempty"` // Default method for binding tasks to allocated CPUs
	Extra *string `json:"extra,omitempty"` // Arbitrary string used for node filtering if extra constraints are enabled
	Features []string `json:"features,omitempty"` // Available features
	FeaturesAct []string `json:"features_act,omitempty"` // Currently active features
	GRES *string `json:"gres,omitempty"` // Generic resources
	Hostname []string `json:"hostname,omitempty"` // NodeHostname
	Name []string `json:"name,omitempty"` // NodeName
	Reason *string `json:"reason,omitempty"` // Reason for node being DOWN or DRAINING
	ReasonUID *string `json:"reason_uid,omitempty"` // User ID to associate with the reason (needed if user root is sending message)
	ResumeAfter *uint32 `json:"resume_after,omitempty"` // Number of seconds after which to automatically resume DOWN or DRAINED node (32...
	State []NodeState `json:"state,omitempty"` // New state to assign to the node
	TopologyStr *string `json:"topology_str,omitempty"` // Topology
	Weight *uint32 `json:"weight,omitempty"` // Weight of the node for scheduling purposes (32 bit integer number with flags)
}
