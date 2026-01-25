// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_43

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jontk/slurm-client/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
)

// UserManagerImpl implements the UserManager interface for v0.0.43
type UserManagerImpl struct {
	client *WrapperClient
}

// NewUserManagerImpl creates a new UserManagerImpl
func NewUserManagerImpl(client *WrapperClient) *UserManagerImpl {
	return &UserManagerImpl{
		client: client,
	}
}

// List retrieves a list of users with optional filtering
func (u *UserManagerImpl) List(ctx context.Context, opts *interfaces.ListUsersOptions) (*interfaces.UserList, error) {
	if u.client == nil || u.client.apiClient == nil || u.client.apiClient.ClientInterface == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Prepare parameters for the API call
	params := &SlurmdbV0043GetUsersParams{}

	// Call the generated OpenAPI client
	resp, err := u.client.apiClient.SlurmdbV0043GetUsersWithResponse(ctx, params)
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
	if resp.JSON200 == nil || resp.JSON200.Users == nil {
		return nil, errors.NewClientError(errors.ErrorCodeServerInternal, "Unexpected response format", "Expected JSON response with users but got nil")
	}

	// Convert the response to our interface types
	users := make([]interfaces.User, 0, len(resp.JSON200.Users))
	for _, apiUser := range resp.JSON200.Users {
		user, err := convertAPIUserToInterface(&apiUser)
		if err != nil {
			conversionErr := errors.NewClientError(errors.ErrorCodeServerInternal, "Failed to convert user data")
			conversionErr.Cause = err
			conversionErr.Details = "Error converting user data"
			return nil, conversionErr
		}
		users = append(users, *user)
	}

	// Apply client-side filtering if options are provided
	if opts != nil {
		users = filterUsers(users, opts)
	}

	return &interfaces.UserList{
		Users: users,
		Total: len(users),
	}, nil
}

// Get retrieves a specific user by name
func (u *UserManagerImpl) Get(ctx context.Context, userName string) (*interfaces.User, error) {
	// Validate input first (cheap check)
	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}

	// Then check client initialization
	if u.client == nil || u.client.apiClient == nil || u.client.apiClient.ClientInterface == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Prepare parameters for the API call
	params := &SlurmdbV0043GetUserParams{}

	// Call the generated OpenAPI client
	resp, err := u.client.apiClient.SlurmdbV0043GetUserWithResponse(ctx, userName, params)
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
	if resp.JSON200 == nil || resp.JSON200.Users == nil || len(resp.JSON200.Users) == 0 {
		return nil, errors.NewClientError(errors.ErrorCodeResourceNotFound, "User not found", fmt.Sprintf("User '%s' not found", userName))
	}

	// Convert the first user in the response
	user, err := convertAPIUserToInterface(&resp.JSON200.Users[0])
	if err != nil {
		conversionErr := errors.NewClientError(errors.ErrorCodeServerInternal, "Failed to convert user data")
		conversionErr.Cause = err
		conversionErr.Details = fmt.Sprintf("Error converting user '%s'", userName)
		return nil, conversionErr
	}

	return user, nil
}

// GetUserAccounts retrieves all account associations for a user
func (u *UserManagerImpl) GetUserAccounts(ctx context.Context, userName string) ([]*interfaces.UserAccount, error) {
	// Validate input first (cheap check)
	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}

	// Then check client initialization
	if u.client == nil || u.client.apiClient == nil || u.client.apiClient.ClientInterface == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Get user details including associations
	user, err := u.Get(ctx, userName)
	if err != nil {
		return nil, err
	}

	// Convert associations to UserAccount format
	userAccounts := make([]*interfaces.UserAccount, 0)
	for _, assoc := range user.Associations {
		userAccount := &interfaces.UserAccount{
			AccountName:   assoc.Account,
			Partition:     assoc.Partition,
			DefaultQoS:    assoc.DefaultQoS,
			IsDefault:     assoc.IsDefault,
			IsActive:      true, // Assuming active associations are returned
			MaxJobs:       assoc.MaxJobs,
			MaxSubmitJobs: assoc.MaxSubmitJobs,
			MaxWallTime:   assoc.MaxWallDuration,
			Priority:      0,          // Not available in association
			GraceTime:     0,          // Not available in association
			Created:       time.Now(), // Not available, using current time
			Modified:      time.Now(), // Not available, using current time
		}
		// Handle QoS array - use first value
		if len(assoc.QoS) > 0 {
			userAccount.QoS = assoc.QoS[0]
		}
		// Convert MaxTRES from map[string]string to map[string]int
		if assoc.MaxTRES != nil {
			userAccount.MaxTRES = make(map[string]int)
			for k, v := range assoc.MaxTRES {
				if intVal, err := stringToInt(v); err == nil {
					userAccount.MaxTRES[k] = intVal
				}
			}
		}
		userAccounts = append(userAccounts, userAccount)
	}

	return userAccounts, nil
}

