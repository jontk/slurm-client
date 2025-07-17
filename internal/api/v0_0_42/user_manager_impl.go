package v0_0_42

import (
	"context"
	"fmt"

	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
)

// UserManagerImpl implements the UserManager interface for v0.0.42
// Some user management features are available but limited compared to v0.0.43
type UserManagerImpl struct {
	client *WrapperClient
}

// NewUserManagerImpl creates a new UserManagerImpl
func NewUserManagerImpl(client *WrapperClient) *UserManagerImpl {
	return &UserManagerImpl{
		client: client,
	}
}

// List lists users with optional filtering
// Limited features compared to v0.0.43
func (u *UserManagerImpl) List(ctx context.Context, opts *interfaces.ListUsersOptions) (*interfaces.UserList, error) {
	if u.client == nil || u.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// v0.0.42 has limited user management support compared to v0.0.43
	// Enhanced features like WithFairShare, WithAssociations are not available
	return nil, errors.NewNotImplementedError("user listing", "v0.0.42")
}

// Get retrieves a specific user by name
// Basic support available in v0.0.42
func (u *UserManagerImpl) Get(ctx context.Context, userName string) (*interfaces.User, error) {
	if u.client == nil || u.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}

	// v0.0.42 has basic user retrieval support
	return nil, errors.NewNotImplementedError("user retrieval", "v0.0.42")
}

// GetUserAccounts retrieves all account associations for a user
// Limited support in v0.0.42
func (u *UserManagerImpl) GetUserAccounts(ctx context.Context, userName string) ([]*interfaces.UserAccount, error) {
	if u.client == nil || u.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}

	// v0.0.42 has limited user account association support
	return nil, errors.NewNotImplementedError("user accounts retrieval not fully supported", "v0.0.42")
}

// GetUserQuotas retrieves quota information for a user
// Limited quota support in v0.0.42
func (u *UserManagerImpl) GetUserQuotas(ctx context.Context, userName string) (*interfaces.UserQuota, error) {
	if u.client == nil || u.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}

	// v0.0.42 has limited user quota information support
	return nil, errors.NewNotImplementedError("user quotas retrieval not fully supported", "v0.0.42")
}

// GetUserDefaultAccount retrieves the default account for a user
// Basic support available in v0.0.42
func (u *UserManagerImpl) GetUserDefaultAccount(ctx context.Context, userName string) (*interfaces.Account, error) {
	if u.client == nil || u.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}

	// v0.0.42 has basic default account support
	return nil, errors.NewNotImplementedError("user default account retrieval not fully supported", "v0.0.42")
}

// GetUserFairShare retrieves fair-share information for a user
// Limited fair-share support in v0.0.42
func (u *UserManagerImpl) GetUserFairShare(ctx context.Context, userName string) (*interfaces.UserFairShare, error) {
	if u.client == nil || u.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}

	// v0.0.42 has limited fair-share information support
	return nil, errors.NewNotImplementedError("user fair-share retrieval not fully supported", "v0.0.42")
}

// CalculateJobPriority calculates job priority for a user and job submission
// Limited priority calculation support in v0.0.42
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

	// Basic validation of job submission
	if jobSubmission.Script == "" && jobSubmission.Command == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "job script or command is required", "jobSubmission.Script/Command", fmt.Sprintf("Script: %s, Command: %s", jobSubmission.Script, jobSubmission.Command), nil)
	}

	// v0.0.42 has limited job priority calculation support
	return nil, errors.NewNotImplementedError("job priority calculation not fully supported", "v0.0.42")
}