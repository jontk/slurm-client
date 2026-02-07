// Code generated from OpenAPI spec. DO NOT EDIT.
// SPDX-FileCopyrightText: 2026 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package api

// JobResCore represents a SLURM JobResCore.
type JobResCore struct {
	Index int32 `json:"index"` // Core index
	Status []JobResCoreStatusValue `json:"status"` // Core status
}


// JobResCoreStatusValue represents possible values for JobResCoreStatus field.
type JobResCoreStatusValue string

// JobResCoreStatusValue constants.
const (
	JobResCoreStatusInvalid JobResCoreStatusValue = "INVALID"
	JobResCoreStatusUnallocated JobResCoreStatusValue = "UNALLOCATED"
	JobResCoreStatusAllocated JobResCoreStatusValue = "ALLOCATED"
	JobResCoreStatusInUse JobResCoreStatusValue = "IN_USE"
)