// GetUserQuotas retrieves quota information for a user
func (u *UserManagerImpl) GetUserQuotas(ctx context.Context, userName string) (*interfaces.UserQuota, error) {
	// Validate input first (cheap check)
	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}

	// Then check client initialization
	if u.client == nil || u.client.apiClient == nil || u.client.apiClient.ClientInterface == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Get user details including associations
	user, err := u.Get(ctx, userName)
	if err != nil {
		return nil, err
	}

	// Aggregate quotas from all associations
	userQuota := &interfaces.UserQuota{
		UserName:       userName,
		DefaultAccount: user.DefaultAccount,
		AccountQuotas:  make(map[string]*interfaces.UserAccountQuota),
		IsActive:       true,
		Enforcement:    "soft",
	}

	// Track effective quotas (most restrictive across all accounts)
	var minMaxJobs *int
	var minMaxSubmitJobs *int
	var minMaxWallTime *int
	var minMaxCPUs *int
	var minMaxNodes *int

	// Process each association
	for _, assoc := range user.Associations {
		accountQuota := &interfaces.UserAccountQuota{
			AccountName:   assoc.Account,
			MaxJobs:       assoc.MaxJobs,
			MaxSubmitJobs: assoc.MaxSubmitJobs,
			MaxWallTime:   assoc.MaxWallDuration,
			Priority:      assoc.Shares, // Using shares as priority
			DefaultQoS:    assoc.DefaultQoS,
		}

		// Handle QoS array
		if len(assoc.QoS) > 0 {
			accountQuota.QoS = assoc.QoS
		}

		// Convert TRES limits
		if assoc.MaxTRES != nil {
			accountQuota.TRESLimits = make(map[string]int)
			for k, v := range assoc.MaxTRES {
				if intVal, err := stringToInt(v); err == nil {
					accountQuota.TRESLimits[k] = intVal
				}
			}
		}

		userQuota.AccountQuotas[assoc.Account] = accountQuota

		// Update effective quotas (most restrictive)
		if assoc.MaxJobs > 0 {
			if minMaxJobs == nil || assoc.MaxJobs < *minMaxJobs {
				minMaxJobs = &assoc.MaxJobs
			}
		}
		if assoc.MaxSubmitJobs > 0 {
			if minMaxSubmitJobs == nil || assoc.MaxSubmitJobs < *minMaxSubmitJobs {
				minMaxSubmitJobs = &assoc.MaxSubmitJobs
			}
		}
		if assoc.MaxWallDuration > 0 {
			if minMaxWallTime == nil || assoc.MaxWallDuration < *minMaxWallTime {
				minMaxWallTime = &assoc.MaxWallDuration
			}
		}

		// Extract CPU and node limits from MaxTRES if available
		if assoc.MaxTRES != nil {
			if cpuLimit, ok := assoc.MaxTRES["cpu"]; ok {
				if cpuInt, err := stringToInt(cpuLimit); err == nil && cpuInt > 0 {
					if minMaxCPUs == nil || cpuInt < *minMaxCPUs {
						minMaxCPUs = &cpuInt
					}
				}
			}
			if nodeLimit, ok := assoc.MaxTRES["node"]; ok {
				if nodeInt, err := stringToInt(nodeLimit); err == nil && nodeInt > 0 {
					if minMaxNodes == nil || nodeInt < *minMaxNodes {
						minMaxNodes = &nodeInt
					}
				}
			}
		}
	}

	// Set effective quotas in the main UserQuota fields
	if minMaxJobs != nil {
		userQuota.MaxJobs = *minMaxJobs
	}
	if minMaxSubmitJobs != nil {
		userQuota.MaxSubmitJobs = *minMaxSubmitJobs
	}
	if minMaxWallTime != nil {
		userQuota.MaxWallTime = *minMaxWallTime
	}
	if minMaxCPUs != nil {
		userQuota.MaxCPUs = *minMaxCPUs
	}
	if minMaxNodes != nil {
		userQuota.MaxNodes = *minMaxNodes
	}

	return userQuota, nil
}

// GetUserDefaultAccount retrieves the default account for a user
func (u *UserManagerImpl) GetUserDefaultAccount(ctx context.Context, userName string) (*interfaces.Account, error) {
	// Validate input first (cheap check)
	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}

	// Then check client initialization
	if u.client == nil || u.client.apiClient == nil || u.client.apiClient.ClientInterface == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Get user accounts
	userAccounts, err := u.GetUserAccounts(ctx, userName)
	if err != nil {
		return nil, err
	}

	// Find the default account
	for _, ua := range userAccounts {
		if ua.IsDefault {
			// Get full account details
			accountMgr := NewAccountManagerImpl(u.client)
			account, err := accountMgr.Get(ctx, ua.AccountName)
			if err != nil {
				return nil, fmt.Errorf("failed to get default account details: %w", err)
			}
			return account, nil
		}
	}

	// If no default account is explicitly set, return the first one
	if len(userAccounts) > 0 {
		accountMgr := NewAccountManagerImpl(u.client)
		account, err := accountMgr.Get(ctx, userAccounts[0].AccountName)
		if err != nil {
			return nil, fmt.Errorf("failed to get account details: %w", err)
		}
		return account, nil
	}

	return nil, errors.NewClientError(errors.ErrorCodeResourceNotFound, "No default account found", fmt.Sprintf("User '%s' has no account associations", userName))
}

