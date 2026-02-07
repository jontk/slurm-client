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

// extendedAccountManager wraps the base adapter to add extended methods
type extendedAccountManager struct {
	adapter            common.AccountAdapter
	associationAdapter common.AssociationAdapter
}

// GetAccountHierarchy retrieves the complete account hierarchy starting from the specified root account
func (m *extendedAccountManager) GetAccountHierarchy(ctx context.Context, rootAccount string) (*types.AccountHierarchy, error) {
	if rootAccount == "" {
		return nil, fmt.Errorf("account name required")
	}

	// Get the root account details
	account, err := m.adapter.Get(ctx, rootAccount)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, errors.NewSlurmError(errors.ErrorCodeResourceNotFound, fmt.Sprintf("account %s not found", rootAccount))
	}

	// Get all associations to build the hierarchy
	associations, err := getAllAssociations(ctx, m.associationAdapter)
	if err != nil {
		return nil, fmt.Errorf("failed to get associations: %w", err)
	}

	// Build the tree
	tree := buildAccountTree(associations)

	// Create the hierarchy starting from root
	hierarchy, err := m.buildHierarchyNode(ctx, rootAccount, tree, 0, []string{})
	if err != nil {
		return nil, fmt.Errorf("failed to build hierarchy: %w", err)
	}

	return hierarchy, nil
}

// buildHierarchyNode recursively builds the hierarchy tree
func (m *extendedAccountManager) buildHierarchyNode(ctx context.Context, accountName string, tree map[string]*accountNode, level int, path []string) (*types.AccountHierarchy, error) {
	node, exists := tree[accountName]

	// Get account details
	account, err := m.adapter.Get(ctx, accountName)
	if err != nil {
		return nil, fmt.Errorf("failed to get account %s: %w", accountName, err)
	}

	hierarchy := &types.AccountHierarchy{
		Account:          account,
		Level:            level,
		Path:             append(path, accountName),
		ChildAccounts:    make([]*types.AccountHierarchy, 0),
		TotalUsers:       0,
		TotalSubAccounts: 0,
	}

	if !exists {
		return hierarchy, nil
	}

	// Count users
	hierarchy.TotalUsers = len(node.users)

	// Build child hierarchies
	for _, childNode := range node.children {
		if childNode == nil {
			continue
		}
		// Find child account name from associations
		for name, n := range tree {
			if n == childNode {
				childHierarchy, err := m.buildHierarchyNode(ctx, name, tree, level+1, hierarchy.Path)
				if err != nil {
					return nil, err
				}
				hierarchy.ChildAccounts = append(hierarchy.ChildAccounts, childHierarchy)
				hierarchy.TotalSubAccounts++
				hierarchy.TotalUsers += childHierarchy.TotalUsers
				hierarchy.TotalSubAccounts += childHierarchy.TotalSubAccounts
				break
			}
		}
	}

	// Extract quota from first association if available
	if len(node.assocs) > 0 {
		hierarchy.AggregateQuota = extractQuotaFromAssociation(node.assocs[0])
		hierarchy.AggregateUsage = extractUsageFromAssociation(node.assocs[0])
	}

	return hierarchy, nil
}

// GetParentAccounts retrieves all parent accounts for the specified account
func (m *extendedAccountManager) GetParentAccounts(ctx context.Context, accountName string) ([]*types.Account, error) {
	if accountName == "" {
		return nil, fmt.Errorf("account name required")
	}

	// Get associations for this account to find its lineage
	associations, err := getAssociationsForAccount(ctx, m.associationAdapter, accountName)
	if err != nil {
		return nil, fmt.Errorf("failed to get associations: %w", err)
	}

	if len(associations) == 0 {
		return []*types.Account{}, nil
	}

	// Parse the lineage from the first association
	var parents []*types.Account
	seen := make(map[string]bool)

	for _, assoc := range associations {
		lineage := parseLineage(assoc.Lineage)
		for _, parentName := range lineage {
			if parentName == accountName || seen[parentName] {
				continue
			}
			seen[parentName] = true

			// Fetch the parent account
			parent, err := m.adapter.Get(ctx, parentName)
			if err != nil {
				continue // Skip if we can't fetch the parent
			}
			if parent != nil {
				parents = append(parents, parent)
			}
		}
	}

	return parents, nil
}

