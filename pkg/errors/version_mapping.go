package errors

import "fmt"

// VersionSpecificErrorMapping provides error mapping for different API versions
type VersionSpecificErrorMapping struct {
	Version            string
	SlurmErrorMappings map[string]ErrorCode
	HTTPStatusMappings map[int]ErrorCode
	FeatureSupport     map[string]bool
}

// GetVersionMapping returns error mapping for a specific API version
func GetVersionMapping(apiVersion string) *VersionSpecificErrorMapping {
	switch apiVersion {
	case "v0.0.40":
		return getV0040Mapping()
	case "v0.0.41":
		return getV0041Mapping()
	case "v0.0.42":
		return getV0042Mapping()
	case "v0.0.43":
		return getV0043Mapping()
	default:
		// Fall back to latest version mapping
		return getV0043Mapping()
	}
}

// getV0040Mapping returns error mapping for API v0.0.40 (Slurm 24.05-25.05)
func getV0040Mapping() *VersionSpecificErrorMapping {
	return &VersionSpecificErrorMapping{
		Version: "v0.0.40",
		SlurmErrorMappings: map[string]ErrorCode{
			// Core Slurm errors
			"SLURM_SUCCESS":                        ErrorCodeUnknown, // Special case - not an error
			"SLURM_ERROR":                          ErrorCodeServerInternal,
			"SLURM_PROTOCOL_VERSION_ERROR":         ErrorCodeVersionMismatch,
			"SLURM_UNEXPECTED_MSG_ERROR":           ErrorCodeInvalidRequest,
			
			// Authentication errors
			"SLURM_AUTHENTICATION_ERROR":           ErrorCodeInvalidCredentials,
			"SLURM_ACCESS_DENIED":                  ErrorCodePermissionDenied,
			
			// Communication errors
			"SLURM_COMMUNICATIONS_CONNECTION_ERROR": ErrorCodeConnectionRefused,
			"SLURM_COMMUNICATIONS_SEND_ERROR":      ErrorCodeNetworkTimeout,
			"SLURM_COMMUNICATIONS_RECEIVE_ERROR":   ErrorCodeNetworkTimeout,
			"SLURM_COMMUNICATIONS_SHUTDOWN_ERROR":  ErrorCodeSlurmDaemonDown,
			
			// Job-related errors
			"SLURM_INVALID_JOB_ID":                 ErrorCodeResourceNotFound,
			"SLURM_JOB_PENDING":                    ErrorCodeResourceExhausted,
			"SLURM_JOB_ALREADY_COMPLETE":           ErrorCodeConflict,
			"SLURM_DUPLICATE_JOB_ID":               ErrorCodeConflict,
			"SLURM_JOB_FINISHED":                   ErrorCodeConflict,
			"SLURM_JOB_SUSPENDED":                  ErrorCodeConflict,
			
			// Node-related errors
			"SLURM_INVALID_NODE_NAME":              ErrorCodeResourceNotFound,
			"SLURM_NODE_NOT_AVAIL":                 ErrorCodeResourceExhausted,
			"SLURM_REQUESTED_NODE_CONFIG_UNAVAILABLE": ErrorCodeResourceExhausted,
			
			// Partition-related errors
			"SLURM_INVALID_PARTITION_NAME":         ErrorCodeResourceNotFound,
			"SLURM_PARTITION_DOWN":                 ErrorCodePartitionUnavailable,
			"SLURM_PARTITION_NOT_AVAIL":            ErrorCodePartitionUnavailable,
			
			// Resource-related errors
			"SLURM_REQUESTED_PART_CONFIG_UNAVAILABLE": ErrorCodeResourceExhausted,
			"SLURM_NODES_BUSY":                     ErrorCodeResourceExhausted,
			"SLURM_INVALID_TIME_LIMIT":             ErrorCodeValidationFailed,
			"SLURM_INVALID_TASK_COUNT":             ErrorCodeValidationFailed,
			
			// Data/configuration errors
			"SLURM_NO_CHANGE_IN_DATA":              ErrorCodeResourceNotFound,
			"SLURM_DATA_CONV_FAILED":               ErrorCodeServerInternal,
			"SLURM_NOT_SUPPORTED":                  ErrorCodeUnsupportedOperation,
			
			// v0.0.40 specific error codes (if any unique ones exist)
			"SLURM_EXCLUSIVE_ACCESS_ERROR":         ErrorCodeConflict,
			"SLURM_OVERSUBSCRIBE_ERROR":            ErrorCodeValidationFailed,
		},
		HTTPStatusMappings: getCommonHTTPMappings(),
		FeatureSupport: map[string]bool{
			"job_submit":           true,
			"job_cancel":           true,
			"job_list":             true,
			"job_get":              true,
			"node_list":            true,
			"node_get":             true,
			"partition_list":       true,
			"partition_get":        true,
			"info_get":             true,
			"reservation_support":  false, // Limited in v0.0.40
			"job_array_support":    true,
			"frontend_mode":        true,  // Available in v0.0.40
		},
	}
}

