// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package builders

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/jontk/slurm-client/internal/common/types"
)

func TestUserBuilder_Basic(t *testing.T) {
	t.Run("simple user", func(t *testing.T) {
		user, err := NewUserBuilder("testuser").
			WithUID(1001).
			WithDefaultAccount("research").
			WithDefaultQoS("normal").
			Build()

		require.NoError(t, err)
		assert.Equal(t, "testuser", user.Name)
		assert.Equal(t, int32(1001), user.UID)
		assert.Equal(t, "research", user.DefaultAccount)
		assert.Equal(t, "normal", user.DefaultQoS)
		assert.Equal(t, types.AdminLevelNone, user.AdminLevel) // Default
	})

	t.Run("empty name fails", func(t *testing.T) {
		_, err := NewUserBuilder("").Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "name is required")
	})

	t.Run("negative values fail", func(t *testing.T) {
		_, err := NewUserBuilder("test").
			WithUID(-1).
			Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must be non-negative")
	})
}

func TestUserBuilder_AdminLevels(t *testing.T) {
	t.Run("admin levels", func(t *testing.T) {
		tests := []struct {
			name  string
			level types.AdminLevel
		}{
			{"none", types.AdminLevelNone},
			{"operator", types.AdminLevelOperator},
			{"administrator", types.AdminLevelAdministrator},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				user, err := NewUserBuilder("testuser").
					WithAdminLevel(tt.level).
					Build()

				require.NoError(t, err)
				assert.Equal(t, tt.level, user.AdminLevel)
			})
		}
	})
}

func TestUserBuilder_AccountsAndQoS(t *testing.T) {
	t.Run("accounts and QoS", func(t *testing.T) {
		user, err := NewUserBuilder("researcher").
			WithDefaultAccount("research").
			WithAccounts("research", "compute", "gpu").
			WithDefaultQoS("normal").
			WithQoSList("debug", "normal", "high").
			WithDefaultWCKey("project1").
			WithWCKeys("project1", "project2").
			Build()

		require.NoError(t, err)
		assert.Equal(t, "research", user.DefaultAccount)
		assert.Equal(t, []string{"research", "compute", "gpu"}, user.Accounts)
		assert.Equal(t, "normal", user.DefaultQoS)
		assert.Equal(t, []string{"debug", "normal", "high"}, user.QoSList)
		assert.Equal(t, "project1", user.DefaultWCKey)
		assert.Equal(t, []string{"project1", "project2"}, user.WCKeys)
	})

	t.Run("default QoS not in list fails", func(t *testing.T) {
		_, err := NewUserBuilder("test").
			WithDefaultQoS("high").
			WithQoSList("debug", "normal").
			Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "default QoS high must be in the allowed QoS list")
	})

	t.Run("default account not in list fails", func(t *testing.T) {
		_, err := NewUserBuilder("test").
			WithDefaultAccount("research").
			WithAccounts("compute", "gpu").
			Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "default account research must be in the allowed accounts list")
	})
}

func TestUserBuilder_Limits(t *testing.T) {
	t.Run("job limits", func(t *testing.T) {
		user, err := NewUserBuilder("power-user").
			WithMaxJobs(1000).
			WithMaxJobsPerAccount(100).
			WithMaxSubmitJobs(2000).
			WithMaxWallTime(10080). // 7 days
			WithMaxCPUTime(20160).  // 14 days
			Build()

		require.NoError(t, err)
		assert.Equal(t, int32(1000), user.MaxJobs)
		assert.Equal(t, int32(100), user.MaxJobsPerAccount)
		assert.Equal(t, int32(2000), user.MaxSubmitJobs)
		assert.Equal(t, int32(10080), user.MaxWallTime)
		assert.Equal(t, int32(20160), user.MaxCPUTime)
	})

	t.Run("resource limits", func(t *testing.T) {
		user, err := NewUserBuilder("hpc-user").
			WithMaxNodes(50).
			WithMaxCPUs(1000).
			WithMaxMemoryGB(2000).
			WithMinPriorityThreshold(100).
			Build()

		require.NoError(t, err)
		assert.Equal(t, int32(50), user.MaxNodes)
		assert.Equal(t, int32(1000), user.MaxCPUs)
		assert.Equal(t, int64(2000*GB), user.MaxMemory)
		assert.Equal(t, int32(100), user.MinPriorityThreshold)
	})
}

