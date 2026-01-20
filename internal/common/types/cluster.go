// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package types

import "time"

// Cluster represents a SLURM cluster configuration
type Cluster struct {
	Name           string                 `json:"name"`
	ControllerHost string                 `json:"controller_host,omitempty"`
	ControllerPort int32                  `json:"controller_port,omitempty"`
	Nodes          string                 `json:"nodes,omitempty"`
	RpcVersion     int32                  `json:"rpc_version,omitempty"`
	SelectPlugin   string                 `json:"select_plugin,omitempty"`
	Flags          []string               `json:"flags,omitempty"`
	TRES           []TRES                 `json:"tres,omitempty"`
	Associations   *AssociationShort      `json:"associations,omitempty"`
	CreatedTime    *time.Time             `json:"created_time,omitempty"`
	ModifiedTime   *time.Time             `json:"modified_time,omitempty"`
	Meta           map[string]interface{} `json:"meta,omitempty"`
}

// ClusterList represents a list of clusters
type ClusterList struct {
	Clusters []Cluster              `json:"clusters"`
	Meta     map[string]interface{} `json:"meta,omitempty"`
}

// ClusterCreate represents a request to create a new cluster
type ClusterCreate struct {
	Name           string   `json:"name"`
	ControllerHost string   `json:"controller_host,omitempty"`
	ControllerPort int32    `json:"controller_port,omitempty"`
	Nodes          string   `json:"nodes,omitempty"`
	RpcVersion     int32    `json:"rpc_version,omitempty"`
	SelectPlugin   string   `json:"select_plugin,omitempty"`
	Flags          []string `json:"flags,omitempty"`
}

// ClusterCreateResponse represents the response from creating a cluster
type ClusterCreateResponse struct {
	Name    string                 `json:"name"`
	Status  string                 `json:"status"`
	Message string                 `json:"message,omitempty"`
	Meta    map[string]interface{} `json:"meta,omitempty"`
}

// ClusterListOptions provides filtering options for cluster listing
type ClusterListOptions struct {
	UpdateTime     *time.Time `json:"update_time,omitempty"`
	Classification string     `json:"classification,omitempty"`
	Flags          string     `json:"flags,omitempty"`
}

// ClusterDeleteOptions provides options for deleting clusters
type ClusterDeleteOptions struct {
	Classification string `json:"classification,omitempty"`
	Flags          string `json:"flags,omitempty"`
}

// AssociationShort represents a simplified association structure
type AssociationShort struct {
	Root *AssocShort `json:"root,omitempty"`
}

// AssocShort represents basic association information
type AssocShort struct {
	Account   string `json:"account,omitempty"`
	Cluster   string `json:"cluster,omitempty"`
	Partition string `json:"partition,omitempty"`
	User      string `json:"user,omitempty"`
}
