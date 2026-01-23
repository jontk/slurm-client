// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_43

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jontk/slurm-client/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
)

// === Manager Methods ===
// Note: All manager methods are now generated in managers.go

// Requeue requeues a job
func (m *JobManager) Requeue(ctx context.Context, jobID string) error {
	if m.impl == nil {
		m.impl = NewJobManagerImpl(m.client)
	}
	return m.impl.Requeue(ctx, jobID)
}

// === Standalone Operations Implementation ===

// GetLicenses retrieves license information
func (c *WrapperClient) GetLicenses(ctx context.Context) (*interfaces.LicenseList, error) {
	if c.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Call the generated OpenAPI client
	resp, err := c.apiClient.SlurmV0043GetLicensesWithResponse(ctx)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
	}

	// Check HTTP status
	if resp.StatusCode() != 200 {
		return nil, errors.NewClientError(
			errors.ErrorCodeResourceNotFound,
			fmt.Sprintf("Operation failed with status %d", resp.StatusCode()),
		)
	}

	if resp.JSON200 == nil {
		return nil, errors.NewClientError(errors.ErrorCodeValidationFailed, "Empty response from licenses API")
	}

	// Check for API errors
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
			var description string
			if apiErr.Description != nil {
				description = *apiErr.Description
			}
			var source string
			if apiErr.Source != nil {
				source = *apiErr.Source
			}
			apiErrors[i] = errors.SlurmAPIErrorDetail{
				ErrorNumber: errorNumber,
				ErrorCode:   errorCode,
				Description: description,
				Source:      source,
			}
		}
		return nil, errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.43", apiErrors).SlurmError
	}

	// Convert response to our interface types
	licenses := make([]interfaces.License, 0)
	if len(resp.JSON200.Licenses) > 0 {
		for _, license := range resp.JSON200.Licenses {
			convertedLicense := interfaces.License{
				Name:       getStringFromPtr(license.LicenseName),
				Total:      getIntFromPtr(license.Total),
				Used:       getIntFromPtr(license.Used),
				Available:  getIntFromPtr(license.Free),
				Reserved:   getIntFromPtr(license.Reserved),
				Remote:     getBoolFromPtr(license.Remote),
				Server:     "", // Not available in this API version
				LastUpdate: time.Unix(getInt64FromPtr(license.LastUpdate), 0),
			}
			// Calculate percentage
			if convertedLicense.Total > 0 {
				convertedLicense.Percent = float64(convertedLicense.Used) / float64(convertedLicense.Total) * 100
			}
			licenses = append(licenses, convertedLicense)
		}
	}

	result := &interfaces.LicenseList{
		Licenses: licenses,
		Meta:     map[string]interface{}{},
	}

	// Add warnings if any
	if resp.JSON200.Warnings != nil && len(*resp.JSON200.Warnings) > 0 {
		warnings := make([]string, len(*resp.JSON200.Warnings))
		for i, warning := range *resp.JSON200.Warnings {
			warnings[i] = getStringFromPtr(warning.Description)
		}
		result.Meta["warnings"] = warnings
	}

	return result, nil
}