func TestUserBuilder_GroupLimits(t *testing.T) {
	t.Run("group limits", func(t *testing.T) {
		user, err := NewUserBuilder("group-user").
			WithGrpJobs(500).
			WithGrpJobsAccrue(100).
			WithGrpNodes(25).
			WithGrpCPUs(500).
			WithGrpMemoryGB(1000).
			WithGrpSubmitJobs(1000).
			WithGrpWallTime(5040).  // 3.5 days
			WithGrpCPUTime(10080).  // 7 days
			Build()

		require.NoError(t, err)
		assert.Equal(t, int32(500), user.GrpJobs)
		assert.Equal(t, int32(100), user.GrpJobsAccrue)
		assert.Equal(t, int32(25), user.GrpNodes)
		assert.Equal(t, int32(500), user.GrpCPUs)
		assert.Equal(t, int64(1000*GB), user.GrpMemory)
		assert.Equal(t, int32(1000), user.GrpSubmitJobs)
		assert.Equal(t, int32(5040), user.GrpWallTime)
		assert.Equal(t, int32(10080), user.GrpCPUTime)
	})
}

func TestUserBuilder_TRES(t *testing.T) {
	t.Run("TRES limits", func(t *testing.T) {
		user, err := NewUserBuilder("tres-user").
			WithGrpTRES("cpu", 500).
			WithGrpTRES("mem", 1000*GB).
			WithGrpTRES("gpu", 5).
			WithGrpTRESMins("cpu", 50000).
			WithGrpTRESRunMins("gpu", 25000).
			WithMaxTRES("cpu", 250).
			WithMaxTRES("mem", 500*GB).
			WithMaxTRESPerNode("gpu", 2).
			WithMinTRES("cpu", 1).
			Build()

		require.NoError(t, err)
		require.NotNil(t, user.GrpTRES)
		assert.Equal(t, int64(500), user.GrpTRES["cpu"])
		assert.Equal(t, int64(1000*GB), user.GrpTRES["mem"])
		assert.Equal(t, int64(5), user.GrpTRES["gpu"])
		
		require.NotNil(t, user.GrpTRESMins)
		assert.Equal(t, int64(50000), user.GrpTRESMins["cpu"])
		
		require.NotNil(t, user.GrpTRESRunMins)
		assert.Equal(t, int64(25000), user.GrpTRESRunMins["gpu"])
		
		require.NotNil(t, user.MaxTRES)
		assert.Equal(t, int64(250), user.MaxTRES["cpu"])
		assert.Equal(t, int64(500*GB), user.MaxTRES["mem"])
		
		require.NotNil(t, user.MaxTRESPerNode)
		assert.Equal(t, int64(2), user.MaxTRESPerNode["gpu"])
		
		require.NotNil(t, user.MinTRES)
		assert.Equal(t, int64(1), user.MinTRES["cpu"])
	})

	t.Run("negative TRES fails", func(t *testing.T) {
		_, err := NewUserBuilder("test").
			WithGrpTRES("cpu", -1).
			Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must be non-negative")
	})
}

