// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_43

import (
	"context"
	stderrors "errors"
	"strconv"
	"strings"
	"time"

	"github.com/jontk/slurm-client/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
)

// AssociationManagerImpl implements the AssociationManager interface for v0.0.43
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
	if a.client == nil || a.client.apiClient == nil || a.client.apiClient.ClientInterface == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Prepare parameters for the API call
	params := &SlurmdbV0043GetAssociationsParams{}

	// Apply filtering options if provided
	if opts != nil {
		if len(opts.Users) > 0 {
			userStr := strings.Join(opts.Users, ",")
			params.User = &userStr
		}
		if len(opts.Accounts) > 0 {
			accountStr := strings.Join(opts.Accounts, ",")
			params.Account = &accountStr
		}
		if len(opts.Clusters) > 0 {
			clusterStr := strings.Join(opts.Clusters, ",")
			params.Cluster = &clusterStr
		}
		if len(opts.Partitions) > 0 {
			partitionStr := strings.Join(opts.Partitions, ",")
			params.Partition = &partitionStr
		}
		if len(opts.ParentAccounts) > 0 {
			parentStr := strings.Join(opts.ParentAccounts, ",")
			params.ParentAccount = &parentStr
		}
		if len(opts.QoS) > 0 {
			qosStr := strings.Join(opts.QoS, ",")
			params.Qos = &qosStr
		}
		if opts.WithDeleted {
			withDeleted := "yes"
			params.IncludeDeletedAssociations = &withDeleted
		}
		if opts.WithUsage {
			withUsage := "yes"
			params.IncludeUsage = &withUsage
		}
		if opts.WithSubAccounts {
			withSubAcct := "yes"
			params.IncludeSubAcctInformation = &withSubAcct
		}
		if opts.OnlyDefaults {
			onlyDefaults := "yes"
			params.FilterToOnlyDefaults = &onlyDefaults
		}
	}

	// Call the generated OpenAPI client
	resp, err := a.client.apiClient.SlurmdbV0043GetAssociationsWithResponse(ctx, params)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
	}

	// Check HTTP status and handle API errors
	if resp.StatusCode() != 200 {
		return nil, handleAPIError(resp.StatusCode(), resp.JSON200)
	}

	// Check for unexpected response format
	if resp.JSON200 == nil {
		return nil, errors.NewClientError(errors.ErrorCodeServerInternal, "Unexpected response format", "Expected JSON response with associations but got nil")
	}

	// Convert the response to our interface types
	associations := make([]*interfaces.Association, 0, len(resp.JSON200.Associations))
	for _, apiAssoc := range resp.JSON200.Associations {
		assoc, err := convertAPIAssociationToInterface(&apiAssoc)
		if err != nil {
			conversionErr := errors.NewClientError(errors.ErrorCodeServerInternal, "Failed to convert association data")
			conversionErr.Cause = err
			conversionErr.Details = "Error converting association data"
			return nil, conversionErr
		}
		associations = append(associations, assoc)
	}

	return &interfaces.AssociationList{
		Associations: associations,
		Total:        len(associations),
		Metadata:     extractMetadata(resp.JSON200.Meta),
	}, nil
}

// Get retrieves a specific association
func (a *AssociationManagerImpl) Get(ctx context.Context, opts *interfaces.GetAssociationOptions) (*interfaces.Association, error) {
	// Validate input first (cheap check)
	if opts == nil || opts.User == "" || opts.Account == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "User and Account are required", "opts", opts, nil)
	}

	// Then check client initialization
	if a.client == nil || a.client.apiClient == nil || a.client.apiClient.ClientInterface == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Prepare parameters for the API call
	params := &SlurmdbV0043GetAssociationParams{
		User:    &opts.User,
		Account: &opts.Account,
	}

	if opts.Cluster != "" {
		params.Cluster = &opts.Cluster
	}
	if opts.Partition != "" {
		params.Partition = &opts.Partition
	}
	if opts.WithUsage {
		withUsage := "yes"
		params.IncludeUsage = &withUsage
	}

	// Call the generated OpenAPI client
	resp, err := a.client.apiClient.SlurmdbV0043GetAssociationWithResponse(ctx, params)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
	}

	// Check HTTP status and handle API errors
	if resp.StatusCode() != 200 {
		return nil, handleAPIError(resp.StatusCode(), resp.JSON200)
	}

	// Check for unexpected response format
	if resp.JSON200 == nil || len(resp.JSON200.Associations) == 0 {
		return nil, errors.NewClientError(errors.ErrorCodeResourceNotFound, "Association not found")
	}

	// Convert the first association from the response
	assoc, err := convertAPIAssociationToInterface(&resp.JSON200.Associations[0])
	if err != nil {
		conversionErr := errors.NewClientError(errors.ErrorCodeServerInternal, "Failed to convert association data")
		conversionErr.Cause = err
		return nil, conversionErr
	}

	return assoc, nil
}

