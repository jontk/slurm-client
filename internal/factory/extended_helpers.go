// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package factory

import (
	"context"
	"strings"
	"time"

	types "github.com/jontk/slurm-client/api"
	"github.com/jontk/slurm-client/internal/adapters/common"
)

// parseLineage converts a lineage string like "/root/parent/child" to a slice ["root", "parent", "child"]
// The lineage represents the path from root account down to the current account.
func parseLineage(lineage *string) []string {
	if lineage == nil || *lineage == "" {
		return nil
	}
	// Remove leading slash and split by "/"
	path := strings.TrimPrefix(*lineage, "/")
	if path == "" {
		return nil
	}
	return strings.Split(path, "/")
}

// accountNode represents an internal node for building account trees
type accountNode struct {
	account  *types.Account
	children []*accountNode
	parent   *accountNode
	users    []string
	assocs   []*types.Association
}

// buildAccountTree creates a map of account names to their nodes from associations
func buildAccountTree(associations []types.Association) map[string]*accountNode {
	tree := make(map[string]*accountNode)

	// First pass: create nodes for all unique accounts
	for i := range associations {
		assoc := &associations[i]
		if assoc.Account == nil {
			continue
		}
		accountName := *assoc.Account
		if _, exists := tree[accountName]; !exists {
			tree[accountName] = &accountNode{
				children: make([]*accountNode, 0),
				users:    make([]string, 0),
				assocs:   make([]*types.Association, 0),
			}
		}
		tree[accountName].assocs = append(tree[accountName].assocs, assoc)
		if assoc.User != "" {
			tree[accountName].users = append(tree[accountName].users, assoc.User)
		}
	}

	// Second pass: establish parent-child relationships
	for i := range associations {
		assoc := &associations[i]
		if assoc.Account == nil || assoc.ParentAccount == nil {
			continue
		}
		accountName := *assoc.Account
		parentName := *assoc.ParentAccount
		if parentName == "" || parentName == accountName {
			continue
		}

		childNode := tree[accountName]
		parentNode, exists := tree[parentName]
		if !exists {
			parentNode = &accountNode{
				children: make([]*accountNode, 0),
				users:    make([]string, 0),
				assocs:   make([]*types.Association, 0),
			}
			tree[parentName] = parentNode
		}

		// Only add if not already a child
		found := false
		for _, c := range parentNode.children {
			if c == childNode {
				found = true
				break
			}
		}
		if !found {
			parentNode.children = append(parentNode.children, childNode)
			childNode.parent = parentNode
		}
	}

	return tree
}

// extractQuotaFromAssociation extracts quota information from an association
func extractQuotaFromAssociation(assoc *types.Association) *types.AccountQuota {
	if assoc == nil {
		return nil
	}

	quota := &types.AccountQuota{
		LastUpdated: time.Now(),
	}

	if assoc.Max != nil {
		if assoc.Max.Jobs != nil {
			if assoc.Max.Jobs.Active != nil {
				quota.MaxJobs = int(*assoc.Max.Jobs.Active)
			}
			if assoc.Max.Jobs.Per != nil {
				if assoc.Max.Jobs.Per.Count != nil {
					quota.MaxJobsPerUser = int(*assoc.Max.Jobs.Per.Count)
				}
				if assoc.Max.Jobs.Per.WallClock != nil {
					quota.MaxWallTime = int(*assoc.Max.Jobs.Per.WallClock)
				}
			}
		}

		// Extract TRES limits
		if assoc.Max.TRES != nil {
			quota.GrpTRES = make(map[string]int)
			quota.MaxTRES = make(map[string]int)

			if assoc.Max.TRES.Total != nil {
				for _, tres := range assoc.Max.TRES.Total {
					if tres.Name != nil && tres.Count != nil {
						quota.GrpTRES[*tres.Name] = int(*tres.Count)
					}
				}
			}

			if assoc.Max.TRES.Per != nil && assoc.Max.TRES.Per.Job != nil {
				for _, tres := range assoc.Max.TRES.Per.Job {
					if tres.Name != nil && tres.Count != nil {
						quota.MaxTRES[*tres.Name] = int(*tres.Count)
					}
				}
			}
		}
	}

	return quota
}

// extractUsageFromAssociation extracts usage information from association accounting
func extractUsageFromAssociation(assoc *types.Association) *types.AccountUsage {
	if assoc == nil {
		return nil
	}

	accountName := ""
	if assoc.Account != nil {
		accountName = *assoc.Account
	}

	usage := &types.AccountUsage{
		AccountName: accountName,
		StartTime:   time.Now().AddDate(0, -1, 0), // Default to last month
		EndTime:     time.Now(),
		TRESUsage:   make(map[string]float64),
	}

	// Process accounting records
	if len(assoc.Accounting) > 0 {
		for _, acct := range assoc.Accounting {
			if acct.Allocated != nil && acct.Allocated.Seconds != nil {
				usage.CPUSeconds += *acct.Allocated.Seconds
			}
		}
	}

	return usage
}