// GetChildAccounts retrieves child accounts up to the specified depth
func (m *extendedAccountManager) GetChildAccounts(ctx context.Context, accountName string, depth int) ([]*types.Account, error) {
	if accountName == "" {
		return nil, fmt.Errorf("account name required")
	}

	if depth <= 0 {
		depth = 1 // Default to immediate children only
	}

	// Get all associations
	associations, err := getAllAssociations(ctx, m.associationAdapter)
	if err != nil {
		return nil, fmt.Errorf("failed to get associations: %w", err)
	}

	var children []*types.Account
	seen := make(map[string]bool)
	currentLevel := []string{accountName}

	for level := 0; level < depth && len(currentLevel) > 0; level++ {
		var nextLevel []string
		for _, parent := range currentLevel {
			childNames := findChildAccounts(associations, parent)
			for _, childName := range childNames {
				if seen[childName] {
					continue
				}
				seen[childName] = true

				// Fetch the child account
				child, err := m.adapter.Get(ctx, childName)
				if err != nil {
					continue
				}
				if child != nil {
					children = append(children, child)
					nextLevel = append(nextLevel, childName)
				}
			}
		}
		currentLevel = nextLevel
	}

	return children, nil
}

// GetAccountQuotas retrieves quota information for the specified account
func (m *extendedAccountManager) GetAccountQuotas(ctx context.Context, accountName string) (*types.AccountQuota, error) {
	if accountName == "" {
		return nil, fmt.Errorf("account name required")
	}

	// Get associations for this account
	associations, err := getAssociationsForAccount(ctx, m.associationAdapter, accountName)
	if err != nil {
		return nil, fmt.Errorf("failed to get associations: %w", err)
	}

	if len(associations) == 0 {
		return nil, errors.NewSlurmError(errors.ErrorCodeResourceNotFound, fmt.Sprintf("no associations found for account %s", accountName))
	}

	// Aggregate quotas from all associations (account-level, not user-level)
	quota := &types.AccountQuota{
		GrpTRES:            make(map[string]int),
		GrpTRESUsed:        make(map[string]int),
		GrpTRESMinutes:     make(map[string]int),
		GrpTRESMinutesUsed: make(map[string]int),
		MaxTRES:            make(map[string]int),
		MaxTRESUsed:        make(map[string]int),
		MaxTRESPerUser:     make(map[string]int),
		LastUpdated:        time.Now(),
	}

	// Find the account-level association (user == "")
	for _, assoc := range associations {
		if assoc.User == "" {
			// This is the account-level association
			extracted := extractQuotaFromAssociation(&assoc)
			if extracted != nil {
				return extracted, nil
			}
			break
		}
	}

	// If no account-level association, aggregate from user associations
	for _, assoc := range associations {
		extracted := extractQuotaFromAssociation(&assoc)
		if extracted == nil {
			continue
		}

		// Take the maximum values
		if extracted.MaxJobs > quota.MaxJobs {
			quota.MaxJobs = extracted.MaxJobs
		}
		if extracted.MaxJobsPerUser > quota.MaxJobsPerUser {
			quota.MaxJobsPerUser = extracted.MaxJobsPerUser
		}
		if extracted.MaxWallTime > quota.MaxWallTime {
			quota.MaxWallTime = extracted.MaxWallTime
		}

		// Merge TRES limits
		for k, v := range extracted.GrpTRES {
			if v > quota.GrpTRES[k] {
				quota.GrpTRES[k] = v
			}
		}
		for k, v := range extracted.MaxTRES {
			if v > quota.MaxTRES[k] {
				quota.MaxTRES[k] = v
			}
		}
	}

	return quota, nil
}

// GetAccountQuotaUsage retrieves usage information for the specified account
func (m *extendedAccountManager) GetAccountQuotaUsage(ctx context.Context, accountName string, timeframe string) (*types.AccountUsage, error) {
	if accountName == "" {
		return nil, fmt.Errorf("account name required")
	}

	// Get associations for this account
	associations, err := getAssociationsForAccount(ctx, m.associationAdapter, accountName)
	if err != nil {
		return nil, fmt.Errorf("failed to get associations: %w", err)
	}

	if len(associations) == 0 {
		return nil, errors.NewSlurmError(errors.ErrorCodeResourceNotFound, fmt.Sprintf("no associations found for account %s", accountName))
	}

	// Aggregate usage from all associations
	usage := &types.AccountUsage{
		AccountName: accountName,
		StartTime:   time.Now().AddDate(0, -1, 0), // Default to last month
		EndTime:     time.Now(),
		TRESUsage:   make(map[string]float64),
	}

	// Parse timeframe if provided
	if timeframe != "" {
		switch timeframe {
		case "day":
			usage.StartTime = time.Now().AddDate(0, 0, -1)
		case "week":
			usage.StartTime = time.Now().AddDate(0, 0, -7)
		case "month":
			usage.StartTime = time.Now().AddDate(0, -1, 0)
		case "year":
			usage.StartTime = time.Now().AddDate(-1, 0, 0)
		}
	}

	// Count unique users
	userSet := make(map[string]bool)

	for _, assoc := range associations {
		if assoc.User != "" {
			userSet[assoc.User] = true
		}

		extracted := extractUsageFromAssociation(&assoc)
		if extracted != nil {
			usage.CPUSeconds += extracted.CPUSeconds
			usage.JobCount += extracted.JobCount

			for k, v := range extracted.TRESUsage {
				usage.TRESUsage[k] += v
			}
		}
	}

	usage.UserCount = int32(len(userSet))

	return usage, nil
}

