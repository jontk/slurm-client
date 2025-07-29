package v0_0_42

import (
	"context"
	"fmt"
	"strings"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	api "github.com/jontk/slurm-client/internal/api/v0_0_42"
)

// AssociationAdapter implements the AssociationAdapter interface for v0.0.42
type AssociationAdapter struct {
	*base.BaseManager
	client *api.ClientWithResponses
}

// NewAssociationAdapter creates a new Association adapter for v0.0.42
func NewAssociationAdapter(client *api.ClientWithResponses) *AssociationAdapter {
	return &AssociationAdapter{
		BaseManager: base.NewBaseManager("v0.0.42", "Association"),
		client:      client,
	}
}

// List retrieves a list of associations
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
	params := &api.SlurmdbV0042GetAssociationsParams{}

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
		if opts.WithSubAccounts {
			withSubAccts := "true"
			params.WithSubAccts = &withSubAccts
		}
		if opts.WithDeleted {
			withDeleted := "true"
			params.WithDeleted = &withDeleted
		}
	}

	// Call the API
	resp, err := a.client.SlurmdbV0042GetAssociationsWithResponse(ctx, params)
	if err != nil {
		return nil, a.WrapError(err, "failed to list associations")
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return nil, a.HandleAPIError(resp.StatusCode(), resp.Body)
	}

	// Check for API response
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response from API")
	}

	// Convert the response to common types
	assocList := &types.AssociationList{
		Associations: make([]*types.Association, 0),
	}

	if resp.JSON200.Associations != nil {
		for _, apiAssoc := range *resp.JSON200.Associations {
			assoc, err := a.convertAPIAssociationToCommon(apiAssoc)
			if err != nil {
				// Log conversion error but continue
				continue
			}
			assocList.Associations = append(assocList.Associations, assoc)
		}
	}

	return assocList, nil
}

// Get retrieves a specific association
func (a *AssociationAdapter) Get(ctx context.Context, account, user, cluster, partition string) (*types.Association, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// v0.0.42 doesn't have a single association get endpoint, use list with filters
	params := &api.SlurmdbV0042GetAssociationsParams{
		Account: &account,
		User:    &user,
	}
	if cluster != "" {
		params.Cluster = &cluster
	}
	if partition != "" {
		params.Partition = &partition
	}

	// Call the API
	resp, err := a.client.SlurmdbV0042GetAssociationsWithResponse(ctx, params)
	if err != nil {
		return nil, a.WrapError(err, "failed to get association")
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return nil, a.HandleAPIError(resp.StatusCode(), resp.Body)
	}

	// Check for API response
	if resp.JSON200 == nil || resp.JSON200.Associations == nil || len(*resp.JSON200.Associations) == 0 {
		return nil, fmt.Errorf("association not found")
	}

	// Find the matching association
	for _, apiAssoc := range *resp.JSON200.Associations {
		if apiAssoc.Account != nil && *apiAssoc.Account == account &&
			apiAssoc.User != nil && *apiAssoc.User == user {
			if (cluster == "" || (apiAssoc.Cluster != nil && *apiAssoc.Cluster == cluster)) &&
				(partition == "" || (apiAssoc.Partition != nil && *apiAssoc.Partition == partition)) {
				return a.convertAPIAssociationToCommon(apiAssoc)
			}
		}
	}

	return nil, fmt.Errorf("association not found")
}

// Create creates a new association
func (a *AssociationAdapter) Create(ctx context.Context, assoc *types.AssociationCreateRequest) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Convert common association to API format
	apiAssoc, err := a.convertCommonAssociationCreateToAPI(assoc)
	if err != nil {
		return a.WrapError(err, "failed to convert association create request")
	}

	// Call the API
	resp, err := a.client.SlurmdbV0042PostAssociationsWithResponse(ctx, apiAssoc)
	if err != nil {
		return a.WrapError(err, "failed to create association")
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return a.HandleAPIError(resp.StatusCode(), resp.Body)
	}

	return nil
}

// Update updates an existing association
func (a *AssociationAdapter) Update(ctx context.Context, assoc *types.AssociationUpdateRequest) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// v0.0.42 doesn't have a direct association update endpoint
	// Updates are done through the create endpoint with update semantics
	apiAssoc, err := a.convertCommonAssociationUpdateToAPI(assoc)
	if err != nil {
		return a.WrapError(err, "failed to convert association update request")
	}

	// Call the API
	resp, err := a.client.SlurmdbV0042PostAssociationsWithResponse(ctx, apiAssoc)
	if err != nil {
		return a.WrapError(err, "failed to update association")
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return a.HandleAPIError(resp.StatusCode(), resp.Body)
	}

	return nil
}

// Delete deletes an association
func (a *AssociationAdapter) Delete(ctx context.Context, account, user, cluster, partition string) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Build parameters
	params := &api.SlurmdbV0042DeleteAssociationParams{
		Account: &account,
		User:    &user,
	}
	if cluster != "" {
		params.Cluster = &cluster
	}
	if partition != "" {
		params.Partition = &partition
	}

	// Call the API
	resp, err := a.client.SlurmdbV0042DeleteAssociationWithResponse(ctx, params)
	if err != nil {
		return a.WrapError(err, "failed to delete association")
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return a.HandleAPIError(resp.StatusCode(), resp.Body)
	}

	return nil
}