// Create creates new associations
func (a *AssociationManagerImpl) Create(ctx context.Context, associations []*interfaces.AssociationCreate) (*interfaces.AssociationCreateResponse, error) {
	// Validate input first (cheap check)
	if len(associations) == 0 {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "At least one association is required", "associations", associations, nil)
	}

	// Then check client initialization
	if a.client == nil || a.client.apiClient == nil || a.client.apiClient.ClientInterface == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Convert interface types to API types
	apiAssocs := make([]V0043Assoc, len(associations))
	for i, assoc := range associations {
		apiAssoc, err := convertInterfaceAssociationCreateToAPI(assoc)
		if err != nil {
			return nil, err
		}
		apiAssocs[i] = *apiAssoc
	}

	// Prepare request body
	reqBody := SlurmdbV0043PostAssociationsJSONRequestBody{
		Associations: apiAssocs,
	}

	// Call the generated OpenAPI client
	resp, err := a.client.apiClient.SlurmdbV0043PostAssociationsWithResponse(ctx, reqBody)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
	}

	// Check HTTP status and handle API errors
	if resp.StatusCode() != 200 {
		// For POST operations, we need to handle errors differently
		if resp.JSON200 != nil && resp.JSON200.Errors != nil {
			var apiErrors []errors.SlurmAPIErrorDetail
			for _, apiErr := range *resp.JSON200.Errors {
				var errorNumber int
				if apiErr.ErrorNumber != nil {
					errorNumber = int(*apiErr.ErrorNumber)
				}
				var errorCode string
				if apiErr.Error != nil {
					errorCode = *apiErr.Error
				}
				var source string
				if apiErr.Source != nil {
					source = *apiErr.Source
				}
				var description string
				if apiErr.Description != nil {
					description = *apiErr.Description
				}

				apiErrors = append(apiErrors, errors.SlurmAPIErrorDetail{
					ErrorNumber: errorNumber,
					ErrorCode:   errorCode,
					Source:      source,
					Description: description,
				})
			}
			if len(apiErrors) > 0 {
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.43", apiErrors)
				return nil, apiError.SlurmError
			}
		}
		return nil, errors.WrapHTTPError(resp.StatusCode(), nil, "v0.0.43")
	}

	// Build response
	createResp := &interfaces.AssociationCreateResponse{
		Created:  0,
		Updated:  0,
		Errors:   []string{},
		Warnings: []string{},
		Metadata: extractMetadata(resp.JSON200.Meta),
	}

	// Extract errors and warnings from response
	if resp.JSON200.Errors != nil {
		for _, apiErr := range *resp.JSON200.Errors {
			if apiErr.Description != nil {
				createResp.Errors = append(createResp.Errors, *apiErr.Description)
			}
		}
	}

	if resp.JSON200.Warnings != nil {
		for _, warn := range *resp.JSON200.Warnings {
			if warn.Description != nil {
				createResp.Warnings = append(createResp.Warnings, *warn.Description)
			}
		}
	}

	// Note: The API doesn't return created associations in the response
	// We can only confirm the operation succeeded
	if resp.StatusCode() == 200 {
		createResp.Created = len(associations)
	}

	return createResp, nil
}

