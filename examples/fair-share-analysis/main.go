// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/interfaces"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/pkg/errors"
)

var (
	baseURL         = flag.String("url", "http://localhost:6820", "SLURM REST API URL")
	apiVersion      = flag.String("version", "v0.0.43", "API version to use")
	analyze         = flag.String("analyze", "user", "Analysis type: user, account, or hierarchy")
	target          = flag.String("target", "", "User or account name to analyze")
	compareUsers    = flag.Bool("compare-users", false, "Compare fair-share across multiple users")
	compareAccounts = flag.Bool("compare-accounts", false, "Compare fair-share across accounts")
	predictPriority = flag.Bool("predict", false, "Predict job priority for different configurations")
	showFactors     = flag.Bool("factors", false, "Show detailed priority factor breakdown")
	outputFormat    = flag.String("format", "table", "Output format: table, csv, or json")
)

func main() {
	flag.Parse()

	if *target == "" && !*compareUsers && !*compareAccounts {
		fmt.Fprintf(os.Stderr, "Error: -target is required unless using -compare-users or -compare-accounts\n")
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

	if err := run(ctx, client); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		_ = client.Close()
		os.Exit(1)
	}
	_ = client.Close()
}

func run(ctx context.Context, client slurm.SlurmClient) error {
	fmt.Printf("Connected to SLURM REST API %s at %s\n\n", client.Version(), *baseURL)

	// Perform analysis based on type
	switch *analyze {
	case "user":
		if *compareUsers {
			compareUserFairShare(ctx, client)
		} else if *target != "" {
			analyzeUserFairShare(ctx, client, *target)
		} else {
			return fmt.Errorf("specify -target or use -compare-users")
		}
	case "account":
		if *compareAccounts {
			compareAccountFairShare(ctx, client)
		} else if *target != "" {
			analyzeAccountFairShare(ctx, client, *target)
		} else {
			return fmt.Errorf("specify -target or use -compare-accounts")
		}
	case "hierarchy":
		analyzeHierarchy(ctx, client, *target)
	default:
		return fmt.Errorf("invalid analysis type. Use 'user', 'account', or 'hierarchy'")
	}

	// Predict priority if requested
	if *predictPriority && *target != "" {
		predictJobPriority(ctx, client, *target)
	}
	return nil
}

func analyzeUserFairShare(ctx context.Context, client interfaces.SlurmClient, userName string) {
	fmt.Printf("=== Fair-Share Analysis for User: %s ===\n\n", userName)

	userManager := client.Users()

	// Get user fair-share information
	fairShare, err := userManager.GetUserFairShare(ctx, userName)
	if err != nil {
		if errors.IsNotImplementedError(err) {
			fmt.Println("Note: User fair-share retrieval not implemented yet")
		} else {
			fmt.Printf("Error retrieving user fair-share: %v\n", err)
		}
		return
	}

	if fairShare == nil {
		fmt.Println("No fair-share information available")
		return
	}

	// Display fair-share information
	displayUserFairShareReport(fairShare)

	// Get user accounts to show per-account fair-share
	accounts, err := userManager.GetUserAccounts(ctx, userName)
	if err == nil && len(accounts) > 0 {
		fmt.Printf("\n\nFair-Share by Account:\n")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "Account\tRole\tFair-Share Factor\tNormalized Shares\tEffective Usage")
		fmt.Fprintln(w, "-------\t----\t-----------------\t-----------------\t---------------")

		// For each account, get detailed fair-share info
		for _, acc := range accounts {
			// In a real implementation, this would get per-account fair-share
			fmt.Fprintf(w, "%s\t%s\t%.4f\t%.4f\t%.4f\n",
				acc.AccountName,
				acc.Partition,
				fairShare.FairShareFactor,
				fairShare.NormalizedShares,
				fairShare.EffectiveUsage)
		}
		w.Flush()
	}

	// Show historical trend if available
	if *showFactors {
		displayFairShareTrend(fairShare)
	}
}

