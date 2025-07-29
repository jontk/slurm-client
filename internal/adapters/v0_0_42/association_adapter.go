package v0_0_42

import (
	"context"
	"fmt"
	"strings"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	"github.com/jontk/slurm-client/pkg/errors"
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
		// Note: v0.0.42 API doesn't support WithSubAccounts or WithDeleted parameters
		// These would be filtered client-side if needed
	}

	// Call the API
	resp, err := a.client.SlurmdbV0042GetAssociationsWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list associations: %w", err)
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return nil, a.HandleAPIError(fmt.Errorf("API error: status %d", resp.StatusCode()))
	}

	// Check for API response
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response from API")
	}

	// Convert the response to common types
	assocList := &types.AssociationList{
		Associations: make([]types.Association, 0),
	}

	if resp.JSON200.Associations != nil {
		for _, apiAssoc := range resp.JSON200.Associations {
			assoc, err := a.convertAPIAssociationToCommon(apiAssoc)
			if err != nil {
				// Log conversion error but continue
				continue
			}
			assocList.Associations = append(assocList.Associations, *assoc)
		}
	}

	return assocList, nil
}

// Get retrieves a specific association
func (a *AssociationAdapter) Get(ctx context.Context, associationID string) (*types.Association, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// v0.0.42 doesn't have a single association get endpoint, use list to find by ID
	params := &api.SlurmdbV0042GetAssociationsParams{}

	// Call the API
	resp, err := a.client.SlurmdbV0042GetAssociationsWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get association: %w", err)
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return nil, a.HandleAPIError(fmt.Errorf("API error: status %d", resp.StatusCode()))
	}

	// Check for API response
	if resp.JSON200 == nil || resp.JSON200.Associations == nil {
		return nil, fmt.Errorf("association not found")
	}

	// Find the matching association by ID
	for _, apiAssoc := range resp.JSON200.Associations {
		assoc, err := a.convertAPIAssociationToCommon(apiAssoc)
		if err != nil {
			continue
		}
		if assoc.ID == associationID {
			return assoc, nil
		}
	}

	return nil, fmt.Errorf("association %s not found", associationID)
}

// Create creates a new association
func (a *AssociationAdapter) Create(ctx context.Context, assoc *types.AssociationCreate) (*types.AssociationCreateResponse, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Convert common association to API format
	apiAssoc, err := a.convertCommonAssociationCreateToAPI(assoc)
	if err != nil {
		return nil, fmt.Errorf("failed to convert association create request: %w", err)
	}

	// Call the API
	resp, err := a.client.SlurmdbV0042PostAssociationsWithResponse(ctx, *apiAssoc)
	if err != nil {
		return nil, fmt.Errorf("failed to create association: %w", err)
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return nil, a.HandleAPIError(fmt.Errorf("API error: status %d", resp.StatusCode()))
	}

	// Return success response - for v0.0.42, we can't get specific association details from create response
	return &types.AssociationCreateResponse{
		AssociationID: fmt.Sprintf("%s_%s", assoc.AccountName, assoc.UserName),
		AccountName:   assoc.AccountName,
		UserName:      assoc.UserName,
		Cluster:       assoc.Cluster,
	}, nil
}

// Update updates an existing association
func (a *AssociationAdapter) Update(ctx context.Context, associationID string, update *types.AssociationUpdate) error {
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
	apiAssoc, err := a.convertCommonAssociationUpdateToAPI(update)
	if err != nil {
		return fmt.Errorf("failed to convert association update request: %w", err)
	}

	// Call the API
	resp, err := a.client.SlurmdbV0042PostAssociationsWithResponse(ctx, *apiAssoc)
	if err != nil {
		return fmt.Errorf("failed to update association: %w", err)
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return a.HandleAPIError(fmt.Errorf("API error: status %d", resp.StatusCode()))
	}

	return nil
}

// Delete deletes an association
func (a *AssociationAdapter) Delete(ctx context.Context, associationID string) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// First get the association to extract the required fields
	assoc, err := a.Get(ctx, associationID)
	if err != nil {
		return err
	}

	// Build parameters using the association details
	params := &api.SlurmdbV0042DeleteAssociationParams{
		Account: &assoc.AccountName,
		User:    &assoc.UserName,
	}
	if assoc.Cluster != "" {
		params.Cluster = &assoc.Cluster
	}
	if assoc.Partition != "" {
		params.Partition = &assoc.Partition
	}

	// Call the API
	resp, err := a.client.SlurmdbV0042DeleteAssociationWithResponse(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to delete association: %w", err)
	}

	// Check response status
	if resp.StatusCode() != 200 {
		return a.HandleAPIError(fmt.Errorf("API error: status %d", resp.StatusCode()))
	}

	return nil
}