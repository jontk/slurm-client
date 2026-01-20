// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/interfaces"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/tests/helpers"
	"github.com/jontk/slurm-client/tests/mocks"
)

// TestUserAccountScenarios tests comprehensive user-account management scenarios
func TestUserAccountScenarios(t *testing.T) {
	testCases := []struct {
		name       string
		apiVersion string
	}{
		{"v0.0.41", "v0.0.41"},
		{"v0.0.42", "v0.0.42"},
		{"v0.0.43", "v0.0.43"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testUserAccountScenariosForVersion(t, tc.apiVersion)
		})
	}
}

func testUserAccountScenariosForVersion(t *testing.T, apiVersion string) {
	// Setup mock server for the specific API version
	mockServer := mocks.NewMockSlurmServerForVersion(apiVersion)
	defer mockServer.Close()

	// Create client
	ctx := helpers.TestContext(t)
	client, err := slurm.NewClientWithVersion(ctx, apiVersion,
		slurm.WithBaseURL(mockServer.URL()),
		slurm.WithAuth(auth.NewNoAuth()),
	)
	require.NoError(t, err)
	defer client.Close()

	// Verify client version
	assert.Equal(t, apiVersion, client.Version())

	t.Run("account_management", func(t *testing.T) {
		testAccountManagement(t, ctx, client, apiVersion)
	})

	t.Run("user_management", func(t *testing.T) {
		testUserManagement(t, ctx, client, apiVersion)
	})

	t.Run("user_account_associations", func(t *testing.T) {
		testUserAccountAssociations(t, ctx, client, apiVersion)
	})

	t.Run("fair_share_operations", func(t *testing.T) {
		testFairShareOperations(t, ctx, client, apiVersion)
	})

	t.Run("hierarchy_navigation", func(t *testing.T) {
		testHierarchyNavigation(t, ctx, client, apiVersion)
	})

	t.Run("quota_monitoring", func(t *testing.T) {
		testQuotaMonitoring(t, ctx, client, apiVersion)
	})
}

func testAccountManagement(t *testing.T, ctx context.Context, client interfaces.SlurmClient, apiVersion string) {
	accountManager := client.Accounts()
	require.NotNil(t, accountManager)

	t.Run("list_accounts", func(t *testing.T) {
		opts := &interfaces.ListAccountsOptions{
			WithUsers:  true,
			WithQuotas: true,
			WithUsage:  true,
		}

		// Test basic account listing
		accountList, err := accountManager.List(ctx, opts)

		// Expected behavior based on version
		switch apiVersion {
		case "v0.0.43":
			// v0.0.43 should support enhanced account features
			// For now, expect NotImplementedError as the actual API integration is pending
			assert.Error(t, err)
			assert.Nil(t, accountList)
		case "v0.0.42":
			// v0.0.42 has limited support
			assert.Error(t, err)
			assert.Nil(t, accountList)
		case "v0.0.41":
			// v0.0.41 has minimal support
			assert.Error(t, err)
			assert.Nil(t, accountList)
		}
	})

	t.Run("get_account_hierarchy", func(t *testing.T) {
		hierarchy, err := accountManager.GetAccountHierarchy(ctx, "root")

		// All versions should handle this consistently
		assert.Error(t, err)
		assert.Nil(t, hierarchy)
	})

	t.Run("get_account_quotas", func(t *testing.T) {
		quotas, err := accountManager.GetAccountQuotas(ctx, "testaccount")

		// All versions should handle this consistently
		assert.Error(t, err)
		assert.Nil(t, quotas)
	})

	t.Run("get_account_fair_share", func(t *testing.T) {
		fairShare, err := accountManager.GetAccountFairShare(ctx, "testaccount")

		// Test version-specific behavior
		switch apiVersion {
		case "v0.0.43":
			assert.Error(t, err)
			assert.Nil(t, fairShare)
		case "v0.0.42":
			assert.Error(t, err)
			assert.Nil(t, fairShare)
		case "v0.0.41":
			assert.Error(t, err)
			assert.Nil(t, fairShare)
		}
	})

	t.Run("get_fair_share_hierarchy", func(t *testing.T) {
		hierarchy, err := accountManager.GetFairShareHierarchy(ctx, "root")

		// Test version-specific behavior
		switch apiVersion {
		case "v0.0.43":
			assert.Error(t, err)
			assert.Nil(t, hierarchy)
		case "v0.0.42":
			assert.Error(t, err)
			assert.Nil(t, hierarchy)
		case "v0.0.41":
			assert.Error(t, err)
			assert.Nil(t, hierarchy)
		}
	})
}