// GetShares retrieves fairshare information with optional filtering
func (c *WrapperClient) GetShares(ctx context.Context, opts *interfaces.GetSharesOptions) (*interfaces.SharesList, error) {
	if c.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Prepare parameters
	params := &SlurmV0043GetSharesParams{}
	if opts != nil {
		if len(opts.Users) > 0 {
			users := strings.Join(opts.Users, ",")
			params.Users = &users
		}
		if len(opts.Accounts) > 0 {
			accounts := strings.Join(opts.Accounts, ",")
			params.Accounts = &accounts
		}
		// Note: Clusters and UpdateTime are not supported by this API version
	}

	// Call the generated OpenAPI client
	resp, err := c.apiClient.SlurmV0043GetSharesWithResponse(ctx, params)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
	}

	// Check HTTP status
	if resp.StatusCode() != 200 {
		return nil, errors.NewClientError(
			errors.ErrorCodeServerInternal,
			fmt.Sprintf("Operation failed with status %d", resp.StatusCode()))
	}

	if resp.JSON200 == nil {
		return nil, errors.NewClientError(errors.ErrorCodeValidationFailed, "Empty response from shares API")
	}

	// Check for API errors
	if resp.JSON200.Errors != nil && len(*resp.JSON200.Errors) > 0 {
		apiErrors := make([]errors.SlurmAPIErrorDetail, len(*resp.JSON200.Errors))
		for i, apiErr := range *resp.JSON200.Errors {
			var errorNumber int
			if apiErr.ErrorNumber != nil {
				errorNumber = int(*apiErr.ErrorNumber)
			}
			apiErrors[i] = errors.SlurmAPIErrorDetail{
				ErrorNumber: errorNumber,
				ErrorCode:   getStringFromPtr(apiErr.Error),
				Description: getStringFromPtr(apiErr.Description),
				Source:      getStringFromPtr(apiErr.Source),
			}
		}
		return nil, errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.43", apiErrors).SlurmError
	}

	// Convert response to our interface types
	shares := make([]interfaces.Share, 0)
	if resp.JSON200.Shares.Shares != nil && len(*resp.JSON200.Shares.Shares) > 0 {
		for _, shareWrap := range *resp.JSON200.Shares.Shares {
			convertedShare := interfaces.Share{
				Name:        getStringFromPtr(shareWrap.Name),
				User:        getStringFromPtr(shareWrap.Name), // API doesn't have separate user field
				Account:     getStringFromPtr(shareWrap.Parent),
				Cluster:     getStringFromPtr(shareWrap.Cluster),
				Partition:   getStringFromPtr(shareWrap.Partition),
				Shares:      getIntFromNoValStruct(shareWrap.Shares),
				RawShares:   getIntFromNoValStruct(shareWrap.Shares), // Use same value
				NormShares:  getFloatFromNoValStruct(shareWrap.SharesNormalized),
				RawUsage:    getIntFromPtr64(shareWrap.Usage),
				NormUsage:   getFloatFromNoValStruct(shareWrap.UsageNormalized),
				EffectUsage: getFloatFromNoValStruct(shareWrap.EffectiveUsage),
				FairShare:   getFairShareFactor(shareWrap.Fairshare),
				LevelFS:     getFairShareLevel(shareWrap.Fairshare),
				Priority:    0.0, // Not available in this API
				Level:       0,   // Not available in this API
				LastUpdate:  time.Now(),
			}
			shares = append(shares, convertedShare)
		}
	}

	result := &interfaces.SharesList{
		Shares: shares,
		Meta:   map[string]interface{}{},
	}

	// Add warnings if any
	if resp.JSON200.Warnings != nil && len(*resp.JSON200.Warnings) > 0 {
		warnings := make([]string, len(*resp.JSON200.Warnings))
		for i, warning := range *resp.JSON200.Warnings {
			warnings[i] = getStringFromPtr(warning.Description)
		}
		result.Meta["warnings"] = warnings
	}

	return result, nil
}

// GetConfig retrieves SLURM configuration
func (c *WrapperClient) GetConfig(ctx context.Context) (*interfaces.Config, error) {
	if c.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Call the generated OpenAPI client
	resp, err := c.apiClient.SlurmdbV0043GetConfigWithResponse(ctx)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
	}

	// Check HTTP status
	if resp.StatusCode() != 200 {
		return nil, errors.NewClientError(
			errors.ErrorCodeServerInternal,
			fmt.Sprintf("Operation failed with status %d", resp.StatusCode()))
	}

	if resp.JSON200 == nil {
		return nil, errors.NewClientError(errors.ErrorCodeValidationFailed, "Empty response from config API")
	}

	// Check for API errors
	if resp.JSON200.Errors != nil && len(*resp.JSON200.Errors) > 0 {
		apiErrors := make([]errors.SlurmAPIErrorDetail, len(*resp.JSON200.Errors))
		for i, apiErr := range *resp.JSON200.Errors {
			var errorNumber int
			if apiErr.ErrorNumber != nil {
				errorNumber = int(*apiErr.ErrorNumber)
			}
			apiErrors[i] = errors.SlurmAPIErrorDetail{
				ErrorNumber: errorNumber,
				ErrorCode:   getStringFromPtr(apiErr.Error),
				Description: getStringFromPtr(apiErr.Description),
				Source:      getStringFromPtr(apiErr.Source),
			}
		}
		return nil, errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.43", apiErrors).SlurmError
	}

	// Convert response to our interface types
	config := &interfaces.Config{
		Parameters: make(map[string]interface{}),
		Nodes:      make([]interfaces.ConfigNode, 0),
		Partitions: make([]interfaces.ConfigPartition, 0),
	}

	// Extract configuration parameters from the response
	// Note: The actual structure depends on the API response format
	// This is a basic implementation that can be extended based on the actual API response
	if resp.JSON200.Clusters != nil {
		config.Parameters["clusters"] = resp.JSON200.Clusters
	}
	if resp.JSON200.Accounts != nil {
		config.Parameters["accounts"] = resp.JSON200.Accounts
	}
	if resp.JSON200.Associations != nil {
		config.Parameters["associations"] = resp.JSON200.Associations
	}
	if resp.JSON200.Qos != nil {
		config.Parameters["qos"] = resp.JSON200.Qos
	}
	if resp.JSON200.Tres != nil {
		config.Parameters["tres"] = resp.JSON200.Tres
	}
	if resp.JSON200.Users != nil {
		config.Parameters["users"] = resp.JSON200.Users
	}

	return config, nil
}

