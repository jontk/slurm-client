package v0_0_41

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
	assert.Contains(t, err.Error(), "API client not initialized")
	// The impl should now be created
	assert.NotNil(t, infoManager.impl)
}

func TestInfoManager_Ping_Structure(t *testing.T) {
	// Test that Ping method properly delegates to implementation
	infoManager := &InfoManager{
		client: &WrapperClient{},
	}

	err := infoManager.Ping(context.Background())

	// We expect an error since there's no real API client
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API client not initialized")
	// The impl should now be created
	assert.NotNil(t, infoManager.impl)
}

func TestInfoManager_Stats_NotImplemented(t *testing.T) {
	// Test that Stats returns not implemented error
	infoManager := &InfoManager{
		client: &WrapperClient{},
	}

	_, err := infoManager.Stats(context.Background())

	// v0.0.41 Stats is not implemented
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not implemented")
	// The impl should now be created
	assert.NotNil(t, infoManager.impl)
}

func TestInfoManager_Version_Structure(t *testing.T) {
	// Test that Version method properly delegates to implementation
	infoManager := &InfoManager{
		client: &WrapperClient{},
	}

	_, err := infoManager.Version(context.Background())

	// We expect an error since there's no real API client
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API client not initialized")
	// The impl should now be created
	assert.NotNil(t, infoManager.impl)
}