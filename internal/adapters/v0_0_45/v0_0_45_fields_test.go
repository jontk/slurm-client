// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_45

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	types "github.com/jontk/slurm-client/api"
	api "github.com/jontk/slurm-client/internal/openapi/v0_0_45"
)

func TestPartitionConversion_V0045Fields(t *testing.T) {
	adapter := NewPartitionAdapter(nil)
	name := "debug"
	excl := "NO"
	over := "YES:4"
	total := int32(8)
	binding := []api.V0045PartitionInfoCpusTaskBinding{
		api.V0045PartitionInfoCpusTaskBindingCPUBINDTOCORES,
		api.V0045PartitionInfoCpusTaskBindingVERBOSE,
	}
	flags := []api.V0045PartitionInfoFlags{api.DEFAULT, api.HIDDEN}
	preempt := []api.V0045PartitionInfoPreemptMode{api.V0045PartitionInfoPreemptModeSUSPEND}
	apiObj := api.V0045PartitionInfo{Name: &name, Flags: &flags, PreemptMode: &preempt}
	apiObj.Cpus = &struct {
		TaskBinding *[]api.V0045PartitionInfoCpusTaskBinding `json:"task_binding,omitempty"`
		Total       *int32                                   `json:"total,omitempty"`
	}{TaskBinding: &binding, Total: &total}
	apiObj.Partition = &struct {
		Exclusive     *string                                 `json:"exclusive,omitempty"`
		Oversubscribe *string                                 `json:"oversubscribe,omitempty"`
		State         *[]api.V0045PartitionInfoPartitionState `json:"state,omitempty"`
	}{Exclusive: &excl, Oversubscribe: &over}

	got := adapter.convertAPIPartitionToCommon(apiObj)
	require.NotNil(t, got)
	require.NotNil(t, got.CPUs)
	assert.Nil(t, got.CPUs.TaskBinding, "integer task_binding does not exist in v0.0.45")
	assert.Equal(t, []types.TaskBindingValue{types.TaskBindingCPUBindToCores, types.TaskBindingVerbose}, got.CPUs.TaskBindingFlags)
	assert.Equal(t, &total, got.CPUs.Total)
	require.NotNil(t, got.Partition)
	assert.Equal(t, &excl, got.Partition.Exclusive)
	assert.Equal(t, &over, got.Partition.Oversubscribe)
	assert.Equal(t, []types.PartitionDefaultsFlagsValue{types.PartitionDefaultsFlagsDefault, types.PartitionDefaultsFlagsHidden}, got.Flags)
	assert.Equal(t, []types.PreemptModeValue{types.PreemptModeSuspend}, got.PreemptMode)
}

func TestNodeConversion_V0045SuspendTime(t *testing.T) {
	adapter := NewNodeAdapter(nil)
	name := "n1"
	set := true
	secs := int32(300)
	apiObj := api.V0045Node{Name: &name, SuspendTime: &api.V0045Uint32NoValStruct{Set: &set, Number: &secs}}

	got := adapter.convertAPINodeToCommon(apiObj)
	require.NotNil(t, got)
	assert.Equal(t, time.Unix(300, 0), got.SuspendTime)
}

func TestJobConversion_V0045Fields(t *testing.T) {
	adapter := NewJobAdapter(nil)
	ct := "docker"
	delay := int32(15)
	margin := int32(10)
	apiObj := api.V0045JobInfo{ContainerType: &ct, MemoryUpdateDelay: &delay, MemoryUpdateMargin: &margin}

	got := adapter.convertAPIJobToCommon(apiObj)
	require.NotNil(t, got)
	assert.Equal(t, &ct, got.ContainerType)
	assert.Equal(t, &delay, got.MemoryUpdateDelay)
	assert.Equal(t, &margin, got.MemoryUpdateMargin)
}

func TestJobCreateConversion_V0045Fields(t *testing.T) {
	ct := "docker"
	cores := int32(4)
	delay := int32(15)
	margin := int32(10)
	in := &types.JobCreate{ContainerType: &ct, CoresPerSocket: &cores, MemoryUpdateDelay: &delay, MemoryUpdateMargin: &margin}

	got := jobWriteConverter.ConvertCommonJobCreateToAPI(in)
	require.NotNil(t, got)
	assert.Equal(t, &ct, got.ContainerType)
	assert.Equal(t, &cores, got.CoresPerSocket)
	assert.Equal(t, &delay, got.MemoryUpdateDelay)
	assert.Equal(t, &margin, got.MemoryUpdateMargin)
}

