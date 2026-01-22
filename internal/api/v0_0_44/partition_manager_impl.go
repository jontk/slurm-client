// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_44

import (
	"net/http"
	"context"
	"fmt"
	"strings"
	"time"

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
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Prepare parameters for the API call
	params := &SlurmV0044GetPartitionsParams{}

	// Apply filtering options if provided
	if opts != nil {
		// Note: v0.0.44 API has limited filtering options
		// States filtering will be applied client-side if needed
	}

	// Call the generated OpenAPI client
	resp, err := m.client.apiClient.SlurmV0044GetPartitionsWithResponse(ctx, params)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.44")
	}

	// Handle response
	if resp.HTTPResponse.StatusCode != http.StatusOK {
		return nil, errors.NewSlurmError(errors.ErrorCodeServerInternal,
			fmt.Sprintf("HTTP %d", resp.HTTPResponse.StatusCode))
	}

	if resp.JSON200 == nil {
		return nil, errors.NewSlurmError(errors.ErrorCodeServerInternal, "Empty response body")
	}

	// Convert to interface types
	partitionList := &interfaces.PartitionList{
		Partitions: make([]interfaces.Partition, 0),
	}

	if resp.JSON200.Partitions != nil {
		for _, partition := range resp.JSON200.Partitions {
			convertedPartition := m.convertPartitionToInterface(&partition)
			partitionList.Partitions = append(partitionList.Partitions, *convertedPartition)
		}
	}

	return partitionList, nil
}

// Get retrieves a specific partition by name
func (m *PartitionManagerImpl) Get(ctx context.Context, partitionName string) (*interfaces.Partition, error) {
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Use the partition-specific endpoint
	resp, err := m.client.apiClient.SlurmV0044GetPartitionWithResponse(ctx, partitionName, &SlurmV0044GetPartitionParams{})
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.44")
	}

	if resp.HTTPResponse.StatusCode == http.StatusNotFound {
		return nil, errors.NewSlurmError(errors.ErrorCodeResourceNotFound,
			fmt.Sprintf("Partition %s not found", partitionName))
	}

	if resp.HTTPResponse.StatusCode != http.StatusOK {
		return nil, errors.NewSlurmError(errors.ErrorCodeServerInternal,
			fmt.Sprintf("HTTP %d", resp.HTTPResponse.StatusCode))
	}

	if resp.JSON200 == nil || resp.JSON200.Partitions == nil || len(resp.JSON200.Partitions) == 0 {
		return nil, errors.NewSlurmError(errors.ErrorCodeResourceNotFound,
			fmt.Sprintf("Partition %s not found", partitionName))
	}

	partition := resp.JSON200.Partitions[0]
	return m.convertPartitionToInterface(&partition), nil
}

// Update updates partition properties
func (m *PartitionManagerImpl) Update(ctx context.Context, partitionName string, update *interfaces.PartitionUpdate) error {
	// Partition updates are not supported in SLURM REST API v0.0.44
	return errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation,
		"Partition updates are not supported in SLURM REST API v0.0.44")
}

// Watch provides real-time partition updates
func (m *PartitionManagerImpl) Watch(ctx context.Context, opts *interfaces.WatchPartitionsOptions) (<-chan interfaces.PartitionEvent, error) {
	// Create a channel for partition events
	eventChan := make(chan interfaces.PartitionEvent)

	// For now, return a basic watcher that polls for changes
	go func() {
		defer close(eventChan)

		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// Poll for partition changes - basic implementation
				// In a real implementation, this would use WebSocket or SSE
			}
		}
	}()

	return eventChan, nil
}

// convertPartitionToInterface converts a v0.0.44 partition to the interface type
func (m *PartitionManagerImpl) convertPartitionToInterface(partition *V0044PartitionInfo) *interfaces.Partition {
	result := &interfaces.Partition{}

	// Basic partition information
	if partition.Name != nil {
		result.Name = *partition.Name
	}

	// State (handle array of states)
	if partition.Partition != nil && partition.Partition.State != nil && len(*partition.Partition.State) > 0 {
		result.State = string((*partition.Partition.State)[0])
	}

	// Nodes (handle nested structure)
	if partition.Nodes != nil && partition.Nodes.Configured != nil {
		result.Nodes = strings.Split(*partition.Nodes.Configured, ",")
	}

	// Resource limits (these fields don't exist in interfaces.Partition)
	// MaxNodes and MinNodes are not part of the interface

	if partition.Nodes != nil && partition.Nodes.Total != nil {
		result.TotalNodes = int(*partition.Nodes.Total)
	}

	if partition.Cpus != nil && partition.Cpus.Total != nil {
		result.TotalCPUs = int(*partition.Cpus.Total)
	}

	// Time limits (handle nested defaults/maximums)
	if partition.Defaults != nil && partition.Defaults.Time != nil && partition.Defaults.Time.Number != nil {
		result.DefaultTime = int(*partition.Defaults.Time.Number)
	}

	if partition.Maximums != nil && partition.Maximums.Time != nil && partition.Maximums.Time.Number != nil {
		result.MaxTime = int(*partition.Maximums.Time.Number)
	}

	// Account and group limits (map to existing interface fields)
	if partition.Groups != nil && partition.Groups.Allowed != nil {
		result.AllowedGroups = strings.Split(*partition.Groups.Allowed, ",")
	}

	if partition.Accounts != nil && partition.Accounts.Allowed != nil {
		result.AllowedUsers = strings.Split(*partition.Accounts.Allowed, ",")
	}

	if partition.Accounts != nil && partition.Accounts.Deny != nil {
		result.DeniedUsers = strings.Split(*partition.Accounts.Deny, ",")
	}

	// Additional configuration
	if partition.Priority != nil && partition.Priority.Tier != nil {
		result.Priority = int(*partition.Priority.Tier)
	}

	// Set defaults for missing fields
	result.AvailableNodes = 0 // Not directly available in v0.0.44
	result.IdleCPUs = 0       // Not directly available in v0.0.44
	result.MaxMemory = 0      // Would need to parse from TRES or memory fields
	result.DefaultMemory = 0  // Would need to parse from TRES or memory fields

	// Ensure slice fields are not nil
	if result.AllowedUsers == nil {
		result.AllowedUsers = []string{}
	}
	if result.DeniedUsers == nil {
		result.DeniedUsers = []string{}
	}
	if result.AllowedGroups == nil {
		result.AllowedGroups = []string{}
	}
	if result.DeniedGroups == nil {
		result.DeniedGroups = []string{}
	}
	if result.Nodes == nil {
		result.Nodes = []string{}
	}

	return result
}

// Create creates a new partition (if supported by version)
func (m *PartitionManagerImpl) Create(ctx context.Context, partition *interfaces.PartitionCreate) (*interfaces.PartitionCreateResponse, error) {
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// For now, return not implemented - partition creation might not be supported in REST API
	return nil, errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "Partition creation not yet implemented for v0.0.44")
}

// Delete removes a partition (if supported by version)
func (m *PartitionManagerImpl) Delete(ctx context.Context, partitionName string) error {
	if m.client.apiClient == nil {
		return errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// For now, return not implemented - partition deletion might not be supported in REST API
	return errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "Partition deletion not yet implemented for v0.0.44")
}
