// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package factory

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jontk/slurm-client/internal/adapters/common"
	v040adapter "github.com/jontk/slurm-client/internal/adapters/v0_0_40"
	v041adapter "github.com/jontk/slurm-client/internal/adapters/v0_0_41"
	v042adapter "github.com/jontk/slurm-client/internal/adapters/v0_0_42"
	v043adapter "github.com/jontk/slurm-client/internal/adapters/v0_0_43"
	v040api "github.com/jontk/slurm-client/internal/api/v0_0_40"
	v041api "github.com/jontk/slurm-client/internal/api/v0_0_41"
	v042api "github.com/jontk/slurm-client/internal/api/v0_0_42"
	v043api "github.com/jontk/slurm-client/internal/api/v0_0_43"
	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/internal/common/types"
)

// AdapterClient wraps a version-specific adapter to implement the SlurmClient interface
type AdapterClient struct {
	adapter common.VersionAdapter
	version string
}

// NewAdapterClient creates a new adapter-based client for the specified version
func NewAdapterClient(version string, config *interfaces.ClientConfig) (SlurmClient, error) {
	switch version {
	case "v0.0.40":
		client, err := v040api.NewClientWithResponses(config.BaseURL, v040api.WithHTTPClient(config.HTTPClient))
		if err != nil {
			return nil, fmt.Errorf("failed to create v0.0.40 client: %w", err)
		}
		adapter := v040adapter.NewAdapter(client)
		return &AdapterClient{
			adapter: adapter,
			version: version,
		}, nil

	case "v0.0.41":
		client, err := v041api.NewClientWithResponses(config.BaseURL, v041api.WithHTTPClient(config.HTTPClient))
		if err != nil {
			return nil, fmt.Errorf("failed to create v0.0.41 client: %w", err)
		}
		adapter := v041adapter.NewAdapter(client)
		return &AdapterClient{
			adapter: adapter,
			version: version,
		}, nil

	case "v0.0.42":
		client, err := v042api.NewClientWithResponses(config.BaseURL, v042api.WithHTTPClient(config.HTTPClient))
		if err != nil {
			return nil, fmt.Errorf("failed to create v0.0.42 client: %w", err)
		}
		adapter := v042adapter.NewAdapter(client)
		return &AdapterClient{
			adapter: adapter,
			version: version,
		}, nil

	case "v0.0.43":
		client, err := v043api.NewClientWithResponses(config.BaseURL, v043api.WithHTTPClient(config.HTTPClient))
		if err != nil {
			return nil, fmt.Errorf("failed to create v0.0.43 client: %w", err)
		}
		adapter := v043adapter.NewAdapter(client)
		return &AdapterClient{
			adapter: adapter,
			version: version,
		}, nil

	default:
		return nil, fmt.Errorf("unsupported version: %s", version)
	}
}

// Version returns the API version
func (c *AdapterClient) Version() string {
	return c.version
}

// Jobs returns the JobManager
func (c *AdapterClient) Jobs() interfaces.JobManager {
	return &adapterJobManager{adapter: c.adapter.GetJobManager()}
}

// Nodes returns the NodeManager
func (c *AdapterClient) Nodes() interfaces.NodeManager {
	return &adapterNodeManager{adapter: c.adapter.GetNodeManager()}
}

// Partitions returns the PartitionManager
func (c *AdapterClient) Partitions() interfaces.PartitionManager {
	return &adapterPartitionManager{adapter: c.adapter.GetPartitionManager()}
}

// Info returns the InfoManager
func (c *AdapterClient) Info() interfaces.InfoManager {
	// For now, return a basic implementation
	return &adapterInfoManager{version: c.version}
}

// Reservations returns the ReservationManager
func (c *AdapterClient) Reservations() interfaces.ReservationManager {
	return &adapterReservationManager{adapter: c.adapter.GetReservationManager()}
}

// QoS returns the QoSManager
func (c *AdapterClient) QoS() interfaces.QoSManager {
	return &adapterQoSManager{adapter: c.adapter.GetQoSManager()}
}

// Accounts returns the AccountManager
func (c *AdapterClient) Accounts() interfaces.AccountManager {
	return &adapterAccountManager{adapter: c.adapter.GetAccountManager()}
}

// Users returns the UserManager
func (c *AdapterClient) Users() interfaces.UserManager {
	return &adapterUserManager{adapter: c.adapter.GetUserManager()}
}

// Clusters returns the ClusterManager
func (c *AdapterClient) Clusters() interfaces.ClusterManager {
	// Clusters are not supported in the adapter pattern yet
	return nil
}

// Associations returns the AssociationManager
func (c *AdapterClient) Associations() interfaces.AssociationManager {
	return &adapterAssociationManager{adapter: c.adapter.GetAssociationManager()}
}

// Close closes the client
func (c *AdapterClient) Close() error {
	// No resources to close
	return nil
}

// === Standalone Operations ===

// GetLicenses retrieves license information
func (c *AdapterClient) GetLicenses(ctx context.Context) (*interfaces.LicenseList, error) {
	return nil, fmt.Errorf("GetLicenses not implemented in adapter")
}

// GetShares retrieves fairshare information with optional filtering
func (c *AdapterClient) GetShares(ctx context.Context, opts *interfaces.GetSharesOptions) (*interfaces.SharesList, error) {
	return nil, fmt.Errorf("GetShares not implemented in adapter")
}

// GetConfig retrieves SLURM configuration
func (c *AdapterClient) GetConfig(ctx context.Context) (*interfaces.Config, error) {
	return nil, fmt.Errorf("GetConfig not implemented in adapter")
}

// GetDiagnostics retrieves SLURM diagnostics information
func (c *AdapterClient) GetDiagnostics(ctx context.Context) (*interfaces.Diagnostics, error) {
	return nil, fmt.Errorf("GetDiagnostics not implemented in adapter")
}

// GetDBDiagnostics retrieves SLURM database diagnostics information
func (c *AdapterClient) GetDBDiagnostics(ctx context.Context) (*interfaces.Diagnostics, error) {
	return nil, fmt.Errorf("GetDBDiagnostics not implemented in adapter")
}

// GetInstance retrieves a specific database instance
func (c *AdapterClient) GetInstance(ctx context.Context, opts *interfaces.GetInstanceOptions) (*interfaces.Instance, error) {
	return nil, fmt.Errorf("GetInstance not implemented in adapter")
}

// GetInstances retrieves multiple database instances with filtering
func (c *AdapterClient) GetInstances(ctx context.Context, opts *interfaces.GetInstancesOptions) (*interfaces.InstanceList, error) {
	return nil, fmt.Errorf("GetInstances not implemented in adapter")
}

// GetTRES retrieves all TRES (Trackable RESources)
func (c *AdapterClient) GetTRES(ctx context.Context) (*interfaces.TRESList, error) {
	return nil, fmt.Errorf("GetTRES not implemented in adapter")
}

// CreateTRES creates a new TRES entry
func (c *AdapterClient) CreateTRES(ctx context.Context, req *interfaces.CreateTRESRequest) (*interfaces.TRES, error) {
	return nil, fmt.Errorf("CreateTRES not implemented in adapter")
}

// Reconfigure triggers a SLURM reconfiguration
func (c *AdapterClient) Reconfigure(ctx context.Context) (*interfaces.ReconfigureResponse, error) {
	return nil, fmt.Errorf("Reconfigure not implemented in adapter")
}

// Manager wrappers to convert between common.types and interfaces types

// adapterJobManager wraps a common.JobAdapter to implement interfaces.JobManager
type adapterJobManager struct {
	adapter common.JobAdapter
}

