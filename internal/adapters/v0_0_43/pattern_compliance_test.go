// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package v0_0_43

import (
	"testing"
	"os"
	"strings"
	"io/ioutil"
	"regexp"
)

// TestPatternCompliance validates that all adapters follow the QoS pattern
func TestPatternCompliance(t *testing.T) {
	adapters := []struct {
		name            string
		adapterFile     string
		converterFile   string
		apiType         string
		shouldHaveTests bool
		isImplemented   bool
	}{
		{
			name:            "QoS",
			adapterFile:     "qos_adapter.go",
			converterFile:   "qos_converters.go",
			apiType:         "V0043Qos",
			shouldHaveTests: true,
			isImplemented:   true,
		},
		{
			name:            "Account",
			adapterFile:     "account_adapter.go",
			converterFile:   "account_converters.go",
			apiType:         "V0043Account", // Should be this, not V0043AccountInfo
			shouldHaveTests: true,
			isImplemented:   true,
		},
		{
			name:            "Job",
			adapterFile:     "job_adapter.go",
			converterFile:   "job_converters.go",
			apiType:         "V0043Job", // Should be this, not V0043JobProperties
			shouldHaveTests: true,
			isImplemented:   true,
		},
		{
			name:            "Partition",
			adapterFile:     "partition_adapter.go",
			converterFile:   "partition_converters.go",
			apiType:         "V0043PartitionInfo",
			shouldHaveTests: true,
			isImplemented:   true,
		},
		{
			name:            "Node",
			adapterFile:     "node_adapter.go",
			converterFile:   "node_converters.go",
			apiType:         "V0043Node", // Should be this, not V0043NodeInfo
			shouldHaveTests: true,
			isImplemented:   true,
		},
		{
			name:            "User",
			adapterFile:     "user_adapter.go",
			converterFile:   "",
			apiType:         "V0043User", // Should be this, not V0043UserInfo
			shouldHaveTests: true,
			isImplemented:   true,
		},
		{
			name:            "Association",
			adapterFile:     "association_adapter.go",
			converterFile:   "",
			apiType:         "V0043Assoc", // Should be this, not V0043AssociationInfo
			shouldHaveTests: true,
			isImplemented:   true,
		},
		{
			name:            "Reservation",
			adapterFile:     "reservation_adapter.go",
			converterFile:   "",
			apiType:         "V0043ReservationInfo",
			shouldHaveTests: true,
			isImplemented:   true,
		},
	}

	for _, adapter := range adapters {
		t.Run(adapter.name+"_PatternCompliance", func(t *testing.T) {
			if !adapter.isImplemented {
				t.Skip("Not yet implemented")
			}

			// Check adapter file exists
			if _, err := os.Stat(adapter.adapterFile); os.IsNotExist(err) {
				t.Errorf("Missing adapter file: %s", adapter.adapterFile)
				return
			}

			// Read adapter file
			content, err := ioutil.ReadFile(adapter.adapterFile)
			if err != nil {
				t.Errorf("Failed to read %s: %v", adapter.adapterFile, err)
				return
			}

			// Check for base manager usage
			if !strings.Contains(string(content), adapter.name+"BaseManager") {
				t.Errorf("%s adapter should use %sBaseManager", adapter.name, adapter.name)
			}

			// Check for proper struct definition
			adapterStructPattern := regexp.MustCompile(`type ` + adapter.name + `Adapter struct \{[^}]+baseManager[^}]+\}`)
			if !adapterStructPattern.Match(content) {
				t.Errorf("%s adapter should embed baseManager", adapter.name)
			}

			// Check for constructor
			constructorPattern := regexp.MustCompile(`func New` + adapter.name + `Adapter\(.*\) \*` + adapter.name + `Adapter`)
			if !constructorPattern.Match(content) {
				t.Errorf("%s adapter should have New%sAdapter constructor", adapter.name, adapter.name)
			}

			// Check converter file if specified
			if adapter.converterFile != "" {
				if _, err := os.Stat(adapter.converterFile); os.IsNotExist(err) {
					// Converters might be in adapter file
					t.Logf("Note: No separate converter file for %s", adapter.name)
				} else {
					// Check converter patterns
					converterContent, err := ioutil.ReadFile(adapter.converterFile)
					if err == nil {
						// Check for convert functions
						if !strings.Contains(string(converterContent), "convertAPI"+adapter.name+"ToCommon") {
							t.Errorf("%s should have convertAPI%sToCommon function", adapter.converterFile, adapter.name)
						}
						if !strings.Contains(string(converterContent), "convertCommon"+adapter.name+"CreateToAPI") &&
						   !strings.Contains(string(converterContent), "convertCommon"+adapter.name+"UpdateToAPI") {
							t.Errorf("%s should have conversion functions for Create/Update", adapter.converterFile)
						}
					}
				}
			}

			// Check for test file if required
			if adapter.shouldHaveTests {
				testFile := strings.TrimSuffix(adapter.adapterFile, ".go") + "_test.go"
				if _, err := os.Stat(testFile); os.IsNotExist(err) {
					t.Errorf("Missing test file: %s", testFile)
				}
			}
		})
	}
}

