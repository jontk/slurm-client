package v0_0_42

import (
	"context"
	"fmt"
	"strings"

	"github.com/jontk/slurm-client/internal/interfaces"
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
		return nil, fmt.Errorf("API client not initialized")
	}
	
	// Prepare parameters for the API call
	params := &SlurmV0042GetPartitionsParams{}
	
	// Set flags to get detailed partition information
	flags := SlurmV0042GetPartitionsParamsFlagsDETAIL
	params.Flags = &flags
	
	// Call the generated OpenAPI client
	resp, err := m.client.apiClient.SlurmV0042GetPartitionsWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list partitions: %w", err)
	}
	
	// Check HTTP status
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode(), resp.Status())
	}
	
	// Check for API errors
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected response format")
	}
	
	if resp.JSON200.Errors != nil && len(*resp.JSON200.Errors) > 0 {
		return nil, fmt.Errorf("API error: %v", (*resp.JSON200.Errors)[0])
	}
	
	// Convert the response to our interface types
	partitions := make([]interfaces.Partition, 0, len(resp.JSON200.Partitions))
	for _, apiPartition := range resp.JSON200.Partitions {
		partition, err := convertAPIPartitionToInterface(apiPartition)
		if err != nil {
			return nil, fmt.Errorf("failed to convert partition data: %w", err)
		}
		partitions = append(partitions, *partition)
	}
	
	// Apply client-side filtering if options are provided
	if opts != nil {
		partitions = filterPartitions(partitions, opts)
	}
	
	return &interfaces.PartitionList{
		Partitions: partitions,
		Total:      len(partitions),
	}, nil
}

// convertAPIPartitionToInterface converts a V0042PartitionInfo to interfaces.Partition
func convertAPIPartitionToInterface(apiPartition V0042PartitionInfo) (*interfaces.Partition, error) {
	partition := &interfaces.Partition{}
	
	// Partition name - simple string
	if apiPartition.Name != nil {
		partition.Name = *apiPartition.Name
	}
	
	// Partition state - nested under Partition.State
	if apiPartition.Partition != nil && apiPartition.Partition.State != nil && len(*apiPartition.Partition.State) > 0 {
		partition.State = (*apiPartition.Partition.State)[0]
	}
	
	// Node counts - nested under Nodes
	if apiPartition.Nodes != nil && apiPartition.Nodes.Total != nil {
		partition.TotalNodes = int(*apiPartition.Nodes.Total)
	}
	
	// Available nodes calculation - we don't have allocated nodes directly
	// For now, we'll assume all nodes are available (this could be enhanced)
	partition.AvailableNodes = partition.TotalNodes
	
	// CPU counts - nested under Cpus.Total
	if apiPartition.Cpus != nil && apiPartition.Cpus.Total != nil {
		partition.TotalCPUs = int(*apiPartition.Cpus.Total)
	}
	
	// Idle CPUs - without allocated CPUs data, assume all are idle
	// This is a simplification; real implementation might query job data
	partition.IdleCPUs = partition.TotalCPUs
	
	// Time limits (in minutes) - nested under Maximums.Time and Defaults.Time
	if apiPartition.Maximums != nil && apiPartition.Maximums.Time != nil && 
		apiPartition.Maximums.Time.Set != nil && *apiPartition.Maximums.Time.Set && 
		apiPartition.Maximums.Time.Number != nil {
		partition.MaxTime = int(*apiPartition.Maximums.Time.Number)
	}
	
	if apiPartition.Defaults != nil && apiPartition.Defaults.Time != nil && 
		apiPartition.Defaults.Time.Set != nil && *apiPartition.Defaults.Time.Set && 
		apiPartition.Defaults.Time.Number != nil {
		partition.DefaultTime = int(*apiPartition.Defaults.Time.Number)
	}
	
	// Memory limits (convert MB to bytes for consistency) - nested under Maximums
	if apiPartition.Maximums != nil && apiPartition.Maximums.MemoryPerCpu != nil {
		partition.MaxMemory = int(*apiPartition.Maximums.MemoryPerCpu) * 1024 * 1024
	}
	
	// Default memory - nested under Defaults
	if apiPartition.Defaults != nil && apiPartition.Defaults.MemoryPerCpu != nil {
		partition.DefaultMemory = int(*apiPartition.Defaults.MemoryPerCpu) * 1024 * 1024
	}
	
	// User and group access lists - nested under Accounts and Groups
	if apiPartition.Accounts != nil && apiPartition.Accounts.Allowed != nil {
		// Parse comma-separated string into slice
		if *apiPartition.Accounts.Allowed != "" {
			partition.AllowedUsers = strings.Split(*apiPartition.Accounts.Allowed, ",")
		} else {
			partition.AllowedUsers = []string{}
		}
	} else {
		partition.AllowedUsers = []string{}
	}
	
	if apiPartition.Accounts != nil && apiPartition.Accounts.Deny != nil {
		// Parse comma-separated string into slice
		if *apiPartition.Accounts.Deny != "" {
			partition.DeniedUsers = strings.Split(*apiPartition.Accounts.Deny, ",")
		} else {
			partition.DeniedUsers = []string{}
		}
	} else {
		partition.DeniedUsers = []string{}
	}
	
	if apiPartition.Groups != nil && apiPartition.Groups.Allowed != nil {
		// Parse comma-separated string into slice
		if *apiPartition.Groups.Allowed != "" {
			partition.AllowedGroups = strings.Split(*apiPartition.Groups.Allowed, ",")
		} else {
			partition.AllowedGroups = []string{}
		}
	} else {
		partition.AllowedGroups = []string{}
	}
	
	// No denied groups field in V0042PartitionInfo
	partition.DeniedGroups = []string{}
	
	// Priority - nested under Priority.Tier
	if apiPartition.Priority != nil && apiPartition.Priority.Tier != nil {
		partition.Priority = int(*apiPartition.Priority.Tier)
	}
	
	// Node list - nested under Nodes.Configured
	if apiPartition.Nodes != nil && apiPartition.Nodes.Configured != nil {
		// Parse node list string into slice (simplified - real parsing might be more complex)
		if *apiPartition.Nodes.Configured != "" {
			partition.Nodes = strings.Split(*apiPartition.Nodes.Configured, ",")
		} else {
			partition.Nodes = []string{}
		}
	} else {
		partition.Nodes = []string{}
	}
	
	return partition, nil
}

