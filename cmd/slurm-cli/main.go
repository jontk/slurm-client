// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/interfaces"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/pkg/config"
	"github.com/spf13/cobra"
)

var (
	// Version information (set at build time)
	Version   = "dev"
	BuildTime = ""
	Commit    = ""

	// Global flags
	baseURL    string
	token      string
	username   string
	password   string
	apiVersion string
	outputFmt  string
	debug      bool

	// Root command
	rootCmd = &cobra.Command{
		Use:     "slurm-cli",
		Short:   "CLI tool for SLURM REST API",
		Long:    `A command-line interface for interacting with SLURM clusters via the REST API.`,
		Version: Version,
	}
)

func init() {
	// Custom version output
	rootCmd.Version = fmt.Sprintf("%s (commit: %s, built: %s)", Version, Commit, BuildTime)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&baseURL, "url", "", "SLURM REST API URL (env: SLURM_REST_URL)")
	rootCmd.PersistentFlags().StringVar(&token, "token", "", "JWT authentication token (env: SLURM_JWT)")
	rootCmd.PersistentFlags().StringVar(&username, "username", "", "Basic auth username (env: SLURM_USERNAME)")
	rootCmd.PersistentFlags().StringVar(&password, "password", "", "Basic auth password (env: SLURM_PASSWORD)")
	rootCmd.PersistentFlags().StringVar(&apiVersion, "api-version", "", "API version (e.g., v0.0.42)")
	rootCmd.PersistentFlags().StringVarP(&outputFmt, "output", "o", "table", "Output format: table, json, yaml")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug logging")

	// Add subcommands
	rootCmd.AddCommand(jobsCmd)
	rootCmd.AddCommand(nodesCmd)
	rootCmd.AddCommand(partitionsCmd)
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(submitCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(docsCmd)
}

// Version command with detailed info
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("slurm-cli version %s\n", Version)
		if BuildTime != "" {
			fmt.Printf("Build Time: %s\n", BuildTime)
		}
		if Commit != "" {
			fmt.Printf("Commit:     %s\n", Commit)
		}
		fmt.Printf("\nSupported API versions:\n")
		for _, v := range slurm.SupportedVersions() {
			fmt.Printf("  - %s", v)
			if v == slurm.StableVersion() {
				fmt.Printf(" (stable)")
			}
			if v == slurm.LatestVersion() {
				fmt.Printf(" (latest)")
			}
			fmt.Println()
		}
	},
}

// createClient creates a SLURM client with the provided configuration
func createClient() (slurm.SlurmClient, error) {
	// Create configuration
	cfg := config.NewDefault()

	// Override with flags or environment variables
	if baseURL != "" {
		cfg.BaseURL = baseURL
	} else if url := os.Getenv("SLURM_REST_URL"); url != "" {
		cfg.BaseURL = url
	}

	if cfg.BaseURL == "" {
		return nil, fmt.Errorf("SLURM REST API URL is required (use --url or set SLURM_REST_URL)")
	}

	cfg.Debug = debug

	// Create authentication provider
	var authProvider auth.Provider
	if token != "" {
		authProvider = auth.NewTokenAuth(token)
	} else if t := os.Getenv("SLURM_JWT"); t != "" {
		authProvider = auth.NewTokenAuth(t)
	} else if username != "" && password != "" {
		authProvider = auth.NewBasicAuth(username, password)
	} else if u := os.Getenv("SLURM_USERNAME"); u != "" {
		p := os.Getenv("SLURM_PASSWORD")
		authProvider = auth.NewBasicAuth(u, p)
	} else {
		authProvider = auth.NewNoAuth()
	}

	// Create client options
	opts := []slurm.ClientOption{
		slurm.WithConfig(cfg),
		slurm.WithAuth(authProvider),
	}

	// Create client
	ctx := context.Background()
	if apiVersion != "" {
		return slurm.NewClientWithVersion(ctx, apiVersion, opts...)
	}
	return slurm.NewClient(ctx, opts...)
}

// printOutput prints data in the requested format
func printOutput(data interface{}) error {
	switch outputFmt {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(data)
	case "yaml":
		// For simplicity, we'll just use JSON for now
		// In a real implementation, you'd use a YAML library
		fmt.Println("# YAML output not implemented, showing JSON:")
		return printOutput(data)
	default:
		// Table format - custom per data type
		return nil
	}
}

// Jobs command
var jobsCmd = &cobra.Command{
	Use:   "jobs",
	Short: "Manage jobs",
	Long:  `List, view, and manage SLURM jobs.`,
}

var jobsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List jobs",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := createClient()
		if err != nil {
			log.Fatal(err)
		}

		// Get flags
		userID, _ := cmd.Flags().GetString("user")
		states, _ := cmd.Flags().GetStringSlice("states")
		partition, _ := cmd.Flags().GetString("partition")
		limit, _ := cmd.Flags().GetInt("limit")

		// Create options
		opts := &interfaces.ListJobsOptions{
			UserID:    userID,
			States:    states,
			Partition: partition,
			Limit:     limit,
		}

		// List jobs
		ctx := context.Background()
		jobList, err := client.Jobs().List(ctx, opts)
		if err != nil {
			log.Fatal(err)
		}

		// Output results
		if outputFmt == "table" {
			fmt.Printf("%-10s %-20s %-15s %-10s %-15s\n", "JOB ID", "NAME", "USER", "STATE", "PARTITION")
			fmt.Println(strings.Repeat("-", 75))
			for _, job := range jobList.Jobs {
				fmt.Printf("%-10s %-20s %-15s %-10s %-15s\n",
					job.ID, job.Name, job.UserID, job.State, job.Partition)
			}
			fmt.Printf("\nTotal: %d jobs\n", jobList.Total)
		} else {
			printOutput(jobList)
		}
	},
}

var jobsGetCmd = &cobra.Command{
	Use:   "get JOB_ID",
	Short: "Get job details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := createClient()
		if err != nil {
			log.Fatal(err)
		}

		jobID := args[0]
		ctx := context.Background()
		job, err := client.Jobs().Get(ctx, jobID)
		if err != nil {
			log.Fatal(err)
		}

		if outputFmt == "table" {
			fmt.Printf("Job ID:      %s\n", job.ID)
			fmt.Printf("Name:        %s\n", job.Name)
			fmt.Printf("User:        %s\n", job.UserID)
			fmt.Printf("State:       %s\n", job.State)
			fmt.Printf("Partition:   %s\n", job.Partition)
			fmt.Printf("CPUs:        %d\n", job.CPUs)
			fmt.Printf("Memory:      %d MB\n", job.Memory)
			fmt.Printf("Time Limit:  %d minutes\n", job.TimeLimit)
			if !job.SubmitTime.IsZero() {
				fmt.Printf("Submit Time: %s\n", job.SubmitTime.Format(time.DateTime))
			}
			if job.StartTime != nil {
				fmt.Printf("Start Time:  %s\n", job.StartTime.Format(time.DateTime))
			}
		} else {
			printOutput(job)
		}
	},
}

var jobsCancelCmd = &cobra.Command{
	Use:   "cancel JOB_ID",
	Short: "Cancel a job",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := createClient()
		if err != nil {
			log.Fatal(err)
		}

		jobID := args[0]
		ctx := context.Background()
		err = client.Jobs().Cancel(ctx, jobID)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Job %s cancelled successfully\n", jobID)
	},
}

func init() {
	// Jobs list flags
	jobsListCmd.Flags().StringP("user", "u", "", "Filter by user ID")
	jobsListCmd.Flags().StringSliceP("states", "s", nil, "Filter by job states (comma-separated)")
	jobsListCmd.Flags().StringP("partition", "p", "", "Filter by partition")
	jobsListCmd.Flags().IntP("limit", "l", 0, "Limit number of results")

	// Add subcommands
	jobsCmd.AddCommand(jobsListCmd)
	jobsCmd.AddCommand(jobsGetCmd)
	jobsCmd.AddCommand(jobsCancelCmd)
}

// Nodes command
var nodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "Manage nodes",
	Long:  `List, view, and manage SLURM compute nodes.`,
}

var nodesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List nodes",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := createClient()
		if err != nil {
			log.Fatal(err)
		}

		// Get flags
		states, _ := cmd.Flags().GetStringSlice("states")
		partition, _ := cmd.Flags().GetString("partition")

		// Create options
		opts := &interfaces.ListNodesOptions{
			States:    states,
			Partition: partition,
		}

		// List nodes
		ctx := context.Background()
		nodeList, err := client.Nodes().List(ctx, opts)
		if err != nil {
			log.Fatal(err)
		}

		// Output results
		if outputFmt == "table" {
			fmt.Printf("%-20s %-15s %-10s %-10s %-30s\n", "NODE", "STATE", "CPUS", "MEMORY", "PARTITIONS")
			fmt.Println(strings.Repeat("-", 90))
			for _, node := range nodeList.Nodes {
				partitions := strings.Join(node.Partitions, ",")
				if len(partitions) > 30 {
					partitions = partitions[:27] + "..."
				}
				fmt.Printf("%-20s %-15s %-10d %-10d %-30s\n",
					node.Name, node.State, node.CPUs, node.Memory, partitions)
			}
			fmt.Printf("\nTotal: %d nodes\n", nodeList.Total)
		} else {
			printOutput(nodeList)
		}
	},
}

