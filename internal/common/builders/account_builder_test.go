// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package builders

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountBuilder_Basic(t *testing.T) {
	t.Run("simple account", func(t *testing.T) {
		account, err := NewAccountBuilder("research-group").
			WithDescription("Research group account").
			WithOrganization("University").
			WithDefaultQoS("normal").
			Build()

		require.NoError(t, err)
		assert.Equal(t, "research-group", account.Name)
		assert.Equal(t, "Research group account", account.Description)
		assert.Equal(t, "University", account.Organization)
		assert.Equal(t, "normal", account.DefaultQoS)
		assert.Equal(t, int32(1), account.FairShare) // Default
		assert.Equal(t, int32(1), account.SharesRaw) // Default
	})

	t.Run("empty name fails", func(t *testing.T) {
		_, err := NewAccountBuilder("").Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "name is required")
	})

	t.Run("negative values fail", func(t *testing.T) {
		_, err := NewAccountBuilder("test").
			WithMaxJobs(-1).
			Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must be non-negative")
	})
}

func TestAccountBuilder_Hierarchy(t *testing.T) {
	t.Run("parent child relationship", func(t *testing.T) {
		account, err := NewAccountBuilder("child-account").
			WithParentAccount("parent-account").
			WithCoordinators("coordinator1", "coordinator2").
			Build()

		require.NoError(t, err)
		assert.Equal(t, "child-account", account.Name)
		assert.Equal(t, "parent-account", account.ParentName)
		assert.Equal(t, []string{"coordinator1", "coordinator2"}, account.Coordinators)
	})

	t.Run("self parent fails", func(t *testing.T) {
		_, err := NewAccountBuilder("self").
			WithParentAccount("self").
			Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be its own parent")
	})
}

func TestAccountBuilder_QoSAndPartitions(t *testing.T) {
	t.Run("QoS management", func(t *testing.T) {
		account, err := NewAccountBuilder("qos-test").
			WithDefaultQoS("normal").
			WithQoSList("debug", "normal", "high").
			Build()

		require.NoError(t, err)
		assert.Equal(t, "normal", account.DefaultQoS)
		assert.Equal(t, []string{"debug", "normal", "high"}, account.QoSList)
	})

	t.Run("partition management", func(t *testing.T) {
		account, err := NewAccountBuilder("partition-test").
			WithDefaultPartition("batch").
			WithAllowedPartitions("debug", "batch", "gpu").
			Build()

		require.NoError(t, err)
		assert.Equal(t, "batch", account.DefaultPartition)
		assert.Equal(t, []string{"debug", "batch", "gpu"}, account.AllowedPartitions)
	})

	t.Run("default QoS not in list fails", func(t *testing.T) {
		_, err := NewAccountBuilder("test").
			WithDefaultQoS("high").
			WithQoSList("debug", "normal").
			Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "default QoS high must be in the allowed QoS list")
	})

	t.Run("default partition not in list fails", func(t *testing.T) {
		_, err := NewAccountBuilder("test").
			WithDefaultPartition("gpu").
			WithAllowedPartitions("debug", "batch").
			Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "default partition gpu must be in the allowed partitions list")
	})
}

func TestAccountBuilder_Limits(t *testing.T) {
	t.Run("job limits", func(t *testing.T) {
		account, err := NewAccountBuilder("limits-test").
			WithMaxJobs(1000).
			WithMaxJobsPerUser(50).
			WithMaxSubmitJobs(2000).
			WithMaxWallTime(10080). // 7 days
			WithMaxCPUTime(20160).  // 14 days
			Build()

		require.NoError(t, err)
		assert.Equal(t, int32(1000), account.MaxJobs)
		assert.Equal(t, int32(50), account.MaxJobsPerUser)
		assert.Equal(t, int32(2000), account.MaxSubmitJobs)
		assert.Equal(t, int32(10080), account.MaxWallTime)
		assert.Equal(t, int32(20160), account.MaxCPUTime)
	})

	t.Run("resource limits", func(t *testing.T) {
		account, err := NewAccountBuilder("resources-test").
			WithMaxNodes(100).
			WithMaxCPUs(2000).
			WithMaxMemoryGB(4000).
			WithMinPriorityThreshold(100).
			Build()

		require.NoError(t, err)
		assert.Equal(t, int32(100), account.MaxNodes)
		assert.Equal(t, int32(2000), account.MaxCPUs)
		assert.Equal(t, int64(4000*GB), account.MaxMemory)
		assert.Equal(t, int32(100), account.MinPriorityThreshold)
	})
}

