// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package types

import "time"

// === License Types ===

// License represents a SLURM license
type License struct {
	Name        string `json:"name"`
	Total       int    `json:"total"`
	Used        int    `json:"used"`
	Free        int    `json:"free"`
	Reserved    int    `json:"reserved"`
	RemoteUsed  int    `json:"remote_used,omitempty"`
}

// LicenseList represents a list of licenses
type LicenseList struct {
	Licenses []License              `json:"licenses"`
	Meta     map[string]interface{} `json:"meta,omitempty"`
}

// === Share Types ===

// Share represents fairshare information
type Share struct {
	Cluster          string                 `json:"cluster,omitempty"`
	Account          string                 `json:"account,omitempty"`
	User             string                 `json:"user,omitempty"`
	Partition        string                 `json:"partition,omitempty"`
	EffectiveUsage   float64                `json:"effective_usage"`
	FairshareLevel   float64                `json:"fairshare_level"`
	FairshareUsage   float64                `json:"fairshare_usage"`
	FairshareShares  int                    `json:"fairshare_shares"`
	NormalizedShares float64                `json:"normalized_shares"`
	NormalizedUsage  float64                `json:"normalized_usage"`
	RawShares        int                    `json:"raw_shares"`
	RawUsage         int64                  `json:"raw_usage"`
	SharesUsed       int64                  `json:"shares_used"`
	RunSeconds       int64                  `json:"run_seconds"`
	AssocID          int                    `json:"assoc_id,omitempty"`
	ParentAccount    string                 `json:"parent_account,omitempty"`
	Meta             map[string]interface{} `json:"meta,omitempty"`
}

// SharesList represents a list of shares
type SharesList struct {
	Shares []Share                `json:"shares"`
	Meta   map[string]interface{} `json:"meta,omitempty"`
}

// GetSharesOptions provides filtering options for shares
type GetSharesOptions struct {
	Users     []string `json:"users,omitempty"`
	Accounts  []string `json:"accounts,omitempty"`
	Clusters  []string `json:"clusters,omitempty"`
	Partition string   `json:"partition,omitempty"`
}

// === Config Types ===

