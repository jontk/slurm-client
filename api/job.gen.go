// Code generated from OpenAPI spec. DO NOT EDIT.
// SPDX-FileCopyrightText: 2026 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package api

import "time"

// Job represents a SLURM Job.
type Job struct {
	Account *string `json:"account,omitempty"` // Account associated with the job
	AccrueTime time.Time `json:"accrue_time,omitempty"` // When the job started accruing age priority (UNIX timestamp) (UNIX timestamp or...
	AdminComment *string `json:"admin_comment,omitempty"` // Arbitrary comment made by administrator
	AllocatingNode *string `json:"allocating_node,omitempty"` // Local node making the resource allocation
	ArrayJobID *uint32 `json:"array_job_id,omitempty"` // Job ID of job array, or 0 if N/A (32 bit integer number with flags)
	ArrayMaxTasks *uint32 `json:"array_max_tasks,omitempty"` // Maximum number of simultaneously running array tasks, 0 if no limit (32 bit...
	ArrayTaskID *uint32 `json:"array_task_id,omitempty"` // Task ID of this task in job array (32 bit integer number with flags)
	ArrayTaskString *string `json:"array_task_string,omitempty"` // String expression of task IDs in this record
	AssociationID *int32 `json:"association_id,omitempty"` // Unique identifier for the association
	BatchFeatures *string `json:"batch_features,omitempty"` // Features required for batch script's node
	BatchFlag *bool `json:"batch_flag,omitempty"` // True if batch job
	BatchHost *string `json:"batch_host,omitempty"` // Name of host running batch script
	BillableTRES *float64 `json:"billable_tres,omitempty"` // Billable TRES (64 bit floating point number with flags)
	BurstBuffer *string `json:"burst_buffer,omitempty"` // Burst buffer specifications
	BurstBufferState *string `json:"burst_buffer_state,omitempty"` // Burst buffer state details
	Cluster *string `json:"cluster,omitempty"` // Cluster name
	ClusterFeatures *string `json:"cluster_features,omitempty"` // List of required cluster features
	Command *string `json:"command,omitempty"` // Executed command
	Comment *string `json:"comment,omitempty"` // Arbitrary comment
	Container *string `json:"container,omitempty"` // Absolute path to OCI container bundle
	ContainerID *string `json:"container_id,omitempty"` // OCI container ID
	Contiguous *bool `json:"contiguous,omitempty"` // True if job requires contiguous nodes
	CoreSpec *int32 `json:"core_spec,omitempty"` // Specialized core count
	CoresPerSocket *uint16 `json:"cores_per_socket,omitempty"` // Cores per socket required (16 bit integer number with flags)
	CPUFrequencyGovernor *uint32 `json:"cpu_frequency_governor,omitempty"` // CPU frequency governor (32 bit integer number with flags)
	CPUFrequencyMaximum *uint32 `json:"cpu_frequency_maximum,omitempty"` // Maximum CPU frequency (32 bit integer number with flags)
	CPUFrequencyMinimum *uint32 `json:"cpu_frequency_minimum,omitempty"` // Minimum CPU frequency (32 bit integer number with flags)
	CPUs *uint32 `json:"cpus,omitempty"` // Minimum number of CPUs required (32 bit integer number with flags)
	CPUsPerTask *uint16 `json:"cpus_per_task,omitempty"` // Number of CPUs required by each task (16 bit integer number with flags)
	CPUsPerTRES *string `json:"cpus_per_tres,omitempty"` // Semicolon delimited list of TRES=# values indicating how many CPUs should be...
	Cron *string `json:"cron,omitempty"` // Time specification for scrontab job
	CurrentWorkingDirectory *string `json:"current_working_directory,omitempty"` // Working directory to use for the job
	Deadline time.Time `json:"deadline,omitempty"` // Latest time that the job may start (UNIX timestamp) (UNIX timestamp or time...
	DelayBoot *uint32 `json:"delay_boot,omitempty"` // Number of seconds after job eligible start that nodes will be rebooted to...
	Dependency *string `json:"dependency,omitempty"` // Other jobs that must meet certain criteria before this job can start
	DerivedExitCode *ExitCode `json:"derived_exit_code,omitempty"` // Highest exit code of all job steps (return code returned by process)
	EligibleTime time.Time `json:"eligible_time,omitempty"` // Time when the job became eligible to run (UNIX timestamp) (UNIX timestamp or...
	EndTime time.Time `json:"end_time,omitempty"` // End time, real or expected (UNIX timestamp) (UNIX timestamp or time string...
	ExcludedNodes *string `json:"excluded_nodes,omitempty"` // Comma-separated list of nodes that may not be used
	ExitCode *ExitCode `json:"exit_code,omitempty"` // Exit code of the job (return code returned by process)
	Extra *string `json:"extra,omitempty"` // Arbitrary string used for node filtering if extra constraints are enabled
	FailedNode *string `json:"failed_node,omitempty"` // Name of node that caused job failure
	Features *string `json:"features,omitempty"` // Comma-separated list of features that are required
	FederationOrigin *string `json:"federation_origin,omitempty"` // Origin cluster's name (when using federation)
	FederationSiblingsActive *string `json:"federation_siblings_active,omitempty"` // Active sibling job names
	FederationSiblingsViable *string `json:"federation_siblings_viable,omitempty"` // Viable sibling job names
	Flags []FlagsValue `json:"flags,omitempty"` // Job flags
	GRESDetail []string `json:"gres_detail,omitempty"` // List of GRES index and counts allocated per node
	GroupID *int32 `json:"group_id,omitempty"` // Group ID of the user that owns the job
	GroupName *string `json:"group_name,omitempty"` // Group name of the user that owns the job
	HetJobID *uint32 `json:"het_job_id,omitempty"` // Heterogeneous job ID, if applicable (32 bit integer number with flags)
	HetJobIDSet *string `json:"het_job_id_set,omitempty"` // Job ID range for all heterogeneous job components
	HetJobOffset *uint32 `json:"het_job_offset,omitempty"` // Unique sequence number applied to this component of the heterogeneous job (32...
	Hold *bool `json:"hold,omitempty"` // Hold (true) or release (false) job (Job held)
	JobID *int32 `json:"job_id,omitempty"` // Job ID
	JobResources *JobResources `json:"job_resources,omitempty"` // Resources used by the job
	JobSizeStr []string `json:"job_size_str,omitempty"` // Number of nodes (in a range) required for this job
	JobState []JobState `json:"job_state,omitempty"` // Current state
	LastSchedEvaluation time.Time `json:"last_sched_evaluation,omitempty"` // Last time job was evaluated for scheduling (UNIX timestamp) (UNIX timestamp or...
	Licenses *string `json:"licenses,omitempty"` // License(s) required by the job
	LicensesAllocated *string `json:"licenses_allocated,omitempty"` // License(s) allocated to the job
	MailType []MailTypeValue `json:"mail_type,omitempty"` // Mail event type(s)
	MailUser *string `json:"mail_user,omitempty"` // User to receive email notifications
	MaxCPUs *uint32 `json:"max_cpus,omitempty"` // Maximum number of CPUs usable by the job (32 bit integer number with flags)
	MaxNodes *uint32 `json:"max_nodes,omitempty"` // Maximum number of nodes usable by the job (32 bit integer number with flags)
	MaximumSwitchWaitTime *int32 `json:"maximum_switch_wait_time,omitempty"` // Maximum time to wait for switches in seconds
	MCSLabel *string `json:"mcs_label,omitempty"` // Multi-Category Security label on the job
	MemoryPerCPU *uint64 `json:"memory_per_cpu,omitempty"` // Minimum memory in megabytes per allocated CPU (64 bit integer number with flags)
	MemoryPerNode *uint64 `json:"memory_per_node,omitempty"` // Minimum memory in megabytes per allocated node (64 bit integer number with...
	MemoryPerTRES *string `json:"memory_per_tres,omitempty"` // Semicolon delimited list of TRES=# values indicating how much memory in...
	MinimumCPUsPerNode *uint16 `json:"minimum_cpus_per_node,omitempty"` // Minimum number of CPUs per node (16 bit integer number with flags)
	MinimumTmpDiskPerNode *uint32 `json:"minimum_tmp_disk_per_node,omitempty"` // Minimum tmp disk space required per node (32 bit integer number with flags)
	Name *string `json:"name,omitempty"` // Job name
	Network *string `json:"network,omitempty"` // Network specs for the job
	Nice *int32 `json:"nice,omitempty"` // Requested job priority change
	NodeCount *uint32 `json:"node_count,omitempty"` // Minimum number of nodes required (32 bit integer number with flags)
	Nodes *string `json:"nodes,omitempty"` // Node(s) allocated to the job
	Partition *string `json:"partition,omitempty"` // Partition assigned to the job
	Power *JobPower `json:"power,omitempty"`
	PreSusTime *uint64 `json:"pre_sus_time,omitempty"` // Total run time prior to last suspend in seconds (UNIX timestamp or time string...
	PreemptTime time.Time `json:"preempt_time,omitempty"` // Time job received preemption signal (UNIX timestamp) (UNIX timestamp or time...
	PreemptableTime time.Time `json:"preemptable_time,omitempty"` // Time job becomes eligible for preemption (UNIX timestamp) (UNIX timestamp or...
	Prefer *string `json:"prefer,omitempty"` // Feature(s) the job requested but that are not required
	Priority *uint32 `json:"priority,omitempty"` // Request specific job priority (32 bit integer number with flags)
	PriorityByPartition []JobPartitionPriority `json:"priority_by_partition,omitempty"` // Prospective job priority in each partition that may be used by this job
	Profile []ProfileValue `json:"profile,omitempty"` // Profile used by the acct_gather_profile plugin
	QoS *string `json:"qos,omitempty"` // Quality of Service assigned to the job, if pending the QOS requested
	Reboot *bool `json:"reboot,omitempty"` // Node reboot requested before start
	Requeue *bool `json:"requeue,omitempty"` // Determines whether the job may be requeued
	RequiredNodes *string `json:"required_nodes,omitempty"` // Comma-separated list of required nodes
	RequiredSwitches *int32 `json:"required_switches,omitempty"` // Maximum number of switches
	ResizeTime time.Time `json:"resize_time,omitempty"` // Time of last size change (UNIX timestamp) (UNIX timestamp or time string...
	RestartCnt *int32 `json:"restart_cnt,omitempty"` // Number of job restarts
	ResvName *string `json:"resv_name,omitempty"` // Name of reservation to use
	ScheduledNodes *string `json:"scheduled_nodes,omitempty"` // List of nodes scheduled to be used for the job
	SegmentSize *int32 `json:"segment_size,omitempty"` // Requested segment size
	SelinuxContext *string `json:"selinux_context,omitempty"` // SELinux context
	Shared []SharedValue `json:"shared,omitempty"` // How the job can share resources with other jobs, if at all
	SocketsPerBoard *int32 `json:"sockets_per_board,omitempty"` // Number of sockets per board required
	SocketsPerNode *uint16 `json:"sockets_per_node,omitempty"` // Number of sockets per node required (16 bit integer number with flags)
	StandardError *string `json:"standard_error,omitempty"` // Path to stderr file
	StandardInput *string `json:"standard_input,omitempty"` // Path to stdin file
	StandardOutput *string `json:"standard_output,omitempty"` // Path to stdout file
	StartTime time.Time `json:"start_time,omitempty"` // Time execution began, or is expected to begin (UNIX timestamp) (UNIX timestamp...
	StateDescription *string `json:"state_description,omitempty"` // Optional details for state_reason
	StateReason *string `json:"state_reason,omitempty"` // Reason for current Pending or Failed state
	StderrExpanded *string `json:"stderr_expanded,omitempty"` // Job stderr with expanded fields
	StdinExpanded *string `json:"stdin_expanded,omitempty"` // Job stdin with expanded fields
	StdoutExpanded *string `json:"stdout_expanded,omitempty"` // Job stdout with expanded fields
	StepID *StepID `json:"step_id,omitempty"` // Job step ID
	SubmitLine *string `json:"submit_line,omitempty"` // Job submit line (e.g. 'sbatch -N3 job.sh job_arg'
	SubmitTime time.Time `json:"submit_time,omitempty"` // Time when the job was submitted (UNIX timestamp) (UNIX timestamp or time string...
	SuspendTime time.Time `json:"suspend_time,omitempty"` // Time the job was last suspended or resumed (UNIX timestamp) (UNIX timestamp or...
	SystemComment *string `json:"system_comment,omitempty"` // Arbitrary comment from slurmctld
	Tasks *uint32 `json:"tasks,omitempty"` // Number of tasks (32 bit integer number with flags)
	TasksPerBoard *uint16 `json:"tasks_per_board,omitempty"` // Number of tasks invoked on each board (16 bit integer number with flags)
	TasksPerCore *uint16 `json:"tasks_per_core,omitempty"` // Number of tasks invoked on each core (16 bit integer number with flags)
	TasksPerNode *uint16 `json:"tasks_per_node,omitempty"` // Number of tasks invoked on each node (16 bit integer number with flags)
	TasksPerSocket *uint16 `json:"tasks_per_socket,omitempty"` // Number of tasks invoked on each socket (16 bit integer number with flags)
	TasksPerTRES *uint16 `json:"tasks_per_tres,omitempty"` // Number of tasks that can assess each GPU (16 bit integer number with flags)
	ThreadSpec *int32 `json:"thread_spec,omitempty"` // Specialized thread count
	ThreadsPerCore *uint16 `json:"threads_per_core,omitempty"` // Number of processor threads per CPU core required (16 bit integer number with...
	TimeLimit *uint32 `json:"time_limit,omitempty"` // Maximum run time in minutes (32 bit integer number with flags)
	TimeMinimum *uint32 `json:"time_minimum,omitempty"` // Minimum run time in minutes (32 bit integer number with flags)
	TRESAllocStr *string `json:"tres_alloc_str,omitempty"` // TRES used by the job
	TRESBind *string `json:"tres_bind,omitempty"` // Task to TRES binding directives
	TRESFreq *string `json:"tres_freq,omitempty"` // TRES frequency directives
	TRESPerJob *string `json:"tres_per_job,omitempty"` // Comma-separated list of TRES=# values to be allocated per job
	TRESPerNode *string `json:"tres_per_node,omitempty"` // Comma-separated list of TRES=# values to be allocated per node
	TRESPerSocket *string `json:"tres_per_socket,omitempty"` // Comma-separated list of TRES=# values to be allocated per socket
	TRESPerTask *string `json:"tres_per_task,omitempty"` // Comma-separated list of TRES=# values to be allocated per task
	TRESReqStr *string `json:"tres_req_str,omitempty"` // TRES requested by the job
	UserID *int32 `json:"user_id,omitempty"` // User ID that owns the job
	UserName *string `json:"user_name,omitempty"` // User name that owns the job
	Wckey *string `json:"wckey,omitempty"` // Workload characterization key
}