func TestAccountBuilder_GroupLimits(t *testing.T) {
	t.Run("group limits", func(t *testing.T) {
		account, err := NewAccountBuilder("group-test").
			WithGrpJobs(500).
			WithGrpJobsAccrue(100).
			WithGrpNodes(50).
			WithGrpCPUs(1000).
			WithGrpMemoryGB(2000).
			WithGrpSubmitJobs(1000).
			WithGrpWallTime(5040).  // 3.5 days
			WithGrpCPUTime(10080).  // 7 days
			Build()

		require.NoError(t, err)
		assert.Equal(t, int32(500), account.GrpJobs)
		assert.Equal(t, int32(100), account.GrpJobsAccrue)
		assert.Equal(t, int32(50), account.GrpNodes)
		assert.Equal(t, int32(1000), account.GrpCPUs)
		assert.Equal(t, int64(2000*GB), account.GrpMemory)
		assert.Equal(t, int32(1000), account.GrpSubmitJobs)
		assert.Equal(t, int32(5040), account.GrpWallTime)
		assert.Equal(t, int32(10080), account.GrpCPUTime)
	})
}

func TestAccountBuilder_TRES(t *testing.T) {
	t.Run("TRES limits", func(t *testing.T) {
		account, err := NewAccountBuilder("tres-test").
			WithGrpTRES("cpu", 1000).
			WithGrpTRES("mem", 2000*GB).
			WithGrpTRES("gpu", 10).
			WithGrpTRESMins("cpu", 100000).
			WithGrpTRESRunMins("gpu", 50000).
			WithMaxTRES("cpu", 500).
			WithMaxTRES("mem", 1000*GB).
			WithMaxTRESPerNode("gpu", 4).
			WithMinTRES("cpu", 1).
			Build()

		require.NoError(t, err)
		require.NotNil(t, account.GrpTRES)
		assert.Equal(t, int64(1000), account.GrpTRES["cpu"])
		assert.Equal(t, int64(2000*GB), account.GrpTRES["mem"])
		assert.Equal(t, int64(10), account.GrpTRES["gpu"])
		
		require.NotNil(t, account.GrpTRESMins)
		assert.Equal(t, int64(100000), account.GrpTRESMins["cpu"])
		
		require.NotNil(t, account.GrpTRESRunMins)
		assert.Equal(t, int64(50000), account.GrpTRESRunMins["gpu"])
		
		require.NotNil(t, account.MaxTRES)
		assert.Equal(t, int64(500), account.MaxTRES["cpu"])
		assert.Equal(t, int64(1000*GB), account.MaxTRES["mem"])
		
		require.NotNil(t, account.MaxTRESPerNode)
		assert.Equal(t, int64(4), account.MaxTRESPerNode["gpu"])
		
		require.NotNil(t, account.MinTRES)
		assert.Equal(t, int64(1), account.MinTRES["cpu"])
	})

	t.Run("negative TRES fails", func(t *testing.T) {
		_, err := NewAccountBuilder("test").
			WithGrpTRES("cpu", -1).
			Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must be non-negative")
	})
}

func TestAccountBuilder_Shares(t *testing.T) {
	t.Run("share settings", func(t *testing.T) {
		account, err := NewAccountBuilder("shares-test").
			WithFairShare(100).
			WithSharesRaw(200).
			WithPriority(500).
			Build()

		require.NoError(t, err)
		assert.Equal(t, int32(100), account.FairShare)
		assert.Equal(t, int32(200), account.SharesRaw)
		assert.Equal(t, int32(500), account.Priority)
	})
}