// Update updates existing associations
func (a *AssociationManagerImpl) Update(ctx context.Context, associations []*interfaces.AssociationUpdate) error {
	// Validate input first (cheap check)
	if len(associations) == 0 {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "At least one association update is required", "associations", associations, nil)
	}

	// Then check client initialization
	if a.client == nil || a.client.apiClient == nil || a.client.apiClient.ClientInterface == nil {
		return errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Convert interface types to API types
	apiAssocs := make([]V0043Assoc, len(associations))
	for i, assoc := range associations {
		apiAssoc, err := convertInterfaceAssociationUpdateToAPI(assoc)
		if err != nil {
			return err
		}
		apiAssocs[i] = *apiAssoc
	}

	// Prepare request body
	reqBody := SlurmdbV0043PostAssociationsJSONRequestBody{
		Associations: apiAssocs,
	}

	// Call the generated OpenAPI client for update (POST is used for both create and update in SLURM)
	resp, err := a.client.apiClient.SlurmdbV0043PostAssociationsWithResponse(ctx, reqBody)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
	}

	// Check HTTP status and handle API errors
	if resp.StatusCode() != 200 {
		// For POST operations, we need to handle errors differently
		if resp.JSON200 != nil && resp.JSON200.Errors != nil {
			var apiErrors []errors.SlurmAPIErrorDetail
			for _, apiErr := range *resp.JSON200.Errors {
				var errorNumber int
				if apiErr.ErrorNumber != nil {
					errorNumber = int(*apiErr.ErrorNumber)
				}
				var errorCode string
				if apiErr.Error != nil {
					errorCode = *apiErr.Error
				}
				var source string
				if apiErr.Source != nil {
					source = *apiErr.Source
				}
				var description string
				if apiErr.Description != nil {
					description = *apiErr.Description
				}

				apiErrors = append(apiErrors, errors.SlurmAPIErrorDetail{
					ErrorNumber: errorNumber,
					ErrorCode:   errorCode,
					Source:      source,
					Description: description,
				})
			}
			if len(apiErrors) > 0 {
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.43", apiErrors)
				return apiError.SlurmError
			}
		}
		return errors.WrapHTTPError(resp.StatusCode(), nil, "v0.0.43")
	}

	return nil
}

// Delete deletes a single association
func (a *AssociationManagerImpl) Delete(ctx context.Context, opts *interfaces.DeleteAssociationOptions) error {
	// Validate input first (cheap check)
	if opts == nil || opts.User == "" || opts.Account == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "User and Account are required", "opts", opts, nil)
	}

	// Then check client initialization
	if a.client == nil || a.client.apiClient == nil || a.client.apiClient.ClientInterface == nil {
		return errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Prepare parameters for the API call
	params := &SlurmdbV0043DeleteAssociationParams{
		User:    &opts.User,
		Account: &opts.Account,
	}

	if opts.Cluster != "" {
		params.Cluster = &opts.Cluster
	}
	if opts.Partition != "" {
		params.Partition = &opts.Partition
	}

	// Call the generated OpenAPI client
	resp, err := a.client.apiClient.SlurmdbV0043DeleteAssociationWithResponse(ctx, params)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
	}

	// Check HTTP status and handle API errors
	if resp.StatusCode() != 200 {
		return handleAPIErrorForDelete(resp.StatusCode(), resp.JSON200)
	}

	return nil
}

// BulkDelete deletes multiple associations
func (a *AssociationManagerImpl) BulkDelete(ctx context.Context, opts *interfaces.BulkDeleteOptions) (*interfaces.BulkDeleteResponse, error) {
	// Validate input first (cheap check)
	if opts == nil {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "BulkDeleteOptions is required", "opts", opts, nil)
	}

	// Then check client initialization
	if a.client == nil || a.client.apiClient == nil || a.client.apiClient.ClientInterface == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Prepare parameters for the API call
	params := &SlurmdbV0043DeleteAssociationsParams{}

	if len(opts.Users) > 0 {
		userStr := strings.Join(opts.Users, ",")
		params.User = &userStr
	}
	if len(opts.Accounts) > 0 {
		accountStr := strings.Join(opts.Accounts, ",")
		params.Account = &accountStr
	}
	if len(opts.Clusters) > 0 {
		clusterStr := strings.Join(opts.Clusters, ",")
		params.Cluster = &clusterStr
	}
	if len(opts.Partitions) > 0 {
		partitionStr := strings.Join(opts.Partitions, ",")
		params.Partition = &partitionStr
	}
	// Note: OnlyIfIdle doesn't seem to be supported in the API params

	// Call the generated OpenAPI client
	resp, err := a.client.apiClient.SlurmdbV0043DeleteAssociationsWithResponse(ctx, params)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
	}

	// Check HTTP status and handle API errors
	if resp.StatusCode() != 200 {
		return nil, handleAPIErrorForDelete(resp.StatusCode(), resp.JSON200)
	}

	// Build response
	deleteResp := &interfaces.BulkDeleteResponse{
		Deleted: 0,
		Failed:  0,
		Errors:  []string{},
	}

	// Extract errors from response
	if resp.JSON200 != nil && resp.JSON200.Errors != nil {
		for _, apiErr := range *resp.JSON200.Errors {
			if apiErr.Description != nil {
				deleteResp.Errors = append(deleteResp.Errors, *apiErr.Description)
			}
		}
		deleteResp.Failed = len(*resp.JSON200.Errors)
	}

	// Extract deleted associations count from response
	if resp.JSON200 != nil {
		deleteResp.Deleted = len(resp.JSON200.RemovedAssociations)
	}

	return deleteResp, nil
}

