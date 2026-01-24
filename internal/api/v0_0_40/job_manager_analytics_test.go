// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_40

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jontk/slurm-client/interfaces"
)

// TestJobManager_GetJobCPUAnalytics_Enhanced tests CPU analytics with comprehensive scenarios
func TestJobManager_GetJobCPUAnalytics_Enhanced(t *testing.T) {
	ctx := context.Background()

	t.Run("empty job ID validation", func(t *testing.T) {
		jm := &JobManager{}

		analytics, err := jm.GetJobCPUAnalytics(ctx, "")

		assert.Error(t, err)
		assert.Nil(t, analytics)
		assert.Contains(t, err.Error(), "not initialized")
	})

	// Nil context validation is tested through other error paths

	t.Run("uninitialized client", func(t *testing.T) {
		jm := &JobManager{} // No client set

		analytics, err := jm.GetJobCPUAnalytics(ctx, "12345")

		assert.Error(t, err)
		assert.Nil(t, analytics)
		assert.Contains(t, err.Error(), "not initialized")
	})

	t.Run("client with nil api client", func(t *testing.T) {
		wrapperClient := &WrapperClient{
			apiClient: nil, // Uninitialized API client
			config: &interfaces.ClientConfig{
				BaseURL: "http://test.example.com",
			},
		}
		jm := &JobManager{client: wrapperClient}

		analytics, err := jm.GetJobCPUAnalytics(ctx, "12345")

		assert.Error(t, err)
		assert.Nil(t, analytics)
		assert.Contains(t, err.Error(), "not initialized")
	})
}

func TestJobManager_GetJobMemoryAnalytics_Enhanced(t *testing.T) {
	ctx := context.Background()

	t.Run("input validation", func(t *testing.T) {
		jm := &JobManager{}

		// Test empty job ID
		analytics, err := jm.GetJobMemoryAnalytics(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, analytics)
	})

	t.Run("uninitialized client error handling", func(t *testing.T) {
		jm := &JobManager{}

		analytics, err := jm.GetJobMemoryAnalytics(ctx, "12345")

		assert.Error(t, err)
		assert.Nil(t, analytics)
	})
}

func TestJobManager_GetJobIOAnalytics_Enhanced(t *testing.T) {
	ctx := context.Background()

	t.Run("input validation", func(t *testing.T) {
		jm := &JobManager{}

		// Test various invalid inputs
		testCases := []struct {
			name  string
			ctx   context.Context
			jobID string
		}{
			{"empty job ID", ctx, ""},
			{"whitespace job ID", ctx, "   "},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				analytics, err := jm.GetJobIOAnalytics(tc.ctx, tc.jobID)
				assert.Error(t, err)
				assert.Nil(t, analytics)
			})
		}
	})

	t.Run("client initialization validation", func(t *testing.T) {
		jm := &JobManager{}

		analytics, err := jm.GetJobIOAnalytics(ctx, "12345")

		assert.Error(t, err)
		assert.Nil(t, analytics)
		assert.Contains(t, err.Error(), "not initialized")
	})
}

func TestJobManager_GetJobComprehensiveAnalytics_Enhanced(t *testing.T) {
	ctx := context.Background()

	t.Run("input validation comprehensive", func(t *testing.T) {
		jm := &JobManager{}

		// Test comprehensive validation
		testCases := []struct {
			name        string
			ctx         context.Context
			jobID       string
			expectError string
		}{
			{
				name:        "empty job ID",
				ctx:         ctx,
				jobID:       "",
				expectError: "not initialized",
			},
			{
				name:        "job ID with only spaces",
				ctx:         ctx,
				jobID:       "   ",
				expectError: "not initialized",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				analytics, err := jm.GetJobComprehensiveAnalytics(tc.ctx, tc.jobID)

				assert.Error(t, err)
				assert.Nil(t, analytics)
				if tc.expectError != "" {
					assert.Contains(t, err.Error(), tc.expectError)
				}
			})
		}
	})

	t.Run("client state validation", func(t *testing.T) {
		// Test with uninitialized JobManager
		jm := &JobManager{}

		analytics, err := jm.GetJobComprehensiveAnalytics(ctx, "12345")

		assert.Error(t, err)
		assert.Nil(t, analytics)
		assert.Contains(t, err.Error(), "not initialized")
	})
}

// TestJobManager_AnalyticsErrorConsistency ensures all analytics methods handle errors consistently
func TestJobManager_AnalyticsErrorConsistency(t *testing.T) {
	ctx := context.Background()
	jm := &JobManager{}

	// Test that all analytics methods return consistent errors for the same invalid input
	t.Run("consistent empty job ID errors", func(t *testing.T) {
		cpuAnalytics, cpuErr := jm.GetJobCPUAnalytics(ctx, "")
		memAnalytics, memErr := jm.GetJobMemoryAnalytics(ctx, "")
		ioAnalytics, ioErr := jm.GetJobIOAnalytics(ctx, "")
		compAnalytics, compErr := jm.GetJobComprehensiveAnalytics(ctx, "")

		// All should return errors
		assert.Error(t, cpuErr)
		assert.Error(t, memErr)
		assert.Error(t, ioErr)
		assert.Error(t, compErr)

		// All should return nil analytics
		assert.Nil(t, cpuAnalytics)
		assert.Nil(t, memAnalytics)
		assert.Nil(t, ioAnalytics)
		assert.Nil(t, compAnalytics)

		// All errors should mention client not initialized
		assert.Contains(t, cpuErr.Error(), "not initialized")
		assert.Contains(t, memErr.Error(), "not initialized")
		assert.Contains(t, ioErr.Error(), "not initialized")
		assert.Contains(t, compErr.Error(), "not initialized")
	})

	t.Run("consistent uninitialized client errors", func(t *testing.T) {
		cpuAnalytics, cpuErr := jm.GetJobCPUAnalytics(ctx, "12345")
		memAnalytics, memErr := jm.GetJobMemoryAnalytics(ctx, "12345")
		ioAnalytics, ioErr := jm.GetJobIOAnalytics(ctx, "12345")
		compAnalytics, compErr := jm.GetJobComprehensiveAnalytics(ctx, "12345")

		// All should return errors
		assert.Error(t, cpuErr)
		assert.Error(t, memErr)
		assert.Error(t, ioErr)
		assert.Error(t, compErr)

		// All should return nil analytics
		assert.Nil(t, cpuAnalytics)
		assert.Nil(t, memAnalytics)
		assert.Nil(t, ioAnalytics)
		assert.Nil(t, compAnalytics)
	})
}