// TestAPITypeCorrectness validates that adapters use correct API types
func TestAPITypeCorrectness(t *testing.T) {
	incorrectPatterns := map[string]string{
		"V0043AccountInfo":     "V0043Account",
		"V0043AssociationInfo": "V0043Assoc",
		"V0043JobProperties":   "V0043Job",
		"V0043NodeInfo":        "V0043Node",
		"V0043UserInfo":        "V0043User",
	}

	files := []string{
		"account_adapter.go",
		"account_converters.go",
		"association_adapter.go",
		"job_adapter.go",
		"job_converters.go",
		"node_adapter.go",
		"node_converters.go",
		"user_adapter.go",
	}

	for _, file := range files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			continue
		}

		content, err := ioutil.ReadFile(file)
		if err != nil {
			t.Errorf("Failed to read %s: %v", file, err)
			continue
		}

		for incorrect, correct := range incorrectPatterns {
			if strings.Contains(string(content), incorrect) {
				t.Errorf("%s uses incorrect type '%s', should use '%s'", file, incorrect, correct)
			}
		}
	}
}

// TestCRUDMethodImplementation validates that all adapters implement CRUD methods correctly
func TestCRUDMethodImplementation(t *testing.T) {
	adapters := []struct {
		name     string
		file     string
		hasList  bool
		hasGet   bool
		hasCreate bool
		hasUpdate bool
		hasDelete bool
	}{
		{"QoS", "qos_adapter.go", true, true, true, true, true},
		{"Account", "account_adapter.go", true, true, true, true, true},
		{"Job", "job_adapter.go", true, true, true, true, true},
		{"Partition", "partition_adapter.go", true, true, false, true, false},
		{"Node", "node_adapter.go", true, true, false, true, false},
		{"User", "user_adapter.go", true, true, true, true, true},
		{"Association", "association_adapter.go", true, true, true, true, true},
		{"Reservation", "reservation_adapter.go", true, true, true, true, true},
	}

	for _, adapter := range adapters {
		t.Run(adapter.name+"_CRUD", func(t *testing.T) {
			if _, err := os.Stat(adapter.file); os.IsNotExist(err) {
				t.Skip("File not found")
			}

			content, err := ioutil.ReadFile(adapter.file)
			if err != nil {
				t.Errorf("Failed to read %s: %v", adapter.file, err)
				return
			}

			contentStr := string(content)

			// Check CRUD methods
			methods := []struct {
				name     string
				pattern  string
				required bool
			}{
				{"List", `func \(.*\) List` + adapter.name, adapter.hasList},
				{"Get", `func \(.*\) Get` + adapter.name, adapter.hasGet},
				{"Create", `func \(.*\) Create` + adapter.name, adapter.hasCreate},
				{"Update", `func \(.*\) Update` + adapter.name, adapter.hasUpdate},
				{"Delete", `func \(.*\) Delete` + adapter.name, adapter.hasDelete},
			}

			for _, method := range methods {
				hasMethod := regexp.MustCompile(method.pattern).MatchString(contentStr)
				if method.required && !hasMethod {
					t.Errorf("%s adapter should implement %s method", adapter.name, method.name)
				} else if !method.required && hasMethod {
					t.Logf("Note: %s adapter implements optional %s method", adapter.name, method.name)
				}
			}
		})
	}
}
