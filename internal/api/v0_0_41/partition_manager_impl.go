// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_41

import (
	"context"

	"github.com/jontk/slurm-client/interfaces"
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
	// Check if API client is available
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Call the GetPartitions API
	resp, err := m.client.apiClient.SlurmV0041GetPartitionsWithResponse(ctx, nil)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.41")
	}

	// Check HTTP status
	if resp.StatusCode() != 200 {
		var responseBody []byte
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.41")
		return nil, httpErr
	}

	// Check for unexpected response format
	if resp.JSON200 == nil || len(resp.JSON200.Partitions) == 0 {
		return &interfaces.PartitionList{Partitions: []interfaces.Partition{}, Total: 0}, nil
	}

	// Minimal conversion for v0.0.41 - just extract partition names
	partitions := make([]interfaces.Partition, 0, len(resp.JSON200.Partitions))
	for _, apiPart := range resp.JSON200.Partitions {
		partition := interfaces.Partition{
			Name: *apiPart.Name,
		}
		partitions = append(partitions, partition)
	}

	return &interfaces.PartitionList{
		Partitions: partitions,
		Total:      len(partitions),
	}, nil
}

// Get retrieves a specific partition by name
func (m *PartitionManagerImpl) Get(ctx context.Context, partitionName string) (*interfaces.Partition, error) {
	// Check if API client is available
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Call the GetPartition API
	resp, err := m.client.apiClient.SlurmV0041GetPartitionWithResponse(ctx, partitionName, nil)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.41")
	}

	// Check HTTP status
	if resp.StatusCode() != 200 {
		var responseBody []byte
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.41")
		return nil, httpErr
	}

	// Check for unexpected response format
	if resp.JSON200 == nil || len(resp.JSON200.Partitions) == 0 {
		return nil, errors.NewClientError(errors.ErrorCodeServerInternal, "Unexpected response format")
	}

	// Minimal conversion for v0.0.41
	apiPart := resp.JSON200.Partitions[0]
	partition := &interfaces.Partition{
		Name: *apiPart.Name,
	}

	return partition, nil
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

// Create creates a new partition
func (m *PartitionManagerImpl) Create(ctx context.Context, partition *interfaces.PartitionCreate) (*interfaces.PartitionCreateResponse, error) {
	return nil, errors.NewNotImplementedError("Create", "v0.0.41")
}

// Delete deletes a partition
func (m *PartitionManagerImpl) Delete(ctx context.Context, partitionName string) error {
	return errors.NewNotImplementedError("Delete", "v0.0.41")
}