func TestPartitionFields_ToAPI(t *testing.T) {
	state := types.PartitionStateUp
	f := partitionFieldsFromCreate(&types.PartitionCreate{
		Name:          "batch",
		Nodes:         "n[1-4]",
		AllowAccounts: []string{"a", "b"},
		MaxTime:       60,
		DefaultTime:   10,
		Priority:      5,
		PreemptMode:   []string{"suspend"},
		State:         state,
		Hidden:        true,
		RootOnly:      true,
		JobDefaults:   map[string]string{"b": "2", "a": "1"},
	})
	got, err := f.toAPI()
	require.NoError(t, err)
	require.NotNil(t, got.Name)
	assert.Equal(t, "batch", *got.Name)
	assert.Equal(t, "n[1-4]", *got.Nodes.Configured)
	assert.Equal(t, "a,b", *got.Accounts.Allowed)
	assert.Equal(t, int32(60), *got.Maximums.Time.Number)
	assert.True(t, *got.Maximums.Time.Set)
	assert.Equal(t, int32(10), *got.Defaults.Time.Number)
	assert.Equal(t, "a=1,b=2", *got.Defaults.Job)
	assert.Equal(t, int32(5), *got.Priority.Tier)
	assert.Equal(t, int32(5), *got.Priority.JobFactor)
	assert.Equal(t, []api.V0045PartitionInfoPreemptMode{api.V0045PartitionInfoPreemptModeSUSPEND}, *got.PreemptMode)
	assert.Equal(t, []api.V0045PartitionInfoPartitionState{api.V0045PartitionInfoPartitionState("UP")}, *got.Partition.State)
	assert.ElementsMatch(t, []api.V0045PartitionInfoFlags{api.HIDDEN, api.ROOTONLY}, *got.Flags)
	assert.Nil(t, got.Minimums)
}

func TestPartitionFields_UpdateOnlySetsGivenFields(t *testing.T) {
	maxTime := int32(0) // explicit zero must still be sent on update
	f := partitionFieldsFromUpdate("batch", &types.PartitionUpdate{MaxTime: &maxTime})
	got, err := f.toAPI()
	require.NoError(t, err)
	require.NotNil(t, got.Maximums)
	assert.Equal(t, int32(0), *got.Maximums.Time.Number)
	assert.Nil(t, got.Nodes)
	assert.Nil(t, got.Accounts)
	assert.Nil(t, got.Partition)
}

func TestPartitionFields_UnsupportedSettingsRejected(t *testing.T) {
	f := partitionFieldsFromCreate(&types.PartitionCreate{Name: "x", TresStr: "cpu=4", PowerDownOnIdle: true})
	_, err := f.toAPI()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "tres_str")
	assert.Contains(t, err.Error(), "power_down_on_idle")
}

func newTestServerAdapterClient(t *testing.T, h http.HandlerFunc) *api.ClientWithResponses {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	c, err := api.NewClientWithResponses(srv.URL)
	require.NoError(t, err)
	return c
}

func TestJobRequeue_UsesRequeueEndpoint(t *testing.T) {
	var method, path string
	client := newTestServerAdapterClient(t, func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"errors":[],"warnings":[]}`))
	})
	require.NoError(t, NewJobAdapter(client).Requeue(context.Background(), 42))
	assert.Equal(t, http.MethodGet, method)
	assert.Equal(t, "/slurm/v0.0.45/job/42/requeue", path)
}

func TestPartitionWrite_Endpoints(t *testing.T) {
	type call struct{ method, path, body string }
	var calls []call
	client := newTestServerAdapterClient(t, func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		calls = append(calls, call{r.Method, r.URL.Path, string(b)})
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"errors":[],"warnings":[]}`))
	})
	a := NewPartitionAdapter(client)
	ctx := context.Background()

	resp, err := a.Create(ctx, &types.PartitionCreate{Name: "p1", Nodes: "n1", MaxTime: 30})
	require.NoError(t, err)
	assert.Equal(t, "p1", resp.PartitionName)
	mt := int32(60)
	require.NoError(t, a.Update(ctx, "p1", &types.PartitionUpdate{MaxTime: &mt}))
	require.NoError(t, a.Delete(ctx, "p1"))

	require.Len(t, calls, 3)
	assert.Equal(t, http.MethodPost, calls[0].method)
	assert.Equal(t, "/slurm/v0.0.45/partitions/", calls[0].path)
	assert.Contains(t, calls[0].body, `"name":"p1"`)
	assert.Contains(t, calls[0].body, `"configured":"n1"`)
	assert.Equal(t, http.MethodPost, calls[1].method)
	assert.Contains(t, calls[1].body, `"number":60`)
	assert.Equal(t, http.MethodDelete, calls[2].method)
	assert.Equal(t, "/slurm/v0.0.45/partition/p1", calls[2].path)
}
