package v0_0_41

import (
	"context"
	"testing"

	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/stretchr/testify/assert"
)

func TestPartitionManager_List_NotImplemented(t *testing.T) {
	// Test that List returns not implemented error
	partitionManager := &PartitionManager{
		client: &WrapperClient{},
	}

	_, err := partitionManager.List(context.Background(), nil)

	// v0.0.41 PartitionManager is not implemented
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not implemented")
	// The impl should now be created
	assert.NotNil(t, partitionManager.impl)
}

func TestPartitionManager_Get_NotImplemented(t *testing.T) {
	// Test that Get returns not implemented error
	partitionManager := &PartitionManager{
		client: &WrapperClient{},
	}

	_, err := partitionManager.Get(context.Background(), "gpu")

	// v0.0.41 PartitionManager is not implemented
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not implemented")
	// The impl should now be created
	assert.NotNil(t, partitionManager.impl)
}

func TestPartitionManager_Update_NotSupported(t *testing.T) {
	// Test that Update returns not supported error
	partitionManager := &PartitionManager{
		client: &WrapperClient{},
	}

	err := partitionManager.Update(context.Background(), "gpu", &interfaces.PartitionUpdate{})

	// Partition updates are not supported in v0.0.41
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not supported")
	// The impl should now be created
	assert.NotNil(t, partitionManager.impl)
}

func TestPartitionManager_Watch_Structure(t *testing.T) {
	// Test that Watch method properly delegates to implementation
	partitionManager := &PartitionManager{
		client: &WrapperClient{},
	}

	_, err := partitionManager.Watch(context.Background(), &interfaces.WatchPartitionsOptions{})

	// We expect an error since there's no real API client
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API client not initialized")
	// The impl should now be created
	assert.NotNil(t, partitionManager.impl)
}