func TestAccountBuilder_Presets(t *testing.T) {
	t.Run("research account preset", func(t *testing.T) {
		account, err := NewAccountBuilder("research").
			AsResearchAccount().
			Build()

		require.NoError(t, err)
		assert.Contains(t, account.Description, "Research account")
		assert.Equal(t, int32(100), account.FairShare)
		assert.Equal(t, int32(1000), account.MaxJobs)
		assert.Equal(t, int32(50), account.MaxJobsPerUser)
		assert.Equal(t, int32(10080), account.MaxWallTime) // 7 days
		assert.Equal(t, int32(500), account.GrpJobs)
		assert.Equal(t, "normal", account.DefaultQoS)
		assert.Contains(t, account.QoSList, "normal")
		assert.Contains(t, account.QoSList, "high")
	})

	t.Run("compute account preset", func(t *testing.T) {
		account, err := NewAccountBuilder("compute").
			AsComputeAccount().
			Build()

		require.NoError(t, err)
		assert.Contains(t, account.Description, "General compute")
		assert.Equal(t, int32(50), account.FairShare)
		assert.Equal(t, int32(500), account.MaxJobs)
		assert.Equal(t, int32(25), account.MaxJobsPerUser)
		assert.Equal(t, int32(2880), account.MaxWallTime) // 2 days
		assert.Equal(t, int32(200), account.GrpJobs)
		assert.Equal(t, "normal", account.DefaultQoS)
		assert.Equal(t, []string{"normal"}, account.QoSList)
	})

	t.Run("student account preset", func(t *testing.T) {
		account, err := NewAccountBuilder("students").
			AsStudentAccount().
			Build()

		require.NoError(t, err)
		assert.Contains(t, account.Description, "Student account")
		assert.Equal(t, int32(10), account.FairShare)
		assert.Equal(t, int32(50), account.MaxJobs)
		assert.Equal(t, int32(10), account.MaxJobsPerUser)
		assert.Equal(t, int32(480), account.MaxWallTime) // 8 hours
		assert.Equal(t, int32(4), account.MaxNodes)
		assert.Equal(t, int32(64), account.MaxCPUs)
		assert.Equal(t, int32(25), account.GrpJobs)
		assert.Equal(t, "normal", account.DefaultQoS)
	})

	t.Run("guest account preset", func(t *testing.T) {
		account, err := NewAccountBuilder("guest").
			AsGuestAccount().
			Build()

		require.NoError(t, err)
		assert.Contains(t, account.Description, "Guest account")
		assert.Equal(t, int32(1), account.FairShare)
		assert.Equal(t, int32(10), account.MaxJobs)
		assert.Equal(t, int32(5), account.MaxJobsPerUser)
		assert.Equal(t, int32(60), account.MaxWallTime) // 1 hour
		assert.Equal(t, int32(1), account.MaxNodes)
		assert.Equal(t, int32(8), account.MaxCPUs)
		assert.Equal(t, int32(5), account.GrpJobs)
		assert.Equal(t, "debug", account.DefaultQoS)
	})

	t.Run("high performance account preset", func(t *testing.T) {
		account, err := NewAccountBuilder("hpc").
			AsHighPerformanceAccount().
			Build()

		require.NoError(t, err)
		assert.Contains(t, account.Description, "High-performance")
		assert.Equal(t, int32(200), account.FairShare)
		assert.Equal(t, int32(2000), account.MaxJobs)
		assert.Equal(t, int32(100), account.MaxJobsPerUser)
		assert.Equal(t, int32(20160), account.MaxWallTime) // 14 days
		assert.Equal(t, int32(100), account.MaxNodes)
		assert.Equal(t, int32(2000), account.MaxCPUs)
		assert.Equal(t, int32(1000), account.GrpJobs)
		assert.Equal(t, int32(50), account.GrpNodes)
		assert.Equal(t, int32(1000), account.GrpCPUs)
		assert.Equal(t, "high", account.DefaultQoS)
		assert.Contains(t, account.QoSList, "normal")
		assert.Contains(t, account.QoSList, "high")
		assert.Contains(t, account.QoSList, "urgent")
		assert.Equal(t, int64(2000), account.MaxTRES["cpu"])
		assert.Equal(t, int64(2000*GB), account.MaxTRES["mem"])
		assert.Equal(t, int64(100), account.MaxTRES["node"])
	})

	t.Run("service account preset", func(t *testing.T) {
		account, err := NewAccountBuilder("service").
			AsServiceAccount().
			Build()

		require.NoError(t, err)
		assert.Contains(t, account.Description, "Service account")
		assert.Equal(t, int32(25), account.FairShare)
		assert.Equal(t, int32(100), account.MaxJobs)
		assert.Equal(t, int32(50), account.MaxJobsPerUser)
		assert.Equal(t, int32(1440), account.MaxWallTime) // 1 day
		assert.Equal(t, int32(50), account.GrpJobs)
		assert.Equal(t, "normal", account.DefaultQoS)
	})
}

