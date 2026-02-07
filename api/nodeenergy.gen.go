// Code generated from OpenAPI spec. DO NOT EDIT.
// SPDX-FileCopyrightText: 2026 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package api

// NodeEnergy represents a SLURM NodeEnergy.
type NodeEnergy struct {
	AverageWatts *int32 `json:"average_watts,omitempty"` // Average power consumption, in watts
	BaseConsumedEnergy *int64 `json:"base_consumed_energy,omitempty"` // The energy consumed between when the node was powered on and the last time it...
	ConsumedEnergy *int64 `json:"consumed_energy,omitempty"` // The energy consumed between the last time the node was registered by the slurmd...
	CurrentWatts *uint32 `json:"current_watts,omitempty"` // The instantaneous power consumption at the time of the last node energy...
	LastCollected *int64 `json:"last_collected,omitempty"` // Time when energy data was last retrieved (UNIX timestamp) (UNIX timestamp or...
	PreviousConsumedEnergy *int64 `json:"previous_consumed_energy,omitempty"` // Previous value of consumed_energy
}
