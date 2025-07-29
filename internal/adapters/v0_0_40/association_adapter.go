package v0_0_40

import (
	"context"
	"strings"

	"github.com/jontk/slurm-client/internal/common"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	api "github.com/jontk/slurm-client/internal/api/v0_0_40"
)

// AssociationAdapter implements the AssociationAdapter interface for v0.0.40
type AssociationAdapter struct {
	*base.BaseManager
	client  *api.ClientWithResponses
	wrapper *api.WrapperClient
}

// NewAssociationAdapter creates a new Association adapter for v0.0.40
func NewAssociationAdapter(client *api.ClientWithResponses) *AssociationAdapter {
	return &AssociationAdapter{
		BaseManager: base.NewBaseManager("v0.0.40", "Association"),
		client:      client,
		wrapper:     nil,
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
	params := &api.SlurmdbV0040GetAssociationsParams{}

	// Apply filters from options
	if opts != nil {
		if len(opts.Accounts) > 0 {
			accountStr := strings.Join(opts.Accounts, ",")
			params.Account = &accountStr
		}
		if len(opts.Clusters) > 0 {
			clusterStr := strings.Join(opts.Clusters, ",")
			params.Cluster = &clusterStr
		}
		if len(opts.Users) > 0 {
			userStr := strings.Join(opts.Users, ",")
			params.User = &userStr
		}
		if len(opts.Partitions) > 0 {
			partitionStr := strings.Join(opts.Partitions, ",")
			params.Partition = &partitionStr
		}
		if opts.WithDeleted {
			withDeleted := "true"
			params.WithDeleted = &withDeleted
		}
		if opts.WithSubAccounts {
			withSubAccts := "true"
			params.WithSubAccts = &withSubAccts
		}
	}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmdbV0040GetAssociationsWithResponse(ctx, params)
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
	if err := a.CheckNilResponse(resp.JSON200, "List Associations"); err != nil {
		return nil, err
	}
	if err := a.CheckNilResponse(resp.JSON200.Associations, "List Associations - associations field"); err != nil {
		return nil, err
	}

	// Convert the response to common types
	associationList := make([]types.Association, 0, len(resp.JSON200.Associations))
	for _, apiAssociation := range resp.JSON200.Associations {
		association, err := a.convertAPIAssociationToCommon(apiAssociation)
		if err != nil {
			return nil, a.HandleConversionError(err, "association")
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

	// v0.0.40 doesn't have a single association GET endpoint
	// We need to list all and filter
	list, err := a.List(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Find the association by constructing ID from fields
	for _, assoc := range list.Associations {
		constructedID := a.constructAssociationID(assoc)
		if constructedID == associationID {
			return &assoc, nil
		}
	}

	return nil, common.NewResourceNotFoundError("Association", associationID)
}

// Create creates a new association
func (a *AssociationAdapter) Create(ctx context.Context, association *types.AssociationCreate) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.validateAssociationCreate(association); err != nil {
		return err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Convert to API format
	apiAssociation, err := a.convertCommonAssociationCreateToAPI(association)
	if err != nil {
		return err
	}

	// Create request body
	reqBody := api.SlurmdbV0040PostAssociationsJSONRequestBody{
		Associations: &[]api.V0040AssociationShort{*apiAssociation},
	}

	// Call the generated OpenAPI client
	resp, err := a.client.SlurmdbV0040PostAssociationsWithResponse(ctx, reqBody)
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

// Update updates an existing association
func (a *AssociationAdapter) Update(ctx context.Context, associationID string, update *types.AssociationUpdate) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.ValidateResourceName(associationID, "associationID"); err != nil {
		return err
	}

	// v0.0.40 may not support association updates directly
	return common.NewNotImplementedError("Update Association is not implemented for v0.0.40")
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

	// v0.0.40 doesn't have a direct association delete endpoint
	// You would typically delete by removing the user from account
	return common.NewNotImplementedError("Delete Association is not implemented for v0.0.40")
}

// constructAssociationID constructs an ID from association fields
func (a *AssociationAdapter) constructAssociationID(assoc types.Association) string {
	parts := []string{assoc.Cluster, assoc.Account, assoc.User}
	if assoc.Partition != "" {
		parts = append(parts, assoc.Partition)
	}
	return strings.Join(parts, ":")
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