// convertAssociationToUserAccountAssociation converts an association to UserAccountAssociation
func convertAssociationToUserAccountAssociation(assoc *types.Association) *types.UserAccountAssociation {
	if assoc == nil {
		return nil
	}

	result := &types.UserAccountAssociation{
		UserName:    assoc.User,
		Cluster:     derefString(assoc.Cluster),
		Partition:   derefString(assoc.Partition),
		Role:        "user",
		Permissions: []string{},
		IsActive:    true,
		Created:     time.Now(),
		Modified:    time.Now(),
	}

	if assoc.Account != nil {
		result.AccountName = *assoc.Account
	}

	if assoc.IsDefault != nil && *assoc.IsDefault {
		result.IsDefault = true
	}

	if assoc.SharesRaw != nil {
		result.SharesRaw = int(*assoc.SharesRaw)
	}

	if assoc.Priority != nil {
		result.Priority = int(*assoc.Priority)
	}

	if assoc.QoS != nil {
		result.QoS = assoc.QoS
	}

	if assoc.Default != nil && assoc.Default.QoS != nil {
		result.DefaultQoS = *assoc.Default.QoS
	}

	// Check for coordinator flags
	for _, flag := range assoc.Flags {
		if flag == types.AssociationDefaultFlagsUsersarecoords {
			result.IsCoordinator = true
			result.Role = "coordinator"
			result.Permissions = append(result.Permissions, "manage_users", "manage_accounts")
		}
	}

	// Extract limits from Max
	if assoc.Max != nil && assoc.Max.Jobs != nil {
		if assoc.Max.Jobs.Active != nil {
			result.MaxJobs = int(*assoc.Max.Jobs.Active)
		}
		if assoc.Max.Jobs.Total != nil {
			result.MaxSubmitJobs = int(*assoc.Max.Jobs.Total)
		}
		if assoc.Max.Jobs.Per != nil && assoc.Max.Jobs.Per.WallClock != nil {
			result.MaxWallTime = int(*assoc.Max.Jobs.Per.WallClock)
		}
	}

	return result
}

// convertAssociationToUserAccount converts an association to UserAccount
func convertAssociationToUserAccount(assoc *types.Association) *types.UserAccount {
	if assoc == nil {
		return nil
	}

	result := &types.UserAccount{
		IsActive: true,
		Created:  time.Now(),
		Modified: time.Now(),
	}

	if assoc.Account != nil {
		result.AccountName = *assoc.Account
	}

	if assoc.Partition != nil {
		result.Partition = *assoc.Partition
	}

	if assoc.IsDefault != nil && *assoc.IsDefault {
		result.IsDefault = true
	}

	if assoc.Priority != nil {
		result.Priority = int(*assoc.Priority)
	}

	if len(assoc.QoS) > 0 {
		result.QoS = assoc.QoS[0]
	}

	if assoc.Default != nil && assoc.Default.QoS != nil {
		result.DefaultQoS = *assoc.Default.QoS
	}

	// Extract limits from Max
	if assoc.Max != nil && assoc.Max.Jobs != nil {
		if assoc.Max.Jobs.Active != nil {
			result.MaxJobs = int(*assoc.Max.Jobs.Active)
		}
		if assoc.Max.Jobs.Total != nil {
			result.MaxSubmitJobs = int(*assoc.Max.Jobs.Total)
		}
		if assoc.Max.Jobs.Per != nil && assoc.Max.Jobs.Per.WallClock != nil {
			result.MaxWallTime = int(*assoc.Max.Jobs.Per.WallClock)
		}
	}

	return result
}

// derefString safely dereferences a string pointer
func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// derefInt32 safely dereferences an int32 pointer
func derefInt32(i *int32) int32 {
	if i == nil {
		return 0
	}
	return *i
}

// derefUint32 safely dereferences a uint32 pointer
func derefUint32(i *uint32) uint32 {
	if i == nil {
		return 0
	}
	return *i
}

// getAssociationsForAccount retrieves all associations for a specific account
func getAssociationsForAccount(ctx context.Context, adapter common.AssociationAdapter, accountName string) ([]types.Association, error) {
	opts := &types.AssociationListOptions{
		Accounts: []string{accountName},
		Limit:    1000,
	}
	result, err := adapter.List(ctx, opts)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return []types.Association{}, nil
	}
	return result.Associations, nil
}

