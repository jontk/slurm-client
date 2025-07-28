package v0_0_43

import (
	"strconv"
	"strings"

	"github.com/jontk/slurm-client/internal/common/types"
	api "github.com/jontk/slurm-client/internal/api/v0_0_43"
)

// convertAPIAccountToCommon converts a v0.0.43 API Account to common Account type
func (a *AccountAdapter) convertAPIAccountToCommon(apiAccount api.V0043AccountInfo) (*types.Account, error) {
	account := &types.Account{}

	// Basic fields
	if apiAccount.Name != nil {
		account.Name = *apiAccount.Name
	}
	if apiAccount.Description != nil {
		account.Description = *apiAccount.Description
	}
	if apiAccount.Organization != nil {
		account.Organization = *apiAccount.Organization
	}

	// Coordinators
	if apiAccount.Coordinators != nil && len(*apiAccount.Coordinators) > 0 {
		coordinators := make([]string, len(*apiAccount.Coordinators))
		for i, coord := range *apiAccount.Coordinators {
			if coord.Name != nil {
				coordinators[i] = *coord.Name
			}
		}
		account.Coordinators = coordinators
	}

	// QoS settings
	if apiAccount.DefaultQos != nil {
		account.DefaultQoS = *apiAccount.DefaultQos
	}
	if apiAccount.QosList != nil {
		account.QoSList = *apiAccount.QosList
	}

	// Parent account
	if apiAccount.Parent != nil {
		account.ParentName = *apiAccount.Parent
	}

	// Shares and priority
	if apiAccount.SharesRaw != nil {
		account.SharesRaw = *apiAccount.SharesRaw
	}
	if apiAccount.Priority != nil && apiAccount.Priority.Number != nil {
		account.Priority = *apiAccount.Priority.Number
	}

	// Job limits
	if apiAccount.MaxJobs != nil && apiAccount.MaxJobs.Number != nil {
		account.MaxJobs = *apiAccount.MaxJobs.Number
	}
	if apiAccount.MaxJobsPerUser != nil && apiAccount.MaxJobsPerUser.Number != nil {
		account.MaxJobsPerUser = *apiAccount.MaxJobsPerUser.Number
	}
	if apiAccount.MaxSubmitJobs != nil && apiAccount.MaxSubmitJobs.Number != nil {
		account.MaxSubmitJobs = *apiAccount.MaxSubmitJobs.Number
	}

	// Time limits
	if apiAccount.MaxWallTime != nil && apiAccount.MaxWallTime.Number != nil {
		account.MaxWallTime = *apiAccount.MaxWallTime.Number
	}

	// Resource limits
	if apiAccount.MaxNodes != nil && apiAccount.MaxNodes.Number != nil {
		account.MaxNodes = *apiAccount.MaxNodes.Number
	}
	if apiAccount.MaxCpus != nil && apiAccount.MaxCpus.Number != nil {
		account.MaxCPUs = *apiAccount.MaxCpus.Number
	}

	// Group limits
	if apiAccount.GrpJobs != nil && apiAccount.GrpJobs.Number != nil {
		account.GrpJobs = *apiAccount.GrpJobs.Number
	}
	if apiAccount.GrpNodes != nil && apiAccount.GrpNodes.Number != nil {
		account.GrpNodes = *apiAccount.GrpNodes.Number
	}
	if apiAccount.GrpCpus != nil && apiAccount.GrpCpus.Number != nil {
		account.GrpCPUs = *apiAccount.GrpCpus.Number
	}
	if apiAccount.GrpSubmitJobs != nil && apiAccount.GrpSubmitJobs.Number != nil {
		account.GrpSubmitJobs = *apiAccount.GrpSubmitJobs.Number
	}
	if apiAccount.GrpWallTime != nil && apiAccount.GrpWallTime.Number != nil {
		account.GrpWallTime = *apiAccount.GrpWallTime.Number
	}

	// Flags
	if apiAccount.Flags != nil {
		for _, flag := range *apiAccount.Flags {
			switch flag {
			case api.V0043AccountInfoFlagsDELETED:
				account.Deleted = true
			case api.V0043AccountInfoFlagsDEFAULT:
				account.IsDefault = true
			}
		}
	}

	return account, nil
}

