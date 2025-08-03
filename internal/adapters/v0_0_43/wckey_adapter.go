// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_43

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jontk/slurm-client/internal/common/types"
	"github.com/jontk/slurm-client/internal/managers/base"
	"github.com/jontk/slurm-client/pkg/errors"
	api "github.com/jontk/slurm-client/internal/api/v0_0_43"
)

// WCKeyAdapter implements the WCKeyAdapter interface for v0.0.43
type WCKeyAdapter struct {
	*base.BaseManager
	client       *api.ClientWithResponses
	errorAdapter *ErrorAdapter
}

// NewWCKeyAdapter creates a new WCKey adapter for v0.0.43
func NewWCKeyAdapter(client *api.ClientWithResponses) *WCKeyAdapter {
	return &WCKeyAdapter{
		BaseManager:  base.NewBaseManager("v0.0.43", "WCKey"),
		client:       client,
		errorAdapter: NewErrorAdapter(),
	}
}

// List retrieves a list of WCKeys with optional filtering
func (a *WCKeyAdapter) List(ctx context.Context, opts *types.WCKeyListOptions) (*types.WCKeyList, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Prepare parameters for the API call
	params := &api.SlurmdbV0043GetWckeysParams{}

	// Apply filters from options
	if opts != nil {
		if len(opts.Users) > 0 {
			users := strings.Join(opts.Users, ",")
			params.User = &users
		}
		if len(opts.Clusters) > 0 {
			clusters := strings.Join(opts.Clusters, ",")
			params.Cluster = &clusters
		}
		if opts.WithDeleted {
			withDeleted := "true"
			params.WithDeleted = &withDeleted
		}
		if opts.OnlyDefaults {
			onlyDefaults := "true"
			params.OnlyDefaults = &onlyDefaults
		}
	}

	// Call the API
	resp, err := a.client.SlurmdbV0043GetWckeysWithResponse(ctx, params)
	if err != nil {
		return nil, a.WrapError(err, "failed to list WCKeys")
	}

	// Handle API response with enhanced error handling
	if err := a.errorAdapter.HandleAPIResponse(resp.StatusCode(), resp.Body, "List WCKeys"); err != nil {
		return nil, err
	}

	if resp.JSON200 == nil {
		return &types.WCKeyList{WCKeys: []types.WCKey{}}, nil
	}

	// Convert API response to common types
	wckeys := make([]types.WCKey, 0)
	if resp.JSON200.Wckeys != nil {
		for _, apiWCKey := range resp.JSON200.Wckeys {
			wckey := a.convertAPIWCKeyToCommon(apiWCKey)
			wckeys = append(wckeys, wckey)
		}
	}

	return &types.WCKeyList{
		WCKeys: wckeys,
		Meta:   a.extractMeta(resp.JSON200.Meta),
	}, nil
}

// Get retrieves a specific WCKey by ID
func (a *WCKeyAdapter) Get(ctx context.Context, wcKeyID string) (*types.WCKey, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if err := a.ValidateResourceID(wcKeyID, "wcKeyID"); err != nil {
		return nil, err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Call the API
	resp, err := a.client.SlurmdbV0043GetWckeyWithResponse(ctx, wcKeyID)
	if err != nil {
		return nil, a.WrapError(err, fmt.Sprintf("failed to get WCKey %s", wcKeyID))
	}

	// Handle API response with enhanced error handling
	if err := a.errorAdapter.HandleAPIResponse(resp.StatusCode(), resp.Body, "Get WCKey"); err != nil {
		return nil, err
	}

	if resp.JSON200 == nil || resp.JSON200.Wckeys == nil || len(resp.JSON200.Wckeys) == 0 {
		return nil, fmt.Errorf("WCKey %s not found", wcKeyID)
	}

	// Return the first WCKey from the response
	wckey := a.convertAPIWCKeyToCommon(resp.JSON200.Wckeys[0])
	return &wckey, nil
}

// Create creates a new WCKey
func (a *WCKeyAdapter) Create(ctx context.Context, wckey *types.WCKeyCreate) (*types.WCKeyCreateResponse, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if err := a.validateWCKeyCreate(wckey); err != nil {
		return nil, err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}

	// Convert common WCKey to API format
	apiWCKey := a.convertCommonWCKeyCreateToAPI(wckey)

	// Create the request body
	apiReq := api.V0043OpenapiWckeyResp{
		Wckeys: &[]api.V0043Wckey{apiWCKey},
	}

	// Call the API
	resp, err := a.client.SlurmdbV0043PostWckeysWithResponse(ctx, apiReq)
	if err != nil {
		return nil, a.WrapError(err, "failed to create WCKey")
	}

	// Handle API response with enhanced error handling
	if err := a.errorAdapter.HandleAPIResponse(resp.StatusCode(), resp.Body, "Create WCKey"); err != nil {
		return nil, err
	}

	// Build response
	result := &types.WCKeyCreateResponse{
		Status:  "success",
		Message: fmt.Sprintf("WCKey '%s' created successfully", wckey.Name),
		Meta:    make(map[string]interface{}),
	}

	if resp.JSON200 != nil {
		result.Meta = a.extractMeta(resp.JSON200.Meta)

		// Extract the created WCKey ID if available
		if resp.JSON200.Wckeys != nil && len(*resp.JSON200.Wckeys) > 0 {
			firstWCKey := (*resp.JSON200.Wckeys)[0]
			if firstWCKey.Id != nil {
				result.ID = fmt.Sprintf("%d", *firstWCKey.Id)
			}
		}
	}

	return result, nil
}

// Delete deletes a WCKey
func (a *WCKeyAdapter) Delete(ctx context.Context, wcKeyID string) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.ValidateResourceID(wcKeyID, "wcKeyID"); err != nil {
		return err
	}

	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}

	// Call the API
	resp, err := a.client.SlurmdbV0043DeleteWckeyWithResponse(ctx, wcKeyID)
	if err != nil {
		return a.WrapError(err, fmt.Sprintf("failed to delete WCKey %s", wcKeyID))
	}

	// Handle API response with enhanced error handling
	if err := a.errorAdapter.HandleAPIResponse(resp.StatusCode(), resp.Body, "Delete WCKey"); err != nil {
		return err
	}

	return nil
}

