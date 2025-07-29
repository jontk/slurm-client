package v0_0_41

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jontk/slurm-client/internal/common/types"
	api "github.com/jontk/slurm-client/internal/api/v0_0_41"
)

// convertAPIJobToCommon converts a v0.0.41 API Job to common Job type
func (a *JobAdapter) convertAPIJobToCommon(apiJob interface{}) (*types.Job, error) {
	// Type assertion to handle the anonymous struct
	jobData, ok := apiJob.(struct {
		Account             *string                                       `json:"account,omitempty"`
		AccrueTime          *api.V0041OpenapiJobInfoRespJobsAccrueTime   `json:"accrue_time,omitempty"`
		AdminComment        *string                                       `json:"admin_comment,omitempty"`
		AllocatingNode      *string                                       `json:"allocating_node,omitempty"`
		ArrayJobId          *api.V0041OpenapiJobInfoRespJobsArrayJobId   `json:"array_job_id,omitempty"`
		ArrayTaskId         *api.V0041OpenapiJobInfoRespJobsArrayTaskId  `json:"array_task_id,omitempty"`
		ArrayMaxTasks       *api.V0041OpenapiJobInfoRespJobsArrayMaxTasks `json:"array_max_tasks,omitempty"`
		ArrayTaskString     *string                                       `json:"array_task_string,omitempty"`
		AssociationId       *api.V0041OpenapiJobInfoRespJobsAssociationId `json:"association_id,omitempty"`
		BatchFeatures       *string                                       `json:"batch_features,omitempty"`
		BatchFlag           *bool                                         `json:"batch_flag,omitempty"`
		BatchHost           *string                                       `json:"batch_host,omitempty"`
		Flags               *[]api.V0041OpenapiJobInfoRespJobsFlags       `json:"flags,omitempty"`
		BurstBuffer         *string                                       `json:"burst_buffer,omitempty"`
		BurstBufferState    *string                                       `json:"burst_buffer_state,omitempty"`
		Cluster             *string                                       `json:"cluster,omitempty"`
		ClusterFeatures     *string                                       `json:"cluster_features,omitempty"`
		Command             *string                                       `json:"command,omitempty"`
		Comment             *string                                       `json:"comment,omitempty"`
		Container           *string                                       `json:"container,omitempty"`
		ContainerId         *string                                       `json:"container_id,omitempty"`
		Contiguous          *bool                                         `json:"contiguous,omitempty"`
		CoreSpec            *int32                                        `json:"core_spec,omitempty"`
		ThreadSpec          *int32                                        `json:"thread_spec,omitempty"`
		CoresPerSocket      *int32                                        `json:"cores_per_socket,omitempty"`
		Billable            *api.V0041OpenapiJobInfoRespJobsBillable      `json:"billable,omitempty"`
		CpusPerTask         *int32                                        `json:"cpus_per_task,omitempty"`
		CpuFrequencyMinimum *api.V0041OpenapiJobInfoRespJobsCpuFrequencyMinimum `json:"cpu_frequency_minimum,omitempty"`
		CpuFrequencyMaximum *api.V0041OpenapiJobInfoRespJobsCpuFrequencyMaximum `json:"cpu_frequency_maximum,omitempty"`
		CpuFrequencyGovernor *api.V0041OpenapiJobInfoRespJobsCpuFrequencyGovernor `json:"cpu_frequency_governor,omitempty"`
		CpusPerTres         *string                                       `json:"cpus_per_tres,omitempty"`
		Credential          *[]string                                     `json:"credential,omitempty"`
		Deadline            *api.V0041OpenapiJobInfoRespJobsDeadline      `json:"deadline,omitempty"`
		DegreeConstraint    *api.V0041OpenapiJobInfoRespJobsDegreeConstraint `json:"degree_constraint,omitempty"`
		DelayBoot           *api.V0041OpenapiJobInfoRespJobsDelayBoot      `json:"delay_boot,omitempty"`
		Dependency          *string                                       `json:"dependency,omitempty"`
		DerivedExitCode     *api.V0041OpenapiJobInfoRespJobsDerivedExitCode `json:"derived_exit_code,omitempty"`
		EligibleTime        *api.V0041OpenapiJobInfoRespJobsEligibleTime   `json:"eligible_time,omitempty"`
		EndTime             *api.V0041OpenapiJobInfoRespJobsEndTime        `json:"end_time,omitempty"`
		Environment         *[]string                                     `json:"environment,omitempty"`
		ExcludedNodes       *string                                       `json:"excluded_nodes,omitempty"`
		ExitCode            *api.V0041OpenapiJobInfoRespJobsExitCode       `json:"exit_code,omitempty"`
		Extra               *string                                       `json:"extra,omitempty"`
		FailedNode          *string                                       `json:"failed_node,omitempty"`
		Features            *string                                       `json:"features,omitempty"`
		FederationOrigin    *string                                       `json:"federation_origin,omitempty"`
		FederationSiblingsActive *api.V0041OpenapiJobInfoRespJobsFederationSiblingsActive `json:"federation_siblings_active,omitempty"`
		FederationSiblingsViable *api.V0041OpenapiJobInfoRespJobsFederationSiblingsViable `json:"federation_siblings_viable,omitempty"`
		GresTotal           *string                                       `json:"gres_total,omitempty"`
		GroupId             *api.V0041OpenapiJobInfoRespJobsGroupId        `json:"group_id,omitempty"`
		GroupName           *string                                       `json:"group_name,omitempty"`
		JobId               *int32                                        `json:"job_id,omitempty"`
		JobResources        *api.V0041OpenapiJobInfoRespJobsJobResources   `json:"job_resources,omitempty"`
		JobState            *[]api.V0041OpenapiJobInfoRespJobsJobState     `json:"job_state,omitempty"`
		LastSchedEvaluation *api.V0041OpenapiJobInfoRespJobsLastSchedEvaluation `json:"last_sched_evaluation,omitempty"`
		Licenses            *string                                       `json:"licenses,omitempty"`
		MaxCpus             *api.V0041OpenapiJobInfoRespJobsMaxCpus        `json:"max_cpus,omitempty"`
		MaxNodes            *api.V0041OpenapiJobInfoRespJobsMaxNodes       `json:"max_nodes,omitempty"`
		McsLabel            *string                                       `json:"mcs_label,omitempty"`
		MemoryPerCpu        *api.V0041OpenapiJobInfoRespJobsMemoryPerCpu   `json:"memory_per_cpu,omitempty"`
		MemoryPerNode       *api.V0041OpenapiJobInfoRespJobsMemoryPerNode  `json:"memory_per_node,omitempty"`
		Name                *string                                       `json:"name,omitempty"`
		Nodes               *string                                       `json:"nodes,omitempty"`
		Nice                *int32                                        `json:"nice,omitempty"`
		TasksPerCore        *api.V0041OpenapiJobInfoRespJobsTasksPerCore   `json:"tasks_per_core,omitempty"`
		TasksPerNode        *api.V0041OpenapiJobInfoRespJobsTasksPerNode   `json:"tasks_per_node,omitempty"`
		TasksPerSocket      *api.V0041OpenapiJobInfoRespJobsTasksPerSocket `json:"tasks_per_socket,omitempty"`
		TasksPerBoard       *api.V0041OpenapiJobInfoRespJobsTasksPerBoard  `json:"tasks_per_board,omitempty"`
		Cpus                *api.V0041OpenapiJobInfoRespJobsCpus           `json:"cpus,omitempty"`
		NodeCount           *api.V0041OpenapiJobInfoRespJobsNodeCount      `json:"node_count,omitempty"`
		Tasks               *api.V0041OpenapiJobInfoRespJobsTasks          `json:"tasks,omitempty"`
		HetJobId            *api.V0041OpenapiJobInfoRespJobsHetJobId       `json:"het_job_id,omitempty"`
		HetJobIdSet         *string                                       `json:"het_job_id_set,omitempty"`
		HetJobOffset        *api.V0041OpenapiJobInfoRespJobsHetJobOffset   `json:"het_job_offset,omitempty"`
		Partition           *string                                       `json:"partition,omitempty"`
		MemoryPerTres       *string                                       `json:"memory_per_tres,omitempty"`
		MinCpus             *api.V0041OpenapiJobInfoRespJobsMinCpus        `json:"min_cpus,omitempty"`
		MinMemoryPerCpu     *api.V0041OpenapiJobInfoRespJobsMinMemoryPerCpu `json:"min_memory_per_cpu,omitempty"`
		MinMemoryPerNode    *api.V0041OpenapiJobInfoRespJobsMinMemoryPerNode `json:"min_memory_per_node,omitempty"`
		MinNodes            *api.V0041OpenapiJobInfoRespJobsMinNodes        `json:"min_nodes,omitempty"`
		MinTime             *api.V0041OpenapiJobInfoRespJobsMinTime         `json:"min_time,omitempty"`
		MinTmpDisk          *api.V0041OpenapiJobInfoRespJobsMinTmpDisk      `json:"min_tmp_disk,omitempty"`
		TrePerJob           *string                                       `json:"tre_per_job,omitempty"`
		TrePerNode          *string                                       `json:"tre_per_node,omitempty"`
		TrePerSocket        *string                                       `json:"tre_per_socket,omitempty"`
		TrePerTask          *string                                       `json:"tre_per_task,omitempty"`
		Qos                 *string                                       `json:"qos,omitempty"`
		PreemptTime         *api.V0041OpenapiJobInfoRespJobsPreemptTime     `json:"preempt_time,omitempty"`
		PreemptableTime     *api.V0041OpenapiJobInfoRespJobsPreemptableTime `json:"preemptable_time,omitempty"`
		Priority            *api.V0041OpenapiJobInfoRespJobsPriority        `json:"priority,omitempty"`
		Profile             *[]api.V0041OpenapiJobInfoRespJobsProfile       `json:"profile,omitempty"`
		Reboot              *bool                                         `json:"reboot,omitempty"`
		RequiredNodes       *string                                       `json:"required_nodes,omitempty"`
		Requeue             *bool                                         `json:"requeue,omitempty"`
		ResizeTime          *api.V0041OpenapiJobInfoRespJobsResizeTime      `json:"resize_time,omitempty"`
		RestartCnt          *api.V0041OpenapiJobInfoRespJobsRestartCnt      `json:"restart_cnt,omitempty"`
		ResvName            *string                                       `json:"resv_name,omitempty"`
		ScheduledNodes      *string                                       `json:"scheduled_nodes,omitempty"`
		SelectJobinfo       *string                                       `json:"select_jobinfo,omitempty"`
		Shared              *[]api.V0041OpenapiJobInfoRespJobsShared        `json:"shared,omitempty"`
		Exclusive           *[]api.V0041OpenapiJobInfoRespJobsExclusive     `json:"exclusive,omitempty"`
		ShowFlags           *[]api.V0041OpenapiJobInfoRespJobsShowFlags     `json:"show_flags,omitempty"`
		SocketsPerBoard     *int32                                        `json:"sockets_per_board,omitempty"`
		SocketsPerNode      *int32                                        `json:"sockets_per_node,omitempty"`
		StartTime           *api.V0041OpenapiJobInfoRespJobsStartTime       `json:"start_time,omitempty"`
		StateDescription    *string                                       `json:"state_description,omitempty"`
		StateReason         *string                                       `json:"state_reason,omitempty"`
		StandardError       *string                                       `json:"standard_error,omitempty"`
		StandardInput       *string                                       `json:"standard_input,omitempty"`
		StandardOutput      *string                                       `json:"standard_output,omitempty"`
		SubmitTime          *api.V0041OpenapiJobInfoRespJobsSubmitTime      `json:"submit_time,omitempty"`
		SuspendTime         *api.V0041OpenapiJobInfoRespJobsSuspendTime     `json:"suspend_time,omitempty"`
		SystemComment       *string                                       `json:"system_comment,omitempty"`
		TimeLimit           *api.V0041OpenapiJobInfoRespJobsTimeLimit       `json:"time_limit,omitempty"`
		TimeMinimum         *api.V0041OpenapiJobInfoRespJobsTimeMinimum     `json:"time_minimum,omitempty"`
		ThreadsPerCore      *int32                                        `json:"threads_per_core,omitempty"`
		TresBind            *string                                       `json:"tres_bind,omitempty"`
		TresFreq            *string                                       `json:"tres_freq,omitempty"`
		TresPerJob          *string                                       `json:"tres_per_job,omitempty"`
		TresPerNode         *string                                       `json:"tres_per_node,omitempty"`
		TresPerSocket       *string                                       `json:"tres_per_socket,omitempty"`
		TresPerTask         *string                                       `json:"tres_per_task,omitempty"`
		TresReqStr          *string                                       `json:"tres_req_str,omitempty"`
		TresAllocStr        *string                                       `json:"tres_alloc_str,omitempty"`
		UserId              *api.V0041OpenapiJobInfoRespJobsUserId         `json:"user_id,omitempty"`
		UserName            *string                                       `json:"user_name,omitempty"`
		Wckey               *string                                       `json:"wckey,omitempty"`
		CurrentWorkingDirectory *string                                  `json:"current_working_directory,omitempty"`
	})
	if !ok {
		return nil, fmt.Errorf("unexpected job data type")
	}

	job := &types.Job{}

	// Basic fields
	if jobData.JobId != nil {
		job.JobID = uint32(*jobData.JobId)
	}
	if jobData.Name != nil {
		job.Name = *jobData.Name
	}
	if jobData.UserId != nil && jobData.UserId.Number != nil {
		job.UserID = uint32(*jobData.UserId.Number)
	}
	if jobData.UserName != nil {
		job.UserName = *jobData.UserName
	}
	if jobData.GroupId != nil && jobData.GroupId.Number != nil {
		job.GroupID = uint32(*jobData.GroupId.Number)
	}
	if jobData.Account != nil {
		job.Account = *jobData.Account
	}
	if jobData.Partition != nil {
		job.Partition = *jobData.Partition
	}
	if jobData.Qos != nil {
		job.QoS = *jobData.Qos
	}

	// Job state
	if jobData.JobState != nil && len(*jobData.JobState) > 0 {
		job.State = types.JobState((*jobData.JobState)[0])
	}
	if jobData.StateReason != nil {
		job.StateReason = *jobData.StateReason
	}

	// Time fields
	if jobData.SubmitTime != nil && jobData.SubmitTime.Number != nil {
		job.SubmitTime = time.Unix(*jobData.SubmitTime.Number, 0)
	}
	if jobData.StartTime != nil && jobData.StartTime.Number != nil {
		job.StartTime = time.Unix(*jobData.StartTime.Number, 0)
	}
	if jobData.EndTime != nil && jobData.EndTime.Number != nil {
		job.EndTime = time.Unix(*jobData.EndTime.Number, 0)
	}
	if jobData.EligibleTime != nil && jobData.EligibleTime.Number != nil {
		job.EligibleTime = time.Unix(*jobData.EligibleTime.Number, 0)
	}

	// Resource requirements
	if jobData.NodeCount != nil && jobData.NodeCount.Number != nil {
		job.NumNodes = *jobData.NodeCount.Number
	}
	if jobData.Cpus != nil && jobData.Cpus.Number != nil {
		job.NumCPUs = *jobData.Cpus.Number
	}
	if jobData.Tasks != nil && jobData.Tasks.Number != nil {
		job.NumTasks = *jobData.Tasks.Number
	}

	// Memory requirements
	if jobData.MemoryPerCpu != nil && jobData.MemoryPerCpu.Number != nil {
		job.MemoryPerCPU = uint64(*jobData.MemoryPerCpu.Number)
	}
	if jobData.MemoryPerNode != nil && jobData.MemoryPerNode.Number != nil {
		job.MemoryPerNode = uint64(*jobData.MemoryPerNode.Number)
	}

	// Time limit
	if jobData.TimeLimit != nil && jobData.TimeLimit.Number != nil {
		job.TimeLimit = time.Duration(*jobData.TimeLimit.Number) * time.Minute
	}

	// Priority
	if jobData.Priority != nil && jobData.Priority.Number != nil {
		job.Priority = *jobData.Priority.Number
	}

	// Node information
	if jobData.Nodes != nil {
		job.NodeList = *jobData.Nodes
	}
	if jobData.ExcludedNodes != nil {
		job.ExcNodeList = *jobData.ExcludedNodes
	}

	// Features and constraints
	if jobData.Features != nil {
		job.Features = *jobData.Features
	}
	if jobData.Dependency != nil {
		job.Dependency = *jobData.Dependency
	}

	// Array job information
	if jobData.ArrayJobId != nil && jobData.ArrayJobId.Number != nil {
		job.ArrayJobID = uint32(*jobData.ArrayJobId.Number)
	}
	if jobData.ArrayTaskId != nil && jobData.ArrayTaskId.Number != nil {
		job.ArrayTaskID = uint32(*jobData.ArrayTaskId.Number)
	}

	// Exit codes
	if jobData.ExitCode != nil {
		if jobData.ExitCode.ReturnCode != nil && jobData.ExitCode.ReturnCode.Number != nil {
			job.ExitCode = int32(*jobData.ExitCode.ReturnCode.Number)
		}
	}

	// Standard I/O
	if jobData.StandardInput != nil {
		job.StdIn = *jobData.StandardInput
	}
	if jobData.StandardOutput != nil {
		job.StdOut = *jobData.StandardOutput
	}
	if jobData.StandardError != nil {
		job.StdErr = *jobData.StandardError
	}

	// Working directory
	if jobData.CurrentWorkingDirectory != nil {
		job.WorkDir = *jobData.CurrentWorkingDirectory
	}

	// Environment
	if jobData.Environment != nil {
		job.Environment = *jobData.Environment
	}

	// Comments
	if jobData.Comment != nil {
		job.Comment = *jobData.Comment
	}
	if jobData.AdminComment != nil {
		job.AdminComment = *jobData.AdminComment
	}

	// GRES
	if jobData.GresTotal != nil {
		job.Gres = *jobData.GresTotal
	}

	// Batch job info
	if jobData.BatchFlag != nil {
		job.BatchFlag = *jobData.BatchFlag
	}
	if jobData.BatchHost != nil {
		job.BatchHost = *jobData.BatchHost
	}

	return job, nil
}

