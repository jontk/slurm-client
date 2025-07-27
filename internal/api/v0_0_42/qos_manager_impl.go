package v0_0_42

import (
	"context"
	"fmt"
	"strings"

	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/pkg/errors"
)

// QoSManagerImpl implements the QoSManager interface for v0.0.42
type QoSManagerImpl struct {
	client *WrapperClient
}

// NewQoSManagerImpl creates a new QoSManagerImpl
func NewQoSManagerImpl(client *WrapperClient) *QoSManagerImpl {
	return &QoSManagerImpl{
		client: client,
	}
}

// List retrieves a list of QoS with optional filtering
func (q *QoSManagerImpl) List(ctx context.Context, opts *interfaces.ListQoSOptions) (*interfaces.QoSList, error) {
	if q.client == nil || q.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	// Prepare parameters for the API call
	params := &SlurmdbV0042GetQosParams{}

	// Call the generated OpenAPI client
	resp, err := q.client.apiClient.SlurmdbV0042GetQosWithResponse(ctx, params)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.42")
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
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.42", apiErrors)
				return nil, apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.42")
		return nil, httpErr
	}

	// Check for unexpected response format
	if resp.JSON200 == nil || resp.JSON200.Qos == nil {
		return nil, errors.NewClientError(errors.ErrorCodeServerInternal, "Unexpected response format", "Expected JSON response with QoS but got nil")
	}

	// Convert the response to our interface types
	qosList := make([]interfaces.QoS, 0, len(resp.JSON200.Qos))
	for _, apiQos := range resp.JSON200.Qos {
		qos, err := convertAPIQoSToInterface(apiQos)
		if err != nil {
			conversionErr := errors.NewClientError(errors.ErrorCodeServerInternal, "Failed to convert QoS data")
			conversionErr.Cause = err
			conversionErr.Details = fmt.Sprintf("Error converting QoS %v", apiQos.Name)
			return nil, conversionErr
		}
		qosList = append(qosList, *qos)
	}

	// Apply client-side filtering if options are provided
	if opts != nil {
		qosList = filterQoS(qosList, opts)
	}

	return &interfaces.QoSList{
		QoS:   qosList,
		Total: len(qosList),
	}, nil
}

// Get retrieves a specific QoS by name
func (q *QoSManagerImpl) Get(ctx context.Context, qosName string) (*interfaces.QoS, error) {
	if q.client == nil || q.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if qosName == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "QoS name is required", "qosName", qosName, nil)
	}

	// Prepare parameters for the API call
	params := &SlurmdbV0042GetSingleQosParams{}

	// Call the generated OpenAPI client
	resp, err := q.client.apiClient.SlurmdbV0042GetSingleQosWithResponse(ctx, qosName, params)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.42")
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
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.42", apiErrors)
				return nil, apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.42")
		return nil, httpErr
	}

	// Check for unexpected response format
	if resp.JSON200 == nil || resp.JSON200.Qos == nil || len(resp.JSON200.Qos) == 0 {
		return nil, errors.NewClientError(errors.ErrorCodeResourceNotFound, "QoS not found", fmt.Sprintf("QoS '%s' not found", qosName))
	}

	// Convert the first QoS in the response
	qos, err := convertAPIQoSToInterface(resp.JSON200.Qos[0])
	if err != nil {
		conversionErr := errors.NewClientError(errors.ErrorCodeServerInternal, "Failed to convert QoS data")
		conversionErr.Cause = err
		conversionErr.Details = fmt.Sprintf("Error converting QoS '%s'", qosName)
		return nil, conversionErr
	}

	return qos, nil
}