func analyzeAccountFairShare(ctx context.Context, client interfaces.SlurmClient, accountName string) {
	fmt.Printf("=== Fair-Share Analysis for Account: %s ===\n\n", accountName)

	accountManager := client.Accounts()

	// Get account fair-share information
	fairShare, err := accountManager.GetAccountFairShare(ctx, accountName)
	if err != nil {
		if errors.IsNotImplementedError(err) {
			fmt.Println("Note: Account fair-share retrieval not implemented yet")
		} else {
			fmt.Printf("Error retrieving account fair-share: %v\n", err)
		}
		return
	}

	if fairShare == nil {
		fmt.Println("No fair-share information available")
		return
	}

	// Display account fair-share information
	displayAccountFairShareReport(fairShare)

	// Get account users to show user distribution
	users, err := accountManager.GetAccountUsers(ctx, accountName, nil)
	if err == nil && len(users) > 0 {
		fmt.Printf("\n\nUser Distribution:\n")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "User\tRole\tShare of Account\tActive")
		fmt.Fprintln(w, "----\t----\t----------------\t------")

		for _, user := range users {
			// Calculate user's share of account (simplified)
			sharePercent := 100.0 / float64(len(users))
			fmt.Fprintf(w, "%s\t%s\t%.2f%%\t%v\n",
				user.UserName,
				user.Role,
				sharePercent,
				user.IsActive)
		}
		w.Flush()
	}
}

func analyzeHierarchy(ctx context.Context, client interfaces.SlurmClient, rootAccount string) {
	if rootAccount == "" {
		rootAccount = "root"
	}

	fmt.Printf("=== Fair-Share Hierarchy Analysis from: %s ===\n\n", rootAccount)

	accountManager := client.Accounts()

	// Get fair-share hierarchy
	hierarchy, err := accountManager.GetFairShareHierarchy(ctx, rootAccount)
	if err != nil {
		if errors.IsNotImplementedError(err) {
			fmt.Println("Note: Fair-share hierarchy retrieval not implemented yet")
		} else {
			fmt.Printf("Error retrieving fair-share hierarchy: %v\n", err)
		}
		return
	}

	if hierarchy == nil {
		fmt.Println("No hierarchy information available")
		return
	}

	// Display hierarchy summary
	fmt.Printf("Hierarchy Summary:\n")
	fmt.Printf("  Root Account: %s\n", hierarchy.RootAccount)
	fmt.Printf("  Cluster: %s\n", hierarchy.Cluster)
	fmt.Printf("  Total Shares: %d\n", hierarchy.TotalShares)
	fmt.Printf("  Total Usage: %.4f\n", hierarchy.TotalUsage)
	fmt.Printf("  Algorithm: %s\n", hierarchy.Algorithm)
	fmt.Printf("  Decay Half-Life: %d days\n", hierarchy.DecayHalfLife)
	fmt.Printf("  Usage Window: %d days\n", hierarchy.UsageWindow)
	fmt.Printf("  Last Update: %v\n\n", hierarchy.LastUpdate.Format(time.RFC3339))

	// Display the hierarchy tree with fair-share values
	if hierarchy.Tree != nil {
		fmt.Println("Fair-Share Tree:")
		displayFairShareTree(hierarchy.Tree, 0, *outputFormat)
	}

	// Calculate and display statistics
	if hierarchy.Tree != nil {
		stats := calculateHierarchyStats(hierarchy.Tree)
		fmt.Printf("\n\nHierarchy Statistics:\n")
		fmt.Printf("  Total Nodes: %d\n", stats.totalNodes)
		fmt.Printf("  Total Accounts: %d\n", stats.totalAccounts)
		fmt.Printf("  Total Users: %d\n", stats.totalUsers)
		fmt.Printf("  Max Depth: %d\n", stats.maxDepth)
		fmt.Printf("  Average Fair-Share Factor: %.4f\n", stats.avgFairShare)
		fmt.Printf("  Min Fair-Share Factor: %.4f\n", stats.minFairShare)
		fmt.Printf("  Max Fair-Share Factor: %.4f\n", stats.maxFairShare)
	}
}