func TestUserBuilder_Presets(t *testing.T) {
	t.Run("administrator preset", func(t *testing.T) {
		user, err := NewUserBuilder("admin").
			AsAdministrator().
			Build()

		require.NoError(t, err)
		assert.Equal(t, types.AdminLevelAdministrator, user.AdminLevel)
		assert.Equal(t, int32(10000), user.MaxJobs)
		assert.Equal(t, int32(1000), user.MaxJobsPerAccount)
		assert.Equal(t, int32(20000), user.MaxSubmitJobs)
		assert.Equal(t, int32(43200), user.MaxWallTime) // 30 days
		assert.Equal(t, "high", user.DefaultQoS)
		assert.Contains(t, user.QoSList, "debug")
		assert.Contains(t, user.QoSList, "normal")
		assert.Contains(t, user.QoSList, "high")
		assert.Contains(t, user.QoSList, "urgent")
	})

	t.Run("operator preset", func(t *testing.T) {
		user, err := NewUserBuilder("operator").
			AsOperator().
			Build()

		require.NoError(t, err)
		assert.Equal(t, types.AdminLevelOperator, user.AdminLevel)
		assert.Equal(t, int32(5000), user.MaxJobs)
		assert.Equal(t, int32(500), user.MaxJobsPerAccount)
		assert.Equal(t, int32(10000), user.MaxSubmitJobs)
		assert.Equal(t, int32(20160), user.MaxWallTime) // 14 days
		assert.Equal(t, "normal", user.DefaultQoS)
		assert.Contains(t, user.QoSList, "debug")
		assert.Contains(t, user.QoSList, "normal")
		assert.Contains(t, user.QoSList, "high")
	})

	t.Run("regular user preset", func(t *testing.T) {
		user, err := NewUserBuilder("regular").
			AsRegularUser().
			Build()

		require.NoError(t, err)
		assert.Equal(t, types.AdminLevelNone, user.AdminLevel)
		assert.Equal(t, int32(1000), user.MaxJobs)
		assert.Equal(t, int32(100), user.MaxJobsPerAccount)
		assert.Equal(t, int32(2000), user.MaxSubmitJobs)
		assert.Equal(t, int32(10080), user.MaxWallTime) // 7 days
		assert.Equal(t, int32(50), user.MaxNodes)
		assert.Equal(t, int32(1000), user.MaxCPUs)
		assert.Equal(t, "normal", user.DefaultQoS)
		assert.Contains(t, user.QoSList, "normal")
		assert.Contains(t, user.QoSList, "high")
	})

	t.Run("student user preset", func(t *testing.T) {
		user, err := NewUserBuilder("student").
			AsStudentUser().
			Build()

		require.NoError(t, err)
		assert.Equal(t, types.AdminLevelNone, user.AdminLevel)
		assert.Equal(t, int32(100), user.MaxJobs)
		assert.Equal(t, int32(25), user.MaxJobsPerAccount)
		assert.Equal(t, int32(200), user.MaxSubmitJobs)
		assert.Equal(t, int32(1440), user.MaxWallTime) // 1 day
		assert.Equal(t, int32(4), user.MaxNodes)
		assert.Equal(t, int32(64), user.MaxCPUs)
		assert.Equal(t, int64(256*GB), user.MaxMemory)
		assert.Equal(t, "normal", user.DefaultQoS)
		assert.Equal(t, []string{"normal"}, user.QoSList)
	})

	t.Run("guest user preset", func(t *testing.T) {
		user, err := NewUserBuilder("guest").
			AsGuestUser().
			Build()

		require.NoError(t, err)
		assert.Equal(t, types.AdminLevelNone, user.AdminLevel)
		assert.Equal(t, int32(10), user.MaxJobs)
		assert.Equal(t, int32(5), user.MaxJobsPerAccount)
		assert.Equal(t, int32(20), user.MaxSubmitJobs)
		assert.Equal(t, int32(240), user.MaxWallTime) // 4 hours
		assert.Equal(t, int32(1), user.MaxNodes)
		assert.Equal(t, int32(8), user.MaxCPUs)
		assert.Equal(t, int64(32*GB), user.MaxMemory)
		assert.Equal(t, "debug", user.DefaultQoS)
		assert.Equal(t, []string{"debug"}, user.QoSList)
	})

	t.Run("service user preset", func(t *testing.T) {
		user, err := NewUserBuilder("service").
			AsServiceUser().
			Build()

		require.NoError(t, err)
		assert.Equal(t, types.AdminLevelNone, user.AdminLevel)
		assert.Equal(t, int32(500), user.MaxJobs)
		assert.Equal(t, int32(100), user.MaxJobsPerAccount)
		assert.Equal(t, int32(1000), user.MaxSubmitJobs)
		assert.Equal(t, int32(2880), user.MaxWallTime) // 2 days
		assert.Equal(t, int32(10), user.MaxNodes)
		assert.Equal(t, int32(200), user.MaxCPUs)
		assert.Equal(t, "normal", user.DefaultQoS)
		assert.Equal(t, []string{"normal"}, user.QoSList)
	})

	t.Run("researcher preset", func(t *testing.T) {
		user, err := NewUserBuilder("researcher").
			AsResearcher().
			Build()

		require.NoError(t, err)
		assert.Equal(t, types.AdminLevelNone, user.AdminLevel)
		assert.Equal(t, int32(2000), user.MaxJobs)
		assert.Equal(t, int32(200), user.MaxJobsPerAccount)
		assert.Equal(t, int32(4000), user.MaxSubmitJobs)
		assert.Equal(t, int32(20160), user.MaxWallTime) // 14 days
		assert.Equal(t, int32(100), user.MaxNodes)
		assert.Equal(t, int32(2000), user.MaxCPUs)
		assert.Equal(t, int64(2000*GB), user.MaxMemory)
		assert.Equal(t, "normal", user.DefaultQoS)
		assert.Contains(t, user.QoSList, "normal")
		assert.Contains(t, user.QoSList, "high")
		assert.Contains(t, user.QoSList, "urgent")
		assert.Equal(t, int64(2000), user.MaxTRES["cpu"])
		assert.Equal(t, int64(2000*GB), user.MaxTRES["mem"])
		assert.Equal(t, int64(10), user.MaxTRES["gpu"])
	})

	t.Run("high performance user preset", func(t *testing.T) {
		user, err := NewUserBuilder("hpc").
			AsHighPerformanceUser().
			Build()

		require.NoError(t, err)
		assert.Equal(t, types.AdminLevelNone, user.AdminLevel)
		assert.Equal(t, int32(5000), user.MaxJobs)
		assert.Equal(t, int32(500), user.MaxJobsPerAccount)
		assert.Equal(t, int32(10000), user.MaxSubmitJobs)
		assert.Equal(t, int32(43200), user.MaxWallTime) // 30 days
		assert.Equal(t, int32(500), user.MaxNodes)
		assert.Equal(t, int32(5000), user.MaxCPUs)
		assert.Equal(t, int64(10000*GB), user.MaxMemory)
		assert.Equal(t, int32(2000), user.GrpJobs)
		assert.Equal(t, int32(200), user.GrpNodes)
		assert.Equal(t, int32(2000), user.GrpCPUs)
		assert.Equal(t, int64(5000*GB), user.GrpMemory)
		assert.Equal(t, "high", user.DefaultQoS)
		assert.Contains(t, user.QoSList, "normal")
		assert.Contains(t, user.QoSList, "high")
		assert.Contains(t, user.QoSList, "urgent")
		assert.Equal(t, int64(5000), user.MaxTRES["cpu"])
		assert.Equal(t, int64(10000*GB), user.MaxTRES["mem"])
		assert.Equal(t, int64(500), user.MaxTRES["node"])
		assert.Equal(t, int64(50), user.MaxTRES["gpu"])
	})
}

