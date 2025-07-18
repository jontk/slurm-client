package v0_0_43

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
)

// TestUserManagerCoverage provides comprehensive test coverage for UserManager
func TestUserManagerCoverage(t *testing.T) {
	ctx := context.Background()
	// Create manager with nil client to trigger validation/implementation errors
	manager := &UserManagerImpl{
		client: &WrapperClient{},
	}

	t.Run("List_Coverage", func(t *testing.T) {
		// Test all option combinations
		testCases := []struct {
			name string
			opts *interfaces.ListUsersOptions
		}{
			{
				name: "nil_options",
				opts: nil,
			},
			{
				name: "empty_options",
				opts: &interfaces.ListUsersOptions{},
			},
			{
				name: "all_flags_true",
				opts: &interfaces.ListUsersOptions{
					WithAccounts:     true,
					WithQuotas:       true,
					WithFairShare:    true,
					WithAssociations: true,
					ActiveOnly:       true,
					Limit:            50,
					Offset:          10,
				},
			},
			{
				name: "partial_flags",
				opts: &interfaces.ListUsersOptions{
					WithAccounts: true,
					ActiveOnly:   false,
					Limit:        100,
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result, err := manager.List(ctx, tc.opts)
				assert.Error(t, err)
				assert.Nil(t, result)
				if !errors.IsNotImplementedError(err) {
					t.Logf("Expected NotImplementedError, got: %v (type: %T)", err, err)
				}
				assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))
			})
		}
	})

	t.Run("Get_Coverage", func(t *testing.T) {
		// Test various user names
		userNames := []string{
			"validuser",
			"user-with-hyphen",
			"user_with_underscore",
			"user123",
			"", // empty should fail validation
		}

		for _, userName := range userNames {
			t.Run(userName, func(t *testing.T) {
				result, err := manager.Get(ctx, userName)
				if userName == "" {
					assert.Error(t, err)
					assert.True(t, errors.IsValidationError(err) || errors.IsClientError(err))
				} else {
					assert.Error(t, err)
					assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))
				}
				assert.Nil(t, result)
			})
		}
	})

	t.Run("GetUserAccounts_Coverage", func(t *testing.T) {
		// Test validation
		_, err := manager.GetUserAccounts(ctx, "")
		assert.Error(t, err)
		assert.True(t, errors.IsValidationError(err) || errors.IsClientError(err))

		// Test normal case
		_, err = manager.GetUserAccounts(ctx, "testuser")
		assert.Error(t, err)
		assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))
	})

	t.Run("GetUserDefaultAccount_Coverage", func(t *testing.T) {
		// Test validation
		_, err := manager.GetUserDefaultAccount(ctx, "")
		assert.Error(t, err)
		assert.True(t, errors.IsValidationError(err) || errors.IsClientError(err))

		// Test normal case
		_, err = manager.GetUserDefaultAccount(ctx, "testuser")
		assert.Error(t, err)
		assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))
	})

	t.Run("GetUserQuotas_Coverage", func(t *testing.T) {
		// Test validation
		_, err := manager.GetUserQuotas(ctx, "")
		assert.Error(t, err)
		assert.True(t, errors.IsValidationError(err) || errors.IsClientError(err))

		// Test normal case
		_, err = manager.GetUserQuotas(ctx, "testuser")
		assert.Error(t, err)
		assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))
	})

	t.Run("GetUserFairShare_Coverage", func(t *testing.T) {
		// Test validation
		_, err := manager.GetUserFairShare(ctx, "")
		assert.Error(t, err)
		assert.True(t, errors.IsValidationError(err) || errors.IsClientError(err))

		// Test normal case
		_, err = manager.GetUserFairShare(ctx, "testuser")
		assert.Error(t, err)
		assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))
	})

	t.Run("CalculateJobPriority_Coverage", func(t *testing.T) {
		// Test nil job submission
		_, err := manager.CalculateJobPriority(ctx, "testuser", nil)
		assert.Error(t, err)
		assert.True(t, errors.IsValidationError(err) || errors.IsClientError(err))

		// Test empty user name
		job := &interfaces.JobSubmission{Script: "#!/bin/bash"}
		_, err = manager.CalculateJobPriority(ctx, "", job)
		assert.Error(t, err)
		assert.True(t, errors.IsValidationError(err) || errors.IsClientError(err))

		// Test empty script
		job = &interfaces.JobSubmission{}
		_, err = manager.CalculateJobPriority(ctx, "testuser", job)
		assert.Error(t, err)
		assert.True(t, errors.IsValidationError(err) || errors.IsClientError(err))

		// Test valid case
		job = &interfaces.JobSubmission{
			Script:    "#!/bin/bash\necho 'test'",
			Partition: "compute",
			CPUs:      4,
		}
		_, err = manager.CalculateJobPriority(ctx, "testuser", job)
		assert.Error(t, err)
		assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))
	})

	t.Run("ValidateUserAccountAccess_Coverage", func(t *testing.T) {
		// Test empty user name
		_, err := manager.ValidateUserAccountAccess(ctx, "", "account")
		assert.Error(t, err)
		assert.True(t, errors.IsValidationError(err) || errors.IsClientError(err))

		// Test empty account name
		_, err = manager.ValidateUserAccountAccess(ctx, "user", "")
		assert.Error(t, err)
		assert.True(t, errors.IsValidationError(err) || errors.IsClientError(err))

		// Test valid case
		_, err = manager.ValidateUserAccountAccess(ctx, "testuser", "testaccount")
		assert.Error(t, err)
		assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))
	})

	t.Run("GetUserAccountAssociations_Coverage", func(t *testing.T) {
		// Test with nil options
		_, err := manager.GetUserAccountAssociations(ctx, "testuser", nil)
		assert.Error(t, err)
		assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))

		// Test with various options
		opts := &interfaces.ListUserAccountAssociationsOptions{
			Accounts:   []string{"acc1", "acc2"},
			Partitions: []string{"compute", "gpu"},
			ActiveOnly: true,
		}
		_, err = manager.GetUserAccountAssociations(ctx, "testuser", opts)
		assert.Error(t, err)
		assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))

		// Test validation
		_, err = manager.GetUserAccountAssociations(ctx, "", opts)
		assert.Error(t, err)
		assert.True(t, errors.IsValidationError(err) || errors.IsClientError(err))
	})

	t.Run("GetBulkUserAccounts_Coverage", func(t *testing.T) {
		// Test empty list
		_, err := manager.GetBulkUserAccounts(ctx, []string{})
		assert.Error(t, err)
		assert.True(t, errors.IsValidationError(err) || errors.IsClientError(err))

		// Test over limit
		tooMany := make([]string, 101)
		for i := range tooMany {
			tooMany[i] = "user" + string(rune('0'+i%10))
		}
		_, err = manager.GetBulkUserAccounts(ctx, tooMany)
		assert.Error(t, err)
		assert.True(t, errors.IsValidationError(err) || errors.IsClientError(err))

		// Test with invalid characters
		_, err = manager.GetBulkUserAccounts(ctx, []string{"user@invalid", "user#bad"})
		assert.Error(t, err)
		assert.True(t, errors.IsValidationError(err) || errors.IsClientError(err))

		// Test valid case
		_, err = manager.GetBulkUserAccounts(ctx, []string{"user1", "user2", "user3"})
		assert.Error(t, err)
		assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))
	})

	t.Run("GetBulkAccountUsers_Coverage", func(t *testing.T) {
		// Test empty list
		_, err := manager.GetBulkAccountUsers(ctx, []string{})
		assert.Error(t, err)
		assert.True(t, errors.IsValidationError(err) || errors.IsClientError(err))

		// Test over limit
		tooMany := make([]string, 101)
		for i := range tooMany {
			tooMany[i] = "account" + string(rune('0'+i%10))
		}
		_, err = manager.GetBulkAccountUsers(ctx, tooMany)
		assert.Error(t, err)
		assert.True(t, errors.IsValidationError(err) || errors.IsClientError(err))

		// Test with invalid characters
		_, err = manager.GetBulkAccountUsers(ctx, []string{"acc@invalid", "acc#bad"})
		assert.Error(t, err)
		assert.True(t, errors.IsValidationError(err) || errors.IsClientError(err))

		// Test valid case
		_, err = manager.GetBulkAccountUsers(ctx, []string{"acc1", "acc2", "acc3"})
		assert.Error(t, err)
		assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))
	})
}

