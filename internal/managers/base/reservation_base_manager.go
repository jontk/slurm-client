// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package base

import (
	"fmt"
	"strings"
	"time"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/pkg/errors"
)

// ReservationBaseManager provides common reservation management functionality
type ReservationBaseManager struct {
	*CRUDManager
}

// NewReservationBaseManager creates a new reservation base manager
func NewReservationBaseManager(version string) *ReservationBaseManager {
	return &ReservationBaseManager{
		CRUDManager: NewCRUDManager(version, "Reservation"),
	}
}

// ValidateReservationCreate validates reservation creation data
// validateReservationTimeConstraints validates time-related fields
func (m *ReservationBaseManager) validateReservationTimeConstraints(reservation *types.ReservationCreate) error {
	if reservation.StartTime.IsZero() {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Start time is required",
			"reservation.StartTime", reservation.StartTime, nil,
		)
	}
	if reservation.Duration == 0 && (reservation.EndTime == nil || reservation.EndTime.IsZero()) {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Either duration or end time must be specified",
			"reservation.Duration", reservation.Duration, nil,
		)
	}
	if reservation.Duration > 0 && reservation.EndTime != nil && !reservation.EndTime.IsZero() {
		expectedEndTime := reservation.StartTime.Add(time.Duration(reservation.Duration) * time.Minute)
		if !expectedEndTime.Equal(*reservation.EndTime) {
			return errors.NewValidationError(
				errors.ErrorCodeValidationFailed,
				"Duration and end time are inconsistent",
				"reservation.Duration", reservation.Duration, nil,
			)
		}
	}
	if reservation.EndTime != nil && !reservation.EndTime.IsZero() && reservation.EndTime.Before(reservation.StartTime) {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"End time cannot be before start time",
			"reservation.EndTime", reservation.EndTime, nil,
		)
	}
	return nil
}

// validateReservationNumericFields validates numeric field constraints
func (m *ReservationBaseManager) validateReservationNumericFields(reservation *types.ReservationCreate) error {
	if err := m.ValidateNonNegative(int(reservation.CoreCount), "reservation.CoreCount"); err != nil {
		return err
	}
	if err := m.ValidateNonNegative(int(reservation.NodeCount), "reservation.NodeCount"); err != nil {
		return err
	}
	if err := m.ValidateNonNegative(int(reservation.MaxStartDelay), "reservation.MaxStartDelay"); err != nil {
		return err
	}
	if reservation.Watts < 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Watts must be non-negative",
			"reservation.Watts", reservation.Watts, nil,
		)
	}
	return nil
}

// validateReservationResources validates account/user/group requirements
func (m *ReservationBaseManager) validateReservationResources(reservation *types.ReservationCreate) error {
	if len(reservation.Accounts) == 0 && len(reservation.Users) == 0 && len(reservation.Groups) == 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"At least one of accounts, users, or groups must be specified",
			"reservation", reservation, nil,
		)
	}
	return nil
}

func (m *ReservationBaseManager) ValidateReservationCreate(reservation *types.ReservationCreate) error {
	if reservation == nil {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Reservation data is required",
			"reservation", reservation, nil,
		)
	}

	if err := m.ValidateResourceName(reservation.Name, "reservation name"); err != nil {
		return err
	}

	// Validate time constraints
	if err := m.validateReservationTimeConstraints(reservation); err != nil {
		return err
	}

	// Validate numeric fields
	if err := m.validateReservationNumericFields(reservation); err != nil {
		return err
	}

	// Validate resource requirements
	if err := m.validateReservationResources(reservation); err != nil {
		return err
	}

	// Validate flags if provided
	if len(reservation.Flags) > 0 {
		if err := m.ValidateReservationFlags(reservation.Flags); err != nil {
			return err
		}
	}

	// Validate TRES map
	if err := m.validateTRESMap(reservation.TRES, "reservation.TRES"); err != nil {
		return err
	}

	// Validate licenses map
	if err := m.validateLicensesMap(reservation.Licenses, "reservation.Licenses"); err != nil {
		return err
	}

	return nil
}

// ValidateReservationUpdate validates reservation update data
func (m *ReservationBaseManager) ValidateReservationUpdate(update *types.ReservationUpdate) error {
	if update == nil {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Update data is required",
			"update", update, nil,
		)
	}

	// Validate time constraints if provided
	if update.StartTime != nil && update.EndTime != nil {
		if update.EndTime.Before(*update.StartTime) {
			return errors.NewValidationError(
				errors.ErrorCodeValidationFailed,
				"End time cannot be before start time",
				"update.EndTime", update.EndTime, nil,
			)
		}
	}

	// Validate numeric fields if provided
	if update.CoreCount != nil {
		if err := m.ValidateNonNegative(int(*update.CoreCount), "update.CoreCount"); err != nil {
			return err
		}
	}
	if update.NodeCount != nil {
		if err := m.ValidateNonNegative(int(*update.NodeCount), "update.NodeCount"); err != nil {
			return err
		}
	}
	if update.MaxStartDelay != nil {
		if err := m.ValidateNonNegative(int(*update.MaxStartDelay), "update.MaxStartDelay"); err != nil {
			return err
		}
	}

	// Validate watts if provided
	if update.Watts != nil && *update.Watts < 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Watts must be non-negative",
			"update.Watts", *update.Watts, nil,
		)
	}

	// Validate flags if provided
	if len(update.Flags) > 0 {
		if err := m.ValidateReservationFlags(update.Flags); err != nil {
			return err
		}
	}

	// Validate TRES map
	if err := m.validateTRESMap(update.TRES, "update.TRES"); err != nil {
		return err
	}

	// Validate licenses map
	if err := m.validateLicensesMap(update.Licenses, "update.Licenses"); err != nil {
		return err
	}

	return nil
}