func compareUserFairShare(ctx context.Context, client interfaces.SlurmClient) {
	fmt.Println("=== User Fair-Share Comparison ===")

	// For demonstration, compare a predefined list of users
	// In practice, this could be read from a file or command line
	users := []string{"user1", "user2", "user3", "user4", "user5"}

	userManager := client.Users()

	type userFairShareData struct {
		userName  string
		fairShare *interfaces.UserFairShare
		accounts  []interfaces.UserAccount
		err       error
	}

	results := make([]userFairShareData, 0, len(users))

	// Collect fair-share data for all users
	for _, userName := range users {
		data := userFairShareData{userName: userName}

		data.fairShare, data.err = userManager.GetUserFairShare(ctx, userName)
		if data.err == nil {
			pointerAccounts, _ := userManager.GetUserAccounts(ctx, userName)
			// Convert from []*UserAccount to []UserAccount
			data.accounts = make([]interfaces.UserAccount, len(pointerAccounts))
			for i, acc := range pointerAccounts {
				data.accounts[i] = *acc
			}
		}

		results = append(results, data)
	}

	// Sort by fair-share factor (descending)
	sort.Slice(results, func(i, j int) bool {
		if results[i].fairShare == nil {
			return false
		}
		if results[j].fairShare == nil {
			return true
		}
		return results[i].fairShare.FairShareFactor > results[j].fairShare.FairShareFactor
	})

	// Display comparison table
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Rank\tUser\tFair-Share\tNorm. Shares\tEff. Usage\tLevel\tPrimary Account")
	fmt.Fprintln(w, "----\t----\t----------\t------------\t----------\t-----\t---------------")

	for i, data := range results {
		if data.err != nil {
			fmt.Fprintf(w, "%d\t%s\tERROR\t-\t-\t-\t-\n", i+1, data.userName)
			continue
		}
		if data.fairShare == nil {
			fmt.Fprintf(w, "%d\t%s\tN/A\t-\t-\t-\t-\n", i+1, data.userName)
			continue
		}

		primaryAccount := data.fairShare.Account
		if len(data.accounts) > 0 {
			primaryAccount = data.accounts[0].AccountName
		}

		fmt.Fprintf(w, "%d\t%s\t%.4f\t%.4f\t%.4f\t%d\t%s\n",
			i+1,
			data.userName,
			data.fairShare.FairShareFactor,
			data.fairShare.NormalizedShares,
			data.fairShare.EffectiveUsage,
			data.fairShare.Level,
			primaryAccount)
	}
	w.Flush()
}

func compareAccountFairShare(ctx context.Context, client interfaces.SlurmClient) {
	fmt.Println("=== Account Fair-Share Comparison ===")

	accountManager := client.Accounts()

	// Get all accounts (with limited fields for efficiency)
	opts := &interfaces.ListAccountsOptions{
		WithQuotas: false,
		WithUsers:  false,
		Limit:      20, // Top 20 accounts
	}

	accountList, err := accountManager.List(ctx, opts)
	if err != nil {
		fmt.Printf("Error listing accounts: %v\n", err)
		return
	}

	type accountFairShareData struct {
		account   interfaces.Account
		fairShare *interfaces.AccountFairShare
		err       error
	}

	results := make([]accountFairShareData, 0, len(accountList.Accounts))

	// Collect fair-share data for each account
	for _, account := range accountList.Accounts {
		data := accountFairShareData{account: account}
		data.fairShare, data.err = accountManager.GetAccountFairShare(ctx, account.Name)
		results = append(results, data)
	}

	// Sort by fair-share factor (descending)
	sort.Slice(results, func(i, j int) bool {
		if results[i].fairShare == nil {
			return false
		}
		if results[j].fairShare == nil {
			return true
		}
		return results[i].fairShare.FairShareFactor > results[j].fairShare.FairShareFactor
	})

	// Display comparison table
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Rank\tAccount\tParent\tShares\tFair-Share\tNorm. Usage\tUsers\tActive")
	fmt.Fprintln(w, "----\t-------\t------\t------\t----------\t-----------\t-----\t------")

	for i, data := range results {
		if data.err != nil || data.fairShare == nil {
			fmt.Fprintf(w, "%d\t%s\t-\t-\tERROR\t-\t-\t-\n", i+1, data.account.Name)
			continue
		}

		fmt.Fprintf(w, "%d\t%s\t%s\t%d\t%.4f\t%.4f\t%d\t%d\n",
			i+1,
			data.account.Name,
			data.fairShare.Parent,
			data.fairShare.Shares,
			data.fairShare.FairShareFactor,
			data.fairShare.EffectiveUsage,
			data.fairShare.UserCount,
			data.fairShare.ActiveUsers)
	}
	w.Flush()
}

