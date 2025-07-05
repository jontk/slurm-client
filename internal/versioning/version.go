package versioning

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// APIVersion represents a Slurm REST API version
type APIVersion struct {
	Major int
	Minor int
	Patch int
	Raw   string
}

// ParseVersion parses a version string like "v0.0.42" into an APIVersion
func ParseVersion(version string) (*APIVersion, error) {
	// Remove 'v' prefix if present
	version = strings.TrimPrefix(version, "v")
	
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid version format: %s (expected x.y.z)", version)
	}
	
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid major version: %s", parts[0])
	}
	
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid minor version: %s", parts[1])
	}
	
	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, fmt.Errorf("invalid patch version: %s", parts[2])
	}
	
	return &APIVersion{
		Major: major,
		Minor: minor,
		Patch: patch,
		Raw:   fmt.Sprintf("v%d.%d.%d", major, minor, patch),
	}, nil
}

// String returns the string representation of the version
func (v *APIVersion) String() string {
	return v.Raw
}

// Compare compares two versions. Returns:
// -1 if v < other
//  0 if v == other
//  1 if v > other
func (v *APIVersion) Compare(other *APIVersion) int {
	if v.Major != other.Major {
		if v.Major < other.Major {
			return -1
		}
		return 1
	}
	
	if v.Minor != other.Minor {
		if v.Minor < other.Minor {
			return -1
		}
		return 1
	}
	
	if v.Patch != other.Patch {
		if v.Patch < other.Patch {
			return -1
		}
		return 1
	}
	
	return 0
}

// IsCompatibleWith checks if this version is compatible with another version
// Based on Slurm's compatibility rules
func (v *APIVersion) IsCompatibleWith(other *APIVersion) bool {
	// Same version is always compatible
	if v.Compare(other) == 0 {
		return true
	}
	
	// Major version must match
	if v.Major != other.Major {
		return false
	}
	
	// Minor version must match
	if v.Minor != other.Minor {
		return false
	}
	
	// Patch versions are generally compatible within same minor version
	return true
}

// SupportedVersions represents the versions supported by this client
var SupportedVersions = []*APIVersion{
	{Major: 0, Minor: 0, Patch: 40, Raw: "v0.0.40"},
	{Major: 0, Minor: 0, Patch: 41, Raw: "v0.0.41"},
	{Major: 0, Minor: 0, Patch: 42, Raw: "v0.0.42"},
	{Major: 0, Minor: 0, Patch: 43, Raw: "v0.0.43"},
}

// LatestVersion returns the latest supported version
func LatestVersion() *APIVersion {
	if len(SupportedVersions) == 0 {
		return nil
	}
	
	latest := SupportedVersions[0]
	for _, v := range SupportedVersions[1:] {
		if v.Compare(latest) > 0 {
			latest = v
		}
	}
	
	return latest
}

// StableVersion returns the stable version (v0.0.42 as per ARCHITECTURE.md)
func StableVersion() *APIVersion {
	stable, _ := ParseVersion("v0.0.42")
	return stable
}

// FindBestVersion finds the best supported version for the given constraint
func FindBestVersion(constraint string) (*APIVersion, error) {
	if constraint == "" || constraint == "latest" {
		return LatestVersion(), nil
	}
	
	if constraint == "stable" {
		return StableVersion(), nil
	}
	
	// Try to parse as exact version
	requested, err := ParseVersion(constraint)
	if err != nil {
		return nil, fmt.Errorf("invalid version constraint: %s", constraint)
	}
	
	// Check if we support this exact version
	for _, supported := range SupportedVersions {
		if supported.Compare(requested) == 0 {
			return supported, nil
		}
	}
	
	// Find compatible version
	var compatible []*APIVersion
	for _, supported := range SupportedVersions {
		if supported.IsCompatibleWith(requested) {
			compatible = append(compatible, supported)
		}
	}
	
	if len(compatible) == 0 {
		return nil, fmt.Errorf("no compatible version found for %s", constraint)
	}
	
	// Return the latest compatible version
	sort.Slice(compatible, func(i, j int) bool {
		return compatible[i].Compare(compatible[j]) > 0
	})
	
	return compatible[0], nil
}

// VersionCompatibilityMatrix defines compatibility between versions
// Based on Slurm's version compatibility rules from ARCHITECTURE.md
type VersionCompatibilityMatrix struct {
	// SlurmVersions maps API versions to compatible Slurm versions
	SlurmVersions map[string][]string
	
	// BreakingChanges maps version transitions to breaking changes
	BreakingChanges map[string][]BreakingChange
}

// BreakingChange represents a breaking change between versions
type BreakingChange struct {
	Type        string `json:"type"`        // "field_rename", "field_removed", "endpoint_changed"
	Description string `json:"description"`
	OldValue    string `json:"old_value,omitempty"`
	NewValue    string `json:"new_value,omitempty"`
	Mitigation  string `json:"mitigation,omitempty"`
}

// DefaultCompatibilityMatrix returns the default compatibility matrix
func DefaultCompatibilityMatrix() *VersionCompatibilityMatrix {
	return &VersionCompatibilityMatrix{
		SlurmVersions: map[string][]string{
			"v0.0.40": {"24.05", "24.11", "25.05"},
			"v0.0.41": {"24.11", "25.05", "25.11"},
			"v0.0.42": {"25.05", "25.11"},
			"v0.0.43": {"25.05", "25.11"},
		},
		BreakingChanges: map[string][]BreakingChange{
			"v0.0.40->v0.0.41": {
				{
					Type:        "field_rename",
					Description: "minimum_switches field renamed to required_switches",
					OldValue:    "minimum_switches",
					NewValue:    "required_switches",
					Mitigation:  "Use compatibility layer for automatic field mapping",
				},
			},
			"v0.0.41->v0.0.42": {
				{
					Type:        "field_removed",
					Description: "exclusive and oversubscribe fields removed from job outputs",
					OldValue:    "exclusive, oversubscribe",
					NewValue:    "",
					Mitigation:  "Fields moved to metadata section",
				},
			},
			"v0.0.42->v0.0.43": {
				{
					Type:        "endpoint_added",
					Description: "Reservation management endpoints added",
					OldValue:    "",
					NewValue:    "/slurm/v0.0.43/reservations",
					Mitigation:  "Feature available only in v0.0.43+",
				},
				{
					Type:        "feature_removed",
					Description: "FrontEnd mode support removed",
					OldValue:    "FrontEnd mode",
					NewValue:    "",
					Mitigation:  "Use standard node management",
				},
			},
		},
	}
}

// GetBreakingChanges returns breaking changes for a version transition
func (m *VersionCompatibilityMatrix) GetBreakingChanges(from, to *APIVersion) []BreakingChange {
	key := fmt.Sprintf("%s->%s", from.String(), to.String())
	return m.BreakingChanges[key]
}

// IsSlurmVersionSupported checks if a Slurm version is supported by an API version
func (m *VersionCompatibilityMatrix) IsSlurmVersionSupported(apiVersion, slurmVersion string) bool {
	supported, exists := m.SlurmVersions[apiVersion]
	if !exists {
		return false
	}
	
	for _, version := range supported {
		if strings.HasPrefix(slurmVersion, version) {
			return true
		}
	}
	
	return false
}