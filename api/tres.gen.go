// Code generated from OpenAPI spec. DO NOT EDIT.
// SPDX-FileCopyrightText: 2026 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package api

// TRES represents a SLURM TRES.
type TRES struct {
	Count *int64 `json:"count,omitempty"` // TRES count (0 if listed generically)
	ID *int32 `json:"id,omitempty"` // ID used in the database
	Name *string `json:"name,omitempty"` // TRES name (if applicable)
	Type string `json:"type"` // TRES type (CPU, MEM, etc)
}
