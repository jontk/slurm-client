package v0_0_42

import (
	"context"
	"fmt"
	"regexp"

	"github.com/jontk/slurm-client/internal/interfaces"
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
	// Check if API client is available
	if m.client.apiClient == nil {
		return nil, fmt.Errorf("API client not initialized")
	}
	
	// Call the ping endpoint to get basic cluster information
	pingResp, err := m.client.apiClient.SlurmV0042GetPingWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster info: %w", err)
	}
	
	// Check HTTP status
	if pingResp.StatusCode() != 200 {
		return nil, fmt.Errorf("API returned status %d: %s", pingResp.StatusCode(), pingResp.Status())
	}
	
	// Check for API errors
	if pingResp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected response format")
	}
	
	if pingResp.JSON200.Errors != nil && len(*pingResp.JSON200.Errors) > 0 {
		return nil, fmt.Errorf("API error: %v", (*pingResp.JSON200.Errors)[0])
	}
	
	// Extract cluster info from ping response
	clusterInfo := &interfaces.ClusterInfo{
		APIVersion: "v0.0.42", // Current API version
	}
	
	if pingResp.JSON200.Meta != nil && pingResp.JSON200.Meta.Slurm != nil {
		if pingResp.JSON200.Meta.Slurm.Version != nil {
			if pingResp.JSON200.Meta.Slurm.Version.Major != nil && 
				pingResp.JSON200.Meta.Slurm.Version.Minor != nil &&
				pingResp.JSON200.Meta.Slurm.Version.Micro != nil {
				clusterInfo.Version = fmt.Sprintf("%s.%s.%s", 
					*pingResp.JSON200.Meta.Slurm.Version.Major,
					*pingResp.JSON200.Meta.Slurm.Version.Minor,
					*pingResp.JSON200.Meta.Slurm.Version.Micro)
			}
		}
		
		if pingResp.JSON200.Meta.Slurm.Release != nil {
			clusterInfo.Release = *pingResp.JSON200.Meta.Slurm.Release
		}
		
		if pingResp.JSON200.Meta.Slurm.Cluster != nil {
			clusterInfo.ClusterName = *pingResp.JSON200.Meta.Slurm.Cluster
		}
	}
	
	// Try to get additional diagnostic information
	diagResp, diagErr := m.client.apiClient.SlurmV0042GetDiagWithResponse(ctx)
	if diagErr == nil && diagResp.StatusCode() == 200 && diagResp.JSON200 != nil &&
		(diagResp.JSON200.Errors == nil || len(*diagResp.JSON200.Errors) == 0) {
		
		// Extract uptime from diagnostic statistics
		if diagResp.JSON200.Statistics.ServerThreadCount != nil {
			// Server thread count can be used as a proxy for activity/uptime
			clusterInfo.Uptime = int(*diagResp.JSON200.Statistics.ServerThreadCount)
		}
	}
	
	return clusterInfo, nil
}

// Ping tests connectivity to the cluster
func (m *InfoManagerImpl) Ping(ctx context.Context) error {
	// Check if API client is available
	if m.client.apiClient == nil {
		return fmt.Errorf("API client not initialized")
	}
	
	// Call the ping endpoint
	resp, err := m.client.apiClient.SlurmV0042GetPingWithResponse(ctx)
	if err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}
	
	// Check HTTP status
	if resp.StatusCode() != 200 {
		return fmt.Errorf("ping returned status %d: %s", resp.StatusCode(), resp.Status())
	}
	
	// Check for API errors
	if resp.JSON200 == nil {
		return fmt.Errorf("ping returned unexpected response format")
	}
	
	if resp.JSON200.Errors != nil && len(*resp.JSON200.Errors) > 0 {
		return fmt.Errorf("ping API error: %v", (*resp.JSON200.Errors)[0])
	}
	
	return nil
}

