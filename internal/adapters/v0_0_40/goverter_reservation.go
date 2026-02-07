// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_40

import (
	types "github.com/jontk/slurm-client/api"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_40"
)

// ReservationConverterGoverter defines the goverter interface for Reservation conversions.
// This converts API V0040ReservationInfo to the common types.Reservation type.
// goverter:converter
// goverter:output:file reservation_converters_goverter.gen.go
// goverter:output:package github.com/jontk/slurm-client/internal/adapters/v0_0_40
// goverter:extend ConvertTimeNoVal
// goverter:extend ConvertUint32NoVal
// goverter:extend ConvertReservationFlagsRead
// goverter:extend ConvertReservationCoreSpec
// goverter:extend ConvertReservationPurgeCompleted
type ReservationConverterGoverter interface {
	// ConvertAPIReservationToCommon converts API V0040ReservationInfo to common Reservation type
	//
	// Time fields (use ConvertTimeNoVal):
	// goverter:map EndTime | ConvertTimeNoVal
	// goverter:map StartTime | ConvertTimeNoVal
	//
	// NoValNumber fields (use ConvertUint32NoVal):
	// goverter:map Watts | ConvertUint32NoVal
	//
	// Field name mappings (Source -> Target):
	// goverter:map Tres TRES
	//
	// Complex type conversions:
	// goverter:map Flags | ConvertReservationFlagsRead
	// goverter:map CoreSpecializations | ConvertReservationCoreSpec
	// goverter:map PurgeCompleted | ConvertReservationPurgeCompleted
	ConvertAPIReservationToCommon(source api.V0040ReservationInfo) *types.Reservation
}