// getV0041Mapping returns error mapping for API v0.0.41 (Slurm 24.11-25.11)
func getV0041Mapping() *VersionSpecificErrorMapping {
	mapping := getV0040Mapping()
	mapping.Version = "v0.0.41"
	
	// v0.0.41 specific additions/changes
	mapping.SlurmErrorMappings["SLURM_RESERVATION_INVALID"] = ErrorCodeResourceNotFound
	mapping.SlurmErrorMappings["SLURM_RESERVATION_ACCESS_DENIED"] = ErrorCodePermissionDenied
	mapping.SlurmErrorMappings["SLURM_RESERVATION_NOT_AVAILABLE"] = ErrorCodeResourceExhausted
	
	// Enhanced job features
	mapping.SlurmErrorMappings["SLURM_JOB_SCRIPT_MISSING"] = ErrorCodeValidationFailed
	mapping.SlurmErrorMappings["SLURM_JOB_DEPENDENCY_ERROR"] = ErrorCodeValidationFailed
	
	// Field changes (breaking changes documented)
	mapping.SlurmErrorMappings["SLURM_FIELD_RENAMED_ERROR"] = ErrorCodeValidationFailed
	
	// Updated feature support
	mapping.FeatureSupport["reservation_support"] = true
	mapping.FeatureSupport["enhanced_job_deps"] = true
	mapping.FeatureSupport["field_renames"] = true // minimum_switches → required_switches
	mapping.FeatureSupport["exclusive_field"] = true      // Still available in v0.0.41
	mapping.FeatureSupport["oversubscribe_field"] = true  // Still available in v0.0.41
	
	return mapping
}

// getV0042Mapping returns error mapping for API v0.0.42 (Slurm 25.05+ - Stable)
func getV0042Mapping() *VersionSpecificErrorMapping {
	mapping := getV0041Mapping()
	mapping.Version = "v0.0.42"
	
	// v0.0.42 specific changes (field removals)
	mapping.SlurmErrorMappings["SLURM_FIELD_REMOVED_ERROR"] = ErrorCodeUnsupportedOperation
	mapping.SlurmErrorMappings["SLURM_DEPRECATED_FIELD_ERROR"] = ErrorCodeValidationFailed
	
	// Enhanced validation
	mapping.SlurmErrorMappings["SLURM_STRICT_VALIDATION_ERROR"] = ErrorCodeValidationFailed
	mapping.SlurmErrorMappings["SLURM_SCHEMA_VALIDATION_ERROR"] = ErrorCodeValidationFailed
	
	// Updated feature support (field removals)
	mapping.FeatureSupport["exclusive_field"] = false     // Removed in v0.0.42
	mapping.FeatureSupport["oversubscribe_field"] = false // Removed in v0.0.42
	mapping.FeatureSupport["strict_validation"] = true
	mapping.FeatureSupport["enhanced_schemas"] = true
	
	return mapping
}

// getV0043Mapping returns error mapping for API v0.0.43 (Slurm 25.05+ - Latest)
func getV0043Mapping() *VersionSpecificErrorMapping {
	mapping := getV0042Mapping()
	mapping.Version = "v0.0.43"
	
	// v0.0.43 specific additions
	mapping.SlurmErrorMappings["SLURM_ADVANCED_RESERVATION_ERROR"] = ErrorCodeResourceNotFound
	mapping.SlurmErrorMappings["SLURM_MULTI_CLUSTER_ERROR"] = ErrorCodeUnsupportedOperation
	mapping.SlurmErrorMappings["SLURM_FEDERATION_ERROR"] = ErrorCodeServerInternal
	
	// Enhanced error reporting
	mapping.SlurmErrorMappings["SLURM_DETAILED_ERROR"] = ErrorCodeServerInternal
	mapping.SlurmErrorMappings["SLURM_TRACE_ERROR"] = ErrorCodeServerInternal
	
	// Updated feature support
	mapping.FeatureSupport["advanced_reservations"] = true
	mapping.FeatureSupport["multi_cluster"] = true
	mapping.FeatureSupport["federation_support"] = true
	mapping.FeatureSupport["detailed_errors"] = true
	mapping.FeatureSupport["frontend_mode"] = false // Removed in v0.0.43
	
	return mapping
}

// getCommonHTTPMappings returns HTTP status code mappings common to all versions
func getCommonHTTPMappings() map[int]ErrorCode {
	return map[int]ErrorCode{
		200: ErrorCodeUnknown,        // Success - not an error
		400: ErrorCodeInvalidRequest,
		401: ErrorCodeUnauthorized,
		403: ErrorCodePermissionDenied,
		404: ErrorCodeResourceNotFound,
		409: ErrorCodeConflict,
		422: ErrorCodeValidationFailed,
		429: ErrorCodeRateLimited,
		500: ErrorCodeServerInternal,
		502: ErrorCodeSlurmDaemonDown,
		503: ErrorCodeSlurmDaemonDown,
		504: ErrorCodeNetworkTimeout,
	}
}