// GetDiagnostics retrieves SLURM diagnostics information
func (c *WrapperClient) GetDiagnostics(ctx context.Context) (*interfaces.Diagnostics, error) {
	if c.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Call the generated OpenAPI client
	resp, err := c.apiClient.SlurmV0043GetDiagWithResponse(ctx)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
	}

	// Check HTTP status
	if resp.StatusCode() != 200 {
		return nil, errors.NewClientError(
			errors.ErrorCodeServerInternal,
			fmt.Sprintf("Operation failed with status %d", resp.StatusCode()))
	}

	if resp.JSON200 == nil {
		return nil, errors.NewClientError(errors.ErrorCodeValidationFailed, "Empty response from diagnostics API")
	}

	// Check for API errors
	if resp.JSON200.Errors != nil && len(*resp.JSON200.Errors) > 0 {
		apiErrors := make([]errors.SlurmAPIErrorDetail, len(*resp.JSON200.Errors))
		for i, apiErr := range *resp.JSON200.Errors {
			var errorNumber int
			if apiErr.ErrorNumber != nil {
				errorNumber = int(*apiErr.ErrorNumber)
			}
			apiErrors[i] = errors.SlurmAPIErrorDetail{
				ErrorNumber: errorNumber,
				ErrorCode:   getStringFromPtr(apiErr.Error),
				Description: getStringFromPtr(apiErr.Description),
				Source:      getStringFromPtr(apiErr.Source),
			}
		}
		return nil, errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.43", apiErrors).SlurmError
	}

	// Convert response to our interface types
	diagnostics := &interfaces.Diagnostics{
		DataCollected:     time.Now(),
		RPCsByMessageType: make(map[string]int),
		RPCsByUser:        make(map[string]int),
		Statistics:        make(map[string]interface{}),
	}

	// Extract diagnostics information from the response
	// The resp.JSON200.Statistics field is of type V0043StatsMsg (not a pointer)
	if resp.JSON200.Statistics.ReqTime != nil && resp.JSON200.Statistics.ReqTime.Number != nil {
		diagnostics.ReqTime = *resp.JSON200.Statistics.ReqTime.Number
	}
	if resp.JSON200.Statistics.ReqTimeStart != nil && resp.JSON200.Statistics.ReqTimeStart.Number != nil {
		diagnostics.ReqTimeStart = *resp.JSON200.Statistics.ReqTimeStart.Number
	}
	if resp.JSON200.Statistics.ServerThreadCount != nil {
		diagnostics.ServerThreadCount = int(*resp.JSON200.Statistics.ServerThreadCount)
	}
	if resp.JSON200.Statistics.AgentCount != nil {
		diagnostics.AgentCount = int(*resp.JSON200.Statistics.AgentCount)
	}
	if resp.JSON200.Statistics.AgentThreadCount != nil {
		diagnostics.AgentThreadCount = int(*resp.JSON200.Statistics.AgentThreadCount)
	}
	if resp.JSON200.Statistics.JobsSubmitted != nil {
		diagnostics.JobsSubmitted = int(*resp.JSON200.Statistics.JobsSubmitted)
	}
	if resp.JSON200.Statistics.JobsStarted != nil {
		diagnostics.JobsStarted = int(*resp.JSON200.Statistics.JobsStarted)
	}
	if resp.JSON200.Statistics.JobsCompleted != nil {
		diagnostics.JobsCompleted = int(*resp.JSON200.Statistics.JobsCompleted)
	}
	if resp.JSON200.Statistics.JobsCanceled != nil {
		diagnostics.JobsCanceled = int(*resp.JSON200.Statistics.JobsCanceled)
	}
	if resp.JSON200.Statistics.JobsFailed != nil {
		diagnostics.JobsFailed = int(*resp.JSON200.Statistics.JobsFailed)
	}

	// Store raw statistics for additional information
	diagnostics.Statistics["raw_statistics"] = resp.JSON200.Statistics

	return diagnostics, nil
}