func predictJobPriority(ctx context.Context, client interfaces.SlurmClient, userName string) {
	fmt.Printf("\n\n=== Job Priority Prediction for User: %s ===\n\n", userName)

	userManager := client.Users()

	// Define test scenarios
	scenarios := []struct {
		name string
		job  interfaces.JobSubmission
	}{
		{
			name: "Small Job (1 CPU, 1GB, 1 hour)",
			job: interfaces.JobSubmission{
				Script:    "#!/bin/bash\necho 'Small job'",
				Partition: "compute",
				CPUs:      1,
				Memory:    1024 * 1024 * 1024,
				TimeLimit: 60, // 1 hour
			},
		},
		{
			name: "Medium Job (8 CPUs, 16GB, 4 hours)",
			job: interfaces.JobSubmission{
				Script:    "#!/bin/bash\necho 'Medium job'",
				Partition: "compute",
				CPUs:      8,
				Memory:    16 * 1024 * 1024 * 1024,
				TimeLimit: 240, // 4 hours
			},
		},
		{
			name: "Large Job (32 CPUs, 64GB, 24 hours)",
			job: interfaces.JobSubmission{
				Script:    "#!/bin/bash\necho 'Large job'",
				Partition: "compute",
				CPUs:      32,
				Memory:    64 * 1024 * 1024 * 1024,
				TimeLimit: 1440, // 24 hours
			},
		},
		{
			name: "GPU Job (4 CPUs, 16GB, 2 GPUs)",
			job: interfaces.JobSubmission{
				Script:    "#!/bin/bash\necho 'GPU job'",
				Partition: "gpu",
				CPUs:      4,
				Memory:    16 * 1024 * 1024 * 1024,
				TimeLimit: 720, // 12 hours
			},
		},
	}

	// Get user's accounts
	accounts, err := userManager.GetUserAccounts(ctx, userName)
	if err != nil {
		fmt.Printf("Error retrieving user accounts: %v\n", err)
		return
	}

	// Run predictions for each scenario and account combination
	fmt.Println("Priority Predictions by Job Type and Account:")

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Job Type\tAccount\tPriority\tTier\tEst. Start\tQueue Pos")
	fmt.Fprintln(w, "--------\t-------\t--------\t----\t----------\t---------")

	for _, scenario := range scenarios {
		// Test with each account
		if len(accounts) > 0 {
			for _, acc := range accounts {
				scenario.job.Partition = acc.AccountName
				priority, err := userManager.CalculateJobPriority(ctx, userName, &scenario.job)

				if err != nil {
					if errors.IsNotImplementedError(err) {
						fmt.Fprintf(w, "%s\t%s\tN/A\t-\t-\t-\n", scenario.name, acc.AccountName)
					} else {
						fmt.Fprintf(w, "%s\t%s\tERROR\t-\t-\t-\n", scenario.name, acc.AccountName)
					}
					continue
				}

				if priority != nil {
					estStart := "Unknown"
					if !priority.EstimatedStart.IsZero() {
						estStart = priority.EstimatedStart.Format("Jan 02 15:04")
					}

					fmt.Fprintf(w, "%s\t%s\t%d\t%s\t%s\t%d\n",
						scenario.name,
						acc.AccountName,
						priority.Priority,
						priority.PriorityTier,
						estStart,
						priority.PositionInQueue)
				}
			}
		} else {
			// Test without account
			priority, err := userManager.CalculateJobPriority(ctx, userName, &scenario.job)

			if err != nil {
				fmt.Fprintf(w, "%s\t(none)\tERROR\t-\t-\t-\n", scenario.name)
				continue
			}

			if priority != nil {
				estStart := "Unknown"
				if !priority.EstimatedStart.IsZero() {
					estStart = priority.EstimatedStart.Format("Jan 02 15:04")
				}

				fmt.Fprintf(w, "%s\t(none)\t%d\t%s\t%s\t%d\n",
					scenario.name,
					priority.Priority,
					priority.PriorityTier,
					estStart,
					priority.PositionInQueue)
			}
		}
	}
	w.Flush()

	// Show priority factor breakdown if requested
	if *showFactors && len(accounts) > 0 {
		fmt.Printf("\n\nDetailed Priority Factor Breakdown:\n")

		// Use medium job as example
		mediumJob := scenarios[1].job
		mediumJob.Partition = accounts[0].AccountName

		priority, err := userManager.CalculateJobPriority(ctx, userName, &mediumJob)
		if err == nil && priority != nil && priority.Factors != nil {
			displayPriorityFactors(priority.Factors)
		}
	}
}

