// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

//go:build realserver

// Package realserver contains tests that run against a live slurmrestd.
//
// Run with:
//
//	SLURM_SERVER_URL=http://127.0.0.1:6820 SLURM_JWT_TOKEN=... [SLURM_USER=slurm] [SLURM_API_VERSION=v0.0.45] \
//	  go test -tags realserver ./tests/realserver/ -v
//
// Without SLURM_SERVER_URL and SLURM_JWT_TOKEN the tests are skipped. Leave SLURM_API_VERSION unset
// to exercise version auto-detection. The tests submit and cancel jobs on the "debug" partition if it exists.
package realserver

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	slurm "github.com/jontk/slurm-client"
)

func newClient(t *testing.T) (slurm.SlurmClient, context.Context) {
	t.Helper()
	url, token := os.Getenv("SLURM_SERVER_URL"), os.Getenv("SLURM_JWT_TOKEN")
	if url == "" || token == "" {
		t.Skip("set SLURM_SERVER_URL and SLURM_JWT_TOKEN to run real-server tests")
	}
	user := os.Getenv("SLURM_USER")
	if user == "" {
		user = "slurm"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	t.Cleanup(cancel)
	opts := []slurm.ClientOption{slurm.WithBaseURL(url), slurm.WithUserToken(user, token)}
	var (
		c   slurm.SlurmClient
		err error
	)
	if v := os.Getenv("SLURM_API_VERSION"); v != "" {
		c, err = slurm.NewClientWithVersion(ctx, v, opts...)
	} else {
		c, err = slurm.NewClient(ctx, opts...)
	}
	require.NoError(t, err)
	t.Cleanup(func() { _ = c.Close() })
	t.Logf("client API version: %s", c.Version())
	return c, ctx
}

func TestRealServer_ReadOnly(t *testing.T) {
	c, ctx := newClient(t)

	require.NoError(t, c.Info().Ping(ctx))
	info, err := c.Info().Get(ctx)
	require.NoError(t, err)
	t.Logf("cluster=%s slurm=%s", info.ClusterName, info.Version)

	parts, err := c.Partitions().List(ctx, nil)
	require.NoError(t, err)
	assert.NotEmpty(t, parts.Partitions)

	nodes, err := c.Nodes().List(ctx, nil)
	require.NoError(t, err)
	assert.NotEmpty(t, nodes.Nodes)

	_, err = c.Jobs().List(ctx, nil)
	require.NoError(t, err)

	// slurmdb-backed resources (require slurmdbd)
	_, err = c.Accounts().List(ctx, nil)
	assert.NoError(t, err)
	_, err = c.Users().List(ctx, nil)
	assert.NoError(t, err)
	_, err = c.QoS().List(ctx, nil)
	assert.NoError(t, err)
	_, err = c.Clusters().List(ctx, nil)
	assert.NoError(t, err)
	_, err = c.Associations().List(ctx, nil)
	assert.NoError(t, err)
}

func TestRealServer_JobLifecycle(t *testing.T) {
	c, ctx := newClient(t)

	name, script, wd := "realserver-test", "#!/bin/bash\nsleep 300\n", "/tmp"
	resp, err := c.Jobs().SubmitRaw(ctx, &slurm.JobCreate{
		Name: &name, Script: &script, CurrentWorkingDirectory: &wd,
		Environment: []string{"PATH=/usr/bin:/bin"},
	})
	require.NoError(t, err)
	id := fmt.Sprint(resp.JobId)
	t.Cleanup(func() { _ = c.Jobs().Cancel(context.Background(), id) })

	job, err := c.Jobs().Get(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, job.Name)
	assert.Equal(t, name, *job.Name)

	require.NoError(t, c.Jobs().Cancel(ctx, id))
}

// TestRealServer_JobRequeue checks that Requeue puts a running job back in the queue
// instead of cancelling it. Only API v0.0.45+ has a real requeue endpoint.
func TestRealServer_JobRequeue(t *testing.T) {
	c, ctx := newClient(t)
	if c.Version() != "v0.0.45" {
		t.Skipf("requeue endpoint requires v0.0.45, client is %s", c.Version())
	}
	name, script, wd := "realserver-requeue", "#!/bin/bash\nsleep 300\n", "/tmp"
	resp, err := c.Jobs().SubmitRaw(ctx, &slurm.JobCreate{
		Name: &name, Script: &script, CurrentWorkingDirectory: &wd,
		Environment: []string{"PATH=/usr/bin:/bin"},
	})
	require.NoError(t, err)
	id := fmt.Sprint(resp.JobId)
	t.Cleanup(func() { _ = c.Jobs().Cancel(context.Background(), id) })

	time.Sleep(3 * time.Second) // let the job start
	require.NoError(t, c.Jobs().Requeue(ctx, id))
	time.Sleep(2 * time.Second)

	job, err := c.Jobs().Get(ctx, id)
	require.NoError(t, err)
	assert.NotContains(t, fmt.Sprint(job.JobState), "CANCELLED", "requeue must not cancel the job")
}

// TestRealServer_PartitionManagement checks partition create/update/delete (API v0.0.45+).
func TestRealServer_PartitionManagement(t *testing.T) {
	c, ctx := newClient(t)
	if c.Version() != "v0.0.45" {
		t.Skipf("partition management requires v0.0.45, client is %s", c.Version())
	}
	nodes, err := c.Nodes().List(ctx, nil)
	require.NoError(t, err)
	require.NotEmpty(t, nodes.Nodes)
	require.NotNil(t, nodes.Nodes[0].Name)

	const part = "realserver-part"
	t.Cleanup(func() { _ = c.Partitions().Delete(context.Background(), part) })

	_, err = c.Partitions().Create(ctx, &slurm.PartitionCreate{
		Name: part, Nodes: *nodes.Nodes[0].Name, MaxTime: 45, Hidden: true, PreemptMode: []string{"suspend"},
	})
	require.NoError(t, err)

	p, err := c.Partitions().Get(ctx, part)
	require.NoError(t, err)
	require.NotNil(t, p.Maximums)
	require.NotNil(t, p.Maximums.Time)
	assert.Equal(t, uint32(45), *p.Maximums.Time)
	assert.Contains(t, fmt.Sprint(p.Flags), "HIDDEN")
	assert.Contains(t, fmt.Sprint(p.PreemptMode), "SUSPEND")

	maxTime := int32(90)
	require.NoError(t, c.Partitions().Update(ctx, part, &slurm.PartitionUpdate{MaxTime: &maxTime}))
	p, err = c.Partitions().Get(ctx, part)
	require.NoError(t, err)
	assert.Equal(t, uint32(90), *p.Maximums.Time)

	require.NoError(t, c.Partitions().Delete(ctx, part))
	_, err = c.Partitions().Get(ctx, part)
	assert.Error(t, err)

	_, err = c.Partitions().Create(ctx, &slurm.PartitionCreate{Name: "bad", TresStr: "cpu=1"})
	assert.Error(t, err, "settings v0.0.45 cannot express must be rejected")
}