// GetDBDiagnostics retrieves SLURM database diagnostics information
func (c *WrapperClient) GetDBDiagnostics(ctx context.Context) (*interfaces.Diagnostics, error) {
	if c.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Call the generated OpenAPI client
	resp, err := c.apiClient.SlurmdbV0043GetDiagWithResponse(ctx)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
	}

	// Check HTTP status
	if resp.StatusCode() != 200 {
		return nil, errors.NewClientError(
			errors.ErrorCodeServerInternal,
			fmt.Sprintf("Operation failed with status %d", resp.StatusCode()))
	}

	if resp.JSON200 == nil {
		return nil, errors.NewClientError(errors.ErrorCodeValidationFailed, "Empty response from DB diagnostics API")
	}

	// Check for API errors
	if resp.JSON200.Errors != nil && len(*resp.JSON200.Errors) > 0 {
		apiErrors := make([]errors.SlurmAPIErrorDetail, len(*resp.JSON200.Errors))
		for i, apiErr := range *resp.JSON200.Errors {
			var errorNumber int
			if apiErr.ErrorNumber != nil {
				errorNumber = int(*apiErr.ErrorNumber)
			}
			apiErrors[i] = errors.SlurmAPIErrorDetail{
				ErrorNumber: errorNumber,
				ErrorCode:   getStringFromPtr(apiErr.Error),
				Description: getStringFromPtr(apiErr.Description),
				Source:      getStringFromPtr(apiErr.Source),
			}
		}
		return nil, errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.43", apiErrors).SlurmError
	}

	// Convert response to our interface types
	diagnostics := &interfaces.Diagnostics{
		DataCollected:     time.Now(),
		RPCsByMessageType: make(map[string]int),
		RPCsByUser:        make(map[string]int),
		Statistics:        make(map[string]interface{}),
	}

	// Extract DB diagnostics information from the response
	// The resp.JSON200.Statistics field is of type V0043StatsRec (not a pointer)
	// Store raw statistics for additional information
	diagnostics.Statistics["raw_db_statistics"] = resp.JSON200.Statistics
	diagnostics.Statistics["source"] = "database"

	return diagnostics, nil
}

