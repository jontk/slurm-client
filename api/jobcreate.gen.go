// Code generated from OpenAPI spec. DO NOT EDIT.
// SPDX-FileCopyrightText: 2026 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package api

// JobCreate represents a SLURM JobCreate.
type JobCreate struct {
	Account *string `json:"account,omitempty"` // Account associated with the job
	AccountGatherFrequency *string `json:"account_gather_frequency,omitempty"` // Job accounting and profiling sampling intervals in seconds
	AdminComment *string `json:"admin_comment,omitempty"` // Arbitrary comment made by administrator
	AllocationNodeList *string `json:"allocation_node_list,omitempty"` // Local node making the resource allocation
	AllocationNodePort *int32 `json:"allocation_node_port,omitempty"` // Port to send allocation confirmation to
	Argv []string `json:"argv,omitempty"` // Arguments to the script. Note: The slurmstepd always overrides argv[0] with the...
	Array *string `json:"array,omitempty"` // Job array index value specification
	BatchFeatures *string `json:"batch_features,omitempty"` // Features required for batch script's node
	BeginTime *uint64 `json:"begin_time,omitempty"` // Defer the allocation of the job until the specified time (UNIX timestamp) (UNIX...
	BurstBuffer *string `json:"burst_buffer,omitempty"` // Burst buffer specifications
	ClusterConstraint *string `json:"cluster_constraint,omitempty"` // Required features that a federated cluster must have to have a sibling job...
	Clusters *string `json:"clusters,omitempty"` // Clusters that a federated job can run on
	Comment *string `json:"comment,omitempty"` // Arbitrary comment made by user
	Constraints *string `json:"constraints,omitempty"` // Comma-separated list of features that are required
	Container *string `json:"container,omitempty"` // Absolute path to OCI container bundle
	ContainerID *string `json:"container_id,omitempty"` // OCI container ID
	Contiguous *bool `json:"contiguous,omitempty"` // True if job requires contiguous nodes
	CoreSpecification *int32 `json:"core_specification,omitempty"` // Specialized core count
	CPUBinding *string `json:"cpu_binding,omitempty"` // Method for binding tasks to allocated CPUs
	CPUBindingFlags []CPUBindingFlagsValue `json:"cpu_binding_flags,omitempty"` // Flags for CPU binding
	CPUFrequency *string `json:"cpu_frequency,omitempty"` // Requested CPU frequency range <p1>[-p2][:p3]
	CPUsPerTask *int32 `json:"cpus_per_task,omitempty"` // Number of CPUs required by each task
	CPUsPerTRES *string `json:"cpus_per_tres,omitempty"` // Semicolon delimited list of TRES=# values values indicating how many CPUs...
	Crontab *CronEntry `json:"crontab,omitempty"` // Specification for scrontab job (crontab entry)
	CurrentWorkingDirectory *string `json:"current_working_directory,omitempty"` // Working directory to use for the job
	Deadline *int64 `json:"deadline,omitempty"` // Latest time that the job may start (UNIX timestamp) (UNIX timestamp or time...
	DelayBoot *int32 `json:"delay_boot,omitempty"` // Number of seconds after job eligible start that nodes will be rebooted to...
	Dependency *string `json:"dependency,omitempty"` // Other jobs that must meet certain criteria before this job can start
	Distribution *string `json:"distribution,omitempty"` // Layout
	DistributionPlaneSize *uint16 `json:"distribution_plane_size,omitempty"` // Plane size specification when distribution specifies plane (16 bit integer...
	EndTime *int64 `json:"end_time,omitempty"` // Expected end time (UNIX timestamp) (UNIX timestamp or time string recognized by...
	Environment []string `json:"environment,omitempty"` // Environment variables to be set for the job
	ExcludedNodes []string `json:"excluded_nodes,omitempty"` // Comma-separated list of nodes that may not be used
	Extra *string `json:"extra,omitempty"` // Arbitrary string used for node filtering if extra constraints are enabled
	Flags []FlagsValue `json:"flags,omitempty"` // Job flags
	GroupID *string `json:"group_id,omitempty"` // Group ID of the user that owns the job
	HetjobGroup *int32 `json:"hetjob_group,omitempty"` // Unique sequence number applied to this component of the heterogeneous job
	Hold *bool `json:"hold,omitempty"` // Hold (true) or release (false) job (Job held)
	Immediate *bool `json:"immediate,omitempty"` // If true, exit if resources are not available within the time period specified
	JobID *int32 `json:"job_id,omitempty"` // Job ID
	KillOnNodeFail *bool `json:"kill_on_node_fail,omitempty"` // If true, kill job on node failure
	KillWarningDelay *uint16 `json:"kill_warning_delay,omitempty"` // Number of seconds before end time to send the warning signal (16 bit integer...
	KillWarningFlags []KillWarningFlagsValue `json:"kill_warning_flags,omitempty"` // Flags related to job signals
	KillWarningSignal *string `json:"kill_warning_signal,omitempty"` // Signal to send when approaching end time (e.g. "10" or "USR1")
	Licenses *string `json:"licenses,omitempty"` // License(s) required by the job
	MailType []MailTypeValue `json:"mail_type,omitempty"` // Mail event type(s)
	MailUser *string `json:"mail_user,omitempty"` // User to receive email notifications
	MaximumCPUs *int32 `json:"maximum_cpus,omitempty"` // Maximum number of CPUs required
	MaximumNodes *int32 `json:"maximum_nodes,omitempty"` // Maximum node count
	MCSLabel *string `json:"mcs_label,omitempty"` // Multi-Category Security label on the job
	MemoryBinding *string `json:"memory_binding,omitempty"` // Binding map for map/mask_cpu
	MemoryBindingType []MemoryBindingTypeValue `json:"memory_binding_type,omitempty"` // Method for binding tasks to memory
	MemoryPerCPU *uint64 `json:"memory_per_cpu,omitempty"` // Minimum memory in megabytes per allocated CPU (64 bit integer number with flags)
	MemoryPerNode *uint64 `json:"memory_per_node,omitempty"` // Minimum memory in megabytes per allocated node (64 bit integer number with...
	MemoryPerTRES *string `json:"memory_per_tres,omitempty"` // Semicolon delimited list of TRES=# values indicating how much memory in...
	MinimumBoardsPerNode *int32 `json:"minimum_boards_per_node,omitempty"` // Boards per node required
	MinimumCPUs *int32 `json:"minimum_cpus,omitempty"` // Minimum number of CPUs required
	MinimumCPUsPerNode *int32 `json:"minimum_cpus_per_node,omitempty"` // Minimum number of CPUs per node
	MinimumNodes *int32 `json:"minimum_nodes,omitempty"` // Minimum node count
	MinimumSocketsPerBoard *int32 `json:"minimum_sockets_per_board,omitempty"` // Sockets per board required
	Name *string `json:"name,omitempty"` // Job name
	Network *string `json:"network,omitempty"` // Network specs for job step
	Nice *int32 `json:"nice,omitempty"` // Requested job priority change
	Nodes *string `json:"nodes,omitempty"` // Node count range specification (e.g. 1-15:4)
	NtasksPerTRES *int32 `json:"ntasks_per_tres,omitempty"` // Number of tasks that can access each GPU
	OomKillStep *int32 `json:"oom_kill_step,omitempty"` // Kill whole step in case of OOM in one of the tasks
	OpenMode []OpenModeValue `json:"open_mode,omitempty"` // Open mode used for stdout and stderr files
	Overcommit *bool `json:"overcommit,omitempty"` // Overcommit resources
	Partition *string `json:"partition,omitempty"` // Partition assigned to the job
	PowerFlags []interface{} `json:"power_flags,omitempty"`
	Prefer *string `json:"prefer,omitempty"` // Comma-separated list of features that are preferred but not required
	Priority *uint32 `json:"priority,omitempty"` // Request specific job priority (32 bit integer number with flags)
	Profile []ProfileValue `json:"profile,omitempty"` // Profile used by the acct_gather_profile plugin
	QoS *string `json:"qos,omitempty"` // Quality of Service assigned to the job
	Reboot *bool `json:"reboot,omitempty"` // Node reboot requested before start
	Requeue *bool `json:"requeue,omitempty"` // Determines whether the job may be requeued
	RequiredNodes []string `json:"required_nodes,omitempty"` // Comma-separated list of required nodes
	RequiredSwitches *uint32 `json:"required_switches,omitempty"` // Maximum number of switches (32 bit integer number with flags)
	Reservation *string `json:"reservation,omitempty"` // Name of reservation to use
	ReservePorts *int32 `json:"reserve_ports,omitempty"` // Port to send various notification msg to
	Rlimits *JobCreateRlimits `json:"rlimits,omitempty"`
	Script *string `json:"script,omitempty"` // Job batch script contents; only the first component in a HetJob is populated or...
	SegmentSize *uint16 `json:"segment_size,omitempty"` // Segment size for topology/block (16 bit integer number with flags)
	SelinuxContext *string `json:"selinux_context,omitempty"` // SELinux context
	Shared []SharedValue `json:"shared,omitempty"` // How the job can share resources with other jobs, if at all
	SiteFactor *int32 `json:"site_factor,omitempty"` // Site-specific priority factor
	SocketsPerNode *int32 `json:"sockets_per_node,omitempty"` // Sockets per node required
	SpankEnvironment []string `json:"spank_environment,omitempty"` // Environment variables for job prolog/epilog scripts as set by SPANK plugins
	StandardError *string `json:"standard_error,omitempty"` // Path to stderr file
	StandardInput *string `json:"standard_input,omitempty"` // Path to stdin file
	StandardOutput *string `json:"standard_output,omitempty"` // Path to stdout file
	StepID *StepID `json:"step_id,omitempty"` // Job step ID
	Tasks *int32 `json:"tasks,omitempty"` // Number of tasks
	TasksPerBoard *int32 `json:"tasks_per_board,omitempty"` // Number of tasks to invoke on each board
	TasksPerCore *int32 `json:"tasks_per_core,omitempty"` // Number of tasks to invoke on each core
	TasksPerNode *int32 `json:"tasks_per_node,omitempty"` // Number of tasks to invoke on each node
	TasksPerSocket *int32 `json:"tasks_per_socket,omitempty"` // Number of tasks to invoke on each socket
	TemporaryDiskPerNode *int32 `json:"temporary_disk_per_node,omitempty"` // Minimum tmp disk space required per node
	ThreadSpecification *int32 `json:"thread_specification,omitempty"` // Specialized thread count
	ThreadsPerCore *int32 `json:"threads_per_core,omitempty"` // Threads per core required
	TimeLimit *uint32 `json:"time_limit,omitempty"` // Maximum run time in minutes (32 bit integer number with flags)
	TimeMinimum *uint32 `json:"time_minimum,omitempty"` // Minimum run time in minutes (32 bit integer number with flags)
	TRESBind *string `json:"tres_bind,omitempty"` // Task to TRES binding directives
	TRESFreq *string `json:"tres_freq,omitempty"` // TRES frequency directives
	TRESPerJob *string `json:"tres_per_job,omitempty"` // Comma-separated list of TRES=# values to be allocated for every job
	TRESPerNode *string `json:"tres_per_node,omitempty"` // Comma-separated list of TRES=# values to be allocated for every node
	TRESPerSocket *string `json:"tres_per_socket,omitempty"` // Comma-separated list of TRES=# values to be allocated for every socket
	TRESPerTask *string `json:"tres_per_task,omitempty"` // Comma-separated list of TRES=# values to be allocated for every task
	UserID *string `json:"user_id,omitempty"` // User ID that owns the job
	WaitAllNodes *bool `json:"wait_all_nodes,omitempty"` // If true, wait to start until after all nodes have booted
	WaitForSwitch *int32 `json:"wait_for_switch,omitempty"` // Maximum time to wait for switches in seconds
	Wckey *string `json:"wckey,omitempty"` // Workload characterization key
	X11 []X11Value `json:"x11,omitempty"` // X11 forwarding options
	X11MagicCookie *string `json:"x11_magic_cookie,omitempty"` // Magic cookie for X11 forwarding
	X11TargetHost *string `json:"x11_target_host,omitempty"` // Hostname or UNIX socket if x11_target_port=0
	X11TargetPort *int32 `json:"x11_target_port,omitempty"` // TCP port
}