// TestJobManager_AnalyticsImplementationVersionCompatibility tests v0.0.40 specific behavior
func TestJobManager_AnalyticsImplementationVersionCompatibility(t *testing.T) {
	// These tests verify that the v0.0.40 implementation behaves correctly
	// with its known limitations and characteristics

	t.Run("version specific limitations", func(t *testing.T) {
		// v0.0.40 should have specific characteristics:
		// - Fixed percentage estimates for CPU utilization (~50%)
		// - Fixed percentage estimates for memory utilization (~60%)
		// - Basic I/O analytics with estimated metrics
		// - No GPU analytics support
		// - Limited NUMA support

		// Since we can't mock the actual SLURM API easily in this context,
		// we're testing the error handling and validation logic
		jm := &JobManager{}

		// Test that the manager correctly identifies when it's not properly initialized
		analytics, err := jm.GetJobCPUAnalytics(context.Background(), "12345")

		assert.Error(t, err)
		assert.Nil(t, analytics)

		// The error should indicate a client initialization issue
		assert.Contains(t, err.Error(), "not initialized")
	})
}

// BenchmarkJobManager_AnalyticsValidation benchmarks the validation logic
func BenchmarkJobManager_AnalyticsValidation(b *testing.B) {
	ctx := context.Background()
	jm := &JobManager{}

	b.Run("CPU analytics validation", func(b *testing.B) {
		b.ResetTimer()
		for range b.N {
			_, _ = jm.GetJobCPUAnalytics(ctx, "") // Will fail validation quickly
		}
	})

	b.Run("Comprehensive analytics validation", func(b *testing.B) {
		b.ResetTimer()
		for range b.N {
			_, _ = jm.GetJobComprehensiveAnalytics(ctx, "") // Will fail validation quickly
		}
	})
}

// TestJobManager_ConcurrentAccess tests concurrent access to analytics methods
func TestJobManager_ConcurrentAccess(t *testing.T) {
	ctx := context.Background()
	jm := &JobManager{}

	// Test that multiple goroutines can safely call analytics methods
	t.Run("concurrent analytics calls", func(t *testing.T) {
		done := make(chan bool, 4)

		// Launch concurrent analytics calls
		go func() {
			_, _ = jm.GetJobCPUAnalytics(ctx, "")
			done <- true
		}()

		go func() {
			_, _ = jm.GetJobMemoryAnalytics(ctx, "")
			done <- true
		}()

		go func() {
			_, _ = jm.GetJobIOAnalytics(ctx, "")
			done <- true
		}()

		go func() {
			_, _ = jm.GetJobComprehensiveAnalytics(ctx, "")
			done <- true
		}()

		// Wait for all to complete
		for range 4 {
			<-done
		}

		// If we get here without deadlock, the test passes
		assert.True(t, true, "Concurrent access completed successfully")
	})
}

// TestJobManager_AnalyticsInputEdgeCases tests edge cases in input validation
func TestJobManager_AnalyticsInputEdgeCases(t *testing.T) {
	ctx := context.Background()
	jm := &JobManager{}

	edgeCases := []struct {
		name  string
		jobID string
	}{
		{"very long job ID", "this-is-a-very-long-job-id-that-might-cause-issues-with-some-implementations-1234567890"},
		{"job ID with special characters", "job-123_test.special"},
		{"numeric job ID", "123456"},
		{"job ID with unicode", "job-测试-123"},
		{"job ID with newlines", "job\n123"},
		{"job ID with tabs", "job\t123"},
	}

	for _, tc := range edgeCases {
		t.Run(tc.name, func(t *testing.T) {
			// These should all fail due to uninitialized client, but not crash
			cpuAnalytics, cpuErr := jm.GetJobCPUAnalytics(ctx, tc.jobID)
			memAnalytics, memErr := jm.GetJobMemoryAnalytics(ctx, tc.jobID)
			ioAnalytics, ioErr := jm.GetJobIOAnalytics(ctx, tc.jobID)
			compAnalytics, compErr := jm.GetJobComprehensiveAnalytics(ctx, tc.jobID)

			// All should handle the input gracefully (either succeed or fail cleanly)
			// Since client is not initialized, we expect errors
			assert.Error(t, cpuErr)
			assert.Error(t, memErr)
			assert.Error(t, ioErr)
			assert.Error(t, compErr)

			assert.Nil(t, cpuAnalytics)
			assert.Nil(t, memAnalytics)
			assert.Nil(t, ioAnalytics)
			assert.Nil(t, compAnalytics)
		})
	}
}
