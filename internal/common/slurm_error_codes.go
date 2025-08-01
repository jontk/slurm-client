// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package common

// SLURM-specific error codes and their meanings
// Based on SLURM documentation and observed API responses

// SlurmErrorCode represents a SLURM-specific error code
type SlurmErrorCode int32

const (
	// Success codes
	SlurmSuccess SlurmErrorCode = 0

	// Job submission errors (2000-2099)
	SlurmErrorInvalidPartition     SlurmErrorCode = 2001
	SlurmErrorInvalidAccount       SlurmErrorCode = 2002
	SlurmErrorInvalidTimeLimit     SlurmErrorCode = 2003
	SlurmErrorInvalidNodeCount     SlurmErrorCode = 2004
	SlurmErrorInvalidMemory        SlurmErrorCode = 2005
	SlurmErrorInvalidCPU           SlurmErrorCode = 2006
	SlurmErrorInvalidGres          SlurmErrorCode = 2007
	SlurmErrorJobPending           SlurmErrorCode = 2008
	SlurmErrorJobQueueFull         SlurmErrorCode = 2009
	SlurmErrorInvalidConstraints   SlurmErrorCode = 2010
	SlurmErrorBatchJobSubmitFailed SlurmErrorCode = 2063 // Job cannot be submitted without current working directory
	SlurmErrorInvalidJobID         SlurmErrorCode = 2064
	SlurmErrorJobAlreadyCompleted  SlurmErrorCode = 2065
	SlurmErrorJobNotFound          SlurmErrorCode = 2066
	SlurmErrorJobCancelled         SlurmErrorCode = 2067
	SlurmErrorJobTimeout           SlurmErrorCode = 2068
	SlurmErrorJobMemoryLimit       SlurmErrorCode = 2069
	SlurmErrorJobNodeFailure       SlurmErrorCode = 2070
	
	// Resource errors (3000-3099)
	SlurmErrorNodeNotAvailable     SlurmErrorCode = 3001
	SlurmErrorPartitionDown        SlurmErrorCode = 3002
	SlurmErrorInsufficientResource SlurmErrorCode = 3003
	SlurmErrorNodeDown             SlurmErrorCode = 3004
	SlurmErrorNodeDrained          SlurmErrorCode = 3005
	SlurmErrorNodeMaintenance      SlurmErrorCode = 3006
	SlurmErrorPartitionInactive    SlurmErrorCode = 3007
	SlurmErrorResourceExhausted    SlurmErrorCode = 3008
	SlurmErrorLicenseUnavailable   SlurmErrorCode = 3009
	SlurmErrorGresUnavailable      SlurmErrorCode = 3010
	
	// Account/User errors (4000-4099)
	SlurmErrorAccountNotFound      SlurmErrorCode = 4001
	SlurmErrorUserNotFound         SlurmErrorCode = 4002
	SlurmErrorInvalidAssociation   SlurmErrorCode = 4003
	SlurmErrorAccountAlreadyExists SlurmErrorCode = 4004
	SlurmErrorUserAlreadyExists    SlurmErrorCode = 4005
	SlurmErrorAssociationNotFound  SlurmErrorCode = 4006
	SlurmErrorAssociationExists    SlurmErrorCode = 4007
	SlurmErrorInvalidUser          SlurmErrorCode = 4008
	SlurmErrorUserPermissionDenied SlurmErrorCode = 4009
	SlurmErrorAccountLocked        SlurmErrorCode = 4010
	
	// QoS errors (5000-5099)
	SlurmErrorQoSNotFound          SlurmErrorCode = 5001
	SlurmErrorQoSAlreadyExists     SlurmErrorCode = 5002
	SlurmErrorInvalidQoS           SlurmErrorCode = 5003
	SlurmErrorQoSLimitExceeded     SlurmErrorCode = 5004
	SlurmErrorQoSPriorityInvalid   SlurmErrorCode = 5005
	SlurmErrorQoSResourceLimit     SlurmErrorCode = 5006
	SlurmErrorQoSTimeLimit         SlurmErrorCode = 5007
	SlurmErrorQoSUsageLimitHit     SlurmErrorCode = 5008
	
	// Reservation errors (6000-6099)
	SlurmErrorReservationNotFound  SlurmErrorCode = 6001
	SlurmErrorReservationInvalid   SlurmErrorCode = 6002
	SlurmErrorReservationBusy      SlurmErrorCode = 6003
	SlurmErrorReservationExists    SlurmErrorCode = 6004
	SlurmErrorReservationExpired   SlurmErrorCode = 6005
	SlurmErrorReservationDenied    SlurmErrorCode = 6006
	SlurmErrorReservationMaintenance SlurmErrorCode = 6007
	SlurmErrorReservationNodeCount   SlurmErrorCode = 6008
	SlurmErrorReservationTimeConflict SlurmErrorCode = 6009
	
	// Authentication errors (7000-7099)
	SlurmErrorAuthenticationFailed SlurmErrorCode = 7001
	SlurmErrorPermissionDenied     SlurmErrorCode = 7002
	SlurmErrorTokenExpired         SlurmErrorCode = 7003
	SlurmErrorTokenInvalid         SlurmErrorCode = 7004
	SlurmErrorTokenMissing         SlurmErrorCode = 7005
	SlurmErrorCredentialsInvalid   SlurmErrorCode = 7006
	SlurmErrorAccessDenied         SlurmErrorCode = 7007
	SlurmErrorRoleInvalid          SlurmErrorCode = 7008
	
	// Configuration errors (8000-8099)
	SlurmErrorConfigInvalid        SlurmErrorCode = 8001
	SlurmErrorConfigMissing        SlurmErrorCode = 8002
	SlurmErrorConfigSyntax         SlurmErrorCode = 8003
	SlurmErrorConfigPermission     SlurmErrorCode = 8004
	SlurmErrorConfigDuplicate      SlurmErrorCode = 8005
	SlurmErrorPluginError          SlurmErrorCode = 8006
	SlurmErrorLibraryMissing       SlurmErrorCode = 8007
	
	// Communication errors (8100-8199)
	SlurmErrorConnectionRefused    SlurmErrorCode = 8101
	SlurmErrorConnectionTimeout    SlurmErrorCode = 8102
	SlurmErrorConnectionLost       SlurmErrorCode = 8103
	SlurmErrorProtocolVersion      SlurmErrorCode = 8104
	SlurmErrorMessageFormat        SlurmErrorCode = 8105
	SlurmErrorNetworkUnavailable   SlurmErrorCode = 8106
	
	// General errors (9000-9099)
	SlurmErrorUnknown              SlurmErrorCode = 9000
	SlurmErrorInvalidRequest       SlurmErrorCode = 9001
	SlurmErrorDatabaseError        SlurmErrorCode = 9002
	SlurmErrorConfigError          SlurmErrorCode = 9003
	SlurmErrorSystemError          SlurmErrorCode = 9004
	SlurmErrorOutOfMemory          SlurmErrorCode = 9005
	SlurmErrorFileNotFound         SlurmErrorCode = 9006
	SlurmErrorFilePermission       SlurmErrorCode = 9007
	SlurmErrorDiskFull             SlurmErrorCode = 9008
	SlurmErrorInternalError        SlurmErrorCode = 9009
)

