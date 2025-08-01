package types

import "time"

// QoS represents a Quality of Service configuration in SLURM
// This is a version-agnostic representation used across all API versions
type QoS struct {
	// Basic identification
	ID          int32
	Name        string
	Description string

	// Priority and scheduling
	Priority       int
	UsageFactor    float64
	UsageThreshold float64

	// Preemption settings
	PreemptMode      string
	PreemptExemptTime int // in minutes

	// Flags
	Flags []string

	// Time limits
	GraceTime int // in seconds

	// Hierarchy
	ParentQoS string

	// Resource limits
	Limits *QoSLimits

	// TRES limits
	MaxTRESPerUser    string
	MaxTRESPerAccount string
	MaxTRESPerJob     string

	// Allowed entities
	AllowedAccounts []string
	AllowedUsers    []string

	// Metadata
	CreatedAt  time.Time
	ModifiedAt time.Time
}

// QoSLimits represents resource limits for a QoS
type QoSLimits struct {
	// Per-user limits
	MaxCPUsPerUser     *int
	MaxJobsPerUser     *int
	MaxNodesPerUser    *int
	MaxSubmitJobsPerUser *int

	// Per-account limits
	MaxCPUsPerAccount  *int
	MaxJobsPerAccount  *int
	MaxNodesPerAccount *int

	// Per-job limits
	MaxCPUsPerJob     *int
	MaxNodesPerJob    *int
	MaxWallTimePerJob *int // in minutes

	// Memory limits (in MB)
	MaxMemoryPerNode *int64
	MaxMemoryPerCPU  *int64

	// Other limits
	MaxBurstBuffer *int64 // in bytes
	MinCPUsPerJob  *int
	MinNodesPerJob *int
}

// QoSCreate represents the data needed to create a new QoS
type QoSCreate struct {
	Name              string
	Description       string
	Priority          int
	Flags             []string
	PreemptMode       []string
	PreemptList       []string
	PreemptExemptTime *int
	GraceTime         int // Changed to non-pointer for validation
	UsageFactor       float64
	UsageThreshold    float64
	ParentQoS         string
	MaxTRESPerUser    string
	MaxTRESPerAccount string
	MaxTRESPerJob     string
	Limits            *QoSLimits
}

// QoSUpdate represents fields that can be updated on a QoS
type QoSUpdate struct {
	Description       *string
	Priority          *int
	Flags             *[]string
	PreemptMode       *[]string
	PreemptList       []string
	PreemptExemptTime *int
	GraceTime         *int
	UsageFactor       *float64
	UsageThreshold    *float64
	ParentQoS         *string
	MaxTRESPerUser    *string
	MaxTRESPerAccount *string
	MaxTRESPerJob     *string
	Limits            *QoSLimits
}

// QoSListOptions represents options for listing QoS entries
type QoSListOptions struct {
	Names    []string
	Accounts []string
	Users    []string
	Limit    int
	Offset   int
}

// QoSList represents a list of QoS entries
type QoSList struct {
	QoS   []QoS
	Total int
}

// QoSCreateResponse represents the response from creating a QoS
type QoSCreateResponse struct {
	QoSName string
}