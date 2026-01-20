// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/interfaces"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/pkg/config"
)

// Example: Account management (v0.0.43+ only)
func main() {
	// Create configuration
	cfg := config.NewDefault()
	cfg.BaseURL = "https://cluster.example.com:6820"

	// Create authentication
	authProvider := auth.NewTokenAuth("your-jwt-token")

	ctx := context.Background()

	// Example 1: Check account support
	fmt.Println("=== Account Support Check ===")
	checkAccountSupport(ctx, cfg, authProvider)

	// Example 2: List accounts
	fmt.Println("\n=== List Accounts ===")
	listAccounts(ctx, cfg, authProvider)

	// Example 3: Create account hierarchy
	fmt.Println("\n=== Create Account Hierarchy ===")
	createAccountHierarchy(ctx, cfg, authProvider)

	// Example 4: Update account limits
	fmt.Println("\n=== Update Account Limits ===")
	updateAccountLimits(ctx, cfg, authProvider)

	// Example 5: Account associations
	fmt.Println("\n=== Account Associations ===")
	demonstrateAccountAssociations(ctx, cfg, authProvider)

	// Example 6: Resource allocation
	fmt.Println("\n=== Resource Allocation ===")
	demonstrateResourceAllocation(ctx, cfg, authProvider)
}

// checkAccountSupport checks if the cluster supports account management
func checkAccountSupport(ctx context.Context, cfg *config.Config, auth auth.Provider) {
	// Try different versions
	versions := []string{"v0.0.40", "v0.0.41", "v0.0.42", "v0.0.43"}

	for _, version := range versions {
		client, err := slurm.NewClientWithVersion(ctx, version,
			slurm.WithConfig(cfg),
			slurm.WithAuth(auth),
		)
		if err != nil {
			log.Printf("Failed to create %s client: %v", version, err)
			continue
		}
		defer client.Close()

		// Check if Accounts is supported
		if client.Accounts() == nil {
			fmt.Printf("%s: Accounts NOT supported\n", version)
		} else {
			fmt.Printf("%s: Accounts supported ✓\n", version)
		}
	}

	fmt.Println("\nNote: Account management requires SLURM REST API v0.0.43 or later")
}

// listAccounts demonstrates listing accounts
func listAccounts(ctx context.Context, cfg *config.Config, auth auth.Provider) {
	// Create v0.0.43 client
	client, err := slurm.NewClientWithVersion(ctx, "v0.0.43",
		slurm.WithConfig(cfg),
		slurm.WithAuth(auth),
	)
	if err != nil {
		log.Printf("Failed to create v0.0.43 client: %v", err)
		return
	}
	defer client.Close()

	if client.Accounts() == nil {
		fmt.Println("Accounts not supported")
		return
	}

	// List all accounts
	fmt.Println("Listing all accounts:")
	accountList, err := client.Accounts().List(ctx, nil)
	if err != nil {
		log.Printf("Failed to list accounts: %v", err)
		return
	}

	if len(accountList.Accounts) == 0 {
		fmt.Println("No accounts found")
		return
	}

	// Display account hierarchy
	fmt.Println("\nAccount Hierarchy:")
	displayAccountHierarchy(accountList.Accounts, "", "")

	// List accounts for specific organizations
	fmt.Println("\nListing accounts for specific organizations:")
	orgAccounts, err := client.Accounts().List(ctx, &interfaces.ListAccountsOptions{
		Organizations: []string{"Engineering", "Research"},
	})
	if err != nil {
		log.Printf("Failed to list organization accounts: %v", err)
		return
	}

	fmt.Printf("Found %d accounts for specified organizations\n",
		len(orgAccounts.Accounts))

	// List accounts with associations
	fmt.Println("\nListing accounts with associations:")
	assocAccounts, err := client.Accounts().List(ctx, &interfaces.ListAccountsOptions{
		WithAssociations: true,
		WithCoordinators: true,
	})
	if err != nil {
		log.Printf("Failed to list accounts with associations: %v", err)
		return
	}

	for _, account := range assocAccounts.Accounts {
		if len(account.Users) > 0 || len(account.CoordinatorUsers) > 0 {
			fmt.Printf("\nAccount: %s\n", account.Name)
			if len(account.CoordinatorUsers) > 0 {
				fmt.Printf("  Coordinators: %v\n", account.CoordinatorUsers)
			}
			if len(account.Users) > 0 {
				fmt.Printf("  Users: %v\n", account.Users)
			}
		}
	}
}