// GetInstance retrieves a specific database instance
func (c *WrapperClient) GetInstance(ctx context.Context, opts *interfaces.GetInstanceOptions) (*interfaces.Instance, error) {
	if c.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Prepare parameters
	params := &SlurmdbV0043GetInstanceParams{}
	if opts != nil {
		if opts.Cluster != "" {
			params.Cluster = &opts.Cluster
		}
		if opts.Extra != "" {
			params.Extra = &opts.Extra
		}
		if opts.Instance != "" {
			// Note: The API doesn't have an Instance field in the params
			// We'll use the instance name for filtering in the response
		}
		if len(opts.NodeList) > 0 {
			nodeList := strings.Join(opts.NodeList, ",")
			params.NodeList = &nodeList
		}
		if opts.TimeStart != nil {
			// Convert time to string format expected by API
			timeStartStr := opts.TimeStart.Format("2006-01-02T15:04:05")
			params.TimeStart = &timeStartStr
		}
		if opts.TimeEnd != nil {
			// Convert time to string format expected by API
			timeEndStr := opts.TimeEnd.Format("2006-01-02T15:04:05")
			params.TimeEnd = &timeEndStr
		}
	}

	// Call the generated OpenAPI client
	resp, err := c.apiClient.SlurmdbV0043GetInstanceWithResponse(ctx, params)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
	}

	// Check HTTP status
	if resp.StatusCode() != 200 {
		return nil, errors.NewClientError(
			errors.ErrorCodeServerInternal,
			fmt.Sprintf("Operation failed with status %d", resp.StatusCode()))
	}

	if resp.JSON200 == nil {
		return nil, errors.NewClientError(errors.ErrorCodeValidationFailed, "Empty response from instance API")
	}

	// Check for API errors
	if resp.JSON200.Errors != nil && len(*resp.JSON200.Errors) > 0 {
		apiErrors := make([]errors.SlurmAPIErrorDetail, len(*resp.JSON200.Errors))
		for i, apiErr := range *resp.JSON200.Errors {
			var errorNumber int
			if apiErr.ErrorNumber != nil {
				errorNumber = int(*apiErr.ErrorNumber)
			}
			apiErrors[i] = errors.SlurmAPIErrorDetail{
				ErrorNumber: errorNumber,
				ErrorCode:   getStringFromPtr(apiErr.Error),
				Description: getStringFromPtr(apiErr.Description),
				Source:      getStringFromPtr(apiErr.Source),
			}
		}
		return nil, errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.43", apiErrors).SlurmError
	}

	// Convert response to our interface types
	// resp.JSON200.Instances is of type V0043InstanceList which is []V0043Instance
	if len(resp.JSON200.Instances) > 0 {
		// Filter by instance name if provided
		for _, inst := range resp.JSON200.Instances {
			if opts != nil && opts.Instance != "" {
				if getStringFromPtr(inst.InstanceId) != opts.Instance {
					continue
				}
			}
			instance := &interfaces.Instance{
				Cluster:   getStringFromPtr(inst.Cluster),
				ExtraInfo: getStringFromPtr(inst.Extra),
				Instance:  getStringFromPtr(inst.InstanceId),
				NodeName:  getStringFromPtr(inst.NodeName),
				TimeEnd:   getTimeFromInstance(inst.Time, false),
				TimeStart: getTimeFromInstance(inst.Time, true),
				TRES:      "", // Not available in this structure
				Created:   time.Now(),
				Modified:  time.Now(),
			}
			return instance, nil
		}
	}

	return nil, errors.NewClientError(errors.ErrorCodeResourceNotFound, "Instance not found")
}