// JobPower is a nested type within its parent.
type JobPower struct {
	Flags []interface{} `json:"flags,omitempty"`
}


// FlagsValue represents possible values for Flags field.
type FlagsValue string

// FlagsValue constants.
const (
	FlagsKillInvalidDependency FlagsValue = "KILL_INVALID_DEPENDENCY"
	FlagsNoKillInvalidDependency FlagsValue = "NO_KILL_INVALID_DEPENDENCY"
	FlagsHasStateDirectory FlagsValue = "HAS_STATE_DIRECTORY"
	FlagsTestingBackfill FlagsValue = "TESTING_BACKFILL"
	FlagsGRESBindingEnforced FlagsValue = "GRES_BINDING_ENFORCED"
	FlagsTestNowOnly FlagsValue = "TEST_NOW_ONLY"
	FlagsSendJobEnvironment FlagsValue = "SEND_JOB_ENVIRONMENT"
	FlagsSpreadJob FlagsValue = "SPREAD_JOB"
	FlagsPreferMinimumNodeCount FlagsValue = "PREFER_MINIMUM_NODE_COUNT"
	FlagsJobKillHurry FlagsValue = "JOB_KILL_HURRY"
	FlagsSkipTRESStringAccounting FlagsValue = "SKIP_TRES_STRING_ACCOUNTING"
	FlagsSiblingClusterUpdateOnly FlagsValue = "SIBLING_CLUSTER_UPDATE_ONLY"
	FlagsHeterogeneousJob FlagsValue = "HETEROGENEOUS_JOB"
	FlagsExactTaskCountRequested FlagsValue = "EXACT_TASK_COUNT_REQUESTED"
	FlagsExactCPUCountRequested FlagsValue = "EXACT_CPU_COUNT_REQUESTED"
	FlagsTestingWholeNodeBackfill FlagsValue = "TESTING_WHOLE_NODE_BACKFILL"
	FlagsTopPriorityJob FlagsValue = "TOP_PRIORITY_JOB"
	FlagsAccrueCountCleared FlagsValue = "ACCRUE_COUNT_CLEARED"
	FlagsGRESBindingDisabled FlagsValue = "GRES_BINDING_DISABLED"
	FlagsJobWasRunning FlagsValue = "JOB_WAS_RUNNING"
	FlagsJobAccrueTimeReset FlagsValue = "JOB_ACCRUE_TIME_RESET"
	FlagsCronJob FlagsValue = "CRON_JOB"
	FlagsExactMemoryRequested FlagsValue = "EXACT_MEMORY_REQUESTED"
	FlagsExternalJob FlagsValue = "EXTERNAL_JOB"
	FlagsUsingDefaultAccount FlagsValue = "USING_DEFAULT_ACCOUNT"
	FlagsUsingDefaultPartition FlagsValue = "USING_DEFAULT_PARTITION"
	FlagsUsingDefaultQoS FlagsValue = "USING_DEFAULT_QOS"
	FlagsUsingDefaultWckey FlagsValue = "USING_DEFAULT_WCKEY"
	FlagsDependent FlagsValue = "DEPENDENT"
	FlagsMagnetic FlagsValue = "MAGNETIC"
	FlagsPartitionAssigned FlagsValue = "PARTITION_ASSIGNED"
	FlagsBackfillAttempted FlagsValue = "BACKFILL_ATTEMPTED"
	FlagsSchedulingAttempted FlagsValue = "SCHEDULING_ATTEMPTED"
	FlagsStepmgrEnabled FlagsValue = "STEPMGR_ENABLED"
	FlagsSpreadSegments FlagsValue = "SPREAD_SEGMENTS"
	FlagsConsolidateSegments FlagsValue = "CONSOLIDATE_SEGMENTS"
	FlagsExpeditedRequeue FlagsValue = "EXPEDITED_REQUEUE"
)