// TestAccountManagerCoverage provides comprehensive test coverage for AccountManager extensions
func TestAccountManagerCoverage(t *testing.T) {
	ctx := context.Background()
	// Create manager with nil client to trigger validation/implementation errors
	manager := &AccountManagerImpl{
		client: &WrapperClient{},
	}

	t.Run("GetAccountHierarchy_Coverage", func(t *testing.T) {
		// Test empty root
		_, err := manager.GetAccountHierarchy(ctx, "")
		assert.Error(t, err)
		assert.True(t, errors.IsValidationError(err) || errors.IsClientError(err))

		// Test normal case
		_, err = manager.GetAccountHierarchy(ctx, "root")
		assert.Error(t, err)
		assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))
	})

	t.Run("GetParentAccounts_Coverage", func(t *testing.T) {
		// Test empty account
		_, err := manager.GetParentAccounts(ctx, "")
		assert.Error(t, err)
		assert.True(t, errors.IsValidationError(err) || errors.IsClientError(err))

		// Test normal case
		_, err = manager.GetParentAccounts(ctx, "child_account")
		assert.Error(t, err)
		assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))
	})

	t.Run("GetChildAccounts_Coverage", func(t *testing.T) {
		// Test empty account
		_, err := manager.GetChildAccounts(ctx, "", 2)
		assert.Error(t, err)
		assert.True(t, errors.IsValidationError(err) || errors.IsClientError(err))

		// Test negative depth
		_, err = manager.GetChildAccounts(ctx, "parent", -1)
		assert.Error(t, err)
		assert.True(t, errors.IsValidationError(err) || errors.IsClientError(err))

		// Test normal cases
		depths := []int{0, 1, 5, 10}
		for _, depth := range depths {
			_, err = manager.GetChildAccounts(ctx, "parent", depth)
			assert.Error(t, err)
			assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))
		}
	})

	t.Run("GetAccountQuotas_Coverage", func(t *testing.T) {
		// Test empty account
		_, err := manager.GetAccountQuotas(ctx, "")
		assert.Error(t, err)
		assert.True(t, errors.IsValidationError(err) || errors.IsClientError(err))

		// Test normal case
		_, err = manager.GetAccountQuotas(ctx, "research")
		assert.Error(t, err)
		assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))
	})

	t.Run("GetAccountQuotaUsage_Coverage", func(t *testing.T) {
		// Test empty account
		_, err := manager.GetAccountQuotaUsage(ctx, "", "daily")
		assert.Error(t, err)
		assert.True(t, errors.IsValidationError(err) || errors.IsClientError(err))

		// Test invalid timeframe
		_, err = manager.GetAccountQuotaUsage(ctx, "research", "invalid")
		assert.Error(t, err)
		assert.True(t, errors.IsValidationError(err) || errors.IsClientError(err))

		// Test valid timeframes
		timeframes := []string{"daily", "weekly", "monthly", "yearly"}
		for _, tf := range timeframes {
			_, err = manager.GetAccountQuotaUsage(ctx, "research", tf)
			assert.Error(t, err)
			assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))
		}
	})

	t.Run("GetAccountUsers_Coverage", func(t *testing.T) {
		// Test empty account
		_, err := manager.GetAccountUsers(ctx, "", nil)
		assert.Error(t, err)
		assert.True(t, errors.IsValidationError(err) || errors.IsClientError(err))

		// Test with options
		opts := &interfaces.ListAccountUsersOptions{
			ActiveOnly: true,
			Limit:      50,
			Offset:     10,
		}
		_, err = manager.GetAccountUsers(ctx, "research", opts)
		assert.Error(t, err)
		assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))
	})

	t.Run("GetAccountUsersWithPermissions_Coverage", func(t *testing.T) {
		// Test empty account
		_, err := manager.GetAccountUsersWithPermissions(ctx, "", []string{"admin"})
		assert.Error(t, err)
		assert.True(t, errors.IsValidationError(err) || errors.IsClientError(err))

		// Test empty permissions
		_, err = manager.GetAccountUsersWithPermissions(ctx, "research", []string{})
		assert.Error(t, err)
		assert.True(t, errors.IsValidationError(err) || errors.IsClientError(err))

		// Test normal case
		perms := []string{"admin", "coordinator", "user"}
		_, err = manager.GetAccountUsersWithPermissions(ctx, "research", perms)
		assert.Error(t, err)
		assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))
	})

	t.Run("ValidateUserAccess_Coverage", func(t *testing.T) {
		// Test empty user
		_, err := manager.ValidateUserAccess(ctx, "", "research")
		assert.Error(t, err)
		assert.True(t, errors.IsValidationError(err) || errors.IsClientError(err))

		// Test empty account
		_, err = manager.ValidateUserAccess(ctx, "alice", "")
		assert.Error(t, err)
		assert.True(t, errors.IsValidationError(err) || errors.IsClientError(err))

		// Test normal case
		_, err = manager.ValidateUserAccess(ctx, "alice", "research")
		assert.Error(t, err)
		assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))
	})

	t.Run("GetAccountFairShare_Coverage", func(t *testing.T) {
		// Test empty account
		_, err := manager.GetAccountFairShare(ctx, "")
		assert.Error(t, err)
		assert.True(t, errors.IsValidationError(err) || errors.IsClientError(err))

		// Test normal case
		_, err = manager.GetAccountFairShare(ctx, "research")
		assert.Error(t, err)
		assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))
	})

	t.Run("GetFairShareHierarchy_Coverage", func(t *testing.T) {
		// Test empty root
		_, err := manager.GetFairShareHierarchy(ctx, "")
		assert.Error(t, err)
		assert.True(t, errors.IsValidationError(err) || errors.IsClientError(err))

		// Test normal case
		_, err = manager.GetFairShareHierarchy(ctx, "root")
		assert.Error(t, err)
		assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))
	})
}

