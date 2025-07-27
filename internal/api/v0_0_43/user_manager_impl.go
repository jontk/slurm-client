package v0_0_43

import (
	"context"
	"fmt"
	"strings"

	"github.com/jontk/slurm-client/internal/interfaces"
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
	if u.client == nil || u.client.apiClient == nil {
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
	users := make([]interfaces.User, 0, len(*resp.JSON200.Users))
	for _, apiUser := range *resp.JSON200.Users {
		user, err := convertAPIUserToInterface(apiUser)
		if err != nil {
			conversionErr := errors.NewClientError(errors.ErrorCodeServerInternal, "Failed to convert user data")
			conversionErr.Cause = err
			conversionErr.Details = fmt.Sprintf("Error converting user %v", apiUser.Name)
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
	if u.client == nil || u.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
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
	if resp.JSON200 == nil || resp.JSON200.Users == nil || len(*resp.JSON200.Users) == 0 {
		return nil, errors.NewClientError(errors.ErrorCodeNotFound, "User not found", fmt.Sprintf("User '%s' not found", userName))
	}

	// Convert the first user in the response
	user, err := convertAPIUserToInterface((*resp.JSON200.Users)[0])
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
	if u.client == nil || u.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
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
			UserName:      userName,
			AccountName:   assoc.Account,
			IsDefault:     assoc.IsDefault,
			Partition:     assoc.Partition,
			Cluster:       assoc.Cluster,
			QoS:           assoc.QoS,
			SharesRaw:     assoc.SharesRaw,
			GrpTRES:       assoc.GrpTRES,
			GrpJobs:       assoc.GrpJobs,
			GrpSubmitJobs: assoc.GrpSubmitJobs,
			GrpWall:       assoc.GrpWall,
			MaxTRES:       assoc.MaxTRES,
			MaxJobs:       assoc.MaxJobs,
			MaxSubmitJobs: assoc.MaxSubmitJobs,
			MaxWallDurationPerJob: assoc.MaxWallDurationPerJob,
			// Priority is available in association
			Priority: assoc.Priority,
		}
		userAccounts = append(userAccounts, userAccount)
	}

	return userAccounts, nil
}

// GetUserQuotas retrieves quota information for a user
func (u *UserManagerImpl) GetUserQuotas(ctx context.Context, userName string) (*interfaces.UserQuota, error) {
	if u.client == nil || u.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}

	// Get user details including associations
	user, err := u.Get(ctx, userName)
	if err != nil {
		return nil, err
	}

	// Aggregate quotas from all associations
	userQuota := &interfaces.UserQuota{
		UserName:        userName,
		AccountQuotas:   make(map[string]*interfaces.AccountSpecificQuota),
		EffectiveQuotas: &interfaces.EffectiveQuotas{},
	}

	// Track effective quotas (most restrictive across all accounts)
	var minMaxJobs *int
	var minMaxSubmitJobs *int
	var minMaxWallTime *int
	var minMaxCPUs *int
	var minMaxNodes *int

	// Process each association
	for _, assoc := range user.Associations {
		accountQuota := &interfaces.AccountSpecificQuota{
			AccountName:           assoc.Account,
			SharesRaw:             assoc.SharesRaw,
			GrpTRES:               assoc.GrpTRES,
			GrpJobs:               assoc.GrpJobs,
			GrpSubmitJobs:         assoc.GrpSubmitJobs,
			GrpWall:               assoc.GrpWall,
			MaxTRES:               assoc.MaxTRES,
			MaxJobs:               assoc.MaxJobs,
			MaxSubmitJobs:         assoc.MaxSubmitJobs,
			MaxWallDurationPerJob: assoc.MaxWallDurationPerJob,
			QoS:                   assoc.QoS,
			Priority:              assoc.Priority,
		}

		userQuota.AccountQuotas[assoc.Account] = accountQuota

		// Update effective quotas (most restrictive)
		if assoc.MaxJobs != nil {
			if minMaxJobs == nil || *assoc.MaxJobs < *minMaxJobs {
				minMaxJobs = assoc.MaxJobs
			}
		}
		if assoc.MaxSubmitJobs != nil {
			if minMaxSubmitJobs == nil || *assoc.MaxSubmitJobs < *minMaxSubmitJobs {
				minMaxSubmitJobs = assoc.MaxSubmitJobs
			}
		}
		if assoc.MaxWallDurationPerJob != nil {
			if minMaxWallTime == nil || *assoc.MaxWallDurationPerJob < *minMaxWallTime {
				minMaxWallTime = assoc.MaxWallDurationPerJob
			}
		}

		// Extract CPU and node limits from MaxTRES if available
		if assoc.MaxTRES != nil {
			if cpuLimit, ok := assoc.MaxTRES["cpu"]; ok {
				cpuInt := int(cpuLimit)
				if minMaxCPUs == nil || cpuInt < *minMaxCPUs {
					minMaxCPUs = &cpuInt
				}
			}
			if nodeLimit, ok := assoc.MaxTRES["node"]; ok {
				nodeInt := int(nodeLimit)
				if minMaxNodes == nil || nodeInt < *minMaxNodes {
					minMaxNodes = &nodeInt
				}
			}
		}
	}

	// Set effective quotas
	userQuota.EffectiveQuotas.MaxJobs = minMaxJobs
	userQuota.EffectiveQuotas.MaxSubmitJobs = minMaxSubmitJobs
	userQuota.EffectiveQuotas.MaxWallDurationPerJob = minMaxWallTime
	userQuota.EffectiveQuotas.MaxCPUsPerJob = minMaxCPUs
	userQuota.EffectiveQuotas.MaxNodesPerJob = minMaxNodes

	return userQuota, nil
}