// convertCommonToAPIJobSubmit converts common JobSubmitOptions to v0.0.41 API request
func (a *JobAdapter) convertCommonToAPIJobSubmit(opts *types.JobSubmitOptions) *api.V0041OpenapiJobSubmitReq {
	req := &api.V0041OpenapiJobSubmitReq{
		Jobs: []struct {
			Account                 *string   `json:"account,omitempty"`
			AccountGatherFrequency  *string   `json:"account_gather_frequency,omitempty"`
			Argv                    *[]string `json:"argv,omitempty"`
			Array                   *string   `json:"array,omitempty"`
			BatchFeatures           *string   `json:"batch_features,omitempty"`
			Begin                   *string   `json:"begin,omitempty"`
			ClusterConstraint       *string   `json:"cluster_constraint,omitempty"`
			Comment                 *string   `json:"comment,omitempty"`
			Constraints             *string   `json:"constraints,omitempty"`
			ContainerImage          *string   `json:"container_image,omitempty"`
			ContiguousNode          *bool     `json:"contiguous_node,omitempty"`
			CoreSpecification       *int32    `json:"core_specification,omitempty"`
			CoresPerSocket          *int32    `json:"cores_per_socket,omitempty"`
			CpuBinding              *string   `json:"cpu_binding,omitempty"`
			CpuBindingFlags         *[]string `json:"cpu_binding_flags,omitempty"`
			CpuFrequency            *string   `json:"cpu_frequency,omitempty"`
			CpusPerGpu              *string   `json:"cpus_per_gpu,omitempty"`
			CpusPerTask             *int32    `json:"cpus_per_task,omitempty"`
			CurrentWorkingDirectory *string   `json:"current_working_directory,omitempty"`
			Deadline                *string   `json:"deadline,omitempty"`
			DelayBoot               *int32    `json:"delay_boot,omitempty"`
			Dependency              *string   `json:"dependency,omitempty"`
			Distribution            *string   `json:"distribution,omitempty"`
			Environment             *[]string `json:"environment,omitempty"`
			Exclusive               *string   `json:"exclusive,omitempty"`
			GetUserEnvironment      *bool     `json:"get_user_environment,omitempty"`
			Gres                    *string   `json:"gres,omitempty"`
			GresFlags               *[]string `json:"gres_flags,omitempty"`
			GpuBinding              *string   `json:"gpu_binding,omitempty"`
			GpuFrequency            *string   `json:"gpu_frequency,omitempty"`
			Gpus                    *string   `json:"gpus,omitempty"`
			GpusPerNode             *string   `json:"gpus_per_node,omitempty"`
			GpusPerSocket           *string   `json:"gpus_per_socket,omitempty"`
			GpusPerTask             *string   `json:"gpus_per_task,omitempty"`
			HealthCheckInterval     *int32    `json:"health_check_interval,omitempty"`
			HeldArrayTaskIds        *string   `json:"held_array_task_ids,omitempty"`
			HoldJob                 *bool     `json:"hold_job,omitempty"`
			HeterogeneousJob        *string   `json:"heterogeneous_job,omitempty"`
			Immediate               *bool     `json:"immediate,omitempty"`
			JobId                   *int32    `json:"job_id,omitempty"`
			KillOnNodeFail          *bool     `json:"kill_on_node_fail,omitempty"`
			Licenses                *string   `json:"licenses,omitempty"`
			MailType                *[]string `json:"mail_type,omitempty"`
			MailUser                *string   `json:"mail_user,omitempty"`
			McsLabel                *string   `json:"mcs_label,omitempty"`
			MemoryBinding           *string   `json:"memory_binding,omitempty"`
			MemoryPerCpu            *int64    `json:"memory_per_cpu,omitempty"`
			MemoryPerGpu            *int64    `json:"memory_per_gpu,omitempty"`
			MemoryPerNode           *int64    `json:"memory_per_node,omitempty"`
			MinimumCpusPerNode      *int32    `json:"minimum_cpus_per_node,omitempty"`
			MinimumNodes            *bool     `json:"minimum_nodes,omitempty"`
			Name                    *string   `json:"name,omitempty"`
			Nice                    *int32    `json:"nice,omitempty"`
			NoKill                  *bool     `json:"no_kill,omitempty"`
			Nodes                   *string   `json:"nodes,omitempty"`
			NodesMax                *int32    `json:"nodes_max,omitempty"`
			NodesMin                *int32    `json:"nodes_min,omitempty"`
			OverrideTimeLimit       *int32    `json:"override_time_limit,omitempty"`
			Partition               *string   `json:"partition,omitempty"`
			Priority                *int32    `json:"priority,omitempty"`
			Qos                     *string   `json:"qos,omitempty"`
			Requeue                 *bool     `json:"requeue,omitempty"`
			Reservation             *string   `json:"reservation,omitempty"`
			Script                  *string   `json:"script,omitempty"`
			Shared                  *[]string `json:"shared,omitempty"`
			SignalNumber            *int32    `json:"signal_number,omitempty"`
			SignalTime              *int32    `json:"signal_time,omitempty"`
			SocketsPerNode          *int32    `json:"sockets_per_node,omitempty"`
			SpreadJob               *bool     `json:"spread_job,omitempty"`
			StandardError           *string   `json:"standard_error,omitempty"`
			StandardInput           *string   `json:"standard_input,omitempty"`
			StandardOutput          *string   `json:"standard_output,omitempty"`
			Tasks                   *int32    `json:"tasks,omitempty"`
			TasksPerCore            *int32    `json:"tasks_per_core,omitempty"`
			TasksPerNode            *int32    `json:"tasks_per_node,omitempty"`
			TasksPerSocket          *int32    `json:"tasks_per_socket,omitempty"`
			ThreadSpecification     *int32    `json:"thread_specification,omitempty"`
			ThreadsPerCore          *int32    `json:"threads_per_core,omitempty"`
			TimeLimit               *int32    `json:"time_limit,omitempty"`
			TimeMinimum             *int32    `json:"time_minimum,omitempty"`
			WaitAllNodes            *bool     `json:"wait_all_nodes,omitempty"`
			WarningFlags            *[]string `json:"warning_flags,omitempty"`
			WarningSignal           *int32    `json:"warning_signal,omitempty"`
			WarningTime             *int32    `json:"warning_time,omitempty"`
			Wckey                   *string   `json:"wckey,omitempty"`
		}{
			{},
		},
	}

	job := &req.Jobs[0]

	// Set basic fields
	if opts.Name != "" {
		job.Name = &opts.Name
	}
	if opts.Account != "" {
		job.Account = &opts.Account
	}
	if opts.Partition != "" {
		job.Partition = &opts.Partition
	}
	if opts.QoS != "" {
		job.Qos = &opts.QoS
	}
	if opts.Script != "" {
		job.Script = &opts.Script
	}
	if opts.WorkDir != "" {
		job.CurrentWorkingDirectory = &opts.WorkDir
	}

	// Set resource requirements
	if opts.NumNodes > 0 {
		nodes := int32(opts.NumNodes)
		job.NodesMin = &nodes
		job.NodesMax = &nodes
	}
	if opts.NumTasks > 0 {
		tasks := int32(opts.NumTasks)
		job.Tasks = &tasks
	}
	if opts.CPUsPerTask > 0 {
		cpus := int32(opts.CPUsPerTask)
		job.CpusPerTask = &cpus
	}
	if opts.MemoryPerCPU > 0 {
		mem := int64(opts.MemoryPerCPU)
		job.MemoryPerCpu = &mem
	}
	if opts.MemoryPerNode > 0 {
		mem := int64(opts.MemoryPerNode)
		job.MemoryPerNode = &mem
	}

	// Set time limit
	if opts.TimeLimit > 0 {
		// Convert duration to minutes
		timeLimitMinutes := int32(opts.TimeLimit.Minutes())
		job.TimeLimit = &timeLimitMinutes
	}

	// Set I/O paths
	if opts.StdOut != "" {
		job.StandardOutput = &opts.StdOut
	}
	if opts.StdErr != "" {
		job.StandardError = &opts.StdErr
	}
	if opts.StdIn != "" {
		job.StandardInput = &opts.StdIn
	}

	// Set environment
	if len(opts.Environment) > 0 {
		job.Environment = &opts.Environment
	}

	// Set constraints and features
	if opts.Features != "" {
		job.Constraints = &opts.Features
	}
	if opts.Dependency != "" {
		job.Dependency = &opts.Dependency
	}
	if opts.ExcludeNodes != "" {
		// v0.0.41 doesn't have direct exclude_nodes, we need to use constraints
		constraint := fmt.Sprintf("exclude:%s", opts.ExcludeNodes)
		job.Constraints = &constraint
	}

	// Set GRES
	if opts.Gres != "" {
		job.Gres = &opts.Gres
	}

	// Set array job parameters
	if opts.ArraySpec != "" {
		job.Array = &opts.ArraySpec
	}

	// Set mail options
	if opts.MailUser != "" {
		job.MailUser = &opts.MailUser
	}
	if opts.MailType != "" {
		mailTypes := strings.Split(opts.MailType, ",")
		job.MailType = &mailTypes
	}

	// Set other options
	if opts.Comment != "" {
		job.Comment = &opts.Comment
	}
	if opts.Priority > 0 {
		priority := int32(opts.Priority)
		job.Priority = &priority
	}
	if opts.Nice > 0 {
		nice := int32(opts.Nice)
		job.Nice = &nice
	}
	if opts.Hold {
		job.HoldJob = &opts.Hold
	}
	if opts.Requeue {
		job.Requeue = &opts.Requeue
	}

	return req
}