func (m *adapterJobManager) List(ctx context.Context, opts *interfaces.ListJobsOptions) (*interfaces.JobList, error) {
	// Convert options
	adapterOpts := &types.JobListOptions{}
	if opts != nil {
		adapterOpts.Users = []string{opts.UserID}  // Convert UserID to array
		adapterOpts.Partitions = []string{opts.Partition}  // Convert Partition to array
		adapterOpts.Limit = opts.Limit
		adapterOpts.Offset = opts.Offset
		// Convert states
		for _, s := range opts.States {
			adapterOpts.States = append(adapterOpts.States, types.JobState(s))
		}
	}

	// Call adapter
	result, err := m.adapter.List(ctx, adapterOpts)
	if err != nil {
		return nil, err
	}

	// Convert result
	jobList := &interfaces.JobList{
		Jobs: make([]interfaces.Job, 0, len(result.Jobs)),
		Total: len(result.Jobs), // Use actual count since Meta may not exist
	}

	for _, job := range result.Jobs {
		jobList.Jobs = append(jobList.Jobs, convertJobToInterface(job))
	}

	return jobList, nil
}

func (m *adapterJobManager) Get(ctx context.Context, jobID string) (*interfaces.Job, error) {
	// Convert string to int32 for adapter
	jobIDInt, err := strconv.ParseInt(jobID, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid job ID: %w", err)
	}
	
	job, err := m.adapter.Get(ctx, int32(jobIDInt))
	if err != nil {
		return nil, err
	}
	result := convertJobToInterface(*job)
	return &result, nil
}

func (m *adapterJobManager) Submit(ctx context.Context, job *interfaces.JobSubmission) (*interfaces.JobSubmitResponse, error) {
	// Convert submission - map from interfaces.JobSubmission to types.JobCreate
	priority := int32(job.Priority)
	submission := &types.JobCreate{
		Name:             job.Name,
		Account:          job.Account,
		Script:           job.Script,
		Command:          job.Command,
		Partition:        job.Partition,  
		CPUs:             int32(job.CPUs),
		TimeLimit:        int32(job.TimeLimit),
		WorkingDirectory: job.WorkingDir,
		Environment:      job.Environment,
		Nodes:            int32(job.Nodes),
		Priority:         &priority,
		ResourceRequests: types.ResourceRequests{
			Memory: int64(job.Memory),
		},
	}

	// Call adapter
	resp, err := m.adapter.Submit(ctx, submission)
	if err != nil {
		return nil, err
	}

	return &interfaces.JobSubmitResponse{
		JobID: fmt.Sprintf("%d", resp.JobID), // Convert int32 to string
	}, nil
}

func (m *adapterJobManager) Update(ctx context.Context, jobID string, update *interfaces.JobUpdate) error {
	// Convert string to int32 for adapter
	jobIDInt, err := strconv.ParseInt(jobID, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid job ID: %w", err)
	}
	// Convert update - only use available fields from interfaces.JobUpdate
	adapterUpdate := &types.JobUpdate{
		Name: update.Name,
	}
	
	// Convert time limit if present
	if update.TimeLimit != nil {
		timeLimit := int32(*update.TimeLimit)
		adapterUpdate.TimeLimit = &timeLimit
	}
	
	// Convert priority if present
	if update.Priority != nil {
		priority := int32(*update.Priority)
		adapterUpdate.Priority = &priority
	}

	return m.adapter.Update(ctx, int32(jobIDInt), adapterUpdate)
}

func (m *adapterJobManager) Cancel(ctx context.Context, jobID string) error {
	// Convert string to int32 for adapter
	jobIDInt, err := strconv.ParseInt(jobID, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid job ID: %w", err)
	}
	return m.adapter.Cancel(ctx, int32(jobIDInt), nil)
}

func (m *adapterJobManager) Watch(ctx context.Context, opts *interfaces.WatchJobsOptions) (<-chan interfaces.JobEvent, error) {
	// Convert WatchJobsOptions to types.JobWatchOptions
	adapterOpts := &types.JobWatchOptions{}
	
	if opts != nil {
		// Convert JobIDs from []string to []int32
		if len(opts.JobIDs) > 0 {
			// Just watch the first job ID for now (adapter expects single job ID)
			jobIDInt, err := strconv.ParseInt(opts.JobIDs[0], 10, 32)
			if err == nil {
				adapterOpts.JobID = int32(jobIDInt)
			}
		}
		
		// Convert state filters
		if len(opts.States) > 0 {
			adapterOpts.EventTypes = opts.States
		}
	}
	
	// Call adapter's Watch method
	adapterEventChan, err := m.adapter.Watch(ctx, adapterOpts)
	if err != nil {
		return nil, err
	}
	
	// Create interface event channel
	interfaceEventChan := make(chan interfaces.JobEvent, 10)
	
	// Start goroutine to convert events
	go func() {
		defer close(interfaceEventChan)
		
		for adapterEvent := range adapterEventChan {
			// Convert types.JobWatchEvent to interfaces.JobEvent
			interfaceEvent := interfaces.JobEvent{
				Type:      adapterEvent.EventType,
				JobID:     strconv.Itoa(int(adapterEvent.JobID)),
				OldState:  string(adapterEvent.PreviousState),
				NewState:  string(adapterEvent.NewState),
				Timestamp: adapterEvent.EventTime,
			}
			
			select {
			case interfaceEventChan <- interfaceEvent:
			case <-ctx.Done():
				return
			}
		}
	}()
	
	return interfaceEventChan, nil
}

func (m *adapterJobManager) AnalyzeBatchJobs(ctx context.Context, jobIDs []string, opts *interfaces.BatchAnalysisOptions) (*interfaces.BatchJobAnalysis, error) {
	// AnalyzeBatchJobs is not implemented in adapters
	return nil, fmt.Errorf("AnalyzeBatchJobs not implemented in adapter")
}

// Analytics methods for resource utilization and performance
func (m *adapterJobManager) Steps(ctx context.Context, jobID string) (*interfaces.JobStepList, error) {
	return nil, fmt.Errorf("Steps not implemented in adapter")
}

func (m *adapterJobManager) GetJobUtilization(ctx context.Context, jobID string) (*interfaces.JobUtilization, error) {
	return nil, fmt.Errorf("GetJobUtilization not implemented in adapter")
}

func (m *adapterJobManager) GetJobEfficiency(ctx context.Context, jobID string) (*interfaces.ResourceUtilization, error) {
	return nil, fmt.Errorf("GetJobEfficiency not implemented in adapter")
}

func (m *adapterJobManager) GetJobPerformance(ctx context.Context, jobID string) (*interfaces.JobPerformance, error) {
	return nil, fmt.Errorf("GetJobPerformance not implemented in adapter")
}

func (m *adapterJobManager) GetJobLiveMetrics(ctx context.Context, jobID string) (*interfaces.JobLiveMetrics, error) {
	return nil, fmt.Errorf("GetJobLiveMetrics not implemented in adapter")
}

func (m *adapterJobManager) WatchJobMetrics(ctx context.Context, jobID string, opts *interfaces.WatchMetricsOptions) (<-chan interfaces.JobMetricsEvent, error) {
	return nil, fmt.Errorf("WatchJobMetrics not implemented in adapter")
}

func (m *adapterJobManager) GetJobResourceTrends(ctx context.Context, jobID string, opts *interfaces.ResourceTrendsOptions) (*interfaces.JobResourceTrends, error) {
	return nil, fmt.Errorf("GetJobResourceTrends not implemented in adapter")
}

func (m *adapterJobManager) GetJobStepDetails(ctx context.Context, jobID string, stepID string) (*interfaces.JobStepDetails, error) {
	return nil, fmt.Errorf("GetJobStepDetails not implemented in adapter")
}