// GetUserFairShare retrieves fair-share information for a user
func (u *UserManagerImpl) GetUserFairShare(ctx context.Context, userName string) (*interfaces.UserFairShare, error) {
	// Validate input first (cheap check)
	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}

	// Then check client initialization
	if u.client == nil || u.client.apiClient == nil || u.client.apiClient.ClientInterface == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Use the shares API to get fair share information
	params := &SlurmV0043GetSharesParams{
		Users: &userName, // Users is a string parameter
	}

	resp, err := u.client.apiClient.SlurmV0043GetSharesWithResponse(ctx, params)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
	}

	// Check HTTP status and handle API errors
	if resp.StatusCode() != 200 {
		if resp.JSON200 != nil && resp.JSON200.Errors != nil && len(*resp.JSON200.Errors) > 0 {
			apiErrors := convertAPIErrors(*resp.JSON200.Errors)
			apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.43", apiErrors)
			return nil, apiError.SlurmError
		}
		return nil, errors.WrapHTTPError(resp.StatusCode(), nil, "v0.0.43")
	}

	// Check for valid response
	if resp.JSON200 == nil || resp.JSON200.Shares.Shares == nil || len(*resp.JSON200.Shares.Shares) == 0 {
		return nil, errors.NewClientError(errors.ErrorCodeResourceNotFound, "Fair share information not found", fmt.Sprintf("No fair-share data for user '%s'", userName))
	}

	// Find the user's share information
	var userShare *V0043AssocSharesObjWrap
	for _, share := range *resp.JSON200.Shares.Shares {
		// Check if this share is for our user - using Name field
		if share.Name != nil && *share.Name == userName {
			userShare = &share
			break
		}
	}

	if userShare == nil {
		return nil, errors.NewClientError(errors.ErrorCodeResourceNotFound, fmt.Sprintf("Fair share information not found for user '%s'", userName))
	}

	// Convert to interface type
	fairShare := &interfaces.UserFairShare{
		UserName: userName,
	}

	// Extract fair share data
	// Account info is not directly available in share object, use userName instead
	fairShare.Account = userName // Default to userName as account is not in share data

	if userShare.Cluster != nil {
		fairShare.Cluster = *userShare.Cluster
	}
	// Partition is not available in V0043AssocSharesObjWrap
	if userShare.Shares != nil && userShare.Shares.Number != nil {
		fairShare.RawShares = int(*userShare.Shares.Number)
	}
	if userShare.SharesNormalized != nil && userShare.SharesNormalized.Number != nil {
		fairShare.NormalizedShares = *userShare.SharesNormalized.Number
	}
	if userShare.Usage != nil {
		// Usage is an int64 pointer
		fairShare.EffectiveUsage = float64(*userShare.Usage)
	}
	if userShare.Fairshare != nil {
		if userShare.Fairshare.Level != nil && userShare.Fairshare.Level.Number != nil {
			fairShare.Level = int(*userShare.Fairshare.Level.Number)
		}
		if userShare.Fairshare.Factor != nil && userShare.Fairshare.Factor.Number != nil {
			fairShare.FairShareFactor = *userShare.Fairshare.Factor.Number
		}
	}

	// Build priority factors with available data
	fairShare.PriorityFactors = &interfaces.JobPriorityFactors{
		Age:       0,                                     // Not available in shares API
		FairShare: int(fairShare.FairShareFactor * 1000), // Convert to int
		JobSize:   0,                                     // Not available in shares API
		Partition: 0,                                     // Not available in shares API
		QoS:       0,                                     // Not available in shares API
		TRES:      0,                                     // Not available in shares API
		Site:      0,                                     // Not available in shares API
		Nice:      0,                                     // Not available in shares API
		Assoc:     0,                                     // Not available in shares API
		Total:     int(fairShare.FairShareFactor * 1000), // Set total to fair share
	}

	return fairShare, nil
}

