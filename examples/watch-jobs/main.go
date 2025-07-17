package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/internal/interfaces"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/pkg/config"
)

func main() {
	// Create configuration
	cfg := config.NewDefault()
	
	// Override with environment variables if needed
	if url := os.Getenv("SLURM_REST_URL"); url != "" {
		cfg.BaseURL = url
	}

	// Create authentication provider
	var authProvider auth.Provider
	if token := os.Getenv("SLURM_JWT"); token != "" {
		authProvider = auth.NewTokenAuth(token)
	} else if username := os.Getenv("SLURM_USERNAME"); username != "" {
		password := os.Getenv("SLURM_PASSWORD")
		authProvider = auth.NewBasicAuth(username, password)
	} else {
		authProvider = auth.NewNoAuth()
	}

	// Create client
	ctx := context.Background()
	client, err := slurm.NewClient(ctx,
		slurm.WithConfig(cfg),
		slurm.WithAuth(authProvider),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Get job manager
	jobManager := client.Jobs()

	// Set up watch options
	watchOpts := &interfaces.WatchJobsOptions{
		// Watch specific user's jobs (optional)
		// UserID: "1000",
		
		// Watch specific states (optional)
		// States: []string{"RUNNING", "PENDING"},
		
		// Watch specific job IDs (optional)
		// JobIDs: []string{"12345", "12346"},
	}

	// Start watching for job events
	fmt.Println("Starting to watch for job events...")
	fmt.Println("Press Ctrl+C to stop")
	
	watchCtx, cancelWatch := context.WithCancel(ctx)
	defer cancelWatch()

	eventChan, err := jobManager.Watch(watchCtx, watchOpts)
	if err != nil {
		log.Fatalf("Failed to start watching jobs: %v", err)
	}

	// Handle interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Process events
	for {
		select {
		case event, ok := <-eventChan:
			if !ok {
				fmt.Println("Event channel closed")
				return
			}

			// Handle different event types
			switch event.Type {
			case "job_new":
				fmt.Printf("[%s] New job detected: ID=%s, State=%s\n",
					event.Timestamp.Format(time.RFC3339),
					event.JobID,
					event.NewState)
				if event.Job != nil {
					fmt.Printf("  User: %s, Partition: %s\n", event.Job.UserID, event.Job.Partition)
				}

			case "job_state_change":
				fmt.Printf("[%s] Job state changed: ID=%s, %s -> %s\n",
					event.Timestamp.Format(time.RFC3339),
					event.JobID,
					event.OldState,
					event.NewState)
				if event.Job != nil {
					fmt.Printf("  User: %s, Partition: %s\n", event.Job.UserID, event.Job.Partition)
				}

			case "job_completed":
				fmt.Printf("[%s] Job completed: ID=%s (was %s)\n",
					event.Timestamp.Format(time.RFC3339),
					event.JobID,
					event.OldState)

			case "error":
				fmt.Printf("[%s] Error: %v\n",
					event.Timestamp.Format(time.RFC3339),
					event.Error)

			default:
				fmt.Printf("[%s] Unknown event type: %s\n",
					event.Timestamp.Format(time.RFC3339),
					event.Type)
			}

		case <-sigChan:
			fmt.Println("\nShutting down...")
			cancelWatch()
			// Give a moment for cleanup
			time.Sleep(100 * time.Millisecond)
			return
		}
	}
}