// Helper functions

func displayUserFairShareReport(fairShare *interfaces.UserFairShare) {
	fmt.Printf("User Fair-Share Report:\n")
	fmt.Printf("  Account: %s\n", fairShare.Account)
	fmt.Printf("  Fair-Share Factor: %.6f\n", fairShare.FairShareFactor)
	fmt.Printf("  Normalized Shares: %.6f\n", fairShare.NormalizedShares)
	fmt.Printf("  Effective Usage: %.6f\n", fairShare.EffectiveUsage)
	// RawUsage field doesn't exist in UserFairShare
	fmt.Printf("  Raw Shares: %d\n", fairShare.RawShares)
	fmt.Printf("  Hierarchical Level: %d\n", fairShare.Level)

	// LastUpdate field doesn't exist in UserFairShare/AccountFairShare
	if !fairShare.LastDecay.IsZero() {
		fmt.Printf("  Last Decay: %v\n", fairShare.LastDecay.Format(time.RFC3339))
	}

	if fairShare.PriorityFactors != nil {
		fmt.Printf("\nCurrent Priority Factors:\n")
		fmt.Printf("  Fair-Share Component: %d\n", fairShare.PriorityFactors.FairShare)
		fmt.Printf("  Age Component: %d\n", fairShare.PriorityFactors.Age)
		fmt.Printf("  Job Size Component: %d\n", fairShare.PriorityFactors.JobSize)
		fmt.Printf("  Partition Component: %d\n", fairShare.PriorityFactors.Partition)
		fmt.Printf("  QoS Component: %d\n", fairShare.PriorityFactors.QoS)
		fmt.Printf("  Total Priority: %d\n", fairShare.PriorityFactors.Total)
	}
}

func displayAccountFairShareReport(fairShare *interfaces.AccountFairShare) {
	fmt.Printf("Account Fair-Share Report:\n")
	fmt.Printf("  Account Name: %s\n", fairShare.AccountName)
	fmt.Printf("  Parent Account: %s\n", fairShare.Parent)
	fmt.Printf("  Shares: %d (Raw: %d)\n", fairShare.Shares, fairShare.RawShares)
	fmt.Printf("  Fair-Share Factor: %.6f\n", fairShare.FairShareFactor)
	fmt.Printf("  Normalized Shares: %.6f\n", fairShare.NormalizedShares)
	fmt.Printf("  Current Usage: %.6f\n", fairShare.Usage)
	fmt.Printf("  Effective Usage: %.6f\n", fairShare.EffectiveUsage)
	fmt.Printf("  Fairshare Level: %d\n", fairShare.Level)
	fmt.Printf("  Total Users: %d (Active: %d)\n", fairShare.UserCount, fairShare.ActiveUsers)
	fmt.Printf("  Total Jobs: %d\n", fairShare.JobCount)

	if len(fairShare.Children) > 0 {
		fmt.Printf("  Child Accounts: %d\n", len(fairShare.Children))
	}

	// LastUpdate field doesn't exist in UserFairShare/AccountFairShare
	if !fairShare.LastDecay.IsZero() {
		fmt.Printf("  Last Decay: %v\n", fairShare.LastDecay.Format(time.RFC3339))
	}
}

func displayFairShareTree(node *interfaces.FairShareNode, level int, format string) {
	if format == "csv" {
		// CSV format
		if level == 0 {
			fmt.Println("Type,Name,Level,Shares,FairShareFactor,Usage,NormalizedShares")
		}
		nodeType := "Account"
		if node.User != "" {
			nodeType = "User"
		}
		fmt.Printf("%s,%s,%d,%d,%.6f,%.6f,%.6f\n",
			nodeType, node.Name, level, node.Shares,
			node.FairShareFactor, node.Usage, node.NormalizedShares)
	} else {
		// Default tree format
		indent := strings.Repeat("  ", level)
		nodeType := "A"
		if node.User != "" {
			nodeType = "U"
		}

		fmt.Printf("%s├─ [%s] %s (shares: %d, factor: %.6f, usage: %.6f)\n",
			indent, nodeType, node.Name, node.Shares,
			node.FairShareFactor, node.Usage)
	}

	// Recursively display children
	for _, child := range node.Children {
		displayFairShareTree(child, level+1, format)
	}
}

