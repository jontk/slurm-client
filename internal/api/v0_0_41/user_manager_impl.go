package v0_0_41

import (
	"context"

	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
)

// UserManagerImpl implements the UserManager interface for v0.0.41
// Most user management features are not available in v0.0.41
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
// This feature is not available in v0.0.41
func (u *UserManagerImpl) List(ctx context.Context, opts *interfaces.ListUsersOptions) (*interfaces.UserList, error) {
	if u.client == nil || u.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// v0.0.41 has very limited user management support
	// Most user features are not available
	return nil, errors.NewNotImplementedError("user listing not supported", "v0.0.41")
}

// Get retrieves a specific user by name
// This feature is not available in v0.0.41
func (u *UserManagerImpl) Get(ctx context.Context, userName string) (*interfaces.User, error) {
	if u.client == nil || u.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}

	// v0.0.41 has minimal user retrieval support
	return nil, errors.NewNotImplementedError("user retrieval not supported", "v0.0.41")
}

// GetUserAccounts retrieves all account associations for a user
// This feature is not available in v0.0.41
func (u *UserManagerImpl) GetUserAccounts(ctx context.Context, userName string) ([]*interfaces.UserAccount, error) {
	if u.client == nil || u.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}

	return nil, errors.NewNotImplementedError("user accounts not supported", "v0.0.41")
}

// GetUserQuotas retrieves quota information for a user
// This feature is not available in v0.0.41
func (u *UserManagerImpl) GetUserQuotas(ctx context.Context, userName string) (*interfaces.UserQuota, error) {
	if u.client == nil || u.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}

	return nil, errors.NewNotImplementedError("user quotas not supported", "v0.0.41")
}

// GetUserDefaultAccount retrieves the default account for a user
// This feature is not available in v0.0.41
func (u *UserManagerImpl) GetUserDefaultAccount(ctx context.Context, userName string) (*interfaces.Account, error) {
	if u.client == nil || u.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}

	return nil, errors.NewNotImplementedError("user default account not supported", "v0.0.41")
}

// GetUserFairShare retrieves fair-share information for a user
// This feature is not available in v0.0.41
func (u *UserManagerImpl) GetUserFairShare(ctx context.Context, userName string) (*interfaces.UserFairShare, error) {
	if u.client == nil || u.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}

	return nil, errors.NewNotImplementedError("user fair-share not supported", "v0.0.41")
}

// CalculateJobPriority calculates job priority for a user and job submission
// This feature is not available in v0.0.41
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

	return nil, errors.NewNotImplementedError("job priority calculation not supported", "v0.0.41")
}