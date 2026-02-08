// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package factory

import (
	"context"
	"fmt"
	"strconv"
	"time"

	types "github.com/jontk/slurm-client/api"
	"github.com/jontk/slurm-client/internal/adapters/common"
	v040adapter "github.com/jontk/slurm-client/internal/adapters/v0_0_40"
	v041adapter "github.com/jontk/slurm-client/internal/adapters/v0_0_41"
	v042adapter "github.com/jontk/slurm-client/internal/adapters/v0_0_42"
	v043adapter "github.com/jontk/slurm-client/internal/adapters/v0_0_43"
	v044adapter "github.com/jontk/slurm-client/internal/adapters/v0_0_44"
	v040api "github.com/jontk/slurm-client/internal/openapi/v0_0_40"
	v041api "github.com/jontk/slurm-client/internal/openapi/v0_0_41"
	v042api "github.com/jontk/slurm-client/internal/openapi/v0_0_42"
	v043api "github.com/jontk/slurm-client/internal/openapi/v0_0_43"
	v044api "github.com/jontk/slurm-client/internal/openapi/v0_0_44"
	"github.com/jontk/slurm-client/pkg/errors"
	"github.com/jontk/slurm-client/pkg/pool"
)

// AdapterClient wraps a version-specific adapter to implement the SlurmClient interface
type AdapterClient struct {
	adapter common.VersionAdapter
	version string
	pool    *pool.HTTPClientPool // optional connection pool for cleanup
}

// NewAdapterClient creates a new adapter-based client for the specified version
func NewAdapterClient(version string, config *types.ClientConfig) (SlurmClient, error) {
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
			return nil, fmt.Errorf("failed to create %s client: %w", version, err)
		}
		adapter := v043adapter.NewAdapter(client)
		return &AdapterClient{
			adapter: adapter,
			version: version,
		}, nil

	case "v0.0.44":
		client, err := v044api.NewClientWithResponses(config.BaseURL, v044api.WithHTTPClient(config.HTTPClient))
		if err != nil {
			return nil, fmt.Errorf("failed to create v0.0.44 client: %w", err)
		}
		adapter := v044adapter.NewAdapter(client)
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

// Capabilities returns the features supported by this client version
func (c *AdapterClient) Capabilities() types.ClientCapabilities {
	return c.adapter.GetCapabilities()
}

// Jobs returns the JobManager
func (c *AdapterClient) Jobs() types.JobManager {
	return &adapterJobManager{adapter: c.adapter.GetJobManager()}
}

// Nodes returns the NodeManager
func (c *AdapterClient) Nodes() types.NodeManager {
	return &adapterNodeManager{adapter: c.adapter.GetNodeManager()}
}

// Partitions returns the PartitionManager
func (c *AdapterClient) Partitions() types.PartitionManager {
	return &adapterPartitionManager{adapter: c.adapter.GetPartitionManager()}
}

// Info returns the InfoManager
func (c *AdapterClient) Info() types.InfoManager {
	return &adapterInfoManager{
		adapter: c.adapter.GetInfoManager(),
		version: c.version,
	}
}

// Reservations returns the ReservationManager
func (c *AdapterClient) Reservations() types.ReservationManager {
	return &adapterReservationManager{adapter: c.adapter.GetReservationManager()}
}

// QoS returns the QoSManager
func (c *AdapterClient) QoS() types.QoSManager {
	return &adapterQoSManager{adapter: c.adapter.GetQoSManager()}
}

// Accounts returns the AccountManager
func (c *AdapterClient) Accounts() types.AccountManager {
	return &adapterAccountManager{
		adapter:            c.adapter.GetAccountManager(),
		associationAdapter: c.adapter.GetAssociationManager(),
	}
}

// Users returns the UserManager
func (c *AdapterClient) Users() types.UserManager {
	return &adapterUserManager{
		adapter:            c.adapter.GetUserManager(),
		accountAdapter:     c.adapter.GetAccountManager(),
		associationAdapter: c.adapter.GetAssociationManager(),
	}
}

// Clusters returns the ClusterManager
func (c *AdapterClient) Clusters() types.ClusterManager {
	return &adapterClusterManager{adapter: c.adapter.GetClusterManager()}
}

// Associations returns the AssociationManager
func (c *AdapterClient) Associations() types.AssociationManager {
	return &adapterAssociationManager{adapter: c.adapter.GetAssociationManager()}
}

// WCKeys returns the WCKeyManager
func (c *AdapterClient) WCKeys() types.WCKeyManager {
	return &adapterWCKeyManager{adapter: c.adapter.GetWCKeyManager()}
}

// Analytics returns the AnalyticsManager (not implemented in current release)
func (c *AdapterClient) Analytics() types.AnalyticsManager {
	// Analytics is not yet implemented - this is a value-added feature
	// that will compute insights from API data in future releases
	return nil
}

// Close closes the client and releases any resources
func (c *AdapterClient) Close() error {
	if c.pool != nil {
		return c.pool.Close()
	}
	return nil
}

// SetPool sets the connection pool for resource cleanup on Close
func (c *AdapterClient) SetPool(p *pool.HTTPClientPool) {
	c.pool = p
}

// === Standalone Operations ===

// GetLicenses retrieves license information
func (c *AdapterClient) GetLicenses(ctx context.Context) (*types.LicenseList, error) {
	standaloneManager := c.adapter.GetStandaloneManager()
	if standaloneManager == nil {
		return nil, fmt.Errorf("standalone operations not supported for version %s", c.version)
	}
	result, err := standaloneManager.GetLicenses(ctx)
	if err != nil {
		return nil, err
	}
	return convertLicenseListToInterface(result), nil
}

// GetShares retrieves fairshare information with optional filtering
func (c *AdapterClient) GetShares(ctx context.Context, opts *types.GetSharesOptions) (*types.SharesList, error) {
	standaloneManager := c.adapter.GetStandaloneManager()
	if standaloneManager == nil {
		return nil, fmt.Errorf("standalone operations not supported for version %s", c.version)
	}
	typesOpts := convertGetSharesOptionsToTypes(opts)
	result, err := standaloneManager.GetShares(ctx, typesOpts)
	if err != nil {
		return nil, err
	}
	return convertSharesListToInterface(result), nil
}

// GetConfig retrieves SLURM configuration
func (c *AdapterClient) GetConfig(ctx context.Context) (*types.Config, error) {
	standaloneManager := c.adapter.GetStandaloneManager()
	if standaloneManager == nil {
		return nil, fmt.Errorf("standalone operations not supported for version %s", c.version)
	}
	result, err := standaloneManager.GetConfig(ctx)
	if err != nil {
		return nil, err
	}
	return convertConfigToInterface(result), nil
}

// GetDiagnostics retrieves SLURM diagnostics information
func (c *AdapterClient) GetDiagnostics(ctx context.Context) (*types.Diagnostics, error) {
	standaloneManager := c.adapter.GetStandaloneManager()
	if standaloneManager == nil {
		return nil, fmt.Errorf("standalone operations not supported for version %s", c.version)
	}
	result, err := standaloneManager.GetDiagnostics(ctx)
	if err != nil {
		return nil, err
	}
	return convertDiagnosticsToInterface(result), nil
}