// CalculateJobPriority calculates job priority for a user and job submission
func (u *UserManagerImpl) CalculateJobPriority(ctx context.Context, userName string, jobSubmission *interfaces.JobSubmission) (*interfaces.JobPriorityInfo, error) {
	// Validate input first (cheap check)
	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}

	if jobSubmission == nil {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "job submission data is required", "jobSubmission", jobSubmission, nil)
	}

	// Validate job submission data
	if jobSubmission.Script == "" && jobSubmission.Command == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "job script or command is required", "jobSubmission.Script/Command", fmt.Sprintf("Script: %s, Command: %s", jobSubmission.Script, jobSubmission.Command), nil)
	}

	// Then check client initialization
	if u.client == nil || u.client.apiClient == nil || u.client.apiClient.ClientInterface == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Get user's fair share information first
	fairShare, err := u.GetUserFairShare(ctx, userName)
	if err != nil {
		// If we can't get fair share, continue with basic priority
		fairShare = &interfaces.UserFairShare{
			UserName:        userName,
			FairShareFactor: 0.5, // Default fair share
		}
	}

	// Create priority info
	priorityInfo := &interfaces.JobPriorityInfo{
		UserName:       userName,
		Account:        fairShare.Account,
		Partition:      jobSubmission.Partition,
		QoS:            "",   // QoS not available in JobSubmission
		Priority:       1000, // Base priority
		EligibleTime:   time.Now(),
		EstimatedStart: time.Now().Add(5 * time.Minute), // Estimate
	}

	// Calculate priority factors
	factors := &interfaces.JobPriorityFactors{
		Age:       100,                                              // Base age factor
		FairShare: int(fairShare.FairShareFactor * 1000),            // Convert to int
		JobSize:   int(calculateJobSizeFactor(jobSubmission) * 100), // Convert to int
		Partition: 100,                                              // Base partition factor
		QoS:       100,                                              // Base QoS factor
		TRES:      100,                                              // Base TRES factor
		Site:      0,                                                // No site adjustment
		Nice:      0,                                                // No nice adjustment
		Assoc:     100,                                              // Base association factor
	}

	// Calculate total priority from factors
	factors.Total = factors.Age + factors.FairShare + factors.JobSize +
		factors.Partition + factors.QoS + factors.TRES + factors.Assoc

	priorityInfo.Factors = factors

	// Calculate total priority (simplified calculation)
	totalPriority := factors.Age +
		factors.FairShare +
		factors.JobSize +
		factors.Partition +
		factors.QoS +
		factors.TRES -
		factors.Nice

	priorityInfo.Priority = totalPriority

	priorityInfo.Priority = totalPriority

	return priorityInfo, nil
}

// ValidateUserAccountAccess validates user access to a specific account
func (u *UserManagerImpl) ValidateUserAccountAccess(ctx context.Context, userName, accountName string) (*interfaces.UserAccessValidation, error) {
	// Validate input first (cheap check)
	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}

	if accountName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name is required", "accountName", accountName, nil)
	}

	// Then check client initialization
	if u.client == nil || u.client.apiClient == nil || u.client.apiClient.ClientInterface == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Use the associations API to check user-account access
	// Convert to CSV string as expected by API
	userCSV := userName
	accountCSV := accountName
	params := &SlurmdbV0043GetAssociationsParams{
		User:    &userCSV,
		Account: &accountCSV,
	}

	resp, err := u.client.apiClient.SlurmdbV0043GetAssociationsWithResponse(ctx, params)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
	}

	// Check HTTP status and handle API errors
	if resp.StatusCode() != 200 {
		if resp.JSON200 != nil && resp.JSON200.Errors != nil && len(*resp.JSON200.Errors) > 0 {
			apiErrors := convertAPIErrors(*resp.JSON200.Errors)
			apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.43", apiErrors)
			return nil, apiError.SlurmError
		}
		return nil, errors.WrapHTTPError(resp.StatusCode(), nil, "v0.0.43")
	}

	// Create validation response
	validation := &interfaces.UserAccessValidation{
		UserName:    userName,
		AccountName: accountName,
		HasAccess:   false,
		AccessLevel: "none",
		Permissions: []string{},
		ValidFrom:   time.Now(),
	}

	// Check if associations exist
	if resp.JSON200 == nil || resp.JSON200.Associations == nil || len(resp.JSON200.Associations) == 0 {
		validation.Reason = "No association found between user and account"
		return validation, nil
	}

	// Find the specific association
	for _, assoc := range resp.JSON200.Associations {
		// Check if this association matches our user and account
		if assoc.User != userName ||
			assoc.Account == nil || *assoc.Account != accountName {
			continue
		}
		validation.HasAccess = true
		validation.AccessLevel = "user" // Default access level

		// Check if user is a coordinator
		if assoc.IsDefault != nil && *assoc.IsDefault {
			validation.Permissions = append(validation.Permissions, "default")
		}

		// Set access level based on flags
		if assoc.Flags != nil && len(*assoc.Flags) > 0 {
			for _, flag := range *assoc.Flags {
				switch flag {
				case V0043AssocFlagsDELETED:
					validation.HasAccess = false
					validation.Reason = "Association is deleted"
				default:
					// Other flags don't affect access
				}
			}
		}

		// Basic user permissions (since admin level is not available in association)
		validation.Permissions = append(validation.Permissions, "view", "submit")

		// Check for any restrictions
		if assoc.Max != nil && assoc.Max.Jobs != nil {
			if assoc.Max.Jobs.Active != nil && assoc.Max.Jobs.Active.Number != nil {
				maxJobs := *assoc.Max.Jobs.Active.Number
				if maxJobs == 0 {
					validation.Restrictions = append(validation.Restrictions, "no_job_submission")
					validation.HasAccess = false
					validation.Reason = "Job submission disabled for this association"
				} else {
					validation.Restrictions = append(validation.Restrictions, fmt.Sprintf("max_jobs:%d", maxJobs))
				}
			}
		}

		// Create association details
		validation.Association = &interfaces.UserAccountAssociation{
			UserName:    userName,
			AccountName: accountName,
			IsActive:    validation.HasAccess,
			IsDefault:   assoc.IsDefault != nil && *assoc.IsDefault,
		}

		if assoc.Cluster != nil {
			validation.Association.Cluster = *assoc.Cluster
		}
		if assoc.Partition != nil {
			validation.Association.Partition = *assoc.Partition
		}

		// Found valid association
		break
	}

	if !validation.HasAccess && validation.Reason == "" {
		validation.Reason = "User does not have active association with the account"
	}

	return validation, nil
}