var nodesGetCmd = &cobra.Command{
	Use:   "get NODE_NAME",
	Short: "Get node details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := createClient()
		if err != nil {
			log.Fatal(err)
		}

		nodeName := args[0]
		ctx := context.Background()
		node, err := client.Nodes().Get(ctx, nodeName)
		if err != nil {
			log.Fatal(err)
		}

		if outputFmt == "table" {
			fmt.Printf("Node Name:     %s\n", node.Name)
			fmt.Printf("State:         %s\n", node.State)
			fmt.Printf("CPUs:          %d\n", node.CPUs)
			fmt.Printf("Memory:        %d MB\n", node.Memory)
			fmt.Printf("Architecture:  %s\n", node.Architecture)
			fmt.Printf("Partitions:    %s\n", strings.Join(node.Partitions, ", "))
			if len(node.Features) > 0 {
				fmt.Printf("Features:      %s\n", strings.Join(node.Features, ", "))
			}
			if node.Reason != "" {
				fmt.Printf("Reason:        %s\n", node.Reason)
			}
		} else {
			printOutput(node)
		}
	},
}

func init() {
	// Nodes list flags
	nodesListCmd.Flags().StringSliceP("states", "s", nil, "Filter by node states")
	nodesListCmd.Flags().StringP("partition", "p", "", "Filter by partition")

	// Add subcommands
	nodesCmd.AddCommand(nodesListCmd)
	nodesCmd.AddCommand(nodesGetCmd)
}

// Partitions command
var partitionsCmd = &cobra.Command{
	Use:   "partitions",
	Short: "Manage partitions",
	Long:  `List and view SLURM partitions.`,
}

var partitionsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List partitions",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := createClient()
		if err != nil {
			log.Fatal(err)
		}

		// Get flags
		states, _ := cmd.Flags().GetStringSlice("states")

		// Create options
		opts := &interfaces.ListPartitionsOptions{
			States: states,
		}

		// List partitions
		ctx := context.Background()
		partitionList, err := client.Partitions().List(ctx, opts)
		if err != nil {
			log.Fatal(err)
		}

		// Output results
		if outputFmt == "table" {
			fmt.Printf("%-20s %-10s %-10s %-10s %-15s %-15s\n",
				"PARTITION", "STATE", "NODES", "CPUS", "MAX TIME", "DEFAULT TIME")
			fmt.Println(strings.Repeat("-", 85))
			for _, partition := range partitionList.Partitions {
				fmt.Printf("%-20s %-10s %-10d %-10d %-15d %-15d\n",
					partition.Name, partition.State, partition.TotalNodes,
					partition.TotalCPUs, partition.MaxTime, partition.DefaultTime)
			}
			fmt.Printf("\nTotal: %d partitions\n", partitionList.Total)
		} else {
			printOutput(partitionList)
		}
	},
}

func init() {
	// Partitions list flags
	partitionsListCmd.Flags().StringSliceP("states", "s", nil, "Filter by partition states")

	// Add subcommands
	partitionsCmd.AddCommand(partitionsListCmd)
}

// Info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get cluster information",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := createClient()
		if err != nil {
			log.Fatal(err)
		}

		ctx := context.Background()
		info, err := client.Info().Get(ctx)
		if err != nil {
			log.Fatal(err)
		}

		if outputFmt == "table" {
			fmt.Printf("Cluster Information\n")
			fmt.Println(strings.Repeat("-", 40))
			fmt.Printf("Version:        %s\n", info.Version)
			fmt.Printf("Release:        %s\n", info.Release)
			fmt.Printf("Cluster Name:   %s\n", info.ClusterName)
			fmt.Printf("API Version:    %s\n", info.APIVersion)
			fmt.Printf("Uptime:         %d seconds\n", info.Uptime)
		} else {
			printOutput(info)
		}
	},
}

// Submit command
var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit a job",
	Long:  `Submit a new job to the SLURM cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := createClient()
		if err != nil {
			log.Fatal(err)
		}

		// Get flags
		name, _ := cmd.Flags().GetString("name")
		command, _ := cmd.Flags().GetString("command")
		partition, _ := cmd.Flags().GetString("partition")
		cpus, _ := cmd.Flags().GetInt("cpus")
		memory, _ := cmd.Flags().GetInt("memory")
		timeLimit, _ := cmd.Flags().GetInt("time")
		workDir, _ := cmd.Flags().GetString("workdir")

		if command == "" {
			log.Fatal("Command is required (--command)")
		}

		// Create job submission
		job := &interfaces.JobSubmission{
			Name:       name,
			Command:    command,
			Partition:  partition,
			CPUs:       cpus,
			Memory:     memory,
			TimeLimit:  timeLimit,
			WorkingDir: workDir,
		}

		// Submit job
		ctx := context.Background()
		resp, err := client.Jobs().Submit(ctx, job)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Job submitted successfully!\n")
		fmt.Printf("Job ID: %s\n", resp.JobID)
	},
}

func init() {
	// Submit flags
	submitCmd.Flags().StringP("name", "n", "", "Job name")
	submitCmd.Flags().StringP("command", "c", "", "Command to run (required)")
	submitCmd.Flags().StringP("partition", "p", "", "Partition name")
	submitCmd.Flags().IntP("cpus", "", 1, "Number of CPUs")
	submitCmd.Flags().IntP("memory", "m", 1024, "Memory in MB")
	submitCmd.Flags().IntP("time", "t", 60, "Time limit in minutes")
	submitCmd.Flags().StringP("workdir", "w", "", "Working directory")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
