// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_45

import (
	"context"
	"fmt"
	"sort"
	"strings"

	types "github.com/jontk/slurm-client/api"
	"github.com/jontk/slurm-client/internal/common"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_45"
	"github.com/jontk/slurm-client/pkg/errors"
)

// Partition create/update/delete were added to the REST API in Slurm 26.05
// (POST /slurm/v0.0.45/partitions, DELETE /slurm/v0.0.45/partition/{name}).
// These methods are listed under manual_methods in tools/codegen/adapter_config.yaml
// so regeneration does not replace them with "unsupported" stubs.

// partitionFields is the version-neutral set of writable partition settings.
// A nil pointer or empty slice means "not specified".
type partitionFields struct {
	name                string
	nodes               *string
	allowAllocNodes     *string
	allowAccounts       []string
	denyAccounts        []string
	allowGroups         []string
	allowQoS            []string
	denyQoS             []string
	qos                 *string
	defaultTime         *int32
	defaultMemPerCPU    *int64
	defaultMemPerNode   *int64
	graceTime           *int32
	maxCPUsPerNode      *int32
	maxMemPerCPU        *int64
	maxMemPerNode       *int64
	maxNodes            *int32
	maxTime             *int32
	minNodes            *int32
	overTimeLimit       *int32
	preemptMode         []string
	priority            *int32
	priorityJobFactor   *int32
	priorityTier        *int32
	state               *types.PartitionState
	billingWeightStr    *string
	jobDefaults         map[string]string
	resumeTimeout       *int32
	suspendTime         *int32
	suspendTimeout      *int32
	hidden              *bool
	exclusiveUser       *bool
	lln                 *bool
	rootOnly            *bool
	reqResv             *bool
	unsupportedSettings []string // settings the caller specified that v0.0.45 cannot express
}

func ptrIfNonZero[T comparable](v T) *T {
	var zero T
	if v == zero {
		return nil
	}
	return &v
}

func strsOrNil(v string) *string { return ptrIfNonZero(v) }

func boolIfTrue(v bool) *bool {
	if !v {
		return nil
	}
	return &v
}

func partitionFieldsFromCreate(c *types.PartitionCreate) partitionFields {
	f := partitionFields{
		name:              c.Name,
		nodes:             strsOrNil(c.Nodes),
		allowAllocNodes:   strsOrNil(c.AllowAllocNodes),
		allowAccounts:     c.AllowAccounts,
		denyAccounts:      c.DenyAccounts,
		allowGroups:       c.AllowGroups,
		allowQoS:          c.AllowQoS,
		denyQoS:           c.DenyQoS,
		qos:               strsOrNil(c.QoS),
		defaultTime:       ptrIfNonZero(c.DefaultTime),
		defaultMemPerCPU:  ptrIfNonZero(c.DefaultMemPerCPU),
		defaultMemPerNode: ptrIfNonZero(c.DefaultMemPerNode),
		graceTime:         ptrIfNonZero(c.GraceTime),
		maxCPUsPerNode:    ptrIfNonZero(c.MaxCPUsPerNode),
		maxMemPerCPU:      ptrIfNonZero(c.MaxMemPerCPU),
		maxMemPerNode:     ptrIfNonZero(c.MaxMemPerNode),
		maxNodes:          ptrIfNonZero(c.MaxNodes),
		maxTime:           ptrIfNonZero(c.MaxTime),
		minNodes:          ptrIfNonZero(c.MinNodes),
		overTimeLimit:     ptrIfNonZero(c.OverTimeLimit),
		preemptMode:       c.PreemptMode,
		priority:          ptrIfNonZero(c.Priority),
		priorityJobFactor: ptrIfNonZero(c.PriorityJobFactor),
		priorityTier:      ptrIfNonZero(c.PriorityTier),
		state:             ptrIfNonZero(c.State),
		billingWeightStr:  strsOrNil(c.BillingWeightStr),
		jobDefaults:       c.JobDefaults,
		resumeTimeout:     ptrIfNonZero(c.ResumeTimeout),
		suspendTime:       ptrIfNonZero(c.SuspendTime),
		suspendTimeout:    ptrIfNonZero(c.SuspendTimeout),
		hidden:            boolIfTrue(c.Hidden),
		exclusiveUser:     boolIfTrue(c.ExclusiveUser),
		lln:               boolIfTrue(c.LLN),
		rootOnly:          boolIfTrue(c.RootOnly),
		reqResv:           boolIfTrue(c.ReqResv),
	}
	// DefMemPerNode is the older spelling of DefaultMemPerNode.
	if f.defaultMemPerNode == nil {
		f.defaultMemPerNode = ptrIfNonZero(c.DefMemPerNode)
	}
	if c.AllocNodes != "" && f.allowAllocNodes == nil {
		f.allowAllocNodes = strsOrNil(c.AllocNodes)
	}
	f.unsupportedSettings = unsupportedPartitionSettings(c.TresStr != "", len(c.SelectTypeParameters) > 0, c.PowerDownOnIdle)
	return f
}