// GetUserAssociations retrieves all associations for a specific user
func (a *AssociationManagerImpl) GetUserAssociations(ctx context.Context, userName string) ([]*interfaces.Association, error) {
	if userName == "" {
		return nil, errors.NewClientError(errors.ErrorCodeValidationFailed, "User name is required")
	}

	opts := &interfaces.ListAssociationsOptions{
		Users:    []string{userName},
		WithTRES: true,
	}

	result, err := a.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return result.Associations, nil
}

// GetAccountAssociations retrieves all associations for a specific account
func (a *AssociationManagerImpl) GetAccountAssociations(ctx context.Context, accountName string) ([]*interfaces.Association, error) {
	if accountName == "" {
		return nil, errors.NewClientError(errors.ErrorCodeValidationFailed, "Account name is required")
	}

	opts := &interfaces.ListAssociationsOptions{
		Accounts: []string{accountName},
		WithTRES: true,
	}

	result, err := a.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return result.Associations, nil
}

// ValidateAssociation checks if a user-account-cluster association exists and is valid
func (a *AssociationManagerImpl) ValidateAssociation(ctx context.Context, user, account, cluster string) (bool, error) {
	if user == "" || account == "" {
		return false, errors.NewClientError(errors.ErrorCodeValidationFailed, "User and Account are required")
	}

	opts := &interfaces.GetAssociationOptions{
		User:    user,
		Account: account,
		Cluster: cluster,
	}

	assoc, err := a.Get(ctx, opts)
	if err != nil {
		// If association not found, it's not valid
		var clientErr *errors.SlurmError
		if stderrors.As(err, &clientErr) && clientErr.Code == errors.ErrorCodeResourceNotFound {
			return false, nil
		}
		return false, err
	}

	// Check if association is deleted
	if assoc.Deleted != nil {
		return false, nil
	}

	return true, nil
}

// Helper function to handle API errors
func handleAPIError(statusCode int, response *V0043OpenapiAssocsResp) error {
	// Extract errors from response if available
	var apiErrors []errors.SlurmAPIErrorDetail

	// Check if response has errors
	if response != nil && response.Errors != nil {
		for _, apiErr := range *response.Errors {
			var errorNumber int
			if apiErr.ErrorNumber != nil {
				errorNumber = int(*apiErr.ErrorNumber)
			}
			var errorCode string
			if apiErr.Error != nil {
				errorCode = *apiErr.Error
			}
			var source string
			if apiErr.Source != nil {
				source = *apiErr.Source
			}
			var description string
			if apiErr.Description != nil {
				description = *apiErr.Description
			}

			apiErrors = append(apiErrors, errors.SlurmAPIErrorDetail{
				ErrorNumber: errorNumber,
				ErrorCode:   errorCode,
				Source:      source,
				Description: description,
			})
		}
	}

	if len(apiErrors) > 0 {
		apiError := errors.NewSlurmAPIError(statusCode, "v0.0.43", apiErrors)
		return apiError.SlurmError
	}

	// Fall back to HTTP error handling
	return errors.WrapHTTPError(statusCode, nil, "v0.0.43")
}

