// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/interfaces"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/pkg/errors"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	baseURL    = flag.String("url", "http://localhost:6820", "SLURM REST API URL")
	apiVersion = flag.String("version", "v0.0.43", "API version to use")
	username   = flag.String("user", "", "Username to analyze (required)")
	account    = flag.String("account", "", "Account to analyze")
	hierarchy  = flag.Bool("hierarchy", false, "Show account hierarchy")
	quotas     = flag.Bool("quotas", false, "Show quota information")
	fairshare  = flag.Bool("fairshare", false, "Show fair-share information")
	validate   = flag.Bool("validate", false, "Validate user-account access")
	bulk       = flag.Bool("bulk", false, "Demonstrate bulk operations")
)

func main() {
	flag.Parse()

	if *username == "" && !*hierarchy && !*bulk {
		fmt.Fprintf(os.Stderr, "Error: -user is required unless using -hierarchy or -bulk\n")
		flag.Usage()
		os.Exit(1)
	}

	// Create client
	ctx := context.Background()
	client, err := slurm.NewClientWithVersion(ctx, *apiVersion,
		slurm.WithBaseURL(*baseURL),
		slurm.WithAuth(auth.NewTokenAuth(os.Getenv("SLURM_JWT_TOKEN"))),
	)
	if err != nil {
		log.Fatal("Failed to create client:", err)
	}
	defer client.Close()

	fmt.Printf("Connected to SLURM REST API %s at %s\n\n", client.Version(), *baseURL)

	// Demonstrate various user-account management features
	if *hierarchy {
		demonstrateHierarchyNavigation(ctx, client)
	}

	if *username != "" {
		if *account != "" && *validate {
			validateUserAccountAccess(ctx, client, *username, *account)
		}

		if *quotas {
			demonstrateQuotaMonitoring(ctx, client, *username, *account)
		}

		if *fairshare {
			demonstrateFairShareAnalysis(ctx, client, *username, *account)
		}

		if !*validate && !*quotas && !*fairshare {
			// Show basic user information
			showUserAccountAssociations(ctx, client, *username)
		}
	}

	if *bulk {
		demonstrateBulkOperations(ctx, client)
	}
}

func demonstrateHierarchyNavigation(ctx context.Context, client interfaces.SlurmClient) {
	fmt.Println("=== Account Hierarchy Navigation ===")

	accountManager := client.Accounts()

	// Get account hierarchy from root
	rootAccount := "root"
	if *account != "" {
		rootAccount = *account
	}

	fmt.Printf("\nRetrieving hierarchy from root account: %s\n", rootAccount)
	hierarchy, err := accountManager.GetAccountHierarchy(ctx, rootAccount)
	if err != nil {
		if errors.IsNotImplementedError(err) {
			fmt.Println("Note: Account hierarchy retrieval not implemented yet")
		} else {
			fmt.Printf("Error retrieving hierarchy: %v\n", err)
		}
		return
	}

	if hierarchy != nil {
		fmt.Printf("\nAccount Hierarchy Structure:\n")
		if hierarchy.Account != nil {
			fmt.Printf("  Account: %s\n", hierarchy.Account.Name)
		}
		fmt.Printf("  Level: %d\n", hierarchy.Level)
		fmt.Printf("  Total Sub-Accounts: %d\n", hierarchy.TotalSubAccounts)
		fmt.Printf("  Total Users: %d\n", hierarchy.TotalUsers)

		// Display tree structure
		if len(hierarchy.ChildAccounts) > 0 {
			fmt.Printf("\nHierarchy Tree:\n")
			printAccountHierarchy(hierarchy, 0)
		}
	}

	// Demonstrate parent/child navigation
	if *account != "" {
		fmt.Printf("\n\nNavigating from account: %s\n", *account)

		// Get parent accounts
		parents, err := accountManager.GetParentAccounts(ctx, *account)
		if err == nil && len(parents) > 0 {
			fmt.Printf("\nParent Accounts (up to root):\n")
			for i, parent := range parents {
				fmt.Printf("  %d. %s\n", i+1, parent.Name)
			}
		}

		// Get child accounts
		children, err := accountManager.GetChildAccounts(ctx, *account, 2) // 2 levels deep
		if err == nil && len(children) > 0 {
			fmt.Printf("\nChild Accounts (2 levels deep):\n")
			for i, child := range children {
				fmt.Printf("  %d. %s\n", i+1, child.Name)
			}
		}
	}
}