// convertCommonToAPIJobUpdate converts common JobUpdate to v0.0.41 API request
func (a *JobAdapter) convertCommonToAPIJobUpdate(update *types.JobUpdate) *api.V0041OpenapiJobUpdateReq {
	req := &api.V0041OpenapiJobUpdateReq{
		Jobs: []struct {
			AdminComment     *string `json:"admin_comment,omitempty"`
			ArrayInx         *string `json:"array_inx,omitempty"`
			ArrayTaskThrottle *int32  `json:"array_task_throttle,omitempty"`
			BatchScript      *string `json:"batch_script,omitempty"`
			BurstBuffer      *string `json:"burst_buffer,omitempty"`
			ClusterFeatures  *string `json:"cluster_features,omitempty"`
			Comment          *string `json:"comment,omitempty"`
			Container        *string `json:"container,omitempty"`
			ContainerId      *string `json:"container_id,omitempty"`
			CoreSpec         *int32  `json:"core_spec,omitempty"`
			CoresPerSocket   *int32  `json:"cores_per_socket,omitempty"`
			CpusPerTask      *int32  `json:"cpus_per_task,omitempty"`
			CpuFreqGov       *string `json:"cpu_freq_gov,omitempty"`
			CpuFreqMax       *int32  `json:"cpu_freq_max,omitempty"`
			CpuFreqMin       *int32  `json:"cpu_freq_min,omitempty"`
			Deadline         *string `json:"deadline,omitempty"`
			DelayBoot        *int32  `json:"delay_boot,omitempty"`
			Dependency       *string `json:"dependency,omitempty"`
			EndTime          *string `json:"end_time,omitempty"`
			Environment      *string `json:"environment,omitempty"`
			ExcNodes         *string `json:"exc_nodes,omitempty"`
			Extra            *string `json:"extra,omitempty"`
			Features         *string `json:"features,omitempty"`
			FedIsLock        *bool   `json:"fed_is_lock,omitempty"`
			GresReq          *string `json:"gres_req,omitempty"`
			Hold             *bool   `json:"hold,omitempty"`
			JobId            *string `json:"job_id,omitempty"`
			Kill_on_node_fail *bool   `json:"kill_on_node_fail,omitempty"`
			Licenses         *string `json:"licenses,omitempty"`
			MailType         *string `json:"mail_type,omitempty"`
			MailUser         *string `json:"mail_user,omitempty"`
			McsLabel         *string `json:"mcs_label,omitempty"`
			MemPerCpu        *int64  `json:"mem_per_cpu,omitempty"`
			MemPerGpu        *int64  `json:"mem_per_gpu,omitempty"`
			MemPerNode       *int64  `json:"mem_per_node,omitempty"`
			MinCpusNode      *int32  `json:"min_cpus_node,omitempty"`
			MinMemoryCpu     *int64  `json:"min_memory_cpu,omitempty"`
			MinMemoryNode    *int64  `json:"min_memory_node,omitempty"`
			MinTmpDisk       *int32  `json:"min_tmp_disk,omitempty"`
			Name             *string `json:"name,omitempty"`
			Network          *string `json:"network,omitempty"`
			Nice             *int32  `json:"nice,omitempty"`
			NodeCnt          *int32  `json:"node_cnt,omitempty"`
			NodeInx          *string `json:"node_inx,omitempty"`
			Nodes            *string `json:"nodes,omitempty"`
			NodesMax         *int32  `json:"nodes_max,omitempty"`
			NodesMin         *int32  `json:"nodes_min,omitempty"`
			NumCpus          *int32  `json:"num_cpus,omitempty"`
			NumNodes         *int32  `json:"num_nodes,omitempty"`
			NumTasks         *int32  `json:"num_tasks,omitempty"`
			Oversubscribe    *bool   `json:"oversubscribe,omitempty"`
			Partition        *string `json:"partition,omitempty"`
			PreSusTime       *int32  `json:"pre_sus_time,omitempty"`
			Priority         *int32  `json:"priority,omitempty"`
			Profile          *string `json:"profile,omitempty"`
			Qos              *string `json:"qos,omitempty"`
			Reboot           *bool   `json:"reboot,omitempty"`
			ReqNodes         *string `json:"req_nodes,omitempty"`
			ReqSwitch        *int32  `json:"req_switch,omitempty"`
			Requeue          *bool   `json:"requeue,omitempty"`
			ResvName         *string `json:"resv_name,omitempty"`
			Sched_nodes      *string `json:"sched_nodes,omitempty"`
			Script           *string `json:"script,omitempty"`
			SiblingsActive   *string `json:"siblings_active,omitempty"`
			SiblingsViable   *string `json:"siblings_viable,omitempty"`
			Shared           *string `json:"shared,omitempty"`
			SocketsPerNode   *int32  `json:"sockets_per_node,omitempty"`
			SpankJobEnv      *string `json:"spank_job_env,omitempty"`
			SpankJobEnvSize  *int32  `json:"spank_job_env_size,omitempty"`
			StartTime        *string `json:"start_time,omitempty"`
			StateDesc        *string `json:"state_desc,omitempty"`
			StdErr           *string `json:"std_err,omitempty"`
			StdIn            *string `json:"std_in,omitempty"`
			StdOut           *string `json:"std_out,omitempty"`
			TasksPerCore     *int32  `json:"tasks_per_core,omitempty"`
			TasksPerNode     *int32  `json:"tasks_per_node,omitempty"`
			TasksPerSocket   *int32  `json:"tasks_per_socket,omitempty"`
			TasksPerBoard    *int32  `json:"tasks_per_board,omitempty"`
			ThreadSpec       *int32  `json:"thread_spec,omitempty"`
			ThreadsPerCore   *int32  `json:"threads_per_core,omitempty"`
			TimeLimit        *string `json:"time_limit,omitempty"`
			TimeMin          *string `json:"time_min,omitempty"`
			TresBind         *string `json:"tres_bind,omitempty"`
			TresFreq         *string `json:"tres_freq,omitempty"`
			TresPerJob       *string `json:"tres_per_job,omitempty"`
			TresPerNode      *string `json:"tres_per_node,omitempty"`
			TresPerSocket    *string `json:"tres_per_socket,omitempty"`
			TresPerTask      *string `json:"tres_per_task,omitempty"`
			TresReqStr       *string `json:"tres_req_str,omitempty"`
			UserId           *string `json:"user_id,omitempty"`
			Wait4switch      *int32  `json:"wait4switch,omitempty"`
			Wckey            *string `json:"wckey,omitempty"`
			WorkDir          *string `json:"work_dir,omitempty"`
		}{
			{},
		},
	}

	job := &req.Jobs[0]

	// Set priority
	if update.Priority != nil {
		priority := int32(*update.Priority)
		job.Priority = &priority
	}

	// Set time limit
	if update.TimeLimit != nil {
		// Convert duration to string format (e.g., "1-00:00:00" for 1 day)
		timeLimitStr := formatDurationForSlurm(*update.TimeLimit)
		job.TimeLimit = &timeLimitStr
	}

	// Set partition
	if update.Partition != nil {
		job.Partition = update.Partition
	}

	// Set QoS
	if update.QoS != nil {
		job.Qos = update.QoS
	}

	// Set node count
	if update.NodeCount != nil {
		nodeCount := int32(*update.NodeCount)
		job.NumNodes = &nodeCount
	}

	// Set features/constraints
	if update.Features != nil {
		job.Features = update.Features
	}

	// Set comment
	if update.Comment != nil {
		job.Comment = update.Comment
	}

	// Set hold state
	if update.Hold != nil {
		job.Hold = update.Hold
	}

	// Set nice value
	if update.Nice != nil {
		nice := int32(*update.Nice)
		job.Nice = &nice
	}

	// Set requeue flag
	if update.Requeue != nil {
		job.Requeue = update.Requeue
	}

	// Set account
	if update.Account != nil {
		// v0.0.41 doesn't support account update through job update
		// This would need to be done through association update
	}

	// Set WCKEY
	if update.WCKey != nil {
		job.Wckey = update.WCKey
	}

	// Set reservation
	if update.Reservation != nil {
		job.ResvName = update.Reservation
	}

	// Set dependency
	if update.Dependency != nil {
		job.Dependency = update.Dependency
	}

	// Set deadline
	if update.Deadline != nil {
		deadlineStr := update.Deadline.Format("2006-01-02T15:04:05")
		job.Deadline = &deadlineStr
	}

	return req
}