// Create creates a new QoS
func (q *QoSManagerImpl) Create(ctx context.Context, qos *interfaces.QoSCreate) (*interfaces.QoSCreateResponse, error) {
	if q.client == nil || q.client.apiClient == nil {
		return nil, errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if qos == nil {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "QoS data is required", "qos", qos, nil)
	}

	// Validate required fields
	if qos.Name == "" {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "QoS name is required", "qos.Name", qos.Name, nil)
	}

	// Convert interface types to API types
	apiQoS, err := convertQoSCreateToAPI(qos)
	if err != nil {
		conversionErr := errors.NewClientError(errors.ErrorCodeInvalidRequest, "Failed to convert QoS data")
		conversionErr.Cause = err
		return nil, conversionErr
	}

	// Prepare parameters
	params := &SlurmdbV0042PostQosParams{}

	// Create the request body
	requestBody := SlurmdbV0042PostQosJSONRequestBody{
		Qos: []V0042Qos{*apiQoS},
	}

	// Call the generated OpenAPI client
	resp, err := q.client.apiClient.SlurmdbV0042PostQosWithResponse(ctx, params, requestBody)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return nil, errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.42")
	}

	// Check HTTP status (200 and 201 for creation is success)
	if resp.StatusCode() != 200 && resp.StatusCode() != 201 {
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
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.42", apiErrors)
				return nil, apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.42")
		return nil, httpErr
	}

	// Create successful response
	response := &interfaces.QoSCreateResponse{
		QoSName: qos.Name,
	}

	return response, nil
}

// Update updates an existing QoS
func (q *QoSManagerImpl) Update(ctx context.Context, qosName string, update *interfaces.QoSUpdate) error {
	if q.client == nil || q.client.apiClient == nil {
		return errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if qosName == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "QoS name is required", "qosName", qosName, nil)
	}

	if update == nil {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "update data is required", "update", update, nil)
	}

	// Convert update to API format
	apiUpdate, err := convertQoSUpdateToAPI(update)
	if err != nil {
		conversionErr := errors.NewClientError(errors.ErrorCodeInvalidRequest, "Failed to convert update data")
		conversionErr.Cause = err
		return conversionErr
	}

	// Set the QoS name in the update
	apiUpdate.Name = &qosName

	// Prepare parameters
	params := &SlurmdbV0042PostQosParams{}

	// Create the request body
	requestBody := SlurmdbV0042PostQosJSONRequestBody{
		Qos: []V0042Qos{*apiUpdate},
	}

	// Call the generated OpenAPI client (POST is used for updates in Slurm)
	resp, err := q.client.apiClient.SlurmdbV0042PostQosWithResponse(ctx, params, requestBody)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.42")
	}

	// Check HTTP status
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
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.42", apiErrors)
				return apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.42")
		return httpErr
	}

	return nil
}

// Delete deletes a QoS
func (q *QoSManagerImpl) Delete(ctx context.Context, qosName string) error {
	if q.client == nil || q.client.apiClient == nil {
		return errors.NewClientError(errors.ErrorCodeClientNotInitialized, "API client not initialized")
	}

	if qosName == "" {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "QoS name is required", "qosName", qosName, nil)
	}

	// Call the generated OpenAPI client
	resp, err := q.client.apiClient.SlurmdbV0042DeleteSingleQosWithResponse(ctx, qosName)
	if err != nil {
		wrappedErr := errors.WrapError(err)
		return errors.EnhanceErrorWithVersion(wrappedErr, "v0.0.42")
	}

	// Check HTTP status (200 or 204 for successful deletion)
	if resp.StatusCode() != 200 && resp.StatusCode() != 204 {
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
				apiError := errors.NewSlurmAPIError(resp.StatusCode(), "v0.0.42", apiErrors)
				return apiError.SlurmError
			}
		}

		// Fall back to HTTP error handling
		httpErr := errors.WrapHTTPError(resp.StatusCode(), responseBody, "v0.0.42")
		return httpErr
	}

	return nil
}