// GetUserDefaultAccount retrieves the default account for a user
func (u *UserManagerImpl) GetUserDefaultAccount(ctx context.Context, userName string) (*interfaces.Account, error) {
	if u.client == nil || u.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
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

	return nil, errors.NewClientError(errors.ErrorCodeNotFound, "No default account found", fmt.Sprintf("User '%s' has no account associations", userName))
}

// GetUserFairShare retrieves fair-share information for a user
func (u *UserManagerImpl) GetUserFairShare(ctx context.Context, userName string) (*interfaces.UserFairShare, error) {
	if u.client == nil || u.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}

	// TODO: Implement actual API call to retrieve user fair-share information
	// This would involve querying SLURM's fair-share database for user priority factors
	// including fair-share factor, normalized shares, effective usage, etc.
	return nil, errors.NewNotImplementedError("user fair-share retrieval", "v0.0.43")
}

// CalculateJobPriority calculates job priority for a user and job submission
func (u *UserManagerImpl) CalculateJobPriority(ctx context.Context, userName string, jobSubmission *interfaces.JobSubmission) (*interfaces.JobPriorityInfo, error) {
	if u.client == nil || u.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

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

	// TODO: Implement actual API call to calculate job priority
	// This would involve:
	// 1. Getting user's fair-share information
	// 2. Getting account and QoS priority factors
	// 3. Calculating age, job size, and other priority components
	// 4. Returning a complete priority breakdown
	return nil, errors.NewNotImplementedError("job priority calculation", "v0.0.43")
}

// ValidateUserAccountAccess validates user access to a specific account
func (u *UserManagerImpl) ValidateUserAccountAccess(ctx context.Context, userName, accountName string) (*interfaces.UserAccessValidation, error) {
	if u.client == nil || u.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}

	if accountName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name is required", "accountName", accountName, nil)
	}

	// TODO: Implement actual API call to validate user-account access
	// This would involve checking association tables, permissions, and active status
	return nil, errors.NewNotImplementedError("user-account access validation", "v0.0.43")
}

// GetUserAccountAssociations retrieves detailed user account associations
func (u *UserManagerImpl) GetUserAccountAssociations(ctx context.Context, userName string, opts *interfaces.ListUserAccountAssociationsOptions) ([]*interfaces.UserAccountAssociation, error) {
	if u.client == nil || u.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}

	// TODO: Implement actual API call to retrieve detailed user account associations
	// This would include roles, permissions, quotas, and usage information
	return nil, errors.NewNotImplementedError("user account associations retrieval", "v0.0.43")
}

