// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	types "github.com/jontk/slurm-client/api"
	"github.com/jontk/slurm-client/internal/factory"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/pkg/config"
)

// This program analyzes adapter functionality by testing wrapper clients
// and documenting what adapter methods would need to be implemented

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <version>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Supported versions: v0.0.40, v0.0.41, v0.0.42, v0.0.43\n")
		os.Exit(1)
	}

	version := os.Args[1]

	// Get JWT token from environment
	jwtToken := os.Getenv("SLURM_JWT")
	if jwtToken == "" {
		log.Fatal("SLURM_JWT environment variable is required")
	}

	// Create configuration
	cfg := config.NewDefault()
	cfg.BaseURL = "http://localhost:6820"
	cfg.Debug = false

	// Create JWT authentication provider
	authProvider := auth.NewTokenAuth(jwtToken)

	// Create factory
	clientFactory, err := factory.NewClientFactory(
		factory.WithConfig(cfg),
		factory.WithAuth(authProvider),
		factory.WithBaseURL(cfg.BaseURL),
	)
	if err != nil {
		log.Fatalf("Failed to create factory: %v", err)
	}

	// Create client with specific version
	client, err := clientFactory.NewClientWithVersion(context.Background(), version)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Printf("Testing API Version: %s\n", version)
	fmt.Println("=====================================")

	// Test each manager
	testJobsManager(client, version)
	testNodesManager(client, version)
	testPartitionsManager(client, version)
	testAccountsManager(client, version)
	testUsersManager(client, version)
	testQoSManager(client, version)
	testReservationsManager(client, version)
	testAssociationsManager(client, version)

	fmt.Println("\nAdapter Implementation Analysis Complete")
}

func testJobsManager(client types.SlurmClient, version string) {
	ctx := context.Background()
	fmt.Println("\n=== Jobs Manager ===")

	fmt.Println("Methods to implement in adapter:")
	fmt.Println("  - List(ctx, opts) (*JobList, error)")
	fmt.Println("  - Get(ctx, jobID) (*Job, error)")
	fmt.Println("  - Submit(ctx, job) (*JobSubmitResponse, error)")
	fmt.Println("  - Update(ctx, jobID, update) error")
	fmt.Println("  - Cancel(ctx, jobID) error")
	fmt.Println("  - Watch(ctx, opts) (<-chan JobEvent, error)")

	// Test List
	fmt.Print("\nTesting List: ")
	jobs, err := client.Jobs().List(ctx, &types.ListJobsOptions{Limit: 5})
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		jobCount := 0
		if jobs != nil && jobs.Jobs != nil {
			jobCount = len(jobs.Jobs)
		}
		fmt.Printf("✅ Success: Found %d jobs\n", jobCount)
	}

	// Test Submit
	fmt.Print("Testing Submit: ")
	submitJob := &types.JobSubmission{
		Name:       fmt.Sprintf("adapter-analysis-%s-%d", version, time.Now().Unix()),
		Partition:  "normal",
		Script:     "#!/bin/bash\necho 'Testing adapter analysis'\nsleep 5",
		TimeLimit:  1,
		Nodes:      1,
		WorkingDir: "/tmp",
		Environment: map[string]string{
			"PATH": "/usr/bin:/bin",
			"USER": "root",
			"HOME": "/tmp",
		},
	}
	submitResp, err := client.Jobs().Submit(ctx, submitJob)
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		fmt.Printf("✅ Success: Job ID %d\n", submitResp.JobId)

		// Test Get
		fmt.Print("Testing Get: ")
		job, err := client.Jobs().Get(ctx, fmt.Sprintf("%d", submitResp.JobId))
		if err != nil {
			fmt.Printf("❌ Failed: %v\n", err)
		} else {
			fmt.Printf("✅ Success: Got job %v\n", job.Name)
		}

		// Test Update
		fmt.Print("Testing Update: ")
		err = client.Jobs().Update(ctx, fmt.Sprintf("%d", submitResp.JobId), &types.JobUpdate{
			Priority: uint32Ptr(100),
		})
		if err != nil {
			fmt.Printf("❌ Failed: %v\n", err)
		} else {
			fmt.Printf("✅ Success\n")
		}

		// Test Cancel
		fmt.Print("Testing Cancel: ")
		err = client.Jobs().Cancel(ctx, fmt.Sprintf("%d", submitResp.JobId))
		if err != nil {
			fmt.Printf("❌ Failed: %v\n", err)
		} else {
			fmt.Printf("✅ Success\n")
		}
	}

	// Test Watch
	fmt.Print("Testing Watch: ")
	_, err = client.Jobs().Watch(ctx, &types.WatchJobsOptions{})
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		fmt.Printf("✅ Success\n")
	}
}

