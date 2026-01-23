// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_44

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jontk/slurm-client/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
)

// InfoManagerImpl provides the actual implementation for InfoManager methods
type InfoManagerImpl struct {
	client *WrapperClient
}

// NewInfoManagerImpl creates a new InfoManager implementation
func NewInfoManagerImpl(client *WrapperClient) *InfoManagerImpl {
	return &InfoManagerImpl{client: client}
}

// Get retrieves cluster information
func (m *InfoManagerImpl) Get(ctx context.Context) (*interfaces.ClusterInfo, error) {
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Call the diag endpoint for cluster information
	resp, err := m.client.apiClient.SlurmV0044GetDiagWithResponse(ctx)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.44")
	}

	if resp.HTTPResponse.StatusCode != http.StatusOK {
		return nil, errors.NewSlurmError(errors.ErrorCodeServerInternal,
			fmt.Sprintf("HTTP %d", resp.HTTPResponse.StatusCode))
	}

	if resp.JSON200 == nil {
		return nil, errors.NewSlurmError(errors.ErrorCodeServerInternal, "Empty response body")
	}

	// Convert to interface type
	clusterInfo := &interfaces.ClusterInfo{
		APIVersion: "v0.0.44",
	}

	// Extract basic cluster information
	// Note: Simplified implementation until we understand the exact v0.0.44 response structure
	clusterInfo.ClusterName = "SLURM Cluster"
	clusterInfo.Version = "25.11.x" // Placeholder

	return clusterInfo, nil
}

// Ping tests connectivity to the cluster
func (m *InfoManagerImpl) Ping(ctx context.Context) error {
	if m.client.apiClient == nil {
		return errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Use the ping endpoint
	resp, err := m.client.apiClient.SlurmV0044GetPingWithResponse(ctx)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.44")
	}

	if resp.HTTPResponse.StatusCode != http.StatusOK {
		return errors.NewSlurmError(errors.ErrorCodeConnectionRefused,
			fmt.Sprintf("Ping failed with HTTP %d", resp.HTTPResponse.StatusCode))
	}

	// Successful ping
	return nil
}

// Stats retrieves cluster statistics
func (m *InfoManagerImpl) Stats(ctx context.Context) (*interfaces.ClusterStats, error) {
	if m.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Call the diag endpoint for statistics
	resp, err := m.client.apiClient.SlurmV0044GetDiagWithResponse(ctx)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.44")
	}

	if resp.HTTPResponse.StatusCode != http.StatusOK {
		return nil, errors.NewSlurmError(errors.ErrorCodeServerInternal,
			fmt.Sprintf("HTTP %d", resp.HTTPResponse.StatusCode))
	}

	if resp.JSON200 == nil {
		return nil, errors.NewSlurmError(errors.ErrorCodeServerInternal, "Empty response body")
	}

	// Convert to interface type
	stats := &interfaces.ClusterStats{
		// Placeholder values until we understand the exact v0.0.44 response structure
		TotalJobs:      0,
		RunningJobs:    0,
		PendingJobs:    0,
		TotalNodes:     0,
		IdleNodes:      0,
		AllocatedNodes: 0,
	}

	return stats, nil
}

// Version retrieves API version information
func (m *InfoManagerImpl) Version(ctx context.Context) (*interfaces.APIVersion, error) {
	// Return static version information for v0.0.44
	return &interfaces.APIVersion{
		Version:     "v0.0.44",
		Release:     "SLURM REST API v0.0.44",
		Description: "SLURM REST API version 0.0.44",
		Deprecated:  false,
	}, nil
}

// PingDatabase tests connectivity to the SLURM database (v0.0.44)
func (m *InfoManagerImpl) PingDatabase(ctx context.Context) error {
	if m.client.apiClient == nil {
		return errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// For now, return not implemented - database ping might use different endpoints
	return errors.NewSlurmError(errors.ErrorCodeUnsupportedOperation, "Database ping not yet implemented for v0.0.44")
}
