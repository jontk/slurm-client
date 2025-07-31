package v0_0_40

import (
	"context"
	"fmt"
	"strings"

	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
)

// AssociationManagerImpl implements the AssociationManager interface for v0.0.40
type AssociationManagerImpl struct {
	client *WrapperClient
}

// NewAssociationManagerImpl creates a new AssociationManagerImpl
func NewAssociationManagerImpl(client *WrapperClient) *AssociationManagerImpl {
	return &AssociationManagerImpl{
		client: client,
	}
}

// List retrieves a list of associations with optional filtering
func (a *AssociationManagerImpl) List(ctx context.Context, opts *interfaces.ListAssociationsOptions) (*interfaces.AssociationList, error) {
	if err := a.client.CheckContext(ctx); err != nil {
		return nil, err
	}

	// Prepare parameters
	params := &SlurmdbV0040GetAssociationsParams{}
	if opts != nil {
		if opts.Accounts != nil && len(opts.Accounts) > 0 {
			accountsStr := strings.Join(opts.Accounts, ",")
			params.Account = &accountsStr
		}
		if opts.Users != nil && len(opts.Users) > 0 {
			usersStr := strings.Join(opts.Users, ",")
			params.User = &usersStr
		}
		if opts.Clusters != nil && len(opts.Clusters) > 0 {
			clustersStr := strings.Join(opts.Clusters, ",")
			params.Cluster = &clustersStr
		}
		if opts.WithDeleted {
			withDeleted := "true"
			params.WithDeleted = &withDeleted
		}
	}

	// Make API call
	resp, err := a.client.apiClient.SlurmdbV0040GetAssociationsWithResponse(ctx, params)
	if err != nil {
		return nil, errors.NewSlurmErrorWithCause(errors.ErrorCodeServerInternal, "failed to list associations", err)
	}

	// Check response
	if resp.StatusCode() != 200 {
		return nil, a.client.HandleErrorResponse(resp.StatusCode(), resp.Body)
	}

	if resp.JSON200 == nil {
		return nil, errors.NewSlurmError(errors.ErrorCodeServerInternal, "received nil response")
	}

	// Convert response
	associationList := &interfaces.AssociationList{
		Associations: make([]*interfaces.Association, 0),
	}

	if resp.JSON200.Associations != nil {
		for _, assoc := range resp.JSON200.Associations {
			association := a.convertV0040AssociationToInterface(assoc)
			associationList.Associations = append(associationList.Associations, association)
		}
	}

	return associationList, nil
}

// Get retrieves a specific association by ID
func (a *AssociationManagerImpl) Get(ctx context.Context, associationID string) (*interfaces.Association, error) {
	if associationID == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "association ID is required", "associationID", associationID, nil)
	}

	if err := a.client.CheckContext(ctx); err != nil {
		return nil, err
	}

	// Make API call
	params := &SlurmdbV0040GetAssociationParams{}
	resp, err := a.client.apiClient.SlurmdbV0040GetAssociationWithResponse(ctx, params)
	if err != nil {
		return nil, errors.NewSlurmErrorWithCause(errors.ErrorCodeServerInternal, "failed to get association", err)
	}

	// Check response
	if resp.StatusCode() != 200 {
		return nil, a.client.HandleErrorResponse(resp.StatusCode(), resp.Body)
	}

	if resp.JSON200 == nil || resp.JSON200.Associations == nil || len(resp.JSON200.Associations) == 0 {
		return nil, errors.NewClientError(errors.ErrorCodeResourceNotFound, "association not found")
	}

	// Convert the first association
	associations := resp.JSON200.Associations
	return a.convertV0040AssociationToInterface(associations[0]), nil
}

// Create creates a new association
func (a *AssociationManagerImpl) Create(ctx context.Context, association *interfaces.AssociationCreate) (*interfaces.AssociationCreateResponse, error) {
	if association == nil {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "association data is required", "association", association, nil)
	}

	if err := a.client.CheckContext(ctx); err != nil {
		return nil, err
	}

	// Convert to API format
	apiAssoc := a.convertInterfaceAssociationCreateToV0040(association)
	
	// Create request body
	reqBody := V0040OpenapiAssocsResp{
		Associations: V0040AssocList{*apiAssoc},
	}

	// Make API call
	resp, err := a.client.apiClient.SlurmdbV0040PostAssociationsWithResponse(ctx, reqBody)
	if err != nil {
		return nil, errors.NewSlurmErrorWithCause(errors.ErrorCodeServerInternal, "failed to create association", err)
	}

	// Check response
	if resp.StatusCode() != 200 {
		return nil, a.client.HandleErrorResponse(resp.StatusCode(), resp.Body)
	}

	// Create a fake association to return
	createdAssoc := &interfaces.Association{
		ID:          0, // ID will be set by server
		Account:     association.Account,
		User:        association.User,
		Cluster:     association.Cluster,
	}
	
	return &interfaces.AssociationCreateResponse{
		Associations: []*interfaces.Association{createdAssoc},
		Created:      1,
		Updated:      0,
	}, nil
}

// Update updates an existing association
func (a *AssociationManagerImpl) Update(ctx context.Context, associationID string, update *interfaces.AssociationUpdate) error {
	if associationID == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "association ID is required", "associationID", associationID, nil)
	}
	if update == nil {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "update data is required", "update", update, nil)
	}

	if err := a.client.CheckContext(ctx); err != nil {
		return err
	}

	// Get existing association first
	existing, err := a.Get(ctx, associationID)
	if err != nil {
		return err
	}

	// Apply updates to existing association
	apiAssoc := a.convertInterfaceAssociationToV0040Update(existing, update)

	// Create request body
	reqBody := V0040OpenapiAssocsResp{
		Associations: V0040AssocList{*apiAssoc},
	}

	// Make API call
	resp, err := a.client.apiClient.SlurmdbV0040PostAssociationsWithResponse(ctx, reqBody)
	if err != nil {
		return errors.NewSlurmErrorWithCause(errors.ErrorCodeServerInternal, "failed to update association", err)
	}

	// Check response
	if resp.StatusCode() != 200 {
		return a.client.HandleErrorResponse(resp.StatusCode(), resp.Body)
	}

	return nil
}

