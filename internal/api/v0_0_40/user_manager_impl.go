package v0_0_40

import (
	"context"

	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
)

// UserManagerImpl implements the UserManager interface for v0.0.40
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
	// v0.0.40 is an older API version with limited support
	// User management was added in later versions
	return nil, errors.NewNotImplementedError("user listing", "v0.0.40")
}

// Get retrieves a specific user by name
func (u *UserManagerImpl) Get(ctx context.Context, userName string) (*interfaces.User, error) {
	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}
	return nil, errors.NewNotImplementedError("user retrieval", "v0.0.40")
}

// GetUserAccounts retrieves all account associations for a user
func (u *UserManagerImpl) GetUserAccounts(ctx context.Context, userName string) ([]*interfaces.UserAccount, error) {
	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}
	return nil, errors.NewNotImplementedError("user accounts retrieval", "v0.0.40")
}

// GetUserQuotas retrieves quota information for a user
func (u *UserManagerImpl) GetUserQuotas(ctx context.Context, userName string) (*interfaces.UserQuota, error) {
	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}
	return nil, errors.NewNotImplementedError("user quotas retrieval", "v0.0.40")
}

// GetUserDefaultAccount retrieves the default account for a user
func (u *UserManagerImpl) GetUserDefaultAccount(ctx context.Context, userName string) (*interfaces.Account, error) {
	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}
	return nil, errors.NewNotImplementedError("user default account retrieval", "v0.0.40")
}

// GetUserFairShare retrieves fair-share information for a user
func (u *UserManagerImpl) GetUserFairShare(ctx context.Context, userName string) (*interfaces.UserFairShare, error) {
	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}
	return nil, errors.NewNotImplementedError("user fair-share retrieval", "v0.0.40")
}

// CalculateJobPriority calculates job priority for a user and job submission
func (u *UserManagerImpl) CalculateJobPriority(ctx context.Context, userName string, jobSubmission *interfaces.JobSubmission) (*interfaces.JobPriorityInfo, error) {
	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}
	if jobSubmission == nil {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "job submission data is required", "jobSubmission", jobSubmission, nil)
	}
	return nil, errors.NewNotImplementedError("job priority calculation", "v0.0.40")
}

// ValidateUserAccountAccess validates user access to a specific account
func (u *UserManagerImpl) ValidateUserAccountAccess(ctx context.Context, userName, accountName string) (*interfaces.UserAccessValidation, error) {
	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}
	if accountName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "account name is required", "accountName", accountName, nil)
	}
	return nil, errors.NewNotImplementedError("user-account access validation", "v0.0.40")
}

// GetUserAccountAssociations retrieves detailed user account associations
func (u *UserManagerImpl) GetUserAccountAssociations(ctx context.Context, userName string, opts *interfaces.ListUserAccountAssociationsOptions) ([]*interfaces.UserAccountAssociation, error) {
	if userName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "user name is required", "userName", userName, nil)
	}
	return nil, errors.NewNotImplementedError("user account associations retrieval", "v0.0.40")
}

// GetBulkUserAccounts retrieves accounts for multiple users in a single call
func (u *UserManagerImpl) GetBulkUserAccounts(ctx context.Context, userNames []string) (map[string][]*interfaces.UserAccount, error) {
	if len(userNames) == 0 {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "at least one user name is required", "userNames", userNames, nil)
	}
	return nil, errors.NewNotImplementedError("bulk user accounts retrieval", "v0.0.40")
}

// GetBulkAccountUsers retrieves users for multiple accounts in a single call
func (u *UserManagerImpl) GetBulkAccountUsers(ctx context.Context, accountNames []string) (map[string][]*interfaces.UserAccountAssociation, error) {
	if len(accountNames) == 0 {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "at least one account name is required", "accountNames", accountNames, nil)
	}
	return nil, errors.NewNotImplementedError("bulk account users retrieval", "v0.0.40")
}