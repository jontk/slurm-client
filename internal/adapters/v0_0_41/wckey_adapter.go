// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_41

import (
	"context"
	"fmt"
	"net/http"

	types "github.com/jontk/slurm-client/api"
	adapterbase "github.com/jontk/slurm-client/internal/adapters/base"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_41"
)

// WCKeyAdapter implements the WCKeyManager interface for v0.0.41
type WCKeyAdapter struct {
	*adapterbase.BaseManager
	client *api.ClientWithResponses
}

// NewWCKeyAdapter creates a new WCKey adapter for v0.0.41
func NewWCKeyAdapter(client *api.ClientWithResponses) *WCKeyAdapter {
	return &WCKeyAdapter{
		BaseManager: adapterbase.NewBaseManager("v0.0.41", "WCKey"),
		client:      client,
	}
}

// List retrieves a list of WCKeys
func (a *WCKeyAdapter) List(ctx context.Context, opts *types.WCKeyListOptions) (*types.WCKeyList, error) {
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}
	params := &api.SlurmdbV0041GetWckeysParams{}
	// Filtering logic not implemented due to limited param support
	resp, err := a.client.SlurmdbV0041GetWckeysWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list wckeys: %w", err)
	}
	if resp.HTTPResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: failed to list wckeys", resp.HTTPResponse.StatusCode)
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected nil response")
	}
	return toInternalWCKeyList(resp.JSON200)
}

// Get retrieves a specific WCKey by ID
// Note: v0.0.41 API uses Get /wckey/{id} endpoint
func (a *WCKeyAdapter) Get(ctx context.Context, wckeyID string) (*types.WCKey, error) {
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if err := a.ValidateResourceName("wckey ID", wckeyID); err != nil {
		return nil, err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}
	resp, err := a.client.SlurmdbV0041GetWckeyWithResponse(ctx, wckeyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wckey %s: %w", wckeyID, err)
	}
	if resp.HTTPResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: failed to get wckey %s", resp.HTTPResponse.StatusCode, wckeyID)
	}
	if resp.JSON200 == nil || len(resp.JSON200.Wckeys) == 0 {
		return nil, fmt.Errorf("wckey %s not found", wckeyID)
	}
	// Convert the first (and should be only) wckey
	list, err := toInternalWCKeyList(resp.JSON200)
	if err != nil {
		return nil, err
	}
	if len(list.WCKeys) == 0 {
		return nil, fmt.Errorf("wckey %s not found", wckeyID)
	}
	return &list.WCKeys[0], nil
}

// Create creates a new WCKey
func (a *WCKeyAdapter) Create(ctx context.Context, wckey *types.WCKeyCreate) (*types.WCKeyCreateResponse, error) {
	return nil, a.HandleNotImplemented("Create", "v0.0.41 wckey adapter")
}

// Update updates an existing WCKey
func (a *WCKeyAdapter) Update(ctx context.Context, wckeyName, user, cluster string, update *types.WCKeyUpdate) error {
	return a.HandleNotImplemented("Update", "v0.0.41 wckey adapter")
}

// Delete deletes a WCKey
func (a *WCKeyAdapter) Delete(ctx context.Context, wckeyID string) error {
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	resp, err := a.client.SlurmdbV0041DeleteWckeyWithResponse(ctx, wckeyID)
	if err != nil {
		return fmt.Errorf("failed to delete wckey %s: %w", wckeyID, err)
	}
	if resp.HTTPResponse.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d: failed to delete wckey %s", resp.HTTPResponse.StatusCode, wckeyID)
	}
	return nil
}
