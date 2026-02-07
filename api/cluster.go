// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

// Package types provides common type definitions for SLURM entities.
// Core entity types (Job, Node, User, etc.) are generated in *.gen.go files.
// This file contains operation types (Create, Update, List, etc.).
package api

import "time"

// ClusterList represents a list of clusters
type ClusterList struct {
	Clusters []Cluster              `json:"clusters"`
	Total    int                    `json:"total"`
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

// ClusterUpdate represents a request to update an existing cluster
type ClusterUpdate struct {
	ControlHost        *string  `json:"control_host,omitempty"`
	ControlPort        *int32   `json:"control_port,omitempty"`
	RPCVersion         *int32   `json:"rpc_version,omitempty"`
	PluginIDSelect     *int32   `json:"plugin_id_select,omitempty"`
	PluginIDAuth       *int32   `json:"plugin_id_auth,omitempty"`
	PluginIDAcct       *int32   `json:"plugin_id_acct,omitempty"`
	TRESList           []string `json:"tres_list,omitempty"`
	Features           []string `json:"features,omitempty"`
	FederationFeatures []string `json:"federation_features,omitempty"`
	FederationState    *string  `json:"federation_state,omitempty"`
}

// ClusterListOptions provides filtering options for cluster listing
type ClusterListOptions struct {
	UpdateTime     *time.Time `json:"update_time,omitempty"`
	Classification string     `json:"classification,omitempty"`
	Flags          string     `json:"flags,omitempty"`

	// Limit specifies the maximum number of clusters to return.
	// WARNING: Due to SLURM REST API limitations, this is CLIENT-SIDE pagination.
	// The full cluster list is fetched from the server, then sliced. Consider using
	// filtering options (Classification, Flags) to reduce the dataset before pagination.
	Limit int `json:"limit,omitempty"`

	// Offset specifies the number of clusters to skip before returning results.
	// WARNING: This is CLIENT-SIDE pagination - see Limit field documentation.
	Offset int `json:"offset,omitempty"`
}

// ClusterDeleteOptions provides options for deleting clusters
type ClusterDeleteOptions struct {
	Classification string `json:"classification,omitempty"`
	Flags          string `json:"flags,omitempty"`
}