func partitionFieldsFromUpdate(name string, u *types.PartitionUpdate) partitionFields {
	f := partitionFields{
		name:              name,
		nodes:             u.Nodes,
		allowAllocNodes:   u.AllowAllocNodes,
		allowAccounts:     u.AllowAccounts,
		denyAccounts:      u.DenyAccounts,
		allowGroups:       u.AllowGroups,
		allowQoS:          u.AllowQoS,
		denyQoS:           u.DenyQoS,
		qos:               u.QoS,
		defaultTime:       u.DefaultTime,
		defaultMemPerCPU:  u.DefaultMemPerCPU,
		defaultMemPerNode: u.DefaultMemPerNode,
		graceTime:         u.GraceTime,
		maxCPUsPerNode:    u.MaxCPUsPerNode,
		maxMemPerCPU:      u.MaxMemPerCPU,
		maxMemPerNode:     u.MaxMemPerNode,
		maxNodes:          u.MaxNodes,
		maxTime:           u.MaxTime,
		minNodes:          u.MinNodes,
		overTimeLimit:     u.OverTimeLimit,
		preemptMode:       u.PreemptMode,
		priority:          u.Priority,
		priorityJobFactor: u.PriorityJobFactor,
		priorityTier:      u.PriorityTier,
		state:             u.State,
		billingWeightStr:  u.BillingWeightStr,
		jobDefaults:       u.JobDefaults,
		resumeTimeout:     u.ResumeTimeout,
		suspendTime:       u.SuspendTime,
		suspendTimeout:    u.SuspendTimeout,
		hidden:            u.Hidden,
		exclusiveUser:     u.ExclusiveUser,
		lln:               u.LLN,
		rootOnly:          u.RootOnly,
		reqResv:           u.ReqResv,
	}
	if f.allowAllocNodes == nil {
		f.allowAllocNodes = u.AllocNodes
	}
	if f.defaultMemPerNode == nil {
		f.defaultMemPerNode = u.DefMemPerNode
	}
	f.unsupportedSettings = unsupportedPartitionSettings(u.TresStr != nil, len(u.SelectTypeParameters) > 0, u.PowerDownOnIdle != nil)
	return f
}

func unsupportedPartitionSettings(tres, selectTypeParams, powerDownOnIdle bool) []string {
	var out []string
	if tres {
		out = append(out, "tres_str")
	}
	if selectTypeParams {
		out = append(out, "select_type_parameters")
	}
	if powerDownOnIdle {
		out = append(out, "power_down_on_idle")
	}
	return out
}

func noVal32(v *int32) *api.V0045Uint32NoValStruct {
	if v == nil {
		return nil
	}
	set := true
	return &api.V0045Uint32NoValStruct{Set: &set, Number: v}
}

