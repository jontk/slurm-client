// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package base

import (
	"strings"
	"time"

	types "github.com/jontk/slurm-client/api"
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
	hasDuration := reservation.Duration != nil && *reservation.Duration > 0
	hasEndTime := !reservation.EndTime.IsZero()
	if !hasDuration && !hasEndTime {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Either duration or end time must be specified",
			"reservation.Duration", reservation.Duration, nil,
		)
	}
	if hasDuration && hasEndTime {
		expectedEndTime := reservation.StartTime.Add(time.Duration(*reservation.Duration) * time.Minute)
		if !expectedEndTime.Equal(reservation.EndTime) {
			return errors.NewValidationError(
				errors.ErrorCodeValidationFailed,
				"Duration and end time are inconsistent",
				"reservation.Duration", reservation.Duration, nil,
			)
		}
	}
	if hasEndTime && reservation.EndTime.Before(reservation.StartTime) {
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
	if reservation.CoreCount != nil {
		if err := m.ValidateNonNegative(int(*reservation.CoreCount), "reservation.CoreCount"); err != nil {
			return err
		}
	}
	if reservation.NodeCount != nil {
		if err := m.ValidateNonNegative(int(*reservation.NodeCount), "reservation.NodeCount"); err != nil {
			return err
		}
	}
	if reservation.MaxStartDelay != nil {
		if err := m.ValidateNonNegative(int(*reservation.MaxStartDelay), "reservation.MaxStartDelay"); err != nil {
			return err
		}
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
	name := ""
	if reservation.Name != nil {
		name = *reservation.Name
	}
	if err := m.ValidateResourceName(name, "reservation name"); err != nil {
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
	// Validate flags if provided (FlagsValue from generated code)
	// Note: Flags are now []FlagsValue, validation is optional
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
	// Note: Flags validation is optional since FlagsValue is generated from OpenAPI spec
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
	if reservation.Flags == nil {
		reservation.Flags = []types.FlagsValue{}
	}
	if reservation.TRES == nil {
		reservation.TRES = []types.TRES{}
	}
	if reservation.Licenses == nil {
		reservation.Licenses = []string{}
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
	// Dereference pointer fields safely
	name := ""
	if reservation.Name != nil {
		name = *reservation.Name
	}
	partition := ""
	if reservation.Partition != nil {
		partition = *reservation.Partition
	}
	// Convert comma-separated strings to slices for filtering
	accounts := parseCommaSeparated(reservation.Accounts)
	users := parseCommaSeparated(reservation.Users)
	return m.checkStringFilter(opts.Names, name, true) &&
		m.checkStringSliceFilter(opts.Accounts, accounts, true) &&
		m.checkStringSliceFilter(opts.Users, users, true) &&
		m.checkStringFilter(opts.Partitions, partition, true) &&
		m.checkReservationTimeRange(reservation, opts)
}

// parseCommaSeparated converts a comma-separated string pointer to a slice
func parseCommaSeparated(s *string) []string {
	if s == nil || *s == "" {
		return nil
	}
	return strings.Split(*s, ",")
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