// getAssociationsForUser retrieves all associations for a specific user
func getAssociationsForUser(ctx context.Context, adapter common.AssociationAdapter, userName string) ([]types.Association, error) {
	opts := &types.AssociationListOptions{
		Users: []string{userName},
		Limit: 1000,
	}
	result, err := adapter.List(ctx, opts)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return []types.Association{}, nil
	}
	return result.Associations, nil
}

// getAllAssociations retrieves all associations (used for hierarchy building)
func getAllAssociations(ctx context.Context, adapter common.AssociationAdapter) ([]types.Association, error) {
	opts := &types.AssociationListOptions{
		Limit: 10000, // Large limit to get all associations
	}
	result, err := adapter.List(ctx, opts)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return []types.Association{}, nil
	}
	return result.Associations, nil
}

// findChildAccounts finds all direct children of an account
func findChildAccounts(associations []types.Association, parentAccountName string) []string {
	seen := make(map[string]bool)
	var children []string

	for _, assoc := range associations {
		if assoc.ParentAccount != nil && *assoc.ParentAccount == parentAccountName {
			if assoc.Account != nil && !seen[*assoc.Account] {
				seen[*assoc.Account] = true
				children = append(children, *assoc.Account)
			}
		}
	}

	return children
}

// extractFairShareFromAssociation extracts fair share info from association
func extractFairShareFromAssociation(assoc *types.Association, accountName string) *types.AccountFairShare {
	if assoc == nil {
		return nil
	}

	result := &types.AccountFairShare{
		AccountName: accountName,
		Cluster:     derefString(assoc.Cluster),
		Parent:      derefString(assoc.ParentAccount),
		Created:     time.Now(),
		Modified:    time.Now(),
	}

	if assoc.SharesRaw != nil {
		result.Shares = int(*assoc.SharesRaw)
		result.RawShares = int(*assoc.SharesRaw)
	}

	return result
}

// extractUserFairShare extracts user fair share from an association
func extractUserFairShare(assoc *types.Association) *types.UserFairShare {
	if assoc == nil {
		return nil
	}

	result := &types.UserFairShare{
		UserName:  assoc.User,
		Account:   derefString(assoc.Account),
		Cluster:   derefString(assoc.Cluster),
		Partition: derefString(assoc.Partition),
		LastDecay: time.Now(),
	}

	if assoc.SharesRaw != nil {
		result.RawShares = int(*assoc.SharesRaw)
	}

	return result
}

// aggregateUserQuotas aggregates quotas across all of a user's associations
func aggregateUserQuotas(associations []types.Association, userName string) *types.UserQuota {
	quota := &types.UserQuota{
		UserName:      userName,
		AccountQuotas: make(map[string]*types.UserAccountQuota),
		TRESLimits:    make(map[string]int),
		IsActive:      true,
	}

	for _, assoc := range associations {
		if assoc.User != userName {
			continue
		}

		accountName := derefString(assoc.Account)

		// Track default account
		if assoc.IsDefault != nil && *assoc.IsDefault {
			quota.DefaultAccount = accountName
		}

		// Create per-account quota
		acctQuota := &types.UserAccountQuota{
			AccountName: accountName,
		}

		if assoc.Max != nil && assoc.Max.Jobs != nil {
			if assoc.Max.Jobs.Active != nil {
				jobs := int(*assoc.Max.Jobs.Active)
				acctQuota.MaxJobs = jobs
				if jobs > quota.MaxJobs {
					quota.MaxJobs = jobs
				}
			}
			if assoc.Max.Jobs.Total != nil {
				submit := int(*assoc.Max.Jobs.Total)
				acctQuota.MaxSubmitJobs = submit
				if submit > quota.MaxSubmitJobs {
					quota.MaxSubmitJobs = submit
				}
			}
			if assoc.Max.Jobs.Per != nil && assoc.Max.Jobs.Per.WallClock != nil {
				wall := int(*assoc.Max.Jobs.Per.WallClock)
				acctQuota.MaxWallTime = wall
				if wall > quota.MaxWallTime {
					quota.MaxWallTime = wall
				}
			}
		}

		if assoc.Priority != nil {
			acctQuota.Priority = int(*assoc.Priority)
		}

		if assoc.QoS != nil {
			acctQuota.QoS = assoc.QoS
		}

		if assoc.Default != nil && assoc.Default.QoS != nil {
			acctQuota.DefaultQoS = *assoc.Default.QoS
		}

		quota.AccountQuotas[accountName] = acctQuota
	}

	return quota
}
