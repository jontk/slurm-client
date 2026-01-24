// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package base

import (
	"fmt"
	"strings"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/pkg/errors"
)

// AssociationBaseManager provides common association management functionality
type AssociationBaseManager struct {
	*CRUDManager
}

// NewAssociationBaseManager creates a new association base manager
func NewAssociationBaseManager(version string) *AssociationBaseManager {
	return &AssociationBaseManager{
		CRUDManager: NewCRUDManager(version, "Association"),
	}
}

// ValidateAssociationCreate validates association creation data
func (m *AssociationBaseManager) ValidateAssociationCreate(association *types.AssociationCreate) error {
	if association == nil {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Association data is required",
			"association", association, nil,
		)
	}

	if err := m.ValidateResourceName(association.AccountName, "Association name"); err != nil {
		return err
	}

	if err := m.ValidateResourceName(association.Cluster, "cluster name"); err != nil {
		return err
	}

	// User name is optional for account associations
	if association.UserName != "" {
		if err := m.ValidateResourceName(association.UserName, "user name"); err != nil {
			return err
		}
	}

	// Validate numeric resource limits
	if err := m.ValidateNonNegative(int(association.SharesRaw), "association.SharesRaw"); err != nil {
		return err
	}
	if err := m.ValidateNonNegative(int(association.Priority), "association.Priority"); err != nil {
		return err
	}

	// Validate job limits
	if err := m.ValidateNonNegative(int(association.MaxJobs), "association.MaxJobs"); err != nil {
		return err
	}
	if err := m.ValidateNonNegative(int(association.MaxJobsAccrue), "association.MaxJobsAccrue"); err != nil {
		return err
	}
	if err := m.ValidateNonNegative(int(association.MaxSubmitJobs), "association.MaxSubmitJobs"); err != nil {
		return err
	}

	// Validate time limits
	if err := m.ValidateNonNegative(int(association.MaxWallTime), "association.MaxWallTime"); err != nil {
		return err
	}
	if err := m.ValidateNonNegative(int(association.MaxCPUTime), "association.MaxCPUTime"); err != nil {
		return err
	}

	// Validate resource limits
	if err := m.ValidateNonNegative(int(association.MaxNodes), "association.MaxNodes"); err != nil {
		return err
	}
	if err := m.ValidateNonNegative(int(association.MaxCPUs), "association.MaxCPUs"); err != nil {
		return err
	}
	if association.MaxMemory < 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Max memory must be non-negative",
			"association.MaxMemory", association.MaxMemory, nil,
		)
	}

	// Validate group limits
	if err := m.ValidateNonNegative(int(association.GrpJobs), "association.GrpJobs"); err != nil {
		return err
	}
	if err := m.ValidateNonNegative(int(association.GrpJobsAccrue), "association.GrpJobsAccrue"); err != nil {
		return err
	}
	if err := m.ValidateNonNegative(int(association.GrpNodes), "association.GrpNodes"); err != nil {
		return err
	}
	if err := m.ValidateNonNegative(int(association.GrpCPUs), "association.GrpCPUs"); err != nil {
		return err
	}
	if association.GrpMemory < 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Group memory must be non-negative",
			"association.GrpMemory", association.GrpMemory, nil,
		)
	}

	// Validate CPU run mins
	if association.GrpCPURunMins < 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Group CPU run mins must be non-negative",
			"association.GrpCPURunMins", association.GrpCPURunMins, nil,
		)
	}

	// Validate TRES maps
	if err := m.validateTRESMap(association.GrpTRES, "association.GrpTRES"); err != nil {
		return err
	}
	if err := m.validateTRESMap(association.MaxTRES, "association.MaxTRES"); err != nil {
		return err
	}
	if err := m.validateTRESMap(association.MaxTRESMins, "association.MaxTRESMins"); err != nil {
		return err
	}
	if err := m.validateTRESMap(association.MinTRES, "association.MinTRES"); err != nil {
		return err
	}

	return nil
}