// displayAccountHierarchy recursively displays account hierarchy
func displayAccountHierarchy(accounts []interfaces.Account, parent, indent string) {
	for _, account := range accounts {
		if account.ParentAccount == parent {
			fmt.Printf("%s%s", indent, account.Name)
			if account.Description != "" {
				fmt.Printf(" - %s", account.Description)
			}
			fmt.Println()

			// Display account details
			if account.Organization != "" {
				fmt.Printf("%s  Organization: %s\n", indent, account.Organization)
			}
			if account.MaxJobs > 0 {
				fmt.Printf("%s  Max Jobs: %d\n", indent, account.MaxJobs)
			}
			if len(account.AllowedPartitions) > 0 {
				fmt.Printf("%s  Partitions: %v\n", indent, account.AllowedPartitions)
			}

			// Recurse for child accounts
			displayAccountHierarchy(accounts, account.Name, indent+"  ")
		}
	}
}

// createAccountHierarchy demonstrates creating an account hierarchy
func createAccountHierarchy(ctx context.Context, cfg *config.Config, auth auth.Provider) {
	client, err := slurm.NewClientWithVersion(ctx, "v0.0.43",
		slurm.WithConfig(cfg),
		slurm.WithAuth(auth),
	)
	if err != nil {
		log.Printf("Failed to create v0.0.43 client: %v", err)
		return
	}
	defer client.Close()

	if client.Accounts() == nil {
		fmt.Println("Accounts not supported")
		return
	}

	// Example 1: Create root organization account
	fmt.Println("Creating root organization account:")

	rootAccount := &interfaces.AccountCreate{
		Name:             "acme-corp",
		Description:      "ACME Corporation - Root Account",
		Organization:     "ACME Corporation",
		CoordinatorUsers: []string{"admin", "finance-admin"},
		MaxJobs:          1000,
		MaxNodes:         100,
		SharesPriority:   1000,
		Flags:            []string{"AllowSubmit"},
	}

	resp, err := client.Accounts().Create(ctx, rootAccount)
	if err != nil {
		log.Printf("Failed to create root account: %v", err)
	} else {
		fmt.Printf("Created account: %s\n", resp.AccountName)
	}

	// Example 2: Create department accounts
	fmt.Println("\nCreating department accounts:")

	departments := []struct {
		name        string
		description string
		maxJobs     int
		partitions  []string
	}{
		{
			name:        "engineering",
			description: "Engineering Department",
			maxJobs:     500,
			partitions:  []string{"compute", "gpu", "highmem"},
		},
		{
			name:        "research",
			description: "Research Department",
			maxJobs:     300,
			partitions:  []string{"compute", "gpu"},
		},
		{
			name:        "analytics",
			description: "Analytics Department",
			maxJobs:     200,
			partitions:  []string{"compute", "highmem"},
		},
	}

	for _, dept := range departments {
		deptAccount := &interfaces.AccountCreate{
			Name:              dept.name,
			Description:       dept.description,
			Organization:      "ACME Corporation",
			ParentAccount:     "acme-corp",
			AllowedPartitions: dept.partitions,
			MaxJobs:           dept.maxJobs,
			MaxJobsPerUser:    50,
			DefaultPartition:  dept.partitions[0],
			SharesPriority:    100,
			Flags:             []string{"AllowSubmit"},
		}

		resp2, err := client.Accounts().Create(ctx, deptAccount)
		if err != nil {
			log.Printf("Failed to create %s account: %v", dept.name, err)
		} else {
			fmt.Printf("  Created department: %s\n", resp2.AccountName)
		}
	}

	// Example 3: Create project accounts under departments
	fmt.Println("\nCreating project accounts:")

	projects := []struct {
		name        string
		parent      string
		description string
		qos         []string
	}{
		{
			name:        "ml-research",
			parent:      "research",
			description: "Machine Learning Research Project",
			qos:         []string{"normal", "high-priority"},
		},
		{
			name:        "platform-dev",
			parent:      "engineering",
			description: "Platform Development Team",
			qos:         []string{"normal", "urgent"},
		},
		{
			name:        "data-pipeline",
			parent:      "analytics",
			description: "Data Pipeline Team",
			qos:         []string{"normal", "scavenger"},
		},
	}

	for _, proj := range projects {
		projAccount := &interfaces.AccountCreate{
			Name:           proj.name,
			Description:    proj.description,
			ParentAccount:  proj.parent,
			AllowedQoS:     proj.qos,
			DefaultQoS:     proj.qos[0],
			MaxJobs:        100,
			MaxJobsPerUser: 20,
			MaxNodes:       25,
			// Set resource limits using TRES
			MaxTRES: map[string]int{
				"cpu": 1000,
				"mem": 4096000, // MB
				"gpu": 10,
			},
			GrpTRES: map[string]int{
				"cpu": 500,
				"mem": 2048000, // MB
				"gpu": 5,
			},
			SharesPriority: 50,
			Flags:          []string{"AllowSubmit"},
		}

		resp3, err := client.Accounts().Create(ctx, projAccount)
		if err != nil {
			log.Printf("Failed to create %s project: %v", proj.name, err)
		} else {
			fmt.Printf("  Created project: %s under %s\n", resp3.AccountName, proj.parent)
		}
	}
}