// GetUserAccountAssociations retrieves detailed user account associations
func (u *UserManagerImpl) GetUserAccountAssociations(ctx context.Context, userName string, opts *interfaces.ListUserAccountAssociationsOptions) ([]*interfaces.UserAccountAssociation, error) {
	// Validate input first (cheap check)
	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}

	// Then check client initialization
	if u.client == nil || u.client.apiClient == nil || u.client.apiClient.ClientInterface == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Build parameters for associations query
	userNameCSV := userName
	params := &SlurmdbV0043GetAssociationsParams{
		User: &userNameCSV,
	}

	// Apply filters if provided
	if opts != nil {
		if len(opts.Accounts) > 0 {
			accountsCSV := strings.Join(opts.Accounts, ",")
			params.Account = &accountsCSV
		}
		if len(opts.Clusters) > 0 {
			clustersCSV := strings.Join(opts.Clusters, ",")
			params.Cluster = &clustersCSV
		}
		// Note: Partitions parameter doesn't exist in SlurmdbV0043GetAssociationsParams
	}

	resp, err := u.client.apiClient.SlurmdbV0043GetAssociationsWithResponse(ctx, params)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
	}

	// Check HTTP status and handle API errors
	if resp.StatusCode() != 200 {
		if resp.JSON200 != nil && resp.JSON200.Errors != nil && len(*resp.JSON200.Errors) > 0 {
			apiErrors := convertAPIErrors(*resp.JSON200.Errors)
			apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.43", apiErrors)
			return nil, apiError.SlurmError
		}
		return nil, errors.WrapHTTPError(resp.StatusCode(), nil, "v0.0.43")
	}

	// Check for valid response
	if resp.JSON200 == nil || resp.JSON200.Associations == nil {
		return []*interfaces.UserAccountAssociation{}, nil
	}

	// Convert associations to interface types
	associations := make([]*interfaces.UserAccountAssociation, 0)

	for _, apiAssoc := range resp.JSON200.Associations {
		// Skip if not for the requested user
		if apiAssoc.User != userName {
			continue
		}

		assoc := &interfaces.UserAccountAssociation{
			UserName: userName,
		}

		// Extract account information
		if apiAssoc.Account != nil {
			assoc.AccountName = *apiAssoc.Account
		}
		if apiAssoc.Cluster != nil {
			assoc.Cluster = *apiAssoc.Cluster
		}
		if apiAssoc.Partition != nil {
			assoc.Partition = *apiAssoc.Partition
		}

		// Set flags
		assoc.IsDefault = apiAssoc.IsDefault != nil && *apiAssoc.IsDefault
		assoc.IsActive = true // Assume active if returned by API

		// Extract role and permissions - admin level is not available in V0043Assoc
		// Check flags for deleted associations
		if apiAssoc.Flags != nil && len(*apiAssoc.Flags) > 0 {
			for _, flag := range *apiAssoc.Flags {
				switch flag {
				case V0043AssocFlagsDELETED:
					assoc.IsActive = false
				default:
					// Other flags don't affect basic permissions
				}
			}
		}

		// Basic user permissions (admin level is not available in V0043Assoc)
		assoc.Role = "user"
		assoc.Permissions = []string{"view", "submit"}

		// Extract limits from Max struct
		if apiAssoc.Max != nil {
			if apiAssoc.Max.Jobs != nil {
				if apiAssoc.Max.Jobs.Active != nil && apiAssoc.Max.Jobs.Active.Number != nil {
					assoc.MaxJobs = int(*apiAssoc.Max.Jobs.Active.Number)
				}
				if apiAssoc.Max.Jobs.Total != nil && apiAssoc.Max.Jobs.Total.Number != nil {
					assoc.MaxSubmitJobs = int(*apiAssoc.Max.Jobs.Total.Number)
				}
			}
			// TRES limits (CPU, Memory, Nodes) are in apiAssoc.Max.Tres
			if apiAssoc.Max.Tres != nil && apiAssoc.Max.Tres.Per != nil && apiAssoc.Max.Tres.Per.Job != nil {
				// Note: TRES limits conversion not implemented yet
				assoc.TRESLimits = make(map[string]int)
			}
		}

		// Extract usage information
		// Note: Shares not directly available in V0043Assoc
		if false { // Disabled
			assoc.SharesRaw = 0 // Shares field not available in V0043Assoc
		}

		// Apply additional filtering if options are provided
		if opts != nil {
			// Filter by role
			if len(opts.Roles) > 0 {
				roleMatch := false
				for _, role := range opts.Roles {
					if assoc.Role == role {
						roleMatch = true
						break
					}
				}
				if !roleMatch {
					continue
				}
			}

			// Filter by permissions
			if len(opts.Permissions) > 0 {
				permMatch := false
				for _, reqPerm := range opts.Permissions {
					for _, assocPerm := range assoc.Permissions {
						if assocPerm == reqPerm {
							permMatch = true
							break
						}
					}
					if permMatch {
						break
					}
				}
				if !permMatch {
					continue
				}
			}

			// Filter by active status
			if opts.ActiveOnly && !assoc.IsActive {
				continue
			}

			// Filter by default status
			if opts.DefaultOnly && !assoc.IsDefault {
				continue
			}

			// Note: CoordinatorsOnly filtering not available in this API version
			// since coordinator status is not directly available in V0043Assoc
		}

		associations = append(associations, assoc)
	}

	return associations, nil
}

