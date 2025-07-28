package v0_0_43

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jontk/slurm-client/internal/common/types"
	api "github.com/jontk/slurm-client/internal/api/v0_0_43"
)

// convertAPIJobToCommon converts a v0.0.43 API Job to common Job type
func (a *JobAdapter) convertAPIJobToCommon(apiJob api.V0043JobInfo) (*types.Job, error) {
	job := &types.Job{}

	// Basic fields
	if apiJob.JobId != nil {
		job.JobID = *apiJob.JobId
	}
	if apiJob.Name != nil {
		job.Name = *apiJob.Name
	}
	if apiJob.UserId != nil {
		job.UserID = *apiJob.UserId
	}
	if apiJob.UserName != nil {
		job.UserName = *apiJob.UserName
	}
	if apiJob.GroupId != nil {
		job.GroupID = *apiJob.GroupId
	}
	if apiJob.Account != nil {
		job.Account = *apiJob.Account
	}
	if apiJob.Partition != nil {
		job.Partition = *apiJob.Partition
	}
	if apiJob.Qos != nil {
		job.QoS = *apiJob.Qos
	}

	// Job state
	if apiJob.JobState != nil && len(*apiJob.JobState) > 0 {
		job.State = types.JobState((*apiJob.JobState)[0])
	}
	if apiJob.StateReason != nil {
		job.StateReason = *apiJob.StateReason
	}

	// Time fields
	if apiJob.TimeLimit != nil && apiJob.TimeLimit.Number != nil {
		job.TimeLimit = *apiJob.TimeLimit.Number
	}
	if apiJob.SubmitTime != nil && apiJob.SubmitTime.Number != nil {
		job.SubmitTime = time.Unix(*apiJob.SubmitTime.Number, 0)
	}
	if apiJob.StartTime != nil && apiJob.StartTime.Number != nil && *apiJob.StartTime.Number > 0 {
		startTime := time.Unix(*apiJob.StartTime.Number, 0)
		job.StartTime = &startTime
	}
	if apiJob.EndTime != nil && apiJob.EndTime.Number != nil && *apiJob.EndTime.Number > 0 {
		endTime := time.Unix(*apiJob.EndTime.Number, 0)
		job.EndTime = &endTime
	}

	// Priority
	if apiJob.Priority != nil && apiJob.Priority.Number != nil {
		job.Priority = *apiJob.Priority.Number
	}

	// Resource allocation
	if apiJob.Cpus != nil && apiJob.Cpus.Number != nil {
		job.CPUs = *apiJob.Cpus.Number
	}
	if apiJob.Nodes != nil {
		job.NodeList = *apiJob.Nodes
	}
	if apiJob.NodeCount != nil && len(*apiJob.NodeCount) > 0 {
		// Take the first node count value
		job.Nodes = (*apiJob.NodeCount)[0]
	}

	// Job specification
	if apiJob.Command != nil && len(*apiJob.Command) > 0 {
		job.Command = (*apiJob.Command)[0]
	}
	if apiJob.WorkingDirectory != nil && len(*apiJob.WorkingDirectory) > 0 {
		job.WorkingDirectory = (*apiJob.WorkingDirectory)[0]
	}
	if apiJob.StandardInput != nil {
		job.StandardInput = *apiJob.StandardInput
	}
	if apiJob.StandardOutput != nil {
		job.StandardOutput = *apiJob.StandardOutput
	}
	if apiJob.StandardError != nil {
		job.StandardError = *apiJob.StandardError
	}

	// Array job information
	if apiJob.ArrayJobId != nil {
		job.ArrayJobID = apiJob.ArrayJobId
	}
	if apiJob.ArrayTaskId != nil {
		job.ArrayTaskID = apiJob.ArrayTaskId
	}
	if apiJob.ArrayTaskStr != nil {
		job.ArrayTaskString = *apiJob.ArrayTaskStr
	}

	// Environment
	if apiJob.Environment != nil && len(*apiJob.Environment) > 0 {
		env := make(map[string]string)
		for _, envVar := range *apiJob.Environment {
			parts := strings.SplitN(envVar, "=", 2)
			if len(parts) == 2 {
				env[parts[0]] = parts[1]
			}
		}
		job.Environment = env
	}

	// Mail settings
	if apiJob.MailType != nil && len(*apiJob.MailType) > 0 {
		mailTypes := make([]string, len(*apiJob.MailType))
		for i, mt := range *apiJob.MailType {
			mailTypes[i] = string(mt)
		}
		job.MailType = mailTypes
	}
	if apiJob.MailUser != nil {
		job.MailUser = *apiJob.MailUser
	}

	// Additional fields
	if apiJob.ExcludeNodes != nil && len(*apiJob.ExcludeNodes) > 0 {
		job.ExcludeNodes = (*apiJob.ExcludeNodes)[0]
	}
	if apiJob.Nice != nil {
		job.Nice = *apiJob.Nice
	}
	if apiJob.Comment != nil {
		job.Comment = *apiJob.Comment
	}

	// Features and GRES
	if apiJob.Features != nil {
		job.Features = *apiJob.Features
	}
	if apiJob.Gres != nil {
		job.Gres = *apiJob.Gres
	}

	// Resource requests
	resourceReq := types.ResourceRequests{}
	if apiJob.Memory != nil && apiJob.Memory.Number != nil {
		resourceReq.Memory = *apiJob.Memory.Number
	}
	if apiJob.MemoryPerCpu != nil && apiJob.MemoryPerCpu.Number != nil {
		resourceReq.MemoryPerCPU = *apiJob.MemoryPerCpu.Number
	}
	if apiJob.TmpDisk != nil {
		resourceReq.TmpDisk = *apiJob.TmpDisk
	}
	if apiJob.CpusPerTask != nil {
		resourceReq.CPUsPerTask = *apiJob.CpusPerTask
	}

	job.ResourceRequests = resourceReq

	// Reservation
	if apiJob.Reservation != nil {
		job.Reservation = *apiJob.Reservation
	}

	return job, nil
}

