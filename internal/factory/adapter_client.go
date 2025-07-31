package factory

import (
	"context"
	"fmt"
	"strconv"

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
	// For adapters, we'll use polling-based watch functionality
	// Since the adapters don't have native watch support, we create a poller
	// that periodically calls List to detect changes
	
	// Note: This is a simple implementation. For production use, you might want to:
	// 1. Make the poll interval configurable
	// 2. Add proper state tracking to detect changes
	// 3. Implement more sophisticated change detection logic
	
	return nil, fmt.Errorf("watch functionality requires polling implementation - not yet implemented for adapters")
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
	// For adapters, we'll use polling-based watch functionality
	// Since the adapters don't have native watch support, we create a poller
	// that periodically calls List to detect changes
	
	// Note: This is a simple implementation. For production use, you might want to:
	// 1. Make the poll interval configurable
	// 2. Add proper state tracking to detect changes
	// 3. Implement more sophisticated change detection logic
	
	return nil, fmt.Errorf("watch functionality requires polling implementation - not yet implemented for adapters")
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
	// Implementation would go here
	return nil, fmt.Errorf("not implemented")
}

func (m *adapterQoSManager) Get(ctx context.Context, qosName string) (*interfaces.QoS, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *adapterQoSManager) Create(ctx context.Context, qos *interfaces.QoSCreate) (*interfaces.QoSCreateResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *adapterQoSManager) Update(ctx context.Context, qosName string, update *interfaces.QoSUpdate) error {
	return fmt.Errorf("not implemented")
}

func (m *adapterQoSManager) Delete(ctx context.Context, qosName string) error {
	return fmt.Errorf("not implemented")
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
	return nil, fmt.Errorf("not implemented")
}

func (m *adapterUserManager) Get(ctx context.Context, userName string) (*interfaces.User, error) {
	return nil, fmt.Errorf("not implemented")
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
	return nil, fmt.Errorf("not implemented")
}

func (m *adapterReservationManager) Get(ctx context.Context, reservationName string) (*interfaces.Reservation, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *adapterReservationManager) Create(ctx context.Context, reservation *interfaces.ReservationCreate) (*interfaces.ReservationCreateResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *adapterReservationManager) Update(ctx context.Context, reservationName string, update *interfaces.ReservationUpdate) error {
	return fmt.Errorf("not implemented")
}

func (m *adapterReservationManager) Delete(ctx context.Context, reservationName string) error {
	return fmt.Errorf("not implemented")
}

type adapterAssociationManager struct {
	adapter common.AssociationAdapter
}

func (m *adapterAssociationManager) List(ctx context.Context, opts *interfaces.ListAssociationsOptions) (*interfaces.AssociationList, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *adapterAssociationManager) Get(ctx context.Context, opts *interfaces.GetAssociationOptions) (*interfaces.Association, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *adapterAssociationManager) Create(ctx context.Context, associations []*interfaces.AssociationCreate) (*interfaces.AssociationCreateResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *adapterAssociationManager) Update(ctx context.Context, associations []*interfaces.AssociationUpdate) error {
	return fmt.Errorf("not implemented")
}

func (m *adapterAssociationManager) Delete(ctx context.Context, opts *interfaces.DeleteAssociationOptions) error {
	return fmt.Errorf("not implemented")
}

func (m *adapterAssociationManager) BulkDelete(ctx context.Context, opts *interfaces.BulkDeleteOptions) (*interfaces.BulkDeleteResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *adapterAssociationManager) GetUserAssociations(ctx context.Context, userName string) ([]*interfaces.Association, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *adapterAssociationManager) GetAccountAssociations(ctx context.Context, accountName string) ([]*interfaces.Association, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *adapterAssociationManager) ValidateAssociation(ctx context.Context, user, account, cluster string) (bool, error) {
	return false, fmt.Errorf("not implemented")
}