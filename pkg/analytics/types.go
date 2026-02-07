// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

// Package analytics provides performance analytics types for HPC job analysis.
// These types are SDK-specific and provide value-added analytics not in the SLURM REST API.
//
// Most types in this package are aliases to types in github.com/jontk/slurm-client/api
// to maintain backward compatibility while eliminating duplication.
package analytics

import "github.com/jontk/slurm-client/api"

// ============================================================================
// Type Aliases to api/analytics_types.go
// These maintain backward compatibility for existing code using pkg/analytics
// ============================================================================

// Core Analytics Types
type JobUtilization = api.JobUtilization
type ResourceUtilization = api.ResourceUtilization

// GPU Analytics
type GPUUtilization = api.GPUUtilization
type GPUDeviceUtilization = api.GPUDeviceUtilization
type GPUProcess = api.GPUProcess

// I/O Analytics
type IOUtilization = api.IOUtilization
type IOStats = api.IOStats

// Network Analytics
type NetworkUtilization = api.NetworkUtilization
type NetworkInterfaceStats = api.NetworkInterfaceStats

// Energy Analytics
type EnergyUsage = api.EnergyUsage

// Performance Analytics
type JobPerformance = api.JobPerformance
type JobStepPerformance = api.JobStepPerformance
type PerformanceBottleneck = api.PerformanceBottleneck

// Live Metrics
type JobLiveMetrics = api.JobLiveMetrics
type LiveResourceMetric = api.LiveResourceMetric
type NodeLiveMetrics = api.NodeLiveMetrics
type PerformanceAlert = api.PerformanceAlert

// Watch Options and Events
type WatchMetricsOptions = api.WatchMetricsOptions
type JobMetricsEvent = api.JobMetricsEvent
type JobStateChange = api.JobStateChange

// Resource Trends
type ResourceTrendsOptions = api.ResourceTrendsOptions
type JobResourceTrends = api.JobResourceTrends
type ResourceTimeSeries = api.ResourceTimeSeries
type IOTimeSeries = api.IOTimeSeries
type NetworkTimeSeries = api.NetworkTimeSeries
type EnergyTimeSeries = api.EnergyTimeSeries
type ResourceAnomaly = api.ResourceAnomaly
type TrendsSummary = api.TrendsSummary

// Job Step Analytics
type JobStepDetails = api.JobStepDetails
type JobStepUtilization = api.JobStepUtilization
type TaskUtilization = api.TaskUtilization
type StepTaskInfo = api.StepTaskInfo
type StepPerformanceMetrics = api.StepPerformanceMetrics
type ListJobStepsOptions = api.ListJobStepsOptions
type JobStepMetricsList = api.JobStepMetricsList
type JobStepWithMetrics = api.JobStepWithMetrics
type JobStepsSummary = api.JobStepsSummary
type StepResourceTrends = api.StepResourceTrends
type ResourceTrendData = api.ResourceTrendData
type StepComparison = api.StepComparison
type StepOptimizationSuggestions = api.StepOptimizationSuggestions
type OptimizationSuggestion = api.OptimizationSuggestion

// Accounting Data Types
type AccountingQueryOptions = api.AccountingQueryOptions
type AccountingJobSteps = api.AccountingJobSteps
type StepAccountingRecord = api.StepAccountingRecord
type JobAccountingSummary = api.JobAccountingSummary
type JobStepAPIData = api.JobStepAPIData
type ProcessInfo = api.ProcessInfo
type SacctQueryOptions = api.SacctQueryOptions
type SacctJobStepData = api.SacctJobStepData
type SacctStepRecord = api.SacctStepRecord

// Comprehensive Analytics
type CPUAnalytics = api.CPUAnalytics
type CPUCoreMetric = api.CPUCoreMetric
type MemoryAnalytics = api.MemoryAnalytics
type NUMANodeMetrics = api.NUMANodeMetrics
type MemoryLeak = api.MemoryLeak
type IOAnalytics = api.IOAnalytics
type StorageDevice = api.StorageDevice
type JobComprehensiveAnalytics = api.JobComprehensiveAnalytics
type CrossResourceAnalysis = api.CrossResourceAnalysis
type OptimalJobConfiguration = api.OptimalJobConfiguration

// Performance History and Trends
type PerformanceHistoryOptions = api.PerformanceHistoryOptions
type JobPerformanceHistory = api.JobPerformanceHistory
type PerformanceSnapshot = api.PerformanceSnapshot
type PerformanceTrendAnalysis = api.PerformanceTrendAnalysis
type TrendInfo = api.TrendInfo
type PerformanceAnomaly = api.PerformanceAnomaly
type TrendAnalysisOptions = api.TrendAnalysisOptions
type PerformanceTrends = api.PerformanceTrends
type TimeRange = api.TimeRange
type UtilizationPoint = api.UtilizationPoint
type EfficiencyPoint = api.EfficiencyPoint
type PartitionTrend = api.PartitionTrend
type ResourceTrend = api.ResourceTrend
type JobSizeTrend = api.JobSizeTrend
type JobDurationTrend = api.JobDurationTrend
type JobCountPoint = api.JobCountPoint
type QueueLengthPoint = api.QueueLengthPoint
type TrendInsight = api.TrendInsight

// User Efficiency Trends
type EfficiencyTrendOptions = api.EfficiencyTrendOptions
type UserEfficiencyTrends = api.UserEfficiencyTrends
type EfficiencyDataPoint = api.EfficiencyDataPoint

// Batch Analysis
type BatchAnalysisOptions = api.BatchAnalysisOptions
type BatchJobAnalysis = api.BatchJobAnalysis
type BatchStatistics = api.BatchStatistics
type JobAnalysisSummary = api.JobAnalysisSummary
type BatchComparison = api.BatchComparison
type BatchPattern = api.BatchPattern
type BatchRecommendation = api.BatchRecommendation
type ResourceWaste = api.ResourceWaste

// Workflow Analysis
type WorkflowAnalysisOptions = api.WorkflowAnalysisOptions
type WorkflowPerformance = api.WorkflowPerformance
type WorkflowStage = api.WorkflowStage
type WorkflowBottleneck = api.WorkflowBottleneck
type WorkflowDependencies = api.WorkflowDependencies
type WorkflowOptimization = api.WorkflowOptimization
type ResourceChange = api.ResourceChange

// Efficiency Reports
type ReportOptions = api.ReportOptions
type EfficiencyReport = api.EfficiencyReport
type ExecutiveSummary = api.ExecutiveSummary
type ClusterOverview = api.ClusterOverview
type PartitionAnalysis = api.PartitionAnalysis
type UserAnalysis = api.UserAnalysis
type ResourceAnalysis = api.ResourceAnalysis
type ResourceTypeAnalysis = api.ResourceTypeAnalysis
type ReportTrendAnalysis = api.ReportTrendAnalysis
type ReportRecommendation = api.ReportRecommendation
type ChartData = api.ChartData

// Extended Diagnostics
type ExtendedDiagnostics = api.ExtendedDiagnostics
type AccountUsageStats = api.AccountUsageStats
