// Code generated from OpenAPI spec. DO NOT EDIT.
// SPDX-FileCopyrightText: 2026 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package api

// AssocShort represents a SLURM AssocShort.
type AssocShort struct {
	Account *string `json:"account,omitempty"` // Account name
	Cluster *string `json:"cluster,omitempty"` // Cluster name
	ID *int32 `json:"id,omitempty"` // Numeric association ID
	Partition *string `json:"partition,omitempty"` // Partition name
	User string `json:"user"` // User name
}