// Config represents SLURM configuration
type Config struct {
	AccountingStorageType      string                 `json:"accounting_storage_type,omitempty"`
	AccountingStorageHost      string                 `json:"accounting_storage_host,omitempty"`
	AccountingStoragePort      int                    `json:"accounting_storage_port,omitempty"`
	AccountingStorageUser      string                 `json:"accounting_storage_user,omitempty"`
	AccountingStorageEnforce   []string               `json:"accounting_storage_enforce,omitempty"`
	ClusterName                string                 `json:"cluster_name"`
	ControlMachine             []string               `json:"control_machine,omitempty"`
	ControlAddr                string                 `json:"control_addr,omitempty"`
	BackupController           string                 `json:"backup_controller,omitempty"`
	BackupAddr                 string                 `json:"backup_addr,omitempty"`
	SlurmUser                  string                 `json:"slurm_user"`
	SlurmUID                   int                    `json:"slurm_uid"`
	SlurmGID                   int                    `json:"slurm_gid"`
	SlurmctldPort              int                    `json:"slurmctld_port"`
	SlurmdPort                 int                    `json:"slurmd_port"`
	FirstJobID                 int                    `json:"first_job_id"`
	MaxJobCount                int                    `json:"max_job_count"`
	MaxJobTime                 int                    `json:"max_job_time,omitempty"`
	MinJobAge                  int                    `json:"min_job_age"`
	StateLocation              string                 `json:"state_location"`
	StateSaveLocation          string                 `json:"state_save_location,omitempty"`
	SlurmctldPidFile           string                 `json:"slurmctld_pid_file"`
	SlurmdPidFile              string                 `json:"slurmd_pid_file"`
	SlurmdSpoolDir             string                 `json:"slurmd_spool_dir"`
	SlurmctldLogFile           string                 `json:"slurmctld_log_file,omitempty"`
	SlurmdLogFile              string                 `json:"slurmd_log_file,omitempty"`
	SlurmctldDebug             string                 `json:"slurmctld_debug,omitempty"`
	SlurmdDebug                string                 `json:"slurmd_debug,omitempty"`
	SchedulerType              string                 `json:"scheduler_type"`
	SchedulerParameters        map[string]interface{} `json:"scheduler_parameters,omitempty"`
	SelectType                 string                 `json:"select_type"`
	PreemptType                string                 `json:"preempt_type,omitempty"`
	PreemptMode                string                 `json:"preempt_mode,omitempty"`
	PriorityType               string                 `json:"priority_type"`
	PriorityDecayHalfLife      int                    `json:"priority_decay_half_life,omitempty"`
	PriorityMaxAge             int                    `json:"priority_max_age,omitempty"`
	PriorityWeightAge          int                    `json:"priority_weight_age,omitempty"`
	PriorityWeightFairshare    int                    `json:"priority_weight_fairshare,omitempty"`
	PriorityWeightJobSize      int                    `json:"priority_weight_job_size,omitempty"`
	PriorityWeightPartition    int                    `json:"priority_weight_partition,omitempty"`
	PriorityWeightQOS          int                    `json:"priority_weight_qos,omitempty"`
	ProctrackType              string                 `json:"proctrack_type"`
	SwitchType                 string                 `json:"switch_type,omitempty"`
	TaskPlugin                 string                 `json:"task_plugin"`
	TopologyPlugin             string                 `json:"topology_plugin,omitempty"`
	TreeWidth                  int                    `json:"tree_width,omitempty"`
	CPUFreqDef                 string                 `json:"cpu_freq_def,omitempty"`
	CPUFreqGovernors           string                 `json:"cpu_freq_governors,omitempty"`
	CompleteWait               int                    `json:"complete_wait,omitempty"`
	DefMemPerCPU               int                    `json:"def_mem_per_cpu,omitempty"`
	EnforcePartLimits          string                 `json:"enforce_part_limits,omitempty"`
	KillOnBadExit              bool                   `json:"kill_on_bad_exit,omitempty"`
	KillWait                   int                    `json:"kill_wait,omitempty"`
	MaxMemPerCPU               int                    `json:"max_mem_per_cpu,omitempty"`
	MailProg                   string                 `json:"mail_prog,omitempty"`
	PluginDir                  string                 `json:"plugin_dir,omitempty"`
	PlugStackConfig            string                 `json:"plug_stack_config,omitempty"`
	PrivateData                []string               `json:"private_data,omitempty"`
	PropagateResourceLimits    string                 `json:"propagate_resource_limits,omitempty"`
	PropagateResourceLimitsAll string                 `json:"propagate_resource_limits_all,omitempty"`
	ResumeProgram              string                 `json:"resume_program,omitempty"`
	SuspendProgram             string                 `json:"suspend_program,omitempty"`
	SuspendTime                int                    `json:"suspend_time,omitempty"`
	ResumeTimeout              int                    `json:"resume_timeout,omitempty"`
	SuspendTimeout             int                    `json:"suspend_timeout,omitempty"`
	TCPTimeout                 int                    `json:"tcp_timeout,omitempty"`
	UnkillableStepTimeout      int                    `json:"unkillable_step_timeout,omitempty"`
	Version                    string                 `json:"version"`
	Meta                       map[string]interface{} `json:"meta,omitempty"`
}

// === Diagnostics Types ===

