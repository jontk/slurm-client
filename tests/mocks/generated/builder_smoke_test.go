// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package generated_test

import (
	"testing"

	builderv0040 "github.com/jontk/slurm-client/tests/mocks/generated/v0_0_40"
	builderv0042 "github.com/jontk/slurm-client/tests/mocks/generated/v0_0_42"
	builderv0043 "github.com/jontk/slurm-client/tests/mocks/generated/v0_0_43"
	builderv0044 "github.com/jontk/slurm-client/tests/mocks/generated/v0_0_44"
)

// TestV0040BuilderCompleteness verifies essential fields exist in v0.0.40 builder
func TestV0040BuilderCompleteness(t *testing.T) {
	job := builderv0040.NewJobInfo().
		WithJobId(1001).
		WithName("test").
		WithCpus(4).
		WithMemoryPerNode(8 * 1024 * 1024 * 1024).
		WithSubmitTime(1234567890).
		WithTimeLimit(60).
		WithJobState("RUNNING").
		Build()

	if job == nil {
		t.Fatal("Builder returned nil")
	}

	// Verify fields are set
	if job.JobId == nil || *job.JobId != 1001 {
		t.Error("JobId not set correctly")
	}
	if job.Name == nil || *job.Name != "test" {
		t.Error("Name not set correctly")
	}
	if job.Cpus == nil || job.Cpus.Number == nil || *job.Cpus.Number != 4 {
		t.Error("Cpus not set correctly")
	}
	if job.MemoryPerNode == nil || job.MemoryPerNode.Number == nil {
		t.Error("MemoryPerNode not set correctly")
	}
	if job.SubmitTime == nil || job.SubmitTime.Number == nil {
		t.Error("SubmitTime not set correctly")
	}
}

// TestV0042BuilderCompleteness verifies essential fields exist in v0.0.42 builder
func TestV0042BuilderCompleteness(t *testing.T) {
	job := builderv0042.NewJobInfo().
		WithJobId(1001).
		WithName("test").
		WithCpus(4).
		WithMemoryPerNode(8 * 1024 * 1024 * 1024).
		WithSubmitTime(1234567890).
		WithTimeLimit(60).
		WithJobState("RUNNING").
		Build()

	if job == nil {
		t.Fatal("Builder returned nil")
	}

	// Verify fields are set
	if job.JobId == nil || *job.JobId != 1001 {
		t.Error("JobId not set correctly")
	}
	if job.Name == nil || *job.Name != "test" {
		t.Error("Name not set correctly")
	}
	if job.Cpus == nil || job.Cpus.Number == nil || *job.Cpus.Number != 4 {
		t.Error("Cpus not set correctly")
	}
	if job.MemoryPerNode == nil || job.MemoryPerNode.Number == nil {
		t.Error("MemoryPerNode not set correctly")
	}
	if job.SubmitTime == nil || job.SubmitTime.Number == nil {
		t.Error("SubmitTime not set correctly")
	}
}

// TestV0043BuilderCompleteness verifies essential fields exist in v0.0.43 builder
func TestV0043BuilderCompleteness(t *testing.T) {
	job := builderv0043.NewJobInfo().
		WithJobId(1001).
		WithName("test").
		WithCpus(4).
		WithMemoryPerNode(8 * 1024 * 1024 * 1024).
		WithSubmitTime(1234567890).
		WithTimeLimit(60).
		WithJobState("RUNNING").
		Build()

	if job == nil {
		t.Fatal("Builder returned nil")
	}

	// Verify fields are set
	if job.JobId == nil || *job.JobId != 1001 {
		t.Error("JobId not set correctly")
	}
	if job.Name == nil || *job.Name != "test" {
		t.Error("Name not set correctly")
	}
	if job.Cpus == nil || job.Cpus.Number == nil || *job.Cpus.Number != 4 {
		t.Error("Cpus not set correctly")
	}
	if job.MemoryPerNode == nil || job.MemoryPerNode.Number == nil {
		t.Error("MemoryPerNode not set correctly")
	}
	if job.SubmitTime == nil || job.SubmitTime.Number == nil {
		t.Error("SubmitTime not set correctly")
	}
}

// TestV0044BuilderCompleteness verifies essential fields exist in v0.0.44 builder
func TestV0044BuilderCompleteness(t *testing.T) {
	job := builderv0044.NewJobInfo().
		WithJobId(1001).
		WithName("test").
		WithCpus(4).
		WithMemoryPerNode(8 * 1024 * 1024 * 1024).
		WithSubmitTime(1234567890).
		WithTimeLimit(60).
		WithJobState("RUNNING").
		Build()

	if job == nil {
		t.Fatal("Builder returned nil")
	}

	// Verify fields are set
	if job.JobId == nil || *job.JobId != 1001 {
		t.Error("JobId not set correctly")
	}
	if job.Name == nil || *job.Name != "test" {
		t.Error("Name not set correctly")
	}
	if job.Cpus == nil || job.Cpus.Number == nil || *job.Cpus.Number != 4 {
		t.Error("Cpus not set correctly")
	}
	if job.MemoryPerNode == nil || job.MemoryPerNode.Number == nil {
		t.Error("MemoryPerNode not set correctly")
	}
	if job.SubmitTime == nil || job.SubmitTime.Number == nil {
		t.Error("SubmitTime not set correctly")
	}
}