func demonstrateQuotaMonitoring(ctx context.Context, client interfaces.SlurmClient, userName, accountName string) {
	fmt.Println("\n=== Quota Monitoring ===")

	userManager := client.Users()
	accountManager := client.Accounts()

	// Get user quotas
	fmt.Printf("\nUser Quotas for %s:\n", userName)
	userQuotas, err := userManager.GetUserQuotas(ctx, userName)
	if err != nil {
		if errors.IsNotImplementedError(err) {
			fmt.Println("Note: User quota retrieval not implemented yet")
		} else {
			fmt.Printf("Error retrieving user quotas: %v\n", err)
		}
	} else if userQuotas != nil {
		displayUserQuotas(userQuotas)
	}

	// Get account quotas if specified
	if accountName != "" {
		fmt.Printf("\nAccount Quotas for %s:\n", accountName)
		accountQuotas, err := accountManager.GetAccountQuotas(ctx, accountName)
		if err != nil {
			if errors.IsNotImplementedError(err) {
				fmt.Println("Note: Account quota retrieval not implemented yet")
			} else {
				fmt.Printf("Error retrieving account quotas: %v\n", err)
			}
		} else if accountQuotas != nil {
			displayAccountQuotas(accountQuotas)
		}

		// Get account usage over different timeframes
		fmt.Printf("\nAccount Usage for %s:\n", accountName)
		timeframes := []string{"daily", "weekly", "monthly"}
		for _, timeframe := range timeframes {
			usage, err := accountManager.GetAccountQuotaUsage(ctx, accountName, timeframe)
			if err == nil && usage != nil {
				fmt.Printf("\n%s Usage:\n", cases.Title(language.English).String(timeframe))
				displayAccountUsage(usage)
			}
		}
	}
}

func demonstrateFairShareAnalysis(ctx context.Context, client interfaces.SlurmClient, userName, accountName string) {
	fmt.Println("\n=== Fair-Share Analysis ===")

	userManager := client.Users()
	accountManager := client.Accounts()

	// Get user fair-share information
	fmt.Printf("\nUser Fair-Share for %s:\n", userName)
	userFairShare, err := userManager.GetUserFairShare(ctx, userName)
	if err != nil {
		if errors.IsNotImplementedError(err) {
			fmt.Println("Note: User fair-share retrieval not implemented yet")
		} else {
			fmt.Printf("Error retrieving user fair-share: %v\n", err)
		}
	} else if userFairShare != nil {
		displayUserFairShare(userFairShare)
	}

	// Get account fair-share if specified
	if accountName != "" {
		fmt.Printf("\nAccount Fair-Share for %s:\n", accountName)
		accountFairShare, err := accountManager.GetAccountFairShare(ctx, accountName)
		if err != nil {
			if errors.IsNotImplementedError(err) {
				fmt.Println("Note: Account fair-share retrieval not implemented yet")
			} else {
				fmt.Printf("Error retrieving account fair-share: %v\n", err)
			}
		} else if accountFairShare != nil {
			displayAccountFairShare(accountFairShare)
		}
	}

	// Demonstrate job priority calculation
	fmt.Printf("\n\nJob Priority Calculation for %s:\n", userName)
	jobSubmission := &interfaces.JobSubmission{
		Script: "#!/bin/bash\necho 'Test job for priority calculation'",
		// Account field doesn't exist in JobSubmission
		Partition: "compute",
		CPUs:      4,
		Memory:    8192,
		TimeLimit: 120, // 2 hours in minutes
	}

	priority, err := userManager.CalculateJobPriority(ctx, userName, jobSubmission)
	if err != nil {
		if errors.IsNotImplementedError(err) {
			fmt.Println("Note: Job priority calculation not implemented yet")
		} else {
			fmt.Printf("Error calculating job priority: %v\n", err)
		}
	} else if priority != nil {
		displayJobPriority(priority)
	}

	// Get fair-share hierarchy
	if accountName != "" {
		fmt.Printf("\n\nFair-Share Hierarchy from %s:\n", accountName)
		hierarchy, err := accountManager.GetFairShareHierarchy(ctx, accountName)
		if err == nil && hierarchy != nil {
			displayFairShareHierarchy(hierarchy)
		}
	}
}

