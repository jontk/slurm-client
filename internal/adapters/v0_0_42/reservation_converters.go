package v0_0_42

import (
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
		reservation.Partition = *apiReservation.Partition
	}

	// Time fields
	if apiReservation.StartTime != nil && apiReservation.StartTime.Number > 0 {
		reservation.StartTime = time.Unix(int64(apiReservation.StartTime.Number), 0)
	}
	if apiReservation.EndTime != nil && apiReservation.EndTime.Number > 0 {
		reservation.EndTime = time.Unix(int64(apiReservation.EndTime.Number), 0)
	}

	// Resource counts
	if apiReservation.NodeCount != nil {
		reservation.NodeCount = uint32(*apiReservation.NodeCount)
	}
	if apiReservation.CoreCount != nil {
		reservation.CoreCount = uint32(*apiReservation.CoreCount)
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
		reservation.Features = *apiReservation.Features
	}

	// Licenses
	if apiReservation.Licenses != nil {
		reservation.Licenses = *apiReservation.Licenses
	}

	// TRES
	if apiReservation.Tres != nil {
		reservation.TRES = *apiReservation.Tres
	}

	// Flags
	if apiReservation.Flags != nil && len(*apiReservation.Flags) > 0 {
		reservation.Flags = *apiReservation.Flags
	}

	// Burst buffer
	if apiReservation.BurstBuffer != nil {
		reservation.BurstBuffer = *apiReservation.BurstBuffer
	}

	// Max start delay
	if apiReservation.MaxStartDelay != nil {
		reservation.MaxStartDelay = time.Duration(*apiReservation.MaxStartDelay) * time.Second
	}

	// Purge time
	if apiReservation.PurgeCompleted != nil && apiReservation.PurgeCompleted.Time != nil {
		purgeTime := time.Duration(apiReservation.PurgeCompleted.Time.Number) * time.Second
		reservation.PurgeTime = &purgeTime
	}

	// Watts
	if apiReservation.Watts != nil && apiReservation.Watts.Number > 0 {
		watts := uint32(apiReservation.Watts.Number)
		reservation.Watts = &watts
	}

	// Core specializations
	if apiReservation.CoreSpecializations != nil && len(*apiReservation.CoreSpecializations) > 0 {
		// Store as extra info since we don't have a direct field for this in common types
		reservation.Extra = make(map[string]string)
		reservation.Extra["core_specializations"] = "configured"
	}

	return reservation, nil
}

// convertCommonReservationCreateToAPI converts common reservation create request to v0.0.42 API format
// Note: v0.0.42 API doesn't have a reservation POST endpoint, using placeholder
func (a *ReservationAdapter) convertCommonReservationCreateToAPI(req *types.ReservationCreateRequest) (*api.V0042ReservationInfo, error) {
	apiReq := &api.V0042ReservationInfo{
		Reservations: &[]api.V0042ReservationInfo{
			{
				Name: &req.Name,
			},
		},
	}

	reservation := &(*apiReq.Reservations)[0]

	// Partition
	if req.Partition != "" {
		reservation.Partition = &req.Partition
	}

	// Time fields
	if !req.StartTime.IsZero() {
		startTime := api.V0042Uint64NoValStruct{
			Set:    true,
			Number: uint64(req.StartTime.Unix()),
		}
		reservation.StartTime = &startTime
	}
	if !req.EndTime.IsZero() {
		endTime := api.V0042Uint64NoValStruct{
			Set:    true,
			Number: uint64(req.EndTime.Unix()),
		}
		reservation.EndTime = &endTime
	}

	// Duration (if specified, calculate end time)
	if req.Duration > 0 && req.StartTime.IsZero() {
		// If duration is specified but not start time, assume start now
		now := time.Now()
		startTime := api.V0042Uint64NoValStruct{
			Set:    true,
			Number: uint64(now.Unix()),
		}
		reservation.StartTime = &startTime

		endTime := api.V0042Uint64NoValStruct{
			Set:    true,
			Number: uint64(now.Add(req.Duration).Unix()),
		}
		reservation.EndTime = &endTime
	} else if req.Duration > 0 && !req.StartTime.IsZero() {
		// If both duration and start time are specified, calculate end time
		endTime := api.V0042Uint64NoValStruct{
			Set:    true,
			Number: uint64(req.StartTime.Add(req.Duration).Unix()),
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
	if req.Features != "" {
		reservation.Features = &req.Features
	}

	// Licenses
	if req.Licenses != "" {
		reservation.Licenses = &req.Licenses
	}

	// TRES
	if req.TRES != "" {
		reservation.Tres = &req.TRES
	}

	// Flags
	if len(req.Flags) > 0 {
		flags := api.V0042ReservationFlags(req.Flags)
		reservation.Flags = &flags
	}

	// Burst buffer
	if req.BurstBuffer != "" {
		reservation.BurstBuffer = &req.BurstBuffer
	}

	// Max start delay
	if req.MaxStartDelay > 0 {
		maxStartDelay := int32(req.MaxStartDelay / time.Second)
		reservation.MaxStartDelay = &maxStartDelay
	}

	// Watts
	if req.Watts > 0 {
		watts := api.V0042Uint32NoValStruct{
			Set:    true,
			Number: uint64(req.Watts),
		}
		reservation.Watts = &watts
	}

	return apiReq, nil
}