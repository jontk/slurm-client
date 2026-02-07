// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package streaming

import (
	"context"

	types "github.com/jontk/slurm-client/api"
)

// Mock implementations for testing

type mockSlurmClient struct {
	jobs       *mockJobManager
	nodes      *mockNodeManager
	partitions *mockPartitionManager
}

func (m *mockSlurmClient) Version() string { return "test" }
func (m *mockSlurmClient) Jobs() types.JobManager {
	return m.jobs
}
func (m *mockSlurmClient) Nodes() types.NodeManager {
	return m.nodes
}
func (m *mockSlurmClient) Partitions() types.PartitionManager {
	return m.partitions
}
func (m *mockSlurmClient) Info() types.InfoManager                           { return nil }
func (m *mockSlurmClient) Reservations() types.ReservationManager            { return nil }
func (m *mockSlurmClient) QoS() types.QoSManager                             { return nil }
func (m *mockSlurmClient) Accounts() types.AccountManager                    { return nil }
func (m *mockSlurmClient) Users() types.UserManager                          { return nil }
func (m *mockSlurmClient) Clusters() types.ClusterManager                    { return nil }
func (m *mockSlurmClient) Associations() types.AssociationManager            { return nil }
func (m *mockSlurmClient) WCKeys() types.WCKeyManager                        { return nil }
func (m *mockSlurmClient) Analytics() types.AnalyticsManager                { return nil }
func (m *mockSlurmClient) GetLicenses(ctx context.Context) (*types.LicenseList, error) {
	return nil, nil
}
func (m *mockSlurmClient) GetShares(ctx context.Context, opts *types.GetSharesOptions) (*types.SharesList, error) {
	return nil, nil
}
func (m *mockSlurmClient) GetConfig(ctx context.Context) (*types.Config, error) {
	return nil, nil
}
func (m *mockSlurmClient) GetDiagnostics(ctx context.Context) (*types.Diagnostics, error) {
	return nil, nil
}
func (m *mockSlurmClient) GetDBDiagnostics(ctx context.Context) (*types.Diagnostics, error) {
	return nil, nil
}
func (m *mockSlurmClient) GetInstance(ctx context.Context, opts *types.GetInstanceOptions) (*types.Instance, error) {
	return nil, nil
}
func (m *mockSlurmClient) GetInstances(ctx context.Context, opts *types.GetInstancesOptions) (*types.InstanceList, error) {
	return nil, nil
}
func (m *mockSlurmClient) GetTRES(ctx context.Context) (*types.TRESList, error) {
	return nil, nil
}
func (m *mockSlurmClient) CreateTRES(ctx context.Context, req *types.CreateTRESRequest) (*types.TRES, error) {
	return nil, nil
}
func (m *mockSlurmClient) Reconfigure(ctx context.Context) (*types.ReconfigureResponse, error) {
	return nil, nil
}
func (m *mockSlurmClient) Capabilities() types.ClientCapabilities {
	return types.ClientCapabilities{
		Version:           "v0.0.42",
		SupportsJobs:      true,
		SupportsNodes:     true,
		SupportsJobSubmit: true,
	}
}
func (m *mockSlurmClient) Close() error { return nil }

type mockJobManager struct {
	watchFunc func(ctx context.Context, opts *types.WatchJobsOptions) (<-chan types.JobEvent, error)
}

func (m *mockJobManager) List(ctx context.Context, opts *types.ListJobsOptions) (*types.JobList, error) {
	return nil, nil
}
func (m *mockJobManager) Get(ctx context.Context, jobID string) (*types.Job, error) {
	return nil, nil
}
func (m *mockJobManager) Submit(ctx context.Context, job *types.JobSubmission) (*types.JobSubmitResponse, error) {
	return nil, nil
}
func (m *mockJobManager) Allocate(ctx context.Context, req *types.JobAllocateRequest) (*types.JobAllocateResponse, error) {
	return nil, nil
}
func (m *mockJobManager) Cancel(ctx context.Context, jobID string) error { return nil }
func (m *mockJobManager) Requeue(ctx context.Context, jobID string) error { return nil }
func (m *mockJobManager) Update(ctx context.Context, jobID string, update *types.JobUpdate) error {
	return nil
}
func (m *mockJobManager) Steps(ctx context.Context, jobID string) (*types.JobStepList, error) {
	return nil, nil
}
func (m *mockJobManager) Watch(ctx context.Context, opts *types.WatchJobsOptions) (<-chan types.JobEvent, error) {
	if m.watchFunc != nil {
		return m.watchFunc(ctx, opts)
	}
	return nil, nil
}
func (m *mockJobManager) Hold(ctx context.Context, jobID string) error    { return nil }
func (m *mockJobManager) Release(ctx context.Context, jobID string) error { return nil }
func (m *mockJobManager) Signal(ctx context.Context, jobID string, signal string) error {
	return nil
}
func (m *mockJobManager) Notify(ctx context.Context, jobID string, message string) error {
	return nil
}

// NOTE: Analytics methods removed - JobManager no longer includes AnalyticsManager
// Analytics is now accessed via client.Analytics() which returns nil for mocks

type mockNodeManager struct {
	watchFunc func(ctx context.Context, opts *types.WatchNodesOptions) (<-chan types.NodeEvent, error)
}

func (m *mockNodeManager) List(ctx context.Context, opts *types.ListNodesOptions) (*types.NodeList, error) {
	return nil, nil
}
func (m *mockNodeManager) Get(ctx context.Context, nodeName string) (*types.Node, error) {
	return nil, nil
}
func (m *mockNodeManager) Update(ctx context.Context, nodeName string, update *types.NodeUpdate) error {
	return nil
}
func (m *mockNodeManager) Delete(ctx context.Context, nodeName string) error {
	return nil
}
func (m *mockNodeManager) Drain(ctx context.Context, nodeName string, reason string) error {
	return nil
}
func (m *mockNodeManager) Resume(ctx context.Context, nodeName string) error {
	return nil
}
func (m *mockNodeManager) Watch(ctx context.Context, opts *types.WatchNodesOptions) (<-chan types.NodeEvent, error) {
	if m.watchFunc != nil {
		return m.watchFunc(ctx, opts)
	}
	return nil, nil
}

type mockPartitionManager struct {
	watchFunc func(ctx context.Context, opts *types.WatchPartitionsOptions) (<-chan types.PartitionEvent, error)
}

func (m *mockPartitionManager) List(ctx context.Context, opts *types.ListPartitionsOptions) (*types.PartitionList, error) {
	return nil, nil
}
func (m *mockPartitionManager) Get(ctx context.Context, partitionName string) (*types.Partition, error) {
	return nil, nil
}
func (m *mockPartitionManager) Create(ctx context.Context, partition *types.PartitionCreate) (*types.PartitionCreateResponse, error) {
	return nil, nil
}
func (m *mockPartitionManager) Update(ctx context.Context, partitionName string, update *types.PartitionUpdate) error {
	return nil
}
func (m *mockPartitionManager) Delete(ctx context.Context, partitionName string) error {
	return nil
}
func (m *mockPartitionManager) Watch(ctx context.Context, opts *types.WatchPartitionsOptions) (<-chan types.PartitionEvent, error) {
	if m.watchFunc != nil {
		return m.watchFunc(ctx, opts)
	}
	return nil, nil
}
