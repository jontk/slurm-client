// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package types

import "time"

// WCKey represents a Workload Characterization Key
type WCKey struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	User        string                 `json:"user,omitempty"`
	Cluster     string                 `json:"cluster"`
	Description string                 `json:"description,omitempty"`
	Active      bool                   `json:"active"`
	CreatedTime *time.Time             `json:"created_time,omitempty"`
	ModifiedTime *time.Time            `json:"modified_time,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
}

// WCKeyList represents a list of WCKeys
type WCKeyList struct {
	WCKeys []WCKey                `json:"wckeys"`
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
}