func noVal16(v *int32) *api.V0045Uint16NoValStruct {
	if v == nil {
		return nil
	}
	set := true
	n := int32(*v)
	return &api.V0045Uint16NoValStruct{Set: &set, Number: &n}
}

func noVal64(v *int64) *api.V0045Uint64NoValStruct {
	if v == nil {
		return nil
	}
	set := true
	n := *v
	return &api.V0045Uint64NoValStruct{Set: &set, Number: &n}
}

func csvOrNil(v []string) *string {
	if len(v) == 0 {
		return nil
	}
	s := strings.Join(v, ",")
	return &s
}

func jobDefaultsString(m map[string]string) *string {
	if len(m) == 0 {
		return nil
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, k+"="+m[k])
	}
	s := strings.Join(parts, ",")
	return &s
}

// toAPI converts the neutral field set to the v0.0.45 request item.
func (f partitionFields) toAPI() (*api.V0045PartitionInfo, error) {
	if len(f.unsupportedSettings) > 0 {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed,
			fmt.Sprintf("partition settings not supported by API v0.0.45: %s", strings.Join(f.unsupportedSettings, ", ")),
			"partition", nil, nil)
	}
	name := f.name
	out := &api.V0045PartitionInfo{Name: &name}

	if f.nodes != nil || f.allowAllocNodes != nil {
		out.Nodes = &struct {
			AllowedAllocation *string `json:"allowed_allocation,omitempty"`
			Configured        *string `json:"configured,omitempty"`
			Total             *int32  `json:"total,omitempty"`
		}{AllowedAllocation: f.allowAllocNodes, Configured: f.nodes}
	}
	if a, d := csvOrNil(f.allowAccounts), csvOrNil(f.denyAccounts); a != nil || d != nil {
		out.Accounts = &struct {
			Allowed *string `json:"allowed,omitempty"`
			Deny    *string `json:"deny,omitempty"`
		}{Allowed: a, Deny: d}
	}
	if g := csvOrNil(f.allowGroups); g != nil {
		out.Groups = &struct {
			Allowed *string `json:"allowed,omitempty"`
		}{Allowed: g}
	}
	if a, d := csvOrNil(f.allowQoS), csvOrNil(f.denyQoS); a != nil || d != nil || f.qos != nil {
		out.Qos = &struct {
			Allowed  *string `json:"allowed,omitempty"`
			Assigned *string `json:"assigned,omitempty"`
			Deny     *string `json:"deny,omitempty"`
		}{Allowed: a, Assigned: f.qos, Deny: d}
	}
	if jd := jobDefaultsString(f.jobDefaults); f.defaultTime != nil || f.defaultMemPerCPU != nil || f.defaultMemPerNode != nil || jd != nil {
		out.Defaults = &struct {
			Job                    *string                     `json:"job,omitempty"`
			MemoryPerCpu           *int64                      `json:"memory_per_cpu,omitempty"`
			PartitionMemoryPerCpu  *api.V0045Uint64NoValStruct `json:"partition_memory_per_cpu,omitempty"`
			PartitionMemoryPerNode *api.V0045Uint64NoValStruct `json:"partition_memory_per_node,omitempty"`
			Time                   *api.V0045Uint32NoValStruct `json:"time,omitempty"`
		}{Job: jd, PartitionMemoryPerCpu: noVal64(f.defaultMemPerCPU), PartitionMemoryPerNode: noVal64(f.defaultMemPerNode), Time: noVal32(f.defaultTime)}
	}
	if f.maxCPUsPerNode != nil || f.maxMemPerCPU != nil || f.maxMemPerNode != nil || f.maxNodes != nil || f.maxTime != nil || f.overTimeLimit != nil {
		out.Maximums = &struct {
			CpusPerNode   *api.V0045Uint32NoValStruct `json:"cpus_per_node,omitempty"`
			CpusPerSocket *api.V0045Uint32NoValStruct `json:"cpus_per_socket,omitempty"`
			MemoryPerCpu  *int64                      `json:"memory_per_cpu,omitempty"`
			Nodes         *api.V0045Uint32NoValStruct `json:"nodes,omitempty"`
			OverTimeLimit *api.V0045Uint16NoValStruct `json:"over_time_limit,omitempty"`
			Oversubscribe *struct {
				Flags *[]api.V0045PartitionInfoMaximumsOversubscribeFlags `json:"flags,omitempty"`
				Jobs  *int32                                              `json:"jobs,omitempty"`
			} `json:"oversubscribe,omitempty"`
			PartitionMemoryPerCpu  *api.V0045Uint64NoValStruct `json:"partition_memory_per_cpu,omitempty"`
			PartitionMemoryPerNode *api.V0045Uint64NoValStruct `json:"partition_memory_per_node,omitempty"`
			Shares                 *int32                      `json:"shares,omitempty"`
			Time                   *api.V0045Uint32NoValStruct `json:"time,omitempty"`
		}{
			CpusPerNode:            noVal32(f.maxCPUsPerNode),
			Nodes:                  noVal32(f.maxNodes),
			OverTimeLimit:          noVal16(f.overTimeLimit),
			PartitionMemoryPerCpu:  noVal64(f.maxMemPerCPU),
			PartitionMemoryPerNode: noVal64(f.maxMemPerNode),
			Time:                   noVal32(f.maxTime),
		}
	}
	if f.minNodes != nil {
		out.Minimums = &struct {
			Nodes *int32 `json:"nodes,omitempty"`
		}{Nodes: f.minNodes}
	}
	if f.graceTime != nil {
		out.GraceTime = f.graceTime
	}
	if len(f.preemptMode) > 0 {
		modes := make([]api.V0045PartitionInfoPreemptMode, len(f.preemptMode))
		for i, m := range f.preemptMode {
			modes[i] = api.V0045PartitionInfoPreemptMode(strings.ToUpper(m))
		}
		out.PreemptMode = &modes
	}
	// "Priority" sets both the tier and the job factor in slurm.conf.
	tier, factor := f.priorityTier, f.priorityJobFactor
	if f.priority != nil {
		if tier == nil {
			tier = f.priority
		}
		if factor == nil {
			factor = f.priority
		}
	}
	if tier != nil || factor != nil {
		out.Priority = &struct {
			JobFactor *int32 `json:"job_factor,omitempty"`
			Tier      *int32 `json:"tier,omitempty"`
		}{JobFactor: factor, Tier: tier}
	}
	if f.state != nil || f.exclusiveUser != nil {
		p := &struct {
			Exclusive     *string                                 `json:"exclusive,omitempty"`
			Oversubscribe *string                                 `json:"oversubscribe,omitempty"`
			State         *[]api.V0045PartitionInfoPartitionState `json:"state,omitempty"`
		}{}
		if f.state != nil {
			st := []api.V0045PartitionInfoPartitionState{api.V0045PartitionInfoPartitionState(strings.ToUpper(string(*f.state)))}
			p.State = &st
		}
		if f.exclusiveUser != nil && *f.exclusiveUser {
			user := "USER"
			p.Exclusive = &user
		}
		out.Partition = p
	}
	var flags []api.V0045PartitionInfoFlags
	for _, fl := range []struct {
		on   *bool
		flag api.V0045PartitionInfoFlags
	}{
		{f.hidden, api.HIDDEN},
		{f.lln, api.LLN},
		{f.rootOnly, api.ROOTONLY},
		{f.reqResv, api.REQRESV},
	} {
		if fl.on != nil && *fl.on {
			flags = append(flags, fl.flag)
		}
	}
	if len(flags) > 0 {
		out.Flags = &flags
	}
	if f.billingWeightStr != nil {
		out.Tres = &struct {
			BillingWeights *string `json:"billing_weights,omitempty"`
			Configured     *string `json:"configured,omitempty"`
		}{BillingWeights: f.billingWeightStr}
	}
	if f.suspendTime != nil {
		out.SuspendTime = noVal32(f.suspendTime)
	}
	if f.resumeTimeout != nil || f.suspendTimeout != nil {
		out.Timeouts = &struct {
			Resume  *api.V0045Uint16NoValStruct `json:"resume,omitempty"`
			Suspend *api.V0045Uint16NoValStruct `json:"suspend,omitempty"`
		}{Resume: noVal16(f.resumeTimeout), Suspend: noVal16(f.suspendTimeout)}
	}
	return out, nil
}

