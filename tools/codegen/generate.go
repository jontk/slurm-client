//go:build ignore

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// Code generation tool for Slurm REST API clients using oapi-codegen
func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run generate.go <version> [spec-file]")
	}

	version := os.Args[1]
	var specFile string

	if len(os.Args) >= 3 {
		specFile = os.Args[2]
	} else {
		specFile = filepath.Join("openapi-specs", fmt.Sprintf("slurm-%s.json", version))
	}

	if err := generateClient(version, specFile); err != nil {
		log.Fatalf("Failed to generate client for %s: %v", version, err)
	}

	fmt.Printf("Successfully generated client for version %s\n", version)
}

func generateClient(version, specFile string) error {
	// Check if spec file exists
	if _, err := os.Stat(specFile); os.IsNotExist(err) {
		return fmt.Errorf("OpenAPI spec file not found: %s", specFile)
	}

	// Ensure output directory exists
	outputDir := filepath.Join("internal", "api", normalizeVersion(version))
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	outputFile := filepath.Join(outputDir, "client.go")

	// Generate client using oapi-codegen
	cmd := exec.Command("oapi-codegen",
		"-package", normalizeVersion(version),
		"-generate", "client,models,spec",
		"-o", outputFile,
		specFile,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("oapi-codegen failed: %w", err)
	}

	return nil
}


func normalizeVersion(version string) string {
	// Convert v0.0.42 to v0_0_42 for package names
	if len(version) > 0 && version[0] == 'v' {
		version = version[1:]
	}

	result := "v"
	for _, char := range version {
		if char == '.' {
			result += "_"
		} else {
			result += string(char)
		}
	}

	return result
}