func (m *adapterJobManager) GetJobStepUtilization(ctx context.Context, jobID string, stepID string) (*interfaces.JobStepUtilization, error) {
	return nil, fmt.Errorf("GetJobStepUtilization not implemented in adapter")
}

func (m *adapterJobManager) ListJobStepsWithMetrics(ctx context.Context, jobID string, opts *interfaces.ListJobStepsOptions) (*interfaces.JobStepMetricsList, error) {
	return nil, fmt.Errorf("ListJobStepsWithMetrics not implemented in adapter")
}

// SLURM Integration Methods
func (m *adapterJobManager) GetJobStepsFromAccounting(ctx context.Context, jobID string, opts *interfaces.AccountingQueryOptions) (*interfaces.AccountingJobSteps, error) {
	return nil, fmt.Errorf("GetJobStepsFromAccounting not implemented in adapter")
}

func (m *adapterJobManager) GetStepAccountingData(ctx context.Context, jobID string, stepID string) (*interfaces.StepAccountingRecord, error) {
	return nil, fmt.Errorf("GetStepAccountingData not implemented in adapter")
}

func (m *adapterJobManager) GetJobStepAPIData(ctx context.Context, jobID string, stepID string) (*interfaces.JobStepAPIData, error) {
	return nil, fmt.Errorf("GetJobStepAPIData not implemented in adapter")
}

func (m *adapterJobManager) ListJobStepsFromSacct(ctx context.Context, jobID string, opts *interfaces.SacctQueryOptions) (*interfaces.SacctJobStepData, error) {
	return nil, fmt.Errorf("ListJobStepsFromSacct not implemented in adapter")
}

// Advanced Analytics Methods
func (m *adapterJobManager) GetJobCPUAnalytics(ctx context.Context, jobID string) (*interfaces.CPUAnalytics, error) {
	return nil, fmt.Errorf("GetJobCPUAnalytics not implemented in adapter")
}

func (m *adapterJobManager) GetJobMemoryAnalytics(ctx context.Context, jobID string) (*interfaces.MemoryAnalytics, error) {
	return nil, fmt.Errorf("GetJobMemoryAnalytics not implemented in adapter")
}

func (m *adapterJobManager) GetJobIOAnalytics(ctx context.Context, jobID string) (*interfaces.IOAnalytics, error) {
	return nil, fmt.Errorf("GetJobIOAnalytics not implemented in adapter")
}

func (m *adapterJobManager) GetJobComprehensiveAnalytics(ctx context.Context, jobID string) (*interfaces.JobComprehensiveAnalytics, error) {
	return nil, fmt.Errorf("GetJobComprehensiveAnalytics not implemented in adapter")
}

// Historical Performance Tracking Methods
func (m *adapterJobManager) GetJobPerformanceHistory(ctx context.Context, jobID string, opts *interfaces.PerformanceHistoryOptions) (*interfaces.JobPerformanceHistory, error) {
	return nil, fmt.Errorf("GetJobPerformanceHistory not implemented in adapter")
}

func (m *adapterJobManager) GetPerformanceTrends(ctx context.Context, opts *interfaces.TrendAnalysisOptions) (*interfaces.PerformanceTrends, error) {
	return nil, fmt.Errorf("GetPerformanceTrends not implemented in adapter")
}

func (m *adapterJobManager) GetUserEfficiencyTrends(ctx context.Context, userID string, opts *interfaces.EfficiencyTrendOptions) (*interfaces.UserEfficiencyTrends, error) {
	return nil, fmt.Errorf("GetUserEfficiencyTrends not implemented in adapter")
}

func (m *adapterJobManager) GetWorkflowPerformance(ctx context.Context, workflowID string, opts *interfaces.WorkflowAnalysisOptions) (*interfaces.WorkflowPerformance, error) {
	return nil, fmt.Errorf("GetWorkflowPerformance not implemented in adapter")
}

func (m *adapterJobManager) GenerateEfficiencyReport(ctx context.Context, opts *interfaces.ReportOptions) (*interfaces.EfficiencyReport, error) {
	return nil, fmt.Errorf("GenerateEfficiencyReport not implemented in adapter")
}

// Helper function to convert types.Job to interfaces.Job
func convertJobToInterface(job types.Job) interfaces.Job {
	// Convert node list to slice
	nodes := []string{}
	if job.NodeList != "" {
		nodes = []string{job.NodeList} // Simple conversion - could be improved to split properly
	}
	
	return interfaces.Job{
		ID:          fmt.Sprintf("%d", job.JobID),
		Name:        job.Name,
		UserID:      fmt.Sprintf("%d", job.UserID),
		GroupID:     fmt.Sprintf("%d", job.GroupID),
		State:       string(job.State),
		Partition:   job.Partition,
		Priority:    int(job.Priority),
		SubmitTime:  job.SubmitTime,
		StartTime:   job.StartTime,
		EndTime:     job.EndTime,
		CPUs:        int(job.CPUs),
		Memory:      int(job.MinMemory), // Use MinMemory as the closest match
		TimeLimit:   int(job.TimeLimit),
		WorkingDir:  job.WorkingDirectory,
		Command:     job.Command,
		Environment: job.Environment,
		Nodes:       nodes,
		ExitCode:    0, // Not available in types.Job
		Metadata:    make(map[string]interface{}),
	}
}

// adapterNodeManager wraps a common.NodeAdapter to implement interfaces.NodeManager
type adapterNodeManager struct {
	adapter common.NodeAdapter
}

func (m *adapterNodeManager) List(ctx context.Context, opts *interfaces.ListNodesOptions) (*interfaces.NodeList, error) {
	// Convert options
	adapterOpts := &types.NodeListOptions{}
	if opts != nil {
		if opts.States != nil {
			for _, state := range opts.States {
				adapterOpts.States = append(adapterOpts.States, types.NodeState(state))
			}
		}
		if opts.Partition != "" {
			adapterOpts.Partitions = []string{opts.Partition}
		}
		adapterOpts.Limit = opts.Limit
		adapterOpts.Offset = opts.Offset
	}

	// Call adapter
	result, err := m.adapter.List(ctx, adapterOpts)
	if err != nil {
		return nil, err
	}

	// Convert result
	nodeList := &interfaces.NodeList{
		Nodes: make([]interfaces.Node, 0, len(result.Nodes)),
		Total: result.Total,
	}

	for _, node := range result.Nodes {
		nodeList.Nodes = append(nodeList.Nodes, convertNodeToInterface(node))
	}

	return nodeList, nil
}

func (m *adapterNodeManager) Get(ctx context.Context, nodeName string) (*interfaces.Node, error) {
	node, err := m.adapter.Get(ctx, nodeName)
	if err != nil {
		return nil, err
	}
	result := convertNodeToInterface(*node)
	return &result, nil
}

func (m *adapterNodeManager) Update(ctx context.Context, nodeName string, update *interfaces.NodeUpdate) error {
	// Convert update
	adapterUpdate := &types.NodeUpdate{
		State:    (*types.NodeState)(update.State),
		Reason:   update.Reason,
		Features: update.Features,
	}

	return m.adapter.Update(ctx, nodeName, adapterUpdate)
}

