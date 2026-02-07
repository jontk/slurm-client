// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_40

import (
	"context"
	"fmt"
	"strings"

	types "github.com/jontk/slurm-client/api"
	adapterbase "github.com/jontk/slurm-client/internal/adapters/base"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_40"
)

// WCKeyAdapter implements the WCKeyAdapter interface for v0.0.40
type WCKeyAdapter struct {
	*adapterbase.BaseManager
	client *api.ClientWithResponses
}

// NewWCKeyAdapter creates a new WCKey adapter for v0.0.40
func NewWCKeyAdapter(client *api.ClientWithResponses) *WCKeyAdapter {
	return &WCKeyAdapter{
		BaseManager: adapterbase.NewBaseManager("v0.0.40", "WCKey"),
		client:      client,
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
	// Build parameters from options
	params := a.buildWCKeyListParams(opts)
	// Call the API
	resp, err := a.client.SlurmdbV0040GetWckeysWithResponse(ctx, params)
	if err != nil {
		return nil, a.WrapError(err, "failed to list WCKeys")
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
	wckeys := make([]types.WCKey, 0)
	if len(resp.JSON200.Wckeys) > 0 {
		for _, apiWCKey := range resp.JSON200.Wckeys {
			wckey := a.convertAPIWCKeyToCommon(apiWCKey)
			wckeys = append(wckeys, *wckey)
		}
	}
	return &types.WCKeyList{
		WCKeys: wckeys,
		Meta:   a.extractMeta(resp.JSON200.Meta),
	}, nil
}

// buildWCKeyListParams constructs API parameters from options
func (a *WCKeyAdapter) buildWCKeyListParams(opts *types.WCKeyListOptions) *api.SlurmdbV0040GetWckeysParams {
	params := &api.SlurmdbV0040GetWckeysParams{}
	if opts == nil {
		return params
	}
	// Build user filter
	if len(opts.Users) > 0 {
		params.User = a.buildCSVString(opts.Users)
	}
	// Build cluster filter
	if len(opts.Clusters) > 0 {
		params.Cluster = a.buildCSVString(opts.Clusters)
	}
	// Build name filter
	if len(opts.Names) > 0 {
		params.Name = a.buildCSVString(opts.Names)
	}
	// Set boolean flags
	if opts.OnlyDefaults {
		onlyDefaultsStr := "true"
		params.OnlyDefaults = &onlyDefaultsStr
	}
	if opts.WithDeleted {
		withDeletedStr := "true"
		params.WithDeleted = &withDeletedStr
	}
	return params
}

// buildCSVString constructs a comma-separated string from a slice
func (a *WCKeyAdapter) buildCSVString(items []string) *string {
	if len(items) == 0 {
		return nil
	}
	result := strings.Join(items, ",")
	return &result
}

// Get retrieves a specific WCKey by ID
func (a *WCKeyAdapter) Get(ctx context.Context, wcKeyID string) (*types.WCKey, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}
	// Call the API
	resp, err := a.client.SlurmdbV0040GetWckeyWithResponse(ctx, wcKeyID)
	if err != nil {
		return nil, a.WrapError(err, "failed to get WCKey "+wcKeyID)
	}
	// Check response status
	if resp.StatusCode() != 200 {
		return nil, a.HandleAPIError(fmt.Errorf("API error: status %d", resp.StatusCode()))
	}
	// Check for API response
	if resp.JSON200 == nil || len(resp.JSON200.Wckeys) == 0 {
		return nil, fmt.Errorf("WCKey %s not found", wcKeyID)
	}
	// Convert the first WCKey in the response
	wckeys := resp.JSON200.Wckeys
	return a.convertAPIWCKeyToCommon(wckeys[0]), nil
}

// Create creates a new WCKey
func (a *WCKeyAdapter) Create(ctx context.Context, wckey *types.WCKeyCreate) (*types.WCKeyCreateResponse, error) {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}
	// Convert common WCKey create to API request
	apiWCKey := &api.V0040Wckey{
		Name:    wckey.Name,
		Cluster: wckey.Cluster,
	}
	if wckey.User != "" {
		apiWCKey.User = wckey.User
	}
	// Create request body
	apiReq := api.V0040OpenapiWckeyResp{
		Wckeys: []api.V0040Wckey{*apiWCKey},
	}
	// Call the API
	params := &api.SlurmdbV0040PostWckeysParams{}
	resp, err := a.client.SlurmdbV0040PostWckeysWithResponse(ctx, params, apiReq)
	if err != nil {
		return nil, a.WrapError(err, "failed to create WCKey")
	}
	// Check response status
	if resp.StatusCode() != 200 {
		return nil, a.HandleAPIError(fmt.Errorf("API error: status %d", resp.StatusCode()))
	}
	// Check for API response
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response from API")
	}
	// Convert response - v0.0.40 POST returns V0040OpenapiResp not V0040OpenapiWckeyResp
	return a.convertAPIWCKeyCreateResponseToCommon(resp.JSON200, wckey.Name)
}

// Delete deletes a WCKey
func (a *WCKeyAdapter) Delete(ctx context.Context, wcKeyID string) error {
	// Use base validation
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	// Check client initialization
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}
	// Call the API
	resp, err := a.client.SlurmdbV0040DeleteWckeyWithResponse(ctx, wcKeyID)
	if err != nil {
		return a.WrapError(err, "failed to delete WCKey "+wcKeyID)
	}
	// Check response status
	if resp.StatusCode() != 200 {
		return a.HandleAPIError(fmt.Errorf("API error: status %d", resp.StatusCode()))
	}
	return nil
}

// convertAPIWCKeyCreateResponseToCommon converts API create response to common type
func (a *WCKeyAdapter) convertAPIWCKeyCreateResponseToCommon(apiResp *api.V0040OpenapiResp, name string) (*types.WCKeyCreateResponse, error) {
	resp := &types.WCKeyCreateResponse{
		Status: "success",
		Meta:   make(map[string]interface{}),
	}
	// V0040OpenapiResp doesn't contain WCKeys - it's a general response
	// We cannot extract ID from the response for v0.0.40
	// Extract metadata if available
	if apiResp.Meta != nil {
		resp.Meta = a.extractMeta(apiResp.Meta)
	}
	// Handle errors in response - V0040OpenapiErrors is []V0040OpenapiError
	if apiResp.Errors != nil && len(*apiResp.Errors) > 0 {
		resp.Status = "error"
		errors := *apiResp.Errors
		if len(errors) > 0 && errors[0].Error != nil {
			resp.Message = *errors[0].Error
		} else {
			resp.Message = "WCKey creation failed"
		}
	} else {
		resp.Message = fmt.Sprintf("WCKey '%s' created successfully", name)
	}
	return resp, nil
}

// extractMeta safely extracts metadata from API response
func (a *WCKeyAdapter) extractMeta(meta *api.V0040OpenapiMeta) map[string]interface{} {
	result := make(map[string]interface{})
	if meta == nil {
		return result
	}
	// Extract basic metadata
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