// postPartition sends one partition description to POST /partitions.
func (a *PartitionAdapter) postPartition(ctx context.Context, f partitionFields) error {
	item, err := f.toAPI()
	if err != nil {
		return err
	}
	resp, err := a.client.SlurmV0045PostPartitionsWithResponse(ctx, api.SlurmV0045PostPartitionsJSONRequestBody{
		Partitions: api.V0045UpdatePartitionMsgList{*item},
	})
	if err != nil {
		return a.HandleAPIError(err)
	}
	var apiErrors *api.V0045OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	} else if resp.JSONDefault != nil {
		apiErrors = resp.JSONDefault.Errors
	}
	return common.HandleAPIResponse(api.NewResponseAdapter(resp.StatusCode(), apiErrors), "v0.0.45")
}

// Create creates a partition (POST /slurm/v0.0.45/partitions).
func (a *PartitionAdapter) Create(ctx context.Context, partition *types.PartitionCreate) (*types.PartitionCreateResponse, error) {
	if err := a.ValidateContext(ctx); err != nil {
		return nil, err
	}
	if partition == nil {
		return nil, errors.NewValidationError(errors.ErrorCodeValidationFailed, "partition is required", "partition", nil, nil)
	}
	if err := a.ValidateResourceName(partition.Name, "name"); err != nil {
		return nil, err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return nil, err
	}
	if err := a.postPartition(ctx, partitionFieldsFromCreate(partition)); err != nil {
		return nil, err
	}
	return &types.PartitionCreateResponse{PartitionName: partition.Name}, nil
}

