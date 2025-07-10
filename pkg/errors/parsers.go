package errors

import (
	"encoding/json"
	"fmt"
	"strings"
)

// SlurmAPIResponse represents the structure of Slurm REST API error responses
type SlurmAPIResponse struct {
	Meta   *SlurmAPIMeta         `json:"meta,omitempty"`
	Errors []SlurmAPIErrorDetail `json:"errors,omitempty"`
	Data   interface{}           `json:"data,omitempty"`
}

// SlurmAPIMeta contains metadata about the API response
type SlurmAPIMeta struct {
	Plugin       *SlurmPlugin  `json:"plugin,omitempty"`
	SlurmVersion *SlurmVersion `json:"Slurm,omitempty"`
}

// SlurmPlugin contains information about the Slurm plugin
type SlurmPlugin struct {
	Type       string `json:"type,omitempty"`
	Name       string `json:"name,omitempty"`
	DataParser string `json:"data_parser,omitempty"`
}

// SlurmVersion contains Slurm version information
type SlurmVersion struct {
	Version struct {
		Major int `json:"major"`
		Micro int `json:"micro"`
		Minor int `json:"minor"`
	} `json:"version"`
	Release string `json:"release"`
}

// parseSlurmAPIError attempts to parse a Slurm API error response
func parseSlurmAPIError(statusCode int, body []byte, apiVersion string) *SlurmAPIError {
	if len(body) == 0 {
		return nil
	}

	// Try to parse as JSON
	var response SlurmAPIResponse
	if err := json.Unmarshal(body, &response); err != nil {
		// If JSON parsing fails, try to extract error information from plain text
		return parsePlainTextError(statusCode, body, apiVersion)
	}

	// Extract error details
	var errors []SlurmAPIErrorDetail
	if len(response.Errors) > 0 {
		errors = response.Errors
	} else {
		// No structured errors, create one from status code
		errors = []SlurmAPIErrorDetail{
			{
				ErrorNumber: statusCode,
				ErrorCode:   fmt.Sprintf("HTTP_%d", statusCode),
				Source:      "http",
				Description: fmt.Sprintf("HTTP %d error", statusCode),
			},
		}
	}

	return NewSlurmAPIError(statusCode, apiVersion, errors)
}

// parsePlainTextError handles non-JSON error responses
func parsePlainTextError(statusCode int, body []byte, apiVersion string) *SlurmAPIError {
	bodyStr := string(body)

	// Only process if it looks like a Slurm error message
	hasKnownSlurmError := strings.Contains(bodyStr, "SLURM_")

	if !hasKnownSlurmError {
		// Return nil so that the caller can handle as generic HTTP error
		return nil
	}

	// Look for common Slurm error patterns in plain text
	var errorCode string
	var description string

	// Try to extract Slurm error codes from text
	switch {
	case strings.Contains(bodyStr, "SLURM_NO_CHANGE_IN_DATA"):
		errorCode = "SLURM_NO_CHANGE_IN_DATA"
		description = "No data changes detected"
	case strings.Contains(bodyStr, "SLURM_PROTOCOL_VERSION_ERROR"):
		errorCode = "SLURM_PROTOCOL_VERSION_ERROR"
		description = "Protocol version mismatch"
	case strings.Contains(bodyStr, "SLURM_AUTHENTICATION_ERROR"):
		errorCode = "SLURM_AUTHENTICATION_ERROR"
		description = "Authentication failed"
	case strings.Contains(bodyStr, "SLURM_ACCESS_DENIED"):
		errorCode = "SLURM_ACCESS_DENIED"
		description = "Access denied"
	case strings.Contains(bodyStr, "SLURM_INVALID_JOB_ID"):
		errorCode = "SLURM_INVALID_JOB_ID"
		description = "Invalid job ID"
	case strings.Contains(bodyStr, "SLURM_INVALID_PARTITION_NAME"):
		errorCode = "SLURM_INVALID_PARTITION_NAME"
		description = "Invalid partition name"
	case strings.Contains(bodyStr, "SLURM_NODE_NOT_AVAIL"):
		errorCode = "SLURM_NODE_NOT_AVAIL"
		description = "Node not available"
	case strings.Contains(bodyStr, "SLURM_JOB_PENDING"):
		errorCode = "SLURM_JOB_PENDING"
		description = "Job is pending"
	case strings.Contains(bodyStr, "SLURM_JOB_ALREADY_COMPLETE"):
		errorCode = "SLURM_JOB_ALREADY_COMPLETE"
		description = "Job already completed"
	default:
		errorCode = fmt.Sprintf("HTTP_%d", statusCode)
		description = bodyStr
		if len(description) > 200 {
			description = description[:200] + "..."
		}
	}

	errors := []SlurmAPIErrorDetail{
		{
			ErrorNumber: statusCode,
			ErrorCode:   errorCode,
			Source:      "text_response",
			Description: description,
		},
	}

	return NewSlurmAPIError(statusCode, apiVersion, errors)
}