func TestAccountBuilder_BusinessRules(t *testing.T) {
	t.Run("max jobs per user exceeds max jobs", func(t *testing.T) {
		_, err := NewAccountBuilder("test").
			WithMaxJobs(10).
			WithMaxJobsPerUser(20).
			Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "max jobs per user (20) cannot exceed max jobs (10)")
	})

	t.Run("group jobs exceeds max jobs", func(t *testing.T) {
		_, err := NewAccountBuilder("test").
			WithMaxJobs(100).
			WithGrpJobs(200).
			Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "group jobs (200) should not exceed max jobs (100)")
	})

	t.Run("TRES CPU consistency", func(t *testing.T) {
		_, err := NewAccountBuilder("test").
			WithMaxCPUs(1000).
			WithMaxTRES("cpu", 500).
			Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "MaxCPUs (1000) and MaxTRES[cpu] (500) should be consistent")
	})

	t.Run("TRES memory consistency", func(t *testing.T) {
		_, err := NewAccountBuilder("test").
			WithMaxMemoryGB(1000).
			WithMaxTRES("mem", 500*GB).
			Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "MaxMemory")
	})

	t.Run("TRES node consistency", func(t *testing.T) {
		_, err := NewAccountBuilder("test").
			WithMaxNodes(100).
			WithMaxTRES("node", 50).
			Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "MaxNodes (100) and MaxTRES[node] (50) should be consistent")
	})

	t.Run("consistent TRES settings succeed", func(t *testing.T) {
		account, err := NewAccountBuilder("test").
			WithMaxCPUs(1000).
			WithMaxTRES("cpu", 1000).
			WithMaxMemoryGB(2000).
			WithMaxTRES("mem", 2000*GB).
			WithMaxNodes(50).
			WithMaxTRES("node", 50).
			Build()
		require.NoError(t, err)
		assert.Equal(t, int32(1000), account.MaxCPUs)
		assert.Equal(t, int64(2000*GB), account.MaxMemory)
		assert.Equal(t, int32(50), account.MaxNodes)
	})
}

func TestAccountBuilder_Clone(t *testing.T) {
	original := NewAccountBuilder("original").
		WithDescription("Original account").
		WithFairShare(100).
		WithCoordinators("coord1").
		WithQoSList("normal").
		WithMaxTRES("cpu", 1000)

	// Clone and modify
	cloned := original.Clone()
	cloned.WithDescription("Cloned account").
		WithFairShare(200).
		WithCoordinators("coord2").
		WithMaxTRES("mem", 2000*GB)

	// Build both
	originalAccount, err := original.Build()
	require.NoError(t, err)
	clonedAccount, err := cloned.Build()
	require.NoError(t, err)

	// Verify original is unchanged
	assert.Equal(t, "Original account", originalAccount.Description)
	assert.Equal(t, int32(100), originalAccount.FairShare)
	assert.Equal(t, []string{"coord1"}, originalAccount.Coordinators)
	assert.Equal(t, int64(1000), originalAccount.MaxTRES["cpu"])
	assert.Equal(t, int64(0), originalAccount.MaxTRES["mem"])

	// Verify clone has modifications
	assert.Equal(t, "Cloned account", clonedAccount.Description)
	assert.Equal(t, int32(200), clonedAccount.FairShare)
	assert.Equal(t, []string{"coord1", "coord2"}, clonedAccount.Coordinators)
	assert.Equal(t, int64(1000), clonedAccount.MaxTRES["cpu"])
	assert.Equal(t, int64(2000*GB), clonedAccount.MaxTRES["mem"])

	// Verify shared attributes are copied
	assert.Equal(t, originalAccount.Name, clonedAccount.Name)
	assert.Equal(t, originalAccount.QoSList, clonedAccount.QoSList)
}

