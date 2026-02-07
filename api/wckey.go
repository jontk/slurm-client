// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

// Package types provides common type definitions for SLURM entities.
// Core entity types (WCKey, etc.) are generated in *.gen.go files.
// This file contains operation types (Create, Update, List, etc.).
package api

// WCKeyList represents a list of WCKeys
type WCKeyList struct {
	WCKeys []WCKey                `json:"wckeys"`
	Total  int                    `json:"total,omitempty"`
	Meta   map[string]interface{} `json:"meta,omitempty"`
}

// WCKeyCreate represents a request to create a new WCKey
type WCKeyCreate struct {
	Name        string `json:"name"`
	User        string `json:"user,omitempty"`
	Cluster     string `json:"cluster"`
	Description string `json:"description,omitempty"`
}

// WCKeyCreateResponse represents the response from creating a WCKey
type WCKeyCreateResponse struct {
	ID      string                 `json:"id"`
	Status  string                 `json:"status"`
	Message string                 `json:"message,omitempty"`
	Meta    map[string]interface{} `json:"meta,omitempty"`
}

// WCKeyListOptions provides filtering options for WCKey listing
type WCKeyListOptions struct {
	Users        []string `json:"users,omitempty"`
	Clusters     []string `json:"clusters,omitempty"`
	Names        []string `json:"names,omitempty"`
	OnlyDefaults bool     `json:"only_defaults,omitempty"`
	WithDeleted  bool     `json:"with_deleted,omitempty"`

	// Limit specifies the maximum number of WCKeys to return.
	// WARNING: Due to SLURM REST API limitations, this is CLIENT-SIDE pagination.
	// The full WCKey list is fetched from the server, then sliced. Consider using
	// filtering options (Users, Clusters, Names) to reduce the dataset before pagination.
	Limit int `json:"limit,omitempty"`

	// Offset specifies the number of WCKeys to skip before returning results.
	// WARNING: This is CLIENT-SIDE pagination - see Limit field documentation.
	Offset int `json:"offset,omitempty"`
}

// WCKeyUpdate represents a request to update an existing WCKey
type WCKeyUpdate struct {
	Description *string `json:"description,omitempty"`
}
