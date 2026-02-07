// Code generated from OpenAPI spec. DO NOT EDIT.
// SPDX-FileCopyrightText: 2026 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package api

// JobResSocket represents a SLURM JobResSocket.
type JobResSocket struct {
	Cores []JobResCore `json:"cores"` // Core in socket
	Index int32 `json:"index"` // Core index
}