// Helper function to handle API errors for delete operations
func handleAPIErrorForDelete(statusCode int, response *V0043OpenapiAssocsRemovedResp) error {
	// Extract errors from response if available
	var apiErrors []errors.SlurmAPIErrorDetail

	// Check if response has errors
	if response != nil && response.Errors != nil {
		for _, apiErr := range *response.Errors {
			var errorNumber int
			if apiErr.ErrorNumber != nil {
				errorNumber = int(*apiErr.ErrorNumber)
			}
			var errorCode string
			if apiErr.Error != nil {
				errorCode = *apiErr.Error
			}
			var source string
			if apiErr.Source != nil {
				source = *apiErr.Source
			}
			var description string
			if apiErr.Description != nil {
				description = *apiErr.Description
			}

			apiErrors = append(apiErrors, errors.SlurmAPIErrorDetail{
				ErrorNumber: errorNumber,
				ErrorCode:   errorCode,
				Source:      source,
				Description: description,
			})
		}
	}

	if len(apiErrors) > 0 {
		apiError := errors.NewSlurmAPIError(statusCode, "v0.0.43", apiErrors)
		return apiError.SlurmError
	}

	// Fall back to HTTP error handling
	return errors.WrapHTTPError(statusCode, nil, "v0.0.43")
}

// Helper function to extract metadata from API response
func extractMetadata(meta interface{}) map[string]interface{} {
	if meta == nil {
		return nil
	}
	// Convert meta to map if needed
	// This implementation depends on the actual structure of the Meta field
	return nil
}

// Conversion helper functions

func convertAPIAssociationToInterface(apiAssoc *V0043Assoc) (*interfaces.Association, error) {
	if apiAssoc == nil {
		return nil, errors.NewClientError(errors.ErrorCodeServerInternal, "API association is nil")
	}

	assoc := &interfaces.Association{
		Metadata: make(map[string]interface{}),
	}

	// Map basic fields
	if apiAssoc.Id != nil {
		assoc.ID = uint32(*apiAssoc.Id)
	}
	// User is a required string field
	assoc.User = apiAssoc.User
	if apiAssoc.Account != nil {
		assoc.Account = *apiAssoc.Account
	}
	if apiAssoc.Cluster != nil {
		assoc.Cluster = *apiAssoc.Cluster
	}
	if apiAssoc.Partition != nil {
		assoc.Partition = *apiAssoc.Partition
	}
	if apiAssoc.ParentAccount != nil {
		assoc.ParentAccount = *apiAssoc.ParentAccount
	}
	if apiAssoc.IsDefault != nil {
		assoc.IsDefault = *apiAssoc.IsDefault
	}
	if apiAssoc.Comment != nil {
		assoc.Comment = *apiAssoc.Comment
	}

	// Map resource limits
	if apiAssoc.SharesRaw != nil {
		assoc.SharesRaw = int(*apiAssoc.SharesRaw)
	}
	if apiAssoc.Priority != nil && apiAssoc.Priority.Number != nil {
		assoc.Priority = uint32(*apiAssoc.Priority.Number)
	}

	// Map job limits
	if apiAssoc.Max != nil && apiAssoc.Max.Jobs != nil && apiAssoc.Max.Jobs.Per != nil && apiAssoc.Max.Jobs.Per.Count != nil && apiAssoc.Max.Jobs.Per.Count.Number != nil {
		val := int(*apiAssoc.Max.Jobs.Per.Count.Number)
		assoc.MaxJobs = &val
	}
	if apiAssoc.Max != nil && apiAssoc.Max.Jobs != nil && apiAssoc.Max.Jobs.Accruing != nil && apiAssoc.Max.Jobs.Accruing.Number != nil {
		val := int(*apiAssoc.Max.Jobs.Accruing.Number)
		assoc.MaxJobsAccrue = &val
	}
	if apiAssoc.Max != nil && apiAssoc.Max.Jobs != nil && apiAssoc.Max.Jobs.Total != nil && apiAssoc.Max.Jobs.Total.Number != nil {
		val := int(*apiAssoc.Max.Jobs.Total.Number)
		assoc.MaxSubmitJobs = &val
	}

	// Map TRES limits
	assoc.MaxTRESPerJob = make(map[string]string)
	assoc.MaxTRESMins = make(map[string]string)
	assoc.GrpTRES = make(map[string]string)
	assoc.GrpTRESMins = make(map[string]string)
	assoc.GrpTRESRunMins = make(map[string]string)

	if apiAssoc.Max != nil && apiAssoc.Max.Tres != nil && apiAssoc.Max.Tres.Total != nil {
		for _, tres := range *apiAssoc.Max.Tres.Total {
			if tres.Type != "" && tres.Count != nil {
				assoc.MaxTRESPerJob[tres.Type] = strconv.FormatInt(*tres.Count, 10)
			}
		}
	}

	// Map QoS
	if apiAssoc.Default != nil && apiAssoc.Default.Qos != nil {
		assoc.DefaultQoS = *apiAssoc.Default.Qos
	}
	if apiAssoc.Qos != nil {
		assoc.QoSList = *apiAssoc.Qos
	}

	// Map flags
	if apiAssoc.Flags != nil {
		assoc.Flags = make([]string, len(*apiAssoc.Flags))
		for i, flag := range *apiAssoc.Flags {
			assoc.Flags[i] = string(flag)
		}
	}

	// Set timestamps (these would need to be extracted from metadata or other fields)
	assoc.Created = time.Now()  // Placeholder
	assoc.Modified = time.Now() // Placeholder

	return assoc, nil
}

