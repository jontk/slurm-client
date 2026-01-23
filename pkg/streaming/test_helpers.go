// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package streaming

import (
	"context"

	"github.com/jontk/slurm-client/interfaces"
)

// Mock implementations for testing

type mockSlurmClient struct {
	jobs       *mockJobManager
	nodes      *mockNodeManager
	partitions *mockPartitionManager
}

func (m *mockSlurmClient) Version() string { return "test" }
func (m *mockSlurmClient) Jobs() interfaces.JobManager {
	return m.jobs
}
func (m *mockSlurmClient) Nodes() interfaces.NodeManager {
	return m.nodes
}
func (m *mockSlurmClient) Partitions() interfaces.PartitionManager {
	return m.partitions
}
func (m *mockSlurmClient) Info() interfaces.InfoManager                           { return nil }
func (m *mockSlurmClient) Reservations() interfaces.ReservationManager            { return nil }
func (m *mockSlurmClient) QoS() interfaces.QoSManager                             { return nil }
func (m *mockSlurmClient) Accounts() interfaces.AccountManager                    { return nil }
func (m *mockSlurmClient) Users() interfaces.UserManager                          { return nil }
func (m *mockSlurmClient) Clusters() interfaces.ClusterManager                    { return nil }
func (m *mockSlurmClient) Associations() interfaces.AssociationManager            { return nil }
func (m *mockSlurmClient) WCKeys() interfaces.WCKeyManager                        { return nil }
func (m *mockSlurmClient) GetLicenses(ctx context.Context) (*interfaces.LicenseList, error) {
	return nil, nil
}
func (m *mockSlurmClient) GetShares(ctx context.Context, opts *interfaces.GetSharesOptions) (*interfaces.SharesList, error) {
	return nil, nil
}
func (m *mockSlurmClient) GetConfig(ctx context.Context) (*interfaces.Config, error) {
	return nil, nil
}
func (m *mockSlurmClient) GetDiagnostics(ctx context.Context) (*interfaces.Diagnostics, error) {
	return nil, nil
}
func (m *mockSlurmClient) GetDBDiagnostics(ctx context.Context) (*interfaces.Diagnostics, error) {
	return nil, nil
}
func (m *mockSlurmClient) GetInstance(ctx context.Context, opts *interfaces.GetInstanceOptions) (*interfaces.Instance, error) {
	return nil, nil
}
func (m *mockSlurmClient) GetInstances(ctx context.Context, opts *interfaces.GetInstancesOptions) (*interfaces.InstanceList, error) {
	return nil, nil
}
func (m *mockSlurmClient) GetTRES(ctx context.Context) (*interfaces.TRESList, error) {
	return nil, nil
}
func (m *mockSlurmClient) CreateTRES(ctx context.Context, req *interfaces.CreateTRESRequest) (*interfaces.TRES, error) {
	return nil, nil
}
func (m *mockSlurmClient) Reconfigure(ctx context.Context) (*interfaces.ReconfigureResponse, error) {
	return nil, nil
}
func (m *mockSlurmClient) Close() error { return nil }

type mockJobManager struct {
	watchFunc func(ctx context.Context, opts *interfaces.WatchJobsOptions) (<-chan interfaces.JobEvent, error)
}