func testNodesManager(client types.SlurmClient, _ string) {
	ctx := context.Background()
	fmt.Println("\n=== Nodes Manager ===")

	fmt.Println("Methods to implement in adapter:")
	fmt.Println("  - List(ctx, opts) (*NodeList, error)")
	fmt.Println("  - Get(ctx, nodeName) (*Node, error)")
	fmt.Println("  - Update(ctx, nodeName, update) error")
	fmt.Println("  - Watch(ctx, opts) (<-chan NodeEvent, error)")

	// Test List
	fmt.Print("\nTesting List: ")
	nodes, err := client.Nodes().List(ctx, &types.ListNodesOptions{})
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		nodeCount := 0
		var firstNodeName string
		if nodes != nil && nodes.Nodes != nil {
			nodeCount = len(nodes.Nodes)
			if nodeCount > 0 && nodes.Nodes[0].Name != nil {
				firstNodeName = *nodes.Nodes[0].Name
			}
		}
		fmt.Printf("✅ Success: Found %d nodes\n", nodeCount)

		if firstNodeName != "" {
			// Test Get
			fmt.Print("Testing Get: ")
			node, err := client.Nodes().Get(ctx, firstNodeName)
			if err != nil {
				fmt.Printf("❌ Failed: %v\n", err)
			} else {
				fmt.Printf("✅ Success: Got node %v\n", node.Name)
			}

			// Test Update
			fmt.Print("Testing Update: ")
			err = client.Nodes().Update(ctx, firstNodeName, &types.NodeUpdate{
				State: []types.NodeState{types.NodeStateIdle},
			})
			if err != nil {
				fmt.Printf("❌ Failed: %v\n", err)
			} else {
				fmt.Printf("✅ Success\n")
			}
		}
	}

	// Test Watch
	fmt.Print("Testing Watch: ")
	_, err = client.Nodes().Watch(ctx, &types.WatchNodesOptions{})
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		fmt.Printf("✅ Success\n")
	}
}

func testPartitionsManager(client types.SlurmClient, _ string) {
	ctx := context.Background()
	fmt.Println("\n=== Partitions Manager ===")

	fmt.Println("Methods to implement in adapter:")
	fmt.Println("  - List(ctx, opts) (*PartitionList, error)")
	fmt.Println("  - Get(ctx, partitionName) (*Partition, error)")
	fmt.Println("  - Update(ctx, partitionName, update) error")

	// Test List
	fmt.Print("\nTesting List: ")
	partitions, err := client.Partitions().List(ctx, &types.ListPartitionsOptions{})
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		partitionCount := 0
		var firstPartitionName string
		if partitions != nil && partitions.Partitions != nil {
			partitionCount = len(partitions.Partitions)
			if partitionCount > 0 && partitions.Partitions[0].Name != nil {
				firstPartitionName = *partitions.Partitions[0].Name
			}
		}
		fmt.Printf("✅ Success: Found %d partitions\n", partitionCount)

		if firstPartitionName != "" {
			// Test Get
			fmt.Print("Testing Get: ")
			partition, err := client.Partitions().Get(ctx, firstPartitionName)
			if err != nil {
				fmt.Printf("❌ Failed: %v\n", err)
			} else {
				fmt.Printf("✅ Success: Got partition %v\n", partition.Name)
			}

			// Test Update
			fmt.Print("Testing Update: ")
			err = client.Partitions().Update(ctx, firstPartitionName, &types.PartitionUpdate{
				MaxTime: int32Ptr(120),
			})
			if err != nil {
				fmt.Printf("❌ Failed: %v\n", err)
			} else {
				fmt.Printf("✅ Success\n")
			}
		}
	}
}

func testAccountsManager(client types.SlurmClient, _ string) {
	ctx := context.Background()
	fmt.Println("\n=== Accounts Manager ===")

	fmt.Println("Methods to implement in adapter:")
	fmt.Println("  - List(ctx) ([]*Account, error)")
	fmt.Println("  - Get(ctx, accountName) (*Account, error)")
	fmt.Println("  - Create(ctx, account) error")
	fmt.Println("  - Update(ctx, accountName, account) error")
	fmt.Println("  - Delete(ctx, accountName) error")

	// Test List
	fmt.Print("\nTesting List: ")
	accounts, err := client.Accounts().List(ctx, &types.ListAccountsOptions{})
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		accountCount := 0
		if accounts != nil && accounts.Accounts != nil {
			accountCount = len(accounts.Accounts)
		}
		fmt.Printf("✅ Success: Found %d accounts\n", accountCount)
	}

	// Test Create
	fmt.Print("Testing Create: ")
	testAccount := &types.AccountCreate{
		Name:        fmt.Sprintf("test-account-%d", time.Now().Unix()),
		Description: "Test account",
	}
	_, err = client.Accounts().Create(ctx, testAccount)
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		fmt.Printf("✅ Success: Created %s\n", testAccount.Name)

		// Test Get
		fmt.Print("Testing Get: ")
		account, err := client.Accounts().Get(ctx, testAccount.Name)
		if err != nil {
			fmt.Printf("❌ Failed: %v\n", err)
		} else {
			fmt.Printf("✅ Success: Got account %s\n", account.Name)
		}

		// Test Update
		fmt.Print("Testing Update: ")
		updateAccount := &types.AccountUpdate{
			Description: stringPtr("Updated test account"),
		}
		err = client.Accounts().Update(ctx, testAccount.Name, updateAccount)
		if err != nil {
			fmt.Printf("❌ Failed: %v\n", err)
		} else {
			fmt.Printf("✅ Success\n")
		}

		// Test Delete
		fmt.Print("Testing Delete: ")
		err = client.Accounts().Delete(ctx, testAccount.Name)
		if err != nil {
			fmt.Printf("❌ Failed: %v\n", err)
		} else {
			fmt.Printf("✅ Success\n")
		}
	}
}