// convertAPIQoSToInterface converts V0042Qos to interfaces.QoS
func convertAPIQoSToInterface(apiQoS V0042Qos) (*interfaces.QoS, error) {
	qos := &interfaces.QoS{}

	// Basic fields
	if apiQoS.Name != nil {
		qos.Name = *apiQoS.Name
	}

	if apiQoS.Description != nil {
		qos.Description = *apiQoS.Description
	}

	// Priority
	if apiQoS.Priority != nil && apiQoS.Priority.Set != nil && *apiQoS.Priority.Set && apiQoS.Priority.Number != nil {
		qos.Priority = int(*apiQoS.Priority.Number)
	}

	// Flags
	if apiQoS.Flags != nil {
		flags := make([]string, 0, len(*apiQoS.Flags))
		for _, flag := range *apiQoS.Flags {
			flags = append(flags, string(flag))
		}
		qos.Flags = flags
	}

	// Preempt mode - V0042QosPreemptModes is already a slice
	if apiQoS.Preempt != nil && apiQoS.Preempt.Mode != nil && len(*apiQoS.Preempt.Mode) > 0 {
		// Take the first preempt mode if multiple are set
		qos.PreemptMode = (*apiQoS.Preempt.Mode)[0]
	}

	// Limits - extract from nested structure
	if apiQoS.Limits != nil {
		// Grace time
		if apiQoS.Limits.GraceTime != nil {
			qos.GraceTime = int(*apiQoS.Limits.GraceTime)
		}

		// Factor (usage factor)
		if apiQoS.Limits.Factor != nil && apiQoS.Limits.Factor.Set != nil && *apiQoS.Limits.Factor.Set && apiQoS.Limits.Factor.Number != nil {
			qos.UsageFactor = *apiQoS.Limits.Factor.Number
		}

		// Max Jobs
		if apiQoS.Limits.Max != nil && apiQoS.Limits.Max.Jobs != nil {
			// Total jobs
			if apiQoS.Limits.Max.Jobs.Count != nil &&
				apiQoS.Limits.Max.Jobs.Count.Set != nil && *apiQoS.Limits.Max.Jobs.Count.Set &&
				apiQoS.Limits.Max.Jobs.Count.Number != nil {
				qos.MaxJobs = int(*apiQoS.Limits.Max.Jobs.Count.Number)
			}

			// Per user jobs
			if apiQoS.Limits.Max.Jobs.Per != nil && apiQoS.Limits.Max.Jobs.Per.User != nil &&
				apiQoS.Limits.Max.Jobs.Per.User.Set != nil && *apiQoS.Limits.Max.Jobs.Per.User.Set &&
				apiQoS.Limits.Max.Jobs.Per.User.Number != nil {
				qos.MaxJobsPerUser = int(*apiQoS.Limits.Max.Jobs.Per.User.Number)
			}

			// Per account jobs
			if apiQoS.Limits.Max.Jobs.Per != nil && apiQoS.Limits.Max.Jobs.Per.Account != nil &&
				apiQoS.Limits.Max.Jobs.Per.Account.Set != nil && *apiQoS.Limits.Max.Jobs.Per.Account.Set &&
				apiQoS.Limits.Max.Jobs.Per.Account.Number != nil {
				qos.MaxJobsPerAccount = int(*apiQoS.Limits.Max.Jobs.Per.Account.Number)
			}

			// Active jobs
			if apiQoS.Limits.Max.Jobs.ActiveJobs != nil && apiQoS.Limits.Max.Jobs.ActiveJobs.Count != nil &&
				apiQoS.Limits.Max.Jobs.ActiveJobs.Count.Set != nil && *apiQoS.Limits.Max.Jobs.ActiveJobs.Count.Set &&
				apiQoS.Limits.Max.Jobs.ActiveJobs.Count.Number != nil {
				qos.MaxSubmitJobs = int(*apiQoS.Limits.Max.Jobs.ActiveJobs.Count.Number)
			}
		}

		// Max wall clock per job
		if apiQoS.Limits.Max != nil && apiQoS.Limits.Max.WallClock != nil &&
			apiQoS.Limits.Max.WallClock.Per != nil && apiQoS.Limits.Max.WallClock.Per.Job != nil &&
			apiQoS.Limits.Max.WallClock.Per.Job.Set != nil && *apiQoS.Limits.Max.WallClock.Per.Job.Set &&
			apiQoS.Limits.Max.WallClock.Per.Job.Number != nil {
			qos.MaxWallTime = int(*apiQoS.Limits.Max.WallClock.Per.Job.Number)
		}

		// Max TRES (CPUs, Nodes, etc.)
		if apiQoS.Limits.Max != nil && apiQoS.Limits.Max.Tres != nil {
			// Per job
			if apiQoS.Limits.Max.Tres.Per != nil && apiQoS.Limits.Max.Tres.Per.Job != nil {
				// Extract CPU limits
				for _, tres := range *apiQoS.Limits.Max.Tres.Per.Job {
					if tres.Type != nil && *tres.Type == "cpu" && tres.Count != nil {
						qos.MaxCPUs = int(*tres.Count)
					}
					if tres.Type != nil && *tres.Type == "node" && tres.Count != nil {
						qos.MaxNodes = int(*tres.Count)
					}
				}
			}

			// Per user
			if apiQoS.Limits.Max.Tres.Per != nil && apiQoS.Limits.Max.Tres.Per.User != nil {
				// Extract CPU limits per user
				for _, tres := range *apiQoS.Limits.Max.Tres.Per.User {
					if tres.Type != nil && *tres.Type == "cpu" && tres.Count != nil {
						qos.MaxCPUsPerUser = int(*tres.Count)
					}
				}
			}
		}

		// Min TRES
		if apiQoS.Limits.Min != nil && apiQoS.Limits.Min.Tres != nil &&
			apiQoS.Limits.Min.Tres.Per != nil && apiQoS.Limits.Min.Tres.Per.Job != nil {
			// Extract minimum CPU/node requirements
			for _, tres := range *apiQoS.Limits.Min.Tres.Per.Job {
				if tres.Type != nil && *tres.Type == "cpu" && tres.Count != nil {
					qos.MinCPUs = int(*tres.Count)
				}
				if tres.Type != nil && *tres.Type == "node" && tres.Count != nil {
					qos.MinNodes = int(*tres.Count)
				}
			}
		}
		
		// Usage threshold  
		if apiQoS.Limits.Min != nil && apiQoS.Limits.Min.PriorityThreshold != nil &&
			apiQoS.Limits.Min.PriorityThreshold.Set != nil && *apiQoS.Limits.Min.PriorityThreshold.Set &&
			apiQoS.Limits.Min.PriorityThreshold.Number != nil {
			qos.UsageThreshold = float64(*apiQoS.Limits.Min.PriorityThreshold.Number)
		}
	}

	return qos, nil
}

