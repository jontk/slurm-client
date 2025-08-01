// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"testing"
)

func TestCLI(t *testing.T) {
	// Simple test to ensure the CLI compiles
	// In a real implementation, you would test individual commands
	if rootCmd == nil {
		t.Fatal("rootCmd is nil")
	}

	// Test that version info is set
	if Version == "" {
		t.Error("Version is not set")
	}

	// Test that subcommands are registered
	expectedCommands := []string{"jobs", "nodes", "partitions", "info", "submit", "version"}
	for _, cmdName := range expectedCommands {
		found := false
		for _, cmd := range rootCmd.Commands() {
			if cmd.Name() == cmdName {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Command %s not found", cmdName)
		}
	}
}

func TestCreateClient(t *testing.T) {
	// Test client creation with missing URL
	baseURL = ""
	_, err := createClient()
	if err == nil {
		t.Error("Expected error when creating client without URL")
	}
}