// SlurmErrorInfo provides detailed information about a SLURM error code
type SlurmErrorInfo struct {
	Code        SlurmErrorCode
	Name        string
	Description string
	Category    string
}

// slurmErrorMap maps error codes to their detailed information
var slurmErrorMap = map[SlurmErrorCode]SlurmErrorInfo{
	SlurmSuccess: {
		Code:        SlurmSuccess,
		Name:        "SUCCESS",
		Description: "Operation completed successfully",
		Category:    "Success",
	},
	
	// Job submission errors
	SlurmErrorInvalidPartition: {
		Code:        SlurmErrorInvalidPartition,
		Name:        "INVALID_PARTITION",
		Description: "The specified partition does not exist or is not available",
		Category:    "Job Submission",
	},
	SlurmErrorInvalidAccount: {
		Code:        SlurmErrorInvalidAccount,
		Name:        "INVALID_ACCOUNT",
		Description: "The specified account does not exist or user has no access",
		Category:    "Job Submission",
	},
	SlurmErrorInvalidTimeLimit: {
		Code:        SlurmErrorInvalidTimeLimit,
		Name:        "INVALID_TIME_LIMIT",
		Description: "The specified time limit is invalid or exceeds partition limits",
		Category:    "Job Submission",
	},
	SlurmErrorInvalidNodeCount: {
		Code:        SlurmErrorInvalidNodeCount,
		Name:        "INVALID_NODE_COUNT",
		Description: "The requested number of nodes is invalid or unavailable",
		Category:    "Job Submission",
	},
	SlurmErrorInvalidMemory: {
		Code:        SlurmErrorInvalidMemory,
		Name:        "INVALID_MEMORY",
		Description: "The requested memory amount is invalid or exceeds limits",
		Category:    "Job Submission",
	},
	SlurmErrorInvalidCPU: {
		Code:        SlurmErrorInvalidCPU,
		Name:        "INVALID_CPU",
		Description: "The requested CPU count is invalid or exceeds limits",
		Category:    "Job Submission",
	},
	SlurmErrorBatchJobSubmitFailed: {
		Code:        SlurmErrorBatchJobSubmitFailed,
		Name:        "BATCH_JOB_SUBMIT_FAILED",
		Description: "Batch job submission failed - often due to missing required fields like working directory",
		Category:    "Job Submission",
	},
	SlurmErrorInvalidJobID: {
		Code:        SlurmErrorInvalidJobID,
		Name:        "INVALID_JOB_ID",
		Description: "The specified job ID is invalid or does not exist",
		Category:    "Job Management",
	},
	SlurmErrorJobAlreadyCompleted: {
		Code:        SlurmErrorJobAlreadyCompleted,
		Name:        "JOB_ALREADY_COMPLETED",
		Description: "The job has already completed and cannot be modified",
		Category:    "Job Management",
	},
	SlurmErrorJobNotFound: {
		Code:        SlurmErrorJobNotFound,
		Name:        "JOB_NOT_FOUND",
		Description: "The requested job was not found",
		Category:    "Job Management",
	},
	
	// Resource errors
	SlurmErrorNodeNotAvailable: {
		Code:        SlurmErrorNodeNotAvailable,
		Name:        "NODE_NOT_AVAILABLE",
		Description: "The requested nodes are not available",
		Category:    "Resource Management",
	},
	SlurmErrorPartitionDown: {
		Code:        SlurmErrorPartitionDown,
		Name:        "PARTITION_DOWN",
		Description: "The partition is currently down or unavailable",
		Category:    "Resource Management",
	},
	SlurmErrorInsufficientResource: {
		Code:        SlurmErrorInsufficientResource,
		Name:        "INSUFFICIENT_RESOURCE",
		Description: "Insufficient resources available to fulfill the request",
		Category:    "Resource Management",
	},
	SlurmErrorResourceExhausted: {
		Code:        SlurmErrorResourceExhausted,
		Name:        "RESOURCE_EXHAUSTED",
		Description: "The requested resources have been exhausted",
		Category:    "Resource Management",
	},
	
	// Account/User errors
	SlurmErrorAccountNotFound: {
		Code:        SlurmErrorAccountNotFound,
		Name:        "ACCOUNT_NOT_FOUND",
		Description: "The requested account does not exist",
		Category:    "Account Management",
	},
	SlurmErrorAccountAlreadyExists: {
		Code:        SlurmErrorAccountAlreadyExists,
		Name:        "ACCOUNT_ALREADY_EXISTS",
		Description: "An account with this name already exists",
		Category:    "Account Management",
	},
	SlurmErrorUserNotFound: {
		Code:        SlurmErrorUserNotFound,
		Name:        "USER_NOT_FOUND",
		Description: "The requested user does not exist",
		Category:    "User Management",
	},
	SlurmErrorUserAlreadyExists: {
		Code:        SlurmErrorUserAlreadyExists,
		Name:        "USER_ALREADY_EXISTS",
		Description: "A user with this name already exists",
		Category:    "User Management",
	},
	SlurmErrorInvalidAssociation: {
		Code:        SlurmErrorInvalidAssociation,
		Name:        "INVALID_ASSOCIATION",
		Description: "The specified association is invalid",
		Category:    "Association Management",
	},
	SlurmErrorAssociationNotFound: {
		Code:        SlurmErrorAssociationNotFound,
		Name:        "ASSOCIATION_NOT_FOUND",
		Description: "The requested association does not exist",
		Category:    "Association Management",
	},
	
	// QoS errors
	SlurmErrorQoSNotFound: {
		Code:        SlurmErrorQoSNotFound,
		Name:        "QOS_NOT_FOUND",
		Description: "The requested QoS does not exist",
		Category:    "QoS Management",
	},
	SlurmErrorQoSAlreadyExists: {
		Code:        SlurmErrorQoSAlreadyExists,
		Name:        "QOS_ALREADY_EXISTS",
		Description: "A QoS with this name already exists",
		Category:    "QoS Management",
	},
	SlurmErrorInvalidQoS: {
		Code:        SlurmErrorInvalidQoS,
		Name:        "INVALID_QOS",
		Description: "The specified QoS is invalid or not available to the user",
		Category:    "QoS Management",
	},
	SlurmErrorQoSLimitExceeded: {
		Code:        SlurmErrorQoSLimitExceeded,
		Name:        "QOS_LIMIT_EXCEEDED",
		Description: "The request exceeds QoS limits",
		Category:    "QoS Management",
	},
	
	// Reservation errors
	SlurmErrorReservationNotFound: {
		Code:        SlurmErrorReservationNotFound,
		Name:        "RESERVATION_NOT_FOUND",
		Description: "The requested reservation does not exist",
		Category:    "Reservation Management",
	},
	SlurmErrorReservationInvalid: {
		Code:        SlurmErrorReservationInvalid,
		Name:        "RESERVATION_INVALID",
		Description: "The reservation request is invalid",
		Category:    "Reservation Management",
	},
	SlurmErrorReservationBusy: {
		Code:        SlurmErrorReservationBusy,
		Name:        "RESERVATION_BUSY",
		Description: "The reservation is currently in use and cannot be modified",
		Category:    "Reservation Management",
	},
	SlurmErrorReservationExists: {
		Code:        SlurmErrorReservationExists,
		Name:        "RESERVATION_EXISTS",
		Description: "A reservation with this name already exists",
		Category:    "Reservation Management",
	},
	
	// Authentication errors
	SlurmErrorAuthenticationFailed: {
		Code:        SlurmErrorAuthenticationFailed,
		Name:        "AUTHENTICATION_FAILED",
		Description: "Authentication failed - check credentials or token",
		Category:    "Authentication",
	},
	SlurmErrorPermissionDenied: {
		Code:        SlurmErrorPermissionDenied,
		Name:        "PERMISSION_DENIED",
		Description: "User does not have permission for this operation",
		Category:    "Authentication",
	},
	SlurmErrorTokenExpired: {
		Code:        SlurmErrorTokenExpired,
		Name:        "TOKEN_EXPIRED",
		Description: "The authentication token has expired",
		Category:    "Authentication",
	},
	SlurmErrorTokenInvalid: {
		Code:        SlurmErrorTokenInvalid,
		Name:        "TOKEN_INVALID",
		Description: "The authentication token is invalid",
		Category:    "Authentication",
	},
	
	// Communication errors
	SlurmErrorConnectionRefused: {
		Code:        SlurmErrorConnectionRefused,
		Name:        "CONNECTION_REFUSED",
		Description: "Connection to SLURM daemon was refused",
		Category:    "Communication",
	},
	SlurmErrorConnectionTimeout: {
		Code:        SlurmErrorConnectionTimeout,
		Name:        "CONNECTION_TIMEOUT",
		Description: "Connection to SLURM daemon timed out",
		Category:    "Communication",
	},
	SlurmErrorProtocolVersion: {
		Code:        SlurmErrorProtocolVersion,
		Name:        "PROTOCOL_VERSION",
		Description: "SLURM protocol version mismatch",
		Category:    "Communication",
	},
	
	// General errors
	SlurmErrorInvalidRequest: {
		Code:        SlurmErrorInvalidRequest,
		Name:        "INVALID_REQUEST",
		Description: "The request is malformed or invalid",
		Category:    "General",
	},
	SlurmErrorDatabaseError: {
		Code:        SlurmErrorDatabaseError,
		Name:        "DATABASE_ERROR",
		Description: "Database operation failed",
		Category:    "General",
	},
	SlurmErrorSystemError: {
		Code:        SlurmErrorSystemError,
		Name:        "SYSTEM_ERROR",
		Description: "System error occurred",
		Category:    "General",
	},
	SlurmErrorInternalError: {
		Code:        SlurmErrorInternalError,
		Name:        "INTERNAL_ERROR",
		Description: "Internal SLURM error occurred",
		Category:    "General",
	},
}