// GetDBDiagnostics retrieves SLURM database diagnostics information
func (c *AdapterClient) GetDBDiagnostics(ctx context.Context) (*types.Diagnostics, error) {
	standaloneManager := c.adapter.GetStandaloneManager()
	if standaloneManager == nil {
		return nil, fmt.Errorf("standalone operations not supported for version %s", c.version)
	}
	result, err := standaloneManager.GetDBDiagnostics(ctx)
	if err != nil {
		return nil, err
	}
	return convertDiagnosticsToInterface(result), nil
}

// GetInstance retrieves a specific database instance
func (c *AdapterClient) GetInstance(ctx context.Context, opts *types.GetInstanceOptions) (*types.Instance, error) {
	standaloneManager := c.adapter.GetStandaloneManager()
	if standaloneManager == nil {
		return nil, fmt.Errorf("standalone operations not supported for version %s", c.version)
	}
	typesOpts := convertGetInstanceOptionsToTypes(opts)
	result, err := standaloneManager.GetInstance(ctx, typesOpts)
	if err != nil {
		return nil, err
	}
	return convertInstanceToInterface(result), nil
}

// GetInstances retrieves multiple database instances with filtering
func (c *AdapterClient) GetInstances(ctx context.Context, opts *types.GetInstancesOptions) (*types.InstanceList, error) {
	standaloneManager := c.adapter.GetStandaloneManager()
	if standaloneManager == nil {
		return nil, fmt.Errorf("standalone operations not supported for version %s", c.version)
	}
	typesOpts := convertGetInstancesOptionsToTypes(opts)
	result, err := standaloneManager.GetInstances(ctx, typesOpts)
	if err != nil {
		return nil, err
	}
	return convertInstanceListToInterface(result), nil
}

// GetTRES retrieves all TRES (Trackable RESources)
func (c *AdapterClient) GetTRES(ctx context.Context) (*types.TRESList, error) {
	standaloneManager := c.adapter.GetStandaloneManager()
	if standaloneManager == nil {
		return nil, fmt.Errorf("standalone operations not supported for version %s", c.version)
	}
	result, err := standaloneManager.GetTRES(ctx)
	if err != nil {
		return nil, err
	}
	return convertTRESListToInterface(result), nil
}

// CreateTRES creates a new TRES entry
func (c *AdapterClient) CreateTRES(ctx context.Context, req *types.CreateTRESRequest) (*types.TRES, error) {
	standaloneManager := c.adapter.GetStandaloneManager()
	if standaloneManager == nil {
		return nil, fmt.Errorf("standalone operations not supported for version %s", c.version)
	}
	typesReq := convertCreateTRESRequestToTypes(req)
	result, err := standaloneManager.CreateTRES(ctx, typesReq)
	if err != nil {
		return nil, err
	}
	return convertTRESToInterface(result), nil
}

// Reconfigure triggers a SLURM reconfiguration
func (c *AdapterClient) Reconfigure(ctx context.Context) (*types.ReconfigureResponse, error) {
	standaloneManager := c.adapter.GetStandaloneManager()
	if standaloneManager == nil {
		return nil, fmt.Errorf("standalone operations not supported for version %s", c.version)
	}
	result, err := standaloneManager.Reconfigure(ctx)
	if err != nil {
		return nil, err
	}
	return convertReconfigureResponseToInterface(result), nil
}

// Manager wrappers to convert between common.types and interfaces types

// adapterJobManager wraps a common.JobAdapter to implement types.JobManager
type adapterJobManager struct {
	adapter common.JobAdapter
}

