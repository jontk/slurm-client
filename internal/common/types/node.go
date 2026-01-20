// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"time"
)

// Node represents a SLURM compute node with common fields across all API versions
type Node struct {
	Name                      string      `json:"name"`
	Arch                      string      `json:"arch,omitempty"`
	BcastAddress              string      `json:"bcast_address,omitempty"`
	Boards                    int32       `json:"boards,omitempty"`
	BootTime                  *time.Time  `json:"boot_time,omitempty"`
	ClusterName               string      `json:"cluster_name,omitempty"`
	Comment                   string      `json:"comment,omitempty"`
	Cores                     int32       `json:"cores,omitempty"`
	CoreSpecCount             int32       `json:"core_spec_count,omitempty"`
	CPUBinding                int32       `json:"cpu_binding,omitempty"`
	CPULoad                   float64     `json:"cpu_load,omitempty"`
	CPUs                      int32       `json:"cpus"`
	CPUsEffective             int32       `json:"cpus_effective,omitempty"`
	Features                  []string    `json:"features,omitempty"`
	ActiveFeatures            []string    `json:"active_features,omitempty"`
	FreeMemory                int64       `json:"free_memory,omitempty"`
	Gres                      string      `json:"gres,omitempty"`
	GresDrained               string      `json:"gres_drained,omitempty"`
	GresUsed                  string      `json:"gres_used,omitempty"`
	LastBusy                  *time.Time  `json:"last_busy,omitempty"`
	MCSLabel                  string      `json:"mcs_label,omitempty"`
	Memory                    int64       `json:"memory"`
	MemorySpecLimit           int64       `json:"memory_spec_limit,omitempty"`
	NextStateAfterReboot      NodeState   `json:"next_state_after_reboot,omitempty"`
	NextStateAfterRebootFlags []string    `json:"next_state_after_reboot_flags,omitempty"`
	NodeAddress               string      `json:"node_address,omitempty"`
	NodeHostname              string      `json:"node_hostname,omitempty"`
	OS                        string      `json:"os,omitempty"`
	Owner                     string      `json:"owner,omitempty"`
	Partitions                []string    `json:"partitions,omitempty"`
	Port                      int32       `json:"port,omitempty"`
	RealMemory                int64       `json:"real_memory,omitempty"`
	Reason                    string      `json:"reason,omitempty"`
	ReasonChangedAt           *time.Time  `json:"reason_changed_at,omitempty"`
	ReasonSetByUser           string      `json:"reason_set_by_user,omitempty"`
	ResumeAfter               *time.Time  `json:"resume_after,omitempty"`
	ResvMemory                int64       `json:"resv_memory,omitempty"`
	SlurmdStartTime           *time.Time  `json:"slurmd_start_time,omitempty"`
	Sockets                   int32       `json:"sockets,omitempty"`
	State                     NodeState   `json:"state"`
	StateFlags                []string    `json:"state_flags,omitempty"`
	StateReason               string      `json:"state_reason,omitempty"`
	ThreadsPerCore            int32       `json:"threads_per_core,omitempty"`
	TmpDisk                   int64       `json:"tmp_disk,omitempty"`
	TresUsed                  string      `json:"tres_used,omitempty"`
	TresFmtStr                string      `json:"tres_fmt_str,omitempty"`
	Version                   string      `json:"version,omitempty"`
	Weight                    int32       `json:"weight,omitempty"`
	AllocCPUs                 int32       `json:"alloc_cpus,omitempty"`
	AllocIdleCPUs             int32       `json:"alloc_idle_cpus,omitempty"`
	AllocMemory               int64       `json:"alloc_memory,omitempty"`
	Energy                    *NodeEnergy `json:"energy,omitempty"`
}

// NodeState represents the state of a node
type NodeState string