func convertInterfaceAssociationCreateToAPI(assoc *interfaces.AssociationCreate) (*V0043Assoc, error) {
	if assoc == nil {
		return nil, errors.NewClientError(errors.ErrorCodeValidationFailed, "Association create request is nil")
	}

	apiAssoc := &V0043Assoc{
		User:    assoc.User,
		Account: &assoc.Account,
	}

	if assoc.Cluster != "" {
		apiAssoc.Cluster = &assoc.Cluster
	}
	if assoc.Partition != "" {
		apiAssoc.Partition = &assoc.Partition
	}
	if assoc.ParentAccount != "" {
		apiAssoc.ParentAccount = &assoc.ParentAccount
	}
	apiAssoc.IsDefault = &assoc.IsDefault
	if assoc.Comment != "" {
		apiAssoc.Comment = &assoc.Comment
	}

	// Map resource limits
	if assoc.SharesRaw != nil {
		shares := int32(*assoc.SharesRaw)
		apiAssoc.SharesRaw = &shares
	}
	if assoc.Priority != nil {
		prio := int32(*assoc.Priority)
		apiAssoc.Priority = &V0043Uint32NoValStruct{Number: &prio}
	}

	// Initialize Max structure if needed
	if assoc.MaxJobs != nil || assoc.MaxJobsAccrue != nil || assoc.MaxSubmitJobs != nil {
		apiAssoc.Max = &struct {
			Jobs *struct {
				Accruing *V0043Uint32NoValStruct `json:"accruing,omitempty"`
				Active   *V0043Uint32NoValStruct `json:"active,omitempty"`
				Per      *struct {
					Accruing  *V0043Uint32NoValStruct `json:"accruing,omitempty"`
					Count     *V0043Uint32NoValStruct `json:"count,omitempty"`
					Submitted *V0043Uint32NoValStruct `json:"submitted,omitempty"`
					WallClock *V0043Uint32NoValStruct `json:"wall_clock,omitempty"`
				} `json:"per,omitempty"`
				Total *V0043Uint32NoValStruct `json:"total,omitempty"`
			} `json:"jobs,omitempty"`
			Per *struct {
				Account *struct {
					WallClock *V0043Uint32NoValStruct `json:"wall_clock,omitempty"`
				} `json:"account,omitempty"`
			} `json:"per,omitempty"`
			Tres *struct {
				Group *struct {
					Active  *V0043TresList `json:"active,omitempty"`
					Minutes *V0043TresList `json:"minutes,omitempty"`
				} `json:"group,omitempty"`
				Minutes *struct {
					Per *struct {
						Job *V0043TresList `json:"job,omitempty"`
					} `json:"per,omitempty"`
					Total *V0043TresList `json:"total,omitempty"`
				} `json:"minutes,omitempty"`
				Per *struct {
					Job  *V0043TresList `json:"job,omitempty"`
					Node *V0043TresList `json:"node,omitempty"`
				} `json:"per,omitempty"`
				Total *V0043TresList `json:"total,omitempty"`
			} `json:"tres,omitempty"`
		}{}
		apiAssoc.Max.Jobs = &struct {
			Accruing *V0043Uint32NoValStruct `json:"accruing,omitempty"`
			Active   *V0043Uint32NoValStruct `json:"active,omitempty"`
			Per      *struct {
				Accruing  *V0043Uint32NoValStruct `json:"accruing,omitempty"`
				Count     *V0043Uint32NoValStruct `json:"count,omitempty"`
				Submitted *V0043Uint32NoValStruct `json:"submitted,omitempty"`
				WallClock *V0043Uint32NoValStruct `json:"wall_clock,omitempty"`
			} `json:"per,omitempty"`
			Total *V0043Uint32NoValStruct `json:"total,omitempty"`
		}{}

		if assoc.MaxJobs != nil {
			jobNum := int32(*assoc.MaxJobs)
			if apiAssoc.Max.Jobs.Per == nil {
				apiAssoc.Max.Jobs.Per = &struct {
					Accruing  *V0043Uint32NoValStruct `json:"accruing,omitempty"`
					Count     *V0043Uint32NoValStruct `json:"count,omitempty"`
					Submitted *V0043Uint32NoValStruct `json:"submitted,omitempty"`
					WallClock *V0043Uint32NoValStruct `json:"wall_clock,omitempty"`
				}{}
			}
			apiAssoc.Max.Jobs.Per.Count = &V0043Uint32NoValStruct{Number: &jobNum}
		}
		if assoc.MaxJobsAccrue != nil {
			jobAccrue := int32(*assoc.MaxJobsAccrue)
			apiAssoc.Max.Jobs.Accruing = &V0043Uint32NoValStruct{Number: &jobAccrue}
		}
		if assoc.MaxSubmitJobs != nil {
			jobSubmit := int32(*assoc.MaxSubmitJobs)
			apiAssoc.Max.Jobs.Total = &V0043Uint32NoValStruct{Number: &jobSubmit}
		}
	}

	// Map QoS
	if assoc.DefaultQoS != "" {
		if apiAssoc.Default == nil {
			apiAssoc.Default = &struct {
				Qos *string `json:"qos,omitempty"`
			}{}
		}
		apiAssoc.Default.Qos = &assoc.DefaultQoS
	}
	if len(assoc.QoSList) > 0 {
		apiAssoc.Qos = &assoc.QoSList
	}

	// Map flags
	if len(assoc.Flags) > 0 {
		flags := make([]V0043AssocFlags, len(assoc.Flags))
		for i, flag := range assoc.Flags {
			flags[i] = V0043AssocFlags(flag)
		}
		apiAssoc.Flags = &flags
	}

	return apiAssoc, nil
}