func (m *adapterNodeManager) Watch(ctx context.Context, opts *interfaces.WatchNodesOptions) (<-chan interfaces.NodeEvent, error) {
	// Convert WatchNodesOptions to types.NodeWatchOptions
	adapterOpts := &types.NodeWatchOptions{}
	
	if opts != nil {
		// Convert node names
		if len(opts.NodeNames) > 0 {
			adapterOpts.NodeNames = opts.NodeNames
		}
		
		// Convert states
		if len(opts.States) > 0 {
			for _, state := range opts.States {
				adapterOpts.States = append(adapterOpts.States, types.NodeState(state))
			}
		}
		
		// Convert partition
		if opts.Partition != "" {
			adapterOpts.Partitions = []string{opts.Partition}
		}
	}
	
	// Call adapter's Watch method
	adapterEventChan, err := m.adapter.Watch(ctx, adapterOpts)
	if err != nil {
		return nil, err
	}
	
	// Create interface event channel
	interfaceEventChan := make(chan interfaces.NodeEvent, 10)
	
	// Start goroutine to convert events
	go func() {
		defer close(interfaceEventChan)
		
		for adapterEvent := range adapterEventChan {
			// Convert types.NodeWatchEvent to interfaces.NodeEvent
			interfaceEvent := interfaces.NodeEvent{
				Type:      adapterEvent.EventType,
				NodeName:  adapterEvent.NodeName,
				OldState:  string(adapterEvent.PreviousState),
				NewState:  string(adapterEvent.NewState),
				Timestamp: adapterEvent.EventTime,
			}
			
			select {
			case interfaceEventChan <- interfaceEvent:
			case <-ctx.Done():
				return
			}
		}
	}()
	
	return interfaceEventChan, nil
}

// Helper function to convert types.Node to interfaces.Node
func convertNodeToInterface(node types.Node) interfaces.Node {
	// Copy features directly
	features := node.Features

	return interfaces.Node{
		Name:         node.Name,
		State:        string(node.State),
		CPUs:         int(node.CPUs),
		Memory:       int(node.Memory),
		Partitions:   node.Partitions,
		Features:     features,
		Reason:       node.Reason,
		LastBusy:     node.LastBusy,
		Architecture: node.Arch,
		Metadata:     make(map[string]interface{}),
	}
}

// adapterPartitionManager wraps a common.PartitionAdapter
type adapterPartitionManager struct {
	adapter common.PartitionAdapter
}

func (m *adapterPartitionManager) List(ctx context.Context, opts *interfaces.ListPartitionsOptions) (*interfaces.PartitionList, error) {
	// Convert options
	adapterOpts := &types.PartitionListOptions{}
	if opts != nil {
		adapterOpts.Names = opts.States // Using States as Names for now
		adapterOpts.Limit = opts.Limit
		adapterOpts.Offset = opts.Offset
	}

	// Call adapter
	result, err := m.adapter.List(ctx, adapterOpts)
	if err != nil {
		return nil, err
	}

	// Convert result
	partitionList := &interfaces.PartitionList{
		Partitions: make([]interfaces.Partition, 0, len(result.Partitions)),
		Total:      result.Total,
	}

	for _, partition := range result.Partitions {
		partitionList.Partitions = append(partitionList.Partitions, convertPartitionToInterface(partition))
	}

	return partitionList, nil
}

func (m *adapterPartitionManager) Get(ctx context.Context, partitionName string) (*interfaces.Partition, error) {
	partition, err := m.adapter.Get(ctx, partitionName)
	if err != nil {
		return nil, err
	}
	result := convertPartitionToInterface(*partition)
	return &result, nil
}

func (m *adapterPartitionManager) Update(ctx context.Context, partitionName string, update *interfaces.PartitionUpdate) error {
	// Convert update request
	adapterUpdate := &types.PartitionUpdate{}
	if update != nil {
		if update.MaxTime != nil {
			maxTime := int32(*update.MaxTime)
			adapterUpdate.MaxTime = &maxTime
		}
		if update.DefaultTime != nil {
			defaultTime := int32(*update.DefaultTime)
			adapterUpdate.DefaultTime = &defaultTime
		}
		if update.State != nil {
			// Convert string to PartitionState
			state := types.PartitionState(*update.State)
			adapterUpdate.State = &state
		}
		// Add other fields as needed
	}

	return m.adapter.Update(ctx, partitionName, adapterUpdate)
}

func (m *adapterPartitionManager) Watch(ctx context.Context, opts *interfaces.WatchPartitionsOptions) (<-chan interfaces.PartitionEvent, error) {
	// Watch is not implemented in adapters
	return nil, fmt.Errorf("watch not implemented in adapter")
}

// Helper function to convert types.Partition to interfaces.Partition
func convertPartitionToInterface(partition types.Partition) interfaces.Partition {
	return interfaces.Partition{
		Name:           partition.Name,
		State:          string(partition.State),
		TotalNodes:     int(partition.TotalNodes),
		AvailableNodes: 0, // Not available in types.Partition
		TotalCPUs:      int(partition.TotalCPUs),
		IdleCPUs:       int(partition.TotalCPUs), // No idle CPUs field in types
		MaxTime:        int(partition.MaxTime),
		DefaultTime:    int(partition.DefaultTime),
		MaxMemory:      int(partition.MaxMemPerNode),
		DefaultMemory:  int(partition.DefMemPerNode),
		AllowedUsers:   []string{}, // Not available in types
		DeniedUsers:    []string{}, // Not available in types
		AllowedGroups:  partition.AllowGroups,
		DeniedGroups:   []string{}, // Not available as DenyGroups
		Priority:       int(partition.Priority),
		Nodes:          convertNodeStringToArray(partition.Nodes), // Convert from string to []string
	}
}

// Helper function to convert node string (comma-separated) to array
func convertNodeStringToArray(nodes string) []string {
	if nodes == "" {
		return []string{}
	}
	// Split by comma or other delimiters as needed
	return []string{nodes}
}

// Helper function to convert types.Account to interfaces.Account
func convertAccountToInterface(account types.Account) interfaces.Account {
	return interfaces.Account{
		Name:              account.Name,
		Description:       account.Description,
		Organization:      account.Organization,
		CoordinatorUsers:  account.Coordinators,
		AllowedPartitions: account.AllowedPartitions,
		DefaultPartition:  account.DefaultPartition,
		AllowedQoS:        account.QoSList,
		DefaultQoS:        account.DefaultQoS,
		CPULimit:          int(account.MaxCPUs),
		MaxJobs:           int(account.MaxJobs),
		MaxJobsPerUser:    int(account.MaxJobsPerUser),
		MaxNodes:          int(account.MaxNodes),
		MaxWallTime:       int(account.MaxWallTime),
		FairShareTRES:     make(map[string]int), // Will need proper conversion if available
		GrpTRES:           make(map[string]int), // Will need proper conversion if available
		GrpTRESMinutes:    make(map[string]int), // Will need proper conversion if available
		MaxTRES:           make(map[string]int), // Will need proper conversion if available
		MaxTRESPerUser:    make(map[string]int), // Will need proper conversion if available
		SharesPriority:    int(account.Priority),
		ParentAccount:     account.ParentName,
	}
}

// Helper function to convert types.User to interfaces.User
func convertUserToInterface(user types.User) interfaces.User {
	// Convert accounts to UserAccount format
	accounts := make([]interfaces.UserAccount, 0, len(user.Accounts))
	for _, accountName := range user.Accounts {
		accounts = append(accounts, interfaces.UserAccount{
			AccountName: accountName,
			// Other fields would need to be populated from associations
		})
	}
	
	// Convert coordinators
	coordinatorAccounts := make([]string, 0, len(user.Coordinators))
	for _, coord := range user.Coordinators {
		coordinatorAccounts = append(coordinatorAccounts, coord.AccountName)
	}

	return interfaces.User{
		Name:                user.Name,
		UID:                 int(user.UID),
		DefaultAccount:      user.DefaultAccount,
		DefaultWCKey:        user.DefaultWCKey,
		AdminLevel:          string(user.AdminLevel),
		CoordinatorAccounts: coordinatorAccounts,
		Accounts:            accounts,
		Quotas:              nil, // Would need proper conversion
		FairShare:           nil, // Would need proper conversion
		Associations:        nil, // Would need proper conversion
		Created:             time.Now(), // Not available in types.User
		Modified:            time.Now(), // Not available in types.User,
		Metadata:            nil,
	}
}