// MapSlurmErrorForVersion maps a Slurm error code to client error code for specific version
func MapSlurmErrorForVersion(slurmErrorCode string, apiVersion string, statusCode int) ErrorCode {
	mapping := GetVersionMapping(apiVersion)
	
	// First try version-specific Slurm error mapping
	if clientCode, exists := mapping.SlurmErrorMappings[slurmErrorCode]; exists {
		return clientCode
	}
	
	// Fall back to HTTP status code mapping
	if clientCode, exists := mapping.HTTPStatusMappings[statusCode]; exists {
		return clientCode
	}
	
	// Last resort: map based on HTTP status code using common patterns
	return mapHTTPStatusToErrorCode(statusCode)
}

// IsFeatureSupportedInVersion checks if a feature is supported in the given API version
func IsFeatureSupportedInVersion(feature, apiVersion string) bool {
	mapping := GetVersionMapping(apiVersion)
	if supported, exists := mapping.FeatureSupport[feature]; exists {
		return supported
	}
	return false // Default to not supported for unknown features
}

// GetBreakingChanges returns breaking changes for version transitions
func GetBreakingChanges(fromVersion, toVersion string) []string {
	changes := []string{}
	
	// v0.0.40 → v0.0.41 changes
	if fromVersion == "v0.0.40" && (toVersion == "v0.0.41" || toVersion == "v0.0.42" || toVersion == "v0.0.43") {
		changes = append(changes, "Field renamed: minimum_switches → required_switches")
		changes = append(changes, "Enhanced reservation support added")
	}
	
	// v0.0.41 → v0.0.42 changes
	if (fromVersion == "v0.0.40" || fromVersion == "v0.0.41") && (toVersion == "v0.0.42" || toVersion == "v0.0.43") {
		changes = append(changes, "Removed fields: exclusive, oversubscribe from job outputs")
		changes = append(changes, "Stricter validation enabled")
	}
	
	// v0.0.42 → v0.0.43 changes
	if (fromVersion == "v0.0.40" || fromVersion == "v0.0.41" || fromVersion == "v0.0.42") && toVersion == "v0.0.43" {
		changes = append(changes, "Removed: FrontEnd mode support")
		changes = append(changes, "Added: Advanced reservation management")
		changes = append(changes, "Enhanced: Multi-cluster and federation support")
	}
	
	return changes
}

// ValidateVersionCompatibility checks compatibility between versions
func ValidateVersionCompatibility(clientVersion, serverVersion string) error {
	// Check for explicitly unsupported versions
	supportedVersions := []string{"v0.0.40", "v0.0.41", "v0.0.42", "v0.0.43"}
	clientSupported := false
	serverSupported := false
	
	for _, version := range supportedVersions {
		if version == clientVersion {
			clientSupported = true
		}
		if version == serverVersion {
			serverSupported = true
		}
	}
	
	if !clientSupported || !serverSupported {
		return NewSlurmError(ErrorCodeVersionMismatch, 
			fmt.Sprintf("Unsupported version combination: client=%s, server=%s", clientVersion, serverVersion))
	}
	
	// Check for breaking changes
	breakingChanges := GetBreakingChanges(clientVersion, serverVersion)
	if len(breakingChanges) > 0 {
		details := fmt.Sprintf("Breaking changes detected: %v", breakingChanges)
		err := NewSlurmError(ErrorCodeVersionMismatch, 
			fmt.Sprintf("Version compatibility issue: client=%s, server=%s", clientVersion, serverVersion))
		err.Details = details
		return err
	}
	
	return nil
}

// EnhanceErrorWithVersion adds version-specific context to errors
func EnhanceErrorWithVersion(err *SlurmError, apiVersion string) *SlurmError {
	if err == nil {
		return nil
	}
	
	_ = GetVersionMapping(apiVersion)
	err.APIVersion = apiVersion
	
	// Add version-specific details for certain error types
	switch err.Code {
	case ErrorCodeUnsupportedOperation:
		if !IsFeatureSupportedInVersion("frontend_mode", apiVersion) {
			err.Details = fmt.Sprintf("FrontEnd mode not supported in API %s", apiVersion)
		}
		if !IsFeatureSupportedInVersion("reservation_support", apiVersion) {
			err.Details = fmt.Sprintf("Reservation management not supported in API %s", apiVersion)
		}
		
	case ErrorCodeValidationFailed:
		if apiVersion == "v0.0.42" || apiVersion == "v0.0.43" {
			err.Details = fmt.Sprintf("Strict validation enforced in API %s", apiVersion)
		}
		
	case ErrorCodeVersionMismatch:
		err.Details = fmt.Sprintf("API version %s compatibility issue", apiVersion)
	}
	
	return err
}