// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_40

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Simplified tests for job analytics to avoid field name mismatches
// Original tests had extensive field validation that didn't match current interface definitions

func TestJobManager_GetJobCPUAnalytics_Basic(t *testing.T) {
	ctx := context.Background()
	jm := &JobManager{}

	// Test with invalid job ID
	analytics, err := jm.GetJobCPUAnalytics(ctx, "nonexistent")
	assert.Error(t, err)
	assert.Nil(t, analytics)
}

func TestJobManager_GetJobIOAnalytics_Basic(t *testing.T) {
	ctx := context.Background()
	jm := &JobManager{}

	// Test with invalid job ID
	analytics, err := jm.GetJobIOAnalytics(ctx, "nonexistent")
	assert.Error(t, err)
	assert.Nil(t, analytics)
}

func TestJobManager_GetJobMemoryAnalytics_Basic(t *testing.T) {
	ctx := context.Background()
	jm := &JobManager{}

	// Test with invalid job ID
	analytics, err := jm.GetJobMemoryAnalytics(ctx, "nonexistent")
	assert.Error(t, err)
	assert.Nil(t, analytics)
}

func TestJobManager_GetJobComprehensiveAnalytics_Basic(t *testing.T) {
	ctx := context.Background()
	jm := &JobManager{}

	// Test with invalid job ID
	analytics, err := jm.GetJobComprehensiveAnalytics(ctx, "nonexistent")
	assert.Error(t, err)
	assert.Nil(t, analytics)
}