// GetBulkUserAccounts retrieves accounts for multiple users in a single call
func (u *UserManagerImpl) GetBulkUserAccounts(ctx context.Context, userNames []string) (map[string][]*interfaces.UserAccount, error) {
	// Validate input first (cheap check)
	if len(userNames) == 0 {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "at least one user name is required", "userNames", userNames, nil)
	}

	// Validate all user names
	for i, userName := range userNames {
		if userName == "" {
			return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, fmt.Sprintf("user name at index %d is empty", i), "userNames["+strconv.Itoa(i)+"]", userName, nil)
		}
		if err := validateUserName(userName); err != nil {
			return nil, err
		}
	}

	// Limit bulk operations to prevent excessive load
	if len(userNames) > 100 {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "bulk operations limited to 100 users maximum", "userNames", len(userNames), nil)
	}

	// Then check client initialization
	if u.client == nil || u.client.apiClient == nil || u.client.apiClient.ClientInterface == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Use associations API to get accounts for all users
	// Convert to CSV string as expected by API
	userNamesCSV := strings.Join(userNames, ",")
	params := &SlurmdbV0043GetAssociationsParams{
		User: &userNamesCSV,
	}

	resp, err := u.client.apiClient.SlurmdbV0043GetAssociationsWithResponse(ctx, params)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
	}

	// Check HTTP status and handle API errors
	if resp.StatusCode() != 200 {
		if resp.JSON200 != nil && resp.JSON200.Errors != nil && len(*resp.JSON200.Errors) > 0 {
			apiErrors := convertAPIErrors(*resp.JSON200.Errors)
			apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.43", apiErrors)
			return nil, apiError.SlurmError
		}
		return nil, errors.WrapHTTPError(resp.StatusCode(), nil, "v0.0.43")
	}

	// Initialize result map
	result := make(map[string][]*interfaces.UserAccount)
	for _, userName := range userNames {
		result[userName] = make([]*interfaces.UserAccount, 0)
	}

	// Check for valid response
	if resp.JSON200 == nil || resp.JSON200.Associations == nil {
		return result, nil
	}

	// Process associations and group by user
	for _, assoc := range resp.JSON200.Associations {
		if assoc.User == "" || assoc.Account == nil {
			continue
		}

		userName := assoc.User

		// Skip if user not in requested list
		found := false
		for _, requestedUser := range userNames {
			if userName == requestedUser {
				found = true
				break
			}
		}
		if !found {
			continue
		}

		// Create user account entry
		userAccount := &interfaces.UserAccount{
			AccountName: *assoc.Account,
			IsDefault:   assoc.IsDefault != nil && *assoc.IsDefault,
			IsActive:    true, // Assume active if returned
		}

		// Add partition if available (Cluster and Role fields don't exist in UserAccount)
		if assoc.Partition != nil {
			userAccount.Partition = *assoc.Partition
		}

		// Note: Admin level and role are not available in V0043Assoc structure

		// Note: Shares information not available in UserAccount structure

		// Extract limits from Max struct if available
		if assoc.Max != nil && assoc.Max.Jobs != nil {
			if assoc.Max.Jobs.Active != nil && assoc.Max.Jobs.Active.Number != nil {
				userAccount.MaxJobs = int(*assoc.Max.Jobs.Active.Number)
			}
			if assoc.Max.Jobs.Total != nil && assoc.Max.Jobs.Total.Number != nil {
				userAccount.MaxSubmitJobs = int(*assoc.Max.Jobs.Total.Number)
			}
		}
		// MaxWallTime is not directly available in this structure

		// Add to user's account list
		result[userName] = append(result[userName], userAccount)
	}

	return result, nil
}

