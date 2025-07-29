package v0_0_40

import (
	"context"
	"strings"

	"github.com/jontk/slurm-client/internal/common"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	api "github.com/jontk/slurm-client/internal/api/v0_0_40"
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
		if len(opts.Names) > 0 {
			nameStr := strings.Join(opts.Names, ",")
			params.PartitionName = &nameStr
		}
		if opts.UpdateTime != nil {
			updateTime := opts.UpdateTime.Unix()
			params.UpdateTime = &updateTime
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
		partition, err := a.convertAPIPartitionToCommon(apiPartition)
		if err != nil {
			return nil, a.HandleConversionError(err, apiPartition.Name)
		}
		partitionList = append(partitionList, *partition)
	}

	// Apply client-side filtering if needed
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
	if err := a.ValidateResourceName(partitionName, "partitionName"); err != nil {
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
	partition, err := a.convertAPIPartitionToCommon(resp.JSON200.Partitions[0])
	if err != nil {
		return nil, a.HandleConversionError(err, partitionName)
	}

	return partition, nil
}

// Create creates a new partition
func (a *PartitionAdapter) Create(ctx context.Context, partition *types.PartitionCreate) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.validatePartitionCreate(partition); err != nil {
		return err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Convert to API format
	apiPartition, err := a.convertCommonPartitionCreateToAPI(partition)
	if err != nil {
		return err
	}

	// Create request body
	reqBody := api.SlurmV0040PostPartitionJSONRequestBody{
		Partitions: &[]api.V0040PartitionInfo{*apiPartition},
	}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmV0040PostPartitionWithResponse(ctx, reqBody)
	if err != nil {
		return a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0040OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	return common.HandleAPIResponse(responseAdapter, "v0.0.40")
}

// Update updates an existing partition
func (a *PartitionAdapter) Update(ctx context.Context, partitionName string, update *types.PartitionUpdate) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.ValidateResourceName(partitionName, "partitionName"); err != nil {
		return err
	}
	if err := a.validatePartitionUpdate(update); err != nil {
		return err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// First, get the existing partition to merge updates
	existingPartition, err := a.Get(ctx, partitionName)
	if err != nil {
		return err
	}

	// Convert to API format and apply updates
	apiPartition, err := a.convertCommonPartitionUpdateToAPI(existingPartition, update)
	if err != nil {
		return err
	}

	// Create request body
	reqBody := api.SlurmV0040PostPartitionJSONRequestBody{
		Partitions: &[]api.V0040PartitionInfo{*apiPartition},
	}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmV0040PostPartitionWithResponse(ctx, reqBody)
	if err != nil {
		return a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0040OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	return common.HandleAPIResponse(responseAdapter, "v0.0.40")
}

// Delete deletes a partition
func (a *PartitionAdapter) Delete(ctx context.Context, partitionName string) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.ValidateResourceName(partitionName, "partitionName"); err != nil {
		return err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmV0040DeletePartitionWithResponse(ctx, partitionName)
	if err != nil {
		return a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0040OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	// Create adapter with special handling for 204 (No Content) status
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)

	// For DELETE operations, 204 is also a success
	if resp.StatusCode() == 204 {
		return nil
	}

	return common.HandleAPIResponse(responseAdapter, "v0.0.40")
}

// filterPartitionList applies client-side filtering to the partition list
func (a *PartitionAdapter) filterPartitionList(partitions []types.Partition, opts *types.PartitionListOptions) []types.Partition {
	filtered := make([]types.Partition, 0, len(partitions))
	
	for _, partition := range partitions {
		// Apply State filter
		if len(opts.States) > 0 {
			found := false
			for _, state := range opts.States {
				if partition.State == state {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Apply Default filter
		if opts.Default != nil {
			if partition.Default != *opts.Default {
				continue
			}
		}

		// Apply AllowGroups filter
		if len(opts.AllowGroups) > 0 {
			found := false
			for _, group := range opts.AllowGroups {
				for _, partGroup := range partition.AllowGroups {
					if group == partGroup {
						found = true
						break
					}
				}
				if found {
					break
				}
			}
			if !found {
				continue
			}
		}

		// Apply Flags filter
		if len(opts.Flags) > 0 {
			hasAllFlags := true
			for _, flag := range opts.Flags {
				found := false
				for _, partFlag := range partition.Flags {
					if flag == partFlag {
						found = true
						break
					}
				}
				if !found {
					hasAllFlags = false
					break
				}
			}
			if !hasAllFlags {
				continue
			}
		}

		// Apply QoS filter
		if opts.QoS != "" && partition.QoS != opts.QoS {
			continue
		}

		filtered = append(filtered, partition)
	}

	return filtered
}

// validatePartitionCreate validates partition creation request
func (a *PartitionAdapter) validatePartitionCreate(partition *types.PartitionCreate) error {
	if partition == nil {
		return common.NewValidationError("partition creation data is required", "partition", nil)
	}
	if partition.Name == "" {
		return common.NewValidationError("partition name is required", "name", partition.Name)
	}
	if len(partition.Nodes) == 0 {
		return common.NewValidationError("at least one node is required", "nodes", partition.Nodes)
	}
	return nil
}

// validatePartitionUpdate validates partition update request
func (a *PartitionAdapter) validatePartitionUpdate(update *types.PartitionUpdate) error {
	if update == nil {
		return common.NewValidationError("partition update data is required", "update", nil)
	}
	// At least one field should be provided for update
	if update.State == nil && update.MaxTime == nil && update.DefaultTime == nil && 
	   update.MaxNodes == nil && update.MinNodes == nil && update.Default == nil && 
	   update.Nodes == nil && update.AllowGroups == nil && update.AllowAccounts == nil && 
	   update.DenyGroups == nil && update.DenyAccounts == nil && update.Flags == nil && 
	   update.QoS == nil {
		return common.NewValidationError("at least one field must be provided for update", "update", update)
	}
	return nil
}