// GetInstances retrieves multiple database instances with filtering
func (c *WrapperClient) GetInstances(ctx context.Context, opts *interfaces.GetInstancesOptions) (*interfaces.InstanceList, error) {
	if c.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Prepare parameters
	params := &SlurmdbV0043GetInstancesParams{}
	if opts != nil {
		if len(opts.Clusters) > 0 {
			clusters := strings.Join(opts.Clusters, ",")
			params.Cluster = &clusters
		}
		if opts.Extra != "" {
			params.Extra = &opts.Extra
		}
		// Note: The API doesn't have an Instance field in the params for filtering multiple instances
		if len(opts.NodeList) > 0 {
			nodeList := strings.Join(opts.NodeList, ",")
			params.NodeList = &nodeList
		}
		if opts.TimeStart != nil {
			// Convert time to string format expected by API
			timeStartStr := opts.TimeStart.Format("2006-01-02T15:04:05")
			params.TimeStart = &timeStartStr
		}
		if opts.TimeEnd != nil {
			// Convert time to string format expected by API
			timeEndStr := opts.TimeEnd.Format("2006-01-02T15:04:05")
			params.TimeEnd = &timeEndStr
		}
	}

	// Call the generated OpenAPI client
	resp, err := c.apiClient.SlurmdbV0043GetInstancesWithResponse(ctx, params)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
	}

	// Check HTTP status
	if resp.StatusCode() != 200 {
		return nil, errors.NewClientError(
			errors.ErrorCodeServerInternal,
			fmt.Sprintf("Operation failed with status %d", resp.StatusCode()))
	}

	if resp.JSON200 == nil {
		return nil, errors.NewClientError(errors.ErrorCodeValidationFailed, "Empty response from instances API")
	}

	// Check for API errors
	if resp.JSON200.Errors != nil && len(*resp.JSON200.Errors) > 0 {
		apiErrors := make([]errors.SlurmAPIErrorDetail, len(*resp.JSON200.Errors))
		for i, apiErr := range *resp.JSON200.Errors {
			var errorNumber int
			if apiErr.ErrorNumber != nil {
				errorNumber = int(*apiErr.ErrorNumber)
			}
			apiErrors[i] = errors.SlurmAPIErrorDetail{
				ErrorNumber: errorNumber,
				ErrorCode:   getStringFromPtr(apiErr.Error),
				Description: getStringFromPtr(apiErr.Description),
				Source:      getStringFromPtr(apiErr.Source),
			}
		}
		return nil, errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.43", apiErrors).SlurmError
	}

	// Convert response to our interface types
	instances := make([]interfaces.Instance, 0)
	// resp.JSON200.Instances is of type V0043InstanceList which is []V0043Instance
	for _, inst := range resp.JSON200.Instances {
		// Filter by instance names if provided
		if opts != nil && len(opts.Instances) > 0 {
			found := false
			instName := getStringFromPtr(inst.InstanceId)
			for _, filterName := range opts.Instances {
				if instName == filterName {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		convertedInstance := interfaces.Instance{
			Cluster:   getStringFromPtr(inst.Cluster),
			ExtraInfo: getStringFromPtr(inst.Extra),
			Instance:  getStringFromPtr(inst.InstanceId),
			NodeName:  getStringFromPtr(inst.NodeName),
			TimeEnd:   getTimeFromInstance(inst.Time, false),
			TimeStart: getTimeFromInstance(inst.Time, true),
			TRES:      "", // Not available in this structure
			Created:   time.Now(),
			Modified:  time.Now(),
		}
		instances = append(instances, convertedInstance)
	}

	result := &interfaces.InstanceList{
		Instances: instances,
		Meta:      map[string]interface{}{},
	}

	// Add warnings if any
	if resp.JSON200.Warnings != nil && len(*resp.JSON200.Warnings) > 0 {
		warnings := make([]string, len(*resp.JSON200.Warnings))
		for i, warning := range *resp.JSON200.Warnings {
			warnings[i] = getStringFromPtr(warning.Description)
		}
		result.Meta["warnings"] = warnings
	}

	return result, nil
}

// GetTRES retrieves all TRES (Trackable RESources)
func (c *WrapperClient) GetTRES(ctx context.Context) (*interfaces.TRESList, error) {
	if c.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Call the generated OpenAPI client
	resp, err := c.apiClient.SlurmdbV0043GetTresWithResponse(ctx)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
	}

	// Check HTTP status
	if resp.StatusCode() != 200 {
		return nil, errors.NewClientError(
			errors.ErrorCodeServerInternal,
			fmt.Sprintf("Operation failed with status %d", resp.StatusCode()))
	}

	if resp.JSON200 == nil {
		return nil, errors.NewClientError(errors.ErrorCodeValidationFailed, "Empty response from TRES API")
	}

	// Check for API errors
	if resp.JSON200.Errors != nil && len(*resp.JSON200.Errors) > 0 {
		apiErrors := make([]errors.SlurmAPIErrorDetail, len(*resp.JSON200.Errors))
		for i, apiErr := range *resp.JSON200.Errors {
			var errorNumber int
			if apiErr.ErrorNumber != nil {
				errorNumber = int(*apiErr.ErrorNumber)
			}
			apiErrors[i] = errors.SlurmAPIErrorDetail{
				ErrorNumber: errorNumber,
				ErrorCode:   getStringFromPtr(apiErr.Error),
				Description: getStringFromPtr(apiErr.Description),
				Source:      getStringFromPtr(apiErr.Source),
			}
		}
		return nil, errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.43", apiErrors).SlurmError
	}

	// Convert response to our interface types
	tresList := make([]interfaces.TRES, 0)
	// resp.JSON200.TRES is of type V0043TresList which is []V0043Tres
	for _, tres := range resp.JSON200.TRES {
		tresID := getIntFromPtr(tres.Id)
		// TRES IDs should be non-negative
		if tresID < 0 {
			tresID = 0
		}
		// #nosec G115 -- tresID is validated to be non-negative before conversion
		convertedTRES := interfaces.TRES{
			ID:          uint64(tresID),
			Type:        tres.Type, // Type is required, not a pointer
			Name:        getStringFromPtr(tres.Name),
			Count:       getInt64FromPtr(tres.Count),
			AllocSecs:   0, // Not available in this structure
			Created:     time.Now(),
			Modified:    time.Now(),
			Description: "", // Not available in this structure
		}
		tresList = append(tresList, convertedTRES)
	}

	result := &interfaces.TRESList{
		TRES: tresList,
		Meta: map[string]interface{}{},
	}

	// Add warnings if any
	if resp.JSON200.Warnings != nil && len(*resp.JSON200.Warnings) > 0 {
		warnings := make([]string, len(*resp.JSON200.Warnings))
		for i, warning := range *resp.JSON200.Warnings {
			warnings[i] = getStringFromPtr(warning.Description)
		}
		result.Meta["warnings"] = warnings
	}

	return result, nil
}

// CreateTRES creates a new TRES entry
func (c *WrapperClient) CreateTRES(ctx context.Context, req *interfaces.CreateTRESRequest) (*interfaces.TRES, error) {
	if c.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if req == nil {
		return nil, errors.NewClientError(errors.ErrorCodeValidationFailed, "Create TRES request cannot be nil")
	}

	// Prepare the request body
	requestBody := SlurmdbV0043PostTresJSONRequestBody{
		TRES: []V0043Tres{
			{
				Type: req.Type,
				Name: &req.Name,
			},
		},
	}

	// Call the generated OpenAPI client
	resp, err := c.apiClient.SlurmdbV0043PostTresWithResponse(ctx, requestBody)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
	}

	// Check HTTP status
	if resp.StatusCode() != 200 && resp.StatusCode() != http.StatusCreated {
		return nil, errors.NewClientError(
			errors.ErrorCodeServerInternal,
			fmt.Sprintf("Operation failed with status %d", resp.StatusCode()))
	}

	if resp.JSON200 == nil {
		return nil, errors.NewClientError(errors.ErrorCodeValidationFailed, "Empty response from create TRES API")
	}

	// Check for API errors
	if resp.JSON200.Errors != nil && len(*resp.JSON200.Errors) > 0 {
		apiErrors := make([]errors.SlurmAPIErrorDetail, len(*resp.JSON200.Errors))
		for i, apiErr := range *resp.JSON200.Errors {
			var errorNumber int
			if apiErr.ErrorNumber != nil {
				errorNumber = int(*apiErr.ErrorNumber)
			}
			apiErrors[i] = errors.SlurmAPIErrorDetail{
				ErrorNumber: errorNumber,
				ErrorCode:   getStringFromPtr(apiErr.Error),
				Description: getStringFromPtr(apiErr.Description),
				Source:      getStringFromPtr(apiErr.Source),
			}
		}
		return nil, errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.43", apiErrors).SlurmError
	}

	// The POST response doesn't return the created TRES details
	// Return a basic response with the requested values
	result := &interfaces.TRES{
		Type:        req.Type,
		Name:        req.Name,
		Description: req.Description,
		Created:     time.Now(),
		Modified:    time.Now(),
	}
	return result, nil
}