// JobCreateRlimits is a nested type within its parent.
type JobCreateRlimits struct {
	As *uint64 `json:"as,omitempty"` // Address space limit (Address space limit.) (64 bit integer number with flags)
	Core *uint64 `json:"core,omitempty"` // Largest core file that can be created, in bytes (Largest core file that can be...
	CPU *uint64 `json:"cpu,omitempty"` // Per-process CPU limit, in seconds (Per-process CPU limit, in seconds.) (64 bit...
	Data *uint64 `json:"data,omitempty"` // Maximum size of data segment, in bytes (Maximum size of data segment, in bytes....
	Fsize *uint64 `json:"fsize,omitempty"` // Largest file that can be created, in bytes (Largest file that can be created,...
	Memlock *uint64 `json:"memlock,omitempty"` // Locked-in-memory address space (Locked-in-memory address space) (64 bit integer...
	Nofile *uint64 `json:"nofile,omitempty"` // Number of open files (Number of open files.) (64 bit integer number with flags)
	Nproc *uint64 `json:"nproc,omitempty"` // Number of processes (Number of processes.) (64 bit integer number with flags)
	Rss *uint64 `json:"rss,omitempty"` // Largest resident set size, in bytes. This affects swapping; processes that are...
	Stack *uint64 `json:"stack,omitempty"` // Maximum size of stack segment, in bytes (Maximum size of stack segment, in...
}


