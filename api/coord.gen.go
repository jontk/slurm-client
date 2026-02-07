// Code generated from OpenAPI spec. DO NOT EDIT.
// SPDX-FileCopyrightText: 2026 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package api

// Coord represents a SLURM Coord.
type Coord struct {
	Direct *bool `json:"direct,omitempty"` // Indicates whether the coordinator was directly assigned to this account
	Name string `json:"name"` // User name
}