func TestUserBuilder_BusinessRules(t *testing.T) {
	t.Run("max jobs per account exceeds max jobs", func(t *testing.T) {
		_, err := NewUserBuilder("test").
			WithMaxJobs(10).
			WithMaxJobsPerAccount(20).
			Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "max jobs per account (20) cannot exceed max jobs (10)")
	})

	t.Run("group jobs exceeds max jobs", func(t *testing.T) {
		_, err := NewUserBuilder("test").
			WithMaxJobs(100).
			WithGrpJobs(200).
			Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "group jobs (200) should not exceed max jobs (100)")
	})

	t.Run("TRES CPU consistency", func(t *testing.T) {
		_, err := NewUserBuilder("test").
			WithMaxCPUs(1000).
			WithMaxTRES("cpu", 500).
			Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "MaxCPUs (1000) and MaxTRES[cpu] (500) should be consistent")
	})

	t.Run("TRES memory consistency", func(t *testing.T) {
		_, err := NewUserBuilder("test").
			WithMaxMemoryGB(1000).
			WithMaxTRES("mem", 500*GB).
			Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "MaxMemory")
	})

	t.Run("TRES node consistency", func(t *testing.T) {
		_, err := NewUserBuilder("test").
			WithMaxNodes(100).
			WithMaxTRES("node", 50).
			Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "MaxNodes (100) and MaxTRES[node] (50) should be consistent")
	})

	t.Run("administrator job limits", func(t *testing.T) {
		_, err := NewUserBuilder("test").
			WithAdminLevel(types.AdminLevelAdministrator).
			WithMaxJobs(100). // Too low for admin
			Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "administrators should have at least 1000 max jobs")
	})

	t.Run("consistent TRES settings succeed", func(t *testing.T) {
		user, err := NewUserBuilder("test").
			WithMaxCPUs(1000).
			WithMaxTRES("cpu", 1000).
			WithMaxMemoryGB(2000).
			WithMaxTRES("mem", 2000*GB).
			WithMaxNodes(50).
			WithMaxTRES("node", 50).
			Build()
		require.NoError(t, err)
		assert.Equal(t, int32(1000), user.MaxCPUs)
		assert.Equal(t, int64(2000*GB), user.MaxMemory)
		assert.Equal(t, int32(50), user.MaxNodes)
	})
}