// Helper function to convert types.QoS to interfaces.QoS
func convertQoSToInterface(qos types.QoS) interfaces.QoS {
	// Extract limits from the QoS
	maxJobs := 0
	maxJobsPerUser := 0
	maxJobsPerAccount := 0
	maxSubmitJobs := 0
	maxCPUs := 0
	maxCPUsPerUser := 0
	maxNodes := 0
	maxWallTime := 0
	minCPUs := 0
	minNodes := 0

	if qos.Limits != nil {
		if qos.Limits.MaxJobsPerUser != nil {
			maxJobsPerUser = *qos.Limits.MaxJobsPerUser
		}
		if qos.Limits.MaxJobsPerAccount != nil {
			maxJobsPerAccount = *qos.Limits.MaxJobsPerAccount
		}
		if qos.Limits.MaxSubmitJobsPerUser != nil {
			maxSubmitJobs = *qos.Limits.MaxSubmitJobsPerUser
		}
		if qos.Limits.MaxCPUsPerUser != nil {
			maxCPUsPerUser = *qos.Limits.MaxCPUsPerUser
		}
		if qos.Limits.MaxCPUsPerJob != nil {
			maxCPUs = *qos.Limits.MaxCPUsPerJob
		}
		if qos.Limits.MaxNodesPerJob != nil {
			maxNodes = *qos.Limits.MaxNodesPerJob
		}
		if qos.Limits.MaxWallTimePerJob != nil {
			maxWallTime = *qos.Limits.MaxWallTimePerJob
		}
		if qos.Limits.MinCPUsPerJob != nil {
			minCPUs = *qos.Limits.MinCPUsPerJob
		}
		if qos.Limits.MinNodesPerJob != nil {
			minNodes = *qos.Limits.MinNodesPerJob
		}
	}

	return interfaces.QoS{
		Name:              qos.Name,
		Description:       qos.Description,
		Priority:          qos.Priority,
		PreemptMode:       qos.PreemptMode,
		GraceTime:         qos.GraceTime,
		MaxJobs:           maxJobs,
		MaxJobsPerUser:    maxJobsPerUser,
		MaxJobsPerAccount: maxJobsPerAccount,
		MaxSubmitJobs:     maxSubmitJobs,
		MaxCPUs:           maxCPUs,
		MaxCPUsPerUser:    maxCPUsPerUser,
		MaxNodes:          maxNodes,
		MaxWallTime:       maxWallTime,
		MinCPUs:           minCPUs,
		MinNodes:          minNodes,
		UsageFactor:       qos.UsageFactor,
		UsageThreshold:    qos.UsageThreshold,
		Flags:             qos.Flags,
		AllowedAccounts:   qos.AllowedAccounts,
		DeniedAccounts:    []string{}, // Not available in types.QoS
		AllowedUsers:      qos.AllowedUsers,
		DeniedUsers:       []string{}, // Not available in types.QoS,
		Metadata:          nil,
	}
}

// adapterInfoManager provides basic info operations
type adapterInfoManager struct {
	version string
}

func (m *adapterInfoManager) Ping(ctx context.Context) error {
	// Basic ping - always succeeds if we get here
	return nil
}

func (m *adapterInfoManager) Get(ctx context.Context) (*interfaces.ClusterInfo, error) {
	// Not implemented
	return nil, fmt.Errorf("cluster info not implemented in adapter")
}

func (m *adapterInfoManager) Stats(ctx context.Context) (*interfaces.ClusterStats, error) {
	// Not implemented
	return nil, fmt.Errorf("cluster stats not implemented in adapter")
}

func (m *adapterInfoManager) Version(ctx context.Context) (*interfaces.APIVersion, error) {
	return &interfaces.APIVersion{
		Version:     m.version,
		Release:     "stable",
		Description: "SLURM REST API",
		Deprecated:  false,
	}, nil
}

// Other manager implementations...

type adapterQoSManager struct {
	adapter common.QoSAdapter
}

func (m *adapterQoSManager) List(ctx context.Context, opts *interfaces.ListQoSOptions) (*interfaces.QoSList, error) {
	// Convert options
	adapterOpts := &types.QoSListOptions{}
	if opts != nil {
		adapterOpts.Limit = opts.Limit
		adapterOpts.Offset = opts.Offset
	}

	// Call adapter
	result, err := m.adapter.List(ctx, adapterOpts)
	if err != nil {
		return nil, err
	}

	// Convert result
	qosList := &interfaces.QoSList{
		QoS:   make([]interfaces.QoS, 0, len(result.QoS)),
		Total: result.Total,
	}

	for _, qos := range result.QoS {
		qosList.QoS = append(qosList.QoS, convertQoSToInterface(qos))
	}

	return qosList, nil
}

func (m *adapterQoSManager) Get(ctx context.Context, qosName string) (*interfaces.QoS, error) {
	qos, err := m.adapter.Get(ctx, qosName)
	if err != nil {
		return nil, err
	}
	result := convertQoSToInterface(*qos)
	return &result, nil
}

func (m *adapterQoSManager) Create(ctx context.Context, qos *interfaces.QoSCreate) (*interfaces.QoSCreateResponse, error) {
	// Convert create request
	adapterCreate := &types.QoSCreate{
		Name:        qos.Name,
		Description: qos.Description,
		Priority:    qos.Priority,
		// Add other fields as needed
	}

	// Call adapter
	resp, err := m.adapter.Create(ctx, adapterCreate)
	if err != nil {
		return nil, err
	}

	return &interfaces.QoSCreateResponse{
		QoSName: resp.QoSName,
	}, nil
}

func (m *adapterQoSManager) Update(ctx context.Context, qosName string, update *interfaces.QoSUpdate) error {
	// Convert update request
	adapterUpdate := &types.QoSUpdate{}
	if update.Description != nil {
		adapterUpdate.Description = update.Description
	}
	if update.Priority != nil {
		adapterUpdate.Priority = update.Priority
	}
	// Add other fields as needed

	return m.adapter.Update(ctx, qosName, adapterUpdate)
}

func (m *adapterQoSManager) Delete(ctx context.Context, qosName string) error {
	return m.adapter.Delete(ctx, qosName)
}

type adapterAccountManager struct {
	adapter common.AccountAdapter
}

func (m *adapterAccountManager) List(ctx context.Context, opts *interfaces.ListAccountsOptions) (*interfaces.AccountList, error) {
	// Convert options
	adapterOpts := &types.AccountListOptions{}
	if opts != nil {
		adapterOpts.Limit = opts.Limit
		adapterOpts.Offset = opts.Offset
		// Note: Some fields may not have direct mappings
	}

	// Call adapter
	result, err := m.adapter.List(ctx, adapterOpts)
	if err != nil {
		return nil, err
	}

	// Convert result
	accountList := &interfaces.AccountList{
		Accounts: make([]interfaces.Account, 0, len(result.Accounts)),
		Total:    result.Total,
	}

	for _, account := range result.Accounts {
		accountList.Accounts = append(accountList.Accounts, convertAccountToInterface(account))
	}

	return accountList, nil
}

func (m *adapterAccountManager) Get(ctx context.Context, accountName string) (*interfaces.Account, error) {
	account, err := m.adapter.Get(ctx, accountName)
	if err != nil {
		return nil, err
	}
	result := convertAccountToInterface(*account)
	return &result, nil
}

