// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package factory

import (
	"context"
	"fmt"
	"time"

	types "github.com/jontk/slurm-client/api"
	"github.com/jontk/slurm-client/internal/adapters/common"
	"github.com/jontk/slurm-client/pkg/errors"
)

// extendedUserManager wraps the base adapter to add extended methods
type extendedUserManager struct {
	adapter            common.UserAdapter
	accountAdapter     common.AccountAdapter
	associationAdapter common.AssociationAdapter
}

// GetUserAccounts retrieves all accounts that a user is associated with
func (m *extendedUserManager) GetUserAccounts(ctx context.Context, userName string) ([]*types.UserAccount, error) {
	if userName == "" {
		return nil, fmt.Errorf("user name required")
	}

	// Get associations for this user
	associations, err := getAssociationsForUser(ctx, m.associationAdapter, userName)
	if err != nil {
		return nil, fmt.Errorf("failed to get associations: %w", err)
	}

	var accounts []*types.UserAccount
	seen := make(map[string]bool)

	for _, assoc := range associations {
		accountName := derefString(assoc.Account)
		if accountName == "" || seen[accountName] {
			continue
		}
		seen[accountName] = true

		userAccount := convertAssociationToUserAccount(&assoc)
		if userAccount != nil {
			accounts = append(accounts, userAccount)
		}
	}

	return accounts, nil
}

// GetUserQuotas retrieves aggregated quota information for a user across all accounts
func (m *extendedUserManager) GetUserQuotas(ctx context.Context, userName string) (*types.UserQuota, error) {
	if userName == "" {
		return nil, fmt.Errorf("user name required")
	}

	// Get associations for this user
	associations, err := getAssociationsForUser(ctx, m.associationAdapter, userName)
	if err != nil {
		return nil, fmt.Errorf("failed to get associations: %w", err)
	}

	if len(associations) == 0 {
		return nil, errors.NewSlurmError(errors.ErrorCodeResourceNotFound, fmt.Sprintf("no associations found for user %s", userName))
	}

	// Aggregate quotas
	quota := aggregateUserQuotas(associations, userName)

	return quota, nil
}

// GetUserDefaultAccount retrieves the user's default account
func (m *extendedUserManager) GetUserDefaultAccount(ctx context.Context, userName string) (*types.Account, error) {
	if userName == "" {
		return nil, fmt.Errorf("user name required")
	}

	// Get the user to find their default account
	user, err := m.adapter.Get(ctx, userName)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.NewSlurmError(errors.ErrorCodeResourceNotFound, fmt.Sprintf("user %s not found", userName))
	}

	// Check if user has a default account set
	if user.Default == nil || user.Default.Account == nil || *user.Default.Account == "" {
		// Try to find default from associations
		associations, err := getAssociationsForUser(ctx, m.associationAdapter, userName)
		if err != nil {
			return nil, fmt.Errorf("failed to get associations: %w", err)
		}

		for _, assoc := range associations {
			if assoc.IsDefault != nil && *assoc.IsDefault {
				if assoc.Account != nil {
					return m.accountAdapter.Get(ctx, *assoc.Account)
				}
			}
		}

		return nil, errors.NewSlurmError(errors.ErrorCodeResourceNotFound, "no default account found for user")
	}

	// Get the account details
	return m.accountAdapter.Get(ctx, *user.Default.Account)
}

// GetUserFairShare retrieves fairshare information for a user
func (m *extendedUserManager) GetUserFairShare(ctx context.Context, userName string) (*types.UserFairShare, error) {
	if userName == "" {
		return nil, fmt.Errorf("user name required")
	}

	// Get associations for this user
	associations, err := getAssociationsForUser(ctx, m.associationAdapter, userName)
	if err != nil {
		return nil, fmt.Errorf("failed to get associations: %w", err)
	}

	if len(associations) == 0 {
		return nil, errors.NewSlurmError(errors.ErrorCodeResourceNotFound, fmt.Sprintf("no associations found for user %s", userName))
	}

	// Use the first association (or the default one if we can find it)
	var targetAssoc *types.Association
	for i := range associations {
		assoc := &associations[i]
		if assoc.IsDefault != nil && *assoc.IsDefault {
			targetAssoc = assoc
			break
		}
	}
	if targetAssoc == nil {
		targetAssoc = &associations[0]
	}

	fairShare := extractUserFairShare(targetAssoc)
	if fairShare == nil {
		return nil, errors.NewSlurmError(errors.ErrorCodeResourceNotFound, "could not extract fairshare data")
	}

	return fairShare, nil
}