// ExtractRequestID attempts to extract a request ID from various sources
func ExtractRequestID(headers map[string][]string, body []byte) string {
	// Check common request ID headers
	requestIDHeaders := []string{
		"X-Request-ID",
		"X-Request-Id",
		"Request-ID",
		"Request-Id",
		"X-Correlation-ID",
		"X-Correlation-Id",
		"X-Trace-ID",
		"X-Trace-Id",
	}

	for _, headerName := range requestIDHeaders {
		if values, exists := headers[headerName]; exists && len(values) > 0 {
			return values[0]
		}
	}

	// Try to extract from response body if it's JSON
	if len(body) > 0 {
		var response map[string]interface{}
		if err := json.Unmarshal(body, &response); err == nil {
			if requestID, exists := response["request_id"]; exists {
				if id, ok := requestID.(string); ok {
					return id
				}
			}
			if meta, exists := response["meta"]; exists {
				if metaMap, ok := meta.(map[string]interface{}); ok {
					if requestID, exists := metaMap["request_id"]; exists {
						if id, ok := requestID.(string); ok {
							return id
						}
					}
				}
			}
		}
	}

	return ""
}

// ParseVersionFromResponse extracts API version from response metadata
func ParseVersionFromResponse(body []byte) string {
	if len(body) == 0 {
		return ""
	}

	var response SlurmAPIResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return ""
	}

	if response.Meta != nil && response.Meta.SlurmVersion != nil {
		return response.Meta.SlurmVersion.Release
	}

	return ""
}

// ErrorContainsPattern checks if error message contains specific patterns
func ErrorContainsPattern(err error, patterns ...string) bool {
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())
	for _, pattern := range patterns {
		if strings.Contains(errStr, strings.ToLower(pattern)) {
			return true
		}
	}
	return false
}

// ExtractJobIDFromError attempts to extract job ID from error messages
func ExtractJobIDFromError(err error) (uint32, bool) {
	if err == nil {
		return 0, false
	}

	errStr := err.Error()

	// Look for patterns like "job 12345" or "job_id: 12345"
	patterns := []string{
		"job ",
		"job_id:",
		"job id:",
		"jobid:",
		"job-id:",
	}

	for _, pattern := range patterns {
		if idx := strings.Index(strings.ToLower(errStr), pattern); idx != -1 {
			start := idx + len(pattern)
			end := start

			// Find the end of the number
			for end < len(errStr) && errStr[end] >= '0' && errStr[end] <= '9' {
				end++
			}

			if end > start {
				var jobID uint32
				if n, err := fmt.Sscanf(errStr[start:end], "%d", &jobID); n == 1 && err == nil {
					return jobID, true
				}
			}
		}
	}

	return 0, false
}

// ExtractNodeNamesFromError attempts to extract node names from error messages
func ExtractNodeNamesFromError(err error) ([]string, bool) {
	if err == nil {
		return nil, false
	}

	errStr := err.Error()

	// Look for patterns like "node compute-01" or "nodes: compute-[01-03]"
	patterns := []string{
		"node ",
		"nodes:",
		"node:",
		"nodelist:",
		"node list:",
	}

	for _, pattern := range patterns {
		if idx := strings.Index(strings.ToLower(errStr), pattern); idx != -1 {
			start := idx + len(pattern)

			// Find the end of the node specification
			end := start
			for end < len(errStr) && errStr[end] != ' ' && errStr[end] != ',' && errStr[end] != '\n' {
				end++
			}

			if end > start {
				nodeSpec := strings.TrimSpace(errStr[start:end])
				if nodeSpec != "" {
					// For now, return as single node; could expand to parse node ranges
					return []string{nodeSpec}, true
				}
			}
		}
	}

	return nil, false
}

// ExtractPartitionFromError attempts to extract partition name from error messages
func ExtractPartitionFromError(err error) (string, bool) {
	if err == nil {
		return "", false
	}

	errStr := err.Error()

	// Look for patterns like "partition debug" or "partition: compute"
	patterns := []string{
		"partition ",
		"partition:",
	}

	for _, pattern := range patterns {
		if idx := strings.Index(strings.ToLower(errStr), pattern); idx != -1 {
			start := idx + len(pattern)

			// Find the end of the partition name
			end := start
			for end < len(errStr) && errStr[end] != ' ' && errStr[end] != ',' && errStr[end] != '\n' {
				end++
			}

			if end > start {
				partition := strings.TrimSpace(errStr[start:end])
				if partition != "" {
					return partition, true
				}
			}
		}
	}

	return "", false
}
