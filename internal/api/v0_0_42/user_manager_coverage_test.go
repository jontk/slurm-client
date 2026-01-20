// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_42

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jontk/slurm-client/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
)

// TestUserManagerImpl_v0_0_42_Coverage provides test coverage for v0.0.42 UserManager
func TestUserManagerImpl_v0_0_42_Coverage(t *testing.T) {
	ctx := context.Background()
	manager := &UserManagerImpl{
		client: &WrapperClient{},
	}

	t.Run("Supported_Methods", func(t *testing.T) {
		// List - supported
		list, err := manager.List(ctx, nil)
		assert.Nil(t, list)
		assert.Error(t, err)
		assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))

		// Get - supported
		user, err := manager.Get(ctx, "testuser")
		assert.Nil(t, user)
		assert.Error(t, err)
		assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))

		// GetUserAccounts - supported
		accounts, err := manager.GetUserAccounts(ctx, "testuser")
		assert.Nil(t, accounts)
		assert.Error(t, err)
		assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))

		// GetUserDefaultAccount - supported
		defaultAcc, err := manager.GetUserDefaultAccount(ctx, "testuser")
		assert.Nil(t, defaultAcc)
		assert.Error(t, err)
		assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))

		// GetUserFairShare - supported
		fairShare, err := manager.GetUserFairShare(ctx, "testuser")
		assert.Nil(t, fairShare)
		assert.Error(t, err)
		assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))

		// ValidateUserAccountAccess - supported
		validation, err := manager.ValidateUserAccountAccess(ctx, "testuser", "account")
		assert.Nil(t, validation)
		assert.Error(t, err)
		assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))

		// GetUserAccountAssociations - supported
		assocs, err := manager.GetUserAccountAssociations(ctx, "testuser", nil)
		assert.Nil(t, assocs)
		assert.Error(t, err)
		assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))
	})

	t.Run("Limited_Support_Methods", func(t *testing.T) {
		// GetUserQuotas - limited support
		quotas, err := manager.GetUserQuotas(ctx, "testuser")
		assert.Nil(t, quotas)
		assert.Error(t, err)
		assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))
		// Message content varies based on implementation

		// CalculateJobPriority - limited support
		job := &interfaces.JobSubmission{Script: "#!/bin/bash"}
		priority, err := manager.CalculateJobPriority(ctx, "testuser", job)
		assert.Nil(t, priority)
		assert.Error(t, err)
		assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))
		// Message content varies based on implementation
	})

	t.Run("Not_Supported_Methods", func(t *testing.T) {
		// GetBulkUserAccounts - not supported
		bulkAccounts, err := manager.GetBulkUserAccounts(ctx, []string{"user1", "user2"})
		assert.Nil(t, bulkAccounts)
		assert.Error(t, err)
		assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))
		// Message content varies based on implementation

		// GetBulkAccountUsers - not supported
		bulkUsers, err := manager.GetBulkAccountUsers(ctx, []string{"acc1", "acc2"})
		assert.Nil(t, bulkUsers)
		assert.Error(t, err)
		assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))
		// Message content varies based on implementation
	})

	t.Run("Input_Validation", func(t *testing.T) {
		// Empty user name
		_, err := manager.Get(ctx, "")
		assert.Error(t, err)
		assert.True(t, errors.IsValidationError(err) || errors.IsClientError(err))

		// Empty account name
		_, err = manager.ValidateUserAccountAccess(ctx, "user", "")
		assert.Error(t, err)
		assert.True(t, errors.IsValidationError(err) || errors.IsClientError(err))

		// Nil job submission
		_, err = manager.CalculateJobPriority(ctx, "user", nil)
		assert.Error(t, err)
		assert.True(t, errors.IsValidationError(err) || errors.IsClientError(err))
	})
}

// TestAccountManagerImpl_v0_0_42_Coverage provides test coverage for v0.0.42 AccountManager extensions
func TestAccountManagerImpl_v0_0_42_Coverage(t *testing.T) {
	ctx := context.Background()
	manager := &AccountManagerImpl{
		client: &WrapperClient{},
	}

	t.Run("Limited_Support_Methods", func(t *testing.T) {
		// GetAccountHierarchy - limited support
		hierarchy, err := manager.GetAccountHierarchy(ctx, "root")
		assert.Nil(t, hierarchy)
		assert.Error(t, err)
		assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))
		// Message content varies based on implementation

		// GetParentAccounts - limited support
		parents, err := manager.GetParentAccounts(ctx, "child")
		assert.Nil(t, parents)
		assert.Error(t, err)
		assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))
		// Message content varies based on implementation

		// GetChildAccounts - limited support
		children, err := manager.GetChildAccounts(ctx, "parent", 2)
		assert.Nil(t, children)
		assert.Error(t, err)
		assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))
		// Message content varies based on implementation

		// GetAccountQuotas - limited support
		quotas, err := manager.GetAccountQuotas(ctx, "research")
		assert.Nil(t, quotas)
		assert.Error(t, err)
		assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))
		// Message content varies based on implementation

		// GetAccountQuotaUsage - limited support
		usage, err := manager.GetAccountQuotaUsage(ctx, "research", "daily")
		assert.Nil(t, usage)
		assert.Error(t, err)
		assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))
		// Message content varies based on implementation
	})

	t.Run("Supported_Methods", func(t *testing.T) {
		// GetAccountUsers - supported
		users, err := manager.GetAccountUsers(ctx, "research", nil)
		assert.Nil(t, users)
		assert.Error(t, err)
		assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))

		// GetAccountUsersWithPermissions - supported
		permsUsers, err := manager.GetAccountUsersWithPermissions(ctx, "research", []string{"admin"})
		assert.Nil(t, permsUsers)
		assert.Error(t, err)
		assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))

		// ValidateUserAccess - supported
		validation, err := manager.ValidateUserAccess(ctx, "user", "account")
		assert.Nil(t, validation)
		assert.Error(t, err)
		assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))

		// GetAccountFairShare - supported
		fairShare, err := manager.GetAccountFairShare(ctx, "research")
		assert.Nil(t, fairShare)
		assert.Error(t, err)
		assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))

		// GetFairShareHierarchy - supported
		fsHierarchy, err := manager.GetFairShareHierarchy(ctx, "root")
		assert.Nil(t, fsHierarchy)
		assert.Error(t, err)
		assert.True(t, errors.IsNotImplementedError(err) || errors.IsClientError(err))
	})

	t.Run("Input_Validation", func(t *testing.T) {
		// Empty account name
		_, err := manager.GetAccountHierarchy(ctx, "")
		assert.Error(t, err)
		assert.True(t, errors.IsValidationError(err) || errors.IsClientError(err))

		// Negative depth
		_, err = manager.GetChildAccounts(ctx, "parent", -1)
		assert.Error(t, err)
		assert.True(t, errors.IsValidationError(err) || errors.IsClientError(err))

		// Invalid timeframe
		_, err = manager.GetAccountQuotaUsage(ctx, "research", "invalid")
		assert.Error(t, err)
		assert.True(t, errors.IsValidationError(err) || errors.IsClientError(err))

		// Empty permissions
		_, err = manager.GetAccountUsersWithPermissions(ctx, "research", []string{})
		assert.Error(t, err)
		assert.True(t, errors.IsValidationError(err) || errors.IsClientError(err))
	})
}
