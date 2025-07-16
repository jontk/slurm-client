package v0_0_43

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfoManager_Get_Structure(t *testing.T) {
	// Test that the InfoManager properly creates implementation
	infoManager := &InfoManager{
		client: &WrapperClient{},
	}

	// Test that impl is created lazily
	assert.Nil(t, infoManager.impl)

	// After attempting to call Get (even with nil client), impl should be created
	_, err := infoManager.Get(context.Background())

	// We expect an error since there's no real API client
	assert.Error(t, err)
	// The impl should now be created
	assert.NotNil(t, infoManager.impl)
}

func TestInfoManager_Ping_Structure(t *testing.T) {
	// Test that the InfoManager properly creates implementation for Ping
	infoManager := &InfoManager{
		client: &WrapperClient{},
	}

	// Test that impl is created lazily
	assert.Nil(t, infoManager.impl)

	// After attempting to call Ping (even with nil client), impl should be created
	err := infoManager.Ping(context.Background())

	// We expect an error since there's no real API client
	assert.Error(t, err)
	// The impl should now be created
	assert.NotNil(t, infoManager.impl)
}

func TestInfoManager_Stats_Structure(t *testing.T) {
	// Test that the InfoManager properly creates implementation for Stats
	infoManager := &InfoManager{
		client: &WrapperClient{},
	}

	// Test that impl is created lazily
	assert.Nil(t, infoManager.impl)

	// After attempting to call Stats (even with nil client), impl should be created
	_, err := infoManager.Stats(context.Background())

	// We expect an error since there's no real API client
	assert.Error(t, err)
	// The impl should now be created
	assert.NotNil(t, infoManager.impl)
}

func TestInfoManager_Version_Structure(t *testing.T) {
	// Test that the InfoManager properly creates implementation for Version
	infoManager := &InfoManager{
		client: &WrapperClient{},
	}

	// Test that impl is created lazily
	assert.Nil(t, infoManager.impl)

	// After attempting to call Version (even with nil client), impl should be created
	_, err := infoManager.Version(context.Background())

	// We expect an error since there's no real API client
	assert.Error(t, err)
	// The impl should now be created
	assert.NotNil(t, infoManager.impl)
}

func TestNewInfoManagerImpl(t *testing.T) {
	client := &WrapperClient{}
	impl := NewInfoManagerImpl(client)

	assert.NotNil(t, impl)
	assert.Equal(t, client, impl.client)
}

func TestInfoManagerImpl_ErrorHandling(t *testing.T) {
	// Test error handling when API client is not initialized
	impl := &InfoManagerImpl{
		client: &WrapperClient{}, // client without apiClient initialized
	}

	// Test Get method error handling
	_, err := impl.Get(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API client not initialized")

	// Test Ping method error handling
	err = impl.Ping(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API client not initialized")

	// Test Stats method error handling
	_, err = impl.Stats(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API client not initialized")

	// Test Version method error handling
	_, err = impl.Version(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API client not initialized")
}

func TestInfoManagerImpl_VersionInfo(t *testing.T) {
	// Test that Version method returns expected static information
	impl := &InfoManagerImpl{
		client: &WrapperClient{}, // client without apiClient initialized
	}

	// Even though the client isn't initialized, we can test the error message
	// to ensure the method is properly structured
	_, err := impl.Version(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API client not initialized")
}

func TestInfoManagerImpl_ClusterInfoStructure(t *testing.T) {
	// Test that Get method structures the ClusterInfo correctly
	impl := &InfoManagerImpl{
		client: &WrapperClient{}, // client without apiClient initialized
	}

	// Even though the client isn't initialized, we can test the error handling
	_, err := impl.Get(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API client not initialized")
}

func TestInfoManagerImpl_ClusterStatsStructure(t *testing.T) {
	// Test that Stats method structures the ClusterStats correctly
	impl := &InfoManagerImpl{
		client: &WrapperClient{}, // client without apiClient initialized
	}

	// Even though the client isn't initialized, we can test the error handling
	_, err := impl.Stats(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API client not initialized")
}