// JobState represents possible values for JobState field.
type JobState string

// JobState constants.
const (
	JobStatePending JobState = "PENDING"
	JobStateRunning JobState = "RUNNING"
	JobStateSuspended JobState = "SUSPENDED"
	JobStateCompleted JobState = "COMPLETED"
	JobStateCancelled JobState = "CANCELLED"
	JobStateFailed JobState = "FAILED"
	JobStateTimeout JobState = "TIMEOUT"
	JobStateNodeFail JobState = "NODE_FAIL"
	JobStatePreempted JobState = "PREEMPTED"
	JobStateBootFail JobState = "BOOT_FAIL"
	JobStateDeadline JobState = "DEADLINE"
	JobStateOutOfMemory JobState = "OUT_OF_MEMORY"
	JobStateLaunchFailed JobState = "LAUNCH_FAILED"
	JobStateRequeued JobState = "REQUEUED"
	JobStateRequeueHold JobState = "REQUEUE_HOLD"
	JobStateSpecialExit JobState = "SPECIAL_EXIT"
	JobStateResizing JobState = "RESIZING"
	JobStateConfiguring JobState = "CONFIGURING"
	JobStateCompleting JobState = "COMPLETING"
	JobStateStopped JobState = "STOPPED"
	JobStateReconfigFail JobState = "RECONFIG_FAIL"
	JobStatePowerUpNode JobState = "POWER_UP_NODE"
	JobStateRevoked JobState = "REVOKED"
	JobStateRequeueFed JobState = "REQUEUE_FED"
	JobStateResvDelHold JobState = "RESV_DEL_HOLD"
	JobStateSignaling JobState = "SIGNALING"
	JobStateStageOut JobState = "STAGE_OUT"
	JobStateExpediting JobState = "EXPEDITING"
)

