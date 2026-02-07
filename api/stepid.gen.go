// Code generated from OpenAPI spec. DO NOT EDIT.
// SPDX-FileCopyrightText: 2026 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package api

// StepID represents a SLURM StepID.
type StepID struct {
	JobID *uint32 `json:"job_id,omitempty"` // Job ID (32 bit integer number with flags)
	Sluid *string `json:"sluid,omitempty"` // SLUID (Slurm Lexicographically-sortable Unique ID)
	StepHetComponent *uint32 `json:"step_het_component,omitempty"` // HetJob component (32 bit integer number with flags)
	StepID *string `json:"step_id,omitempty"` // Job step ID
}
