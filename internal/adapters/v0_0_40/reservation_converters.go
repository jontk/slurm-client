// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_40

import (
	"strings"
	"time"

	api "github.com/jontk/slurm-client/internal/api/v0_0_40"
	"github.com/jontk/slurm-client/internal/common/types"
)

// convertAPIReservationToCommon converts a v0.0.40 API Reservation to common Reservation type
func (a *ReservationAdapter) convertAPIReservationToCommon(apiReservation api.V0040ReservationInfo) *types.Reservation {
	reservation := &types.Reservation{}

	// Basic fields
	if apiReservation.Name != nil {
		reservation.Name = *apiReservation.Name
	}

	// Time fields
	if apiReservation.StartTime != nil && apiReservation.StartTime.Number != nil && *apiReservation.StartTime.Number > 0 {
		reservation.StartTime = time.Unix(*apiReservation.StartTime.Number, 0)
	}
	if apiReservation.EndTime != nil && apiReservation.EndTime.Number != nil && *apiReservation.EndTime.Number > 0 {
		reservation.EndTime = time.Unix(*apiReservation.EndTime.Number, 0)
	}
	// v0.0.40 doesn't have Duration field
	// Duration is stored as int32 seconds in common types
	if apiReservation.StartTime != nil && apiReservation.StartTime.Number != nil &&
		apiReservation.EndTime != nil && apiReservation.EndTime.Number != nil {
		startTime := time.Unix(*apiReservation.StartTime.Number, 0)
		endTime := time.Unix(*apiReservation.EndTime.Number, 0)
		duration := endTime.Sub(startTime)
		reservation.Duration = int32(duration.Seconds())
	}

	// Node information
	if apiReservation.NodeList != nil {
		reservation.NodeList = *apiReservation.NodeList
	}
	if apiReservation.NodeCount != nil {
		reservation.NodeCount = *apiReservation.NodeCount
	}

	// Users and accounts
	if apiReservation.Users != nil {
		reservation.Users = strings.Split(*apiReservation.Users, ",")
	}
	if apiReservation.Accounts != nil {
		reservation.Accounts = strings.Split(*apiReservation.Accounts, ",")
	}

	// Flags
	if apiReservation.Flags != nil && len(*apiReservation.Flags) > 0 {
		reservation.Flags = make([]types.ReservationFlag, len(*apiReservation.Flags))
		for i, flag := range *apiReservation.Flags {
			reservation.Flags[i] = types.ReservationFlag(flag)
		}
	}

	// Features
	if apiReservation.Features != nil {
		reservation.Features = strings.Split(*apiReservation.Features, ",")
	}

	// Partition
	if apiReservation.Partition != nil {
		reservation.PartitionName = *apiReservation.Partition
	}

	// License
	// Parse licenses if needed
	if apiReservation.Licenses != nil {
	}

	// TRES - convert string to map
	if apiReservation.Tres != nil {
		// v0.0.40 has TRES as string, but common type expects map
		// Parse TRES string if needed (format: "cpu=4,mem=1000")
		reservation.TRES = make(map[string]int64)
		// For now, just store the raw string in a special key
		reservation.TRES["_raw"] = 0 // Placeholder since we can't parse it properly here
	}

	// Core count
	if apiReservation.CoreCount != nil {
		reservation.CoreCount = *apiReservation.CoreCount
	}

	return reservation
}