// GetAccountUsers retrieves all users associated with the specified account
func (m *extendedAccountManager) GetAccountUsers(ctx context.Context, accountName string, opts *types.ListAccountUsersOptions) ([]*types.UserAccountAssociation, error) {
	if accountName == "" {
		return nil, fmt.Errorf("account name required")
	}

	// Get associations for this account
	associations, err := getAssociationsForAccount(ctx, m.associationAdapter, accountName)
	if err != nil {
		return nil, fmt.Errorf("failed to get associations: %w", err)
	}

	var users []*types.UserAccountAssociation

	for _, assoc := range associations {
		// Skip account-level associations (no user)
		if assoc.User == "" {
			continue
		}

		userAssoc := convertAssociationToUserAccountAssociation(&assoc)
		if userAssoc == nil {
			continue
		}

		// Apply filters from opts
		if opts != nil {
			if opts.ActiveOnly && !userAssoc.IsActive {
				continue
			}
			if opts.CoordinatorsOnly && !userAssoc.IsCoordinator {
				continue
			}
		}

		users = append(users, userAssoc)
	}

	// Apply pagination
	if opts != nil {
		if opts.Offset > 0 && opts.Offset < len(users) {
			users = users[opts.Offset:]
		}
		if opts.Limit > 0 && opts.Limit < len(users) {
			users = users[:opts.Limit]
		}
	}

	return users, nil
}

// ValidateUserAccess validates whether a user has access to an account
func (m *extendedAccountManager) ValidateUserAccess(ctx context.Context, userName, accountName string) (*types.UserAccessValidation, error) {
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

	// Get associations for this account
	associations, err := getAssociationsForAccount(ctx, m.associationAdapter, accountName)
	if err != nil {
		validation.Reason = fmt.Sprintf("failed to get associations: %v", err)
		return validation, nil
	}

	// Check if user has an association with this account
	for _, assoc := range associations {
		if assoc.User == userName {
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

			return validation, nil
		}
	}

	validation.Reason = "no association found for user with this account"
	return validation, nil
}

// GetAccountUsersWithPermissions retrieves users with specific permissions
func (m *extendedAccountManager) GetAccountUsersWithPermissions(ctx context.Context, accountName string, permissions []string) ([]*types.UserAccountAssociation, error) {
	if accountName == "" {
		return nil, fmt.Errorf("account name required")
	}

	// First, get the account to check coordinators
	account, err := m.adapter.Get(ctx, accountName)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, errors.NewSlurmError(errors.ErrorCodeResourceNotFound, fmt.Sprintf("account %s not found", accountName))
	}

	// Get all users for this account
	allUsers, err := m.GetAccountUsers(ctx, accountName, nil)
	if err != nil {
		return nil, err
	}

	// Check permissions for each user
	var result []*types.UserAccountAssociation

	// Create a permission set for quick lookup
	permSet := make(map[string]bool)
	for _, perm := range permissions {
		permSet[perm] = true
	}

	for _, user := range allUsers {
		// Check if user is a coordinator (has extra permissions)
		if user.IsCoordinator {
			// Coordinators have manage permissions
			if permSet["manage_users"] || permSet["manage_accounts"] || permSet["coordinator"] {
				result = append(result, user)
				continue
			}
		}

		// Check if user has any of the requested permissions
		for _, userPerm := range user.Permissions {
			if permSet[userPerm] {
				result = append(result, user)
				break
			}
		}
	}

	return result, nil
}