// convertCommonAccountCreateToAPI converts common AccountCreate type to v0.0.43 API format
func (a *AccountAdapter) convertCommonAccountCreateToAPI(create *types.AccountCreate) (*api.V0043AccountInfo, error) {
	apiAccount := &api.V0043AccountInfo{}

	// Required fields
	apiAccount.Name = &create.Name

	// Basic fields
	if create.Description != "" {
		apiAccount.Description = &create.Description
	}
	if create.Organization != "" {
		apiAccount.Organization = &create.Organization
	}

	// Coordinators
	if len(create.Coordinators) > 0 {
		coordinators := make([]api.V0043AccountCoordinator, len(create.Coordinators))
		for i, coord := range create.Coordinators {
			coordinators[i] = api.V0043AccountCoordinator{
				Name: &coord,
			}
		}
		apiAccount.Coordinators = &coordinators
	}

	// QoS settings
	if create.DefaultQoS != "" {
		apiAccount.DefaultQos = &create.DefaultQoS
	}
	if len(create.QoSList) > 0 {
		apiAccount.QosList = &create.QoSList
	}

	// Parent account
	if create.ParentName != "" {
		apiAccount.Parent = &create.ParentName
	}

	// Shares and priority
	if create.SharesRaw > 0 {
		apiAccount.SharesRaw = &create.SharesRaw
	}
	if create.Priority > 0 {
		setTrue := true
		priority := int32(create.Priority)
		apiAccount.Priority = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &priority,
		}
	}

	// Job limits
	if create.MaxJobs > 0 {
		setTrue := true
		maxJobs := int32(create.MaxJobs)
		apiAccount.MaxJobs = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &maxJobs,
		}
	}
	if create.MaxJobsPerUser > 0 {
		setTrue := true
		maxJobsPerUser := int32(create.MaxJobsPerUser)
		apiAccount.MaxJobsPerUser = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &maxJobsPerUser,
		}
	}
	if create.MaxSubmitJobs > 0 {
		setTrue := true
		maxSubmitJobs := int32(create.MaxSubmitJobs)
		apiAccount.MaxSubmitJobs = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &maxSubmitJobs,
		}
	}

	// Time limits
	if create.MaxWallTime > 0 {
		setTrue := true
		maxWallTime := int32(create.MaxWallTime)
		apiAccount.MaxWallTime = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &maxWallTime,
		}
	}

	// Resource limits
	if create.MaxNodes > 0 {
		setTrue := true
		maxNodes := int32(create.MaxNodes)
		apiAccount.MaxNodes = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &maxNodes,
		}
	}
	if create.MaxCPUs > 0 {
		setTrue := true
		maxCPUs := int32(create.MaxCPUs)
		apiAccount.MaxCpus = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &maxCPUs,
		}
	}

	// Group limits
	if create.GrpJobs > 0 {
		setTrue := true
		grpJobs := int32(create.GrpJobs)
		apiAccount.GrpJobs = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &grpJobs,
		}
	}
	if create.GrpNodes > 0 {
		setTrue := true
		grpNodes := int32(create.GrpNodes)
		apiAccount.GrpNodes = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &grpNodes,
		}
	}
	if create.GrpCPUs > 0 {
		setTrue := true
		grpCPUs := int32(create.GrpCPUs)
		apiAccount.GrpCpus = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &grpCPUs,
		}
	}
	if create.GrpSubmitJobs > 0 {
		setTrue := true
		grpSubmitJobs := int32(create.GrpSubmitJobs)
		apiAccount.GrpSubmitJobs = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &grpSubmitJobs,
		}
	}
	if create.GrpWallTime > 0 {
		setTrue := true
		grpWallTime := int32(create.GrpWallTime)
		apiAccount.GrpWallTime = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &grpWallTime,
		}
	}

	return apiAccount, nil
}

// convertCommonAccountUpdateToAPI converts common AccountUpdate to v0.0.43 API format
func (a *AccountAdapter) convertCommonAccountUpdateToAPI(existing *types.Account, update *types.AccountUpdate) (*api.V0043AccountInfo, error) {
	apiAccount := &api.V0043AccountInfo{}

	// Always include the account name for updates
	apiAccount.Name = &existing.Name

	// Apply updates to fields
	description := existing.Description
	if update.Description != nil {
		description = *update.Description
	}
	if description != "" {
		apiAccount.Description = &description
	}

	organization := existing.Organization
	if update.Organization != nil {
		organization = *update.Organization
	}
	if organization != "" {
		apiAccount.Organization = &organization
	}

	// Coordinators
	coordinators := existing.Coordinators
	if len(update.Coordinators) > 0 {
		coordinators = update.Coordinators
	}
	if len(coordinators) > 0 {
		apiCoordinators := make([]api.V0043AccountCoordinator, len(coordinators))
		for i, coord := range coordinators {
			apiCoordinators[i] = api.V0043AccountCoordinator{
				Name: &coord,
			}
		}
		apiAccount.Coordinators = &apiCoordinators
	}

	// QoS settings
	defaultQoS := existing.DefaultQoS
	if update.DefaultQoS != nil {
		defaultQoS = *update.DefaultQoS
	}
	if defaultQoS != "" {
		apiAccount.DefaultQos = &defaultQoS
	}

	qosList := existing.QoSList
	if len(update.QoSList) > 0 {
		qosList = update.QoSList
	}
	if len(qosList) > 0 {
		apiAccount.QosList = &qosList
	}

	// Priority
	priority := existing.Priority
	if update.Priority != nil {
		priority = *update.Priority
	}
	if priority > 0 {
		setTrue := true
		priorityInt32 := int32(priority)
		apiAccount.Priority = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &priorityInt32,
		}
	}

	// Job limits
	maxJobs := existing.MaxJobs
	if update.MaxJobs != nil {
		maxJobs = *update.MaxJobs
	}
	if maxJobs > 0 {
		setTrue := true
		maxJobsInt32 := int32(maxJobs)
		apiAccount.MaxJobs = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &maxJobsInt32,
		}
	}

	maxWallTime := existing.MaxWallTime
	if update.MaxWallTime != nil {
		maxWallTime = *update.MaxWallTime
	}
	if maxWallTime > 0 {
		setTrue := true
		maxWallTimeInt32 := int32(maxWallTime)
		apiAccount.MaxWallTime = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &maxWallTimeInt32,
		}
	}

	return apiAccount, nil
}