// Stats retrieves cluster statistics
func (m *InfoManagerImpl) Stats(ctx context.Context) (*interfaces.ClusterStats, error) {
	// Check if API client is available
	if m.client.apiClient == nil {
		return nil, fmt.Errorf("API client not initialized")
	}
	
	// Call the diagnostic endpoint to get cluster statistics
	resp, err := m.client.apiClient.SlurmV0042GetDiagWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster stats: %w", err)
	}
	
	// Check HTTP status
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode(), resp.Status())
	}
	
	// Check for API errors
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected response format")
	}
	
	if resp.JSON200.Errors != nil && len(*resp.JSON200.Errors) > 0 {
		return nil, fmt.Errorf("API error: %v", (*resp.JSON200.Errors)[0])
	}
	
	stats := &interfaces.ClusterStats{}
	
	// Extract statistics from diagnostic response
	// Job statistics
	if resp.JSON200.Statistics.JobsSubmitted != nil {
		stats.TotalJobs = int(*resp.JSON200.Statistics.JobsSubmitted)
	}
	
	if resp.JSON200.Statistics.JobsPending != nil {
		stats.PendingJobs = int(*resp.JSON200.Statistics.JobsPending)
	}
	
	if resp.JSON200.Statistics.JobsRunning != nil {
		stats.RunningJobs = int(*resp.JSON200.Statistics.JobsRunning)
	}
	
	if resp.JSON200.Statistics.JobsCompleted != nil {
		stats.CompletedJobs = int(*resp.JSON200.Statistics.JobsCompleted)
	}
	
	// Server thread count and other metrics can provide additional insights
	// but aren't directly mappable to our interface. We'll use what we have.
	
	// Note: Node and CPU statistics aren't directly available in the diag endpoint
	// These would typically require separate calls to the nodes endpoint
	// For now, we'll leave them as 0 or implement them as separate calls if needed
	
	return stats, nil
}

// Version retrieves API version information
func (m *InfoManagerImpl) Version(ctx context.Context) (*interfaces.APIVersion, error) {
	// Check if API client is available
	if m.client.apiClient == nil {
		return nil, fmt.Errorf("API client not initialized")
	}
	
	// Call the ping endpoint to get version information
	resp, err := m.client.apiClient.SlurmV0042GetPingWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get version info: %w", err)
	}
	
	// Check HTTP status
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode(), resp.Status())
	}
	
	// Check for API errors
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected response format")
	}
	
	if resp.JSON200.Errors != nil && len(*resp.JSON200.Errors) > 0 {
		return nil, fmt.Errorf("API error: %v", (*resp.JSON200.Errors)[0])
	}
	
	apiVersion := &interfaces.APIVersion{
		Version:     "v0.0.42",
		Description: "Slurm REST API v0.0.42",
		Deprecated:  false, // v0.0.42 is currently stable
	}
	
	// Extract Slurm version from ping response
	if resp.JSON200.Meta != nil && resp.JSON200.Meta.Slurm != nil {
		if resp.JSON200.Meta.Slurm.Version != nil {
			if resp.JSON200.Meta.Slurm.Version.Major != nil && 
				resp.JSON200.Meta.Slurm.Version.Minor != nil &&
				resp.JSON200.Meta.Slurm.Version.Micro != nil {
				apiVersion.Release = fmt.Sprintf("%s.%s.%s", 
					*resp.JSON200.Meta.Slurm.Version.Major,
					*resp.JSON200.Meta.Slurm.Version.Minor,
					*resp.JSON200.Meta.Slurm.Version.Micro)
			}
		}
		
		if resp.JSON200.Meta.Slurm.Release != nil {
			// Extract more detailed release info if available
			release := *resp.JSON200.Meta.Slurm.Release
			if release != "" {
				apiVersion.Release = release
				
				// Check if this is a pre-release or development version
				if matched, _ := regexp.MatchString(`(alpha|beta|rc|dev)`, release); matched {
					apiVersion.Description += " (Development/Pre-release)"
				}
			}
		}
	}
	
	return apiVersion, nil
}