// CalculateJobPriority calculates the estimated job priority for a user
func (m *extendedUserManager) CalculateJobPriority(ctx context.Context, userName string, jobSubmission *types.JobSubmission) (*types.JobPriorityInfo, error) {
	if userName == "" {
		return nil, fmt.Errorf("user name required")
	}

	// Get user's fairshare info
	fairShare, err := m.GetUserFairShare(ctx, userName)
	if err != nil {
		// If we can't get fairshare, continue with default values
		fairShare = &types.UserFairShare{
			UserName:        userName,
			FairShareFactor: 0.5, // Default middle value
		}
	}

	// Determine account and partition from job submission
	account := ""
	partition := ""
	qos := ""
	if jobSubmission != nil {
		account = jobSubmission.Account
		partition = jobSubmission.Partition
	}

	if account == "" && fairShare.Account != "" {
		account = fairShare.Account
	}

	// Build priority info
	priorityInfo := &types.JobPriorityInfo{
		UserName:        userName,
		Account:         account,
		Partition:       partition,
		QoS:             qos,
		EligibleTime:    time.Now(),
		EstimatedStart:  time.Now().Add(5 * time.Minute), // Rough estimate
		PriorityTier:    "normal",
		PositionInQueue: 1, // Placeholder
	}

	// Calculate priority factors
	// These are estimates based on typical SLURM configurations
	factors := &types.JobPriorityFactors{
		Age:       0,
		FairShare: int(fairShare.FairShareFactor * 1000), // Scale to int
		JobSize:   100,                                   // Default job size factor
		Partition: 100,                                   // Default partition factor
		QoS:       100,                                   // Default QoS factor
		TRES:      0,                                     // Default TRES factor
		Site:      0,
		Nice:      0,
		Assoc:     0,
	}

	// Calculate total priority (simplified formula)
	factors.Total = factors.Age + factors.FairShare + factors.JobSize + factors.Partition + factors.QoS

	priorityInfo.Factors = factors
	priorityInfo.Priority = factors.Total

	// Determine priority tier based on total
	switch {
	case factors.Total >= 800:
		priorityInfo.PriorityTier = "high"
	case factors.Total >= 400:
		priorityInfo.PriorityTier = "normal"
	default:
		priorityInfo.PriorityTier = "low"
	}

	return priorityInfo, nil
}

// ValidateUserAccountAccess validates whether a user has access to a specific account
func (m *extendedUserManager) ValidateUserAccountAccess(ctx context.Context, userName, accountName string) (*types.UserAccessValidation, error) {
	if userName == "" {
		return nil, fmt.Errorf("user name required")
	}
	if accountName == "" {
		return nil, fmt.Errorf("account name required")
	}

	validation := &types.UserAccessValidation{
		UserName:       userName,
		AccountName:    accountName,
		HasAccess:      false,
		AccessLevel:    "none",
		Permissions:    []string{},
		ValidationTime: time.Now(),
	}

	// Get associations for this user
	associations, err := getAssociationsForUser(ctx, m.associationAdapter, userName)
	if err != nil {
		validation.Reason = fmt.Sprintf("failed to get associations: %v", err)
		return validation, nil
	}

	// Check if user has an association with this account
	for _, assoc := range associations {
		if assoc.Account != nil && *assoc.Account == accountName {
			validation.HasAccess = true
			validation.AccessLevel = "user"
			validation.Permissions = []string{"submit_jobs", "view_jobs"}
			validation.ValidFrom = time.Now()
			validation.Association = convertAssociationToUserAccountAssociation(&assoc)

			// Check if user is coordinator
			for _, flag := range assoc.Flags {
				if flag == types.AssociationDefaultFlagsUsersarecoords {
					validation.AccessLevel = "coordinator"
					validation.Permissions = append(validation.Permissions, "manage_users", "manage_accounts")
				}
			}

			// Extract quota limits
			if assoc.Max != nil {
				validation.QuotaLimits = &types.UserAccountQuota{
					AccountName: accountName,
				}
				if assoc.Max.Jobs != nil {
					if assoc.Max.Jobs.Active != nil {
						validation.QuotaLimits.MaxJobs = int(*assoc.Max.Jobs.Active)
					}
					if assoc.Max.Jobs.Total != nil {
						validation.QuotaLimits.MaxSubmitJobs = int(*assoc.Max.Jobs.Total)
					}
				}
			}

			return validation, nil
		}
	}

	validation.Reason = "no association found for user with this account"
	return validation, nil
}

