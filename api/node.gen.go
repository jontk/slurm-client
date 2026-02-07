// Code generated from OpenAPI spec. DO NOT EDIT.
// SPDX-FileCopyrightText: 2026 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package api

import "time"

// Node represents a SLURM Node.
type Node struct {
	ActiveFeatures []string `json:"active_features,omitempty"` // Currently active features
	Address *string `json:"address,omitempty"` // NodeAddr, used to establish a communication path
	AllocCPUs *int32 `json:"alloc_cpus,omitempty"` // Total number of CPUs currently allocated for jobs
	AllocIdleCPUs *int32 `json:"alloc_idle_cpus,omitempty"` // Total number of idle CPUs
	AllocMemory *int64 `json:"alloc_memory,omitempty"` // Total memory in MB currently allocated for jobs
	Architecture *string `json:"architecture,omitempty"` // Computer architecture
	Boards *int32 `json:"boards,omitempty"` // Number of Baseboards in nodes with a baseboard controller
	BootTime time.Time `json:"boot_time,omitempty"` // Time when the node booted (UNIX timestamp) (UNIX timestamp or time string...
	BurstbufferNetworkAddress *string `json:"burstbuffer_network_address,omitempty"` // Alternate network path to be used for sbcast network traffic
	CertFlags []CertFlagsValue `json:"cert_flags,omitempty"` // Certmgr status flags
	ClusterName *string `json:"cluster_name,omitempty"` // Cluster name (only set in federated environments)
	Comment *string `json:"comment,omitempty"` // Arbitrary comment
	Cores *int32 `json:"cores,omitempty"` // Number of cores in a single physical processor socket
	CPUBinding *int32 `json:"cpu_binding,omitempty"` // Default method for binding tasks to allocated CPUs
	CPULoad *int32 `json:"cpu_load,omitempty"` // CPU load as reported by the OS
	CPUs *int32 `json:"cpus,omitempty"` // Total CPUs, including cores and threads
	EffectiveCPUs *int32 `json:"effective_cpus,omitempty"` // Number of effective CPUs (excluding specialized CPUs)
	Energy *NodeEnergy `json:"energy,omitempty"` // Energy usage data
	ExternalSensors map[string]interface{} `json:"external_sensors,omitempty"`
	Extra *string `json:"extra,omitempty"` // Arbitrary string used for node filtering if extra constraints are enabled
	Features []string `json:"features,omitempty"` // Available features
	FreeMem *uint64 `json:"free_mem,omitempty"` // Total memory in MB currently free as reported by the OS (64 bit integer number...
	GPUSpec *string `json:"gpu_spec,omitempty"` // CPU cores reserved for jobs that also use a GPU
	GRES *string `json:"gres,omitempty"` // Generic resources
	GRESDrained *string `json:"gres_drained,omitempty"` // Drained generic resources
	GRESUsed *string `json:"gres_used,omitempty"` // Generic resources currently in use
	Hostname *string `json:"hostname,omitempty"` // NodeHostname
	InstanceID *string `json:"instance_id,omitempty"` // Cloud instance ID
	InstanceType *string `json:"instance_type,omitempty"` // Cloud instance type
	LastBusy time.Time `json:"last_busy,omitempty"` // Time when the node was last busy (UNIX timestamp) (UNIX timestamp or time...
	MCSLabel *string `json:"mcs_label,omitempty"` // Multi-Category Security label
	Name *string `json:"name,omitempty"` // NodeName
	NextStateAfterReboot []NodeState `json:"next_state_after_reboot,omitempty"` // The state the node will be assigned after rebooting
	OperatingSystem *string `json:"operating_system,omitempty"` // Operating system reported by the node
	Owner *string `json:"owner,omitempty"` // User allowed to run jobs on this node (unset if no restriction)
	Partitions []string `json:"partitions,omitempty"` // Partitions containing this node
	Port *int32 `json:"port,omitempty"` // TCP port number of the slurmd
	Power map[string]interface{} `json:"power,omitempty"`
	RealMemory *int64 `json:"real_memory,omitempty"` // Total memory in MB on the node
	Reason *string `json:"reason,omitempty"` // Describes why the node is in a "DOWN", "DRAINED", "DRAINING", "FAILING" or...
	ReasonChangedAt time.Time `json:"reason_changed_at,omitempty"` // When the reason changed (UNIX timestamp) (UNIX timestamp or time string...
	ReasonSetByUser *string `json:"reason_set_by_user,omitempty"` // User who set the reason
	ResCoresPerGPU *int32 `json:"res_cores_per_gpu,omitempty"` // Number of CPU cores per GPU restricted to GPU jobs
	Reservation *string `json:"reservation,omitempty"` // Name of reservation containing this node
	ResumeAfter *uint64 `json:"resume_after,omitempty"` // Number of seconds after the node's state is updated to "DOWN" or "DRAIN" before...
	SlurmdStartTime time.Time `json:"slurmd_start_time,omitempty"` // Time when the slurmd started (UNIX timestamp) (UNIX timestamp or time string...
	Sockets *int32 `json:"sockets,omitempty"` // Number of physical processor sockets/chips on the node
	SpecializedCores *int32 `json:"specialized_cores,omitempty"` // Number of cores reserved for system use
	SpecializedCPUs *string `json:"specialized_cpus,omitempty"` // Abstract CPU IDs on this node reserved for exclusive use by slurmd and...
	SpecializedMemory *int64 `json:"specialized_memory,omitempty"` // Combined memory limit, in MB, for Slurm compute node daemons
	State []NodeState `json:"state,omitempty"` // Node state(s) applicable to this node
	TemporaryDisk *int32 `json:"temporary_disk,omitempty"` // Total size in MB of temporary disk storage in TmpFS
	Threads *int32 `json:"threads,omitempty"` // Number of logical threads in a single physical core
	TLSCertLastRenewal time.Time `json:"tls_cert_last_renewal,omitempty"` // Time when TLS certificate was created (UNIX timestamp or time string recognized...
	Topology *string `json:"topology,omitempty"` // Topology
	TRES *string `json:"tres,omitempty"` // Configured trackable resources
	TRESUsed *string `json:"tres_used,omitempty"` // Trackable resources currently allocated for jobs
	TRESWeighted *float64 `json:"tres_weighted,omitempty"` // Ignored. Was weighted number of billable trackable resources allocated
	Version *string `json:"version,omitempty"` // Slurmd version
	Weight *int32 `json:"weight,omitempty"` // Weight of the node for scheduling purposes
}


