// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package helpers

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestContext returns a test context with timeout
func TestContext(t *testing.T) context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute) // 10 minute timeout
	t.Cleanup(cancel)
	return ctx
}

// AssertNoError is a helper that fails the test if error is not nil
func AssertNoError(t *testing.T, err error) {
	t.Helper()
	assert.NoError(t, err)
}

// RequireNoError is a helper that fails the test immediately if error is not nil
func RequireNoError(t *testing.T, err error) {
	t.Helper()
	require.NoError(t, err)
}

// AssertEqual is a helper for equality assertions
func AssertEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()
	assert.Equal(t, expected, actual)
}

// RequireEqual is a helper for equality assertions that fails immediately
func RequireEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()
	require.Equal(t, expected, actual)
}

// AssertNotNil is a helper for nil checks
func AssertNotNil(t *testing.T, obj interface{}) {
	t.Helper()
	assert.NotNil(t, obj)
}

// RequireNotNil is a helper for nil checks that fails immediately
func RequireNotNil(t *testing.T, obj interface{}) {
	t.Helper()
	require.NotNil(t, obj)
}

// IntPtr returns a pointer to an int value
func IntPtr(v int) *int {
	return &v
}

// Int32Ptr returns a pointer to an int32 value
func Int32Ptr(v int32) *int32 {
	return &v
}

// Int64Ptr returns a pointer to an int64 value
func Int64Ptr(v int64) *int64 {
	return &v
}

// StringPtr returns a pointer to a string value
func StringPtr(v string) *string {
	return &v
}

// BoolPtr returns a pointer to a bool value
func BoolPtr(v bool) *bool {
	return &v
}