// GetUserAccountAssociations retrieves all user-account associations for a user
func (m *extendedUserManager) GetUserAccountAssociations(ctx context.Context, userName string, opts *types.ListUserAccountAssociationsOptions) ([]*types.UserAccountAssociation, error) {
	if userName == "" {
		return nil, fmt.Errorf("user name required")
	}

	// Get associations for this user
	associations, err := getAssociationsForUser(ctx, m.associationAdapter, userName)
	if err != nil {
		return nil, fmt.Errorf("failed to get associations: %w", err)
	}

	var result []*types.UserAccountAssociation

	for _, assoc := range associations {
		userAssoc := convertAssociationToUserAccountAssociation(&assoc)
		if userAssoc == nil {
			continue
		}

		// Apply filters from opts
		if opts != nil {
			// Filter by accounts
			if len(opts.Accounts) > 0 {
				found := false
				for _, acc := range opts.Accounts {
					if userAssoc.AccountName == acc {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}

			// Filter by clusters
			if len(opts.Clusters) > 0 {
				found := false
				for _, cluster := range opts.Clusters {
					if userAssoc.Cluster == cluster {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}

			// Filter by partitions
			if len(opts.Partitions) > 0 {
				found := false
				for _, part := range opts.Partitions {
					if userAssoc.Partition == part {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}

			// Filter by active only
			if opts.ActiveOnly && !userAssoc.IsActive {
				continue
			}

			// Filter by default only
			if opts.DefaultOnly && !userAssoc.IsDefault {
				continue
			}

			// Filter by coordinator roles
			if opts.CoordinatorRoles && !userAssoc.IsCoordinator {
				continue
			}
		}

		result = append(result, userAssoc)
	}

	// Apply pagination
	if opts != nil {
		if opts.Offset > 0 && opts.Offset < len(result) {
			result = result[opts.Offset:]
		}
		if opts.Limit > 0 && opts.Limit < len(result) {
			result = result[:opts.Limit]
		}
	}

	return result, nil
}

// GetBulkUserAccounts retrieves accounts for multiple users in a single call
func (m *extendedUserManager) GetBulkUserAccounts(ctx context.Context, userNames []string) (map[string][]*types.UserAccount, error) {
	if len(userNames) == 0 {
		return nil, fmt.Errorf("at least one user name required")
	}

	result := make(map[string][]*types.UserAccount)

	// Get all associations at once (more efficient than per-user calls)
	opts := &types.AssociationListOptions{
		Users: userNames,
		Limit: 10000,
	}
	assocList, err := m.associationAdapter.List(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get associations: %w", err)
	}

	// Initialize result map with empty slices for all requested users
	for _, userName := range userNames {
		result[userName] = []*types.UserAccount{}
	}

	// Group associations by user
	for _, assoc := range assocList.Associations {
		userName := assoc.User
		if userName == "" {
			continue
		}

		// Only process if this user was requested
		if _, exists := result[userName]; !exists {
			continue
		}

		userAccount := convertAssociationToUserAccount(&assoc)
		if userAccount != nil {
			result[userName] = append(result[userName], userAccount)
		}
	}

	return result, nil
}

// GetBulkAccountUsers retrieves users for multiple accounts in a single call
func (m *extendedUserManager) GetBulkAccountUsers(ctx context.Context, accountNames []string) (map[string][]*types.UserAccountAssociation, error) {
	if len(accountNames) == 0 {
		return nil, fmt.Errorf("at least one account name required")
	}

	result := make(map[string][]*types.UserAccountAssociation)

	// Get all associations at once (more efficient than per-account calls)
	opts := &types.AssociationListOptions{
		Accounts: accountNames,
		Limit:    10000,
	}
	assocList, err := m.associationAdapter.List(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get associations: %w", err)
	}

	// Initialize result map with empty slices for all requested accounts
	for _, accountName := range accountNames {
		result[accountName] = []*types.UserAccountAssociation{}
	}

	// Group associations by account
	for _, assoc := range assocList.Associations {
		accountName := derefString(assoc.Account)
		if accountName == "" {
			continue
		}

		// Skip account-level associations (no user)
		if assoc.User == "" {
			continue
		}

		// Only process if this account was requested
		if _, exists := result[accountName]; !exists {
			continue
		}

		userAssoc := convertAssociationToUserAccountAssociation(&assoc)
		if userAssoc != nil {
			result[accountName] = append(result[accountName], userAssoc)
		}
	}

	return result, nil
}