func testUsersManager(client types.SlurmClient, _ string) {
	ctx := context.Background()
	fmt.Println("\n=== Users Manager ===")

	fmt.Println("Methods to implement in adapter:")
	fmt.Println("  - List(ctx, opts) (*UserList, error)")
	fmt.Println("  - Get(ctx, userName) (*User, error)")
	fmt.Println("  - GetUserAccounts(ctx, userName) ([]*UserAccount, error)")
	fmt.Println("  - GetUserQuotas(ctx, userName) (*UserQuota, error)")
	fmt.Println("  - GetUserDefaultAccount(ctx, userName) (*Account, error)")

	// Test List
	fmt.Print("\nTesting List: ")
	users, err := client.Users().List(ctx, &types.ListUsersOptions{})
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		userCount := 0
		var firstUserName string
		if users != nil && users.Users != nil {
			userCount = len(users.Users)
			if userCount > 0 {
				firstUserName = users.Users[0].Name
			}
		}
		fmt.Printf("✅ Success: Found %d users\n", userCount)

		if firstUserName != "" {
			// Test Get
			fmt.Print("Testing Get: ")
			user, err := client.Users().Get(ctx, firstUserName)
			if err != nil {
				fmt.Printf("❌ Failed: %v\n", err)
			} else {
				fmt.Printf("✅ Success: Got user %s\n", user.Name)
			}

			// Note: GetUserAccounts, GetUserQuotas, GetUserDefaultAccount
			// are not part of the core UserManager interface.
			// Use Associations().List() to get user-account relationships.
		}
	}
}

func testQoSManager(client types.SlurmClient, _ string) {
	ctx := context.Background()
	fmt.Println("\n=== QoS Manager ===")

	fmt.Println("Methods to implement in adapter:")
	fmt.Println("  - List(ctx) ([]*QoS, error)")
	fmt.Println("  - Get(ctx, qosName) (*QoS, error)")
	fmt.Println("  - Create(ctx, qos) error")
	fmt.Println("  - Update(ctx, qosName, qos) error")
	fmt.Println("  - Delete(ctx, qosName) error")

	// Test List
	fmt.Print("\nTesting List: ")
	qosList, err := client.QoS().List(ctx, &types.ListQoSOptions{})
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		qosCount := 0
		var firstQoSName string
		if qosList != nil && qosList.QoS != nil {
			qosCount = len(qosList.QoS)
			if qosCount > 0 && qosList.QoS[0].Name != nil {
				firstQoSName = *qosList.QoS[0].Name
			}
		}
		fmt.Printf("✅ Success: Found %d QoS\n", qosCount)

		if firstQoSName != "" {
			// Test Get
			fmt.Print("Testing Get: ")
			qos, err := client.QoS().Get(ctx, firstQoSName)
			if err != nil {
				fmt.Printf("❌ Failed: %v\n", err)
			} else {
				fmt.Printf("✅ Success: Got QoS %v\n", qos.Name)
			}
		}
	}

	// Test Create
	fmt.Print("Testing Create: ")
	testQoS := &types.QoSCreate{
		Name:        fmt.Sprintf("test-qos-%d", time.Now().Unix()),
		Description: "Test QoS",
		Priority:    10,
	}
	_, err = client.QoS().Create(ctx, testQoS)
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		fmt.Printf("✅ Success: Created %s\n", testQoS.Name)

		// Test Update
		fmt.Print("Testing Update: ")
		updateQoS := &types.QoSUpdate{
			Priority: intPtr(20),
		}
		err = client.QoS().Update(ctx, testQoS.Name, updateQoS)
		if err != nil {
			fmt.Printf("❌ Failed: %v\n", err)
		} else {
			fmt.Printf("✅ Success\n")
		}

		// Test Delete
		fmt.Print("Testing Delete: ")
		err = client.QoS().Delete(ctx, testQoS.Name)
		if err != nil {
			fmt.Printf("❌ Failed: %v\n", err)
		} else {
			fmt.Printf("✅ Success\n")
		}
	}
}