// CPUBindingFlagsValue represents possible values for CPUBindingFlags field.
type CPUBindingFlagsValue string

// CPUBindingFlagsValue constants.
const (
	CPUBindingFlagsCPUBindToThreads CPUBindingFlagsValue = "CPU_BIND_TO_THREADS"
	CPUBindingFlagsCPUBindToCores CPUBindingFlagsValue = "CPU_BIND_TO_CORES"
	CPUBindingFlagsCPUBindToSockets CPUBindingFlagsValue = "CPU_BIND_TO_SOCKETS"
	CPUBindingFlagsCPUBindToLdoms CPUBindingFlagsValue = "CPU_BIND_TO_LDOMS"
	CPUBindingFlagsCPUBindNone CPUBindingFlagsValue = "CPU_BIND_NONE"
	CPUBindingFlagsCPUBindRank CPUBindingFlagsValue = "CPU_BIND_RANK"
	CPUBindingFlagsCPUBindMap CPUBindingFlagsValue = "CPU_BIND_MAP"
	CPUBindingFlagsCPUBindMask CPUBindingFlagsValue = "CPU_BIND_MASK"
	CPUBindingFlagsCPUBindLdrank CPUBindingFlagsValue = "CPU_BIND_LDRANK"
	CPUBindingFlagsCPUBindLdmap CPUBindingFlagsValue = "CPU_BIND_LDMAP"
	CPUBindingFlagsCPUBindLdmask CPUBindingFlagsValue = "CPU_BIND_LDMASK"
	CPUBindingFlagsVerbose CPUBindingFlagsValue = "VERBOSE"
	CPUBindingFlagsCPUBindOneThreadPerCore CPUBindingFlagsValue = "CPU_BIND_ONE_THREAD_PER_CORE"
)