// filterPartitions applies client-side filtering to partition list
func filterPartitions(partitions []interfaces.Partition, opts *interfaces.ListPartitionsOptions) []interfaces.Partition {
	var filtered []interfaces.Partition
	
	// If no options provided, return all partitions
	if opts == nil {
		return partitions
	}
	
	for _, partition := range partitions {
		// Filter by states
		if len(opts.States) > 0 {
			stateMatch := false
			for _, state := range opts.States {
				if strings.EqualFold(partition.State, state) {
					stateMatch = true
					break
				}
			}
			if !stateMatch {
				continue
			}
		}
		
		filtered = append(filtered, partition)
	}
	
	// Apply limit and offset
	if opts.Offset > 0 {
		if opts.Offset >= len(filtered) {
			return []interfaces.Partition{}
		}
		filtered = filtered[opts.Offset:]
	}
	
	if opts.Limit > 0 && len(filtered) > opts.Limit {
		filtered = filtered[:opts.Limit]
	}
	
	return filtered
}

// Get retrieves a specific partition by name
func (m *PartitionManagerImpl) Get(ctx context.Context, partitionName string) (*interfaces.Partition, error) {
	// Check if API client is available
	if m.client.apiClient == nil {
		return nil, fmt.Errorf("API client not initialized")
	}
	
	// Prepare parameters for the API call
	params := &SlurmV0042GetPartitionParams{}
	
	// Set flags to get detailed partition information
	flags := SlurmV0042GetPartitionParamsFlagsDETAIL
	params.Flags = &flags
	
	// Call the generated OpenAPI client
	resp, err := m.client.apiClient.SlurmV0042GetPartitionWithResponse(ctx, partitionName, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get partition %s: %w", partitionName, err)
	}
	
	// Check HTTP status
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("API returned status %d for partition %s: %s", resp.StatusCode(), partitionName, resp.Status())
	}
	
	// Check for API errors
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected response format for partition %s", partitionName)
	}
	
	if resp.JSON200.Errors != nil && len(*resp.JSON200.Errors) > 0 {
		return nil, fmt.Errorf("API error for partition %s: %v", partitionName, (*resp.JSON200.Errors)[0])
	}
	
	// Convert the response to our interface types
	if len(resp.JSON200.Partitions) == 0 {
		return nil, fmt.Errorf("partition %s not found", partitionName)
	}
	
	if len(resp.JSON200.Partitions) > 1 {
		return nil, fmt.Errorf("unexpected: multiple partitions returned for name %s", partitionName)
	}
	
	partition, err := convertAPIPartitionToInterface(resp.JSON200.Partitions[0])
	if err != nil {
		return nil, fmt.Errorf("failed to convert partition %s data: %w", partitionName, err)
	}
	
	return partition, nil
}