// MailTypeValue represents possible values for MailType field.
type MailTypeValue string

// MailTypeValue constants.
const (
	MailTypeBegin MailTypeValue = "BEGIN"
	MailTypeEnd MailTypeValue = "END"
	MailTypeFail MailTypeValue = "FAIL"
	MailTypeRequeue MailTypeValue = "REQUEUE"
	MailTypeTime100 MailTypeValue = "TIME=100%"
	MailTypeTime90 MailTypeValue = "TIME=90%"
	MailTypeTime80 MailTypeValue = "TIME=80%"
	MailTypeTime50 MailTypeValue = "TIME=50%"
	MailTypeStageOut MailTypeValue = "STAGE_OUT"
	MailTypeArrayTasks MailTypeValue = "ARRAY_TASKS"
	MailTypeInvalidDependency MailTypeValue = "INVALID_DEPENDENCY"
)


// ProfileValue represents possible values for Profile field.
type ProfileValue string

// ProfileValue constants.
const (
	ProfileNotSet ProfileValue = "NOT_SET"
	ProfileNone ProfileValue = "NONE"
	ProfileEnergy ProfileValue = "ENERGY"
	ProfileLustre ProfileValue = "LUSTRE"
	ProfileNetwork ProfileValue = "NETWORK"
	ProfileTask ProfileValue = "TASK"
)

// SharedValue represents possible values for Shared field.
type SharedValue string

// SharedValue constants.
const (
	SharedNone SharedValue = "none"
	SharedOversubscribe SharedValue = "oversubscribe"
	SharedUser SharedValue = "user"
	SharedMCS SharedValue = "mcs"
	SharedTopo SharedValue = "topo"
)