// ValidateReservationFlags validates reservation flags
func (m *ReservationBaseManager) ValidateReservationFlags(flags []types.ReservationFlag) error {
	validFlags := []types.ReservationFlag{
		types.ReservationFlagMaintenance,
		types.ReservationFlagOverlap,
		types.ReservationFlagIgnoreJobs,
		types.ReservationFlagDaily,
		types.ReservationFlagWeekly,
		types.ReservationFlagAnyNodes,
		types.ReservationFlagStatic,
		types.ReservationFlagPartNodes,
		types.ReservationFlagFirstCores,
		types.ReservationFlagTimeFLoat,
		types.ReservationFlagReplace,
		types.ReservationFlagLicenseOnly,
		types.ReservationFlagNoLicenseOnly,
		types.ReservationFlagPrompt,
		types.ReservationFlagNoHoldJobsAfter,
		types.ReservationFlagPurgeCompleted,
		types.ReservationFlagWeekend,
		types.ReservationFlagFlexible,
		types.ReservationFlagMagneticCores,
		types.ReservationFlagForce,
		types.ReservationFlagSkipProlog,
		types.ReservationFlagSkipEpilog,
	}

	for _, flag := range flags {
		found := false
		for _, validFlag := range validFlags {
			if flag == validFlag {
				found = true
				break
			}
		}
		if !found {
			return errors.NewValidationError(
				errors.ErrorCodeValidationFailed,
				fmt.Sprintf("Invalid reservation flag: %s", flag),
				"flags", flags, nil,
			)
		}
	}

	return nil
}

// ApplyReservationDefaults applies default values to reservation create request
func (m *ReservationBaseManager) ApplyReservationDefaults(reservation *types.ReservationCreate) *types.ReservationCreate {
	// Ensure slice fields are initialized
	if reservation.Accounts == nil {
		reservation.Accounts = []string{}
	}
	if reservation.Users == nil {
		reservation.Users = []string{}
	}
	if reservation.Groups == nil {
		reservation.Groups = []string{}
	}
	if reservation.Features == nil {
		reservation.Features = []string{}
	}
	if reservation.Flags == nil {
		reservation.Flags = []types.ReservationFlag{}
	}

	// Initialize TRES map
	if reservation.TRES == nil {
		reservation.TRES = make(map[string]int64)
	}

	// Initialize licenses map
	if reservation.Licenses == nil {
		reservation.Licenses = make(map[string]int32)
	}

	return reservation
}

// FilterReservationList applies filtering to a reservation list
func (m *ReservationBaseManager) FilterReservationList(items []types.Reservation, opts *types.ReservationListOptions) []types.Reservation {
	if opts == nil {
		return items
	}

	filtered := make([]types.Reservation, 0, len(items))
	for _, reservation := range items {
		if m.matchesReservationFilters(reservation, opts) {
			filtered = append(filtered, reservation)
		}
	}

	return filtered
}

// matchesReservationFilters checks if a reservation matches the given filters
func (m *ReservationBaseManager) matchesReservationFilters(reservation types.Reservation, opts *types.ReservationListOptions) bool {
	return m.checkStringFilter(opts.Names, reservation.Name, true) &&
		m.checkStateFilter(opts.States, reservation.State) &&
		m.checkStringSliceFilter(opts.Accounts, reservation.Accounts, true) &&
		m.checkStringSliceFilter(opts.Users, reservation.Users, true) &&
		m.checkStringFilter(opts.Partitions, reservation.PartitionName, true) &&
		m.checkReservationTimeRange(reservation, opts)
}

func (m *ReservationBaseManager) checkStringFilter(filters []string, value string, caseInsensitive bool) bool {
	if len(filters) == 0 {
		return true
	}
	for _, filter := range filters {
		if caseInsensitive {
			if strings.EqualFold(value, filter) {
				return true
			}
		} else if value == filter {
			return true
		}
	}
	return false
}

func (m *ReservationBaseManager) checkStateFilter(states []types.ReservationState, resState types.ReservationState) bool {
	if len(states) == 0 {
		return true
	}
	for _, state := range states {
		if resState == state {
			return true
		}
	}
	return false
}

func (m *ReservationBaseManager) checkStringSliceFilter(filters []string, values []string, caseInsensitive bool) bool {
	if len(filters) == 0 {
		return true
	}
	for _, filter := range filters {
		for _, value := range values {
			if caseInsensitive {
				if strings.EqualFold(value, filter) {
					return true
				}
			} else if value == filter {
				return true
			}
		}
	}
	return false
}

func (m *ReservationBaseManager) checkReservationTimeRange(reservation types.Reservation, opts *types.ReservationListOptions) bool {
	if opts.StartTime != nil && reservation.EndTime.Before(*opts.StartTime) {
		return false
	}
	if opts.EndTime != nil && reservation.StartTime.After(*opts.EndTime) {
		return false
	}
	return true
}

// validateTRESMap validates that TRES values are non-negative
func (m *ReservationBaseManager) validateTRESMap(tres map[string]int64, fieldName string) error {
	for resource, value := range tres {
		if value < 0 {
			return errors.NewValidationError(
				errors.ErrorCodeValidationFailed,
				fmt.Sprintf("TRES value for %s must be non-negative", resource),
				fieldName, tres, nil,
			)
		}
	}
	return nil
}

// validateLicensesMap validates that license values are non-negative
func (m *ReservationBaseManager) validateLicensesMap(licenses map[string]int32, fieldName string) error {
	for license, count := range licenses {
		if count < 0 {
			return errors.NewValidationError(
				errors.ErrorCodeValidationFailed,
				fmt.Sprintf("License count for %s must be non-negative", license),
				fieldName, licenses, nil,
			)
		}
	}
	return nil
}
