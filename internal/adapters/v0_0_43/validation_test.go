// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_43

import (
	"testing"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"fmt"
	"os"
	"path/filepath"
)

// TestValidateAPITypeUsage validates that all adapters use correct API types
func TestValidateAPITypeUsage(t *testing.T) {
	// Map of incorrect type names to correct ones based on actual API
	typeCorrections := map[string]string{
		"V0043AccountInfo":     "V0043Account",
		"V0043AssociationInfo": "V0043Assoc",
		"V0043JobProperties":   "V0043Job",
		"V0043NodeInfo":        "V0043Node",
		"V0043PartitionInfo":   "V0043PartitionInfo", // This one seems correct
		"V0043ReservationInfo": "V0043ReservationInfo", // This one seems correct
		"V0043UserInfo":        "V0043User",
	}

	// Files to check
	adapterFiles := []string{
		"account_converters.go",
		"association_adapter.go",
		"job_converters.go",
		"node_converters.go",
		"partition_converters.go",
		"reservation_adapter.go",
		"user_adapter.go",
	}

	var errors []string

	for _, file := range adapterFiles {
		filePath := filepath.Join(".", file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			continue // Skip files that don't exist yet
		}

		// Parse the file
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
		if err != nil {
			t.Errorf("Failed to parse %s: %v", file, err)
			continue
		}

		// Check for incorrect type usage
		ast.Inspect(node, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.SelectorExpr:
				if ident, ok := x.X.(*ast.Ident); ok && ident.Name == "api" {
					typeName := x.Sel.Name
					if correctType, needsCorrection := typeCorrections[typeName]; needsCorrection {
						pos := fset.Position(x.Pos())
						errors = append(errors, fmt.Sprintf(
							"%s:%d:%d: Using incorrect type '%s', should be '%s'",
							file, pos.Line, pos.Column, typeName, correctType,
						))
					}
				}
			}
			return true
		})
	}

	if len(errors) > 0 {
		t.Errorf("Found %d API type mismatches:\n%s", len(errors), strings.Join(errors, "\n"))
	}
}

// TestValidateAdapterImplementations validates that all adapters follow the correct pattern
func TestValidateAdapterImplementations(t *testing.T) {
	requiredAdapters := []struct {
		name          string
		adapterFile   string
		converterFile string
		baseManager   string
		commonType    string
	}{
		{
			name:          "Account",
			adapterFile:   "account_adapter.go",
			converterFile: "account_converters.go",
			baseManager:   "AccountBaseManager",
			commonType:    "Account",
		},
		{
			name:          "Association",
			adapterFile:   "association_adapter.go",
			converterFile: "association_converters.go",
			baseManager:   "AssociationBaseManager",
			commonType:    "Association",
		},
		{
			name:          "Job",
			adapterFile:   "job_adapter.go",
			converterFile: "job_converters.go",
			baseManager:   "JobBaseManager",
			commonType:    "Job",
		},
		{
			name:          "Node",
			adapterFile:   "node_adapter.go",
			converterFile: "node_converters.go",
			baseManager:   "NodeBaseManager",
			commonType:    "Node",
		},
		{
			name:          "Partition",
			adapterFile:   "partition_adapter.go",
			converterFile: "partition_converters.go",
			baseManager:   "PartitionBaseManager",
			commonType:    "Partition",
		},
		{
			name:          "Reservation",
			adapterFile:   "reservation_adapter.go",
			converterFile: "reservation_converters.go",
			baseManager:   "ReservationBaseManager",
			commonType:    "Reservation",
		},
		{
			name:          "User",
			adapterFile:   "user_adapter.go",
			converterFile: "user_converters.go",
			baseManager:   "UserBaseManager",
			commonType:    "User",
		},
	}

	for _, adapter := range requiredAdapters {
		t.Run(adapter.name, func(t *testing.T) {
			// Check adapter file exists
			if _, err := os.Stat(adapter.adapterFile); os.IsNotExist(err) {
				t.Errorf("Missing adapter file: %s", adapter.adapterFile)
			}

			// Check converter file exists (might not exist for all)
			if adapter.converterFile != "" {
				if _, err := os.Stat(adapter.converterFile); os.IsNotExist(err) {
					// Some adapters might have converters in the adapter file itself
					t.Logf("Note: No separate converter file for %s", adapter.name)
				}
			}

			// TODO: Add more validation logic here
			// - Check that adapter uses base manager
			// - Check that converter methods follow naming convention
			// - Check that all CRUD operations are implemented
		})
	}
}

// TestValidateBuildErrors checks if the package can build successfully
func TestValidateBuildErrors(t *testing.T) {
	// This test will fail if there are build errors
	// Run: go test -c ./internal/adapters/v0_0_43
	t.Log("Run 'go build ./internal/adapters/v0_0_43' to check for build errors")
}