// Update updates a partition. The API applies the same POST used for creation,
// so only the fields set on update are sent (Slurm applies them to the named partition).
func (a *PartitionAdapter) Update(ctx context.Context, partitionName string, update *types.PartitionUpdate) error {
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.ValidateResourceName(partitionName, "partitionName"); err != nil {
		return err
	}
	if update == nil {
		return errors.NewValidationError(errors.ErrorCodeValidationFailed, "update is required", "update", nil, nil)
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}
	return a.postPartition(ctx, partitionFieldsFromUpdate(partitionName, update))
}

// Delete deletes a partition (DELETE /slurm/v0.0.45/partition/{name}).
func (a *PartitionAdapter) Delete(ctx context.Context, partitionName string) error {
	if err := a.ValidateContext(ctx); err != nil {
		return err
	}
	if err := a.ValidateResourceName(partitionName, "partitionName"); err != nil {
		return err
	}
	if err := a.CheckClientInitialized(a.client); err != nil {
		return err
	}
	resp, err := a.client.SlurmV0045DeletePartitionWithResponse(ctx, partitionName)
	if err != nil {
		return a.HandleAPIError(err)
	}
	var apiErrors *api.V0045OpenapiErrors
	if resp.JSON200 != nil {
		apiErrors = resp.JSON200.Errors
	} else if resp.JSONDefault != nil {
		apiErrors = resp.JSONDefault.Errors
	}
	if resp.StatusCode() == 204 {
		return nil
	}
	return common.HandleAPIResponse(api.NewResponseAdapter(resp.StatusCode(), apiErrors), "v0.0.45")
}