// validateWCKeyCreate validates WCKey creation request
func (a *WCKeyAdapter) validateWCKeyCreate(wckey *types.WCKeyCreate) error {
	if wckey == nil {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "WCKey creation data is required", "wckey", nil, nil)
	}
	if wckey.Name == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "WCKey name is required", "name", wckey.Name, nil)
	}
	if wckey.Cluster == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "cluster is required for WCKey creation", "cluster", wckey.Cluster, nil)
	}
	return nil
}

// convertAPIWCKeyToCommon converts API WCKey to common type
func (a *WCKeyAdapter) convertAPIWCKeyToCommon(apiWCKey api.V0043Wckey) types.WCKey {
	wckey := types.WCKey{
		Meta: make(map[string]interface{}),
	}

	if apiWCKey.Id != nil {
		wckey.ID = fmt.Sprintf("%d", *apiWCKey.Id)
	}
	if apiWCKey.Name != nil {
		wckey.Name = *apiWCKey.Name
	}
	if apiWCKey.User != nil {
		wckey.User = *apiWCKey.User
	}
	if apiWCKey.Cluster != nil {
		wckey.Cluster = *apiWCKey.Cluster
	}

	// Set active status (WCKeys are typically active by default)
	wckey.Active = true

	return wckey
}

// convertCommonWCKeyCreateToAPI converts common WCKey create to API format
func (a *WCKeyAdapter) convertCommonWCKeyCreateToAPI(wckey *types.WCKeyCreate) api.V0043Wckey {
	apiWCKey := api.V0043Wckey{}

	if wckey.Name != "" {
		apiWCKey.Name = &wckey.Name
	}
	if wckey.User != "" {
		apiWCKey.User = &wckey.User
	}
	if wckey.Cluster != "" {
		apiWCKey.Cluster = &wckey.Cluster
	}

	return apiWCKey
}

// extractMeta safely extracts metadata from API response
func (a *WCKeyAdapter) extractMeta(meta *api.V0043OpenapiMeta) map[string]interface{} {
	result := make(map[string]interface{})

	if meta == nil {
		return result
	}

	// Extract basic metadata following existing patterns
	if meta.Client != nil {
		clientInfo := make(map[string]interface{})
		if meta.Client.Source != nil {
			clientInfo["source"] = *meta.Client.Source
		}
		if meta.Client.User != nil {
			clientInfo["user"] = *meta.Client.User
		}
		if meta.Client.Group != nil {
			clientInfo["group"] = *meta.Client.Group
		}
		if len(clientInfo) > 0 {
			result["client"] = clientInfo
		}
	}

	if meta.Plugin != nil {
		pluginInfo := make(map[string]interface{})
		if meta.Plugin.AccountingStorage != nil {
			pluginInfo["accounting_storage"] = *meta.Plugin.AccountingStorage
		}
		if len(pluginInfo) > 0 {
			result["plugin"] = pluginInfo
		}
	}

	return result
}