// GetAccountFairShare retrieves fairshare information for the specified account
func (m *extendedAccountManager) GetAccountFairShare(ctx context.Context, accountName string) (*types.AccountFairShare, error) {
	if accountName == "" {
		return nil, fmt.Errorf("account name required")
	}

	// Get associations for this account
	associations, err := getAssociationsForAccount(ctx, m.associationAdapter, accountName)
	if err != nil {
		return nil, fmt.Errorf("failed to get associations: %w", err)
	}

	if len(associations) == 0 {
		return nil, errors.NewSlurmError(errors.ErrorCodeResourceNotFound, fmt.Sprintf("no associations found for account %s", accountName))
	}

	// Find the account-level association first
	for _, assoc := range associations {
		if assoc.User == "" {
			fairShare := extractFairShareFromAssociation(&assoc, accountName)
			if fairShare != nil {
				// Count users
				for _, a := range associations {
					if a.User != "" {
						fairShare.UserCount++
						fairShare.ActiveUsers++ // Assume all are active for now
					}
				}
				return fairShare, nil
			}
		}
	}

	// Fall back to first user association for shares info
	fairShare := extractFairShareFromAssociation(&associations[0], accountName)
	if fairShare != nil {
		// Count users
		userSet := make(map[string]bool)
		for _, a := range associations {
			if a.User != "" {
				userSet[a.User] = true
			}
		}
		fairShare.UserCount = len(userSet)
		fairShare.ActiveUsers = len(userSet)
	}

	return fairShare, nil
}

// GetFairShareHierarchy retrieves the complete fairshare hierarchy from the specified root
func (m *extendedAccountManager) GetFairShareHierarchy(ctx context.Context, rootAccount string) (*types.FairShareHierarchy, error) {
	if rootAccount == "" {
		return nil, fmt.Errorf("account name required")
	}

	// Get all associations
	associations, err := getAllAssociations(ctx, m.associationAdapter)
	if err != nil {
		return nil, fmt.Errorf("failed to get associations: %w", err)
	}

	hierarchy := &types.FairShareHierarchy{
		RootAccount: rootAccount,
		LastUpdate:  time.Now(),
		Accounts:    make([]*types.AccountFairShare, 0),
		Users:       make([]*types.UserFairShare, 0),
	}

	// Build account fair share list
	accountShares := make(map[string]*types.AccountFairShare)
	for _, assoc := range associations {
		accountName := derefString(assoc.Account)
		if accountName == "" {
			continue
		}

		// Create or update account fair share
		if _, exists := accountShares[accountName]; !exists {
			accountShares[accountName] = extractFairShareFromAssociation(&assoc, accountName)
		}

		// Accumulate shares
		if assoc.SharesRaw != nil {
			accountShares[accountName].Shares += int(*assoc.SharesRaw)
			hierarchy.TotalShares += int(*assoc.SharesRaw)
		}

		// Add user fair share
		if assoc.User != "" {
			userFS := extractUserFairShare(&assoc)
			if userFS != nil {
				hierarchy.Users = append(hierarchy.Users, userFS)
			}
			accountShares[accountName].UserCount++
		}
	}

	// Convert map to slice
	for _, acctFS := range accountShares {
		hierarchy.Accounts = append(hierarchy.Accounts, acctFS)
	}

	// Determine cluster from first association
	if len(associations) > 0 && associations[0].Cluster != nil {
		hierarchy.Cluster = *associations[0].Cluster
	}

	// Build the tree structure
	hierarchy.Tree = m.buildFairShareTree(rootAccount, associations)

	return hierarchy, nil
}

// buildFairShareTree builds a FairShareNode tree from associations
func (m *extendedAccountManager) buildFairShareTree(rootAccount string, associations []types.Association) *types.FairShareNode {
	// Build a map of account -> associations
	accountMap := make(map[string][]types.Association)
	for _, assoc := range associations {
		if assoc.Account == nil {
			continue
		}
		accountMap[*assoc.Account] = append(accountMap[*assoc.Account], assoc)
	}

	// Build the root node
	return m.buildFairShareNode(rootAccount, accountMap, associations, 0)
}

func (m *extendedAccountManager) buildFairShareNode(accountName string, accountMap map[string][]types.Association, allAssocs []types.Association, level int) *types.FairShareNode {
	node := &types.FairShareNode{
		Name:     accountName,
		Account:  accountName,
		Level:    level,
		Children: make([]*types.FairShareNode, 0),
	}

	// Get shares from associations
	if assocs, exists := accountMap[accountName]; exists {
		for _, assoc := range assocs {
			if assoc.SharesRaw != nil {
				node.Shares += int(*assoc.SharesRaw)
			}
		}
	}

	// Find parent
	for _, assoc := range allAssocs {
		if assoc.Account != nil && *assoc.Account == accountName && assoc.ParentAccount != nil {
			node.Parent = *assoc.ParentAccount
			break
		}
	}

	// Find children
	childAccounts := findChildAccounts(allAssocs, accountName)
	for _, childName := range childAccounts {
		childNode := m.buildFairShareNode(childName, accountMap, allAssocs, level+1)
		node.Children = append(node.Children, childNode)
	}

	return node
}
