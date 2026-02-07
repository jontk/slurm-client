// Code generated from OpenAPI spec. DO NOT EDIT.
// SPDX-FileCopyrightText: 2026 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package api

// WCKey represents a SLURM WCKey.
type WCKey struct {
	Accounting []Accounting `json:"accounting,omitempty"` // Accounting records containing related resource usage
	Cluster string `json:"cluster"` // Cluster name
	Flags []WCKeyFlagsValue `json:"flags,omitempty"` // Flags associated with this WCKey
	ID *int32 `json:"id,omitempty"` // Unique ID for this user-cluster-wckey combination
	Name string `json:"name"` // WCKey name
	User string `json:"user"` // User name
}


// WCKeyFlagsValue represents possible values for WCKeyFlags field.
type WCKeyFlagsValue string

// WCKeyFlagsValue constants.
const (
	WCKeyFlagsDeleted WCKeyFlagsValue = "DELETED"
)
