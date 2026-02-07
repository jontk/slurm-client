// Code generated from OpenAPI spec. DO NOT EDIT.
// SPDX-FileCopyrightText: 2026 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package api

// ReservationCoreSpec represents a SLURM ReservationCoreSpec.
type ReservationCoreSpec struct {
	Core *string `json:"core,omitempty"` // IDs of reserved cores
	Node *string `json:"node,omitempty"` // Name of reserved node
}