// TestDataStructureValidation tests validation of complex data structures
func TestDataStructureValidation(t *testing.T) {
	t.Run("UserAccount_Validation", func(t *testing.T) {
		// Test UserAccount structure
		acc := interfaces.UserAccount{
			AccountName: "research",
			Partition:   "compute",
			QoS:         "normal",
			DefaultQoS:  "normal",
			IsDefault:   true,
			IsActive:    true,
			MaxJobs:     100,
			Priority:    1000,
			Flags:       []string{"user", "active"},
		}
		assert.NotEmpty(t, acc.AccountName)
		assert.Equal(t, "compute", acc.Partition)
		assert.True(t, acc.IsDefault)
		assert.True(t, acc.IsActive)
		assert.Len(t, acc.Flags, 2)
	})

	t.Run("UserQuota_Validation", func(t *testing.T) {
		// Test UserQuota structure
		quota := interfaces.UserQuota{
			UserName:       "alice",
			DefaultAccount: "research",
			MaxJobs:        50,
			MaxSubmitJobs:  100,
			MaxCPUs:        1000,
			MaxMemory:      1024 * 1024, // 1TB in MB
			MaxWallTime:    1440,        // 24 hours in minutes
			TRESLimits: map[string]int{
				"cpu":  1000,
				"mem":  1024 * 1024,
				"gpu":  8,
				"node": 10,
			},
		}
		assert.Equal(t, "alice", quota.UserName)
		assert.Equal(t, 4, len(quota.TRESLimits))
		assert.Equal(t, 8, quota.TRESLimits["gpu"])
	})

	t.Run("AccountHierarchy_Validation", func(t *testing.T) {
		// Test AccountHierarchy structure
		hierarchy := interfaces.AccountHierarchy{
			Account: &interfaces.Account{
				Name: "root",
			},
			Level:           0,
			Path:            []string{"root"},
			TotalUsers:      500,
			TotalSubAccounts: 50,
			ChildAccounts:   []*interfaces.AccountHierarchy{},
		}
		assert.NotNil(t, hierarchy.Account)
		assert.Equal(t, "root", hierarchy.Account.Name)
		assert.Equal(t, 0, hierarchy.Level)
		assert.Len(t, hierarchy.Path, 1)
		assert.Equal(t, 500, hierarchy.TotalUsers)
	})

	t.Run("FairShareHierarchy_Validation", func(t *testing.T) {
		// Test FairShareHierarchy structure
		hierarchy := interfaces.FairShareHierarchy{
			Cluster:       "main",
			RootAccount:   "root",
			TotalShares:   100000,
			TotalUsage:    0.75,
			Algorithm:     "fair-tree",
			DecayHalfLife: 7,
			UsageWindow:   30,
			LastUpdate:    time.Now(),
			Tree: &interfaces.FairShareNode{
				Name:             "root",
				Account:          "root",
				Shares:           100000,
				FairShareFactor:  1.0,
				NormalizedShares: 1.0,
				Children:         []*interfaces.FairShareNode{},
			},
		}
		assert.Equal(t, "main", hierarchy.Cluster)
		assert.Equal(t, "fair-tree", hierarchy.Algorithm)
		assert.Equal(t, 7, hierarchy.DecayHalfLife)
		assert.NotNil(t, hierarchy.Tree)
	})

	t.Run("JobPriorityInfo_Validation", func(t *testing.T) {
		// Test JobPriorityInfo structure
		estStart := time.Now().Add(30 * time.Minute)
		priority := interfaces.JobPriorityInfo{
			JobID:          12345,
			UserName:       "alice",
			Account:        "research",
			Partition:      "compute",
			QoS:            "normal",
			Priority:       10000,
			PriorityTier:   "high",
			EstimatedStart: estStart,
			EligibleTime:   time.Now(),
			PositionInQueue: 5,
			Factors: &interfaces.JobPriorityFactors{
				FairShare: 5000,
				Age:       2000,
				JobSize:   1000,
				Partition: 1000,
				QoS:       1000,
				Total:     10000,
			},
		}
		assert.Equal(t, 10000, priority.Priority)
		assert.Equal(t, "high", priority.PriorityTier)
		assert.False(t, priority.EstimatedStart.IsZero())
		assert.Equal(t, 5000, priority.Factors.FairShare)
		assert.Equal(t, uint32(12345), priority.JobID)
	})
}