func validateUserAccountAccess(ctx context.Context, client interfaces.SlurmClient, userName, accountName string) {
	fmt.Printf("\n=== Validating Access: %s → %s ===\n", userName, accountName)

	userManager := client.Users()
	accountManager := client.Accounts()

	// Validate from user perspective
	fmt.Printf("\nValidating from user perspective:\n")
	userValidation, err := userManager.ValidateUserAccountAccess(ctx, userName, accountName)
	if err != nil {
		if errors.IsNotImplementedError(err) {
			fmt.Println("Note: User-account validation not implemented yet")
		} else {
			fmt.Printf("Error validating access: %v\n", err)
		}
	} else if userValidation != nil {
		// IsAllowed field doesn't exist in UserAccessValidation
		// fmt.Printf("  Access Allowed: %v\n", userValidation.IsAllowed)
		fmt.Printf("  Reason: %s\n", userValidation.Reason)
		if len(userValidation.Permissions) > 0 {
			fmt.Printf("  Permissions: %v\n", userValidation.Permissions)
		}
	}

	// Validate from account perspective
	fmt.Printf("\nValidating from account perspective:\n")
	accountValidation, err := accountManager.ValidateUserAccess(ctx, userName, accountName)
	if err != nil {
		if errors.IsNotImplementedError(err) {
			fmt.Println("Note: Account-user validation not implemented yet")
		} else {
			fmt.Printf("Error validating access: %v\n", err)
		}
	} else if accountValidation != nil {
		// IsAllowed field doesn't exist in UserAccessValidation
		// fmt.Printf("  Access Allowed: %v\n", accountValidation.IsAllowed)
		fmt.Printf("  Has Access: %v\n", accountValidation.HasAccess)
		fmt.Printf("  Access Level: %s\n", accountValidation.AccessLevel)
	}
}

func showUserAccountAssociations(ctx context.Context, client interfaces.SlurmClient, userName string) {
	fmt.Printf("\n=== User Account Associations for %s ===\n", userName)

	userManager := client.Users()

	// Get all user accounts
	accounts, err := userManager.GetUserAccounts(ctx, userName)
	if err != nil {
		if errors.IsNotImplementedError(err) {
			fmt.Println("Note: User account retrieval not implemented yet")
		} else {
			fmt.Printf("Error retrieving user accounts: %v\n", err)
		}
		return
	}

	if len(accounts) == 0 {
		fmt.Println("No accounts found for user")
		return
	}

	// Display accounts in a table
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Account\tRole\tDefault\tPartitions\tQoS")
	fmt.Fprintln(w, "-------\t----\t-------\t----------\t---")

	for _, acc := range accounts {
		fmt.Fprintf(w, "%s\t%s\t%v\t%s\t%s\n",
			acc.AccountName,
			"User", // Role field doesn't exist
			acc.IsDefault,
			acc.Partition, // Partitions is a single string
			acc.QoS)       // QoS is a single string)
	}
	w.Flush()

	// Get default account
	defaultAccount, err := userManager.GetUserDefaultAccount(ctx, userName)
	if err == nil && defaultAccount != nil {
		fmt.Printf("\nDefault Account: %s\n", defaultAccount.Name)
	}

	// Get detailed associations with filtering
	fmt.Printf("\n\nDetailed Account Associations:\n")
	opts := &interfaces.ListUserAccountAssociationsOptions{
		ActiveOnly: true,
	}

	associations, err := userManager.GetUserAccountAssociations(ctx, userName, opts)
	if err == nil && len(associations) > 0 {
		for i, assoc := range associations {
			fmt.Printf("\n%d. %s/%s:\n", i+1, assoc.AccountName, assoc.Partition)
			fmt.Printf("   Role: %s\n", assoc.Role)
			fmt.Printf("   Permissions: %v\n", assoc.Permissions)
			fmt.Printf("   Max Jobs: %d\n", assoc.MaxJobs)
			fmt.Printf("   Priority: %d\n", assoc.Priority)
			// CreatedAt field doesn't exist in UserAccountAssociation
		}
	}
}