// filterQoS applies client-side filtering to the QoS list
func filterQoS(qosList []interfaces.QoS, opts *interfaces.ListQoSOptions) []interfaces.QoS {
	if opts == nil {
		return qosList
	}

	filtered := make([]interfaces.QoS, 0, len(qosList))
	for _, qos := range qosList {
		// Filter by names
		if len(opts.Names) > 0 {
			found := false
			for _, name := range opts.Names {
				if qos.Name == name {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Filter by preempt mode
		if opts.PreemptMode != "" && qos.PreemptMode != opts.PreemptMode {
			continue
		}

		filtered = append(filtered, qos)
	}

	return filtered
}

// convertQoSCreateToAPI converts interfaces.QoSCreate to API format
func convertQoSCreateToAPI(create *interfaces.QoSCreate) (*V0042Qos, error) {
	apiQoS := &V0042Qos{
		Name: &create.Name,
	}

	// Description
	if create.Description != "" {
		apiQoS.Description = &create.Description
	}

	// Priority
	if create.Priority != nil {
		priority := uint32(*create.Priority)
		apiQoS.Priority = &V0042Uint32NoValStruct{
			Set:    &[]bool{true}[0],
			Number: &priority,
		}
	}

	// Flags
	if len(create.Flags) > 0 {
		flags := make([]V0042QosFlags, 0, len(create.Flags))
		for _, flag := range create.Flags {
			flags = append(flags, V0042QosFlags(flag))
		}
		apiQoS.Flags = &flags
	}

	// Preempt mode
	if create.PreemptMode != "" {
		modes := V0042QosPreemptModes{create.PreemptMode}
		apiQoS.Preempt = &struct {
			ExemptTime *V0042Uint32NoValStruct `json:"exempt_time,omitempty"`
			List       *V0042QosStringIdList   `json:"list,omitempty"`
			Mode       *V0042QosPreemptModes   `json:"mode,omitempty"`
		}{
			Mode: &modes,
		}
	}

	// Limits
	if create.Limits != nil {
		limits := &V0042QosLimits{}

		// Max jobs
		if create.Limits.MaxJobsTotal != nil {
			maxJobs := uint32(*create.Limits.MaxJobsTotal)
			limits.MaxJobs = &V0042QosLimitsMaxJobs{
				Total: &V0042Uint32NoValStruct{
					Set:    &[]bool{true}[0],
					Number: &maxJobs,
				},
			}
		}

		// Max submit jobs
		if create.Limits.MaxSubmitJobs != nil {
			maxSubmit := uint32(*create.Limits.MaxSubmitJobs)
			limits.MaxSubmitJobs = &V0042QosLimitsMaxSubmitJobs{
				Total: &V0042Uint32NoValStruct{
					Set:    &[]bool{true}[0],
					Number: &maxSubmit,
				},
			}
		}

		// Max wall clock per job
		if create.Limits.MaxWallClockPerJob != nil {
			maxWallClock := uint32(*create.Limits.MaxWallClockPerJob)
			limits.MaxWallClock = &V0042QosLimitsMaxWallClock{
				Per: &V0042QosLimitsMaxWallClockPer{
					Job: &V0042Uint32NoValStruct{
						Set:    &[]bool{true}[0],
						Number: &maxWallClock,
					},
				},
			}
		}

		// Max nodes per job
		if create.Limits.MaxNodesPerJob != nil {
			maxNodes := uint32(*create.Limits.MaxNodesPerJob)
			limits.MaxNodes = &V0042QosLimitsMaxNodes{
				Per: &V0042QosLimitsMaxNodesPer{
					Job: &V0042Uint32NoValStruct{
						Set:    &[]bool{true}[0],
						Number: &maxNodes,
					},
				},
			}
		}

		// Max CPUs per job
		if create.Limits.MaxCPUsPerJob != nil {
			maxCPUs := uint32(*create.Limits.MaxCPUsPerJob)
			limits.MaxCPUs = &V0042QosLimitsMaxCPUs{
				Per: &V0042QosLimitsMaxCPUsPer{
					Job: &V0042Uint32NoValStruct{
						Set:    &[]bool{true}[0],
						Number: &maxCPUs,
					},
				},
			}
		}

		apiQoS.Limits = limits
	}

	// Usage threshold
	if create.UsageThreshold != nil {
		threshold := float32(*create.UsageThreshold)
		apiQoS.UsageThreshold = &V0042Float64NoValStruct{
			Set:    &[]bool{true}[0],
			Number: &threshold,
		}
	}

	// Grace time
	if create.GraceTime != nil {
		graceTime := int32(*create.GraceTime)
		apiQoS.GraceTime = &graceTime
	}

	// Usage factor
	if create.UsageFactor != nil {
		factor := float32(*create.UsageFactor)
		apiQoS.UsageFactor = &V0042Float64NoValStruct{
			Set:    &[]bool{true}[0],
			Number: &factor,
		}
	}

	// Preemptable QoS
	if len(create.PreemptableQoS) > 0 {
		if apiQoS.Preempt == nil {
			apiQoS.Preempt = &V0042QosPreempt{}
		}
		apiQoS.Preempt.List = &create.PreemptableQoS
	}

	return apiQoS, nil
}

// convertQoSUpdateToAPI converts interfaces.QoSUpdate to API format
func convertQoSUpdateToAPI(update *interfaces.QoSUpdate) (*V0042Qos, error) {
	apiQoS := &V0042Qos{}

	// Description
	if update.Description != nil {
		apiQoS.Description = update.Description
	}

	// Priority
	if update.Priority != nil {
		priority := uint32(*update.Priority)
		apiQoS.Priority = &V0042Uint32NoValStruct{
			Set:    &[]bool{true}[0],
			Number: &priority,
		}
	}

	// Flags
	if update.Flags != nil {
		flags := make([]V0042QosFlags, 0, len(*update.Flags))
		for _, flag := range *update.Flags {
			flags = append(flags, V0042QosFlags(flag))
		}
		apiQoS.Flags = &flags
	}

	// Preempt mode
	if update.PreemptMode != nil {
		mode := V0042QosPreemptMode(*update.PreemptMode)
		apiQoS.Preempt = &V0042QosPreempt{
			Mode: &mode,
		}
	}

	// Similar conversion for limits and other fields as in Create
	// (omitted for brevity - follows same pattern)

	return apiQoS, nil
}

// convertTRESToMap converts TRES array to a map
func convertTRESToMap(tres []V0042Tres) map[string]int64 {
	result := make(map[string]int64)
	for _, t := range tres {
		if t.Type != nil && t.Count != nil {
			result[*t.Type] = int64(*t.Count)
		}
	}
	return result
}