// Delete deletes an association
func (a *AssociationManagerImpl) Delete(ctx context.Context, associationID string) error {
	if associationID == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "association ID is required", "associationID", associationID, nil)
	}

	if err := a.client.CheckContext(ctx); err != nil {
		return err
	}

	// Make API call
	params := &SlurmdbV0040DeleteAssociationParams{}
	resp, err := a.client.apiClient.SlurmdbV0040DeleteAssociationWithResponse(ctx, params)
	if err != nil {
		return errors.NewSlurmErrorWithCause(errors.ErrorCodeServerInternal, "failed to delete association", err)
	}

	// Check response
	if resp.StatusCode() != 200 && resp.StatusCode() != 204 {
		return a.client.HandleErrorResponse(resp.StatusCode(), resp.Body)
	}

	return nil
}


// GetUserAssociations retrieves all associations for a specific user
func (a *AssociationManagerImpl) GetUserAssociations(ctx context.Context, userName string) ([]*interfaces.Association, error) {
	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}

	// Use List with user filter
	listOpts := &interfaces.ListAssociationsOptions{
		Users: []string{userName},
	}

	result, err := a.List(ctx, listOpts)
	if err != nil {
		return nil, err
	}

	associations := make([]*interfaces.Association, len(result.Associations))
	for i := range result.Associations {
		associations[i] = result.Associations[i]
	}

	return associations, nil
}

// GetAccountAssociations retrieves all associations for a specific account
func (a *AssociationManagerImpl) GetAccountAssociations(ctx context.Context, accountName string) ([]*interfaces.Association, error) {
	if accountName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name is required", "accountName", accountName, nil)
	}

	// Use List with account filter
	listOpts := &interfaces.ListAssociationsOptions{
		Accounts: []string{accountName},
	}

	result, err := a.List(ctx, listOpts)
	if err != nil {
		return nil, err
	}

	associations := make([]*interfaces.Association, len(result.Associations))
	for i := range result.Associations {
		associations[i] = result.Associations[i]
	}

	return associations, nil
}

// GetClusterAssociations retrieves all associations for a specific cluster
func (a *AssociationManagerImpl) GetClusterAssociations(ctx context.Context, clusterName string) ([]*interfaces.Association, error) {
	if clusterName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "cluster name is required", "clusterName", clusterName, nil)
	}

	// Use List with cluster filter
	listOpts := &interfaces.ListAssociationsOptions{
		Clusters: []string{clusterName},
	}

	result, err := a.List(ctx, listOpts)
	if err != nil {
		return nil, err
	}

	associations := make([]*interfaces.Association, len(result.Associations))
	for i := range result.Associations {
		associations[i] = result.Associations[i]
	}

	return associations, nil
}


// convertV0040AssociationToInterface converts v0.0.40 association to interface format
func (a *AssociationManagerImpl) convertV0040AssociationToInterface(assoc V0040Assoc) *interfaces.Association {
	association := &interfaces.Association{}
	
	if assoc.Account != nil {
		association.Account = *assoc.Account
	}
	// User is a string, not *string in V0040Assoc
	association.User = assoc.User
	if assoc.Cluster != nil {
		association.Cluster = *assoc.Cluster
	}
	if assoc.Partition != nil {
		association.Partition = *assoc.Partition
	}
	if assoc.Id != nil && assoc.Id.Id != nil {
		association.ID = uint32(*assoc.Id.Id)
	}
	
	// TODO: Add more field conversions as needed
	
	return association
}

// convertInterfaceAssociationCreateToV0040 converts interface association create to v0.0.40 format
func (a *AssociationManagerImpl) convertInterfaceAssociationCreateToV0040(assoc *interfaces.AssociationCreate) *V0040Assoc {
	apiAssoc := &V0040Assoc{
		Account: &assoc.Account,
		User:    assoc.User, // User is string, not *string
		Cluster: &assoc.Cluster,
	}
	
	if assoc.Partition != "" {
		apiAssoc.Partition = &assoc.Partition
	}
	
	// TODO: Add more field conversions as needed
	
	return apiAssoc
}

// convertInterfaceAssociationToV0040Update converts interface association with updates to v0.0.40 format
func (a *AssociationManagerImpl) convertInterfaceAssociationToV0040Update(existing *interfaces.Association, update *interfaces.AssociationUpdate) *V0040Assoc {
	apiAssoc := &V0040Assoc{
		Account: &existing.Account,
		User:    existing.User, // User is string, not *string
		Cluster: &existing.Cluster,
	}
	
	if existing.Partition != "" {
		apiAssoc.Partition = &existing.Partition
	}
	
	// Apply updates
	// TODO: Add more field conversions as needed based on what fields can be updated
	
	return apiAssoc
}

// BulkDelete deletes multiple associations
func (a *AssociationManagerImpl) BulkDelete(ctx context.Context, opts *interfaces.BulkDeleteOptions) (*interfaces.BulkDeleteResponse, error) {
	if opts == nil || len(opts.IDs) == 0 {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "at least one association ID is required", "opts.IDs", opts, nil)
	}
	return nil, errors.NewNotImplementedError("bulk association deletion", "v0.0.40")
}