// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_42

import (
	"context"
	"fmt"
	"strings"

	api "github.com/jontk/slurm-client/internal/api/v0_0_42"
	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
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
		Status:  "success",
		Message: "Association created successfully",
		Meta: map[string]interface{}{
			"association_id": fmt.Sprintf("%s_%s", assoc.AccountName, assoc.UserName),
			"account_name":   assoc.AccountName,
			"user_name":      assoc.UserName,
			"cluster":        assoc.Cluster,
		},
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

// convertAPIAssociationToCommon converts API association to common type
func (a *AssociationAdapter) convertAPIAssociationToCommon(apiAssoc api.V0042Assoc) (*types.Association, error) {
	assoc := &types.Association{}

	// Set basic fields
	if apiAssoc.Account != nil {
		assoc.AccountName = *apiAssoc.Account
	}

	// User is a string in V0042Assoc, not a pointer
	assoc.UserName = apiAssoc.User

	if apiAssoc.Cluster != nil {
		assoc.Cluster = *apiAssoc.Cluster
	}

	if apiAssoc.Partition != nil {
		assoc.Partition = *apiAssoc.Partition
	}

	if apiAssoc.ParentAccount != nil {
		assoc.ParentAccount = *apiAssoc.ParentAccount
	}

	// Default QoS is nested in Default.Qos
	if apiAssoc.Default != nil && apiAssoc.Default.Qos != nil {
		assoc.DefaultQoS = *apiAssoc.Default.Qos
	}

	// Shares are in SharesRaw field
	if apiAssoc.SharesRaw != nil {
		assoc.SharesRaw = *apiAssoc.SharesRaw
	}

	// Priority is nested
	if apiAssoc.Priority != nil && apiAssoc.Priority.Set != nil && *apiAssoc.Priority.Set && apiAssoc.Priority.Number != nil {
		assoc.Priority = *apiAssoc.Priority.Number
	}

	// Max jobs are nested in Max.Jobs structure
	if apiAssoc.Max != nil && apiAssoc.Max.Jobs != nil {
		if apiAssoc.Max.Jobs.Active != nil && apiAssoc.Max.Jobs.Active.Set != nil && *apiAssoc.Max.Jobs.Active.Set && apiAssoc.Max.Jobs.Active.Number != nil {
			assoc.MaxJobs = *apiAssoc.Max.Jobs.Active.Number
		}
		if apiAssoc.Max.Jobs.Total != nil && apiAssoc.Max.Jobs.Total.Set != nil && *apiAssoc.Max.Jobs.Total.Set && apiAssoc.Max.Jobs.Total.Number != nil {
			assoc.MaxSubmitJobs = *apiAssoc.Max.Jobs.Total.Number
		}
	}

	// TRES limits are nested in Max.Tres structure
	if apiAssoc.Max != nil && apiAssoc.Max.Tres != nil {
		assoc.MaxTRES = make(map[string]int64)
		if apiAssoc.Max.Tres.Total != nil {
			for _, tres := range *apiAssoc.Max.Tres.Total {
				if tres.Count != nil {
					// V0042Tres has Type as string, not pointer
					assoc.MaxTRES[tres.Type] = *tres.Count
				}
			}
		}
	}

	// Set ID as a combination of account and user
	if assoc.AccountName != "" && assoc.UserName != "" {
		assoc.ID = fmt.Sprintf("%s_%s", assoc.AccountName, assoc.UserName)
	} else if assoc.AccountName != "" {
		assoc.ID = assoc.AccountName
	}

	// Set other fields if available
	if apiAssoc.Id != nil {
		assoc.ID = fmt.Sprintf("%d", *apiAssoc.Id)
	}

	if apiAssoc.Comment != nil {
		assoc.Comment = *apiAssoc.Comment
	}

	if apiAssoc.IsDefault != nil {
		assoc.IsDefault = *apiAssoc.IsDefault
	}

	return assoc, nil
}

// convertCommonAssociationCreateToAPI converts common association create to API format
func (a *AssociationAdapter) convertCommonAssociationCreateToAPI(assocCreate *types.AssociationCreate) (*api.V0042OpenapiAssocsResp, error) {
	if assocCreate == nil {
		return nil, fmt.Errorf("association create request cannot be nil")
	}

	apiAssoc := api.V0042Assoc{}

	if assocCreate.AccountName != "" {
		apiAssoc.Account = &assocCreate.AccountName
	}

	// User is a string in V0042Assoc, not a pointer
	apiAssoc.User = assocCreate.UserName

	if assocCreate.Cluster != "" {
		apiAssoc.Cluster = &assocCreate.Cluster
	}

	if assocCreate.Partition != "" {
		apiAssoc.Partition = &assocCreate.Partition
	}

	if assocCreate.ParentAccount != "" {
		apiAssoc.ParentAccount = &assocCreate.ParentAccount
	}

	// Default QoS is nested in Default.Qos
	if assocCreate.DefaultQoS != "" {
		apiAssoc.Default = &struct {
			Qos *string `json:"qos,omitempty"`
		}{
			Qos: &assocCreate.DefaultQoS,
		}
	}

	// Shares are in SharesRaw field
	if assocCreate.SharesRaw != 0 {
		apiAssoc.SharesRaw = &assocCreate.SharesRaw
	}

	if assocCreate.Comment != "" {
		apiAssoc.Comment = &assocCreate.Comment
	}

	apiAssoc.IsDefault = &assocCreate.IsDefault

	return &api.V0042OpenapiAssocsResp{
		Associations: []api.V0042Assoc{apiAssoc},
	}, nil
}

// convertCommonAssociationUpdateToAPI converts common association update to API format
func (a *AssociationAdapter) convertCommonAssociationUpdateToAPI(assocUpdate *types.AssociationUpdate) (*api.V0042OpenapiAssocsResp, error) {
	if assocUpdate == nil {
		return nil, fmt.Errorf("association update request cannot be nil")
	}

	apiAssoc := api.V0042Assoc{}

	// Default QoS is nested in Default.Qos
	if assocUpdate.DefaultQoS != nil {
		apiAssoc.Default = &struct {
			Qos *string `json:"qos,omitempty"`
		}{
			Qos: assocUpdate.DefaultQoS,
		}
	}

	// Shares are in SharesRaw field
	if assocUpdate.SharesRaw != nil {
		apiAssoc.SharesRaw = assocUpdate.SharesRaw
	}

	if assocUpdate.Comment != nil {
		apiAssoc.Comment = assocUpdate.Comment
	}

	if assocUpdate.IsDefault != nil {
		apiAssoc.IsDefault = assocUpdate.IsDefault
	}

	return &api.V0042OpenapiAssocsResp{
		Associations: []api.V0042Assoc{apiAssoc},
	}, nil
}
