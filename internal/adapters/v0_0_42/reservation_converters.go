// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_42

import (
	"fmt"
	"strings"
	"time"

	"github.com/jontk/slurm-client/internal/common/types"
	api "github.com/jontk/slurm-client/internal/api/v0_0_42"
)

// convertAPIReservationToCommon converts a v0.0.42 API Reservation to common Reservation type
func (a *ReservationAdapter) convertAPIReservationToCommon(apiReservation api.V0042ReservationInfo) (*types.Reservation, error) {
	reservation := &types.Reservation{}

	// Basic fields
	if apiReservation.Name != nil {
		reservation.Name = *apiReservation.Name
	}
	if apiReservation.Partition != nil {
		reservation.PartitionName = *apiReservation.Partition
	}

	// Time fields
	if apiReservation.StartTime != nil && apiReservation.StartTime.Number != nil && *apiReservation.StartTime.Number > 0 {
		reservation.StartTime = time.Unix(*apiReservation.StartTime.Number, 0)
	}
	if apiReservation.EndTime != nil && apiReservation.EndTime.Number != nil && *apiReservation.EndTime.Number > 0 {
		reservation.EndTime = time.Unix(*apiReservation.EndTime.Number, 0)
	}

	// Resource counts
	if apiReservation.NodeCount != nil {
		reservation.NodeCount = int32(*apiReservation.NodeCount)
	}
	if apiReservation.CoreCount != nil {
		reservation.CoreCount = int32(*apiReservation.CoreCount)
	}

	// Node list
	if apiReservation.NodeList != nil {
		reservation.NodeList = *apiReservation.NodeList
	}

	// Access lists
	if apiReservation.Users != nil {
		reservation.Users = strings.Split(*apiReservation.Users, ",")
	}
	if apiReservation.Accounts != nil {
		reservation.Accounts = strings.Split(*apiReservation.Accounts, ",")
	}
	if apiReservation.Groups != nil {
		reservation.Groups = strings.Split(*apiReservation.Groups, ",")
	}

	// Features
	if apiReservation.Features != nil {
		reservation.Features = strings.Split(*apiReservation.Features, ",")
	}

	// Licenses
	if apiReservation.Licenses != nil {
		// v0.0.42 has licenses as string, but common type expects map
		// Parse licenses string if needed (format: "lic1:4,lic2:8")
		reservation.Licenses = make(map[string]int32)
		// For now, skip parsing the license string
	}

	// TRES
	if apiReservation.Tres != nil {
		// v0.0.42 has TRES as string, but common type expects map
		// Parse TRES string if needed (format: "cpu=4,mem=1000")
		reservation.TRES = make(map[string]int64)
		// For now, skip parsing the TRES string
	}

	// Flags
	if apiReservation.Flags != nil && len(*apiReservation.Flags) > 0 {
		reservation.Flags = make([]types.ReservationFlag, len(*apiReservation.Flags))
		for i, flag := range *apiReservation.Flags {
			reservation.Flags[i] = types.ReservationFlag(flag)
		}
	}

	// Burst buffer
	if apiReservation.BurstBuffer != nil {
		reservation.BurstBuffer = *apiReservation.BurstBuffer
	}

	// Max start delay
	if apiReservation.MaxStartDelay != nil {
		reservation.MaxStartDelay = *apiReservation.MaxStartDelay
	}

	// Purge time - field doesn't exist in common types
	// Skip purge time handling

	// Watts - field doesn't exist in common types
	// Skip watts handling

	// Core specializations - field doesn't exist in common types
	// Skip core specializations handling

	return reservation, nil
}