func (m *adapterAccountManager) Create(ctx context.Context, account *interfaces.AccountCreate) (*interfaces.AccountCreateResponse, error) {
	// Convert create request
	adapterCreate := &types.AccountCreate{
		Name:        account.Name,
		Description: account.Description,
		Organization: account.Organization,
		// Add other fields as needed
	}

	// Call adapter
	resp, err := m.adapter.Create(ctx, adapterCreate)
	if err != nil {
		return nil, err
	}

	return &interfaces.AccountCreateResponse{
		AccountName: resp.AccountName,
	}, nil
}

func (m *adapterAccountManager) Update(ctx context.Context, accountName string, update *interfaces.AccountUpdate) error {
	// Convert update request
	adapterUpdate := &types.AccountUpdate{
		Description: update.Description,
		Organization: update.Organization,
		// Add other fields as needed
	}

	return m.adapter.Update(ctx, accountName, adapterUpdate)
}

func (m *adapterAccountManager) Delete(ctx context.Context, accountName string) error {
	return m.adapter.Delete(ctx, accountName)
}

func (m *adapterAccountManager) GetAccountHierarchy(ctx context.Context, rootAccount string) (*interfaces.AccountHierarchy, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *adapterAccountManager) GetParentAccounts(ctx context.Context, accountName string) ([]*interfaces.Account, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *adapterAccountManager) GetChildAccounts(ctx context.Context, accountName string, depth int) ([]*interfaces.Account, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *adapterAccountManager) GetAccountQuotas(ctx context.Context, accountName string) (*interfaces.AccountQuota, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *adapterAccountManager) GetAccountQuotaUsage(ctx context.Context, accountName string, timeframe string) (*interfaces.AccountUsage, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *adapterAccountManager) GetAccountUsers(ctx context.Context, accountName string, opts *interfaces.ListAccountUsersOptions) ([]*interfaces.UserAccountAssociation, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *adapterAccountManager) ValidateUserAccess(ctx context.Context, userName, accountName string) (*interfaces.UserAccessValidation, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *adapterAccountManager) GetAccountUsersWithPermissions(ctx context.Context, accountName string, permissions []string) ([]*interfaces.UserAccountAssociation, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *adapterAccountManager) GetAccountFairShare(ctx context.Context, accountName string) (*interfaces.AccountFairShare, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *adapterAccountManager) GetFairShareHierarchy(ctx context.Context, rootAccount string) (*interfaces.FairShareHierarchy, error) {
	return nil, fmt.Errorf("not implemented")
}

type adapterUserManager struct {
	adapter common.UserAdapter
}

func (m *adapterUserManager) List(ctx context.Context, opts *interfaces.ListUsersOptions) (*interfaces.UserList, error) {
	// Convert options
	adapterOpts := &types.UserListOptions{}
	if opts != nil {
		adapterOpts.Limit = opts.Limit
		adapterOpts.Offset = opts.Offset
		// Note: Some fields may not have direct mappings
	}

	// Call adapter
	result, err := m.adapter.List(ctx, adapterOpts)
	if err != nil {
		return nil, err
	}

	// Convert result
	userList := &interfaces.UserList{
		Users: make([]interfaces.User, 0, len(result.Users)),
		Total: result.Total,
	}

	for _, user := range result.Users {
		userList.Users = append(userList.Users, convertUserToInterface(user))
	}

	return userList, nil
}

func (m *adapterUserManager) Get(ctx context.Context, userName string) (*interfaces.User, error) {
	user, err := m.adapter.Get(ctx, userName)
	if err != nil {
		return nil, err
	}
	result := convertUserToInterface(*user)
	return &result, nil
}

func (m *adapterUserManager) GetUserAccounts(ctx context.Context, userName string) ([]*interfaces.UserAccount, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *adapterUserManager) GetUserQuotas(ctx context.Context, userName string) (*interfaces.UserQuota, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *adapterUserManager) GetUserDefaultAccount(ctx context.Context, userName string) (*interfaces.Account, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *adapterUserManager) GetUserFairShare(ctx context.Context, userName string) (*interfaces.UserFairShare, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *adapterUserManager) CalculateJobPriority(ctx context.Context, userName string, jobSubmission *interfaces.JobSubmission) (*interfaces.JobPriorityInfo, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *adapterUserManager) ValidateUserAccountAccess(ctx context.Context, userName, accountName string) (*interfaces.UserAccessValidation, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *adapterUserManager) GetUserAccountAssociations(ctx context.Context, userName string, opts *interfaces.ListUserAccountAssociationsOptions) ([]*interfaces.UserAccountAssociation, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *adapterUserManager) GetBulkUserAccounts(ctx context.Context, userNames []string) (map[string][]*interfaces.UserAccount, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *adapterUserManager) GetBulkAccountUsers(ctx context.Context, accountNames []string) (map[string][]*interfaces.UserAccountAssociation, error) {
	return nil, fmt.Errorf("not implemented")
}

type adapterReservationManager struct {
	adapter common.ReservationAdapter
}

func (m *adapterReservationManager) List(ctx context.Context, opts *interfaces.ListReservationsOptions) (*interfaces.ReservationList, error) {
	// Convert options
	adapterOpts := &types.ReservationListOptions{}
	if opts != nil {
		adapterOpts.Limit = opts.Limit
		adapterOpts.Offset = opts.Offset
	}

	// Call adapter
	result, err := m.adapter.List(ctx, adapterOpts)
	if err != nil {
		return nil, err
	}

	// Convert result
	reservationList := &interfaces.ReservationList{
		Reservations: make([]interfaces.Reservation, 0, len(result.Reservations)),
		Total:        result.Total,
	}

	for _, reservation := range result.Reservations {
		reservationList.Reservations = append(reservationList.Reservations, convertReservationToInterface(reservation))
	}

	return reservationList, nil
}

func (m *adapterReservationManager) Get(ctx context.Context, reservationName string) (*interfaces.Reservation, error) {
	// Call adapter
	result, err := m.adapter.Get(ctx, reservationName)
	if err != nil {
		return nil, err
	}

	// Convert result
	reservation := convertReservationToInterface(*result)
	return &reservation, nil
}

func (m *adapterReservationManager) Create(ctx context.Context, reservation *interfaces.ReservationCreate) (*interfaces.ReservationCreateResponse, error) {
	// Convert to adapter type
	adapterReservation := &types.ReservationCreate{
		Name:          reservation.Name,
		StartTime:     reservation.StartTime,
		EndTime:       &reservation.EndTime,
		Duration:      int32(reservation.Duration),
		Users:         reservation.Users,
		Accounts:      reservation.Accounts,
		NodeList:      strings.Join(reservation.Nodes, ","),
		NodeCount:     int32(reservation.NodeCount),
		CoreCount:     int32(reservation.CoreCount),
		PartitionName: reservation.PartitionName,
		Features:      reservation.Features,
		BurstBuffer:   reservation.BurstBuffer,
		Comment:       "", // Not available in interfaces.ReservationCreate
	}

	// Convert flags from []string to []ReservationFlag
	if len(reservation.Flags) > 0 {
		adapterReservation.Flags = make([]types.ReservationFlag, len(reservation.Flags))
		for i, flag := range reservation.Flags {
			adapterReservation.Flags[i] = types.ReservationFlag(flag)
		}
	}

	// Convert licenses from map[string]int to map[string]int32
	if len(reservation.Licenses) > 0 {
		adapterReservation.Licenses = make(map[string]int32)
		for k, v := range reservation.Licenses {
			adapterReservation.Licenses[k] = int32(v)
		}
	}

	// Call adapter
	result, err := m.adapter.Create(ctx, adapterReservation)
	if err != nil {
		return nil, err
	}

	// Convert result
	return &interfaces.ReservationCreateResponse{
		ReservationName: result.ReservationName,
	}, nil
}

