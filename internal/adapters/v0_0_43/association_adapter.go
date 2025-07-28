package v0_0_43

import (
	"context"
	"strings"

	"github.com/jontk/slurm-client/internal/common"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	api "github.com/jontk/slurm-client/internal/api/v0_0_43"
)

// AssociationAdapter implements the AssociationAdapter interface for v0.0.43
type AssociationAdapter struct {
	*base.BaseManager
	client  *api.ClientWithResponses
	wrapper *api.WrapperClient
}

// NewAssociationAdapter creates a new Association adapter for v0.0.43
func NewAssociationAdapter(client *api.ClientWithResponses) *AssociationAdapter {
	return &AssociationAdapter{
		BaseManager: base.NewBaseManager("v0.0.43", "Association"),
		client:      client,
		wrapper:     nil, // We'll implement this later
	}
}

// List retrieves a list of associations with optional filtering
func (a *AssociationAdapter) List(ctx context.Context, opts *types.AssociationListOptions) (*types.AssociationList, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Prepare parameters for the API call
	params := &api.SlurmdbV0043GetAssociationsParams{}

	// Apply filters from options
	if opts != nil {
		if len(opts.Accounts) > 0 {
			accountStr := strings.Join(opts.Accounts, ",")
			params.Account = &accountStr
		}
		if len(opts.Users) > 0 {
			userStr := strings.Join(opts.Users, ",")
			params.User = &userStr
		}
		if len(opts.Clusters) > 0 {
			clusterStr := strings.Join(opts.Clusters, ",")
			params.Cluster = &clusterStr
		}
		if len(opts.Partitions) > 0 {
			partitionStr := strings.Join(opts.Partitions, ",")
			params.Partition = &partitionStr
		}
		if opts.WithDeleted {
			withDeleted := "true"
			params.WithDeleted = &withDeleted
		}
	}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmdbV0043GetAssociationsWithResponse(ctx, params)
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
	if err := a.CheckNilResponse(resp.JSON200, "List Associations"); err != nil {
		return nil, err
	}
	if err := a.CheckNilResponse(resp.JSON200.Associations, "List Associations - associations field"); err != nil {
		return nil, err
	}

	// Convert the response to common types
	associationList := make([]types.Association, 0, len(*resp.JSON200.Associations))
	for _, apiAssociation := range *resp.JSON200.Associations {
		association, err := a.convertAPIAssociationToCommon(apiAssociation)
		if err != nil {
			return nil, a.HandleConversionError(err, apiAssociation.Id)
		}
		associationList = append(associationList, *association)
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
	if start >= len(associationList) {
		return &types.AssociationList{
			Associations: []types.Association{},
			Total:        len(associationList),
		}, nil
	}

	end := len(associationList)
	if listOpts.Limit > 0 {
		end = start + listOpts.Limit
		if end > len(associationList) {
			end = len(associationList)
		}
	}

	return &types.AssociationList{
		Associations: associationList[start:end],
		Total:        len(associationList),
	}, nil
}

// Get retrieves a specific association by ID
func (a *AssociationAdapter) Get(ctx context.Context, associationID string) (*types.Association, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if err := a.ValidateResourceName(associationID, "associationID"); err != nil {
		return nil, err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Prepare parameters for the API call
	params := &api.SlurmdbV0043GetSingleAssociationParams{}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmdbV0043GetSingleAssociationWithResponse(ctx, associationID, params)
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
	if err := a.CheckNilResponse(resp.JSON200, "Get Association"); err != nil {
		return nil, err
	}
	if err := a.CheckNilResponse(resp.JSON200.Associations, "Get Association - associations field"); err != nil {
		return nil, err
	}

	// Check if we got any association entries
	if len(*resp.JSON200.Associations) == 0 {
		return nil, common.NewResourceNotFoundError("Association", associationID)
	}

	// Convert the first association (should be the only one)
	association, err := a.convertAPIAssociationToCommon((*resp.JSON200.Associations)[0])
	if err != nil {
		return nil, a.HandleConversionError(err, associationID)
	}

	return association, nil
}

// Create creates a new association
func (a *AssociationAdapter) Create(ctx context.Context, association *types.AssociationCreate) (*types.AssociationCreateResponse, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if err := a.validateAssociationCreate(association); err != nil {
		return nil, err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Convert to API format
	apiAssociation, err := a.convertCommonAssociationCreateToAPI(association)
	if err != nil {
		return nil, err
	}

	// Create request body
	reqBody := api.SlurmdbV0043PostAssociationsJSONRequestBody{
		Associations: []api.V0043AssociationInfo{*apiAssociation},
	}

	// Prepare parameters for the API call
	params := &api.SlurmdbV0043PostAssociationsParams{}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmdbV0043PostAssociationsWithResponse(ctx, params, reqBody)
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

	return &types.AssociationCreateResponse{
		AssociationID: association.Account + ":" + association.User + ":" + association.Cluster,
	}, nil
}

// Update updates an existing association
func (a *AssociationAdapter) Update(ctx context.Context, associationID string, update *types.AssociationUpdate) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.ValidateResourceName(associationID, "associationID"); err != nil {
		return err
	}
	if err := a.validateAssociationUpdate(update); err != nil {
		return err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// First, get the existing association to merge updates
	existingAssociation, err := a.Get(ctx, associationID)
	if err != nil {
		return err
	}

	// Convert to API format and apply updates
	apiAssociation, err := a.convertCommonAssociationUpdateToAPI(existingAssociation, update)
	if err != nil {
		return err
	}

	// Create request body
	reqBody := api.SlurmdbV0043PostAssociationsJSONRequestBody{
		Associations: []api.V0043AssociationInfo{*apiAssociation},
	}

	// Prepare parameters for the API call
	params := &api.SlurmdbV0043PostAssociationsParams{}

	// Call the generated OpenAPI client (POST is used for updates in SLURM API)
	resp, err := a.client.SlurmdbV0043PostAssociationsWithResponse(ctx, params, reqBody)
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

// Delete deletes an association
func (a *AssociationAdapter) Delete(ctx context.Context, associationID string) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.ValidateResourceName(associationID, "associationID"); err != nil {
		return err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmdbV0043DeleteSingleAssociationWithResponse(ctx, associationID)
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

// validateAssociationCreate validates association creation request
func (a *AssociationAdapter) validateAssociationCreate(association *types.AssociationCreate) error {
	if association == nil {
		return common.NewValidationError("association creation data is required", "association", nil)
	}
	if association.Account == "" {
		return common.NewValidationError("account is required", "account", association.Account)
	}
	if association.User == "" {
		return common.NewValidationError("user is required", "user", association.User)
	}
	if association.Cluster == "" {
		return common.NewValidationError("cluster is required", "cluster", association.Cluster)
	}
	return nil
}

// validateAssociationUpdate validates association update request
func (a *AssociationAdapter) validateAssociationUpdate(update *types.AssociationUpdate) error {
	if update == nil {
		return common.NewValidationError("association update data is required", "update", nil)
	}
	// At least one field should be provided for update
	if update.DefaultQoS == nil && len(update.QoSList) == 0 && update.FairShare == nil &&
	   update.MaxJobs == nil && update.MaxWallTime == nil {
		return common.NewValidationError("at least one field must be provided for update", "update", update)
	}
	return nil
}

// Simplified converter methods for association management
func (a *AssociationAdapter) convertAPIAssociationToCommon(apiAssociation api.V0043AssociationInfo) (*types.Association, error) {
	association := &types.Association{}
	if apiAssociation.Account != nil {
		association.Account = *apiAssociation.Account
	}
	if apiAssociation.User != nil {
		association.User = *apiAssociation.User
	}
	if apiAssociation.Cluster != nil {
		association.Cluster = *apiAssociation.Cluster
	}
	if apiAssociation.Partition != nil {
		association.Partition = *apiAssociation.Partition
	}
	// TODO: Add more field conversions as needed
	return association, nil
}

func (a *AssociationAdapter) convertCommonAssociationCreateToAPI(create *types.AssociationCreate) (*api.V0043AssociationInfo, error) {
	apiAssociation := &api.V0043AssociationInfo{}
	apiAssociation.Account = &create.Account
	apiAssociation.User = &create.User
	apiAssociation.Cluster = &create.Cluster
	if create.Partition != "" {
		apiAssociation.Partition = &create.Partition
	}
	// TODO: Add more field conversions as needed
	return apiAssociation, nil
}

func (a *AssociationAdapter) convertCommonAssociationUpdateToAPI(existing *types.Association, update *types.AssociationUpdate) (*api.V0043AssociationInfo, error) {
	apiAssociation := &api.V0043AssociationInfo{}
	apiAssociation.Account = &existing.Account
	apiAssociation.User = &existing.User
	apiAssociation.Cluster = &existing.Cluster
	if existing.Partition != "" {
		apiAssociation.Partition = &existing.Partition
	}
	// TODO: Add more field conversions as needed
	return apiAssociation, nil
}