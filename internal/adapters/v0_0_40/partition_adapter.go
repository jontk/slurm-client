// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_40

import (
	"context"
	"strconv"
	"strings"

	api "github.com/jontk/slurm-client/internal/api/v0_0_40"
	"github.com/jontk/slurm-client/internal/common"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	"github.com/jontk/slurm-client/pkg/errors"
)

// PartitionAdapter implements the PartitionAdapter interface for v0.0.40
type PartitionAdapter struct {
	*base.BaseManager
	client  *api.ClientWithResponses
	wrapper *api.WrapperClient
}

// NewPartitionAdapter creates a new Partition adapter for v0.0.40
func NewPartitionAdapter(client *api.ClientWithResponses) *PartitionAdapter {
	return &PartitionAdapter{
		BaseManager: base.NewBaseManager("v0.0.40", "Partition"),
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
	params := &api.SlurmV0040GetPartitionsParams{}

	// Apply filters from options
	if opts != nil {
		// v0.0.40 doesn't support partition name filtering in params
		if opts.UpdateTime != nil {
			updateTimeStr := strconv.FormatInt(opts.UpdateTime.Unix(), 10)
			params.UpdateTime = &updateTimeStr
		}
	}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmV0040GetPartitionsWithResponse(ctx, params)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0040OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.40"); err != nil {
		return nil, err
	}

	// Check for unexpected response format
	if err := a.CheckNilResponse(resp.JSON200, "List Partitions"); err != nil {
		return nil, err
	}
	if err := a.CheckNilResponse(resp.JSON200.Partitions, "List Partitions - partitions field"); err != nil {
		return nil, err
	}

	// Convert the response to common types
	partitionList := make([]types.Partition, 0, len(resp.JSON200.Partitions))
	for _, apiPartition := range resp.JSON200.Partitions {
		partition := a.convertAPIPartitionToCommon(apiPartition)
		partitionList = append(partitionList, *partition)
	}

	// Apply client-side filtering if needed (since API doesn't support all filters)
	if opts != nil {
		partitionList = a.filterPartitionList(partitionList, opts)
	}

	// Apply pagination
	listOpts := base.ListOptions{}
	if opts != nil {
		listOpts.Limit = opts.Limit
		listOpts.Offset = opts.Offset
	}

	// Apply pagination
	start := listOpts.Offset
	if start < 0 {
		start = 0
	}
	if start >= len(partitionList) {
		return &types.PartitionList{
			Partitions: []types.Partition{},
			Total:      len(partitionList),
		}, nil
	}

	end := len(partitionList)
	if listOpts.Limit > 0 {
		end = start + listOpts.Limit
		if end > len(partitionList) {
			end = len(partitionList)
		}
	}

	return &types.PartitionList{
		Partitions: partitionList[start:end],
		Total:      len(partitionList),
	}, nil
}

// Get retrieves a specific partition by name
func (a *PartitionAdapter) Get(ctx context.Context, partitionName string) (*types.Partition, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if err := a.ValidateResourceName(partitionName, "partition name"); err != nil {
		return nil, err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Prepare parameters for the API call
	params := &api.SlurmV0040GetPartitionParams{}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmV0040GetPartitionWithResponse(ctx, partitionName, params)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0040OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.40"); err != nil {
		return nil, err
	}

	// Check for unexpected response format
	if err := a.CheckNilResponse(resp.JSON200, "Get Partition"); err != nil {
		return nil, err
	}
	if err := a.CheckNilResponse(resp.JSON200.Partitions, "Get Partition - partitions field"); err != nil {
		return nil, err
	}

	// Check if we got any partition entries
	if len(resp.JSON200.Partitions) == 0 {
		return nil, common.NewResourceNotFoundError("Partition", partitionName)
	}

	// Convert the first partition (should be the only one)
	partition := a.convertAPIPartitionToCommon(resp.JSON200.Partitions[0])

	return partition, nil
}

// Create creates a new partition
func (a *PartitionAdapter) Create(ctx context.Context, partition *types.PartitionCreate) (*types.PartitionCreateResponse, error) {
	// v0.0.40 doesn't support partition creation
	return nil, errors.NewNotImplementedError("partition creation", "v0.0.40")
}

// Update updates an existing partition
func (a *PartitionAdapter) Update(ctx context.Context, partitionName string, update *types.PartitionUpdate) error {
	// v0.0.40 doesn't support partition updates
	return errors.NewNotImplementedError("partition updates", "v0.0.40")
}

// Delete deletes a partition
func (a *PartitionAdapter) Delete(ctx context.Context, partitionName string) error {
	// v0.0.40 doesn't support partition deletion
	return errors.NewNotImplementedError("partition deletion", "v0.0.40")
}

// convertAPIPartitionToCommon converts a v0.0.40 API Partition to common Partition type
func (a *PartitionAdapter) convertAPIPartitionToCommon(apiPartition api.V0040PartitionInfo) *types.Partition {
	partition := &types.Partition{}

	// Basic fields
	if apiPartition.Name != nil {
		partition.Name = *apiPartition.Name
	}
	if apiPartition.Nodes != nil && apiPartition.Nodes.Configured != nil {
		partition.Nodes = *apiPartition.Nodes.Configured
	}

	// Limits and timeouts - v0.0.40 structure is different
	// Just set basic fields for now

	return partition
}

// filterPartitionList applies client-side filtering to partition list
func (a *PartitionAdapter) filterPartitionList(partitions []types.Partition, opts *types.PartitionListOptions) []types.Partition {
	if opts == nil {
		return partitions
	}

	filtered := make([]types.Partition, 0, len(partitions))
	for _, partition := range partitions {
		if a.matchesPartitionFilters(partition, opts) {
			filtered = append(filtered, partition)
		}
	}

	return filtered
}

// matchesPartitionFilters checks if a partition matches the given filters
func (a *PartitionAdapter) matchesPartitionFilters(partition types.Partition, opts *types.PartitionListOptions) bool {
	// Filter by names (client-side since API doesn't support it)
	if len(opts.Names) > 0 {
		found := false
		for _, name := range opts.Names {
			if strings.EqualFold(partition.Name, name) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Filter by states
	if len(opts.States) > 0 {
		found := false
		for _, state := range opts.States {
			if partition.State == state {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}
