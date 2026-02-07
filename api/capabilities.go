// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package api

// ClientCapabilities describes the features supported by a specific API version.
// Callers should check these capabilities before invoking optional features
// to avoid runtime errors.
type ClientCapabilities struct {
	// Version is the API version string (e.g., "v0.0.41")
	Version string

	// Resource Manager Support
	SupportsJobs         bool
	SupportsNodes        bool
	SupportsPartitions   bool
	SupportsReservations bool

	// Database Manager Support
	SupportsAccounts     bool
	SupportsUsers        bool
	SupportsQoS          bool
	SupportsClusters     bool
	SupportsAssociations bool
	SupportsWCKeys       bool

	// Write Operations Support
	SupportsJobSubmit        bool
	SupportsJobUpdate        bool
	SupportsJobCancel        bool
	SupportsNodeUpdate       bool
	SupportsPartitionWrite   bool
	SupportsReservationWrite bool

	// Database Write Operations
	SupportsAccountWrite     bool
	SupportsUserWrite        bool
	SupportsQoSWrite         bool
	SupportsClusterWrite     bool
	SupportsAssociationWrite bool
	SupportsWCKeyWrite       bool

	// Advanced Features
	SupportsTRES        bool
	SupportsInstances   bool
	SupportsReconfigure bool
	SupportsDiagnostics bool
	SupportsShares      bool
	SupportsLicenses    bool

	// Extended Features (not implemented in all versions)
	// IMPORTANT: Check these capabilities before calling the corresponding methods.
	// Methods return "not implemented" error when capability is false.
	SupportsJobSteps       bool // Jobs().Steps() method - returns job step info
	SupportsJobWatch       bool // Jobs().Watch() method - real-time job events
	SupportsNodeWatch      bool // Nodes().Watch() method - real-time node events
	SupportsPartitionWatch bool // Partitions().Watch() method - real-time partition events
	SupportsAnalytics      bool // Analytics() returns non-nil manager (computed insights, NOT part of SLURM REST API)

	// Extended Account/User Operations (PLANNED - NOT YET IMPLEMENTED)
	// These helper methods require database queries beyond the base adapter.
	// All currently return "not implemented" error regardless of version.
	// See examples/user-account-management for usage patterns when implemented.
	SupportsAccountHierarchy bool // GetAccountHierarchy, GetParentAccounts, GetChildAccounts
	SupportsAccountQuotas    bool // GetAccountQuotas, GetAccountQuotaUsage
	SupportsUserHelpers      bool // GetUserAccounts, GetUserQuotas, GetUserDefaultAccount
	SupportsFairShare        bool // GetAccountFairShare, GetUserFairShare, GetFairShareHierarchy

	// Cluster Operations (limited in adapter pattern)
	SupportsClusterCreate bool // Cluster().Create()
	SupportsClusterUpdate bool // Cluster().Update()
	SupportsClusterDelete bool // Cluster().Delete()

	// Bulk Operations
	SupportsAssociationBulkDelete bool // Associations().BulkDelete()
}