func TestAccountBuilder_BuildForUpdate(t *testing.T) {
	t.Run("only modified fields", func(t *testing.T) {
		update, err := NewAccountBuilder("test").
			WithDescription("Updated description").
			WithMaxJobs(500).
			WithPriority(100).
			BuildForUpdate()

		require.NoError(t, err)
		require.NotNil(t, update.Description)
		assert.Equal(t, "Updated description", *update.Description)
		require.NotNil(t, update.MaxJobs)
		assert.Equal(t, int32(500), *update.MaxJobs)
		require.NotNil(t, update.Priority)
		assert.Equal(t, int32(100), *update.Priority)

		// Default values should not be included
		assert.Nil(t, update.FairShare)  // Still default 1
		assert.Nil(t, update.SharesRaw)  // Still default 1
	})

	t.Run("with collections and maps", func(t *testing.T) {
		update, err := NewAccountBuilder("test").
			WithCoordinators("coord1", "coord2").
			WithQoSList("normal", "high").
			WithMaxTRES("cpu", 1000).
			WithMaxTRES("mem", 2000*GB).
			BuildForUpdate()

		require.NoError(t, err)
		assert.Equal(t, []string{"coord1", "coord2"}, update.Coordinators)
		assert.Equal(t, []string{"normal", "high"}, update.QoSList)
		require.NotNil(t, update.MaxTRES)
		assert.Equal(t, int64(1000), update.MaxTRES["cpu"])
		assert.Equal(t, int64(2000*GB), update.MaxTRES["mem"])
	})
}