func TestUserBuilder_Clone(t *testing.T) {
	original := NewUserBuilder("original").
		WithUID(1001).
		WithDefaultAccount("research").
		WithAccounts("research").
		WithQoSList("normal").
		WithMaxTRES("cpu", 1000)

	// Clone and modify
	cloned := original.Clone()
	cloned.WithUID(1002).
		WithDefaultAccount("compute").
		WithAccounts("compute").
		WithMaxTRES("mem", 2000*GB)

	// Build both
	originalUser, err := original.Build()
	require.NoError(t, err)
	clonedUser, err := cloned.Build()
	require.NoError(t, err)

	// Verify original is unchanged
	assert.Equal(t, int32(1001), originalUser.UID)
	assert.Equal(t, "research", originalUser.DefaultAccount)
	assert.Equal(t, []string{"research"}, originalUser.Accounts)
	assert.Equal(t, int64(1000), originalUser.MaxTRES["cpu"])
	assert.Equal(t, int64(0), originalUser.MaxTRES["mem"])

	// Verify clone has modifications
	assert.Equal(t, int32(1002), clonedUser.UID)
	assert.Equal(t, "compute", clonedUser.DefaultAccount)
	assert.Equal(t, []string{"research", "compute"}, clonedUser.Accounts)
	assert.Equal(t, int64(1000), clonedUser.MaxTRES["cpu"])
	assert.Equal(t, int64(2000*GB), clonedUser.MaxTRES["mem"])

	// Verify shared attributes are copied
	assert.Equal(t, originalUser.Name, clonedUser.Name)
	assert.Equal(t, originalUser.QoSList, clonedUser.QoSList)
}