// ValidateAssociationUpdate validates association update data
func (m *AssociationBaseManager) ValidateAssociationUpdate(update *types.AssociationUpdate) error {
	if update == nil {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Update data is required",
			"update", update, nil,
		)
	}

	// Validate numeric fields if provided
	if update.SharesRaw != nil {
		if err := m.ValidateNonNegative(int(*update.SharesRaw), "update.SharesRaw"); err != nil {
			return err
		}
	}
	if update.Priority != nil {
		if err := m.ValidateNonNegative(int(*update.Priority), "update.Priority"); err != nil {
			return err
		}
	}

	// Validate job limits if provided
	if update.MaxJobs != nil {
		if err := m.ValidateNonNegative(int(*update.MaxJobs), "update.MaxJobs"); err != nil {
			return err
		}
	}
	if update.MaxJobsAccrue != nil {
		if err := m.ValidateNonNegative(int(*update.MaxJobsAccrue), "update.MaxJobsAccrue"); err != nil {
			return err
		}
	}
	if update.MaxSubmitJobs != nil {
		if err := m.ValidateNonNegative(int(*update.MaxSubmitJobs), "update.MaxSubmitJobs"); err != nil {
			return err
		}
	}

	// Validate time limits if provided
	if update.MaxWallTime != nil {
		if err := m.ValidateNonNegative(int(*update.MaxWallTime), "update.MaxWallTime"); err != nil {
			return err
		}
	}
	if update.MaxCPUTime != nil {
		if err := m.ValidateNonNegative(int(*update.MaxCPUTime), "update.MaxCPUTime"); err != nil {
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

	// Validate CPU run mins if provided
	if update.GrpCPURunMins != nil && *update.GrpCPURunMins < 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Group CPU run mins must be non-negative",
			"update.GrpCPURunMins", *update.GrpCPURunMins, nil,
		)
	}

	// Validate TRES maps if provided
	if err := m.validateTRESMap(update.GrpTRES, "update.GrpTRES"); err != nil {
		return err
	}
	if err := m.validateTRESMap(update.MaxTRES, "update.MaxTRES"); err != nil {
		return err
	}
	if err := m.validateTRESMap(update.MaxTRESMins, "update.MaxTRESMins"); err != nil {
		return err
	}
	if err := m.validateTRESMap(update.MinTRES, "update.MinTRES"); err != nil {
		return err
	}

	return nil
}

// ApplyAssociationDefaults applies default values to association create request
func (m *AssociationBaseManager) ApplyAssociationDefaults(association *types.AssociationCreate) *types.AssociationCreate {
	// Ensure slice fields are initialized
	if association.QoSList == nil {
		association.QoSList = []string{}
	}

	// Initialize TRES maps
	if association.GrpTRES == nil {
		association.GrpTRES = make(map[string]int64)
	}
	if association.GrpTRESMins == nil {
		association.GrpTRESMins = make(map[string]int64)
	}
	if association.GrpTRESRunMins == nil {
		association.GrpTRESRunMins = make(map[string]int64)
	}
	if association.MaxTRES == nil {
		association.MaxTRES = make(map[string]int64)
	}
	if association.MaxTRESPerNode == nil {
		association.MaxTRESPerNode = make(map[string]int64)
	}
	if association.MaxTRESMins == nil {
		association.MaxTRESMins = make(map[string]int64)
	}
	if association.MinTRES == nil {
		association.MinTRES = make(map[string]int64)
	}

	return association
}

// FilterAssociationList applies filtering to an association list
func (m *AssociationBaseManager) FilterAssociationList(items []types.Association, opts *types.AssociationListOptions) []types.Association {
	if opts == nil {
		return items
	}

	filtered := make([]types.Association, 0, len(items))
	for _, association := range items {
		if m.matchesAssociationFilters(association, opts) {
			filtered = append(filtered, association)
		}
	}

	return filtered
}

// matchesAssociationFilters checks if an association matches the given filters
func (m *AssociationBaseManager) matchesAssociationFilters(association types.Association, opts *types.AssociationListOptions) bool {
	// Filter by accounts
	if len(opts.Accounts) > 0 {
		found := false
		for _, account := range opts.Accounts {
			if strings.EqualFold(association.AccountName, account) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Filter by clusters
	if len(opts.Clusters) > 0 {
		found := false
		for _, cluster := range opts.Clusters {
			if strings.EqualFold(association.Cluster, cluster) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Filter by users
	if len(opts.Users) > 0 {
		found := false
		for _, user := range opts.Users {
			if strings.EqualFold(association.UserName, user) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Filter by partitions
	if len(opts.Partitions) > 0 {
		found := false
		for _, partition := range opts.Partitions {
			if strings.EqualFold(association.Partition, partition) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Filter only defaults
	if opts.OnlyDefaults && !association.IsDefault {
		return false
	}

	// Filter deleted associations
	if !opts.WithDeleted && association.Deleted {
		return false
	}

	// Filter by update time
	// This would require API support to track update times
	// For now, we'll accept all items
	if opts.UpdateTime != nil {
	}

	return true
}

// validateTRESMap validates that TRES values are non-negative
func (m *AssociationBaseManager) validateTRESMap(tres map[string]int64, fieldName string) error {
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
