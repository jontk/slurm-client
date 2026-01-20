// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package base

import (
	"fmt"
	"strings"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/pkg/errors"
)

// AccountBaseManager provides common account management functionality
type AccountBaseManager struct {
	*CRUDManager
}

// NewAccountBaseManager creates a new account base manager
func NewAccountBaseManager(version string) *AccountBaseManager {
	return &AccountBaseManager{
		CRUDManager: NewCRUDManager(version, "Account"),
	}
}

// ValidateAccountCreate validates account creation data
func (m *AccountBaseManager) ValidateAccountCreate(account *types.AccountCreate) error {
	if account == nil {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Account data is required",
			"account", account, nil,
		)
	}

	if err := m.ValidateResourceName(account.Name, "account name"); err != nil {
		return err
	}

	// Validate numeric resource limits
	if err := m.ValidateNonNegative(int(account.Priority), "account.Priority"); err != nil {
		return err
	}
	if err := m.ValidateNonNegative(int(account.FairShare), "account.FairShare"); err != nil {
		return err
	}
	if err := m.ValidateNonNegative(int(account.SharesRaw), "account.SharesRaw"); err != nil {
		return err
	}

	// Validate job limits
	if err := m.ValidateNonNegative(int(account.MaxJobs), "account.MaxJobs"); err != nil {
		return err
	}
	if err := m.ValidateNonNegative(int(account.MaxJobsPerUser), "account.MaxJobsPerUser"); err != nil {
		return err
	}
	if err := m.ValidateNonNegative(int(account.MaxSubmitJobs), "account.MaxSubmitJobs"); err != nil {
		return err
	}

	// Validate time limits
	if err := m.ValidateNonNegative(int(account.MaxWallTime), "account.MaxWallTime"); err != nil {
		return err
	}
	if err := m.ValidateNonNegative(int(account.MaxCPUTime), "account.MaxCPUTime"); err != nil {
		return err
	}

	// Validate resource limits
	if err := m.ValidateNonNegative(int(account.MaxNodes), "account.MaxNodes"); err != nil {
		return err
	}
	if err := m.ValidateNonNegative(int(account.MaxCPUs), "account.MaxCPUs"); err != nil {
		return err
	}
	if account.MaxMemory < 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Max memory must be non-negative",
			"account.MaxMemory", account.MaxMemory, nil,
		)
	}

	// Validate group limits
	if err := m.ValidateNonNegative(int(account.GrpJobs), "account.GrpJobs"); err != nil {
		return err
	}
	if err := m.ValidateNonNegative(int(account.GrpNodes), "account.GrpNodes"); err != nil {
		return err
	}
	if err := m.ValidateNonNegative(int(account.GrpCPUs), "account.GrpCPUs"); err != nil {
		return err
	}
	if account.GrpMemory < 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Group memory must be non-negative",
			"account.GrpMemory", account.GrpMemory, nil,
		)
	}

	// Validate TRES maps
	if err := m.validateTRESMap(account.GrpTRES, "account.GrpTRES"); err != nil {
		return err
	}
	if err := m.validateTRESMap(account.MaxTRES, "account.MaxTRES"); err != nil {
		return err
	}

	return nil
}

// ValidateAccountUpdate validates account update data
func (m *AccountBaseManager) ValidateAccountUpdate(update *types.AccountUpdate) error {
	if update == nil {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Update data is required",
			"update", update, nil,
		)
	}

	// Validate numeric fields if provided
	if update.Priority != nil {
		if err := m.ValidateNonNegative(int(*update.Priority), "update.Priority"); err != nil {
			return err
		}
	}
	if update.FairShare != nil {
		if err := m.ValidateNonNegative(int(*update.FairShare), "update.FairShare"); err != nil {
			return err
		}
	}
	if update.SharesRaw != nil {
		if err := m.ValidateNonNegative(int(*update.SharesRaw), "update.SharesRaw"); err != nil {
			return err
		}
	}

	// Validate job limits if provided
	if update.MaxJobs != nil {
		if err := m.ValidateNonNegative(int(*update.MaxJobs), "update.MaxJobs"); err != nil {
			return err
		}
	}
	if update.MaxJobsPerUser != nil {
		if err := m.ValidateNonNegative(int(*update.MaxJobsPerUser), "update.MaxJobsPerUser"); err != nil {
			return err
		}
	}

	// Validate memory limits if provided
	if update.MaxMemory != nil && *update.MaxMemory < 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Max memory must be non-negative",
			"update.MaxMemory", *update.MaxMemory, nil,
		)
	}
	if update.GrpMemory != nil && *update.GrpMemory < 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Group memory must be non-negative",
			"update.GrpMemory", *update.GrpMemory, nil,
		)
	}

	// Validate TRES maps if provided
	if err := m.validateTRESMap(update.GrpTRES, "update.GrpTRES"); err != nil {
		return err
	}
	if err := m.validateTRESMap(update.MaxTRES, "update.MaxTRES"); err != nil {
		return err
	}

	return nil
}

// ApplyAccountDefaults applies default values to account create request
func (m *AccountBaseManager) ApplyAccountDefaults(account *types.AccountCreate) *types.AccountCreate {
	// Ensure slice fields are initialized
	if account.Coordinators == nil {
		account.Coordinators = []string{}
	}
	if account.QoSList == nil {
		account.QoSList = []string{}
	}
	if account.AllowedPartitions == nil {
		account.AllowedPartitions = []string{}
	}

	// Initialize TRES maps
	if account.GrpTRES == nil {
		account.GrpTRES = make(map[string]int64)
	}
	if account.GrpTRESMins == nil {
		account.GrpTRESMins = make(map[string]int64)
	}
	if account.GrpTRESRunMins == nil {
		account.GrpTRESRunMins = make(map[string]int64)
	}
	if account.MaxTRES == nil {
		account.MaxTRES = make(map[string]int64)
	}
	if account.MaxTRESPerNode == nil {
		account.MaxTRESPerNode = make(map[string]int64)
	}
	if account.MinTRES == nil {
		account.MinTRES = make(map[string]int64)
	}

	return account
}

// FilterAccountList applies filtering to an account list
func (m *AccountBaseManager) FilterAccountList(items []types.Account, opts *types.AccountListOptions) []types.Account {
	if opts == nil {
		return items
	}

	filtered := make([]types.Account, 0, len(items))
	for _, account := range items {
		if m.matchesAccountFilters(account, opts) {
			filtered = append(filtered, account)
		}
	}

	return filtered
}

// matchesAccountFilters checks if an account matches the given filters
func (m *AccountBaseManager) matchesAccountFilters(account types.Account, opts *types.AccountListOptions) bool {
	// Filter by names
	if len(opts.Names) > 0 {
		found := false
		for _, name := range opts.Names {
			if strings.EqualFold(account.Name, name) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Filter by descriptions
	if len(opts.Descriptions) > 0 {
		found := false
		for _, desc := range opts.Descriptions {
			if strings.Contains(strings.ToLower(account.Description), strings.ToLower(desc)) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Filter by organizations
	if len(opts.Organizations) > 0 {
		found := false
		for _, org := range opts.Organizations {
			if strings.EqualFold(account.Organization, org) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Filter deleted accounts
	if !opts.WithDeleted && account.Deleted {
		return false
	}

	// Filter by update time
	if opts.UpdateTime != nil {
		// This would require API support to track update times
		// For now, we'll accept all items
	}

	return true
}

// validateTRESMap validates that TRES values are non-negative
func (m *AccountBaseManager) validateTRESMap(tres map[string]int64, fieldName string) error {
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