// Reconfigure triggers a SLURM reconfiguration
func (c *WrapperClient) Reconfigure(ctx context.Context) (*interfaces.ReconfigureResponse, error) {
	if c.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Call the generated OpenAPI client
	resp, err := c.apiClient.SlurmV0043GetReconfigureWithResponse(ctx)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
	}

	// Check HTTP status
	if resp.StatusCode() != 200 {
		return nil, errors.NewClientError(
			errors.ErrorCodeServerInternal,
			fmt.Sprintf("Operation failed with status %d", resp.StatusCode()))
	}

	if resp.JSON200 == nil {
		return nil, errors.NewClientError(errors.ErrorCodeValidationFailed, "Empty response from reconfigure API")
	}

	// Check for API errors
	if resp.JSON200.Errors != nil && len(*resp.JSON200.Errors) > 0 {
		apiErrors := make([]errors.SlurmAPIErrorDetail, len(*resp.JSON200.Errors))
		for i, apiErr := range *resp.JSON200.Errors {
			var errorNumber int
			if apiErr.ErrorNumber != nil {
				errorNumber = int(*apiErr.ErrorNumber)
			}
			apiErrors[i] = errors.SlurmAPIErrorDetail{
				ErrorNumber: errorNumber,
				ErrorCode:   getStringFromPtr(apiErr.Error),
				Description: getStringFromPtr(apiErr.Description),
				Source:      getStringFromPtr(apiErr.Source),
			}
		}
		return nil, errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.43", apiErrors).SlurmError
	}

	// Convert response to our interface types
	result := &interfaces.ReconfigureResponse{
		Status:   "success",
		Message:  "Reconfiguration completed successfully",
		Changes:  make([]string, 0),
		Warnings: make([]string, 0),
		Errors:   make([]string, 0),
		Meta:     map[string]interface{}{},
	}

	// Add warnings if any
	if resp.JSON200.Warnings != nil && len(*resp.JSON200.Warnings) > 0 {
		warnings := make([]string, len(*resp.JSON200.Warnings))
		for i, warning := range *resp.JSON200.Warnings {
			warnings[i] = getStringFromPtr(warning.Description)
		}
		result.Warnings = warnings
	}

	// Store additional response information
	result.Meta["response_data"] = resp.JSON200

	return result, nil
}