// updateAccountLimits demonstrates updating account resource limits
func updateAccountLimits(ctx context.Context, cfg *config.Config, auth auth.Provider) {
	client, err := slurm.NewClientWithVersion(ctx, "v0.0.43",
		slurm.WithConfig(cfg),
		slurm.WithAuth(auth),
	)
	if err != nil {
		log.Printf("Failed to create v0.0.43 client: %v", err)
		return
	}
	defer client.Close()

	if client.Accounts() == nil {
		fmt.Println("Accounts not supported")
		return
	}

	accountName := "ml-research"

	// Get current account
	fmt.Printf("Getting account %s:\n", accountName)
	current, err := client.Accounts().Get(ctx, accountName)
	if err != nil {
		log.Printf("Failed to get account: %v", err)
		return
	}

	fmt.Printf("Current settings:\n")
	fmt.Printf("  Max Jobs: %d\n", current.MaxJobs)
	fmt.Printf("  Max Nodes: %d\n", current.MaxNodes)
	if current.MaxTRES != nil {
		fmt.Printf("  Max TRES: %v\n", current.MaxTRES)
	}

	// Update account - increase limits for a big project
	newMaxJobs := 200
	newMaxNodes := 50
	update := &interfaces.AccountUpdate{
		MaxJobs:     &newMaxJobs,
		MaxNodes:    &newMaxNodes,
		Description: stringPtr("Machine Learning Research Project - Expanded Resources"),
		MaxTRES: map[string]int{
			"cpu": 2000,
			"mem": 8192000, // 8TB
			"gpu": 20,
		},
		GrpTRES: map[string]int{
			"cpu": 1000,
			"mem": 4096000, // 4TB
			"gpu": 10,
		},
		// Update allowed partitions
		AllowedPartitions: []string{"compute", "gpu", "highmem", "gpu-large"},
		// Add more QoS options
		AllowedQoS: []string{"normal", "high-priority", "urgent"},
	}

	fmt.Println("\nUpdating account limits...")
	err = client.Accounts().Update(ctx, accountName, update)
	if err != nil {
		log.Printf("Failed to update account: %v", err)
		return
	}

	fmt.Println("Account updated successfully")

	// Verify update
	updated, err := client.Accounts().Get(ctx, accountName)
	if err != nil {
		log.Printf("Failed to get updated account: %v", err)
		return
	}

	fmt.Printf("New settings:\n")
	fmt.Printf("  Max Jobs: %d\n", updated.MaxJobs)
	fmt.Printf("  Max Nodes: %d\n", updated.MaxNodes)
	if updated.MaxTRES != nil {
		fmt.Printf("  Max TRES: %v\n", updated.MaxTRES)
	}
	fmt.Printf("  Allowed Partitions: %v\n", updated.AllowedPartitions)
}