// CertFlagsValue represents possible values for CertFlags field.
type CertFlagsValue string

// CertFlagsValue constants.
const (
	CertFlagsTokenSet CertFlagsValue = "TOKEN_SET"
)

// NodeState represents possible values for NodeState field.
type NodeState string

// NodeState constants.
const (
	NodeStateInvalid NodeState = "INVALID"
	NodeStateUnknown NodeState = "UNKNOWN"
	NodeStateDown NodeState = "DOWN"
	NodeStateIdle NodeState = "IDLE"
	NodeStateAllocated NodeState = "ALLOCATED"
	NodeStateError NodeState = "ERROR"
	NodeStateMixed NodeState = "MIXED"
	NodeStateFuture NodeState = "FUTURE"
	NodeStateExternal NodeState = "EXTERNAL"
	NodeStateReserved NodeState = "RESERVED"
	NodeStateUndrain NodeState = "UNDRAIN"
	NodeStateCloud NodeState = "CLOUD"
	NodeStateResume NodeState = "RESUME"
	NodeStateDrain NodeState = "DRAIN"
	NodeStateCompleting NodeState = "COMPLETING"
	NodeStateNotResponding NodeState = "NOT_RESPONDING"
	NodeStatePoweredDown NodeState = "POWERED_DOWN"
	NodeStateFail NodeState = "FAIL"
	NodeStatePoweringUp NodeState = "POWERING_UP"
	NodeStateMaintenance NodeState = "MAINTENANCE"
	NodeStateRebootRequested NodeState = "REBOOT_REQUESTED"
	NodeStateRebootCanceled NodeState = "REBOOT_CANCELED"
	NodeStatePoweringDown NodeState = "POWERING_DOWN"
	NodeStateDynamicFuture NodeState = "DYNAMIC_FUTURE"
	NodeStateRebootIssued NodeState = "REBOOT_ISSUED"
	NodeStatePlanned NodeState = "PLANNED"
	NodeStateInvalidReg NodeState = "INVALID_REG"
	NodeStatePowerDown NodeState = "POWER_DOWN"
	NodeStatePowerUp NodeState = "POWER_UP"
	NodeStatePowerDrain NodeState = "POWER_DRAIN"
	NodeStateDynamicNorm NodeState = "DYNAMIC_NORM"
	NodeStateBlocked NodeState = "BLOCKED"
)
