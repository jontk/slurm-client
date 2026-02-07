// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

// This file re-exports interface types from pkg/types.
// Concrete types (structs) are defined in types.go.
// This follows Go community standards for a single public package.

package slurm

import (
	types "github.com/jontk/slurm-client/api"
)

// ============================================================================
// Core Client Interface
// ============================================================================

// SlurmClient represents a version-agnostic Slurm REST API client
type SlurmClient = types.SlurmClient

// ============================================================================
// Job Interfaces
// ============================================================================

// JobReader provides read-only job operations
type JobReader = types.JobReader

// JobWriter provides job mutation operations
type JobWriter = types.JobWriter

// JobController provides job control operations
type JobController = types.JobController

// JobWatcher provides real-time job operations
type JobWatcher = types.JobWatcher

// JobManager combines all core job operations
type JobManager = types.JobManager

// ============================================================================
// Resource Manager Interfaces
// ============================================================================

// NodeManager provides node operations
type NodeManager = types.NodeManager

// PartitionManager provides partition operations
type PartitionManager = types.PartitionManager

// InfoManager provides cluster info operations
type InfoManager = types.InfoManager

// ReservationManager provides reservation operations
type ReservationManager = types.ReservationManager

// QoSManager provides QoS operations
type QoSManager = types.QoSManager

// AccountManager provides account operations
type AccountManager = types.AccountManager

// UserManager provides user operations
type UserManager = types.UserManager

// ClusterManager provides cluster operations
type ClusterManager = types.ClusterManager

// AssociationManager provides association operations
type AssociationManager = types.AssociationManager

// WCKeyManager provides WCKey operations
type WCKeyManager = types.WCKeyManager

// ============================================================================
// Analytics Interface
// ============================================================================

// AnalyticsManager provides advanced performance analytics
type AnalyticsManager = types.AnalyticsManager