// KillWarningFlagsValue represents possible values for KillWarningFlags field.
type KillWarningFlagsValue string

// KillWarningFlagsValue constants.
const (
	KillWarningFlagsBatchJob KillWarningFlagsValue = "BATCH_JOB"
	KillWarningFlagsArrayTask KillWarningFlagsValue = "ARRAY_TASK"
	KillWarningFlagsFullStepsOnly KillWarningFlagsValue = "FULL_STEPS_ONLY"
	KillWarningFlagsFullJob KillWarningFlagsValue = "FULL_JOB"
	KillWarningFlagsFederationRequeue KillWarningFlagsValue = "FEDERATION_REQUEUE"
	KillWarningFlagsHurry KillWarningFlagsValue = "HURRY"
	KillWarningFlagsOutOfMemory KillWarningFlagsValue = "OUT_OF_MEMORY"
	KillWarningFlagsNoSiblingJobs KillWarningFlagsValue = "NO_SIBLING_JOBS"
	KillWarningFlagsReservationJob KillWarningFlagsValue = "RESERVATION_JOB"
	KillWarningFlagsVerbose KillWarningFlagsValue = "VERBOSE"
	KillWarningFlagsCronJobs KillWarningFlagsValue = "CRON_JOBS"
	KillWarningFlagsWarningSent KillWarningFlagsValue = "WARNING_SENT"
)

// MemoryBindingTypeValue represents possible values for MemoryBindingType field.
type MemoryBindingTypeValue string

// MemoryBindingTypeValue constants.
const (
	MemoryBindingTypeNone MemoryBindingTypeValue = "NONE"
	MemoryBindingTypeRank MemoryBindingTypeValue = "RANK"
	MemoryBindingTypeMap MemoryBindingTypeValue = "MAP"
	MemoryBindingTypeMask MemoryBindingTypeValue = "MASK"
	MemoryBindingTypeLocal MemoryBindingTypeValue = "LOCAL"
	MemoryBindingTypeVerbose MemoryBindingTypeValue = "VERBOSE"
	MemoryBindingTypePrefer MemoryBindingTypeValue = "PREFER"
)

// OpenModeValue represents possible values for OpenMode field.
type OpenModeValue string

// OpenModeValue constants.
const (
	OpenModeAppend OpenModeValue = "APPEND"
	OpenModeTruncate OpenModeValue = "TRUNCATE"
)


// X11Value represents possible values for X11 field.
type X11Value string

// X11Value constants.
const (
	X11ForwardAllNodes X11Value = "FORWARD_ALL_NODES"
	X11BatchNode X11Value = "BATCH_NODE"
	X11FirstNode X11Value = "FIRST_NODE"
	X11LastNode X11Value = "LAST_NODE"
)