// demonstrateAccountAssociations shows account-user associations
func demonstrateAccountAssociations(ctx context.Context, cfg *config.Config, auth auth.Provider) {
	client, err := slurm.NewClientWithVersion(ctx, "v0.0.43",
		slurm.WithConfig(cfg),
		slurm.WithAuth(auth),
	)
	if err != nil {
		log.Printf("Failed to create v0.0.43 client: %v", err)
		return
	}
	defer client.Close()

	if client.Accounts() == nil {
		fmt.Println("Accounts not supported")
		return
	}

	// Example 1: Add coordinators to an account
	fmt.Println("Adding coordinators to engineering account:")

	engUpdate := &interfaces.AccountUpdate{
		CoordinatorUsers: []string{"eng-lead", "eng-manager", "tech-lead"},
		Description:      stringPtr("Engineering Department - Updated Coordinators"),
	}

	err = client.Accounts().Update(ctx, "engineering", engUpdate)
	if err != nil {
		log.Printf("Failed to update coordinators: %v", err)
	} else {
		fmt.Println("Coordinators updated successfully")
	}

	// Example 2: Update default settings for users
	fmt.Println("\nUpdating default settings for platform-dev:")

	defaultCPULimit := 100
	defaultWallTime := 86400 // 24 hours
	platUpdate := &interfaces.AccountUpdate{
		DefaultPartition: stringPtr("compute"),
		DefaultQoS:       stringPtr("normal"),
		CPULimit:         &defaultCPULimit,
		MaxWallTime:      &defaultWallTime,
		// Set fair share allocation
		FairShareTRES: map[string]int{
			"cpu": 1000,
			"mem": 2048000, // 2TB
		},
		SharesPriority: intPtr(75),
	}

	err = client.Accounts().Update(ctx, "platform-dev", platUpdate)
	if err != nil {
		log.Printf("Failed to update platform-dev: %v", err)
	} else {
		fmt.Println("Default settings updated")
	}

	// Example 3: Create a special account for training/education
	fmt.Println("\nCreating training account with restrictions:")

	trainingAccount := &interfaces.AccountCreate{
		Name:              "training",
		Description:       "Training and Education Account",
		Organization:      "ACME Corporation",
		ParentAccount:     "acme-corp",
		CoordinatorUsers:  []string{"trainer1", "trainer2"},
		AllowedPartitions: []string{"training", "compute"},
		DefaultPartition:  "training",
		AllowedQoS:        []string{"training-qos"},
		DefaultQoS:        "training-qos",
		MaxJobs:           10,
		MaxJobsPerUser:    2,
		MaxNodes:          2,
		MaxWallTime:       3600, // 1 hour max
		CPULimit:          4,
		// Strict resource limits for training
		MaxTRES: map[string]int{
			"cpu": 16,
			"mem": 32768, // 32GB
		},
		MaxTRESPerUser: map[string]int{
			"cpu": 4,
			"mem": 8192, // 8GB
		},
		Flags: []string{"AllowSubmit", "RequireAssoc"},
	}

	resp, err := client.Accounts().Create(ctx, trainingAccount)
	if err != nil {
		log.Printf("Failed to create training account: %v", err)
	} else {
		fmt.Printf("Created training account: %s\n", resp.AccountName)
		fmt.Println("  - Limited to 2 jobs per user")
		fmt.Println("  - Maximum 1 hour runtime")
		fmt.Println("  - Restricted to training partition")
	}
}

