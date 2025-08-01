package fixtures

import (
	"github.com/jontk/slurm-client/internal/common/types"
	"time"
)

// QoSFixtures provides pre-configured QoS objects for testing
type QoSFixtures struct{}

// NewQoSFixtures creates a new QoS fixtures instance
func NewQoSFixtures() *QoSFixtures {
	return &QoSFixtures{}
}

// SimpleQoS returns a basic QoS with minimal fields
func (f *QoSFixtures) SimpleQoS(name string) *types.QoS {
	return &types.QoS{
		ID:          1,
		Name:        name,
		Description: "Test QoS " + name,
		Priority:    100,
		UsageFactor: 1.0,
		CreatedAt:   time.Now().Add(-24 * time.Hour),
		ModifiedAt:  time.Now(),
	}
}

// QoSWithLimits returns a QoS with comprehensive limits
func (f *QoSFixtures) QoSWithLimits(name string) *types.QoS {
	qos := f.SimpleQoS(name)
	qos.Limits = &types.QoSLimits{
		MaxCPUsPerUser:      intPtr(100),
		MaxJobsPerUser:      intPtr(10),
		MaxNodesPerUser:     intPtr(5),
		MaxSubmitJobsPerUser: intPtr(20),
		MaxCPUsPerAccount:   intPtr(1000),
		MaxJobsPerAccount:   intPtr(100),
		MaxNodesPerAccount:  intPtr(50),
		MaxCPUsPerJob:       intPtr(32),
		MaxNodesPerJob:      intPtr(2),
		MaxWallTimePerJob:   intPtr(1440), // 24 hours
		MaxMemoryPerNode:    int64Ptr(64000), // 64GB
		MaxMemoryPerCPU:     int64Ptr(4000),  // 4GB
	}
	return qos
}

// HighPriorityQoS returns a high-priority QoS configuration
func (f *QoSFixtures) HighPriorityQoS() *types.QoS {
	qos := f.QoSWithLimits("high-priority")
	qos.Priority = 1000
	qos.Flags = []string{"DenyOnLimit", "RequiresReservation"}
	qos.PreemptMode = "cluster"
	qos.PreemptExemptTime = 300 // 5 minutes
	qos.UsageFactor = 2.0
	qos.UsageThreshold = 0.95
	return qos
}

// BatchQoS returns a QoS suitable for batch jobs
func (f *QoSFixtures) BatchQoS() *types.QoS {
	qos := f.QoSWithLimits("batch")
	qos.Priority = 10
	qos.Flags = []string{"NoReserve"}
	qos.UsageFactor = 0.5
	qos.GraceTime = 3600 // 1 hour
	// More relaxed limits for batch
	qos.Limits.MaxWallTimePerJob = intPtr(10080) // 7 days
	qos.Limits.MaxJobsPerUser = intPtr(100)
	return qos
}

// QoSCreateRequest returns a QoSCreate request for testing
func (f *QoSFixtures) QoSCreateRequest(name string) *types.QoSCreate {
	return &types.QoSCreate{
		Name:           name,
		Description:    "Test QoS creation",
		Priority:       100,
		Flags:          []string{"DenyOnLimit"},
		PreemptMode:    []string{"cluster"},
		GraceTime:      300, // Changed to non-pointer
		UsageFactor:    1.5,
		UsageThreshold: 0.8,
		Limits: &types.QoSLimits{
			MaxCPUsPerUser: intPtr(50),
			MaxJobsPerUser: intPtr(5),
		},
	}
}

// QoSUpdateRequest returns a QoSUpdate request for testing
func (f *QoSFixtures) QoSUpdateRequest() *types.QoSUpdate {
	return &types.QoSUpdate{
		Description:    stringPtr("Updated description"),
		Priority:       intPtr(200),
		Flags:          &[]string{"DenyOnLimit", "RequiresReservation"},
		PreemptMode:    &[]string{"suspend"},
		UsageFactor:    float64Ptr(2.0),
		UsageThreshold: float64Ptr(0.9),
		Limits: &types.QoSLimits{
			MaxCPUsPerUser: intPtr(100),
			MaxJobsPerUser: intPtr(10),
		},
	}
}

// QoSList returns a list of diverse QoS entries
func (f *QoSFixtures) QoSList() []types.QoS {
	return []types.QoS{
		*f.SimpleQoS("normal"),
		*f.HighPriorityQoS(),
		*f.BatchQoS(),
		*f.QoSWithLimits("gpu-jobs"),
		*f.SimpleQoS("low-priority"),
	}
}

// Helper functions
func intPtr(i int) *int {
	return &i
}

func int64Ptr(i int64) *int64 {
	return &i
}

func stringPtr(s string) *string {
	return &s
}

func float64Ptr(f float64) *float64 {
	return &f
}