// GetErrorInfo returns detailed information about a SLURM error code
func GetErrorInfo(code int32) *SlurmErrorInfo {
	if info, exists := slurmErrorMap[SlurmErrorCode(code)]; exists {
		return &info
	}
	return &SlurmErrorInfo{
		Code:        SlurmErrorCode(code),
		Name:        "UNKNOWN_ERROR",
		Description: "Unknown SLURM error code",
		Category:    "Unknown",
	}
}

// IsKnownError checks if the error code is a known SLURM error
func IsKnownError(code int32) bool {
	_, exists := slurmErrorMap[SlurmErrorCode(code)]
	return exists
}

// GetErrorCategory returns the category of the error
func GetErrorCategory(code int32) string {
	if info := GetErrorInfo(code); info != nil {
		return info.Category
	}
	return "Unknown"
}

// GetErrorDescription returns a human-readable description of the error
func GetErrorDescription(code int32) string {
	if info := GetErrorInfo(code); info != nil {
		return info.Description
	}
	return "Unknown error"
}

// EnhanceErrorMessage enhances an error message with SLURM-specific information
func EnhanceErrorMessage(errorCode int32, originalMessage string) string {
	info := GetErrorInfo(errorCode)
	if info.Name != "UNKNOWN_ERROR" {
		return info.Description + " (SLURM Error " + info.Name + ")"
	}
	return originalMessage
}