func testUserManagement(t *testing.T, ctx context.Context, client interfaces.SlurmClient, apiVersion string) {
	userManager := client.Users()
	require.NotNil(t, userManager)

	t.Run("list_users", func(t *testing.T) {
		opts := &interfaces.ListUsersOptions{
			WithAccounts:     true,
			WithQuotas:       true,
			WithFairShare:    true,
			WithAssociations: true,
		}

		userList, err := userManager.List(ctx, opts)

		// Expected behavior based on version
		switch apiVersion {
		case "v0.0.43":
			assert.Error(t, err)
			assert.Nil(t, userList)
		case "v0.0.42":
			assert.Error(t, err)
			assert.Nil(t, userList)
		case "v0.0.41":
			assert.Error(t, err)
			assert.Nil(t, userList)
		}
	})

	t.Run("get_user", func(t *testing.T) {
		user, err := userManager.Get(ctx, "testuser")

		// All versions should handle this consistently
		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("get_user_accounts", func(t *testing.T) {
		accounts, err := userManager.GetUserAccounts(ctx, "testuser")

		// All versions should handle this consistently
		assert.Error(t, err)
		assert.Nil(t, accounts)
	})

	t.Run("get_user_quotas", func(t *testing.T) {
		quotas, err := userManager.GetUserQuotas(ctx, "testuser")

		// All versions should handle this consistently
		assert.Error(t, err)
		assert.Nil(t, quotas)
	})

	t.Run("get_user_fair_share", func(t *testing.T) {
		fairShare, err := userManager.GetUserFairShare(ctx, "testuser")

		// All versions should handle this consistently
		assert.Error(t, err)
		assert.Nil(t, fairShare)
	})

	t.Run("calculate_job_priority", func(t *testing.T) {
		jobSubmission := &interfaces.JobSubmission{
			Script:    "#!/bin/bash\necho 'test job'",
			Account:   "testaccount",
			Partition: "compute",
			CPUs:      1,
		}

		priority, err := userManager.CalculateJobPriority(ctx, "testuser", jobSubmission)

		// v0.0.43+ supports job priority calculation
		switch apiVersion {
		case "v0.0.43":
			assert.NoError(t, err)
			assert.NotNil(t, priority)
			assert.Equal(t, "testuser", priority.UserName)
			assert.Greater(t, priority.Priority, 0)
		default:
			// Earlier versions return NotImplementedError
			assert.Error(t, err)
			assert.Nil(t, priority)
		}
	})
}

func testUserAccountAssociations(t *testing.T, ctx context.Context, client interfaces.SlurmClient, apiVersion string) {
	userManager := client.Users()
	accountManager := client.Accounts()
	require.NotNil(t, userManager)
	require.NotNil(t, accountManager)

	t.Run("validate_user_account_access", func(t *testing.T) {
		// Test from UserManager perspective
		validation, err := userManager.ValidateUserAccountAccess(ctx, "testuser", "testaccount")
		assert.Error(t, err)
		assert.Nil(t, validation)

		// Test from AccountManager perspective
		validation2, err := accountManager.ValidateUserAccess(ctx, "testuser", "testaccount")
		assert.Error(t, err)
		assert.Nil(t, validation2)
	})

	t.Run("get_user_account_associations", func(t *testing.T) {
		opts := &interfaces.ListUserAccountAssociationsOptions{
			Accounts:   []string{"testaccount1", "testaccount2"},
			ActiveOnly: true,
		}

		associations, err := userManager.GetUserAccountAssociations(ctx, "testuser", opts)
		assert.Error(t, err)
		assert.Nil(t, associations)
	})

	t.Run("get_account_users", func(t *testing.T) {
		opts := &interfaces.ListAccountUsersOptions{
			ActiveOnly: true,
		}

		users, err := accountManager.GetAccountUsers(ctx, "testaccount", opts)
		assert.Error(t, err)
		assert.Nil(t, users)
	})

	t.Run("get_account_users_with_permissions", func(t *testing.T) {
		permissions := []string{"read", "write", "admin"}
		users, err := accountManager.GetAccountUsersWithPermissions(ctx, "testaccount", permissions)
		assert.Error(t, err)
		assert.Nil(t, users)
	})

	t.Run("bulk_operations", func(t *testing.T) {
		// Test bulk user accounts
		userNames := []string{"user1", "user2", "user3"}
		bulkAccounts, err := userManager.GetBulkUserAccounts(ctx, userNames)
		assert.Error(t, err)
		assert.Nil(t, bulkAccounts)

		// Test bulk account users
		accountNames := []string{"account1", "account2", "account3"}
		bulkUsers, err := userManager.GetBulkAccountUsers(ctx, accountNames)
		assert.Error(t, err)
		assert.Nil(t, bulkUsers)
	})
}

func testFairShareOperations(t *testing.T, ctx context.Context, client interfaces.SlurmClient, apiVersion string) {
	userManager := client.Users()
	accountManager := client.Accounts()
	require.NotNil(t, userManager)
	require.NotNil(t, accountManager)

	t.Run("user_fair_share", func(t *testing.T) {
		fairShare, err := userManager.GetUserFairShare(ctx, "testuser")

		// All versions should handle this consistently with NotImplementedError
		assert.Error(t, err)
		assert.Nil(t, fairShare)
	})

	t.Run("account_fair_share", func(t *testing.T) {
		fairShare, err := accountManager.GetAccountFairShare(ctx, "testaccount")

		// Version-specific behavior
		switch apiVersion {
		case "v0.0.43":
			assert.Error(t, err) // NotImplementedError expected
		case "v0.0.42":
			assert.Error(t, err) // Limited support
		case "v0.0.41":
			assert.Error(t, err) // Not supported
		}
		assert.Nil(t, fairShare)
	})

	t.Run("fair_share_hierarchy", func(t *testing.T) {
		hierarchy, err := accountManager.GetFairShareHierarchy(ctx, "root")

		// Version-specific behavior
		switch apiVersion {
		case "v0.0.43":
			assert.Error(t, err) // NotImplementedError expected
		case "v0.0.42":
			assert.Error(t, err) // Limited support
		case "v0.0.41":
			assert.Error(t, err) // Not supported
		}
		assert.Nil(t, hierarchy)
	})

	t.Run("job_priority_calculation", func(t *testing.T) {
		jobSubmission := &interfaces.JobSubmission{
			Script:    "#!/bin/bash\necho 'priority test'",
			Account:   "testaccount",
			Partition: "compute",
			CPUs:      4,
			Memory:    8192,
			TimeLimit: 60,
		}

		priority, err := userManager.CalculateJobPriority(ctx, "testuser", jobSubmission)

		// v0.0.43+ supports job priority calculation
		switch apiVersion {
		case "v0.0.43":
			assert.NoError(t, err)
			assert.NotNil(t, priority)
			assert.Equal(t, "testuser", priority.UserName)
			assert.Greater(t, priority.Priority, 0)
			assert.NotNil(t, priority.Factors)
		default:
			// Earlier versions return NotImplementedError
			assert.Error(t, err)
			assert.Nil(t, priority)
		}
	})
}

func testHierarchyNavigation(t *testing.T, ctx context.Context, client interfaces.SlurmClient, apiVersion string) {
	accountManager := client.Accounts()
	require.NotNil(t, accountManager)

	t.Run("account_hierarchy_navigation", func(t *testing.T) {
		// Test getting complete hierarchy
		hierarchy, err := accountManager.GetAccountHierarchy(ctx, "root")
		assert.Error(t, err)
		assert.Nil(t, hierarchy)

		// Test getting parent accounts
		parents, err := accountManager.GetParentAccounts(ctx, "child_account")
		assert.Error(t, err)
		assert.Nil(t, parents)

		// Test getting child accounts with different depths
		children, err := accountManager.GetChildAccounts(ctx, "parent_account", 1)
		assert.Error(t, err)
		assert.Nil(t, children)

		// Test unlimited depth
		allChildren, err := accountManager.GetChildAccounts(ctx, "parent_account", 0)
		assert.Error(t, err)
		assert.Nil(t, allChildren)
	})

	t.Run("fair_share_tree_navigation", func(t *testing.T) {
		// Test fair-share hierarchy from different root accounts
		mainHierarchy, err := accountManager.GetFairShareHierarchy(ctx, "main")
		assert.Error(t, err)
		assert.Nil(t, mainHierarchy)

		researchHierarchy, err := accountManager.GetFairShareHierarchy(ctx, "research")
		assert.Error(t, err)
		assert.Nil(t, researchHierarchy)
	})
}

func testQuotaMonitoring(t *testing.T, ctx context.Context, client interfaces.SlurmClient, apiVersion string) {
	userManager := client.Users()
	accountManager := client.Accounts()
	require.NotNil(t, userManager)
	require.NotNil(t, accountManager)

	t.Run("account_quota_monitoring", func(t *testing.T) {
		// Test account quotas
		quotas, err := accountManager.GetAccountQuotas(ctx, "testaccount")
		assert.Error(t, err)
		assert.Nil(t, quotas)

		// Test account quota usage with different timeframes
		dailyUsage, err := accountManager.GetAccountQuotaUsage(ctx, "testaccount", "daily")
		assert.Error(t, err)
		assert.Nil(t, dailyUsage)

		weeklyUsage, err := accountManager.GetAccountQuotaUsage(ctx, "testaccount", "weekly")
		assert.Error(t, err)
		assert.Nil(t, weeklyUsage)

		monthlyUsage, err := accountManager.GetAccountQuotaUsage(ctx, "testaccount", "monthly")
		assert.Error(t, err)
		assert.Nil(t, monthlyUsage)
	})

	t.Run("user_quota_monitoring", func(t *testing.T) {
		// Test user quotas
		quotas, err := userManager.GetUserQuotas(ctx, "testuser")
		assert.Error(t, err)
		assert.Nil(t, quotas)

		// Test user default account
		defaultAccount, err := userManager.GetUserDefaultAccount(ctx, "testuser")
		assert.Error(t, err)
		assert.Nil(t, defaultAccount)
	})

	t.Run("cross_version_quota_behavior", func(t *testing.T) {
		// Test that quota operations behave consistently across versions
		quotas, err := accountManager.GetAccountQuotas(ctx, "testaccount")

		switch apiVersion {
		case "v0.0.43":
			assert.Error(t, err) // NotImplementedError expected
		case "v0.0.42":
			assert.Error(t, err) // Limited support
		case "v0.0.41":
			assert.Error(t, err) // Not supported
		}
		assert.Nil(t, quotas)
	})
}

// TestUserAccountValidation tests input validation scenarios
func TestUserAccountValidation(t *testing.T) {
	// Setup mock server for v0.0.43 (most complete version)
	mockServer := mocks.NewMockSlurmServerForVersion("v0.0.43")
	defer mockServer.Close()

	// Create client
	ctx := helpers.TestContext(t)
	client, err := slurm.NewClientWithVersion(ctx, "v0.0.43",
		slurm.WithBaseURL(mockServer.URL()),
		slurm.WithAuth(auth.NewNoAuth()),
	)
	require.NoError(t, err)
	defer client.Close()

	userManager := client.Users()
	accountManager := client.Accounts()

	t.Run("empty_parameters", func(t *testing.T) {
		// Test empty user names
		_, err := userManager.Get(ctx, "")
		assert.Error(t, err)

		_, err = userManager.GetUserAccounts(ctx, "")
		assert.Error(t, err)

		// Test empty account names
		_, err = accountManager.Get(ctx, "")
		assert.Error(t, err)

		_, err = accountManager.GetAccountQuotas(ctx, "")
		assert.Error(t, err)
	})

	t.Run("invalid_bulk_operations", func(t *testing.T) {
		// Test empty lists
		_, err := userManager.GetBulkUserAccounts(ctx, []string{})
		assert.Error(t, err)

		_, err = userManager.GetBulkAccountUsers(ctx, []string{})
		assert.Error(t, err)

		// Test too many items (over 100 limit)
		tooManyUsers := make([]string, 101)
		for i := range tooManyUsers {
			tooManyUsers[i] = "user" + string(rune('0'+i%10))
		}
		_, err = userManager.GetBulkUserAccounts(ctx, tooManyUsers)
		assert.Error(t, err)
	})

	t.Run("invalid_job_priority_calculation", func(t *testing.T) {
		// Test nil job submission
		_, err := userManager.CalculateJobPriority(ctx, "testuser", nil)
		assert.Error(t, err)

		// Test job submission with no script or command
		emptyJob := &interfaces.JobSubmission{}
		_, err = userManager.CalculateJobPriority(ctx, "testuser", emptyJob)
		assert.Error(t, err)
	})

	t.Run("invalid_hierarchy_parameters", func(t *testing.T) {
		// Test negative depth
		_, err := accountManager.GetChildAccounts(ctx, "testaccount", -1)
		assert.Error(t, err)

		// Test empty root account for hierarchy
		_, err = accountManager.GetAccountHierarchy(ctx, "")
		assert.Error(t, err)

		_, err = accountManager.GetFairShareHierarchy(ctx, "")
		assert.Error(t, err)
	})
}