func (m *mockJobManager) List(ctx context.Context, opts *interfaces.ListJobsOptions) (*interfaces.JobList, error) {
	return nil, nil
}
func (m *mockJobManager) Get(ctx context.Context, jobID string) (*interfaces.Job, error) {
	return nil, nil
}
func (m *mockJobManager) Submit(ctx context.Context, job *interfaces.JobSubmission) (*interfaces.JobSubmitResponse, error) {
	return nil, nil
}
func (m *mockJobManager) Allocate(ctx context.Context, req *interfaces.JobAllocateRequest) (*interfaces.JobAllocateResponse, error) {
	return nil, nil
}
func (m *mockJobManager) Cancel(ctx context.Context, jobID string) error { return nil }
func (m *mockJobManager) Requeue(ctx context.Context, jobID string) error { return nil }
func (m *mockJobManager) Update(ctx context.Context, jobID string, update *interfaces.JobUpdate) error {
	return nil
}
func (m *mockJobManager) Steps(ctx context.Context, jobID string) (*interfaces.JobStepList, error) {
	return nil, nil
}
func (m *mockJobManager) Watch(ctx context.Context, opts *interfaces.WatchJobsOptions) (<-chan interfaces.JobEvent, error) {
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
func (m *mockJobManager) GetJobUtilization(ctx context.Context, jobID string) (*interfaces.JobUtilization, error) {
	return nil, nil
}
func (m *mockJobManager) GetJobEfficiency(ctx context.Context, jobID string) (*interfaces.ResourceUtilization, error) {
	return nil, nil
}
func (m *mockJobManager) GetJobPerformance(ctx context.Context, jobID string) (*interfaces.JobPerformance, error) {
	return nil, nil
}
func (m *mockJobManager) GetJobLiveMetrics(ctx context.Context, jobID string) (*interfaces.JobLiveMetrics, error) {
	return nil, nil
}
func (m *mockJobManager) WatchJobMetrics(ctx context.Context, jobID string, opts *interfaces.WatchMetricsOptions) (<-chan interfaces.JobMetricsEvent, error) {
	return nil, nil
}
func (m *mockJobManager) GetJobResourceTrends(ctx context.Context, jobID string, opts *interfaces.ResourceTrendsOptions) (*interfaces.JobResourceTrends, error) {
	return nil, nil
}
func (m *mockJobManager) GetJobStepDetails(ctx context.Context, jobID string, stepID string) (*interfaces.JobStepDetails, error) {
	return nil, nil
}
func (m *mockJobManager) GetJobStepUtilization(ctx context.Context, jobID string, stepID string) (*interfaces.JobStepUtilization, error) {
	return nil, nil
}
func (m *mockJobManager) ListJobStepsWithMetrics(ctx context.Context, jobID string, opts *interfaces.ListJobStepsOptions) (*interfaces.JobStepMetricsList, error) {
	return nil, nil
}
func (m *mockJobManager) GetJobStepsFromAccounting(ctx context.Context, jobID string, opts *interfaces.AccountingQueryOptions) (*interfaces.AccountingJobSteps, error) {
	return nil, nil
}
func (m *mockJobManager) GetStepAccountingData(ctx context.Context, jobID string, stepID string) (*interfaces.StepAccountingRecord, error) {
	return nil, nil
}
func (m *mockJobManager) GetJobStepAPIData(ctx context.Context, jobID string, stepID string) (*interfaces.JobStepAPIData, error) {
	return nil, nil
}
func (m *mockJobManager) ListJobStepsFromSacct(ctx context.Context, jobID string, opts *interfaces.SacctQueryOptions) (*interfaces.SacctJobStepData, error) {
	return nil, nil
}
func (m *mockJobManager) GetJobCPUAnalytics(ctx context.Context, jobID string) (*interfaces.CPUAnalytics, error) {
	return nil, nil
}
func (m *mockJobManager) GetJobMemoryAnalytics(ctx context.Context, jobID string) (*interfaces.MemoryAnalytics, error) {
	return nil, nil
}
func (m *mockJobManager) GetJobIOAnalytics(ctx context.Context, jobID string) (*interfaces.IOAnalytics, error) {
	return nil, nil
}
func (m *mockJobManager) GetJobComprehensiveAnalytics(ctx context.Context, jobID string) (*interfaces.JobComprehensiveAnalytics, error) {
	return nil, nil
}
func (m *mockJobManager) GetJobPerformanceHistory(ctx context.Context, jobID string, opts *interfaces.PerformanceHistoryOptions) (*interfaces.JobPerformanceHistory, error) {
	return nil, nil
}
func (m *mockJobManager) GetPerformanceTrends(ctx context.Context, opts *interfaces.TrendAnalysisOptions) (*interfaces.PerformanceTrends, error) {
	return nil, nil
}
func (m *mockJobManager) GetUserEfficiencyTrends(ctx context.Context, userID string, opts *interfaces.EfficiencyTrendOptions) (*interfaces.UserEfficiencyTrends, error) {
	return nil, nil
}
func (m *mockJobManager) AnalyzeBatchJobs(ctx context.Context, jobIDs []string, opts *interfaces.BatchAnalysisOptions) (*interfaces.BatchJobAnalysis, error) {
	return nil, nil
}
func (m *mockJobManager) GetWorkflowPerformance(ctx context.Context, workflowID string, opts *interfaces.WorkflowAnalysisOptions) (*interfaces.WorkflowPerformance, error) {
	return nil, nil
}
func (m *mockJobManager) GenerateEfficiencyReport(ctx context.Context, opts *interfaces.ReportOptions) (*interfaces.EfficiencyReport, error) {
	return nil, nil
}

type mockNodeManager struct {
	watchFunc func(ctx context.Context, opts *interfaces.WatchNodesOptions) (<-chan interfaces.NodeEvent, error)
}

func (m *mockNodeManager) List(ctx context.Context, opts *interfaces.ListNodesOptions) (*interfaces.NodeList, error) {
	return nil, nil
}
func (m *mockNodeManager) Get(ctx context.Context, nodeName string) (*interfaces.Node, error) {
	return nil, nil
}
func (m *mockNodeManager) Update(ctx context.Context, nodeName string, update *interfaces.NodeUpdate) error {
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
func (m *mockNodeManager) Watch(ctx context.Context, opts *interfaces.WatchNodesOptions) (<-chan interfaces.NodeEvent, error) {
	if m.watchFunc != nil {
		return m.watchFunc(ctx, opts)
	}
	return nil, nil
}

type mockPartitionManager struct {
	watchFunc func(ctx context.Context, opts *interfaces.WatchPartitionsOptions) (<-chan interfaces.PartitionEvent, error)
}

func (m *mockPartitionManager) List(ctx context.Context, opts *interfaces.ListPartitionsOptions) (*interfaces.PartitionList, error) {
	return nil, nil
}
func (m *mockPartitionManager) Get(ctx context.Context, partitionName string) (*interfaces.Partition, error) {
	return nil, nil
}
func (m *mockPartitionManager) Create(ctx context.Context, partition *interfaces.PartitionCreate) (*interfaces.PartitionCreateResponse, error) {
	return nil, nil
}
func (m *mockPartitionManager) Update(ctx context.Context, partitionName string, update *interfaces.PartitionUpdate) error {
	return nil
}
func (m *mockPartitionManager) Delete(ctx context.Context, partitionName string) error {
	return nil
}
func (m *mockPartitionManager) Watch(ctx context.Context, opts *interfaces.WatchPartitionsOptions) (<-chan interfaces.PartitionEvent, error) {
	if m.watchFunc != nil {
		return m.watchFunc(ctx, opts)
	}
	return nil, nil
}