const (
	NodeStateUnknown      NodeState = "UNKNOWN"
	NodeStateDown         NodeState = "DOWN"
	NodeStateIdle         NodeState = "IDLE"
	NodeStateAllocated    NodeState = "ALLOCATED"
	NodeStateError        NodeState = "ERROR"
	NodeStateMixed        NodeState = "MIXED"
	NodeStateFuture       NodeState = "FUTURE"
	NodeStateReserved     NodeState = "RESERVED"
	NodeStateUndrained    NodeState = "UNDRAINED"
	NodeStateCloud        NodeState = "CLOUD"
	NodeStateDraining     NodeState = "DRAINING"
	NodeStateDrained      NodeState = "DRAINED"
	NodeStateResuming     NodeState = "RESUMING"
	NodeStateFail         NodeState = "FAIL"
	NodeStateFailing      NodeState = "FAILING"
	NodeStateMaintenance  NodeState = "MAINTENANCE"
	NodeStateRebooting    NodeState = "REBOOTING"
	NodeStateCancelling   NodeState = "CANCELLING"
	NodeStatePoweredDown  NodeState = "POWERED_DOWN"
	NodeStatePoweringDown NodeState = "POWERING_DOWN"
	NodeStatePoweringUp   NodeState = "POWERING_UP"
	NodeStatePlanned      NodeState = "PLANNED"
)

// NodeEnergy represents energy consumption data for a node
type NodeEnergy struct {
	AveWatts           int64     `json:"ave_watts,omitempty"`
	BaseConsumedEnergy int64     `json:"base_consumed_energy,omitempty"`
	ConsumedEnergy     int64     `json:"consumed_energy,omitempty"`
	CurrentWatts       int64     `json:"current_watts,omitempty"`
	ExtSensorsJoules   int64     `json:"ext_sensors_joules,omitempty"`
	ExtSensorsWatts    int64     `json:"ext_sensors_watts,omitempty"`
	ExtSensorsTemp     int32     `json:"ext_sensors_temp,omitempty"`
	LastCollected      time.Time `json:"last_collected,omitempty"`
}

// NodeUpdate represents the data needed to update a node
type NodeUpdate struct {
	Comment              *string           `json:"comment,omitempty"`
	CPUBinding           *int32            `json:"cpu_binding,omitempty"`
	Features             []string          `json:"features,omitempty"`
	ActiveFeatures       []string          `json:"active_features,omitempty"`
	Gres                 *string           `json:"gres,omitempty"`
	NextStateAfterReboot *NodeState        `json:"next_state_after_reboot,omitempty"`
	Reason               *string           `json:"reason,omitempty"`
	ResumeAfter          *time.Time        `json:"resume_after,omitempty"`
	State                *NodeState        `json:"state,omitempty"`
	Weight               *int32            `json:"weight,omitempty"`
	Extra                map[string]string `json:"extra,omitempty"`
}

// NodeListOptions represents options for listing nodes
type NodeListOptions struct {
	Names      []string    `json:"names,omitempty"`
	States     []NodeState `json:"states,omitempty"`
	Partitions []string    `json:"partitions,omitempty"`
	UpdateTime *time.Time  `json:"update_time,omitempty"`
	Reasons    []string    `json:"reasons,omitempty"`
	Limit      int         `json:"limit,omitempty"`
	Offset     int         `json:"offset,omitempty"`
}

// NodeList represents a list of nodes
type NodeList struct {
	Nodes []Node `json:"nodes"`
	Total int    `json:"total"`
}

// NodeAllocation represents the current allocation status of a node
type NodeAllocation struct {
	Name        string    `json:"name"`
	JobID       int32     `json:"job_id,omitempty"`
	JobName     string    `json:"job_name,omitempty"`
	UserName    string    `json:"user_name,omitempty"`
	AllocCPUs   int32     `json:"alloc_cpus"`
	AllocMemory int64     `json:"alloc_memory"`
	AllocGres   string    `json:"alloc_gres,omitempty"`
	StartTime   time.Time `json:"start_time,omitempty"`
	EndTime     time.Time `json:"end_time,omitempty"`
}

// NodeMaintenanceRequest represents a request to put nodes into maintenance
type NodeMaintenanceRequest struct {
	Nodes         []string   `json:"nodes"`
	Reason        string     `json:"reason"`
	StartTime     *time.Time `json:"start_time,omitempty"`
	EndTime       *time.Time `json:"end_time,omitempty"`
	FixedDuration int32      `json:"fixed_duration,omitempty"`
}

// NodePowerRequest represents a request to change node power state
type NodePowerRequest struct {
	Nodes        []string       `json:"nodes"`
	PowerState   NodePowerState `json:"power_state"`
	Asynchronous bool           `json:"asynchronous,omitempty"`
	Force        bool           `json:"force,omitempty"`
}

// NodePowerState represents node power states
type NodePowerState string

const (
	NodePowerDown NodePowerState = "POWER_DOWN"
	NodePowerUp   NodePowerState = "POWER_UP"
	NodePowerSave NodePowerState = "POWER_SAVE"
)