func demonstrateBulkOperations(ctx context.Context, client interfaces.SlurmClient) {
	fmt.Println("\n=== Bulk Operations Demo ===")

	userManager := client.Users()

	// Demonstrate bulk user accounts retrieval
	userNames := []string{"user1", "user2", "user3", "user4", "user5"}
	fmt.Printf("\nRetrieving accounts for %d users in bulk...\n", len(userNames))

	bulkAccounts, err := userManager.GetBulkUserAccounts(ctx, userNames)
	if err != nil {
		if errors.IsNotImplementedError(err) {
			fmt.Println("Note: Bulk user accounts retrieval not implemented yet")
		} else {
			fmt.Printf("Error retrieving bulk accounts: %v\n", err)
		}
	} else {
		for userName, accounts := range bulkAccounts {
			fmt.Printf("\n%s has %d accounts:\n", userName, len(accounts))
			for _, acc := range accounts {
				fmt.Printf("  - %s\n", acc.AccountName)
			}
		}
	}

	// Demonstrate bulk account users retrieval
	accountNames := []string{"research", "engineering", "finance"}
	fmt.Printf("\n\nRetrieving users for %d accounts in bulk...\n", len(accountNames))

	bulkUsers, err := userManager.GetBulkAccountUsers(ctx, accountNames)
	if err != nil {
		if errors.IsNotImplementedError(err) {
			fmt.Println("Note: Bulk account users retrieval not implemented yet")
		} else {
			fmt.Printf("Error retrieving bulk users: %v\n", err)
		}
	} else {
		for accountName, users := range bulkUsers {
			fmt.Printf("\n%s has %d users:\n", accountName, len(users))
			for _, user := range users {
				fmt.Printf("  - %s (role: %s, permissions: %v)\n",
					user.UserName, user.Role, user.Permissions)
			}
		}
	}
}

// Helper functions for displaying data

func printAccountHierarchy(hierarchy *interfaces.AccountHierarchy, depth int) {
	indent := strings.Repeat("  ", depth)
	if hierarchy.Account != nil {
		fmt.Printf("%s├─ %s (users: %d)\n", indent, hierarchy.Account.Name, hierarchy.TotalUsers)
	}

	for _, child := range hierarchy.ChildAccounts {
		printAccountHierarchy(child, depth+1)
	}
}

func displayUserQuotas(quotas *interfaces.UserQuota) {
	fmt.Printf("  Default Account: %s\n", quotas.DefaultAccount)
	fmt.Printf("  Max Jobs: %d\n", quotas.MaxJobs)
	fmt.Printf("  Max Submit Jobs: %d\n", quotas.MaxSubmitJobs)
	fmt.Printf("  Max CPUs: %d\n", quotas.MaxCPUs)
	fmt.Printf("  Max Memory: %d MB\n", quotas.MaxMemory)
	fmt.Printf("  Max Wall Time: %d minutes\n", quotas.MaxWallTime)

	if len(quotas.TRESLimits) > 0 {
		fmt.Printf("  TRES Limits:\n")
		for tres, limit := range quotas.TRESLimits {
			fmt.Printf("    %s: %d\n", tres, limit)
		}
	}
}

func displayAccountQuotas(quotas *interfaces.AccountQuota) {
	// AccountName and Description fields don't exist in AccountQuota
	fmt.Printf("  CPU Limit: %d\n", quotas.CPULimit)
	fmt.Printf("  Max Jobs: %d\n", quotas.MaxJobs)
	fmt.Printf("  Jobs Used: %d\n", quotas.JobsUsed)
	fmt.Printf("  Max Wall Time: %d minutes\n", quotas.MaxWallTime)

	if len(quotas.MaxTRES) > 0 {
		fmt.Printf("  Max TRES:\n")
		for tres, limit := range quotas.MaxTRES {
			fmt.Printf("    %s: %d\n", tres, limit)
		}
	}
}

func displayAccountUsage(usage *interfaces.AccountUsage) {
	fmt.Printf("  CPU Hours: %.2f\n", usage.CPUHours)
	fmt.Printf("  Jobs Completed: %d\n", usage.JobsCompleted)
	fmt.Printf("  Jobs Failed: %d\n", usage.JobsFailed)
	fmt.Printf("  Total Jobs: %d\n", usage.JobCount)
	fmt.Printf("  Period: %s to %s\n",
		usage.StartTime.Format("2006-01-02"),
		usage.EndTime.Format("2006-01-02"))
}

