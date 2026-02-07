// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package builders

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJobBuilder_Basic(t *testing.T) {
	t.Run("simple job with script", func(t *testing.T) {
		job, err := NewJobBuilder("echo 'Hello World'").
			WithName("test-job").
			WithAccount("test-account").
			WithPartition("debug").
			Build()

		require.NoError(t, err)
		require.NotNil(t, job.Name)
		require.NotNil(t, job.Script)
		require.NotNil(t, job.Account)
		require.NotNil(t, job.Partition)
		assert.Equal(t, "test-job", *job.Name)
		assert.Equal(t, "echo 'Hello World'", *job.Script)
		assert.Equal(t, "test-account", *job.Account)
		assert.Equal(t, "debug", *job.Partition)
		assert.Equal(t, uint32(30), *job.TimeLimit) // Default
		assert.Equal(t, int32(1), *job.MinimumCPUs) // Default
	})

	t.Run("job with batch script", func(t *testing.T) {
		script := `#!/bin/bash
echo "Running batch job"
sleep 10`
		job, err := NewJobBuilderFromScript(script).
			WithName("batch-job").
			Build()

		require.NoError(t, err)
		require.NotNil(t, job.Name)
		require.NotNil(t, job.Script)
		assert.Equal(t, "batch-job", *job.Name)
		assert.Equal(t, script, *job.Script)
	})

	t.Run("empty name returns error", func(t *testing.T) {
		_, err := NewJobBuilder("test").
			WithName("").
			Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "job name cannot be empty")
	})
}

func TestJobBuilder_Resources(t *testing.T) {
	t.Run("basic resources", func(t *testing.T) {
		job, err := NewJobBuilder("./my_program").
			WithName("compute-job").
			WithCPUs(8).
			WithNodes(2).
			WithTasks(16).
			WithTimeLimit(120).
			Build()

		require.NoError(t, err)
		require.NotNil(t, job.MinimumCPUs)
		require.NotNil(t, job.MinimumNodes)
		require.NotNil(t, job.Tasks)
		require.NotNil(t, job.TimeLimit)
		assert.Equal(t, int32(8), *job.MinimumCPUs)
		assert.Equal(t, int32(2), *job.MinimumNodes)
		assert.Equal(t, int32(16), *job.Tasks)
		assert.Equal(t, uint32(120), *job.TimeLimit)
	})

	t.Run("memory resources", func(t *testing.T) {
		job, err := NewJobBuilder("./memory_intensive").
			WithName("memory-job").
			WithMemoryPerNode(32768). // 32GB in MB
			WithMemoryPerCPU(4096).   // 4GB in MB
			Build()

		require.NoError(t, err)
		require.NotNil(t, job.MemoryPerNode)
		require.NotNil(t, job.MemoryPerCPU)
		assert.Equal(t, uint64(32768), *job.MemoryPerNode)
		assert.Equal(t, uint64(4096), *job.MemoryPerCPU)
	})

	t.Run("negative CPUs fail", func(t *testing.T) {
		_, err := NewJobBuilder("test").
			WithCPUs(-1).
			Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must be positive")
	})

	t.Run("negative time limit fail", func(t *testing.T) {
		_, err := NewJobBuilder("test").
			WithTimeLimit(-1).
			Build()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must be positive")
	})
}

func TestJobBuilder_Environment(t *testing.T) {
	t.Run("environment variables", func(t *testing.T) {
		env := map[string]string{
			"OMP_NUM_THREADS":      "8",
			"CUDA_VISIBLE_DEVICES": "0,1",
		}

		job, err := NewJobBuilder("./parallel_app").
			WithName("parallel-job").
			WithEnvironment(env).
			WithEnvironmentVariable("CUSTOM_VAR", "custom_value").
			Build()

		require.NoError(t, err)
		require.NotNil(t, job.Environment)
		// Environment is now []string in "KEY=VALUE" format
		assert.Contains(t, job.Environment, "CUSTOM_VAR=custom_value")
		// Note: map iteration order is not guaranteed, so we check contains
		foundOMP := false
		foundCUDA := false
		for _, e := range job.Environment {
			if e == "OMP_NUM_THREADS=8" {
				foundOMP = true
			}
			if e == "CUDA_VISIBLE_DEVICES=0,1" {
				foundCUDA = true
			}
		}
		assert.True(t, foundOMP, "OMP_NUM_THREADS not found")
		assert.True(t, foundCUDA, "CUDA_VISIBLE_DEVICES not found")
	})
}

func TestJobBuilder_Dependencies(t *testing.T) {
	t.Run("job dependency string", func(t *testing.T) {
		job, err := NewJobBuilder("./dependent_job").
			WithName("dependent").
			WithDependency("afterok:12345:12346").
			Build()

		require.NoError(t, err)
		require.NotNil(t, job.Dependency)
		assert.Equal(t, "afterok:12345:12346", *job.Dependency)
	})
}