func (m *adapterReservationManager) Update(ctx context.Context, reservationName string, update *interfaces.ReservationUpdate) error {
	// Convert to adapter type
	adapterUpdate := &types.ReservationUpdate{}
	if update != nil {
		if update.Users != nil {
			adapterUpdate.Users = update.Users
		}
		if update.Accounts != nil {
			adapterUpdate.Accounts = update.Accounts
		}
		if update.Nodes != nil && len(update.Nodes) > 0 {
			nodeList := strings.Join(update.Nodes, ",")
			adapterUpdate.NodeList = &nodeList
		}
		if update.NodeCount != nil {
			nodeCount := int32(*update.NodeCount)
			adapterUpdate.NodeCount = &nodeCount
		}
		if update.StartTime != nil {
			adapterUpdate.StartTime = update.StartTime
		}
		if update.EndTime != nil {
			adapterUpdate.EndTime = update.EndTime
		}
		if update.Duration != nil {
			duration := int32(*update.Duration)
			adapterUpdate.Duration = &duration
		}
		if update.Flags != nil && len(update.Flags) > 0 {
			adapterUpdate.Flags = make([]types.ReservationFlag, len(update.Flags))
			for i, flag := range update.Flags {
				adapterUpdate.Flags[i] = types.ReservationFlag(flag)
			}
		}
		if update.Features != nil {
			adapterUpdate.Features = update.Features
		}
	}

	return m.adapter.Update(ctx, reservationName, adapterUpdate)
}

func (m *adapterReservationManager) Delete(ctx context.Context, reservationName string) error {
	return m.adapter.Delete(ctx, reservationName)
}

type adapterAssociationManager struct {
	adapter common.AssociationAdapter
}

func (m *adapterAssociationManager) List(ctx context.Context, opts *interfaces.ListAssociationsOptions) (*interfaces.AssociationList, error) {
	// Convert options
	adapterOpts := &types.AssociationListOptions{}
	if opts != nil {
		adapterOpts.Limit = opts.Limit
		adapterOpts.Offset = opts.Offset
		// Note: Some filter fields may not map directly
	}

	// Call adapter
	result, err := m.adapter.List(ctx, adapterOpts)
	if err != nil {
		return nil, err
	}

	// Convert result
	associationList := &interfaces.AssociationList{
		Associations: make([]*interfaces.Association, 0, len(result.Associations)),
		Total:        result.Total,
	}

	for _, association := range result.Associations {
		converted := convertAssociationToInterface(association)
		associationList.Associations = append(associationList.Associations, &converted)
	}

	return associationList, nil
}

func (m *adapterAssociationManager) Get(ctx context.Context, opts *interfaces.GetAssociationOptions) (*interfaces.Association, error) {
	// Build association ID from options
	if opts == nil || (opts.User == "" && opts.Account == "") {
		return nil, fmt.Errorf("user or account must be specified")
	}

	// Create a composite ID or use the first matching association
	associationID := ""
	if opts.User != "" && opts.Account != "" {
		associationID = fmt.Sprintf("%s:%s", opts.User, opts.Account)
	} else if opts.User != "" {
		associationID = opts.User
	} else {
		associationID = opts.Account
	}

	// Call adapter
	result, err := m.adapter.Get(ctx, associationID)
	if err != nil {
		return nil, err
	}

	// Convert result
	association := convertAssociationToInterface(*result)
	return &association, nil
}

func (m *adapterAssociationManager) Create(ctx context.Context, associations []*interfaces.AssociationCreate) (*interfaces.AssociationCreateResponse, error) {
	// For now, handle single association creation
	if len(associations) == 0 {
		return nil, fmt.Errorf("no associations provided")
	}

	// Convert first association (adapter interface expects single association)
	assoc := associations[0]
	adapterAssoc := &types.AssociationCreate{
		AccountName:   assoc.Account,
		Cluster:       assoc.Cluster,
		UserName:      assoc.User,
		Partition:     assoc.Partition,
		ParentAccount: assoc.ParentAccount,
		IsDefault:     assoc.IsDefault,
		DefaultQoS:    assoc.DefaultQoS,
		QoSList:       assoc.QoSList,
		Comment:       assoc.Comment,
	}

	// Convert pointer values to non-pointer values for the adapter
	if assoc.SharesRaw != nil {
		adapterAssoc.SharesRaw = int32(*assoc.SharesRaw)
	}
	if assoc.Priority != nil {
		adapterAssoc.Priority = int32(*assoc.Priority)
	}
	if assoc.MaxJobs != nil {
		adapterAssoc.MaxJobs = int32(*assoc.MaxJobs)
	}
	if assoc.MaxJobsAccrue != nil {
		adapterAssoc.MaxJobsAccrue = int32(*assoc.MaxJobsAccrue)
	}
	if assoc.MaxSubmitJobs != nil {
		adapterAssoc.MaxSubmitJobs = int32(*assoc.MaxSubmitJobs)
	}
	if assoc.MaxWallDuration != nil {
		adapterAssoc.MaxWallTime = int32(*assoc.MaxWallDuration)
	}
	if assoc.GrpJobs != nil {
		adapterAssoc.GrpJobs = int32(*assoc.GrpJobs)
	}
	if assoc.GrpJobsAccrue != nil {
		adapterAssoc.GrpJobsAccrue = int32(*assoc.GrpJobsAccrue)
	}
	if assoc.GrpSubmitJobs != nil {
		adapterAssoc.GrpSubmitJobs = int32(*assoc.GrpSubmitJobs)
	}
	if assoc.GrpWall != nil {
		adapterAssoc.GrpWallTime = int32(*assoc.GrpWall)
	}

	// Call adapter
	result, err := m.adapter.Create(ctx, adapterAssoc)
	if err != nil {
		return nil, err
	}

	// Convert result - create a basic association from the response
	association := &interfaces.Association{
		User:    result.UserName,
		Account: result.AccountName,
		Cluster: result.Cluster,
	}
	
	return &interfaces.AssociationCreateResponse{
		Associations: []*interfaces.Association{association},
		Created:      1,
		Updated:      0,
		Errors:       nil,
		Warnings:     nil,
	}, nil
}

func (m *adapterAssociationManager) Update(ctx context.Context, associations []*interfaces.AssociationUpdate) error {
	// For now, handle single association update
	if len(associations) == 0 {
		return fmt.Errorf("no associations provided")
	}

	// Process first association
	assoc := associations[0]
	associationID := ""
	if assoc.User != "" && assoc.Account != "" {
		associationID = fmt.Sprintf("%s:%s", assoc.User, assoc.Account)
	} else if assoc.User != "" {
		associationID = assoc.User
	} else if assoc.Account != "" {
		associationID = assoc.Account
	} else {
		return fmt.Errorf("user or account must be specified")
	}

	// Convert to adapter type
	adapterUpdate := &types.AssociationUpdate{}
	if assoc.IsDefault != nil {
		adapterUpdate.IsDefault = assoc.IsDefault
	}
	if assoc.Comment != nil {
		adapterUpdate.Comment = assoc.Comment
	}
	if assoc.DefaultQoS != nil {
		adapterUpdate.DefaultQoS = assoc.DefaultQoS
	}
	if assoc.QoSList != nil {
		adapterUpdate.QoSList = assoc.QoSList
	}

	// Convert int pointers to int32 pointers
	if assoc.SharesRaw != nil {
		sharesRaw := int32(*assoc.SharesRaw)
		adapterUpdate.SharesRaw = &sharesRaw
	}
	if assoc.Priority != nil {
		priority := int32(*assoc.Priority)
		adapterUpdate.Priority = &priority
	}
	if assoc.MaxJobs != nil {
		maxJobs := int32(*assoc.MaxJobs)
		adapterUpdate.MaxJobs = &maxJobs
	}
	if assoc.MaxJobsAccrue != nil {
		maxJobsAccrue := int32(*assoc.MaxJobsAccrue)
		adapterUpdate.MaxJobsAccrue = &maxJobsAccrue
	}
	if assoc.MaxSubmitJobs != nil {
		maxSubmitJobs := int32(*assoc.MaxSubmitJobs)
		adapterUpdate.MaxSubmitJobs = &maxSubmitJobs
	}
	if assoc.MaxWallDuration != nil {
		maxWallTime := int32(*assoc.MaxWallDuration)
		adapterUpdate.MaxWallTime = &maxWallTime
	}
	if assoc.GrpJobs != nil {
		grpJobs := int32(*assoc.GrpJobs)
		adapterUpdate.GrpJobs = &grpJobs
	}
	if assoc.GrpJobsAccrue != nil {
		grpJobsAccrue := int32(*assoc.GrpJobsAccrue)
		adapterUpdate.GrpJobsAccrue = &grpJobsAccrue
	}
	if assoc.GrpSubmitJobs != nil {
		grpSubmitJobs := int32(*assoc.GrpSubmitJobs)
		adapterUpdate.GrpSubmitJobs = &grpSubmitJobs
	}
	if assoc.GrpWall != nil {
		grpWallTime := int32(*assoc.GrpWall)
		adapterUpdate.GrpWallTime = &grpWallTime
	}

	return m.adapter.Update(ctx, associationID, adapterUpdate)
}

