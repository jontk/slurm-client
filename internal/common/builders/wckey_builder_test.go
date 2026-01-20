// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package builders

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWCKeyBuilder_Basic(t *testing.T) {
	wckey, err := NewWCKeyBuilder("ml-training").
		WithUser("researcher1").
		WithCluster("gpu-cluster").
		Build()

	require.NoError(t, err)
	assert.Equal(t, "ml-training", wckey.Name)
	assert.Equal(t, "researcher1", wckey.User)
	assert.Equal(t, "gpu-cluster", wckey.Cluster)
}

func TestWCKeyBuilder_ValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		builder func() *WCKeyBuilder
		errMsg  string
	}{
		{
			name: "missing user",
			builder: func() *WCKeyBuilder {
				return NewWCKeyBuilder("test").WithCluster("cluster1")
			},
			errMsg: "WCKey user is required",
		},
		{
			name: "missing cluster",
			builder: func() *WCKeyBuilder {
				return NewWCKeyBuilder("test").WithUser("user1")
			},
			errMsg: "WCKey cluster is required",
		},
		{
			name: "empty name",
			builder: func() *WCKeyBuilder {
				return NewWCKeyBuilder("").WithUser("user1").WithCluster("cluster1")
			},
			errMsg: "WCKey name is required",
		},
		{
			name: "empty user",
			builder: func() *WCKeyBuilder {
				return NewWCKeyBuilder("test").WithUser("").WithCluster("cluster1")
			},
			errMsg: "user cannot be empty",
		},
		{
			name: "empty cluster",
			builder: func() *WCKeyBuilder {
				return NewWCKeyBuilder("test").WithUser("user1").WithCluster("")
			},
			errMsg: "cluster cannot be empty",
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

func TestWCKeyBuilder_Clone(t *testing.T) {
	original := NewWCKeyBuilder("original").
		WithUser("user1").
		WithCluster("cluster1")

	clone := original.Clone()

	// Modify clone
	clone.WithUser("user2")

	// Build both
	originalWCKey, err := original.Build()
	require.NoError(t, err)

	cloneWCKey, err := clone.Build()
	require.NoError(t, err)

	// Verify they're different
	assert.Equal(t, "user1", originalWCKey.User)
	assert.Equal(t, "user2", cloneWCKey.User)

	// But share same name and cluster
	assert.Equal(t, originalWCKey.Name, cloneWCKey.Name)
	assert.Equal(t, originalWCKey.Cluster, cloneWCKey.Cluster)
}

func TestWCKeyBuilder_Reset(t *testing.T) {
	builder := NewWCKeyBuilder("test").
		WithUser("user1").
		WithCluster("cluster1")

	// Reset the builder
	builder.Reset()

	// Should fail validation after reset
	err := builder.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "WCKey name is required")
}

func TestWCKeyBuilder_Validate(t *testing.T) {
	// Valid builder
	builder := NewWCKeyBuilder("test").
		WithUser("user1").
		WithCluster("cluster1")

	err := builder.Validate()
	assert.NoError(t, err)

	// Invalid builder - missing user
	builder2 := NewWCKeyBuilder("test").WithCluster("cluster1")
	err = builder2.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "WCKey user is required")
}

func TestWCKeyBuilder_String(t *testing.T) {
	builder := NewWCKeyBuilder("ml-training").
		WithUser("researcher1").
		WithCluster("gpu-cluster")

	str := builder.String()
	assert.Equal(t, "WCKey{Name: ml-training, User: researcher1, Cluster: gpu-cluster}", str)
}

func TestWCKeyBuilder_MultipleErrors(t *testing.T) {
	builder := NewWCKeyBuilder("test").
		WithUser("").   // This will add an error
		WithCluster("") // This will add another error

	_, err := builder.Build()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "user cannot be empty")
	assert.Contains(t, err.Error(), "cluster cannot be empty")
}

func TestWCKeyBuilder_Chaining(t *testing.T) {
	// Test that all methods return the builder for chaining
	builder := NewWCKeyBuilder("test")

	assert.Equal(t, builder, builder.WithUser("user1"))
	assert.Equal(t, builder, builder.WithCluster("cluster1"))
	assert.Equal(t, builder, builder.Reset())

	// Even with errors, chaining should work
	assert.Equal(t, builder, builder.WithUser(""))
	assert.Equal(t, builder, builder.WithCluster(""))
}