func displayUserFairShare(fairShare *interfaces.UserFairShare) {
	fmt.Printf("  Account: %s\n", fairShare.Account)
	fmt.Printf("  Fair-Share Factor: %.4f\n", fairShare.FairShareFactor)
	fmt.Printf("  Normalized Shares: %.4f\n", fairShare.NormalizedShares)
	fmt.Printf("  Effective Usage: %.4f\n", fairShare.EffectiveUsage)
	fmt.Printf("  Raw Shares: %d\n", fairShare.RawShares)
	fmt.Printf("  Level: %d\n", fairShare.Level)

	if fairShare.PriorityFactors != nil {
		fmt.Printf("  Priority Factors:\n")
		fmt.Printf("    Fair-Share: %d\n", fairShare.PriorityFactors.FairShare)
		fmt.Printf("    Age: %d\n", fairShare.PriorityFactors.Age)
		fmt.Printf("    Job Size: %d\n", fairShare.PriorityFactors.JobSize)
	}
}

func displayAccountFairShare(fairShare *interfaces.AccountFairShare) {
	fmt.Printf("  Account: %s\n", fairShare.AccountName)
	fmt.Printf("  Parent: %s\n", fairShare.Parent)
	fmt.Printf("  Shares: %d (Raw: %d)\n", fairShare.Shares, fairShare.RawShares)
	fmt.Printf("  Fair-Share Factor: %.4f\n", fairShare.FairShareFactor)
	fmt.Printf("  Normalized Shares: %.4f\n", fairShare.NormalizedShares)
	fmt.Printf("  Usage: %.4f (Effective: %.4f)\n", fairShare.Usage, fairShare.EffectiveUsage)
	fmt.Printf("  User Count: %d (Active: %d)\n", fairShare.UserCount, fairShare.ActiveUsers)
	fmt.Printf("  Job Count: %d\n", fairShare.JobCount)
}

func displayJobPriority(priority *interfaces.JobPriorityInfo) {
	fmt.Printf("  Calculated Priority: %d\n", priority.Priority)
	fmt.Printf("  Priority Tier: %s\n", priority.PriorityTier)
	fmt.Printf("  Estimated Start: %v\n", priority.EstimatedStart)
	fmt.Printf("  Position in Queue: %d\n", priority.PositionInQueue)

	if priority.Factors != nil {
		fmt.Printf("  Priority Breakdown:\n")
		fmt.Printf("    Fair-Share: %d\n", priority.Factors.FairShare)
		fmt.Printf("    Age: %d\n", priority.Factors.Age)
		fmt.Printf("    Job Size: %d\n", priority.Factors.JobSize)
		fmt.Printf("    Partition: %d\n", priority.Factors.Partition)
		fmt.Printf("    QoS: %d\n", priority.Factors.QoS)
		fmt.Printf("    Total: %d\n", priority.Factors.Total)
	}
}

func displayFairShareHierarchy(hierarchy *interfaces.FairShareHierarchy) {
	fmt.Printf("  Root Account: %s\n", hierarchy.RootAccount)
	fmt.Printf("  Total Shares: %d\n", hierarchy.TotalShares)
	fmt.Printf("  Total Usage: %.4f\n", hierarchy.TotalUsage)
	fmt.Printf("  Algorithm: %s\n", hierarchy.Algorithm)
	fmt.Printf("  Decay Half-Life: %d\n", hierarchy.DecayHalfLife)
	fmt.Printf("  Usage Window: %d\n", hierarchy.UsageWindow)

	if hierarchy.Tree != nil {
		fmt.Printf("\nFair-Share Tree:\n")
		printFairShareTree(hierarchy.Tree, 0)
	}
}

func printFairShareTree(node *interfaces.FairShareNode, level int) {
	indent := strings.Repeat("  ", level)
	nodeType := "Account"
	if node.User != "" {
		nodeType = "User"
	}

	fmt.Printf("%s├─ %s: %s (shares: %d, factor: %.4f)\n",
		indent, nodeType, node.Name, node.Shares, node.FairShareFactor)

	for _, child := range node.Children {
		printFairShareTree(child, level+1)
	}
}