func (m *adapterAssociationManager) Delete(ctx context.Context, opts *interfaces.DeleteAssociationOptions) error {
	if opts == nil || (opts.User == "" && opts.Account == "") {
		return fmt.Errorf("user or account must be specified")
	}

	// Build association ID from options
	associationID := ""
	if opts.User != "" && opts.Account != "" {
		associationID = fmt.Sprintf("%s:%s", opts.User, opts.Account)
	} else if opts.User != "" {
		associationID = opts.User
	} else {
		associationID = opts.Account
	}

	return m.adapter.Delete(ctx, associationID)
}

func (m *adapterAssociationManager) BulkDelete(ctx context.Context, opts *interfaces.BulkDeleteOptions) (*interfaces.BulkDeleteResponse, error) {
	// Not supported in base adapter interface
	return nil, fmt.Errorf("bulk delete not supported in adapter implementation")
}

func (m *adapterAssociationManager) GetUserAssociations(ctx context.Context, userName string) ([]*interfaces.Association, error) {
	// List associations and filter by user
	opts := &types.AssociationListOptions{
		Limit: 1000, // Get all associations for the user
	}

	result, err := m.adapter.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	// Filter associations for the specified user
	userAssociations := make([]*interfaces.Association, 0)
	for _, assoc := range result.Associations {
		if assoc.UserName == userName {
			converted := convertAssociationToInterface(assoc)
			userAssociations = append(userAssociations, &converted)
		}
	}

	return userAssociations, nil
}

func (m *adapterAssociationManager) GetAccountAssociations(ctx context.Context, accountName string) ([]*interfaces.Association, error) {
	// List associations and filter by account
	opts := &types.AssociationListOptions{
		Limit: 1000, // Get all associations for the account
	}

	result, err := m.adapter.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	// Filter associations for the specified account
	accountAssociations := make([]*interfaces.Association, 0)
	for _, assoc := range result.Associations {
		if assoc.AccountName == accountName {
			converted := convertAssociationToInterface(assoc)
			accountAssociations = append(accountAssociations, &converted)
		}
	}

	return accountAssociations, nil
}

func (m *adapterAssociationManager) ValidateAssociation(ctx context.Context, user, account, cluster string) (bool, error) {
	// Try to get the specific association
	opts := &interfaces.GetAssociationOptions{
		User:    user,
		Account: account,
		Cluster: cluster,
	}

	_, err := m.Get(ctx, opts)
	if err != nil {
		// Association doesn't exist or error occurred
		return false, nil
	}

	// Association exists
	return true, nil
}

// Helper function to convert types.Reservation to interfaces.Reservation
func convertReservationToInterface(reservation types.Reservation) interfaces.Reservation {
	// Convert flags from []ReservationFlag to []string
	flags := make([]string, len(reservation.Flags))
	for i, flag := range reservation.Flags {
		flags[i] = string(flag)
	}

	// Convert licenses from map[string]int32 to map[string]int
	licenses := make(map[string]int)
	for k, v := range reservation.Licenses {
		licenses[k] = int(v)
	}

	// Convert node list to array if needed
	nodes := []string{}
	if reservation.NodeList != "" {
		nodes = strings.Split(reservation.NodeList, ",")
	}

	return interfaces.Reservation{
		Name:          reservation.Name,
		State:         string(reservation.State),
		StartTime:     reservation.StartTime,
		EndTime:       reservation.EndTime,
		Duration:      int(reservation.Duration),
		Users:         reservation.Users,
		Accounts:      reservation.Accounts,
		Nodes:         nodes,
		NodeCount:     int(reservation.NodeCount),
		CoreCount:     int(reservation.CoreCount),
		PartitionName: reservation.PartitionName,
		Flags:         flags,
		Features:      reservation.Features,
		Licenses:      licenses,
		BurstBuffer:   reservation.BurstBuffer,
		Metadata:      nil,
	}
}

// Helper function to convert types.Association to interfaces.Association
func convertAssociationToInterface(association types.Association) interfaces.Association {
	// Convert pointers for limits
	var maxJobs, maxSubmitJobs *int
	if association.MaxJobs > 0 {
		jobs := int(association.MaxJobs)
		maxJobs = &jobs
	}
	if association.MaxSubmitJobs > 0 {
		submitJobs := int(association.MaxSubmitJobs)
		maxSubmitJobs = &submitJobs
	}

	var maxWallDuration *int
	if association.MaxWallTime > 0 {
		wallTime := int(association.MaxWallTime)
		maxWallDuration = &wallTime
	}

	var grpJobs, grpSubmitJobs, grpWall *int
	if association.GrpJobs > 0 {
		jobs := int(association.GrpJobs)
		grpJobs = &jobs
	}
	if association.GrpSubmitJobs > 0 {
		submitJobs := int(association.GrpSubmitJobs)
		grpSubmitJobs = &submitJobs
	}
	if association.GrpWallTime > 0 {
		wallTime := int(association.GrpWallTime)
		grpWall = &wallTime
	}

	return interfaces.Association{
		ID:              0, // Not available in types.Association
		User:            association.UserName,
		Account:         association.AccountName,
		Cluster:         association.Cluster,
		Partition:       association.Partition,
		ParentAccount:   association.ParentAccount,
		IsDefault:       false, // Not available in types.Association
		Comment:         association.Comment,
		SharesRaw:       int(association.SharesRaw),
		Priority:        uint32(association.Priority),
		MaxJobs:         maxJobs,
		MaxJobsAccrue:   nil, // Not available in types.Association
		MaxSubmitJobs:   maxSubmitJobs,
		MaxWallDuration: maxWallDuration,
		GrpJobs:         grpJobs,
		GrpJobsAccrue:   nil, // Not available in types.Association
		GrpSubmitJobs:   grpSubmitJobs,
		GrpWall:         grpWall,
		MaxTRESPerJob:   make(map[string]string), // Would need proper conversion
		MaxTRESMins:     make(map[string]string), // Would need proper conversion
		GrpTRES:         make(map[string]string), // Would need proper conversion
		GrpTRESMins:     make(map[string]string), // Would need proper conversion
		GrpTRESRunMins:  make(map[string]string), // Would need proper conversion
		DefaultQoS:      association.DefaultQoS,
		QoSList:         association.QoSList,
		Flags:           []string{}, // Not available in types.Association
	}
}