// === Helper Functions ===

// getStringFromPtr safely extracts string value from pointer
func getStringFromPtr(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

// getIntFromPtr safely extracts int value from int32 pointer
func getIntFromPtr(ptr *int32) int {
	if ptr == nil {
		return 0
	}
	return int(*ptr)
}

// getInt64FromPtr safely extracts int64 value from pointer
func getInt64FromPtr(ptr *int64) int64 {
	if ptr == nil {
		return 0
	}
	return *ptr
}

// getBoolFromPtr safely extracts bool value from pointer
func getBoolFromPtr(ptr *bool) bool {
	if ptr == nil {
		return false
	}
	return *ptr
}

// getIntFromPtr64 safely extracts int value from int64 pointer
func getIntFromPtr64(ptr *int64) int {
	if ptr == nil {
		return 0
	}
	return int(*ptr)
}

// getIntFromNoValStruct safely extracts int value from V0043Uint32NoValStruct
func getIntFromNoValStruct(noVal *V0043Uint32NoValStruct) int {
	if noVal == nil || noVal.Number == nil {
		return 0
	}
	return int(*noVal.Number)
}

// getFloatFromNoValStruct safely extracts float64 value from V0043Float64NoValStruct
func getFloatFromNoValStruct(noVal *V0043Float64NoValStruct) float64 {
	if noVal == nil || noVal.Number == nil {
		return 0.0
	}
	return *noVal.Number
}

// getFairShareFactor safely extracts fairshare factor
func getFairShareFactor(fairshare *struct {
	Factor *V0043Float64NoValStruct `json:"factor,omitempty"`
	Level  *V0043Float64NoValStruct `json:"level,omitempty"`
}) float64 {
	if fairshare == nil || fairshare.Factor == nil {
		return 0.0
	}
	return getFloatFromNoValStruct(fairshare.Factor)
}

// getFairShareLevel safely extracts fairshare level
func getFairShareLevel(fairshare *struct {
	Factor *V0043Float64NoValStruct `json:"factor,omitempty"`
	Level  *V0043Float64NoValStruct `json:"level,omitempty"`
}) float64 {
	if fairshare == nil || fairshare.Level == nil {
		return 0.0
	}
	return getFloatFromNoValStruct(fairshare.Level)
}

// getTimeFromInstance safely extracts time from instance time structure
func getTimeFromInstance(timeStruct *struct {
	TimeEnd   *int64 `json:"time_end,omitempty"`
	TimeStart *int64 `json:"time_start,omitempty"`
}, isStart bool) int64 {
	if timeStruct == nil {
		return 0
	}
	if isStart && timeStruct.TimeStart != nil {
		return *timeStruct.TimeStart
	}
	if !isStart && timeStruct.TimeEnd != nil {
		return *timeStruct.TimeEnd
	}
	return 0
}
