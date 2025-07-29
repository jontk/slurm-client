package v0_0_40

import (
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
	if apiReservation.Duration != nil && apiReservation.Duration.Number != nil {
		reservation.Duration = time.Duration(*apiReservation.Duration.Number) * time.Second
	}

	// Node information
	if apiReservation.NodeNames != nil {
		reservation.Nodes = strings.Split(*apiReservation.NodeNames, ",")
	}
	if apiReservation.NodeCount != nil {
		reservation.NodeCount = *apiReservation.NodeCount
	}

	// Users and accounts
	if apiReservation.Users != nil {
		reservation.Users = strings.Split(*apiReservation.Users, ",")
	}
	if apiReservation.Accounts != nil && len(*apiReservation.Accounts) > 0 {
		reservation.Accounts = *apiReservation.Accounts
	}

	// Flags
	if apiReservation.Flags != nil && len(*apiReservation.Flags) > 0 {
		reservation.Flags = make([]string, len(*apiReservation.Flags))
		for i, flag := range *apiReservation.Flags {
			reservation.Flags[i] = string(flag)
		}
	}

	// Features
	if apiReservation.Features != nil {
		reservation.Features = *apiReservation.Features
	}

	// Partition
	if apiReservation.Partition != nil {
		reservation.Partition = *apiReservation.Partition
	}

	// License
	if apiReservation.Licenses != nil {
		// Parse licenses if needed
	}

	// TRES
	if apiReservation.Tres != nil {
		reservation.TRES = *apiReservation.Tres
	}

	// Core count
	if apiReservation.CoreCount != nil && len(*apiReservation.CoreCount) > 0 {
		reservation.CoreCount = (*apiReservation.CoreCount)[0]
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

	if reservation.Duration > 0 {
		duration := api.V0040Uint32NoVal{
			Set:    boolPtr(true),
			Number: int32Ptr(int32(reservation.Duration.Seconds())),
		}
		apiReservation.Duration = &duration
	}

	// Node information
	if len(reservation.Nodes) > 0 {
		nodeNames := strings.Join(reservation.Nodes, ",")
		apiReservation.NodeNames = &nodeNames
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
		apiReservation.Accounts = &reservation.Accounts
	}

	// Flags
	if len(reservation.Flags) > 0 {
		flags := make([]api.V0040ReservationInfoFlags, len(reservation.Flags))
		for i, flag := range reservation.Flags {
			flags[i] = api.V0040ReservationInfoFlags(flag)
		}
		apiReservation.Flags = &flags
	}

	// Features
	if reservation.Features != "" {
		apiReservation.Features = &reservation.Features
	}

	// Partition
	if reservation.Partition != "" {
		apiReservation.Partition = &reservation.Partition
	}

	// TRES
	if reservation.TRES != "" {
		apiReservation.Tres = &reservation.TRES
	}

	// Core count
	if reservation.CoreCount > 0 {
		coreCount := []int32{reservation.CoreCount}
		apiReservation.CoreCount = &coreCount
	}

	return apiReservation, nil
}