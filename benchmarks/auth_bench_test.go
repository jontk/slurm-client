// SPDX-FileCopyrightText: 2025 Jon Thor Kristinsson
// SPDX-License-Identifier: Apache-2.0

package benchmarks

import (
	"context"
	"net/http"
	"testing"

	"github.com/jontk/slurm-client/pkg/auth"
)

func BenchmarkTokenAuth(b *testing.B) {
	ta := auth.NewTokenAuth("test-jwt-token-with-reasonable-length")
	ctx := context.Background()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://example.com", http.NoBody)

	b.ResetTimer()
	for range b.N {
		ta.Authenticate(ctx, req)
	}
}

func BenchmarkBasicAuth(b *testing.B) {
	ba := auth.NewBasicAuth("username", "password")
	ctx := context.Background()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://example.com", http.NoBody)

	b.ResetTimer()
	for range b.N {
		ba.Authenticate(ctx, req)
	}
}

func BenchmarkNoAuth(b *testing.B) {
	na := auth.NewNoAuth()
	ctx := context.Background()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://example.com", http.NoBody)

	b.ResetTimer()
	for range b.N {
		na.Authenticate(ctx, req)
	}
}