// demonstrateResourceAllocation shows advanced resource allocation
func demonstrateResourceAllocation(ctx context.Context, cfg *config.Config, auth auth.Provider) {
	client, err := slurm.NewClientWithVersion(ctx, "v0.0.43",
		slurm.WithConfig(cfg),
		slurm.WithAuth(auth),
	)
	if err != nil {
		log.Printf("Failed to create v0.0.43 client: %v", err)
		return
	}
	defer client.Close()

	if client.Accounts() == nil {
		fmt.Println("Accounts not supported")
		return
	}

	// Example 1: Set up tiered resource allocation
	fmt.Println("Setting up tiered resource allocation:")

	tiers := []struct {
		account    string
		shares     int
		grpTRES    map[string]int
		grpMinutes map[string]int
	}{
		{
			account: "research",
			shares:  200,
			grpTRES: map[string]int{
				"cpu":  2000,
				"mem":  8192000, // 8TB
				"gpu":  20,
				"node": 50,
			},
			grpMinutes: map[string]int{
				"cpu": 120000000, // CPU-minutes per month
			},
		},
		{
			account: "engineering",
			shares:  150,
			grpTRES: map[string]int{
				"cpu":  1500,
				"mem":  6144000, // 6TB
				"gpu":  10,
				"node": 40,
			},
			grpMinutes: map[string]int{
				"cpu": 90000000, // CPU-minutes per month
			},
		},
		{
			account: "analytics",
			shares:  100,
			grpTRES: map[string]int{
				"cpu":  1000,
				"mem":  4096000, // 4TB
				"gpu":  5,
				"node": 30,
			},
			grpMinutes: map[string]int{
				"cpu": 60000000, // CPU-minutes per month
			},
		},
	}

	for _, tier := range tiers {
		update := &interfaces.AccountUpdate{
			SharesPriority: &tier.shares,
			GrpTRES:        tier.grpTRES,
			GrpTRESMinutes: tier.grpMinutes,
			Description:    stringPtr(fmt.Sprintf("Updated with tiered allocation - %d shares", tier.shares)),
		}

		err := client.Accounts().Update(ctx, tier.account, update)
		if err != nil {
			log.Printf("Failed to update %s: %v", tier.account, err)
		} else {
			fmt.Printf("  Updated %s: %d shares, TRES limits set\n", tier.account, tier.shares)
		}
	}

	// Example 2: Create a burst account for temporary high resource usage
	fmt.Println("\nCreating burst account for temporary projects:")

	burstAccount := &interfaces.AccountCreate{
		Name:              "burst-projects",
		Description:       "Temporary burst capacity for urgent projects",
		Organization:      "ACME Corporation",
		ParentAccount:     "acme-corp",
		CoordinatorUsers:  []string{"ops-lead", "cto"},
		AllowedPartitions: []string{"compute", "gpu", "highmem"},
		AllowedQoS:        []string{"urgent", "executive"},
		MaxJobs:           50,
		MaxJobsPerUser:    10,
		MaxNodes:          100,
		MaxWallTime:       14400, // 4 hours max
		// High resource limits but with usage tracking
		MaxTRES: map[string]int{
			"cpu": 5000,
			"mem": 20480000, // 20TB
			"gpu": 50,
		},
		// Group limits to prevent monopolization
		GrpTRES: map[string]int{
			"cpu": 2500,
			"mem": 10240000, // 10TB
			"gpu": 25,
		},
		// Time-based limits (monthly quota)
		GrpTRESMinutes: map[string]int{
			"cpu": 30000000, // Limited monthly CPU-minutes
			"gpu": 1000000,  // Limited monthly GPU-minutes
		},
		SharesPriority: 500, // High priority for burst
		Flags:          []string{"AllowSubmit", "RequireAssoc"},
	}

	resp, err := client.Accounts().Create(ctx, burstAccount)
	if err != nil {
		log.Printf("Failed to create burst account: %v", err)
	} else {
		fmt.Printf("Created burst account: %s\n", resp.AccountName)
		fmt.Println("  - High resource limits for urgent work")
		fmt.Println("  - Monthly quotas to prevent overuse")
		fmt.Println("  - Requires association to track usage")
	}

	// Example 3: Resource allocation summary
	fmt.Println("\nResource Allocation Summary:")
	fmt.Println("  Research:    200 shares, 2000 CPUs, 20 GPUs")
	fmt.Println("  Engineering: 150 shares, 1500 CPUs, 10 GPUs")
	fmt.Println("  Analytics:   100 shares, 1000 CPUs, 5 GPUs")
	fmt.Println("  Burst:       500 shares, 5000 CPUs (quota limited)")
	fmt.Println("\nFair share will distribute resources based on:")
	fmt.Println("  - Account share priority")
	fmt.Println("  - Current usage vs allocation")
	fmt.Println("  - Group TRES limits")
	fmt.Println("  - Time-based quotas")
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
