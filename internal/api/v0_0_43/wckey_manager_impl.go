// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_43

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/jontk/slurm-client/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
)

// WCKeyManagerImpl provides the actual implementation for WCKeyManager methods
type WCKeyManagerImpl struct {
	client *WrapperClient
}

// NewWCKeyManagerImpl creates a new WCKeyManager implementation
func NewWCKeyManagerImpl(client *WrapperClient) *WCKeyManagerImpl {
	return &WCKeyManagerImpl{client: client}
}

// List retrieves WCKeys with optional filtering
func (m *WCKeyManagerImpl) List(ctx context.Context, opts *interfaces.WCKeyListOptions) (*interfaces.WCKeyList, error) {
	// Check if API client is available
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Prepare parameters for the API call
	params := &SlurmdbV0043GetWckeysParams{}

	if opts != nil {
		if len(opts.Names) > 0 {
			nameStr := joinStringSlice(opts.Names, ",")
			params.Name = &nameStr
		}
		if len(opts.Users) > 0 {
			userStr := joinStringSlice(opts.Users, ",")
			params.User = &userStr
		}
		if len(opts.Clusters) > 0 {
			clusterStr := joinStringSlice(opts.Clusters, ",")
			params.Cluster = &clusterStr
		}
	}

	// Call the generated OpenAPI client
	resp, err := m.client.apiClient.SlurmdbV0043GetWckeysWithResponse(ctx, params)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
	}

	// Check HTTP status and handle API errors
	if resp.StatusCode() != 200 {
		var responseBody []byte
		if resp.JSON200 != nil {
			// Try to extract error details from response
			if resp.JSON200.Errors != nil && len(*resp.JSON200.Errors) > 0 {
				apiErrors := make([]errors.SlurmAPIErrorDetail, len(*resp.JSON200.Errors))
				for i, apiErr := range *resp.JSON200.Errors {
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

					apiErrors[i] = errors.SlurmAPIErrorDetail{
						ErrorNumber: errorNumber,
						ErrorCode:   errorCode,
						Source:      source,
						Description: description,
					}
				}
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.43", apiErrors)
				return nil, apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.43")
		return nil, httpErr
	}

	// Check for unexpected response format
	if resp.JSON200 == nil {
		return nil, errors.NewClientError(errors.ErrorCodeServerInternal, "Unexpected response format", "Expected JSON response but got nil")
	}

	// Convert the response to our interface types
	wckeys := make([]interfaces.WCKey, 0, len(resp.JSON200.Wckeys))

	for _, apiWCKey := range resp.JSON200.Wckeys {
		wckey := interfaces.WCKey{
			Name:    apiWCKey.Name,
			User:    apiWCKey.User,
			Cluster: apiWCKey.Cluster,
		}
		wckeys = append(wckeys, wckey)
	}

	return &interfaces.WCKeyList{
		WCKeys: wckeys,
		Total:  len(wckeys),
	}, nil
}

// Get retrieves a specific WCKey
func (m *WCKeyManagerImpl) Get(ctx context.Context, wckeyName, user, cluster string) (*interfaces.WCKey, error) {
	// Check if API client is available
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Create WCKey ID from components
	wckeyID := fmt.Sprintf("%s_%s_%s", wckeyName, user, cluster)

	// Call the generated OpenAPI client
	resp, err := m.client.apiClient.SlurmdbV0043GetWckeyWithResponse(ctx, wckeyID)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
	}

	// Check HTTP status and handle API errors
	if resp.StatusCode() != 200 {
		var responseBody []byte
		if resp.JSON200 != nil {
			// Try to extract error details from response
			if resp.JSON200.Errors != nil && len(*resp.JSON200.Errors) > 0 {
				apiErrors := make([]errors.SlurmAPIErrorDetail, len(*resp.JSON200.Errors))
				for i, apiErr := range *resp.JSON200.Errors {
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

					apiErrors[i] = errors.SlurmAPIErrorDetail{
						ErrorNumber: errorNumber,
						ErrorCode:   errorCode,
						Source:      source,
						Description: description,
					}
				}
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.43", apiErrors)
				return nil, apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.43")
		return nil, httpErr
	}

	// Check for unexpected response format
	if resp.JSON200 == nil || len(resp.JSON200.Wckeys) == 0 {
		return nil, errors.NewClientError(errors.ErrorCodeResourceNotFound, "WCKey not found")
	}

	// Return the first WCKey from the response
	apiWCKey := resp.JSON200.Wckeys[0]
	return &interfaces.WCKey{
		Name:    apiWCKey.Name,
		User:    apiWCKey.User,
		Cluster: apiWCKey.Cluster,
	}, nil
}

// Create creates a new WCKey
func (m *WCKeyManagerImpl) Create(ctx context.Context, wckey *interfaces.WCKeyCreate) (*interfaces.WCKeyCreateResponse, error) {
	// Check if API client is available
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Convert interface WCKeyCreate to API format
	apiWCKey := V0043Wckey{
		Name:    wckey.Name,
		User:    wckey.User,
		Cluster: wckey.Cluster,
	}

	// Create the request body
	requestBody := V0043OpenapiWckeyResp{
		Wckeys: []V0043Wckey{apiWCKey},
	}

	// Prepare parameters (can be empty for creation)
	params := &SlurmdbV0043PostWckeysParams{}

	// Call the generated OpenAPI client
	resp, err := m.client.apiClient.SlurmdbV0043PostWckeysWithResponse(ctx, params, requestBody)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
	}

	// Check HTTP status (200 and 201 for creation is success)
	if resp.StatusCode() != 200 && resp.StatusCode() != http.StatusCreated {
		var responseBody []byte
		// Handle error response similar to other methods
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.43")
		return nil, httpErr
	}

	return &interfaces.WCKeyCreateResponse{
		WCKeyName: wckey.Name,
	}, nil
}

// Update updates an existing WCKey
func (m *WCKeyManagerImpl) Update(ctx context.Context, wckeyName, user, cluster string, update *interfaces.WCKeyUpdate) error {
	// WCKey updates are typically not supported in SLURM - WCKeys are usually just created/deleted
	return errors.NewClientError(errors.ErrorCodeUnsupportedOperation, "WCKey updates are not supported by SLURM")
}

// Delete deletes a WCKey
func (m *WCKeyManagerImpl) Delete(ctx context.Context, wckeyID string) error {
	// Check if API client is available
	if m.client.apiClient == nil {
		return errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Call the generated OpenAPI client
	resp, err := m.client.apiClient.SlurmdbV0043DeleteWckeyWithResponse(ctx, wckeyID)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
	}

	// Check HTTP status
	if resp.StatusCode() != 200 && resp.StatusCode() != 204 {
		var responseBody []byte
		if resp.JSON200 != nil {
			// Try to extract error details from response
			if resp.JSON200.Errors != nil && len(*resp.JSON200.Errors) > 0 {
				apiErrors := make([]errors.SlurmAPIErrorDetail, len(*resp.JSON200.Errors))
				for i, apiErr := range *resp.JSON200.Errors {
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

					apiErrors[i] = errors.SlurmAPIErrorDetail{
						ErrorNumber: errorNumber,
						ErrorCode:   errorCode,
						Source:      source,
						Description: description,
					}
				}
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.43", apiErrors)
				return apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.43")
		return httpErr
	}

	return nil
}

// joinStringSlice joins a string slice with the given separator
func joinStringSlice(slice []string, sep string) string {
	if len(slice) == 0 {
		return ""
	}

	result := slice[0]
	var resultSb309 strings.Builder
	for i := 1; i < len(slice); i++ {
		resultSb309.WriteString(sep + slice[i])
	}
	result += resultSb309.String()
	return result
}