func TestUserBuilder_BuildForUpdate(t *testing.T) {
	t.Run("only modified fields", func(t *testing.T) {
		update, err := NewUserBuilder("test").
			WithDefaultAccount("new-account").
			WithMaxJobs(500).
			WithAdminLevel(types.AdminLevelOperator).
			BuildForUpdate()

		require.NoError(t, err)
		require.NotNil(t, update.DefaultAccount)
		assert.Equal(t, "new-account", *update.DefaultAccount)
		require.NotNil(t, update.MaxJobs)
		assert.Equal(t, int32(500), *update.MaxJobs)
		require.NotNil(t, update.AdminLevel)
		assert.Equal(t, types.AdminLevelOperator, *update.AdminLevel)

		// Default values should not be included
		assert.Nil(t, update.MinPriorityThreshold) // Still default 0
	})

	t.Run("with collections and maps", func(t *testing.T) {
		update, err := NewUserBuilder("test").
			WithAccounts("account1", "account2").
			WithQoSList("normal", "high").
			WithWCKeys("key1", "key2").
			WithMaxTRES("cpu", 1000).
			WithMaxTRES("mem", 2000*GB).
			BuildForUpdate()

		require.NoError(t, err)
		assert.Equal(t, []string{"account1", "account2"}, update.Accounts)
		assert.Equal(t, []string{"normal", "high"}, update.QoSList)
		assert.Equal(t, []string{"key1", "key2"}, update.WCKeys)
		require.NotNil(t, update.MaxTRES)
		assert.Equal(t, int64(1000), update.MaxTRES["cpu"])
		assert.Equal(t, int64(2000*GB), update.MaxTRES["mem"])
	})
}

