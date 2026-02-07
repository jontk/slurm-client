// Code generated from OpenAPI spec. DO NOT EDIT.
// SPDX-FileCopyrightText: 2026 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package api

// JobPartitionPriority represents a SLURM JobPartitionPriority.
type JobPartitionPriority struct {
	Partition *string `json:"partition,omitempty"` // Partition name
	Priority *int32 `json:"priority,omitempty"` // Prospective job priority if it runs in this partition
}
