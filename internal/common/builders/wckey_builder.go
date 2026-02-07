// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package builders

import (
	"fmt"

	types "github.com/jontk/slurm-client/api"
)

// WCKeyBuilder provides a fluent interface for building WCKey objects
type WCKeyBuilder struct {
	wckey  *types.WCKeyCreate
	errors []error
}

// NewWCKeyBuilder creates a new WCKey builder
func NewWCKeyBuilder(name string) *WCKeyBuilder {
	return &WCKeyBuilder{
		wckey: &types.WCKeyCreate{
			Name: name,
		},
		errors: []error{},
	}
}

// WithUser sets the user for the WCKey
func (b *WCKeyBuilder) WithUser(user string) *WCKeyBuilder {
	if user == "" {
		b.errors = append(b.errors, fmt.Errorf("user cannot be empty"))
		return b
	}
	b.wckey.User = user
	return b
}

// WithCluster sets the cluster for the WCKey
func (b *WCKeyBuilder) WithCluster(cluster string) *WCKeyBuilder {
	if cluster == "" {
		b.errors = append(b.errors, fmt.Errorf("cluster cannot be empty"))
		return b
	}
	b.wckey.Cluster = cluster
	return b
}

// Build creates the final WCKey object
func (b *WCKeyBuilder) Build() (*types.WCKeyCreate, error) {
	// Check for accumulated errors
	if len(b.errors) > 0 {
		return nil, fmt.Errorf("WCKey builder errors: %v", b.errors)
	}

	// Validate required fields
	if b.wckey.Name == "" {
		return nil, fmt.Errorf("WCKey name is required")
	}
	if b.wckey.User == "" {
		return nil, fmt.Errorf("WCKey user is required")
	}
	if b.wckey.Cluster == "" {
		return nil, fmt.Errorf("WCKey cluster is required")
	}

	// Return a copy to prevent external modifications
	result := &types.WCKeyCreate{
		Name:    b.wckey.Name,
		User:    b.wckey.User,
		Cluster: b.wckey.Cluster,
	}

	return result, nil
}

// Clone creates a copy of the current builder state
func (b *WCKeyBuilder) Clone() *WCKeyBuilder {
	if b.wckey == nil {
		return &WCKeyBuilder{
			wckey:  &types.WCKeyCreate{},
			errors: append([]error{}, b.errors...),
		}
	}

	return &WCKeyBuilder{
		wckey: &types.WCKeyCreate{
			Name:    b.wckey.Name,
			User:    b.wckey.User,
			Cluster: b.wckey.Cluster,
		},
		errors: append([]error{}, b.errors...),
	}
}

// Reset clears the builder state
func (b *WCKeyBuilder) Reset() *WCKeyBuilder {
	b.wckey = &types.WCKeyCreate{}
	b.errors = []error{}
	return b
}

// Validate checks if the current state would build successfully
func (b *WCKeyBuilder) Validate() error {
	if b.wckey.Name == "" {
		return fmt.Errorf("WCKey name is required")
	}
	if b.wckey.User == "" {
		return fmt.Errorf("WCKey user is required")
	}
	if b.wckey.Cluster == "" {
		return fmt.Errorf("WCKey cluster is required")
	}

	if len(b.errors) > 0 {
		return fmt.Errorf("builder has errors: %v", b.errors)
	}

	return nil
}

// String returns a string representation of the WCKey being built
func (b *WCKeyBuilder) String() string {
	return fmt.Sprintf("WCKey{Name: %s, User: %s, Cluster: %s}",
		b.wckey.Name, b.wckey.User, b.wckey.Cluster)
}