func (m *adapterJobManager) List(ctx context.Context, opts *types.ListJobsOptions) (*types.JobList, error) {
	// Convert options
	adapterOpts := &types.JobListOptions{}
	if opts != nil {
		// Only add UserID if non-empty (avoid creating slice with single empty string)
		if opts.UserID != "" {
			adapterOpts.Users = []string{opts.UserID}
		}
		// Only add Partition if non-empty (avoid creating slice with single empty string)
		if opts.Partition != "" {
			adapterOpts.Partitions = []string{opts.Partition}
		}
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

	// Check if result is nil before accessing it
	if result == nil {
		return &types.JobList{
			Jobs:  []types.Job{},
			Total: 0,
		}, nil
	}

	// Convert result
	jobList := &types.JobList{
		Jobs:  append([]types.Job{}, result.Jobs...),
		Total: result.Total, // Total from adapter is full count before pagination
	}

	return jobList, nil
}

func (m *adapterJobManager) Get(ctx context.Context, jobID string) (*types.Job, error) {
	// Convert string to int32 for adapter
	jobIDInt, err := strconv.ParseInt(jobID, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid job JobId: %w", err)
	}

	job, err := m.adapter.Get(ctx, int32(jobIDInt))
	if err != nil {
		return nil, err
	}
	return job, nil
}

// Helper functions for creating pointers
func ptrString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func ptrInt32(i int32) *int32    { return &i }
func ptrUint32(i uint32) *uint32 { return &i }
func ptrInt64(i int64) *int64    { return &i }

// convertMapToEnvList converts map[string]string to []string in "KEY=VALUE" format
func convertMapToEnvList(env map[string]string) []string {
	if len(env) == 0 {
		return nil
	}
	result := make([]string, 0, len(env))
	for k, v := range env {
		result = append(result, k+"="+v)
	}
	return result
}

func (m *adapterJobManager) Submit(ctx context.Context, job *types.JobSubmission) (*types.JobSubmitResponse, error) {
	// Convert submission - map from types.JobSubmission to types.JobCreate
	submission := &types.JobCreate{
		Name:                    ptrString(job.Name),
		Account:                 ptrString(job.Account),
		Script:                  ptrString(job.Script),
		Partition:               ptrString(job.Partition),
		MinimumCPUs:             ptrInt32(int32(job.CPUs)),
		TimeLimit:               ptrUint32(uint32(job.TimeLimit)),
		CurrentWorkingDirectory: ptrString(job.WorkingDir),
		Environment:             convertMapToEnvList(job.Environment),
		MinimumNodes:            ptrInt32(int32(job.Nodes)),
		Priority:                ptrUint32(uint32(job.Priority)),
	}

	// Set memory if provided
	if job.Memory > 0 {
		submission.MemoryPerNode = func() *uint64 { v := uint64(job.Memory); return &v }()
	}

	// Call adapter
	resp, err := m.adapter.Submit(ctx, submission)
	if err != nil {
		return nil, err
	}

	return &types.JobSubmitResponse{
		JobId: resp.JobId,
	}, nil
}

func (m *adapterJobManager) Update(ctx context.Context, jobID string, update *types.JobUpdate) error {
	// Convert string to int32 for adapter
	jobIDInt, err := strconv.ParseInt(jobID, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid job JobId: %w", err)
	}
	// JobUpdate is an alias for JobCreate - pass it directly
	// The adapter will use only the fields that are set
	return m.adapter.Update(ctx, int32(jobIDInt), update)
}

func (m *adapterJobManager) Cancel(ctx context.Context, jobID string) error {
	// Convert string to int32 for adapter
	jobIDInt, err := strconv.ParseInt(jobID, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid job JobId: %w", err)
	}
	return m.adapter.Cancel(ctx, int32(jobIDInt), nil)
}

// Hold holds a job (prevents it from running)
func (m *adapterJobManager) Hold(ctx context.Context, jobID string) error {
	// Convert string to int32 for adapter
	jobIDInt, err := strconv.ParseInt(jobID, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid job JobId: %w", err)
	}
	// Create hold request (hold = true)
	req := &types.JobHoldRequest{
		JobId: int32(jobIDInt),
		Hold:  true,
	}
	return m.adapter.Hold(ctx, req)
}

// Release releases a held job (allows it to run)
func (m *adapterJobManager) Release(ctx context.Context, jobID string) error {
	// Convert string to int32 for adapter
	jobIDInt, err := strconv.ParseInt(jobID, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid job JobId: %w", err)
	}
	// Create hold request (hold = false to release)
	req := &types.JobHoldRequest{
		JobId: int32(jobIDInt),
		Hold:  false,
	}
	return m.adapter.Hold(ctx, req)
}

// Signal sends a signal to a job
func (m *adapterJobManager) Signal(ctx context.Context, jobID string, signal string) error {
	// Convert string to int32 for adapter
	jobIDInt, err := strconv.ParseInt(jobID, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid job JobId: %w", err)
	}
	req := &types.JobSignalRequest{
		JobId:  int32(jobIDInt),
		Signal: signal,
	}
	return m.adapter.Signal(ctx, req)
}

// Notify sends a message to a job
func (m *adapterJobManager) Notify(ctx context.Context, jobID string, message string) error {
	// Convert string to int32 for adapter
	jobIDInt, err := strconv.ParseInt(jobID, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid job JobId: %w", err)
	}
	req := &types.JobNotifyRequest{
		JobId:   int32(jobIDInt),
		Message: message,
	}
	return m.adapter.Notify(ctx, req)
}

// Requeue requeues a job
func (m *adapterJobManager) Requeue(ctx context.Context, jobID string) error {
	// Convert string to int32 for adapter
	jobIDInt, err := strconv.ParseInt(jobID, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid job JobId: %w", err)
	}
	return m.adapter.Requeue(ctx, int32(jobIDInt))
}

func (m *adapterJobManager) Watch(ctx context.Context, opts *types.WatchJobsOptions) (<-chan types.JobEvent, error) {
	// Convert WatchJobsOptions to types.JobWatchOptions
	adapterOpts := &types.JobWatchOptions{}

	if opts != nil {
		// Convert JobIDs from []string to []int32
		if len(opts.JobIDs) > 0 {
			// Just watch the first job ID for now (adapter expects single job ID)
			jobIDInt, err := strconv.ParseInt(opts.JobIDs[0], 10, 32)
			if err == nil {
				adapterOpts.JobId = int32(jobIDInt)
			}
		}

		// Note: WatchJobsOptions.States (job states like RUNNING, PENDING) are not
		// mapped to JobWatchOptions.EventTypes (event names like start, end, fail)
		// as they represent fundamentally different concepts.
		// State filtering must be implemented at a higher level through event filtering.
	}

	// Call adapter's Watch method
	adapterEventChan, err := m.adapter.Watch(ctx, adapterOpts)
	if err != nil {
		return nil, err
	}

	// Create interface event channel
	interfaceEventChan := make(chan types.JobEvent, 10)

	// Start goroutine to convert events
	go func() {
		defer close(interfaceEventChan)

		for adapterEvent := range adapterEventChan {
			// Convert types.JobWatchEvent to types.JobEvent
			interfaceEvent := types.JobEvent{
				EventType:     adapterEvent.EventType,
				JobId:         adapterEvent.JobId,
				PreviousState: adapterEvent.PreviousState,
				NewState:      adapterEvent.NewState,
				EventTime:     adapterEvent.EventTime,
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

// Allocate allocates resources for a job
func (m *adapterJobManager) Allocate(ctx context.Context, req *types.JobAllocateRequest) (*types.JobAllocateResponse, error) {
	// Convert types.JobAllocateRequest to types.JobAllocateRequest
	adapterReq := &types.JobAllocateRequest{
		Name:      req.Name,
		Account:   req.Account,
		Partition: req.Partition,
		Nodes:     req.Nodes,
		Cpus:      int32(req.Cpus),
		TimeLimit: int32(req.TimeLimit), // Time limit in minutes
		QoS:       req.QoS,
	}

	// Call the adapter's Allocate method
	result, err := m.adapter.Allocate(ctx, adapterReq)
	if err != nil {
		return nil, err
	}

	// Convert types.JobAllocateResponse to types.JobAllocateResponse
	// Since types.JobAllocateResponse = types.JobAllocateResponse, just return it
	return result, nil
}

// adapterNodeManager wraps a common.NodeAdapter to implement types.NodeManager
type adapterNodeManager struct {
	adapter common.NodeAdapter
}

func (m *adapterNodeManager) List(ctx context.Context, opts *types.ListNodesOptions) (*types.NodeList, error) {
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

	// Check if result is nil before accessing it
	if result == nil {
		return &types.NodeList{
			Nodes: []types.Node{},
			Total: 0,
		}, nil
	}

	// Convert result
	nodeList := &types.NodeList{
		Nodes: append([]types.Node{}, result.Nodes...),
		Total: result.Total,
	}

	return nodeList, nil
}

func (m *adapterNodeManager) Get(ctx context.Context, nodeName string) (*types.Node, error) {
	node, err := m.adapter.Get(ctx, nodeName)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (m *adapterNodeManager) Update(ctx context.Context, nodeName string, update *types.NodeUpdate) error {
	// The common NodeUpdate type is now used directly
	return m.adapter.Update(ctx, nodeName, update)
}

func (m *adapterNodeManager) Watch(ctx context.Context, opts *types.WatchNodesOptions) (<-chan types.NodeEvent, error) {
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
	interfaceEventChan := make(chan types.NodeEvent, 10)

	// Start goroutine to convert events
	go func() {
		defer close(interfaceEventChan)

		for adapterEvent := range adapterEventChan {
			// Since types.NodeEvent = types.NodeEvent, just pass it through
			interfaceEvent := adapterEvent

			select {
			case interfaceEventChan <- interfaceEvent:
			case <-ctx.Done():
				return
			}
		}
	}()

	return interfaceEventChan, nil
}

// Delete deletes a node
func (m *adapterNodeManager) Delete(ctx context.Context, nodeName string) error {
	return m.adapter.Delete(ctx, nodeName)
}

// Drain drains a node, preventing new jobs from being scheduled on it
func (m *adapterNodeManager) Drain(ctx context.Context, nodeName string, reason string) error {
	return m.adapter.Drain(ctx, nodeName, reason)
}

// Resume resumes a drained node, allowing new jobs to be scheduled on it
func (m *adapterNodeManager) Resume(ctx context.Context, nodeName string) error {
	return m.adapter.Resume(ctx, nodeName)
}

// Helper function to convert types.Node to types.Node
// adapterPartitionManager wraps a common.PartitionAdapter
type adapterPartitionManager struct {
	adapter common.PartitionAdapter
}

func (m *adapterPartitionManager) List(ctx context.Context, opts *types.ListPartitionsOptions) (*types.PartitionList, error) {
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

	// Check if result is nil before accessing it
	if result == nil {
		return &types.PartitionList{
			Partitions: []types.Partition{},
			Total:      0,
		}, nil
	}

	// Convert result
	partitionList := &types.PartitionList{
		Partitions: append([]types.Partition{}, result.Partitions...),
		Total:      result.Total,
	}

	return partitionList, nil
}

func (m *adapterPartitionManager) Get(ctx context.Context, partitionName string) (*types.Partition, error) {
	partition, err := m.adapter.Get(ctx, partitionName)
	if err != nil {
		return nil, err
	}
	return partition, nil
}

func (m *adapterPartitionManager) Update(ctx context.Context, partitionName string, update *types.PartitionUpdate) error {
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

func (m *adapterPartitionManager) Watch(ctx context.Context, opts *types.WatchPartitionsOptions) (<-chan types.PartitionEvent, error) {
	// Implement polling-based watch since the adapter layer doesn't have Watch
	eventChan := make(chan types.PartitionEvent, 10)

	go func() {
		defer close(eventChan)

		// Track previous partition states
		prevPartitions := make(map[string]*types.Partition)
		pollInterval := 5 * time.Second

		ticker := time.NewTicker(pollInterval)
		defer ticker.Stop()

		// Helper to get state from partition
		getState := func(p *types.Partition) string {
			if p != nil && p.Partition != nil && len(p.Partition.State) > 0 {
				return string(p.Partition.State[0])
			}
			return ""
		}

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// Fetch current partitions
				result, err := m.adapter.List(ctx, &types.PartitionListOptions{})
				if err != nil {
					continue // Skip this poll on error
				}

				currentPartitions := make(map[string]*types.Partition)
				for i := range result.Partitions {
					p := &result.Partitions[i]
					name := ""
					if p.Name != nil {
						name = *p.Name
					}
					if name == "" {
						continue
					}
					// Filter by partition names if specified
					if opts != nil && len(opts.PartitionNames) > 0 {
						found := false
						for _, pn := range opts.PartitionNames {
							if pn == name {
								found = true
								break
							}
						}
						if !found {
							continue
						}
					}
					currentPartitions[name] = p
				}

				// Detect changes
				now := time.Now()

				// Check for new or changed partitions
				for name, current := range currentPartitions {
					prev, existed := prevPartitions[name]
					if !existed {
						// New partition
						if len(prevPartitions) > 0 { // Skip first poll
							event := types.PartitionEvent{
								EventTime:     now,
								EventType:     "created",
								PartitionName: name,
								Partition:     current,
								NewState:      types.PartitionState(getState(current)),
							}
							select {
							case eventChan <- event:
							case <-ctx.Done():
								return
							}
						}
					} else {
						// Check for state change
						prevState := getState(prev)
						currState := getState(current)
						if prevState != currState {
							event := types.PartitionEvent{
								EventTime:     now,
								EventType:     "state_change",
								PartitionName: name,
								PreviousState: types.PartitionState(prevState),
								NewState:      types.PartitionState(currState),
								Partition:     current,
							}
							select {
							case eventChan <- event:
							case <-ctx.Done():
								return
							}
						}
					}
				}

				// Check for deleted partitions
				for name := range prevPartitions {
					if _, exists := currentPartitions[name]; !exists {
						event := types.PartitionEvent{
							EventTime:     now,
							EventType:     "deleted",
							PartitionName: name,
						}
						select {
						case eventChan <- event:
						case <-ctx.Done():
							return
						}
					}
				}

				prevPartitions = currentPartitions
			}
		}
	}()

	return eventChan, nil
}

// Create creates a new partition
func (m *adapterPartitionManager) Create(ctx context.Context, partition *types.PartitionCreate) (*types.PartitionCreateResponse, error) {
	// Since types.PartitionCreate = types.PartitionCreate, no conversion needed
	resp, err := m.adapter.Create(ctx, partition)
	if err != nil {
		return nil, err
	}

	// Since types.PartitionCreateResponse = types.PartitionCreateResponse, just return it
	return resp, nil
}

// Delete deletes a partition
func (m *adapterPartitionManager) Delete(ctx context.Context, partitionName string) error {
	return m.adapter.Delete(ctx, partitionName)
}

// Helper function to convert types.Account to types.Account

// Helper function to convert types.User to types.User

// Helper function to convert types.QoS to types.QoS

// adapterInfoManager provides info operations via the adapter
type adapterInfoManager struct {
	adapter common.InfoAdapter
	version string
}

func (m *adapterInfoManager) Ping(ctx context.Context) error {
	return m.adapter.Ping(ctx)
}

func (m *adapterInfoManager) Get(ctx context.Context) (*types.ClusterInfo, error) {
	result, err := m.adapter.Get(ctx)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errors.NewSlurmError(errors.ErrorCodeResourceNotFound, "cluster info not found")
	}
	return convertTypesClusterInfoToInterface(result), nil
}

func (m *adapterInfoManager) Stats(ctx context.Context) (*types.ClusterStats, error) {
	result, err := m.adapter.Stats(ctx)
	if err != nil {
		return nil, err
	}
	return convertTypesClusterStatsToInterface(result), nil
}

func (m *adapterInfoManager) Version(ctx context.Context) (*types.APIVersion, error) {
	result, err := m.adapter.Version(ctx)
	if err != nil {
		return nil, err
	}
	return convertTypesAPIVersionToInterface(result), nil
}

func (m *adapterInfoManager) PingDatabase(ctx context.Context) error {
	return m.adapter.PingDatabase(ctx)
}

// Type converters for Info types
func convertTypesClusterInfoToInterface(t *types.ClusterInfo) *types.ClusterInfo {
	if t == nil {
		return nil
	}
	return &types.ClusterInfo{
		ClusterName: t.ClusterName,
		Version:     t.Version,
		Release:     t.Release,
		APIVersion:  t.APIVersion,
		Uptime:      t.Uptime,
	}
}

func convertTypesClusterStatsToInterface(t *types.ClusterStats) *types.ClusterStats {
	if t == nil {
		return nil
	}
	return &types.ClusterStats{
		TotalNodes:     t.TotalNodes,
		IdleNodes:      t.IdleNodes,
		AllocatedNodes: t.AllocatedNodes,
		TotalCPUs:      t.TotalCPUs,
		IdleCPUs:       t.IdleCPUs,
		AllocatedCPUs:  t.AllocatedCPUs,
		TotalJobs:      t.TotalJobs,
		RunningJobs:    t.RunningJobs,
		PendingJobs:    t.PendingJobs,
		CompletedJobs:  t.CompletedJobs,
	}
}

func convertTypesAPIVersionToInterface(t *types.APIVersion) *types.APIVersion {
	if t == nil {
		return nil
	}
	return &types.APIVersion{
		Version:     t.Version,
		Release:     t.Release,
		Description: t.Description,
		Deprecated:  t.Deprecated,
	}
}

// Other manager implementations...

type adapterQoSManager struct {
	adapter common.QoSAdapter
}

func (m *adapterQoSManager) List(ctx context.Context, opts *types.ListQoSOptions) (*types.QoSList, error) {
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

	// Check if result is nil before accessing it
	if result == nil {
		return &types.QoSList{
			QoS:   []types.QoS{},
			Total: 0,
		}, nil
	}

	// Convert result
	qosList := &types.QoSList{
		QoS:   append([]types.QoS{}, result.QoS...),
		Total: result.Total,
	}

	return qosList, nil
}

func (m *adapterQoSManager) Get(ctx context.Context, qosName string) (*types.QoS, error) {
	qos, err := m.adapter.Get(ctx, qosName)
	if err != nil {
		return nil, err
	}
	return qos, nil
}

func (m *adapterQoSManager) Create(ctx context.Context, qos *types.QoSCreate) (*types.QoSCreateResponse, error) {
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

	return &types.QoSCreateResponse{
		QoSName: resp.QoSName,
	}, nil
}

func (m *adapterQoSManager) Update(ctx context.Context, qosName string, update *types.QoSUpdate) error {
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
	adapter            common.AccountAdapter
	associationAdapter common.AssociationAdapter
}

func (m *adapterAccountManager) List(ctx context.Context, opts *types.ListAccountsOptions) (*types.AccountList, error) {
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

	// Check if result is nil before accessing it
	if result == nil {
		return &types.AccountList{
			Accounts: []types.Account{},
			Total:    0,
		}, nil
	}

	// Convert result
	accountList := &types.AccountList{
		Accounts: append([]types.Account{}, result.Accounts...),
		Total:    result.Total,
	}

	return accountList, nil
}

func (m *adapterAccountManager) Get(ctx context.Context, accountName string) (*types.Account, error) {
	account, err := m.adapter.Get(ctx, accountName)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (m *adapterAccountManager) Create(ctx context.Context, account *types.AccountCreate) (*types.AccountCreateResponse, error) {
	// Convert create request
	adapterCreate := &types.AccountCreate{
		Name:         account.Name,
		Description:  account.Description,
		Organization: account.Organization,
		// Add other fields as needed
	}

	// Call adapter
	resp, err := m.adapter.Create(ctx, adapterCreate)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m *adapterAccountManager) Update(ctx context.Context, accountName string, update *types.AccountUpdate) error {
	// Convert update request
	adapterUpdate := &types.AccountUpdate{
		Description:  update.Description,
		Organization: update.Organization,
		// Add other fields as needed
	}

	return m.adapter.Update(ctx, accountName, adapterUpdate)
}

func (m *adapterAccountManager) Delete(ctx context.Context, accountName string) error {
	return m.adapter.Delete(ctx, accountName)
}

func (m *adapterAccountManager) GetAccountHierarchy(ctx context.Context, rootAccount string) (*types.AccountHierarchy, error) {
	ext := &extendedAccountManager{adapter: m.adapter, associationAdapter: m.associationAdapter}
	return ext.GetAccountHierarchy(ctx, rootAccount)
}

func (m *adapterAccountManager) GetParentAccounts(ctx context.Context, accountName string) ([]*types.Account, error) {
	ext := &extendedAccountManager{adapter: m.adapter, associationAdapter: m.associationAdapter}
	return ext.GetParentAccounts(ctx, accountName)
}

func (m *adapterAccountManager) GetChildAccounts(ctx context.Context, accountName string, depth int) ([]*types.Account, error) {
	ext := &extendedAccountManager{adapter: m.adapter, associationAdapter: m.associationAdapter}
	return ext.GetChildAccounts(ctx, accountName, depth)
}

func (m *adapterAccountManager) GetAccountQuotas(ctx context.Context, accountName string) (*types.AccountQuota, error) {
	ext := &extendedAccountManager{adapter: m.adapter, associationAdapter: m.associationAdapter}
	return ext.GetAccountQuotas(ctx, accountName)
}

func (m *adapterAccountManager) GetAccountQuotaUsage(ctx context.Context, accountName string, timeframe string) (*types.AccountUsage, error) {
	ext := &extendedAccountManager{adapter: m.adapter, associationAdapter: m.associationAdapter}
	return ext.GetAccountQuotaUsage(ctx, accountName, timeframe)
}

func (m *adapterAccountManager) GetAccountUsers(ctx context.Context, accountName string, opts *types.ListAccountUsersOptions) ([]*types.UserAccountAssociation, error) {
	ext := &extendedAccountManager{adapter: m.adapter, associationAdapter: m.associationAdapter}
	return ext.GetAccountUsers(ctx, accountName, opts)
}

func (m *adapterAccountManager) ValidateUserAccess(ctx context.Context, userName, accountName string) (*types.UserAccessValidation, error) {
	ext := &extendedAccountManager{adapter: m.adapter, associationAdapter: m.associationAdapter}
	return ext.ValidateUserAccess(ctx, userName, accountName)
}

func (m *adapterAccountManager) GetAccountUsersWithPermissions(ctx context.Context, accountName string, permissions []string) ([]*types.UserAccountAssociation, error) {
	ext := &extendedAccountManager{adapter: m.adapter, associationAdapter: m.associationAdapter}
	return ext.GetAccountUsersWithPermissions(ctx, accountName, permissions)
}

func (m *adapterAccountManager) GetAccountFairShare(ctx context.Context, accountName string) (*types.AccountFairShare, error) {
	ext := &extendedAccountManager{adapter: m.adapter, associationAdapter: m.associationAdapter}
	return ext.GetAccountFairShare(ctx, accountName)
}

func (m *adapterAccountManager) GetFairShareHierarchy(ctx context.Context, rootAccount string) (*types.FairShareHierarchy, error) {
	ext := &extendedAccountManager{adapter: m.adapter, associationAdapter: m.associationAdapter}
	return ext.GetFairShareHierarchy(ctx, rootAccount)
}

type adapterUserManager struct {
	adapter            common.UserAdapter
	accountAdapter     common.AccountAdapter
	associationAdapter common.AssociationAdapter
}

func (m *adapterUserManager) List(ctx context.Context, opts *types.ListUsersOptions) (*types.UserList, error) {
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

	// Check if result is nil before accessing it
	if result == nil {
		return &types.UserList{
			Users: []types.User{},
			Total: 0,
		}, nil
	}

	// Convert result
	userList := &types.UserList{
		Users: append([]types.User{}, result.Users...),
		Total: result.Total,
	}

	return userList, nil
}

func (m *adapterUserManager) Get(ctx context.Context, userName string) (*types.User, error) {
	user, err := m.adapter.Get(ctx, userName)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (m *adapterUserManager) GetUserAccounts(ctx context.Context, userName string) ([]*types.UserAccount, error) {
	ext := &extendedUserManager{adapter: m.adapter, accountAdapter: m.accountAdapter, associationAdapter: m.associationAdapter}
	return ext.GetUserAccounts(ctx, userName)
}

func (m *adapterUserManager) GetUserQuotas(ctx context.Context, userName string) (*types.UserQuota, error) {
	ext := &extendedUserManager{adapter: m.adapter, accountAdapter: m.accountAdapter, associationAdapter: m.associationAdapter}
	return ext.GetUserQuotas(ctx, userName)
}

func (m *adapterUserManager) GetUserDefaultAccount(ctx context.Context, userName string) (*types.Account, error) {
	ext := &extendedUserManager{adapter: m.adapter, accountAdapter: m.accountAdapter, associationAdapter: m.associationAdapter}
	return ext.GetUserDefaultAccount(ctx, userName)
}

func (m *adapterUserManager) GetUserFairShare(ctx context.Context, userName string) (*types.UserFairShare, error) {
	ext := &extendedUserManager{adapter: m.adapter, accountAdapter: m.accountAdapter, associationAdapter: m.associationAdapter}
	return ext.GetUserFairShare(ctx, userName)
}

func (m *adapterUserManager) CalculateJobPriority(ctx context.Context, userName string, jobSubmission *types.JobSubmission) (*types.JobPriorityInfo, error) {
	ext := &extendedUserManager{adapter: m.adapter, accountAdapter: m.accountAdapter, associationAdapter: m.associationAdapter}
	return ext.CalculateJobPriority(ctx, userName, jobSubmission)
}

func (m *adapterUserManager) ValidateUserAccountAccess(ctx context.Context, userName, accountName string) (*types.UserAccessValidation, error) {
	ext := &extendedUserManager{adapter: m.adapter, accountAdapter: m.accountAdapter, associationAdapter: m.associationAdapter}
	return ext.ValidateUserAccountAccess(ctx, userName, accountName)
}

func (m *adapterUserManager) GetUserAccountAssociations(ctx context.Context, userName string, opts *types.ListUserAccountAssociationsOptions) ([]*types.UserAccountAssociation, error) {
	ext := &extendedUserManager{adapter: m.adapter, accountAdapter: m.accountAdapter, associationAdapter: m.associationAdapter}
	return ext.GetUserAccountAssociations(ctx, userName, opts)
}

func (m *adapterUserManager) GetBulkUserAccounts(ctx context.Context, userNames []string) (map[string][]*types.UserAccount, error) {
	ext := &extendedUserManager{adapter: m.adapter, accountAdapter: m.accountAdapter, associationAdapter: m.associationAdapter}
	return ext.GetBulkUserAccounts(ctx, userNames)
}

func (m *adapterUserManager) GetBulkAccountUsers(ctx context.Context, accountNames []string) (map[string][]*types.UserAccountAssociation, error) {
	ext := &extendedUserManager{adapter: m.adapter, accountAdapter: m.accountAdapter, associationAdapter: m.associationAdapter}
	return ext.GetBulkAccountUsers(ctx, accountNames)
}

// Create creates a new user
func (m *adapterUserManager) Create(ctx context.Context, user *types.UserCreate) (*types.UserCreateResponse, error) {
	// Since types.UserCreate = types.UserCreate, no conversion needed
	resp, err := m.adapter.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	// Since types.UserCreateResponse = types.UserCreateResponse, just return it
	return resp, nil
}

// Update updates a user
func (m *adapterUserManager) Update(ctx context.Context, userName string, update *types.UserUpdate) error {
	// Convert to adapter type
	adapterUpdate := &types.UserUpdate{}
	if update != nil {
		if update.DefaultAccount != nil {
			adapterUpdate.DefaultAccount = update.DefaultAccount
		}
		if update.DefaultWCKey != nil {
			adapterUpdate.DefaultWCKey = update.DefaultWCKey
		}
		if update.AdminLevel != nil {
			adminLevel := types.AdminLevel(*update.AdminLevel)
			adapterUpdate.AdminLevel = &adminLevel
		}
		// Add other fields as needed
	}

	return m.adapter.Update(ctx, userName, adapterUpdate)
}

// Delete deletes a user
func (m *adapterUserManager) Delete(ctx context.Context, userName string) error {
	return m.adapter.Delete(ctx, userName)
}

type adapterReservationManager struct {
	adapter common.ReservationAdapter
}

func (m *adapterReservationManager) List(ctx context.Context, opts *types.ListReservationsOptions) (*types.ReservationList, error) {
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

	// Check if result is nil before accessing it
	if result == nil {
		return &types.ReservationList{
			Reservations: []types.Reservation{},
			Total:        0,
		}, nil
	}

	// Convert result
	reservationList := &types.ReservationList{
		Reservations: append([]types.Reservation{}, result.Reservations...),
		Total:        result.Total,
	}

	return reservationList, nil
}

func (m *adapterReservationManager) Get(ctx context.Context, reservationName string) (*types.Reservation, error) {
	// Call adapter
	result, err := m.adapter.Get(ctx, reservationName)
	if err != nil {
		return nil, err
	}

	// Check if result is nil before dereferencing
	if result == nil {
		return nil, errors.NewSlurmError(errors.ErrorCodeResourceNotFound, fmt.Sprintf("reservation %s not found", reservationName))
	}

	return result, nil
}

func (m *adapterReservationManager) Create(ctx context.Context, reservation *types.ReservationCreate) (*types.ReservationCreateResponse, error) {
	// Since types.ReservationCreate = types.ReservationCreate, no conversion needed
	result, err := m.adapter.Create(ctx, reservation)
	if err != nil {
		return nil, err
	}

	// Convert result
	return &types.ReservationCreateResponse{
		ReservationName: result.ReservationName,
	}, nil
}

func (m *adapterReservationManager) Update(ctx context.Context, reservationName string, update *types.ReservationUpdate) error {
	// Since types.ReservationUpdate = types.ReservationUpdate, no conversion needed
	return m.adapter.Update(ctx, reservationName, update)
}

func (m *adapterReservationManager) Delete(ctx context.Context, reservationName string) error {
	return m.adapter.Delete(ctx, reservationName)
}

type adapterAssociationManager struct {
	adapter common.AssociationAdapter
}

func (m *adapterAssociationManager) List(ctx context.Context, opts *types.ListAssociationsOptions) (*types.AssociationList, error) {
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

	// Check if result is nil before accessing it
	if result == nil {
		return &types.AssociationList{
			Associations: []types.Association{},
			Total:        0,
		}, nil
	}

	// Since types.AssociationList = types.AssociationList, just return it
	return result, nil
}

func (m *adapterAssociationManager) Get(ctx context.Context, associationID string) (*types.Association, error) {
	if associationID == "" {
		return nil, fmt.Errorf("association ID must be specified")
	}

	// Call adapter
	result, err := m.adapter.Get(ctx, associationID)
	if err != nil {
		return nil, err
	}

	// Check if result is nil before dereferencing
	if result == nil {
		return nil, errors.NewSlurmError(errors.ErrorCodeResourceNotFound, fmt.Sprintf("association %s not found", associationID))
	}

	return result, nil
}

func (m *adapterAssociationManager) Create(ctx context.Context, associations []*types.AssociationCreate) (*types.AssociationCreateResponse, error) {
	// For now, handle single association creation
	if len(associations) == 0 {
		return nil, fmt.Errorf("no associations provided")
	}

	// Since types.AssociationCreate = types.AssociationCreate, no conversion needed
	resp, err := m.adapter.Create(ctx, associations[0])
	if err != nil {
		return nil, err
	}

	// Since types.AssociationCreateResponse = types.AssociationCreateResponse, just return it
	return resp, nil
}

func (m *adapterAssociationManager) Update(ctx context.Context, associations []*types.AssociationUpdate) error {
	// For now, handle single association update
	if len(associations) == 0 {
		return fmt.Errorf("no associations provided")
	}

	assoc := associations[0]

	// Extract the association ID from the update request
	if assoc.ID == nil || *assoc.ID == 0 {
		return fmt.Errorf("associationID is required")
	}
	associationID := fmt.Sprintf("%d", *assoc.ID)

	return m.adapter.Update(ctx, associationID, assoc)
}

func (m *adapterAssociationManager) Delete(ctx context.Context, associationID string) error {
	if associationID == "" {
		return fmt.Errorf("association ID must be specified")
	}
	return m.adapter.Delete(ctx, associationID)
}

func (m *adapterAssociationManager) GetUserAssociations(ctx context.Context, userName string) ([]*types.Association, error) {
	// List associations and filter by user
	opts := &types.AssociationListOptions{
		Limit: 1000, // Get all associations for the user
	}

	result, err := m.adapter.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	// Filter associations for the specified user
	userAssociations := []*types.Association{}
	for _, assoc := range result.Associations {
		if assoc.User == userName {
			// Need to create a copy to take address
			a := assoc
			userAssociations = append(userAssociations, &a)
		}
	}

	return userAssociations, nil
}

func (m *adapterAssociationManager) GetAccountAssociations(ctx context.Context, accountName string) ([]*types.Association, error) {
	// List associations and filter by account
	opts := &types.AssociationListOptions{
		Limit: 1000, // Get all associations for the account
	}

	result, err := m.adapter.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	// Filter associations for the specified account
	accountAssociations := []*types.Association{}
	for _, assoc := range result.Associations {
		if assoc.Account != nil && *assoc.Account == accountName {
			a := assoc
			accountAssociations = append(accountAssociations, &a)
		}
	}

	return accountAssociations, nil
}

func (m *adapterAssociationManager) ValidateAssociation(ctx context.Context, user, account, cluster string) (bool, error) {
	// Build association ID from components
	associationID := fmt.Sprintf("%s:%s", user, account)
	if cluster != "" {
		associationID = fmt.Sprintf("%s@%s", associationID, cluster)
	}

	_, err := m.Get(ctx, associationID)
	if err != nil {
		// Association doesn't exist or error occurred
		return false, err
	}

	// Association exists
	return true, nil
}

// Helper function to convert types.Reservation to types.Reservation

// Helper function to convert types.Association to types.Association

// adapterClusterManager wraps a common.ClusterAdapter to implement types.ClusterManager
type adapterClusterManager struct {
	adapter common.ClusterAdapter
}

func (m *adapterClusterManager) List(ctx context.Context, opts *types.ListClustersOptions) (*types.ClusterList, error) {
	// Call adapter (types.ClusterListOptions doesn't have Names field, filtering done after)
	adapterOpts := &types.ClusterListOptions{}

	// Call adapter
	result, err := m.adapter.List(ctx, adapterOpts)
	if err != nil {
		return nil, err
	}

	// Check if result is nil before accessing it
	if result == nil {
		return &types.ClusterList{
			Clusters: []types.Cluster{},
			Total:    0,
		}, nil
	}

	// Apply client-side filtering by names if specified
	if opts != nil && len(opts.Names) > 0 {
		filtered := []types.Cluster{}
		for _, cluster := range result.Clusters {
			found := false
			for _, name := range opts.Names {
				if cluster.Name != nil && *cluster.Name == name {
					found = true
					break
				}
			}
			if found {
				filtered = append(filtered, cluster)
			}
		}
		return &types.ClusterList{
			Clusters: filtered,
			Total:    len(filtered),
		}, nil
	}

	// No filtering needed, return result directly (since types.ClusterList = types.ClusterList)
	return result, nil
}

func (m *adapterClusterManager) Get(ctx context.Context, clusterName string) (*types.Cluster, error) {
	cluster, err := m.adapter.Get(ctx, clusterName)
	if err != nil {
		return nil, err
	}
	return cluster, nil
}

func (m *adapterClusterManager) Create(ctx context.Context, cluster *types.ClusterCreate) (*types.ClusterCreateResponse, error) {
	return m.adapter.Create(ctx, cluster)
}

func (m *adapterClusterManager) Delete(ctx context.Context, clusterName string) error {
	return m.adapter.Delete(ctx, clusterName)
}

// convertClusterToInterface converts types.Cluster to types.Cluster

// adapterWCKeyManager wraps a common.WCKeyAdapter to implement types.WCKeyManager
type adapterWCKeyManager struct {
	adapter common.WCKeyAdapter
}

func (m *adapterWCKeyManager) List(ctx context.Context, opts *types.WCKeyListOptions) (*types.WCKeyList, error) {
	// Convert options
	adapterOpts := &types.WCKeyListOptions{}
	if opts != nil {
		adapterOpts.Names = opts.Names
		adapterOpts.Users = opts.Users
		adapterOpts.Clusters = opts.Clusters
	}

	// Call adapter
	result, err := m.adapter.List(ctx, adapterOpts)
	if err != nil {
		return nil, err
	}

	// Check if result is nil before accessing it
	if result == nil {
		return &types.WCKeyList{
			WCKeys: []types.WCKey{},
			Total:  0,
		}, nil
	}

	// Convert result
	wckeys := make([]types.WCKey, 0, len(result.WCKeys))
	for _, wckey := range result.WCKeys {
		wckeys = append(wckeys, types.WCKey{
			Name:    wckey.Name,
			User:    wckey.User,
			Cluster: wckey.Cluster,
		})
	}

	return &types.WCKeyList{
		WCKeys: wckeys,
		Total:  len(wckeys),
	}, nil
}

func (m *adapterWCKeyManager) Get(ctx context.Context, wckeyName, user, cluster string) (*types.WCKey, error) {
	// For the adapter Get method, we pass the ID constructed from name, user, cluster
	wcKeyID := fmt.Sprintf("%s:%s:%s", wckeyName, user, cluster)

	result, err := m.adapter.Get(ctx, wcKeyID)
	if err != nil {
		return nil, err
	}

	// Check if result is nil before accessing fields
	if result == nil {
		return nil, errors.NewSlurmError(errors.ErrorCodeResourceNotFound, fmt.Sprintf("WCKey %s not found", wcKeyID))
	}

	return &types.WCKey{
		Name:    result.Name,
		User:    result.User,
		Cluster: result.Cluster,
	}, nil
}

func (m *adapterWCKeyManager) Create(ctx context.Context, wckey *types.WCKeyCreate) (*types.WCKeyCreateResponse, error) {
	// Convert request
	adapterReq := &types.WCKeyCreate{
		Name:    wckey.Name,
		User:    wckey.User,
		Cluster: wckey.Cluster,
	}

	// Call adapter
	resp, err := m.adapter.Create(ctx, adapterReq)
	if err != nil {
		return nil, err
	}

	// Since types.WCKeyCreateResponse = types.WCKeyCreateResponse, just return it
	return resp, nil
}

func (m *adapterWCKeyManager) Delete(ctx context.Context, wckeyID string) error {
	return m.adapter.Delete(ctx, wckeyID)
}

// === Type Conversion Helpers for Standalone Operations ===

// convertLicenseListToInterface converts types.LicenseList to types.LicenseList
func convertLicenseListToInterface(list *types.LicenseList) *types.LicenseList {
	if list == nil {
		return nil
	}

	interfaceLicenses := make([]types.License, len(list.Licenses))
	for i, license := range list.Licenses {
		interfaceLicenses[i] = types.License{
			Name:       license.Name,
			Total:      license.Total,
			Used:       license.Used,
			Free:       license.Free,
			Reserved:   license.Reserved,
			RemoteUsed: license.RemoteUsed,
		}
	}

	return &types.LicenseList{
		Licenses: interfaceLicenses,
		Meta:     list.Meta,
	}
}

// convertSharesListToInterface converts types.SharesList to types.SharesList
func convertSharesListToInterface(list *types.SharesList) *types.SharesList {
	if list == nil {
		return nil
	}

	interfaceShares := make([]types.Share, len(list.Shares))
	for i, share := range list.Shares {
		interfaceShares[i] = types.Share{
			Cluster:          share.Cluster,
			Account:          share.Account,
			User:             share.User,
			Partition:        share.Partition,
			EffectiveUsage:   share.EffectiveUsage,
			FairshareLevel:   share.FairshareLevel,
			FairshareUsage:   share.FairshareUsage,
			FairshareShares:  share.FairshareShares,
			NormalizedShares: share.NormalizedShares,
			NormalizedUsage:  share.NormalizedUsage,
			RawShares:        share.RawShares,
			RawUsage:         share.RawUsage,
			SharesUsed:       share.SharesUsed,
			RunSeconds:       share.RunSeconds,
			AssocID:          share.AssocID,
			ParentAccount:    share.ParentAccount,
			Meta:             share.Meta,
		}
	}

	return &types.SharesList{
		Shares: interfaceShares,
		Meta:   list.Meta,
	}
}

// convertConfigToInterface converts types.Config to types.Config
// Since both types are identical, this is a pass-through
func convertConfigToInterface(cfg *types.Config) *types.Config {
	return cfg
}

// convertDiagnosticsToInterface converts types.Diagnostics to types.Diagnostics
// Since both types are identical, this is a pass-through
func convertDiagnosticsToInterface(diag *types.Diagnostics) *types.Diagnostics {
	return diag
}

// convertInstanceToInterface converts types.Instance to types.Instance
// Since both types are identical, this is a pass-through
func convertInstanceToInterface(inst *types.Instance) *types.Instance {
	return inst
}

// convertInstanceListToInterface converts types.InstanceList to types.InstanceList
func convertInstanceListToInterface(list *types.InstanceList) *types.InstanceList {
	if list == nil {
		return nil
	}

	interfaceInstances := make([]types.Instance, len(list.Instances))
	for i, instance := range list.Instances {
		interfaceInstances[i] = *convertInstanceToInterface(&instance)
	}

	return &types.InstanceList{
		Instances: interfaceInstances,
		Meta:      list.Meta,
	}
}

// convertTRESListToInterface converts types.TRESList to types.TRESList
// Since both types are identical, this is a pass-through
func convertTRESListToInterface(list *types.TRESList) *types.TRESList {
	return list
}

// convertTRESToInterface converts types.TRES to types.TRES
// Since both types are identical, this is a pass-through
func convertTRESToInterface(t *types.TRES) *types.TRES {
	return t
}

// convertReconfigureResponseToInterface converts types.ReconfigureResponse to types.ReconfigureResponse
func convertReconfigureResponseToInterface(resp *types.ReconfigureResponse) *types.ReconfigureResponse {
	if resp == nil {
		return nil
	}

	return &types.ReconfigureResponse{
		Status:   resp.Status,
		Message:  resp.Message,
		Warnings: resp.Warnings,
		Errors:   resp.Errors,
		Meta:     resp.Meta,
	}
}

// === Input Conversion Helpers ===

// convertGetSharesOptionsToTypes converts types.GetSharesOptions to types.GetSharesOptions
func convertGetSharesOptionsToTypes(opts *types.GetSharesOptions) *types.GetSharesOptions {
	if opts == nil {
		return nil
	}

	return &types.GetSharesOptions{
		Users:     opts.Users,
		Accounts:  opts.Accounts,
		Clusters:  opts.Clusters,
	}
}

// convertGetInstanceOptionsToTypes converts types.GetInstanceOptions to types.GetInstanceOptions
// Since both types are identical, this is a pass-through
func convertGetInstanceOptionsToTypes(opts *types.GetInstanceOptions) *types.GetInstanceOptions {
	return opts
}

// convertGetInstancesOptionsToTypes converts types.GetInstancesOptions to types.GetInstancesOptions
// Since both types are identical, this is a pass-through
func convertGetInstancesOptionsToTypes(opts *types.GetInstancesOptions) *types.GetInstancesOptions {
	return opts
}

// convertCreateTRESRequestToTypes converts types.CreateTRESRequest to types.CreateTRESRequest
func convertCreateTRESRequestToTypes(req *types.CreateTRESRequest) *types.CreateTRESRequest {
	if req == nil {
		return nil
	}

	return &types.CreateTRESRequest{
		Type:        req.Type,
		Name:        req.Name,
		Description: req.Description,
	}
}

// calculatePercentage calculates percentage with safe division
func calculatePercentage(used, total int) float64 {
	if total == 0 {
		return 0.0
	}
	return float64(used) / float64(total) * 100.0
}