// GetBulkUserAccounts retrieves accounts for multiple users in a single call
func (u *UserManagerImpl) GetBulkUserAccounts(ctx context.Context, userNames []string) (map[string][]*interfaces.UserAccount, error) {
	if u.client == nil || u.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if len(userNames) == 0 {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "at least one user name is required", "userNames", userNames, nil)
	}

	// Validate all user names
	for i, userName := range userNames {
		if userName == "" {
			return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, fmt.Sprintf("user name at index %d is empty", i), "userNames[" + fmt.Sprintf("%d", i) + "]", userName, nil)
		}
		if err := validateUserName(userName); err != nil {
			return nil, err
		}
	}

	// Limit bulk operations to prevent excessive load
	if len(userNames) > 100 {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "bulk operations limited to 100 users maximum", "userNames", len(userNames), nil)
	}

	// TODO: Implement actual API call for bulk user account retrieval
	return nil, errors.NewNotImplementedError("bulk user accounts retrieval", "v0.0.43")
}

// GetBulkAccountUsers retrieves users for multiple accounts in a single call
func (u *UserManagerImpl) GetBulkAccountUsers(ctx context.Context, accountNames []string) (map[string][]*interfaces.UserAccountAssociation, error) {
	if u.client == nil || u.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if len(accountNames) == 0 {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "at least one account name is required", "accountNames", accountNames, nil)
	}

	// Validate all account names
	for i, accountName := range accountNames {
		if accountName == "" {
			return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, fmt.Sprintf("account name at index %d is empty", i), "accountNames[" + fmt.Sprintf("%d", i) + "]", accountName, nil)
		}
		if err := validateAccountContext(accountName); err != nil {
			return nil, err
		}
	}

	// Limit bulk operations to prevent excessive load
	if len(accountNames) > 100 {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "bulk operations limited to 100 accounts maximum", "accountNames", len(accountNames), nil)
	}

	// TODO: Implement actual API call for bulk account users retrieval
	return nil, errors.NewNotImplementedError("bulk account users retrieval", "v0.0.43")
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

// convertAPIUserToInterface converts V0043UserShort to interfaces.User
func convertAPIUserToInterface(apiUser V0043UserShort) (*interfaces.User, error) {
	user := &interfaces.User{}

	// V0043UserShort has minimal fields, we need to extract the name from somewhere else
	// For now, we'll use default values since UserShort doesn't have name
	
	// Admin level
	if apiUser.Adminlevel != nil && len(*apiUser.Adminlevel) > 0 {
		user.AdminLevel = string((*apiUser.Adminlevel)[0])
	}

	// Default account and WCKEY
	if apiUser.Defaultaccount != nil {
		user.DefaultAccount = *apiUser.Defaultaccount
	}
	if apiUser.Defaultwckey != nil {
		user.DefaultWCKey = *apiUser.Defaultwckey
	}

	// V0043UserShort doesn't have associations, flags, or name
	// These would need to be fetched separately with a more detailed API call

	return user, nil
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

		// Filter by default account
		if opts.DefaultAccount != "" && user.DefaultAccount != opts.DefaultAccount {
			continue
		}

		// Filter by admin level
		if opts.AdminLevel != "" && user.AdminLevel != opts.AdminLevel {
			continue
		}

		// Filter by include deleted
		if !opts.IncludeDeleted {
			// Check if user has DELETED flag
			hasDeletedFlag := false
			for _, flag := range user.Flags {
				if flag == "DELETED" {
					hasDeletedFlag = true
					break
				}
			}
			if hasDeletedFlag {
				continue
			}
		}

		filtered = append(filtered, user)
	}

	return filtered
}


// convertTRESArrayToMap converts TRES array to a map
func convertTRESArrayToMap(tres []V0043Tres) map[string]int64 {
	result := make(map[string]int64)
	for _, t := range tres {
		if t.Type != nil && t.Count != nil {
			result[*t.Type] = int64(*t.Count)
		}
	}
	return result
}