// convertCommonReservationCreateToAPI converts common reservation create request to v0.0.42 API format
// Note: v0.0.42 API doesn't have a reservation POST endpoint, using placeholder
func (a *ReservationAdapter) convertCommonReservationCreateToAPI(req *types.ReservationCreate) (*api.V0042ReservationInfo, error) {
	reservation := &api.V0042ReservationInfo{
		Name: &req.Name,
	}

	// Partition
	if req.PartitionName != "" {
		reservation.Partition = &req.PartitionName
	}

	// Time fields
	if !req.StartTime.IsZero() {
		set := true
		num := req.StartTime.Unix()
		startTime := api.V0042Uint64NoValStruct{
			Set:    &set,
			Number: &num,
		}
		reservation.StartTime = &startTime
	}
	if req.EndTime != nil && !req.EndTime.IsZero() {
		set := true
		num := req.EndTime.Unix()
		endTime := api.V0042Uint64NoValStruct{
			Set:    &set,
			Number: &num,
		}
		reservation.EndTime = &endTime
	}

	// Duration (if specified, calculate end time)
	if req.Duration > 0 && req.StartTime.IsZero() {
		// If duration is specified but not start time, assume start now
		now := time.Now()
		set := true
		startNum := now.Unix()
		startTime := api.V0042Uint64NoValStruct{
			Set:    &set,
			Number: &startNum,
		}
		reservation.StartTime = &startTime

		endNum := now.Add(time.Duration(req.Duration) * time.Second).Unix()
		endTime := api.V0042Uint64NoValStruct{
			Set:    &set,
			Number: &endNum,
		}
		reservation.EndTime = &endTime
	} else if req.Duration > 0 && !req.StartTime.IsZero() {
		// If both duration and start time are specified, calculate end time
		set := true
		endNum := req.StartTime.Add(time.Duration(req.Duration) * time.Second).Unix()
		endTime := api.V0042Uint64NoValStruct{
			Set:    &set,
			Number: &endNum,
		}
		reservation.EndTime = &endTime
	}

	// Resources
	if req.NodeCount > 0 {
		nodeCount := int32(req.NodeCount)
		reservation.NodeCount = &nodeCount
	}
	if req.CoreCount > 0 {
		coreCount := int32(req.CoreCount)
		reservation.CoreCount = &coreCount
	}
	if req.NodeList != "" {
		reservation.NodeList = &req.NodeList
	}

	// Access lists
	if len(req.Users) > 0 {
		users := strings.Join(req.Users, ",")
		reservation.Users = &users
	}
	if len(req.Accounts) > 0 {
		accounts := strings.Join(req.Accounts, ",")
		reservation.Accounts = &accounts
	}
	if len(req.Groups) > 0 {
		groups := strings.Join(req.Groups, ",")
		reservation.Groups = &groups
	}

	// Features
	if len(req.Features) > 0 {
		features := strings.Join(req.Features, ",")
		reservation.Features = &features
	}

	// Licenses - convert map to string format
	if len(req.Licenses) > 0 {
		licenseStrs := make([]string, 0, len(req.Licenses))
		for lic, count := range req.Licenses {
			licenseStrs = append(licenseStrs, fmt.Sprintf("%s:%d", lic, count))
		}
		licenses := strings.Join(licenseStrs, ",")
		reservation.Licenses = &licenses
	}

	// TRES - convert map to string format
	if len(req.TRES) > 0 {
		tresStrs := make([]string, 0, len(req.TRES))
		for tres, count := range req.TRES {
			tresStrs = append(tresStrs, fmt.Sprintf("%s=%d", tres, count))
		}
		tres := strings.Join(tresStrs, ",")
		reservation.Tres = &tres
	}

	// Flags
	if len(req.Flags) > 0 {
		// Convert common ReservationFlag to API flags (which is []string)
		apiFlags := make([]string, len(req.Flags))
		for i, flag := range req.Flags {
			apiFlags[i] = string(flag)
		}
		reservation.Flags = &apiFlags
	}

	// Burst buffer
	if req.BurstBuffer != "" {
		reservation.BurstBuffer = &req.BurstBuffer
	}

	// Max start delay
	if req.MaxStartDelay > 0 {
		reservation.MaxStartDelay = &req.MaxStartDelay
	}

	// Watts - not in ReservationCreate type
	// Skip watts handling

	return reservation, nil
}
