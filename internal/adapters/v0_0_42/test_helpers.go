// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0
package v0_0_42

import (
	types "github.com/jontk/slurm-client/api"
)

// Helper functions for tests
func ptrString(s string) *string {
	return &s
}
func ptrInt64(i int64) *int64 {
	return &i
}
func ptrUint32(u uint32) *uint32 {
	return &u
}
func ptrNodeState(s types.NodeState) *types.NodeState {
	return &s
}
