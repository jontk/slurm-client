// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/jontk/slurm-client"
	"github.com/jontk/slurm-client/pkg/auth"
	"github.com/jontk/slurm-client/pkg/config"
)

// Example: Resource allocation patterns and constraints
func main() {
	// Create configuration
	cfg := config.NewDefault()
	cfg.BaseURL = "https://cluster.example.com:6820"

	// Create authentication
	authProvider := auth.NewTokenAuth("your-jwt-token")

	// Create client
	ctx := context.Background()
	client, err := slurm.NewClient(ctx,
		slurm.WithConfig(cfg),
		slurm.WithAuth(authProvider),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Example 1: GPU allocation
	fmt.Println("=== GPU Resource Allocation ===")
	allocateGPUResources(ctx, client)

	// Example 2: Memory-intensive jobs
	fmt.Println("\n=== Memory-Intensive Jobs ===")
	allocateHighMemoryJobs(ctx, client)

	// Example 3: Node-specific constraints
	fmt.Println("\n=== Node-Specific Constraints ===")
	allocateWithNodeConstraints(ctx, client)

	// Example 4: Resource sharing and exclusive allocation
	fmt.Println("\n=== Resource Sharing Patterns ===")
	demonstrateResourceSharing(ctx, client)

	// Example 5: Dynamic resource discovery
	fmt.Println("\n=== Dynamic Resource Discovery ===")
	discoverAndAllocateResources(ctx, client)
}

// allocateGPUResources demonstrates GPU resource allocation patterns
func allocateGPUResources(ctx context.Context, client slurm.SlurmClient) {
	// Single GPU job
	singleGPUJob := &slurm.JobCreate{
		Name: ptrString("single-gpu-job"),
		Script: ptrString(`#!/bin/bash
#SBATCH --gres=gpu:1
#SBATCH --constraint=gpu_mem_32gb

echo "Running on node: $SLURMD_NODENAME"
echo "GPU allocated: $CUDA_VISIBLE_DEVICES"

# Show GPU info
nvidia-smi

# Run GPU workload
python3 train_model.py --gpu 0 --batch-size 128
`),
		Partition:     ptrString("gpu"),
		MinimumCPUs:   ptrInt32(8),
		MemoryPerNode: ptrUint64(32768), // 32GB
		TimeLimit:     ptrUint32(120),
		// Note: GRES and constraints would be specified in the SBATCH script
		// as JobCreate doesn't have a Metadata field
	}

	resp1, err := client.Jobs().SubmitRaw(ctx, singleGPUJob)
	if err != nil {
		log.Printf("Failed to submit single GPU job: %v", err)
	} else {
		fmt.Printf("Single GPU job submitted: %s\n", fmt.Sprintf("%d", resp1.JobId))
	}

	// Multi-GPU job
	multiGPUJob := &slurm.JobCreate{
		Name: ptrString("multi-gpu-job"),
		Script: ptrString(`#!/bin/bash
#SBATCH --gres=gpu:4
#SBATCH --constraint=v100
#SBATCH --ntasks=4
#SBATCH --cpus-per-task=8

echo "Running distributed training on 4 GPUs"
echo "GPUs allocated: $CUDA_VISIBLE_DEVICES"

# Run distributed training
srun python3 distributed_train.py \
    --world-size 4 \
    --rank $SLURM_PROCID \
    --master-addr $SLURM_LAUNCH_NODE_IPADDR \
    --master-port 29500
`),
		Partition:     ptrString("gpu"),
		MinimumCPUs:   ptrInt32(32),     // 4 tasks * 8 CPUs
		MemoryPerNode: ptrUint64(131072), // 128GB total
		TimeLimit:     ptrUint32(360),
		MinimumNodes:  ptrInt32(1), // All GPUs on same node
		// Metadata would be in SBATCH directives
		// 		// Removed: Metadata: map[string]interface{}{
		// 			"gres":       "gpu:4",
		// 			"constraint": "v100",
		// 			"ntasks":     4,
		// 		},
	}

	resp2, err := client.Jobs().SubmitRaw(ctx, multiGPUJob)
	if err != nil {
		log.Printf("Failed to submit multi-GPU job: %v", err)
	} else {
		fmt.Printf("Multi-GPU job submitted: %s\n", fmt.Sprintf("%d", resp2.JobId))
	}

	// GPU type-specific job
	gpuTypeJob := &slurm.JobCreate{
		Name: ptrString("gpu-type-specific"),
		Script: ptrString(`#!/bin/bash
#SBATCH --gres=gpu:a100:2
#SBATCH --partition=gpu-a100

echo "Running on A100 GPUs"
nvidia-smi --query-gpu=name,memory.total --format=csv

# Run workload optimized for A100
python3 inference.py --model large_model.pt --precision fp16
`),
		Partition:     ptrString("gpu-a100"),
		MinimumCPUs:   ptrInt32(16),
		MemoryPerNode: ptrUint64(65536),
		TimeLimit:     ptrUint32(60),
		// Metadata would be in SBATCH directives
		// 		// Removed: Metadata: map[string]interface{}{
		// 			"gres": "gpu:a100:2",
		// 		},
	}

	resp3, err := client.Jobs().SubmitRaw(ctx, gpuTypeJob)
	if err != nil {
		log.Printf("Failed to submit GPU type-specific job: %v", err)
	} else {
		fmt.Printf("GPU type-specific job submitted: %s\n", fmt.Sprintf("%d", resp3.JobId))
	}
}

// allocateHighMemoryJobs demonstrates memory-intensive job patterns
func allocateHighMemoryJobs(ctx context.Context, client slurm.SlurmClient) {
	// Standard memory job
	standardMemJob := &slurm.JobCreate{
		Name: ptrString("standard-memory-job"),
		Script: ptrString(`#!/bin/bash
echo "Allocated memory: ${SLURM_MEM_PER_NODE}MB"
echo "Memory per CPU: ${SLURM_MEM_PER_CPU}MB"

# Run memory-aware application
python3 analyze_data.py --max-memory ${SLURM_MEM_PER_NODE}
`),
		Partition:     ptrString("compute"),
		MinimumCPUs:   ptrInt32(4),
		MemoryPerNode: ptrUint64(16384), // 16GB total
		TimeLimit:     ptrUint32(30),
	}

	resp1, err := client.Jobs().SubmitRaw(ctx, standardMemJob)
	if err != nil {
		log.Printf("Failed to submit standard memory job: %v", err)
	} else {
		fmt.Printf("Standard memory job submitted: %s (4GB per CPU)\n", fmt.Sprintf("%d", resp1.JobId))
	}

	// High memory job with specific memory-per-cpu
	highMemJob := &slurm.JobCreate{
		Name: ptrString("high-memory-job"),
		Script: ptrString(`#!/bin/bash
#SBATCH --mem-per-cpu=32G
#SBATCH --constraint=highmem

echo "Running high-memory analysis"
free -h

# Process large dataset in memory
python3 process_large_dataset.py --input /data/huge_file.csv
`),
		Partition:     ptrString("highmem"),
		MinimumCPUs:   ptrInt32(8),
		MemoryPerNode: ptrUint64(262144), // 256GB total (32GB per CPU)
		TimeLimit:     ptrUint32(120),
		// Metadata would be in SBATCH directives
		// 		// Removed: Metadata: map[string]interface{}{
		// 			"mem-per-cpu": "32G",
		// 			"constraint":  "highmem",
		// 		},
	}

	resp2, err := client.Jobs().SubmitRaw(ctx, highMemJob)
	if err != nil {
		log.Printf("Failed to submit high memory job: %v", err)
	} else {
		fmt.Printf("High memory job submitted: %s (32GB per CPU)\n", fmt.Sprintf("%d", resp2.JobId))
	}

	// Memory reservation pattern
	memReserveJob := &slurm.JobCreate{
		Name: ptrString("memory-reservation"),
		Script: ptrString(`#!/bin/bash
#SBATCH --mem=0  # Request all available memory on the node
#SBATCH --exclusive

echo "Exclusive node allocation with all memory"
echo "Total node memory: $(free -h | grep Mem | awk '{print $2}')"

# Run memory-intensive workload
./memory_intensive_app --use-all-available-memory
`),
		Partition:     ptrString("compute"),
		MinimumCPUs:   ptrInt32(48), // Full node
		MemoryPerNode: ptrUint64(0), // All available
		TimeLimit:     ptrUint32(60),
		// Metadata would be in SBATCH directives
		// 		// Removed: Metadata: map[string]interface{}{
		// 			"exclusive": true,
		// 			"mem":       "0",
		// 		},
	}

	resp3, err := client.Jobs().SubmitRaw(ctx, memReserveJob)
	if err != nil {
		log.Printf("Failed to submit memory reservation job: %v", err)
	} else {
		fmt.Printf("Memory reservation job submitted: %s (exclusive node)\n", fmt.Sprintf("%d", resp3.JobId))
	}
}

// allocateWithNodeConstraints demonstrates node-specific resource constraints
func allocateWithNodeConstraints(ctx context.Context, client slurm.SlurmClient) {
	// CPU architecture constraint
	archJob := &slurm.JobCreate{
		Name: ptrString("arch-specific-job"),
		Script: ptrString(`#!/bin/bash
#SBATCH --constraint="haswell|broadwell"

echo "Running on CPU architecture: $(lscpu | grep 'Model name' | cut -d: -f2)"
echo "Node: $SLURMD_NODENAME"

# Run architecture-optimized code
./optimized_binary --arch $(lscpu | grep 'Architecture' | cut -d: -f2 | xargs)
`),
		Partition:     ptrString("compute"),
		MinimumCPUs:   ptrInt32(16),
		MemoryPerNode: ptrUint64(32768),
		TimeLimit:     ptrUint32(45),
		// Metadata would be in SBATCH directives
		// 		// Removed: Metadata: map[string]interface{}{
		// 			"constraint": "haswell|broadwell",
		// 		},
	}

	resp1, err := client.Jobs().SubmitRaw(ctx, archJob)
	if err != nil {
		log.Printf("Failed to submit architecture-specific job: %v", err)
	} else {
		fmt.Printf("Architecture-specific job submitted: %s\n", fmt.Sprintf("%d", resp1.JobId))
	}

	// Network topology constraint
	networkJob := &slurm.JobCreate{
		Name: ptrString("network-topology-job"),
		Script: ptrString(`#!/bin/bash
#SBATCH --constraint="ib&rack3"
#SBATCH --switches=1

echo "Running on InfiniBand-enabled nodes in rack 3"
echo "Ensuring minimal network hops for MPI communication"

# Run MPI job with optimal network topology
mpirun -np $SLURM_NTASKS ./mpi_application
`),
		Partition:     ptrString("compute"),
		MinimumCPUs:   ptrInt32(64),
		MemoryPerNode: ptrUint64(131072),
		MinimumNodes:  ptrInt32(4),
		TimeLimit:     ptrUint32(180),
		// Metadata would be in SBATCH directives
		// 		// Removed: Metadata: map[string]interface{}{
		// 			"constraint": "ib&rack3",
		// 			"switches":   1,
		// 			"ntasks":     64,
		// 		},
	}

	resp2, err := client.Jobs().SubmitRaw(ctx, networkJob)
	if err != nil {
		log.Printf("Failed to submit network topology job: %v", err)
	} else {
		fmt.Printf("Network topology job submitted: %s\n", fmt.Sprintf("%d", resp2.JobId))
	}

	// Feature-based constraint
	featureJob := &slurm.JobCreate{
		Name: ptrString("feature-based-job"),
		Script: ptrString(`#!/bin/bash
#SBATCH --constraint="ssd&gpu&centos7"

echo "Running on nodes with SSD, GPU, and CentOS 7"
echo "Features available: $SLURM_JOB_CONSTRAINTS"

# Utilize node features
df -h | grep -E 'ssd|nvme'  # Show SSD storage
nvidia-smi --list-gpus       # Show GPUs
cat /etc/redhat-release      # Show OS version
`),
		Partition:     ptrString("mixed"),
		MinimumCPUs:   ptrInt32(8),
		MemoryPerNode: ptrUint64(16384),
		TimeLimit:     ptrUint32(30),
		// Metadata would be in SBATCH directives
		// 		// Removed: Metadata: map[string]interface{}{
		// 			"constraint": "ssd&gpu&centos7",
		// 			"gres":       "gpu:1",
		// 		},
	}

	resp3, err := client.Jobs().SubmitRaw(ctx, featureJob)
	if err != nil {
		log.Printf("Failed to submit feature-based job: %v", err)
	} else {
		fmt.Printf("Feature-based job submitted: %s\n", fmt.Sprintf("%d", resp3.JobId))
	}
}

// demonstrateResourceSharing shows different resource sharing patterns
func demonstrateResourceSharing(ctx context.Context, client slurm.SlurmClient) {
	// Exclusive node allocation
	exclusiveJob := &slurm.JobCreate{
		Name: ptrString("exclusive-node"),
		Script: ptrString(`#!/bin/bash
#SBATCH --exclusive
#SBATCH --nodes=2

echo "Exclusive allocation of 2 nodes"
echo "No other jobs can run on these nodes"
scontrol show node $SLURM_JOB_NODELIST
`),
		Partition:     ptrString("compute"),
		MinimumCPUs:   ptrInt32(96), // 2 nodes * 48 CPUs
		MemoryPerNode: ptrUint64(0), // All available
		MinimumNodes:  ptrInt32(2),
		TimeLimit:     ptrUint32(60),
		// Metadata would be in SBATCH directives
		// 		// Removed: Metadata: map[string]interface{}{
		// 			"exclusive": true,
		// 		},
	}

	resp1, err := client.Jobs().SubmitRaw(ctx, exclusiveJob)
	if err != nil {
		log.Printf("Failed to submit exclusive job: %v", err)
	} else {
		fmt.Printf("Exclusive node job submitted: %s\n", fmt.Sprintf("%d", resp1.JobId))
	}

	// Shared node allocation
	sharedJob := &slurm.JobCreate{
		Name: ptrString("shared-resources"),
		Script: ptrString(`#!/bin/bash
#SBATCH --oversubscribe
#SBATCH --mem=8G

echo "Shared node allocation"
echo "Other jobs can run on the same node"
echo "Allocated CPUs: $SLURM_CPUS_ON_NODE"
echo "Allocated memory: ${SLURM_MEM_PER_NODE}MB"
`),
		Partition:     ptrString("shared"),
		MinimumCPUs:   ptrInt32(2),
		MemoryPerNode: ptrUint64(8192),
		TimeLimit:     ptrUint32(30),
		// Metadata would be in SBATCH directives
		// 		// Removed: Metadata: map[string]interface{}{
		// 			"oversubscribe": true,
		// 		},
	}

	resp2, err := client.Jobs().SubmitRaw(ctx, sharedJob)
	if err != nil {
		log.Printf("Failed to submit shared job: %v", err)
	} else {
		fmt.Printf("Shared resource job submitted: %s\n", fmt.Sprintf("%d", resp2.JobId))
	}

	// Core binding pattern
	coreBindJob := &slurm.JobCreate{
		Name: ptrString("core-binding"),
		Script: ptrString(`#!/bin/bash
#SBATCH --cpu-bind=cores
#SBATCH --ntasks=4
#SBATCH --cpus-per-task=2

echo "Core binding for NUMA optimization"
srun --cpu-bind=verbose hostname
srun numactl --show

# Run NUMA-aware application
srun ./numa_optimized_app
`),
		Partition:     ptrString("compute"),
		MinimumCPUs:   ptrInt32(8), // 4 tasks * 2 CPUs
		MemoryPerNode: ptrUint64(16384),
		TimeLimit:     ptrUint32(45),
		// Metadata would be in SBATCH directives
		// 		// Removed: Metadata: map[string]interface{}{
		// 			"cpu-bind":       "cores",
		// 			"ntasks":         4,
		// 			"cpus-per-task":  2,
		// 		},
	}

	resp3, err := client.Jobs().SubmitRaw(ctx, coreBindJob)
	if err != nil {
		log.Printf("Failed to submit core binding job: %v", err)
	} else {
		fmt.Printf("Core binding job submitted: %s\n", fmt.Sprintf("%d", resp3.JobId))
	}
}

// discoverAndAllocateResources dynamically discovers and allocates resources
func discoverAndAllocateResources(ctx context.Context, client slurm.SlurmClient) {
	// First, discover available resources
	fmt.Println("Discovering available resources...")

	// Get partition information
	partitions, err := client.Partitions().List(ctx, nil)
	if err != nil {
		log.Printf("Failed to list partitions: %v", err)
		return
	}

	// Find best partition for our needs
	var bestPartition *slurm.Partition
	for _, p := range partitions.Partitions {
		// Check if partition has sufficient resources
		totalCPUs := int32(0)
		if p.CPUs != nil && p.CPUs.Total != nil {
			totalCPUs = *p.CPUs.Total
		}
		totalNodes := int32(0)
		if p.Nodes != nil && p.Nodes.Total != nil {
			totalNodes = *p.Nodes.Total
		}

		if totalCPUs >= 16 && totalNodes >= 1 {
			if bestPartition == nil {
				pCopy := p
				bestPartition = &pCopy
			} else {
				bestCPUs := int32(0)
				if bestPartition.CPUs != nil && bestPartition.CPUs.Total != nil {
					bestCPUs = *bestPartition.CPUs.Total
				}
				if totalCPUs > bestCPUs {
					pCopy := p
					bestPartition = &pCopy
				}
			}
		}
	}

	if bestPartition == nil {
		log.Println("No suitable partition found")
		return
	}

	name := ""
	if bestPartition.Name != nil {
		name = *bestPartition.Name
	}
	totalCPUs := int32(0)
	if bestPartition.CPUs != nil && bestPartition.CPUs.Total != nil {
		totalCPUs = *bestPartition.CPUs.Total
	}
	totalNodes := int32(0)
	if bestPartition.Nodes != nil && bestPartition.Nodes.Total != nil {
		totalNodes = *bestPartition.Nodes.Total
	}
	fmt.Printf("Selected partition: %s (CPUs: %d, Nodes: %d)\n", name, totalCPUs, totalNodes)

	// Get node information for the partition
	partitionName := ""
	if bestPartition.Name != nil {
		partitionName = *bestPartition.Name
	}
	nodes, err := client.Nodes().List(ctx, &slurm.ListNodesOptions{
		Partition: partitionName,
		States:    []string{"IDLE", "MIXED"},
	})
	if err != nil {
		log.Printf("Failed to list nodes: %v", err)
		return
	}

	// Find nodes with specific features
	var gpuNodes []string
	var highMemNodes []string

	for _, node := range nodes.Nodes {
		nodeName := ""
		if node.Name != nil {
			nodeName = *node.Name
		}

		// Check for GPU nodes
		for _, feature := range node.Features {
			if strings.Contains(strings.ToLower(feature), "gpu") {
				gpuNodes = append(gpuNodes, nodeName)
				break
			}
		}

		// Check for high memory nodes (>256GB)
		if node.RealMemory != nil && *node.RealMemory > 262144 {
			highMemNodes = append(highMemNodes, nodeName)
		}
	}

	fmt.Printf("Found %d GPU nodes: %v\n", len(gpuNodes), gpuNodes)
	fmt.Printf("Found %d high-memory nodes: %v\n", len(highMemNodes), highMemNodes)

	// Submit job based on discovered resources
	if len(gpuNodes) > 0 {
		partName := ""
		if bestPartition.Name != nil {
			partName = *bestPartition.Name
		}
		gpuJob := &slurm.JobCreate{
			Name: ptrString("dynamic-gpu-job"),
			Script: ptrString(fmt.Sprintf(`#!/bin/bash
#SBATCH --nodelist=%s

echo "Running on discovered GPU node: $SLURMD_NODENAME"
nvidia-smi
python3 gpu_workload.py
`, gpuNodes[0])),
			Partition:     ptrString(partName),
			MinimumCPUs:   ptrInt32(8),
			MemoryPerNode: ptrUint64(32768),
			TimeLimit:     ptrUint32(60),
			// Metadata would be in SBATCH directives
			// 		// Removed: Metadata: map[string]interface{}{
			// 				"gres":     "gpu:1",
			// 				"nodelist": gpuNodes[0],
			// 			},
		}

		resp, err := client.Jobs().SubmitRaw(ctx, gpuJob)
		if err != nil {
			log.Printf("Failed to submit GPU job: %v", err)
		} else {
			fmt.Printf("GPU job submitted to discovered node: %s\n", fmt.Sprintf("%d", resp.JobId))
		}
	}

	// Submit job to high memory node if available
	if len(highMemNodes) > 0 {
		partName := ""
		if bestPartition.Name != nil {
			partName = *bestPartition.Name
		}
		memJob := &slurm.JobCreate{
			Name: ptrString("dynamic-highmem-job"),
			Script: ptrString(fmt.Sprintf(`#!/bin/bash
#SBATCH --nodelist=%s

echo "Running on high-memory node: $SLURMD_NODENAME"
free -h
python3 memory_analysis.py --use-all-memory
`, highMemNodes[0])),
			Partition:     ptrString(partName),
			MinimumCPUs:   ptrInt32(16),
			MemoryPerNode: ptrUint64(262144), // 256GB
			TimeLimit:     ptrUint32(90),
			// Metadata would be in SBATCH directives
			// 		// Removed: Metadata: map[string]interface{}{
			// 				"nodelist": highMemNodes[0],
			// 			},
		}

		resp, err := client.Jobs().SubmitRaw(ctx, memJob)
		if err != nil {
			log.Printf("Failed to submit high-memory job: %v", err)
		} else {
			fmt.Printf("High-memory job submitted to discovered node: %s\n", fmt.Sprintf("%d", resp.JobId))
		}
	}
}

func ptrString(s string) *string { return &s }
func ptrInt32(i int32) *int32    { return &i }
func ptrUint32(i uint32) *uint32 { return &i }
func ptrUint64(i uint64) *uint64 { return &i }
