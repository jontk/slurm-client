// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_40

import (
	"fmt"
	"strings"
	"time"

	"github.com/jontk/slurm-client/internal/common/types"
	api "github.com/jontk/slurm-client/internal/api/v0_0_40"
)

// convertAPIReservationToCommon converts a v0.0.40 API Reservation to common Reservation type
func (a *ReservationAdapter) convertAPIReservationToCommon(apiReservation api.V0040ReservationInfo) (*types.Reservation, error) {
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
	if apiReservation.Licenses != nil {
		// Parse licenses if needed
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

	return reservation, nil
}

// convertCommonReservationCreateToAPI converts common ReservationCreate to v0.0.40 API format
func (a *ReservationAdapter) convertCommonReservationCreateToAPI(reservation *types.ReservationCreate) (*api.V0040ReservationInfo, error) {
	apiReservation := &api.V0040ReservationInfo{}

	// Basic fields
	apiReservation.Name = &reservation.Name

	// Time fields
	startTime := api.V0040Uint64NoVal{
		Set:    boolPtr(true),
		Number: int64Ptr(reservation.StartTime.Unix()),
	}
	apiReservation.StartTime = &startTime

	endTime := api.V0040Uint64NoVal{
		Set:    boolPtr(true),
		Number: int64Ptr(reservation.EndTime.Unix()),
	}
	apiReservation.EndTime = &endTime

	// v0.0.40 doesn't have Duration field in API
	// Skip duration handling

	// Node information
	if reservation.NodeList != "" {
		apiReservation.NodeList = &reservation.NodeList
	}
	if reservation.NodeCount > 0 {
		apiReservation.NodeCount = &reservation.NodeCount
	}

	// Users and accounts
	if len(reservation.Users) > 0 {
		users := strings.Join(reservation.Users, ",")
		apiReservation.Users = &users
	}
	if len(reservation.Accounts) > 0 {
		accounts := strings.Join(reservation.Accounts, ",")
		apiReservation.Accounts = &accounts
	}

	// Flags
	if len(reservation.Flags) > 0 {
		flags := make([]string, len(reservation.Flags))
		for i, flag := range reservation.Flags {
			flags[i] = string(flag)
		}
		apiReservation.Flags = &flags
	}

	// Features
	if len(reservation.Features) > 0 {
		features := strings.Join(reservation.Features, ",")
		apiReservation.Features = &features
	}

	// Partition
	if reservation.PartitionName != "" {
		apiReservation.Partition = &reservation.PartitionName
	}

	// TRES - convert map to string format
	if len(reservation.TRES) > 0 {
		tresStrs := make([]string, 0, len(reservation.TRES))
		for tres, count := range reservation.TRES {
			tresStrs = append(tresStrs, fmt.Sprintf("%s=%d", tres, count))
		}
		tres := strings.Join(tresStrs, ",")
		apiReservation.Tres = &tres
	}

	// Core count
	if reservation.CoreCount > 0 {
		apiReservation.CoreCount = &reservation.CoreCount
	}

	return apiReservation, nil
}