func TestUserBuilder_Validation(t *testing.T) {
	tests := []struct {
		name    string
		builder func() *UserBuilder
		errMsg  string
	}{
		{
			name: "negative UID",
			builder: func() *UserBuilder {
				return NewUserBuilder("test").WithUID(-1)
			},
			errMsg: "must be non-negative",
		},
		{
			name: "negative max jobs",
			builder: func() *UserBuilder {
				return NewUserBuilder("test").WithMaxJobs(-1)
			},
			errMsg: "must be non-negative",
		},
		{
			name: "negative max memory",
			builder: func() *UserBuilder {
				return NewUserBuilder("test").WithMaxMemory(-1)
			},
			errMsg: "must be non-negative",
		},
		{
			name: "negative group TRES",
			builder: func() *UserBuilder {
				return NewUserBuilder("test").WithGrpTRES("cpu", -1)
			},
			errMsg: "must be non-negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.builder().Build()
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}

func TestUserBuilder_ComplexScenario(t *testing.T) {
	// Build a complex user with all features
	user, err := NewUserBuilder("complex-user").
		WithUID(2001).
		WithDefaultAccount("research").
		WithAccounts("research", "compute", "gpu", "highmem").
		WithDefaultWCKey("project-alpha").
		WithWCKeys("project-alpha", "project-beta", "shared").
		WithAdminLevel(types.AdminLevelOperator).
		WithDefaultQoS("normal").
		WithQoSList("debug", "normal", "high", "urgent").
		WithMaxJobs(3000).
		WithMaxJobsPerAccount(300).
		WithMaxSubmitJobs(6000).
		WithMaxWallTime(30240). // 21 days
		WithMaxCPUTime(60480).  // 42 days
		WithMaxNodes(150).
		WithMaxCPUs(3000).
		WithMaxMemoryGB(6000).
		WithMinPriorityThreshold(75).
		WithGrpJobs(1500).
		WithGrpJobsAccrue(300).
		WithGrpNodes(75).
		WithGrpCPUs(1500).
		WithGrpMemoryGB(3000).
		WithGrpSubmitJobs(3000).
		WithGrpWallTime(15120). // 10.5 days
		WithGrpCPUTime(30240).  // 21 days
		WithGrpTRES("cpu", 1500).
		WithGrpTRES("mem", 3000*GB).
		WithGrpTRES("gpu", 15).
		WithGrpTRESMins("cpu", 150000).
		WithGrpTRESMins("gpu", 75000).
		WithGrpTRESRunMins("cpu", 100000).
		WithGrpTRESRunMins("gpu", 50000).
		WithMaxTRES("cpu", 3000).
		WithMaxTRES("mem", 6000*GB).
		WithMaxTRES("node", 150).
		WithMaxTRES("gpu", 30).
		WithMaxTRESPerNode("gpu", 4).
		WithMaxTRESPerNode("mem", 256*GB).
		WithMinTRES("cpu", 1).
		WithMinTRES("mem", 1*GB).
		Build()

	require.NoError(t, err)

	// Verify all fields
	assert.Equal(t, "complex-user", user.Name)
	assert.Equal(t, int32(2001), user.UID)
	assert.Equal(t, "research", user.DefaultAccount)
	assert.Equal(t, []string{"research", "compute", "gpu", "highmem"}, user.Accounts)
	assert.Equal(t, "project-alpha", user.DefaultWCKey)
	assert.Equal(t, []string{"project-alpha", "project-beta", "shared"}, user.WCKeys)
	assert.Equal(t, types.AdminLevelOperator, user.AdminLevel)
	assert.Equal(t, "normal", user.DefaultQoS)
	assert.Equal(t, []string{"debug", "normal", "high", "urgent"}, user.QoSList)
	assert.Equal(t, int32(3000), user.MaxJobs)
	assert.Equal(t, int32(300), user.MaxJobsPerAccount)
	assert.Equal(t, int32(6000), user.MaxSubmitJobs)
	assert.Equal(t, int32(30240), user.MaxWallTime)
	assert.Equal(t, int32(60480), user.MaxCPUTime)
	assert.Equal(t, int32(150), user.MaxNodes)
	assert.Equal(t, int32(3000), user.MaxCPUs)
	assert.Equal(t, int64(6000*GB), user.MaxMemory)
	assert.Equal(t, int32(75), user.MinPriorityThreshold)
	assert.Equal(t, int32(1500), user.GrpJobs)
	assert.Equal(t, int32(300), user.GrpJobsAccrue)
	assert.Equal(t, int32(75), user.GrpNodes)
	assert.Equal(t, int32(1500), user.GrpCPUs)
	assert.Equal(t, int64(3000*GB), user.GrpMemory)
	assert.Equal(t, int32(3000), user.GrpSubmitJobs)
	assert.Equal(t, int32(15120), user.GrpWallTime)
	assert.Equal(t, int32(30240), user.GrpCPUTime)

	// Verify TRES maps
	require.NotNil(t, user.GrpTRES)
	assert.Equal(t, int64(1500), user.GrpTRES["cpu"])
	assert.Equal(t, int64(3000*GB), user.GrpTRES["mem"])
	assert.Equal(t, int64(15), user.GrpTRES["gpu"])

	require.NotNil(t, user.GrpTRESMins)
	assert.Equal(t, int64(150000), user.GrpTRESMins["cpu"])
	assert.Equal(t, int64(75000), user.GrpTRESMins["gpu"])

	require.NotNil(t, user.GrpTRESRunMins)
	assert.Equal(t, int64(100000), user.GrpTRESRunMins["cpu"])
	assert.Equal(t, int64(50000), user.GrpTRESRunMins["gpu"])

	require.NotNil(t, user.MaxTRES)
	assert.Equal(t, int64(3000), user.MaxTRES["cpu"])
	assert.Equal(t, int64(6000*GB), user.MaxTRES["mem"])
	assert.Equal(t, int64(150), user.MaxTRES["node"])
	assert.Equal(t, int64(30), user.MaxTRES["gpu"])

	require.NotNil(t, user.MaxTRESPerNode)
	assert.Equal(t, int64(4), user.MaxTRESPerNode["gpu"])
	assert.Equal(t, int64(256*GB), user.MaxTRESPerNode["mem"])

	require.NotNil(t, user.MinTRES)
	assert.Equal(t, int64(1), user.MinTRES["cpu"])
	assert.Equal(t, int64(1*GB), user.MinTRES["mem"])
}
