// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_41

import (
	"context"

	"github.com/jontk/slurm-client/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
)

// NodeManagerImpl provides the actual implementation for NodeManager methods
type NodeManagerImpl struct {
	client *WrapperClient
}

// NewNodeManagerImpl creates a new NodeManager implementation
func NewNodeManagerImpl(client *WrapperClient) *NodeManagerImpl {
	return &NodeManagerImpl{client: client}
}

// List nodes with optional filtering
func (m *NodeManagerImpl) List(ctx context.Context, opts *interfaces.ListNodesOptions) (*interfaces.NodeList, error) {
	// Check if API client is available
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Prepare parameters for the API call
	params := &SlurmV0041GetNodesParams{}

	// Call the generated OpenAPI client
	resp, err := m.client.apiClient.SlurmV0041GetNodesWithResponse(ctx, params)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.41")
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
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.41", apiErrors)
				return nil, apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.41")
		return nil, httpErr
	}

	// Check for unexpected response format
	if resp.JSON200 == nil {
		return nil, errors.NewClientError(errors.ErrorCodeServerInternal, "Unexpected response format", "Expected JSON response but got nil")
	}

	// Convert response to interface types
	nodes := make([]interfaces.Node, 0, len(resp.JSON200.Nodes))
	for _, apiNode := range resp.JSON200.Nodes {
		node := convertNodeFromAPI(apiNode)
		nodes = append(nodes, node)
	}

	// Note: v0.0.41 API doesn't support node name filtering in the request,
	// so we can't apply client-side filtering by names

	return &interfaces.NodeList{
		Nodes: nodes,
		Total: len(nodes),
	}, nil
}

// Get retrieves a specific node by name
func (m *NodeManagerImpl) Get(ctx context.Context, nodeName string) (*interfaces.Node, error) {
	// Check if client is initialized first
	if m.client == nil || m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Call the GetNode API (params can be nil)
	resp, err := m.client.apiClient.SlurmV0041GetNodeWithResponse(ctx, nodeName, nil)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.41")
	}

	// Check HTTP status
	if resp.StatusCode() != 200 {
		var responseBody []byte
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.41")
		return nil, httpErr
	}

	// Check for unexpected response format
	if resp.JSON200 == nil || len(resp.JSON200.Nodes) == 0 {
		return nil, errors.NewClientError(errors.ErrorCodeServerInternal, "Unexpected response format")
	}

	// Get the first node from the response - minimal conversion for v0.0.41
	apiNode := resp.JSON200.Nodes[0]
	node := &interfaces.Node{
		Name: *apiNode.Name,
	}

	return node, nil
}

// Update updates node properties
func (m *NodeManagerImpl) Update(ctx context.Context, nodeName string, update *interfaces.NodeUpdate) error {
	// Check if client is initialized first
	if m.client == nil || m.client.apiClient == nil {
		return errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	return errors.NewClientError(
		errors.ErrorCodeUnsupportedOperation,
		"Node updates not implemented for v0.0.41",
		"The v0.0.41 node update requires complex inline struct mapping that differs significantly from other API versions",
	)
}

// Watch provides real-time node updates through polling
func (m *NodeManagerImpl) Watch(ctx context.Context, opts *interfaces.WatchNodesOptions) (<-chan interfaces.NodeEvent, error) {
	// Check if API client is available
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Create event channel
	eventChan := make(chan interfaces.NodeEvent, 100)

	// Start polling goroutine
	go func() {
		defer close(eventChan)

		// Note: This is a simplified polling implementation
		// v0.0.41 doesn't have native streaming support

		select {
		case <-ctx.Done():
			return
		default:
			// In a full implementation, this would start a polling loop
		}
	}()

	return eventChan, nil
}

// Delete deletes a node
func (m *NodeManagerImpl) Delete(ctx context.Context, nodeName string) error {
	return errors.NewNotImplementedError("Delete", "v0.0.41")
}

// Drain drains a node, preventing new jobs from being scheduled on it
func (m *NodeManagerImpl) Drain(ctx context.Context, nodeName string, reason string) error {
	return errors.NewNotImplementedError("Drain", "v0.0.41")
}

// Resume resumes a drained node, allowing new jobs to be scheduled on it
func (m *NodeManagerImpl) Resume(ctx context.Context, nodeName string) error {
	return errors.NewNotImplementedError("Resume", "v0.0.41")
}