func TestJobBuilder_TimeLimit(t *testing.T) {
	t.Run("time limit duration", func(t *testing.T) {
		job, err := NewJobBuilder("test").
			WithName("duration-job").
			WithTimeLimitDuration(2 * time.Hour).
			Build()

		require.NoError(t, err)
		require.NotNil(t, job.TimeLimit)
		assert.Equal(t, uint32(120), *job.TimeLimit) // 2 hours = 120 minutes
	})
}

func TestJobBuilder_IO(t *testing.T) {
	t.Run("standard io paths", func(t *testing.T) {
		job, err := NewJobBuilder("./my_script.sh").
			WithName("io-job").
			WithWorkingDirectory("/scratch/user").
			WithStandardOutput("/scratch/user/output-%j.txt").
			WithStandardError("/scratch/user/error-%j.txt").
			WithStandardInput("/scratch/user/input.txt").
			Build()

		require.NoError(t, err)
		require.NotNil(t, job.CurrentWorkingDirectory)
		require.NotNil(t, job.StandardOutput)
		require.NotNil(t, job.StandardError)
		require.NotNil(t, job.StandardInput)
		assert.Equal(t, "/scratch/user", *job.CurrentWorkingDirectory)
		assert.Equal(t, "/scratch/user/output-%j.txt", *job.StandardOutput)
		assert.Equal(t, "/scratch/user/error-%j.txt", *job.StandardError)
		assert.Equal(t, "/scratch/user/input.txt", *job.StandardInput)
	})
}

func TestJobBuilder_Features(t *testing.T) {
	t.Run("feature constraints", func(t *testing.T) {
		job, err := NewJobBuilder("./my_script.sh").
			WithName("feature-job").
			WithFeatures("gpu", "nvme", "avx2").
			Build()

		require.NoError(t, err)
		require.NotNil(t, job.Constraints)
		// Features joined with & for SLURM constraint syntax
		assert.Equal(t, "gpu&nvme&avx2", *job.Constraints)
	})

	t.Run("constraint string", func(t *testing.T) {
		job, err := NewJobBuilder("./my_script.sh").
			WithName("constraint-job").
			WithConstraints("(gpu|a100)&nvme").
			Build()

		require.NoError(t, err)
		require.NotNil(t, job.Constraints)
		assert.Equal(t, "(gpu|a100)&nvme", *job.Constraints)
	})
}

func TestJobBuilder_Array(t *testing.T) {
	t.Run("array job", func(t *testing.T) {
		job, err := NewJobBuilder("./array_task.sh").
			WithName("array-job").
			WithArray("1-100%10"). // Indexes 1-100 with max 10 concurrent
			Build()

		require.NoError(t, err)
		require.NotNil(t, job.Array)
		assert.Equal(t, "1-100%10", *job.Array)
	})
}

func TestJobBuilder_Mail(t *testing.T) {
	t.Run("mail notifications", func(t *testing.T) {
		job, err := NewJobBuilder("./my_script.sh").
			WithName("mail-job").
			WithMailUser("user@example.com").
			Build()

		require.NoError(t, err)
		require.NotNil(t, job.MailUser)
		assert.Equal(t, "user@example.com", *job.MailUser)
	})
}

func TestJobBuilder_Scheduling(t *testing.T) {
	t.Run("reservation", func(t *testing.T) {
		job, err := NewJobBuilder("./my_script.sh").
			WithName("reserved-job").
			WithReservation("my-reservation").
			Build()

		require.NoError(t, err)
		require.NotNil(t, job.Reservation)
		assert.Equal(t, "my-reservation", *job.Reservation)
	})

	t.Run("excluded nodes", func(t *testing.T) {
		job, err := NewJobBuilder("./my_script.sh").
			WithName("exclude-job").
			WithExcludedNodes("node1", "node2", "node3").
			Build()

		require.NoError(t, err)
		require.NotNil(t, job.ExcludedNodes)
		assert.Equal(t, []string{"node1", "node2", "node3"}, job.ExcludedNodes)
	})

	t.Run("immediate mode", func(t *testing.T) {
		job, err := NewJobBuilder("./my_script.sh").
			WithName("immediate-job").
			WithImmediate(true).
			Build()

		require.NoError(t, err)
		require.NotNil(t, job.Immediate)
		assert.True(t, *job.Immediate)
	})
}

func TestJobBuilder_ErrorAccumulation(t *testing.T) {
	t.Run("multiple errors accumulate", func(t *testing.T) {
		builder := NewJobBuilder("test").
			WithCPUs(-1).
			WithNodes(-1).
			WithTimeLimit(-1)

		assert.True(t, builder.HasErrors())
		assert.Len(t, builder.Errors(), 3)

		_, err := builder.Build()
		require.Error(t, err)
	})
}

func TestJobBuilder_MustBuild(t *testing.T) {
	t.Run("MustBuild panics on error", func(t *testing.T) {
		assert.Panics(t, func() {
			NewJobBuilder("test").
				WithCPUs(-1).
				MustBuild()
		})
	})

	t.Run("MustBuild succeeds without error", func(t *testing.T) {
		assert.NotPanics(t, func() {
			job := NewJobBuilder("test").
				WithName("valid-job").
				MustBuild()
			assert.NotNil(t, job)
		})
	})
}