// convertCommonJobCreateToAPI converts common JobCreate type to v0.0.43 API format
func (a *JobAdapter) convertCommonJobCreateToAPI(create *types.JobCreate) (*api.V0043JobProperties, error) {
	apiJob := &api.V0043JobProperties{}

	// Basic fields
	if create.Name != "" {
		apiJob.Name = &create.Name
	}
	if create.Account != "" {
		apiJob.Account = &create.Account
	}
	if create.Partition != "" {
		apiJob.Partition = &create.Partition
	}
	if create.QoS != "" {
		apiJob.Qos = &create.QoS
	}

	// Time limit
	if create.TimeLimit > 0 {
		setTrue := true
		timeLimit := int32(create.TimeLimit)
		apiJob.TimeLimit = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &timeLimit,
		}
	}

	// Priority
	if create.Priority != nil {
		setTrue := true
		priority := *create.Priority
		apiJob.Priority = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &priority,
		}
	}

	// Resource allocation
	if create.CPUs > 0 {
		setTrue := true
		cpus := int32(create.CPUs)
		apiJob.Cpus = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &cpus,
		}
	}
	if create.Nodes > 0 {
		setTrue := true
		nodes := int32(create.Nodes)
		apiJob.NodeCount = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &nodes,
		}
	}
	if create.Tasks > 0 {
		setTrue := true
		tasks := int32(create.Tasks)
		apiJob.Tasks = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &tasks,
		}
	}

	// Command or script
	if create.Script != "" {
		apiJob.Script = &create.Script
	} else if create.Command != "" {
		commands := []string{create.Command}
		apiJob.Argv = &commands
	}

	// Working directory
	if create.WorkingDirectory != "" {
		workingDirs := []string{create.WorkingDirectory}
		apiJob.CurrentWorkingDirectory = &workingDirs
	}

	// Standard I/O
	if create.StandardInput != "" {
		apiJob.StandardInput = &create.StandardInput
	}
	if create.StandardOutput != "" {
		apiJob.StandardOutput = &create.StandardOutput
	}
	if create.StandardError != "" {
		apiJob.StandardError = &create.StandardError
	}

	// Array job
	if create.ArrayString != "" {
		apiJob.ArrayInx = &create.ArrayString
	}

	// Environment
	if len(create.Environment) > 0 {
		envVars := make([]string, 0, len(create.Environment))
		for key, value := range create.Environment {
			envVars = append(envVars, fmt.Sprintf("%s=%s", key, value))
		}
		apiJob.Environment = &envVars
	}

	// Mail settings
	if len(create.MailType) > 0 {
		mailTypes := make([]api.V0043JobMailType, len(create.MailType))
		for i, mt := range create.MailType {
			mailTypes[i] = api.V0043JobMailType(mt)
		}
		apiJob.MailType = &mailTypes
	}
	if create.MailUser != "" {
		apiJob.MailUser = &create.MailUser
	}

	// Additional settings
	if create.ExcludeNodes != "" {
		excludeNodes := []string{create.ExcludeNodes}
		apiJob.ExcludeNodes = &excludeNodes
	}
	if create.Nice != 0 {
		apiJob.Nice = &create.Nice
	}
	if create.Comment != "" {
		apiJob.Comment = &create.Comment
	}

	// Features and GRES
	if len(create.Features) > 0 {
		apiJob.Features = &create.Features
	}
	if create.Gres != "" {
		apiJob.Gres = &create.Gres
	}

	// Resource requests
	if create.ResourceRequests.Memory > 0 {
		setTrue := true
		memory := create.ResourceRequests.Memory
		apiJob.Memory = &api.V0043Uint64NoValStruct{
			Set:    &setTrue,
			Number: &memory,
		}
	}
	if create.ResourceRequests.MemoryPerCPU > 0 {
		setTrue := true
		memPerCpu := create.ResourceRequests.MemoryPerCPU
		apiJob.MemoryPerCpu = &api.V0043Uint64NoValStruct{
			Set:    &setTrue,
			Number: &memPerCpu,
		}
	}
	if create.ResourceRequests.TmpDisk > 0 {
		apiJob.TmpDisk = &create.ResourceRequests.TmpDisk
	}
	if create.ResourceRequests.CPUsPerTask > 0 {
		apiJob.CpusPerTask = &create.ResourceRequests.CPUsPerTask
	}

	// Begin time
	if create.BeginTime != nil {
		setTrue := true
		beginTime := create.BeginTime.Unix()
		apiJob.BeginTime = &api.V0043Uint64NoValStruct{
			Set:    &setTrue,
			Number: &beginTime,
		}
	}

	// Deadline
	if create.Deadline != nil {
		setTrue := true  
		deadline := create.Deadline.Unix()
		apiJob.Deadline = &api.V0043Uint64NoValStruct{
			Set:    &setTrue,
			Number: &deadline,
		}
	}

	// Reservation
	if create.Reservation != "" {
		apiJob.Reservation = &create.Reservation
	}

	// Dependencies
	if len(create.Dependencies) > 0 {
		depStrs := make([]string, len(create.Dependencies))
		for i, dep := range create.Dependencies {
			depStr := dep.Type
			if len(dep.JobIDs) > 0 {
				jobIDStrs := make([]string, len(dep.JobIDs))
				for j, id := range dep.JobIDs {
					jobIDStrs[j] = strconv.FormatInt(int64(id), 10)
				}
				depStr += ":" + strings.Join(jobIDStrs, ",")
			}
			if dep.State != "" {
				depStr += "?" + dep.State
			}
			depStrs[i] = depStr
		}
		dependency := strings.Join(depStrs, ",")
		apiJob.Dependency = &dependency
	}

	// Minimum resource requirements
	if create.MinCPUs > 0 {
		apiJob.MinCpus = &create.MinCPUs
	}
	if create.MinMemory > 0 {
		apiJob.MinMemoryPerNode = &create.MinMemory
	}
	if create.MinNodes > 0 {
		apiJob.MinNodes = &create.MinNodes
	}
	if create.MaxNodes > 0 {
		apiJob.MaxNodes = &create.MaxNodes
	}

	return apiJob, nil
}

