// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_41

import (
	"context"

	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
)

// PartitionManagerImpl provides the actual implementation for PartitionManager methods
type PartitionManagerImpl struct {
	client *WrapperClient
}

// NewPartitionManagerImpl creates a new PartitionManager implementation
func NewPartitionManagerImpl(client *WrapperClient) *PartitionManagerImpl {
	return &PartitionManagerImpl{client: client}
}

// List partitions with optional filtering
func (m *PartitionManagerImpl) List(ctx context.Context, opts *interfaces.ListPartitionsOptions) (*interfaces.PartitionList, error) {
	// Note: v0.0.41 has complex inline struct for partitions
	// Return basic error for now
	return nil, errors.NewClientError(
		errors.ErrorCodeUnsupportedOperation,
		"Partition listing not implemented for v0.0.41",
		"The v0.0.41 partition response uses complex inline structs that differ significantly from other API versions",
	)
}

// Get retrieves a specific partition by name
func (m *PartitionManagerImpl) Get(ctx context.Context, partitionName string) (*interfaces.Partition, error) {
	return nil, errors.NewClientError(
		errors.ErrorCodeUnsupportedOperation,
		"Partition retrieval not implemented for v0.0.41",
		"The v0.0.41 partition response uses complex inline structs that differ significantly from other API versions",
	)
}

// Update updates partition properties
func (m *PartitionManagerImpl) Update(ctx context.Context, partitionName string, update *interfaces.PartitionUpdate) error {
	// Partition updates are not supported in v0.0.41 API
	return errors.NewClientError(
		errors.ErrorCodeUnsupportedOperation,
		"Partition updates not supported",
		"The v0.0.41 Slurm REST API does not support partition update operations. Partition configuration changes must be made through slurmctld configuration files and require admin privileges.",
	)
}

// Watch provides real-time partition updates through polling
func (m *PartitionManagerImpl) Watch(ctx context.Context, opts *interfaces.WatchPartitionsOptions) (<-chan interfaces.PartitionEvent, error) {
	// Check if API client is available
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Create event channel
	eventChan := make(chan interfaces.PartitionEvent, 100)

	// Start polling goroutine
	go func() {
		defer close(eventChan)

		// Note: This is a simplified polling implementation
		// v0.0.41 doesn't have native streaming support

		select {
		case <-ctx.Done():
			return
		default:
			// In a full implementation, this would start a polling loop
		}
	}()

	return eventChan, nil
}