func testReservationsManager(client types.SlurmClient, _ string) {
	ctx := context.Background()
	fmt.Println("\n=== Reservations Manager ===")

	fmt.Println("Methods to implement in adapter:")
	fmt.Println("  - List(ctx, opts) (*ReservationList, error)")
	fmt.Println("  - Get(ctx, reservationName) (*Reservation, error)")
	fmt.Println("  - Create(ctx, reservation) error")
	fmt.Println("  - Update(ctx, reservationName, update) error")
	fmt.Println("  - Delete(ctx, reservationName) error")

	// Test List
	fmt.Print("\nTesting List: ")
	reservations, err := client.Reservations().List(ctx, &types.ListReservationsOptions{})
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		reservationCount := 0
		if reservations != nil && reservations.Reservations != nil {
			reservationCount = len(reservations.Reservations)
		}
		fmt.Printf("✅ Success: Found %d reservations\n", reservationCount)
	}

	// Test Create
	fmt.Print("Testing Create: ")
	startTime := time.Now().Add(1 * time.Hour)
	endTime := startTime.Add(1 * time.Hour)
	resName := fmt.Sprintf("test-res-%d", time.Now().Unix())
	testReservation := &types.ReservationCreate{
		Name:      stringPtr(resName),
		StartTime: startTime,
		EndTime:   endTime,
		NodeCount: uint32Ptr(1),
		Users:     []string{"root"},
	}
	_, err = client.Reservations().Create(ctx, testReservation)
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		fmt.Printf("✅ Success: Created %s\n", resName)

		// Test Get
		fmt.Print("Testing Get: ")
		reservation, err := client.Reservations().Get(ctx, resName)
		if err != nil {
			fmt.Printf("❌ Failed: %v\n", err)
		} else {
			fmt.Printf("✅ Success: Got reservation %v\n", reservation.Name)
		}

		// Test Update
		fmt.Print("Testing Update: ")
		err = client.Reservations().Update(ctx, resName, &types.ReservationUpdate{
			NodeCount: int32Ptr(2),
		})
		if err != nil {
			fmt.Printf("❌ Failed: %v\n", err)
		} else {
			fmt.Printf("✅ Success\n")
		}

		// Test Delete
		fmt.Print("Testing Delete: ")
		err = client.Reservations().Delete(ctx, resName)
		if err != nil {
			fmt.Printf("❌ Failed: %v\n", err)
		} else {
			fmt.Printf("✅ Success\n")
		}
	}
}

func testAssociationsManager(client types.SlurmClient, _ string) {
	ctx := context.Background()
	fmt.Println("\n=== Associations Manager ===")

	fmt.Println("Methods to implement in adapter:")
	fmt.Println("  - List(ctx, opts) (*AssociationList, error)")
	fmt.Println("  - Get(ctx, opts) (*Association, error)")
	fmt.Println("  - Create(ctx, associations) ([]*Association, error)")
	fmt.Println("  - Update(ctx, associations) error")
	fmt.Println("  - Delete(ctx, opts) error")

	// Test List
	fmt.Print("\nTesting List: ")
	associations, err := client.Associations().List(ctx, &types.ListAssociationsOptions{})
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		assocCount := 0
		if associations != nil && associations.Associations != nil {
			assocCount = len(associations.Associations)
		}
		fmt.Printf("✅ Success: Found %d associations\n", assocCount)

		if assocCount > 0 {
			// Test Get by ID
			fmt.Print("Testing Get: ")
			firstAssoc := associations.Associations[0]
			if firstAssoc.ID != nil {
				assocID := fmt.Sprintf("%d", *firstAssoc.ID)
				association, err := client.Associations().Get(ctx, assocID)
				if err != nil {
					fmt.Printf("❌ Failed: %v\n", err)
				} else {
					fmt.Printf("✅ Success: Got association for user %s\n", association.User)
				}
			} else {
				fmt.Printf("⚠️ Skipped: First association has no ID\n")
			}
		}
	}

	// Test Create
	fmt.Print("Testing Create: ")
	testAssociation := &types.AssociationCreate{
		Account:   "root",
		User:      "root",
		Partition: "normal",
	}
	_, err = client.Associations().Create(ctx, []*types.AssociationCreate{testAssociation})
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		fmt.Printf("✅ Success\n")
	}
}

// Helper functions
func intPtr(i int) *int {
	return &i
}

func int32Ptr(i int32) *int32 {
	return &i
}

func uint32Ptr(i uint32) *uint32 {
	return &i
}

func stringPtr(s string) *string {
	return &s
}