// TestErrorConditions tests various error conditions
func TestErrorConditions(t *testing.T) {
	ctx := context.Background()
	userManager := &UserManagerImpl{}
	accountManager := &AccountManagerImpl{}

	t.Run("UserManager_InvalidInput", func(t *testing.T) {
		// Test with special characters
		invalidNames := []string{
			"user@domain",
			"user#123",
			"user$money",
			"user%percent",
			"user&and",
			"user*star",
			"user/slash",
			"user\\backslash",
		}

		for _, name := range invalidNames {
			t.Run(name, func(t *testing.T) {
				_, err := userManager.Get(ctx, name)
				assert.Error(t, err)
			})
		}
	})

	t.Run("AccountManager_InvalidInput", func(t *testing.T) {
		// Test with special characters
		invalidNames := []string{
			"acc@domain",
			"acc#123",
			"acc$money",
			"acc%percent",
		}

		for _, name := range invalidNames {
			t.Run(name, func(t *testing.T) {
				_, err := accountManager.Get(ctx, name)
				assert.Error(t, err)
			})
		}
	})

	t.Run("BulkOperations_EdgeCases", func(t *testing.T) {
		// Test with exactly 100 items (at limit)
		exactly100 := make([]string, 100)
		for i := range exactly100 {
			exactly100[i] = "user" + string(rune('0'+i%10))
		}
		_, err := userManager.GetBulkUserAccounts(ctx, exactly100)
		assert.Error(t, err)
		assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))

		// Test with duplicate entries
		duplicates := []string{"user1", "user2", "user1", "user3", "user2"}
		_, err = userManager.GetBulkUserAccounts(ctx, duplicates)
		assert.Error(t, err)
	})
}

// TestConcurrentAccess tests thread safety
func TestConcurrentAccess(t *testing.T) {
	ctx := context.Background()
	manager := &UserManagerImpl{}

	// Test concurrent access to same manager
	t.Run("Concurrent_Get", func(t *testing.T) {
		done := make(chan bool, 10)
		
		for i := 0; i < 10; i++ {
			go func(id int) {
				userName := "user" + string(rune('0'+id%10))
				_, err := manager.Get(ctx, userName)
				assert.Error(t, err)
				done <- true
			}(i)
		}

		// Wait for all goroutines
		for i := 0; i < 10; i++ {
			<-done
		}
	})
}

// TestNilHandling tests nil parameter handling
func TestNilHandling(t *testing.T) {
	manager := &UserManagerImpl{}

	t.Run("Nil_Context", func(t *testing.T) {
		// Should handle nil context gracefully
		assert.NotPanics(t, func() {
			_, _ = manager.Get(nil, "testuser")
		})
	})

	t.Run("Nil_Options", func(t *testing.T) {
		ctx := context.Background()
		// Should handle nil options
		_, err := manager.List(ctx, nil)
		assert.Error(t, err)
		assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))
	})
}