// GetBulkAccountUsers retrieves users for multiple accounts in a single call
func (u *UserManagerImpl) GetBulkAccountUsers(ctx context.Context, accountNames []string) (map[string][]*interfaces.UserAccountAssociation, error) {
	// Validate input first (cheap check)
	if len(accountNames) == 0 {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "at least one account name is required", "accountNames", accountNames, nil)
	}

	// Validate all account names
	for i, accountName := range accountNames {
		if accountName == "" {
			return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, fmt.Sprintf("account name at index %d is empty", i), "accountNames["+strconv.Itoa(i)+"]", accountName, nil)
		}
		if err := validateAccountContext(accountName); err != nil {
			return nil, err
		}
	}

	// Limit bulk operations to prevent excessive load
	if len(accountNames) > 100 {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "bulk operations limited to 100 accounts maximum", "accountNames", len(accountNames), nil)
	}

	// Then check client initialization
	if u.client == nil || u.client.apiClient == nil || u.client.apiClient.ClientInterface == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Use associations API to get users for all accounts
	// Convert to CSV string as expected by API
	accountNamesCSV := strings.Join(accountNames, ",")
	params := &SlurmdbV0043GetAssociationsParams{
		Account: &accountNamesCSV,
	}

	resp, err := u.client.apiClient.SlurmdbV0043GetAssociationsWithResponse(ctx, params)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.43")
	}

	// Check HTTP status and handle API errors
	if resp.StatusCode() != 200 {
		if resp.JSON200 != nil && resp.JSON200.Errors != nil && len(*resp.JSON200.Errors) > 0 {
			apiErrors := convertAPIErrors(*resp.JSON200.Errors)
			apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.43", apiErrors)
			return nil, apiError.SlurmError
		}
		return nil, errors.WrapHTTPError(resp.StatusCode(), nil, "v0.0.43")
	}

	// Initialize result map
	result := make(map[string][]*interfaces.UserAccountAssociation)
	for _, accountName := range accountNames {
		result[accountName] = make([]*interfaces.UserAccountAssociation, 0)
	}

	// Check for valid response
	if resp.JSON200 == nil || resp.JSON200.Associations == nil {
		return result, nil
	}

	// Process associations and group by account
	for _, apiAssoc := range resp.JSON200.Associations {
		if apiAssoc.Account == nil || apiAssoc.User == "" {
			continue
		}

		accountName := *apiAssoc.Account

		// Skip if account not in requested list
		found := false
		for _, requestedAccount := range accountNames {
			if accountName == requestedAccount {
				found = true
				break
			}
		}
		if !found {
			continue
		}

		// Create user account association
		assoc := &interfaces.UserAccountAssociation{
			UserName:    apiAssoc.User,
			AccountName: accountName,
			IsActive:    true, // Assume active if returned
			IsDefault:   apiAssoc.IsDefault != nil && *apiAssoc.IsDefault,
		}

		// Add cluster and partition if available
		if apiAssoc.Cluster != nil {
			assoc.Cluster = *apiAssoc.Cluster
		}
		if apiAssoc.Partition != nil {
			assoc.Partition = *apiAssoc.Partition
		}

		// Extract role and permissions - admin level is not available in V0043Assoc
		// Check flags for deleted associations
		if apiAssoc.Flags != nil && len(*apiAssoc.Flags) > 0 {
			for _, flag := range *apiAssoc.Flags {
				switch flag {
				case V0043AssocFlagsDELETED:
					assoc.IsActive = false
				default:
					// Other flags don't affect basic permissions
				}
			}
		}

		// Basic user permissions (admin level is not available in V0043Assoc)
		assoc.Role = "user"
		assoc.Permissions = []string{"view", "submit"}

		// Extract limits from Max struct
		if apiAssoc.Max != nil {
			if apiAssoc.Max.Jobs != nil {
				if apiAssoc.Max.Jobs.Active != nil && apiAssoc.Max.Jobs.Active.Number != nil {
					assoc.MaxJobs = int(*apiAssoc.Max.Jobs.Active.Number)
				}
				if apiAssoc.Max.Jobs.Total != nil && apiAssoc.Max.Jobs.Total.Number != nil {
					assoc.MaxSubmitJobs = int(*apiAssoc.Max.Jobs.Total.Number)
				}
			}
			// TRES limits (CPU, Memory, Nodes) are in apiAssoc.Max.Tres
			if apiAssoc.Max.Tres != nil && apiAssoc.Max.Tres.Per != nil && apiAssoc.Max.Tres.Per.Job != nil {
				// Note: TRES limits conversion not implemented yet
				assoc.TRESLimits = make(map[string]int)
			}
		}

		// Extract share information (SharesAllocated not available in UserAccountAssociation)
		if apiAssoc.SharesRaw != nil {
			assoc.SharesRaw = int(*apiAssoc.SharesRaw)
		}

		// Add to account's user list
		result[accountName] = append(result[accountName], assoc)
	}

	return result, nil
}

// Helper function to validate user name format
func validateUserName(userName string) error {
	if userName == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name cannot be empty", "userName", userName, nil)
	}

	// Basic validation - user names should be alphanumeric with underscores and hyphens
	for _, char := range userName {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '_' || char == '-') {
			return errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name contains invalid characters", "userName", userName, nil)
		}
	}

	return nil
}

// Helper function to validate account name in user context
func validateAccountContext(accountName string) error {
	if accountName == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name cannot be empty", "accountName", accountName, nil)
	}

	// Account names should follow similar validation as user names
	for _, char := range accountName {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '_' || char == '-') {
			return errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name contains invalid characters", "accountName", accountName, nil)
		}
	}

	return nil
}

