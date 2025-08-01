// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_41

import (
	"context"
	"fmt"
	"strings"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	api "github.com/jontk/slurm-client/internal/api/v0_0_41"
)

// PartitionAdapter implements the PartitionAdapter interface for v0.0.41
type PartitionAdapter struct {
	*base.BaseManager
	client  *api.ClientWithResponses
	wrapper *api.WrapperClient
}

// NewPartitionAdapter creates a new Partition adapter for v0.0.41
func NewPartitionAdapter(client *api.ClientWithResponses) *PartitionAdapter {
	return &PartitionAdapter{
		BaseManager: base.NewBaseManager("v0.0.41", "Partition"),
		client:      client,
		wrapper:     nil, // We'll implement this later
	}
}

// List retrieves a list of partitions with optional filtering
func (a *PartitionAdapter) List(ctx context.Context, opts *types.PartitionListOptions) (*types.PartitionList, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Prepare parameters for the API call
	params := &api.SlurmV0041GetPartitionsParams{}

	// Apply filters from options
	// Note: v0.0.41 API has limited filtering support, implement client-side filtering
	if opts != nil {
		// v0.0.41 GetPartitions doesn't support filtering parameters like PartitionName
		// We'll filter results after fetching
		_ = opts
	}

	// Make the API call
	resp, err := a.client.SlurmV0041GetPartitionsWithResponse(ctx, params)
	if err != nil {
		return nil, a.WrapError(err, "failed to list partitions")
	}

	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return nil, err
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected nil response")
	}

	// Convert response to common types
	partitionList := &types.PartitionList{
		Partitions: make([]types.Partition, 0, len(resp.JSON200.Partitions)),
		Total:      len(resp.JSON200.Partitions),
	}

	for _, apiPartition := range resp.JSON200.Partitions {
		partition, err := a.convertAPIPartitionToCommon(apiPartition)
		if err != nil {
			// Log the error but continue processing other partitions
			continue
		}
		partitionList.Partitions = append(partitionList.Partitions, *partition)
	}

	// Update total count after filtering
	partitionList.Total = len(partitionList.Partitions)

	// Extract warning and error messages if any (but PartitionList doesn't have Meta)
	// Warnings and errors are ignored for now as PartitionList structure doesn't support them
	if resp.JSON200.Warnings != nil {
		// Log warnings if needed
		_ = resp.JSON200.Warnings
	}
	if resp.JSON200.Errors != nil {
		// Log errors if needed  
		_ = resp.JSON200.Errors
	}

	return partitionList, nil
}

// Get retrieves a specific partition by name
func (a *PartitionAdapter) Get(ctx context.Context, name string) (*types.Partition, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	// Validate name
	if err := a.ValidateResourceName("partition name", name); err != nil {
		return nil, err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Make the API call
	params := &api.SlurmV0041GetPartitionParams{}
	resp, err := a.client.SlurmV0041GetPartitionWithResponse(ctx, name, params)
	if err != nil {
		return nil, a.WrapError(err, fmt.Sprintf("failed to get partition %s", name))
	}

	// Handle response
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return nil, err
	}

	if resp.JSON200 == nil || len(resp.JSON200.Partitions) == 0 {
		return nil, a.HandleNotFound(fmt.Sprintf("partition %s", name))
	}

	// Convert the first partition in the response
	partition, err := a.convertAPIPartitionToCommon(resp.JSON200.Partitions[0])
	if err != nil {
		return nil, a.WrapError(err, fmt.Sprintf("failed to convert partition %s", name))
	}

	return partition, nil
}

// Create creates a new partition
func (a *PartitionAdapter) Create(ctx context.Context, req *types.PartitionCreate) (*types.PartitionCreateResponse, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	// Validate request
	if req == nil {
		return nil, a.HandleValidationError("partition create request cannot be nil")
	}
	if err := a.ValidateResourceName("partition name", req.Name); err != nil {
		return nil, err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// v0.0.41 doesn't support partition creation through the API
	return nil, fmt.Errorf("partition creation is not supported in API v0.0.41")
}

// Update updates an existing partition
func (a *PartitionAdapter) Update(ctx context.Context, name string, update *types.PartitionUpdate) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Validate name
	if err := a.ValidateResourceName("partition name", name); err != nil {
		return err
	}

	// Validate update
	if update == nil {
		return a.HandleValidationError("partition update cannot be nil")
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// v0.0.41 doesn't support partition updates through the API
	return fmt.Errorf("partition update is not supported in API v0.0.41")
}

// Delete deletes a partition
func (a *PartitionAdapter) Delete(ctx context.Context, name string) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Validate name
	if err := a.ValidateResourceName("partition name", name); err != nil {
		return err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// v0.0.41 doesn't support partition deletion through the API
	return fmt.Errorf("partition deletion is not supported in API v0.0.41")
}

// GetNodeList gets the list of nodes for a partition
func (a *PartitionAdapter) GetNodeList(ctx context.Context, name string) ([]string, error) {
	// Get the partition first
	partition, err := a.Get(ctx, name)
	if err != nil {
		return nil, err
	}

	// Extract node names from the partition
	if partition.Nodes == "" {
		return []string{}, nil
	}

	// Parse the node list (Slurm format can be "node[001-005,007]")
	nodes := parseNodeList(partition.Nodes)
	return nodes, nil
}

// parseNodeList parses a Slurm node list string into individual node names
func parseNodeList(nodeStr string) []string {
	// Simple implementation - in production, this would need to handle
	// Slurm's node range notation properly
	nodes := strings.Split(nodeStr, ",")
	var result []string
	for _, node := range nodes {
		node = strings.TrimSpace(node)
		if node != "" {
			result = append(result, node)
		}
	}
	return result
}
