package v0_0_43

import (
	"context"
	"strings"

	"github.com/jontk/slurm-client/internal/common"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	api "github.com/jontk/slurm-client/internal/api/v0_0_43"
)

// PartitionAdapter implements the PartitionAdapter interface for v0.0.43
type PartitionAdapter struct {
	*base.BaseManager
	client  *api.ClientWithResponses
	wrapper *api.WrapperClient
}

// NewPartitionAdapter creates a new Partition adapter for v0.0.43
func NewPartitionAdapter(client *api.ClientWithResponses) *PartitionAdapter {
	return &PartitionAdapter{
		BaseManager: base.NewBaseManager("v0.0.43", "Partition"),
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
	params := &api.SlurmV0043GetPartitionsParams{}

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
	resp, err := a.client.SlurmV0043GetPartitionsWithResponse(ctx, params)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0043OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.43"); err != nil {
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
	partitionList := make([]types.Partition, 0, len(*resp.JSON200.Partitions))
	for _, apiPartition := range *resp.JSON200.Partitions {
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
	params := &api.SlurmV0043GetPartitionParams{}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmV0043GetPartitionWithResponse(ctx, partitionName, params)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0043OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.43"); err != nil {
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
	if len(*resp.JSON200.Partitions) == 0 {
		return nil, common.NewResourceNotFoundError("Partition", partitionName)
	}

	// Convert the first partition (should be the only one)
	partition, err := a.convertAPIPartitionToCommon((*resp.JSON200.Partitions)[0])
	if err != nil {
		return nil, a.HandleConversionError(err, partitionName)
	}

	return partition, nil
}

// Create creates a new partition
func (a *PartitionAdapter) Create(ctx context.Context, partition *types.PartitionCreate) (*types.PartitionCreateResponse, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if err := a.validatePartitionCreate(partition); err != nil {
		return nil, err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Convert to API format
	apiPartition, err := a.convertCommonPartitionCreateToAPI(partition)
	if err != nil {
		return nil, err
	}

	// Create request body
	reqBody := api.SlurmV0043PostPartitionJSONRequestBody{
		Partitions: []api.V0043PartitionInfo{*apiPartition},
	}

	// Prepare parameters for the API call
	params := &api.SlurmV0043PostPartitionParams{}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmV0043PostPartitionWithResponse(ctx, params, reqBody)
	if err != nil {
		return nil, a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0043OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	if err := common.HandleAPIResponse(responseAdapter, "v0.0.43"); err != nil {
		return nil, err
	}

	return &types.PartitionCreateResponse{
		PartitionName: partition.Name,
	}, nil
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
	reqBody := api.SlurmV0043PostPartitionJSONRequestBody{
		Partitions: []api.V0043PartitionInfo{*apiPartition},
	}

	// Prepare parameters for the API call
	params := &api.SlurmV0043PostPartitionParams{}

	// Call the generated OpenAPI client (POST is used for updates in SLURM API)
	resp, err := a.client.SlurmV0043PostPartitionWithResponse(ctx, params, reqBody)
	if err != nil {
		return a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0043OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)
	return common.HandleAPIResponse(responseAdapter, "v0.0.43")
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
	resp, err := a.client.SlurmV0043DeletePartitionWithResponse(ctx, partitionName)
	if err != nil {
		return a.HandleAPIError(err)
	}

	// Use common response error handling
	var apiErrors *api.V0043OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	}

	// Create adapter with special handling for 204 (No Content) status
	responseAdapter := api.NewResponseAdapter(resp.StatusCode(), apiErrors)

	// For DELETE operations, 204 is also a success
	if resp.StatusCode() == 204 {
		return nil
	}

	return common.HandleAPIResponse(responseAdapter, "v0.0.43")
}

// validatePartitionCreate validates partition creation request
func (a *PartitionAdapter) validatePartitionCreate(partition *types.PartitionCreate) error {
	if partition == nil {
		return common.NewValidationError("partition creation data is required", "partition", nil)
	}
	if partition.Name == "" {
		return common.NewValidationError("partition name is required", "name", partition.Name)
	}
	// Validate numeric fields
	if partition.MaxNodes < 0 {
		return common.NewValidationError("max nodes must be non-negative", "maxNodes", partition.MaxNodes)
	}
	if partition.MinNodes < 0 {
		return common.NewValidationError("min nodes must be non-negative", "minNodes", partition.MinNodes)
	}
	if partition.MinNodes > 0 && partition.MaxNodes > 0 && partition.MinNodes > partition.MaxNodes {
		return common.NewValidationError("min nodes cannot be greater than max nodes", "minNodes/maxNodes", nil)
	}
	if partition.DefaultTime < 0 {
		return common.NewValidationError("default time must be non-negative", "defaultTime", partition.DefaultTime)
	}
	if partition.MaxTime < 0 {
		return common.NewValidationError("max time must be non-negative", "maxTime", partition.MaxTime)
	}
	if partition.Priority < 0 {
		return common.NewValidationError("priority must be non-negative", "priority", partition.Priority)
	}
	return nil
}

// validatePartitionUpdate validates partition update request
func (a *PartitionAdapter) validatePartitionUpdate(update *types.PartitionUpdate) error {
	if update == nil {
		return common.NewValidationError("partition update data is required", "update", nil)
	}
	// At least one field should be provided for update
	if update.State == nil && update.AllowAccounts == nil && update.DenyAccounts == nil &&
	   update.AllowQoS == nil && update.DenyQoS == nil && update.MaxNodes == nil &&
	   update.MinNodes == nil && update.DefaultTime == nil && update.MaxTime == nil &&
	   update.Priority == nil && update.Hidden == nil && update.RootOnly == nil {
		return common.NewValidationError("at least one field must be provided for update", "update", update)
	}
	
	// Validate numeric fields if provided
	if update.MaxNodes != nil && *update.MaxNodes < 0 {
		return common.NewValidationError("max nodes must be non-negative", "maxNodes", *update.MaxNodes)
	}
	if update.MinNodes != nil && *update.MinNodes < 0 {
		return common.NewValidationError("min nodes must be non-negative", "minNodes", *update.MinNodes)
	}
	if update.MinNodes != nil && update.MaxNodes != nil && *update.MinNodes > *update.MaxNodes {
		return common.NewValidationError("min nodes cannot be greater than max nodes", "minNodes/maxNodes", nil)
	}
	if update.DefaultTime != nil && *update.DefaultTime < 0 {
		return common.NewValidationError("default time must be non-negative", "defaultTime", *update.DefaultTime)
	}
	if update.MaxTime != nil && *update.MaxTime < 0 {
		return common.NewValidationError("max time must be non-negative", "maxTime", *update.MaxTime)
	}
	if update.Priority != nil && *update.Priority < 0 {
		return common.NewValidationError("priority must be non-negative", "priority", *update.Priority)
	}
	return nil
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
	// Filter by names (already handled by API, but included for completeness)
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