func displayFairShareTrend(fairShare *interfaces.UserFairShare) {
	fmt.Printf("\n\nFair-Share Trend Analysis:\n")
	fmt.Println("(Historical data would be displayed here if available)")

	// In a real implementation, this would show:
	// - Fair-share factor over time (daily/weekly)
	// - Usage patterns
	// - Share allocation changes
	// - Priority ranking changes

	// For now, show a simple representation
	fmt.Printf("\nCurrent Status:\n")
	fmt.Printf("  Fair-Share Factor: %.6f", fairShare.FairShareFactor)

	if fairShare.FairShareFactor > 0.5 {
		fmt.Println(" (Above Average)")
	} else if fairShare.FairShareFactor > 0.1 {
		fmt.Println(" (Average)")
	} else {
		fmt.Println(" (Below Average - Heavy Usage)")
	}

	fmt.Printf("  Usage vs Allocation: %.2f%%\n",
		(fairShare.EffectiveUsage/fairShare.NormalizedShares)*100)
}

func displayPriorityFactors(factors *interfaces.JobPriorityFactors) {
	fmt.Println("\nPriority Factor Breakdown:")

	// Calculate percentages
	total := float64(factors.Total)
	if total == 0 {
		total = 1 // Avoid division by zero
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Factor\tValue\tPercentage\tDescription")
	fmt.Fprintln(w, "------\t-----\t----------\t-----------")

	fmt.Fprintf(w, "Fair-Share\t%d\t%.1f%%\tBased on historical usage vs allocation\n",
		factors.FairShare, (float64(factors.FairShare)/total)*100)
	fmt.Fprintf(w, "Age\t%d\t%.1f%%\tTime since job submission\n",
		factors.Age, (float64(factors.Age)/total)*100)
	fmt.Fprintf(w, "Job Size\t%d\t%.1f%%\tFavors smaller jobs\n",
		factors.JobSize, (float64(factors.JobSize)/total)*100)
	fmt.Fprintf(w, "Partition\t%d\t%.1f%%\tPartition-specific priority\n",
		factors.Partition, (float64(factors.Partition)/total)*100)
	fmt.Fprintf(w, "QoS\t%d\t%.1f%%\tQuality of Service level\n",
		factors.QoS, (float64(factors.QoS)/total)*100)
	fmt.Fprintf(w, "TOTAL\t%d\t100.0%%\tFinal job priority\n", factors.Total)

	w.Flush()
}

type hierarchyStats struct {
	totalNodes    int
	totalAccounts int
	totalUsers    int
	maxDepth      int
	avgFairShare  float64
	minFairShare  float64
	maxFairShare  float64
	fairShareSum  float64
}

func calculateHierarchyStats(node *interfaces.FairShareNode) hierarchyStats {
	stats := hierarchyStats{
		minFairShare: 1.0,
		maxFairShare: 0.0,
	}

	calculateNodeStats(node, 0, &stats)

	if stats.totalNodes > 0 {
		stats.avgFairShare = stats.fairShareSum / float64(stats.totalNodes)
	}

	return stats
}

func calculateNodeStats(node *interfaces.FairShareNode, depth int, stats *hierarchyStats) {
	stats.totalNodes++
	stats.fairShareSum += node.FairShareFactor

	if node.FairShareFactor < stats.minFairShare {
		stats.minFairShare = node.FairShareFactor
	}
	if node.FairShareFactor > stats.maxFairShare {
		stats.maxFairShare = node.FairShareFactor
	}

	if depth > stats.maxDepth {
		stats.maxDepth = depth
	}

	if node.User != "" {
		stats.totalUsers++
	} else {
		stats.totalAccounts++
	}

	for _, child := range node.Children {
		calculateNodeStats(child, depth+1, stats)
	}
}
