// Code generated from OpenAPI spec. DO NOT EDIT.
// SPDX-FileCopyrightText: 2026 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package api

// Accounting represents a SLURM Accounting.
type Accounting struct {
	TRES *TRES `json:"TRES,omitempty"` // Trackable resources
	Allocated *AccountingAllocated `json:"allocated,omitempty"`
	ID *int32 `json:"id,omitempty"` // Association ID or Workload characterization key ID
	IDAlt *int32 `json:"id_alt,omitempty"` // Alternate ID (not currently used)
	Start *int64 `json:"start,omitempty"` // When the record was started (UNIX timestamp) (UNIX timestamp or time string...
}


// AccountingAllocated is a nested type within its parent.
type AccountingAllocated struct {
	Seconds *int64 `json:"seconds,omitempty"` // Number of seconds allocated
}
