package base

import (
	"fmt"
	"strings"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/pkg/errors"
)

// UserBaseManager provides common user management functionality
type UserBaseManager struct {
	*CRUDManager
}

// NewUserBaseManager creates a new user base manager
func NewUserBaseManager(version string) *UserBaseManager {
	return &UserBaseManager{
		CRUDManager: NewCRUDManager(version, "User"),
	}
}

// ValidateUserCreate validates user creation data
func (m *UserBaseManager) ValidateUserCreate(user *types.UserCreate) error {
	if user == nil {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"User data is required",
			"user", user, nil,
		)
	}

	if err := m.ValidateResourceName(user.Name, "user.Name"); err != nil {
		return err
	}

	// Validate UID if provided
	if user.UID < 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"UID must be non-negative",
			"user.UID", user.UID, nil,
		)
	}

	// Validate admin level
	if user.AdminLevel != "" {
		if err := m.ValidateAdminLevel(user.AdminLevel); err != nil {
			return err
		}
	}

	// Validate numeric resource limits
	if err := m.ValidateNonNegative(int(user.MaxJobs), "user.MaxJobs"); err != nil {
		return err
	}
	if err := m.ValidateNonNegative(int(user.MaxJobsPerAccount), "user.MaxJobsPerAccount"); err != nil {
		return err
	}
	if err := m.ValidateNonNegative(int(user.MaxSubmitJobs), "user.MaxSubmitJobs"); err != nil {
		return err
	}

	// Validate time limits
	if err := m.ValidateNonNegative(int(user.MaxWallTime), "user.MaxWallTime"); err != nil {
		return err
	}
	if err := m.ValidateNonNegative(int(user.MaxCPUTime), "user.MaxCPUTime"); err != nil {
		return err
	}

	// Validate resource limits
	if err := m.ValidateNonNegative(int(user.MaxNodes), "user.MaxNodes"); err != nil {
		return err
	}
	if err := m.ValidateNonNegative(int(user.MaxCPUs), "user.MaxCPUs"); err != nil {
		return err
	}
	if user.MaxMemory < 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Max memory must be non-negative",
			"user.MaxMemory", user.MaxMemory, nil,
		)
	}

	// Validate group limits
	if err := m.ValidateNonNegative(int(user.GrpJobs), "user.GrpJobs"); err != nil {
		return err
	}
	if err := m.ValidateNonNegative(int(user.GrpNodes), "user.GrpNodes"); err != nil {
		return err
	}
	if err := m.ValidateNonNegative(int(user.GrpCPUs), "user.GrpCPUs"); err != nil {
		return err
	}
	if user.GrpMemory < 0 {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Group memory must be non-negative",
			"user.GrpMemory", user.GrpMemory, nil,
		)
	}

	// Validate TRES maps
	if err := m.validateTRESMap(user.GrpTRES, "user.GrpTRES"); err != nil {
		return err
	}
	if err := m.validateTRESMap(user.MaxTRES, "user.MaxTRES"); err != nil {
		return err
	}

	return nil
}

// ValidateUserUpdate validates user update data
func (m *UserBaseManager) ValidateUserUpdate(update *types.UserUpdate) error {
	if update == nil {
		return errors.NewValidationError(
			errors.ErrorCodeValidationFailed,
			"Update data is required",
			"update", update, nil,
		)
	}

	// Validate admin level if provided
	if update.AdminLevel != nil {
		if err := m.ValidateAdminLevel(*update.AdminLevel); err != nil {
			return err
		}
	}

	// Validate numeric fields if provided
	if update.MaxJobs != nil {
		if err := m.ValidateNonNegative(int(*update.MaxJobs), "update.MaxJobs"); err != nil {
			return err
		}
	}
	if update.MaxJobsPerAccount != nil {
		if err := m.ValidateNonNegative(int(*update.MaxJobsPerAccount), "update.MaxJobsPerAccount"); err != nil {
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

	// Validate TRES maps if provided
	if err := m.validateTRESMap(update.GrpTRES, "update.GrpTRES"); err != nil {
		return err
	}
	if err := m.validateTRESMap(update.MaxTRES, "update.MaxTRES"); err != nil {
		return err
	}

	return nil
}

// ValidateAdminLevel validates user admin level
func (m *UserBaseManager) ValidateAdminLevel(level types.AdminLevel) error {
	validLevels := []types.AdminLevel{
		types.AdminLevelNone,
		types.AdminLevelOperator,
		types.AdminLevelAdministrator,
	}

	for _, validLevel := range validLevels {
		if level == validLevel {
			return nil
		}
	}

	return errors.NewValidationError(
		errors.ErrorCodeValidationFailed,
		fmt.Sprintf("Invalid admin level: %s", level),
		"adminLevel", level, nil,
	)
}

// ApplyUserDefaults applies default values to user create request
func (m *UserBaseManager) ApplyUserDefaults(user *types.UserCreate) *types.UserCreate {
	// Apply default admin level if not set
	if user.AdminLevel == "" {
		user.AdminLevel = types.AdminLevelNone
	}

	// Ensure slice fields are initialized
	if user.Accounts == nil {
		user.Accounts = []string{}
	}
	if user.QoSList == nil {
		user.QoSList = []string{}
	}
	if user.WCKeys == nil {
		user.WCKeys = []string{}
	}

	// Initialize TRES maps
	if user.GrpTRES == nil {
		user.GrpTRES = make(map[string]int64)
	}
	if user.GrpTRESMins == nil {
		user.GrpTRESMins = make(map[string]int64)
	}
	if user.GrpTRESRunMins == nil {
		user.GrpTRESRunMins = make(map[string]int64)
	}
	if user.MaxTRES == nil {
		user.MaxTRES = make(map[string]int64)
	}
	if user.MaxTRESPerNode == nil {
		user.MaxTRESPerNode = make(map[string]int64)
	}
	if user.MinTRES == nil {
		user.MinTRES = make(map[string]int64)
	}

	return user
}

// FilterUserList applies filtering to a user list
func (m *UserBaseManager) FilterUserList(items []types.User, opts *types.UserListOptions) []types.User {
	if opts == nil {
		return items
	}

	filtered := make([]types.User, 0, len(items))
	for _, user := range items {
		if m.matchesUserFilters(user, opts) {
			filtered = append(filtered, user)
		}
	}

	return filtered
}

// matchesUserFilters checks if a user matches the given filters
func (m *UserBaseManager) matchesUserFilters(user types.User, opts *types.UserListOptions) bool {
	// Filter by names
	if len(opts.Names) > 0 {
		found := false
		for _, name := range opts.Names {
			if strings.EqualFold(user.Name, name) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Filter by default account
	if opts.DefaultAccount != "" {
		if !strings.EqualFold(user.DefaultAccount, opts.DefaultAccount) {
			return false
		}
	}

	// Filter deleted users
	if !opts.WithDeleted && user.Deleted {
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
func (m *UserBaseManager) validateTRESMap(tres map[string]int64, fieldName string) error {
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