// Diagnostics represents SLURM diagnostics information
type Diagnostics struct {
	DataCollected    time.Time              `json:"data_collected"`
	ReqTime          int64                  `json:"req_time"`
	ReqTimeStart     int64                  `json:"req_time_start"`
	ServerThreadCount int                   `json:"server_thread_count"`
	AgentQueueSize   int                    `json:"agent_queue_size"`
	AgentCount       int                    `json:"agent_count"`
	AgentThreadCount int                    `json:"agent_thread_count"`
	DBDAgentCount    int                    `json:"dbd_agent_count"`
	GittosCount      int                    `json:"gittos_count"`
	GittosTime       int64                  `json:"gittos_time"`
	ScheduleCycleLast int64                 `json:"schedule_cycle_last"`
	ScheduleCycleMax  int64                 `json:"schedule_cycle_max"`
	ScheduleCycleMean int64                 `json:"schedule_cycle_mean"`
	ScheduleCycleSum  int64                 `json:"schedule_cycle_sum"`
	ScheduleCycleCounter int                `json:"schedule_cycle_counter"`
	ScheduleCycleDepth int                 `json:"schedule_cycle_depth"`
	ScheduleQueueLen  int                   `json:"schedule_queue_len"`
	JobsSubmitted     int                   `json:"jobs_submitted"`
	JobsStarted       int                   `json:"jobs_started"`
	JobsCompleted     int                   `json:"jobs_completed"`
	JobsCanceled      int                   `json:"jobs_canceled"`
	JobsFailed        int                   `json:"jobs_failed"`
	JobsPending       int                   `json:"jobs_pending"`
	JobsRunning       int                   `json:"jobs_running"`
	JobStatesTs       time.Time             `json:"job_states_ts"`
	BFActive          bool                  `json:"bf_active"`
	BFBackfilledJobs  int                   `json:"bf_backfilled_jobs"`
	BFCycle           int                   `json:"bf_cycle"`
	BFCycleMean       int64                 `json:"bf_cycle_mean"`
	BFCycleMax        int64                 `json:"bf_cycle_max"`
	BFDepth           int                   `json:"bf_depth"`
	BFDepthMean       int                   `json:"bf_depth_mean"`
	BFDepthSum        int                   `json:"bf_depth_sum"`
	BFQueueLen        int                   `json:"bf_queue_len"`
	BFQueueLenMean    int                   `json:"bf_queue_len_mean"`
	BFQueueLenSum     int                   `json:"bf_queue_len_sum"`
	BFTableSize       int                   `json:"bf_table_size"`
	BFTableSizeMean   int                   `json:"bf_table_size_mean"`
	BFTableSizeSum    int                   `json:"bf_table_size_sum"`
	BFWhenLastCycle   time.Time             `json:"bf_when_last_cycle"`
	BFActive2         bool                  `json:"bf_active2"`
	RPCsQueued        int                   `json:"rpcs_queued"`
	RPCsDropped       int                   `json:"rpcs_dropped"`
	RPCsCompleted     int                   `json:"rpcs_completed"`
	RPCsQueued2       int                   `json:"rpcs_queued2"`
	Meta              map[string]interface{} `json:"meta,omitempty"`
}

// === Instance Types ===

// Instance represents a SLURM database instance
type Instance struct {
	Cluster     string    `json:"cluster"`
	ExtraInfo   string    `json:"extra,omitempty"`
	Instance    string    `json:"instance"`
	InstanceID  string    `json:"instance_id"`
	InstanceType string   `json:"instance_type"`
	NodeCount   int       `json:"node_count,omitempty"`
	TimeStart   time.Time `json:"time_start,omitempty"`
	TimeEnd     time.Time `json:"time_end,omitempty"`
}

// InstanceList represents a list of instances
type InstanceList struct {
	Instances []Instance             `json:"instances"`
	Meta      map[string]interface{} `json:"meta,omitempty"`
}

// GetInstanceOptions provides filtering options for instances
type GetInstanceOptions struct {
	Cluster   string     `json:"cluster,omitempty"`
	Extra     string     `json:"extra,omitempty"`
	Format    string     `json:"format,omitempty"`
	Instance  string     `json:"instance,omitempty"`
	NodeList  string     `json:"node_list,omitempty"`
	TimeStart *time.Time `json:"time_start,omitempty"`
	TimeEnd   *time.Time `json:"time_end,omitempty"`
}

// GetInstancesOptions provides filtering options for multiple instances
type GetInstancesOptions struct {
	Clusters  []string   `json:"clusters,omitempty"`
	Extra     string     `json:"extra,omitempty"`
	Format    string     `json:"format,omitempty"`
	Instance  string     `json:"instance,omitempty"`
	NodeList  string     `json:"node_list,omitempty"`
	TimeStart *time.Time `json:"time_start,omitempty"`
	TimeEnd   *time.Time `json:"time_end,omitempty"`
}

// === TRES Types ===

// TRES represents a Trackable RESource
type TRES struct {
	ID    int    `json:"id"`
	Type  string `json:"type"`
	Name  string `json:"name"`
	Count int64  `json:"count,omitempty"`
}

// TRESList represents a list of TRES
type TRESList struct {
	TRES []TRES                 `json:"tres"`
	Meta map[string]interface{} `json:"meta,omitempty"`
}

// CreateTRESRequest represents a request to create TRES
type CreateTRESRequest struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Count       int64  `json:"count,omitempty"`
}

// === Reconfigure Types ===

// ReconfigureResponse represents the response from a reconfigure operation
type ReconfigureResponse struct {
	Status   string                 `json:"status"`
	Message  string                 `json:"message,omitempty"`
	Changes  []string               `json:"changes,omitempty"`
	Warnings []string               `json:"warnings,omitempty"`
	Errors   []string               `json:"errors,omitempty"`
	Meta     map[string]interface{} `json:"meta,omitempty"`
}