// formatDurationForSlurm converts a Go duration to Slurm time format
func formatDurationForSlurm(d time.Duration) string {
	totalMinutes := int(d.Minutes())
	days := totalMinutes / 1440
	hours := (totalMinutes % 1440) / 60
	minutes := totalMinutes % 60

	if days > 0 {
		return fmt.Sprintf("%d-%02d:%02d:00", days, hours, minutes)
	}
	return fmt.Sprintf("%02d:%02d:00", hours, minutes)
}

// parseNumberField safely extracts int32 from a NumberField
func parseNumberField(field interface{}) (int32, bool) {
	if field == nil {
		return 0, false
	}

	// Handle different number field types
	switch v := field.(type) {
	case *struct {
		Infinite *bool  `json:"infinite,omitempty"`
		Number   *int32 `json:"number,omitempty"`
		Set      *bool  `json:"set,omitempty"`
	}:
		if v != nil && v.Set != nil && *v.Set && v.Number != nil {
			return *v.Number, true
		}
	case *struct {
		Infinite *bool  `json:"infinite,omitempty"`
		Number   *int64 `json:"number,omitempty"`
		Set      *bool  `json:"set,omitempty"`
	}:
		if v != nil && v.Set != nil && *v.Set && v.Number != nil {
			return int32(*v.Number), true
		}
	}

	return 0, false
}

// parseNumberField64 safely extracts int64 from a NumberField
func parseNumberField64(field interface{}) (int64, bool) {
	if field == nil {
		return 0, false
	}

	// Handle different number field types
	switch v := field.(type) {
	case *struct {
		Infinite *bool  `json:"infinite,omitempty"`
		Number   *int64 `json:"number,omitempty"`
		Set      *bool  `json:"set,omitempty"`
	}:
		if v != nil && v.Set != nil && *v.Set && v.Number != nil {
			return *v.Number, true
		}
	case *struct {
		Infinite *bool  `json:"infinite,omitempty"`
		Number   *int32 `json:"number,omitempty"`
		Set      *bool  `json:"set,omitempty"`
	}:
		if v != nil && v.Set != nil && *v.Set && v.Number != nil {
			return int64(*v.Number), true
		}
	}

	return 0, false
}