func convertInterfaceAssociationUpdateToAPI(assoc *interfaces.AssociationUpdate) (*V0043Assoc, error) {
	if assoc == nil {
		return nil, errors.NewClientError(errors.ErrorCodeValidationFailed, "Association update request is nil")
	}

	apiAssoc := &V0043Assoc{
		User:    assoc.User,
		Account: &assoc.Account,
	}

	if assoc.Cluster != "" {
		apiAssoc.Cluster = &assoc.Cluster
	}
	if assoc.Partition != "" {
		apiAssoc.Partition = &assoc.Partition
	}

	// Map updateable fields
	if assoc.IsDefault != nil {
		apiAssoc.IsDefault = assoc.IsDefault
	}
	if assoc.Comment != nil {
		apiAssoc.Comment = assoc.Comment
	}

	// Map resource limits
	if assoc.SharesRaw != nil {
		shares := int32(*assoc.SharesRaw)
		apiAssoc.SharesRaw = &shares
	}
	if assoc.Priority != nil {
		prio := int32(*assoc.Priority)
		apiAssoc.Priority = &V0043Uint32NoValStruct{Number: &prio}
	}

	// Initialize Max structure if needed
	if assoc.MaxJobs != nil || assoc.MaxJobsAccrue != nil || assoc.MaxSubmitJobs != nil {
		apiAssoc.Max = &struct {
			Jobs *struct {
				Accruing *V0043Uint32NoValStruct `json:"accruing,omitempty"`
				Active   *V0043Uint32NoValStruct `json:"active,omitempty"`
				Per      *struct {
					Accruing  *V0043Uint32NoValStruct `json:"accruing,omitempty"`
					Count     *V0043Uint32NoValStruct `json:"count,omitempty"`
					Submitted *V0043Uint32NoValStruct `json:"submitted,omitempty"`
					WallClock *V0043Uint32NoValStruct `json:"wall_clock,omitempty"`
				} `json:"per,omitempty"`
				Total *V0043Uint32NoValStruct `json:"total,omitempty"`
			} `json:"jobs,omitempty"`
			Per *struct {
				Account *struct {
					WallClock *V0043Uint32NoValStruct `json:"wall_clock,omitempty"`
				} `json:"account,omitempty"`
			} `json:"per,omitempty"`
			Tres *struct {
				Group *struct {
					Active  *V0043TresList `json:"active,omitempty"`
					Minutes *V0043TresList `json:"minutes,omitempty"`
				} `json:"group,omitempty"`
				Minutes *struct {
					Per *struct {
						Job *V0043TresList `json:"job,omitempty"`
					} `json:"per,omitempty"`
					Total *V0043TresList `json:"total,omitempty"`
				} `json:"minutes,omitempty"`
				Per *struct {
					Job  *V0043TresList `json:"job,omitempty"`
					Node *V0043TresList `json:"node,omitempty"`
				} `json:"per,omitempty"`
				Total *V0043TresList `json:"total,omitempty"`
			} `json:"tres,omitempty"`
		}{}
		apiAssoc.Max.Jobs = &struct {
			Accruing *V0043Uint32NoValStruct `json:"accruing,omitempty"`
			Active   *V0043Uint32NoValStruct `json:"active,omitempty"`
			Per      *struct {
				Accruing  *V0043Uint32NoValStruct `json:"accruing,omitempty"`
				Count     *V0043Uint32NoValStruct `json:"count,omitempty"`
				Submitted *V0043Uint32NoValStruct `json:"submitted,omitempty"`
				WallClock *V0043Uint32NoValStruct `json:"wall_clock,omitempty"`
			} `json:"per,omitempty"`
			Total *V0043Uint32NoValStruct `json:"total,omitempty"`
		}{}

		if assoc.MaxJobs != nil {
			jobNum := int32(*assoc.MaxJobs)
			if apiAssoc.Max.Jobs.Per == nil {
				apiAssoc.Max.Jobs.Per = &struct {
					Accruing  *V0043Uint32NoValStruct `json:"accruing,omitempty"`
					Count     *V0043Uint32NoValStruct `json:"count,omitempty"`
					Submitted *V0043Uint32NoValStruct `json:"submitted,omitempty"`
					WallClock *V0043Uint32NoValStruct `json:"wall_clock,omitempty"`
				}{}
			}
			apiAssoc.Max.Jobs.Per.Count = &V0043Uint32NoValStruct{Number: &jobNum}
		}
		if assoc.MaxJobsAccrue != nil {
			jobAccrue := int32(*assoc.MaxJobsAccrue)
			apiAssoc.Max.Jobs.Accruing = &V0043Uint32NoValStruct{Number: &jobAccrue}
		}
		if assoc.MaxSubmitJobs != nil {
			jobSubmit := int32(*assoc.MaxSubmitJobs)
			apiAssoc.Max.Jobs.Total = &V0043Uint32NoValStruct{Number: &jobSubmit}
		}
	}

	// Map QoS
	if assoc.DefaultQoS != nil {
		if apiAssoc.Default == nil {
			apiAssoc.Default = &struct {
				Qos *string `json:"qos,omitempty"`
			}{}
		}
		apiAssoc.Default.Qos = assoc.DefaultQoS
	}
	if len(assoc.QoSList) > 0 {
		apiAssoc.Qos = &assoc.QoSList
	}

	// Map flags
	if len(assoc.Flags) > 0 {
		flags := make([]V0043AssocFlags, len(assoc.Flags))
		for i, flag := range assoc.Flags {
			flags[i] = V0043AssocFlags(flag)
		}
		apiAssoc.Flags = &flags
	}

	return apiAssoc, nil
}
