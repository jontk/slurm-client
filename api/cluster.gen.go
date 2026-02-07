// Code generated from OpenAPI spec. DO NOT EDIT.
// SPDX-FileCopyrightText: 2026 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package api

// Cluster represents a SLURM Cluster.
type Cluster struct {
	Associations *ClusterAssociations `json:"associations,omitempty"`
	Controller *ClusterController `json:"controller,omitempty"`
	Flags []ClusterControllerFlagsValue `json:"flags,omitempty"` // Flags
	Name *string `json:"name,omitempty"` // ClusterName
	Nodes *string `json:"nodes,omitempty"` // Node names
	RpcVersion *int32 `json:"rpc_version,omitempty"` // RPC version used in the cluster
	SelectPlugin *string `json:"select_plugin,omitempty"`
	TRES []TRES `json:"tres,omitempty"` // Trackable resources
}


// ClusterAssociations is a nested type within its parent.
type ClusterAssociations struct {
	Root *AssocShort `json:"root,omitempty"` // Root association information
}


// ClusterController is a nested type within its parent.
type ClusterController struct {
	Host *string `json:"host,omitempty"` // ControlHost
	Port *int32 `json:"port,omitempty"` // ControlPort
}


// ClusterControllerFlagsValue represents possible values for ClusterControllerFlags field.
type ClusterControllerFlagsValue string

// ClusterControllerFlagsValue constants.
const (
	ClusterControllerFlagsDeleted ClusterControllerFlagsValue = "DELETED"
	ClusterControllerFlagsRegistering ClusterControllerFlagsValue = "REGISTERING"
	ClusterControllerFlagsMultipleSlurmd ClusterControllerFlagsValue = "MULTIPLE_SLURMD"
	ClusterControllerFlagsFederation ClusterControllerFlagsValue = "FEDERATION"
	ClusterControllerFlagsExternal ClusterControllerFlagsValue = "EXTERNAL"
)
