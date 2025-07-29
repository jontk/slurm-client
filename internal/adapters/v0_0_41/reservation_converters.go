package v0_0_41

import (
	"fmt"
	"strings"
	"time"

	"github.com/jontk/slurm-client/internal/common/types"
	api "github.com/jontk/slurm-client/internal/api/v0_0_41"
)

// convertAPIReservationToCommon converts a v0.0.41 API Reservation to common Reservation type
func (a *ReservationAdapter) convertAPIReservationToCommon(apiRes interface{}) (*types.Reservation, error) {
	// Type assertion to handle the anonymous struct
	resData, ok := apiRes.(struct {
		Accounts    *string `json:"accounts,omitempty"`
		BurstBuffer *string `json:"burst_buffer,omitempty"`
		CoreCount   *int32  `json:"core_count,omitempty"`
		CoreSpecCnt *int32  `json:"core_spec_cnt,omitempty"`
		EndTime     *api.V0041OpenapiReservationRespReservationsEndTime `json:"end_time,omitempty"`
		Features    *string `json:"features,omitempty"`
		Flags       *[]api.V0041OpenapiReservationRespReservationsFlags `json:"flags,omitempty"`
		Groups      *string `json:"groups,omitempty"`
		Licenses    *string `json:"licenses,omitempty"`
		MaxStartDelay *api.V0041OpenapiReservationRespReservationsMaxStartDelay `json:"max_start_delay,omitempty"`
		Name        *string `json:"name,omitempty"`
		NodeCount   *int32  `json:"node_count,omitempty"`
		NodeList    *string `json:"node_list,omitempty"`
		Partition   *string `json:"partition,omitempty"`
		PurgeCompleted *struct {
			Time *api.V0041OpenapiReservationRespReservationsPurgeCompletedTime `json:"time,omitempty"`
		} `json:"purge_completed,omitempty"`
		StartTime   *api.V0041OpenapiReservationRespReservationsStartTime `json:"start_time,omitempty"`
		Watts       *api.V0041OpenapiReservationRespReservationsWatts `json:"watts,omitempty"`
		Tres        *string `json:"tres,omitempty"`
		Users       *string `json:"users,omitempty"`
	})
	if !ok {
		return nil, fmt.Errorf("unexpected reservation data type")
	}

	res := &types.Reservation{}

	// Basic fields
	if resData.Name != nil {
		res.Name = *resData.Name
	}
	if resData.NodeList != nil {
		res.Nodes = *resData.NodeList
	}
	if resData.NodeCount != nil {
		res.NodeCnt = uint32(*resData.NodeCount)
	}
	if resData.CoreCount != nil {
		res.CoreCnt = uint32(*resData.CoreCount)
	}
	if resData.Partition != nil {
		res.Partition = *resData.Partition
	}

	// Time fields
	if resData.StartTime != nil && resData.StartTime.Number != nil {
		res.StartTime = time.Unix(*resData.StartTime.Number, 0)
	}
	if resData.EndTime != nil && resData.EndTime.Number != nil {
		res.EndTime = time.Unix(*resData.EndTime.Number, 0)
	}

	// Access control
	if resData.Accounts != nil {
		res.Accounts = strings.Split(*resData.Accounts, ",")
	}
	if resData.Users != nil {
		res.Users = strings.Split(*resData.Users, ",")
	}
	if resData.Groups != nil {
		res.Groups = strings.Split(*resData.Groups, ",")
	}

	// Features and licenses
	if resData.Features != nil {
		res.Features = *resData.Features
	}
	if resData.Licenses != nil {
		res.Licenses = *resData.Licenses
	}

	// Flags
	if resData.Flags != nil {
		var flags []string
		for _, flag := range *resData.Flags {
			flags = append(flags, string(flag))
		}
		res.Flags = flags
	}

	// TRES
	if resData.Tres != nil {
		res.TRES = *resData.Tres
	}

	// Burst buffer
	if resData.BurstBuffer != nil {
		res.BurstBuffer = *resData.BurstBuffer
	}

	// Watts
	if resData.Watts != nil && resData.Watts.Number != nil {
		res.Watts = uint32(*resData.Watts.Number)
	}

	// Max start delay
	if resData.MaxStartDelay != nil && resData.MaxStartDelay.Number != nil {
		res.MaxStartDelay = uint32(*resData.MaxStartDelay.Number)
	}

	return res, nil
}