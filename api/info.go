// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package api

// ClusterInfo represents cluster information (mirrors interfaces.ClusterInfo)
type ClusterInfo struct {
	Version     string `json:"version"`
	Release     string `json:"release"`
	ClusterName string `json:"cluster_name"`
	APIVersion  string `json:"api_version"`
	Uptime      int    `json:"uptime"`
}

// ClusterStats represents cluster statistics (mirrors interfaces.ClusterStats)
type ClusterStats struct {
	TotalNodes     int `json:"total_nodes"`
	IdleNodes      int `json:"idle_nodes"`
	AllocatedNodes int `json:"allocated_nodes"`
	TotalCPUs      int `json:"total_cpus"`
	IdleCPUs       int `json:"idle_cpus"`
	AllocatedCPUs  int `json:"allocated_cpus"`
	TotalJobs      int `json:"total_jobs"`
	RunningJobs    int `json:"running_jobs"`
	PendingJobs    int `json:"pending_jobs"`
	CompletedJobs  int `json:"completed_jobs"`
}

// APIVersion represents API version information (mirrors interfaces.APIVersion)
type APIVersion struct {
	Version     string `json:"version"`
	Release     string `json:"release"`
	Description string `json:"description"`
	Deprecated  bool   `json:"deprecated"`
}