// convertCommonJobUpdateToAPI converts common JobUpdate to v0.0.43 API format
func (a *JobAdapter) convertCommonJobUpdateToAPI(existing *types.Job, update *types.JobUpdate) (*api.V0043JobProperties, error) {
	apiJob := &api.V0043JobProperties{}

	// Always include the job ID for updates
	apiJob.JobId = &existing.JobID

	// Apply updates to fields
	name := existing.Name
	if update.Name != nil {
		name = *update.Name
	}
	if name != "" {
		apiJob.Name = &name
	}

	account := existing.Account
	if update.Account != nil {
		account = *update.Account
	}
	if account != "" {
		apiJob.Account = &account
	}

	partition := existing.Partition
	if update.Partition != nil {
		partition = *update.Partition
	}
	if partition != "" {
		apiJob.Partition = &partition
	}

	qos := existing.QoS
	if update.QoS != nil {
		qos = *update.QoS
	}
	if qos != "" {
		apiJob.Qos = &qos
	}

	// Time limit
	timeLimit := existing.TimeLimit
	if update.TimeLimit != nil {
		timeLimit = *update.TimeLimit
	}
	if timeLimit > 0 {
		setTrue := true
		apiJob.TimeLimit = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &timeLimit,
		}
	}

	// Priority
	priority := existing.Priority
	if update.Priority != nil {
		priority = *update.Priority
	}
	if priority > 0 {
		setTrue := true
		apiJob.Priority = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &priority,
		}
	}

	// Nice value
	nice := existing.Nice
	if update.Nice != nil {
		nice = *update.Nice
	}
	if nice != 0 {
		apiJob.Nice = &nice
	}

	// Comment
	comment := existing.Comment
	if update.Comment != nil {
		comment = *update.Comment
	}
	if comment != "" {
		apiJob.Comment = &comment
	}

	// Deadline
	if update.Deadline != nil {
		setTrue := true
		deadline := update.Deadline.Unix()
		apiJob.Deadline = &api.V0043Uint64NoValStruct{
			Set:    &setTrue,
			Number: &deadline,
		}
	}

	// Features
	features := existing.Features
	if len(update.Features) > 0 {
		features = update.Features
	}
	if len(features) > 0 {
		apiJob.Features = &features
	}

	// GRES
	gres := existing.Gres
	if update.Gres != nil {
		gres = *update.Gres
	}
	if gres != "" {
		apiJob.Gres = &gres
	}

	// Begin time
	if update.BeginTime != nil {
		setTrue := true
		beginTime := update.BeginTime.Unix()
		apiJob.BeginTime = &api.V0043Uint64NoValStruct{
			Set:    &setTrue,
			Number: &beginTime,
		}
	}

	// Reservation
	reservation := existing.Reservation
	if update.Reservation != nil {
		reservation = *update.Reservation
	}
	if reservation != "" {
		apiJob.Reservation = &reservation
	}

	// Exclude nodes
	excludeNodes := existing.ExcludeNodes
	if update.ExcludeNodes != nil {
		excludeNodes = *update.ExcludeNodes
	}
	if excludeNodes != "" {
		excludeNodesSlice := []string{excludeNodes}
		apiJob.ExcludeNodes = &excludeNodesSlice
	}

	// Node list (for job updates that specify exact nodes)
	if update.NodeList != nil {
		apiJob.ReqNodes = update.NodeList
	}

	// Min/Max nodes
	if update.MinNodes != nil {
		apiJob.MinNodes = update.MinNodes
	}
	if update.MaxNodes != nil {
		apiJob.MaxNodes = update.MaxNodes
	}

	// Requeue priority
	if update.RequeuePriority != nil {
		setTrue := true
		requeuePriority := *update.RequeuePriority
		apiJob.Priority = &api.V0043Uint32NoValStruct{
			Set:    &setTrue,
			Number: &requeuePriority,
		}
	}

	return apiJob, nil
}