// convertAPIErrors converts SLURM API errors to our error format
func convertAPIErrors(apiErrors []V0043OpenapiError) []errors.SlurmAPIErrorDetail {
	result := make([]errors.SlurmAPIErrorDetail, len(apiErrors))
	for i, apiErr := range apiErrors {
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

		result[i] = errors.SlurmAPIErrorDetail{
			ErrorNumber: errorNumber,
			ErrorCode:   errorCode,
			Source:      source,
			Description: description,
		}
	}
	return result
}

// calculateJobSizeFactor calculates a job size factor based on resources
func calculateJobSizeFactor(jobSubmission *interfaces.JobSubmission) float64 {
	factor := 100.0 // Base factor

	// Adjust based on CPUs
	if jobSubmission.CPUs > 0 {
		factor += float64(jobSubmission.CPUs) * 10
	}

	// Adjust based on memory (in MB)
	if jobSubmission.Memory > 0 {
		factor += float64(jobSubmission.Memory) / 1024 // Convert to GB and add
	}

	// Adjust based on time limit (in minutes)
	if jobSubmission.TimeLimit > 0 {
		factor += float64(jobSubmission.TimeLimit) / 60 // Convert to hours and add
	}

	return factor
}

// convertAPIUserToInterface converts V0043User to interfaces.User
func convertAPIUserToInterface(apiUser *V0043User) (*interfaces.User, error) {
	user := &interfaces.User{}

	// Basic user information
	if apiUser.Name != "" {
		user.Name = apiUser.Name
	}

	// Admin level
	if apiUser.AdministratorLevel != nil && len(*apiUser.AdministratorLevel) > 0 {
		user.AdminLevel = string((*apiUser.AdministratorLevel)[0])
	}

	// Default account and WCKEY
	if apiUser.Default != nil {
		if apiUser.Default.Account != nil {
			user.DefaultAccount = *apiUser.Default.Account
		}
		if apiUser.Default.Wckey != nil {
			user.DefaultWCKey = *apiUser.Default.Wckey
		}
	}

	// Associations
	if apiUser.Associations != nil {
		associations := make([]interfaces.UserAssociation, 0, len(*apiUser.Associations))
		for _, assoc := range *apiUser.Associations {
			userAssoc := interfaces.UserAssociation{
				Account:       getStringValue(assoc.Account),
				Cluster:       getStringValue(assoc.Cluster),
				Partition:     getStringValue(assoc.Partition),
				QoS:           []string{}, // QoS not available in V0043AssocShort
				DefaultQoS:    "",         // Default not available in V0043AssocShort
				IsDefault:     false,      // IsDefault not available in V0043AssocShort
				MaxJobs:       0,          // Max not available in V0043AssocShort
				MaxSubmitJobs: 0,          // Max not available in V0043AssocShort
			}
			// TRES limits not available in V0043AssocShort structure
			associations = append(associations, userAssoc)
		}
		user.Associations = associations
	}

	// User flags
	// Note: Flags field not available in interfaces.User
	// Flag conversion code removed as the result is not used

	return user, nil
}

// Helper functions for safe value extraction
func getStringValue(ptr *string) string {
	if ptr != nil {
		return *ptr
	}
	return ""
}

// filterUsers applies client-side filtering to the user list
func filterUsers(users []interfaces.User, opts *interfaces.ListUsersOptions) []interfaces.User {
	if opts == nil {
		return users
	}

	filtered := make([]interfaces.User, 0, len(users))
	for _, user := range users {
		// Filter by names
		if len(opts.Names) > 0 {
			found := false
			for _, name := range opts.Names {
				if user.Name == name {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Note: DefaultAccount filtering not available in ListUsersOptions

		// Note: AdminLevel filtering not available in ListUsersOptions

		// Note: IncludeDeleted and Flags fields not available in interfaces

		filtered = append(filtered, user)
	}

	return filtered
}

// stringToInt converts string to int
func stringToInt(s string) (int, error) {
	// Try to parse the string as an integer
	if s == "" {
		return 0, nil
	}
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}

// Create creates a new user
func (u *UserManagerImpl) Create(ctx context.Context, user *interfaces.UserCreate) (*interfaces.UserCreateResponse, error) {
	return nil, errors.NewNotImplementedError("Create", "v0.0.43")
}

// Update updates a user
func (u *UserManagerImpl) Update(ctx context.Context, userName string, update *interfaces.UserUpdate) error {
	return errors.NewNotImplementedError("Update", "v0.0.43")
}

// Delete deletes a user
func (u *UserManagerImpl) Delete(ctx context.Context, userName string) error {
	return errors.NewNotImplementedError("Delete", "v0.0.43")
}

// CreateAssociation creates a user-account association
func (u *UserManagerImpl) CreateAssociation(ctx context.Context, accountName string, opts *interfaces.AssociationOptions) (*interfaces.AssociationCreateResponse, error) {
	return nil, errors.NewNotImplementedError("CreateAssociation", "v0.0.43")
}
