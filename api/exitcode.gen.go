// Code generated from OpenAPI spec. DO NOT EDIT.
// SPDX-FileCopyrightText: 2026 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package api

// ExitCode represents a SLURM ExitCode.
type ExitCode struct {
	ReturnCode *uint32 `json:"return_code,omitempty"` // Process return code (numeric) (32 bit integer number with flags)
	Signal *ExitCodeSignal `json:"signal,omitempty"`
	Status []StatusValue `json:"status,omitempty"` // Status given by return code
}


// ExitCodeSignal is a nested type within its parent.
type ExitCodeSignal struct {
	ID *uint16 `json:"id,omitempty"` // Signal sent to process (numeric) (16 bit integer number with flags)
	Name *string `json:"name,omitempty"` // Signal sent to process (name)
}


// StatusValue represents possible values for Status field.
type StatusValue string

// StatusValue constants.
const (
	StatusInvalid StatusValue = "INVALID"
	StatusPending StatusValue = "PENDING"
	StatusSuccess StatusValue = "SUCCESS"
	StatusError StatusValue = "ERROR"
	StatusSignaled StatusValue = "SIGNALED"
	StatusCoreDumped StatusValue = "CORE_DUMPED"
)