func TestAccountBuilder_Validation(t *testing.T) {
	tests := []struct {
		name    string
		builder func() *AccountBuilder
		errMsg  string
	}{
		{
			name: "negative fair share",
			builder: func() *AccountBuilder {
				return NewAccountBuilder("test").WithFairShare(-1)
			},
			errMsg: "must be non-negative",
		},
		{
			name: "negative max jobs",
			builder: func() *AccountBuilder {
				return NewAccountBuilder("test").WithMaxJobs(-1)
			},
			errMsg: "must be non-negative",
		},
		{
			name: "negative max memory",
			builder: func() *AccountBuilder {
				return NewAccountBuilder("test").WithMaxMemory(-1)
			},
			errMsg: "must be non-negative",
		},
		{
			name: "negative group TRES",
			builder: func() *AccountBuilder {
				return NewAccountBuilder("test").WithGrpTRES("cpu", -1)
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

func TestAccountBuilder_ComplexScenario(t *testing.T) {
	// Build a complex account with all features
	account, err := NewAccountBuilder("complex-account").
		WithDescription("Complex research account with all features").
		WithOrganization("Research University").
		WithCoordinators("pi1", "pi2", "admin1").
		WithDefaultQoS("normal").
		WithQoSList("debug", "normal", "high", "urgent").
		WithParentAccount("root-account").
		WithAllowedPartitions("debug", "batch", "gpu", "highmem").
		WithDefaultPartition("batch").
		WithFairShare(150).
		WithSharesRaw(300).
		WithPriority(100).
		WithMaxJobs(2000).
		WithMaxJobsPerUser(100).
		WithMaxSubmitJobs(4000).
		WithMaxWallTime(20160). // 14 days
		WithMaxCPUTime(40320).  // 28 days
		WithMaxNodes(200).
		WithMaxCPUs(4000).
		WithMaxMemoryGB(8000).
		WithMinPriorityThreshold(50).
		WithGrpJobs(1000).
		WithGrpJobsAccrue(200).
		WithGrpNodes(100).
		WithGrpCPUs(2000).
		WithGrpMemoryGB(4000).
		WithGrpSubmitJobs(2000).
		WithGrpWallTime(10080). // 7 days
		WithGrpCPUTime(20160).  // 14 days
		WithGrpTRES("cpu", 2000).
		WithGrpTRES("mem", 4000*GB).
		WithGrpTRES("gpu", 20).
		WithGrpTRESMins("cpu", 200000).
		WithGrpTRESMins("gpu", 100000).
		WithGrpTRESRunMins("cpu", 150000).
		WithGrpTRESRunMins("gpu", 75000).
		WithMaxTRES("cpu", 4000).
		WithMaxTRES("mem", 8000*GB).
		WithMaxTRES("node", 200).
		WithMaxTRES("gpu", 40).
		WithMaxTRESPerNode("gpu", 8).
		WithMaxTRESPerNode("mem", 512*GB).
		WithMinTRES("cpu", 1).
		WithMinTRES("mem", 1*GB).
		Build()

	require.NoError(t, err)

	// Verify all fields
	assert.Equal(t, "complex-account", account.Name)
	assert.Equal(t, "Complex research account with all features", account.Description)
	assert.Equal(t, "Research University", account.Organization)
	assert.Equal(t, []string{"pi1", "pi2", "admin1"}, account.Coordinators)
	assert.Equal(t, "normal", account.DefaultQoS)
	assert.Equal(t, []string{"debug", "normal", "high", "urgent"}, account.QoSList)
	assert.Equal(t, "root-account", account.ParentName)
	assert.Equal(t, []string{"debug", "batch", "gpu", "highmem"}, account.AllowedPartitions)
	assert.Equal(t, "batch", account.DefaultPartition)
	assert.Equal(t, int32(150), account.FairShare)
	assert.Equal(t, int32(300), account.SharesRaw)
	assert.Equal(t, int32(100), account.Priority)
	assert.Equal(t, int32(2000), account.MaxJobs)
	assert.Equal(t, int32(100), account.MaxJobsPerUser)
	assert.Equal(t, int32(4000), account.MaxSubmitJobs)
	assert.Equal(t, int32(20160), account.MaxWallTime)
	assert.Equal(t, int32(40320), account.MaxCPUTime)
	assert.Equal(t, int32(200), account.MaxNodes)
	assert.Equal(t, int32(4000), account.MaxCPUs)
	assert.Equal(t, int64(8000*GB), account.MaxMemory)
	assert.Equal(t, int32(50), account.MinPriorityThreshold)
	assert.Equal(t, int32(1000), account.GrpJobs)
	assert.Equal(t, int32(200), account.GrpJobsAccrue)
	assert.Equal(t, int32(100), account.GrpNodes)
	assert.Equal(t, int32(2000), account.GrpCPUs)
	assert.Equal(t, int64(4000*GB), account.GrpMemory)
	assert.Equal(t, int32(2000), account.GrpSubmitJobs)
	assert.Equal(t, int32(10080), account.GrpWallTime)
	assert.Equal(t, int32(20160), account.GrpCPUTime)

	// Verify TRES maps
	require.NotNil(t, account.GrpTRES)
	assert.Equal(t, int64(2000), account.GrpTRES["cpu"])
	assert.Equal(t, int64(4000*GB), account.GrpTRES["mem"])
	assert.Equal(t, int64(20), account.GrpTRES["gpu"])

	require.NotNil(t, account.GrpTRESMins)
	assert.Equal(t, int64(200000), account.GrpTRESMins["cpu"])
	assert.Equal(t, int64(100000), account.GrpTRESMins["gpu"])

	require.NotNil(t, account.GrpTRESRunMins)
	assert.Equal(t, int64(150000), account.GrpTRESRunMins["cpu"])
	assert.Equal(t, int64(75000), account.GrpTRESRunMins["gpu"])

	require.NotNil(t, account.MaxTRES)
	assert.Equal(t, int64(4000), account.MaxTRES["cpu"])
	assert.Equal(t, int64(8000*GB), account.MaxTRES["mem"])
	assert.Equal(t, int64(200), account.MaxTRES["node"])
	assert.Equal(t, int64(40), account.MaxTRES["gpu"])

	require.NotNil(t, account.MaxTRESPerNode)
	assert.Equal(t, int64(8), account.MaxTRESPerNode["gpu"])
	assert.Equal(t, int64(512*GB), account.MaxTRESPerNode["mem"])

	require.NotNil(t, account.MinTRES)
	assert.Equal(t, int64(1), account.MinTRES["cpu"])
	assert.Equal(t, int64(1*GB), account.MinTRES["mem"])
}
