// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_42

import (
	"context"
	"fmt"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	api "github.com/jontk/slurm-client/internal/api/v0_0_42"
)

// PartitionAdapter implements the PartitionAdapter interface for v0.0.42
type PartitionAdapter struct {
	*base.BaseManager
	client *api.ClientWithResponses
}

// NewPartitionAdapter creates a new Partition adapter for v0.0.42
func NewPartitionAdapter(client *api.ClientWithResponses) *PartitionAdapter {
	return &PartitionAdapter{
		BaseManager: base.NewBaseManager("v0.0.42", "Partition"),
		client:      client,
	}
}

// List retrieves a list of partitions
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
	params := &api.SlurmV0042GetPartitionsParams{}

	// Set flags to get detailed partition information
	flags := api.SlurmV0042GetPartitionsParamsFlagsDETAIL
	params.Flags = &flags

	// Apply filters from options
	if opts != nil && len(opts.Names) > 0 {
		// v0.0.42 doesn't support partition name filtering in the API params,
		// we'll need to filter client-side
	}

	// Call the API
	resp, err := a.client.SlurmV0042GetPartitionsWithResponse(ctx, params)
	if err != nil {
		return nil, a.WrapError(err, "failed to list partitions")
	}

	// Handle response - use HandleHTTPResponse instead
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return nil, err
	}

	// Check for API response
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response from API")
	}

	// Convert the response to common types
	partitionList := &types.PartitionList{
		Partitions: make([]types.Partition, 0),
	}

	if resp.JSON200.Partitions != nil {
		for _, apiPartition := range resp.JSON200.Partitions {
			partition, err := a.convertAPIPartitionToCommon(apiPartition)
			if err != nil {
				// Log conversion error but continue
				continue
			}
			
			// Apply client-side filtering if needed
			if opts != nil && len(opts.Names) > 0 {
				found := false
				for _, name := range opts.Names {
					if partition.Name == name {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}
			
			partitionList.Partitions = append(partitionList.Partitions, *partition)
		}
	}

	return partitionList, nil
}

// Get retrieves a specific partition by name
func (a *PartitionAdapter) Get(ctx context.Context, name string) (*types.Partition, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Prepare parameters
	params := &api.SlurmV0042GetPartitionParams{}
	flags := api.SlurmV0042GetPartitionParamsFlagsDETAIL
	params.Flags = &flags

	// Call the API
	resp, err := a.client.SlurmV0042GetPartitionWithResponse(ctx, name, params)
	if err != nil {
		return nil, a.WrapError(err, fmt.Sprintf("failed to get partition %s", name))
	}

	// Handle response - use HandleHTTPResponse instead
	if err := a.HandleHTTPResponse(resp.HTTPResponse, resp.Body); err != nil {
		return nil, err
	}

	// Check for API response
	if resp.JSON200 == nil || resp.JSON200.Partitions == nil || len(resp.JSON200.Partitions) == 0 {
		return nil, fmt.Errorf("partition %s not found", name)
	}

	// Convert the first partition in the response
	partitions := resp.JSON200.Partitions
	return a.convertAPIPartitionToCommon(partitions[0])
}

// Create creates a new partition
func (a *PartitionAdapter) Create(ctx context.Context, partition *types.PartitionCreate) (*types.PartitionCreateResponse, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// v0.0.42 doesn't have a direct partition create endpoint
	// This would typically be done through slurmctld configuration
	return nil, fmt.Errorf("partition creation not supported via v0.0.42 API")
}

// Update updates an existing partition
func (a *PartitionAdapter) Update(ctx context.Context, name string, updates *types.PartitionUpdateRequest) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// v0.0.42 doesn't have a partition update endpoint
	// This would typically be done through slurmctld reconfiguration
	return fmt.Errorf("partition update not supported via v0.0.42 API")
}

// Delete deletes a partition
func (a *PartitionAdapter) Delete(ctx context.Context, name string) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// v0.0.42 doesn't have a partition delete endpoint
	// This would typically be done through slurmctld configuration
	return fmt.Errorf("partition deletion not supported via v0.0.42 API")
}

// convertAPIPartitionToCommon converts API partition to common type
func (a *PartitionAdapter) convertAPIPartitionToCommon(apiPartition api.V0042PartitionInfo) (*types.Partition, error) {
	partition := &types.Partition{}

	// Set basic fields
	if apiPartition.Name != nil {
		partition.Name = *apiPartition.Name
	}

	// Convert state - check if State field exists and has a State field inside
	if apiPartition.State != nil {
		states := make([]types.PartitionState, len(*apiPartition.State))
		for i, state := range *apiPartition.State {
			states[i] = types.PartitionState(state)
		}
		partition.States = states
	}

	// Convert node information
	if apiPartition.Nodes != nil {
		if apiPartition.Nodes.AllowedAllocation != nil {
			partition.AllowAllocNodes = *apiPartition.Nodes.AllowedAllocation
		}
		if apiPartition.Nodes.Config != nil {
			partition.Nodes = *apiPartition.Nodes.Config
		}
	}

	// Convert resource limits from Maximums
	if apiPartition.Maximums != nil {
		if apiPartition.Maximums.CpusPerNode != nil && apiPartition.Maximums.CpusPerNode.Set {
			partition.MaxCPUsPerNode = apiPartition.Maximums.CpusPerNode.Number
		}
		if apiPartition.Maximums.MemoryPerCpu != nil {
			partition.MaxMemPerCPU = *apiPartition.Maximums.MemoryPerCpu
		}
		if apiPartition.Maximums.Nodes != nil && apiPartition.Maximums.Nodes.Set {
			partition.MaxNodes = apiPartition.Maximums.Nodes.Number
		}
	}

	// Convert defaults
	if apiPartition.Defaults != nil {
		if apiPartition.Defaults.MemoryPerCpu != nil {
			partition.DefaultMemPerCPU = *apiPartition.Defaults.MemoryPerCpu
		}
		if apiPartition.Defaults.Time != nil && apiPartition.Defaults.Time.Set {
			partition.DefaultTime = apiPartition.Defaults.Time.Number
		}
	}

	// Convert flags
	if apiPartition.Flags != nil {
		for _, flag := range *apiPartition.Flags {
			if string(flag) == "DEFAULT" {
				partition.Default = true
				break
			}
		}
	}

	// Convert priority
	if apiPartition.Priority != nil && apiPartition.Priority.JobFactor != nil && apiPartition.Priority.JobFactor.Set {
		partition.Priority = apiPartition.Priority.JobFactor.Number
	}

	// Convert CPU total
	if apiPartition.Cpus != nil && apiPartition.Cpus.Total != nil {
		partition.TotalCPUs = *apiPartition.Cpus.Total
	}

	return partition, nil
}
