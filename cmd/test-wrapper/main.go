// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/pkg/config"
)

// userTokenAuth implements authentication with both username and token headers
type userTokenAuth struct {
	username string
	token    string
}

func (u *userTokenAuth) Authenticate(_ context.Context, req *http.Request) error {
	req.Header.Set("X-SLURM-USER-NAME", u.username)
	req.Header.Set("X-SLURM-USER-TOKEN", u.token)
	return nil
}

func (u *userTokenAuth) Type() string {
	return "user-token"
}

func main() {
	// Check if JWT token is provided as environment variable
	jwtToken := os.Getenv("SLURM_JWT")
	if jwtToken == "" {
		fmt.Println("Error: SLURM_JWT environment variable not set")
		fmt.Println("Usage: SLURM_JWT='your-jwt-token' go run test-wrapper.go")
		os.Exit(1)
	}

	username := os.Getenv("SLURM_USER")
	if username == "" {
		username = "root" // Default username for testing
	}

	// Create client
	client, err := createClient(jwtToken, username)
	if err != nil {
		fmt.Printf("Failed to create client: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	runTests(ctx, client)
	fmt.Println("\n=== All Tests Completed ===")
}

func createClient(jwtToken string, username string) (slurm.SlurmClient, error) {
	cfg := &config.Config{
		BaseURL: "http://localhost:6820/slurm",
		Debug:   true,
	}

	authProvider := &userTokenAuth{
		username: username,
		token:    jwtToken,
	}
	ctx := context.Background()
	return slurm.NewClientWithVersion(ctx, "v0.0.40",
		slurm.WithConfig(cfg),
		slurm.WithAuth(authProvider),
	)
}

func runTests(ctx context.Context, client slurm.SlurmClient) {
	runPingTest(ctx, client)
	runVersionTest(ctx, client)
	runClusterInfoTest(ctx, client)
	runJobListTest(ctx, client)
	runNodeListTest(ctx, client)
	runPartitionListTest(ctx, client)
	runAccountListTest(ctx, client)
	runUserListTest(ctx, client)
}

func runPingTest(ctx context.Context, client slurm.SlurmClient) {
	fmt.Println("\n=== Test 1: Ping ===")
	infoMgr := client.Info()
	err := infoMgr.Ping(ctx)
	if err != nil {
		fmt.Printf("Ping failed: %v\n", err)
	} else {
		fmt.Printf("Ping successful!\n")
	}
}

func runVersionTest(_ context.Context, client slurm.SlurmClient) {
	fmt.Println("\n=== Test 2: Get API Version ===")
	fmt.Printf("Client API Version: %s\n", client.Version())
}

func runClusterInfoTest(ctx context.Context, client slurm.SlurmClient) {
	fmt.Println("\n=== Test 2b: Get Cluster Info ===")
	infoMgr := client.Info()
	info, err := infoMgr.Get(ctx)
	if err != nil {
		fmt.Printf("Get cluster info failed: %v\n", err)
	} else {
		fmt.Printf("Cluster Info: Name=%s, Version=%s, API=%s\n", info.ClusterName, info.Version, info.APIVersion)
	}
}

func runJobListTest(ctx context.Context, client slurm.SlurmClient) {
	fmt.Println("\n=== Test 3: List Jobs ===")
	jobMgr := client.Jobs()
	jobs, err := jobMgr.List(ctx, nil)
	if err != nil {
		fmt.Printf("List jobs failed: %v\n", err)
		return
	}
	fmt.Printf("Found %d jobs\n", len(jobs.Jobs))
	for i, job := range jobs.Jobs {
		if i < 5 {
			jobID := "N/A"
			if job.JobID != nil {
				jobID = fmt.Sprintf("%d", *job.JobID)
			}
			jobName := "N/A"
			if job.Name != nil {
				jobName = *job.Name
			}
			jobState := "N/A"
			if len(job.JobState) > 0 {
				jobState = string(job.JobState[0])
			}
			fmt.Printf("  Job: ID=%s, Name=%s, State=%s\n", jobID, jobName, jobState)
		}
	}
	if len(jobs.Jobs) > 5 {
		fmt.Printf("  ... and %d more jobs\n", len(jobs.Jobs)-5)
	}
}

func runNodeListTest(ctx context.Context, client slurm.SlurmClient) {
	fmt.Println("\n=== Test 4: List Nodes ===")
	nodeMgr := client.Nodes()
	nodes, err := nodeMgr.List(ctx, nil)
	if err != nil {
		fmt.Printf("List nodes failed: %v\n", err)
		return
	}
	fmt.Printf("Found %d nodes\n", len(nodes.Nodes))
	for i, node := range nodes.Nodes {
		if i < 5 {
			fmt.Printf("  Node: Name=%v, State=%v, CPUs=%d\n", node.Name, node.State, node.CPUs)
		}
	}
	if len(nodes.Nodes) > 5 {
		fmt.Printf("  ... and %d more nodes\n", len(nodes.Nodes)-5)
	}
}

func runPartitionListTest(ctx context.Context, client slurm.SlurmClient) {
	fmt.Println("\n=== Test 5: List Partitions ===")
	partMgr := client.Partitions()
	partitions, err := partMgr.List(ctx, nil)
	if err != nil {
		fmt.Printf("List partitions failed: %v\n", err)
		return
	}
	fmt.Printf("Found %d partitions\n", len(partitions.Partitions))
	for _, partition := range partitions.Partitions {
		partName := "N/A"
		if partition.Name != nil {
			partName = *partition.Name
		}
		partState := "N/A"
		if partition.Partition != nil && len(partition.Partition.State) > 0 {
			partState = string(partition.Partition.State[0])
		}
		partNodes := "N/A"
		if partition.Nodes != nil && partition.Nodes.Configured != nil {
			partNodes = *partition.Nodes.Configured
		}
		fmt.Printf("  Partition: Name=%s, State=%s, Nodes=%s\n",
			partName, partState, partNodes)
	}
}

func runAccountListTest(ctx context.Context, client slurm.SlurmClient) {
	fmt.Println("\n=== Test 6: List Accounts ===")
	acctMgr := client.Accounts()
	accounts, err := acctMgr.List(ctx, nil)
	if err != nil {
		fmt.Printf("List accounts failed: %v\n", err)
		return
	}
	fmt.Printf("Found %d accounts\n", len(accounts.Accounts))
	for i, account := range accounts.Accounts {
		if i < 5 {
			fmt.Printf("  Account: Name=%s, Description=%s\n",
				account.Name, account.Description)
		}
	}
	if len(accounts.Accounts) > 5 {
		fmt.Printf("  ... and %d more accounts\n", len(accounts.Accounts)-5)
	}
}

func runUserListTest(ctx context.Context, client slurm.SlurmClient) {
	fmt.Println("\n=== Test 7: List Users ===")
	userMgr := client.Users()
	users, err := userMgr.List(ctx, nil)
	if err != nil {
		fmt.Printf("List users failed: %v\n", err)
		return
	}
	fmt.Printf("Found %d users\n", len(users.Users))
	for i, user := range users.Users {
		if i < 5 {
			defaultAccount := "N/A"
			if user.Default != nil && user.Default.Account != nil {
				defaultAccount = *user.Default.Account
			}
			fmt.Printf("  User: Name=%s, DefaultAccount=%s\n",
				user.Name, defaultAccount)
		}
	}
	if len(users.Users) > 5 {
		fmt.Printf("  ... and %d more users\n", len(users.Users)-5)
	}
}
